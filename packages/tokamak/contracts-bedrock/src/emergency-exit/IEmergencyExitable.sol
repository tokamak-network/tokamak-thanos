// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title IEmergencyExitable
/// @notice Interface for L2 contracts that support emergency self-withdrawal to L1.
///         Implemented by app contracts (DEX pools, staking contracts, games)
///         to allow users to force-exit locked positions to L1 via Force TX.
interface IEmergencyExitable {
    /// @notice Emitted when a user successfully initiates an emergency exit.
    event EmergencyExited(address indexed user, address indexed token, uint256 amount);

    /// @notice Execute emergency withdrawal for the caller to L1.
    /// @dev Must be callable without operator gates (no whenNotPaused, onlyActive, etc.).
    ///      Implementations should treat msg.sender as the L2 owner identity. OptimismPortal
    ///      aliases only L1 contract senders, so an EOA Force TX arrives un-aliased; undoing the
    ///      alias unconditionally would corrupt the caller. Do not apply undoL1ToL2Alias.
    function emergencyExit() external;

    /// @notice Returns the L1/L2 token pair this contract manages for exits.
    /// @return l1Token  The ERC-20 address on L1 (the destination of funds).
    /// @return l2Token  The ERC-20 address on L2 (the source of funds).
    function exitToken() external view returns (address l1Token, address l2Token);

    /// @notice Returns the exit amount a user would receive right now.
    /// @param user The address to query.
    /// @return amount The user's exitable balance.
    function exitPreview(address user) external view returns (uint256 amount);
}
