package immutables

import (
	"fmt"
	"log"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/deployer"
)

// PredeploysImmutableConfig represents the set of L2 predeploys. It includes all
// L2 predeploys - not just ones with immutable values. This is to be very explicit
// about the configuration of the predeploys. It is important that the inner struct
// fields are in the same order as the constructor arguments in the solidity code.
type PredeploysImmutableConfig struct {
	LegacyERC20NativeToken struct {
		RemoteToken common.Address
		Bridge      common.Address
		Decimals    uint8
	}
	L2ToL1MessagePasser    struct{}
	DeployerWhitelist      struct{}
	WNativeToken           struct{}
	L2CrossDomainMessenger struct{}
	L2StandardBridge       struct{}
	SequencerFeeVault      struct {
		Recipient           common.Address
		MinWithdrawalAmount *big.Int
		WithdrawalNetwork   uint8
	}
	OptimismMintableERC20Factory  struct{}
	L1BlockNumber                 struct{}
	GasPriceOracle                struct{}
	L1Block                       struct{}
	GovernanceToken               struct{}
	LegacyMessagePasser           struct{}
	L2ERC721Bridge                struct{}
	OptimismMintableERC721Factory struct {
		Bridge        common.Address
		RemoteChainId *big.Int
	}
	ProxyAdmin   struct{}
	BaseFeeVault struct {
		Recipient           common.Address
		MinWithdrawalAmount *big.Int
		WithdrawalNetwork   uint8
	}
	L1FeeVault struct {
		Recipient           common.Address
		MinWithdrawalAmount *big.Int
		WithdrawalNetwork   uint8
	}
	SchemaRegistry struct{}
	EAS            struct {
		Name string
	}
	ETH struct {
		RemoteToken common.Address
		Bridge      common.Address
		Decimals    uint8
	}
	L2UsdcBridge     struct{}
	SignatureChecker struct{}
	MasterMinter     struct {
		MinterManager common.Address
	}
	FiatTokenV2_2 struct{}
	QuoterV2      struct {
		Factory common.Address
		WETH9   common.Address
	}
	SwapRouter02 struct {
		FactoryV2       common.Address
		FactoryV3       common.Address
		PositionManager common.Address
		WETH9           common.Address
	}
	UniswapV3Factory           struct{}
	NFTDescriptor              struct{}
	NonfungiblePositionManager struct {
		Factory          common.Address
		WETH9            common.Address
		TokenDescriptor_ common.Address
	}
	NonfungibleTokenPositionDescriptor struct {
		WETH9                    common.Address
		NativeCurrencyLabelBytes [32]byte
	}
	TickLens                  struct{}
	UniswapInterfaceMulticall struct{}
	UniversalRouter           struct {
		Permit2                     common.Address
		Weth9                       common.Address
		SeaportV15                  common.Address
		SeaportV14                  common.Address
		OpenseaConduit              common.Address
		NftxZap                     common.Address
		X2y2                        common.Address
		Foundation                  common.Address
		Sudoswap                    common.Address
		ElementMarket               common.Address
		Nft20Zap                    common.Address
		Cryptopunks                 common.Address
		LooksRareV2                 common.Address
		RouterRewardsDistributor    common.Address
		LooksRareRewardsDistributor common.Address
		LooksRareToken              common.Address
		V2Factory                   common.Address
		V3Factory                   common.Address
		PairInitCodeHash            [32]byte
		PoolInitCodeHash            [32]byte
	}
	Create2Deployer              struct{}
	MultiCall3                   struct{}
	Safe_v130                    struct{}
	SafeL2_v130                  struct{}
	MultiSendCallOnly_v130       struct{}
	SafeSingletonFactory         struct{}
	DeterministicDeploymentProxy struct{}
	MultiSend_v130               struct{}
	Permit2                      struct{}
	SenderCreator                struct{}
	EntryPoint                   struct{}
}

