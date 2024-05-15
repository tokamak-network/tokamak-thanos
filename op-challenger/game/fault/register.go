package fault

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/config"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/alphabet"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/cannon"
	faultTypes "github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/scheduler"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/types"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/metrics"
	"github.com/tokamak-network/tokamak-thanos/op-service/txmgr"
)

var (
	cannonGameType   = uint8(0)
	alphabetGameType = uint8(255)
)

type Registry interface {
	RegisterGameType(gameType uint8, creator scheduler.PlayerCreator)
}

func RegisterGameTypes(
	registry Registry,
	ctx context.Context,
	logger log.Logger,
	m metrics.Metricer,
	cfg *config.Config,
	txMgr txmgr.TxManager,
	client bind.ContractCaller,
) {
	if cfg.TraceTypeEnabled(config.TraceTypeCannon) {
		resourceCreator := func(addr common.Address, gameDepth uint64, dir string) (faultTypes.TraceProvider, faultTypes.OracleUpdater, error) {
			provider, err := cannon.NewTraceProvider(ctx, logger, m, cfg, client, dir, addr, gameDepth)
			if err != nil {
				return nil, nil, fmt.Errorf("create cannon trace provider: %w", err)
			}
			updater, err := cannon.NewOracleUpdater(ctx, logger, txMgr, addr, client)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create the cannon updater: %w", err)
			}
			return provider, updater, nil
		}
		playerCreator := func(game types.GameMetadata, dir string) (scheduler.GamePlayer, error) {
			return NewGamePlayer(ctx, logger, m, cfg, dir, game.Proxy, txMgr, client, resourceCreator)
		}
		registry.RegisterGameType(cannonGameType, playerCreator)
	}
	if cfg.TraceTypeEnabled(config.TraceTypeAlphabet) {
		resourceCreator := func(addr common.Address, gameDepth uint64, dir string) (faultTypes.TraceProvider, faultTypes.OracleUpdater, error) {
			provider := alphabet.NewTraceProvider(cfg.AlphabetTrace, gameDepth)
			updater := alphabet.NewOracleUpdater(logger)
			return provider, updater, nil
		}
		playerCreator := func(game types.GameMetadata, dir string) (scheduler.GamePlayer, error) {
			return NewGamePlayer(ctx, logger, m, cfg, dir, game.Proxy, txMgr, client, resourceCreator)
		}
		registry.RegisterGameType(alphabetGameType, playerCreator)
	}
}
