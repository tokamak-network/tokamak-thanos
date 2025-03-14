package genesis

import (
	"fmt"
	"math/big"

	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/foundry"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis/beacondeposit"
	"github.com/tokamak-network/tokamak-thanos/op-service/predeploys"
)

// BuildL1DeveloperGenesis will create a L1 genesis block after creating
// all of the state required for an Optimism network to function.
// It is expected that the dump contains all of the required state to bootstrap
// the L1 chain.
func BuildL1DeveloperGenesis(config *DeployConfig, dump *foundry.ForgeAllocs, l1Deployments *L1Deployments) (*core.Genesis, error) {
	log.Info("Building developer L1 genesis block")
	genesis, err := NewL1Genesis(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create L1 developer genesis: %w", err)
	}

	if len(genesis.Alloc) != 0 {
		panic("Did not expect NewL1Genesis to generate non-empty state") // sanity check for dev purposes.
	}

	// copy, for safety when the dump is reused (like in e2e testing)
	genesis.Alloc = dump.Copy().Accounts
	if config.FundDevAccounts {
		FundDevAccounts(genesis)
	}
	SetPrecompileBalances(genesis)

	l1Deployments.ForEach(func(name string, addr common.Address) {
		acc, ok := genesis.Alloc[addr]
		if ok {
			log.Info("Included L1 deployment", "name", name, "address", addr, "balance", acc.Balance, "storage", len(acc.Storage), "nonce", acc.Nonce)
		} else {
			log.Info("Excluded L1 deployment", "name", name, "address", addr)
		}
	})

	beaconDepositAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	if err := beacondeposit.InsertEmptyBeaconDepositContract(genesis, beaconDepositAddr); err != nil {
		return nil, fmt.Errorf("failed to insert beacon deposit contract into L1 dev genesis: %w", err)
	}

	// For 4788, make sure the 4788 beacon-roots contract is there.
	// (required to be there before L1 Dencun activation)
	genesis.Alloc[predeploys.EIP4788ContractAddr] = types.Account{
		Balance: new(big.Int),
		Nonce:   1,
		Code:    predeploys.EIP4788ContractCode,
	}
	// Also record the virtual deployer address
	genesis.Alloc[predeploys.EIP4788ContractDeployer] = types.Account{
		Balance: new(big.Int),
		Nonce:   1,
	}

	return genesis, nil
}
