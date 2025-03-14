package constraints

import (
	"github.com/ethereum/go-ethereum/log"

	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/system"
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/types"
)

type WalletConstraint interface {
	CheckWallet(wallet system.Wallet) bool
}

type WalletConstraintFunc func(wallet system.Wallet) bool

func (f WalletConstraintFunc) CheckWallet(wallet system.Wallet) bool {
	return f(wallet)
}

func WithBalance(amount types.Balance) WalletConstraint {
	return WalletConstraintFunc(func(wallet system.Wallet) bool {
		balance := wallet.Balance()
		log.Debug("checking balance", "wallet", wallet.Address(), "balance", balance, "needed", amount)
		return balance.GreaterThan(amount)
	})
}
