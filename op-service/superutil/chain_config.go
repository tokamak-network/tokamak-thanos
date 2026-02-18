package superutil

import (
	"fmt"

	"github.com/ethereum/go-ethereum/params"
)

func LoadOPStackChainConfigFromChainID(chainID uint64) (*params.ChainConfig, error) {
	cfg, err := params.LoadOPStackChainConfig(chainID)
	if err != nil {
		return nil, fmt.Errorf("unable to load chain config for chain %d: %w", chainID, err)
	}
	return cfg, nil
}
