package genesis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/stretchr/testify/require"

	"github.com/tokamak-network/tokamak-thanos/op-bindings/predeploys"
)

func TestConfigDataMarshalUnmarshal(t *testing.T) {
	b, err := os.ReadFile("testdata/test-deploy-config-full.json")
	require.NoError(t, err)

	dec := json.NewDecoder(bytes.NewReader(b))
	decoded := new(DeployConfig)
	require.NoError(t, dec.Decode(decoded))
	require.EqualValues(t, "non-default value", string(decoded.L2GenesisBlockExtraData))

	encoded, err := json.MarshalIndent(decoded, "", "  ")
	require.NoError(t, err)
	require.JSONEq(t, string(b), string(encoded))
}

func TestUnmarshalL1StartingBlockTag(t *testing.T) {
	decoded := new(DeployConfig)
	require.NoError(t, json.Unmarshal([]byte(`{"l1StartingBlockTag": "earliest"}`), decoded))
	require.EqualValues(t, rpc.EarliestBlockNumber, *decoded.L1StartingBlockTag.BlockNumber)
	h := "0x86c7263d87140ca7cd9bf1bc9e95a435a7a0efc0ae2afaf64920c5b59a6393d4"
	require.NoError(t, json.Unmarshal([]byte(fmt.Sprintf(`{"l1StartingBlockTag": "%s"}`, h)), decoded))
	require.EqualValues(t, common.HexToHash(h), *decoded.L1StartingBlockTag.BlockHash)
}

func TestRegolithTimeZero(t *testing.T) {
	regolithOffset := hexutil.Uint64(0)
	config := &DeployConfig{L2GenesisRegolithTimeOffset: &regolithOffset}
	require.Equal(t, uint64(0), *config.RegolithTime(1234))
}

func TestRegolithTimeAsOffset(t *testing.T) {
	regolithOffset := hexutil.Uint64(1500)
	config := &DeployConfig{L2GenesisRegolithTimeOffset: &regolithOffset}
	require.Equal(t, uint64(1500+5000), *config.RegolithTime(5000))
}

func TestDeployConfigPresetID(t *testing.T) {
	cfg := &DeployConfig{}
	require.Equal(t, "defi", cfg.PresetID()) // empty string → "defi" for backward compat

	cfg.Preset = "general"
	require.Equal(t, "general", cfg.PresetID())

	cfg.Preset = "defi"
	require.Equal(t, "defi", cfg.PresetID())

	cfg.Preset = "gaming"
	require.Equal(t, "gaming", cfg.PresetID())

	cfg.Preset = "full"
	require.Equal(t, "full", cfg.PresetID())
}

func TestDeployConfigPresetValidation(t *testing.T) {
	cfg := &DeployConfig{}

	for _, id := range []string{"", "general", "defi", "gaming", "full"} {
		cfg.Preset = id
		require.NoError(t, cfg.validatePreset())
	}

	cfg.Preset = "invalid"
	require.Error(t, cfg.validatePreset())
}

func TestCanyonTimeZero(t *testing.T) {
	canyonOffset := hexutil.Uint64(0)
	config := &DeployConfig{L2GenesisCanyonTimeOffset: &canyonOffset}
	require.Equal(t, uint64(0), *config.CanyonTime(1234))
}

func TestCanyonTimeOffset(t *testing.T) {
	canyonOffset := hexutil.Uint64(1500)
	config := &DeployConfig{L2GenesisCanyonTimeOffset: &canyonOffset}
	require.Equal(t, uint64(1234+1500), *config.CanyonTime(1234))
}

// TestCopy will copy a DeployConfig and ensure that the copy is equal to the original.
func TestCopy(t *testing.T) {
	b, err := os.ReadFile("testdata/test-deploy-config-full.json")
	require.NoError(t, err)

	decoded := new(DeployConfig)
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(decoded))

	cpy := decoded.Copy()
	require.EqualValues(t, decoded, cpy)

	offset := hexutil.Uint64(100)
	cpy.L2GenesisRegolithTimeOffset = &offset
	require.NotEqual(t, decoded, cpy)
}

