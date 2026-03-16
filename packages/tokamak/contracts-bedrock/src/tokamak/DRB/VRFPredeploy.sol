// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Initializable } from "@openzeppelin/contracts/proxy/utils/Initializable.sol";

/// @title VRFPredeploy
/// @notice L2 predeploy contract for requesting verifiable random words.
///         DRB nodes listen to RandomWordsRequested events and fulfill via VRFCoordinator.
contract VRFPredeploy is Initializable {
    address public coordinator;

    event RandomWordsRequested(
        uint256 indexed requestId,
        address indexed requester,
        uint32 numWords,
        uint256 callbackGasLimit
    );

    /// @notice Initializer (called once by proxy).
    function initialize(address _coordinator) external initializer {
        coordinator = _coordinator;
    }

    /// @notice Request random words from DRB nodes.
    /// @param numWords     Number of random words requested (max 10).
    /// @param callbackGasLimit Gas limit for the consumer callback.
    /// @return requestId   Unique ID for this request.
    function requestRandomWords(
        uint32 numWords,
        uint256 callbackGasLimit
    ) external returns (uint256 requestId) {
        requestId = IVRFCoordinator(coordinator).requestRandomWords(
            msg.sender, numWords, callbackGasLimit
        );
        emit RandomWordsRequested(requestId, msg.sender, numWords, callbackGasLimit);
    }

    /// @notice Get the status and result of a randomness request.
    function getRequestStatus(uint256 requestId)
        external
        view
        returns (bool fulfilled, uint256[] memory randomWords)
    {
        return IVRFCoordinator(coordinator).getRequestStatus(requestId);
    }
}

interface IVRFCoordinator {
    function requestRandomWords(address requester, uint32 numWords, uint256 callbackGasLimit)
        external returns (uint256 requestId);
    function getRequestStatus(uint256 requestId)
        external view returns (bool fulfilled, uint256[] memory randomWords);
}
