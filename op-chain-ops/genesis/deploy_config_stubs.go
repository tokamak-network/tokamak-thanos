package genesis

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

// Stub types for FP upstream compatibility.
// These deploy config sub-types are used by op-deployer but don't exist in old tokamak-thanos.
// They provide the nested struct interface expected by op-deployer while the actual DeployConfig
// keeps its flat field layout.

type AltDADeployConfig struct {
	UseAltDA      bool           `json:"useAltDA,omitempty"`
	DACommitmentType string      `json:"daCommitmentType,omitempty"`
	DAChallengeWindow uint64     `json:"daChallengeWindow,omitempty"`
	DAResolveWindow   uint64     `json:"daResolveWindow,omitempty"`
	DABondSize        uint64     `json:"daBondSize,omitempty"`
	DAResolverRefundPercentage uint64 `json:"daResolverRefundPercentage,omitempty"`
}

func (c AltDADeployConfig) Check(_ *log.Logger) error {
	if !c.UseAltDA {
		return nil
	}
	if c.DAChallengeWindow == 0 {
		return fmt.Errorf("DAChallengeWindow must be set when UseAltDA is true")
	}
	if c.DAResolveWindow == 0 {
		return fmt.Errorf("DAResolveWindow must be set when UseAltDA is true")
	}
	return nil
}

type DevDeployConfig struct {
	FundDevAccounts bool `json:"fundDevAccounts,omitempty"`
}

type EIP1559DeployConfig struct {
	EIP1559Denominator       uint64 `json:"eip1559Denominator"`
	EIP1559DenominatorCanyon uint64 `json:"eip1559DenominatorCanyon"`
	EIP1559Elasticity        uint64 `json:"eip1559Elasticity"`
}

type FaultProofDeployConfig struct {
	UseFaultProofs                  bool   `json:"useFaultProofs,omitempty"`
	FaultGameWithdrawalDelay        uint64 `json:"faultGameWithdrawalDelay,omitempty"`
	PreimageOracleMinProposalSize   uint64 `json:"preimageOracleMinProposalSize,omitempty"`
	PreimageOracleChallengePeriod   uint64 `json:"preimageOracleChallengePeriod,omitempty"`
	ProofMaturityDelaySeconds       uint64 `json:"proofMaturityDelaySeconds,omitempty"`
	DisputeGameFinalityDelaySeconds uint64 `json:"disputeGameFinalityDelaySeconds,omitempty"`
}

type FeeMarketConfig struct {
	MinBaseFee           uint64 `json:"minBaseFee,omitempty"`
	DAFootprintGasScalar uint16 `json:"daFootprintGasScalar,omitempty"`
}

type GasPriceOracleDeployConfig struct {
	GasPriceOracleOverhead          uint64 `json:"gasPriceOracleOverhead"`
	GasPriceOracleScalar            uint64 `json:"gasPriceOracleScalar"`
	GasPriceOracleBaseFeeScalar     uint32 `json:"gasPriceOracleBaseFeeScalar"`
	GasPriceOracleBlobBaseFeeScalar uint32 `json:"gasPriceOracleBlobBaseFeeScalar"`
	GasPriceOracleOperatorFeeScalar   uint32 `json:"gasPriceOracleOperatorFeeScalar,omitempty"`
	GasPriceOracleOperatorFeeConstant uint64 `json:"gasPriceOracleOperatorFeeConstant,omitempty"`
}

type GasTokenDeployConfig struct {
	UseCustomGasToken          bool           `json:"useCustomGasToken,omitempty"`
	GasPayingTokenName         string         `json:"gasPayingTokenName,omitempty"`
	GasPayingTokenSymbol       string         `json:"gasPayingTokenSymbol,omitempty"`
	NativeAssetLiquidityAmount *hexutil.Big   `json:"nativeAssetLiquidityAmount,omitempty"`
	LiquidityControllerOwner   common.Address `json:"liquidityControllerOwner,omitempty"`
}

type GovernanceDeployConfig struct {
	EnableGovernance      bool           `json:"enableGovernance,omitempty"`
	GovernanceTokenSymbol string         `json:"governanceTokenSymbol,omitempty"`
	GovernanceTokenName   string         `json:"governanceTokenName,omitempty"`
	GovernanceTokenOwner  common.Address `json:"governanceTokenOwner,omitempty"`
}

type L1DependenciesConfig struct {
	L1StandardBridgeProxy       common.Address `json:"l1StandardBridgeProxy,omitempty"`
	L1CrossDomainMessengerProxy common.Address `json:"l1CrossDomainMessengerProxy,omitempty"`
	L1ERC721BridgeProxy         common.Address `json:"l1ERC721BridgeProxy,omitempty"`
	SystemConfigProxy           common.Address `json:"systemConfigProxy,omitempty"`
	OptimismPortalProxy         common.Address `json:"optimismPortalProxy,omitempty"`
	ProtocolVersionsProxy       common.Address `json:"protocolVersionsProxy,omitempty"`
	DAChallengeProxy            common.Address `json:"daChallengeProxy,omitempty"`
}

