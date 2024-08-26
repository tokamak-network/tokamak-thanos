package contracts

import (
	"context"
	_ "embed"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	"github.com/tokamak-network/tokamak-thanos/op-service/txmgr"
)

//go:embed abis/FaultDisputeGame-1.1.1.json
var faultDisputeGameAbi111 []byte

type FaultDisputeGameContract111 struct {
	FaultDisputeGameContractLatest
}

func (f *FaultDisputeGameContract111) AttackTx(ctx context.Context, parent types.Claim, pivot common.Hash) (txmgr.TxCandidate, error) {
	call := f.contract.Call(methodAttack, big.NewInt(int64(parent.ContractIndex)), pivot)
	return f.txWithBond(ctx, parent.Position.Attack(), call)
}

func (f *FaultDisputeGameContract111) DefendTx(ctx context.Context, parent types.Claim, pivot common.Hash) (txmgr.TxCandidate, error) {
	call := f.contract.Call(methodDefend, big.NewInt(int64(parent.ContractIndex)), pivot)
	return f.txWithBond(ctx, parent.Position.Defend(), call)
}
