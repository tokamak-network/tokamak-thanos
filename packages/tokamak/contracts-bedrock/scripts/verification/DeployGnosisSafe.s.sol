// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "safe-contracts/GnosisSafe.sol";
import "safe-contracts/proxies/GnosisSafeProxyFactory.sol";
import "forge-std/console.sol";

contract DeployGnosisSafeScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // Deploy singleton/implementation contract
        GnosisSafe safeSingleton = new GnosisSafe();
        console.log("GnosisSafe singleton deployed at:", address(safeSingleton));

        // Deploy proxy factory
        GnosisSafeProxyFactory proxyFactory = new GnosisSafeProxyFactory();
        console.log("GnosisSafeProxyFactory deployed at:", address(proxyFactory));

        // Get private keys from .env and derive addresses
        uint256 ownerPrivateKey1 = vm.envUint("MULTISIG_OWNER_1");
        uint256 ownerPrivateKey2 = vm.envUint("MULTISIG_OWNER_2");
        uint256 ownerPrivateKey3 = vm.envUint("MULTISIG_OWNER_3");

        // Derive addresses from private keys
        address owner1 = vm.addr(ownerPrivateKey1);
        address owner2 = vm.addr(ownerPrivateKey2);
        address owner3 = vm.addr(ownerPrivateKey3);

        console.log("Owner 1 address:", owner1);
        console.log("Owner 2 address:", owner2);
        console.log("Owner 3 address:", owner3);

        // Setup parameters for a new Safe
        address[] memory owners = new address[](3);
        owners[0] = owner1;
        owners[1] = owner2;
        owners[2] = owner3;

        uint256 threshold = 2; // Require 2 of 3 owners to confirm a transaction
        address fallbackHandler = vm.envAddress("FALLBACK_HANDLER");

        // Create proxy with setup call encoded
        bytes memory setupData = abi.encodeWithSelector(
            GnosisSafe.setup.selector,
            owners,
            threshold,
            address(0), // to: No delegate call during setup
            bytes(""),  // data: No delegate call during setup
            fallbackHandler,
            address(0), // paymentToken: ETH
            0,          // payment: No payment
            payable(address(0)) // paymentReceiver: No payment
        );

        GnosisSafeProxy safeProxy = proxyFactory.createProxy(
            address(safeSingleton),
            setupData
        );

        console.log("GnosisSafe proxy deployed at:", address(safeProxy));

        vm.stopBroadcast();
    }
}