// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

/// @title MptProofBuilder
/// @notice Test helper that builds real single-leaf Merkle-Patricia trie proofs, using the same
///         encoding a live Ethereum node produces:
///           - account trie:  key = keccak256(account),        value = RLP([nonce, balance, storageRoot, codeHash])
///           - storage trie:  key = keccak256(storageSlotKey),  value = RLP(slotValue)
///         Both tries are "secure" (keys hashed with keccak256). Consumers verify that
///         EmergencyExitProof / AppExitCoordinator hash the keys and RLP-decode the storage value.
library MptProofBuilder {
    struct Fixture {
        bytes32 stateRoot;
        address account;
        bytes accountStateRlp;
        bytes[] accountProof;
        bytes32 storageKey;
        bytes[] storageProof;
    }

    /// @notice Builds a coupled account+storage proof for `user`'s balance at `slotIndex` of `account`.
    function build(address account, address user, uint256 slotIndex, uint256 balance)
        internal
        pure
        returns (Fixture memory f)
    {
        f.account = account;
        f.storageKey = keccak256(abi.encode(user, slotIndex));

        // ── Storage trie (single leaf): secure key = keccak256(storageKey), value = RLP(balance) ──
        bytes32 secureStorageKey = keccak256(abi.encodePacked(f.storageKey));
        bytes memory storageLeaf = _leaf(secureStorageKey, _rlpStr(_trim(balance)));
        bytes32 storageRoot = keccak256(storageLeaf);
        f.storageProof = new bytes[](1);
        f.storageProof[0] = storageLeaf;

        // ── Account trie (single leaf): RLP([nonce=0, balance=0, storageRoot, codeHash]) ──
        f.accountStateRlp = _rlpList(
            bytes.concat(
                _rlpStr(""), // nonce
                _rlpStr(""), // balance
                _rlpStr(abi.encodePacked(storageRoot)),
                _rlpStr(abi.encodePacked(keccak256("")))
            )
        );
        bytes memory accountLeaf = _leaf(keccak256(abi.encodePacked(account)), f.accountStateRlp);
        f.stateRoot = keccak256(accountLeaf);
        f.accountProof = new bytes[](1);
        f.accountProof[0] = accountLeaf;
    }

    /// @notice Encodes a single even-length leaf node: RLP([0x20 || key32, valueField]).
    function _leaf(bytes32 secureKey, bytes memory valueField) private pure returns (bytes memory) {
        bytes memory path = bytes.concat(hex"20", secureKey); // HP prefix 0x20 = leaf, even length
        return _rlpList(bytes.concat(_rlpStr(path), _rlpStr(valueField)));
    }

    // ── Minimal RLP encoder (lengths < 256, sufficient for single-leaf tries) ──

    function _rlpStr(bytes memory b) private pure returns (bytes memory) {
        if (b.length == 1 && uint8(b[0]) < 0x80) return b;
        if (b.length <= 55) return bytes.concat(bytes1(uint8(0x80 + b.length)), b);
        return bytes.concat(hex"b8", bytes1(uint8(b.length)), b);
    }

    function _rlpList(bytes memory items) private pure returns (bytes memory) {
        if (items.length <= 55) return bytes.concat(bytes1(uint8(0xc0 + items.length)), items);
        return bytes.concat(hex"f8", bytes1(uint8(items.length)), items);
    }

    /// @notice Big-endian byte representation of `v` with leading zero bytes stripped.
    function _trim(uint256 v) private pure returns (bytes memory out) {
        if (v == 0) return out;
        uint256 len;
        uint256 t = v;
        while (t != 0) {
            len++;
            t >>= 8;
        }
        out = new bytes(len);
        for (uint256 i; i < len; i++) {
            out[len - 1 - i] = bytes1(uint8(v >> (8 * i)));
        }
    }
}
