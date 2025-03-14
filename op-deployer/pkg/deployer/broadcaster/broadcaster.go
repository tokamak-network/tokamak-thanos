package broadcaster

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/script"
)

type Broadcaster interface {
	Broadcast(ctx context.Context) ([]BroadcastResult, error)
	Hook(bcast script.Broadcast)
}

type BroadcastResult struct {
	Broadcast script.Broadcast `json:"broadcast"`
	TxHash    common.Hash      `json:"txHash"`
	Receipt   *types.Receipt   `json:"receipt"`
	Err       error            `json:"-"`
}