type L2CoreDeployConfig struct {
	L1ChainID                 uint64         `json:"l1ChainID"`
	L2ChainID                 uint64         `json:"l2ChainID"`
	L2BlockTime               uint64         `json:"l2BlockTime"`
	FinalizationPeriodSeconds uint64         `json:"finalizationPeriodSeconds"`
	MaxSequencerDrift         uint64         `json:"maxSequencerDrift"`
	SequencerWindowSize       uint64         `json:"sequencerWindowSize"`
	ChannelTimeoutBedrock     uint64         `json:"channelTimeoutBedrock,omitempty"`
	SystemConfigStartBlock    uint64         `json:"systemConfigStartBlock,omitempty"`
	BatchInboxAddress         common.Address `json:"batchInboxAddress"`
}

type L2GenesisBlockDeployConfig struct {
	L2GenesisBlockNonce         uint64         `json:"l2GenesisBlockNonce"`
	L2GenesisBlockGasLimit      hexutil.Uint64 `json:"l2GenesisBlockGasLimit"`
	L2GenesisBlockDifficulty    uint64         `json:"l2GenesisBlockDifficulty"`
	L2GenesisBlockMixHash       string         `json:"l2GenesisBlockMixHash"`
	L2GenesisBlockNumber        uint64         `json:"l2GenesisBlockNumber"`
	L2GenesisBlockGasUsed       uint64         `json:"l2GenesisBlockGasUsed"`
	L2GenesisBlockParentHash    string         `json:"l2GenesisBlockParentHash"`
	L2GenesisBlockBaseFeePerGas *hexutil.Big   `json:"l2GenesisBlockBaseFeePerGas"`
}

type L2InitializationConfig struct {
	DevDeployConfig
	L2GenesisBlockDeployConfig
	L2VaultsDeployConfig
	GasPriceOracleDeployConfig
	EIP1559DeployConfig
	GovernanceDeployConfig
	GasTokenDeployConfig
	RevenueShareDeployConfig
	UpgradeScheduleDeployConfig
	L2CoreDeployConfig
	OperatorDeployConfig
	OwnershipDeployConfig
	FeeMarketConfig
}

type L2VaultsDeployConfig struct {
	BaseFeeVaultWithdrawalNetwork            string         `json:"baseFeeVaultWithdrawalNetwork,omitempty"`
	L1FeeVaultWithdrawalNetwork              string         `json:"l1FeeVaultWithdrawalNetwork,omitempty"`
	SequencerFeeVaultWithdrawalNetwork       string         `json:"sequencerFeeVaultWithdrawalNetwork,omitempty"`
	OperatorFeeVaultWithdrawalNetwork        string         `json:"operatorFeeVaultWithdrawalNetwork,omitempty"`
	SequencerFeeVaultMinimumWithdrawalAmount *hexutil.Big   `json:"sequencerFeeVaultMinimumWithdrawalAmount,omitempty"`
	BaseFeeVaultMinimumWithdrawalAmount      *hexutil.Big   `json:"baseFeeVaultMinimumWithdrawalAmount,omitempty"`
	L1FeeVaultMinimumWithdrawalAmount        *hexutil.Big   `json:"l1FeeVaultMinimumWithdrawalAmount,omitempty"`
	OperatorFeeVaultMinimumWithdrawalAmount  *hexutil.Big   `json:"operatorFeeVaultMinimumWithdrawalAmount,omitempty"`
	BaseFeeVaultRecipient                    common.Address `json:"baseFeeVaultRecipient,omitempty"`
	L1FeeVaultRecipient                      common.Address `json:"l1FeeVaultRecipient,omitempty"`
	SequencerFeeVaultRecipient               common.Address `json:"sequencerFeeVaultRecipient,omitempty"`
	OperatorFeeVaultRecipient                common.Address `json:"operatorFeeVaultRecipient,omitempty"`
}

type OperatorDeployConfig struct {
	BatchSenderAddress  common.Address `json:"batchSenderAddress,omitempty"`
	P2PSequencerAddress common.Address `json:"p2pSequencerAddress,omitempty"`
}

type OutputOracleDeployConfig struct {
	L2OutputOracleSubmissionInterval uint64         `json:"l2OutputOracleSubmissionInterval"`
	L2OutputOracleStartingTimestamp  int            `json:"l2OutputOracleStartingTimestamp"`
	L2OutputOracleProposer           common.Address `json:"l2OutputOracleProposer"`
	L2OutputOracleChallenger         common.Address `json:"l2OutputOracleChallenger"`
}

type OwnershipDeployConfig struct {
	ProxyAdminOwner  common.Address `json:"proxyAdminOwner,omitempty"`
	FinalSystemOwner common.Address `json:"finalSystemOwner,omitempty"`
}

type RevenueShareDeployConfig struct {
	UseRevenueShare    bool           `json:"useRevenueShare,omitempty"`
	ChainFeesRecipient common.Address `json:"chainFeesRecipient,omitempty"`
}

type SuperchainL1DeployConfig struct {
	SuperchainConfigGuardian common.Address `json:"superchainConfigGuardian,omitempty"`
}

// UpgradeScheduleDeployConfig is defined in upgrade_schedule.go
