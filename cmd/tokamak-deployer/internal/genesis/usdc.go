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
	usdcProxyAddress           = "0x4200000000000000000000000000000000000778"
	proxyAdminAddress          = "0x4200000000000000000000000000000000000018"
	zeppelinImplementationSlot = "0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3"
	zeppelinAdminSlot          = "0x10d6a54a4754c8869d6886b5f5d7fbfa5b4522237ea5c60d11bc4e7a1ff9390b"
)

// bytecodeFile represents the JSON structure for pre-extracted bytecode files.
type bytecodeFile struct {
	Bytecode string `json:"bytecode"`
}

// encodeSolidityShortString encodes a short string (< 32 bytes) into a Solidity
// storage slot value. The string bytes are left-aligned and the last byte stores length*2.
func encodeSolidityShortString(s string) common.Hash {
	if len(s) >= 32 {
		// Should not happen for our use cases; return zero hash as safety
		return common.Hash{}
	}
	var slot [32]byte
	copy(slot[:], []byte(s))
	slot[31] = byte(len(s) * 2)
	return common.BytesToHash(slot[:])
}

// injectUSDCIntoGenesis injects the FiatTokenV2_2 USDC predeploy into genesis.json.
// For DeFi/Full presets where USDC is already present, the function is idempotent and skips.
// For General/Gaming presets, it loads bytecodes from embedded artifacts
// and creates proxy + implementation entries with Zeppelin proxy slots and minimal storage.
func injectUSDCIntoGenesis(genesisPath string, artifactsFS fs.FS) error {
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

	proxyKey := formatAddr(usdcProxyAddress)

	// Idempotency check: skip if proxy already has non-empty code
	if existing, ok := alloc[proxyKey]; ok {
		var entry map[string]json.RawMessage
		if err := json.Unmarshal(existing, &entry); err == nil {
			if codeRaw, ok := entry["code"]; ok {
				var code string
				if err := json.Unmarshal(codeRaw, &code); err == nil {
					if code != "" && code != "0x" {
						fmt.Println("USDC predeploy already present at", usdcProxyAddress, ", skipping injection")
						return nil
					}
				}
			}
		}
	}

	// Load bytecodes from embedded artifacts
	implBytecode, err := loadBytecodeFromFS(artifactsFS, "deploy-artifacts/FiatTokenV2_2.json")
	if err != nil {
		// If not found, skip with warning
		fmt.Println("Warning: FiatTokenV2_2.json not found in embedded artifacts, skipping USDC injection")
		return nil
	}

	proxyBytecode, err := loadBytecodeFromFS(artifactsFS, "deploy-artifacts/FiatTokenV2_2Proxy.json")
	if err != nil {
		// If not found, skip with warning
		fmt.Println("Warning: FiatTokenV2_2Proxy.json not found in embedded artifacts, skipping USDC injection")
		return nil
	}

	// Set implementation at code namespace address
	proxyAddr := common.HexToAddress(usdcProxyAddress)
	codeAddr := predeployToCodeNamespace(proxyAddr)
	codeAddrKey := formatAddr(codeAddr.Hex())

	implEntry, err := json.Marshal(map[string]interface{}{
		"code":    implBytecode,
		"balance": "0x0",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal implementation entry: %w", err)
	}
	alloc[codeAddrKey] = implEntry

	// Set proxy at predeploy address with Zeppelin slots and minimal storage
	adminAddr := common.HexToAddress(proxyAdminAddress)

	storage := map[string]string{
		// Zeppelin proxy slots
		zeppelinImplementationSlot: common.BytesToHash(codeAddr.Bytes()).Hex(),
		zeppelinAdminSlot:          common.BytesToHash(adminAddr.Bytes()).Hex(),
		// FiatTokenV2_2 storage slots
		"0x0000000000000000000000000000000000000000000000000000000000000004": encodeSolidityShortString("USD Coin").Hex(),  // name
		"0x0000000000000000000000000000000000000000000000000000000000000005": encodeSolidityShortString("USDC.e").Hex(),    // symbol
		"0x0000000000000000000000000000000000000000000000000000000000000006": common.BytesToHash([]byte{6}).Hex(),           // decimals
		"0x0000000000000000000000000000000000000000000000000000000000000007": encodeSolidityShortString("USD").Hex(),        // currency
		"0x0000000000000000000000000000000000000000000000000000000000000009": common.BytesToHash([]byte{1}).Hex(),           // initialized
	}

	proxyEntry, err := json.Marshal(map[string]interface{}{
		"code":    proxyBytecode,
		"balance": "0x0",
		"storage": storage,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal proxy entry: %w", err)
	}
	alloc[proxyKey] = proxyEntry

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

	fmt.Println("Injected USDC (FiatTokenV2_2) predeploy into genesis at", usdcProxyAddress)
	return os.WriteFile(genesisPath, output, 0644)
}

// loadBytecodeFromFS reads a bytecode JSON file from embedded FS and returns the bytecode string.
// Tries both forge artifact format (deployedBytecode.object) and bytecode file format (bytecode).
func loadBytecodeFromFS(artifactsFS fs.FS, path string) (string, error) {
	data, err := fs.ReadFile(artifactsFS, path)
	if err != nil {
		return "", fmt.Errorf("failed to read bytecode file %s: %w", path, err)
	}

	// Try forge artifact format first
	var artifact forgeArtifact
	if err := json.Unmarshal(data, &artifact); err == nil && artifact.DeployedBytecode.Object != "" {
		return artifact.DeployedBytecode.Object, nil
	}

	// Fall back to bytecodeFile format
	var bf bytecodeFile
	if err := json.Unmarshal(data, &bf); err != nil {
		return "", fmt.Errorf("failed to parse bytecode file %s: %w", path, err)
	}

	if bf.Bytecode == "" {
		return "", fmt.Errorf("empty bytecode in %s", path)
	}

	return bf.Bytecode, nil
}
