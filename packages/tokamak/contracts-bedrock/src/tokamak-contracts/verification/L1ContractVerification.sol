// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@openzeppelin/contracts/access/Ownable.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import './interface/IL1ContractVerification.sol';

contract L1ContractVerification is IL1ContractVerification, Ownable {
  // Additional structs not in the interface
  struct ProxyData {
    address implementation;
    address admin;
  }

  struct ChainConfig {
    mapping(bytes32 => ContractConfig) contractConfigs;
    SafeConfig safeConfig;
    address expectedNativeToken;
    bool safeVerificationRequired;
  }

  // Consolidated mappings
  mapping(uint256 => ChainConfig) public chainConfigs;
  mapping(address => ProxyData) private proxyData;

  // Bridge registry address
  address public L1BridgeRegistryV1_1Address;

  // Contract IDs
  bytes32 public constant L1_STANDARD_BRIDGE = keccak256('L1_STANDARD_BRIDGE');
  bytes32 public constant L1_CROSS_DOMAIN_MESSENGER = keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 public constant OPTIMISM_PORTAL = keccak256('OPTIMISM_PORTAL');
  bytes32 public constant SYSTEM_CONFIG = keccak256('SYSTEM_CONFIG');

  // Events
  event ImplementationAddressSet(address indexed proxyAddress, address indexed implementationAddress);
  event ProxyAdminSet(address indexed proxyAddress, address indexed adminAddress);
  event SafeVerificationRequiredSet(uint256 indexed chainId, bool required);

  constructor() Ownable() {}

  // External functions
  function setContractConfig(
    uint256 chainId,
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash
  ) external onlyOwner {
    _setContractConfig(
      chainId,
      contractId,
      implementationHash,
      proxyHash
    );
  }

  function setImplementationAddress(
    address proxyAddress,
    address implementationAddress
  ) external onlyOwner {
    require(proxyAddress != address(0), "Invalid proxy address");
    require(implementationAddress != address(0), "Invalid implementation address");
    proxyData[proxyAddress].implementation = implementationAddress;
    emit ImplementationAddressSet(proxyAddress, implementationAddress);
  }

  function getImplementationAddress(address proxyAddress) external view returns (address) {
    return proxyData[proxyAddress].implementation;
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
    require(tokenAddress != address(0), "Invalid token address");
    chainConfigs[chainId].expectedNativeToken = tokenAddress;
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
    emit VerificationSuccess(msg.sender, chainId);
    return true;
  }

  function verifyAndRegisterRollupConfig(
    uint256 chainId,
    address systemConfigProxy,
    address safeWallet,
    address rollupConfig,
    uint8 _type,
    address _l2TON,
    string calldata _name
  ) external returns (bool) {
    require(_type == 2, 'Registration allowed only for TON tokens');
    require(
      _verifyL1Contracts(chainId, systemConfigProxy, safeWallet),
      'Contracts verification failed'
    );

    ISystemConfig systemConfig = ISystemConfig(systemConfigProxy);
    require(
      systemConfig.nativeTokenAddress() == chainConfigs[chainId].expectedNativeToken,
      'The native token you are using is not TON.'
    );

    IL1BridgeRegistryV1_1 bridgeRegistry = IL1BridgeRegistryV1_1(
      L1BridgeRegistryV1_1Address
    );
    bridgeRegistry.registerRollupConfig(rollupConfig, _type, _l2TON, _name);

    emit RegistrationSuccess(msg.sender, chainId);
    return true;
  }

  function setBridgeRegistryAddress(
    address _bridgeRegistry
  ) external onlyOwner {
    require(_bridgeRegistry != address(0), 'Invalid bridge registry address');
    L1BridgeRegistryV1_1Address = _bridgeRegistry;
    emit BridgeRegistryUpdated(_bridgeRegistry);
  }

  function setProxyAdmin(
    address proxyAddress,
    address adminAddress
  ) external onlyOwner {
    require(proxyAddress != address(0), "Invalid proxy address");
    require(adminAddress != address(0), "Invalid admin address");
    proxyData[proxyAddress].admin = adminAddress;
    emit ProxyAdminSet(proxyAddress, adminAddress);
  }

  function getProxyAdmin(address proxyAddress) external view returns (address) {
    return proxyData[proxyAddress].admin;
  }

  function setSafeVerificationRequired(uint256 chainId, bool required) external onlyOwner {
    chainConfigs[chainId].safeVerificationRequired = required;
    emit SafeVerificationRequiredSet(chainId, required);
  }

  // Internal functions
  function _setContractConfig(
    uint256 chainId,
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash
  ) internal {
    chainConfigs[chainId].contractConfigs[contractId] = ContractConfig({
      implementationHash: implementationHash,
      proxyHash: proxyHash
    });

    emit ConfigurationSet(chainId, contractId);
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
    address admin = proxyData[proxyAddress].admin;
    if (admin == address(0)) {
      return false;
    }
    return admin == expectedAdmin;
  }

  function _verifyImplementation(
    address proxyAddress,
    bytes32 expectedHash
  ) internal view returns (bool) {
    address implementation = proxyData[proxyAddress].implementation;
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
    require(expectedHash != bytes32(0), "Expected hash cannot be zero");
    return proxyAddress.codehash == expectedHash;
  }

  function _verifyContract(
    address proxyAddress,
    bytes32 contractId,
    uint256 chainId
  ) internal view returns (bool) {
    ContractConfig storage config = chainConfigs[chainId].contractConfigs[contractId];
    return _verifyImplementation(proxyAddress, config.implementationHash) &&
           _verifyProxyAdmin(proxyAddress, proxyData[proxyAddress].admin) &&
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
    if (!_verifyContract(l1StandardBridge, L1_STANDARD_BRIDGE, chainId) ||
        !_verifyContract(l1CrossDomainMessenger, L1_CROSS_DOMAIN_MESSENGER, chainId) ||
        !_verifyContract(optimismPortal, OPTIMISM_PORTAL, chainId)) {
      return false;
    }

    // Verify Safe configuration if required
    if (chainConfigs[chainId].safeVerificationRequired && !_verifySafe(safeWallet, chainId)) {
      return false;
    }

    return true;
  }

    // Add this function to L1ContractVerification.sol
function getContractConfig(uint256 chainId, bytes32 contractId)
    external
    view
    returns (bytes32 implementationHash, bytes32 proxyHash)
{
    ContractConfig storage config = chainConfigs[chainId].contractConfigs[contractId];
    return (config.implementationHash, config.proxyHash);
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
