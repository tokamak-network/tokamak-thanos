// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "../L1ContractVerification.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

// Mock contracts for testing
contract MockSystemConfig {
    address public l1StandardBridge;
    address public l1CrossDomainMessenger;
    address public optimismPortal;

    constructor(address _l1StandardBridge, address _l1CrossDomainMessenger, address _optimismPortal) {
        l1StandardBridge = _l1StandardBridge;
        l1CrossDomainMessenger = _l1CrossDomainMessenger;
        optimismPortal = _optimismPortal;
    }
}

contract MockL1StandardBridge {
    // Empty implementation
}

contract MockL1CrossDomainMessenger {
    // Empty implementation
}

contract MockOptimismPortal {
    // Empty implementation
}

contract MockGnosisSafe {
    address[] private owners;
    uint256 private threshold;

    constructor(address[] memory _owners, uint256 _threshold) {
        owners = _owners;
        threshold = _threshold;
    }

    function getOwners() external view returns (address[] memory) {
        return owners;
    }

    function getThreshold() external view returns (uint256) {
        return threshold;
    }
}

contract MockBridgeRegistry {
    event RollupConfigRegistered(address rollupConfig, uint8 tokenType, address l2Token, string name);

    function registerRollupConfig(address rollupConfig, uint8 tokenType, address l2Token, string calldata name) external {
        emit RollupConfigRegistered(rollupConfig, tokenType, l2Token, name);
    }
}