type RouterParameters struct {
	Permit2                     common.Address
	Weth9                       common.Address
	SeaportV15                  common.Address
	SeaportV14                  common.Address
	OpenseaConduit              common.Address
	NftxZap                     common.Address
	X2y2                        common.Address
	Foundation                  common.Address
	Sudoswap                    common.Address
	ElementMarket               common.Address
	Nft20Zap                    common.Address
	Cryptopunks                 common.Address
	LooksRareV2                 common.Address
	RouterRewardsDistributor    common.Address
	LooksRareRewardsDistributor common.Address
	LooksRareToken              common.Address
	V2Factory                   common.Address
	V3Factory                   common.Address
	PairInitCodeHash            [32]byte
	PoolInitCodeHash            [32]byte
}

// Check will ensure that the required fields are set on the config.
// An error returned by `GetImmutableReferences` means that the solc compiler
// output for the contract has no immutables in it.
func (c *PredeploysImmutableConfig) Check() error {
	return c.ForEach(func(name string, values any) error {
		val := reflect.ValueOf(values)
		if val.NumField() == 0 {
			return nil
		}

		has, err := bindings.HasImmutableReferences(name)
		exists := err == nil && has
		isZero := val.IsZero()

		// There are immutables defined in the solc output and
		// the config is not empty.
		if exists && !isZero {
			return nil
		}
		// There are no immutables defined in the solc output and
		// the config is empty
		if !exists && isZero {
			return nil
		}

		return fmt.Errorf("invalid immutables config: field %s: %w", name, err)
	})
}

// ForEach will iterate over each of the fields in the config and call the callback
// with the value of the field as well as the field's name.
func (c *PredeploysImmutableConfig) ForEach(cb func(string, any) error) error {
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		internalVal := reflect.ValueOf(field.Interface())
		if err := cb(typ.Field(i).Name, internalVal.Interface()); err != nil {
			return err
		}
	}
	return nil
}

// DeploymentResults represents the output of deploying each of the
// contracts so that the immutables can be set properly in the bytecode.
type DeploymentResults map[string]hexutil.Bytes

// Deploy will deploy L2 predeploys that include immutables. This is to prevent the need
// for parsing the solc output to find the correct immutable offsets and splicing in the values.
// Skip any predeploys that do not have immutables as their bytecode will be directly inserted
// into the state. This does not currently support recursive structs.
func Deploy(config *PredeploysImmutableConfig) (DeploymentResults, error) {
	if err := config.Check(); err != nil {
		return DeploymentResults{}, err
	}
	deployments := make([]deployer.Constructor, 0)

	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if reflect.ValueOf(field.Interface()).IsZero() {
			continue
		}

		deployment := deployer.Constructor{
			Name: typ.Field(i).Name,
			Args: []any{},
		}

		internalVal := reflect.ValueOf(field.Interface())
		if deployment.Name == "UniversalRouter" {
			// UniversalRouter에 대한 올바른 RouterParameters 인스턴스 생성
			routerParams := RouterParameters{
				Permit2:                     common.HexToAddress("0x000000000022D473030F116dDEE9F6B43aC78BA3"),
				Weth9:                       common.HexToAddress("0x4200000000000000000000000000000000000006"),
				SeaportV15:                  common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				SeaportV14:                  common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				OpenseaConduit:              common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				NftxZap:                     common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				X2y2:                        common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				Foundation:                  common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				Sudoswap:                    common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				ElementMarket:               common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				Nft20Zap:                    common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				Cryptopunks:                 common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				LooksRareV2:                 common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				RouterRewardsDistributor:    common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				LooksRareRewardsDistributor: common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				LooksRareToken:              common.HexToAddress("0x819B9E61F02Bdb8841e90Af300d5064AD1a30D84"),
				V2Factory:                   common.HexToAddress("0x0000000000000000000000000000000000000000"),
				V3Factory:                   common.HexToAddress("0x4200000000000000000000000000000000000502"),
				PairInitCodeHash:            [32]byte{0x96, 0xe8, 0xac, 0x42, 0x77, 0x19, 0x8f, 0xf8, 0xb6, 0xf7, 0x85, 0x47, 0x8a, 0xa9, 0xa3, 0x9f, 0x40, 0x3c, 0xb7, 0x68, 0xdd, 0x02, 0xcb, 0xee, 0x32, 0x6c, 0x3e, 0x7d, 0xa3, 0x48, 0x84, 0x5f},
				PoolInitCodeHash:            [32]byte{0xe3, 0x4f, 0x19, 0x9b, 0x19, 0xb2, 0xb4, 0xf4, 0x7f, 0x68, 0x44, 0x26, 0x19, 0xd5, 0x55, 0x52, 0x7d, 0x24, 0x4f, 0x78, 0xa3, 0x29, 0x7e, 0xa8, 0x93, 0x25, 0xf8, 0x43, 0xf8, 0x7b, 0x8b, 0x54},
			}
			deployment.Args = append(deployment.Args, routerParams)
		} else {
			for j := 0; j < internalVal.NumField(); j++ {
				internalField := internalVal.Field(j)
				deployment.Args = append(deployment.Args, internalField.Interface())
			}
		}

		deployments = append(deployments, deployment)
	}

	results, err := deployContractsWithImmutables(deployments)
	if err != nil {
		return nil, fmt.Errorf("cannot deploy contracts with immutables: %w", err)
	}
	return results, nil
}

