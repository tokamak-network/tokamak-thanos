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
  function nativeTokenAddress() external view returns (address);
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
  // Struct definitions
  struct ContractConfig {
    bytes32 implementationHash;
    bytes32 proxyHash;
  }

  struct SafeConfig {
    address tokamakDAO;
    address foundation;
    uint256 requiredThreshold;
  }

  // Events
  event ConfigurationSet(uint256 indexed chainId, bytes32 indexed contractId);
  event VerificationSuccess(address indexed verifier, uint256 indexed chainId);
  event VerificationFailure(
    address indexed operator,
    uint256 chainId,
    string reason
  );
  event RegistrationSuccess(address indexed verifier, uint256 indexed chainId);
  event BridgeRegistryUpdated(address indexed bridgeRegistry);

  // Functions
  function setContractConfig(
    uint256 chainId,
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash
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
    returns (address tokamakDAO, address foundation, uint256 requiredThreshold);

  /**
   * @notice Get the contract configuration for a specific chain and contract ID
   * @param chainId The chain ID to get the contract config for
   * @param contractId The contract ID to get the config for
   * @return implementationHash The hash of the implementation contract code
   * @return proxyHash The hash of the proxy contract code
   */
  function getContractConfig(uint256 chainId, bytes32 contractId)
    external
    view
    returns (bytes32 implementationHash, bytes32 proxyHash);
}
