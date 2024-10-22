// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title Burn
/// @notice Utilities for burning stuff.
library Burn {
    /// @notice Burns a given amount of the native token.
    /// @param _amount Amount of the native token to burn.
    function nativeToken(uint256 _amount) internal {
        new Burner{ value: _amount }();
    }

    /// @notice Burns a given amount of gas.
    /// @param _amount Amount of gas to burn.
    function gas(uint256 _amount) internal view {
        uint256 i = 0;
        uint256 initialGas = gasleft();
        while (initialGas - gasleft() < _amount) {
            ++i;
        }
    }
}

/// @title Burner
/// @notice Burner self-destructs on creation and sends all native tokens to itself, removing all native tokens given to
///         the contract from the circulating supply. Self-destructing is the only way to remove native tokens
///         from the circulating supply.
contract Burner {
    constructor() payable {
        selfdestruct(payable(address(this)));
    }
}
