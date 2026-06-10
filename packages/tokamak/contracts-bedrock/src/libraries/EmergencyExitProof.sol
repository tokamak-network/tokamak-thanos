// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { MerkleTrie } from "./trie/MerkleTrie.sol";
import { RLPReader } from "./rlp/RLPReader.sol";

/// @title EmergencyExitProof
/// @notice Verifies MPT storage proofs against a finalized L2 state root.
///         Used by AppExitCoordinator to validate user balance claims in the resolver exit path.
/// @dev Leverages the existing MerkleTrie library for Merkle-Patricia trie traversal
///      and verification. Two-stage verification:
///      1. Verify account state RLP against stateRoot using accountProof
///      2. Extract storageRoot from account state, verify storage value using storageProof
library EmergencyExitProof {
    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;

    /// @notice Error thrown when MPT proof verification fails.
    error ProofInvalid();

    /// @notice Error thrown when storage value is less than expected.
    /// @param actual   The actual storage value found.
    /// @param expected The value the caller claimed.
    error InsufficientBalance(uint256 actual, uint256 expected);

    /// @notice Computes the storage slot hash for a mapping(address => uint256).
    ///         This is the standard Solidity mapping layout: keccak256(abi.encode(user, slotIndex)).
    /// @param user       The user address to look up.
    /// @param slotIndex  The mapping slot index (e.g., 0 for a `mapping(address => uint256) balances` declaration).
    /// @return slotHash  The keccak256 hash representing the storage slot key.
    function balanceSlot(address user, uint256 slotIndex) external pure returns (bytes32 slotHash) {
        return keccak256(abi.encode(user, slotIndex));
    }

    /// @notice Verifies a user's L2 storage balance against a finalized state root.
    ///         Two-stage verification:
    ///         1. Account state RLP → stateRoot via accountProof
    ///         2. storageRoot (extracted from account state) → storage slot value via storageProof
    /// @param stateRoot        Finalized L2 state root (from L2StateOracle on L1).
    /// @param account          L2 contract address whose storage is being verified.
    /// @param accountStateRlp  RLP-encoded account state: [nonce, balance, storageRoot, codeHash].
    /// @param accountProof     MPT proof to verify accountStateRlp against stateRoot.
    /// @param storageKey       keccak256(abi.encode(user, slotIndex)) — the storage slot to verify.
    /// @param storageProof     MPT proof to verify the storage value at storageKey.
    /// @param expectedValue    The balance value the caller claims to have.
    /// @return actualValue     The actual storage value found (may differ from expectedValue).
    function verifyProofs(
        bytes32 stateRoot,
        address account,
        bytes memory accountStateRlp,
        bytes[] memory accountProof,
        bytes32 storageKey,
        bytes[] memory storageProof,
        uint256 expectedValue
    ) external view returns (uint256 actualValue) {
        // ── Step 1: Verify account state against stateRoot ─────────────────────
        // The account key in the state trie is keccak256(accountAddress).
        bytes memory accountKey = abi.encodePacked(keccak256(abi.encodePacked(account)));

        if (
            !MerkleTrie.verifyInclusionProof(
                accountKey,
                accountStateRlp,
                accountProof,
                stateRoot
            )
        ) {
            revert ProofInvalid();
        }

        // ── Step 2: Extract storageRoot from account state ────────────────────
        // OP Stack account state RLP format: [nonce, balance, storageRoot, codeHash]
        // storageRoot is at index 2.
        RLPReader.RLPItem[] memory accountFields = accountStateRlp.readList();
        if (accountFields.length < 4) revert ProofInvalid();

        bytes memory storageRootBytes = accountFields[2].readBytes();
        if (storageRootBytes.length != 32) revert ProofInvalid();
        bytes32 storageRoot = bytes32(storageRootBytes);

        // ── Step 3: Verify storage slot value against storageRoot ─────────────
        // storageProof is an MPT proof for the specific storage slot key.
        // MerkleTrie.get() returns the value as raw bytes (already RLP-decoded).
        bytes memory valueBytes = MerkleTrie.get(
            abi.encodePacked(storageKey),
            storageProof,
            storageRoot
        );

        if (valueBytes.length == 0) revert ProofInvalid();

        // Convert bytes to uint256 (the trie stores values as raw bytes after RLP decode).
        // Storage values are always 32 bytes when padded from the L2 node.
        if (valueBytes.length > 32) revert ProofInvalid();
        assembly {
            let free := mload(0x40)
            // Zero out 32 bytes
            mstore(free, 0)
            // Copy value bytes to the front
            let len := mload(valueBytes)
            pop(
                staticcall(
                    gas(),
                    0x4, // copy precompile isn't needed; just use mload/mstore
                    add(valueBytes, 0x20),
                    len,
                    free,
                    len
                )
            )
            actualValue := shr(mul(8, sub(32, len)), mload(free))
        }

        if (actualValue < expectedValue) {
            revert InsufficientBalance(actualValue, expectedValue);
        }
    }
}
