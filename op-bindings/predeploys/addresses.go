package predeploys

import "github.com/ethereum/go-ethereum/common"

// TODO - we should get a single toml yaml or json file source of truth in @eth-optimism/bedrock package
// This needs to be kept in sync with @eth-optimism/contracts-ts/wagmi.config.ts which also specifies this
// To improve robustness and maintainability contracts-bedrock should export all addresses
const (
	L2ToL1MessagePasser                = "0x4200000000000000000000000000000000000016"
	DeployerWhitelist                  = "0x4200000000000000000000000000000000000002"
	WTON                               = "0x4200000000000000000000000000000000000006"
	L2CrossDomainMessenger             = "0x4200000000000000000000000000000000000007"
	L2StandardBridge                   = "0x4200000000000000000000000000000000000010"
	SequencerFeeVault                  = "0x4200000000000000000000000000000000000011"
	OptimismMintableERC20Factory       = "0x4200000000000000000000000000000000000012"
	L1BlockNumber                      = "0x4200000000000000000000000000000000000013"
	GasPriceOracle                     = "0x420000000000000000000000000000000000000F"
	L1Block                            = "0x4200000000000000000000000000000000000015"
	GovernanceToken                    = "0x4200000000000000000000000000000000000042"
	LegacyMessagePasser                = "0x4200000000000000000000000000000000000000"
	L2ERC721Bridge                     = "0x4200000000000000000000000000000000000014"
	OptimismMintableERC721Factory      = "0x4200000000000000000000000000000000000017"
	ProxyAdmin                         = "0x4200000000000000000000000000000000000018"
	BaseFeeVault                       = "0x4200000000000000000000000000000000000019"
	L1FeeVault                         = "0x420000000000000000000000000000000000001a"
	SchemaRegistry                     = "0x4200000000000000000000000000000000000020"
	EAS                                = "0x4200000000000000000000000000000000000021"
	WETH                               = "0x4200000000000000000000000000000000000022"
	Permit2                            = "0x4200000000000000000000000000000000000501"
	QuoterV2                           = "0x4200000000000000000000000000000000000502"
	SwapRouter02                       = "0x4200000000000000000000000000000000000503"
	UniswapV3Factory                   = "0x4200000000000000000000000000000000000504"
	NFTDescriptor                      = "0x4200000000000000000000000000000000000505"
	NonfungiblePositionManager         = "0x4200000000000000000000000000000000000506"
	NonfungibleTokenPositionDescriptor = "0x4200000000000000000000000000000000000507"
	TickLens                           = "0x4200000000000000000000000000000000000508"
	UniswapInterfaceMulticall          = "0x4200000000000000000000000000000000000509"
	L2UsdcBridge                       = "0x4200000000000000000000000000000000000775"
	MasterMinter                       = "0x4200000000000000000000000000000000000777"
	FiatTokenV2_2                      = "0x4200000000000000000000000000000000000778"
)

