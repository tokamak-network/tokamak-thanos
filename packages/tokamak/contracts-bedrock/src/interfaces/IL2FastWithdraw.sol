// SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.15;

interface IL2FastWithdraw {
    function claimFW(
        address _from,
        address _to,
        uint256 _amount,
        uint256 _saleCount
    )
        external
        payable;
}
