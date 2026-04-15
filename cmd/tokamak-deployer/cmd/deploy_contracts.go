package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tokamak-network/tokamak-thanos/cmd/tokamak-deployer/internal/deployer"
)

var (
	flagL1RPC      string
	flagPrivateKey string
	flagChainID    uint64
	flagOut        string
)

var deployContractsCmd = &cobra.Command{
	Use:   "deploy-contracts",
	Short: "Deploy L1 contracts and write deploy-output.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := deployer.DeployConfig{
			L1RPCURL:   flagL1RPC,
			PrivateKey: flagPrivateKey,
			L2ChainID:  flagChainID,
		}
		output, err := deployer.Deploy(cmd.Context(), cfg, DeployArtifactsFS)
		if err != nil {
			return fmt.Errorf("deployment failed: %w", err)
		}
		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(flagOut, data, 0644)
	},
}

func init() {
	deployContractsCmd.Flags().StringVar(&flagL1RPC, "l1-rpc", "", "L1 RPC URL (required)")
	deployContractsCmd.Flags().StringVar(&flagPrivateKey, "private-key", "", "Deployer private key (required)")
	deployContractsCmd.Flags().Uint64Var(&flagChainID, "chain-id", 0, "L2 chain ID (required)")
	deployContractsCmd.Flags().StringVar(&flagOut, "out", "./deploy-output.json", "Output file path")
	_ = deployContractsCmd.MarkFlagRequired("l1-rpc")
	_ = deployContractsCmd.MarkFlagRequired("private-key")
	_ = deployContractsCmd.MarkFlagRequired("chain-id")
}
