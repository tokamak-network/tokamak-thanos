package devnet

import (
	"fmt"
	"log/slog"
)

// NewForked starts an anvil instance forked from the given RPC URL.
// Stub: not implemented in tokamak-thanos.
func NewForked(lgr *slog.Logger, rpcURL string) (string, func(), error) {
	return "", nil, fmt.Errorf("devnet.NewForked not implemented")
}

// NewForkedSepolia starts an anvil instance forked from Sepolia.
func NewForkedSepolia(lgr *slog.Logger) (string, func(), error) {
	return "", nil, fmt.Errorf("devnet.NewForkedSepolia not implemented")
}
