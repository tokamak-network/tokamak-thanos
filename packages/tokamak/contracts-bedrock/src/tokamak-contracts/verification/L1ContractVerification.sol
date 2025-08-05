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
  // Custom errors for gas optimization
  error ZeroAddress(string parameter);
  error InvalidThreshold();
  error ContractNotRegistered();
  error NoSafeWalletConfigured();
  error NativeTokenNotTON();
  error RollupConfigNotAvailable();
  error ProxyAdminInvalidCodehash();
  error InvalidProxyAdminAddress();
  error GetProxyImplementationFailed();
  error SafeWalletAddressMismatch();
  error SafeWalletInvalidProxyCodehash();
  error SafeWalletInvalidImplCodehash();
  error SafeWalletUnauthorizedModules();
  error SafeWalletInvalidFallbackHandler();
  error SafeWalletInvalidThreshold();
  error SafeWalletWrongOwnerCount();
  error SafeWalletMissingRequiredOwners();
  error SystemConfigVerificationFailed();
  error ProxyAdminNotSystemConfigAdmin();
  error ProxyAdminNotOptimismPortalAdmin();
  error ProxyAdminNotL1StandardBridgeAdmin();
  error ProxyAdminNotL1CrossDomainMessengerAdmin();
  error L1StandardBridgeVerificationFailed();
  error L1CrossDomainMessengerVerificationFailed();
  error OptimismPortalVerificationFailed();
  error ProxyAdminCodehashZero();
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

  address internal constant L2_TON_ADDRESS = address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000);

  uint8 internal constant TYPE = 2;

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
    if (_tokenAddress == address(0)) revert ZeroAddress("tokenAddress");
    if (_initialAdmin == address(0)) revert ZeroAddress("initialAdmin");

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
  function addAdmin(address _admin) external {
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
  function removeAdmin(address _admin) external {
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
    IProxyAdmin _proxyAdmin
  ) external onlyRole(ADMIN_ROLE) {
    if (_systemConfigProxy == address(0)) revert ZeroAddress("systemConfigProxy");
    if (address(_proxyAdmin) == address(0)) revert ZeroAddress("proxyAdmin");

    // Set Proxy Admin Codehash
    _setProxyAdminCodehash(address(_proxyAdmin));

    // Get contract addresses from SystemConfig
    ISystemConfig config = ISystemConfig(_systemConfigProxy);
    address l1StandardBridgeAddress = config.l1StandardBridge();
    address l1CrossDomainMessengerAddress = config.l1CrossDomainMessenger();
    address optimismPortalAddress = config.optimismPortal();

    // Set SystemConfig info
    systemConfig.logicAddress = _proxyAdmin.getProxyImplementation(
      _systemConfigProxy
    );
    systemConfig.proxyCodehash = _systemConfigProxy.codehash;

    // Set L1StandardBridge info
    l1StandardBridge.logicAddress = _proxyAdmin.getProxyImplementation(
      l1StandardBridgeAddress
    );
    l1StandardBridge.proxyCodehash = l1StandardBridgeAddress.codehash;

    // Set L1CrossDomainMessenger info
    l1CrossDomainMessenger.logicAddress = _proxyAdmin.getProxyImplementation(
      l1CrossDomainMessengerAddress
    );
    l1CrossDomainMessenger.proxyCodehash = l1CrossDomainMessengerAddress
      .codehash;

    // Set OptimismPortal info
    optimismPortal.logicAddress = _proxyAdmin.getProxyImplementation(
      optimismPortalAddress
    );
    optimismPortal.proxyCodehash = optimismPortalAddress.codehash;

    // Emit events
    emit ProxyAdminCodehashSet(proxyAdminCodehash);
    emit LogicContractConfigured('SystemConfig', _systemConfigProxy, systemConfig.logicAddress, systemConfig.proxyCodehash);
    emit LogicContractConfigured('L1StandardBridge', l1StandardBridgeAddress, l1StandardBridge.logicAddress, l1StandardBridge.proxyCodehash);
    emit LogicContractConfigured('L1CrossDomainMessenger', l1CrossDomainMessengerAddress, l1CrossDomainMessenger.logicAddress, l1CrossDomainMessenger.proxyCodehash);
    emit LogicContractConfigured('OptimismPortal', optimismPortalAddress, optimismPortal.logicAddress, optimismPortal.proxyCodehash);
  }

  /**
   * @notice Set the Safe wallet configuration
   * @param _tokamakDAO The address of the tokamakDAO owner
   * @param _foundation The address of the foundation owner
   * @param _threshold The required threshold for the safe wallet
   * @param _implementationCodehash The codehash of the implementation contract
   * @param _proxyCodehash The codehash of the proxy contract
   */
  function setSafeConfig(
    address _tokamakDAO,
    address _foundation,
    uint256 _threshold,
    bytes32 _implementationCodehash,
    bytes32 _proxyCodehash
  ) external onlyRole(ADMIN_ROLE) {
    if (_tokamakDAO == address(0)) revert ZeroAddress("tokamakDAO");
    if (_foundation == address(0)) revert ZeroAddress("foundation");
    if (_threshold != 3) revert InvalidThreshold();

    // Set common safe wallet configuration
    safeWalletConfig = SafeWalletInfo({
      tokamakDAO: _tokamakDAO,
      foundation: _foundation,
      implementationCodehash: _implementationCodehash,
      proxyCodehash: _proxyCodehash,
      requiredThreshold: _threshold,
      ownersCount: 3
    });

    emit SafeConfigSet(_tokamakDAO, _foundation, _threshold, _implementationCodehash, _proxyCodehash);
  }

  /**
   * @notice Set the bridge registry address
   * @param _bridgeRegistry The address of the bridge registry
   * @dev The bridge registry is used when registering rollup configurations
   */
  function setBridgeRegistryAddress(
    address _bridgeRegistry
  ) external onlyRole(ADMIN_ROLE) {
    if (_bridgeRegistry == address(0)) revert ZeroAddress("bridgeRegistry");
    l1BridgeRegistryAddress = _bridgeRegistry;
    emit BridgeRegistryUpdated(_bridgeRegistry);
  }

  /**
   * @notice Verify L1 contracts and register rollup configuration
   * @param _systemConfigProxy The address of the SystemConfig proxy
   * @param _proxyAdmin The address of the ProxyAdmin
   * @param _name The name of the rollup configuration
   * @param _safeWalletAddress The address of the safe wallet to verify for
   * @dev Performs verification and additionally registers the rollup with the bridge registry
   */
  function verifyAndRegisterRollupConfig(
    address _systemConfigProxy,
    IProxyAdmin _proxyAdmin,
    string calldata _name,
    address _safeWalletAddress
  ) external {
    if (!isVerificationPossible) revert ContractNotRegistered();

    // Get operator's safe wallet address
    if (_safeWalletAddress == address(0)) revert ZeroAddress("safeWalletAddress");

    // Verify proxy admin
    _verifyProxyAdmin(_proxyAdmin, _safeWalletAddress);

    // Verify L1 contracts
    _verifyL1Contracts(_systemConfigProxy, _proxyAdmin, _safeWalletAddress);

    // Emit verification success event
    emit VerificationSuccess(
      _safeWalletAddress,
      _systemConfigProxy,
      address(_proxyAdmin),
      block.timestamp
    );

        // Verify native token first
    if (ISystemConfig(_systemConfigProxy).nativeTokenAddress() != expectedNativeToken) revert NativeTokenNotTON();


    IL1BridgeRegistry bridgeRegistry = IL1BridgeRegistry(
      l1BridgeRegistryAddress
    );

    // Check if the bridge registry is available for registration
    bool isAvailable = bridgeRegistry.availableForRegistration(
      _systemConfigProxy,
      TYPE
    );
    if (!isAvailable) revert RollupConfigNotAvailable();

    bridgeRegistry.registerRollupConfig(
      _systemConfigProxy,
      TYPE,
      L2_TON_ADDRESS,
      _name
    );

    emit RegistrationSuccess(_safeWalletAddress);
  }

  /**
   * @notice Verify the ProxyAdmin contract for a specific operator's safe wallet
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @param _safeWalletAddress The safe wallet address for the operator
   * @dev Ensures the ProxyAdmin has the expected codehash and is owned by the operator's safe wallet
   */
  function _verifyProxyAdmin(
    IProxyAdmin _proxyAdmin,
    address _safeWalletAddress
  ) private view {
    // 1. Verify that ProxyAdmin contract has the expected codehash
    if (address(_proxyAdmin).codehash != proxyAdminCodehash) revert ProxyAdminInvalidCodehash();

    // 2. Verify the ProxyAdmin is owned by the operator's safe wallet
    address ownerAddress = _proxyAdmin.owner();

    if (ownerAddress != _safeWalletAddress) revert InvalidProxyAdminAddress();
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
    IProxyAdmin _proxyAdmin
  ) private view returns (bool) {
    try _proxyAdmin.getProxyImplementation(_proxyAddress) returns (
      address fetchedImpl
    ) {
      return fetchedImpl == _expectedImplementation;
    } catch {
      revert GetProxyImplementationFailed();
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
   * @dev Reverts with specific error messages if verification fails
   */
  function _verifySafe(
    IProxyAdmin _proxyAdmin,
    address _safeWalletAddress
  ) private view {
    // Get safe wallet address from ProxyAdmin.owner()
    address safeWalletAddress = _proxyAdmin.owner();

    // Check if the safe wallet address is the same as the expected safe wallet address
    if (safeWalletAddress != _safeWalletAddress) revert SafeWalletAddressMismatch();

    // Check if the proxy codehash is the same as the expected proxy codehash
    if (safeWalletAddress.codehash != safeWalletConfig.proxyCodehash) revert SafeWalletInvalidProxyCodehash();

    // Get the implementation from the masterCopy function of the safe wallet
    address implementation = IGnosisSafe(safeWalletAddress).masterCopy();

    // Check if the implementation codehash is the same as the expected implementation codehash
    if (implementation.codehash != safeWalletConfig.implementationCodehash) revert SafeWalletInvalidImplCodehash();

    // Check if the modules are not set
    if (IGnosisSafe(safeWalletAddress).getModulesPaginated(SENTINEL_MODULES, 1).length != 0) revert SafeWalletUnauthorizedModules();

    // verify fallback handler
    if (IGnosisSafe(safeWalletAddress).getFallbackHandler() != address(0)) revert SafeWalletInvalidFallbackHandler();

    // Verify threshold
    if (IGnosisSafe(safeWalletAddress).getThreshold() != safeWalletConfig.requiredThreshold) revert SafeWalletInvalidThreshold();

    // Verify owners (tokamakDAO, foundation must be included)
    address[] memory owners = IGnosisSafe(safeWalletAddress).getOwners();

    // Verify number of owners doesn't exceed the maximum
    if (owners.length != safeWalletConfig.ownersCount) revert SafeWalletWrongOwnerCount();

    bool foundTokamakDAO;
    bool foundFoundation;

    for (uint i = 0; i < owners.length; i++) {
      if (owners[i] == safeWalletConfig.tokamakDAO) foundTokamakDAO = true;
      if (owners[i] == safeWalletConfig.foundation) foundFoundation = true;
      if (foundTokamakDAO && foundFoundation) break;
    }

    // Both tokamakDAO and foundation must be present
    if (!foundTokamakDAO || !foundFoundation) revert SafeWalletMissingRequiredOwners();
  }

  /**
   * @notice Verify the L1 contracts and the safe wallet for a specific operator
   * @param _systemConfigProxy The address of the SystemConfig proxy
   * @param _proxyAdmin The address of the ProxyAdmin contract
   * @param _safeWalletAddress The safe wallet address for the operator
   */
  function _verifyL1Contracts(
    address _systemConfigProxy,
    IProxyAdmin _proxyAdmin,
    address _safeWalletAddress
  ) private view {
    // Step 1: Verify SystemConfig
    if (!_verifyImplementation(
        _systemConfigProxy,
        systemConfig.logicAddress,
        _proxyAdmin
      ) || !_verifyProxyHash(_systemConfigProxy, systemConfig.proxyCodehash)) revert SystemConfigVerificationFailed();

    // Verify proxy admin relationship for SystemConfig
    if (!_verifyProxyOwner(_systemConfigProxy, _proxyAdmin)) revert ProxyAdminNotSystemConfigAdmin();

    // Get contract addresses from SystemConfig
    address l1StandardBridgeAddress = ISystemConfig(_systemConfigProxy).l1StandardBridge();
    address l1CrossDomainMessengerAddress = ISystemConfig(_systemConfigProxy).l1CrossDomainMessenger();
    address optimismPortalAddress = ISystemConfig(_systemConfigProxy).optimismPortal();

        // Verify proxy admin relationship for OptimismPortal
    if (!_verifyProxyOwner(optimismPortalAddress, _proxyAdmin)) revert ProxyAdminNotOptimismPortalAdmin();

    // Verify proxy admin relationship for L1StandardBridge
    if (!_verifyProxyOwner(l1StandardBridgeAddress, _proxyAdmin)) revert ProxyAdminNotL1StandardBridgeAdmin();

    // Verify proxy admin relationship for L1CrossDomainMessenger
    if (!_verifyProxyOwner(l1CrossDomainMessengerAddress, _proxyAdmin)) revert ProxyAdminNotL1CrossDomainMessengerAdmin();


    // Step 2: Verify L1StandardBridge
    if (!_verifyImplementation(
        l1StandardBridgeAddress,
        l1StandardBridge.logicAddress,
        _proxyAdmin
      ) || !_verifyProxyHash(
          l1StandardBridgeAddress,
          l1StandardBridge.proxyCodehash
        )) revert L1StandardBridgeVerificationFailed();

    // Step 3: Verify L1CrossDomainMessenger
    if (!_verifyImplementation(
        l1CrossDomainMessengerAddress,
        l1CrossDomainMessenger.logicAddress,
        _proxyAdmin
      ) || !_verifyProxyHash(
          l1CrossDomainMessengerAddress,
          l1CrossDomainMessenger.proxyCodehash
        )) revert L1CrossDomainMessengerVerificationFailed();

    // Step 4: Verify OptimismPortal
    if (!_verifyImplementation(
        optimismPortalAddress,
        optimismPortal.logicAddress,
        _proxyAdmin
      ) || !_verifyProxyHash(optimismPortalAddress, optimismPortal.proxyCodehash)) revert OptimismPortalVerificationFailed();

    // Step 5: Verify Safe wallet
    _verifySafe(_proxyAdmin, _safeWalletAddress);
  }

  function _setProxyAdminCodehash(address _proxyAdmin) private {
    if (_proxyAdmin.codehash == bytes32(0)) revert ProxyAdminCodehashZero();
    proxyAdminCodehash = _proxyAdmin.codehash;
  }

    /**
     * @notice Verifies that the provided proxy admin is the actual admin of a proxy contract
     * @param _proxyAddress The address of the proxy contract to check
     * @param _proxyAdmin The address of the claimed proxy admin
     * @return Returns true if the _proxyAdmin is the admin of _proxyAddress
     * @dev CRITICAL SECURITY DEPENDENCY: This function MUST only be called after _verifyProxyAdmin
     *      has been executed to validate the ProxyAdmin contract's codehash and ownership.
     */
    function _verifyProxyOwner(address _proxyAddress, IProxyAdmin _proxyAdmin) private view returns (bool) {
        // SECURITY: This function assumes _proxyAdmin has been validated by _verifyProxyAdmin
        address actualAdmin = _proxyAdmin.getProxyAdmin(payable(_proxyAddress));
        return actualAdmin == address(_proxyAdmin);
    }
}
