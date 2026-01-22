// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IL1UsdcBridge {
  function deposits(address l1Token, address l2Token) external view returns (uint256);
  function l1Usdc() external view returns (address);
  function l2Usdc() external view returns (address);
}
