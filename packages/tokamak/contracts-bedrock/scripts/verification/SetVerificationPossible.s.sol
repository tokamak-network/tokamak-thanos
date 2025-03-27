// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';

/**
 * @title SetVerificationPossible
 * @notice Script to update the verification possible flag in the L1ContractVerification contract
 * @dev This script is used to enable or disable the verification functionality
 */
contract SetVerificationPossible is Script {
  // Address of the L1ContractVerification contract
  address public verificationContract;

  // New value for the isVerificationPossible flag
  bool public newValue;

  function setUp() public {
    // Get the verification contract address from environment variable
    verificationContract = vm.envAddress('VERIFICATION_CONTRACT');

    // Get the new value from environment variable (true/false)
    string memory newValueStr = vm.envString('VERIFICATION_POSSIBLE');
    newValue =
      keccak256(abi.encodePacked(newValueStr)) ==
      keccak256(abi.encodePacked('true'));
  }

  function run() public {
    // Ensure the verification contract address is set
    require(
      verificationContract != address(0),
      'Verification contract address not set'
    );

    // Log the current action
    console.log('Setting isVerificationPossible to:', newValue);
    console.log('Verification contract:', verificationContract);

    L1ContractVerification verifier = L1ContractVerification(
      verificationContract
    );

    // Start the broadcast
    vm.startBroadcast();

    // Call the setVerificationPossible function
    verifier.setVerificationPossible(newValue);

    // End the broadcast
    vm.stopBroadcast();

    console.log('Successfully updated isVerificationPossible flag');
  }
}
