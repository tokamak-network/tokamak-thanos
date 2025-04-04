// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockGnosisSafe {
  address[] private owners;
  uint256 private threshold;

  constructor(address[] memory _owners, uint256 _threshold) {
    owners = _owners;
    threshold = _threshold;
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
}