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
}

interface IL1ContractVerification {
  // Struct definitions
  struct LogicContractInfo {
    address logicAddress;
    bytes32 proxyCodehash;
  }

  struct SafeWalletInfo {
    address safeWalletAddress;
    address tokamakDAO;
    address foundation;
    bytes32 implementationCodehash;
    bytes32 proxyCodehash;
    uint256 requiredThreshold;
  }

  // Events
  event ConfigurationSet(string indexed contractName);
  event VerificationSuccess(address indexed verifier);
  event RegistrationSuccess(address indexed verifier);
  event BridgeRegistryUpdated(address indexed bridgeRegistry);
  event SafeConfigSet(address tokamakDAO, address foundation, uint256 threshold);
  event NativeTokenSet(address indexed tokenAddress);
  event ProxyAdminCodehashSet(bytes32 codehash);

  // Functions
  function setLogicContractInfo(
    address _systemConfigProxy,
    address _proxyAdmin
  ) external;

  function setSafeWalletInfo(
    address _tokamakDAO,
    address _foundation,
    uint256 _threshold,
    address _proxyAdmin,
    bytes32 _implementationCodehash,
    bytes32 _proxyCodehash
  ) external;

  function setBridgeRegistryAddress(address _bridgeRegistry) external;

  function verifyAndRegisterRollupConfig(
    address _systemConfigProxy,
    address _proxyAdmin,
    uint8 _type,
    address _l2TON,
    string calldata _name
  ) external returns (bool);

}
