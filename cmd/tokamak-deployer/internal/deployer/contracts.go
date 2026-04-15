package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type artifact struct {
	ABI      json.RawMessage `json:"abi"`
	Bytecode struct {
		Object string `json:"object"`
	} `json:"bytecode"`
}

func loadArtifact(artifactsFS fs.FS, name string) (*artifact, error) {
	data, err := fs.ReadFile(artifactsFS, "deploy-artifacts/"+name+".json")
	if err != nil {
		return nil, fmt.Errorf("artifact %s not found: %w", name, err)
	}
	var a artifact
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("invalid artifact %s: %w", name, err)
	}
	return &a, nil
}

func deployRawContract(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, nonce *uint64, bytecodeWithArgs []byte) (common.Address, error) {
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Address{}, err
	}
	tx := types.NewContractCreation(*nonce, common.Big0, 5_000_000, gasPrice, bytecodeWithArgs)
	*nonce++
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("sign tx: %w", err)
	}
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return common.Address{}, fmt.Errorf("send tx: %w", err)
	}
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return common.Address{}, fmt.Errorf("wait mined: %w", err)
	}
	if receipt.Status == 0 {
		return common.Address{}, fmt.Errorf("deployment reverted")
	}
	return receipt.ContractAddress, nil
}

func deployContract(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, nonce *uint64, a *artifact, constructorArgs ...interface{}) (common.Address, error) {
	parsedABI, err := abi.JSON(strings.NewReader(string(a.ABI)))
	if err != nil {
		return common.Address{}, fmt.Errorf("parse ABI: %w", err)
	}
	bytecode := common.FromHex(a.Bytecode.Object)
	if len(constructorArgs) > 0 {
		packed, err := parsedABI.Pack("", constructorArgs...)
		if err != nil {
			return common.Address{}, fmt.Errorf("pack constructor args: %w", err)
		}
		bytecode = append(bytecode, packed...)
	}
	return deployRawContract(ctx, client, auth, nonce, bytecode)
}

func callContract(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, nonce *uint64, to common.Address, a *artifact, method string, args ...interface{}) error {
	parsedABI, err := abi.JSON(strings.NewReader(string(a.ABI)))
	if err != nil {
		return fmt.Errorf("parse ABI for %s: %w", method, err)
	}
	data, err := parsedABI.Pack(method, args...)
	if err != nil {
		return fmt.Errorf("pack %s: %w", method, err)
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("get gas price in callContract: %w", err)
	}
	tx := types.NewTransaction(*nonce, to, common.Big0, 3_000_000, gasPrice, data)
	*nonce++
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return fmt.Errorf("sign tx: %w", err)
	}
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("send %s tx: %w", method, err)
	}
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return fmt.Errorf("wait mined in callContract: %w", err)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("%s call reverted", method)
	}
	return nil
}

func setProxyType(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, nonce *uint64, proxyAdminAddr common.Address, proxyAddr common.Address, proxyAdminArtifact *artifact) error {
	// ProxyType.ERC1967 = 0
	return callContract(ctx, client, auth, nonce, proxyAdminAddr, proxyAdminArtifact, "setProxyType", proxyAddr, uint8(0))
}

func upgradeProxyViaAdmin(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, nonce *uint64, proxyAdminAddr common.Address, proxyAddr common.Address, implAddr common.Address, proxyAdminArtifact *artifact) error {
	// Try upgrade() first (simpler, no initialization)
	return callContract(ctx, client, auth, nonce, proxyAdminAddr, proxyAdminArtifact, "upgrade", proxyAddr, implAddr)
}

