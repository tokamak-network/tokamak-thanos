// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import '@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol';

/**
 * @title UpgradeL1ContractVerification
 * @notice Script to upgrade the L1ContractVerification contract implementation
 * @dev This script deploys a new implementation and updates the proxy to point to it
 */
contract UpgradeL1ContractVerification is Script {
  // Environment variables - private variables with leading underscore
  address private _proxyAddress;
  address private _proxyAdminAddress;

  /**
   * @notice Load configuration from environment variables
   * @dev Called before the main execution to set up required addresses
   */
  function setUp() public {
    // Load environment variables for addresses
    _proxyAddress = vm.envAddress('L1_CONTRACT_VERIFICATION_PROXY');
    _proxyAdminAddress = vm.envAddress('L1_CONTRACT_VERIFICATION_PROXY_ADMIN');
  }

  /**
   * @notice Main execution function
   * @dev Deploys a new implementation and upgrades the proxy
   */
  function run() external {
    setUp();

    // Start broadcasting transactions
    vm.startBroadcast();

    // Deploy the new implementation contract
    L1ContractVerification newImplementation = new L1ContractVerification();
    console.log('New L1ContractVerification implementation deployed at:', address(newImplementation));

    // Get the ProxyAdmin contract
    ProxyAdmin proxyAdmin = ProxyAdmin(_proxyAdminAddress);

    // Get current implementation for logging
    // Cast proxy address to the expected TransparentUpgradeableProxy type
    TransparentUpgradeableProxy proxy = TransparentUpgradeableProxy(payable(_proxyAddress));
    address currentImplementation = proxyAdmin.getProxyImplementation(proxy);
    console.log('Current implementation:', currentImplementation);

    // Upgrade the proxy to the new implementation
    proxyAdmin.upgrade(
      proxy,
      address(newImplementation)
    );

    console.log('Proxy upgraded to new implementation');

    // Verify the upgrade was successful
    address newImplementationAddress = proxyAdmin.getProxyImplementation(proxy);
    console.log('New implementation confirmed:', newImplementationAddress);

    // If additional initialization or configuration is needed after upgrade,
    // call those functions here through the proxy

    vm.stopBroadcast();
  }
}