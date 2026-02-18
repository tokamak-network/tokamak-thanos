package eth

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
