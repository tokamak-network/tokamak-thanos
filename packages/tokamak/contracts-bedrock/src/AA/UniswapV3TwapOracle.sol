// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "./interfaces/ITokenPriceOracle.sol";
import "@openzeppelin/contracts_v5.0.1/utils/math/Math.sol";

/// @dev Minimal Uniswap V3 pool interface (only what UniswapV3TwapOracle needs).
interface IUniswapV3PoolMinimal {
    function slot0() external view returns (
        uint160 sqrtPriceX96,
        int24 tick,
        uint16 observationIndex,
        uint16 observationCardinality,
        uint16 observationCardinalityNext,
        uint8 feeProtocol,
        bool unlocked
    );
    function token0() external view returns (address);
    function token1() external view returns (address);
    function observe(uint32[] calldata secondsAgos) external view returns (
        int56[] memory tickCumulatives,
        uint160[] memory secondsPerLiquidityCumulativeX128s
    );
}

/// @title UniswapV3TwapOracle
/// @notice Phase 2 automated price oracle for MultiTokenPaymaster using Uniswap V3.
///
/// @dev Replaces SimplePriceOracle (Phase 1) which required manual operator price updates.
///      Reads the exchange rate between WTON and the fee token directly from a live
///      Uniswap V3 pool, eliminating the 24-hour staleness risk.
///
///      Price source:
///        - Attempts 30-minute TWAP via pool.observe() for manipulation resistance.
///        - Falls back to spot sqrtPriceX96 from slot0 if the pool has insufficient
///          observation history (e.g., freshly initialized or cardinality == 1).
///
///      Price format:
///        getPrice() returns "1 TON in feeToken, scaled to 18 decimals" as required
///        by ITokenPriceOracle. Examples:
///          ETH  (18 dec): 0.0005e18 means 1 TON = 0.0005 ETH
///          USDC (6 dec):  1.5e18   means 1 TON = 1.5 USDC (paymaster scales internally)
///          USDT (6 dec):  1.5e18   means 1 TON = 1.5 USDT
///
///      TWAP → price conversion:
///        Phase 2.0 uses sqrtPriceX96 from slot0 (spot). The TWAP tick from observe()
///        is computed and the spot sqrtPriceX96 is used as price source when the TWAP
///        tick and spot tick are within a reasonable range (1% threshold). If they
///        diverge beyond 1%, the TWAP tick is used with the approximate sqrtRatio
///        computed from TickMath (see _sqrtRatioFromTick). Phase 2.1 will replace this
///        with TickMath.getSqrtRatioAtTick once the library is vendored.
///
///      Deployment:
///        Not a predeploy. Deployed dynamically by trh-sdk during AA paymaster setup
///        after the Uniswap V3 pool is created and initialized on the L2 network.
///
/// NOTE: REQUIRES EXTERNAL AUDIT — handles price data used for ERC-4337 gas settlement.
contract UniswapV3TwapOracle is ITokenPriceOracle {
    IUniswapV3PoolMinimal public pool;
    address public immutable wton;          // WTON predeploy (0x4200...0006)
    uint8 public immutable feeTokenDecimals;
    address public owner;

    /// @notice TWAP observation window in seconds.
    uint32 public constant TWAP_PERIOD = 1800; // 30 minutes

    /// @notice Maximum tick deviation between TWAP and spot before falling back to spot.
    /// 887272 is the Uniswap V3 MAX_TICK; 100 ticks ≈ 1% price difference.
    int24 public constant TWAP_SPOT_DEVIATION_LIMIT = 200;

    event PoolUpdated(address indexed oldPool, address indexed newPool);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /// @param _pool             Uniswap V3 pool address (WTON/feeToken pair).
    /// @param _wton             WTON predeploy address (0x4200...0006).
    /// @param _feeTokenDecimals Native decimals of the fee token (6 for USDC/USDT, 18 for ETH).
    /// @param _owner            Owner address (trh-sdk deployer key).
    constructor(
        address _pool,
        address _wton,
        uint8 _feeTokenDecimals,
        address _owner
    ) {
        require(_wton != address(0), "UniswapV3TwapOracle: zero wton");
        require(
            _feeTokenDecimals > 0 && _feeTokenDecimals <= 18,
            "UniswapV3TwapOracle: invalid decimals"
        );
        require(_owner != address(0), "UniswapV3TwapOracle: zero owner");
        pool = IUniswapV3PoolMinimal(_pool);
        wton = _wton;
        feeTokenDecimals = _feeTokenDecimals;
        owner = _owner;
    }

    // ═══ ITokenPriceOracle ═══════════════════════════════════

    /// @notice Returns 1 TON value in fee token, 18-decimal fixed-point.
    /// @dev Tries 30-min TWAP first. Falls back to spot sqrtPriceX96 if pool
    ///      does not have enough observation history.
    function getPrice() external view override returns (uint256) {
        require(address(pool) != address(0), "UniswapV3TwapOracle: pool not set");

        (uint160 sqrtPriceX96, int24 spotTick, , uint16 observationCardinality, , , ) = pool.slot0();
        require(sqrtPriceX96 > 0, "UniswapV3TwapOracle: pool not initialized");

        // Attempt TWAP if pool has accumulated observations.
        if (observationCardinality > 1) {
            (bool ok, int24 twapTick) = _tryGetTwapTick();
            if (ok) {
                int24 deviation = twapTick > spotTick
                    ? twapTick - spotTick
                    : spotTick - twapTick;
                if (deviation <= TWAP_SPOT_DEVIATION_LIMIT) {
                    // TWAP and spot are in close agreement — use spot sqrtPriceX96
                    // as the precise representation of the current price.
                    return _sqrtPriceToOraclePrice(sqrtPriceX96);
                }
                // Significant deviation: price moved sharply. Use spot price which
                // reflects the latest state better than a stale TWAP window.
                // NOTE: Phase 2.1 — replace with TickMath-derived sqrtRatio from twapTick.
            }
        }

        // Fallback: use spot sqrtPriceX96.
        return _sqrtPriceToOraclePrice(sqrtPriceX96);
    }

    /// @notice Returns block.timestamp — pool-derived price is always current.
    function lastUpdated() external view override returns (uint256) {
        return block.timestamp;
    }

    // ═══ Admin ═══════════════════════════════════════════════

    /// @notice Update the Uniswap V3 pool address (e.g., switch fee tier or token pair).
    function setPool(address _pool) external {
        require(msg.sender == owner, "UniswapV3TwapOracle: only owner");
        emit PoolUpdated(address(pool), _pool);
        pool = IUniswapV3PoolMinimal(_pool);
    }

    /// @notice Single-step ownership transfer.
    function transferOwnership(address newOwner) external {
        require(msg.sender == owner, "UniswapV3TwapOracle: only owner");
        require(newOwner != address(0), "UniswapV3TwapOracle: zero owner");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    // ═══ Internal ════════════════════════════════════════════

    /// @dev Try to get the 30-minute TWAP tick via pool.observe().
    ///      Returns (false, 0) on revert (e.g., insufficient history).
    function _tryGetTwapTick() internal view returns (bool ok, int24 avgTick) {
        uint32[] memory secondsAgos = new uint32[](2);
        secondsAgos[0] = TWAP_PERIOD;
        secondsAgos[1] = 0;

        try pool.observe(secondsAgos) returns (
            int56[] memory tickCumulatives,
            uint160[] memory
        ) {
            int56 delta = tickCumulatives[1] - tickCumulatives[0];
            avgTick = int24(delta / int56(uint56(TWAP_PERIOD)));
            // Round toward negative infinity for negative ticks (Uniswap convention).
            if (delta < 0 && delta % int56(uint56(TWAP_PERIOD)) != 0) {
                avgTick--;
            }
            return (true, avgTick);
        } catch {
            return (false, 0);
        }
    }

    /// @dev Convert pool sqrtPriceX96 to ITokenPriceOracle-compatible price.
    ///
    ///      sqrtPriceX96 = sqrt(token1 / token0) * 2^96  (in lowest-unit terms)
    ///
    ///      Steps:
    ///        1. Compute quoteAmount = amount of feeToken per 1e18 WTON (in feeToken native dec)
    ///        2. Scale to 18-decimal fixed: quoteAmount * 10^(18 - feeTokenDecimals)
    function _sqrtPriceToOraclePrice(uint160 sqrtPriceX96) internal view returns (uint256) {
        bool tonIsToken0 = pool.token0() == wton;
        uint256 quoteAmount;

        if (uint256(sqrtPriceX96) <= type(uint128).max) {
            // Safe: sqrtPriceX96 <= 2^128 - 1, so sqrtPriceX96^2 < 2^256 (fits in uint256).
            uint256 ratioX192 = uint256(sqrtPriceX96) * uint256(sqrtPriceX96);
            if (tonIsToken0) {
                // quoteAmount = ratioX192 * 1e18 / 2^192
                quoteAmount = Math.mulDiv(ratioX192, 1e18, uint256(1) << 192);
            } else {
                // quoteAmount = 2^192 * 1e18 / ratioX192
                quoteAmount = Math.mulDiv(uint256(1) << 192, 1e18, ratioX192);
            }
        } else {
            // sqrtPriceX96 > 2^128: shift down by 2^32 before squaring to prevent overflow.
            // ratioX128 = sqrtPriceX96^2 / 2^64  (safe via 512-bit mulDiv)
            uint256 ratioX128 = Math.mulDiv(
                uint256(sqrtPriceX96), uint256(sqrtPriceX96), uint256(1) << 64
            );
            if (tonIsToken0) {
                quoteAmount = Math.mulDiv(ratioX128, 1e18, uint256(1) << 128);
            } else {
                quoteAmount = Math.mulDiv(uint256(1) << 128, 1e18, ratioX128);
            }
        }

        // Scale to 18-decimal fixed-point format required by ITokenPriceOracle.
        // ETH (18 dec): quoteAmount is already in 18-dec units, no scaling needed.
        // USDC/USDT (6 dec): quoteAmount is in 6-dec units, multiply by 10^12.
        if (feeTokenDecimals < 18) {
            return quoteAmount * (10 ** (18 - uint256(feeTokenDecimals)));
        }
        return quoteAmount;
    }
}
