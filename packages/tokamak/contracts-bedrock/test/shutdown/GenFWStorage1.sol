// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title GenFWStorage1
/// @notice Test contract for storing force withdrawal hashes with dynamic function dispatch
/// @dev Uses fallback with payable to allow hash retrieval via function calls
contract GenFWStorage1 {
    mapping(bytes4 => bytes32) private hashes;

    function setHash(bytes4 functionSig, bytes32 value) external {
        hashes[functionSig] = value;
    }

    function getHash(bytes4 functionSig) external view returns (bytes32) {
        return hashes[functionSig];
    }

    fallback() external payable {
        bytes4 sig = msg.sig;
        bytes32 value = hashes[sig];

        assembly {
            mstore(0, value)
            return(0, 32)
        }
    }
}
