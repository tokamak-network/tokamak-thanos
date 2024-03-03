// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Types } from "src/libraries/Types.sol";

interface IL2OutputOracle   {
    function  SUBMISSION_INTERVAL() external view returns (uint256);
    function  L2_BLOCK_TIME() external view returns (uint256);
    function  FINALIZATION_PERIOD_SECONDS() external view returns (uint256);
    function  startingBlockNumber() external view returns (uint256);
    function  startingTimestamp() external view returns (uint256);
    function  challenger() external view returns (address);
    function  proposer() external view returns (address);
    function  version() external view returns (string memory);

    function initialize(
        uint256 _startingBlockNumber,
        uint256 _startingTimestamp,
        address _proposer,
        address _challenger
    )
        external ;

    function submissionInterval() external view returns (uint256);

    function l2BlockTime() external view returns (uint256);

    /// @notice Getter for the finalization period.
    function finalizationPeriodSeconds() external view returns (uint256) ;
    function CHALLENGER() external view returns (address)  ;
    function PROPOSER() external view returns (address);

    function deleteL2Outputs(uint256 _l2OutputIndex) external ;

    function proposeL2Output(
        bytes32 _outputRoot,
        uint256 _l2BlockNumber,
        bytes32 _l1BlockHash,
        uint256 _l1BlockNumber
    )
        external ;

    function getL2Output(uint256 _l2OutputIndex) external view returns (Types.OutputProposal memory)  ;

    function getL2OutputIndexAfter(uint256 _l2BlockNumber) external view returns (uint256)  ;


    function getL2OutputAfter(uint256 _l2BlockNumber) external view returns (Types.OutputProposal memory)   ;


    function latestOutputIndex() external view returns (uint256)  ;

    function nextOutputIndex() external view returns (uint256)  ;


    function latestBlockNumber() external view returns (uint256)  ;

    function nextBlockNumber() external view returns (uint256) ;


    function computeL2Timestamp(uint256 _l2BlockNumber) external view returns (uint256)  ;
}
