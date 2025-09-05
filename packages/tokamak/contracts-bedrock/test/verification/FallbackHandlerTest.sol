// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { console } from "forge-std/console.sol";
import { Safe } from "safe-contracts/Safe.sol";
import { SafeProxy } from "safe-contracts/proxies/SafeProxy.sol";
import { Enum as SafeOps } from "safe-contracts/libraries/Enum.sol";
import { SafeExtender } from "../../src/Safe/SafeExtender.sol";
import { Vm } from "forge-std/Vm.sol";
/**
 * @title MockFallbackHandler
 * @dev Simple mock contract to serve as a fallback handler
 */
contract MockFallbackHandler {
    event FallbackCalled(address caller, bytes data);

    // This function will be called when fallback is triggered
    function handle() external returns (bool) {
        emit FallbackCalled(msg.sender, msg.data);
        return true;
    }

    // This is called through the fallback
    fallback() external {
        // The msg.sender address is appended to the calldata by the Safe's fallback mechanism
        // We just emit an event to verify the call came through
        emit FallbackCalled(msg.sender, msg.data);
    }
}

/**
 * @title MockInterface
 * @dev For mocking external contract interfaces
 */
interface IMasterCopy {
    function masterCopy() external view returns (address);
}

/**
 * @title FallbackHandlerTest
 * @dev Tests the SafeExtender's fallback handler functionality
 */
contract FallbackHandlerTest is Test {
    // Test contracts
    Safe private safeImplementation;
    SafeProxy private safeProxy;
    SafeExtender private safeInstance;
    MockFallbackHandler private fallbackHandler;

    address private owner;
    address private secondOwner;
    address private thirdOwner;
    uint256 private ownerKey;

    function setUp() public {
        // Generate random keys and addresses
        ownerKey = uint256(keccak256(abi.encodePacked("owner key")));
        owner = vm.addr(ownerKey);
        secondOwner = vm.addr(uint256(keccak256(abi.encodePacked("second owner key"))));
        thirdOwner = vm.addr(uint256(keccak256(abi.encodePacked("third owner key"))));

        console.log("Owner address:", owner);

        // Deploy the mock fallback handler
        fallbackHandler = new MockFallbackHandler();

        // Setup owner with ETH
        vm.deal(owner, 10 ether);
        vm.deal(secondOwner, 10 ether);
        vm.deal(thirdOwner, 10 ether);

        vm.startPrank(owner);

        // Create standard Safe implementation
        safeImplementation = new SafeExtender();

        // Create SafeProxy with the implementation
        safeProxy = new SafeProxy(address(safeImplementation));

        // Cast the proxy to Safe for easier interaction
        safeInstance = SafeExtender(payable(address(safeProxy)));

        // Setup owners
        address[] memory safeOwners = new address[](1);
        safeOwners[0] = owner;

        // Initialize the Safe with setup
        bytes memory setupData = abi.encodeWithSelector(
            Safe.setup.selector,
            safeOwners,        // owners
            1,                 // threshold
            address(0),        // to
            bytes(""),         // data
            address(fallbackHandler), // fallbackHandler
            address(0),        // paymentToken
            0,                 // payment
            address(0)         // paymentReceiver
        );

        // Direct call to the proxy
        (bool success, ) = address(safeProxy).call(setupData);
        require(success, "Safe setup failed");
        assertEq(safeInstance.getFallbackHandler(), address(fallbackHandler), "Fallback handler not set correctly");
        vm.stopPrank();

        console.log("Safe proxy deployed at:", address(safeProxy));
    }

    function test_fallbackHandlerExecution() public {
        // Prepare calldata for a method that doesn't exist in Safe but exists in the fallback handler
        bytes memory callData = abi.encodeWithSignature("handle()");

        // Execute call through the proxy which should trigger the fallback mechanism
        vm.recordLogs();
        vm.prank(owner);
        (bool success, ) = address(safeProxy).call(callData);

        assertTrue(success, "Fallback handler call failed");

        Vm.Log[] memory entries = vm.getRecordedLogs();
        bool eventFound = false;
        for (uint i = 0; i < entries.length; i++) {
            // FallbackCalled event has the signature: FallbackCalled(address,bytes)
            if (entries[i].topics[0] == keccak256("FallbackCalled(address,bytes)")) {
                // Check it came from the expected handler contract
                assertEq(entries[i].emitter, address(fallbackHandler), "Event from wrong contract");
                eventFound = true;
                break;
            }
        }
        assertTrue(eventFound, "FallbackCalled event not emitted");
    }

    function test_setFallbackHandler() public {
        MockFallbackHandler newHandler = new MockFallbackHandler();

        // Prepare transaction to change the fallback handler
        bytes memory data = abi.encodeWithSignature("setFallbackHandler(address)", address(newHandler));

        // Get transaction hash for signature
        bytes32 txHash = safeInstance.getTransactionHash(
            address(safeInstance), // to: self-call to change fallback handler
            0,                     // value
            data,                  // data
            SafeOps.Operation.Call, // operation
            0,                     // safeTxGas
            0,                     // baseGas
            0,                     // gasPrice
            address(0),            // gasToken
            payable(address(0)),   // refundReceiver
            safeInstance.nonce()   // nonce
        );

        // Sign transaction with owner
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerKey, txHash);
        bytes memory signature = abi.encodePacked(r, s, v);

        // Execute transaction
        vm.prank(owner);
        bool success = safeInstance.execTransaction(
            address(safeInstance),  // to
            0,                     // value
            data,                  // data
            SafeOps.Operation.Call, // operation
            0,                     // safeTxGas
            0,                     // baseGas
            0,                     // gasPrice
            address(0),            // gasToken
            payable(address(0)),   // refundReceiver
            signature              // signature
        );

        assertTrue(success, "Failed to change fallback handler");
        assertEq(safeInstance.getFallbackHandler(), address(newHandler), "Fallback handler not changed correctly");
        // Test with the new handler
        bytes memory callData = abi.encodeWithSignature("handle()");
        vm.recordLogs();
        vm.prank(owner);
        (success, ) = address(safeProxy).call(callData);

        assertTrue(success, "Call to new fallback handler failed");
    }
}
