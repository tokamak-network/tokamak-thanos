package pipeline

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/standard"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/env"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/broadcaster"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/opcm"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"

	"github.com/ethereum/go-ethereum/common"
)

func GenerateL2Genesis(pEnv *Env, intent *state.Intent, bundle ArtifactsBundle, st *state.State, chainID common.Hash) error {
	lgr := pEnv.Logger.New("stage", "generate-l2-genesis")

	thisIntent, err := intent.Chain(chainID)
	if err != nil {
		return fmt.Errorf("failed to get chain intent: %w", err)
	}

	thisChainState, err := st.Chain(chainID)
	if err != nil {
		return fmt.Errorf("failed to get chain state: %w", err)
	}

	if !shouldGenerateL2Genesis(thisChainState) {
		lgr.Info("L2 genesis generation not needed")
		return nil
	}

	lgr.Info("generating L2 genesis", "id", chainID.Hex())

	initCfg, err := state.CombineDeployConfig(intent, thisIntent, st, thisChainState)
	if err != nil {
		return fmt.Errorf("failed to combine L2 init config: %w", err)
	}

	host, err := env.DefaultScriptHost(
		broadcaster.NoopBroadcaster(),
		pEnv.Logger,
		pEnv.Deployer,
		bundle.L2,
	)
	if err != nil {
		return fmt.Errorf("failed to create L2 script host: %w", err)
	}

	// This is an ugly hack to support holocene. The v1.7.0 predeploy contracts do not support setting the allocs
	// mode as Holocene, even though there are no predeploy changes in Holocene. The v1.7.0 changes are the "official"
	// release of the predeploy contracts, so we need to set the allocs mode to "granite" to avoid having to backport
	// Holocene support into the predeploy contracts.
	var overrideAllocsMode string
	if intent.L2ContractsLocator.IsTag() && intent.L2ContractsLocator.Tag == standard.ContractsV170Beta1L2Tag {
		overrideAllocsMode = "granite"
	}

	if err := opcm.L2Genesis(host, &opcm.L2GenesisInput{
		L1Deployments: opcm.L1Deployments{
			L1CrossDomainMessengerProxy: thisChainState.L1CrossDomainMessengerProxyAddress,
			L1StandardBridgeProxy:       thisChainState.L1StandardBridgeProxyAddress,
			L1ERC721BridgeProxy:         thisChainState.L1ERC721BridgeProxyAddress,
		},
		L2Config:           initCfg.L2InitializationConfig,
		OverrideAllocsMode: overrideAllocsMode,
	}); err != nil {
		return fmt.Errorf("failed to call L2Genesis script: %w", err)
	}

	host.Wipe(pEnv.Deployer)

	dump, err := host.StateDump()
	if err != nil {
		return fmt.Errorf("failed to dump state: %w", err)
	}

	thisChainState.Allocs = &state.GzipData[foundry.ForgeAllocs]{
		Data: dump,
	}

	return nil
}

func shouldGenerateL2Genesis(thisChainState *state.ChainState) bool {
	return thisChainState.Allocs == nil
}
