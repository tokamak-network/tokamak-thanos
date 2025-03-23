// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

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
      nativeToken
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
      address(mockProxyAdmin)
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Get implementation address via masterCopy function
    address implementation = safeWallet.masterCopy();

    // Correctly set the safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementation.codehash,
      proxyCodehash
    );

    // Log information for debugging
    console.log('-----------------------------------------------------------');
    console.log('SystemConfig Proxy:', address(systemConfigProxy));
    console.log('SystemConfig Implementation:', address(systemConfigImpl));
    console.log(
      'L1StandardBridge from SystemConfig:',
      MockSystemConfig(address(systemConfigProxy)).l1StandardBridge()
    );
    console.log(
      'L1CrossDomainMessenger from SystemConfig:',
      MockSystemConfig(address(systemConfigProxy)).l1CrossDomainMessenger()
    );
    console.log(
      'OptimismPortal from SystemConfig:',
      MockSystemConfig(address(systemConfigProxy)).optimismPortal()
    );
    console.log(
      'NativeToken from SystemConfig:',
      MockSystemConfig(address(systemConfigProxy)).nativeTokenAddress()
    );
    console.log('ProxyAdmin Owner:', mockProxyAdmin.owner());
    console.log('-----------------------------------------------------------');

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should succeed
    bool result = verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Verification should succeed
    assertTrue(result);

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
      address(mockProxyAdmin)
    );

    // Store the actual codehashes of the deployed safe wallet
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Get implementation address via masterCopy function
    address implementation = safeWallet.masterCopy();

    // Correctly set the safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementation.codehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification and registration should succeed
    bool result = verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      address(mockProxyAdmin),
      2, // TON token
      address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000), //L2 Ton address
      'TestRollup'
    );

    // Verification and registration should succeed
    assertTrue(result);

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
  function testVerifyL1ContractsFailInvalidProxyAdmin() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Passing a different contract so that verification fails
    // Even if we redeployed the correct proxy admin contract for different address,
    // contract will match the codehash and it will pass
    MockL1StandardBridge wrongProxyAdmin = new MockL1StandardBridge();

    // Verification should fail with ProxyAdmin verification error
    vm.expectRevert('ProxyAdmin verification failed: invalid codehash');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(wrongProxyAdmin)
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
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info with the current safe wallet address
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    // Create a new proxy admin with a different owner
    MockProxyAdmin differentProxyAdmin = new MockProxyAdmin(owner);

    // Set the owner of the different proxy admin to the safe wallet
    differentProxyAdmin.setOwner(address(differentProxyAdmin));

    // Set up the different proxy admin with the same implementations
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

    vm.stopPrank();

    vm.startPrank(user);

    vm.expectRevert('Invalid proxy admin address');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(differentProxyAdmin) // Using different proxy admin
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with mismatched safe wallet address
   * @dev Sets up a different safe wallet address in the configuration
   */
  function testVerifyL1ContractsFailMismatchedSafeWalletAddress() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Create a new safe wallet
    address[] memory safeOwners = new address[](3);
    safeOwners[0] = tokamakDAO;
    safeOwners[1] = foundation;
    safeOwners[2] = thirdOwner;
    MockGnosisSafe differentSafeWallet = new MockGnosisSafe(safeOwners, 3);

    // Create a new proxy admin for the different safe wallet
    MockProxyAdmin differentProxyAdmin = new MockProxyAdmin(owner);
    differentProxyAdmin.setOwner(address(differentSafeWallet));

    // Set up the different proxy admin with the same implementations
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

    // Get codehashes for the different safe wallet
    bytes32 implementationCodehash = address(differentSafeWallet).codehash;
    bytes32 proxyCodehash = address(differentSafeWallet).codehash;

    // Setup safe wallet info with the DIFFERENT proxy admin
    // This should set the expected safe wallet address to differentSafeWallet
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(differentProxyAdmin), // Use the different proxy admin here
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    vm.expectRevert('Invalid proxy admin address');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin) // Use the original proxy admin here
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
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    // Deploy a different SystemConfig implementation with different addresses
    MockSystemConfig differentSystemConfigImpl = new MockSystemConfig();

    // Update the SystemConfig proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(systemConfigProxy),
      address(differentSystemConfigImpl)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('SystemConfig verification failed');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
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
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
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

    vm.stopPrank();

    vm.startPrank(user);

    vm.expectRevert('SystemConfig verification failed');
    verifier.verifyL1Contracts(
      address(differentProxy), // Using different proxy with correct implementation
      address(mockProxyAdmin)
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
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
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

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('L1StandardBridge verification failed');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
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
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
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

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('L1CrossDomainMessenger verification failed');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
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
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
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

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('OptimismPortal verification failed');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    vm.stopPrank();
  }

  /**
   * @dev Safe wallet verification Failures
  /**
   * @notice Test verification failure with invalid safe wallet
   * @dev Creates a different safe wallet with incorrect owners to trigger verification failure
   */
  function testVerifyL1ContractsFailInvalidSafe() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      makeAddr('random'), //Set different foundation address
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Safe wallet verification failed: missing required owners');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with invalid safe wallet
   * @dev Get error using invalid safe wallet codehash
   */
  function testVerifyL1ContractsFailInvalidSafeImplementationCodehash() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Sets wrong implementation codehash
    bytes32 implementationCodehash = address(foundation).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Safe wallet verification failed: invalid implementation codehash');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with invalid safe wallet
   * @dev Get error using invalid safe wallet proxy codehash
   */
  function testVerifyL1ContractsFailInvalidSafeProxyCodehash() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Sets wrong implementation codehash
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(foundation).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Safe wallet verification failed: invalid proxy codehash');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with invalid safe wallet
   * @dev Get error using invalid safe wallet threshold
   */
  function testVerifyL1ContractsFailInvalidSafeThreshold() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Sets wrong implementation codehash
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      4, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Safe wallet verification failed: invalid threshold');
    verifier.verifyL1Contracts(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    vm.stopPrank();
  }

  /**
   * @dev Other Failures
  /**
   * @notice Test verification failure with invalid token type
   * @dev Attempts registration with an invalid token type
   */
  function testVerifyAndRegisterRollupConfigFailInvalidType() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail due to invalid token type
    vm.expectRevert('Registration allowed only for TON tokens');
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      address(mockProxyAdmin),
      1, // Invalid token type
      address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000), //L2 Ton address
      'TestRollup'
    );

    vm.stopPrank();
  }

  /**
   * @notice Test verification failure with invalid L2 TON address
   * @dev Attempts registration with an incorrect L2 TON address
   */
  function testVerifyAndRegisterRollupConfigFailInvalidL2TON() public {
    vm.startPrank(owner);

    // Setup all contract information
    verifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    vm.stopPrank();

    vm.startPrank(user);

    address invalidL2TON = makeAddr('invalidL2TON'); // Different from nativeToken

    // Verification should fail due to invalid L2 TON token address
    vm.expectRevert('Provided L2 TON Token address is not correct');
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      address(mockProxyAdmin),
      2, // Valid token type
      invalidL2TON, // Invalid L2TON address
      'TestRollup'
    );

    vm.stopPrank();
  }

  /**
   * @notice Test access controls on owner-only functions
   * @dev Attempts to call owner functions from a non-owner account
   */
  function testOnlyOwnerFunctions() public {
    vm.startPrank(user);

    // All these should revert because user is not the owner
    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setLogicContractInfo(address(0), address(0));

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setSafeConfig(
      address(0),
      address(0),
      0,
      address(0),
      bytes32(0),
      bytes32(0)
    );

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setBridgeRegistryAddress(address(0));

    vm.stopPrank();
  }

  /**
   * @notice Test bridge registry address validation
   * @dev Attempts to set a zero address for the bridge registry
   */
  function testSetBridgeRegistryAddressFailZeroAddress() public {
    vm.startPrank(owner);

    vm.expectRevert('Bridge registry address cannot be zero');
    verifier.setBridgeRegistryAddress(address(0));

    vm.stopPrank();
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
      address(mockProxyAdmin)
    );

    // Get codehashes for safe wallet
    bytes32 implementationCodehash = address(safeWallet).codehash;
    bytes32 proxyCodehash = address(safeWallet).codehash;

    // Setup safe wallet info
    verifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3, // threshold
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    address differentNativeToken = vm.addr(1234);
    vm.label(differentNativeToken, 'differentNativeToken');

    // Initialize SystemConfig with a different native token
    bytes memory initData = abi.encodeWithSelector(
      MockSystemConfig.initialize.selector,
      address(l1StandardBridgeProxy),
      address(l1CrossDomainMessengerProxy),
      address(optimismPortalProxy),
      differentNativeToken
    );

    // Re-initialize the SystemConfig with a different token
    (bool success, ) = address(systemConfigProxy).call(initData);
    require(success, 'SystemConfig re-initialization failed');

    vm.stopPrank();

    vm.startPrank(user);

    // Create a new L1ContractVerification instance that expects the original native token
    L1ContractVerification newVerifierImpl = new L1ContractVerification();

    // Create a proxy admin
    ProxyAdmin newVerifierProxyAdmin = new ProxyAdmin();

    // Initialize with the original native token
    bytes memory newVerifierInitData = abi.encodeWithSelector(
      L1ContractVerification.initialize.selector,
      nativeToken
    );

    // Create proxy for the new verifier
    TransparentUpgradeableProxy newVerifierProxy = new TransparentUpgradeableProxy(
      address(newVerifierImpl),
      address(newVerifierProxyAdmin),
      newVerifierInitData
    );

    // Create a reference to the new verifier
    L1ContractVerification newVerifier = L1ContractVerification(address(newVerifierProxy));

    // Set contract info on the new verifier
    newVerifier.setBridgeRegistryAddress(address(bridgeRegistry));
    newVerifier.setLogicContractInfo(
      address(systemConfigProxy),
      address(mockProxyAdmin)
    );

    // Setup safe wallet info on the new verifier
    newVerifier.setSafeConfig(
      tokamakDAO,
      foundation,
      3,
      address(mockProxyAdmin),
      implementationCodehash,
      proxyCodehash
    );

    // Verification should fail due to using a different native token
    vm.expectRevert('The native token you are using is not TON');
    newVerifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      address(mockProxyAdmin),
      2, // TON token type
      address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000), //L2 Ton address
      'TestRollup'
    );

    vm.stopPrank();
  }
}
