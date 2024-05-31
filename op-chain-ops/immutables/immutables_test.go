package immutables_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/predeploys"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/immutables"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBuildOptimism(t *testing.T) {
	cfg := immutables.PredeploysImmutableConfig{
		LegacyERC20NativeToken: struct {
			RemoteToken common.Address
			Bridge      common.Address
			Decimals    uint8
		}{
			RemoteToken: common.Address{},
			Bridge:      predeploys.L2StandardBridgeAddr,
			Decimals:    18,
		},
		L2ToL1MessagePasser:    struct{}{},
		DeployerWhitelist:      struct{}{},
		WNativeToken:           struct{}{},
		L2CrossDomainMessenger: struct{}{},
		L2StandardBridge:       struct{}{},
		SequencerFeeVault: struct {
			Recipient           common.Address
			MinWithdrawalAmount *big.Int
			WithdrawalNetwork   uint8
		}{
			Recipient:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			MinWithdrawalAmount: big.NewInt(100),
			WithdrawalNetwork:   0,
		},
		L1BlockNumber:       struct{}{},
		GasPriceOracle:      struct{}{},
		L1Block:             struct{}{},
		GovernanceToken:     struct{}{},
		LegacyMessagePasser: struct{}{},
		L2ERC721Bridge:      struct{}{},
		OptimismMintableERC721Factory: struct {
			Bridge        common.Address
			RemoteChainId *big.Int
		}{
			Bridge:        predeploys.L2StandardBridgeAddr,
			RemoteChainId: big.NewInt(1),
		},
		OptimismMintableERC20Factory: struct{}{},
		ProxyAdmin:                   struct{}{},
		BaseFeeVault: struct {
			Recipient           common.Address
			MinWithdrawalAmount *big.Int
			WithdrawalNetwork   uint8
		}{
			Recipient:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			MinWithdrawalAmount: big.NewInt(200),
			WithdrawalNetwork:   0,
		},
		L1FeeVault: struct {
			Recipient           common.Address
			MinWithdrawalAmount *big.Int
			WithdrawalNetwork   uint8
		}{
			Recipient:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
			MinWithdrawalAmount: big.NewInt(200),
			WithdrawalNetwork:   1,
		},
		SchemaRegistry: struct{}{},
		EAS: struct{ Name string }{
			Name: "EAS",
		},
		ETH: struct {
			RemoteToken common.Address
			Bridge      common.Address
			Decimals    uint8
		}{
			RemoteToken: common.Address{},
			Bridge:      predeploys.L2StandardBridgeAddr,
			Decimals:    18,
		},
		L2UsdcBridge:     struct{}{},
		SignatureChecker: struct{}{},
		MasterMinter: struct {
			MinterManager common.Address
		}{
			MinterManager: predeploys.FiatTokenV2_2Addr,
		},
		FiatTokenV2_2: struct{}{},
		QuoterV2: struct {
			Factory common.Address
			WETH9   common.Address
		}{
			Factory: predeploys.UniswapV3FactoryAddr,
			WETH9:   predeploys.WNativeTokenAddr,
		},
		SwapRouter02: struct {
			FactoryV2       common.Address
			FactoryV3       common.Address
			PositionManager common.Address
			WETH9           common.Address
		}{
			FactoryV2:       common.HexToAddress("0x0000000000000000000000000000000000000000"),
			FactoryV3:       predeploys.UniswapV3FactoryAddr,
			PositionManager: predeploys.NonfungiblePositionManagerAddr,
			WETH9:           predeploys.WNativeTokenAddr,
		},
		UniswapV3Factory: struct{}{},
		NFTDescriptor:    struct{}{},
		NonfungiblePositionManager: struct {
			Factory          common.Address
			WETH9            common.Address
			TokenDescriptor_ common.Address
		}{
			Factory:          predeploys.UniswapV3FactoryAddr,
			WETH9:            predeploys.WNativeTokenAddr,
			TokenDescriptor_: predeploys.NonfungibleTokenPositionDescriptorAddr,
		},
		NonfungibleTokenPositionDescriptor: struct {
			WETH9                    common.Address
			NativeCurrencyLabelBytes [32]byte
		}{
			WETH9:                    predeploys.WNativeTokenAddr,
			NativeCurrencyLabelBytes: [32]byte{84, 87, 79, 78, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		TickLens:                  struct{}{},
		UniswapInterfaceMulticall: struct{}{},
		UniversalRouter: struct {
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
		}{
			Permit2:                     predeploys.Permit2Addr,
			Weth9:                       predeploys.WNativeTokenAddr,
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
			V3Factory:                   predeploys.UniswapV3FactoryAddr,
			PairInitCodeHash:            common.HexToHash("0x96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f"),
			PoolInitCodeHash:            common.HexToHash("0xe34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b54"),
		},

		Create2Deployer:              struct{}{},
		MultiCall3:                   struct{}{},
		Safe_v130:                    struct{}{},
		SafeL2_v130:                  struct{}{},
		MultiSendCallOnly_v130:       struct{}{},
		SafeSingletonFactory:         struct{}{},
		DeterministicDeploymentProxy: struct{}{},
		MultiSend_v130:               struct{}{},
		Permit2:                      struct{}{},
		SenderCreator:                struct{}{},
		EntryPoint:                   struct{}{},
	}

	require.NoError(t, cfg.Check())
	results, err := immutables.Deploy(&cfg)
	require.NoError(t, err)
	require.NotNil(t, results)

	// Build a mapping of all of the predeploys
	all := map[string]bool{}
	// Build a mapping of the predeploys with immutable config
	withConfig := map[string]bool{}

	require.NoError(t, cfg.ForEach(func(name string, predeployConfig any) error {
		all[name] = true

		// If a predeploy has no config, it needs to have no immutable references in the solc output.
		if reflect.ValueOf(predeployConfig).IsZero() {
			ref, _ := bindings.HasImmutableReferences(name)
			require.Zero(t, ref, "found immutable reference for %s", name)
			return nil
		}
		withConfig[name] = true
		return nil
	}))

	// Ensure that the PredeploysImmutableConfig is kept up to date
	for name := range predeploys.Predeploys {
		require.Truef(t, all[name], "predeploy %s not in set of predeploys", name)

		ref, err := bindings.HasImmutableReferences(name)
		// If there is predeploy config, there should be an immutable reference
		if withConfig[name] {
			require.NoErrorf(t, err, "error getting immutable reference for %s", name)
			require.NotZerof(t, ref, "no immutable reference for %s", name)
		} else {
			require.Zero(t, ref, "found immutable reference for %s", name)
		}
	}

	// Only the exact contracts that we care about are being modified
	require.Equal(t, len(results), len(withConfig))

	for name, bytecode := range results {
		// There is bytecode there
		require.Greater(t, len(bytecode), 0)
		// It is in the set of contracts that we care about
		require.Truef(t, withConfig[name], "contract %s not in set of contracts", name)
		// The immutable reference is present
		ref, err := bindings.HasImmutableReferences(name)
		require.NoErrorf(t, err, "cannot get immutable reference for %s", name)
		require.NotZerof(t, ref, "contract %s has no immutable reference", name)
	}
}
