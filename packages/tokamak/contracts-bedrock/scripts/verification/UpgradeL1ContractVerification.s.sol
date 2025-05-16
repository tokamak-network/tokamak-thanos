// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import '@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol';
import {IMultiSigWallet} from '../../src/tokamak-contracts/verification/interface/IMultiSigWallet.sol';

/**
 * @title UpgradeL1ContractVerification
 * @notice Script to upgrade the L1ContractVerification contract implementation
 * @dev This script deploys a new implementation and updates the proxy to point to it
 */
contract UpgradeL1ContractVerification is Script {
    // Environment variables - private variables with leading underscore
    address private _proxyAddress;
    address private _verificationContractProxyAdminAddress;
    uint256 private _multisigOwner1;
    uint256 private _multisigOwner2;
    uint256 private _multisigOwner3;
    address private _multisigWallet;

    /**
     * @notice Load configuration from environment variables
     * @dev Called before the main execution to set up required addresses
     */
    function setUp() public {
        // Load environment variables for addresses
        _proxyAddress = vm.envAddress('L1_CONTRACT_VERIFICATION_PROXY');
        _verificationContractProxyAdminAddress = vm.envAddress('L1_CONTRACT_VERIFICATION_PROXY_ADMIN');
        _multisigOwner1 = vm.envUint('MULTISIG_OWNER_1');
        _multisigOwner2 = vm.envUint('MULTISIG_OWNER_2');
        _multisigOwner3 = vm.envUint('MULTISIG_OWNER_3');
        _multisigWallet = vm.envAddress('MULTISIG_WALLET');
    }

    /**
     * @notice Main execution function
     * @dev Deploys a new implementation and upgrades the proxy through multisig
     */
    function run() external {
        setUp();

        // Get the ProxyAdmin contract
        ProxyAdmin proxyAdmin = ProxyAdmin(_verificationContractProxyAdminAddress);

        // Log the current owner
        address currentOwner = proxyAdmin.owner();
        console.log("Current ProxyAdmin owner:", currentOwner);
        console.log("Expected multisig address:", _multisigWallet);

        require(currentOwner == _multisigWallet, "ProxyAdmin not owned by multisig!");

        // Deploy new implementation
        vm.broadcast(_multisigOwner1);
        L1ContractVerification newImplementation = new L1ContractVerification();
        console.log("New L1ContractVerification implementation deployed at:", address(newImplementation));

        // Get current implementation for comparison
        address currentImplementation = proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(_proxyAddress)));
        console.log("Current implementation:", currentImplementation);

        // Prepare the upgrade transaction data
        bytes memory upgradeData = abi.encodeWithSelector(
            ProxyAdmin.upgrade.selector,
            _proxyAddress,
            address(newImplementation)
        );

        // Get the multisig contract
        IMultiSigWallet multiSig = IMultiSigWallet(_multisigWallet);

        // Owner 1 submits the transaction
        vm.broadcast(_multisigOwner1);
        multiSig.submitTransaction(
            _verificationContractProxyAdminAddress,
            0,
            upgradeData
        );
        uint txIndex = multiSig.getTransactionCount() - 1;
        console.log("Transaction submitted by owner 1, txIndex:", txIndex);

        // Owner 2 confirms the transaction
        vm.broadcast(_multisigOwner2);
        multiSig.confirmTransaction(txIndex);
        console.log("Transaction confirmed by owner 2");

        // Owner 3 confirms the transaction
        vm.broadcast(_multisigOwner3);
        multiSig.confirmTransaction(txIndex);
        console.log("Transaction confirmed by owner 3");

        // Owner 1 executes the transaction
        vm.broadcast(_multisigOwner1);
        multiSig.executeTransaction(txIndex);
        console.log("Transaction executed by owner 1");

        // Verify the upgrade was successful
        address newImplementationAddress = proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(_proxyAddress)));
        console.log("New implementation confirmed:", newImplementationAddress);
        require(newImplementationAddress == address(newImplementation), "Upgrade failed");
    }
}