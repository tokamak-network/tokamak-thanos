// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { ReentrancyGuard } from "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import { EmergencyExitProof } from "../libraries/EmergencyExitProof.sol";

/// @title AppExitCoordinator
/// @notice L1 coordinator for the resolver exit path.
///         Handles existing ERC-20 contracts that cannot be modified to implement
///         IEmergencyExitable. Users submit MPT storage proofs to verify their
///         L2 balance against a finalized state root.
///
///         Flow:
///         1. Tokamak registrar registers L2 tokens with their balance storage slot
///         2. After 48h timelock, anyone can activate the token
///         3. User collects proof off-chain via eth_getProof
///         4. User calls exitViaProof() with proof → EmergencyReserve pays out immediately
///
/// @dev This contract uses L2StateOracle to read finalized L2 state roots.
///      The L2StateOracle is an L1 contract that mirrors L2 block state roots.
contract AppExitCoordinator is ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ── Errors ────────────────────────────────────────────────────────────────
    error NotActive();
    error ProofInvalid();
    error AlreadyClaimed();
    error InvalidParameter();
    error NotRegistrar();
    error NotGuardian();
    error TimelockNotElapsed();
    error UnknownBlock();
    error EmergencyNotDeclared();

    // ── State ─────────────────────────────────────────────────────────────────
    struct TokenRecord {
        address l1Token;
        uint256 storageSlot;
        bool active;
        uint64 activatesAt;
    }

    /// @notice L2StateOracle on L1 — provides finalized L2 state roots.
    address public immutable l2StateOracle;

    mapping(address => TokenRecord) public registry;

    address public registrar;
    address public guardian;
    uint256 public constant REGISTRATION_DELAY = 48 hours;

    /// @notice The single finalized L2 block height that all proofs are checked against.
    ///         Zero until the guardian declares the emergency snapshot. Pinning every exit to one
    ///         protocol-set block (instead of a caller-supplied block) is what prevents the same
    ///         L2 balance from being proven and claimed at many different blocks.
    uint256 public emergencyBlockNumber;

    mapping(bytes32 => bool) public exitClaims;

    // ── Events ────────────────────────────────────────────────────────────────
    event TokenRegistered(
        address indexed l2Token,
        address indexed l1Token,
        uint256 storageSlot,
        uint64 activatedAt
    );
    event TokenActivated(address indexed l2Token, uint64 activatedAt);
    event TokenCancelled(address indexed l2Token);
    event ExitClaimed(address indexed l2Token, address indexed user, uint256 amount);
    event ReserveRecovered(address indexed token, address indexed to, uint256 amount);
    event RegistrarUpdated(address indexed newRegistrar);
    event GuardianUpdated(address indexed newGuardian);
    event EmergencyDeclared(uint256 blockNumber);

    // ── Modifiers ─────────────────────────────────────────────────────────────
    modifier onlyRegistrar() {
        if (msg.sender != registrar) revert NotRegistrar();
        _;
    }

    modifier onlyGuardian() {
        if (msg.sender != guardian) revert NotGuardian();
        _;
    }

    // ── Constructor ───────────────────────────────────────────────────────────
    constructor(
        address _l2StateOracle,
        address _registrar,
        address _guardian
    ) {
        if (_l2StateOracle == address(0)) revert InvalidParameter();
        if (_registrar == address(0)) revert InvalidParameter();
        if (_guardian == address(0)) revert InvalidParameter();

        l2StateOracle = _l2StateOracle;
        registrar = _registrar;
        guardian = _guardian;
    }

    // ── Token Registration ────────────────────────────────────────────────────

    function registerToken(
        address l2Token,
        address l1Token,
        uint256 storageSlot
    ) external onlyRegistrar {
        TokenRecord storage rec = registry[l2Token];
        if (rec.active) revert NotActive();

        uint64 activatedAt = uint64(block.timestamp) + uint64(REGISTRATION_DELAY);
        rec.l1Token = l1Token;
        rec.storageSlot = storageSlot;
        rec.activatesAt = activatedAt;

        emit TokenRegistered(l2Token, l1Token, storageSlot, activatedAt);
    }

    function cancelRegistration(address l2Token) external onlyGuardian {
        delete registry[l2Token];
        emit TokenCancelled(l2Token);
    }

    function activateToken(address l2Token) external {
        TokenRecord storage rec = registry[l2Token];
        if (rec.active) revert NotActive();
        if (block.timestamp < rec.activatesAt) revert TimelockNotElapsed();

        rec.active = true;
        rec.activatesAt = uint64(block.timestamp);
        emit TokenActivated(l2Token, rec.activatesAt);
    }

    // ── Emergency Exit ────────────────────────────────────────────────────────

    /// @notice Declare the emergency snapshot block that all proofs are checked against.
    ///         Guardian-only. Reverts unless the L2StateOracle already exposes a non-zero state
    ///         root for `blockNumber`, so exits can only ever be enabled against a finalized block.
    /// @param blockNumber The finalized L2 block height to freeze exits at.
    function declareEmergency(uint256 blockNumber) external onlyGuardian {
        _getStateRoot(blockNumber); // reverts UnknownBlock if the root is unknown / zero
        emergencyBlockNumber = blockNumber;
        emit EmergencyDeclared(blockNumber);
    }

    /// @notice Bundle proof data to reduce exitViaProof stack depth.
    /// @dev No block number: all proofs are checked against the protocol-set `emergencyBlockNumber`.
    struct ExitProofRequest {
        address user;
        uint256 amount;
        bytes accountStateRlp;
        bytes[] accountProof;
        bytes[] storageProof;
    }

    /// @notice Execute an exit using a Merkle proof of the user's L2 balance.
    ///         Verifies the proof against a finalized state root and pays out
    ///         from the emergency reserve immediately (no 7-day DTD).
    ///
    ///         Off-chain flow:
    ///         1. User calls eth_getProof(l2Token, [balanceSlot], emergencyBlockNumber) on L2 RPC
    ///         2. Response includes accountProof, storageProof, and account state
    ///         3. User submits exitViaProof() on L1 with the proof data
    ///
    /// @dev The request struct bundles all proof parameters to keep stack depth low.
    function exitViaProof(
        address l2Token,
        ExitProofRequest calldata req
    ) external nonReentrant {
        TokenRecord storage rec = registry[l2Token];
        if (!rec.active) revert NotActive();
        if (emergencyBlockNumber == 0) revert EmergencyNotDeclared();

        // 1. Get the finalized state root for the protocol-set emergency snapshot block.
        bytes32 stateRoot = _getStateRoot(emergencyBlockNumber);

        // 2. Verify MPT proofs (account state + storage slot)
        bytes32 storageKey = keccak256(abi.encode(req.user, rec.storageSlot));
        EmergencyExitProof.verifyProofs(
            stateRoot,
            l2Token,
            req.accountStateRlp,
            req.accountProof,
            storageKey,
            req.storageProof,
            req.amount
        );

        // 3. Prevent double claims. The key is (token, user) only — the snapshot block is fixed,
        //    so the same balance cannot be re-proven at a different block to bypass this guard.
        bytes32 claimKey = keccak256(abi.encode(l2Token, req.user));
        if (exitClaims[claimKey]) revert AlreadyClaimed();
        exitClaims[claimKey] = true;

        // 4. Pay out from reserve
        IERC20(rec.l1Token).safeTransfer(req.user, req.amount);

        emit ExitClaimed(l2Token, req.user, req.amount);
    }

    // ── Reserve Management ────────────────────────────────────────────────────

    /// @notice Health check for the emergency reserve of a token.
    /// @return balance  Current L1 token balance held by this contract.
    /// @return minimum  Minimum required reserve (placeholder: zero).
    /// @return healthy  True if reserve balance >= minimum.
    function reserveHealthCheck(address l1Token)
        external
        view
        returns (uint256 balance, uint256 minimum, bool healthy)
    {
        balance = IERC20(l1Token).balanceOf(address(this));
        minimum = 0;
        healthy = balance >= minimum;
    }

    /// @notice Recover unused reserve tokens after an emergency event.
    /// @param token  The L1 ERC-20 token to recover.
    /// @param amount The amount to recover.
    function recoverReserve(address token, uint256 amount) external onlyRegistrar {
        IERC20(token).safeTransfer(msg.sender, amount);
        emit ReserveRecovered(token, msg.sender, amount);
    }

    // ── Admin ─────────────────────────────────────────────────────────────────

    function setRegistrar(address newRegistrar) external {
        if (msg.sender != registrar) revert NotRegistrar();
        registrar = newRegistrar;
        emit RegistrarUpdated(newRegistrar);
    }

    function setGuardian(address newGuardian) external {
        if (msg.sender != registrar) revert NotRegistrar();
        guardian = newGuardian;
        emit GuardianUpdated(newGuardian);
    }

    // ── Internal ──────────────────────────────────────────────────────────────

    function _getStateRoot(uint256 blockNumber) internal view returns (bytes32 stateRoot) {
        (bool success, bytes memory data) = l2StateOracle.staticcall(
            abi.encodeWithSignature("getStateRoot(uint256)", blockNumber)
        );
        if (!success || data.length < 32) revert UnknownBlock();
        stateRoot = abi.decode(data, (bytes32));
        if (stateRoot == bytes32(0)) revert UnknownBlock();
    }
}
