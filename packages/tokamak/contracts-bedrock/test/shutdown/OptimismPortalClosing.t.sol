// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";

import { OptimismPortalClosing } from "src/shutdown/OptimismPortalClosing.sol";

contract OptimismPortalClosing_ImplMock {
    uint256 public value;

    function setValue(uint256 _value) external {
        value = _value;
    }

    function ping() external pure returns (uint256) {
        return 123;
    }
}

contract OptimismPortalClosing_Test is Test {
    OptimismPortalClosing internal closing;
    OptimismPortalClosing_ImplMock internal impl;

    function setUp() public {
        impl = new OptimismPortalClosing_ImplMock();
        closing = new OptimismPortalClosing(address(impl));
    }

    function test_constructor_revertsOnZeroImplementation() public {
        vm.expectRevert("OptimismPortalClosing: implementation is zero");
        new OptimismPortalClosing(address(0));
    }

    function test_version_isClosing() public {
        assertEq(closing.version(), "2.8.1-closing");
    }

    function test_depositTransaction_reverts() public {
        vm.expectRevert("OptimismPortal: deposits are disabled due to chain shutdown");
        closing.depositTransaction(address(1), 0, 0, 0, false, hex"");
    }

    function test_onApprove_reverts() public {
        vm.expectRevert("OptimismPortal: deposits are disabled due to chain shutdown");
        closing.onApprove(address(1), address(2), 0, hex"");
    }

    function test_receive_reverts() public {
        vm.deal(address(this), 1 ether);
        vm.expectRevert("OptimismPortal: deposits are disabled due to chain shutdown");
        payable(address(closing)).transfer(1 wei);
    }

    function test_fallback_delegatesToImplementation() public {
        (bool okSet, ) = address(closing).call(abi.encodeWithSignature("setValue(uint256)", 456));
        assertTrue(okSet);

        (bool okGet, bytes memory data) = address(closing).call(abi.encodeWithSignature("value()"));
        assertTrue(okGet);
        assertEq(abi.decode(data, (uint256)), 456);

        (bool okPing, bytes memory pingData) = address(closing).call(abi.encodeWithSignature("ping()"));
        assertTrue(okPing);
        assertEq(abi.decode(pingData, (uint256)), 123);
    }
}
