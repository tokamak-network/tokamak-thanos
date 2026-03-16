// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "forge-std/Test.sol";
import "../../../src/tokamak/AA/EntryPoint.sol";

contract EntryPointTest is Test {
    EntryPoint entryPoint;

    function setUp() public {
        entryPoint = new EntryPoint();
    }

    function test_DepositTo() public {
        address alice = address(0xABCD);
        vm.deal(alice, 1 ether);
        vm.prank(alice);
        entryPoint.depositTo{value: 0.5 ether}(alice);
        IStakeManager.DepositInfo memory info = entryPoint.getDepositInfo(alice);
        assertEq(info.deposit, 0.5 ether);
    }

    function test_NativeTokenIsAccepted() public {
        // Thanos: msg.value = TON. This test verifies native token (TON) is deposited.
        vm.deal(address(this), 1 ether);
        entryPoint.depositTo{value: 1 ether}(address(this));
        IStakeManager.DepositInfo memory info = entryPoint.getDepositInfo(address(this));
        assertEq(info.deposit, 1 ether);
    }
}
