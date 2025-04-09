// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IGnosisSafe {
  function getThreshold() external view returns (uint256);
  function getOwners() external view returns (address[] memory);
  function masterCopy() external view returns (address);
}

interface ISystemConfig {
  function l1CrossDomainMessenger() external view returns (address);
  function l1StandardBridge() external view returns (address);
  function optimismPortal() external view returns (address);
  function nativeTokenAddress() external view returns (address);
}

interface IL1BridgeRegistry {
  function registerRollupConfig(
    address rollupConfig,
    uint8 _type,
    address _l2TON,
    string calldata _name
  ) external;

  function availableForRegistration(
    address rollupConfig,
    uint8 _type
  ) external view returns (bool);
}

interface IL1ContractVerification {
  // Struct definitions
  struct LogicContractInfo {
    address logicAddress;
    bytes32 proxyCodehash;
  }

  struct SafeWalletInfo {
    address tokamakDAO;
    address foundation;
    bytes32 implementationCodehash;
    bytes32 proxyCodehash;
    uint256 requiredThreshold;
  }

  // Events
  event ConfigurationSet(string indexed contractName);
  event VerificationSuccess(
    address indexed safeWalletAddress,
    address indexed systemConfigProxy,
    address indexed proxyAdmin,
    uint256 timestamp
  );
  event RegistrationSuccess(address indexed safeWalletAddress);
  event BridgeRegistryUpdated(address indexed bridgeRegistry);
  event SafeConfigSet(
    address indexed tokamakDAO,
    address indexed foundation,
    uint256 indexed threshold
  );
  event NativeTokenSet(address indexed tokenAddress);
  event ProxyAdminCodehashSet(bytes32 indexed codehash);
  event VerificationPossibleSet(bool indexed isVerificationPossible);
  // Functions
  function setLogicContractInfo(
    address _systemConfigProxy,
    address _proxyAdmin
  ) external;

  function setSafeConfig(
    address _tokamakDAO,
    address _foundation,
    uint256 _threshold,
    bytes32 _implementationCodehash,
    bytes32 _proxyCodehash
  ) external;

  function setBridgeRegistryAddress(address _bridgeRegistry) external;

  function verifyAndRegisterRollupConfig(
    address _systemConfigProxy,
    address _proxyAdmin,
    uint8 _type,
    address _l2TON,
    string calldata _name,
    address _safeWalletAddress
  ) external returns (bool);

  function setVerificationPossible(bool _isVerificationPossible) external;

  function verifyL1Contracts(
    address systemConfigProxy,
    address proxyAdmin,
    address safeWalletAddress
  ) external view returns (bool);
}
