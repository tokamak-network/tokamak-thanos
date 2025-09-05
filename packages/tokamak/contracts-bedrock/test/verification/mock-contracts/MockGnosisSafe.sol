// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract MockGnosisSafe {
  address[] private owners;
  uint256 private threshold;

  // Storage slot for fallback handler to match Safe v1.5.0
  // bytes32 private constant FALLBACK_HANDLER_STORAGE_SLOT = 0x6c9a6c4a39284e37ed1cf53d337577d14212a4870fb976a4366c693b939918d5;
  address private fallbackHandler;

  address[] private modules;

  constructor(address[] memory _owners, uint256 _threshold) {
    owners = _owners;
    threshold = _threshold;
    // Default fallback handler to address(0)
    fallbackHandler = address(0);
  }

  function setFallbackHandler(address _fallbackHandler) external {
    fallbackHandler = _fallbackHandler;
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
    // Return empty modules array and SENTINEL_MODULES as next address to match Safe expectations
    // SENTINEL_MODULES = address(0x1) as defined in the verification contract
    return (modules, address(0x1)); 
  }

  function addModule(address _module) external {
    modules.push(_module);
  }

  function getFallbackHandler() external view returns (address) {
    return fallbackHandler;
  }

  // New function to support Safe v1.5.0 storage access pattern
  function getStorageAt(uint256 offset, uint256 length) external view returns (bytes memory) {
    // Simple implementation for testing - in real Safe this reads from storage slots
    if (offset == uint256(0x6c9a6c4a39284e37ed1cf53d337577d14212a4870fb976a4366c693b939918d5) && length == 1) {
      // Return the fallback handler address for the specific storage slot
      bytes memory result = new bytes(32);
      assembly {
        mstore(add(result, 0x20), sload(fallbackHandler.slot))
      }
      return result;
    }
    return new bytes(length * 32);
  }
}
