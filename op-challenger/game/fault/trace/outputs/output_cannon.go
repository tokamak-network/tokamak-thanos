package outputs

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/config"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/contracts"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/cannon"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/split"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/metrics"
)

func NewOutputCannonTraceAccessor(
	logger log.Logger,
	m metrics.Metricer,
	cfg *config.Config,
	l2Client cannon.L2HeaderSource,
	contract cannon.L1HeadSource,
	prestateProvider types.PrestateProvider,
	rollupClient OutputRootProvider,
	dir string,
	splitDepth types.Depth,
	prestateBlock uint64,
	poststateBlock uint64,
) (*trace.Accessor, error) {
	outputProvider := NewTraceProviderFromInputs(logger, prestateProvider, rollupClient, splitDepth, prestateBlock, poststateBlock)
	cannonCreator := func(ctx context.Context, localContext common.Hash, depth types.Depth, agreed contracts.Proposal, claimed contracts.Proposal) (types.TraceProvider, error) {
		logger := logger.New("pre", agreed.OutputRoot, "post", claimed.OutputRoot, "localContext", localContext)
		subdir := filepath.Join(dir, localContext.Hex())
		localInputs, err := cannon.FetchLocalInputsFromProposals(ctx, contract, l2Client, agreed, claimed)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch cannon local inputs: %w", err)
		}
		provider := cannon.NewTraceProvider(logger, m, cfg, localInputs, subdir, depth)
		return provider, nil
	}

	cache := NewProviderCache(m, "output_cannon_provider", cannonCreator)
	selector := split.NewSplitProviderSelector(outputProvider, splitDepth, OutputRootSplitAdapter(outputProvider, cache.GetOrCreate))
	return trace.NewAccessor(selector), nil
}
