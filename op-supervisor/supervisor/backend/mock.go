package backend

import (
	"context"
	"errors"
	"io"
	"sync/atomic"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/frontend"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MockBackend struct {
	started atomic.Bool
}

var _ frontend.Backend = (*MockBackend)(nil)

var _ io.Closer = (*MockBackend)(nil)

func NewMockBackend() *MockBackend {
	return &MockBackend{}
}

func (m *MockBackend) Start(ctx context.Context) error {
	if !m.started.CompareAndSwap(false, true) {
		return errors.New("already started")
	}
	return nil
}

func (m *MockBackend) Stop(ctx context.Context) error {
	if !m.started.CompareAndSwap(true, false) {
		return errors.New("already stopped")
	}
	return nil
}

func (m *MockBackend) AllSafeDerivedAt(ctx context.Context, source eth.BlockID) (derived map[eth.ChainID]eth.BlockID, err error) {
	return nil, nil
}

func (m *MockBackend) AddL2RPC(ctx context.Context, rpc string, jwtSecret eth.Bytes32) error {
	return nil
}

func (m *MockBackend) CheckMessage(ctx context.Context, identifier types.Identifier, payloadHash common.Hash, executingDescriptor types.ExecutingDescriptor) (types.SafetyLevel, error) {
	return types.CrossUnsafe, nil
}

func (m *MockBackend) CheckMessages(ctx context.Context, messages []types.Message, minSafety types.SafetyLevel) error {
	return nil
}

func (m *MockBackend) CheckMessagesV2(ctx context.Context, messages []types.Message, minSafety types.SafetyLevel, executingDescriptor types.ExecutingDescriptor) error {
	return nil
}

func (m *MockBackend) LocalUnsafe(ctx context.Context, chainID eth.ChainID) (eth.BlockID, error) {
	return eth.BlockID{}, nil
}

func (m *MockBackend) CrossSafe(ctx context.Context, chainID eth.ChainID) (types.DerivedIDPair, error) {
	return types.DerivedIDPair{}, nil
}

func (m *MockBackend) Finalized(ctx context.Context, chainID eth.ChainID) (eth.BlockID, error) {
	return eth.BlockID{}, nil
}

func (m *MockBackend) FinalizedL1(ctx context.Context) (eth.BlockRef, error) {
	return eth.BlockRef{}, nil
}

func (m *MockBackend) CrossDerivedToSource(ctx context.Context, chainID eth.ChainID, derived eth.BlockID) (source eth.BlockRef, err error) {
	return eth.BlockRef{}, nil
}

func (m *MockBackend) SuperRootAtTimestamp(ctx context.Context, timestamp hexutil.Uint64) (eth.SuperRootResponse, error) {
	return eth.SuperRootResponse{}, nil
}

func (m *MockBackend) SyncStatus(ctx context.Context) (eth.SupervisorSyncStatus, error) {
	return eth.SupervisorSyncStatus{}, nil
}

func (m *MockBackend) Close() error {
	return nil
}
