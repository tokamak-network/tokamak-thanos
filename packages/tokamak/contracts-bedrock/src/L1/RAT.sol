// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Contracts
import { ProxyAdminOwnedBase } from "src/L1/ProxyAdminOwnedBase.sol";
import { ReinitializableBase } from "src/universal/ReinitializableBase.sol";
import { ReentrancyGuard } from "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import { Initializable } from "@openzeppelin/contracts/proxy/utils/Initializable.sol";

// Libraries
import { GameId, LibGameId } from "src/dispute/lib/Types.sol";

// Interfaces
import { ISemver } from "interfaces/universal/ISemver.sol";
import { IDisputeGameFactory } from "interfaces/dispute/IDisputeGameFactory.sol";

/// @custom:proxied true
/// @title RAT (Original Version)
/// @notice Randomized Attention Test contract for challenger monitoring and testing - Original Implementation
contract RAT is ProxyAdminOwnedBase, ReinitializableBase, Initializable, ReentrancyGuard, ISemver {
    /// @notice Challenger information structure
    /// @dev Packed to minimize storage slots
    struct ChallengerInfo {
        uint256 stakingAmount;      // Slot 1: 32 bytes
        uint256 totalSlashedAmount; // Slot 2: 32 bytes (total slashed amount for this challenger)
        uint32 validatorIndex;      // Slot 3: 4 bytes
        bool isValid;               // Slot 3: 1 byte (packed)
    }

    /// @notice Attention test information structure
    /// @dev Packed to minimize storage slots - 3 slots total
    struct AttentionInfo {
        bytes32 stateRoot;          // Slot 1: 32 bytes
        uint96 bondAmount;          // Slot 2: 12 bytes (packed with challengerAddress)
        address challengerAddress;  // Slot 2: 20 bytes (packed with bondAmount)
        uint64 l1BlockNumber;       // Slot 3: 8 bytes
        bool evidenceSubmitted;     // Slot 3: 1 byte (packed)
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

    /// @notice Semantic version
    /// @custom:semver 1.0.0-beta.1
    string public constant version = "1.0.0-beta.1";

    /// @notice DisputeGameFactory contract address
    IDisputeGameFactory public disputeGameFactory;

    /// @notice Bond amount per attention test
    uint256 public perTestBondAmount;

    /// @notice Evidence submission period in blocks
    uint256 public evidenceSubmissionPeriod;

    /// @notice Minimum staking balance required
    uint256 public minimumStakingBalance;

    /// @notice Maximum probability value (100,000 for extended probability range)
    uint256 private constant MAX_PROBABILITY = 100_000;

    /// @notice Probability for triggering RAT (0-50400, where 50400 = weekly)
    uint256 public ratTriggerProbability;

    /// @notice Manager address that can modify RAT parameters
    address public ratManager;

    /// @notice Mapping from challenger address to challenger info
    mapping(address => ChallengerInfo) public challengers;

    /// @notice Mapping from FaultDisputeGame address to attention info
    mapping(address => AttentionInfo) public attentionTests;

    /// @notice Array of valid challengers (index 0 is reserved for "not found")
    address[] public validChallengers;


    /// @notice Error thrown when caller is not the DisputeGameFactory
    error NotDisputeGameFactory();

    /// @notice Error thrown when challenger does not exist
    error ChallengerNotExists();

    /// @notice Error thrown when evidence submission period has expired
    error EvidenceSubmissionExpired();

    /// @notice Error thrown when caller is not the ratManager
    error NotRatManager();

    /// @notice Error thrown when evidence has already been submitted
    error EvidenceAlreadySubmitted();

    /// @notice Error thrown when caller is not the correct challenger
    error InvalidChallengerAddress();



    /// @notice Error thrown when proof verification fails
    error ProofVerificationFailed();

    /// @notice Error thrown when attention test does not exist
    error AttentionTestNotExists();

    /// @notice Error thrown when insufficient staking amount
    error InsufficientStakingAmount();

    /// @notice Error thrown when there are no valid challengers available
    error NoValidChallengers();

    /// @notice Modifier to restrict access to DisputeGameFactory only
    modifier onlyDisputeGameFactory() {
        if (msg.sender != address(disputeGameFactory)) revert NotDisputeGameFactory();
        _;
    }

    /// @notice Modifier to restrict access to manager only
    modifier onlyManager() {
        if (msg.sender != ratManager) revert NotRatManager();
        _;
    }

    /// @notice Modifier to restrict access to RAT manager only
    modifier onlyRatManager() {
        if (msg.sender != ratManager) revert NotRatManager();
        _;
    }

    /// @notice Constructs the RAT contract
    constructor() ReinitializableBase(2) {
        _disableInitializers();
    }

    /// @notice Initializes the contract
    /// @param _disputeGameFactory Address of the DisputeGameFactory contract
    /// @param _perTestBondAmount Bond amount per attention test
    /// @param _evidenceSubmissionPeriod Evidence submission period in blocks
    /// @param _minimumStakingBalance Minimum staking balance required
    /// @param _ratTriggerProbability Initial RAT trigger probability (0-50400)
    /// @param _manager Address of the manager who can modify RAT parameters
    function initialize(
        IDisputeGameFactory _disputeGameFactory,
        uint256 _perTestBondAmount,
        uint256 _evidenceSubmissionPeriod,
        uint256 _minimumStakingBalance,
        uint256 _ratTriggerProbability,
        address _manager
    )
        public
        payable
        reinitializer(initVersion())
    {
        require(_perTestBondAmount < type(uint96).max, "Bond amount exceeds uint96 maximum");
        require(_perTestBondAmount <= _minimumStakingBalance, "Bond amount cannot exceed minimum staking balance");
        require(_ratTriggerProbability <= MAX_PROBABILITY, "Invalid probability");
        disputeGameFactory = _disputeGameFactory;
        perTestBondAmount = _perTestBondAmount;
        evidenceSubmissionPeriod = _evidenceSubmissionPeriod;
        minimumStakingBalance = _minimumStakingBalance;
        ratTriggerProbability = _ratTriggerProbability;
        ratManager = _manager;

        // Initialize validChallengers with a dummy element at index 0
        validChallengers.push(address(0));
    }

    /// @notice Allows challengers to stake ETH
    function stake() external payable nonReentrant {
        require(msg.value > 0, "Must stake positive amount");

        ChallengerInfo storage challenger = challengers[msg.sender];

        challenger.stakingAmount += msg.value;
        // Check if challenger meets per-test bond requirement
        if (!challenger.isValid && (challenger.stakingAmount >= perTestBondAmount)) {
            challenger.isValid = true;
            challenger.validatorIndex = uint32(validChallengers.length);
            validChallengers.push(msg.sender);
        }

        emit ChallengerStaked(msg.sender, msg.value);
    }

    /// @notice Gets challenger information
    /// @param _challenger Address of the challenger
    /// @return Challenger information
    function getChallengerInfo(address _challenger) external view returns (ChallengerInfo memory) {
        return challengers[_challenger];
    }



    /// @notice Gets number of valid challengers
    /// @return Number of valid challengers
    function getValidChallengerCount() external view returns (uint256) {
        return validChallengers.length;
    }



    /// @notice Triggers attention test (called by DisputeGameFactory)
    /// @param _gameAddress Game contract address
    /// @param _stateRoot State root to be verified
    /// @param _blockHash Block hash for validator selection
    function triggerAttentionTest(
        address _gameAddress,
        bytes32 _stateRoot,
        bytes32 _blockHash
    )
        external
        onlyDisputeGameFactory
    {
        // First check if RAT should be triggered based on probability
        if (!shouldTriggerRAT()) return;

        uint256 validChallengersLength = validChallengers.length;
        if (validChallengersLength > 1) {

            // Optimize challenger selection
            uint256 selectedIndex = validChallengersLength == 2 ? 1 :
               ((uint256(keccak256(abi.encodePacked(_blockHash, block.timestamp))) & 0xFFFF) % (validChallengersLength-1) )+1; // -1 to exclude the dummy address(0)

            address selectedChallenger = validChallengers[selectedIndex];

            ChallengerInfo storage challengerInfo = challengers[selectedChallenger];

            // Calculate bond amount and update challenger (gas-optimized)
            uint256 stakingAmount = challengerInfo.stakingAmount;
            uint256 bondAmount = stakingAmount < perTestBondAmount ? stakingAmount : perTestBondAmount;
            uint256 newStakingAmount = stakingAmount - bondAmount;

            // Bond the challenger
            challengerInfo.stakingAmount = newStakingAmount;
            challengerInfo.totalSlashedAmount += bondAmount;

            // Validate block number
            require(block.number <= type(uint64).max, "Block number too large");

            // Store attention test info
            attentionTests[_gameAddress] = AttentionInfo({
                stateRoot: _stateRoot,
                bondAmount: uint96(bondAmount),
                challengerAddress: selectedChallenger,
                l1BlockNumber: uint64(block.number),
                evidenceSubmitted: false
            });

            // Check if challenger is still valid (gas-optimized)
            bool shouldBeValid = newStakingAmount >= perTestBondAmount;
            bool currentIsValid = challengerInfo.isValid;

            if (currentIsValid != shouldBeValid) {
                if (shouldBeValid) {
                    // invalid → valid
                    challengerInfo.isValid = true;
                    _addToValidChallengers(selectedChallenger, uint32(validChallengers.length));
                } else {
                    // valid → invalid
                    challengerInfo.isValid = false;
                    _removeFromValidChallengers(selectedChallenger, challengerInfo.validatorIndex);
                }
            }

            emit AttentionTriggered(_gameAddress, selectedChallenger);
        }
    }

    /// @notice Submits correct evidence for attention test
    /// @param _gameAddress Game contract address
    /// @param _proofLV Left child state value
    /// @param _proofRV Right child state value
    function submitCorrectEvidence(
        address _gameAddress,
        bytes32 _proofLV,
        bytes32 _proofRV
    )
        external
    {
        AttentionInfo storage attentionTest = attentionTests[_gameAddress];

        // Early validation with cached values (gas optimization)
        address challengerAddress = attentionTest.challengerAddress;
        if (challengerAddress == address(0)) revert AttentionTestNotExists();
        if (challengerAddress != msg.sender) revert InvalidChallengerAddress();
        if (attentionTest.evidenceSubmitted) revert EvidenceAlreadySubmitted();

        // Time validation with overflow protection (gas optimized)
        uint256 submissionDeadline = attentionTest.l1BlockNumber + evidenceSubmissionPeriod;
        if (submissionDeadline < attentionTest.l1BlockNumber) revert("Deadline overflow");
        if (block.number >= submissionDeadline) revert EvidenceSubmissionExpired();

        // Verify proof (gas optimized - single hash operation)
        if (keccak256(abi.encodePacked(_proofLV, _proofRV)) != attentionTest.stateRoot) revert ProofVerificationFailed();

        // Cache values for gas optimization
        uint256 bond = uint256(attentionTest.bondAmount);
        attentionTest.evidenceSubmitted = true;

        // Update challenger staking amount
        ChallengerInfo storage challengerInfo = challengers[msg.sender];
        challengerInfo.stakingAmount += bond;

        // Update challenger validity (gas optimized)
        bool currentIsValid = challengerInfo.isValid;
        bool shouldBeValid = challengerInfo.stakingAmount >= perTestBondAmount;

        if (!currentIsValid && shouldBeValid) {
            // invalid → valid
            challengerInfo.isValid = true;
            _addToValidChallengers(msg.sender, uint32(validChallengers.length));
        } else if (currentIsValid && !shouldBeValid) {
            // valid → invalid
            challengerInfo.isValid = false;
            _removeFromValidChallengers(msg.sender, challengerInfo.validatorIndex);
        }

        emit CorrectEvidenceSubmitted(
            _gameAddress,
            challengerAddress,
            bond
        );
    }

    /// @notice Called when a claim is resolved in FaultDisputeGame
    /// @param _claimant Address receiving the bond refund
    function resolveClaim(address _claimant) external {
        // Early validation with caching
        address challengerAddress = attentionTests[msg.sender].challengerAddress;
        if (challengerAddress != address(0) && challengerAddress == _claimant) {
            // bool evidenceSubmitted = attentionTests[msg.sender].evidenceSubmitted;
            if (!attentionTests[msg.sender].evidenceSubmitted) {
                AttentionInfo storage attentionTest = attentionTests[msg.sender];

                // Mark evidence as submitted and refund bond amount
                attentionTest.evidenceSubmitted = true;
                uint256 bond = uint256(attentionTest.bondAmount);

                ChallengerInfo storage challengerInfo = challengers[_claimant];
                challengerInfo.stakingAmount += bond;

                // Update challenger validity
                bool currentIsValid = challengerInfo.isValid;
                bool shouldBeValid = challengerInfo.stakingAmount >= perTestBondAmount;

                if (!currentIsValid && shouldBeValid) {
                    // invalid → valid
                    challengerInfo.isValid = true;
                    _addToValidChallengers(_claimant, uint32(validChallengers.length));
                } else if (currentIsValid && !shouldBeValid) {
                    // valid → invalid
                    challengerInfo.isValid = false;
                    _removeFromValidChallengers(_claimant, challengerInfo.validatorIndex);
                }

                emit BondRefunded(msg.sender, challengerAddress, bond);
            }
        }
    }

    /// @notice Sets the per-test bond amount (only proxy admin owner)
    /// @param _amount New bond amount
    function setPerTestBondAmount(uint256 _amount) external {
        _assertOnlyProxyAdminOwner();
        require(_amount > 0, "Bond amount must be positive");
        require(_amount <= type(uint96).max, "Bond amount exceeds uint96 maximum");
        require(_amount <= minimumStakingBalance, "Bond amount cannot exceed minimum staking balance");
        perTestBondAmount = _amount;
    }

    /// @notice Sets the evidence submission period (only proxy admin owner)
    /// @param _period New submission period in blocks
    function setEvidenceSubmissionPeriod(uint256 _period) external {
        _assertOnlyProxyAdminOwner();
        require(_period > 0, "Period must be positive");
        require(_period <= 50400, "Period too long");
        evidenceSubmissionPeriod = _period;
    }

    /// @notice Sets the minimum staking balance (only proxy admin owner)
    /// @param _balance New minimum staking balance
    function setMinimumStakingBalance(uint256 _balance) external {
        _assertOnlyProxyAdminOwner();
        require(_balance > 0, "Balance must be positive");
        require(_balance <= 1000 ether, "Balance too large");
        minimumStakingBalance = _balance;
    }

    /// @notice Sets the RAT trigger probability (only manager)
    /// @param _probability New trigger probability (0-50400)
    function setRatTriggerProbability(uint256 _probability) external onlyRatManager {
        require(_probability <= MAX_PROBABILITY, "Invalid probability");
        ratTriggerProbability = _probability;
    }

    /// @notice Check if RAT should trigger based on probability
    function shouldTriggerRAT() internal view returns (bool) {
        uint256 prob = ratTriggerProbability;
        if (prob == 0) return false;
        if (prob >= MAX_PROBABILITY) return true;
        return uint256(blockhash(block.number - 1)) % MAX_PROBABILITY < prob;
    }



    /// @notice Internal function to add challenger to valid list
    /// @param _challenger Address of the challenger
    /// @param _index Index to assign to the challenger
    function _addToValidChallengers(address _challenger, uint32 _index) internal {
        challengers[_challenger].validatorIndex = _index;
        validChallengers.push(_challenger);
    }

    /// @notice Internal function to remove challenger from valid list
    /// @param _challenger Address of the challenger
    /// @param _index Index of the challenger in validChallengers array
    function _removeFromValidChallengers(address _challenger, uint256 _index) internal {
        if (_index > 0 && _index < validChallengers.length && validChallengers[_index] == _challenger) {
            // Replace with last element and pop (gas-optimized)
            uint256 lastIndex = validChallengers.length - 1;
            if (_index != lastIndex) {
                address lastChallenger = validChallengers[lastIndex];
                validChallengers[_index] = lastChallenger;
                challengers[lastChallenger].validatorIndex = uint32(_index);
            }
            validChallengers.pop();
        }
    }


}