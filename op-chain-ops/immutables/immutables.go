package immutables

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
	"github.com/ethereum-optimism/optimism/op-chain-ops/deployer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ImmutableValues represents the values to be set in immutable code.
// The key is the name of the variable and the value is the value to set in
// immutable code.
type ImmutableValues map[string]any

// ImmutableConfig represents the immutable configuration for the L2 predeploy
// contracts.
type ImmutableConfig map[string]ImmutableValues

// Check does a sanity check that the specific values that
// Optimism uses are set inside of the ImmutableConfig.
func (i ImmutableConfig) Check() error {
	if _, ok := i["L2CrossDomainMessenger"]["otherMessenger"]; !ok {
		return errors.New("L2CrossDomainMessenger otherMessenger not set")
	}
	if _, ok := i["L2StandardBridge"]["otherBridge"]; !ok {
		return errors.New("L2StandardBridge otherBridge not set")
	}
	if _, ok := i["L2ERC721Bridge"]["messenger"]; !ok {
		return errors.New("L2ERC721Bridge messenger not set")
	}
	if _, ok := i["L2ERC721Bridge"]["otherBridge"]; !ok {
		return errors.New("L2ERC721Bridge otherBridge not set")
	}
	if _, ok := i["OptimismMintableERC721Factory"]["bridge"]; !ok {
		return errors.New("OptimismMintableERC20Factory bridge not set")
	}
	if _, ok := i["OptimismMintableERC721Factory"]["remoteChainId"]; !ok {
		return errors.New("OptimismMintableERC20Factory remoteChainId not set")
	}
	if _, ok := i["SequencerFeeVault"]["recipient"]; !ok {
		return errors.New("SequencerFeeVault recipient not set")
	}
	if _, ok := i["L1FeeVault"]["recipient"]; !ok {
		return errors.New("L1FeeVault recipient not set")
	}
	if _, ok := i["BaseFeeVault"]["recipient"]; !ok {
		return errors.New("BaseFeeVault recipient not set")
	}
	return nil
}

// DeploymentResults represents the output of deploying each of the
// contracts so that the immutables can be set properly in the bytecode.
type DeploymentResults map[string]hexutil.Bytes

// BuildOptimism will deploy the L2 predeploys so that their immutables are set
// correctly.
func BuildOptimism(immutable ImmutableConfig) (DeploymentResults, error) {
	if err := immutable.Check(); err != nil {
		return DeploymentResults{}, err
	}
	factoryV2Addr := common.HexToAddress("0x0000000000000000000000000000000000000000")

	deployments := []deployer.Constructor{
		{
			Name: "GasPriceOracle",
		},
		{
			Name: "L1Block",
		},
		{
			Name: "L2CrossDomainMessenger",
			Args: []interface{}{
				immutable["L2CrossDomainMessenger"]["otherMessenger"],
			},
		},
		{
			Name: "L2StandardBridge",
			Args: []interface{}{
				immutable["L2StandardBridge"]["otherBridge"],
			},
		},
		{
			Name: "L2ToL1MessagePasser",
		},
		{
			Name: "SequencerFeeVault",
			Args: []interface{}{
				immutable["SequencerFeeVault"]["recipient"],
				immutable["SequencerFeeVault"]["minimumWithdrawalAmount"],
				immutable["SequencerFeeVault"]["withdrawalNetwork"],
			},
		},
		{
			Name: "BaseFeeVault",
			Args: []interface{}{
				immutable["BaseFeeVault"]["recipient"],
				immutable["BaseFeeVault"]["minimumWithdrawalAmount"],
				immutable["BaseFeeVault"]["withdrawalNetwork"],
			},
		},
		{
			Name: "L1FeeVault",
			Args: []interface{}{
				immutable["L1FeeVault"]["recipient"],
				immutable["L1FeeVault"]["minimumWithdrawalAmount"],
				immutable["L1FeeVault"]["withdrawalNetwork"],
			},
		},
		{
			Name: "OptimismMintableERC20Factory",
		},
		{
			Name: "DeployerWhitelist",
		},
		{
			Name: "LegacyMessagePasser",
		},
		{
			Name: "L1BlockNumber",
		},
		{
			Name: "L2ERC721Bridge",
			Args: []interface{}{
				immutable["L2ERC721Bridge"]["otherBridge"],
			},
		},
		{
			Name: "OptimismMintableERC721Factory",
			Args: []interface{}{
				predeploys.L2ERC721BridgeAddr,
				immutable["OptimismMintableERC721Factory"]["remoteChainId"],
			},
		},
		{
			Name: "LegacyERC20ETH",
		},
		{
			Name: "EAS",
		},
		{
			Name: "SchemaRegistry",
		},
		{
			Name: "WETH",
		},
		{
			Name: "Permit2",
		},
		{
			Name: "QuoterV2",
			Args: []interface{}{
				predeploys.UniswapV3FactoryAddr,
				predeploys.WETHAddr,
			},
		},
		{
			Name: "SwapRouter02",
			Args: []interface{}{
				factoryV2Addr,
				predeploys.UniswapV3FactoryAddr,
				predeploys.NonfungiblePositionManagerAddr,
				predeploys.WETHAddr,
			},
		},
		{
			Name: "UniswapV3Factory",
			Args: []interface{}{
				predeploys.UniswapV3FactoryAddr,
				immutable["UniswapV3Factory"]["feeAmountTickSpacing500"],
				immutable["UniswapV3Factory"]["feeAmountTickSpacing3000"],
				immutable["UniswapV3Factory"]["feeAmountTickSpacing10000"],
			},
		},
		{
			Name: "NonfungiblePositionManager",
			Args: []interface{}{
				predeploys.UniswapV3FactoryAddr,
				predeploys.WETHAddr,
				predeploys.NonfungibleTokenPositionDescriptorAddr,
			},
		},
		{
			Name: "NonfungibleTokenPositionDescriptor",
			Args: []interface{}{
				predeploys.NonfungibleTokenPositionDescriptorAddr,
			},
		},
		{
			Name: "TickLens",
		},
		{
			Name: "UniswapInterfaceMulticall",
		},
		{
			Name: "L2UsdcBridge",
		},
		{
			Name: "SignatureChecker",
		},
		{
			Name: "MasterMinter",
			Args: []interface{}{
				immutable["MasterMinter"]["_minterManager"],
			},
		},
		{
			Name: "FiatTokenV2_2",
		},
	}
	return BuildL2(deployments)
}

