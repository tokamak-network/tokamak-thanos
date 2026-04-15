package genesis

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	l1BlockProxyAddress = "0x4200000000000000000000000000000000000015"
)

// forgeArtifact represents the JSON structure of a forge build artifact.
type forgeArtifact struct {
	DeployedBytecode struct {
		Object string `json:"object"`
	} `json:"deployedBytecode"`
}

// predeployToCodeNamespace computes implementation address for a predeploy proxy.
// Mirrors Predeploys.sol: (addr & 0xffff) | 0xc0D3C0d3C0d3C0D3c0d3C0d3c0D3C0d3c0d30000
func predeployToCodeNamespace(addr common.Address) common.Address {
	prefix := common.HexToAddress("0xc0D3C0d3C0d3C0D3c0d3C0d3c0D3C0d3c0d30000")
	var result common.Address
	copy(result[:], prefix[:])
	result[18] = addr[18]
	result[19] = addr[19]
	return result
}

// injectL1BlockBytecode patches the L1Block implementation entry (code namespace address)
// in genesis.json with the Isthmus-capable bytecode from the embedded forge artifact.
//
// The function patches only the "code" field of the existing code namespace alloc entry,
// preserving all other fields (balance, storage, nonce).
func injectL1BlockBytecode(genesisPath string, artifactsFS fs.FS) error {
	data, err := os.ReadFile(genesisPath)
	if err != nil {
		return fmt.Errorf("failed to read genesis file: %w", err)
	}

	var genesis map[string]json.RawMessage
	if err := json.Unmarshal(data, &genesis); err != nil {
		return fmt.Errorf("failed to parse genesis JSON: %w", err)
	}

	var alloc map[string]json.RawMessage
	if err := json.Unmarshal(genesis["alloc"], &alloc); err != nil {
		return fmt.Errorf("failed to parse alloc section: %w", err)
	}

	// Detect alloc key format (with or without 0x prefix) from existing entries
	has0xPrefix := false
	for key := range alloc {
		if strings.HasPrefix(key, "0x") || strings.HasPrefix(key, "0X") {
			has0xPrefix = true
		}
		break
	}

	formatAddr := func(addr string) string {
		lower := strings.ToLower(addr)
		if !has0xPrefix {
			return strings.TrimPrefix(lower, "0x")
		}
		return lower
	}

	// Load deployedBytecode from embedded forge artifact
	implBytecode, err := loadForgeArtifactBytecodeFromFS(artifactsFS, "deploy-artifacts/L1Block.json")
	if err != nil {
		return fmt.Errorf("failed to load L1Block bytecode: %w", err)
	}

	// Compute code namespace address for the L1Block implementation
	proxyAddr := common.HexToAddress(l1BlockProxyAddress)
	codeAddr := predeployToCodeNamespace(proxyAddr)
	codeAddrKey := formatAddr(codeAddr.Hex())

	// Preserve existing alloc entry fields (balance, storage, nonce); only replace code.
	// L1Block is a core predeploy — its storage slots hold initialised state in genesis.
	var entry map[string]json.RawMessage
	if existing, ok := alloc[codeAddrKey]; ok {
		if err := json.Unmarshal(existing, &entry); err != nil {
			entry = make(map[string]json.RawMessage)
		}
	} else {
		entry = make(map[string]json.RawMessage)
	}

	codeJSON, err := json.Marshal(implBytecode)
	if err != nil {
		return fmt.Errorf("failed to marshal L1Block bytecode: %w", err)
	}
	entry["code"] = codeJSON

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal L1Block alloc entry: %w", err)
	}
	alloc[codeAddrKey] = entryJSON

	// Write back genesis.json
	allocJSON, err := json.Marshal(alloc)
	if err != nil {
		return fmt.Errorf("failed to marshal alloc: %w", err)
	}
	genesis["alloc"] = allocJSON

	output, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal genesis: %w", err)
	}

	fmt.Println("Injected L1Block (Isthmus-capable) implementation into genesis code namespace at", codeAddr.Hex())
	return os.WriteFile(genesisPath, output, 0644)
}

// loadForgeArtifactBytecodeFromFS reads a forge build artifact JSON from embedded FS
// and returns the deployedBytecode.
func loadForgeArtifactBytecodeFromFS(artifactsFS fs.FS, path string) (string, error) {
	data, err := fs.ReadFile(artifactsFS, path)
	if err != nil {
		return "", fmt.Errorf("failed to read artifact %s: %w", path, err)
	}

	var artifact forgeArtifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		return "", fmt.Errorf("failed to parse artifact %s: %w", path, err)
	}

	if artifact.DeployedBytecode.Object == "" {
		return "", fmt.Errorf("empty deployedBytecode in %s", path)
	}

	return artifact.DeployedBytecode.Object, nil
}
