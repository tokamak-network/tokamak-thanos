// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test, stdStorage, StdStorage } from "forge-std/Test.sol";
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import { ShutdownOptimismPortal2 } from "src/shutdown/ShutdownOptimismPortal2.sol";

contract ShutdownOptimismPortal2_AdminOwnerMock {
    address public owner;

    constructor(address _owner) {
        owner = _owner;
    }
}

contract ShutdownOptimismPortal2_AdminNoOwnerMock {}

contract ShutdownOptimismPortal2_TokenMock is ERC20 {
    constructor() ERC20("Native Token", "NATIVE") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract ShutdownOptimismPortal2_SystemConfigMock {
    address internal token;

    constructor(address _token) {
        token = _token;
    }

    function nativeTokenAddress() external view returns (address) {
        return token;
    }

    function setNativeTokenAddress(address _token) external {
        token = _token;
    }
}

contract ShutdownOptimismPortal2_Test is Test {
    using stdStorage for StdStorage;

    bytes32 internal constant ADMIN_SLOT =
        0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    ShutdownOptimismPortal2 internal portal;
    ShutdownOptimismPortal2_TokenMock internal token;

    address internal admin;
    address internal recipient;

    function setUp() public {
        portal = new ShutdownOptimismPortal2(0, 0);
        token = new ShutdownOptimismPortal2_TokenMock();
        admin = makeAddr("admin");
        recipient = makeAddr("recipient");
    }

    function test_sweepNativeToken_revertsWhenAdminNotSet() public {
        vm.expectRevert("ShutdownOptimismPortal2: admin not set");
        portal.sweepNativeToken(recipient, 1);
    }

    function test_sweepNativeToken_revertsWhenNativeTokenIsETH() public {
        _setProxyAdmin(address(portal), admin);
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal2_SystemConfigMock(address(0))));

        vm.prank(admin);
        vm.expectRevert("ShutdownOptimismPortal2: native token is ETH");
        portal.sweepNativeToken(recipient, 1);
    }

    function test_sweepNativeToken_succeedsForAdminEOA() public {
        _setProxyAdmin(address(portal), admin);
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal2_SystemConfigMock(address(token))));
        token.mint(address(portal), 100);

        vm.prank(admin);
        portal.sweepNativeToken(recipient, 60);

        assertEq(token.balanceOf(recipient), 60);
        assertEq(token.balanceOf(address(portal)), 40);
    }

    function test_sweepNativeToken_succeedsForAdminOwnerContract() public {
        ShutdownOptimismPortal2_AdminOwnerMock adminContract = new ShutdownOptimismPortal2_AdminOwnerMock(admin);
        _setProxyAdmin(address(portal), address(adminContract));
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal2_SystemConfigMock(address(token))));
        token.mint(address(portal), 75);

        vm.prank(admin);
        portal.sweepNativeToken(recipient, 75);

        assertEq(token.balanceOf(recipient), 75);
        assertEq(token.balanceOf(address(portal)), 0);
    }

    function test_sweepNativeToken_revertsForNonOwner() public {
        ShutdownOptimismPortal2_AdminOwnerMock adminContract = new ShutdownOptimismPortal2_AdminOwnerMock(admin);
        _setProxyAdmin(address(portal), address(adminContract));
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal2_SystemConfigMock(address(token))));

        vm.prank(makeAddr("notOwner"));
        vm.expectRevert("ShutdownOptimismPortal2: unauthorized");
        portal.sweepNativeToken(recipient, 1);
    }

    function test_sweepNativeToken_revertsWhenAdminOwnerLookupFails() public {
        ShutdownOptimismPortal2_AdminNoOwnerMock adminContract = new ShutdownOptimismPortal2_AdminNoOwnerMock();
        _setProxyAdmin(address(portal), address(adminContract));
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal2_SystemConfigMock(address(token))));

        vm.prank(admin);
        vm.expectRevert("ShutdownOptimismPortal2: admin owner lookup failed");
        portal.sweepNativeToken(recipient, 1);
    }

    function _setProxyAdmin(address target, address admin_) internal {
        vm.store(target, ADMIN_SLOT, bytes32(uint256(uint160(admin_))));
    }

    function _setSystemConfig(address target, address config_) internal {
        stdstore.target(target).sig("systemConfig()").checked_write(config_);
    }
}
