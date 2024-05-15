package alphabet

import (
	"context"
	"testing"

	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

// TestAlphabetUpdater tests the [alphabetUpdater].
func TestAlphabetUpdater(t *testing.T) {
	logger := testlog.Logger(t, log.LvlInfo)
	updater := NewOracleUpdater(logger)
	require.Nil(t, updater.UpdateOracle(context.Background(), &types.PreimageOracleData{}))
}
