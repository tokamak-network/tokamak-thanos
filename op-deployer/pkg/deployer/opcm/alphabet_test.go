package opcm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-deployer/pkg/deployer/broadcaster"
	"github.com/tokamak-network/tokamak-thanos/op-deployer/pkg/deployer/testutil"
	"github.com/tokamak-network/tokamak-thanos/op-deployer/pkg/env"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"
)

func TestDeployAlphabetVM(t *testing.T) {
	t.Parallel()

	_, artifacts := testutil.LocalArtifacts(t)

	host, err := env.DefaultScriptHost(
		broadcaster.NoopBroadcaster(),
		testlog.Logger(t, log.LevelInfo),
		common.Address{'D'},
		artifacts,
	)
	require.NoError(t, err)

	input := DeployAlphabetVMInput{
		AbsolutePrestate: common.Hash{'A'},
		PreimageOracle:   common.Address{'O'},
	}

	output, err := DeployAlphabetVM(host, input)
	require.NoError(t, err)

	require.NotEmpty(t, output.AlphabetVM)
}
