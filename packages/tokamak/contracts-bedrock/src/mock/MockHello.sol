// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract MockHello {

    uint256 public blockNumber ;
    string public message ;

    constructor() {
    }

    function say(string memory _msg) external {
        message = _msg;
        blockNumber = block.number;
    }

}