// deployContractsWithImmutables will deploy contracts to a simulated backend so that their immutables
// can be properly set. The bytecode returned in the results is suitable to be
// inserted into the state via state surgery.
func deployContractsWithImmutables(constructors []deployer.Constructor) (DeploymentResults, error) {
	// L2 백엔드 인스턴스를 생성합니다.
	log.Printf("Starting to create L2 backend instance.")
	backend, err := deployer.NewL2Backend()
	if err != nil {
		log.Printf("백엔드 생성 실패: %v", err)
		return nil, err
	}

	// 컨트랙트를 배포합니다.
	log.Printf("Starting to deploy contracts.")
	deployments, err := deployer.Deploy(backend, constructors, l2ImmutableDeployer)
	if err != nil {
		log.Printf("컨트랙트 배포 실패: %v", err)
		return nil, err
	}

	// 결과 맵을 생성합니다.
	results := make(DeploymentResults)
	for _, dep := range deployments {
		results[dep.Name] = dep.Bytecode
		log.Printf("배포된 컨트랙트: %s, 바이트코드: %s", dep.Name, dep.Bytecode)
	}

	log.Printf("Deployments completed. Results: %v", results)
	return results, nil
}

// l2ImmutableDeployer will deploy L2 predeploys that contain immutables to the simulated backend.
// It only needs to care about the predeploys that have immutables so that the deployed bytecode
// has the dynamic value set at the correct location in the bytecode.
func l2ImmutableDeployer(backend *backends.SimulatedBackend, opts *bind.TransactOpts, deployment deployer.Constructor) (*types.Transaction, error) {
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

	if has, err := bindings.HasImmutableReferences(deployment.Name); err != nil || !has {
		return nil, fmt.Errorf("%s does not have immutables: %w", deployment.Name, err)
	}

	switch deployment.Name {
	case "LegacyERC20NativeToken":
		_, tx, _, err = bindings.DeployLegacyERC20NativeToken(opts, backend)
	case "SequencerFeeVault":
		recipient, minimumWithdrawalAmount, withdrawalNetwork, err = prepareFeeVaultArguments(deployment)
		if err != nil {
			return nil, err
		}
		_, tx, _, err = bindings.DeploySequencerFeeVault(opts, backend, recipient, minimumWithdrawalAmount, withdrawalNetwork)
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
	case "EAS":
		_, tx, _, err = bindings.DeployEAS(opts, backend)
	case "ETH":
		_, tx, _, err = bindings.DeployETH(opts, backend)
	case "MasterMinter":
		_minterManager, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for _minterManager")
		}
		_, tx, _, err = bindings.DeployMasterMinter(opts, backend, _minterManager)
	case "FiatTokenV2_2":
		_, tx, _, err = bindings.DeployFiatTokenV22(opts, backend)
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
		_factory, _WETH9, _tokenDescriptor_, err = PrepareNonfungiblePositionManager(deployment)
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
		// 각 컨트랙트 배포 전 파라미터의 타입을 로깅
	case "UniversalRouter":
		localParams, ok := deployment.Args[0].(RouterParameters)
		if !ok {
			log.Printf("UniversalRouter의 RouterParameters 타입 불일치: 받은 타입: %T, 값: %#v", deployment.Args[0], deployment.Args[0])
			return nil, fmt.Errorf("invalid type for RouterParameters: received type %T", deployment.Args[0])
		}
		convertedParams := ConvertRouterParameters(localParams)
		_, tx, _, err = bindings.DeployUniversalRouter(opts, backend, convertedParams)
		if err != nil {
			log.Printf("UniversalRouter 배포 실패: %v", err)
		}

	// 기타 컨트랙트들도 필요에 따라 로깅 추가 가능
	default:
		log.Printf("알 수 없는 컨트랙트: %s", deployment.Name)
		return nil, fmt.Errorf("unknown contract: %s", deployment.Name)
	}

	if err != nil {
		log.Printf("%s 배포 중 오류 발생: %v", deployment.Name, err)
	}

	return tx, err
}