// TestL1Deployments ensures that NewL1Deployments can read a JSON file
// from disk and deserialize all of the key/value pairs correctly.
func TestL1Deployments(t *testing.T) {
	deployments, err := NewL1Deployments("testdata/l1-deployments.json")
	require.NoError(t, err)

	require.NotEqual(t, deployments.AddressManager, common.Address{})
	require.NotEqual(t, deployments.DisputeGameFactory, common.Address{})
	require.NotEqual(t, deployments.DisputeGameFactoryProxy, common.Address{})
	require.NotEqual(t, deployments.L1CrossDomainMessenger, common.Address{})
	require.NotEqual(t, deployments.L1CrossDomainMessengerProxy, common.Address{})
	require.NotEqual(t, deployments.L1ERC721Bridge, common.Address{})
	require.NotEqual(t, deployments.L1ERC721BridgeProxy, common.Address{})
	require.NotEqual(t, deployments.L1StandardBridge, common.Address{})
	require.NotEqual(t, deployments.L1StandardBridgeProxy, common.Address{})
	require.NotEqual(t, deployments.L2OutputOracle, common.Address{})
	require.NotEqual(t, deployments.L2OutputOracleProxy, common.Address{})
	require.NotEqual(t, deployments.OptimismMintableERC20Factory, common.Address{})
	require.NotEqual(t, deployments.OptimismMintableERC20FactoryProxy, common.Address{})
	require.NotEqual(t, deployments.OptimismPortal, common.Address{})
	require.NotEqual(t, deployments.OptimismPortalProxy, common.Address{})
	require.NotEqual(t, deployments.ProxyAdmin, common.Address{})
	require.NotEqual(t, deployments.SystemConfig, common.Address{})
	require.NotEqual(t, deployments.SystemConfigProxy, common.Address{})
	require.NotEqual(t, deployments.ProtocolVersions, common.Address{})
	require.NotEqual(t, deployments.ProtocolVersionsProxy, common.Address{})

	require.Equal(t, "AddressManager", deployments.GetName(deployments.AddressManager))
	require.Equal(t, "OptimismPortalProxy", deployments.GetName(deployments.OptimismPortalProxy))
	// One that doesn't exist returns empty string
	require.Equal(t, "", deployments.GetName(common.Address{19: 0xff}))
}

func TestOpChainConfig_Defaults(t *testing.T) {
	cfg := &DeployConfig{}
	opCfg := cfg.OpChainConfig()

	require.Equal(t, uint64(50), opCfg.EIP1559Denominator)
	require.Equal(t, uint64(250), opCfg.EIP1559DenominatorCanyon)
	require.Equal(t, uint64(10), opCfg.EIP1559Elasticity)
}

func TestOpChainConfig_UsesConfiguredValues(t *testing.T) {
	cfg := &DeployConfig{
		EIP1559Denominator:       9,
		EIP1559DenominatorCanyon: 99,
		EIP1559Elasticity:        19,
	}
	opCfg := cfg.OpChainConfig()

	require.Equal(t, uint64(9), opCfg.EIP1559Denominator)
	require.Equal(t, uint64(99), opCfg.EIP1559DenominatorCanyon)
	require.Equal(t, uint64(19), opCfg.EIP1559Elasticity)
}

func TestNewL2StorageConfigGeneralPreset(t *testing.T) {
	config, err := NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.NoError(t, err)
	config.Preset = PresetGeneral

	block := types.NewBlockWithHeader(&types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(1000000000),
	})

	storage, err := NewL2StorageConfig(config, block)
	require.NoError(t, err)

	_, hasUsdcBridge := storage["L2UsdcBridge"]
	require.False(t, hasUsdcBridge, "General preset must not have L2UsdcBridge storage")

	_, hasMasterMinter := storage["MasterMinter"]
	require.False(t, hasMasterMinter, "General preset must not have MasterMinter storage")

	_, hasFiatToken := storage["FiatTokenV2_2"]
	require.False(t, hasFiatToken, "General preset must not have FiatTokenV2_2 storage")

	_, hasUniswap := storage["UniswapV3Factory"]
	require.False(t, hasUniswap, "General preset must not have UniswapV3Factory storage")

	_, hasVRFCoordinator := storage["VRFCoordinator"]
	require.False(t, hasVRFCoordinator, "General preset must not have VRFCoordinator storage")

	_, hasVRFPredeploy := storage["VRFPredeploy"]
	require.False(t, hasVRFPredeploy, "General preset must not have VRFPredeploy storage")
}

func TestNewL2StorageConfigDefiPreset(t *testing.T) {
	config, err := NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.NoError(t, err)
	config.Preset = PresetDeFi

	block := types.NewBlockWithHeader(&types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(1000000000),
	})

	storage, err := NewL2StorageConfig(config, block)
	require.NoError(t, err)

	_, hasUsdcBridge := storage["L2UsdcBridge"]
	require.True(t, hasUsdcBridge, "DeFi preset must have L2UsdcBridge storage")

	_, hasMasterMinter := storage["MasterMinter"]
	require.True(t, hasMasterMinter, "DeFi preset must have MasterMinter storage")

	_, hasFiatToken := storage["FiatTokenV2_2"]
	require.True(t, hasFiatToken, "DeFi preset must have FiatTokenV2_2 storage")

	_, hasUniswap := storage["UniswapV3Factory"]
	require.True(t, hasUniswap, "DeFi preset must have UniswapV3Factory storage")

	_, hasPaymaster := storage["VerifyingPaymasterPredeploy"]
	require.False(t, hasPaymaster, "DeFi preset must not have VerifyingPaymasterPredeploy storage")

	_, hasVRFCoordinator := storage["VRFCoordinator"]
	require.False(t, hasVRFCoordinator, "DeFi preset must not have VRFCoordinator storage")

	_, hasVRFPredeploy := storage["VRFPredeploy"]
	require.False(t, hasVRFPredeploy, "DeFi preset must not have VRFPredeploy storage")
}

