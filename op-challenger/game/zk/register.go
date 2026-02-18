package zk

import (
	"context"
	"fmt"
	"time"

	"github.com/tokamak-network/tokamak-thanos/op-challenger/config"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/client"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/claims"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/contracts"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/generic"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/scheduler"
	gameTypes "github.com/tokamak-network/tokamak-thanos/op-challenger/game/types"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/metrics"
	"github.com/tokamak-network/tokamak-thanos/op-service/txmgr"
	"github.com/ethereum/go-ethereum/log"
)

type ClockReader interface {
	Now() time.Time
}

type Registry interface {
	RegisterGameType(gameType gameTypes.GameType, creator scheduler.PlayerCreator)
	RegisterBondContract(gameType gameTypes.GameType, creator claims.BondContractCreator)
}

type TxSender interface {
	SendAndWaitSimple(txPurpose string, txs ...txmgr.TxCandidate) error
}

func RegisterGameTypes(
	ctx context.Context,
	l1Clock ClockReader,
	logger log.Logger,
	m metrics.Metricer,
	cfg *config.Config,
	registry Registry,
	txSender TxSender,
	clients *client.Provider,
	gameStatusProvider GameStatusProvider,
) error {
	if cfg.GameTypeEnabled(gameTypes.OptimisticZKGameType) {
		registry.RegisterGameType(gameTypes.OptimisticZKGameType, func(game gameTypes.GameMetadata, dir string) (scheduler.GamePlayer, error) {
			rollupClient, syncValidator, err := clients.RollupClients()
			if err != nil {
				return nil, fmt.Errorf("failed to create rollup clients: %w", err)
			}
			contract, err := contracts.NewOptimisticZKDisputeGameContract(m, game.Proxy, clients.MultiCaller())
			if err != nil {
				return nil, fmt.Errorf("failed to create optimistic zk dispute game bindings: %w", err)
			}
			return generic.NewGenericGamePlayer(
				ctx,
				logger,
				game.Proxy,
				contract,
				syncValidator,
				nil,
				clients.L1Client(),
				ActorCreator(l1Clock, rollupClient, gameStatusProvider, contract, txSender),
			)
		})
		registry.RegisterBondContract(gameTypes.OptimisticZKGameType, func(game gameTypes.GameMetadata) (claims.BondContract, error) {
			return contracts.NewOptimisticZKDisputeGameContract(m, game.Proxy, clients.MultiCaller())
		})
	}
	return nil
}
