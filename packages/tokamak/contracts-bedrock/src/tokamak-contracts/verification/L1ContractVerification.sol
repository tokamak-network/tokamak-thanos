// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol';
import '@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol';
import './interface/IL1ContractVerification.sol';
import './interface/IProxyAdmin.sol';

/**
 * @title L1ContractVerification
 * @notice This contract verifies the integrity of critical L1 contracts for Tokamak rollups
 * @dev This contract is designed as an upgradeable contract that provides security
 *      guarantees for Tokamak Layer 2 operators. It ensures that the L1 contracts
 *      they interact with have the correct implementations and configurations.
 *
 * @custom:security-model The security model of this contract follows these principles:
 *      1. Proxy Pattern Verification: Validates both proxy and implementation contracts
 *      2. Ownership Verification: Ensures ownership is held by the expected multisig wallets
 *      3. Address Registry: Acts as a trusted registry of verified contract addresses
 *      4. Trust Minimization: Reduces trust required in L2 operators by verifying their setups
 *
 * @custom:upgrade-safety This contract uses the TransparentUpgradeableProxy pattern for upgrades.
 *      Upgrades should be carefully vetted to ensure storage layout compatibility.
 */
contract L1ContractVerification is
  IL1ContractVerification,
  Initializable,
  AccessControlUpgradeable
{
  // Role definitions
  /**
   * @notice Admin role for managing configuration and performing operations
   * @dev Has access to all contract functions
   */
  bytes32 public constant ADMIN_ROLE = keccak256('ADMIN_ROLE');

  address internal constant SENTINEL_MODULES = address(0x1);

  // The expected native token (TON) address
  address public expectedNativeToken;

  // Bridge registry address
  address public l1BridgeRegistryAddress;

  // The codehash of the ProxyAdmin contract
  bytes32 public proxyAdminCodehash;

  // Flag to control if verification is possible
  bool public isVerificationPossible;

  // Logic contract information storage
  LogicContractInfo public systemConfig;
  LogicContractInfo public l1StandardBridge;
  LogicContractInfo public l1CrossDomainMessenger;
  LogicContractInfo public optimismPortal;

  // Common safe wallet configuration
  SafeWalletInfo public safeWalletConfig;

  /**
   * @custom:oz-upgrades-unsafe-allow constructor
   */
  constructor() {
    _disableInitializers();
  }

  /**
   * @notice Initialize the contract (replaces constructor for upgradeable contracts)
   * @param _tokenAddress The address of the native token (TON)
   * @param _initialAdmin The address that will be granted the admin role
   */
  function initialize(
    address _tokenAddress,
    address _initialAdmin
  ) public initializer {
    __AccessControl_init();

    // Set up roles
    _setupRole(DEFAULT_ADMIN_ROLE, _initialAdmin);
    _setupRole(ADMIN_ROLE, _initialAdmin);

    expectedNativeToken = _tokenAddress;
    isVerificationPossible = false;
    emit NativeTokenSet(_tokenAddress);
  }

  /**
   * @notice Add an admin
   * @param _admin The address to grant the admin role to
   * @dev Only callable by existing admins
   */
  function addAdmin(address _admin) external onlyRole(ADMIN_ROLE) {
    grantRole(ADMIN_ROLE, _admin);
  }

  /**
   * @notice Set whether verification is possible
   * @param _isVerificationPossible Boolean flag to enable/disable verification
   * @dev Only callable by admins
   */
  function setVerificationPossible(
    bool _isVerificationPossible
  ) external onlyRole(ADMIN_ROLE) {
    isVerificationPossible = _isVerificationPossible;
    emit VerificationPossibleSet(_isVerificationPossible);
  }

  /**
   * @notice Remove an admin
   * @param _admin The address to revoke the admin role from
   * @dev Only callable by existing admins
   */
  function removeAdmin(address _admin) external onlyRole(ADMIN_ROLE) {
    revokeRole(ADMIN_ROLE, _admin);
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
  ) external onlyRole(ADMIN_ROLE) {
    require(
      _systemConfigProxy != address(0),
      'SystemConfig proxy address cannot be zero'
    );
    require(_proxyAdmin != address(0), 'ProxyAdmin address cannot be zero');

    // Set Proxy Admin Codehash
    _setProxyAdminCodehash(_proxyAdmin);

    // Get contract addresses from SystemConfig
    ISystemConfig config = ISystemConfig(_systemConfigProxy);
    address l1StandardBridgeAddress = config.l1StandardBridge();
    address l1CrossDomainMessengerAddress = config.l1CrossDomainMessenger();
    address optimismPortalAddress = config.optimismPortal();

    IProxyAdmin proxyAdmin = IProxyAdmin(_proxyAdmin);

    // Set SystemConfig info
    systemConfig.logicAddress = proxyAdmin.getProxyImplementation(
      _systemConfigProxy
    );
    systemConfig.proxyCodehash = _systemConfigProxy.codehash;

    // Set L1StandardBridge info
    l1StandardBridge.logicAddress = proxyAdmin.getProxyImplementation(
      l1StandardBridgeAddress
    );
    l1StandardBridge.proxyCodehash = l1StandardBridgeAddress.codehash;

    // Set L1CrossDomainMessenger info
    l1CrossDomainMessenger.logicAddress = proxyAdmin.getProxyImplementation(
      l1CrossDomainMessengerAddress
    );
    l1CrossDomainMessenger.proxyCodehash = l1CrossDomainMessengerAddress
      .codehash;

    // Set OptimismPortal info
    optimismPortal.logicAddress = proxyAdmin.getProxyImplementation(
      optimismPortalAddress
    );
    optimismPortal.proxyCodehash = optimismPortalAddress.codehash;

    // Emit events
    emit ProxyAdminCodehashSet(proxyAdminCodehash);
    emit ConfigurationSet('SystemConfig');
    emit ConfigurationSet('L1StandardBridge');
    emit ConfigurationSet('L1CrossDomainMessenger');
    emit ConfigurationSet('OptimismPortal');
  }

  /**
   * @notice Set the Safe wallet configuration
   * @param _tokamakDAO The address of the tokamakDAO owner
   * @param _foundation The address of the foundation owner
   * @param _threshold The required threshold for the safe wallet
   * @param _implementationCodehash The codehash of the implementation contract
   * @param _proxyCodehash The codehash of the proxy contract
   * @param _ownersCount The number of owners allowed for the safe wallet
   */
  function setSafeConfig(
    address _tokamakDAO,
    address _foundation,
    uint256 _threshold,
    bytes32 _implementationCodehash,
    bytes32 _proxyCodehash,
    uint256 _ownersCount
  ) external onlyRole(ADMIN_ROLE) {
    require(_tokamakDAO != address(0), 'TokamakDAO address cannot be zero');
    require(_foundation != address(0), 'Foundation address cannot be zero');
    require(_threshold > 0, 'Threshold must be greater than zero');
    require(_ownersCount == 3, 'Owners must be 3');

    // Set common safe wallet configuration
    safeWalletConfig = SafeWalletInfo({
      tokamakDAO: _tokamakDAO,
      foundation: _foundation,
      implementationCodehash: _implementationCodehash,
      proxyCodehash: _proxyCodehash,
      requiredThreshold: _threshold,
      ownersCount: _ownersCount
    });

    emit SafeConfigSet(_tokamakDAO, _foundation, _threshold);
  }

  /**
   * @notice Set the bridge registry address
   * @param _bridgeRegistry The address of the bridge registry
   * @dev The bridge registry is used when registering rollup configurations
   */
  function setBridgeRegistryAddress(
    address _bridgeRegistry
  ) external onlyRole(ADMIN_ROLE) {
    require(
      _bridgeRegistry != address(0),
      'Bridge registry address cannot be zero'
    );
    l1BridgeRegistryAddress = _bridgeRegistry;
    emit BridgeRegistryUpdated(_bridgeRegistry);
  }

  /**
   * @notice Verify L1 contracts for a specific safe wallet
   * @param systemConfigProxy The address of the SystemConfig proxy
   * @param proxyAdmin The address of the ProxyAdmin
   * @param safeWalletAddress The address of the safe wallet to verify
   * @return Returns true if verification succeeds, otherwise reverts
   */
  function verifyL1Contracts(
    address systemConfigProxy,
    address proxyAdmin,
    address safeWalletAddress
  ) external view returns (bool) {
    require(isVerificationPossible, 'Contract not registered as registerant');
    ISystemConfig systemConfigContract = ISystemConfig(systemConfigProxy);
    require(
      systemConfigContract.nativeTokenAddress() == expectedNativeToken,
      'The native token you are using is not TON'
    );

    require(
      safeWalletAddress != address(0),
      'Safe wallet address cannot be zero'
    );

    // Additional security check: Ensure no delegate calls are allowed
    require(
        msg.sender.code.length == 0,
       'Calls from contract accounts not allowed for verification'
    );

    // Verify proxy admin
    _verifyProxyAdmin(proxyAdmin, safeWalletAddress);

        // Verify L1 contracts
    _verifyL1Contracts(systemConfigProxy, proxyAdmin, safeWalletAddress);

    return true;
  }

  /**
   * @notice Verify L1 contracts and register rollup configuration
   * @param _systemConfigProxy The address of the SystemConfig proxy
   * @param _proxyAdmin The address of the ProxyAdmin
   * @param _type Token type (must be 2 for TON)
   * @param _l2TON The address of the L2 TON token
   * @param _name The name of the rollup configuration
   * @param _safeWalletAddress The address of the safe wallet to verify for
   * @return Returns true if verification and registration succeeds, otherwise reverts
   * @dev Performs verification and additionally registers the rollup with the bridge registry
   */
  function verifyAndRegisterRollupConfig(
    address _systemConfigProxy,
    address _proxyAdmin,
    uint8 _type,
    address _l2TON,
    string calldata _name,
    address _safeWalletAddress
  ) external returns (bool) {
    require(isVerificationPossible, 'Contract not registered as registerant');
    require(
      _l2TON == address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000),
      'Provided L2 TON Token address is not correct'
    );
    require(_type == 2, 'Registration allowed only for TON tokens');

    // Additional security check: Ensure no delegate calls are allowed
    require(
        msg.sender.code.length == 0,
       'Calls from contract accounts not allowed for verification'
    );

    // Verify native token first
    ISystemConfig systemConfigContract = ISystemConfig(_systemConfigProxy);
    require(
      systemConfigContract.nativeTokenAddress() == expectedNativeToken,
      'The native token you are using is not TON'
    );

    // Get operator's safe wallet address
    require(
      _safeWalletAddress != address(0),
      'No safe wallet configured for operator'
    );

    // Verify proxy admin
    _verifyProxyAdmin(_proxyAdmin, _safeWalletAddress);

    // Verify L1 contracts
    _verifyL1Contracts(_systemConfigProxy, _proxyAdmin, _safeWalletAddress);

    // Emit verification success event
    emit VerificationSuccess(
      _safeWalletAddress,
      _systemConfigProxy,
      _proxyAdmin,
      block.timestamp
    );

    IL1BridgeRegistry bridgeRegistry = IL1BridgeRegistry(
      l1BridgeRegistryAddress
    );

    // Check if the bridge registry is available for registration
    bool isAvailable = bridgeRegistry.availableForRegistration(
      _systemConfigProxy,
      _type
    );
    require(
      isAvailable,
      'Rollup configuration already registered or not available for registration'
    );

    bridgeRegistry.registerRollupConfig(
      _systemConfigProxy,
      _type,
      _l2TON,
      _name
    );

    emit RegistrationSuccess(_safeWalletAddress);
    return true;
  }

  /**
   * @notice Verify the ProxyAdmin contract for a specific operator's safe wallet
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @param _safeWalletAddress The safe wallet address for the operator
   * @dev Ensures the ProxyAdmin has the expected codehash and is owned by the operator's safe wallet
   */
  function _verifyProxyAdmin(
    address _proxyAdmin,
    address _safeWalletAddress
  ) private view {
    // 1. Verify that ProxyAdmin contract has the expected codehash
    require(
      _proxyAdmin.codehash == proxyAdminCodehash,
      'ProxyAdmin verification failed: invalid codehash'
    );

    // 2. Verify the ProxyAdmin is owned by the operator's safe wallet
    IProxyAdmin proxyAdminContract = IProxyAdmin(_proxyAdmin);
    address ownerAddress = proxyAdminContract.owner();

    require(
      ownerAddress == _safeWalletAddress,
      'Invalid proxy admin address'
    );
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
  ) private view returns (bool) {
    IProxyAdmin proxyAdmin = IProxyAdmin(_proxyAdmin);

    try proxyAdmin.getProxyImplementation(_proxyAddress) returns (
      address fetchedImpl
    ) {
      return fetchedImpl == _expectedImplementation;
    } catch {
      revert("Failed to call getProxyImplementation on ProxyAdmin");
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
  ) private view returns (bool) {
    return _proxyAddress.codehash == _expectedHash;
  }

  /**
   * @notice Verify the Safe wallet for a specific operator
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @param _safeWalletAddress The safe wallet address for the operator
   * @return Returns true if verification succeeds, otherwise reverts with a specific error message
   */
  function _verifySafe(
    address _proxyAdmin,
    address _safeWalletAddress
  ) private view returns (bool) {
    // Get safe wallet address from ProxyAdmin.owner()
    IProxyAdmin proxyAdminContract = IProxyAdmin(_proxyAdmin);
    address safeWalletAddress = proxyAdminContract.owner();

    // Check if the safe wallet address is the same as the expected safe wallet address
    require(
      safeWalletAddress == _safeWalletAddress,
      'Safe wallet verification failed: address mismatch'
    );

    IGnosisSafe safe = IGnosisSafe(safeWalletAddress);

    // Get the implementation from the masterCopy function of the safe wallet
    address implementation = safe.masterCopy();

    // Check if the implementation codehash is the same as the expected implementation codehash
    require(
      implementation.codehash == safeWalletConfig.implementationCodehash,
      'Safe wallet verification failed: invalid implementation codehash'
    );

    // Check if the modules are not set
    require(safe.getModulesPaginated(SENTINEL_MODULES, 100).length == 0, "Safe wallet verification failed: unauthorized modules");

    // Check if the proxy codehash is the same as the expected proxy codehash
    require(
      safeWalletAddress.codehash == safeWalletConfig.proxyCodehash,
      'Safe wallet verification failed: invalid proxy codehash'
    );

    // verify fallback handler
    require(safe.getFallbackHandler() == address(0), "Safe wallet verification failed: fallback handler should be set to ZERO address");

    // Verify threshold
    require(
      safe.getThreshold() == safeWalletConfig.requiredThreshold,
      'Safe wallet verification failed: invalid threshold'
    );

    // Verify owners (tokamakDAO, foundation must be included)
    address[] memory owners = safe.getOwners();

    // Verify number of owners doesn't exceed the maximum
    require(
      owners.length == safeWalletConfig.ownersCount,
      'Safe wallet verification failed: wrong number of owners'
    );

    bool foundTokamakDAO = false;
    bool foundFoundation = false;

    for (uint i = 0; i < owners.length; i++) {
      if (owners[i] == safeWalletConfig.tokamakDAO) foundTokamakDAO = true;
      if (owners[i] == safeWalletConfig.foundation) foundFoundation = true;
      if (foundTokamakDAO && foundFoundation) break;
    }

    // Both tokamakDAO and foundation must be present
    require(
      foundTokamakDAO && foundFoundation,
      'Safe wallet verification failed: missing required owners'
    );

    return true;
  }

  /**
   * @notice Verify the L1 contracts and the safe wallet for a specific operator
   * @param _systemConfigProxy The address of the SystemConfig proxy
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @param _safeWalletAddress The safe wallet address for the operator
   */
  function _verifyL1Contracts(
    address _systemConfigProxy,
    address _proxyAdmin,
    address _safeWalletAddress
  ) private view {
    // Step 1: Verify SystemConfig
    require(
      _verifyImplementation(
        _systemConfigProxy,
        systemConfig.logicAddress,
        _proxyAdmin
      ) && _verifyProxyHash(_systemConfigProxy, systemConfig.proxyCodehash),
      'SystemConfig verification failed'
    );

    // Get contract addresses from SystemConfig
    ISystemConfig systemConfigContract = ISystemConfig(_systemConfigProxy);
    address l1StandardBridgeAddress = systemConfigContract.l1StandardBridge();
    address l1CrossDomainMessengerAddress = systemConfigContract
      .l1CrossDomainMessenger();
    address optimismPortalAddress = systemConfigContract.optimismPortal();

        // Verify proxy admin relationship for OptimismPortal
    require(
        _verifyProxyOwner(optimismPortalAddress, _proxyAdmin),
        'ProxyAdmin is not the admin of OptimismPortal'
    );

    // Verify proxy admin relationship for L1StandardBridge
    require(
        _verifyProxyOwner(l1StandardBridgeAddress, _proxyAdmin),
        'ProxyAdmin is not the admin of L1StandardBridge'
    );

    // Verify proxy admin relationship for L1CrossDomainMessenger
    require(
        _verifyProxyOwner(l1CrossDomainMessengerAddress, _proxyAdmin),
        'ProxyAdmin is not the admin of L1CrossDomainMessenger'
    );


    // Step 2: Verify L1StandardBridge
    require(
      _verifyImplementation(
        l1StandardBridgeAddress,
        l1StandardBridge.logicAddress,
        _proxyAdmin
      ) &&
        _verifyProxyHash(
          l1StandardBridgeAddress,
          l1StandardBridge.proxyCodehash
        ),
      'L1StandardBridge verification failed'
    );

    // Step 3: Verify L1CrossDomainMessenger
    require(
      _verifyImplementation(
        l1CrossDomainMessengerAddress,
        l1CrossDomainMessenger.logicAddress,
        _proxyAdmin
      ) &&
        _verifyProxyHash(
          l1CrossDomainMessengerAddress,
          l1CrossDomainMessenger.proxyCodehash
        ),
      'L1CrossDomainMessenger verification failed'
    );

    // Step 4: Verify OptimismPortal
    require(
      _verifyImplementation(
        optimismPortalAddress,
        optimismPortal.logicAddress,
        _proxyAdmin
      ) &&
        _verifyProxyHash(optimismPortalAddress, optimismPortal.proxyCodehash),
      'OptimismPortal verification failed'
    );

    // Step 5: Verify Safe wallet
    require(_verifySafe(_proxyAdmin, _safeWalletAddress), 'Safe wallet verification failed');
  }

  function _setProxyAdminCodehash(address _proxyAdmin) private {
    bytes32 proxyAdminHash = _proxyAdmin.codehash;
    require(proxyAdminHash != bytes32(0), 'ProxyAdmin codehash cannot be zero');
    proxyAdminCodehash = proxyAdminHash;
  }

    /**
    * @notice Verifies that the provided proxy admin is the actual admin of a proxy contract
    * @param _proxyAddress The address of the proxy contract to check
    * @param _proxyAdmin The address of the claimed proxy admin
    * @return Returns true if the _proxyAdmin is the admin of _proxyAddress
    */
    function _verifyProxyOwner(address _proxyAddress, address _proxyAdmin) private view returns (bool) {

        IProxyAdmin proxyAdminContract = IProxyAdmin(_proxyAdmin);

        address actualAdmin = proxyAdminContract.getProxyAdmin(payable(_proxyAddress));
        return actualAdmin == _proxyAdmin;
    }
}