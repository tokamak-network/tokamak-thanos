// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IGnosisSafe {
  function getThreshold() external view returns (uint256);
  function getOwners() external view returns (address[] memory);
}

interface ISystemConfig {
  function l1CrossDomainMessenger() external view returns (address);
  function l1StandardBridge() external view returns (address);
  function optimismPortal() external view returns (address);
}

interface IL1BridgeRegistryV1_1 {
  function registerRollupConfig(
    address rollupConfig,
    uint8 _type,
    address _l2TON,
    string calldata _name
  ) external;
}

interface IL1ContractVerification {
  // Events
  event ConfigurationSet(uint256 chainId, bytes32 contractId);
  event VerificationSuccess(address indexed operator, uint256 chainId);
  event VerificationFailure(
    address indexed operator,
    uint256 chainId,
    string reason
  );
  event RegistrationSuccess(address indexed operator, uint256 chainId);
  event BridgeRegistryUpdated(address indexed bridgeRegistry);

  // Structs
  struct ContractConfig {
    bytes32 implementationHash;
    bytes32 proxyHash;
    address expectedProxyAdmin;
  }

  struct SafeConfig {
    address tokamakDAO;
    address foundation;
    uint256 requiredThreshold;
  }

  // Functions
  function setContractConfig(
    uint256 chainId,
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash,
    address expectedProxyAdmin
  ) external;

  function setSafeConfig(
    uint256 chainId,
    address tokamakDAO,
    address foundation,
    uint256 threshold
  ) external;

  function verifyL1Contracts(
    uint256 chainId,
    address systemConfigProxy,
    address safeWallet
  ) external returns (bool);

  function verifyAndRegisterRollupConfig(
    uint256 chainId,
    address systemConfigProxy,
    address safeWallet,
    address rollupConfig,
    uint8 _type,
    address _l2TON,
    string calldata _name
  ) external returns (bool);

  function setBridgeRegistryAddress(address _bridgeRegistry) external;
}
