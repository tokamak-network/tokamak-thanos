// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IL1StandardBridge {
  function deposits(
    address l1Token,
    address l2Token
  ) external view returns (uint256);
}
