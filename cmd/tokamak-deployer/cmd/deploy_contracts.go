package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"
	"github.com/tokamak-network/tokamak-thanos/cmd/tokamak-deployer/internal/deployer"
)

var (
	flagL1RPC              string
	flagPrivateKey         string
	flagChainID            uint64
	flagOut                string
	flagGasPrice           string
	flagGasPriceMultiplier int
	flagGasPriceFloor      string
	flagGasPriceCeil       string
)

const envGasPrice = "TOKAMAK_DEPLOY_GAS_PRICE"

var deployContractsCmd = &cobra.Command{
	Use:   "deploy-contracts",
	Short: "Deploy L1 contracts and write deploy-output.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		fixedGasPrice, err := parseWeiFlag(flagGasPrice, envGasPrice)
		if err != nil {
			return fmt.Errorf("--gas-price: %w", err)
		}
		floor, err := parseWeiFlag(flagGasPriceFloor, "")
		if err != nil {
			return fmt.Errorf("--gas-price-floor: %w", err)
		}
		ceil, err := parseWeiFlag(flagGasPriceCeil, "")
		if err != nil {
			return fmt.Errorf("--gas-price-ceil: %w", err)
		}

		cfg := deployer.DeployConfig{
			L1RPCURL:           flagL1RPC,
			PrivateKey:         flagPrivateKey,
			L2ChainID:          flagChainID,
			FixedGasPrice:      fixedGasPrice,
			GasPriceMultiplier: flagGasPriceMultiplier,
			GasPriceFloor:      floor,
			GasPriceCeil:       ceil,
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

// parseWeiFlag converts a decimal wei string to *big.Int. Empty string + empty
// envVar returns nil (no override). If the flag is empty, the env var (when
// provided) is used as a fallback — mirrors the OLD forge path's
// TOKAMAK_DEPLOY_GAS_PRICE pattern.
func parseWeiFlag(flag, envVar string) (*big.Int, error) {
	s := flag
	if s == "" && envVar != "" {
		s = os.Getenv(envVar)
	}
	if s == "" {
		return nil, nil
	}
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid wei value %q (must be a decimal integer)", s)
	}
	if v.Sign() < 0 {
		return nil, fmt.Errorf("wei value must be non-negative, got %s", v)
	}
	return v, nil
}

func init() {
	deployContractsCmd.Flags().StringVar(&flagL1RPC, "l1-rpc", "", "L1 RPC URL (required)")
	deployContractsCmd.Flags().StringVar(&flagPrivateKey, "private-key", "", "Deployer private key (required)")
	deployContractsCmd.Flags().Uint64Var(&flagChainID, "chain-id", 0, "L2 chain ID (required)")
	deployContractsCmd.Flags().StringVar(&flagOut, "out", "./deploy-output.json", "Output file path")
	deployContractsCmd.Flags().StringVar(&flagGasPrice, "gas-price", "",
		"Fixed gas price in wei reused for every TX (empty = auto; overridable via TOKAMAK_DEPLOY_GAS_PRICE)")
	deployContractsCmd.Flags().IntVar(&flagGasPriceMultiplier, "gas-price-multiplier", 0,
		"Percent of SuggestGasPrice to use when --gas-price is not set (default 200 = 2x)")
	deployContractsCmd.Flags().StringVar(&flagGasPriceFloor, "gas-price-floor", "",
		"Minimum gas price in wei applied to the resolved price (default 1 Gwei)")
	deployContractsCmd.Flags().StringVar(&flagGasPriceCeil, "gas-price-ceil", "",
		"Maximum gas price in wei applied to the resolved price (default 100 Gwei)")
	_ = deployContractsCmd.MarkFlagRequired("l1-rpc")
	_ = deployContractsCmd.MarkFlagRequired("private-key")
	_ = deployContractsCmd.MarkFlagRequired("chain-id")
}
