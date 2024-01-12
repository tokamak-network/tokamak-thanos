// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Predeploys } from "src/libraries/Predeploys.sol";
import { OptimismMintableERC20 } from "src/universal/OptimismMintableERC20.sol";

/// @custom:proxied
/// @custom:predeploy 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000
/// @title WETH
/// @notice WETH is a contract that Wrap ETH
contract WETH is OptimismMintableERC20 {
    /// @notice Initializes the contract as an Optimism Mintable ERC20.
    constructor() OptimismMintableERC20(Predeploys.L2_STANDARD_BRIDGE, address(0), "Wrapped Ether", "WETH", 18) { }

}
