// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "./interfaces/ITokenPriceOracle.sol";
import "@openzeppelin/contracts_v5.0.1/proxy/utils/Initializable.sol";

/// @title SimplePriceOracle
/// @notice Phase 1 manual price oracle. Owner updates price; reverts if stale (>24h).
/// @dev Phase 2: replace with UniswapV3TwapOracle.
///      Uses Initializable for genesis predeploy pattern (0x4200...0066).
///      At genesis, owner is set via storage; operator must call updatePrice() before use.
/// NOTE: REQUIRES EXTERNAL AUDIT — handles price data used for gas settlement.
contract SimplePriceOracle is ITokenPriceOracle, Initializable {
    uint256 public override lastUpdated;
    uint256 private _price;
    address public owner;

    uint256 private constant STALE_THRESHOLD = 86400; // strictly < 86400s valid (exactly 24h = stale)

    event PriceUpdated(uint256 newPrice, uint256 timestamp);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    // Empty constructor for genesis predeploy pattern (0x4200...0066).
    // State is pre-set via genesis storage; operator calls updatePrice() before use.
    constructor() {}

    /**
     * Initialize the oracle. For non-genesis deployments only.
     * Genesis predeploys use direct storage initialization instead.
     * @param initialPrice - Initial price (must be > 0).
     * @param _owner       - Owner address for price updates.
     */
    function initialize(uint256 initialPrice, address _owner) external initializer {
        require(initialPrice > 0, "SimplePriceOracle: zero price");
        require(_owner != address(0), "SimplePriceOracle: zero owner");
        owner = _owner;
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
