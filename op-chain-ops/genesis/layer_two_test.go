package genesis_test

import (
	"context"
	"encoding/json"
	"flag"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/predeploys"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// usdcAddrs contains the 4 USDC Bridge predeploy addresses.
var usdcAddrs = []common.Address{
	predeploys.L2UsdcBridgeAddr,
	predeploys.SignatureCheckerAddr,
	predeploys.MasterMinterAddr,
	predeploys.FiatTokenV2_2Addr,
}

// uniswapAddrs contains the 10 Uniswap V3 predeploy addresses.
var uniswapAddrs = []common.Address{
	predeploys.QuoterV2Addr,
	predeploys.SwapRouter02Addr,
	predeploys.UniswapV3FactoryAddr,
	predeploys.NFTDescriptorAddr,
	predeploys.NonfungiblePositionManagerAddr,
	predeploys.NonfungibleTokenPositionDescriptorAddr,
	predeploys.TickLensAddr,
	predeploys.UniswapInterfaceMulticallAddr,
	predeploys.UniversalRouterAddr,
	predeploys.UnsupportedProtocolAddr,
}

// vrfAddrs contains the 2 VRF predeploy addresses.
var vrfAddrs = []common.Address{
	predeploys.VRFPredeployAddr,
	predeploys.VRFCoordinatorAddr,
}

// aaAddrs contains the 3 AA predeploy addresses.
var aaAddrs = []common.Address{
	predeploys.AAEntryPointAddr,
	predeploys.VerifyingPaymasterPredeployAddr,
	predeploys.Simple7702AccountAddr,
}

var writeFile bool

func init() {
	flag.BoolVar(&writeFile, "write-file", false, "write the genesis file to disk")
}

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

// Tests the BuildL2MainnetGenesis factory with the provided config.
func testBuildL2Genesis(t *testing.T, config *genesis.DeployConfig) *core.Genesis {
	backend := simulated.NewBackend(
		types.GenesisAlloc{
			crypto.PubkeyToAddress(testKey.PublicKey): {Balance: big.NewInt(10000000000000000)},
		},
	)
	defer backend.Close()

	client := backend.Client()
	block, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)

	gen, err := genesis.BuildL2Genesis(config, block)
	require.Nil(t, err)
	require.NotNil(t, gen)

	proxyBytecode, err := bindings.GetDeployedBytecode("Proxy")
	require.NoError(t, err)

	transparentUpgradeableProxyBytecode, err := bindings.GetDeployedBytecode("TransparentUpgradeableProxy")
	require.NoError(t, err)

	l2UsdcBridgeProxyBytecode, err := bindings.GetDeployedBytecode("L2UsdcBridgeProxy")
	require.NoError(t, err)

	fiatTokenProxyBytecode, err := bindings.GetDeployedBytecode("FiatTokenProxy")
	require.NoError(t, err)

	for name, predeploy := range predeploys.Predeploys {
		// If the predeploy has an Enabled function and it returns false for this config,
		// assert it's absent from genesis and skip further checks.
		if predeploy.Enabled != nil && !predeploy.Enabled(config) {
			_, ok := gen.Alloc[predeploy.Address]
			require.False(t, ok, "disabled predeploy should not be in genesis: %s", name)
			continue
		}

		addr := predeploy.Address

		account, ok := gen.Alloc[addr]
		require.Equal(t, true, ok, name)
		require.Greater(t, len(account.Code), 0)

		adminSlot, ok := account.Storage[genesis.AdminSlot]
		isProxy := !predeploy.ProxyDisabled ||
			(!config.EnableGovernance && addr == predeploys.GovernanceTokenAddr)

		if isProxy {
			switch addr {
			case predeploys.L2UsdcBridgeAddr:
				require.Equal(t, true, ok, name)
				require.Equal(t, eth.AddressAsLeftPaddedHash(predeploys.ProxyAdminAddr), adminSlot)
				require.Equal(t, l2UsdcBridgeProxyBytecode, account.Code)
			case predeploys.FiatTokenV2_2Addr:
				adminSlotForZepplin, ok := account.Storage[genesis.AdminSlotForZepplin]
				require.Equal(t, true, ok, name)
				require.Equal(t, eth.AddressAsLeftPaddedHash(predeploys.ProxyAdminAddr), adminSlotForZepplin)
				require.Equal(t, fiatTokenProxyBytecode, account.Code)
			case predeploys.NonfungibleTokenPositionDescriptorAddr:
				require.Equal(t, true, ok, name)
				require.Equal(t, eth.AddressAsLeftPaddedHash(predeploys.ProxyAdminAddr), adminSlot)
				require.Equal(t, transparentUpgradeableProxyBytecode, account.Code)

			default:
				require.Equal(t, true, ok, name)
				require.Equal(t, eth.AddressAsLeftPaddedHash(predeploys.ProxyAdminAddr), adminSlot)
				require.Equal(t, proxyBytecode, account.Code)
			}
		} else {
			require.Equal(t, false, ok, name)
			require.NotEqual(t, proxyBytecode, account.Code, name)
		}
	}

	// All of the precompile addresses should be funded with a single wei
	if config.SetPrecompileBalances {
		for i := 0; i < genesis.PrecompileCount; i++ {
			addr := common.BytesToAddress([]byte{byte(i)})
			require.Equal(t, common.Big1, gen.Alloc[addr].Balance)
		}
	}

	create2Deployer := gen.Alloc[predeploys.Create2DeployerAddr]
	codeHash := crypto.Keccak256Hash(create2Deployer.Code)
	require.Equal(t, codeHash, bindings.Create2DeployerCodeHash)

	if writeFile {
		file, _ := json.MarshalIndent(gen, "", " ")
		_ = os.WriteFile("genesis.json", file, 0644)
	}
	return gen
}

