// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../../../src/tokamak/AA/MultiTokenPaymaster.sol";
import "../../../src/tokamak/AA/SimplePriceOracle.sol";
import "../../../src/tokamak/AA/EntryPoint.sol";
import "../../../src/tokamak/AA/interfaces/IPaymaster.sol";

// Mock ERC-20 with 18 decimals (L2 ETH)
contract MockETH is ERC20 {
    constructor() ERC20("L2 ETH", "ETH") {}
    function mint(address to, uint256 amount) external { _mint(to, amount); }
}

// Mock ERC-20 with 6 decimals (USDC/USDT)
contract MockUSDC is ERC20 {
    constructor() ERC20("USD Coin", "USDC") {}
    function decimals() public pure override returns (uint8) { return 6; }
    function mint(address to, uint256 amount) external { _mint(to, amount); }
}

contract MultiTokenPaymasterTest is Test {
    EntryPoint entryPoint;
    MockETH l2Eth;
    MockUSDC usdc;
    SimplePriceOracle ethOracle;  // 1 TON = 0.0005 ETH (price = 0.0005e18)
    SimplePriceOracle usdcOracle; // 1 TON = 0.65 USDC (price = 0.65e18, then scale to 6 dec)
    MultiTokenPaymaster paymaster;
    address owner = address(this);
    address user   = makeAddr("user");

    function setUp() public {
        entryPoint = new EntryPoint();
        l2Eth  = new MockETH();
        usdc   = new MockUSDC();

        // ETH oracle: price = 0.0005e18 means 1 TON = 0.0005 ETH
        ethOracle  = new SimplePriceOracle(0.0005e18);
        // USDC oracle: price = 0.65e18 means 1 TON = 0.65 USDC (after 18→6 dec scaling)
        usdcOracle = new SimplePriceOracle(0.65e18);

        paymaster = new MultiTokenPaymaster(IEntryPoint(address(entryPoint)));

        // Register ETH (18 decimals, markup 5%)
        paymaster.addToken(address(l2Eth), ITokenPriceOracle(address(ethOracle)), 5, 18);
        // Register USDC (6 decimals, markup 3%)
        paymaster.addToken(address(usdc), ITokenPriceOracle(address(usdcOracle)), 3, 6);

        vm.deal(address(paymaster), 10 ether);
        paymaster.deposit{value: 5 ether}();
    }

    // ── Admin ─────────────────────────────────────────────

    function test_AddToken_SetsConfig() public view {
        (bool enabled, ITokenPriceOracle oracle, uint256 markup, uint8 dec) =
            _getTokenConfig(address(l2Eth));
        assertTrue(enabled);
        assertEq(address(oracle), address(ethOracle));
        assertEq(markup, 5);
        assertEq(dec, 18);
    }

    function test_AddToken_AlreadyEnabledReverts() public {
        vm.expectRevert("already enabled");
        paymaster.addToken(address(l2Eth), ITokenPriceOracle(address(ethOracle)), 5, 18);
    }

    function test_AddToken_MarkupTooHighReverts() public {
        vm.expectRevert("markup too high");
        paymaster.addToken(makeAddr("newToken"), ITokenPriceOracle(address(ethOracle)), 51, 18);
    }

    function test_RemoveToken_DisablesToken() public {
        paymaster.removeToken(address(l2Eth));
        (bool enabled,,,) = _getTokenConfig(address(l2Eth));
        assertFalse(enabled);
    }

    function test_RemoveToken_NotEnabledReverts() public {
        paymaster.removeToken(address(l2Eth)); // first removal
        vm.expectRevert("not enabled");
        paymaster.removeToken(address(l2Eth)); // second removal
    }

    function test_UpdateTokenConfig_ChangesMarkup() public {
        paymaster.updateTokenConfig(address(l2Eth), ITokenPriceOracle(address(ethOracle)), 10);
        (, , uint256 markup,) = _getTokenConfig(address(l2Eth));
        assertEq(markup, 10);
    }

    // ── Decimal conversion (critical T2, T9) ─────────────

    function test_TonToToken_ETH_18Decimals() public view {
        // 10 TON at 0.0005 ETH/TON = 0.005 ETH
        uint256 result = paymaster.estimateTokenCostPublic(address(l2Eth), 10 ether);
        assertEq(result, 0.005 ether, "10 TON * 0.0005 = 0.005 ETH");
    }

    function test_TonToToken_USDC_6Decimals() public view {
        // 10 TON at 0.65 USDC/TON = 6.5 USDC = 6_500_000 (6 decimals)
        uint256 result = paymaster.estimateTokenCostPublic(address(usdc), 10 ether);
        assertEq(result, 6_500_000, "10 TON * 0.65 USDC = 6.5 USDC = 6_500_000 units");
    }

    function test_TonToToken_USDC_Markup() public view {
        // T9: 10 TON, USDC/TON=0.65, markup=3%
        // 6.5 USDC * 1.03 = 6.695 USDC = 6_695_000 units
        (, uint256 withMarkup) = paymaster.estimateTokenCost(address(usdc), 10 ether);
        assertEq(withMarkup, 6_695_000, "6.5 USDC * 1.03 markup = 6.695 USDC");
    }

    // ── validatePaymasterUserOp ───────────────────────────

    function test_ValidatePaymasterUserOp_UnsupportedTokenReverts() public {
        // T4: token not in supportedTokens
        PackedUserOperation memory userOp = _buildUserOp(user, makeAddr("unsupportedToken"));
        vm.prank(address(entryPoint));
        vm.expectRevert("PM: token not supported");
        paymaster.validatePaymasterUserOp(userOp, bytes32(0), 1 ether);
    }

    function test_ValidatePaymasterUserOp_InsufficientAllowanceReverts() public {
        PackedUserOperation memory userOp = _buildUserOp(user, address(l2Eth));
        vm.prank(address(entryPoint));
        vm.expectRevert("PM: insufficient allowance");
        paymaster.validatePaymasterUserOp(userOp, bytes32(0), 1 ether);
    }

    function test_ValidatePaymasterUserOp_PreChargesETH() public {
        // T1: ETH (18 dec) pre-charge
        uint256 maxCost = 1 ether; // 1 TON
        // 1 TON * 0.0005 = 0.0005 ETH, markup 5% → 0.000525 ETH
        uint256 expectedCharge = 0.0005 ether * 105 / 100;

        l2Eth.mint(user, expectedCharge);
        vm.prank(user);
        l2Eth.approve(address(paymaster), expectedCharge);

        PackedUserOperation memory userOp = _buildUserOp(user, address(l2Eth));
        vm.prank(address(entryPoint));
        (bytes memory context,) = paymaster.validatePaymasterUserOp(userOp, bytes32(0), maxCost);

        assertEq(l2Eth.balanceOf(user), 0, "User's ETH fully pre-charged");
        assertEq(l2Eth.balanceOf(address(paymaster)), expectedCharge);

        (address ctxSender, address ctxToken, uint256 ctxCharged,) =
            abi.decode(context, (address, address, uint256, uint256));
        assertEq(ctxSender, user);
        assertEq(ctxToken, address(l2Eth));
        assertEq(ctxCharged, expectedCharge);
    }

    function test_ValidatePaymasterUserOp_PreChargesUSDC() public {
        // T2: USDC (6 dec) pre-charge
        uint256 maxCost = 10 ether; // 10 TON
        // 6.5 USDC * 1.03 = 6.695 USDC = 6_695_000 units
        uint256 expectedCharge = 6_695_000;

        usdc.mint(user, expectedCharge);
        vm.prank(user);
        usdc.approve(address(paymaster), expectedCharge);

        PackedUserOperation memory userOp = _buildUserOp(user, address(usdc));
        vm.prank(address(entryPoint));
        paymaster.validatePaymasterUserOp(userOp, bytes32(0), maxCost);

        assertEq(usdc.balanceOf(user), 0, "User's USDC fully pre-charged");
    }

    // ── postOp ────────────────────────────────────────────

    function test_PostOp_RefundsExcessETH() public {
        uint256 maxCost = 1 ether;
        uint256 preCharge = 0.0005 ether * 105 / 100; // ~0.000525 ETH

        l2Eth.mint(user, preCharge);
        vm.prank(user);
        l2Eth.approve(address(paymaster), preCharge);

        PackedUserOperation memory userOp = _buildUserOp(user, address(l2Eth));
        vm.prank(address(entryPoint));
        (bytes memory context,) = paymaster.validatePaymasterUserOp(userOp, bytes32(0), maxCost);

        // Actual gas = 0.5 TON (half of max)
        uint256 actualCost = 0.5 ether;
        uint256 actualCharge = (actualCost * 0.0005e18 / 1e18) * 105 / 100;
        uint256 expectedRefund = preCharge - actualCharge;

        vm.prank(address(entryPoint));
        paymaster.postOp(IPaymaster.PostOpMode.opSucceeded, context, actualCost, 0);

        assertEq(l2Eth.balanceOf(user), expectedRefund, "Excess refunded to user");
    }

    function test_PostOp_CollectedFeesUpdated() public {
        uint256 maxCost = 1 ether;
        uint256 preCharge = 0.0005 ether * 105 / 100;

        l2Eth.mint(user, preCharge);
        vm.prank(user);
        l2Eth.approve(address(paymaster), preCharge);

        PackedUserOperation memory userOp = _buildUserOp(user, address(l2Eth));
        vm.prank(address(entryPoint));
        (bytes memory context,) = paymaster.validatePaymasterUserOp(userOp, bytes32(0), maxCost);

        vm.prank(address(entryPoint));
        paymaster.postOp(IPaymaster.PostOpMode.opSucceeded, context, maxCost, 0); // full charge

        assertEq(paymaster.collectedFees(address(l2Eth)), preCharge, "All pre-charge collected");
    }

    function test_PostOp_RefundsOnOpReverted() public {
        // T5 variant: opReverted must still refund unused portion
        uint256 maxCost = 1 ether;
        uint256 preCharge = 0.0005 ether * 105 / 100;

        l2Eth.mint(user, preCharge);
        vm.prank(user);
        l2Eth.approve(address(paymaster), preCharge);

        PackedUserOperation memory userOp = _buildUserOp(user, address(l2Eth));
        vm.prank(address(entryPoint));
        (bytes memory context,) = paymaster.validatePaymasterUserOp(userOp, bytes32(0), maxCost);

        uint256 actualCost = 0.3 ether;
        vm.prank(address(entryPoint));
        paymaster.postOp(IPaymaster.PostOpMode.opReverted, context, actualCost, 0);

        assertGt(l2Eth.balanceOf(user), 0, "opReverted: still refunds unused gas");
    }

    // ── T5: 오라클 장애 ───────────────────────────────────

    function test_ValidatePaymasterUserOp_StaleOracleReverts() public {
        // T5: stale oracle causes revert
        vm.warp(block.timestamp + 86401); // 24h+1s stale

        uint256 preCharge = 0.0005 ether * 105 / 100;
        l2Eth.mint(user, preCharge);
        vm.prank(user);
        l2Eth.approve(address(paymaster), preCharge);

        PackedUserOperation memory userOp = _buildUserOp(user, address(l2Eth));
        vm.prank(address(entryPoint));
        vm.expectRevert("SimplePriceOracle: stale price");
        paymaster.validatePaymasterUserOp(userOp, bytes32(0), 1 ether);
    }

    // ── T8: 수거된 수수료 인출 ──────────────────────────

    function test_WithdrawCollectedFees_OwnerOnly() public {
        // First, accumulate fees via a full postOp
        _runFullUserOp(address(l2Eth), 1 ether, 0.0005 ether * 105 / 100);

        address to = makeAddr("operator");
        uint256 collected = paymaster.collectedFees(address(l2Eth));
        paymaster.withdrawCollectedFees(address(l2Eth), to, collected);

        assertEq(l2Eth.balanceOf(to), collected);
        assertEq(paymaster.collectedFees(address(l2Eth)), 0);
    }

    function test_WithdrawCollectedFees_InsufficientReverts() public {
        vm.expectRevert("insufficient collected");
        paymaster.withdrawCollectedFees(address(l2Eth), makeAddr("to"), 1);
    }

    function test_WithdrawCollectedFees_NonOwnerReverts() public {
        vm.prank(makeAddr("attacker"));
        vm.expectRevert();
        paymaster.withdrawCollectedFees(address(l2Eth), makeAddr("to"), 0);
    }

    // ── Helpers ──────────────────────────────────────────

    function _getTokenConfig(address token) internal view
        returns (bool enabled, ITokenPriceOracle oracle, uint256 markup, uint8 dec) {
        MultiTokenPaymaster.TokenConfig memory cfg = paymaster.getTokenConfig(token);
        return (cfg.enabled, cfg.oracle, cfg.markup, cfg.decimals);
    }

    // paymasterAndData = [paymaster(20)] [token(20)] — Phase 1 format (40 bytes, no signature)
    // Phase 2+: [paymaster(20)][token(20)][validUntil(6)][validAfter(6)][sig(65)] = 117 bytes (Appendix A)
    function _buildUserOp(address sender, address token) internal view
        returns (PackedUserOperation memory) {
        return PackedUserOperation({
            sender: sender,
            nonce: 0,
            initCode: "",
            callData: "",
            accountGasLimits: bytes32(uint256(100000) << 128 | uint256(100000)),
            preVerificationGas: 21000,
            gasFees: bytes32(uint256(1 gwei) << 128 | uint256(1 gwei)),
            paymasterAndData: abi.encodePacked(address(paymaster), token),
            signature: ""
        });
    }

    function _runFullUserOp(address token, uint256 maxCost, uint256 preCharge) internal {
        MockETH(token).mint(user, preCharge);
        vm.prank(user);
        IERC20(token).approve(address(paymaster), preCharge);
        PackedUserOperation memory userOp = _buildUserOp(user, token);
        vm.prank(address(entryPoint));
        (bytes memory context,) = paymaster.validatePaymasterUserOp(userOp, bytes32(0), maxCost);
        vm.prank(address(entryPoint));
        paymaster.postOp(IPaymaster.PostOpMode.opSucceeded, context, maxCost, 0);
    }
}
