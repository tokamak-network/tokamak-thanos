// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import { ERC721Bridge_Initializer } from "test/CommonTest.t.sol";
import {
    Messenger_Initializer, Reverter, ConfigurableCaller, SuperchainConfig_Initializer
} from "test/CommonTest.t.sol";
import { L2OutputOracle_Initializer } from "test/L2OutputOracle.t.sol";
import "forge-std/console.sol";
// Libraries
import { AddressAliasHelper } from "src/vendor/AddressAliasHelper.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Hashing } from "src/libraries/Hashing.sol";
import { Encoding } from "src/libraries/Encoding.sol";

// Target contract dependencies
import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";

// Target contract
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract L1CrossDomainMessenger_Test is Messenger_Initializer {
    /// @dev The receiver address
    address recipient = address(0xabbaacdc);

    /// @dev The storage slot of the l2Sender
    uint256 constant senderSlotIndex = 50;

    /// @dev Tests that the implementation is initialized correctly.

    /// @dev Tests that the proxy is initialized correctly.
    function test_initialize_succeeds() external {
        assertEq(address(L1Messenger.superchainConfig()), address(sc));
        assertEq(address(L1Messenger.PORTAL()), address(op));
        assertEq(address(L1Messenger.portal()), address(op));
        assertEq(address(L1Messenger.OTHER_MESSENGER()), Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        assertEq(address(L1Messenger.otherMessenger()), Predeploys.L2_CROSS_DOMAIN_MESSENGER);
    }

    /// @dev Tests that the version can be decoded from the message nonce.
    function test_messageVersion_succeeds() external {
        (, uint16 version) = Encoding.decodeVersionedNonce(L1Messenger.messageNonce());
        assertEq(version, L1Messenger.MESSAGE_VERSION());
    }

    /// @dev Tests that the sendMessage function is able to send a single message.
    /// TODO: this same test needs to be done with the legacy message type
    ///       by setting the message version to 0
    function test_sendNativeTokenMessage_succeeds() external {
        // faucet to alice's balance
        dealL2NativeToken(alice, NON_ZERO_VALUE);

        vm.prank(alice);
        token.approve(address(L1Messenger), NON_ZERO_VALUE);

        // deposit transaction on the optimism portal should be called
        vm.expectCall(
            address(op),
            abi.encodeWithSelector(
                OptimismPortal.depositTransaction.selector,
                Predeploys.L2_CROSS_DOMAIN_MESSENGER,
                NON_ZERO_VALUE,
                NON_ZERO_VALUE,
                L1Messenger.baseGas(hex"ff", 100),
                false,
                Encoding.encodeCrossDomainMessage(
                    L1Messenger.messageNonce(), alice, recipient, NON_ZERO_VALUE, 100, hex"ff"
                )
            )
        );

        // TransactionDeposited event
        vm.expectEmit(true, true, true, true);
        emitTransactionDeposited(
            AddressAliasHelper.applyL1ToL2Alias(address(L1Messenger)),
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            NON_ZERO_VALUE,
            NON_ZERO_VALUE,
            L1Messenger.baseGas(hex"ff", 100),
            false,
            Encoding.encodeCrossDomainMessage(
                L1Messenger.messageNonce(), alice, recipient, NON_ZERO_VALUE, 100, hex"ff"
            )
        );

        // SentMessage event
        vm.expectEmit(true, true, true, true);
        emit SentMessage(recipient, alice, hex"ff", L1Messenger.messageNonce(), 100);

        // SentMessageExtension1 event
        vm.expectEmit(true, true, true, true);
        emit SentMessageExtension1(alice, NON_ZERO_VALUE);

        vm.prank(alice);
        L1Messenger.sendNativeTokenMessage(recipient, NON_ZERO_VALUE, hex"ff", uint32(100));
    }

    /// @dev Tests that the sendMessage function is able to send
    ///      the same message twice.
    function test_sendNativeTokenMessage_twice_succeeds() external {
        dealL2NativeToken(address(this), NON_ZERO_VALUE * 2);
        token.approve(address(L1Messenger), NON_ZERO_VALUE * 2);
        uint256 nonce = L1Messenger.messageNonce();
        L1Messenger.sendNativeTokenMessage(recipient, NON_ZERO_VALUE, hex"aa", uint32(500_000));
        L1Messenger.sendNativeTokenMessage(recipient, NON_ZERO_VALUE, hex"aa", uint32(500_000));
        // the nonce increments for each message sent
        assertEq(nonce + 2, L1Messenger.messageNonce());
    }

    /// @dev Tests that the sendMessage function is able to send a single message.
    /// TODO: this same test needs to be done with the legacy message type
    ///       by setting the message version to 0
    function test_sendMessage_succeeds() external {
        // deposit transaction on the optimism portal should be called
        vm.expectCall(
            address(op),
            abi.encodeWithSelector(
                OptimismPortal.depositTransaction.selector,
                Predeploys.L2_CROSS_DOMAIN_MESSENGER,
                0,
                0,
                L1Messenger.baseGas(hex"ff", 100),
                false,
                Encoding.encodeCrossDomainMessage(L1Messenger.messageNonce(), alice, recipient, 0, 100, hex"ff")
            )
        );

        // TransactionDeposited event
        vm.expectEmit(true, true, true, true);
        emitTransactionDeposited(
            AddressAliasHelper.applyL1ToL2Alias(address(L1Messenger)),
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            0,
            0,
            L1Messenger.baseGas(hex"ff", 100),
            false,
            Encoding.encodeCrossDomainMessage(L1Messenger.messageNonce(), alice, recipient, 0, 100, hex"ff")
        );

        // SentMessage event
        vm.expectEmit(true, true, true, true);
        emit SentMessage(recipient, alice, hex"ff", L1Messenger.messageNonce(), 100);

        // SentMessageExtension1 event
        vm.expectEmit(true, true, true, true);
        emit SentMessageExtension1(alice, 0);

        vm.prank(alice);
        L1Messenger.sendMessage(recipient, hex"ff", uint32(100));
    }

    /// @dev Tests that the sendMessage function is able to send
    ///      the same message twice.
    function test_sendMessage_twice_succeeds() external {
        uint256 nonce = L1Messenger.messageNonce();
        L1Messenger.sendMessage(recipient, hex"aa", uint32(500_000));
        L1Messenger.sendMessage(recipient, hex"aa", uint32(500_000));
        // the nonce increments for each message sent
        assertEq(nonce + 2, L1Messenger.messageNonce());
    }

    /// @dev Tests that the xDomainMessageSender reverts when not set.
    function test_xDomainSender_notSet_reverts() external {
        vm.expectRevert("CrossDomainMessenger: xDomainMessageSender is not set");
        L1Messenger.xDomainMessageSender();
    }

    /// @dev Tests that the relayMessage function reverts when
    ///      the message version is not 0 or 1.
    function test_relayMessage_v2_reverts() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;

        // Set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));

        // Expect a revert.
        vm.expectRevert("CrossDomainMessenger: only version 0 or 1 messages are supported at this time");

        // Try to relay a v2 message.
        vm.prank(address(op));
        L2Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 2 }), // nonce
            sender,
            target,
            0, // value
            0,
            hex"1111"
        );
    }

    /// @dev Tests that the relayMessage function is able to relay a message
    ///      successfully by calling the target contract.
    function test_relayMessage_succeeds() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;

        vm.expectCall(target, hex"1111");

        // set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));
        vm.prank(address(op));

        vm.expectEmit(true, true, true, true);

        bytes32 hash = Hashing.hashCrossDomainMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), sender, target, 0, 0, hex"1111"
        );

        emit RelayedMessage(hash);

        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), // nonce
            sender,
            target,
            0, // value
            0,
            hex"1111"
        );

        // the message hash is in the successfulMessages mapping
        assert(L1Messenger.successfulMessages(hash));
        // it is not in the received messages mapping
        assertEq(L1Messenger.failedMessages(hash), false);
    }

    /// @dev Tests that relayMessage reverts if attempting to relay a message
    ///      sent to an L1 system contract.
    function test_relayMessage_toSystemContract_reverts() external {
        // set the target to be the OptimismPortal
        address target = address(op);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        bytes memory message = hex"1111";

        vm.prank(address(op));
        vm.expectRevert("CrossDomainMessenger: message cannot be replayed");
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), sender, target, 0, 0, message
        );

        vm.store(address(op), 0, bytes32(abi.encode(sender)));
        vm.expectRevert("CrossDomainMessenger: message cannot be replayed");
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), sender, target, 0, 0, message
        );
    }

    /// @dev Tests that the relayMessage function reverts if eth is
    ///      sent from a contract other than the standard bridge.
    function test_replayMessage_withValue_reverts() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        bytes memory message = hex"1111";

        vm.expectRevert("CrossDomainMessenger: value must be zero");
        L1Messenger.relayMessage{ value: 100 }(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), sender, target, 0, 0, message
        );
    }

    /// @dev Tests that the xDomainMessageSender is reset to the original value
    ///      after a message is relayed.
    function test_xDomainMessageSender_reset_succeeds() external {
        vm.expectRevert("CrossDomainMessenger: xDomainMessageSender is not set");
        L1Messenger.xDomainMessageSender();

        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;

        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));
        vm.prank(address(op));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), address(0), address(0), 0, 0, hex""
        );

        vm.expectRevert("CrossDomainMessenger: xDomainMessageSender is not set");
        L1Messenger.xDomainMessageSender();
    }

    /// @dev Tests that relayMessage should successfully call the target contract after
    ///      the first message fails and NativeToken is stuck, but the second message succeeds
    ///      with a version 1 message.
    function test_relayMessage_retryAfterFailure_succeeds() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        uint256 value = 100;

        dealL2NativeToken(address(op), 2 * value);
        // Approve the L1Messenger contract
        vm.prank(address(op));
        IERC20(token).approve(address(L1Messenger), type(uint256).max);

        vm.expectCall(target, hex"1111");

        bytes32 hash = Hashing.hashCrossDomainMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), sender, target, value, 0, hex"1111"
        );

        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));
        vm.etch(target, address(new Reverter()).code);
        vm.prank(address(op));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        assertEq(address(L1Messenger).balance, 0);
        assertEq(address(target).balance, 0);
        assertEq(L1Messenger.successfulMessages(hash), false);
        assertEq(L1Messenger.failedMessages(hash), true);

        vm.expectEmit(true, true, true, true);

        emit RelayedMessage(hash);

        vm.etch(target, address(0).code);
        vm.prank(address(sender));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        assertEq(IERC20(token).balanceOf(address(L1Messenger)), value);
        assertEq(IERC20(token).balanceOf(address(target)), 0);
        assertEq(address(L1Messenger).balance, 0);
        assertEq(address(target).balance, 0);
        assertEq(L1Messenger.successfulMessages(hash), true);
        assertEq(L1Messenger.failedMessages(hash), true);
    }

    /// @dev Tests that relayMessage should successfully call the target contract after
    ///      the first message fails and ETH is stuck, but the second message succeeds
    ///      with a legacy message.
    function test_relayMessage_legacy_succeeds() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;

        // Compute the message hash.
        bytes32 hash = Hashing.hashCrossDomainMessageV1(
            // Using a legacy nonce with version 0.
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }),
            sender,
            target,
            0,
            0,
            hex"1111"
        );

        // Set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));

        // Target should be called with expected data.
        vm.expectCall(target, hex"1111");

        // Expect RelayedMessage event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit RelayedMessage(hash);

        // Relay the message.
        vm.prank(address(op));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            0, // value
            0,
            hex"1111"
        );

        // Message was successfully relayed.
        assertEq(L1Messenger.successfulMessages(hash), true);
        assertEq(L1Messenger.failedMessages(hash), false);
    }

    /// @dev Tests that relayMessage should revert if the message is already replayed.
    function test_relayMessage_legacyOldReplay_reverts() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;

        // Compute the message hash.
        bytes32 hash = Hashing.hashCrossDomainMessageV1(
            // Using a legacy nonce with version 0.
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }),
            sender,
            target,
            0,
            0,
            hex"1111"
        );

        // Set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));

        // Mark legacy message as already relayed.
        uint256 successfulMessagesSlot = 203;
        bytes32 oldHash = Hashing.hashCrossDomainMessageV0(target, sender, hex"1111", 0);
        bytes32 slot = keccak256(abi.encode(oldHash, successfulMessagesSlot));
        vm.store(address(L1Messenger), slot, bytes32(uint256(1)));

        // Expect revert.
        vm.expectRevert("CrossDomainMessenger: legacy withdrawal already relayed");

        // Relay the message.
        vm.prank(address(op));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            0, // value
            0,
            hex"1111"
        );

        // Message was not relayed.
        assertEq(L1Messenger.successfulMessages(hash), false);
        assertEq(L1Messenger.failedMessages(hash), false);
    }

    /// @dev Tests that relayMessage can be retried after a failure with a legacy message.
    function test_relayMessage_legacyRetryAfterFailure_succeeds() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        uint256 value = 100;

        dealL2NativeToken(address(op), 2 * value);
        // Approve the L1Messenger contract
        vm.prank(address(op));
        IERC20(token).approve(address(L1Messenger), type(uint256).max);

        // Compute the message hash.
        bytes32 hash = Hashing.hashCrossDomainMessageV1(
            // Using a legacy nonce with version 0.
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }),
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));

        // Turn the target into a Reverter.
        vm.etch(target, address(new Reverter()).code);

        // Target should be called with expected data.
        vm.expectCall(target, hex"1111");

        // Expect FailedRelayedMessage event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit FailedRelayedMessage(hash);

        // Relay the message.
        vm.prank(address(op));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Message failed.
        assertEq(IERC20(token).balanceOf(address(L1Messenger)), value);
        assertEq(IERC20(token).balanceOf(address(target)), 0);
        assertEq(address(L1Messenger).balance, 0);
        assertEq(address(target).balance, 0);
        assertEq(L1Messenger.successfulMessages(hash), false);
        assertEq(L1Messenger.failedMessages(hash), true);

        // Make the target not revert anymore.
        vm.etch(target, address(0).code);

        // Target should be called with expected data.
        vm.expectCall(target, hex"1111");

        // Expect RelayedMessage event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit RelayedMessage(hash);

        // Retry the message.
        vm.prank(address(sender));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Message was successfully relayed.
        assertEq(IERC20(token).balanceOf(address(L1Messenger)), value);
        assertEq(IERC20(token).balanceOf(address(target)), 0);
        assertEq(address(L1Messenger).balance, 0);
        assertEq(address(target).balance, 0);
        assertEq(L1Messenger.successfulMessages(hash), true);
        assertEq(L1Messenger.failedMessages(hash), true);
    }

    /// @dev Tests that relayMessage cannot be retried after success with a legacy message.
    function test_relayMessage_legacyRetryAfterSuccess_reverts() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        uint256 value = 100;

        dealL2NativeToken(address(op), value);
        // Approve the L1Messenger contract
        vm.prank(address(op));
        IERC20(token).approve(address(L1Messenger), type(uint256).max);

        // Compute the message hash.
        bytes32 hash = Hashing.hashCrossDomainMessageV1(
            // Using a legacy nonce with version 0.
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }),
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));

        // Target should be called with expected data.
        vm.expectCall(target, hex"1111");

        // Expect RelayedMessage event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit RelayedMessage(hash);

        // Relay the message.
        vm.prank(address(op));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Message was successfully relayed.
        assertEq(IERC20(token).balanceOf(address(L1Messenger)), value);
        assertEq(IERC20(token).balanceOf(address(target)), 0);
        assertEq(address(target).balance, 0);
        assertEq(address(target).balance, 0);
        assertEq(L1Messenger.successfulMessages(hash), true);
        assertEq(L1Messenger.failedMessages(hash), false);

        // // Expect a revert.
        vm.expectRevert("CrossDomainMessenger: message cannot be replayed");

        // // Retry the message.
        vm.prank(address(sender));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );
    }

    function test_relayMessage_legacyRetryAfterFailureThenSuccess_reverts() external {
        address target = address(0xabcd);
        address sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        uint256 value = 100;

        dealL2NativeToken(address(op), 2 * value);

        // Approve the L1Messenger contract
        vm.prank(address(op));
        IERC20(token).approve(address(L1Messenger), type(uint256).max);

        // Compute the message hash.
        bytes32 hash = Hashing.hashCrossDomainMessageV1(
            // Using a legacy nonce with version 0.
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }),
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));

        // Turn the target into a Reverter.
        vm.etch(target, address(new Reverter()).code);

        // Target should be called with expected data.
        vm.expectCall(target, hex"1111");

        // Relay the message.
        vm.prank(address(op));

        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Message failed.
        assertEq(IERC20(token).balanceOf(address(L1Messenger)), value);
        assertEq(IERC20(token).balanceOf(address(target)), 0);
        assertEq(address(target).balance, 0);
        assertEq(address(L1Messenger).balance, 0);
        assertEq(L1Messenger.successfulMessages(hash), false);
        assertEq(L1Messenger.failedMessages(hash), true);

        // Make the target not revert anymore.
        vm.etch(target, address(0).code);

        // Target should be called with expected data.
        vm.expectCall(target, hex"1111");

        // Expect RelayedMessage event to be emitted.
        vm.expectEmit(true, true, true, true);
        emit RelayedMessage(hash);

        // Retry the message
        vm.prank(address(sender));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );

        // Message was successfully relayed.
        assertEq(address(L1Messenger).balance, 0);
        assertEq(address(target).balance, 0);
        assertEq(IERC20(token).balanceOf(address(L1Messenger)), value);
        assertEq(IERC20(token).balanceOf(address(target)), 0);

        assertEq(L1Messenger.successfulMessages(hash), true);
        assertEq(L1Messenger.failedMessages(hash), true);

        // Expect a revert.
        vm.expectRevert("CrossDomainMessenger: message has already been relayed");

        // Retry the message again.
        vm.prank(address(sender));
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 0 }), // nonce
            sender,
            target,
            value,
            0,
            hex"1111"
        );
    }

    /// @dev Tests that the superchain config is called by the messengers paused function
    function test_pause_callsSuperchainConfig_succeeds() external {
        vm.expectCall(address(sc), abi.encodeWithSelector(SuperchainConfig.paused.selector));
        L1Messenger.paused();
    }

    /// @dev Tests that changing the superchain config paused status changes the return value of the messenger
    function test_pause_matchesSuperchainConfig_succeeds() external {
        assertFalse(L1Messenger.paused());
        assertEq(L1Messenger.paused(), sc.paused());

        vm.prank(sc.guardian());
        sc.pause("identifier");

        assertTrue(L1Messenger.paused());
        assertEq(L1Messenger.paused(), sc.paused());
    }
}

