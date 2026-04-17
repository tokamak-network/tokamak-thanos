package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tokamak-network/tokamak-thanos/cmd/tokamak-deployer/internal/genesis"
)

var (
	flagDeployOutput string
	flagConfig       string
	flagGenesisOut   string
	flagBaseGenesis  string
	flagRollupOut    string
	flagPreset       string
	flagGenesisL1RPC        string
	flagL2Allocs     string
	flagOpNodeBin    string
)

var generateGenesisCmd = &cobra.Command{
	Use:   "generate-genesis",
	Short: "Generate genesis.json from deploy output and apply post-processing",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := genesis.Config{
			Preset:       flagPreset,
			L1RPCURL:     flagGenesisL1RPC,
			L2AllocsPath: flagL2Allocs,
			OpNodeBinary: flagOpNodeBin,
		}
		return genesis.Generate(
			flagDeployOutput,
			flagConfig,
			flagBaseGenesis,
			flagGenesisOut,
			flagRollupOut,
			DeployArtifactsFS,
			cfg,
		)
	},
}

func init() {
	generateGenesisCmd.Flags().StringVar(&flagDeployOutput, "deploy-output", "./deploy-output.json", "deploy-contracts output file")
	generateGenesisCmd.Flags().StringVar(&flagConfig, "config", "./rollup-config.json", "Rollup config file")
	generateGenesisCmd.Flags().StringVar(&flagGenesisOut, "out", "./genesis.json", "Genesis output file path")
	generateGenesisCmd.Flags().StringVar(&flagBaseGenesis, "base-genesis", "", "Skip op-node and use this as base genesis (for testing or when you've already run op-node genesis l2 externally)")
	generateGenesisCmd.Flags().StringVar(&flagRollupOut, "rollup-out", "", "Rollup output file path (default: same dir as --out)")
	generateGenesisCmd.Flags().StringVar(&flagPreset, "preset", "general", "Preset type: general, gaming, full, defi")
	generateGenesisCmd.Flags().StringVar(&flagGenesisL1RPC, "l1-rpc", "", "L1 RPC URL (required unless --base-genesis is set; passed through to op-node)")
	generateGenesisCmd.Flags().StringVar(&flagL2Allocs, "l2-allocs", "", "L2 genesis allocs file (required unless --base-genesis is set; produced by forge L2Genesis.s.sol)")
	generateGenesisCmd.Flags().StringVar(&flagOpNodeBin, "op-node-bin", "", "Path to op-node binary (default: 'op-node' from $PATH)")
	_ = generateGenesisCmd.MarkFlagRequired("deploy-output")
	_ = generateGenesisCmd.MarkFlagRequired("config")
}
