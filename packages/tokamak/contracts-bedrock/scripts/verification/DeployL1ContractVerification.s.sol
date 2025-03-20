// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';

/**
 * @title SetupL1ContractVerification
 * @notice Script to deploy and configure the L1ContractVerification contract
 * @dev This script is used by the TRH team to deploy and set up the verification contract
 *      before releasing the trh-sdk
 */
contract SetupL1ContractVerification is Script {
  // Environment variables - private variables with leading underscore
  address private _systemConfigProxy;
  address private _proxyAdmin;
  address private _tokamakDAO;
  address private _foundation;
  address private _nativeToken; // TON token address
  address private _bridgeRegistry;

  // Constants - private constants with leading underscore
  uint256 private constant _SAFE_THRESHOLD = 3;

  /**
   * @notice Load configuration from environment variables
   * @dev Called before the main execution to set up required addresses
   */
  function setUp() public {
    // Load environment variables for contract addresses
    _systemConfigProxy = vm.envAddress('SYSTEM_CONFIG_PROXY');
    _proxyAdmin = vm.envAddress('PROXY_ADMIN_ADDRESS');

    _tokamakDAO = vm.envAddress('TOKAMAK_DAO_ADDRESS');
    _foundation = vm.envAddress('FOUNDATION_ADDRESS');

    _nativeToken = vm.envAddress('NATIVE_TOKEN_ADDRESS');
    _bridgeRegistry = vm.envAddress('L1_BRIDGE_REGISTRY_ADDRESS');
  }

  /**
   * @notice Main execution function
   * @dev Deploys and configures the L1ContractVerification contract
   */
  function run() external {
    setUp();

    // Start broadcasting transactions
    vm.startBroadcast();

    // Deploy L1ContractVerification
    L1ContractVerification verifier = new L1ContractVerification(_nativeToken);

    console.log('L1ContractVerification deployed at:', address(verifier));

    // Set logic contract info
    verifier.setLogicContractInfo(_systemConfigProxy, _proxyAdmin);

    // Get the safe wallet address from the proxy admin
    IProxyAdmin proxyAdminContract = IProxyAdmin(_proxyAdmin);
    address safeWalletAddress = proxyAdminContract.owner();

    // Get the implementation address of the safe wallet using masterCopy function
    address implementation = IGnosisSafe(safeWalletAddress).masterCopy();

    // Set Safe wallet info with the implementation address and codehash
    verifier.setSafeConfig(
      _tokamakDAO,
      _foundation,
      _SAFE_THRESHOLD,
      _proxyAdmin,
      implementation.codehash,
      safeWalletAddress.codehash
    );

    // Set bridge registry address
    verifier.setBridgeRegistryAddress(_bridgeRegistry);

    console.log('L1ContractVerification configuration complete');
    vm.stopBroadcast();
  }
}
