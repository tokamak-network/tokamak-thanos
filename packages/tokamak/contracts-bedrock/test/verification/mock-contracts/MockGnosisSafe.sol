// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockGnosisSafe {
  address[] private owners;
  uint256 private threshold;
  address private fallbackHandler;
  // Storage slot is defined in L1ContractVerification as:
  // bytes32 private constant FALLBACK_HANDLER_STORAGE_SLOT = 0x6c9a6c4a39284e37ed1cf53d337577d14212a4870fb976a4366c693b939918d5;

  // We'll store the fallback handler directly in storage to match the real implementation
  // rather than using a regular state variable

  address[] private modules;

  constructor(address[] memory _owners, uint256 _threshold) {
    owners = _owners;
    threshold = _threshold;
    // Default fallback handler to address(0)
    // We don't need to initialize the storage slot - it defaults to 0
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

  function getModulesPaginated(address, uint256) external view returns (address[] memory) {
    return modules; // Returns an empty array by default
  }

  function addModule(address _module) external {
    modules.push(_module);
  }

  function getFallbackHandler() external view returns (address) {
    return fallbackHandler;
  }
}
