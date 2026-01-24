// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test, stdStorage, StdStorage } from "forge-std/Test.sol";
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import { ShutdownOptimismPortal } from "src/shutdown/ShutdownOptimismPortal.sol";

contract ShutdownOptimismPortal_AdminOwnerMock {
    address public owner;

    constructor(address _owner) {
        owner = _owner;
    }
}

contract ShutdownOptimismPortal_AdminNoOwnerMock {}

contract ShutdownOptimismPortal_TokenMock is ERC20 {
    constructor() ERC20("Native Token", "NATIVE") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract ShutdownOptimismPortal_SystemConfigMock {
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

contract ShutdownOptimismPortal_Test is Test {
    using stdStorage for StdStorage;

    bytes32 internal constant ADMIN_SLOT =
        0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    ShutdownOptimismPortal internal portal;
    ShutdownOptimismPortal_TokenMock internal token;

    address internal admin;
    address internal recipient;

    function setUp() public {
        portal = new ShutdownOptimismPortal();
        token = new ShutdownOptimismPortal_TokenMock();
        admin = makeAddr("admin");
        recipient = makeAddr("recipient");
    }

    function test_sweepNativeToken_revertsWhenAdminNotSet() public {
        vm.expectRevert("ShutdownOptimismPortal: admin not set");
        portal.sweepNativeToken(recipient, 1);
    }

    function test_sweepNativeToken_revertsWhenNativeTokenIsETH() public {
        _setProxyAdmin(address(portal), admin);
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal_SystemConfigMock(address(0))));

        vm.prank(admin);
        vm.expectRevert("ShutdownOptimismPortal: native token is ETH");
        portal.sweepNativeToken(recipient, 1);
    }

    function test_sweepNativeToken_succeedsForAdminEOA() public {
        _setProxyAdmin(address(portal), admin);
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal_SystemConfigMock(address(token))));
        token.mint(address(portal), 100);

        vm.prank(admin);
        portal.sweepNativeToken(recipient, 60);

        assertEq(token.balanceOf(recipient), 60);
        assertEq(token.balanceOf(address(portal)), 40);
    }

    function test_sweepNativeToken_succeedsForAdminOwnerContract() public {
        ShutdownOptimismPortal_AdminOwnerMock adminContract = new ShutdownOptimismPortal_AdminOwnerMock(admin);
        _setProxyAdmin(address(portal), address(adminContract));
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal_SystemConfigMock(address(token))));
        token.mint(address(portal), 75);

        vm.prank(admin);
        portal.sweepNativeToken(recipient, 75);

        assertEq(token.balanceOf(recipient), 75);
        assertEq(token.balanceOf(address(portal)), 0);
    }

    function test_sweepNativeToken_revertsForNonOwner() public {
        ShutdownOptimismPortal_AdminOwnerMock adminContract = new ShutdownOptimismPortal_AdminOwnerMock(admin);
        _setProxyAdmin(address(portal), address(adminContract));
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal_SystemConfigMock(address(token))));

        vm.prank(makeAddr("notOwner"));
        vm.expectRevert("ShutdownOptimismPortal: unauthorized");
        portal.sweepNativeToken(recipient, 1);
    }

    function test_sweepNativeToken_revertsWhenAdminOwnerLookupFails() public {
        ShutdownOptimismPortal_AdminNoOwnerMock adminContract = new ShutdownOptimismPortal_AdminNoOwnerMock();
        _setProxyAdmin(address(portal), address(adminContract));
        _setSystemConfig(address(portal), address(new ShutdownOptimismPortal_SystemConfigMock(address(token))));

        vm.prank(admin);
        vm.expectRevert("ShutdownOptimismPortal: admin owner lookup failed");
        portal.sweepNativeToken(recipient, 1);
    }

    function _setProxyAdmin(address target, address admin_) internal {
        vm.store(target, ADMIN_SLOT, bytes32(uint256(uint160(admin_))));
    }

    function _setSystemConfig(address target, address config_) internal {
        stdstore.target(target).sig("systemConfig()").checked_write(config_);
    }
}
