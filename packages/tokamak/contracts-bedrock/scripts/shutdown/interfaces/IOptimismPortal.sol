// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IOptimismPortal {
  function finalizedWithdrawals(bytes32 withdrawalHash) external view returns (bool);

  function proofSubmitters(
    bytes32 withdrawalHash,
    uint256 index
  ) external view returns (address);
}
