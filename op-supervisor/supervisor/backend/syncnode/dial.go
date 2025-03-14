package syncnode

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	gn "github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/tokamak-network/tokamak-thanos/op-service/client"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

type RPCDialSetup struct {
	JWTSecret eth.Bytes32
	Endpoint  string
}

var _ SyncNodeSetup = (*RPCDialSetup)(nil)

func (r *RPCDialSetup) Setup(ctx context.Context, logger log.Logger) (SyncNode, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	auth := rpc.WithHTTPAuth(gn.NewJWTAuth(r.JWTSecret))
	opts := []client.RPCOption{
		client.WithGethRPCOptions(auth),
		client.WithDialAttempts(10),
	}
	rpcCl, err := client.NewRPC(ctx, logger, r.Endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return &RPCSyncNode{
		name: fmt.Sprintf("RPCSyncSource(%s)", r.Endpoint),
		cl:   rpcCl,
	}, nil
}
