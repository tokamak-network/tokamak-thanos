package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis"
)

// LoadTokamakDeployment loads the Tokamak deployment addresses from the .deploy file
// and creates an L1Deployments structure
func LoadTokamakDeployment(root string) (*genesis.L1Deployments, error) {
	deployPath := path.Join(root, "packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy")

	data, err := os.ReadFile(deployPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read deployment file: %w", err)
	}

	var addresses map[string]string
	if err := json.Unmarshal(data, &addresses); err != nil {
		return nil, fmt.Errorf("failed to parse deployment file: %w", err)
	}

	return mapToL1Deployments(addresses), nil
}

// mapToL1Deployments converts the deployment addresses map to L1Deployments struct
func mapToL1Deployments(addresses map[string]string) *genesis.L1Deployments {
	deployments := &genesis.L1Deployments{}

	// Map each address to the corresponding field
	if addr, ok := addresses["AddressManager"]; ok {
		deployments.AddressManager = common.HexToAddress(addr)
	}
	if addr, ok := addresses["L1CrossDomainMessengerProxy"]; ok {
		deployments.L1CrossDomainMessengerProxy = common.HexToAddress(addr)
	}
	if addr, ok := addresses["L1StandardBridgeProxy"]; ok {
		deployments.L1StandardBridgeProxy = common.HexToAddress(addr)
	}
	if addr, ok := addresses["L2OutputOracleProxy"]; ok {
		deployments.L2OutputOracleProxy = common.HexToAddress(addr)
	}
	if addr, ok := addresses["OptimismPortalProxy"]; ok {
		deployments.OptimismPortalProxy = common.HexToAddress(addr)
	}
	if addr, ok := addresses["SystemConfigProxy"]; ok {
		deployments.SystemConfigProxy = common.HexToAddress(addr)
	}
	if addr, ok := addresses["ProxyAdmin"]; ok {
		deployments.ProxyAdmin = common.HexToAddress(addr)
	}
	// SuperchainConfigProxy is not in L1Deployments struct in this version
	// if addr, ok := addresses["SuperchainConfigProxy"]; ok {
	//     deployments.SuperchainConfigProxy = common.HexToAddress(addr)
	// }
	if addr, ok := addresses["L1ERC721BridgeProxy"]; ok {
		deployments.L1ERC721BridgeProxy = common.HexToAddress(addr)
	}
	if addr, ok := addresses["ProtocolVersionsProxy"]; ok {
		deployments.ProtocolVersionsProxy = common.HexToAddress(addr)
	}
	if addr, ok := addresses["OptimismMintableERC20FactoryProxy"]; ok {
		deployments.OptimismMintableERC20FactoryProxy = common.HexToAddress(addr)
	}

	// Dispute game related addresses
	if addr, ok := addresses["DisputeGameFactoryProxy"]; ok {
		deployments.DisputeGameFactoryProxy = common.HexToAddress(addr)
	}
	// These fields are not in L1Deployments struct in this version
	// if addr, ok := addresses["DelayedWETHProxy"]; ok {
	//     deployments.DelayedWETHProxy = common.HexToAddress(addr)
	// }
	// if addr, ok := addresses["PermissionedDelayedWETHProxy"]; ok {
	//     deployments.PermissionedDelayedWETHProxy = common.HexToAddress(addr)
	// }
	// if addr, ok := addresses["AnchorStateRegistryProxy"]; ok {
	//     deployments.AnchorStateRegistryProxy = common.HexToAddress(addr)
	// }
	// if addr, ok := addresses["PreimageOracle"]; ok {
	//     deployments.PreimageOracle = common.HexToAddress(addr)
	// }
	// if addr, ok := addresses["Mips"]; ok {
	//     deployments.Mips = common.HexToAddress(addr)
	// }

	// Tokamak specific contracts
	if _, ok := addresses["L2NativeToken"]; ok {
		// Add Tokamak-specific field if it exists in L1Deployments
		// For now, store it in a custom way if needed
	}
	if _, ok := addresses["L1UsdcBridgeProxy"]; ok {
		// Add USDC bridge if needed
	}

	return deployments
}

// CreateMinimalL1State creates minimal L1 state for testing without full state-dump
// This is useful when state-dump files are not available but contract addresses are known
func CreateMinimalL1State(deployments *genesis.L1Deployments) types.GenesisAlloc {
	allocs := make(types.GenesisAlloc)

	// Add minimal proxy code for each deployed contract
	// This is just placeholder code that allows the contracts to exist on-chain
	minimalProxyCode := common.FromHex("0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063f851a44014610030575b600080fd5b61003861004e565b6040516001600160a01b03909116815260200160405180910390f35b60006100767f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b905090565b600080fd5b5056fea26469706673582212208b0e7b7c7d89f3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c64736f6c63430008110033")

	// Helper function to add contract with minimal code
	addContract := func(addr common.Address) {
		if addr != (common.Address{}) {
			allocs[addr] = types.Account{
				Code:    minimalProxyCode,
				Balance: big.NewInt(0),
				Storage: make(map[common.Hash]common.Hash),
			}
		}
	}

	// Add all contract addresses
	addContract(deployments.AddressManager)
	addContract(deployments.L1CrossDomainMessengerProxy)
	addContract(deployments.L1StandardBridgeProxy)
	addContract(deployments.L2OutputOracleProxy)
	addContract(deployments.OptimismPortalProxy)
	addContract(deployments.SystemConfigProxy)
	addContract(deployments.ProxyAdmin)
	// addContract(deployments.SuperchainConfigProxy) // Field not in struct
	addContract(deployments.L1ERC721BridgeProxy)
	addContract(deployments.ProtocolVersionsProxy)
	addContract(deployments.OptimismMintableERC20FactoryProxy)
	addContract(deployments.DisputeGameFactoryProxy)
	// addContract(deployments.DelayedWETHProxy) // Field not in struct
	// addContract(deployments.PermissionedDelayedWETHProxy) // Field not in struct
	// addContract(deployments.AnchorStateRegistryProxy) // Field not in struct
	// addContract(deployments.PreimageOracle) // Field not in struct
	// addContract(deployments.Mips) // Field not in struct

	return allocs
}