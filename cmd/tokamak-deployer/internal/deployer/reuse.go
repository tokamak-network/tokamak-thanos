package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// bytecodeEqual reports whether two bytecode byte slices are identical.
// Uses keccak256 to short-circuit byte-by-byte comparison cost on long bytecode,
// matching the on-chain runtime code returned by eth_getCode.
func bytecodeEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	return crypto.Keccak256Hash(a) == crypto.Keccak256Hash(b)
}

// Registry is the on-disk and embedded format describing previously-deployed
// L1 implementation addresses for a given L1 chain.
type Registry struct {
	TokamakDeployerVersion string            `json:"tokamakDeployerVersion"`
	L1ChainID              uint64            `json:"l1ChainId"`
	Comment                string            `json:"comment,omitempty"`
	Implementations        map[string]string `json:"implementations"`
}

// loadRegistry returns a Registry per the precedence: CLI override > embedded > empty.
//
// An empty registry (not nil, not error) is returned when:
//   - cfg.ReuseDeployment is false
//   - cfg.RegistryPath is empty AND no embedded registry exists for l1ChainID
//
// Errors are returned for:
//   - cfg.RegistryPath set but file unreadable
//   - JSON parse failure
//   - registry.L1ChainID != l1ChainID
func loadRegistry(cfg DeployConfig, l1ChainID uint64, registryFS fs.FS) (*Registry, error) {
	if !cfg.ReuseDeployment {
		return &Registry{Implementations: map[string]string{}}, nil
	}

	var raw []byte
	var src string
	if cfg.RegistryPath != "" {
		b, err := os.ReadFile(cfg.RegistryPath)
		if err != nil {
			return nil, fmt.Errorf("read --reuse-impls %q: %w", cfg.RegistryPath, err)
		}
		raw, src = b, cfg.RegistryPath
	} else {
		path := fmt.Sprintf("registry/%d.json", l1ChainID)
		b, err := fs.ReadFile(registryFS, path)
		if err != nil {
			log.Printf("[deployer] No embedded registry for chainId=%d, reuse will be a no-op", l1ChainID)
			return &Registry{Implementations: map[string]string{}}, nil
		}
		raw, src = b, "embed:"+path
	}

	var r Registry
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("parse registry %s: %w", src, err)
	}
	log.Printf("[deployer] Registry %s declares tokamakDeployerVersion=%q (informational)",
		src, r.TokamakDeployerVersion)
	if r.L1ChainID != l1ChainID {
		return nil, fmt.Errorf("registry %s l1ChainId %d != actual %d",
			src, r.L1ChainID, l1ChainID)
	}
	if r.Implementations == nil {
		r.Implementations = map[string]string{}
	}
	log.Printf("[deployer] Loaded reuse registry from %s (%d entries)", src, len(r.Implementations))
	return &r, nil
}

// ReuseTargets enumerates the implementations eligible for reuse.
// Stable across releases — adding/removing requires bumping registry semantics.
var ReuseTargets = []string{
	"SuperchainConfig",
	"OptimismPortal",
	"SystemConfig",
	"L1StandardBridge",
	"L1CrossDomainMessenger",
	"OptimismMintableERC20Factory",
	"L1ERC721Bridge",
	"L2OutputOracle",
	"DisputeGameFactory",
}

// reuseTable holds verified registry entries — only addresses whose on-chain
// runtime bytecode matches the embedded artifact's deployedBytecode.object.
type reuseTable struct {
	addrs map[string]common.Address
}

func (t *reuseTable) lookup(name string) common.Address {
	if t == nil {
		return common.Address{}
	}
	return t.addrs[name]
}

func (t *reuseTable) size() int {
	if t == nil {
		return 0
	}
	return len(t.addrs)
}

// verify checks each registry entry against on-chain code; returns a reuseTable
// containing only the entries that pass both extcodesize > 0 and bytecode-hash
// equality against the embedded `deployedBytecode.object`.
//
// strict=true causes the first verification failure to abort with an error;
// strict=false silently skips failed entries (caller falls back to fresh deploy).
func (r *Registry) verify(ctx context.Context, client *ethclient.Client, artifactsFS fs.FS, strict bool) (*reuseTable, error) {
	table := &reuseTable{addrs: make(map[string]common.Address)}
	if len(r.Implementations) == 0 {
		return table, nil
	}
	for _, name := range ReuseTargets {
		addrHex, ok := r.Implementations[name]
		if !ok {
			continue
		}
		if !common.IsHexAddress(addrHex) {
			err := fmt.Errorf("registry: %s has invalid address %q", name, addrHex)
			if strict {
				return nil, err
			}
			log.Printf("[deployer] WARN %v — falling back to fresh deploy", err)
			continue
		}
		addr := common.HexToAddress(addrHex)
		code, err := client.CodeAt(ctx, addr, nil)
		if err != nil {
			return nil, fmt.Errorf("get code at %s for %s: %w", addr.Hex(), name, err)
		}
		if len(code) == 0 {
			err := fmt.Errorf("registry: %s @ %s has no code on-chain", name, addr.Hex())
			if strict {
				return nil, err
			}
			log.Printf("[deployer] WARN %v — falling back to fresh deploy", err)
			continue
		}
		a, err := loadArtifact(artifactsFS, name)
		if err != nil {
			return nil, fmt.Errorf("load artifact for %s: %w", name, err)
		}
		expected := common.FromHex(a.DeployedBytecode.Object)
		if !bytecodeEqual(code, expected) {
			err := fmt.Errorf("registry: %s @ %s bytecode mismatch (on-chain %d bytes, expected %d bytes)",
				name, addr.Hex(), len(code), len(expected))
			if strict {
				return nil, err
			}
			log.Printf("[deployer] WARN %v — falling back to fresh deploy", err)
			continue
		}
		table.addrs[name] = addr
		log.Printf("[deployer] Registry preflight: %s @ %s ✓", name, addr.Hex())
	}
	log.Printf("[deployer] Reuse preflight: %d/%d implementations reusable",
		len(table.addrs), len(ReuseTargets))
	return table, nil
}

// deployOrReuse returns the registry-supplied address when the table holds a
// verified entry for `name`; otherwise it falls back to deployContract using
// the caller-supplied signer + nonce + gas price.
//
// On hit: no transaction is sent and no nonce is consumed.
// On miss: identical to calling deployContract(...) directly.
func deployOrReuse(
	ctx context.Context,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	nonce *uint64,
	gasPrice *big.Int,
	name string,
	a *artifact,
	reuse *reuseTable,
) (common.Address, error) {
	if addr := reuse.lookup(name); (addr != common.Address{}) {
		log.Printf("[deployer] ♻ %s impl reused: %s", name, addr.Hex())
		return addr, nil
	}
	addr, err := deployContract(ctx, client, auth, nonce, gasPrice, a)
	if err != nil {
		return common.Address{}, err
	}
	log.Printf("[deployer] ✓ %s impl deployed: %s", name, addr.Hex())
	return addr, nil
}
