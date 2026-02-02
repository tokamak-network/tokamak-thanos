// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title GenFWStorage
/// @notice Contract for storing force withdrawal hashes with dynamic function dispatch
/// @dev Uses fallback for hash retrieval via function calls. ETH transfers are explicitly rejected.
///      Only the owner can set hashes. Ownership can be renounced after setup.
contract GenFWStorage {
    /// @notice Contract owner who can set hashes
    address public owner;

    /// @notice Mapping from function signature to hash value
    mapping(bytes4 => bytes32) private hashes;

    /// @notice Emitted when a hash is set
    event HashSet(bytes4 indexed functionSig, bytes32 value);

    /// @notice Emitted when ownership is transferred
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /// @notice Restricts function access to owner only
    modifier onlyOwner() {
        require(msg.sender == owner, "GenFWStorage: caller is not owner");
        _;
    }

    /// @notice Sets the deployer as the initial owner
    constructor() {
        owner = msg.sender;
        emit OwnershipTransferred(address(0), msg.sender);
    }

    /// @notice Sets a hash value for a given function signature
    /// @param functionSig The 4-byte function signature
    /// @param value The hash value to store
    function setHash(bytes4 functionSig, bytes32 value) external onlyOwner {
        hashes[functionSig] = value;
        emit HashSet(functionSig, value);
    }

    /// @notice Returns the hash value for a given function signature
    /// @param functionSig The 4-byte function signature
    /// @return The stored hash value
    function getHash(bytes4 functionSig) external view returns (bytes32) {
        return hashes[functionSig];
    }

    /// @notice Transfers ownership to a new address
    /// @param newOwner The address of the new owner
    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "GenFWStorage: new owner is zero address");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    /// @notice Renounces ownership, making the contract immutable
    /// @dev After calling this, no more hashes can be set
    function renounceOwnership() external onlyOwner {
        emit OwnershipTransferred(owner, address(0));
        owner = address(0);
    }

    /// @notice Returns stored hash for the given function signature
    /// @dev Called via staticcall from ForceWithdrawBridge, no ETH handling needed.
    ///      ETH transfers are rejected because fallback is non-payable.
    fallback() external {
        bytes4 sig = msg.sig;
        bytes32 value = hashes[sig];

        assembly {
            mstore(0, value)
            return(0, 32)
        }
    }
}
