// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract MockSystemConfig {
  address public l1StandardBridge;
  address public l1CrossDomainMessenger;
  address public optimismPortal;
  address public nativeTokenAddress;
  address public _owner;

  constructor() {
    // Constructor is not called when using TransparentUpgradeableProxy
  }

  function setNativeTokenAddress(address _nativeTokenAddress) external {
    nativeTokenAddress = _nativeTokenAddress;
  }

  // Add an initialize function to set all values through the proxy
  function initialize(
    address _l1StandardBridge,
    address _l1CrossDomainMessenger,
    address _optimismPortal,
    address _nativeTokenAddress
  ) external {
    l1StandardBridge = _l1StandardBridge;
    l1CrossDomainMessenger = _l1CrossDomainMessenger;
    optimismPortal = _optimismPortal;
    nativeTokenAddress = _nativeTokenAddress;

    // Set the owner to msg.sender (the caller through the proxy)
    _owner = msg.sender;
  }

  // Implement owner function for the verification contract
  function owner() external view returns (address) {
    return _owner;
  }

  // Function to set owner for testing purposes
  function setOwner(address newOwner) external {
    _owner = newOwner;
  }
}
