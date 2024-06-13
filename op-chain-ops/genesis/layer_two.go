package genesis

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/predeploys"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/deployer"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/immutables"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/squash"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/state"
	"github.com/tokamak-network/tokamak-thanos/op-node/rollup/derive"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// BuildL2Genesis will build the L2 genesis block.
func BuildL2Genesis(config *DeployConfig, l1StartBlock *types.Block) (*core.Genesis, error) {
	genspec, err := NewL2Genesis(config, l1StartBlock)
	if err != nil {
		return nil, err
	}

	db := state.NewMemoryStateDB(genspec)
	if config.FundDevAccounts {
		log.Info("Funding developer accounts in L2 genesis")
		FundDevAccounts(db)
	}

	if config.SetPrecompileBalances {
		log.Info("Setting precompile balances in L2 genesis")
		SetPrecompileBalances(db)
	}

	storage, err := NewL2StorageConfig(config, l1StartBlock)
	if err != nil {
		return nil, err
	}

	immutableConfig, err := NewL2ImmutableConfig(config, l1StartBlock)
	if err != nil {
		return nil, err
	}

	// Set up the LegacyERC20NativeToken
	db.CreateAccount(predeploys.LegacyERC20NativeTokenAddr)

	// Set up the proxies
	err = setProxies(db, predeploys.ProxyAdminAddr, BigL2PredeployNamespace, 2048)
	if err != nil {
		return nil, err
	}

	// Set up the implementations that contain immutables
	deployResults, err := immutables.Deploy(immutableConfig)
	if err != nil {
		return nil, err
	}
	for name, predeploy := range predeploys.Predeploys {
		if predeploy.Enabled != nil && !predeploy.Enabled(config) {
			log.Warn("Skipping disabled predeploy.", "name", name, "address", predeploy.Address)
			continue
		}

		codeAddr := predeploy.Address
		switch name {
		case "Permit2":
			deployerAddressBytes, err := bindings.GetDeployerAddress(name)
			if err != nil {
				return nil, err
			}
			deployerAddress := common.BytesToAddress(deployerAddressBytes)
			predeploys := map[string]*common.Address{
				"DeterministicDeploymentProxy": &deployerAddress,
			}
			backend, err := deployer.NewL2BackendWithChainIDAndPredeploys(
				new(big.Int).SetUint64(config.L2ChainID),
				predeploys,
			)
			if err != nil {
				return nil, err
			}
			deployedBin, err := deployer.DeployWithDeterministicDeployer(backend, name)
			if err != nil {
				return nil, err
			}
			deployResults[name] = deployedBin
			fallthrough
		case "MultiCall3", "Create2Deployer", "Safe_v130",
			"SafeL2_v130", "MultiSendCallOnly_v130", "SafeSingletonFactory",
			"DeterministicDeploymentProxy", "MultiSend_v130", "SenderCreator", "EntryPoint":
			db.CreateAccount(codeAddr)
		case "SignatureChecker":
			bytecode, err := bindings.GetDeployedBytecode(name)
			if err != nil {
				return nil, err
			}
			a := bytecode[:1]
			b := bytecode[21:]
			c := hexutil.Bytes{}
			d, _ := hex.DecodeString("4200000000000000000000000000000000000776")
			c = append(c, a...)
			c = append(c, d...)
			c = append(c, b...)
			bytecode = c
			deployResults[name] = bytecode
		case "FiatTokenV2_2":
			codeAddr, err = AddressToCodeNamespace(predeploy.Address)
			if err != nil {
				return nil, fmt.Errorf("error converting to code namespace: %w", err)
			}
			db.CreateAccount(codeAddr)
			db.SetState(predeploy.Address, ImplementationSlotForZepplin, eth.AddressAsLeftPaddedHash(codeAddr))
			log.Info("Set proxy for FiatTokenV2_2", "name", name, "address", predeploy.Address, "implementation", codeAddr)
		case "UniswapV3Factory":
			dep, err := bindings.GetDeployedBytecode(name)
			if err != nil {
				return nil, err
			}
			originalCode := "5050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000"
			newCode := "5050565b306001600160a01b037f0000000000000000000000004200000000000000000000000000000000000502"
			originalBytes, _ := hex.DecodeString(originalCode)
			newBytes, _ := hex.DecodeString(newCode)
			startIndex := bytes.Index(dep, originalBytes)
			if startIndex != -1 {
				a := dep[:startIndex]
				b := dep[startIndex+len(originalBytes):]
				c := hexutil.Bytes{}
				c = append(c, a...)
				c = append(c, newBytes...)
				c = append(c, b...)
				dep = c
			}
			deployResults[name] = dep
		case "NFTDescriptor":
			bytecode, err := bindings.GetDeployedBytecode(name)
			if err != nil {
				return nil, err
			}
			a := bytecode[:1]
			b := bytecode[21:]
			c := hexutil.Bytes{}
			d, _ := hex.DecodeString("4200000000000000000000000000000000000503")
			c = append(c, a...)
			c = append(c, d...)
			c = append(c, b...)
			bytecode = c
			deployResults[name] = bytecode
			fallthrough
		default:
			if !predeploy.ProxyDisabled {
				codeAddr, err = AddressToCodeNamespace(predeploy.Address)
				if err != nil {
					return nil, fmt.Errorf("error converting to code namespace: %w", err)
				}
				db.CreateAccount(codeAddr)
				db.SetState(predeploy.Address, ImplementationSlot, eth.AddressAsLeftPaddedHash(codeAddr))
				log.Info("Set proxy", "name", name, "address", predeploy.Address, "implementation", codeAddr)
			}
		}

		if predeploy.ProxyDisabled && db.Exist(predeploy.Address) {
			db.DeleteState(predeploy.Address, AdminSlot)
		}

		if err := setupPredeploy(db, deployResults, storage, name, predeploy.Address, codeAddr); err != nil {
			return nil, err
		}
		code := db.GetCode(codeAddr)
		if len(code) == 0 {
			return nil, fmt.Errorf("code not set for %s", name)
		}
	}

	if err := PerformUpgradeTxs(db); err != nil {
		return nil, fmt.Errorf("failed to perform upgrade txs: %w", err)
	}

	return db.Genesis(), nil
}

func PerformUpgradeTxs(db *state.MemoryStateDB) error {
	// Only the Ecotone upgrade is performed with upgrade-txs.
	if !db.Genesis().Config.IsEcotone(db.Genesis().Timestamp) {
		return nil
	}
	sim := squash.NewSimulator(db)
	ecotone, err := derive.EcotoneNetworkUpgradeTransactions()
	if err != nil {
		return fmt.Errorf("failed to build ecotone upgrade txs: %w", err)
	}
	if err := sim.AddUpgradeTxs(ecotone); err != nil {
		return fmt.Errorf("failed to apply ecotone upgrade txs: %w", err)
	}
	return nil
}
