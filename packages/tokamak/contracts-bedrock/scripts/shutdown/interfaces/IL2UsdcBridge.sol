// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IL2UsdcBridge {
  function l1Usdc() external view returns (address);
  function l2Usdc() external view returns (address);
}
