package genesis

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core"
)

// updateRollupGenesisHash recomputes L2 genesis block hash from genesis.json
// and updates rollup.json.genesis.l2.hash to match.
func updateRollupGenesisHash(genesisPath, rollupPath string) error {
	// Read and parse genesis.json
	genesisData, err := os.ReadFile(genesisPath)
	if err != nil {
		return fmt.Errorf("failed to read genesis file: %w", err)
	}

	var genesisConfig core.Genesis
	if err := json.Unmarshal(genesisData, &genesisConfig); err != nil {
		return fmt.Errorf("failed to parse genesis JSON: %w", err)
	}

	// Compute genesis block hash
	genesisBlock := genesisConfig.ToBlock()
	genesisHash := genesisBlock.Hash()

	// Read and parse rollup.json
	rollupData, err := os.ReadFile(rollupPath)
	if err != nil {
		return fmt.Errorf("failed to read rollup file: %w", err)
	}

	var rollupConfig map[string]interface{}
	if err := json.Unmarshal(rollupData, &rollupConfig); err != nil {
		return fmt.Errorf("failed to parse rollup JSON: %w", err)
	}

	// Update genesis.l2.hash
	genesis, ok := rollupConfig["genesis"].(map[string]interface{})
	if !ok {
		genesis = make(map[string]interface{})
		rollupConfig["genesis"] = genesis
	}

	l2Config, ok := genesis["l2"].(map[string]interface{})
	if !ok {
		l2Config = make(map[string]interface{})
		genesis["l2"] = l2Config
	}

	l2Config["hash"] = genesisHash.Hex()

	// Write back rollup.json
	output, err := json.MarshalIndent(rollupConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal rollup config: %w", err)
	}

	fmt.Println("Updated rollup.json genesis.l2.hash to", genesisHash.Hex())
	return os.WriteFile(rollupPath, output, 0644)
}