var (
	L2ToL1MessagePasserAddr                = common.HexToAddress(L2ToL1MessagePasser)
	DeployerWhitelistAddr                  = common.HexToAddress(DeployerWhitelist)
	WTONAddr                               = common.HexToAddress(WTON)
	L2CrossDomainMessengerAddr             = common.HexToAddress(L2CrossDomainMessenger)
	L2StandardBridgeAddr                   = common.HexToAddress(L2StandardBridge)
	SequencerFeeVaultAddr                  = common.HexToAddress(SequencerFeeVault)
	OptimismMintableERC20FactoryAddr       = common.HexToAddress(OptimismMintableERC20Factory)
	L1BlockNumberAddr                      = common.HexToAddress(L1BlockNumber)
	GasPriceOracleAddr                     = common.HexToAddress(GasPriceOracle)
	L1BlockAddr                            = common.HexToAddress(L1Block)
	GovernanceTokenAddr                    = common.HexToAddress(GovernanceToken)
	LegacyMessagePasserAddr                = common.HexToAddress(LegacyMessagePasser)
	L2ERC721BridgeAddr                     = common.HexToAddress(L2ERC721Bridge)
	OptimismMintableERC721FactoryAddr      = common.HexToAddress(OptimismMintableERC721Factory)
	ProxyAdminAddr                         = common.HexToAddress(ProxyAdmin)
	BaseFeeVaultAddr                       = common.HexToAddress(BaseFeeVault)
	L1FeeVaultAddr                         = common.HexToAddress(L1FeeVault)
	SchemaRegistryAddr                     = common.HexToAddress(SchemaRegistry)
	EASAddr                                = common.HexToAddress(EAS)
	WETHAddr                               = common.HexToAddress(WETH)
	Permit2Addr                            = common.HexToAddress(Permit2)
	QuoterV2Addr                           = common.HexToAddress(QuoterV2)
	SwapRouter02Addr                       = common.HexToAddress(SwapRouter02)
	UniswapV3FactoryAddr                   = common.HexToAddress(UniswapV3Factory)
	NFTDescriptorAddr                      = common.HexToAddress(NFTDescriptor)
	NonfungiblePositionManagerAddr         = common.HexToAddress(NonfungiblePositionManager)
	NonfungibleTokenPositionDescriptorAddr = common.HexToAddress(NonfungibleTokenPositionDescriptor)
	TickLensAddr                           = common.HexToAddress(TickLens)
	UniswapInterfaceMulticallAddr          = common.HexToAddress(UniswapInterfaceMulticall)
	L2UsdcBridgeAddr                       = common.HexToAddress(L2UsdcBridge)
	MasterMinterAddr                       = common.HexToAddress(MasterMinter)
	FiatTokenV2_2Addr                      = common.HexToAddress(FiatTokenV2_2)

	Predeploys = make(map[string]*common.Address)
)

// IsProxied returns true for predeploys that will sit behind a proxy contract
func IsProxied(predeployAddr common.Address) bool {
	switch predeployAddr {
	case WTONAddr:
	case GovernanceTokenAddr:
	case MasterMinterAddr:
	default:
		return true
	}
	return false
}

func init() {
	Predeploys["L2ToL1MessagePasser"] = &L2ToL1MessagePasserAddr
	Predeploys["DeployerWhitelist"] = &DeployerWhitelistAddr
	Predeploys["WTON"] = &WTONAddr
	Predeploys["L2CrossDomainMessenger"] = &L2CrossDomainMessengerAddr
	Predeploys["L2StandardBridge"] = &L2StandardBridgeAddr
	Predeploys["SequencerFeeVault"] = &SequencerFeeVaultAddr
	Predeploys["OptimismMintableERC20Factory"] = &OptimismMintableERC20FactoryAddr
	Predeploys["L1BlockNumber"] = &L1BlockNumberAddr
	Predeploys["GasPriceOracle"] = &GasPriceOracleAddr
	Predeploys["L1Block"] = &L1BlockAddr
	Predeploys["GovernanceToken"] = &GovernanceTokenAddr
	Predeploys["LegacyMessagePasser"] = &LegacyMessagePasserAddr
	Predeploys["L2ERC721Bridge"] = &L2ERC721BridgeAddr
	Predeploys["OptimismMintableERC721Factory"] = &OptimismMintableERC721FactoryAddr
	Predeploys["ProxyAdmin"] = &ProxyAdminAddr
	Predeploys["BaseFeeVault"] = &BaseFeeVaultAddr
	Predeploys["L1FeeVault"] = &L1FeeVaultAddr
	Predeploys["SchemaRegistry"] = &SchemaRegistryAddr
	Predeploys["EAS"] = &EASAddr
	Predeploys["WETH"] = &WETHAddr
	Predeploys["Permit2"] = &Permit2Addr
	Predeploys["QuoterV2"] = &QuoterV2Addr
	Predeploys["SwapRouter02"] = &SwapRouter02Addr
	Predeploys["UniswapV3Factory"] = &UniswapV3FactoryAddr
	Predeploys["NFTDescriptor"] = &NFTDescriptorAddr
	Predeploys["NonfungiblePositionManager"] = &NonfungiblePositionManagerAddr
	Predeploys["NonfungibleTokenPositionDescriptor"] = &NonfungibleTokenPositionDescriptorAddr
	Predeploys["TickLens"] = &TickLensAddr
	Predeploys["UniswapInterfaceMulticall"] = &UniswapInterfaceMulticallAddr
	Predeploys["L2UsdcBridge"] = &L2UsdcBridgeAddr
	Predeploys["MasterMinter"] = &MasterMinterAddr
	Predeploys["FiatTokenV2_2"] = &FiatTokenV2_2Addr
}
