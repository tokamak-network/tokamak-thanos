package dial

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-node/rollup"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// RollupClientInterface is an interface for providing a RollupClient
// It does not describe all of the functions a RollupClient has, only the ones used by the L2 Providers and their callers
type RollupClientInterface interface {
	SyncStatusProvider
	OutputAtBlock(ctx context.Context, blockNum uint64) (*eth.OutputResponse, error)
	RollupConfig(ctx context.Context) (*rollup.Config, error)
	StartSequencer(ctx context.Context, unsafeHead common.Hash) error
	SequencerActive(ctx context.Context) (bool, error)
	Close()
}

// SyncStatusProvider is the interface of a rollup client from which its sync status
// can be queried.
type SyncStatusProvider interface {
	SyncStatus(ctx context.Context) (*eth.SyncStatus, error)
}