// prepareFeeVaultArguments is a helper function that parses the arguments for the fee vault contracts.
func prepareFeeVaultArguments(deployment deployer.Constructor) (common.Address, *big.Int, uint8, error) {
	recipient, ok := deployment.Args[0].(common.Address)
	if !ok {
		return common.Address{}, nil, 0, fmt.Errorf("invalid type for recipient")
	}
	minimumWithdrawalAmountHex, ok := deployment.Args[1].(*big.Int)
	if !ok {
		return common.Address{}, nil, 0, fmt.Errorf("invalid type for minimumWithdrawalAmount")
	}
	withdrawalNetwork, ok := deployment.Args[2].(uint8)
	if !ok {
		return common.Address{}, nil, 0, fmt.Errorf("invalid type for withdrawalNetwork")
	}
	return recipient, minimumWithdrawalAmountHex, withdrawalNetwork, nil
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

func PrepareNonfungiblePositionManager(deployment deployer.Constructor) (common.Address, common.Address, common.Address, error) {
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

// ConvertRouterParameters converts local RouterParameters to bindings RouterParameters.
func ConvertRouterParameters(localParams RouterParameters) bindings.RouterParameters {
	return bindings.RouterParameters{
		Permit2:                     localParams.Permit2,
		Weth9:                       localParams.Weth9,
		SeaportV15:                  localParams.SeaportV15,
		SeaportV14:                  localParams.SeaportV14,
		OpenseaConduit:              localParams.OpenseaConduit,
		NftxZap:                     localParams.NftxZap,
		X2y2:                        localParams.X2y2,
		Foundation:                  localParams.Foundation,
		Sudoswap:                    localParams.Sudoswap,
		ElementMarket:               localParams.ElementMarket,
		Nft20Zap:                    localParams.Nft20Zap,
		Cryptopunks:                 localParams.Cryptopunks,
		LooksRareV2:                 localParams.LooksRareV2,
		RouterRewardsDistributor:    localParams.RouterRewardsDistributor,
		LooksRareRewardsDistributor: localParams.LooksRareRewardsDistributor,
		LooksRareToken:              localParams.LooksRareToken,
		V2Factory:                   localParams.V2Factory,
		V3Factory:                   localParams.V3Factory,
		PairInitCodeHash:            localParams.PairInitCodeHash,
		PoolInitCodeHash:            localParams.PoolInitCodeHash,
	}
}
