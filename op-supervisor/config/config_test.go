package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"

	"github.com/tokamak-network/tokamak-thanos/op-service/metrics"
	"github.com/tokamak-network/tokamak-thanos/op-service/oppprof"
	"github.com/tokamak-network/tokamak-thanos/op-service/rpc"
	"github.com/tokamak-network/tokamak-thanos/op-supervisor/supervisor/backend/depset"
	"github.com/tokamak-network/tokamak-thanos/op-supervisor/supervisor/backend/syncnode"
)

func TestDefaultConfigIsValid(t *testing.T) {
	cfg := validConfig()
	require.NoError(t, cfg.Check())
}

func TestRequireSyncSources(t *testing.T) {
	cfg := validConfig()
	cfg.SyncSources = nil
	require.ErrorIs(t, cfg.Check(), ErrMissingSyncSources)
}

func TestRequireDependencySet(t *testing.T) {
	cfg := validConfig()
	cfg.DependencySetSource = nil
	require.ErrorIs(t, cfg.Check(), ErrMissingDependencySet)
}

func TestRequireDatadir(t *testing.T) {
	cfg := validConfig()
	cfg.Datadir = ""
	require.ErrorIs(t, cfg.Check(), ErrMissingDatadir)
}

func TestValidateMetricsConfig(t *testing.T) {
	cfg := validConfig()
	cfg.MetricsConfig.Enabled = true
	cfg.MetricsConfig.ListenPort = -1
	require.ErrorIs(t, cfg.Check(), metrics.ErrInvalidPort)
}

func TestValidatePprofConfig(t *testing.T) {
	cfg := validConfig()
	cfg.PprofConfig.ListenEnabled = true
	cfg.PprofConfig.ListenPort = -1
	require.ErrorIs(t, cfg.Check(), oppprof.ErrInvalidPort)
}

func TestValidateRPCConfig(t *testing.T) {
	cfg := validConfig()
	cfg.RPC.ListenPort = -1
	require.ErrorIs(t, cfg.Check(), rpc.ErrInvalidPort)
}

func validConfig() *Config {
	depSet, err := depset.NewStaticConfigDependencySet(map[eth.ChainID]*depset.StaticConfigDependency{
		eth.ChainIDFromUInt64(900): &depset.StaticConfigDependency{
			ChainIndex:     900,
			ActivationTime: 0,
			HistoryMinTime: 0,
		},
	})
	if err != nil {
		panic(err)
	}
	// Should be valid using only the required arguments passed in via the constructor.
	return NewConfig("http://localhost:8545", &syncnode.CLISyncNodes{}, depSet, "./supervisor_testdir")
}