func TestBuildL2MainnetGenesis(t *testing.T) {
	config, err := genesis.NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.Nil(t, err)
	config.EnableGovernance = true
	config.FundDevAccounts = false
	gen := testBuildL2Genesis(t, config)
	require.Equal(t, 2076, len(gen.Alloc))
}

func TestBuildL2MainnetNoGovernanceGenesis(t *testing.T) {
	config, err := genesis.NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.Nil(t, err)
	config.EnableGovernance = false
	config.FundDevAccounts = false
	gen := testBuildL2Genesis(t, config)
	// GovernanceToken is excluded when EnableGovernance=false, so alloc count is 1 less than with governance.
	require.Equal(t, 2075, len(gen.Alloc))
}

func TestBuildL2GenesisGeneralPreset(t *testing.T) {
	config, err := genesis.NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.Nil(t, err)
	config.Preset = genesis.PresetGeneral
	config.EnableGovernance = false
	config.FundDevAccounts = false

	gen := testBuildL2Genesis(t, config)

	// General preset: USDC Bridge predeploys must NOT be in genesis
	for _, addr := range usdcAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "USDC address should not be in General preset: %s", addr.Hex())
	}

	// General preset: Uniswap V3 predeploys must NOT be in genesis
	for _, addr := range uniswapAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "Uniswap address should not be in General preset: %s", addr.Hex())
	}

	// General preset: VRF predeploys must NOT be in genesis
	for _, addr := range vrfAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "VRF address should not be in General preset: %s", addr.Hex())
	}

	// General preset: AA predeploys must NOT be in genesis
	for _, addr := range aaAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "AA address should not be in General preset: %s", addr.Hex())
	}
}

func TestBuildL2GenesisGamingPreset(t *testing.T) {
	config, err := genesis.NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.Nil(t, err)
	config.Preset = genesis.PresetGaming
	config.VRFAdmin = common.HexToAddress("0x1234567890123456789012345678901234567890")
	config.AAPaymasterSigner = common.HexToAddress("0x0000000000000000000000000000000000000002")
	config.EnableGovernance = false
	config.FundDevAccounts = false

	gen := testBuildL2Genesis(t, config)

	// Gaming preset: VRF predeploys must be present
	for _, addr := range vrfAddrs {
		require.Contains(t, gen.Alloc, addr, "Gaming preset must have VRF predeploy: %s", addr.Hex())
	}

	// Gaming preset: AA predeploys must be present
	for _, addr := range aaAddrs {
		require.Contains(t, gen.Alloc, addr, "Gaming preset must have AA predeploy: %s", addr.Hex())
	}

	// Gaming preset: USDC Bridge predeploys must NOT be in genesis
	for _, addr := range usdcAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "USDC address should not be in Gaming preset: %s", addr.Hex())
	}

	// Gaming preset: Uniswap V3 predeploys must NOT be in genesis
	for _, addr := range uniswapAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "Uniswap address should not be in Gaming preset: %s", addr.Hex())
	}
}

func TestBuildL2GenesisDefiPreset(t *testing.T) {
	config, err := genesis.NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.Nil(t, err)
	config.Preset = genesis.PresetDeFi
	config.EnableGovernance = false
	config.FundDevAccounts = false

	gen := testBuildL2Genesis(t, config)

	// DeFi preset: USDC Bridge predeploys must be present
	for _, addr := range usdcAddrs {
		require.Contains(t, gen.Alloc, addr, "DeFi preset must have USDC predeploy: %s", addr.Hex())
	}

	// DeFi preset: Uniswap V3 predeploys must be present
	for _, addr := range uniswapAddrs {
		require.Contains(t, gen.Alloc, addr, "DeFi preset must have Uniswap predeploy: %s", addr.Hex())
	}

	// DeFi preset: VRF predeploys must NOT be in genesis
	for _, addr := range vrfAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "VRF address should not be in DeFi preset: %s", addr.Hex())
	}

	// DeFi preset: AA predeploys must NOT be in genesis
	for _, addr := range aaAddrs {
		_, ok := gen.Alloc[addr]
		require.False(t, ok, "AA address should not be in DeFi preset: %s", addr.Hex())
	}
}

func TestBuildL2GenesisFullPreset(t *testing.T) {
	config, err := genesis.NewDeployConfig("./testdata/test-deploy-config-devnet-l1.json")
	require.Nil(t, err)
	config.Preset = genesis.PresetFull
	config.VRFAdmin = common.HexToAddress("0x1234567890123456789012345678901234567890")
	config.AAPaymasterSigner = common.HexToAddress("0x0000000000000000000000000000000000000002")
	config.EnableGovernance = false
	config.FundDevAccounts = false

	gen := testBuildL2Genesis(t, config)

	// Full preset: VRF predeploys must be present
	for _, addr := range vrfAddrs {
		require.Contains(t, gen.Alloc, addr, "Full preset must have VRF predeploy: %s", addr.Hex())
	}

	// Full preset: AA predeploys must be present
	for _, addr := range aaAddrs {
		require.Contains(t, gen.Alloc, addr, "Full preset must have AA predeploy: %s", addr.Hex())
	}

	// Full preset: USDC Bridge predeploys must be present
	for _, addr := range usdcAddrs {
		require.Contains(t, gen.Alloc, addr, "Full preset must have USDC predeploy: %s", addr.Hex())
	}

	// Full preset: Uniswap V3 predeploys must be present
	for _, addr := range uniswapAddrs {
		require.Contains(t, gen.Alloc, addr, "Full preset must have Uniswap predeploy: %s", addr.Hex())
	}
}
