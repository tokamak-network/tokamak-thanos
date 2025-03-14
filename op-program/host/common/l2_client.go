package common

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-service/client"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources/caching"
)

type L2Client struct {
	*sources.L2Client
}

type L2ClientConfig struct {
	*sources.L2ClientConfig
}

func NewL2Client(client client.RPC, log log.Logger, metrics caching.Metrics, config *L2ClientConfig) (*L2Client, error) {
	l2Client, err := sources.NewL2Client(client, log, metrics, config.L2ClientConfig)
	if err != nil {
		return nil, err
	}
	return &L2Client{
		L2Client: l2Client,
	}, nil
}

func (s *L2Client) OutputByRoot(ctx context.Context, blockRoot common.Hash) (eth.Output, error) {
	return s.OutputV0AtBlock(ctx, blockRoot)
}

func (s *L2Client) OutputByNumber(ctx context.Context, blockNum uint64) (eth.Output, error) {
	return s.OutputV0AtBlockNumber(ctx, blockNum)
}
