// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import 'forge-std/Test.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import '@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol';
import {L1ChugSplashProxy} from 'src/legacy/L1ChugSplashProxy.sol';
import './mock-contracts/LegacyProxy.sol';
import './mock-contracts/MockContracts.sol';
import './mock-contracts/MockGnosisSafe.sol';
import './mock-contracts/MockProxyAdmin.sol';
import './mock-contracts/MockSystemConfig.sol';
import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts-upgradeable/utils/StringsUpgradeable.sol';
import {IProxyAdmin} from 'src/tokamak-contracts/verification/interface/IProxyAdmin.sol';

/**
 * @title L1ContractVerificationTest
 * @notice Test contract for L1ContractVerification functionality
 * @dev Tests various verification scenarios including success and failure cases
 */
contract L1ContractVerificationTest is Test {
  // Constants
  bytes32 constant L1_STANDARD_BRIDGE_ID = keccak256('L1_STANDARD_BRIDGE');
  bytes32 constant L1_CROSS_DOMAIN_MESSENGER_ID =
    keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 constant OPTIMISM_PORTAL_ID = keccak256('OPTIMISM_PORTAL');
  bytes32 constant SYSTEM_CONFIG_ID = keccak256('SYSTEM_CONFIG');
  bytes32 constant IMPLEMENTATION_SLOT =
    0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

  // Contracts
  L1ContractVerification verifierImpl;
  L1ContractVerification verifier;
  ProxyAdmin verifierProxyAdmin;
  TransparentUpgradeableProxy verifierProxy;

  MockProxyAdmin mockProxyAdmin;

  // Mock implementations
  MockSystemConfig systemConfigImpl;
  MockL1StandardBridge l1StandardBridgeImpl;
  MockL1CrossDomainMessenger l1CrossDomainMessengerImpl;
  MockOptimismPortal optimismPortalImpl;

  // Proxies - with correct types
  TransparentUpgradeableProxy systemConfigProxy; // Transparent proxy
  L1ChugSplashProxy l1StandardBridgeProxy; // Chug splash proxy
  LegacyProxy l1CrossDomainMessengerProxy; // Legacy proxy
  TransparentUpgradeableProxy optimismPortalProxy; // Transparent proxy

  // Safe
  MockGnosisSafe safeWallet;

  // Bridge registry
  MockBridgeRegistry bridgeRegistry;

  // Addresses
  address owner;
  address tokamakDAO;
  address foundation;
  address thirdOwner;
  address user;
  address nativeToken;
  address l2TONAddress;

    // Add the randomUser address:
  address randomUser;

  /**
   * @notice Set up the test environment
   * @dev Deploys all required contracts and configures them for testing
   */
  function setUp() public {
    // Setup accounts
    owner = makeAddr('owner');
    tokamakDAO = makeAddr('tokamakDAO');
    foundation = makeAddr('foundation');
    thirdOwner = makeAddr('thirdOwner');
    user = makeAddr('user');
    nativeToken = makeAddr('nativeToken');
    l2TONAddress = address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000); // Use the expected L2 TON address

        // Initialize the randomUser
    randomUser = makeAddr('randomUser');

    vm.startPrank(owner);

    // Deploy proxy admin
    mockProxyAdmin = new MockProxyAdmin(owner);

    // Deploy implementations
    l1StandardBridgeImpl = new MockL1StandardBridge();
    l1CrossDomainMessengerImpl = new MockL1CrossDomainMessenger();
    optimismPortalImpl = new MockOptimismPortal();

    // Deploy proxies (without initializing SystemConfig yet)
    l1StandardBridgeProxy = new L1ChugSplashProxy(
      address(l1StandardBridgeImpl)
    );

    l1CrossDomainMessengerProxy = new LegacyProxy(
      address(l1CrossDomainMessengerImpl),
      address(mockProxyAdmin)
    );

    optimismPortalProxy = new TransparentUpgradeableProxy(
      address(optimismPortalImpl),
      address(mockProxyAdmin),
      ''
    );

    // Now deploy SystemConfig implementation with the proxy addresses
    systemConfigImpl = new MockSystemConfig();

    systemConfigProxy = new TransparentUpgradeableProxy(
      address(systemConfigImpl),
      address(mockProxyAdmin),
      ''
    );

    // Set up the mock ProxyAdmin
    mockProxyAdmin.setImplementation(
      address(systemConfigProxy),
      address(systemConfigImpl)
    );
    mockProxyAdmin.setImplementation(
      address(l1StandardBridgeProxy),
      address(l1StandardBridgeImpl)
    );
    mockProxyAdmin.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(l1CrossDomainMessengerImpl)
    );
    mockProxyAdmin.setImplementation(
      address(optimismPortalProxy),
      address(optimismPortalImpl)
    );

    mockProxyAdmin.setAdmin(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );
    mockProxyAdmin.setAdmin(
      address(l1StandardBridgeProxy),
      address(mockProxyAdmin)
    );
    mockProxyAdmin.setAdmin(
      address(l1CrossDomainMessengerProxy),
      address(mockProxyAdmin)
    );
    mockProxyAdmin.setAdmin(
      address(optimismPortalProxy),
      address(mockProxyAdmin)
    );

    // Initialize the SystemConfig proxy with the correct addresses
    bytes memory initData = abi.encodeWithSelector(
      MockSystemConfig.initialize.selector,
      address(l1StandardBridgeProxy),
      address(l1CrossDomainMessengerProxy),
      address(optimismPortalProxy),
      nativeToken
    );

    // Call initialize through the proxy
    (bool success, ) = address(systemConfigProxy).call(initData);
    require(success, 'SystemConfig initialization failed');

    // Verify the addresses were set correctly
    MockSystemConfig config = MockSystemConfig(address(systemConfigProxy));
    require(
      config.l1StandardBridge() == address(l1StandardBridgeProxy),
      'L1StandardBridge not set correctly'
    );
    require(
      config.l1CrossDomainMessenger() == address(l1CrossDomainMessengerProxy),
      'L1CrossDomainMessenger not set correctly'
    );
    require(
      config.optimismPortal() == address(optimismPortalProxy),
      'OptimismPortal not set correctly'
    );
    require(
      config.nativeTokenAddress() == nativeToken,
      'NativeToken not set correctly'
    );

    // Setup Safe wallet
    address[] memory safeOwners = new address[](3);
    safeOwners[0] = tokamakDAO;
    safeOwners[1] = foundation;
    safeOwners[2] = thirdOwner;
    safeWallet = new MockGnosisSafe(safeOwners, 3);

    // Deploy bridge registry
    bridgeRegistry = new MockBridgeRegistry();

    // Deploy verifier contract as upgradeable
    // First the implementation
    verifierImpl = new L1ContractVerification();

    // The proxy admin to manage the proxy
    verifierProxyAdmin = new ProxyAdmin();

    // The initialization data
    bytes memory verifierInitData = abi.encodeWithSelector(
      L1ContractVerification.initialize.selector,
      nativeToken,
      owner
    );

    // Deploy the proxy with the implementation
    verifierProxy = new TransparentUpgradeableProxy(
      address(verifierImpl),
      address(verifierProxyAdmin),
      verifierInitData
    );

    // Create a reference to interact with the proxy
    verifier = L1ContractVerification(address(verifierProxy));

    // Set bridge registry address
    verifier.setBridgeRegistryAddress(address(bridgeRegistry));

    // Set the verification possible to true for tests
    verifier.setVerificationPossible(true);

    // Set the owner of the proxy admin to the safe wallet for verification
    mockProxyAdmin.setOwner(address(safeWallet));

    vm.stopPrank();
  }

  // Sucess tests

  /**
   * @notice Test successful verification of L1 contracts
   * @dev Sets up contract info and safe wallet info, then verifies all contracts
   */
  function testVerifyL1ContractsSuccess() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Get implementation address via masterCopy function
    address implementation = safeWallet.masterCopy();

    // Correctly set the safe wallet info
    verifier.setSafeConfig(
      tokamakDAO, // _tokamakDAO
      foundation, // _foundation
      3, // _threshold
      implementation.codehash,
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user so they can verify
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should succeed
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      'Test Rollup',
      address(safeWallet)
    );

    vm.stopPrank();
  }

  function testFallbackHandler() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Get implementation address via masterCopy function
    address implementation = safeWallet.masterCopy();

    // Correctly set the safe wallet info
    verifier.setSafeConfig(
      tokamakDAO, // _tokamakDAO
      foundation, // _foundation
      3, // _threshold
      implementation.codehash,
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user so they can verify
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should succeed
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    // Verification should succeed

    address fallbackHandler = safeWallet.getFallbackHandler();
    assertEq(fallbackHandler, address(0));

    fallbackHandler = address(0x123);
    safeWallet.setFallbackHandler(fallbackHandler);

    // Verification should fail
    vm.expectRevert(L1ContractVerification.SafeWalletInvalidFallbackHandler.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );
    fallbackHandler = safeWallet.getFallbackHandler();
    assertEq(fallbackHandler, address(0x123));

    // change it back to 0 and should succeed
    safeWallet.setFallbackHandler(address(0));
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test successful verification and registration of rollup configuration
   * @dev Verifies L1 contracts and registers a rollup configuration
   */
  function testVerifyAndRegisterRollupConfig() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Get implementation address via masterCopy function
    address implementation = safeWallet.masterCopy();

    // Correctly set the safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementation.codehash,
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user so they can verify and register
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification and registration should succeed
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      'TestRollup',
      address(safeWallet) // Pass safe wallet address
    );

    vm.stopPrank();
  }

  /**
   * @notice Test setting bridge registry address
   * @dev Sets a new bridge registry address and verifies it was updated
   */
  function testSetBridgeRegistryAddress() public {
    vm.startPrank(owner);

    address newBridgeRegistry = makeAddr('newBridgeRegistry');
    verifier.setBridgeRegistryAddress(newBridgeRegistry);

    assertEq(verifier.l1BridgeRegistryAddress(), newBridgeRegistry);

    vm.stopPrank();
  }

  // Failure tests

  /**
   * @dev ProxyAdmin verification Failures
   */

  /**
   * @notice Test verification failure with invalid ProxyAdmin
   * @dev Attempts verification with a contract that's not the expected ProxyAdmin
   */
  function testVerifyL1ContractsFailInvalidProxyAdminContract() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO, // _tokamakDAO
      foundation, // _foundation
      3, // _threshold
      implementationCodehash, // _implementationCodehash
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Passing a different contract so that verification fails
    MockL1StandardBridge wrongProxyAdmin = new MockL1StandardBridge();

    // Verification should fail with ProxyAdmin verification error
    vm.expectRevert(L1ContractVerification.ProxyAdminInvalidCodehash.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(wrongProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with incorrect proxy admin address
   * @dev Sets up a different proxy admin than what's expected in the safe wallet config
   */
  function testVerifyL1ContractsFailInvalidProxyAdminAddress() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info with the current safe wallet address
    verifier.setSafeConfig(
      tokamakDAO, // _tokamakDAO
      foundation, // _foundation
      3, // _threshold
      implementationCodehash, // _implementationCodehash
      proxyCodehash
    );

    MockProxyAdmin differentProxyAdmin = new MockProxyAdmin(owner);
    // Set the owner of the different proxy admin to the safe wallet
    differentProxyAdmin.setOwner(address(mockProxyAdmin));

    // Set implementations for the different proxy admin
    differentProxyAdmin.setImplementation(
      address(systemConfigProxy),
      address(systemConfigImpl)
    );
    differentProxyAdmin.setImplementation(
      address(l1StandardBridgeProxy),
      address(l1StandardBridgeImpl)
    );
    differentProxyAdmin.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(l1CrossDomainMessengerImpl)
    );
    differentProxyAdmin.setImplementation(
      address(optimismPortalProxy),
      address(optimismPortalImpl)
    );

    // Set admin relationships for each proxy
    differentProxyAdmin.setAdmin(
      address(systemConfigProxy),
      address(differentProxyAdmin)
    );
    differentProxyAdmin.setAdmin(
      address(l1StandardBridgeProxy),
      address(differentProxyAdmin)
    );
    differentProxyAdmin.setAdmin(
      address(l1CrossDomainMessengerProxy),
      address(differentProxyAdmin)
    );
    differentProxyAdmin.setAdmin(
      address(optimismPortalProxy),
      address(differentProxyAdmin)
    );

    vm.stopPrank();

    vm.startPrank(user);

    vm.expectRevert(L1ContractVerification.InvalidProxyAdminAddress.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(differentProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure when safe wallet address doesn't match proxy admin owner
   * @dev Verifies that verification fails if the provided safe wallet address doesn't match proxy admin owner
   */
  function testVerifyL1ContractsFailMismatchedSafeWalletAddress() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
        address(systemConfigProxy),
        IProxyAdmin(address(mockProxyAdmin))
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info with correct owners
    verifier.setSafeConfig(
        tokamakDAO,
        foundation,
        3, // threshold from safeWallet constructor
        implementationCodehash,
        proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail because we're providing a different safe wallet address
    // than what's configured as the owner in the proxy admin, Failing on the ProxyAdmin verification because the owner is different
    vm.expectRevert(L1ContractVerification.InvalidProxyAdminAddress.selector);
    verifier.verifyAndRegisterRollupConfig(
        address(systemConfigProxy),
        IProxyAdmin(address(mockProxyAdmin)),
        "Test Rollup",
        makeAddr("differentSafeWallet")
    );

    vm.stopPrank();
  }

  /**
   * @dev SystemConfig verification Failures
   */
  /**
   * @notice Test verification failure with invalid SystemConfig
   * @dev Modifies the SystemConfig implementation to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidSystemConfig() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO, // _tokamakDAO
      foundation, // _foundation
      3, // _threshold
      implementationCodehash, // _implementationCodehash
    proxyCodehash
    );

    // Deploy a different SystemConfig implementation with different addresses
    MockSystemConfig differentSystemConfigImpl = new MockSystemConfig();

    // Update the SystemConfig proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(systemConfigProxy),
      address(differentSystemConfigImpl)
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert(L1ContractVerification.SystemConfigVerificationFailed.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }


  /**
   * @notice Test verification failure with wrong SystemConfig proxy that has correct native token and implementation
   * @dev Uses a SystemConfig with the correct native token and implementation but wrong proxy codehash
   */
  function testVerifyL1ContractsFailWrongSystemConfigProxyWithCorrectImplementation()
    public
  {
    vm.startPrank(owner);

    // Setup all contract information with the original SystemConfig proxy
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Create a new proxy that points to the original implementation
    TransparentUpgradeableProxy differentProxy = new TransparentUpgradeableProxy(
        address(systemConfigImpl),
        address(mockProxyAdmin),
        abi.encodeWithSelector(
          MockSystemConfig.initialize.selector,
          address(l1StandardBridgeProxy),
          address(l1CrossDomainMessengerProxy),
          address(optimismPortalProxy),
          nativeToken
        )
      );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    vm.expectRevert(L1ContractVerification.GetProxyImplementationFailed.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(differentProxy), // Using different proxy with correct implementation
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @dev L1StandardBridge verification Failures
   */

  /**
   * @notice Test verification failure with invalid L1StandardBridge
   * @dev Modifies the L1StandardBridge implementation to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidL1StandardBridge() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Deploy a different L1StandardBridge implementation
    MockL1StandardBridge differentL1StandardBridgeImpl = new MockL1StandardBridge();

    // Update the L1StandardBridge proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(l1StandardBridgeProxy),
      address(differentL1StandardBridgeImpl)
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert(L1ContractVerification.L1StandardBridgeVerificationFailed.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @dev L1CrossDomainMessenger verification Failures
   */

  /**
   * @notice Test verification failure with invalid L1CrossDomainMessenger
   * @dev Modifies the L1CrossDomainMessenger implementation to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidL1CrossDomainMessenger() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Deploy a different L1CrossDomainMessenger implementation
    MockL1CrossDomainMessenger differentL1CrossDomainMessengerImpl = new MockL1CrossDomainMessenger();

    // Update the L1CrossDomainMessenger proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(differentL1CrossDomainMessengerImpl)
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert(L1ContractVerification.L1CrossDomainMessengerVerificationFailed.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @dev OptimismPortal verification Failures
   */

  /**
   * @notice Test verification failure with invalid OptimismPortal
   * @dev Modifies the OptimismPortal implementation to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidOptimismPortal() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Deploy a different OptimismPortal implementation
    MockOptimismPortal differentOptimismPortalImpl = new MockOptimismPortal();

    // Update the OptimismPortal proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(optimismPortalProxy),
      address(differentOptimismPortalImpl)
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert(L1ContractVerification.OptimismPortalVerificationFailed.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @dev Safe wallet verification Failures
   */
  /**
   * @notice Test verification failure with invalid safe wallet
   * @dev Creates a safe wallet with incorrect owners to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidSafeOwner() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info with different owners than what's in the safe
    verifier.setSafeConfig(
      makeAddr('wrongTokamakDAO'), // Different tokamakDAO address
      makeAddr('wrongFoundation'), // Different foundation address
      3, // _threshold
      implementationCodehash, // _implementationCodehash
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail due to missing required owners
    vm.expectRevert(L1ContractVerification.SafeWalletMissingRequiredOwners.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with invalid safe wallet implementation codehash
   * @dev Sets wrong implementation codehash to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidSafeImplementationCodehash() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Sets wrong implementation codehash
    bytes32 implementationCodehash = address(foundation).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert(L1ContractVerification.SafeWalletInvalidImplCodehash.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with invalid safe wallet proxy codehash
   * @dev Sets wrong proxy codehash to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidSafeProxyCodehash() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Sets wrong proxy codehash
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(foundation).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert(L1ContractVerification.SafeWalletInvalidProxyCodehash.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      "Test Rollup",
      address(safeWallet)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with invalid safe wallet threshold
   * @dev Sets wrong threshold to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidSafeThreshold() public {
    vm.startPrank(owner);

    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    vm.expectRevert(L1ContractVerification.InvalidThreshold.selector);
    verifier.setSafeConfig(
      tokamakDAO, // _tokamakDAO
      foundation, // _foundation
      0, // Zero threshold should cause revert
      implementationCodehash, // _implementationCodehash
      proxyCodehash
    );

    vm.stopPrank();
  }

  /**
   * @notice Test access controls for admin functions
   * @dev Attempts to call admin functions from a non-admin account
   */
  function testOnlyAdminFunctions() public {
    vm.startPrank(randomUser); // Not an admin

    // Each of these should revert due to caller not having ADMIN_ROLE
    vm.expectRevert();
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      bytes32(0),
    bytes32(0)
    );

    vm.expectRevert();
    verifier.setBridgeRegistryAddress(address(0x1));

    // Also test the new admin functions
    vm.expectRevert();
    verifier.addAdmin(address(0x3));

    vm.expectRevert();
    verifier.removeAdmin(address(0x4));

    vm.stopPrank();
  }

  /**
   * @notice Test successful addition of a new admin
   * @dev Owner adds a new admin and verifies they can perform admin operations
   */
  function testAddAdminSuccess() public {
    address newAdmin = makeAddr('newAdmin');

    // Initial state - new admin should not have ADMIN_ROLE
    assertFalse(verifier.hasRole(verifier.ADMIN_ROLE(), newAdmin));

    // Owner grants ADMIN_ROLE to the new admin
    vm.startPrank(owner);
    verifier.addAdmin(newAdmin);
    vm.stopPrank();

    // New admin should now have ADMIN_ROLE
    assertTrue(verifier.hasRole(verifier.ADMIN_ROLE(), newAdmin));

    // Verify that the new admin can perform admin operations
    vm.startPrank(newAdmin);

    // Set bridge registry address as the new admin
    address newBridgeRegistry = makeAddr('newBridgeRegistry');
    verifier.setBridgeRegistryAddress(newBridgeRegistry);

    // Verify the change was successful
    assertEq(verifier.l1BridgeRegistryAddress(), newBridgeRegistry);

    vm.stopPrank();
  }

  /**
   * @notice Test successful removal of an admin
   * @dev Owner removes an admin and verifies they can no longer perform admin operations
   */
  function testRemoveAdminSuccess() public {
    address tempAdmin = makeAddr('tempAdmin');

    // First add the temporary admin
    vm.startPrank(owner);
    verifier.addAdmin(tempAdmin);
    vm.stopPrank();

    // Verify the admin was added correctly
    assertTrue(verifier.hasRole(verifier.ADMIN_ROLE(), tempAdmin));

    // Verify the admin can perform admin operations
    vm.startPrank(tempAdmin);
    address newBridgeRegistry = makeAddr('newBridgeRegistry');
    verifier.setBridgeRegistryAddress(newBridgeRegistry);
    vm.stopPrank();

    // Now remove the admin
    vm.startPrank(owner);
    verifier.removeAdmin(tempAdmin);
    vm.stopPrank();

    // Verify the admin no longer has ADMIN_ROLE
    assertFalse(verifier.hasRole(verifier.ADMIN_ROLE(), tempAdmin));

    // Verify the removed admin can no longer perform admin operations
    vm.startPrank(tempAdmin);

    // This should revert since tempAdmin is no longer an admin
    bytes32 adminRole = verifier.ADMIN_ROLE();
    vm.expectRevert(
      abi.encodePacked(
        "AccessControl: account ",
        StringsUpgradeable.toHexString(uint160(tempAdmin), 20),
        " is missing role ",
        StringsUpgradeable.toHexString(uint256(adminRole), 32)
      )
    );
    verifier.setBridgeRegistryAddress(makeAddr('anotherRegistry'));

    vm.stopPrank();
  }

  /**
   * @notice Test that a non-admin cannot add an admin
   * @dev Non-admin attempts to add a new admin, which should fail
   */
  function testAddAdminFailureNonAdmin() public {
    address nonAdmin = makeAddr('nonAdmin');
    address anotherAccount = makeAddr('anotherAccount');

    // Verify non-admin doesn't have ADMIN_ROLE
    assertFalse(verifier.hasRole(verifier.ADMIN_ROLE(), nonAdmin));

    // Non-admin attempts to add another account as admin
    vm.startPrank(nonAdmin);

    // This should revert with AccessControl error
    vm.expectRevert();
    verifier.addAdmin(anotherAccount);

    vm.stopPrank();

    // Verify that anotherAccount did not get ADMIN_ROLE
    assertFalse(verifier.hasRole(verifier.ADMIN_ROLE(), anotherAccount));
  }

  /**
   * @notice Test that a non-admin cannot remove an admin
   * @dev Non-admin attempts to remove an existing admin, which should fail
   */
  function testRemoveAdminFailureNonAdmin() public {
    address nonAdmin = makeAddr('nonAdmin');

    // Verify non-admin doesn't have ADMIN_ROLE
    assertFalse(verifier.hasRole(verifier.ADMIN_ROLE(), nonAdmin));
    // Verify owner has ADMIN_ROLE
    assertTrue(verifier.hasRole(verifier.ADMIN_ROLE(), owner));

    // Non-admin attempts to remove owner as admin
    vm.startPrank(nonAdmin);

    // This should revert with AccessControl error
    vm.expectRevert();
    verifier.removeAdmin(owner);

    vm.stopPrank();

    // Verify that owner still has ADMIN_ROLE
    assertTrue(verifier.hasRole(verifier.ADMIN_ROLE(), owner));
  }

  /**
   * @notice Test verification failure with different native token
   * @dev Reconfigures the SystemConfig with a different native token to trigger verification failure
   */
  function testVerifyAndRegisterRollupConfigFailInvalidNativeToken() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Reuse existing variables instead of declaring new ones
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Use existing variables where possible
    address differentNativeToken = makeAddr('differentNativeToken');

    // Initialize SystemConfig with a different native token
    (bool success, ) = address(systemConfigProxy).call(
      abi.encodeWithSelector(
        MockSystemConfig.initialize.selector,
        address(l1StandardBridgeProxy),
        address(l1CrossDomainMessengerProxy),
        address(optimismPortalProxy),
        differentNativeToken
      )
    );
    require(success, 'SystemConfig re-initialization failed');

    vm.stopPrank();

    vm.startPrank(user);

    // Create and setup new verifier in a more optimized way
    L1ContractVerification newVerifier = L1ContractVerification(
      address(
        new TransparentUpgradeableProxy(
          address(new L1ContractVerification()),
          address(new ProxyAdmin()),
          abi.encodeWithSelector(
            L1ContractVerification.initialize.selector,
            nativeToken,
            user
          )
        )
      )
    );

    // Setup the new verifier
    newVerifier.setBridgeRegistryAddress(address(bridgeRegistry));
    newVerifier.setLogicContractInfo(address(systemConfigProxy), IProxyAdmin(address(mockProxyAdmin)));
    newVerifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
       proxyCodehash
    );
    newVerifier.setVerificationPossible(true);

    // Test the verification
    vm.expectRevert(L1ContractVerification.NativeTokenNotTON.selector);
    newVerifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      'TestRollup',
      address(safeWallet) // Pass safe wallet address
    );

    vm.stopPrank();
  }

  /**
   * @notice Test safe wallet verification with large owner set
   * @dev Creates a safe with many owners to test the verification works with large owner sets
   */
  function testVerifyL1ContractsWithLargeOwnerSet() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Create a large set of owners (20 addresses)
    address[] memory largeOwnerSet = new address[](20);
    largeOwnerSet[0] = tokamakDAO;
    largeOwnerSet[1] = foundation;
    for (uint i = 2; i < 20; i++) {
      largeOwnerSet[i] = address(uint160(0x1000 + i));
    }

    // Create a new safe with many owners
    MockGnosisSafe largeSafe = new MockGnosisSafe(largeOwnerSet, 3);

    // Create a new proxy admin for the large safe
    MockProxyAdmin largeProxyAdmin = new MockProxyAdmin(owner);
    largeProxyAdmin.setOwner(address(largeSafe));

    // Set up implementations for the large proxy admin
    largeProxyAdmin.setImplementation(
      address(systemConfigProxy),
      address(systemConfigImpl)
    );
    largeProxyAdmin.setImplementation(
      address(l1StandardBridgeProxy),
      address(l1StandardBridgeImpl)
    );
    largeProxyAdmin.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(l1CrossDomainMessengerImpl)
    );
    largeProxyAdmin.setImplementation(
      address(optimismPortalProxy),
      address(optimismPortalImpl)
    );

    // Set admin relationships for each proxy
    largeProxyAdmin.setAdmin(
      address(systemConfigProxy),
      address(largeProxyAdmin)
    );
    largeProxyAdmin.setAdmin(
      address(l1StandardBridgeProxy),
      address(largeProxyAdmin)
    );
    largeProxyAdmin.setAdmin(
      address(l1CrossDomainMessengerProxy),
      address(largeProxyAdmin)
    );
    largeProxyAdmin.setAdmin(
      address(optimismPortalProxy),
      address(largeProxyAdmin)
    );

    // Get safe wallet codehashes
    bytes32 implementationCodehash = address(largeSafe).codehash;
    bytes32 proxyCodehash = address(largeSafe).codehash;

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(largeProxyAdmin))
    );

    // Set safe wallet configuration
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    // Grant ADMIN_ROLE to the user
    verifier.addAdmin(user);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail due to large number of owners
    vm.expectRevert(L1ContractVerification.SafeWalletWrongOwnerCount.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(largeProxyAdmin)),
      "Test Rollup",
      address(largeSafe) // Use largeSafe instead of safeWallet
    );

    vm.stopPrank();
  }

  /**
   * @notice Test bridge registry address validation
   * @dev Attempts to set a zero address for the bridge registry
   */
  function testSetBridgeRegistryAddressFailZeroAddress() public {
    vm.startPrank(owner);

    vm.expectRevert(abi.encodeWithSelector(L1ContractVerification.ZeroAddress.selector, "bridgeRegistry"));
    verifier.setBridgeRegistryAddress(address(0));

    vm.stopPrank();
  }

  /**
   * @notice Test suite for upgrade functionality
   */
  function testUpgradeWithoutMultisigShouldFail() public {
    vm.startPrank(randomUser);

    // Deploy new implementation
    L1ContractVerification newImplementation = new L1ContractVerification();

    // Attempt to upgrade without multisig approval
    vm.expectRevert("Ownable: caller is not the owner");
    verifierProxyAdmin.upgrade(
        verifierProxy,
        address(newImplementation)
    );

    vm.stopPrank();
  }

  function testUpgradeWithMultisigSuccess() public {
    vm.startPrank(owner);

    // Deploy new implementation
    L1ContractVerification newImplementation = new L1ContractVerification();

    // Store the current implementation address
    address currentImplementation = verifierProxyAdmin.getProxyImplementation(verifierProxy);

    // Transfer ProxyAdmin ownership to multisig
    verifierProxyAdmin.transferOwnership(address(safeWallet));

    vm.stopPrank();

    // Simulate multisig wallet call
    vm.prank(address(safeWallet));

    // Upgrade to new implementation
    verifierProxyAdmin.upgrade(
        verifierProxy,
        address(newImplementation)
    );

    // Verify the upgrade was successful
    address newImplementationAddress = verifierProxyAdmin.getProxyImplementation(verifierProxy);
    assertNotEq(currentImplementation, newImplementationAddress, "Implementation should have changed");
    assertEq(newImplementationAddress, address(newImplementation), "New implementation not set correctly");
  }

  function testUpgradeWithoutOwnershipTransferShouldSucceed() public {
    vm.startPrank(owner);

    // Deploy new implementation
    L1ContractVerification newImplementation = new L1ContractVerification();

    // Store the current implementation address
    address currentImplementation = verifierProxyAdmin.getProxyImplementation(verifierProxy);

    // Attempt to upgrade without transferring ownership
    verifierProxyAdmin.upgrade(
        verifierProxy,
        address(newImplementation)
    );

    // This should succeed as ownership hasn't been transferred to multisig yet
    address newImplementationAddress = verifierProxyAdmin.getProxyImplementation(verifierProxy);
    assertNotEq(currentImplementation, newImplementationAddress, "Implementation should have changed");
    assertEq(newImplementationAddress, address(newImplementation), "New implementation not set correctly");

    vm.stopPrank();
  }

  function testVerifyProxyAdminOwnership() public {
    vm.startPrank(owner);

    // Transfer ownership to multisig
    verifierProxyAdmin.transferOwnership(address(safeWallet));

    // Verify ownership
    assertEq(verifierProxyAdmin.owner(), address(safeWallet), "Ownership not transferred correctly");

    // Try to upgrade after ownership transfer (should fail)
    L1ContractVerification newImplementation = new L1ContractVerification();

    vm.expectRevert("Ownable: caller is not the owner");
    verifierProxyAdmin.upgrade(
        verifierProxy,
        address(newImplementation)
    );

    vm.stopPrank();
  }

  function testMultisigOwnershipRequirement() public {
    vm.startPrank(owner);

    // Transfer ownership to multisig
    verifierProxyAdmin.transferOwnership(address(safeWallet));

    // Deploy new implementation
    L1ContractVerification newImplementation = new L1ContractVerification();

    vm.stopPrank();

    // Try to upgrade from each safe owner individually (should fail)
    address[] memory owners = new address[](3);
    owners[0] = tokamakDAO;
    owners[1] = foundation;
    owners[2] = thirdOwner;

    for(uint i = 0; i < owners.length; i++) {
        vm.prank(owners[i]);
        vm.expectRevert("Ownable: caller is not the owner");
        verifierProxyAdmin.upgrade(
            verifierProxy,
            address(newImplementation)
        );
    }
  }

  function testProxyAdminOwnershipAfterDeployment() public {
    // Deploy contracts as in the deployment script
    vm.startPrank(owner);

    // Deploy implementation
    L1ContractVerification implementationContract = new L1ContractVerification();

    // Deploy ProxyAdmin
    ProxyAdmin verificationContractProxyAdmin = new ProxyAdmin();

    // Deploy proxy
    bytes memory initData = abi.encodeWithSelector(
        L1ContractVerification.initialize.selector,
        nativeToken,
        owner
    );

    TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
        address(implementationContract),
        address(verificationContractProxyAdmin),
        initData
    );

    // Transfer ownership to multisig
    verificationContractProxyAdmin.transferOwnership(address(safeWallet));

    vm.stopPrank();

    // Try to upgrade with original owner - should fail
    vm.startPrank(owner);
    L1ContractVerification newImplementation = new L1ContractVerification();

    vm.expectRevert("Ownable: caller is not the owner");
    verificationContractProxyAdmin.upgrade(
        proxy,
        address(newImplementation)
    );

    vm.stopPrank();

    // Verify owner is actually the multisig
    assertEq(verificationContractProxyAdmin.owner(), address(safeWallet), "ProxyAdmin owner should be multisig");
  }

  function testSetSafeConfigFailZeroThreshold() public {
    vm.startPrank(owner);

    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    vm.expectRevert(L1ContractVerification.InvalidThreshold.selector);
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      0, // threshold from safeWallet constructor
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();
  }

  /**
   * @notice Test that different safe wallets can be verified
   * @dev Sets up multiple safe wallets and verifies them independently
   */
  function testMultipleOperatorSafeWallets() public {
    vm.startPrank(owner);

    // Setup contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Create first safe wallet and proxy admin
    address[] memory owners1 = new address[](3);
    owners1[0] = tokamakDAO;
    owners1[1] = foundation;
    owners1[2] = thirdOwner;
    MockGnosisSafe safe1 = new MockGnosisSafe(owners1, 3);

    MockProxyAdmin proxyAdmin1 = new MockProxyAdmin(owner);
    proxyAdmin1.setOwner(address(safe1));

    // Create second safe wallet and proxy admin
    address[] memory owners2 = new address[](3);
    owners2[0] = tokamakDAO;
    owners2[1] = foundation;
    owners2[2] = makeAddr('differentOwner');
    MockGnosisSafe safe2 = new MockGnosisSafe(owners2, 3);

    MockProxyAdmin proxyAdmin2 = new MockProxyAdmin(owner);
    proxyAdmin2.setOwner(address(safe2));

    // Set up implementations for both proxy admins
    proxyAdmin1.setImplementation(
      address(systemConfigProxy),
      address(systemConfigImpl)
    );
    proxyAdmin1.setImplementation(
      address(l1StandardBridgeProxy),
      address(l1StandardBridgeImpl)
    );
    proxyAdmin1.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(l1CrossDomainMessengerImpl)
    );
    proxyAdmin1.setImplementation(
      address(optimismPortalProxy),
      address(optimismPortalImpl)
    );

    proxyAdmin2.setImplementation(
      address(systemConfigProxy),
      address(systemConfigImpl)
    );
    proxyAdmin2.setImplementation(
      address(l1StandardBridgeProxy),
      address(l1StandardBridgeImpl)
    );
    proxyAdmin2.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(l1CrossDomainMessengerImpl)
    );
    proxyAdmin2.setImplementation(
      address(optimismPortalProxy),
      address(optimismPortalImpl)
    );

    // Set up admins for both proxy admins
    proxyAdmin1.setAdmin(address(systemConfigProxy), address(proxyAdmin1));
    proxyAdmin1.setAdmin(address(l1StandardBridgeProxy), address(proxyAdmin1));
    proxyAdmin1.setAdmin(address(l1CrossDomainMessengerProxy), address(proxyAdmin1));
    proxyAdmin1.setAdmin(address(optimismPortalProxy), address(proxyAdmin1));

    proxyAdmin2.setAdmin(address(systemConfigProxy), address(proxyAdmin2));
    proxyAdmin2.setAdmin(address(l1StandardBridgeProxy), address(proxyAdmin2));
    proxyAdmin2.setAdmin(address(l1CrossDomainMessengerProxy), address(proxyAdmin2));
    proxyAdmin2.setAdmin(address(optimismPortalProxy), address(proxyAdmin2));

    // Set safe wallet configurations
    verifier.setSafeConfig(
      tokamakDAO, // _tokamakDAO
      foundation, // _foundation
      3, // _threshold
      safe1.masterCopy().codehash, // _implementationCodehash
      address(safe1).codehash // _proxyCodehash
    );

    // Grant admin roles
    verifier.addAdmin(user);

    vm.stopPrank();

    // Test verification with first safe wallet
    vm.startPrank(user);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(proxyAdmin1)),
      "Test Rollup",
      address(safe1)
    );
    vm.stopPrank();

    // Test verification with second safe wallet
    vm.startPrank(user);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(proxyAdmin2)),
      "Test Rollup 2",
      address(safe2)
    );
    vm.stopPrank();

    // Test that first safe wallet cannot use second safe wallet's proxy admin
    vm.startPrank(user);
    vm.expectRevert(L1ContractVerification.InvalidProxyAdminAddress.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(proxyAdmin2)),
      "Test Rollup 3",
      address(safe1)
    );
    vm.stopPrank();
  }

  /**
   * @notice Test verification failure when modules are set on the safe wallet
   * @dev Verifies that a safe wallet with modules will fail verification
   */
  function testVerifyL1ContractsFailWithModules() public {
    vm.startPrank(owner);

    // Setup contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Create a safe wallet with a module
    address[] memory safeOwners = new address[](3);
    safeOwners[0] = tokamakDAO;
    safeOwners[1] = foundation;
    safeOwners[2] = thirdOwner;
    MockGnosisSafe maliciousSafe = new MockGnosisSafe(safeOwners, 3);

    // Create a module address
    address moduleAddress = makeAddr("maliciousModule");

    // Create a proxy admin owned by the safe
    MockProxyAdmin safeProxyAdmin = new MockProxyAdmin(owner);
    safeProxyAdmin.setOwner(address(maliciousSafe));

    // Set up implementations
    safeProxyAdmin.setImplementation(
      address(systemConfigProxy),
      address(systemConfigImpl)
    );
    safeProxyAdmin.setImplementation(
      address(l1StandardBridgeProxy),
      address(l1StandardBridgeImpl)
    );
    safeProxyAdmin.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(l1CrossDomainMessengerImpl)
    );
    safeProxyAdmin.setImplementation(
      address(optimismPortalProxy),
      address(optimismPortalImpl)
    );

    // Set up admin relationships
    safeProxyAdmin.setAdmin(address(systemConfigProxy), address(safeProxyAdmin));
    safeProxyAdmin.setAdmin(address(l1StandardBridgeProxy), address(safeProxyAdmin));
    safeProxyAdmin.setAdmin(address(l1CrossDomainMessengerProxy), address(safeProxyAdmin));
    safeProxyAdmin.setAdmin(address(optimismPortalProxy), address(safeProxyAdmin));

    // Set safe wallet configurations
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3,
      maliciousSafe.masterCopy().codehash,
      address(maliciousSafe).codehash
    );

    // Grant admin roles
    verifier.addAdmin(user);

    vm.stopPrank();

    // Create a modules array with one module
    address[] memory modules = new address[](1);
    modules[0] = moduleAddress;

    // Mock the getModulesPaginated function to return modules
    vm.mockCall(
      address(maliciousSafe),
      abi.encodeWithSignature("getModulesPaginated(address,uint256)", address(1), 1),
      abi.encode(modules, address(1))
    );

    // Test that verification fails when modules are present
    vm.startPrank(user);
    vm.expectRevert(L1ContractVerification.SafeWalletUnauthorizedModules.selector);
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      IProxyAdmin(address(safeProxyAdmin)),
      "Test Rollup",
      address(maliciousSafe)
    );
    vm.stopPrank();
  }

  /**
   * @notice Test verification success when called from a contract
   * @dev Creates a caller contract to test that contract calls are allowed
   */
  function testVerifyL1ContractsFailContractCall() public {
    vm.startPrank(owner);

    // Setup contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin))
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3,
      implementationCodehash,
      proxyCodehash
    );

    // Deploy a caller contract
    ContractCaller callerContract = new ContractCaller(address(verifier));

    // Grant ADMIN_ROLE to the caller contract
    verifier.addAdmin(address(callerContract));
    verifier.setVerificationPossible(true);

    vm.stopPrank();

    // Call verify through the caller contract - this should succeed
    callerContract.callVerify(
      address(systemConfigProxy),
      IProxyAdmin(address(mockProxyAdmin)),
      address(safeWallet)
    );
  }
}

// Helper contract to test contract calls
contract ContractCaller {
    IL1ContractVerification private verifier;

    constructor(address _verifier) {
        verifier = IL1ContractVerification(_verifier);
    }

    function callVerify(
        address systemConfigProxy,
        IProxyAdmin proxyAdmin,
        address safeWalletAddress
    ) external {
        verifier.verifyAndRegisterRollupConfig(
            systemConfigProxy,
            proxyAdmin,
            "Test Rollup",
            safeWalletAddress
        );
    }
}
