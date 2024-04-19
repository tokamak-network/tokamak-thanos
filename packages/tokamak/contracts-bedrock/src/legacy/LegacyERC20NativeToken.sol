// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Predeploys } from "src/libraries/Predeploys.sol";
import { OptimismMintableERC20 } from "src/universal/OptimismMintableERC20.sol";

/// @custom:legacy
/// @custom:proxied
/// @custom:predeploy 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000
/// @title LegacyERC20NativeToken
/// @notice LegacyERC20NativeToken is a legacy contract that held NativeToken balances before the Bedrock upgrade.
///         All NativeToken balances held within this contract were migrated to the state trie as part of
///         the Bedrock upgrade. Functions within this contract that mutate state were already
///         disabled as part of the EVM equivalence upgrade.
contract LegacyERC20NativeToken is OptimismMintableERC20 {
    /// @notice Initializes the contract as an Optimism Mintable ERC20.
    constructor() OptimismMintableERC20(Predeploys.L2_STANDARD_BRIDGE, address(0), "Ether", "ETH", 18) { }

    /// @notice Returns the NativeToken balance of the target account. Overrides the base behavior of the
    ///         contract to preserve the invariant that the balance within this contract always
    ///         matches the balance in the state trie.
    /// @param _who Address of the account to query.
    /// @return The NativeToken balance of the target account.
    function balanceOf(address _who) public view virtual override returns (uint256) {
        return address(_who).balance;
    }

    /// @custom:blocked
    /// @notice Mints some amount of NativeToken.
    function mint(address, uint256) public virtual override {
        revert("LegacyERC20NativeToken: mint is disabled");
    }

    /// @custom:blocked
    /// @notice Burns some amount of NativeToken.
    function burn(address, uint256) public virtual override {
        revert("LegacyERC20NativeToken: burn is disabled");
    }

    /// @custom:blocked
    /// @notice Transfers some amount of NativeToken.
    function transfer(address, uint256) public virtual override returns (bool) {
        revert("LegacyERC20NativeToken: transfer is disabled");
    }

    /// @custom:blocked
    /// @notice Approves a spender to spend some amount of NativeToken.
    function approve(address, uint256) public virtual override returns (bool) {
        revert("LegacyERC20NativeToken: approve is disabled");
    }

    /// @custom:blocked
    /// @notice Transfers funds from some sender account.
    function transferFrom(address, address, uint256) public virtual override returns (bool) {
        revert("LegacyERC20NativeToken: transferFrom is disabled");
    }

    /// @custom:blocked
    /// @notice Increases the allowance of a spender.
    function increaseAllowance(address, uint256) public virtual override returns (bool) {
        revert("LegacyERC20NativeToken: increaseAllowance is disabled");
    }

    /// @custom:blocked
    /// @notice Decreases the allowance of a spender.
    function decreaseAllowance(address, uint256) public virtual override returns (bool) {
        revert("LegacyERC20NativeToken: decreaseAllowance is disabled");
    }
}
