// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import '@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol';

/**
 * @title VerifyContractSetConfig
 * @notice Script to deploy and configure the L1ContractVerification contract as upgradeable
 * @dev This script is used by the TRH team to deploy and set up the verification contract
 *      before releasing the trh-sdk
 */
contract VerifyContractSetConfig is Script {
  // Environment variables - private variables with leading underscore
  address private _systemConfigProxy;
  address private _proxyAdmin;
  address private _tokamakDAO;
  address private _foundation;
  address private _nativeToken; // TON token address
  address private _bridgeRegistry;

  // Constants - private constants with leading underscore
  uint256 private constant _SAFE_THRESHOLD = 3;

  function setUp() public {
    // Load environment variables for contract addresses
    _systemConfigProxy = vm.envAddress('SYSTEM_CONFIG_PROXY');
    _proxyAdmin = vm.envAddress('PROXY_ADMIN_ADDRESS');

    _tokamakDAO = vm.envAddress('TOKAMAK_DAO_ADDRESS');
    _foundation = vm.envAddress('FOUNDATION_ADDRESS');

    _nativeToken = vm.envAddress('NATIVE_TOKEN_ADDRESS');
    _bridgeRegistry = vm.envAddress('L1_BRIDGE_REGISTRY_ADDRESS');
  }

  function run() external {
    setUp();

    // Start broadcasting transactions
    vm.startBroadcast();

    // First deploy the implementation contract
    L1ContractVerification implementationContract = new L1ContractVerification();
    console.log('L1ContractVerification implementation deployed at:', address(implementationContract));

    // Deploy a ProxyAdmin for managing the proxy
    ProxyAdmin proxyAdmin = new ProxyAdmin();
    console.log('ProxyAdmin deployed at:', address(proxyAdmin));

    // Prepare initialization data for the proxy
    bytes memory initData = abi.encodeWithSelector(
      L1ContractVerification.initialize.selector,
      _nativeToken
    );

    // Deploy the TransparentUpgradeableProxy with the implementation and initialization data
    TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
      address(implementationContract),
      address(proxyAdmin),
      initData
    );
    console.log('L1ContractVerification proxy deployed at:', address(proxy));

    // Create a reference to the proxy contract to call its functions
    L1ContractVerification verifier = L1ContractVerification(address(proxy));

    // Set logic contract info
    verifier.setLogicContractInfo(_systemConfigProxy, _proxyAdmin);

    // Get the safe wallet address from the proxy admin
    IProxyAdmin proxyAdminContract = IProxyAdmin(_proxyAdmin);
    address safeWalletAddress = proxyAdminContract.owner();

    // Get the implementation address of the safe wallet using masterCopy function
    address implementation = IGnosisSafe(safeWalletAddress).masterCopy();

    // Set Safe wallet info
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

    // // Verify and register rollup config
    // bool verifyAndRegisterRollupConfigResult = verifier
    //   .verifyAndRegisterRollupConfig(
    //     _systemConfigProxy,
    //     _proxyAdmin,
    //     2,
    //     address(0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000), //L2 Ton address
    //     'TestRollup'
    //   );
    // console.log(
    //   'Verify and register rollup config result:',
    //   verifyAndRegisterRollupConfigResult
    // );

    vm.stopBroadcast();
  }
}
