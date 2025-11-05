package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

// L1ChainConfigByChainID returns the chain config for the given chain ID,
// if it is in the set of known chain IDs (Mainnet, Sepolia, Holesky).
// If the chain ID is not known, it returns nil.
func L1ChainConfigByChainID(chainID *big.Int) *params.ChainConfig {
	if chainID == nil {
		return nil
	}

	// Check known L1 chains
	if chainID.Cmp(params.MainnetChainConfig.ChainID) == 0 {
		return params.MainnetChainConfig
	}
	if chainID.Cmp(params.SepoliaChainConfig.ChainID) == 0 {
		return params.SepoliaChainConfig
	}
	// Holesky check - may not exist in older geth versions
	if params.HoleskyChainConfig != nil && chainID.Cmp(params.HoleskyChainConfig.ChainID) == 0 {
		return params.HoleskyChainConfig
	}

	return nil
}

