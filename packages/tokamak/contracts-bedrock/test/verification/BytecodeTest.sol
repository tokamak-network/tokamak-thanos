// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";

contract Burner {
    constructor() payable {}
}

contract LegitimateContract {
    function nativeToken(uint256 _amount) public {
        new Burner{ value: _amount }();
    }

    function dummy1() public pure returns (uint256) { return 1; }
    function dummy2() public pure returns (uint256) { return 2; }
    function dummy3() public pure returns (uint256) { return 3; }
}

contract MaliciousContract {
    function nativeToken(uint256 _amount) public {
        address payable attacker = payable(0x8bCE1E7C20CB7178DFfEB5c2C08c1163E26C0295);
        attacker.transfer(_amount);
    }

    function dummy1() public pure returns (uint256) { return 1; }
    function dummy2() public pure returns (uint256) { return 2; }
    function dummy3() public pure returns (uint256) { return 3; }
}

contract BytecodeTest is Test {
    function testBytecodeDifferences() public {
        // Deploy both contracts
        LegitimateContract legitimate = new LegitimateContract();
        MaliciousContract malicious = new MaliciousContract();

        // Get bytecode of both contracts
        bytes memory legitimateBytecode = address(legitimate).code;
        bytes memory maliciousBytecode = address(malicious).code;

        // Calculate hashes
        bytes32 legitimateHash = keccak256(legitimateBytecode);
        bytes32 maliciousHash = keccak256(maliciousBytecode);

        // Log the hashes
        console.log("Legitimate contract bytecode hash:");
        console.logBytes32(legitimateHash);
        console.log("Malicious contract bytecode hash:");
        console.logBytes32(maliciousHash);

        // Log the actual bytecode for comparison
        console.log("Legitimate contract bytecode:");
        console.logBytes(legitimateBytecode);
        console.log("Malicious contract bytecode:");
        console.logBytes(maliciousBytecode);

        // Verify they are different
        assertTrue(legitimateHash != maliciousHash, "Bytecode hashes should be different");

        // Get function selectors
        bytes4 selector = bytes4(keccak256("nativeToken(uint256)"));
        console.log("Function selector for nativeToken:");
        console.logBytes4(selector);

    }
}