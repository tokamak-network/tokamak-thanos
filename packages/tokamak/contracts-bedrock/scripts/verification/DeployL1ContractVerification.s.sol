// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';

contract DeployL1ContractVerification is Script {
  function run() external {
    // Load environment variables
    address proxyAdminAddress = vm.envAddress("PROXY_ADMIN_CONTRACT_ADDRESS");

    // Start broadcasting transactions
    vm.startBroadcast();

    // Deploy the contract with the ProxyAdmin address
    L1ContractVerification verificationContract = new L1ContractVerification(proxyAdminAddress);

    console.log(
      'L1ContractVerification deployed to:',
      address(verificationContract)
    );

    // End broadcasting transactions
    vm.stopBroadcast();
  }
}
