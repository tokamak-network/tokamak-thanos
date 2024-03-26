// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import { CommonTest } from "test/CommonTest.t.sol";

// Target contract dependencies
import { Predeploys } from "src/libraries/Predeploys.sol";

// Target contract
import { LegacyERC20NativeToken } from "src/legacy/LegacyERC20NativeToken.sol";

contract LegacyERC20NativeToken_Test is CommonTest {
    LegacyERC20NativeToken eth;

    /// @dev Sets up the test suite.
    function setUp() public virtual override {
        super.setUp();
        eth = new LegacyERC20NativeToken();
    }

    /// @dev Tests that the default metadata was set correctly.
    function test_metadata_succeeds() external {
        assertEq(eth.name(), "Ether");
        assertEq(eth.symbol(), "ETH");
        assertEq(eth.decimals(), 18);
    }

    /// @dev Tests that `l2Bridge` and `l1Token` return the correct values.
    function test_crossDomain_succeeds() external {
        assertEq(eth.l2Bridge(), Predeploys.L2_STANDARD_BRIDGE);
        assertEq(eth.l1Token(), address(0));
    }

    /// @dev Tests that `transfer` reverts since it does not exist.
    function test_transfer_doesNotExist_reverts() external {
        vm.expectRevert("LegacyERC20NativeToken: transfer is disabled");
        eth.transfer(alice, 100);
    }

    /// @dev Tests that `approve` reverts since it does not exist.
    function test_approve_doesNotExist_reverts() external {
        vm.expectRevert("LegacyERC20NativeToken: approve is disabled");
        eth.approve(alice, 100);
    }

    /// @dev Tests that `transferFrom` reverts since it does not exist.
    function test_transferFrom_doesNotExist_reverts() external {
        vm.expectRevert("LegacyERC20NativeToken: transferFrom is disabled");
        eth.transferFrom(bob, alice, 100);
    }

    /// @dev Tests that `increaseAllowance` reverts since it does not exist.
    function test_increaseAllowance_doesNotExist_reverts() external {
        vm.expectRevert("LegacyERC20NativeToken: increaseAllowance is disabled");
        eth.increaseAllowance(alice, 100);
    }

    /// @dev Tests that `decreaseAllowance` reverts since it does not exist.
    function test_decreaseAllowance_doesNotExist_reverts() external {
        vm.expectRevert("LegacyERC20NativeToken: decreaseAllowance is disabled");
        eth.decreaseAllowance(alice, 100);
    }

    /// @dev Tests that `mint` reverts since it does not exist.
    function test_mint_doesNotExist_reverts() external {
        vm.expectRevert("LegacyERC20NativeToken: mint is disabled");
        eth.mint(alice, 100);
    }

    /// @dev Tests that `burn` reverts since it does not exist.
    function test_burn_doesNotExist_reverts() external {
        vm.expectRevert("LegacyERC20NativeToken: burn is disabled");
        eth.burn(alice, 100);
    }
}
