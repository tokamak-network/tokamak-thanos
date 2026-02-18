package sources

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/tokamak-network/tokamak-thanos/op-service/client"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// FollowStatus tracks the follow state of a peer node.
type FollowStatus struct {
	CurrentL1   eth.L1BlockRef `json:"current_l1"`
	SafeL2      eth.L2BlockRef `json:"safe_l2"`
	UnsafeL2    eth.L2BlockRef `json:"unsafe_l2"`
	FinalizedL2 eth.L2BlockRef `json:"finalized_l2"`
	CrossSafe   eth.L2BlockRef `json:"cross_safe_l2"`
}

// FollowClient is a client for following the status of a peer node.
type FollowClient struct {
	rpc client.RPC
	log log.Logger
}

func NewFollowClient(rpc client.RPC, log log.Logger) *FollowClient {
	return &FollowClient{rpc: rpc, log: log}
}

func (c *FollowClient) GetFollowStatus(ctx context.Context) (*FollowStatus, error) {
	var status FollowStatus
	err := c.rpc.CallContext(ctx, &status, "optimism_syncStatus")
	if err != nil {
		return nil, fmt.Errorf("failed to get follow status: %w", err)
	}
	return &status, nil
}

func (c *FollowClient) Close() {
	c.rpc.Close()
}
