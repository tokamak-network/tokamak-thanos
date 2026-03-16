// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "./lib/BasePaymaster.sol";
import "./interfaces/ITokenPriceOracle.sol";

/// @title MultiTokenPaymaster
/// @notice Users pay gas fees in ERC-20 tokens (ETH, USDC, USDT, etc.).
///         The paymaster fronts the gas in TON (native) from its EntryPoint deposit
///         and collects the equivalent ERC-20 amount from the user.
/// @dev Post-deployment module (NOT a genesis predeploy). Deployed in trh-sdk deploy Phase 6.
///      REQUIRES EXTERNAL AUDIT — multi-token, decimal conversion, price oracle.
contract MultiTokenPaymaster is BasePaymaster {
    using SafeERC20 for IERC20;
    using UserOperationLib for PackedUserOperation;

    struct TokenConfig {
        bool enabled;
        ITokenPriceOracle oracle;
        uint256 markup;    // e.g. 5 = 5%. Max 50.
        uint8 decimals;    // 18 for ETH, 6 for USDC/USDT
    }

    mapping(address => TokenConfig) public supportedTokens;
    address[] public tokenList;
    mapping(address => uint256) public collectedFees;

    /// @notice Per-token minimum charge (rounding protection for 6-decimal tokens)
    /// @dev Default 0: no minimum. Set per token after deployment if needed.
    mapping(address => uint256) public minCharge;

    event TokenAdded(address indexed token, address oracle, uint256 markup, uint8 decimals);
    event TokenRemoved(address indexed token);
    event TokenConfigUpdated(address indexed token, address oracle, uint256 markup);
    event FeesCollected(address indexed token, address indexed sender, uint256 amount);
    event FeesWithdrawn(address indexed token, address indexed to, uint256 amount);

    constructor(IEntryPoint _entryPoint) BasePaymaster() {
        // BasePaymaster() is empty constructor (Thanos modification: supports initialize() pattern)
        _setEntryPoint(_entryPoint);
    }

    // ═══ Admin ═══════════════════════════════════════════

    function addToken(
        address token,
        ITokenPriceOracle oracle,
        uint256 markupPercent,
        uint8 decimals
    ) external onlyOwner {
        require(!supportedTokens[token].enabled, "already enabled");
        require(address(oracle) != address(0), "zero oracle");
        require(markupPercent <= 50, "markup too high");

        supportedTokens[token] = TokenConfig({
            enabled: true,
            oracle: oracle,
            markup: markupPercent,
            decimals: decimals
        });
        tokenList.push(token);
        emit TokenAdded(token, address(oracle), markupPercent, decimals);
    }

    function removeToken(address token) external onlyOwner {
        require(supportedTokens[token].enabled, "not enabled");
        supportedTokens[token].enabled = false;
        emit TokenRemoved(token);
    }

    function updateTokenConfig(
        address token,
        ITokenPriceOracle oracle,
        uint256 markupPercent
    ) external onlyOwner {
        require(supportedTokens[token].enabled, "not enabled");
        require(markupPercent <= 50, "markup too high");
        supportedTokens[token].oracle = oracle;
        supportedTokens[token].markup = markupPercent;
        emit TokenConfigUpdated(token, address(oracle), markupPercent);
    }

    function withdrawCollectedFees(address token, address to, uint256 amount) external onlyOwner {
        require(collectedFees[token] >= amount, "insufficient collected");
        collectedFees[token] -= amount;
        IERC20(token).safeTransfer(to, amount);
        emit FeesWithdrawn(token, to, amount);
    }

    function setMinCharge(address token, uint256 amount) external onlyOwner {
        minCharge[token] = amount;
    }

    // ═══ View helpers (frontend / SDK) ═══════════════════

    function estimateTokenCost(address token, uint256 estimatedTonGasCost)
        external view returns (uint256 tokenCost, uint256 tokenCostWithMarkup)
    {
        TokenConfig memory cfg = supportedTokens[token];
        require(cfg.enabled, "token not supported");
        tokenCost = _tonToTokenSafe(estimatedTonGasCost, token, cfg);
        tokenCostWithMarkup = tokenCost * (100 + cfg.markup) / 100;
    }

    /// @notice Returns full TokenConfig struct for a token (for Solidity callers / tests)
    function getTokenConfig(address token) external view returns (TokenConfig memory) {
        return supportedTokens[token];
    }

    /// @dev Exposed for test access to _tonToTokenSafe (no markup applied)
    function estimateTokenCostPublic(address token, uint256 tonAmount)
        external view returns (uint256)
    {
        TokenConfig memory cfg = supportedTokens[token];
        return _tonToTokenSafe(tonAmount, token, cfg);
    }

    // ═══ Paymaster core ═══════════════════════════════════

    function _validatePaymasterUserOp(
        PackedUserOperation calldata userOp,
        bytes32, /*userOpHash*/
        uint256 maxCost
    ) internal override returns (bytes memory context, uint256 validationData) {
        // paymasterAndData Phase 1 format: [paymaster(20)][token(20)] = 40 bytes total (no signature)
        // Phase 2+: will include validUntil/validAfter/sig (see docs/TRH_MultiToken_Fee_Design.md Appendix A)
        // validationData = 0: no signature verification in Phase 1
        address token = address(bytes20(userOp.paymasterAndData[20:40]));

        TokenConfig memory cfg = supportedTokens[token];
        require(cfg.enabled, "PM: token not supported");

        uint256 tokenCost      = _tonToTokenSafe(maxCost, token, cfg);
        uint256 costWithMarkup = tokenCost * (100 + cfg.markup) / 100;

        address sender = userOp.sender;
        require(
            IERC20(token).allowance(sender, address(this)) >= costWithMarkup,
            "PM: insufficient allowance"
        );

        IERC20(token).safeTransferFrom(sender, address(this), costWithMarkup);

        context = abi.encode(sender, token, costWithMarkup, maxCost);
        validationData = 0;
    }

    function _postOp(
        PostOpMode, /*mode — refund regardless of opReverted per spec Section 9*/
        bytes calldata context,
        uint256 actualGasCost,
        uint256 /*actualUserOpFeePerGas*/
    ) internal override {
        (
            address sender,
            address token,
            uint256 preCharged,
            /*uint256 maxCost*/
        ) = abi.decode(context, (address, address, uint256, uint256));

        TokenConfig memory cfg = supportedTokens[token];

        uint256 actualCost       = _tonToTokenSafe(actualGasCost, token, cfg);
        uint256 actualWithMarkup = actualCost * (100 + cfg.markup) / 100;

        // ERC-4337: actualGasCost is always ≤ maxCost (EntryPoint guarantee).
        // Therefore actualWithMarkup ≤ preCharged always. Undercharge case cannot occur in Phase 1.
        // The `preCharged > actualWithMarkup` check is defensive (safe for future extensions).
        if (preCharged > actualWithMarkup) {
            IERC20(token).safeTransfer(sender, preCharged - actualWithMarkup);
        }

        uint256 actualCollected = preCharged > actualWithMarkup ? actualWithMarkup : preCharged;
        collectedFees[token] += actualCollected;

        emit FeesCollected(token, sender, actualCollected);
    }

    // ═══ Conversion ═══════════════════════════════════════

    /// @notice TON → token (18 or 6 decimals). See Section 3.2 of design doc.
    function _tonToToken(uint256 tonAmount, TokenConfig memory cfg)
        internal view returns (uint256)
    {
        // oracle.getPrice(): 1 TON value expressed in target token (18 decimals fixed)
        // e.g. ETH oracle: 0.0005e18 → 1 TON = 0.0005 ETH
        //      USDC oracle: 0.65e18  → 1 TON = 0.65 USDC (before 18→6 scaling)
        uint256 price = cfg.oracle.getPrice();

        if (cfg.decimals == 18) {
            return tonAmount * price / 1e18;
        } else {
            // e.g. USDC (6 dec): tonAmount(18) * price(18) / 1e18 → 18dec result
            // 18dec → 6dec: / 10^(18-6) = / 1e12
            return tonAmount * price / 1e18 / (10 ** (18 - cfg.decimals));
        }
    }

    /// @notice _tonToToken with minimum charge guard (rounding protection, Section 9.2)
    /// @dev minCharge[token] is a zero-result fallback, NOT a general floor.
    ///      Intent: 6-decimal token amounts can round down to 0 for tiny gas costs.
    ///      Only replaces with min when result == 0 (pure zero-rounding protection).
    ///      If result > 0 but result < min, the actual result is used (by design).
    ///      Use case: setMinCharge(usdc, 1) to prevent 0-USDC charges.
    function _tonToTokenSafe(uint256 tonAmount, address token, TokenConfig memory cfg)
        internal view returns (uint256)
    {
        uint256 result = _tonToToken(tonAmount, cfg);
        uint256 min = minCharge[token];
        // Return min only if result rounded to exactly 0 and min is set (> 0)
        return (result > 0 || min == 0) ? result : min;
    }
}
