package eth

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

// Bytes8 is an 8-byte array.
type Bytes8 [8]byte

func (b *Bytes8) UnmarshalJSON(text []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(b), text, b[:])
}

func (b *Bytes8) UnmarshalText(text []byte) error {
	return hexutil.UnmarshalFixedText("Bytes8", text, b[:])
}

func (b Bytes8) MarshalText() ([]byte, error) {
	return hexutil.Bytes(b[:]).MarshalText()
}

func (b Bytes8) String() string {
	return hexutil.Encode(b[:])
}

func (b Bytes8) TerminalString() string {
	return fmt.Sprintf("%x", b[:])
}

// ExecutionWitness represents block execution witness data.
// In upstream op-geth this is types.ExecutionWitness; stub for tokamak-thanos-geth compat.
type ExecutionWitness struct {
	State   map[string]string `json:"state,omitempty"`
	Headers []byte            `json:"headers,omitempty"`
}

// SupervisorSyncStatus represents the sync status of the supervisor.
type SupervisorSyncStatus struct {
	MinSyncedL1        BlockID `json:"minSyncedL1"`
	FinalizedTimestamp uint64  `json:"finalizedTimestamp"`
}

// ETH represents an amount of ETH (in wei).
type ETH uint64

// SyncTesterSession is used for sync testing.
type SyncTesterSession struct {
	ID        string  `json:"id"`
	Status    string  `json:"status"`
}

// Bytes65 is a 65-byte array (e.g. for signatures).
type Bytes65 [65]byte


// IsDepositsOnly returns true if this payload only has deposit transactions.
func (attrs *PayloadAttributes) IsDepositsOnly() bool {
	return attrs.NoTxPool
}

// BlockRef returns the L1BlockRef-equivalent from an L2BlockRef.
func (ref L2BlockRef) BlockRef() L1BlockRef {
	return L1BlockRef{
		Hash:       ref.Hash,
		Number:     ref.Number,
		ParentHash: ref.ParentHash,
		Time:       ref.Time,
	}
}

// IsEngineError returns true if this is an engine API error.
func (c ErrorCode) IsEngineError() bool {
	return int(c) >= -38099 && int(c) <= -38000
}

// WithdrawalsRoot computes the withdrawals root from the payload's withdrawals.
func (payload *ExecutionPayload) WithdrawalsRoot() *common.Hash {
	if payload.Withdrawals == nil {
		return nil
	}
	h := types.DeriveSha(*payload.Withdrawals, trie.NewStackTrie(nil))
	return &h
}

// BlockRef returns a BlockRef from the ExecutionPayload.
func (p *ExecutionPayload) BlockRef() L1BlockRef {
	return L1BlockRef{
		Hash:       p.BlockHash,
		Number:     uint64(p.BlockNumber),
		ParentHash: p.ParentHash,
		Time:       uint64(p.Timestamp),
	}
}
