package fullseq

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-test-sequencer/sequencer/backend/work"
	"github.com/ethereum-optimism/optimism/op-test-sequencer/sequencer/seqtypes"
)

type Config struct {
	ChainID eth.ChainID `yaml:"chainID"`

	Builder   seqtypes.BuilderID   `yaml:"builder"`
	Signer    seqtypes.SignerID    `yaml:"signer,omitempty"`
	Committer seqtypes.CommitterID `yaml:"committer,omitempty"`
	Publisher seqtypes.PublisherID `yaml:"publisher,omitempty"`

	// SequencerConfDepth is the distance to keep from the L1 head as origin when sequencing new L2 blocks.
	// If this distance is too large, the sequencer may:
	// - not adopt a L1 origin within the allowed time (rollup.Config.MaxSequencerDrift)
	// - not adopt a L1 origin that can be included on L1 within the allowed range (rollup.Config.SeqWindowSize)
	// and thus fail to produce a block with anything more than deposits.
	SequencerConfDepth uint64 `json:"sequencer_conf_depth"`

	// SequencerEnabled is true when the sequencer is operational.
	SequencerEnabled bool `json:"sequencer_enabled"`

	// SequencerStopped is false when the sequencer should not be auto-sequencing at startup.
	SequencerStopped bool `json:"sequencer_stopped"`

	// SequencerMaxSafeLag is the maximum number of L2 blocks for restricting the distance between L2 safe and unsafe.
	// Disabled if 0.
	SequencerMaxSafeLag uint64 `json:"sequencer_max_safe_lag"`
}

func (c *Config) Start(ctx context.Context, id seqtypes.SequencerID, opts *work.ServiceOpts) (work.Sequencer, error) {

	builder := opts.Services.Builder(c.Builder)
	signer := opts.Services.Signer(c.Signer)
	committer := opts.Services.Committer(c.Committer)
	publisher := opts.Services.Publisher(c.Publisher)

	// TODO(#14129) load persisted sequencer state (add config var for peristence path + use op-node persistence code)

	seq := &Sequencer{
		id: id,

		chainID: c.ChainID,

		log: opts.Log.New("chain", c.ChainID),
		m:   opts.Metrics,

		builder:   builder,
		signer:    signer,
		committer: committer,
		publisher: publisher,
	}

	// TODO(#14129) check persisted state, to determine if we should really start or stop
	//if c.SequencerEnabled && !c.SequencerStopped {
	//	if err := seq.forceStart(); err != nil {
	//		return nil, fmt.Errorf("failed to start sequencer at startup phase: %w", err)
	//	}
	//}
	return seq, nil
}
