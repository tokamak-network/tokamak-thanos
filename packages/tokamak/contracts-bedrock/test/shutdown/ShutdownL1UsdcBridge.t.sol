// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test, stdStorage, StdStorage } from "forge-std/Test.sol";
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import { ShutdownL1UsdcBridge } from "src/shutdown/ShutdownL1UsdcBridge.sol";

contract ShutdownL1UsdcBridge_AdminOwnerMock {
    address public owner;

    constructor(address _owner) {
        owner = _owner;
    }
}

contract ShutdownL1UsdcBridge_AdminNoOwnerMock {}

contract ShutdownL1UsdcBridge_TokenMock is ERC20 {
    constructor() ERC20("USDC Mock", "USDC") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract ShutdownL1UsdcBridge_Test is Test {
    using stdStorage for StdStorage;

    bytes32 internal constant ADMIN_SLOT =
        0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    ShutdownL1UsdcBridge internal bridge;
    ShutdownL1UsdcBridge_TokenMock internal token;

    address internal admin;
    address internal recipient;

    function setUp() public {
        bridge = new ShutdownL1UsdcBridge();
        token = new ShutdownL1UsdcBridge_TokenMock();

        admin = makeAddr("admin");
        recipient = makeAddr("recipient");
    }

    function test_sweepUSDC_revertsWhenAdminNotSet() public {
        vm.expectRevert("ShutdownL1UsdcBridge: admin not set");
        bridge.sweepUSDC(recipient, 1);
    }

    function test_sweepUSDC_revertsWhenL1UsdcUnset() public {
        _setProxyAdmin(address(bridge), admin);

        vm.prank(admin);
        vm.expectRevert("ShutdownL1UsdcBridge: l1Usdc not set");
        bridge.sweepUSDC(recipient, 1);
    }

    function test_sweepUSDC_succeedsForAdminEOA() public {
        _setProxyAdmin(address(bridge), admin);
        _setL1Usdc(address(bridge), address(token));
        token.mint(address(bridge), 100);

        vm.prank(admin);
        bridge.sweepUSDC(recipient, 40);

        assertEq(token.balanceOf(recipient), 40);
        assertEq(token.balanceOf(address(bridge)), 60);
    }

    function test_sweepUSDC_succeedsForAdminOwnerContract() public {
        ShutdownL1UsdcBridge_AdminOwnerMock adminContract = new ShutdownL1UsdcBridge_AdminOwnerMock(admin);
        _setProxyAdmin(address(bridge), address(adminContract));
        _setL1Usdc(address(bridge), address(token));
        token.mint(address(bridge), 50);

        vm.prank(admin);
        bridge.sweepUSDC(recipient, 50);

        assertEq(token.balanceOf(recipient), 50);
        assertEq(token.balanceOf(address(bridge)), 0);
    }

    function test_sweepUSDC_revertsForNonOwner() public {
        ShutdownL1UsdcBridge_AdminOwnerMock adminContract = new ShutdownL1UsdcBridge_AdminOwnerMock(admin);
        _setProxyAdmin(address(bridge), address(adminContract));
        _setL1Usdc(address(bridge), address(token));

        vm.prank(makeAddr("notOwner"));
        vm.expectRevert("ShutdownL1UsdcBridge: unauthorized");
        bridge.sweepUSDC(recipient, 1);
    }

    function test_sweepUSDC_revertsWhenAdminOwnerLookupFails() public {
        ShutdownL1UsdcBridge_AdminNoOwnerMock adminContract = new ShutdownL1UsdcBridge_AdminNoOwnerMock();
        _setProxyAdmin(address(bridge), address(adminContract));
        _setL1Usdc(address(bridge), address(token));

        vm.prank(admin);
        vm.expectRevert("ShutdownL1UsdcBridge: admin owner lookup failed");
        bridge.sweepUSDC(recipient, 1);
    }

    function _setProxyAdmin(address target, address admin_) internal {
        vm.store(target, ADMIN_SLOT, bytes32(uint256(uint160(admin_))));
    }

    function _setL1Usdc(address target, address token_) internal {
        stdstore.target(target).sig("l1Usdc()").checked_write(token_);
    }
}