// BuildL2 will deploy contracts to a simulated backend so that their immutables
// can be properly set. The bytecode returned in the results is suitable to be
// inserted into the state via state surgery.
func BuildL2(constructors []deployer.Constructor) (DeploymentResults, error) {
	log.Info("Creating L2 state")
	deployments, err := deployer.Deploy(deployer.NewL2Backend(), constructors, l2Deployer)
	if err != nil {
		return nil, err
	}
	results := make(DeploymentResults)
	for _, dep := range deployments {
		results[dep.Name] = dep.Bytecode
	}
	return results, nil
}

func l2Deployer(backend *backends.SimulatedBackend, opts *bind.TransactOpts, deployment deployer.Constructor) (*types.Transaction, error) {
	var tx *types.Transaction
	var recipient common.Address
	var minimumWithdrawalAmount *big.Int
	var withdrawalNetwork uint8
	var _factoryV2 common.Address
	var factoryV3 common.Address
	var _positionManager common.Address
	var _WETH9 common.Address
	var _factory common.Address
	var _tokenDescriptor_ common.Address
	var err error
	switch deployment.Name {
	case "GasPriceOracle":
		_, tx, _, err = bindings.DeployGasPriceOracle(opts, backend)
	case "L1Block":
		// No arguments required for the L1Block contract
		_, tx, _, err = bindings.DeployL1Block(opts, backend)
	case "L2CrossDomainMessenger":
		otherMessenger, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for otherMessenger")
		}
		_, tx, _, err = bindings.DeployL2CrossDomainMessenger(opts, backend, otherMessenger)
	case "L2StandardBridge":
		otherBridge, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for otherBridge")
		}
		_, tx, _, err = bindings.DeployL2StandardBridge(opts, backend, otherBridge)
	case "L2ToL1MessagePasser":
		// No arguments required for L2ToL1MessagePasser
		_, tx, _, err = bindings.DeployL2ToL1MessagePasser(opts, backend)
	case "SequencerFeeVault":
		recipient, minimumWithdrawalAmount, withdrawalNetwork, err = prepareFeeVaultArguments(deployment)
		if err != nil {
			return nil, err
		}
		_, tx, _, err = bindings.DeploySequencerFeeVault(opts, backend, recipient, minimumWithdrawalAmount, withdrawalNetwork)
	case "BaseFeeVault":
		recipient, minimumWithdrawalAmount, withdrawalNetwork, err = prepareFeeVaultArguments(deployment)
		if err != nil {
			return nil, err
		}
		_, tx, _, err = bindings.DeployBaseFeeVault(opts, backend, recipient, minimumWithdrawalAmount, withdrawalNetwork)
	case "L1FeeVault":
		recipient, minimumWithdrawalAmount, withdrawalNetwork, err = prepareFeeVaultArguments(deployment)
		if err != nil {
			return nil, err
		}
		_, tx, _, err = bindings.DeployL1FeeVault(opts, backend, recipient, minimumWithdrawalAmount, withdrawalNetwork)
	case "OptimismMintableERC20Factory":
		_, tx, _, err = bindings.DeployOptimismMintableERC20Factory(opts, backend)
	case "DeployerWhitelist":
		_, tx, _, err = bindings.DeployDeployerWhitelist(opts, backend)
	case "LegacyMessagePasser":
		_, tx, _, err = bindings.DeployLegacyMessagePasser(opts, backend)
	case "L1BlockNumber":
		_, tx, _, err = bindings.DeployL1BlockNumber(opts, backend)
	case "L2ERC721Bridge":
		otherBridge, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for otherBridge")
		}
		_, tx, _, err = bindings.DeployL2ERC721Bridge(opts, backend, otherBridge)
	case "OptimismMintableERC721Factory":
		bridge, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for bridge")
		}
		remoteChainId, ok := deployment.Args[1].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid type for remoteChainId")
		}
		_, tx, _, err = bindings.DeployOptimismMintableERC721Factory(opts, backend, bridge, remoteChainId)
	case "LegacyERC20ETH":
		_, tx, _, err = bindings.DeployLegacyERC20ETH(opts, backend)
	case "EAS":
		_, tx, _, err = bindings.DeployEAS(opts, backend)
	case "SchemaRegistry":
		_, tx, _, err = bindings.DeploySchemaRegistry(opts, backend)
	case "WETH":
		_, tx, _, err = bindings.DeployWETH(opts, backend)
	case "Permit2":
		_, tx, _, err = bindings.DeployPermit2(opts, backend)

	case "QuoterV2":
		_factory, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for _factory")
		}
		_WETH9, ok := deployment.Args[1].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for _WETH9")
		}
		_, tx, _, err = bindings.DeployQuoterV2(opts, backend, _factory, _WETH9)
	case "SwapRouter02":
		_factoryV2, factoryV3, _positionManager, _WETH9, err = prepareSwapRouter02(deployment)
		if err != nil {
			return nil, err
		}
		_, tx, _, err = bindings.DeploySwapRouter02(opts, backend, _factoryV2, factoryV3, _positionManager, _WETH9)
	case "UniswapV3Factory":
		_, tx, _, err = bindings.DeployUniswapV3Factory(opts, backend)
	case "NFTDescriptor":
		_, tx, _, err = bindings.DeployNFTDescriptor(opts, backend)
	case "NonfungiblePositionManager":
		_factory, _WETH9, _tokenDescriptor_, err = prepareonfungiblePositionManager(deployment)
		if err != nil {
			return nil, err
		}
		_, tx, _, err = bindings.DeployNonfungiblePositionManager(opts, backend, _factory, _WETH9, _tokenDescriptor_)
	case "NonfungibleTokenPositionDescriptor":
		_WETH9, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for _WETH9")
		}
		_nativeCurrencyLabelBytes, ok := deployment.Args[1].([32]byte)
		if !ok {
			return nil, fmt.Errorf("invalid type for _nativeCurrencyLabelBytes")
		}
		_, tx, _, err = bindings.DeployNonfungibleTokenPositionDescriptor(opts, backend, _WETH9, _nativeCurrencyLabelBytes)
	case "TickLens":
		_, tx, _, err = bindings.DeployTickLens(opts, backend)
	case "UniswapInterfaceMulticall":
		_, tx, _, err = bindings.DeployUniswapInterfaceMulticall(opts, backend)
	case "L2UsdcBridge":
		_, tx, _, err = bindings.DeployL2UsdcBridge(opts, backend)
	case "SignatureChecker":
		_, tx, _, err = bindings.DeploySignatureChecker(opts, backend)
	case "MasterMinter":
		_minterManager, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for _minterManager")
		}
		_, tx, _, err = bindings.DeployMasterMinter(opts, backend, _minterManager)
	case "FiatTokenV2_2":
		_, tx, _, err = bindings.DeployFiatTokenV22(opts, backend)
	default:
		return tx, fmt.Errorf("unknown contract: %s", deployment.Name)
	}

	return tx, err
}

