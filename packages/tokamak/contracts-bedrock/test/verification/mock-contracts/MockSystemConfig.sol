// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract MockSystemConfig {
  address public l1StandardBridge;
  address public l1CrossDomainMessenger;
  address public optimismPortal;
  address public nativeTokenAddress;

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
  }
}
