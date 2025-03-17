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
contract VerifyContractSetConfig is Script {
  // Environment variables
  address systemConfigProxy;
  address proxyAdmin;
  address tokamakDAO;
  address foundation;
  address nativeToken;  // TON token address
  address bridgeRegistry;

  // Constants
  uint256 constant SAFE_THRESHOLD = 3;
  bytes32 constant IMPLEMENTATION_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

  function setUp() public {
    // Load environment variables for contract addresses
    systemConfigProxy = vm.envAddress('SYSTEM_CONFIG_PROXY');
    proxyAdmin = vm.envAddress('PROXY_ADMIN_ADDRESS');

    tokamakDAO = vm.envAddress('TOKAMAK_DAO_ADDRESS');
    foundation = vm.envAddress('FOUNDATION_ADDRESS');

    nativeToken = vm.envAddress('NATIVE_TOKEN_ADDRESS');
    bridgeRegistry = vm.envAddress('L1_BRIDGE_REGISTRY_ADDRESS');
  }

  function run() external {
    setUp();

    // Start broadcasting transactions
    vm.startBroadcast();

    // Deploy L1ContractVerification
    L1ContractVerification verifier = new L1ContractVerification(nativeToken);

    console.log("L1ContractVerification deployed at:", address(verifier));

    // Set logic contract info
    verifier.setLogicContractInfo(
      systemConfigProxy,
      proxyAdmin
    );


    IProxyAdmin proxyAdminContract = IProxyAdmin(proxyAdmin);
    address safeWalletAddress = proxyAdminContract.owner();

    // Get the implementation address of the safe wallet
    address implementation = IGnosisSafe(safeWalletAddress).masterCopy();

    // Set Safe wallet info
    verifier.setSafeWalletInfo(
      tokamakDAO,
      foundation,
      SAFE_THRESHOLD,
      proxyAdmin,
      implementation.codehash,
      safeWalletAddress.codehash
    );

    // Set bridge registry address
    verifier.setBridgeRegistryAddress(bridgeRegistry);

    console.log("L1ContractVerification configuration complete");

    // Verify and register rollup config
    bool verifyAndRegisterRollupConfigResult = verifier
      .verifyAndRegisterRollupConfig(
        systemConfigProxy,
        proxyAdmin,
        2,
        nativeToken,
        'TestRollup'
      );
    console.log(
      'Verify and register rollup config result:',
      verifyAndRegisterRollupConfigResult
    );

    vm.stopBroadcast();
  }
}