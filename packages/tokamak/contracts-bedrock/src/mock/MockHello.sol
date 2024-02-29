// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract MockHello {

    uint256 public blockNumber ;
    string public message ;

    constructor() {
    }

    receive() external payable {

    }

    function say(string memory _msg) external {
        message = _msg;
        blockNumber = block.number;
    }

    function sayPayable(string memory _msg) external payable {
        message = _msg;
        blockNumber = block.number;
    }
}