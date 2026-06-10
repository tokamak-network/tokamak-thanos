// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { IEmergencyExitable } from "./IEmergencyExitable.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { ReentrancyGuard } from "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import { Predeploys } from "../libraries/Predeploys.sol";

/// @title EmergencyExitableBase
/// @notice Abstract base contract that implements IEmergencyExitable for L2 app contracts.
///         Provides operator-gate-free emergency withdrawal path via Force TX:
///         1. Resolve L1→L2 alias (if called via Force TX)
///         2. Unwind user position to tokens + amount
///         3. Push tokens to L2StandardBridge for L1 withdrawal
/// @dev Subcontracts must implement _unwindPosition to define contract-specific
///      position resolution logic (DEX LP burn, staking unstake, etc.).
abstract contract EmergencyExitableBase is IEmergencyExitable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    /// @notice Error thrown when user has no balance to exit.
    error ZeroBalance();

    /// @notice Error thrown when position unwind returns zero amount.
    error ExitEmpty();

    /// @dev L2StandardBridge predeploy address — immutable constant, no import dependency.
    address internal constant L2_STANDARD_BRIDGE = Predeploys.L2_STANDARD_BRIDGE;

    /**
     * @notice Unwind the user's position and return (token, amount).
     *         Must transfer the user's position tokens into this contract before returning.
     * @param user The address whose position to unwind (already resolved for alias).
     * @return token   The L2 ERC-20 token address to bridge out.
     * @return amount  The amount of token to withdraw.
     */
    function _unwindPosition(address user)
        internal
        virtual
        returns (address token, uint256 amount);

    /**
     * @notice Execute emergency withdrawal for the caller.
     *         Unwinds the caller's position and pushes tokens to L2StandardBridge.
     * @dev This function intentionally has NO operator controls (no whenNotPaused, onlyActive).
     *      The caller is msg.sender as-is. OptimismPortal aliases only L1 *contract* senders
     *      (from = applyL1ToL2Alias(sender) iff sender != tx.origin), so an EOA Force TX arrives
     *      un-aliased and a direct L2 call is already the owner; in both cases msg.sender is the
     *      correct L2 identity. For an aliased L1 contract sender, the aliased address is its L2
     *      identity. No alias undo is applied (undoing it unconditionally corrupts EOA callers).
     */
    function emergencyExit() external override nonReentrant {
        address caller = msg.sender;

        // Unwind position — this should:
        // 1. Transfer user's position tokens (LP tokens, staking tokens, etc.) into this contract
        // 2. Resolve them to underlying assets
        // 3. Return the token address and amount to bridge out
        (address token, uint256 amount) = _unwindPosition(caller);
        if (amount == 0) revert ExitEmpty();

        // Approve L2StandardBridge to pull the tokens, then trigger L1 withdrawal.
        // The bridge will create a withdrawal message that can be proved + finalized on L1
        // after the standard Fault Proof challenge period (7 days on mainnet).
        IERC20(token).approve(L2_STANDARD_BRIDGE, amount);
        IL2StandardBridge(L2_STANDARD_BRIDGE).withdraw(
            token,
            amount,
            200_000, // gasLimit for the withdrawal message
            ""       // calldata (not used for standard ERC-20 withdrawals)
        );

        emit EmergencyExited(caller, token, amount);
    }

    /**
     * @notice Returns the L1/L2 token pair for this contract.
     *         Subclasses should override with their specific token mapping.
     */
    function exitToken() external view virtual override returns (address l1Token, address l2Token) {
        l1Token = address(0);
        l2Token = address(0);
    }

    /**
     * @notice Returns the exit preview for a user.
     *         Subclasses should override with their specific balance calculation.
     */
    function exitPreview(address user) external view virtual override returns (uint256 amount) {
        amount = 0;
    }
}

/// @notice Interface for the L2StandardBridge predeploy.
///         Minimal ABI for the withdraw() function only.
interface IL2StandardBridge {
    /// @notice Initiates a withdrawal of tokens from L2 to L1.
    /// @param _token     L2 ERC-20 token address.
    /// @param _amount    Amount to withdraw.
    /// @param _gasLimit  Gas limit for the L1 bridge call.
    /// @param _data      Additional calldata (unused for standard withdrawals).
    function withdraw(address _token, uint256 _amount, uint256 _gasLimit, bytes calldata _data) external;
}
