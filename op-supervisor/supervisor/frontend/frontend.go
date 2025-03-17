package frontend

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/types"
)

type Backend interface {
	sources.SupervisorAdminAPI
	sources.SupervisorQueryAPI
}

type QueryFrontend struct {
	Supervisor sources.SupervisorQueryAPI
}

var _ sources.SupervisorQueryAPI = (*QueryFrontend)(nil)

// CheckMessage checks the safety-level of an individual message.
// The payloadHash references the hash of the message-payload of the message.
func (q *QueryFrontend) CheckMessage(ctx context.Context, identifier types.Identifier, payloadHash common.Hash, executingDescriptor types.ExecutingDescriptor) (types.SafetyLevel, error) {
	return q.Supervisor.CheckMessage(ctx, identifier, payloadHash, executingDescriptor)
}

// CheckMessagesV2 checks the safety-level of a collection of messages,
// and returns if the minimum safety-level is met for all messages.
func (q *QueryFrontend) CheckMessagesV2(
	ctx context.Context,
	messages []types.Message,
	minSafety types.SafetyLevel,
	executingDescriptor types.ExecutingDescriptor) error {
	return q.Supervisor.CheckMessagesV2(ctx, messages, minSafety, executingDescriptor)
}

// CheckMessages checks the safety-level of a collection of messages,
// and returns if the minimum safety-level is met for all messages.
// Deprecated: This method does not check for message expiry.
func (q *QueryFrontend) CheckMessages(
	ctx context.Context,
	messages []types.Message,
	minSafety types.SafetyLevel) error {
	return q.Supervisor.CheckMessages(ctx, messages, minSafety)
}

func (q *QueryFrontend) LocalUnsafe(ctx context.Context, chainID eth.ChainID) (eth.BlockID, error) {
	return q.Supervisor.LocalUnsafe(ctx, chainID)
}

func (q *QueryFrontend) CrossSafe(ctx context.Context, chainID eth.ChainID) (types.DerivedIDPair, error) {
	return q.Supervisor.CrossSafe(ctx, chainID)
}

func (q *QueryFrontend) Finalized(ctx context.Context, chainID eth.ChainID) (eth.BlockID, error) {
	return q.Supervisor.Finalized(ctx, chainID)
}

func (q *QueryFrontend) FinalizedL1(ctx context.Context) (eth.BlockRef, error) {
	return q.Supervisor.FinalizedL1(ctx)
}

// CrossDerivedFrom is deprecated, but remains for backwards compatibility to callers
// it is equivalent to CrossDerivedToSource
func (q *QueryFrontend) CrossDerivedFrom(ctx context.Context, chainID eth.ChainID, derived eth.BlockID) (derivedFrom eth.BlockRef, err error) {
	return q.Supervisor.CrossDerivedToSource(ctx, chainID, derived)
}

func (q *QueryFrontend) CrossDerivedToSource(ctx context.Context, chainID eth.ChainID, derived eth.BlockID) (derivedFrom eth.BlockRef, err error) {
	return q.Supervisor.CrossDerivedToSource(ctx, chainID, derived)
}

func (q *QueryFrontend) SuperRootAtTimestamp(ctx context.Context, timestamp hexutil.Uint64) (eth.SuperRootResponse, error) {
	return q.Supervisor.SuperRootAtTimestamp(ctx, timestamp)
}

func (q *QueryFrontend) AllSafeDerivedAt(ctx context.Context, derivedFrom eth.BlockID) (derived map[eth.ChainID]eth.BlockID, err error) {
	return q.Supervisor.AllSafeDerivedAt(ctx, derivedFrom)
}

func (q *QueryFrontend) SyncStatus(ctx context.Context) (eth.SupervisorSyncStatus, error) {
	return q.Supervisor.SyncStatus(ctx)
}

type AdminFrontend struct {
	Supervisor Backend
}

var _ sources.SupervisorAdminAPI = (*AdminFrontend)(nil)

// Start starts the service, if it was previously stopped.
func (a *AdminFrontend) Start(ctx context.Context) error {
	return a.Supervisor.Start(ctx)
}

// Stop stops the service, if it was previously started.
func (a *AdminFrontend) Stop(ctx context.Context) error {
	return a.Supervisor.Stop(ctx)
}

// AddL2RPC adds a new L2 chain to the supervisor backend
func (a *AdminFrontend) AddL2RPC(ctx context.Context, rpc string, jwtSecret eth.Bytes32) error {
	return a.Supervisor.AddL2RPC(ctx, rpc, jwtSecret)
}
