// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "forge-std/Test.sol";
import "../../../src/tokamak/AA/VerifyingPaymaster.sol";
import "../../../src/tokamak/AA/EntryPoint.sol";

contract VerifyingPaymasterTest is Test {
    EntryPoint entryPoint;
    VerifyingPaymaster paymaster;
    address signer = address(0x5166E12);

    function setUp() public {
        entryPoint = new EntryPoint();
        paymaster = new VerifyingPaymaster();
        paymaster.initialize(IEntryPoint(address(entryPoint)), signer);
    }

    function test_Initialize_SetsEntryPoint() public view {
        assertEq(address(paymaster.entryPoint()), address(entryPoint));
    }

    function test_Initialize_SetsSigner() public view {
        assertEq(paymaster.verifyingSigner(), signer);
    }

    function test_Initialize_CannotCallTwice() public {
        vm.expectRevert();
        paymaster.initialize(IEntryPoint(address(entryPoint)), signer);
    }

    function test_Initialize_ZeroSignerReverts() public {
        VerifyingPaymaster pm2 = new VerifyingPaymaster();
        vm.expectRevert("VerifyingPaymaster: zero signer");
        pm2.initialize(IEntryPoint(address(entryPoint)), address(0));
    }
}
