package alphabet

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
)

// alphabetUpdater is a [types.OracleUpdater] that exposes a
// method to update onchain oracles with required data.
type alphabetUpdater struct {
	logger log.Logger
}

// NewOracleUpdater returns a new updater.
func NewOracleUpdater(logger log.Logger) *alphabetUpdater {
	return &alphabetUpdater{
		logger: logger,
	}
}

// UpdateOracle updates the oracle with the given data.
func (u *alphabetUpdater) UpdateOracle(ctx context.Context, data *types.PreimageOracleData) error {
	u.logger.Info("alphabet oracle updater called")
	return nil
}