func prepareFeeVaultArguments(deployment deployer.Constructor) (common.Address, *big.Int, uint8, error) {
	recipient, ok := deployment.Args[0].(common.Address)
	if !ok {
		return common.Address{}, nil, 0, fmt.Errorf("invalid type for recipient")
	}
	minimumWithdrawalAmountHex, ok := deployment.Args[1].(*hexutil.Big)
	if !ok {
		return common.Address{}, nil, 0, fmt.Errorf("invalid type for minimumWithdrawalAmount")
	}
	withdrawalNetwork, ok := deployment.Args[2].(uint8)
	if !ok {
		return common.Address{}, nil, 0, fmt.Errorf("invalid type for withdrawalNetwork")
	}
	return recipient, minimumWithdrawalAmountHex.ToInt(), withdrawalNetwork, nil
}

func prepareSwapRouter02(deployment deployer.Constructor) (common.Address, common.Address, common.Address, common.Address, error) {
	_factoryV2, ok := deployment.Args[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("invalid type for _factoryV2")
	}
	factoryV3, ok := deployment.Args[1].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("invalid type for factoryV3")
	}
	_positionManager, ok := deployment.Args[2].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("invalid type for _positionManager")
	}
	_WETH9, ok := deployment.Args[3].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("invalid type for _WETH9")
	}
	return _factoryV2, factoryV3, _positionManager, _WETH9, nil
}

func prepareonfungiblePositionManager(deployment deployer.Constructor) (common.Address, common.Address, common.Address, error) {
	_factory, ok := deployment.Args[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("invalid type for _factory")
	}
	_WETH9, ok := deployment.Args[1].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("invalid type for _WETH9")
	}
	_tokenDescriptor_, ok := deployment.Args[2].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("invalid type for _tokenDescriptor_")
	}

	return _factory, _WETH9, _tokenDescriptor_, nil
}
