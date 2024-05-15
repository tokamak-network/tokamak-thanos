package op_challenger

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/config"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game"
)

// Main is the programmatic entry-point for running op-challenger
func Main(ctx context.Context, logger log.Logger, cfg *config.Config) error {
	if err := cfg.Check(); err != nil {
		return err
	}
	service, err := game.NewService(ctx, logger, cfg)
	if err != nil {
		return fmt.Errorf("failed to create the fault service: %w", err)
	}

	return service.MonitorGame(ctx)
}