func Deploy(ctx context.Context, cfg DeployConfig, artifactsFS fs.FS) (*DeployOutput, error) {
	client, err := ethclient.DialContext(ctx, cfg.L1RPCURL)
	if err != nil {
		return nil, fmt.Errorf("connect L1: %w", err)
	}
	defer client.Close()

	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	nonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	output := &DeployOutput{L2ChainID: cfg.L2ChainID, L1ChainID: chainID.Uint64()}

	// Load artifacts
	proxyArtifact, err := loadArtifact(artifactsFS, "Proxy")
	if err != nil {
		return nil, err
	}

	// 1. Deploy AddressManager
	addressManagerArtifact, err := loadArtifact(artifactsFS, "AddressManager")
	if err != nil {
		return nil, err
	}
	addressManagerAddr, err := deployContract(ctx, client, auth, &nonce, addressManagerArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy AddressManager: %w", err)
	}
	output.AddressManager = addressManagerAddr.Hex()

	// 2. Deploy ProxyAdmin
	proxyAdminArtifact, err := loadArtifact(artifactsFS, "ProxyAdmin")
	if err != nil {
		return nil, err
	}
	proxyAdminAddr, err := deployContract(ctx, client, auth, &nonce, proxyAdminArtifact, auth.From)
	if err != nil {
		return nil, fmt.Errorf("deploy ProxyAdmin: %w", err)
	}
	output.ProxyAdmin = proxyAdminAddr.Hex()

	// 3. Deploy SuperchainConfigProxy
	superchainConfigProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy SuperchainConfigProxy: %w", err)
	}
	output.SuperchainConfigProxy = superchainConfigProxyAddr.Hex()

	// 4. Deploy SuperchainConfig implementation
	superchainConfigArtifact, err := loadArtifact(artifactsFS, "SuperchainConfig")
	if err != nil {
		return nil, err
	}
	superchainConfigImplAddr, err := deployContract(ctx, client, auth, &nonce, superchainConfigArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy SuperchainConfig impl: %w", err)
	}

	// 5. Upgrade SuperchainConfigProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, superchainConfigProxyAddr, superchainConfigImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade SuperchainConfigProxy: %w", err)
	}

	// 6. Deploy OptimismPortalProxy
	optimismPortalProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy OptimismPortalProxy: %w", err)
	}
	output.OptimismPortalProxy = optimismPortalProxyAddr.Hex()

	// 7. Deploy OptimismPortal implementation
	optimismPortalArtifact, err := loadArtifact(artifactsFS, "OptimismPortal")
	if err != nil {
		return nil, err
	}
	optimismPortalImplAddr, err := deployContract(ctx, client, auth, &nonce, optimismPortalArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy OptimismPortal impl: %w", err)
	}

	// 8. Upgrade OptimismPortalProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, optimismPortalProxyAddr, optimismPortalImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade OptimismPortalProxy: %w", err)
	}

	// 9. Deploy SystemConfigProxy
	systemConfigProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy SystemConfigProxy: %w", err)
	}
	output.SystemConfigProxy = systemConfigProxyAddr.Hex()

	// 10. Deploy SystemConfig implementation
	systemConfigArtifact, err := loadArtifact(artifactsFS, "SystemConfig")
	if err != nil {
		return nil, err
	}
	systemConfigImplAddr, err := deployContract(ctx, client, auth, &nonce, systemConfigArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy SystemConfig impl: %w", err)
	}

	// 11. Upgrade SystemConfigProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, systemConfigProxyAddr, systemConfigImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade SystemConfigProxy: %w", err)
	}

	// 12. Deploy L1StandardBridgeProxy
	l1StandardBridgeProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy L1StandardBridgeProxy: %w", err)
	}
	output.L1StandardBridgeProxy = l1StandardBridgeProxyAddr.Hex()

	// 13. Deploy L1StandardBridge implementation
	l1StandardBridgeArtifact, err := loadArtifact(artifactsFS, "L1StandardBridge")
	if err != nil {
		return nil, err
	}
	l1StandardBridgeImplAddr, err := deployContract(ctx, client, auth, &nonce, l1StandardBridgeArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy L1StandardBridge impl: %w", err)
	}

	// 14. Upgrade L1StandardBridgeProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, l1StandardBridgeProxyAddr, l1StandardBridgeImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade L1StandardBridgeProxy: %w", err)
	}

	// 15. Deploy L1CrossDomainMessengerProxy
	l1CDMProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy L1CrossDomainMessengerProxy: %w", err)
	}
	output.L1CrossDomainMessengerProxy = l1CDMProxyAddr.Hex()

	// 16. Deploy L1CrossDomainMessenger implementation
	l1CDMArtifact, err := loadArtifact(artifactsFS, "L1CrossDomainMessenger")
	if err != nil {
		return nil, err
	}
	l1CDMImplAddr, err := deployContract(ctx, client, auth, &nonce, l1CDMArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy L1CrossDomainMessenger impl: %w", err)
	}

	// 17. Upgrade L1CrossDomainMessengerProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, l1CDMProxyAddr, l1CDMImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade L1CrossDomainMessengerProxy: %w", err)
	}

	// 18. Deploy OptimismMintableERC20FactoryProxy
	optimismMintableERC20FactoryProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy OptimismMintableERC20FactoryProxy: %w", err)
	}
	output.OptimismMintableERC20FactoryProxy = optimismMintableERC20FactoryProxyAddr.Hex()

	// 19. Deploy OptimismMintableERC20Factory implementation
	optimismMintableERC20FactoryArtifact, err := loadArtifact(artifactsFS, "OptimismMintableERC20Factory")
	if err != nil {
		return nil, err
	}
	optimismMintableERC20FactoryImplAddr, err := deployContract(ctx, client, auth, &nonce, optimismMintableERC20FactoryArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy OptimismMintableERC20Factory impl: %w", err)
	}

	// 20. Upgrade OptimismMintableERC20FactoryProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, optimismMintableERC20FactoryProxyAddr, optimismMintableERC20FactoryImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade OptimismMintableERC20FactoryProxy: %w", err)
	}

	// 21. Deploy L1ERC721BridgeProxy
	l1ERC721BridgeProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy L1ERC721BridgeProxy: %w", err)
	}
	output.L1ERC721BridgeProxy = l1ERC721BridgeProxyAddr.Hex()

	// 22. Deploy L1ERC721Bridge implementation
	l1ERC721BridgeArtifact, err := loadArtifact(artifactsFS, "L1ERC721Bridge")
	if err != nil {
		return nil, err
	}
	l1ERC721BridgeImplAddr, err := deployContract(ctx, client, auth, &nonce, l1ERC721BridgeArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy L1ERC721Bridge impl: %w", err)
	}

	// 23. Upgrade L1ERC721BridgeProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, l1ERC721BridgeProxyAddr, l1ERC721BridgeImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade L1ERC721BridgeProxy: %w", err)
	}

	// 24. Deploy L2OutputOracleProxy
	l2OutputOracleProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("deploy L2OutputOracleProxy: %w", err)
	}
	output.L2OutputOracleProxy = l2OutputOracleProxyAddr.Hex()

	// 25. Deploy L2OutputOracle implementation
	l2OutputOracleArtifact, err := loadArtifact(artifactsFS, "L2OutputOracle")
	if err != nil {
		return nil, err
	}
	l2OutputOracleImplAddr, err := deployContract(ctx, client, auth, &nonce, l2OutputOracleArtifact)
	if err != nil {
		return nil, fmt.Errorf("deploy L2OutputOracle impl: %w", err)
	}

	// 26. Upgrade L2OutputOracleProxy to implementation
	if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, l2OutputOracleProxyAddr, l2OutputOracleImplAddr, proxyAdminArtifact); err != nil {
		return nil, fmt.Errorf("upgrade L2OutputOracleProxy: %w", err)
	}

	// 27-33. Deploy DisputeGameFactory and AnchorStateRegistry only if EnableFaultProof is set
	if cfg.EnableFaultProof {
		// 27. Deploy DisputeGameFactoryProxy
		disputeGameFactoryProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
		if err != nil {
			return nil, fmt.Errorf("deploy DisputeGameFactoryProxy: %w", err)
		}
		output.DisputeGameFactoryProxy = disputeGameFactoryProxyAddr.Hex()

		// 28. Deploy DisputeGameFactory implementation
		disputeGameFactoryArtifact, err := loadArtifact(artifactsFS, "DisputeGameFactory")
		if err != nil {
			return nil, err
		}
		disputeGameFactoryImplAddr, err := deployContract(ctx, client, auth, &nonce, disputeGameFactoryArtifact)
		if err != nil {
			return nil, fmt.Errorf("deploy DisputeGameFactory impl: %w", err)
		}

		// 29. Upgrade DisputeGameFactoryProxy to implementation
		if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, disputeGameFactoryProxyAddr, disputeGameFactoryImplAddr, proxyAdminArtifact); err != nil {
			return nil, fmt.Errorf("upgrade DisputeGameFactoryProxy: %w", err)
		}

		// 30. Deploy AnchorStateRegistryProxy
		anchorStateRegistryProxyAddr, err := deployContract(ctx, client, auth, &nonce, proxyArtifact, proxyAdminAddr)
		if err != nil {
			return nil, fmt.Errorf("deploy AnchorStateRegistryProxy: %w", err)
		}
		output.AnchorStateRegistryProxy = anchorStateRegistryProxyAddr.Hex()

		// 31. Deploy AnchorStateRegistry implementation (constructor takes DisputeGameFactory address)
		anchorStateRegistryArtifact, err := loadArtifact(artifactsFS, "AnchorStateRegistry")
		if err != nil {
			return nil, err
		}
		anchorStateRegistryImplAddr, err := deployContract(ctx, client, auth, &nonce, anchorStateRegistryArtifact, disputeGameFactoryImplAddr)
		if err != nil {
			return nil, fmt.Errorf("deploy AnchorStateRegistry impl: %w", err)
		}

		// 32. Upgrade AnchorStateRegistryProxy to implementation
		if err := upgradeProxyViaAdmin(ctx, client, auth, &nonce, proxyAdminAddr, anchorStateRegistryProxyAddr, anchorStateRegistryImplAddr, proxyAdminArtifact); err != nil {
			return nil, fmt.Errorf("upgrade AnchorStateRegistryProxy: %w", err)
		}
	}

	return output, nil
}
