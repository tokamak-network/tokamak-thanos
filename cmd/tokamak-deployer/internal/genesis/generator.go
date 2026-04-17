package genesis

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// Config controls which optional genesis post-processing steps to run.
type Config struct {
	// Preset controls optional injections. Values: "general", "gaming", "full", "defi"
	// "gaming" and "full" enable DRB injection.
	Preset string

	// L1RPCURL is passed to `op-node genesis l2 --l1-rpc` when BaseGenesisPath is empty.
	// Required by op-node even when l1StartingBlockTag is present in the deploy config.
	L1RPCURL string

	// L2AllocsPath is passed to `op-node genesis l2 --l2-allocs`. Contains the L2
	// genesis state dump produced by forge's L2Genesis.s.sol script. Required by
	// op-node when BaseGenesisPath is empty.
	L2AllocsPath string

	// OpNodeBinary overrides the op-node binary lookup. Defaults to "op-node" in $PATH.
	OpNodeBinary string
}

// Generate creates genesis.json from deploy-output and rollup config, then applies
// 5 post-processing steps: DRB inject (if preset=gaming|full), USDC inject,
// MultiTokenPaymaster inject, L1Block bytecode patch, rollup hash update.
//
// Parameters:
//   - deployOutputPath: path to deploy-output.json (from deploy-contracts)
//   - configPath: path to rollup-config.json (deploy config)
//   - baseGenesisPath: if non-empty, copy this file to outPath and skip op-node genesis generation
//   - outPath: path to write the final genesis.json
//   - rollupOutPath: path to rollup.json (for hash update); if empty, inferred as same dir as outPath
//   - artifactsFS: embedded deploy-artifacts FS (for L1Block, MultiTokenPaymaster, USDC bytecodes)
//   - cfg: optional post-processing config (including L1RPCURL / L2AllocsPath for op-node)
func Generate(
	deployOutputPath, configPath, baseGenesisPath, outPath, rollupOutPath string,
	artifactsFS fs.FS,
	cfg Config,
) error {
	// Step 1: Generate base genesis
	if baseGenesisPath != "" {
		// Use provided base genesis instead of running op-node
		data, err := os.ReadFile(baseGenesisPath)
		if err != nil {
			return fmt.Errorf("read base genesis: %w", err)
		}
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			return fmt.Errorf("copy base genesis to output: %w", err)
		}
	} else {
		if cfg.L1RPCURL == "" {
			return fmt.Errorf("op-node genesis l2 requires --l1-rpc (pass via Config.L1RPCURL or use --base-genesis)")
		}
		if cfg.L2AllocsPath == "" {
			return fmt.Errorf("op-node genesis l2 requires --l2-allocs (pass via Config.L2AllocsPath — produced by forge L2Genesis.s.sol — or use --base-genesis)")
		}

		opNodeBin := cfg.OpNodeBinary
		if opNodeBin == "" {
			opNodeBin = "op-node"
		}

		// Infer rollup out path if not specified
		inferredRollupPath := rollupOutPath
		if inferredRollupPath == "" {
			inferredRollupPath = filepath.Join(filepath.Dir(outPath), "rollup.json")
		}

		cmd := exec.Command(opNodeBin, "genesis", "l2",
			"--deploy-config", configPath,
			"--l1-deployments", deployOutputPath,
			"--l2-allocs", cfg.L2AllocsPath,
			"--outfile.l2", outPath,
			"--outfile.rollup", inferredRollupPath,
			"--l1-rpc", cfg.L1RPCURL,
		)

		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("op-node genesis l2 failed: %w\n%s", err, out)
		}
	}

	// Step 2: DRB inject (only for gaming/full presets)
	if cfg.Preset == "gaming" || cfg.Preset == "full" {
		if err := injectDRBIntoGenesis(context.Background(), outPath); err != nil {
			return fmt.Errorf("DRB inject: %w", err)
		}
	}

	// Step 3: USDC inject
	if err := injectUSDCIntoGenesis(outPath, artifactsFS); err != nil {
		return fmt.Errorf("USDC inject: %w", err)
	}

	// Step 4: MultiTokenPaymaster inject
	if err := injectMultiTokenPaymasterBytecode(outPath, artifactsFS); err != nil {
		return fmt.Errorf("MultiTokenPaymaster inject: %w", err)
	}

	// Step 5: L1Block bytecode patch
	if err := injectL1BlockBytecode(outPath, artifactsFS); err != nil {
		return fmt.Errorf("L1Block patch: %w", err)
	}

	// Step 6: Rollup hash update
	inferredRollupPath := rollupOutPath
	if inferredRollupPath == "" {
		inferredRollupPath = filepath.Join(filepath.Dir(outPath), "rollup.json")
	}
	if _, err := os.Stat(inferredRollupPath); err == nil {
		if err := updateRollupGenesisHash(outPath, inferredRollupPath); err != nil {
			return fmt.Errorf("rollup hash update: %w", err)
		}
	}
	// else: rollup.json doesn't exist, skip silently

	return nil
}
