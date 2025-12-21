// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract MockGnosisSafe {
    address[] private owners;
    uint256 private threshold;
    address private fallbackHandler;
    // Storage slot is defined in L1ContractVerification as:
    // bytes32 private constant FALLBACK_HANDLER_STORAGE_SLOT =
    // 0x6c9a6c4a39284e37ed1cf53d337577d14212a4870fb976a4366c693b939918d5;

    // We'll store the fallback handler directly in storage to match the real implementation
    // rather than using a regular state variable

    address[] private modules;

    // Sentinel value used by Gnosis Safe
    address private constant SENTINEL_MODULES = address(0x1);

    // Storage slot for fallback handler
    bytes32 private constant FALLBACK_HANDLER_STORAGE_SLOT =
        0x6c9a6c4a39284e37ed1cf53d337577d14212a4870fb976a4366c693b939918d5;

    constructor(address[] memory _owners, uint256 _threshold) {
        owners = _owners;
        threshold = _threshold;
        // Default fallback handler to address(0)
        // We don't need to initialize the storage slot - it defaults to 0
    }

    function setFallbackHandler(address _fallbackHandler) external {
        // Store in the correct storage slot
        assembly {
            sstore(FALLBACK_HANDLER_STORAGE_SLOT, _fallbackHandler)
        }
    }

    function getOwners() external view returns (address[] memory) {
        return owners;
    }

    function getThreshold() external view returns (uint256) {
        return threshold;
    }

    function masterCopy() external view returns (address) {
        return address(this);
    }

    function getModulesPaginated(address, uint256) external view returns (address[] memory, address) {
        // When there are no modules, return empty array and SENTINEL_MODULES as next pointer
        // When there are modules, return modules array and address(0) as next (indicating there are more)
        if (modules.length == 0) {
            return (modules, SENTINEL_MODULES);
        } else {
            return (modules, address(0));
        }
    }

    function addModule(address _module) external {
        modules.push(_module);
    }

    function getFallbackHandler() external view returns (address) {
        address handler;
        assembly {
            handler := sload(FALLBACK_HANDLER_STORAGE_SLOT)
        }
        return handler;
    }

    /**
     * @notice Read storage at a given slot
     * @dev Implements IStorageAccessible to allow reading from storage slots
     * @param offset The storage slot to read from
     * @param length Number of 32-byte words to read
     * @return The data read from storage
     */
    function getStorageAt(uint256 offset, uint256 length) external view returns (bytes memory) {
        bytes memory result = new bytes(length * 32);

        assembly {
            let ptr := add(result, 0x20)
            for { let i := 0 } lt(i, length) { i := add(i, 1) } {
                let value := sload(add(offset, i))
                mstore(add(ptr, mul(i, 32)), value)
            }
        }

        return result;
    }
}
