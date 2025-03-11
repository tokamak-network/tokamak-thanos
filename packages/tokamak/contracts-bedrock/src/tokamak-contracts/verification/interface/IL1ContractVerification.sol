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
    address expectedProxyAdmin;
  }

  struct SafeConfig {
    address tokamakDAO;
    address foundation;
    uint256 requiredThreshold;
  }

  // Events
  event ConfigurationSet(bytes32 indexed contractId);
  event VerificationSuccess(address indexed verifier);
  event VerificationFailure(
    address indexed operator,
    string reason
  );
  event RegistrationSuccess(address indexed verifier);
  event BridgeRegistryUpdated(address indexed bridgeRegistry);

  // Functions
  function setContractConfig(
    bytes32 contractId,
    bytes32 implementationHash,
    bytes32 proxyHash,
    address expectedProxyAdmin
  ) external;

  function setSafeConfig(
    address tokamakDAO,
    address foundation,
    uint256 threshold
  ) external;

  function verifyL1Contracts(
    address systemConfigProxy,
    address safeWallet
  ) external returns (bool);

  function verifyAndRegisterRollupConfig(
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
   * @return tokamakDAO The address of the TokamakDAO owner
   * @return foundation The address of the Foundation owner
   * @return requiredThreshold The required threshold for the safe
   */
  function getSafeConfig()
    external
    view
    returns (address tokamakDAO, address foundation, uint256 requiredThreshold);

  /**
   * @notice Get the contract configuration for a specific chain and contract ID
   * @param contractId The contract ID to get the config for
   * @return implementationHash The hash of the implementation contract code
   * @return proxyHash The hash of the proxy contract code
   * @return expectedProxyAdmin The expected admin address for the proxy
   */
  function getContractConfig(bytes32 contractId)
    external
    view
    returns (bytes32 implementationHash, bytes32 proxyHash, address expectedProxyAdmin);

  /**
   * @notice Set whether safe verification is required for a specific chain
   * @param chainId The chain ID to set the requirement for
   * @param required Whether safe verification is required
   */
  function setSafeVerificationRequired(uint256 chainId, bool required) external;

  /**
   * @notice Set the expected native token for a specific chain
   * @param chainId The chain ID to set the token for
   * @param tokenAddress The address of the expected native token
   */
  function setExpectedNativeToken(uint256 chainId, address tokenAddress) external;

  /**
   * @notice Get the ProxyAdmin contract address
   * @return The address of the ProxyAdmin contract
   */
  function PROXY_ADMIN_ADDRESS() external view returns (address);
}
