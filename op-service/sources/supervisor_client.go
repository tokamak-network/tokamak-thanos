package sources

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/tokamak-network/tokamak-thanos/op-service/client"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// SupervisorClient is a client for the supervisor RPC API.
type SupervisorClient struct {
	rpc client.RPC
	log log.Logger
}

func NewSupervisorClient(rpc client.RPC, log log.Logger) *SupervisorClient {
	return &SupervisorClient{rpc: rpc, log: log}
}

func NewSupervisorClientFromRPCClient(rpcClient *rpc.Client, log log.Logger) *SupervisorClient {
	return NewSupervisorClient(client.NewBaseRPCClient(rpcClient), log)
}

func DialSupervisorClient(ctx context.Context, log log.Logger, url string) (*SupervisorClient, error) {
	rpcClient, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to dial supervisor at %q: %w", url, err)
	}
	return NewSupervisorClientFromRPCClient(rpcClient, log), nil
}

func (c *SupervisorClient) SyncStatus(ctx context.Context) (eth.SupervisorSyncStatus, error) {
	var status eth.SupervisorSyncStatus
	err := c.rpc.CallContext(ctx, &status, "supervisor_syncStatus")
	return status, err
}

func (c *SupervisorClient) SuperRootAtTimestamp(ctx context.Context, timestamp hexutil.Uint64) (eth.SuperRootResponse, error) {
	var resp eth.SuperRootResponse
	err := c.rpc.CallContext(ctx, &resp, "supervisor_superRootAtTimestamp", timestamp)
	return resp, err
}

func (c *SupervisorClient) AllSafeDerivedAt(ctx context.Context, derivedFrom eth.BlockID) (map[eth.ChainID]eth.BlockID, error) {
	var result map[eth.ChainID]eth.BlockID
	err := c.rpc.CallContext(ctx, &result, "supervisor_allSafeDerivedAt", derivedFrom)
	return result, err
}

func (c *SupervisorClient) Close() {
	c.rpc.Close()
}
