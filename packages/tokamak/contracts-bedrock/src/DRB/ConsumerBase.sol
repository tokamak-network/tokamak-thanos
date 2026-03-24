// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {ICommitReveal2} from "./ICommitReveal2.sol";

/**
 * @notice Abstract contract for contracts using CommitReveal2 DRB randomness
 * ********************************************************************************
 *
 * @dev Consumer contracts must inherit from ConsumerBase, and can
 * @dev initialize Coordinator address in their constructor.
 * @dev Call the `_requestRandomNumber` function to request a random number.
 * @dev Consumers must implement the fullfillRandomWords function, which will be called during
 * @dev fulfillment with the randomness result.
 */
abstract contract ConsumerBase {
    error OnlyCoordinatorCanFulfill(address have, address want);
    error InsufficientBalance();
    /// @dev The RNGCoordinator contract

    ICommitReveal2 internal s_commitreveal2;

    /**
     * @param rngCoordinator The address of the RNGCoordinator contract
     */
    constructor(address rngCoordinator) {
        s_commitreveal2 = ICommitReveal2(rngCoordinator);
    }

    receive() external payable virtual {}

    /**
     * @return requestId The ID of the request
     * @dev Request Randomness to the Coordinator
     */
    function _requestRandomNumber(uint32 callbackGasLimit) internal returns (uint256, uint256) {
        uint256 requestFee = s_commitreveal2.estimateRequestPrice(callbackGasLimit, tx.gasprice);
        require(requestFee <= address(this).balance, InsufficientBalance());
        uint256 requestId = s_commitreveal2.requestRandomNumber{value: requestFee}(callbackGasLimit);
        return (requestId, requestFee);
    }

    function _refund(uint256 round) internal {
        s_commitreveal2.refund(round);
    }

    /**
     * @param round The round of the randomness
     * @param randomNumber the random number
     * @dev Callback function for the Coordinator to call after the request is fulfilled.  Override this function in your contract
     */
    function fulfillRandomRandomNumber(uint256 round, uint256 randomNumber) internal virtual;

    /**
     * @param round The round of the randomness
     * @param randomNumber The random number
     * @dev Callback function for the Coordinator to call after the request is fulfilled. This function is called by the Coordinator, 0x00fc98b8
     */
    function rawFulfillRandomNumber(uint256 round, uint256 randomNumber) external {
        require(msg.sender == address(s_commitreveal2), OnlyCoordinatorCanFulfill(msg.sender, address(s_commitreveal2)));
        fulfillRandomRandomNumber(round, randomNumber);
    }
}
