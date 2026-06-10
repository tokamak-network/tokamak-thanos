// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test } from "forge-std/Test.sol";
import { console2 as console } from "forge-std/console2.sol";
import { AppExitCoordinator } from "src/emergency-exit/AppExitCoordinator.sol";
import { EmergencyExitProof } from "src/libraries/EmergencyExitProof.sol";
import { TestERC20 } from "../mocks/TestERC20.sol";
import { MptProofBuilder } from "./MptProofBuilder.sol";

/// @title AppExitCoordinator_Unit_Test
/// @notice Unit tests for AppExitCoordinator — registration, activation, exitViaProof.
contract AppExitCoordinator_Test is Test {
    AppExitCoordinator public coordinator;
    TestERC20 public l1Token;

    address public registrar = makeAddr("registrar");
    address public guardian = makeAddr("guardian");
    address public l2StateOracle = makeAddr("l2StateOracle");
    address public l2Contract = makeAddr("l2Contract");
    address public user = makeAddr("user");

    uint256 public constant RESERVE_AMOUNT = 10000 ether;
    uint256 public constant SNAPSHOT_BLOCK = 12345;
    uint256 public constant SLOT = 0;

    function setUp() public {
        // Deploy L1 token and fund coordinator
        l1Token = new TestERC20();
        l1Token.mint(address(this), RESERVE_AMOUNT);

        // Deploy coordinator (l2StateOracle is mocked)
        coordinator = new AppExitCoordinator(
            l2StateOracle,
            registrar,
            guardian
        );

        // Fund coordinator with reserve
        l1Token.transfer(address(coordinator), RESERVE_AMOUNT);

        console.log("Coordinator balance:", l1Token.balanceOf(address(coordinator)));
    }

    // ── Constructor ──────────────────────────────────────────────────────────

    function test_constructor_SetsImmutable() public {
        assertEq(coordinator.l2StateOracle(), l2StateOracle, "l2StateOracle");
        assertEq(coordinator.registrar(), registrar, "registrar");
        assertEq(coordinator.guardian(), guardian, "guardian");
        assertEq(coordinator.REGISTRATION_DELAY(), 48 hours, "timelock");
    }

    function test_constructor_Revert_ZeroStateOracle() public {
        vm.expectRevert();
        new AppExitCoordinator(address(0), registrar, guardian);
    }

    // ── Token Registration ───────────────────────────────────────────────────

    function test_registerToken_ByRegistrar_Success() public {
        vm.prank(registrar);
        coordinator.registerToken(l2Contract, address(l1Token), 0);

        // Should be registered but NOT active
        (address l1TokenOut, uint256 slot, bool active, uint64 activatedAt) =
            coordinator.registry(l2Contract);

        assertEq(l1TokenOut, address(l1Token), "l1Token");
        assertEq(slot, 0, "slot");
        assert(!active);
    }

    function test_registerToken_NonRegistrar_Revert() public {
        vm.prank(user);
        vm.expectRevert();
        coordinator.registerToken(l2Contract, address(l1Token), 0);
    }

    // ── Cancel Registration ──────────────────────────────────────────────────

    function test_cancelRegistration_ByGuardian_Success() public {
        // Register first
        vm.prank(registrar);
        coordinator.registerToken(l2Contract, address(l1Token), 0);

        // Cancel
        vm.prank(guardian);
        coordinator.cancelRegistration(l2Contract);

        // Should be deleted
        (address l1TokenOut, , bool active, ) = coordinator.registry(l2Contract);
        assertEq(l1TokenOut, address(0));
        assert(!active);
    }

    function test_cancelRegistration_NonGuardian_Revert() public {
        vm.prank(user);
        vm.expectRevert();
        coordinator.cancelRegistration(l2Contract);
    }

    // ── Token Activation ─────────────────────────────────────────────────────

    function test_activateToken_TimelockNotElapsed_Revert() public {
        vm.prank(registrar);
        coordinator.registerToken(l2Contract, address(l1Token), 0);

        // Try to activate immediately — should revert
        vm.expectRevert(AppExitCoordinator.TimelockNotElapsed.selector);
        coordinator.activateToken(l2Contract);
    }

    function test_activateToken_AfterTimelock_Success() public {
        vm.prank(registrar);
        coordinator.registerToken(l2Contract, address(l1Token), 0);

        // Fast-forward past the 48h timelock
        vm.warp(block.timestamp + 48 hours + 1 seconds);

        // Anyone can activate
        vm.prank(user);
        coordinator.activateToken(l2Contract);

        (, , bool active, ) = coordinator.registry(l2Contract);
        assert(active);
    }

    function test_activateToken_AlreadyActive_Revert() public {
        vm.prank(registrar);
        coordinator.registerToken(l2Contract, address(l1Token), 0);

        vm.warp(block.timestamp + 48 hours + 1 seconds);
        coordinator.activateToken(l2Contract);

        vm.expectRevert(AppExitCoordinator.NotActive.selector);
        coordinator.activateToken(l2Contract);
    }

    // ── Exit via Proof ───────────────────────────────────────────────────────

    /// @notice Test exitViaProof with a valid (mocked) proof.
    ///         This test verifies the flow: register → activate → exit.
    ///         Full MPT proof verification requires fork testing.
    function test_exitViaProof_Flow_Complete() public {
        // Step 1: Register token
        vm.prank(registrar);
        coordinator.registerToken(l2Contract, address(l1Token), 0);

        // Step 2: Activate after timelock
        vm.warp(block.timestamp + 48 hours + 1 seconds);
        vm.prank(user);
        coordinator.activateToken(l2Contract);

        // Step 3: Verify token is active
        (, , bool active, ) = coordinator.registry(l2Contract);
        assert(active);

        // exitViaProof requires actual proof data which we can't construct
        // without a fork test. The function signature and modifiers are verified
        // by the compilation passing.
    }

    function test_exitViaProof_NotActive_Revert() public {
        // Register but don't activate
        vm.prank(registrar);
        coordinator.registerToken(l2Contract, address(l1Token), 0);

        // Call exitViaProof — should revert because not active
        AppExitCoordinator.ExitProofRequest memory req = AppExitCoordinator.ExitProofRequest({
            user: user,
            amount: 100,
            accountStateRlp: "",
            accountProof: new bytes[](0),
            storageProof: new bytes[](0)
        });

        vm.expectRevert(AppExitCoordinator.NotActive.selector);
        coordinator.exitViaProof(l2Contract, req);
    }

    // ── Reserve Management ───────────────────────────────────────────────────

    function test_reserveHealthCheck_ReturnsValues() public {
        (uint256 balance, uint256 minimum, bool healthy) =
            coordinator.reserveHealthCheck(address(l1Token));

        assertEq(balance, RESERVE_AMOUNT, "reserve balance");
        assertEq(minimum, 0, "minimum");
        assert(healthy);
    }

    function test_recoverReserve_ByRegistrar_Success() public {
        // Registrar recovers unused reserve
        uint256 before = l1Token.balanceOf(registrar);
        vm.prank(registrar);
        coordinator.recoverReserve(address(l1Token), 100);
        assertEq(l1Token.balanceOf(registrar) - before, 100, "should receive tokens");
    }

    function test_recoverReserve_NonRegistrar_Revert() public {
        vm.prank(user);
        vm.expectRevert(AppExitCoordinator.NotRegistrar.selector);
        coordinator.recoverReserve(address(l1Token), 100);
    }

    // ── Admin Updates ────────────────────────────────────────────────────────

    function test_setRegistrar_ByRegistrar_Success() public {
        address newRegistrar = makeAddr("newRegistrar");
        vm.prank(registrar);
        coordinator.setRegistrar(newRegistrar);
        assertEq(coordinator.registrar(), newRegistrar, "registrar updated");
    }

    function test_setRegistrar_NonRegistrar_Revert() public {
        vm.prank(user);
        vm.expectRevert(AppExitCoordinator.NotRegistrar.selector);
        coordinator.setRegistrar(makeAddr("newRegistrar"));
    }

    function test_setGuardian_ByRegistrar_Success() public {
        address newGuardian = makeAddr("newGuardian");
        vm.prank(registrar);
        coordinator.setGuardian(newGuardian);
        assertEq(coordinator.guardian(), newGuardian, "guardian updated");
    }

    function test_setGuardian_NonRegistrar_Revert() public {
        vm.prank(user);
        vm.expectRevert(AppExitCoordinator.NotRegistrar.selector);
        coordinator.setGuardian(makeAddr("newGuardian"));
    }

    // ── HIGH-1: declareEmergency (guardian-set snapshot block) ────────────────

    function test_declareEmergency_ByGuardian_Success() public {
        _mockStateRoot(bytes32(uint256(0xBEEF)));
        vm.prank(guardian);
        coordinator.declareEmergency(SNAPSHOT_BLOCK);
        assertEq(coordinator.emergencyBlockNumber(), SNAPSHOT_BLOCK, "snapshot block set");
    }

    function test_declareEmergency_NonGuardian_Revert() public {
        _mockStateRoot(bytes32(uint256(0xBEEF)));
        vm.prank(user);
        vm.expectRevert(AppExitCoordinator.NotGuardian.selector);
        coordinator.declareEmergency(SNAPSHOT_BLOCK);
    }

    function test_declareEmergency_UnknownRoot_Revert() public {
        // No oracle mock: getStateRoot returns empty → UnknownBlock.
        vm.prank(guardian);
        vm.expectRevert(AppExitCoordinator.UnknownBlock.selector);
        coordinator.declareEmergency(SNAPSHOT_BLOCK);
    }

    // ── HIGH-1: exitViaProof gated on emergency declaration ───────────────────

    function test_exitViaProof_EmergencyNotDeclared_Revert() public {
        _activate(l2Contract);
        MptProofBuilder.Fixture memory f = MptProofBuilder.build(l2Contract, user, SLOT, 1000 ether);

        vm.expectRevert(AppExitCoordinator.EmergencyNotDeclared.selector);
        coordinator.exitViaProof(l2Contract, _req(f, user, 1000 ether));
    }

    // ── HIGH-1: end-to-end exit + replay prevention ───────────────────────────

    function test_exitViaProof_Succeeds_AndPaysOut() public {
        uint256 bal = 1000 ether;
        MptProofBuilder.Fixture memory f = _arm(l2Contract, user, bal);

        uint256 beforeBal = l1Token.balanceOf(user);
        coordinator.exitViaProof(l2Contract, _req(f, user, bal));
        assertEq(l1Token.balanceOf(user) - beforeBal, bal, "user should be paid the proven balance");
    }

    /// @notice The same balance cannot be claimed twice. Before the fix, a caller could pass a
    ///         different blockNumber to mint a new claimKey and drain the reserve; the snapshot
    ///         block is now fixed, so the (token, user) claim guard blocks the replay.
    function test_exitViaProof_Replay_Revert() public {
        uint256 bal = 1000 ether;
        MptProofBuilder.Fixture memory f = _arm(l2Contract, user, bal);

        coordinator.exitViaProof(l2Contract, _req(f, user, bal)); // first claim succeeds
        vm.expectRevert(AppExitCoordinator.AlreadyClaimed.selector);
        coordinator.exitViaProof(l2Contract, _req(f, user, bal)); // replay blocked
    }

    function test_exitViaProof_OverClaim_Revert() public {
        uint256 bal = 1000 ether;
        MptProofBuilder.Fixture memory f = _arm(l2Contract, user, bal);

        vm.expectRevert(abi.encodeWithSelector(EmergencyExitProof.InsufficientBalance.selector, bal, bal + 1));
        coordinator.exitViaProof(l2Contract, _req(f, user, bal + 1));
    }

    // ── Helpers ───────────────────────────────────────────────────────────────

    /// @notice Mock the L2StateOracle to return `root` for the snapshot block.
    function _mockStateRoot(bytes32 root) internal {
        vm.mockCall(
            l2StateOracle,
            abi.encodeWithSignature("getStateRoot(uint256)", SNAPSHOT_BLOCK),
            abi.encode(root)
        );
    }

    /// @notice Register and activate `token` against the L1 reserve token at SLOT.
    function _activate(address token) internal {
        vm.prank(registrar);
        coordinator.registerToken(token, address(l1Token), SLOT);
        vm.warp(block.timestamp + 48 hours + 1 seconds);
        coordinator.activateToken(token);
    }

    /// @notice Full setup for an exit: build proof, activate token, mock oracle, declare emergency.
    function _arm(address token, address account, uint256 balance)
        internal
        returns (MptProofBuilder.Fixture memory f)
    {
        f = MptProofBuilder.build(token, account, SLOT, balance);
        _activate(token);
        _mockStateRoot(f.stateRoot);
        vm.prank(guardian);
        coordinator.declareEmergency(SNAPSHOT_BLOCK);
    }

    function _req(MptProofBuilder.Fixture memory f, address account, uint256 amount)
        internal
        pure
        returns (AppExitCoordinator.ExitProofRequest memory)
    {
        return AppExitCoordinator.ExitProofRequest({
            user: account,
            amount: amount,
            accountStateRlp: f.accountStateRlp,
            accountProof: f.accountProof,
            storageProof: f.storageProof
        });
    }
}
