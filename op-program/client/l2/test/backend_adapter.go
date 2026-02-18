package test

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/tokamak-network/tokamak-thanos/op-service/compat/stateless"
)

// BlockChainAdapter wraps *core.BlockChain to satisfy the EngineBackend interface
// which expects the new-style InsertBlockWithoutSetHead signature.
type BlockChainAdapter struct {
	*core.BlockChain
}

func (a *BlockChainAdapter) InsertBlockWithoutSetHead(block *types.Block, makeWitness bool) (*stateless.Witness, error) {
	err := a.BlockChain.InsertBlockWithoutSetHead(block)
	return nil, err
}
