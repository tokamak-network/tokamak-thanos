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
	multiTokenPaymasterAddress = "0x4200000000000000000000000000000000000067"
)

// injectMultiTokenPaymasterBytecode injects the MultiTokenPaymaster implementation
// bytecode into genesis.json at the predeploy code namespace address.
//
// MultiTokenPaymaster (0x4200...0067) is a Transparent Proxy. The implementation lives at
// the code namespace address (0xc0d3...0067).
func injectMultiTokenPaymasterBytecode(genesisPath string, artifactsFS fs.FS) error {
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
	implBytecode, err := loadForgeArtifactBytecodeFromFS(artifactsFS, "deploy-artifacts/MultiTokenPaymaster.json")
	if err != nil {
		// If not found, skip with warning
		fmt.Println("Warning: MultiTokenPaymaster.json not found in embedded artifacts, skipping injection")
		return nil
	}

	// Set implementation at code namespace address
	proxyAddr := common.HexToAddress(multiTokenPaymasterAddress)
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

	fmt.Println("Injected MultiTokenPaymaster implementation into genesis code namespace at", codeAddr.Hex())
	return os.WriteFile(genesisPath, output, 0644)
}