func TestAAPredeployAddressConstants(t *testing.T) {
	// ⚠️ existing predeploys.EntryPoint = 0x5FF... (v0.7 canonical) is separate.
	//    v0.8 predeploy uses AAEntryPoint constant name.
	require.Equal(t, "0x4200000000000000000000000000000000000063",
		predeploys.AAEntryPoint)
	require.Equal(t, "0x4200000000000000000000000000000000000064",
		predeploys.VerifyingPaymasterPredeploy)
	require.Equal(t, "0x4200000000000000000000000000000000000065",
		predeploys.Simple7702Account)
	require.NotNil(t, predeploys.AAEntryPointAddr)
	require.NotNil(t, predeploys.VerifyingPaymasterPredeployAddr)
	require.NotNil(t, predeploys.Simple7702AccountAddr)
}

func TestNewL2StorageConfigGamingPreset_HasAAStorage(t *testing.T) {
	config, err := NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.NoError(t, err)
	config.Preset = PresetGaming
	config.VRFAdmin = common.HexToAddress("0x0000000000000000000000000000000000000001")
	config.AAPaymasterSigner = common.HexToAddress("0x0000000000000000000000000000000000000002")

	block := types.NewBlockWithHeader(&types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(1000000000),
	})

	storage, err := NewL2StorageConfig(config, block)
	require.NoError(t, err)

	_, hasPaymaster := storage["VerifyingPaymasterPredeploy"]
	require.True(t, hasPaymaster, "Gaming preset must have VerifyingPaymasterPredeploy storage")
}

func TestNewL2StorageConfigGeneralPreset_NoAAStorage(t *testing.T) {
	config, err := NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.NoError(t, err)
	config.Preset = PresetGeneral

	block := types.NewBlockWithHeader(&types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(1000000000),
	})

	storage, err := NewL2StorageConfig(config, block)
	require.NoError(t, err)

	_, hasPaymaster := storage["VerifyingPaymasterPredeploy"]
	require.False(t, hasPaymaster, "General preset must not have VerifyingPaymasterPredeploy storage")
}

func TestNewL2StorageConfigFullPreset(t *testing.T) {
	config, err := NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.NoError(t, err)
	config.Preset = PresetFull
	config.VRFAdmin = common.HexToAddress("0x0000000000000000000000000000000000000001")
	config.AAPaymasterSigner = common.HexToAddress("0x0000000000000000000000000000000000000002")

	block := types.NewBlockWithHeader(&types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(1000000000),
	})

	storage, err := NewL2StorageConfig(config, block)
	require.NoError(t, err)

	_, hasUsdcBridge := storage["L2UsdcBridge"]
	require.True(t, hasUsdcBridge, "Full preset must have L2UsdcBridge storage")

	_, hasMasterMinter := storage["MasterMinter"]
	require.True(t, hasMasterMinter, "Full preset must have MasterMinter storage")

	_, hasFiatToken := storage["FiatTokenV2_2"]
	require.True(t, hasFiatToken, "Full preset must have FiatTokenV2_2 storage")

	_, hasUniswap := storage["UniswapV3Factory"]
	require.True(t, hasUniswap, "Full preset must have UniswapV3Factory storage")

	_, hasPaymaster := storage["VerifyingPaymasterPredeploy"]
	require.True(t, hasPaymaster, "Full preset must have VerifyingPaymasterPredeploy storage")

	_, hasVRFCoordinator := storage["VRFCoordinator"]
	require.True(t, hasVRFCoordinator, "Full preset must have VRFCoordinator storage")

	_, hasVRFPredeploy := storage["VRFPredeploy"]
	require.True(t, hasVRFPredeploy, "Full preset must have VRFPredeploy storage")
}

func TestNewL2StorageConfigGamingPreset_HasVRFStorage(t *testing.T) {
	config, err := NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.NoError(t, err)
	config.Preset = PresetGaming
	config.VRFAdmin = common.HexToAddress("0x0000000000000000000000000000000000000001")
	config.AAPaymasterSigner = common.HexToAddress("0x0000000000000000000000000000000000000002")

	block := types.NewBlockWithHeader(&types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(1000000000),
	})

	storage, err := NewL2StorageConfig(config, block)
	require.NoError(t, err)

	_, hasVRFCoordinator := storage["VRFCoordinator"]
	require.True(t, hasVRFCoordinator, "Gaming preset must have VRFCoordinator storage")

	_, hasVRFPredeploy := storage["VRFPredeploy"]
	require.True(t, hasVRFPredeploy, "Gaming preset must have VRFPredeploy storage")
}

func TestRollupConfig_SetsChainOpConfig(t *testing.T) {
	b, err := os.ReadFile("testdata/test-deploy-config-full.json")
	require.NoError(t, err)

	cfg := new(DeployConfig)
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(cfg))

	l1StartBlock := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(1),
		Time:   100,
	})
	rollupCfg, err := cfg.RollupConfig(l1StartBlock, common.HexToHash("0x1234"), 0)
	require.NoError(t, err)
	require.NotNil(t, rollupCfg.ChainOpConfig)
	require.Equal(t, cfg.OpChainConfig(), rollupCfg.ChainOpConfig)
}
