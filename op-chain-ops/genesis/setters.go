package genesis

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"

	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/predeploys"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/immutables"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/state"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// PrecompileCount represents the number of precompile addresses
// starting from `address(0)` to PrecompileCount that are funded
// with a single wei in the genesis state.
const PrecompileCount = 256

// FundDevAccounts will fund each of the development accounts.
func FundDevAccounts(gen *core.Genesis) {
	for _, account := range DevAccounts {
		acc := gen.Alloc[account]
		if acc.Balance == nil {
			acc.Balance = new(big.Int)
		}
		acc.Balance = acc.Balance.Add(acc.Balance, devBalance)
		gen.Alloc[account] = acc
	}
}

func setProxies(db vm.StateDB, proxyAdminAddr common.Address, namespace *big.Int, count uint64) error {
	depBytecode, err := bindings.GetDeployedBytecode("Proxy")
	if err != nil {
		return err
	}
	if len(depBytecode) == 0 {
		return errors.New("the contract 'Proxy' has empty bytecode")
	}

	l2UsdcBridgeProxyBytecode, err := bindings.GetDeployedBytecode("L2UsdcBridgeProxy")
	if err != nil {
		return err
	}
	if len(l2UsdcBridgeProxyBytecode) == 0 {
		return errors.New("the contract L2UsdcBridgeProxy has empty bytecode")
	}

	fiatTokenProxyBytecode, err := bindings.GetDeployedBytecode("FiatTokenProxy")
	if err != nil {
		return err
	}
	if len(fiatTokenProxyBytecode) == 0 {
		return errors.New("the contract FiatTokenProxy has empty bytecode")
	}
	TransparentUpgradeableProxyBytecode, err := bindings.GetDeployedBytecode("TransparentUpgradeableProxy")
	if err != nil {
		return err
	}
	if len(TransparentUpgradeableProxyBytecode) == 0 {
		return errors.New("the contract TransparentUpgradeableProxy has empty bytecode")
	}

	for i := uint64(0); i <= count; i++ {
		bigAddr := new(big.Int).Or(namespace, new(big.Int).SetUint64(i))
		addr := common.BigToAddress(bigAddr)

		if !db.Exist(addr) {
			db.CreateAccount(addr)
		}

		switch addr {
		case predeploys.L2UsdcBridgeAddr:
			db.SetCode(addr, l2UsdcBridgeProxyBytecode)
			db.SetState(addr, AdminSlot, eth.AddressAsLeftPaddedHash(proxyAdminAddr))
		case predeploys.FiatTokenV2_2Addr:
			db.SetCode(addr, fiatTokenProxyBytecode)
			db.SetState(addr, AdminSlotForZepplin, eth.AddressAsLeftPaddedHash(proxyAdminAddr))
		case predeploys.NonfungibleTokenPositionDescriptorAddr:
			db.SetCode(addr, TransparentUpgradeableProxyBytecode)
			db.SetState(addr, AdminSlot, eth.AddressAsLeftPaddedHash(proxyAdminAddr))
		default:
			db.SetCode(addr, depBytecode)
			db.SetState(addr, AdminSlot, eth.AddressAsLeftPaddedHash(proxyAdminAddr))
		}

		log.Trace("Set proxy", "address", addr, "admin", proxyAdminAddr)
	}

	return nil
}

// SetPrecompileBalances will set a single wei at each precompile address.
// This is an optimization to make calling them cheaper.
func SetPrecompileBalances(gen *core.Genesis) {
	for i := 0; i < PrecompileCount; i++ {
		addr := common.BytesToAddress([]byte{byte(i)})
		acc := gen.Alloc[addr]
		if acc.Balance == nil {
			acc.Balance = new(big.Int)
		}
		acc.Balance = acc.Balance.Add(acc.Balance, big.NewInt(1))
		gen.Alloc[addr] = acc
	}
}

func setupPredeploy(db vm.StateDB, deployResults immutables.DeploymentResults, storage state.StorageConfig, name string, proxyAddr common.Address, implAddr common.Address) error {
	// Use the generated bytecode when there are immutables
	// otherwise use the artifact deployed bytecode
	if bytecode, ok := deployResults[name]; ok {
		log.Info("Setting deployed bytecode with immutables", "name", name, "address", implAddr)
		db.SetCode(implAddr, bytecode)
	} else {
		depBytecode, err := bindings.GetDeployedBytecode(name)
		if err != nil {
			return err
		}
		log.Info("Setting deployed bytecode from solc compiler output", "name", name, "address", implAddr)
		db.SetCode(implAddr, depBytecode)
	}

	// Set the storage values
	if storageConfig, ok := storage[name]; ok {
		log.Info("Setting storage", "name", name, "address", proxyAddr)
		if err := state.SetStorage(name, proxyAddr, storageConfig, db); err != nil {
			return err
		}
	}

	return nil
}
