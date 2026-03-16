package e2eutils

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/tokamak-network/tokamak-thanos/op-e2e/config"
	"github.com/tokamak-network/tokamak-thanos/op-service/predeploys"
)

func TestWriteDefaultJWT(t *testing.T) {
	jwtPath := WriteDefaultJWT(t)
	data, err := os.ReadFile(jwtPath)
	require.NoError(t, err)
	require.Equal(t, "0x"+hex.EncodeToString(testingJWTSecret[:]), string(data))
}

func TestSetup(t *testing.T) {
	tp := &TestParams{
		MaxSequencerDrift:   40,
		SequencerWindowSize: 120,
		ChannelTimeout:      120,
		L1BlockTime:         15,
	}
	dp := MakeDeployParams(t, tp)
	alloc := &AllocParams{PrefundTestUsers: true}
	sd := Setup(t, dp, alloc)
	require.Contains(t, sd.L1Cfg.Alloc, dp.Addresses.Alice)
	require.Equal(t, sd.L1Cfg.Alloc[dp.Addresses.Alice].Balance, Ether(1e12))

	require.Contains(t, sd.L2Cfg.Alloc, dp.Addresses.Alice)
	require.Equal(t, sd.L2Cfg.Alloc[dp.Addresses.Alice].Balance, Ether(1e12))

	require.Contains(t, sd.L1Cfg.Alloc, config.L1Deployments.OptimismPortalProxy)
	require.Contains(t, sd.L2Cfg.Alloc, predeploys.L1BlockAddr)
}

func TestSetupPresets(t *testing.T) {
	if config.DeployConfig == nil {
		t.Skip("deploy config not available, run make devnet-allocs")
	}

	type presetCase struct {
		name    string
		preset  string
		present []common.Address
		absent  []common.Address
		setup   func(*DeployParams)
	}

	testCases := []presetCase{
		{
			name:   "General",
			preset: "general",
			present: []common.Address{
				predeploys.L1BlockAddr,
			},
			absent: []common.Address{
				predeploys.L2UsdcBridgeAddr,
				predeploys.UniswapV3FactoryAddr,
				predeploys.VRFPredeployAddr,
				predeploys.AAEntryPointAddr,
			},
		},
		{
			name:   "DeFi",
			preset: "defi",
			present: []common.Address{
				predeploys.L1BlockAddr,
				predeploys.L2UsdcBridgeAddr,
				predeploys.UniswapV3FactoryAddr,
			},
			absent: []common.Address{
				predeploys.VRFPredeployAddr,
				predeploys.AAEntryPointAddr,
			},
		},
		{
			name:   "Gaming",
			preset: "gaming",
			present: []common.Address{
				predeploys.L1BlockAddr,
				predeploys.VRFPredeployAddr,
				predeploys.AAEntryPointAddr,
			},
			absent: []common.Address{
				predeploys.L2UsdcBridgeAddr,
				predeploys.UniswapV3FactoryAddr,
			},
			setup: func(dp *DeployParams) {
				dp.DeployConfig.VRFAdmin = common.HexToAddress("0x1234567890123456789012345678901234567890")
				dp.DeployConfig.AAPaymasterSigner = common.HexToAddress("0x0000000000000000000000000000000000000002")
			},
		},
		{
			name:   "Full",
			preset: "full",
			present: []common.Address{
				predeploys.L1BlockAddr,
				predeploys.L2UsdcBridgeAddr,
				predeploys.UniswapV3FactoryAddr,
				predeploys.VRFPredeployAddr,
				predeploys.AAEntryPointAddr,
			},
			absent: []common.Address{},
			setup: func(dp *DeployParams) {
				dp.DeployConfig.VRFAdmin = common.HexToAddress("0x1234567890123456789012345678901234567890")
				dp.DeployConfig.AAPaymasterSigner = common.HexToAddress("0x0000000000000000000000000000000000000002")
			},
		},
	}

	tp := &TestParams{
		MaxSequencerDrift:   40,
		SequencerWindowSize: 120,
		ChannelTimeout:      120,
		L1BlockTime:         15,
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			dp := MakeDeployParams(t, tp)
			dp.DeployConfig.Preset = tt.preset
			if tt.setup != nil {
				tt.setup(dp)
			}
			sd := Setup(t, dp, &AllocParams{PrefundTestUsers: true})

			for _, addr := range tt.present {
				_, ok := sd.L2Cfg.Alloc[addr]
				require.True(t, ok, "expected predeploy %s to be present for preset %s", addr, tt.preset)
			}
			for _, addr := range tt.absent {
				_, ok := sd.L2Cfg.Alloc[addr]
				require.False(t, ok, "expected predeploy %s to be absent for preset %s", addr, tt.preset)
			}
		})
	}
}
