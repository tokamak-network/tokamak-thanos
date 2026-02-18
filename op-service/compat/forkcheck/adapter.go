package forkcheck

import "github.com/ethereum/go-ethereum/params"

// ForkChecker wraps a ChainConfig to implement the eip1559.ForkChecker interface.
type ForkChecker struct {
	C *params.ChainConfig
}

func (f ForkChecker) IsHolocene(time uint64) bool { return IsHolocene(f.C, time) }
func (f ForkChecker) IsJovian(time uint64) bool    { return IsJovian(f.C, time) }

// Wrap returns a ForkChecker wrapping the given ChainConfig.
func Wrap(c *params.ChainConfig) ForkChecker { return ForkChecker{C: c} }
