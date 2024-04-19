// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Predeploys } from "src/libraries/Predeploys.sol";
import { OptimismMintableERC20 } from "src/universal/OptimismMintableERC20.sol";

/// @custom:proxied
/// @custom:predeploy 0x4200000000000000000000000000000000000486
/// @title ETH
/// @notice ETH is a contract that Wrap ETH
contract ETH is OptimismMintableERC20 {
    /// @notice Initializes the contract as an Optimism Mintable ERC20.
    constructor() OptimismMintableERC20(Predeploys.L2_STANDARD_BRIDGE, address(0), "Ether", "ETH", 18) { }
}
