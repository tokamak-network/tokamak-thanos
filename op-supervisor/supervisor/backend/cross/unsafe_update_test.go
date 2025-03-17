package cross

import (
	"errors"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/backend/depset"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func TestCrossUnsafeUpdate(t *testing.T) {
	t.Run("CrossUnsafe returns error", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		usd.crossUnsafeFn = func(chainID eth.ChainID) (types.BlockSeal, error) {
			return types.BlockSeal{}, errors.New("some error")
		}
		usd.deps = mockDependencySet{}
		// when an error is returned by CrossUnsafe,
		// the error is returned
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.ErrorContains(t, err, "some error")
	})
	t.Run("CrossUnsafe returns ErrFuture", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		usd.crossUnsafeFn = func(chainID eth.ChainID) (types.BlockSeal, error) {
			return types.BlockSeal{}, types.ErrFuture
		}
		usd.deps = mockDependencySet{}
		// when a ErrFuture is returned by CrossUnsafe,
		// no error is returned
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.NoError(t, err)
	})
	t.Run("OpenBlock returns error", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		usd.openBlockFn = func(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error) {
			return eth.BlockRef{}, 0, nil, errors.New("some error")
		}
		usd.deps = mockDependencySet{}
		// when an error is returned by OpenBlock,
		// the error is returned
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.ErrorContains(t, err, "some error")
	})
	t.Run("opened block parent hash does not match", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		crossUnsafe := types.BlockSeal{Hash: common.Hash{0x11}}
		usd.crossUnsafeFn = func(chainID eth.ChainID) (types.BlockSeal, error) {
			return crossUnsafe, nil
		}
		bl := eth.BlockRef{ParentHash: common.Hash{0x01}}
		usd.openBlockFn = func(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error) {
			return bl, 0, nil, nil
		}
		usd.deps = mockDependencySet{}
		// when the parent hash of the opened block does not match the cross-unsafe block,
		// an ErrConflict is returned
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.ErrorIs(t, err, types.ErrConflict)
	})
	t.Run("CrossUnsafeHazards returns error", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		crossUnsafe := types.BlockSeal{Hash: common.Hash{0x01}}
		usd.crossUnsafeFn = func(chainID eth.ChainID) (types.BlockSeal, error) {
			return crossUnsafe, nil
		}
		bl := eth.BlockRef{ParentHash: common.Hash{0x01}}
		usd.openBlockFn = func(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error) {
			// include one executing message to trigger the CanExecuteAt check
			return bl, 0, map[uint32]*types.ExecutingMessage{1: {}}, nil
		}
		usd.deps = mockDependencySet{}
		// make CrossSafeHazards return an error by setting CanExecuteAtfn to return an error
		usd.deps.canExecuteAtfn = func() (bool, error) {
			return false, errors.New("some error")
		}
		// when CrossSafeHazards returns an error,
		// the error is returned
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.ErrorContains(t, err, "some error")
	})
	t.Run("HazardUnsafeFrontierChecks returns error", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		crossUnsafe := types.BlockSeal{Hash: common.Hash{0x01}}
		usd.crossUnsafeFn = func(chainID eth.ChainID) (types.BlockSeal, error) {
			return crossUnsafe, nil
		}
		bl := eth.BlockRef{ParentHash: common.Hash{0x01}, Time: 1}
		em1 := &types.ExecutingMessage{Timestamp: 1}
		usd.openBlockFn = func(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error) {
			// include one executing message to ensure one hazard is returned
			return bl, 0, map[uint32]*types.ExecutingMessage{1: em1}, nil
		}
		usd.deps = mockDependencySet{}
		count := 0
		// make HazardUnsafeFrontierChecks return an error by failing the third ChainIDFromIndex call
		// (the first two are in CrossSafeHazards)
		usd.deps.chainIDFromIndexfn = func() (eth.ChainID, error) {
			defer func() { count++ }()
			if count < 2 {
				return eth.ChainID{}, errors.New("some error")
			}
			return eth.ChainID{}, nil
		}
		// when HazardUnsafeFrontierChecks returns an error,
		// the error is returned
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.ErrorContains(t, err, "some error")
	})
	t.Run("HazardCycleChecks returns error", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		crossUnsafe := types.BlockSeal{Hash: common.Hash{0x01}}
		usd.crossUnsafeFn = func(chainID eth.ChainID) (types.BlockSeal, error) {
			return crossUnsafe, nil
		}
		bl := eth.BlockRef{ParentHash: common.Hash{0x01}, Number: 1, Time: 1}
		em1 := &types.ExecutingMessage{Chain: types.ChainIndex(0), Timestamp: 1, LogIdx: 2}
		em2 := &types.ExecutingMessage{Chain: types.ChainIndex(0), Timestamp: 1, LogIdx: 1}
		usd.openBlockFn = func(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error) {
			return bl, 3, map[uint32]*types.ExecutingMessage{1: em1, 2: em2}, nil
		}
		usd.checkFn = func(chainID eth.ChainID, blockNum uint64, timestamp uint64, logIdx uint32, logHash common.Hash) (types.BlockSeal, error) {
			return types.BlockSeal{Number: 1, Timestamp: 1}, nil
		}
		usd.deps = mockDependencySet{}

		// HazardCycleChecks returns an error with appropriate wrapping
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.ErrorContains(t, err, "cycle detected")
		require.ErrorContains(t, err, "failed to verify block")
	})
	t.Run("successful update", func(t *testing.T) {
		logger := testlog.Logger(t, log.LevelDebug)
		chainID := eth.ChainIDFromUInt64(0)
		usd := &mockCrossUnsafeDeps{}
		crossUnsafe := types.BlockSeal{Hash: common.Hash{0x01}}
		usd.crossUnsafeFn = func(chainID eth.ChainID) (types.BlockSeal, error) {
			return crossUnsafe, nil
		}
		bl := eth.BlockRef{ParentHash: common.Hash{0x01}, Time: 1, Number: 1}
		em1 := &types.ExecutingMessage{Timestamp: 1}
		usd.openBlockFn = func(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error) {
			// include one executing message to ensure one hazard is returned
			if blockNum == 1 {
				return bl, 2, map[uint32]*types.ExecutingMessage{1: em1}, nil
			}
			return eth.BlockRef{
				Hash: crossUnsafe.Hash,
			}, 0, nil, nil
		}
		usd.checkFn = func(chainID eth.ChainID, blockNum uint64, timestamp uint64, logIdx uint32, logHash common.Hash) (types.BlockSeal, error) {
			return crossUnsafe, nil
		}
		usd.deps = mockDependencySet{}
		var updatingChainID eth.ChainID
		var updatingBlock types.BlockSeal
		usd.updateCrossUnsafeFn = func(chain eth.ChainID, crossUnsafe types.BlockSeal) error {
			updatingChainID = chain
			updatingBlock = crossUnsafe
			return nil
		}
		// when there are no errors, the cross-unsafe block is updated
		// the updated block is the block opened in OpenBlock
		err := CrossUnsafeUpdate(logger, chainID, usd)
		require.NoError(t, err)
		require.Equal(t, chainID, updatingChainID)
		require.Equal(t, types.BlockSealFromRef(bl), updatingBlock)
	})
}

