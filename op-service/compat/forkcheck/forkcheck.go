// Package forkcheck provides Optimism fork-check helper functions for ChainConfig
// types that may not have the newer fork methods (Holocene, Isthmus, Jovian).
// These forks are not yet supported in the tokamak-thanos-geth fork.
package forkcheck

import (
	"github.com/ethereum/go-ethereum/params"
)

// IsHolocene checks if the Holocene fork is active. Always false for tokamak-thanos-geth.
func IsHolocene(_ *params.ChainConfig, _ uint64) bool { return false }

// IsIsthmus checks if the Isthmus fork is active. Always false for tokamak-thanos-geth.
func IsIsthmus(_ *params.ChainConfig, _ uint64) bool { return false }

// IsJovian checks if the Jovian fork is active. Always false for tokamak-thanos-geth.
func IsJovian(_ *params.ChainConfig, _ uint64) bool { return false }

// IsOptimismHolocene checks if OP Holocene is active. Always false for tokamak-thanos-geth.
func IsOptimismHolocene(_ *params.ChainConfig, _ uint64) bool { return false }

// IsOptimismJovian checks if OP Jovian is active. Always false for tokamak-thanos-geth.
func IsOptimismJovian(_ *params.ChainConfig, _ uint64) bool { return false }

// BaseFeeChangeDenominator returns the EIP1559 denominator for the chain.
func BaseFeeChangeDenominator(c *params.ChainConfig, _ uint64) uint64 {
	if c.Optimism != nil && c.Optimism.EIP1559Denominator != 0 {
		return c.Optimism.EIP1559Denominator
	}
	return 8 // default EIP-1559 denominator
}

// ElasticityMultiplier returns the EIP1559 elasticity multiplier.
func ElasticityMultiplier(c *params.ChainConfig) uint64 {
	if c.Optimism != nil && c.Optimism.EIP1559Elasticity != 0 {
		return c.Optimism.EIP1559Elasticity
	}
	return 2 // default
}
