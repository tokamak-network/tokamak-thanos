// Package tracing provides shims for go-ethereum/core/tracing types
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

// Hooks is a stub tracing hooks struct for EVM execution tracing.
// In the real geth 1.14+ this is a rich struct; here it's a placeholder.
type Hooks struct{}
