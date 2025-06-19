// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IMultiSigWallet {
    function submitTransaction(
        address _to,
        uint _value,
        bytes memory _data
    ) external;

    function confirmTransaction(uint _txIndex) external;

    function executeTransaction(uint _txIndex) external;

    function getTransactionCount() external view returns (uint);
}