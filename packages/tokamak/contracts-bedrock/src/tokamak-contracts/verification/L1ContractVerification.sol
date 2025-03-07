// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@openzeppelin/contracts/access/Ownable.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import './interface/IL1ContractVerification.sol';

contract L1ContractVerification is IL1ContractVerification, Ownable {
  // State variables
  mapping(uint256 => mapping(bytes32 => ContractConfig)) public contractConfigs;
  mapping(uint256 => SafeConfig) public safeConfigs;
  mapping(uint256 => address) public expectedNativeToken;

  //Bridge registry address
  address public L1BridgeRegistryV1_1Address;
  // Contract IDs
  bytes32 public constant L1_STANDARD_BRIDGE = keccak256('L1_STANDARD_BRIDGE');
  bytes32 public constant L1_CROSS_DOMAIN_MESSENGER =
    keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 public constant OPTIMISM_PORTAL = keccak256('OPTIMISM_PORTAL');
  bytes32 public constant SYSTEM_CONFIG = keccak256('SYSTEM_CONFIG');

  constructor() Ownable() {}

  // External functions
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

  // Setter of native token according to network
  function setExpectedNativeToken(
    uint256 chainId,
    address tokenAddress
  ) external onlyOwner {
    require(tokenAddress != address(0), "Invalid token address");
    expectedNativeToken[chainId] = tokenAddress;
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
    // First verify L1 contracts
    // TODO: Check condition to allow only TON tokens, if it's correct
    require(_type == 2, 'Registration allowed only for TON tokens');
    require(
      _verifyL1Contracts(chainId, systemConfigProxy, safeWallet),
      'Contracts verification failed'
    );

    // Check whether L2 operator use TON as a native token
    ISystemConfig systemConfig = ISystemConfig(systemConfigProxy);
    require(
      systemConfig.nativeTokenAddress() == expectedNativeToken[chainId],
      'The native token you are using is not TON.'
    );

    // Then register the rollup configuration
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

  // Internal functions
  function _setContractConfig(
    uint256 chainId,
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash,
    address expectedProxyAdmin
  ) internal {
    contractConfigs[chainId][contractId] = ContractConfig({
      implementationHash: implementationHash,
      proxyHash: proxyHash,
      expectedProxyAdmin: expectedProxyAdmin
    });

    emit ConfigurationSet(chainId, contractId);
  }

  function _setSafeConfig(
    uint256 chainId,
    address tokamakDAO,
    address foundation,
    uint256 threshold
  ) internal {
    safeConfigs[chainId] = SafeConfig({
      tokamakDAO: tokamakDAO,
      foundation: foundation,
      requiredThreshold: threshold
    });
  }

  function _verifySystemConfig(
    address systemConfigProxy,
    uint256 chainId
  ) internal returns (bool) {
    return
      _verifyImplementation(
        systemConfigProxy,
        contractConfigs[chainId][SYSTEM_CONFIG].implementationHash
      ) &&
      _verifyProxyAdmin(
        systemConfigProxy,
        contractConfigs[chainId][SYSTEM_CONFIG].expectedProxyAdmin
      );
  }

  function _verifyProxyAdmin(
    address proxyAddress,
    address expectedAdmin
  ) internal returns (bool) {
    TransparentUpgradeableProxy proxy = TransparentUpgradeableProxy(
      payable(proxyAddress)
    );

    try proxy.admin() returns (address admin) {
      return admin == expectedAdmin;
    } catch {
      return false;
    }
  }

  function _verifyImplementation(
    address proxyAddress,
    bytes32 expectedHash
  ) internal returns (bool) {
    TransparentUpgradeableProxy proxy = TransparentUpgradeableProxy(
      payable(proxyAddress)
    );

    try proxy.implementation() returns (address implementation) {
      return implementation.codehash == expectedHash;
    } catch {
      return false;
    }
  }

  function _verifySafe(
    address safeWallet,
    uint256 chainId
  ) internal view returns (bool) {
    IGnosisSafe safe = IGnosisSafe(safeWallet);
    SafeConfig storage safeConfig = safeConfigs[chainId];

    if (safe.getThreshold() != safeConfig.requiredThreshold) {
      return false;
    }

    address[] memory owners = safe.getOwners();
    bool foundTokamakDAO = false;
    bool foundFoundation = false;

    for (uint i = 0; i < owners.length; i++) {
      if (owners[i] == safeConfig.tokamakDAO) foundTokamakDAO = true;
      if (owners[i] == safeConfig.foundation) foundFoundation = true;
    }

    return foundTokamakDAO && foundFoundation;
  }

  function _verifyL1Contracts(
    uint256 chainId,
    address systemConfigProxy,
    address safeWallet
  ) internal returns (bool) {
    require(
      _verifySystemConfig(systemConfigProxy, chainId),
      'Invalid SystemConfig'
    );

    // Get contract addresses from SystemConfig
    ISystemConfig systemConfig = ISystemConfig(systemConfigProxy);
    address l1StandardBridge = systemConfig.l1StandardBridge();
    address l1CrossDomainMessenger = systemConfig.l1CrossDomainMessenger();
    address optimismPortal = systemConfig.optimismPortal();

    // Verify implementations and proxy admins
    require(
      _verifyImplementation(
        l1StandardBridge,
        contractConfigs[chainId][L1_STANDARD_BRIDGE].implementationHash
      ),
      'Invalid L1StandardBridge implementation'
    );
    require(
      _verifyProxyAdmin(
        l1StandardBridge,
        contractConfigs[chainId][L1_STANDARD_BRIDGE].expectedProxyAdmin
      ),
      'Invalid L1StandardBridge proxy admin'
    );

    require(
      _verifyImplementation(
        l1CrossDomainMessenger,
        contractConfigs[chainId][L1_CROSS_DOMAIN_MESSENGER].implementationHash
      ),
      'Invalid L1CrossDomainMessenger implementation'
    );
    require(
      _verifyProxyAdmin(
        l1CrossDomainMessenger,
        contractConfigs[chainId][L1_CROSS_DOMAIN_MESSENGER].expectedProxyAdmin
      ),
      'Invalid L1CrossDomainMessenger proxy admin'
    );

    require(
      _verifyImplementation(
        optimismPortal,
        contractConfigs[chainId][OPTIMISM_PORTAL].implementationHash
      ),
      'Invalid OptimismPortal implementation'
    );
    require(
      _verifyProxyAdmin(
        optimismPortal,
        contractConfigs[chainId][OPTIMISM_PORTAL].expectedProxyAdmin
      ),
      'Invalid OptimismPortal proxy admin'
    );

    // Verify Safe configuration
    require(_verifySafe(safeWallet, chainId), 'Invalid Safe configuration');
    return true;
  }
}