/// @dev A regression test against a reentrancy vulnerability in the CrossDomainMessenger contract, which
///      was possible by intercepting and sandwhiching a signed Safe Transaction to upgrade it.
contract L1CrossDomainMessenger_ReinitReentryTest is Messenger_Initializer {
    bool attacked;

    // Common values used across functions
    uint256 constant messageValue = 50;
    bytes constant selector = abi.encodeWithSelector(this.reinitAndReenter.selector);
    address sender;
    bytes32 hash;
    address target;

    function setUp() public override {
        super.setUp();
        target = address(this);
        sender = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        hash = Hashing.hashCrossDomainMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), sender, target, messageValue, 0, selector
        );
        vm.deal(address(L1Messenger), messageValue * 2);
    }

    /// @dev This method will be called by the relayed message, and will attempt to reenter the relayMessage function
    ///      exactly once.
    function reinitAndReenter() public payable {
        // only attempt the attack once
        if (!attacked) {
            attacked = true;
            // set initialized to false
            vm.store(address(L1Messenger), 0, bytes32(uint256(0)));

            // call the initializer function
            L1Messenger.initialize(SuperchainConfig(sc), OptimismPortal(op), address(token));

            // attempt to re-replay the withdrawal
            vm.expectEmit(address(L1Messenger));
            emit FailedRelayedMessage(hash);
            L1Messenger.relayMessage(
                Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), // nonce
                sender,
                target,
                messageValue,
                0,
                selector
            );
        }
    }

    /// @dev Tests that the relayMessage function cannot be reentered by calling the `initialize()` function within the
    ///      relayed message.
    function test_relayMessage_replayStraddlingReinit_reverts() external {
        uint256 balanceBeforeThis = address(this).balance;
        uint256 balanceBeforeMessenger = address(L1Messenger).balance;

        // A requisite for the attack is that the message has already been attempted and written to the failedMessages
        // mapping, so that it can be replayed.
        vm.store(address(L1Messenger), keccak256(abi.encode(hash, 206)), bytes32(uint256(1)));
        assertTrue(L1Messenger.failedMessages(hash));

        vm.expectEmit(address(L1Messenger));
        emit FailedRelayedMessage(hash);
        L1Messenger.relayMessage(
            Encoding.encodeVersionedNonce({ _nonce: 0, _version: 1 }), // nonce
            sender,
            target,
            messageValue,
            0,
            selector
        );

        // The message hash is not in the successfulMessages mapping.
        assertFalse(L1Messenger.successfulMessages(hash));
        // The balance of this contract is unchanged.
        assertEq(address(this).balance, balanceBeforeThis);
        // The balance of the messenger contract is unchanged.
        assertEq(address(L1Messenger).balance, balanceBeforeMessenger);
    }
}
