// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { SuperchainConfig_Initializer } from "test/CommonTest.t.sol";

// Libraries
import { Types } from "src/libraries/Types.sol";
import { Hashing } from "src/libraries/Hashing.sol";

// Target contract dependencies
import { Proxy } from "src/universal/Proxy.sol";

// Target contract
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";

contract SuperchainConfig_Init_Test is SuperchainConfig_Initializer {
    /// @dev Tests that initialization sets the correct values. These are defined in CommonTest.sol.
    function test_initialize_unpaused_succeeds() external {
        assertFalse(sc.paused());
    }

    /// @dev Tests that it can be intialized as paused.
    function test_initialize_paused_succeeds() external {
        Proxy newProxy = new Proxy(alice);
        SuperchainConfig newImpl = new SuperchainConfig();

        vm.startPrank(alice);
        newProxy.upgradeToAndCall(
            address(newImpl), abi.encodeWithSelector(SuperchainConfig.initialize.selector, guardian, true)
        );

        assertTrue(SuperchainConfig(address(newProxy)).paused());
        assertEq(SuperchainConfig(address(newProxy)).guardian(), guardian);
    }
}

contract SuperchainConfig_Pause_TestFail is SuperchainConfig_Initializer {
    /// @dev Tests that `pause` reverts when called by a non-guardian.
    function test_pause_notGuardian_reverts() external {
        assertFalse(sc.paused());

        assertTrue(sc.guardian() != alice);
        vm.expectRevert("SuperchainConfig: only guardian can pause");
        vm.prank(alice);
        sc.pause("identifier");

        assertFalse(sc.paused());
    }
}

contract SuperchainConfig_Pause_Test is SuperchainConfig_Initializer {
    /// @dev Tests that `pause` successfully pauses
    ///      when called by the guardian.
    function test_pause_succeeds() external {
        assertFalse(sc.paused());

        vm.expectEmit(address(sc));
        emit Paused("identifier");

        vm.prank(sc.guardian());
        sc.pause("identifier");

        assertTrue(sc.paused());
    }
}

contract SuperchainConfig_Unpause_TestFail is SuperchainConfig_Initializer {
    /// @dev Tests that `unpause` reverts when called by a non-guardian.
    function test_unpause_notGuardian_reverts() external {
        vm.prank(sc.guardian());
        sc.pause("identifier");
        assertEq(sc.paused(), true);

        assertTrue(sc.guardian() != alice);
        vm.expectRevert("SuperchainConfig: only guardian can unpause");
        vm.prank(alice);
        sc.unpause();

        assertTrue(sc.paused());
    }
}

contract SuperchainConfig_Unpause_Test is SuperchainConfig_Initializer {
    /// @dev Tests that `unpause` successfully unpauses
    ///      when called by the guardian.
    function test_unpause_succeeds() external {
        vm.startPrank(sc.guardian());
        sc.pause("identifier");
        assertEq(sc.paused(), true);

        vm.expectEmit(address(sc));
        emit Unpaused();
        sc.unpause();

        assertFalse(sc.paused());
    }
}
