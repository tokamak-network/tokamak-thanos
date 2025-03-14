package client

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/contracts/bindings"
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/interfaces"
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/types"
)

type superchainWETHBinding struct {
	contractAddress types.Address
	client          *ethclient.Client
	binding         *bindings.SuperchainWETH
}

var _ interfaces.SuperchainWETH = (*superchainWETHBinding)(nil)

func (b *superchainWETHBinding) BalanceOf(addr types.Address) types.ReadInvocation[types.Balance] {
	return &superchainWETHBalanceOfImpl{
		contract: b,
		addr:     addr,
	}
}

type superchainWETHBalanceOfImpl struct {
	contract *superchainWETHBinding
	addr     types.Address
}

func (i *superchainWETHBalanceOfImpl) Call(ctx context.Context) (types.Balance, error) {
	balance, err := i.contract.binding.BalanceOf(nil, i.addr)
	if err != nil {
		return types.Balance{}, err
	}
	return types.NewBalance(balance), nil
}
