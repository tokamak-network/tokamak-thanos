package standardpublisher

import (
	"context"

	"github.com/tokamak-network/tokamak-thanos/op-service/client"
	"github.com/tokamak-network/tokamak-thanos/op-service/endpoint"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer/backend/work"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer/seqtypes"
)

type Config struct {
	// RPC to publish block to using op-stack RPC
	RPC endpoint.MustRPC `yaml:"rpc"`
}

func (c *Config) Start(ctx context.Context, id seqtypes.PublisherID, opts *work.ServiceOpts) (work.Publisher, error) {
	rpcCl, err := client.NewRPC(ctx, opts.Log, c.RPC.Value.RPC(), client.WithLazyDial())
	if err != nil {
		return nil, err
	}
	cl := sources.NewOPStackClient(rpcCl)
	return &Publisher{
		id:      id,
		log:     opts.Log,
		api:     cl,
		onClose: rpcCl.Close,
	}, nil
}
