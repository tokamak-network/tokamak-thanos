// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@openzeppelin/contracts/access/Ownable.sol';
import './interface/IL1ContractVerification.sol';
import './interface/IProxyAdmin.sol';
import "forge-std/console.sol";

/**
 * @title L1ContractVerification
 * @notice This contract verifies the integrity of critical L1 contracts for Tokamak rollups
 * @dev This contract is deployed to L1 by the TRH team before the trh-sdk release and helps
 *      ensure that L2 operators deploy contracts with the correct implementations
 */
contract L1ContractVerification is IL1ContractVerification, Ownable {

  // The expected native token (TON) address
  address public immutable expectedNativeToken;

  // Bridge registry address
  address public l1BridgeRegistryAddress;

  // The codehash of the ProxyAdmin contract
  bytes32 public proxyAdminCodehash;

  // Logic contract information storage
  LogicContractInfo public systemConfig;
  LogicContractInfo public l1StandardBridge;
  LogicContractInfo public l1CrossDomainMessenger;
  LogicContractInfo public optimismPortal;

  // Safe wallet information storage
  SafeWalletInfo public safeWallet;

  /**
   * @notice Constructor
   * @param _tokenAddress The address of the native token (TON)
   */
  constructor(address _tokenAddress) Ownable() {
        expectedNativeToken = _tokenAddress;
  }

  /**
   * @notice Set all logic contract info in one call using a deployed SystemConfig contract
   * @param _systemConfigProxy The address of the SystemConfig proxy
   * @param _proxyAdmin The address of the ProxyAdmin
   * @dev This function records implementation addresses and codehashes for all key contracts
   */
  function setLogicContractInfo(
    address _systemConfigProxy,
    address _proxyAdmin
  ) external onlyOwner {
    require(_systemConfigProxy != address(0), 'SystemConfig proxy address cannot be zero');
    require(_proxyAdmin != address(0), 'ProxyAdmin address cannot be zero');

    // Verify ProxyAdmin
    bytes32 proxyAdminHash = _proxyAdmin.codehash;
    require(proxyAdminHash != bytes32(0), 'ProxyAdmin codehash cannot be zero');
    proxyAdminCodehash = proxyAdminHash;

    // Get contract addresses from SystemConfig
    ISystemConfig config = ISystemConfig(_systemConfigProxy);
    address l1StandardBridgeAddress = config.l1StandardBridge();
    address l1CrossDomainMessengerAddress = config.l1CrossDomainMessenger();
    address optimismPortalAddress = config.optimismPortal();

    IProxyAdmin proxyAdmin = IProxyAdmin(_proxyAdmin);

    // Set SystemConfig info
    systemConfig.logicAddress = proxyAdmin.getProxyImplementation(_systemConfigProxy);
    systemConfig.proxyCodehash = _systemConfigProxy.codehash;

    // Set L1StandardBridge info
    l1StandardBridge.logicAddress = proxyAdmin.getProxyImplementation(l1StandardBridgeAddress);
    l1StandardBridge.proxyCodehash = l1StandardBridgeAddress.codehash;

    // Set L1CrossDomainMessenger info
    l1CrossDomainMessenger.logicAddress = proxyAdmin.getProxyImplementation(l1CrossDomainMessengerAddress);
    l1CrossDomainMessenger.proxyCodehash = l1CrossDomainMessengerAddress.codehash;

    // Set OptimismPortal info
    optimismPortal.logicAddress = proxyAdmin.getProxyImplementation(optimismPortalAddress);
    optimismPortal.proxyCodehash = optimismPortalAddress.codehash;

    // Emit events
    emit ProxyAdminCodehashSet(proxyAdminHash);
    emit ConfigurationSet("SystemConfig");
    emit ConfigurationSet("L1StandardBridge");
    emit ConfigurationSet("L1CrossDomainMessenger");
    emit ConfigurationSet("OptimismPortal");
  }

  /**
   * @notice Set the Safe wallet configuration
   * @param _tokamakDAO The address of the tokamakDAO owner
   * @param _foundation The address of the foundation owner
   * @param _threshold The required threshold for the safe wallet
   * @param _proxyAdmin The address of the ProxyAdmin
   * @param _implementationCodehash The codehash of the implementation contract
   * @param _proxyCodehash The codehash of the proxy contract
   * @dev Records information about the Gnosis Safe wallet that owns the ProxyAdmin
   */
  function setSafeWalletInfo(
    address _tokamakDAO,
    address _foundation,
    uint256 _threshold,
    address _proxyAdmin,
    bytes32 _implementationCodehash,
    bytes32 _proxyCodehash
  ) external onlyOwner {
    require(_tokamakDAO != address(0), 'TokamakDAO address cannot be zero');
    require(_foundation != address(0), 'Foundation address cannot be zero');
    require(_threshold > 0, 'Threshold must be greater than zero');

    IProxyAdmin proxyAdmin = IProxyAdmin(_proxyAdmin);
    address safeWalletAddress = proxyAdmin.owner();

    safeWallet = SafeWalletInfo({
      safeWalletAddress: safeWalletAddress,
      tokamakDAO: _tokamakDAO,
      foundation: _foundation,
      implementationCodehash: _implementationCodehash,
      proxyCodehash: _proxyCodehash,
      requiredThreshold: _threshold
    });

    emit SafeConfigSet(_tokamakDAO, _foundation, _threshold);
  }

  /**
   * @notice Set the bridge registry address
   * @param _bridgeRegistry The address of the bridge registry
   * @dev The bridge registry is used when registering rollup configurations
   */
  function setBridgeRegistryAddress(address _bridgeRegistry) external onlyOwner {
    require(_bridgeRegistry != address(0), 'Bridge registry address cannot be zero');
    l1BridgeRegistryAddress = _bridgeRegistry;
    emit BridgeRegistryUpdated(_bridgeRegistry);
  }

  /**
   * @notice Verify L1 contracts
   * @param systemConfigProxy The address of the SystemConfig proxy
   * @param proxyAdmin The address of the ProxyAdmin
   * @return Returns true if verification succeeds, otherwise reverts
   * @dev First checks that the native token matches expected value, then verifies
   *      the ProxyAdmin and all contract implementations
   */
  function verifyL1Contracts(
    address systemConfigProxy,
    address proxyAdmin
  ) external returns (bool) {
    // Verify native token first
    ISystemConfig systemConfigContract = ISystemConfig(systemConfigProxy);
    require(
      systemConfigContract.nativeTokenAddress() == expectedNativeToken,
      'The native token you are using is not TON'
    );

    // Verify proxy admin
    _verifyProxyAdmin(proxyAdmin);

    // Verify L1 contracts
    _verifyL1Contracts(systemConfigProxy, proxyAdmin);

    emit VerificationSuccess(msg.sender);
    return true;
  }

  /**
   * @notice Verify L1 contracts and register rollup configuration
   * @param _systemConfigProxy The address of the SystemConfig proxy
   * @param _proxyAdmin The address of the ProxyAdmin
   * @param _type Token type (must be 2 for TON)
   * @param _l2TON The address of the L2 TON token
   * @param _name The name of the rollup configuration
   * @return Returns true if verification and registration succeeds, otherwise reverts
   * @dev Performs verification and additionally registers the rollup with the bridge registry
   */
  function verifyAndRegisterRollupConfig(
    address _systemConfigProxy,
    address _proxyAdmin,
    uint8 _type,
    address _l2TON,
    string calldata _name
  ) external returns (bool) {
    require(
      _l2TON == expectedNativeToken,
      'Provided L2 TON Token address is not correct'
    );
    require(_type == 2, 'Registration allowed only for TON tokens');

    // Verify native token first
    ISystemConfig systemConfigContract = ISystemConfig(_systemConfigProxy);
    require(
      systemConfigContract.nativeTokenAddress() == expectedNativeToken,
      'The native token you are using is not TON'
    );

    // Verify proxy admin
    _verifyProxyAdmin(_proxyAdmin);

    // Verify L1 contracts
    _verifyL1Contracts(_systemConfigProxy, _proxyAdmin);

    // Register rollup configuration
    IL1BridgeRegistry bridgeRegistry = IL1BridgeRegistry(
      l1BridgeRegistryAddress
    );
    bridgeRegistry.registerRollupConfig(
      _systemConfigProxy,
      _type,
      _l2TON,
      _name
    );

    emit RegistrationSuccess(msg.sender);
    return true;
  }

  /**
   * @notice Verify the ProxyAdmin contract
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @dev Ensures the ProxyAdmin has the expected codehash
   */
  function _verifyProxyAdmin(address _proxyAdmin) internal view {
    // Verify that ProxyAdmin contract has the expected codehash
    require(_proxyAdmin.codehash == proxyAdminCodehash, 'ProxyAdmin verification failed');
  }

  /**
   * @notice Verify the implementation of a proxy contract
   * @param _proxyAddress The address of the proxy contract
   * @param _expectedImplementation The expected implementation address
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @return Returns true if verification succeeds, otherwise false
   * @dev Uses the ProxyAdmin to get the implementation of the proxy and compares it
   */
  function _verifyImplementation(
    address _proxyAddress,
    address _expectedImplementation,
    address _proxyAdmin
  ) internal view returns (bool) {
    IProxyAdmin proxyAdmin = IProxyAdmin(_proxyAdmin);

    try proxyAdmin.getProxyImplementation(_proxyAddress) returns (
      address fetchedImpl
    ) {
      return fetchedImpl == _expectedImplementation;
    } catch {
      return false;
    }
  }

  /**
   * @notice Verify the codehash of a proxy contract
   * @param _proxyAddress The address of the proxy contract
   * @param _expectedHash The expected codehash
   * @return Returns true if verification succeeds, otherwise false
   * @dev Compares the codehash of the proxy contract with the expected hash
   */
  function _verifyProxyHash(
    address _proxyAddress,
    bytes32 _expectedHash
  ) internal view returns (bool) {
    return _proxyAddress.codehash == _expectedHash;
  }

  /**
   * @notice Verify the Safe wallet
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @return Returns true if verification succeeds, otherwise false
   * @dev Verifies that the safe wallet has the correct address, implementation,
   *      proxy codehash, threshold, and includes both tokamakDAO and foundation as owners
   */
  function _verifySafe(address _proxyAdmin) internal view returns (bool) {
    // Get safe wallet address from ProxyAdmin.owner()
    IProxyAdmin proxyAdminContract = IProxyAdmin(_proxyAdmin);
    address safeWalletAddress = proxyAdminContract.owner();
    IGnosisSafe safe = IGnosisSafe(safeWalletAddress);

    // Check if the safe wallet address is the same as the expected safe wallet address
    if (safeWalletAddress != safeWallet.safeWalletAddress) {
      return false;
    }

    // Get the implementation from the masterCopy function of the safe wallet
    address implementation = safe.masterCopy();

    // Check if the implementation codehash is the same as the expected implementation codehash
    if (implementation.codehash != safeWallet.implementationCodehash) {
      return false;
    }

    // Check if the proxy codehash is the same as the expected proxy codehash
    if (safeWalletAddress.codehash != safeWallet.proxyCodehash) {
      return false;
    }

    // Verify threshold
    if (safe.getThreshold() != safeWallet.requiredThreshold) {
      return false;
    }

    // Verify owners (tokamakDAO, foundation must be included)
    address[] memory owners = safe.getOwners();

    bool foundTokamakDAO = false;
    bool foundFoundation = false;

    for (uint i = 0; i < owners.length; i++) {
      if (owners[i] == safeWallet.tokamakDAO) foundTokamakDAO = true;
      if (owners[i] == safeWallet.foundation) foundFoundation = true;
      if (foundTokamakDAO && foundFoundation) break;
    }

    // Both tokamakDAO and foundation must be present
    return foundTokamakDAO && foundFoundation;
  }

  /**
   * @notice Verify the L1 contracts and the safe wallet
   * @param _systemConfigProxy The address of the SystemConfig proxy
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @dev Verifies each contract in sequence: SystemConfig, L1StandardBridge,
   *      L1CrossDomainMessenger, OptimismPortal, and the Safe wallet
   */
  function _verifyL1Contracts(
    address _systemConfigProxy,
    address _proxyAdmin
  ) internal view {
    // Step 1: Verify SystemConfig
    require(
      _verifyImplementation(_systemConfigProxy, systemConfig.logicAddress, _proxyAdmin) &&
      _verifyProxyHash(_systemConfigProxy, systemConfig.proxyCodehash),
      'SystemConfig verification failed'
    );

    // Get contract addresses from SystemConfig
    ISystemConfig systemConfigContract = ISystemConfig(_systemConfigProxy);
    address l1StandardBridgeAddress = systemConfigContract.l1StandardBridge();
    address l1CrossDomainMessengerAddress = systemConfigContract.l1CrossDomainMessenger();
    address optimismPortalAddress = systemConfigContract.optimismPortal();

    // Step 2: Verify L1StandardBridge
    require(
      _verifyImplementation(l1StandardBridgeAddress, l1StandardBridge.logicAddress, _proxyAdmin) &&
      _verifyProxyHash(l1StandardBridgeAddress, l1StandardBridge.proxyCodehash),
      'L1StandardBridge verification failed'
    );

    // Step 3: Verify L1CrossDomainMessenger
    require(
      _verifyImplementation(l1CrossDomainMessengerAddress, l1CrossDomainMessenger.logicAddress, _proxyAdmin) &&
      _verifyProxyHash(l1CrossDomainMessengerAddress, l1CrossDomainMessenger.proxyCodehash),
      'L1CrossDomainMessenger verification failed'
    );

    // Step 4: Verify OptimismPortal
    require(
      _verifyImplementation(optimismPortalAddress, optimismPortal.logicAddress, _proxyAdmin) &&
      _verifyProxyHash(optimismPortalAddress, optimismPortal.proxyCodehash),
      'OptimismPortal verification failed'
    );

    // Step 5: Verify Safe wallet
    require(_verifySafe(_proxyAdmin), 'Safe wallet verification failed');
  }
}
