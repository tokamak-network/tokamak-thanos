// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockSystemConfig {
  address public l1StandardBridge;
  address public l1CrossDomainMessenger;
  address public optimismPortal;
  address public nativeTokenAddress;

  constructor(
    address _l1StandardBridge,
    address _l1CrossDomainMessenger,
    address _optimismPortal
  ) {
    l1StandardBridge = _l1StandardBridge;
    l1CrossDomainMessenger = _l1CrossDomainMessenger;
    optimismPortal = _optimismPortal;
    nativeTokenAddress = address(0x123); // Default value
  }

  function setNativeTokenAddress(address _nativeTokenAddress) external {
    nativeTokenAddress = _nativeTokenAddress;
  }
}