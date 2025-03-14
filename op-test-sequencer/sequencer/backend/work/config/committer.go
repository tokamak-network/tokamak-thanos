package config

import (
	"context"

	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer/backend/work"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer/backend/work/committers/noopcommitter"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer/seqtypes"
)

type CommitterEntry struct {
	Noop *noopcommitter.Config `yaml:"noop,omitempty"`
}

func (b *CommitterEntry) Start(ctx context.Context, id seqtypes.CommitterID, opts *work.ServiceOpts) (work.Committer, error) {
	switch {
	case b.Noop != nil:
		return b.Noop.Start(ctx, id, opts)
	default:
		return nil, seqtypes.ErrUnknownKind
	}
}
