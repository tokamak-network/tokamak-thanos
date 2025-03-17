// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Test.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';

import './mock-contracts/MockContracts.sol';
import './mock-contracts/MockGenosisSafe.sol';
import './mock-contracts/MockProxyAdmin.sol';
import './mock-contracts/MockSystemConfig.sol';

contract L1ContractVerificationTest is Test {
  // Constants
  uint256 constant CHAIN_ID = 11155111; // Sepolia
  bytes32 constant L1_STANDARD_BRIDGE_ID = keccak256('L1_STANDARD_BRIDGE');
  bytes32 constant L1_CROSS_DOMAIN_MESSENGER_ID =
    keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 constant OPTIMISM_PORTAL_ID = keccak256('OPTIMISM_PORTAL');
  bytes32 constant SYSTEM_CONFIG_ID = keccak256('SYSTEM_CONFIG');

  // Contracts
  L1ContractVerification verifier;
  MockProxyAdmin mockProxyAdmin;

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
  address thirdOwner;
  address user;
  address nativeToken;
  address l2TONAddress;

  function setUp() public {
    // Setup accounts
    owner = makeAddr('owner');
    tokamakDAO = makeAddr('tokamakDAO');
    foundation = makeAddr('foundation');
    thirdOwner = makeAddr('thirdOwner');
    user = makeAddr('user');
    nativeToken = makeAddr('nativeToken');
    l2TONAddress = address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000);

    // Set chain ID for testing
    vm.chainId(CHAIN_ID);

    vm.startPrank(owner);

    // Deploy proxy admin
    mockProxyAdmin = new MockProxyAdmin(owner);

    // Deploy implementations
    l1StandardBridgeImpl = new MockL1StandardBridge();
    l1CrossDomainMessengerImpl = new MockL1CrossDomainMessenger();
    optimismPortalImpl = new MockOptimismPortal();

    // Deploy proxies (without initializing SystemConfig yet)
    l1StandardBridgeProxy = new TransparentUpgradeableProxy(
      address(l1StandardBridgeImpl),
      address(mockProxyAdmin),
      ''
    );

    l1CrossDomainMessengerProxy = new TransparentUpgradeableProxy(
      address(l1CrossDomainMessengerImpl),
      address(mockProxyAdmin),
      ''
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

    // Deploy verifier contract with proxy admin address
    verifier = new L1ContractVerification(
      address(mockProxyAdmin),
      l2TONAddress
    );

    // Set bridge registry address
    verifier.setBridgeRegistryAddress(address(bridgeRegistry));

    // Set expected native token
    verifier.setExpectedNativeToken(nativeToken);

    // Set safe verification required to true for testing
    verifier.setSafeVerificationRequired(true);

    // Set the owner of the proxy admin to the safe wallet for verification
    mockProxyAdmin.setOwner(address(safeWallet));

    vm.stopPrank();
  }

  function testSetContractConfig() public {
    vm.startPrank(owner);

    // Set contract config for SystemConfig
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      address(systemConfigImpl).codehash,
      address(systemConfigProxy).codehash,
      address(mockProxyAdmin)
    );

    // Verify the config was set correctly
    (
      bytes32 implementationHash,
      bytes32 proxyHash,
      address expectedProxyAdmin
    ) = verifier.getContractConfig(SYSTEM_CONFIG_ID);

    assertEq(implementationHash, address(systemConfigImpl).codehash);
    assertEq(proxyHash, address(systemConfigProxy).codehash);
    assertEq(expectedProxyAdmin, address(mockProxyAdmin));

    vm.stopPrank();
  }

  function testSetSafeConfig() public {
    vm.startPrank(owner);

    // Set safe config
    verifier.setSafeConfig(tokamakDAO, foundation, thirdOwner, 3);

    // Verify the config was set correctly
    (
      address configTokamakDAO,
      address configFoundation,
      address configThirdOwner,
      uint256 requiredThreshold
    ) = verifier.getSafeConfig();

    assertEq(configTokamakDAO, tokamakDAO);
    assertEq(configFoundation, foundation);
    assertEq(configThirdOwner, thirdOwner);
    assertEq(requiredThreshold, 3);

    vm.stopPrank();
  }

  function testVerifyL1ContractsSuccess() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

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

    // Verify L1 contracts
    bool result = verifier.verifyL1Contracts(address(systemConfigProxy));

    // Verification should succeed
    assertTrue(result);

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailInvalidSystemConfig() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Deploy a different SystemConfig implementation with different addresses
    MockSystemConfig differentSystemConfigImpl = new MockSystemConfig();

    // Update the SystemConfig proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(systemConfigProxy),
      address(differentSystemConfigImpl)
    );

    // Set an incorrect implementation hash in the contract config
    // This will cause the verification to fail because the actual implementation hash
    // won't match what's expected in the config
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      bytes32(uint256(0x123456)), // Incorrect implementation hash
      address(systemConfigProxy).codehash,
      address(mockProxyAdmin)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('SystemConfig verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailInvalidSafe() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Create a different safe with wrong threshold and owners
    address[] memory safeOwners = new address[](2);
    safeOwners[0] = tokamakDAO;
    safeOwners[1] = foundation;
    // Missing thirdOwner and wrong threshold
    MockGnosisSafe differentSafe = new MockGnosisSafe(safeOwners, 1);

    // Set the owner of the proxy admin to the different safe wallet
    mockProxyAdmin.setOwner(address(differentSafe));

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Safe verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyAndRegisterRollupConfig() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    vm.stopPrank();

    vm.startPrank(user);

    address l2TON = address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000);

    // Verify and register rollup config
    bool result = verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      2, // TON token type
      l2TON,
      'TestRollup'
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

    address l2TON = address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000);

    // Verification should fail due to invalid token type
    vm.expectRevert('Registration allowed only for TON tokens');
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      1, // Invalid token type
      l2TON,
      'TestRollup'
    );

    vm.stopPrank();
  }

  function testOnlyOwnerFunctions() public {
    vm.startPrank(user);

    // All these should revert because user is not the owner
    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      bytes32(0),
      bytes32(0),
      address(0)
    );

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setSafeConfig(address(0), address(0), address(0), 0);

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setBridgeRegistryAddress(address(0));

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setExpectedNativeToken(address(0));

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setSafeVerificationRequired(false);

    // Test chainId versions too
    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      bytes32(0),
      bytes32(0),
      address(0)
    );

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setSafeConfig(address(0), address(0), address(0), 0);

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setExpectedNativeToken(address(0));

    vm.expectRevert('Ownable: caller is not the owner');
    verifier.setSafeVerificationRequired(false);

    vm.stopPrank();
  }

  function testSetBridgeRegistryAddress() public {
    vm.startPrank(owner);

    address newBridgeRegistry = makeAddr('newBridgeRegistry');
    verifier.setBridgeRegistryAddress(newBridgeRegistry);

    assertEq(verifier.L1BridgeRegistryV1_1Address(), newBridgeRegistry);

    vm.stopPrank();
  }

  function testSetBridgeRegistryAddressFailZeroAddress() public {
    vm.startPrank(owner);

    vm.expectRevert('Invalid bridge registry address');
    verifier.setBridgeRegistryAddress(address(0));

    vm.stopPrank();
  }

  function testSetExpectedNativeToken() public {
    vm.startPrank(owner);

    address newToken = makeAddr('newToken');
    verifier.setExpectedNativeToken(newToken);

    // We don't have a getter for this, so we'll test it indirectly in other tests

    vm.stopPrank();
  }

  function testSetExpectedNativeTokenFailZeroAddress() public {
    vm.startPrank(owner);

    vm.expectRevert('Invalid token address');
    verifier.setExpectedNativeToken(address(0));

    vm.stopPrank();
  }

  function testSetSafeVerificationRequired() public {
    vm.startPrank(owner);

    // Set to false
    verifier.setSafeVerificationRequired(false);

    // Set back to true
    verifier.setSafeVerificationRequired(true);

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailInvalidL1StandardBridge() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Deploy a different L1StandardBridge implementation
    MockL1StandardBridge differentL1StandardBridgeImpl = new MockL1StandardBridge();

    // Update the L1StandardBridge proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(l1StandardBridgeProxy),
      address(differentL1StandardBridgeImpl)
    );

    // Set an incorrect implementation hash in the contract config
    verifier.setContractConfig(
      L1_STANDARD_BRIDGE_ID,
      bytes32(uint256(0x123456)), // Incorrect implementation hash
      address(l1StandardBridgeProxy).codehash,
      address(mockProxyAdmin)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('L1 Standard Bridge verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailInvalidL1CrossDomainMessenger() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Deploy a different L1CrossDomainMessenger implementation
    MockL1CrossDomainMessenger differentL1CrossDomainMessengerImpl = new MockL1CrossDomainMessenger();

    // Update the L1CrossDomainMessenger proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(l1CrossDomainMessengerProxy),
      address(differentL1CrossDomainMessengerImpl)
    );

    // Set an incorrect implementation hash in the contract config
    verifier.setContractConfig(
      L1_CROSS_DOMAIN_MESSENGER_ID,
      bytes32(uint256(0x123456)), // Incorrect implementation hash
      address(l1CrossDomainMessengerProxy).codehash,
      address(mockProxyAdmin)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('L1 Cross Domain Messenger verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailInvalidOptimismPortal() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Deploy a different OptimismPortal implementation
    MockOptimismPortal differentOptimismPortalImpl = new MockOptimismPortal();

    // Update the OptimismPortal proxy to point to the different implementation
    mockProxyAdmin.setImplementation(
      address(optimismPortalProxy),
      address(differentOptimismPortalImpl)
    );

    // Set an incorrect implementation hash in the contract config
    verifier.setContractConfig(
      OPTIMISM_PORTAL_ID,
      bytes32(uint256(0x123456)), // Incorrect implementation hash
      address(optimismPortalProxy).codehash,
      address(mockProxyAdmin)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Optimism Portal verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailInvalidProxyAdmin() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Set an incorrect proxy admin in the contract config
    address differentProxyAdmin = makeAddr('differentProxyAdmin');
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      address(systemConfigImpl).codehash,
      address(systemConfigProxy).codehash,
      differentProxyAdmin
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('SystemConfig verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailInvalidProxyHash() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Set an incorrect proxy hash in the contract config
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      address(systemConfigImpl).codehash,
      bytes32(uint256(0x123456)), // Incorrect proxy hash
      address(mockProxyAdmin)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('SystemConfig verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyAndRegisterRollupConfigFailInvalidNativeToken() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Set a different native token in the system config
    address differentNativeToken = makeAddr('differentNativeToken');

    // Re-initialize the SystemConfig proxy with a different native token
    bytes memory initData = abi.encodeWithSelector(
      MockSystemConfig.initialize.selector,
      address(l1StandardBridgeProxy),
      address(l1CrossDomainMessengerProxy),
      address(optimismPortalProxy),
      differentNativeToken
    );

    // Call initialize through the proxy
    (bool success, ) = address(systemConfigProxy).call(initData);
    require(success, 'SystemConfig initialization failed');

    vm.stopPrank();

    vm.startPrank(user);

    address l2TON = makeAddr('l2TON');

    // Verification should fail due to incorrect L2 Ton Address
    vm.expectRevert('Provided L2 TON Token address is not correct');
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      2, // TON token type
      l2TON,
      'TestRollup'
    );

    // Verification should fail due to native token mismatch
    vm.expectRevert('The native token you are using is not TON.');
    verifier.verifyAndRegisterRollupConfig(
      address(systemConfigProxy),
      2, // TON token type
      address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000),
      'TestRollup'
    );

    vm.stopPrank();
  }

  function testVerifyL1ContractsWithSafeVerificationDisabled() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Create a different safe with wrong threshold and owners
    address[] memory safeOwners = new address[](2);
    safeOwners[0] = tokamakDAO;
    safeOwners[1] = foundation;
    // Missing thirdOwner and wrong threshold
    MockGnosisSafe differentSafe = new MockGnosisSafe(safeOwners, 1);

    // Set the owner of the proxy admin to the different safe wallet
    mockProxyAdmin.setOwner(address(differentSafe));

    // Disable safe verification
    verifier.setSafeVerificationRequired(false);

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should succeed because safe verification is disabled
    bool result = verifier.verifyL1Contracts(address(systemConfigProxy));
    assertTrue(result);

    vm.stopPrank();
  }

  function testConstructorFailZeroAddress() public {
    vm.expectRevert('Invalid proxy admin address');
    new L1ContractVerification(address(0), address(0));
  }

  function testVerifyL1ContractsFailZeroImplementation() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Mock the proxy admin to return zero address for implementation
    MockProxyAdmin mockProxyAdminZeroImpl = new MockProxyAdmin(owner);

    // Deploy a new verifier with the mock proxy admin
    L1ContractVerification newVerifier = new L1ContractVerification(
      address(mockProxyAdminZeroImpl),
      l2TONAddress
    );

    // Set up the new verifier with the same configs
    newVerifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      address(systemConfigImpl).codehash,
      address(systemConfigProxy).codehash,
      address(mockProxyAdminZeroImpl)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail because implementation is zero
    vm.expectRevert('SystemConfig verification failed');
    newVerifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyL1ContractsFailZeroAdmin() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    // Create a custom MockProxyAdmin that returns zero address for admin
    MockProxyAdmin mockProxyAdminZeroAdmin = new MockProxyAdmin(owner);

    // Deploy a new verifier with the mock proxy admin
    L1ContractVerification newVerifier = new L1ContractVerification(
      address(mockProxyAdminZeroAdmin),
      l2TONAddress
    );

    // Set up the new verifier with the same configs
    newVerifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      address(systemConfigImpl).codehash,
      address(systemConfigProxy).codehash,
      address(mockProxyAdminZeroAdmin)
    );

    // Create a custom proxy that will be used for testing
    address customProxy = makeAddr('customProxy');

    // Set implementation for the custom proxy but don't set admin
    // This will cause getProxyAdmin to revert with 'Admin not set for proxy'
    mockProxyAdminZeroAdmin.setImplementation(
      customProxy,
      address(systemConfigImpl)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail because admin check will fail
    vm.expectRevert('SystemConfig verification failed');
    newVerifier.verifyL1Contracts(customProxy);

    vm.stopPrank();
  }

  function testVerifyProxyHashFailZeroHash() public {
    vm.startPrank(owner);

    // Set all required configs but with zero proxy hash
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      address(systemConfigImpl).codehash,
      bytes32(0), // Zero proxy hash
      address(mockProxyAdmin)
    );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail because proxy hash is zero
    vm.expectRevert('Expected hash cannot be zero');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  // Helper function to set up all required configurations
  function _setupAllConfigs() internal {
    // Set contract config for SystemConfig
    verifier.setContractConfig(
      SYSTEM_CONFIG_ID,
      address(systemConfigImpl).codehash,
      address(systemConfigProxy).codehash,
      address(mockProxyAdmin)
    );

    // Set contract config for L1StandardBridge
    verifier.setContractConfig(
      L1_STANDARD_BRIDGE_ID,
      address(l1StandardBridgeImpl).codehash,
      address(l1StandardBridgeProxy).codehash,
      address(mockProxyAdmin)
    );

    // Set contract config for L1CrossDomainMessenger
    verifier.setContractConfig(
      L1_CROSS_DOMAIN_MESSENGER_ID,
      address(l1CrossDomainMessengerImpl).codehash,
      address(l1CrossDomainMessengerProxy).codehash,
      address(mockProxyAdmin)
    );

    // Set contract config for OptimismPortal
    verifier.setContractConfig(
      OPTIMISM_PORTAL_ID,
      address(optimismPortalImpl).codehash,
      address(optimismPortalProxy).codehash,
      address(mockProxyAdmin)
    );

    // Set safe config
    verifier.setSafeConfig(tokamakDAO, foundation, thirdOwner, 3);
  }
}
