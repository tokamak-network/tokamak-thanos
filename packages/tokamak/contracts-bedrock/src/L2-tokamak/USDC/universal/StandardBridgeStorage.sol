// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract StandardBridgeStorage  {
    address public messenger;
    address public otherBridge;
    address public l1Usdc;
    address public l2Usdc;

    modifier onlyL1Usdc(address token) {
        require(token == l1Usdc, "not L1 usdc");
        _;
    }

    modifier onlyL2Usdc(address token) {
        require(token == l2Usdc, "not L2 usdc");
        _;
    }
}