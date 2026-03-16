// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "./interfaces/ITokenPriceOracle.sol";

/// @title SimplePriceOracle
/// @notice Phase 1 manual price oracle. Owner updates price; reverts if stale (>24h).
/// @dev Phase 2: replace with UniswapV3TwapOracle.
/// NOTE: REQUIRES EXTERNAL AUDIT — handles price data used for gas settlement.
contract SimplePriceOracle is ITokenPriceOracle {
    uint256 public override lastUpdated;
    uint256 private _price;
    address public owner;

    uint256 private constant STALE_THRESHOLD = 86400; // strictly < 86400s valid (exactly 24h = stale)

    event PriceUpdated(uint256 newPrice, uint256 timestamp);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    constructor(uint256 initialPrice) {
        require(initialPrice > 0, "SimplePriceOracle: zero price");
        owner = msg.sender;
        _price = initialPrice;
        lastUpdated = block.timestamp;
    }

    function getPrice() external view override returns (uint256) {
        require(block.timestamp - lastUpdated < STALE_THRESHOLD, "SimplePriceOracle: stale price");
        return _price;
    }

    function updatePrice(uint256 newPrice) external {
        require(msg.sender == owner, "only owner");
        require(newPrice > 0, "SimplePriceOracle: zero price");
        _price = newPrice;
        lastUpdated = block.timestamp;
        emit PriceUpdated(newPrice, block.timestamp);
    }

    // NOTE: single-step ownership transfer (no pendingOwner two-step).
    // Known limitation: address typo causes permanent loss of control.
    // Phase 2 improvement: adopt two-step pattern.
    // External audit required before production deployment.
    function transferOwnership(address newOwner) external {
        require(msg.sender == owner, "only owner");
        require(newOwner != address(0), "SimplePriceOracle: zero owner");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }
}
