package pipeline

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/addresses"
	"github.com/tokamak-network/tokamak-thanos/op-service/jsonutil"

	"github.com/tokamak-network/tokamak-thanos/op-deployer/pkg/deployer/opcm"
	"github.com/tokamak-network/tokamak-thanos/op-deployer/pkg/deployer/standard"
	"github.com/tokamak-network/tokamak-thanos/op-deployer/pkg/deployer/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// uint96MaxValue represents the maximum value for uint96 (2^96 - 1)
var uint96MaxValue = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 96), big.NewInt(1))

func mustHexBigFromHex(value string) *big.Int {
	// Handle both hex strings (with 0x prefix) and decimal strings
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		num := hexutil.MustDecodeBig(value)
		return num
	} else {
		// Parse as decimal string
		num, ok := new(big.Int).SetString(value, 10)
		if !ok {
			panic(fmt.Sprintf("invalid number format: %s", value))
		}
		return num
	}
}

// mustBigIntWithUint96Limit ensures *big.Int doesn't exceed uint96 max value for contract compatibility
func mustBigIntWithUint96Limit(value *big.Int, fieldName string) *big.Int {
	if value.Cmp(uint96MaxValue) > 0 {
		panic(fmt.Sprintf("%s value %s exceeds uint96 maximum value %s (contract limitation)", fieldName, value.String(), uint96MaxValue.String()))
	}
	return value
}

func DeployOPChain(env *Env, intent *state.Intent, st *state.State, chainID common.Hash) error {
	lgr := env.Logger.New("stage", "deploy-opchain")

	if !shouldDeployOPChain(st, chainID) {
		lgr.Info("opchain deployment not needed")
		return nil
	}

	thisIntent, err := intent.Chain(chainID)
	if err != nil {
		return fmt.Errorf("failed to get chain intent: %w", err)
	}

	var dco opcm.DeployOPChainOutput
	lgr.Info("deploying OP chain using local allocs", "id", chainID.Hex())

	dci, err := makeDCI(intent, thisIntent, chainID, st)
	if err != nil {
		return fmt.Errorf("error making deploy OP chain input: %w", err)
	}

	dco, err = opcm.DeployOPChain(env.L1ScriptHost, dci)
	if err != nil {
		return fmt.Errorf("error deploying OP chain: %w", err)
	}

	st.Chains = append(st.Chains, makeChainState(chainID, dco))

	var release string
	if intent.L1ContractsLocator.IsEmbedded() {
		release = standard.CurrentTag
	} else {
		release = "dev"
	}

	readInput := opcm.ReadImplementationAddressesInput{
		DeployOPChainOutput: dco,
		Opcm:                dci.Opcm,
		Release:             release,
	}
	impls, err := opcm.ReadImplementationAddresses(env.L1ScriptHost, readInput)
	if err != nil {
		return fmt.Errorf("failed to read implementation addresses: %w", err)
	}

	st.ImplementationsDeployment.DelayedWethImpl = impls.DelayedWETH
	st.ImplementationsDeployment.OptimismPortalImpl = impls.OptimismPortal
	st.ImplementationsDeployment.EthLockboxImpl = impls.ETHLockbox
	st.ImplementationsDeployment.SystemConfigImpl = impls.SystemConfig
	st.ImplementationsDeployment.L1CrossDomainMessengerImpl = impls.L1CrossDomainMessenger
	st.ImplementationsDeployment.L1Erc721BridgeImpl = impls.L1ERC721Bridge
	st.ImplementationsDeployment.L1StandardBridgeImpl = impls.L1StandardBridge
	st.ImplementationsDeployment.OptimismMintableErc20FactoryImpl = impls.OptimismMintableERC20Factory
	st.ImplementationsDeployment.DisputeGameFactoryImpl = impls.DisputeGameFactory
	st.ImplementationsDeployment.MipsImpl = impls.MipsSingleton
	st.ImplementationsDeployment.PreimageOracleImpl = impls.PreimageOracleSingleton

	return nil
}

