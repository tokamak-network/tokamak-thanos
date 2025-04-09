// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import '@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol';

interface IMultiSigWallet {
    function submitTransaction(
        address _to,
        uint _value,
        bytes memory _data
    ) external;

    function confirmTransaction(uint _txIndex) external;

    function executeTransaction(uint _txIndex) external;

    function getTransactionCount() external view returns (uint);
}

/**
 * @title VerifyContractSetConfig
 * @notice Script to verify L1 contracts using multisig signatures
 * @dev This script is used to verify contracts after deployment
 */
contract VerifyContractSetConfig is Script {
    // Environment variables
    address private _safeWalletAddress;
    address private _proxyAddress;
    address private _systemConfigProxy;
    address private _l1ProxyAdmin;
    uint256 private _multisigOwner1;
    uint256 private _multisigOwner2;
    uint256 private _multisigOwner3;
    address private _multisigWallet;

    function setUp() public {
        // Load environment variables
        _proxyAddress = vm.envAddress('L1_CONTRACT_VERIFICATION_PROXY');
        _systemConfigProxy = vm.envAddress('SYSTEM_CONFIG_PROXY');
        _l1ProxyAdmin = vm.envAddress('PROXY_ADMIN_ADDRESS');
        _multisigOwner1 = vm.envUint('MULTISIG_OWNER_1');
        _multisigOwner2 = vm.envUint('MULTISIG_OWNER_2');
        _multisigOwner3 = vm.envUint('MULTISIG_OWNER_3');
        _multisigWallet = vm.envAddress('MULTISIG_WALLET');
        _safeWalletAddress = vm.envAddress('SAFE_WALLET');
    }

    function run() external {
        setUp();
        vm.startBroadcast();
        // Create reference to the deployed proxy
        L1ContractVerification verifier = L1ContractVerification(_proxyAddress);

        verifier.verifyL1Contracts(_systemConfigProxy, _l1ProxyAdmin, _safeWalletAddress);

        console.log('L1 contracts verification completed successfully');
        vm.stopBroadcast();
    }
}
