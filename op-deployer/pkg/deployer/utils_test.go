package deployer

import (
	"testing"

	"github.com/tokamak-network/tokamak-thanos/op-deployer/pkg/deployer/flags"
	"github.com/stretchr/testify/require"
)

func TestEnsureDefaultCacheDir(t *testing.T) {
	cacheDir := flags.DefaultCacheDir()
	require.NotNil(t, cacheDir)
}
