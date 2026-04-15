package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tokamak-deployer",
	Short: "L1 contract deployer and genesis generator for tokamak-thanos",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(deployContractsCmd)
	rootCmd.AddCommand(generateGenesisCmd)
}
