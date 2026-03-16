// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "../../../src/tokamak/AA/Simple7702Account.sol";

contract Simple7702AccountTest is Test {
    Simple7702Account impl;

    // Thanos predeploy address for EntryPoint v0.8
    address constant ENTRY_POINT_PREDEPLOY = 0x4200000000000000000000000000000000000063;

    function setUp() public {
        impl = new Simple7702Account();
    }

    function test_EntryPointIsThanosPredeployAddress() public view {
        // CRITICAL: must return Thanos predeploy, not eth-infinitism mainnet address
        assertEq(address(impl.entryPoint()), ENTRY_POINT_PREDEPLOY);
    }

    function test_OnlyEntryPointOrOwnerCanExecute() public {
        address nonEntryPoint = address(0xDEAD);
        vm.prank(nonEntryPoint);
        vm.expectRevert();
        impl.execute(address(0), 0, "");
    }

    function test_ExecuteBatchFromEntryPoint() public {
        BaseAccount.Call[] memory calls = new BaseAccount.Call[](1);
        calls[0] = BaseAccount.Call({
            target: address(0x1234),
            value: 0,
            data: ""
        });

        vm.prank(ENTRY_POINT_PREDEPLOY);
        impl.executeBatch(calls);
    }
}
