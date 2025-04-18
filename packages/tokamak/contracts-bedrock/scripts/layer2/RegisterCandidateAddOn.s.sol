// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/layer2/interfaces/ILayer2Manager.sol';
import '@openzeppelin/contracts/token/ERC20/IERC20.sol';

/**
 * @title RegisterCandidateAddOn
 * @notice Script to register a candidate add-on with the Layer2Manager contract
 * @dev This script registers a new candidate add-on by depositing TON tokens
 */
contract RegisterCandidateAddOn is Script {
    // Environment variables
    address private _layer2Manager;
    address private _rollupConfig;
    // address private _ton;
    uint256 private _depositAmount;
    string private _memo;

    function setUp() public {
        // Load environment variables
        _layer2Manager = vm.envAddress('LAYER2_MANAGER');
        _rollupConfig = vm.envAddress('ROLLUP_CONFIG');
        // _ton = vm.envAddress('TON');
        _depositAmount = vm.envUint('DEPOSIT_AMOUNT');
        _memo = vm.envString('MEMO');
    }

    function run() external {
        setUp();

        vm.startBroadcast();

        // Get the Layer2Manager contract
        ILayer2Manager layer2Manager = ILayer2Manager(_layer2Manager);

        // // Check if registration is available
        // bool isAvailable = layer2Manager.availableRegister(_rollupConfig);
        // require(isAvailable, "Registration not available for this rollup config");

        // // Check minimum deposit amount
        // (bool success, uint256 tvl) = layer2Manager.checkLayer2TVL(_rollupConfig);
        // require(success, "Failed to check TVL");
        // console.log("Current TVL:", tvl);

        // Approve TON tokens for the Layer2Manager
        bool approvehash = IERC20(0x33a66929dE3559315c928556FcFF449b3E708c62).approve(_layer2Manager, _depositAmount);
        console.log("Approve hash:", approvehash);
        console.log("Approved TON tokens for deposit");

        // // Register the candidate add-on
        // layer2Manager.registerCandidateAddOn(
        //     _rollupConfig,
        //     _depositAmount,
        //     true, // flagTon = true means we're using TON
        //     _memo
        // );
        // console.log("Candidate add-on registered");

        // // Get the operator and candidate add-on addresses
        // address operator = layer2Manager.operatorOfRollupConfig(_rollupConfig);
        // address candidateAddOn = layer2Manager.candidateAddOnOfOperator(operator);
        // console.log("Operator:", operator);
        // console.log("Candidate Add-On:", candidateAddOn);

        vm.stopBroadcast();
    }
}
