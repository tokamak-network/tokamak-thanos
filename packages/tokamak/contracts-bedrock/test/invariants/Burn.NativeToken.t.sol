// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { StdUtils } from "forge-std/StdUtils.sol";
import { Test } from "forge-std/Test.sol";
import { Vm } from "forge-std/Vm.sol";

import { StdInvariant } from "forge-std/StdInvariant.sol";
import { Burn } from "src/libraries/Burn.sol";
import { InvariantTest } from "test/invariants/InvariantTest.sol";

contract Burn_NativeTokenBurner is StdUtils {
    Vm internal vm;
    bool public failedNativeTokenBurn;

    constructor(Vm _vm) {
        vm = _vm;
    }

    /// @notice Takes an integer amount of native tokens to burn through the Burn library and
    ///         updates the contract state if an incorrect amount of native tokens moved from the contract
    function burnNativeToken(uint256 _value) external {
        uint256 preBurnvalue = bound(_value, 0, type(uint128).max);

        // Give the burner some native tokens for gas being used
        vm.deal(address(this), preBurnvalue);

        // cache the contract's native tokens balance
        uint256 preBurnBalance = address(this).balance;

        uint256 value = bound(preBurnvalue, 0, preBurnBalance);

        // execute a burn of _value native tokens
        Burn.nativeToken(value);

        // check that exactly value native tokens was transferred from the contract
        unchecked {
            if (address(this).balance != preBurnBalance - value) {
                failedNativeTokenBurn = true;
            }
        }
    }
}

contract Burn_BurnNativeToken_Invariant is StdInvariant, InvariantTest {
    Burn_NativeTokenBurner internal actor;

    function setUp() public override {
        super.setUp();
        // Create a native token burner actor.

        actor = new Burn_NativeTokenBurner(vm);

        targetContract(address(actor));

        bytes4[] memory selectors = new bytes4[](1);
        selectors[0] = actor.burnNativeToken.selector;
        FuzzSelector memory selector = FuzzSelector({ addr: address(actor), selectors: selectors });
        targetSelector(selector);
    }

    /// @custom:invariant `nativeToken(uint256)` always burns the exact amount of native token passed.
    ///
    ///                   Asserts that when `Burn.nativeToken(uint256)` is called, it always
    ///                   burns the exact amount of native tokens passed to the function.
    function invariant_burn_native_token() external view {
        // ASSERTION: The amount burned should always match the amount passed exactly
        assertEq(actor.failedNativeTokenBurn(), false);
    }
}