contract L1ContractVerificationTest is Test {
    // Constants
    uint256 constant CHAIN_ID = 11155111; // Sepolia
    bytes32 constant L1_STANDARD_BRIDGE_ID = keccak256("L1_STANDARD_BRIDGE");
    bytes32 constant L1_CROSS_DOMAIN_MESSENGER_ID = keccak256("L1_CROSS_DOMAIN_MESSENGER");
    bytes32 constant OPTIMISM_PORTAL_ID = keccak256("OPTIMISM_PORTAL");
    bytes32 constant SYSTEM_CONFIG_ID = keccak256("SYSTEM_CONFIG");

    // Contracts
    L1ContractVerification verifier;
    ProxyAdmin proxyAdmin;

    // Mock implementations
    MockSystemConfig systemConfigImpl;
    MockL1StandardBridge l1StandardBridgeImpl;
    MockL1CrossDomainMessenger l1CrossDomainMessengerImpl;
    MockOptimismPortal optimismPortalImpl;

    // Proxies
    TransparentUpgradeableProxy systemConfigProxy;
    TransparentUpgradeableProxy l1StandardBridgeProxy;
    TransparentUpgradeableProxy l1CrossDomainMessengerProxy;
    TransparentUpgradeableProxy optimismPortalProxy;

    // Safe
    MockGnosisSafe safeWallet;

    // Bridge registry
    MockBridgeRegistry bridgeRegistry;

    // Addresses
    address owner;
    address tokamakDAO;
    address foundation;
    address user;

    function setUp() public {
        // Setup accounts
        owner = makeAddr("owner");
        tokamakDAO = makeAddr("tokamakDAO");
        foundation = makeAddr("foundation");
        user = makeAddr("user");

        vm.startPrank(owner);

        // Deploy proxy admin
        proxyAdmin = new ProxyAdmin();

        // Deploy implementations
        l1StandardBridgeImpl = new MockL1StandardBridge();
        l1CrossDomainMessengerImpl = new MockL1CrossDomainMessenger();
        optimismPortalImpl = new MockOptimismPortal();

        // Deploy proxies (without initializing SystemConfig yet)
        l1StandardBridgeProxy = new TransparentUpgradeableProxy(
            address(l1StandardBridgeImpl),
            address(proxyAdmin),
            ""
        );

        l1CrossDomainMessengerProxy = new TransparentUpgradeableProxy(
            address(l1CrossDomainMessengerImpl),
            address(proxyAdmin),
            ""
        );

        optimismPortalProxy = new TransparentUpgradeableProxy(
            address(optimismPortalImpl),
            address(proxyAdmin),
            ""
        );

        // Now deploy SystemConfig implementation with the proxy addresses
        systemConfigImpl = new MockSystemConfig(
            address(l1StandardBridgeProxy),
            address(l1CrossDomainMessengerProxy),
            address(optimismPortalProxy)
        );

        // Deploy SystemConfig proxy
        systemConfigProxy = new TransparentUpgradeableProxy(
            address(systemConfigImpl),
            address(proxyAdmin),
            ""
        );

        // Setup Safe wallet
        address[] memory safeOwners = new address[](2);
        safeOwners[0] = tokamakDAO;
        safeOwners[1] = foundation;
        safeWallet = new MockGnosisSafe(safeOwners, 2);

        // Deploy bridge registry
        bridgeRegistry = new MockBridgeRegistry();

        // Deploy verifier contract
        verifier = new L1ContractVerification();

        // Set bridge registry address
        verifier.setBridgeRegistryAddress(address(bridgeRegistry));

        vm.stopPrank();
    }

    function testSetContractConfig() public {
        vm.startPrank(owner);

        // Set contract config for SystemConfig
        verifier.setContractConfig(
            CHAIN_ID,
            SYSTEM_CONFIG_ID,
            address(systemConfigImpl).codehash,
            address(systemConfigProxy).codehash,
            address(proxyAdmin)
        );

        // Verify the config was set correctly
        (bytes32 implementationHash, bytes32 proxyHash, address expectedProxyAdmin) = verifier.contractConfigs(CHAIN_ID, SYSTEM_CONFIG_ID);

        assertEq(implementationHash, address(systemConfigImpl).codehash);
        assertEq(proxyHash, address(systemConfigProxy).codehash);
        assertEq(expectedProxyAdmin, address(proxyAdmin));

        vm.stopPrank();
    }

    function testSetSafeConfig() public {
        vm.startPrank(owner);

        // Set safe config
        verifier.setSafeConfig(CHAIN_ID, tokamakDAO, foundation, 2);

        // Verify the config was set correctly
        (address configTokamakDAO, address configFoundation, uint256 configThreshold) = verifier.safeConfigs(CHAIN_ID);

        assertEq(configTokamakDAO, tokamakDAO);
        assertEq(configFoundation, foundation);
        assertEq(configThreshold, 2);

        vm.stopPrank();
    }

    function testVerifyL1ContractsSuccess() public {
        vm.startPrank(owner);

        // Set all required configs
        _setupAllConfigs();

        vm.stopPrank();

        vm.startPrank(user);

        // Verify L1 contracts
        bool result = verifier.verifyL1Contracts(
            CHAIN_ID,
            address(systemConfigProxy),
            address(safeWallet)
        );

        // Verification should succeed
        assertTrue(result);

        vm.stopPrank();
    }

    function testVerifyL1ContractsFailInvalidSystemConfig() public {
        vm.startPrank(owner);

        // Set all required configs
        _setupAllConfigs();

        // Deploy a different SystemConfig implementation
        MockSystemConfig differentSystemConfigImpl = new MockSystemConfig(
            address(0),
            address(0),
            address(0)
        );

        // Update the SystemConfig proxy to point to the different implementation
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(systemConfigProxy))),
            address(differentSystemConfigImpl)
        );

        vm.stopPrank();

        vm.startPrank(user);

        // Verification should fail
        vm.expectRevert("Invalid SystemConfig");
        verifier.verifyL1Contracts(
            CHAIN_ID,
            address(systemConfigProxy),
            address(safeWallet)
        );

        vm.stopPrank();
    }

    function testVerifyL1ContractsFailInvalidSafe() public {
        vm.startPrank(owner);

        // Set all required configs
        _setupAllConfigs();

        // Create a different safe with wrong threshold
        address[] memory safeOwners = new address[](2);
        safeOwners[0] = tokamakDAO;
        safeOwners[1] = foundation;
        MockGnosisSafe differentSafe = new MockGnosisSafe(safeOwners, 1); // Wrong threshold

        vm.stopPrank();

        vm.startPrank(user);

        // Verification should fail
        vm.expectRevert("Invalid Safe configuration");
        verifier.verifyL1Contracts(
            CHAIN_ID,
            address(systemConfigProxy),
            address(differentSafe)
        );

        vm.stopPrank();
    }

    function testVerifyAndRegisterRollupConfig() public {
        vm.startPrank(owner);

        // Set all required configs
        _setupAllConfigs();

        vm.stopPrank();

        vm.startPrank(user);

        address rollupConfig = makeAddr("rollupConfig");
        address l2TON = makeAddr("l2TON");

        // Expect the bridge registry to be called
        vm.expectEmit(true, true, true, true);
        emit MockBridgeRegistry.RollupConfigRegistered(rollupConfig, 2, l2TON, "TestRollup");

        // Verify and register rollup config
        bool result = verifier.verifyAndRegisterRollupConfig(
            CHAIN_ID,
            address(systemConfigProxy),
            address(safeWallet),
            rollupConfig,
            2, // TON token type
            l2TON,
            "TestRollup"
        );

        // Verification and registration should succeed
        assertTrue(result);

        vm.stopPrank();
    }

    function testVerifyAndRegisterRollupConfigFailInvalidType() public {
        vm.startPrank(owner);

        // Set all required configs
        _setupAllConfigs();

        vm.stopPrank();

        vm.startPrank(user);

        address rollupConfig = makeAddr("rollupConfig");
        address l2TON = makeAddr("l2TON");

        // Verification should fail due to invalid token type
        vm.expectRevert("Registration allowed only for TON tokens");
        verifier.verifyAndRegisterRollupConfig(
            CHAIN_ID,
            address(systemConfigProxy),
            address(safeWallet),
            rollupConfig,
            1, // Invalid token type
            l2TON,
            "TestRollup"
        );

        vm.stopPrank();
    }

    function testOnlyOwnerFunctions() public {
        vm.startPrank(user);

        // All these should revert because user is not the owner
        vm.expectRevert("Ownable: caller is not the owner");
        verifier.setContractConfig(
            CHAIN_ID,
            SYSTEM_CONFIG_ID,
            bytes32(0),
            bytes32(0),
            address(0)
        );

        vm.expectRevert("Ownable: caller is not the owner");
        verifier.setSafeConfig(
            CHAIN_ID,
            address(0),
            address(0),
            0
        );

        vm.expectRevert("Ownable: caller is not the owner");
        verifier.setBridgeRegistryAddress(address(0));

        vm.stopPrank();
    }

    function testSetBridgeRegistryAddress() public {
        vm.startPrank(owner);

        address newBridgeRegistry = makeAddr("newBridgeRegistry");
        verifier.setBridgeRegistryAddress(newBridgeRegistry);

        assertEq(verifier.L1BridgeRegistryV1_1Address(), newBridgeRegistry);

        vm.stopPrank();
    }

    function testSetBridgeRegistryAddressFailZeroAddress() public {
        vm.startPrank(owner);

        vm.expectRevert("Invalid bridge registry address");
        verifier.setBridgeRegistryAddress(address(0));

        vm.stopPrank();
    }

    // Helper function to set up all required configurations
    function _setupAllConfigs() internal {
        // Set contract config for SystemConfig
        verifier.setContractConfig(
            CHAIN_ID,
            SYSTEM_CONFIG_ID,
            address(systemConfigImpl).codehash,
            address(systemConfigProxy).codehash,
            address(proxyAdmin)
        );

        // Set contract config for L1StandardBridge
        verifier.setContractConfig(
            CHAIN_ID,
            L1_STANDARD_BRIDGE_ID,
            address(l1StandardBridgeImpl).codehash,
            address(l1StandardBridgeProxy).codehash,
            address(proxyAdmin)
        );

        // Set contract config for L1CrossDomainMessenger
        verifier.setContractConfig(
            CHAIN_ID,
            L1_CROSS_DOMAIN_MESSENGER_ID,
            address(l1CrossDomainMessengerImpl).codehash,
            address(l1CrossDomainMessengerProxy).codehash,
            address(proxyAdmin)
        );

        // Set contract config for OptimismPortal
        verifier.setContractConfig(
            CHAIN_ID,
            OPTIMISM_PORTAL_ID,
            address(optimismPortalImpl).codehash,
            address(optimismPortalProxy).codehash,
            address(proxyAdmin)
        );

        // Set safe config
        verifier.setSafeConfig(CHAIN_ID, tokamakDAO, foundation, 2);
    }
}