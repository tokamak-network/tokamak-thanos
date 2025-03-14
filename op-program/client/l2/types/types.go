package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// OracleHinter is an optional interface that can be implemented by an Oracle to provide hints
// to access state preimages. This interface only implements hints that are sent proactively
// instead of in preparation for a specific request.
type OracleHinter interface {
	HintBlockExecution(parentBlockHash common.Hash, attr eth.PayloadAttributes, chainID eth.ChainID)
	HintWithdrawalsRoot(blockHash common.Hash, chainID eth.ChainID)
}
