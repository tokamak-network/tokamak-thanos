// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@openzeppelin/contracts/access/Ownable.sol';
import './interface/IL1ContractVerification.sol';
import './interface/IProxyAdmin.sol';

contract L1ContractVerification is IL1ContractVerification, Ownable {
  // ProxyAdmin contract address
  address public immutable PROXY_ADMIN_ADDRESS;

  struct ChainConfig {
    mapping(bytes32 => ContractConfig) contractConfigs;
    SafeConfig safeConfig;
    address expectedNativeToken;
    bool safeVerificationRequired;
  }

  // Consolidated mappings
  mapping(uint256 => ChainConfig) public chainConfigs;

  // Bridge registry address
  address public L1BridgeRegistryV1_1Address;

  // Contract IDs
  bytes32 public constant L1_STANDARD_BRIDGE = keccak256('L1_STANDARD_BRIDGE');
  bytes32 public constant L1_CROSS_DOMAIN_MESSENGER =
    keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 public constant OPTIMISM_PORTAL = keccak256('OPTIMISM_PORTAL');
  bytes32 public constant SYSTEM_CONFIG = keccak256('SYSTEM_CONFIG');

  // Events
  event SafeVerificationRequiredSet(uint256 indexed chainId, bool required);

  constructor(address _proxyAdminAddress) Ownable() {
    require(_proxyAdminAddress != address(0), "Invalid proxy admin address");
    PROXY_ADMIN_ADDRESS = _proxyAdminAddress;
  }

  // External functions
  function setContractConfig(
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash,
    address expectedProxyAdmin
  ) external onlyOwner {
    _setContractConfig(
      block.chainid,
      contractId,
      implementationHash,
      proxyHash,
      expectedProxyAdmin
    );
  }

  function setSafeConfig(
    address tokamakDAO,
    address foundation,
    uint256 threshold
  ) external onlyOwner {
    _setSafeConfig(block.chainid, tokamakDAO, foundation, threshold);
  }

  function setExpectedNativeToken(address tokenAddress) external onlyOwner {
    require(tokenAddress != address(0), 'Invalid token address');
    chainConfigs[block.chainid].expectedNativeToken = tokenAddress;
  }

  function verifyL1Contracts(
    address systemConfigProxy,
    address safeWallet
  ) external returns (bool) {
    require(
      _verifyL1Contracts(block.chainid, systemConfigProxy, safeWallet),
      'Contracts verification failed'
    );
    emit VerificationSuccess(msg.sender);
    return true;
  }

  function verifyAndRegisterRollupConfig(
    address systemConfigProxy,
    address safeWallet,
    address rollupConfig,
    uint8 _type,
    address _l2TON,
    string calldata _name
  ) external returns (bool) {
    require(_type == 2, 'Registration allowed only for TON tokens');
    require(
      _verifyL1Contracts(block.chainid, systemConfigProxy, safeWallet),
      'Contracts verification failed'
    );

    ISystemConfig systemConfig = ISystemConfig(systemConfigProxy);
    require(
      systemConfig.nativeTokenAddress() ==
        chainConfigs[block.chainid].expectedNativeToken,
      'The native token you are using is not TON.'
    );

    IL1BridgeRegistryV1_1 bridgeRegistry = IL1BridgeRegistryV1_1(
      L1BridgeRegistryV1_1Address
    );
    bridgeRegistry.registerRollupConfig(rollupConfig, _type, _l2TON, _name);

    emit RegistrationSuccess(msg.sender);
    return true;
  }

  function setBridgeRegistryAddress(
    address _bridgeRegistry
  ) external onlyOwner {
    require(_bridgeRegistry != address(0), 'Invalid bridge registry address');
    L1BridgeRegistryV1_1Address = _bridgeRegistry;
    emit BridgeRegistryUpdated(_bridgeRegistry);
  }

  function setSafeVerificationRequired(bool required) external onlyOwner {
    chainConfigs[block.chainid].safeVerificationRequired = required;
    emit SafeVerificationRequiredSet(block.chainid, required);
  }

  // Implementation of interface functions with chainId parameter
  function setContractConfig(
    uint256 chainId,
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash,
    address expectedProxyAdmin
  ) external onlyOwner {
    _setContractConfig(
      chainId,
      contractId,
      implementationHash,
      proxyHash,
      expectedProxyAdmin
    );
  }

  function setSafeConfig(
    uint256 chainId,
    address tokamakDAO,
    address foundation,
    uint256 threshold
  ) external onlyOwner {
    _setSafeConfig(chainId, tokamakDAO, foundation, threshold);
  }

  function setExpectedNativeToken(
    uint256 chainId,
    address tokenAddress
  ) external onlyOwner {
    require(tokenAddress != address(0), 'Invalid token address');
    chainConfigs[chainId].expectedNativeToken = tokenAddress;
  }

  function setSafeVerificationRequired(
    uint256 chainId,
    bool required
  ) external onlyOwner {
    chainConfigs[chainId].safeVerificationRequired = required;
    emit SafeVerificationRequiredSet(chainId, required);
  }

  function verifyL1Contracts(
    uint256 chainId,
    address systemConfigProxy,
    address safeWallet
  ) external returns (bool) {
    require(
      _verifyL1Contracts(chainId, systemConfigProxy, safeWallet),
      'Contracts verification failed'
    );
    emit VerificationSuccess(msg.sender);
    return true;
  }

  /**
   * @notice Get the contract configuration for the current chain and contract ID
   * @param contractId The contract ID to get the config for
   * @return implementationHash The hash of the implementation contract code
   * @return proxyHash The hash of the proxy contract code
   * @return expectedProxyAdmin The expected admin address for the proxy
   */
  function getContractConfig(bytes32 contractId)
    external
    view
    returns (bytes32 implementationHash, bytes32 proxyHash, address expectedProxyAdmin)
  {
    ContractConfig storage config = chainConfigs[block.chainid].contractConfigs[contractId];
    return (config.implementationHash, config.proxyHash, config.expectedProxyAdmin);
  }

  /**
   * @notice Get the safe configuration for the current chain
   * @return tokamakDAO The address of the TokamakDAO owner
   * @return foundation The address of the Foundation owner
   * @return requiredThreshold The required threshold for the safe
   */
  function getSafeConfig()
    external
    view
    returns (address tokamakDAO, address foundation, uint256 requiredThreshold)
  {
    SafeConfig storage config = chainConfigs[block.chainid].safeConfig;
    return (config.tokamakDAO, config.foundation, config.requiredThreshold);
  }

  // Internal functions
  function _setContractConfig(
    uint256 chainId,
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash,
    address expectedProxyAdmin
  ) internal {
    chainConfigs[chainId].contractConfigs[contractId] = ContractConfig({
      implementationHash: implementationHash,
      proxyHash: proxyHash,
      expectedProxyAdmin: expectedProxyAdmin
    });

    emit ConfigurationSet(contractId);
  }

  function _setSafeConfig(
    uint256 chainId,
    address tokamakDAO,
    address foundation,
    uint256 threshold
  ) internal {
    chainConfigs[chainId].safeConfig = SafeConfig({
      tokamakDAO: tokamakDAO,
      foundation: foundation,
      requiredThreshold: threshold
    });
  }

  function _verifyProxyAdmin(
    address proxyAddress,
    address expectedAdmin
  ) internal view returns (bool) {
    IProxyAdmin proxyAdmin = IProxyAdmin(PROXY_ADMIN_ADDRESS);
    address admin;

    // Convert to payable address for compatibility with IProxyAdmin interface
    address payable payableProxyAddress = payable(proxyAddress);

    try proxyAdmin.getProxyAdmin(payableProxyAddress) returns (address fetchedAdmin) {
      admin = fetchedAdmin;
    } catch {
      return false;
    }

    if (admin == address(0)) {
      return false;
    }

    return admin == expectedAdmin;
  }

  function _verifyImplementation(
    address proxyAddress,
    bytes32 expectedHash
  ) internal view returns (bool) {
    IProxyAdmin proxyAdmin = IProxyAdmin(PROXY_ADMIN_ADDRESS);
    address implementation;

    try proxyAdmin.getProxyImplementation(proxyAddress) returns (address fetchedImpl) {
      implementation = fetchedImpl;
    } catch {
      return false;
    }

    if (implementation == address(0)) {
      return false;
    }

    return implementation.codehash == expectedHash;
  }

  function _verifySafe(
    address safeWallet,
    uint256 chainId
  ) internal view returns (bool) {
    IGnosisSafe safe = IGnosisSafe(safeWallet);
    SafeConfig storage safeConfig = chainConfigs[chainId].safeConfig;

    if (safe.getThreshold() != safeConfig.requiredThreshold) {
      return false;
    }

    address[] memory owners = safe.getOwners();
    bool foundTokamakDAO = false;
    bool foundFoundation = false;

    for (uint i = 0; i < owners.length; i++) {
      if (owners[i] == safeConfig.tokamakDAO) foundTokamakDAO = true;
      if (owners[i] == safeConfig.foundation) foundFoundation = true;
      if (foundTokamakDAO && foundFoundation) break; // Exit early if both found
    }

    return foundTokamakDAO && foundFoundation;
  }

  function _verifyProxyHash(
    address proxyAddress,
    bytes32 expectedHash
  ) internal view returns (bool) {
    require(expectedHash != bytes32(0), 'Expected hash cannot be zero');
    return proxyAddress.codehash == expectedHash;
  }

  function _verifyContract(
    address proxyAddress,
    bytes32 contractId,
    uint256 chainId
  ) internal view returns (bool) {
    ContractConfig storage config = chainConfigs[chainId].contractConfigs[contractId];

    return _verifyImplementation(proxyAddress, config.implementationHash) &&
           _verifyProxyAdmin(proxyAddress, config.expectedProxyAdmin) &&
           _verifyProxyHash(proxyAddress, config.proxyHash);
  }

  function _verifyL1Contracts(
    uint256 chainId,
    address systemConfigProxy,
    address safeWallet
  ) internal view returns (bool) {
    // Verify SystemConfig
    if (!_verifyContract(systemConfigProxy, SYSTEM_CONFIG, chainId)) {
      return false;
    }

    // Get contract addresses from SystemConfig
    ISystemConfig systemConfig = ISystemConfig(systemConfigProxy);
    address l1StandardBridge = systemConfig.l1StandardBridge();
    address l1CrossDomainMessenger = systemConfig.l1CrossDomainMessenger();
    address optimismPortal = systemConfig.optimismPortal();

    // Verify other contracts
    if (
      !_verifyContract(l1StandardBridge, L1_STANDARD_BRIDGE, chainId) ||
      !_verifyContract(
        l1CrossDomainMessenger,
        L1_CROSS_DOMAIN_MESSENGER,
        chainId
      ) ||
      !_verifyContract(optimismPortal, OPTIMISM_PORTAL, chainId)
    ) {
      return false;
    }

    // Verify Safe configuration if required
    if (
      chainConfigs[chainId].safeVerificationRequired &&
      !_verifySafe(safeWallet, chainId)
    ) {
      return false;
    }

    return true;
  }

  /**
   * @notice Get the contract configuration for a specific chain and contract ID
   * @param chainId The chain ID to get the contract config for
   * @param contractId The contract ID to get the config for
   * @return implementationHash The hash of the implementation contract code
   * @return proxyHash The hash of the proxy contract code
   * @return expectedProxyAdmin The expected admin address for the proxy
   */
  function getContractConfig(uint256 chainId, bytes32 contractId)
    external
    view
    returns (bytes32 implementationHash, bytes32 proxyHash, address expectedProxyAdmin)
  {
    ContractConfig storage config = chainConfigs[chainId].contractConfigs[contractId];
    return (config.implementationHash, config.proxyHash, config.expectedProxyAdmin);
  }

  /**
   * @notice Get the safe configuration for a specific chain
   * @param chainId The chain ID to get the safe config for
   * @return tokamakDAO The address of the TokamakDAO owner
   * @return foundation The address of the Foundation owner
   * @return requiredThreshold The required threshold for the safe
   */
  function getSafeConfig(uint256 chainId)
    external
    view
    returns (address tokamakDAO, address foundation, uint256 requiredThreshold)
  {
    SafeConfig storage config = chainConfigs[chainId].safeConfig;
    return (config.tokamakDAO, config.foundation, config.requiredThreshold);
  }
}
