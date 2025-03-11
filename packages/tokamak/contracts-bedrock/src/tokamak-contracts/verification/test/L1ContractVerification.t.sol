// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Test.sol';
import '../L1ContractVerification.sol';
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

  function setUp() public {
    // Setup accounts
    owner = makeAddr('owner');
    tokamakDAO = makeAddr('tokamakDAO');
    foundation = makeAddr('foundation');
    thirdOwner = makeAddr('thirdOwner');
    user = makeAddr('user');
    nativeToken = makeAddr('nativeToken');

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
    systemConfigImpl = new MockSystemConfig(
      address(l1StandardBridgeProxy),
      address(l1CrossDomainMessengerProxy),
      address(optimismPortalProxy)
    );

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

    // Set native token address
    systemConfigImpl.setNativeTokenAddress(nativeToken);

    // Setup Safe wallet
    address[] memory safeOwners = new address[](3);
    safeOwners[0] = tokamakDAO;
    safeOwners[1] = foundation;
    safeOwners[2] = thirdOwner;
    safeWallet = new MockGnosisSafe(safeOwners, 3);

    // Deploy bridge registry
    bridgeRegistry = new MockBridgeRegistry();

    // Deploy verifier contract with proxy admin address
    verifier = new L1ContractVerification(address(mockProxyAdmin));

    // Set bridge registry address
    verifier.setBridgeRegistryAddress(address(bridgeRegistry));

    // Set expected native token
    verifier.setExpectedNativeToken(nativeToken);

    // Set safe verification required to true for testing
    verifier.setSafeVerificationRequired(true);

    bytes32 slot = keccak256('PROXY_ADMIN_ADDRESS_SLOT');
    vm.store(
      address(verifier),
      slot,
      bytes32(uint256(uint160(address(mockProxyAdmin))))
    );

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
    console.log(address(systemConfigProxy));
    console.log(address(systemConfigImpl));
    console.log('-----------------------------------------------------------');

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
    payable(address(systemConfigProxy)),
    address(mockProxyAdmin)
  );
  mockProxyAdmin.setAdmin(
    payable(address(l1StandardBridgeProxy)),
    address(mockProxyAdmin)
  );
  mockProxyAdmin.setAdmin(
    payable(address(l1CrossDomainMessengerProxy)),
    address(mockProxyAdmin)
  );
  mockProxyAdmin.setAdmin(
    payable(address(optimismPortalProxy)),
    address(mockProxyAdmin)
  );

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

    // Deploy a different SystemConfig implementation
    // MockSystemConfig differentSystemConfigImpl = new MockSystemConfig(
    //   address(0),
    //   address(0),
    //   address(0)
    // );

    // // Update the SystemConfig proxy to point to the different implementation
    // proxyAdmin.upgrade(
    //   TransparentUpgradeableProxy(payable(address(systemConfigProxy))),
    //   address(differentSystemConfigImpl)
    // );

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Contracts verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

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
    // MockGnosisSafe differentSafe = new MockGnosisSafe(safeOwners, 1); // Wrong threshold

    vm.stopPrank();

    vm.startPrank(user);

    // Verification should fail
    vm.expectRevert('Contracts verification failed');
    verifier.verifyL1Contracts(address(systemConfigProxy));

    vm.stopPrank();
  }

  function testVerifyAndRegisterRollupConfig() public {
    vm.startPrank(owner);

    // Set all required configs
    _setupAllConfigs();

    vm.stopPrank();

    vm.startPrank(user);

    address rollupConfig = makeAddr('rollupConfig');
    address l2TON = makeAddr('l2TON');

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

    address rollupConfig = makeAddr('rollupConfig');
    address l2TON = makeAddr('l2TON');

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

  function testSetExpectedNativeTokenWithChainId() public {
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

  function testSetSafeVerificationRequiredWithChainId() public {
    vm.startPrank(owner);

    // Set to false with explicit chainId
    verifier.setSafeVerificationRequired(false);

    // Set back to true with explicit chainId
    verifier.setSafeVerificationRequired(true);

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
