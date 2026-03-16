// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title VRFConsumerBase
/// @notice Abstract base contract for consuming VRF randomness.
///         Game contracts should inherit this and implement fulfillRandomWords().
abstract contract VRFConsumerBase {
    address private immutable VRF_COORDINATOR;

    constructor(address _vrfCoordinator) {
        VRF_COORDINATOR = _vrfCoordinator;
    }

    /// @notice Called by VRFCoordinator when randomness is fulfilled.
    /// @param requestId The ID of the randomness request.
    /// @param randomWords The random values generated.
    function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal virtual;

    /// @notice Called by VRFCoordinator to deliver randomness. Only callable by coordinator.
    function rawFulfillRandomWords(uint256 requestId, uint256[] memory randomWords) external {
        require(msg.sender == VRF_COORDINATOR, "VRFConsumerBase: only coordinator");
        fulfillRandomWords(requestId, randomWords);
    }
}
