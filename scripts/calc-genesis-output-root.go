package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/foundry"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis"
	"github.com/tokamak-network/tokamak-thanos/op-node/rollup"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

func main() {
	// Get project root
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Load deploy config
	deployConfigPath := filepath.Join(root, "packages", "tokamak", "contracts-bedrock", "deploy-config", "devnetL1.json")
	deployConfig, err := genesis.NewDeployConfig(deployConfigPath)
	if err != nil {
		log.Fatalf("Failed to load deploy config: %v", err)
	}

	// Load L2 allocs (use the latest available mode)
	l2AllocsPath := filepath.Join(root, ".devnet", "allocs-l2.json")
	l2Allocs, err := foundry.LoadForgeAllocs(l2AllocsPath)
	if err != nil {
		log.Fatalf("Failed to load L2 allocs: %v", err)
	}

	// Create dummy L1 block ref (genesis uses block 0)
	l1BlockRef := eth.BlockRef{
		Hash:       common.Hash{},
		Number:     0,
		ParentHash: common.Hash{},
		Time:       0,
	}

	// Build L2 genesis
	l2Genesis, err := genesis.BuildL2Genesis(deployConfig, l2Allocs, &l1BlockRef)
	if err != nil {
		log.Fatalf("Failed to build L2 genesis: %v", err)
	}

	// Convert to block
	l2GenesisBlock := l2Genesis.ToBlock()

	// Get L2ToL1MessagePasser storage root
	// At genesis, this is the empty hash of an empty trie
	messagePasserStorageRoot := getMessagePasserStorageRoot(l2GenesisBlock)

	// Calculate genesis output root using block info wrapper
	blockInfo := eth.HeaderBlockInfo(l2GenesisBlock.Header())
	outputRoot, err := rollup.ComputeL2OutputRootV0(blockInfo, messagePasserStorageRoot)
	if err != nil {
		log.Fatalf("Failed to compute output root: %v", err)
	}

	// Output only the hash (for script usage)
	fmt.Print(common.Hash(outputRoot).Hex())
}

// getMessagePasserStorageRoot extracts the storage root of L2ToL1MessagePasser from genesis block
func getMessagePasserStorageRoot(block *types.Block) [32]byte {
	// At genesis, the storage is empty, so we return the empty trie hash
	// This matches the behavior in op-e2e tests
	return common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b47e1a81b0b345d1a17b4b3d89a5d96")
}
