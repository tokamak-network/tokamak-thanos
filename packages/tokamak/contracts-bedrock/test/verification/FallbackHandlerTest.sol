// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { console } from "forge-std/console.sol";
import { GnosisSafe as Safe } from "safe-contracts/GnosisSafe.sol";
import { GnosisSafeProxy as SafeProxy } from "safe-contracts/proxies/GnosisSafeProxy.sol";
import { Enum as SafeOps } from "safe-contracts/common/Enum.sol";
import { Vm } from "forge-std/Vm.sol";

/**
 * @title IStorageAccessible
 * @dev Interface for reading storage slots from Safe contracts
 */
interface IStorageAccessible {
    function getStorageAt(uint256 offset, uint256 length) external view returns (bytes memory);
}
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
 * @dev Tests the Safe's fallback handler functionality using standard Safe contracts
 */
contract FallbackHandlerTest is Test {
    // Test contracts
    Safe private safeImplementation;
    SafeProxy private safeProxy;
    Safe private safeInstance;
    MockFallbackHandler private fallbackHandler;

    address private owner;
    address private secondOwner;
    address private thirdOwner;
    uint256 private ownerKey;

    // Storage slot for Safe's fallback handler
    bytes32 internal constant FALLBACK_HANDLER_STORAGE_SLOT =
        0x6c9a6c4a39284e37ed1cf53d337577d14212a4870fb976a4366c693b939918d5;

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
        safeImplementation = new Safe();

        // Create SafeProxy with the implementation
        safeProxy = new SafeProxy(address(safeImplementation));

        // Cast the proxy to Safe for easier interaction
        safeInstance = Safe(payable(address(safeProxy)));

        // Setup owners
        address[] memory safeOwners = new address[](1);
        safeOwners[0] = owner;

        // Initialize the Safe with setup
        bytes memory setupData = abi.encodeWithSelector(
            Safe.setup.selector,
            safeOwners, // owners
            1, // threshold
            address(0), // to
            bytes(""), // data
            address(fallbackHandler), // fallbackHandler
            address(0), // paymentToken
            0, // payment
            address(0) // paymentReceiver
        );

        // Direct call to the proxy
        (bool success,) = address(safeProxy).call(setupData);
        require(success, "Safe setup failed");
        assertEq(_getFallbackHandler(address(safeInstance)), address(fallbackHandler), "Fallback handler not set correctly");
        vm.stopPrank();

        console.log("Safe proxy deployed at:", address(safeProxy));
    }

    function test_fallbackHandlerExecution() public {
        // Prepare calldata for a method that doesn't exist in Safe but exists in the fallback handler
        bytes memory callData = abi.encodeWithSignature("handle()");

        // Execute call through the proxy which should trigger the fallback mechanism
        vm.recordLogs();
        vm.prank(owner);
        (bool success,) = address(safeProxy).call(callData);

        assertTrue(success, "Fallback handler call failed");

        Vm.Log[] memory entries = vm.getRecordedLogs();
        bool eventFound = false;
        for (uint256 i = 0; i < entries.length; i++) {
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
            0, // value
            data, // data
            SafeOps.Operation.Call, // operation
            0, // safeTxGas
            0, // baseGas
            0, // gasPrice
            address(0), // gasToken
            payable(address(0)), // refundReceiver
            safeInstance.nonce() // nonce
        );

        // Sign transaction with owner
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerKey, txHash);
        bytes memory signature = abi.encodePacked(r, s, v);

        // Execute transaction
        vm.prank(owner);
        bool success = safeInstance.execTransaction(
            address(safeInstance), // to
            0, // value
            data, // data
            SafeOps.Operation.Call, // operation
            0, // safeTxGas
            0, // baseGas
            0, // gasPrice
            address(0), // gasToken
            payable(address(0)), // refundReceiver
            signature // signature
        );

        assertTrue(success, "Failed to change fallback handler");
        assertEq(_getFallbackHandler(address(safeInstance)), address(newHandler), "Fallback handler not changed correctly");
        // Test with the new handler
        bytes memory callData = abi.encodeWithSignature("handle()");
        vm.recordLogs();
        vm.prank(owner);
        (success,) = address(safeProxy).call(callData);

        assertTrue(success, "Call to new fallback handler failed");
    }

    function test_getFallbackHandlerReturnsCorrectAddress() public {
        // Verify that _getFallbackHandler correctly returns the fallback handler address
        // that was set during setUp
        address retrievedHandler = _getFallbackHandler(address(safeInstance));

        assertEq(retrievedHandler, address(fallbackHandler), "Retrieved handler does not match expected handler");
        assertTrue(retrievedHandler != address(0), "Fallback handler should not be zero address");

        console.log("Fallback handler correctly retrieved:", retrievedHandler);
        console.log("Expected handler address:", address(fallbackHandler));
    }

    function test_getFallbackHandlerReturnsZeroWhenNotSet() public {
        // Create a new Safe without a fallback handler
        Safe newSafeImpl = new Safe();
        SafeProxy newSafeProxy = new SafeProxy(address(newSafeImpl));
        Safe newSafeInstance = Safe(payable(address(newSafeProxy)));

        address[] memory safeOwners = new address[](1);
        safeOwners[0] = owner;

        // Initialize the Safe WITHOUT a fallback handler (set to address(0))
        bytes memory setupData = abi.encodeWithSelector(
            Safe.setup.selector,
            safeOwners, // owners
            1, // threshold
            address(0), // to
            bytes(""), // data
            address(0), // fallbackHandler - explicitly set to zero
            address(0), // paymentToken
            0, // payment
            address(0) // paymentReceiver
        );

        vm.prank(owner);
        (bool success,) = address(newSafeProxy).call(setupData);
        require(success, "Safe setup failed");

        // Verify that _getFallbackHandler returns address(0)
        address retrievedHandler = _getFallbackHandler(address(newSafeInstance));

        assertEq(retrievedHandler, address(0), "Handler should be zero address when not set");

        console.log("Fallback handler correctly returned as zero:", retrievedHandler);
    }

    function test_detectsFallbackHandlerChange() public {
        // Verify initial handler
        address initialHandler = _getFallbackHandler(address(safeInstance));
        assertEq(initialHandler, address(fallbackHandler), "Initial handler mismatch");

        // Create a new handler
        MockFallbackHandler newHandler = new MockFallbackHandler();

        // Change the fallback handler
        bytes memory data = abi.encodeWithSignature("setFallbackHandler(address)", address(newHandler));
        bytes32 txHash = safeInstance.getTransactionHash(
            address(safeInstance),
            0,
            data,
            SafeOps.Operation.Call,
            0, 0, 0,
            address(0),
            payable(address(0)),
            safeInstance.nonce()
        );

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerKey, txHash);
        bytes memory signature = abi.encodePacked(r, s, v);

        vm.prank(owner);
        bool success = safeInstance.execTransaction(
            address(safeInstance),
            0,
            data,
            SafeOps.Operation.Call,
            0, 0, 0,
            address(0),
            payable(address(0)),
            signature
        );

        assertTrue(success, "Failed to change handler");

        // Verify the handler changed
        address updatedHandler = _getFallbackHandler(address(safeInstance));
        assertEq(updatedHandler, address(newHandler), "Handler not updated correctly");
        assertTrue(updatedHandler != initialHandler, "Handler should have changed");

        console.log("Handler successfully changed from:", initialHandler);
        console.log("To new handler:", updatedHandler);
    }

    function test_detectsArbitraryFallbackHandlerAddress() public {
        // Create a new Safe and set a specific arbitrary address as the fallback handler
        Safe newSafeImpl = new Safe();
        SafeProxy newSafeProxy = new SafeProxy(address(newSafeImpl));
        Safe newSafeInstance = Safe(payable(address(newSafeProxy)));

        address[] memory safeOwners = new address[](1);
        safeOwners[0] = owner;

        // Use a specific arbitrary address as fallback handler
        address arbitraryHandler = address(0x123123123);

        // Initialize the Safe WITH the arbitrary fallback handler
        bytes memory setupData = abi.encodeWithSelector(
            Safe.setup.selector,
            safeOwners, // owners
            1, // threshold
            address(0), // to
            bytes(""), // data
            arbitraryHandler, // fallbackHandler - set to arbitrary address
            address(0), // paymentToken
            0, // payment
            address(0) // paymentReceiver
        );

        vm.prank(owner);
        (bool success,) = address(newSafeProxy).call(setupData);
        require(success, "Safe setup failed");

        // Verify that _getFallbackHandler returns the exact arbitrary address we set
        address retrievedHandler = _getFallbackHandler(address(newSafeInstance));

        assertEq(retrievedHandler, arbitraryHandler, "Retrieved handler does not match arbitrary handler");
        assertEq(retrievedHandler, address(0x123123123), "Handler should be exactly 0x123123123");
        assertTrue(retrievedHandler != address(0), "Handler should not be zero address");

        console.log("Arbitrary fallback handler correctly retrieved:", retrievedHandler);
        console.log("Expected arbitrary address:", arbitraryHandler);
    }

    function test_directGetStorageAtCall() public {
        // This test explicitly demonstrates calling getStorageAt on a real Safe contract
        // to prove that StorageAccessible is available and works

        console.log("=== Direct getStorageAt Test ===");
        console.log("Safe address:", address(safeInstance));
        console.log("Storage slot:", uint256(FALLBACK_HANDLER_STORAGE_SLOT));

        // Direct call to getStorageAt - this proves the function exists and is callable
        bytes memory storageData = IStorageAccessible(address(safeInstance)).getStorageAt(
            uint256(FALLBACK_HANDLER_STORAGE_SLOT),
            1  // Read 1 word (32 bytes)
        );

        console.log("Raw storage data length:", storageData.length);
        console.logBytes(storageData);

        // Parse the address from storage
        address parsedHandler;
        assembly {
            parsedHandler := mload(add(storageData, 0x20))
        }

        console.log("Parsed handler address:", parsedHandler);
        console.log("Expected handler:", address(fallbackHandler));

        // Verify it matches
        assertEq(parsedHandler, address(fallbackHandler), "Direct storage read should match expected handler");
        assertEq(storageData.length, 32, "Should return exactly 32 bytes (1 word)");

        console.log("SUCCESS: Called getStorageAt and retrieved fallback handler!");
    }

    /**
     * @notice Helper function to get fallback handler from Safe storage
     * @param safeAddress The address of the Safe to query
     * @return The address of the fallback handler
     */
    function _getFallbackHandler(address safeAddress) private view returns (address) {
        bytes memory result = IStorageAccessible(safeAddress).getStorageAt(
            uint256(FALLBACK_HANDLER_STORAGE_SLOT),
            1
        );

        return abi.decode(result, (address));
    }
}
