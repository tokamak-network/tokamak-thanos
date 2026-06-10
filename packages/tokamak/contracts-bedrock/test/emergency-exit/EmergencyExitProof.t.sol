// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test } from "forge-std/Test.sol";
import { EmergencyExitProof } from "src/libraries/EmergencyExitProof.sol";
import { MptProofBuilder } from "./MptProofBuilder.sol";

/// @title EmergencyExitProof_Test
/// @notice Tests for EmergencyExitProof — slot computation and two-stage MPT proof verification.
/// @dev verifyProofs is exercised against real single-leaf secure tries built by MptProofBuilder,
///      so the tests fail unless the keys are keccak256-hashed (SecureMerkleTrie) AND the storage
///      value is RLP-decoded one extra level.
contract EmergencyExitProof_Test is Test {
    address internal l2Token = makeAddr("l2Token");
    address internal user = makeAddr("user");

    // ── balanceSlot ──────────────────────────────────────────────────────────

    function test_balanceSlot_Deterministic() public {
        address u = address(0x1234);
        bytes32 result1 = EmergencyExitProof.balanceSlot(u, 0);
        assertEq(result1, EmergencyExitProof.balanceSlot(u, 0), "should be deterministic");
        assert(result1 != EmergencyExitProof.balanceSlot(u, 1));
        assert(result1 != EmergencyExitProof.balanceSlot(address(0x5678), 0));
    }

    function test_balanceSlot_Correctness() public {
        assertEq(
            EmergencyExitProof.balanceSlot(address(0xABCD), 5),
            keccak256(abi.encode(address(0xABCD), uint256(5))),
            "should match keccak256(abi.encode)"
        );
    }

    // ── verifyProofs — real two-stage MPT verification ────────────────────────

    /// @notice A valid proof for a multi-byte balance returns the exact balance.
    ///         balance = 1000 (0x03E8) → storage leaf value = RLP(0x03E8) = 0x8203E8.
    ///         A naive parser reading 0x8203E8 as an integer yields 8_521_704, so this passes
    ///         only when the storage value is RLP-decoded correctly.
    function test_verifyProofs_ValidProof_ReturnsExactBalance() public {
        uint256 balance = 1000;
        MptProofBuilder.Fixture memory f = MptProofBuilder.build(l2Token, user, 0, balance);

        uint256 got = EmergencyExitProof.verifyProofs(
            f.stateRoot, f.account, f.accountStateRlp, f.accountProof, f.storageKey, f.storageProof, balance
        );
        assertEq(got, balance, "must return the exact, RLP-decoded balance");
    }

    /// @notice Claiming more than the proven balance reverts (no over-claim).
    function test_verifyProofs_OverClaim_Revert() public {
        uint256 balance = 1000;
        MptProofBuilder.Fixture memory f = MptProofBuilder.build(l2Token, user, 0, balance);

        vm.expectRevert(
            abi.encodeWithSelector(EmergencyExitProof.InsufficientBalance.selector, balance, balance + 1)
        );
        EmergencyExitProof.verifyProofs(
            f.stateRoot, f.account, f.accountStateRlp, f.accountProof, f.storageKey, f.storageProof, balance + 1
        );
    }

    /// @notice A tampered storage proof (root mismatch) reverts.
    function test_verifyProofs_TamperedStorageProof_Revert() public {
        MptProofBuilder.Fixture memory f = MptProofBuilder.build(l2Token, user, 0, 1000);
        f.storageProof[0][f.storageProof[0].length - 1] ^= bytes1(0x01);

        vm.expectRevert();
        EmergencyExitProof.verifyProofs(
            f.stateRoot, f.account, f.accountStateRlp, f.accountProof, f.storageKey, f.storageProof, 1000
        );
    }

    /// @notice A realistic 18-decimal balance is decoded exactly.
    function test_verifyProofs_LargeBalance_ReturnsExactBalance() public {
        uint256 balance = 1e18;
        MptProofBuilder.Fixture memory f = MptProofBuilder.build(l2Token, user, 3, balance);

        uint256 got = EmergencyExitProof.verifyProofs(
            f.stateRoot, f.account, f.accountStateRlp, f.accountProof, f.storageKey, f.storageProof, balance
        );
        assertEq(got, balance, "must decode an 18-decimal balance exactly");
    }
}