func makeDCI(intent *state.Intent, thisIntent *state.ChainIntent, chainID common.Hash, st *state.State) (opcm.DeployOPChainInput, error) {
	proofParams, err := jsonutil.MergeJSON(
		state.ChainProofParams{
			DisputeGameType:         standard.DisputeGameType,
			DisputeAbsolutePrestate: standard.DisputeAbsolutePrestate,
			DisputeMaxGameDepth:     standard.DisputeMaxGameDepth,
			DisputeSplitDepth:       standard.DisputeSplitDepth,
			DisputeClockExtension:   standard.DisputeClockExtension,
			DisputeMaxClockDuration: standard.DisputeMaxClockDuration,
			// RAT defaults
			DeployRAT:                standard.DeployRAT,
			PerTestBondAmount:        standard.PerTestBondAmount,
			EvidenceSubmissionPeriod: standard.EvidenceSubmissionPeriod,
			MinimumStakingBalance:    standard.MinimumStakingBalance,
			RatTriggerProbability:    standard.RatTriggerProbability,
			RatManager:               common.Address{}, // Default to zero address
		},
		intent.GlobalDeployOverrides,
		thisIntent.DeployOverrides,
	)
	if err != nil {
		return opcm.DeployOPChainInput{}, fmt.Errorf("error merging proof params from overrides: %w", err)
	}

	return opcm.DeployOPChainInput{
		OpChainProxyAdminOwner:       thisIntent.Roles.L1ProxyAdminOwner,
		SystemConfigOwner:            thisIntent.Roles.SystemConfigOwner,
		Batcher:                      thisIntent.Roles.Batcher,
		UnsafeBlockSigner:            thisIntent.Roles.UnsafeBlockSigner,
		Proposer:                     thisIntent.Roles.Proposer,
		Challenger:                   thisIntent.Roles.Challenger,
		BasefeeScalar:                standard.BasefeeScalar,
		BlobBaseFeeScalar:            standard.BlobBaseFeeScalar,
		L2ChainId:                    chainID.Big(),
		Opcm:                         st.ImplementationsDeployment.OpcmImpl,
		SaltMixer:                    st.Create2Salt.String(), // passing through salt generated at state initialization
		GasLimit:                     standard.GasLimit,
		DisputeGameType:              proofParams.DisputeGameType,
		DisputeAbsolutePrestate:      proofParams.DisputeAbsolutePrestate,
		DisputeMaxGameDepth:          proofParams.DisputeMaxGameDepth,
		DisputeSplitDepth:            proofParams.DisputeSplitDepth,
		DisputeClockExtension:        proofParams.DisputeClockExtension,   // 3 hours (input in seconds)
		DisputeMaxClockDuration:      proofParams.DisputeMaxClockDuration, // 3.5 days (input in seconds)
		AllowCustomDisputeParameters: proofParams.DangerouslyAllowCustomDisputeParameters,
		OperatorFeeScalar:            thisIntent.OperatorFeeScalar,
		OperatorFeeConstant:          thisIntent.OperatorFeeConstant,
		// RAT configuration - TODO: Add RAT parameters from intent
		DeployRAT:                proofParams.DeployRAT, // Default to false for now
		PerTestBondAmount:        mustBigIntWithUint96Limit(mustHexBigFromHex(proofParams.PerTestBondAmount), "PerTestBondAmount"),
		EvidenceSubmissionPeriod: big.NewInt(int64(proofParams.EvidenceSubmissionPeriod)),
		MinimumStakingBalance:    mustHexBigFromHex(proofParams.MinimumStakingBalance),
		RatTriggerProbability:    mustHexBigFromHex(proofParams.RatTriggerProbability),
		RatManager:               proofParams.RatManager,
	}, nil
}

func makeChainState(chainID common.Hash, dco opcm.DeployOPChainOutput) *state.ChainState {
	opChainContracts := addresses.OpChainContracts{}
	opChainContracts.OpChainProxyAdminImpl = dco.OpChainProxyAdmin
	opChainContracts.AddressManagerImpl = dco.AddressManager
	opChainContracts.L1Erc721BridgeProxy = dco.L1ERC721BridgeProxy
	opChainContracts.SystemConfigProxy = dco.SystemConfigProxy
	opChainContracts.OptimismMintableErc20FactoryProxy = dco.OptimismMintableERC20FactoryProxy
	opChainContracts.L1StandardBridgeProxy = dco.L1StandardBridgeProxy
	opChainContracts.L1CrossDomainMessengerProxy = dco.L1CrossDomainMessengerProxy
	opChainContracts.OptimismPortalProxy = dco.OptimismPortalProxy
	opChainContracts.EthLockboxProxy = dco.ETHLockboxProxy
	opChainContracts.DisputeGameFactoryProxy = dco.DisputeGameFactoryProxy
	opChainContracts.AnchorStateRegistryProxy = dco.AnchorStateRegistryProxy
	opChainContracts.FaultDisputeGameImpl = dco.FaultDisputeGame
	opChainContracts.PermissionedDisputeGameImpl = dco.PermissionedDisputeGame
	opChainContracts.DelayedWethPermissionedGameProxy = dco.DelayedWETHPermissionedGameProxy
	opChainContracts.DelayedWethPermissionlessGameProxy = dco.DelayedWETHPermissionlessGameProxy
	opChainContracts.RATProxy = dco.RATProxy

	return &state.ChainState{
		ID:               chainID,
		OpChainContracts: opChainContracts,
	}
}

func shouldDeployOPChain(st *state.State, chainID common.Hash) bool {
	for _, chain := range st.Chains {
		if chain.ID == chainID {
			return false
		}
	}

	return true
}
