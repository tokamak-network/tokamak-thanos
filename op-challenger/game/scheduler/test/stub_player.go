package test

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/types"
)

type StubGamePlayer struct {
	Addr          common.Address
	ProgressCount int
	StatusValue   types.GameStatus
	Dir           string
}

func (g *StubGamePlayer) ProgressGame(_ context.Context) types.GameStatus {
	g.ProgressCount++
	return g.StatusValue
}

func (g *StubGamePlayer) Status() types.GameStatus {
	return g.StatusValue
}
