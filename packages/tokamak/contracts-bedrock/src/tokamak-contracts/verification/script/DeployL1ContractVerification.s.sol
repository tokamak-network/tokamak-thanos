// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import '../L1ContractVerification.sol';

contract DeployL1ContractVerification is Script {
  function run() external {
    // Start broadcasting transactions
    vm.startBroadcast();

    // Deploy the contract
    L1ContractVerification verificationContract = new L1ContractVerification();

    console.log(
      'L1ContractVerification deployed to:',
      address(verificationContract)
    );

    // End broadcasting transactions
    vm.stopBroadcast();
  }
}