type mockCrossUnsafeDeps struct {
	deps                mockDependencySet
	messageExpiryWindow uint64
	crossUnsafeFn       func(chainID eth.ChainID) (types.BlockSeal, error)
	openBlockFn         func(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error)
	updateCrossUnsafeFn func(chain eth.ChainID, crossUnsafe types.BlockSeal) error
	checkFn             func(chainID eth.ChainID, blockNum uint64, timestamp uint64, logIdx uint32, logHash common.Hash) (types.BlockSeal, error)
}

func (m *mockCrossUnsafeDeps) CrossUnsafe(chainID eth.ChainID) (derived types.BlockSeal, err error) {
	if m.crossUnsafeFn != nil {
		return m.crossUnsafeFn(chainID)
	}
	return types.BlockSeal{}, nil
}

func (m *mockCrossUnsafeDeps) DependencySet() depset.DependencySet {
	return m.deps
}

func (m *mockCrossUnsafeDeps) MessageExpiryWindow() uint64 {
	return m.messageExpiryWindow
}

func (m *mockCrossUnsafeDeps) Contains(chainID eth.ChainID, q types.ContainsQuery) (types.BlockSeal, error) {
	if m.checkFn != nil {
		return m.checkFn(chainID, q.BlockNum, q.Timestamp, q.LogIdx, q.LogHash)
	}
	return types.BlockSeal{}, nil
}

func (m *mockCrossUnsafeDeps) OpenBlock(chainID eth.ChainID, blockNum uint64) (ref eth.BlockRef, logCount uint32, execMsgs map[uint32]*types.ExecutingMessage, err error) {
	if m.openBlockFn != nil {
		return m.openBlockFn(chainID, blockNum)
	}
	return eth.BlockRef{}, 0, nil, nil
}

func (m *mockCrossUnsafeDeps) UpdateCrossUnsafe(chain eth.ChainID, block types.BlockSeal) error {
	if m.updateCrossUnsafeFn != nil {
		return m.updateCrossUnsafeFn(chain, block)
	}
	return nil
}

func (m *mockCrossUnsafeDeps) IsCrossUnsafe(chainID eth.ChainID, blockNum eth.BlockID) error {
	return nil
}

func (m *mockCrossUnsafeDeps) IsLocalUnsafe(chainID eth.ChainID, blockNum eth.BlockID) error {
	return nil
}

func (m *mockCrossUnsafeDeps) ParentBlock(chainID eth.ChainID, blockNum eth.BlockID) (eth.BlockID, error) {
	return eth.BlockID{}, nil
}
