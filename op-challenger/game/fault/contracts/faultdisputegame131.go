package contracts

import (
	"context"
	_ "embed"

	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources/batching/rpcblock"
)

//go:embed abis/FaultDisputeGame-1.3.1.json
var faultDisputeGameAbi131 []byte

type FaultDisputeGameContract131 struct {
	FaultDisputeGameContractLatest
}

func (f *FaultDisputeGameContract131) GetBondDistributionMode(ctx context.Context, block rpcblock.Block) (types.BondDistributionMode, error) {
	return types.LegacyDistributionMode, nil
}
