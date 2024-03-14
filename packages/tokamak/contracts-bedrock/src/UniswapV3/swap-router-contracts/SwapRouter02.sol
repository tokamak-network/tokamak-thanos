// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity =0.7.6;
pragma abicoder v2;

import '@v3-periphery/contracts/base/SelfPermit.sol';
import '@v3-periphery/contracts/base/PeripheryImmutableState.sol';

import '@swap-router-contracts/contracts/interfaces/ISwapRouter02.sol';
import '@swap-router-contracts/contracts/V2SwapRouter.sol';
import '@swap-router-contracts/contracts/V3SwapRouter.sol';
import '@swap-router-contracts/contracts/base/ApproveAndCall.sol';
import '@swap-router-contracts/contracts/base/MulticallExtended.sol';

/// @title Uniswap V2 and V3 Swap Router
contract SwapRouter02 is ISwapRouter02, V2SwapRouter, V3SwapRouter, ApproveAndCall, MulticallExtended, SelfPermit {
    constructor(
        address _factoryV2,
        address factoryV3,
        address _positionManager,
        address _WETH9
    ) ImmutableState(_factoryV2, _positionManager) PeripheryImmutableState(factoryV3, _WETH9) {}
}
