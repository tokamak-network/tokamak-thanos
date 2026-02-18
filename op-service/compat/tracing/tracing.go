// Package tracing provides shims for go-ethereum/core/tracing constants
// that don't exist in the tokamak-thanos-geth fork (based on geth 1.13).
package tracing

// BalanceChangeReason is used by StateDB.SetBalance in newer geth.
type BalanceChangeReason byte

const (
	BalanceChangeUnspecified BalanceChangeReason = 0
)

// CodeChangeReason is used by StateDB.SetCode in newer geth.
type CodeChangeReason byte

const (
	CodeChangeUnspecified CodeChangeReason = 0
)

// NonceChangeReason is used by StateDB.SetNonce in newer geth.
type NonceChangeReason byte

const (
	NonceChangeUnspecified NonceChangeReason = 0
)

// OpContext provides context about the current EVM operation.
type OpContext interface {
	MemoryData() []byte
	StackData() []uint256
	Address() common.Address
	CallValue() *uint256.Int
	CallInput() []byte
}

// Hooks is a tracing hooks struct for EVM execution tracing.
type Hooks struct {
	OnTxStart   func(vm *VMContext, tx *types.Transaction, from common.Address)
	OnTxEnd     func(receipt *types.Receipt, err error)
	OnEnter     func(depth int, typ byte, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int)
	OnExit      func(depth int, output []byte, gasUsed uint64, err error, reverted bool)
	OnOpcode    func(pc uint64, op byte, gas, cost uint64, scope OpContext, rData []byte, depth int, err error)
	OnFault     func(pc uint64, op byte, gas, cost uint64, scope OpContext, depth int, err error)
	OnLog       func(log *types.Log)
}

// VMContext provides context about the executing VM.
type VMContext struct {
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        uint64
	Random      *common.Hash
	StateDB     StateDB
}

// StateDB is a minimal interface for state access.
type StateDB interface {
	GetBalance(addr common.Address) *uint256.Int
}
