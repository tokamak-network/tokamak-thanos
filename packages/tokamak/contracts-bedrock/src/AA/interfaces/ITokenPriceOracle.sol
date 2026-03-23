// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

/// @title ITokenPriceOracle
/// @notice TON/ETH exchange rate oracle interface for MultiTokenPaymaster
interface ITokenPriceOracle {
    /// @notice Returns 1 TON price in ETH (18 decimals)
    /// @return price e.g. 0.0005e18 means 1 TON = 0.0005 ETH
    function getPrice() external view returns (uint256 price);

    /// @notice Returns timestamp of last price update
    function lastUpdated() external view returns (uint256);
}
