// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IL1CrossDomainMessenger {
    function sendNativeTokenMessage(
        address _target,
        uint256 _amount,
        bytes calldata _message,
        uint32 _minGasLimit
    )
        external;

    function sendMessage(
        address _target, 
        bytes calldata _message, 
        uint32 _minGasLimit
    ) 
        external 
        payable;
}