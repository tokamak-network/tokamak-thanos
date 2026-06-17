// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { SecureMerkleTrie } from "./trie/SecureMerkleTrie.sol";
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
        // SecureMerkleTrie hashes the key (keccak256(account)), matching the Ethereum state trie.
        if (
            !SecureMerkleTrie.verifyInclusionProof(
                abi.encodePacked(account),
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
        // SecureMerkleTrie hashes the slot key (keccak256(storageKey)), matching the storage trie.
        bytes memory valueBytes = SecureMerkleTrie.get(
            abi.encodePacked(storageKey),
            storageProof,
            storageRoot
        );

        if (valueBytes.length == 0) revert ProofInvalid();

        // The storage trie leaf stores RLP(slotValue), so decode one more level to recover the raw
        // big-endian value. Parsing the RLP-encoded blob directly would inflate the result by the
        // length-prefix byte for any value >= 0x80, allowing an over-claim.
        bytes memory slotValue = valueBytes.readBytes();
        if (slotValue.length > 32) revert ProofInvalid();
        assembly {
            let free := mload(0x40)
            // Zero out 32 bytes
            mstore(free, 0)
            // Copy the trimmed value bytes to the front, then right-align into a uint256.
            let len := mload(slotValue)
            pop(
                staticcall(
                    gas(),
                    0x4, // identity precompile: copy len bytes
                    add(slotValue, 0x20),
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
