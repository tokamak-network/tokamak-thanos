// SPDX-License-Identifier: MIT
pragma solidity 0.8.22;

contract Test {
    address public myAddress;
    string public myString;
    uint256 public myNumber;

    constructor(address a, string memory b, uint256 c) {
        myAddress = a;
        myString = b;
        myNumber = c;
    }
}
