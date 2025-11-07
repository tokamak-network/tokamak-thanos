// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { GameId } from "src/dispute/lib/Types.sol";
import { IProxyAdminOwnedBase } from "interfaces/L1/IProxyAdminOwnedBase.sol";
import { IReinitializableBase } from "interfaces/universal/IReinitializableBase.sol";
import { IDisputeGameFactory } from "interfaces/dispute/IDisputeGameFactory.sol";

/// @title IRAT
/// @notice Interface for the Randomized Attention Test contract
interface IRAT is IProxyAdminOwnedBase, IReinitializableBase {
    /// @notice Constructor function
    function __constructor__() external;
    /// @notice Challenger information structure
    /// @dev Packed to minimize storage slots
    struct ChallengerInfo {
        uint256 stakingAmount;      // Slot 1: 32 bytes
        uint256 slashedAmount;      // Slot 2: 32 bytes
        uint32 validatorIndex;      // Slot 3: 4 bytes
        bool isValid;               // Slot 3: 1 byte (packed)
    }

    /// @notice Attention test information structure
    /// @dev Packed to minimize storage slots
    struct AttentionInfo {
        // GameId removed - using mapping key instead
        bytes32 stateRoot;          // Slot 2: 32 bytes
        uint256 slashedAmount;      // Slot 3: 32 bytes
        address challengerAddress;  // Slot 4: 20 bytes
        uint64 l1BlockNumber;       // Slot 4: 8 bytes (packed with address)
        bool evidenceSubmitted;     // Slot 4: 1 byte (packed)
    }

    /// @notice Emitted when a challenger stakes ETH
    event ChallengerStaked(address indexed challenger, uint256 amount);

    /// @notice Emitted when attention test is triggered
    event AttentionTriggered(address indexed gameAddress, address indexed challenger);

    /// @notice Emitted when correct evidence is submitted
    event CorrectEvidenceSubmitted(
        address indexed gameAddress,
        address indexed challenger,
        uint256 restoredAmount
    );

    /// @notice Emitted when bonded amount is refunded through claim resolution
    event BondRefunded(address indexed gameAddress, address indexed challenger, uint256 refundedAmount);

    /// @notice Allows challengers to stake ETH
    function stake() external payable;

    /// @notice Gets challenger information
    /// @param _challenger Address of the challenger
    /// @return Challenger information
    function getChallengerInfo(address _challenger) external view returns (ChallengerInfo memory);



    /// @notice Gets number of valid challengers
    /// @return Number of valid challengers
    function getValidChallengerCount() external view returns (uint256);

    /// @notice Gets number of invalid challengers
    /// @return Number of invalid challengers
    function getInvalidChallengerCount() external view returns (uint256);

    /// @notice Triggers attention test (called by DisputeGameFactory)
    /// @param _gameAddress Game contract address
    /// @param _stateRoot State root to be verified
    /// @param _blockHash Block hash for validator selection
    function triggerAttentionTest(address _gameAddress, bytes32 _stateRoot, bytes32 _blockHash) external;

    /// @notice Submits correct evidence for attention test
    /// @param _gameAddress Game contract address
    /// @param _proofLV Left child state value
    /// @param _proofRV Right child state value
    function submitCorrectEvidence(address _gameAddress, bytes32 _proofLV, bytes32 _proofRV) external;

    /// @notice Called when a claim is resolved in FaultDisputeGame
    /// @param _claimant Address receiving the bond refund
    function resolveClaim(address _claimant) external;

    /// @notice Sets the per-test slashing amount
    /// @param _amount New slashing amount
    function setPerTestSlashingAmount(uint256 _amount) external;

    /// @notice Sets the evidence submission period
    /// @param _period New submission period in blocks
    function setEvidenceSubmissionPeriod(uint256 _period) external;

    /// @notice Sets the minimum staking balance
    /// @param _balance New minimum staking balance
    function setMinimumStakingBalance(uint256 _balance) external;

    /// @notice Returns the version
    function version() external view returns (string memory);

    /// @notice Initializes the RAT contract
    /// @param _disputeGameFactory The DisputeGameFactory contract address
    /// @param _perTestBondAmount The bond amount required per test
    /// @param _evidenceSubmissionPeriod The period for evidence submission
    /// @param _minimumStakingBalance The minimum staking balance required
    /// @param _ratTriggerProbability The probability of triggering RAT
    /// @param _manager The manager address
    function initialize(
        IDisputeGameFactory _disputeGameFactory,
        uint256 _perTestBondAmount,
        uint256 _evidenceSubmissionPeriod,
        uint256 _minimumStakingBalance,
        uint256 _ratTriggerProbability,
        address _manager
    ) external payable;
}