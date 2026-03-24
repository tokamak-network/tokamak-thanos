// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {FailLogics} from "./FailLogics.sol";

contract CommitReveal2 is FailLogics {
    constructor(
        uint256 activationThreshold,
        uint256 flatFee,
        string memory name,
        string memory version,
        uint256 offChainSubmissionPeriod,
        uint256 requestOrSubmitOrFailDecisionPeriod,
        uint256 onChainSubmissionPeriod,
        uint256 offChainSubmissionPeriodPerOperator,
        uint256 onChainSubmissionPeriodPerOperator
    ) payable FailLogics(name, version) {
        require(msg.value >= activationThreshold);
        s_depositAmount[msg.sender] = msg.value;
        s_activationThreshold = activationThreshold;
        s_flatFee = flatFee;
        s_offChainSubmissionPeriod = offChainSubmissionPeriod;
        s_requestOrSubmitOrFailDecisionPeriod = requestOrSubmitOrFailDecisionPeriod;
        s_onChainSubmissionPeriod = onChainSubmissionPeriod;
        s_offChainSubmissionPeriodPerOperator = offChainSubmissionPeriodPerOperator;
        s_onChainSubmissionPeriodPerOperator = onChainSubmissionPeriodPerOperator;
        s_isInProcess = COMPLETED;
    }

    function proposeEconomicParameters(uint256 activationThreshold, uint256 flatFee) external onlyOwner {
        assembly ("memory-safe") {
            sstore(s_pendingActivationThreshold.slot, activationThreshold)
            sstore(s_pendingFlatFee.slot, flatFee)
            let effectiveTimestamp := add(timestamp(), SET_DELAY_TIME)
            sstore(s_economicParamsEffectiveTimestamp.slot, effectiveTimestamp)
            mstore(0x00, activationThreshold)
            mstore(0x20, flatFee)
            mstore(0x40, effectiveTimestamp)
            log1(0x00, 0x60, 0xdcf23dfc5bc14859d1943fd156abd0fb732347e70c61c56215bbd728307234e2) // EconomicParametersProposed(uint256 activationThreshold, uint256 flatFee, uint256 effectiveTimestamp)
        }
    }

    function executeSetEconomicParameters() external notInProcess {
        assembly ("memory-safe") {
            if lt(timestamp(), sload(s_economicParamsEffectiveTimestamp.slot)) {
                mstore(0, 0x085de625) // selector for TooEarly()
                revert(0x1c, 0x04)
            }
            let activationThreshold := sload(s_pendingActivationThreshold.slot)
            let flatFee := sload(s_pendingFlatFee.slot)
            sstore(s_activationThreshold.slot, activationThreshold)
            sstore(s_flatFee.slot, flatFee)
            sstore(s_economicParamsEffectiveTimestamp.slot, 0)
            mstore(0x00, activationThreshold)
            mstore(0x20, flatFee)
            log1(0x00, 0x40, 0x08f0774e7eb69e2d6a7cf2192cbf9c6f519a40bcfa16ff60d3f18496585e46dc) // EconomicParametersSet
        }
    }

    function setPeriods(
        uint256 offChainSubmissionPeriod,
        uint256 requestOrSubmitOrFailDecisionPeriod,
        uint256 onChainSubmissionPeriod,
        uint256 offChainSubmissionPeriodPerOperator,
        uint256 onChainSubmissionPeriodPerOperator
    ) external onlyOwner notInProcess {
        assembly ("memory-safe") {
            sstore(s_offChainSubmissionPeriod.slot, offChainSubmissionPeriod)
            sstore(s_requestOrSubmitOrFailDecisionPeriod.slot, requestOrSubmitOrFailDecisionPeriod)
            sstore(s_onChainSubmissionPeriod.slot, onChainSubmissionPeriod)
            sstore(s_offChainSubmissionPeriodPerOperator.slot, offChainSubmissionPeriodPerOperator)
            sstore(s_onChainSubmissionPeriodPerOperator.slot, onChainSubmissionPeriodPerOperator)
            mstore(0x00, offChainSubmissionPeriod)
            mstore(0x20, requestOrSubmitOrFailDecisionPeriod)
            mstore(0x40, onChainSubmissionPeriod)
            mstore(0x60, offChainSubmissionPeriodPerOperator)
            mstore(0x80, onChainSubmissionPeriodPerOperator)
            log1(0x00, 0xa0, 0xe0fd8eabd2cc23ea87b43a00ac588c61789ad28d3edfeb76613f623fa1f6bd08) // event PeriodsSet(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator)
        }
    }

    function proposeGasParameters(
        uint128 gasUsedMerkleRootSubAndGenRandNumA,
        uint128 gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead,
        uint256 maxCallbackGasLimit,
        uint48 getL1UpperBoundGasUsedWhenCalldataSize4,
        uint48 failToRequestCvOrSubmitMerkleRootGasUsed,
        uint48 failToSubmitMerkleRootAfterDisputeGasUsed,
        uint48 failToRequestSOrGenerateRandomNumberGasUsed,
        uint48 failToSubmitSGasUsed,
        uint32 failToSubmitCoGasUsedBaseA,
        uint32 failToSubmitCvGasUsedBaseA,
        uint32 failToSubmitGasUsedBaseB,
        uint32 perOperatorIncreaseGasUsedA,
        uint32 perOperatorIncreaseGasUsedB,
        uint32 perAdditionalDidntSubmitGasUsedA,
        uint32 perAdditionalDidntSubmitGasUsedB,
        uint32 perRequestedIncreaseGasUsed
    ) external onlyOwner {
        assembly ("memory-safe") {
            // Pack as: low 128 bits = A, high 128 bits = B
            sstore(
                s_pendingGasUsedMerkleRootSubAndGenRandNumA.slot,
                or(gasUsedMerkleRootSubAndGenRandNumA, shl(128, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead))
            )
            sstore(s_pendingMaxCallbackGasLimit.slot, maxCallbackGasLimit)

            sstore(
                s_pendingGetL1UpperBoundGasUsedWhenCalldataSize4.slot,
                or(
                    getL1UpperBoundGasUsedWhenCalldataSize4,
                    or(
                        shl(FAILTOREQUESTSUBMITCV_OR_SUBMITMEKRLEROOT_OFFSET, failToRequestCvOrSubmitMerkleRootGasUsed),
                        or(
                            shl(FAILTOSUBMITMERKLEROOTAFTERDISPUTE_OFFSET, failToSubmitMerkleRootAfterDisputeGasUsed),
                            or(
                                shl(
                                    FAILTOREQUESTS_OR_GENERATERANDOMNUMBER_OFFSET,
                                    failToRequestSOrGenerateRandomNumberGasUsed
                                ),
                                shl(FAILTOSUBMITS_OFFSET, failToSubmitSGasUsed)
                            )
                        )
                    )
                )
            )
            sstore(
                s_pendingFailToSubmitCoGasUsedBaseA.slot,
                or(
                    failToSubmitCoGasUsedBaseA,
                    or(
                        shl(FAILTOSUBMITCVGASUSEDBASEA_OFFSET, failToSubmitCvGasUsedBaseA),
                        or(
                            shl(FAILTOSUBMITGASUSEDBASEB_OFFSET, failToSubmitGasUsedBaseB),
                            or(
                                shl(PEROPERATORINCREASEGASUSEDA_OFFSET, perOperatorIncreaseGasUsedA),
                                or(
                                    shl(PEROPERATORINCREASEGASUSEDB_OFFSET, perOperatorIncreaseGasUsedB),
                                    or(
                                        shl(PERADDITIONALDIDNTSUBMITGASUSEDA_OFFSET, perAdditionalDidntSubmitGasUsedA),
                                        or(
                                            shl(
                                                PERADDITIONALDIDNTSUBMITGASUSEDB_OFFSET,
                                                perAdditionalDidntSubmitGasUsedB
                                            ),
                                            shl(PERREQUESTEDINCREASEGASUSED_OFFSET, perRequestedIncreaseGasUsed)
                                        )
                                    )
                                )
                            )
                        )
                    )
                )
            )
            let effectiveTimestamp := add(timestamp(), SET_DELAY_TIME)
            sstore(s_gasParamsEffectiveTimestamp.slot, effectiveTimestamp)

            mstore(0x00, gasUsedMerkleRootSubAndGenRandNumA)
            mstore(0x20, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead)
            mstore(0x40, maxCallbackGasLimit)
            mstore(0x60, getL1UpperBoundGasUsedWhenCalldataSize4)
            mstore(0x80, failToRequestCvOrSubmitMerkleRootGasUsed)
            mstore(0xa0, failToSubmitMerkleRootAfterDisputeGasUsed)
            mstore(0xc0, failToRequestSOrGenerateRandomNumberGasUsed)
            mstore(0xe0, failToSubmitSGasUsed)
            mstore(0x100, failToSubmitCoGasUsedBaseA)
            mstore(0x120, failToSubmitCvGasUsedBaseA)
            mstore(0x140, failToSubmitGasUsedBaseB)
            mstore(0x160, perOperatorIncreaseGasUsedA)
            mstore(0x180, perOperatorIncreaseGasUsedB)
            mstore(0x1a0, perAdditionalDidntSubmitGasUsedA)
            mstore(0x1c0, perAdditionalDidntSubmitGasUsedB)
            mstore(0x1e0, perRequestedIncreaseGasUsed)
            mstore(0x200, effectiveTimestamp)
            log1(0x00, 0x220, 0xac29dedddb8466e143ff09a21b0181b73354eae633cc2787fb6dd4c3b50dfbe2) // event GasParametersProposed(...)
        }
    }

    function executeSetGasParameters() external notInProcess {
        assembly ("memory-safe") {
            if lt(timestamp(), sload(s_gasParamsEffectiveTimestamp.slot)) {
                mstore(0, 0x085de625) // selector for TooEarly()
                revert(0x1c, 0x04)
            }
            let packedData := sload(s_pendingGasUsedMerkleRootSubAndGenRandNumA.slot)
            sstore(s_gasUsedMerkleRootSubAndGenRandNumA.slot, packedData)
            mstore(0x00, and(packedData, GASUSED_MERKLEROOTSUB_GENRANDNUM_MASK))
            mstore(0x20, shr(128, packedData))
            packedData := sload(s_pendingMaxCallbackGasLimit.slot)
            sstore(s_maxCallbackGasLimit.slot, packedData)
            mstore(0x40, packedData)

            packedData := sload(s_pendingGetL1UpperBoundGasUsedWhenCalldataSize4.slot)
            sstore(s_getL1UpperBoundGasUsedWhenCalldataSize4.slot, packedData)
            mstore(0x60, and(packedData, FAILTOSUBMIT_MASK))
            mstore(0x80, and(shr(FAILTOREQUESTSUBMITCV_OR_SUBMITMEKRLEROOT_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0xa0, and(shr(FAILTOSUBMITMERKLEROOTAFTERDISPUTE_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0xc0, and(shr(FAILTOREQUESTS_OR_GENERATERANDOMNUMBER_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0xe0, and(shr(FAILTOSUBMITS_OFFSET, packedData), FAILTOSUBMIT_MASK))

            packedData := sload(s_pendingFailToSubmitCoGasUsedBaseA.slot)
            sstore(s_failToSubmitCoGasUsedBaseA.slot, packedData)
            mstore(0x100, and(packedData, FAILTOSUBMIT_MASK))
            mstore(0x120, and(shr(FAILTOSUBMITCVGASUSEDBASEA_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0x140, and(shr(FAILTOSUBMITGASUSEDBASEB_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0x160, and(shr(PEROPERATORINCREASEGASUSEDA_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0x180, and(shr(PEROPERATORINCREASEGASUSEDB_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0x1a0, and(shr(PERADDITIONALDIDNTSUBMITGASUSEDA_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0x1c0, and(shr(PERADDITIONALDIDNTSUBMITGASUSEDB_OFFSET, packedData), FAILTOSUBMIT_MASK))
            mstore(0x1e0, and(shr(PERREQUESTEDINCREASEGASUSED_OFFSET, packedData), FAILTOSUBMIT_MASK))
            // clear effective timestamp after execution
            sstore(s_gasParamsEffectiveTimestamp.slot, 0)
            log1(0x00, 0x200, 0x8d09171105499771f96d6d39dcdda061a70fd18e5eafd65881c2158c55f94e1d) // event GasParametersSet(...)
        }
    }

    function estimateRequestPrice(uint32 callbackGasLimit, uint256 gasPrice) external view returns (uint256) {
        return _calculateRequestPrice(callbackGasLimit, gasPrice, s_activatedOperators.length);
    }

    function estimateRequestPriceWithNumOfOperators(uint32 callbackGasLimit, uint256 gasPrice, uint256 numOfOperators)
        external
        view
        returns (uint256)
    {
        return _calculateRequestPrice(callbackGasLimit, gasPrice, numOfOperators);
    }

    function requestRandomNumber(uint32 callbackGasLimit) external payable virtual returns (uint256) {
        uint256 activatedOperatorsLength = s_activatedOperators.length;
        // ** check if the fee amount is enough
        require(
            msg.value >= _calculateRequestPrice(callbackGasLimit, tx.gasprice, activatedOperatorsLength),
            InsufficientAmount()
        );
        assembly ("memory-safe") {
            // ** check if the callbackGasLimit is within the limit
            if gt(callbackGasLimit, sload(s_maxCallbackGasLimit.slot)) {
                mstore(0, 0x1cf7ab79) // selector for ExceedCallbackGasLimit()
                revert(0x1c, 0x04)
            }
            // ** check if there are enough activated operators
            if lt(activatedOperatorsLength, 2) {
                mstore(0, 0x77599fd9) // selector for NotEnoughActivatedOperators()
                revert(0x1c, 0x04)
            }
            // ** check if the leader has enough deposit
            mstore(0x00, sload(_OWNER_SLOT))
            mstore(0x20, s_depositAmount.slot)
            if lt(sload(keccak256(0x00, 0x40)), sload(s_activationThreshold.slot)) {
                mstore(0, 0xc0013a5a) // selector for LeaderLowDeposit()
                revert(0x1c, 0x04)
            }
            let newRound := sload(s_requestCount.slot)
            sstore(s_requestCount.slot, add(newRound, 1)) // update the request count
            if gt(sub(newRound, sload(s_currentRound.slot)), 2000) {
                mstore(0, 0x02cd147b) // selector for TooManyRequestsQueued()
                revert(0x1c, 0x04)
            }

            // ** set the round bit
            // calculate the storage slot corresponding to the round
            // wordPos = round >> 8
            mstore(0, shr(8, newRound))
            mstore(0x20, s_roundBitmap.slot)
            // the slot of self[wordPos] is keccak256(abi.encode(wordPos, self.slot))
            let slot := keccak256(0, 0x40)
            // mask = 1 << bitPos = 1 << (round & 0xff)
            // self[wordPos] |= mask
            sstore(slot, or(sload(slot), shl(and(newRound, 0xff), 1)))
            let startTime
            // ** check if the current round is completed
            // ** if the current round is completed, start a new round
            let currentState := sload(s_isInProcess.slot)
            if eq(currentState, HALTED) {
                mstore(0, 0x2caa910c) // selector for CannotRequestWhenHalted()
                revert(0x1c, 0x04)
            }
            if eq(currentState, COMPLETED) {
                startTime := timestamp()
                sstore(s_currentRound.slot, newRound)
                sstore(s_isInProcess.slot, IN_PROGRESS)
                mstore(0, newRound)
                mstore(0x20, 0) // trialNum is 0 for the first trial
                mstore(0x40, IN_PROGRESS)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
            // *** store the request info
            mstore(0x00, newRound)
            mstore(0x20, s_requestInfo.slot)
            let requestInfoSlot := keccak256(0x00, 0x40)
            sstore(requestInfoSlot, or(shl(96, caller()), callbackGasLimit))
            sstore(add(requestInfoSlot, 1), startTime)
            sstore(add(requestInfoSlot, 2), callvalue())
            return(0x00, 0x20)
        }
    }

    function _calculateRequestPrice(uint32 callbackGasLimit, uint256 gasPrice, uint256 numOfOperators)
        internal
        view
        virtual
        returns (uint256 requestFee)
    {
        assembly ("memory-safe") {
            let gasUsedMerkleRootSubAndGenRandNum := sload(s_gasUsedMerkleRootSubAndGenRandNumA.slot)
            requestFee :=
                add(
                    mul(
                        gasPrice,
                        add(
                            callbackGasLimit,
                            add(
                                mul(
                                    and(gasUsedMerkleRootSubAndGenRandNum, GASUSED_MERKLEROOTSUB_GENRANDNUM_MASK),
                                    numOfOperators
                                ),
                                shr(128, gasUsedMerkleRootSubAndGenRandNum) // gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead
                            )
                        )
                    ),
                    sload(s_flatFee.slot)
                )
        }
    }

    function submitMerkleRoot(bytes32 merkleRoot) external onlyOwner {
        assembly ("memory-safe") {
            // * get trialNum
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            mstore(0x20, sload(keccak256(0x40, 0x40))) // trialNum
            // * get merkleRootSubmittedTimestamp
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x40, keccak256(0x40, 0x40))
            let merkleRootSubmittedTimestampSlot := keccak256(0x20, 0x40)
            if gt(sload(merkleRootSubmittedTimestampSlot), 0) {
                mstore(0, 0xa34402b2) // selector for MerkleRootAlreadySubmitted()
                revert(0x1c, 0x04)
            }
            sstore(s_merkleRoot.slot, merkleRoot)
            sstore(merkleRootSubmittedTimestampSlot, timestamp())
            // * emit event MerkleRootSubmitted
            mstore(0x00, curRound)
            // 0x20 already has trialNum
            mstore(0x40, merkleRoot)
            log1(0x00, 0x60, 0x45b19880b523c6750f7f39fca8d77d51101b315495adc482994a4fa2a8294466) // emit event MerkleRootSubmitted(uint256 round, uint256 trialNum, bytes32 merkleRoot)
        }
    }

    function generateRandomNumber(
        SecretAndSigRS[] calldata secretSigRSs,
        uint256, // packedVs
        uint256 packedRevealOrders
    ) external virtual {
        bytes32 domainSeparator = _domainSeparatorV4();
        assembly ("memory-safe") {
            let activatedOperatorsLength := sload(s_activatedOperators.slot)
            // ** check if all secrets are submitted
            if iszero(eq(activatedOperatorsLength, secretSigRSs.length)) {
                mstore(0, 0xe0767fa4) // selector for InvalidSecretLength()
                revert(0x1c, 0x04)
            }
            // ** initialize cos and cvs arrays memory, without length data
            let activatedOperatorsLengthInBytes := shl(5, activatedOperatorsLength)
            let cos := mload(0x40)
            let cvs := add(cos, activatedOperatorsLengthInBytes)
            let secrets := add(cvs, activatedOperatorsLengthInBytes)
            mstore(0x40, add(secrets, activatedOperatorsLengthInBytes)) // update the free memory pointer

            // ** get cos and cvs
            for { let i } lt(i, activatedOperatorsLengthInBytes) { i := add(i, 0x20) } {
                let secretMemP := add(secrets, i)
                mstore(secretMemP, calldataload(add(secretSigRSs.offset, mul(i, 3)))) // secret
                let cosMemP := add(cos, i)
                mstore(cosMemP, keccak256(secretMemP, 0x20))
                mstore(add(cvs, i), keccak256(cosMemP, 0x20))
            }
            // ** verify reveal order
            let index := and(packedRevealOrders, 0xff) // first reveal index
            let revealBitmap := shl(index, 1)
            mstore(0x00, keccak256(cos, activatedOperatorsLengthInBytes)) // rv
            mstore(0x20, mload(add(cvs, shl(5, index))))
            let before := keccak256(0x00, 0x40)
            // revealOrdersOffset = 0x44
            for { let i := 1 } lt(i, activatedOperatorsLength) { i := add(i, 1) } {
                index := and(calldataload(sub(0x44, i)), 0xff)
                revealBitmap := or(revealBitmap, shl(index, 1))
                mstore(0x20, mload(add(cvs, shl(5, index))))
                let after := keccak256(0x00, 0x40)
                if lt(before, after) {
                    mstore(0, 0x24f1948e) // selector for RevealNotInDescendingOrder()
                    revert(0x1c, 0x04)
                }
                before := after
            }
            if iszero(eq(revealBitmap, sub(shl(activatedOperatorsLength, 1), 1))) {
                mstore(0, 0x06efcba4) // selector for RevealOrderHasDuplicates()
                revert(0x1c, 0x04)
            }
            // ** Create Merkle Root and verify it
            let hashCountInBytes := sub(activatedOperatorsLengthInBytes, 0x20)
            let fmp := mload(0x40) // used to store the hashes
            let cvsPosInBytes
            let hashPosInBytes
            for { let i } lt(i, hashCountInBytes) { i := add(i, 0x20) } {
                switch lt(cvsPosInBytes, activatedOperatorsLengthInBytes)
                case 1 {
                    mstore(0x00, mload(add(cvs, cvsPosInBytes)))
                    cvsPosInBytes := add(cvsPosInBytes, 0x20)
                }
                default {
                    mstore(0x00, mload(add(fmp, hashPosInBytes)))
                    hashPosInBytes := add(hashPosInBytes, 0x20)
                }
                switch lt(cvsPosInBytes, activatedOperatorsLengthInBytes)
                case 1 {
                    mstore(0x20, mload(add(cvs, cvsPosInBytes)))
                    cvsPosInBytes := add(cvsPosInBytes, 0x20)
                }
                default {
                    mstore(0x20, mload(add(fmp, hashPosInBytes)))
                    hashPosInBytes := add(hashPosInBytes, 0x20)
                }
                mstore(add(fmp, i), keccak256(0x00, 0x40))
            }
            // ** check if the merkle root is submitted
            let round := sload(s_currentRound.slot)
            mstore(0x20, round)
            mstore(0x40, s_trialNum.slot)
            let trialNum := sload(keccak256(0x20, 0x40))
            mstore(0x00, trialNum)
            mstore(0x40, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x20, 0x40))
            if iszero(sload(keccak256(0x00, 0x40))) {
                mstore(0, 0x8e56b845) // selector for MerkleRootNotSubmitted()
                revert(0x1c, 0x04)
            }
            // ** verify the merkle root
            if iszero(eq(mload(add(fmp, sub(hashCountInBytes, 0x20))), sload(s_merkleRoot.slot))) {
                mstore(0, 0x624dc351) // selector for MerkleVerificationFailed()
                revert(0x1c, 0x04)
            }
            // ** verify signatures
            mstore(fmp, MESSAGE_TYPEHASH_DIRECT) // typehash, overwrite the previous value, which is not used anymore
            mstore(add(fmp, 0x20), round)
            mstore(add(fmp, 0x40), trialNum)
            mstore(add(fmp, 0x80), hex"1901") // prefix and version
            mstore(add(fmp, 0x82), domainSeparator)
            for { let i } lt(i, activatedOperatorsLengthInBytes) { i := add(i, 0x20) } {
                // signature malleability prevention
                let rSOffset := add(secretSigRSs.offset, add(mul(i, 3), 0x20))
                let s := calldataload(add(rSOffset, 0x20))
                if gt(s, SECP256K1_CURVE_ORDER) {
                    mstore(0, 0xbf4bf5b8) // selector for InvalidSignatureS()
                    revert(0x1c, 0x04)
                }
                mstore(add(fmp, 0x60), mload(add(cvs, i))) // cv
                mstore(add(fmp, 0xa2), keccak256(fmp, 0x80)) // structHash
                mstore(0x00, keccak256(add(fmp, 0x80), 0x42)) // digest hash
                mstore(0x20, and(calldataload(sub(0x24, shr(5, i))), 0xff)) // v, 0x24: packedVsOffset
                mstore(0x40, calldataload(rSOffset)) // r
                mstore(0x60, s) // s
                let operatorAddress := mload(staticcall(gas(), 1, 0x00, 0x80, 0x01, 0x20))
                // `returndatasize()` will be `0x20` upon success, and `0x00` otherwise.
                if iszero(returndatasize()) {
                    mstore(0x00, 0x8baa579f) // selector for InvalidSignature()
                    revert(0x1c, 0x04)
                }
                mstore(0x00, operatorAddress)
                mstore(0x20, s_activatedOperatorIndex1Based.slot)
                if iszero(sload(keccak256(0x00, 0x40))) {
                    mstore(0x00, 0x1b256530) // selector for NotActivatedOperator()
                    revert(0x1c, 0x04)
                }
            }

            // ** create random number
            let randomNumber := keccak256(secrets, activatedOperatorsLengthInBytes)
            let nextRound := add(round, 1)
            let requestCount := sload(s_requestCount.slot)
            switch eq(nextRound, requestCount)
            case 1 {
                // there is no next round
                if eq(sload(s_isInProcess.slot), COMPLETED) {
                    mstore(0x00, 0x195332a5) // selector for AlreadyCompleted()
                    revert(0x1c, 0x04)
                }
                sstore(s_isInProcess.slot, COMPLETED)
                mstore(0x00, round)
                mstore(0x20, trialNum)
                mstore(0x40, COMPLETED)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
            default {
                // get next round
                // https://github.com/Uniswap/v4-core/blob/59d3ecf53afa9264a16bba0e38f4c5d2231f80bc/src/libraries/BitMath.sol#L31
                function leastSignificantBit(x) -> r {
                    x := and(x, sub(0, x))
                    r :=
                        shl(
                            5,
                            shr(
                                252,
                                shl(
                                    shl(
                                        2,
                                        shr(250, mul(x, 0xb6db6db6ddddddddd34d34d349249249210842108c6318c639ce739cffffffff))
                                    ),
                                    0x8040405543005266443200005020610674053026020000107506200176117077
                                )
                            )
                        )
                    r :=
                        or(
                            r,
                            byte(
                                and(div(0xd76453e0, shr(r, x)), 0x1f),
                                0x001f0d1e100c1d070f090b19131c1706010e11080a1a141802121b1503160405
                            )
                        )
                }
                function nextRequestedRound(_round) -> _next, _requested {
                    let wordPos := shr(8, _round)
                    let bitPos := and(_round, 0xff)
                    let mask := not(sub(shl(bitPos, 1), 1))
                    mstore(0x00, wordPos)
                    mstore(0x20, s_roundBitmap.slot)
                    let masked := and(sload(keccak256(0x00, 0x40)), mask)
                    _requested := gt(masked, 0)
                    switch _requested
                    case 1 { _next := sub(add(_round, leastSignificantBit(masked)), bitPos) }
                    default { _next := sub(add(_round, 255), bitPos) }
                }
                let requested
                for { let i } lt(i, 10) { i := add(i, 1) } {
                    nextRound, requested := nextRequestedRound(nextRound)
                    if requested {
                        mstore(0x00, nextRound) // round
                        mstore(0x20, s_requestInfo.slot)
                        sstore(add(keccak256(0x00, 0x40), 1), timestamp()) // startTime
                        sstore(s_currentRound.slot, nextRound)
                        mstore(0x20, 0) // trialNum is 0 for the first trial
                        mstore(0x40, IN_PROGRESS)
                        log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
                        break
                    }
                    if iszero(lt(nextRound, requestCount)) {
                        if eq(sload(s_isInProcess.slot), COMPLETED) {
                            mstore(0x00, 0x195332a5) // selector for AlreadyCompleted()
                            revert(0x1c, 0x04)
                        }
                        sstore(s_isInProcess.slot, COMPLETED)
                        let lastRound := sub(requestCount, 1)
                        sstore(s_currentRound.slot, lastRound)
                        mstore(0x00, lastRound)
                        mstore(0x20, 0) // trialNum is 0 for the first trial
                        mstore(0x40, COMPLETED)
                        log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
                        break
                    }
                    nextRound := add(nextRound, 1)
                }
            }
            // ** reward the flatFee to last revealer
            // ** reward the leaderNode (requestFee - flatFee) for submitMerkleRoot and generateRandomNumber
            mstore(0x00, s_activatedOperators.slot)
            mstore(
                0x00,
                sload(
                    add(
                        keccak256(0x00, 0x20), // s_activatedOperators first data slot
                        and(calldataload(sub(0x44, sub(activatedOperatorsLength, 1))), 0xff) // last revealer index, 0x44: revealOrdersOffset
                    )
                )
            ) // last revealer address
            mstore(0x20, s_depositAmount.slot)
            let depositSlot := keccak256(0x00, 0x40) // last revealer
            let flatFee := sload(s_flatFee.slot)
            sstore(depositSlot, add(sload(depositSlot), flatFee))
            // reward sload(add(currentRequestInfoSlot, 2)) - flatFee to the leader
            mstore(0x00, sload(_OWNER_SLOT))
            depositSlot := keccak256(0x00, 0x40) // leader

            mstore(0x00, 0x00fc98b8) // rawFulfillRandomNumber(uint256,uint256) selector
            mstore(0x20, round)
            mstore(0x40, s_requestInfo.slot)
            let currentRequestInfoSlot := keccak256(0x20, 0x40)
            // * update the leader's deposit
            sstore(depositSlot, add(sload(depositSlot), sub(sload(add(currentRequestInfoSlot, 2)), flatFee)))
            mstore(0x40, randomNumber)

            let g := gas()
            // Compute g -= GAS_FOR_CALL_EXACT_CHECK and check for underflow
            // The gas actually passed to the callee is min(gasAmount, 63//64*gas available)
            // We want to ensure that we revert if gasAmount > 63//64*gas available
            // as we do not want to provide them with less, however that check itself costs
            // gas. GAS_FOR_CALL_EXACT_CHECK ensures we have at least enough gas to be able to revert
            // if gasAmount > 63//64*gas available.
            if lt(g, GAS_FOR_CALL_EXACT_CHECK) { revert(0, 0) }
            g := sub(g, GAS_FOR_CALL_EXACT_CHECK)
            // if g - g//64 <= gas
            // we subtract g//64 because of EIP-150
            g := sub(g, div(g, 64))
            let consumerAndCallbackGasLimitPacked := sload(currentRequestInfoSlot)
            let callbackGasLimit := and(consumerAndCallbackGasLimitPacked, 0xffffffff)
            if iszero(gt(sub(g, div(g, 64)), callbackGasLimit)) { revert(0, 0) }
            // solidity calls check that a contract actually exists at the destination, so we do the same
            let consumer := shr(96, consumerAndCallbackGasLimitPacked)
            if gt(extcodesize(consumer), 0) {
                // call and return whether we succeeded. ignore return data
                // call(gas, addr, value, argsOffset,argsLength,retOffset,retLength)
                pop(call(callbackGasLimit, consumer, 0, 0x1c, 0x44, 0, 0))
            }
        }
    }

    function refund(uint256 round) external {
        assembly ("memory-safe") {
            // ** check if the contract is halted
            if iszero(eq(sload(s_isInProcess.slot), HALTED)) {
                mstore(0, 0x78b19eb2) // selector for NotHalted()
                revert(0x1c, 0x04)
            }
            // ** check if the round is valid
            if iszero(lt(round, sload(s_requestCount.slot))) {
                mstore(0, 0x905deff6) // selector for NonExistentRound()
                revert(0x1c, 0x04)
            }
            if lt(round, sload(s_currentRound.slot)) {
                mstore(0, 0x5cafea8c) // selector for RoundAlreadyProcessed()
                revert(0x1c, 0x04)
            }
            mstore(0x00, round)
            mstore(0x20, s_requestInfo.slot)
            let consumerSlot := keccak256(0x00, 0x40)
            // ** check if the caller is the consumer
            if iszero(eq(shr(96, sload(consumerSlot)), caller())) {
                mstore(0, 0x8c7dc13d) // selector for NotConsumer()
                revert(0x1c, 0x04)
            }

            // ** flip the roundBitmap 1 -> 0
            // calculate the storage slot corresponding to the round
            // wordPos = round >> 8
            mstore(0x00, shr(8, round))
            mstore(0x20, s_roundBitmap.slot)
            // the slot of self[wordPos] is keccak256(abi.encode(wordPos, self.slot))
            let slot := keccak256(0, 0x40)
            // mask = 1 << bitPos = 1 << (round & 0xff)
            // self[wordPos] ^= mask
            sstore(slot, xor(sload(slot), shl(and(round, 0xff), 1)))

            // ** refund
            slot := add(consumerSlot, 2) // cost
            let cost := sload(slot)
            if iszero(cost) {
                mstore(0, 0xa85e6f1a) // selector for AlreadyRefunded()
                revert(0x1c, 0x04)
            }
            sstore(slot, 0)
            // Transfer the ETH and check if it succeeded or not.
            if iszero(call(gas(), caller(), cost, 0x00, 0x00, 0x00, 0x00)) {
                mstore(0x00, 0xb12d13eb) // `ETHTransferFailed()`.
                revert(0x1c, 0x04)
            }
        }
    }

    function resume() external payable onlyOwner {
        assembly ("memory-safe") {
            if iszero(eq(sload(s_isInProcess.slot), HALTED)) {
                mstore(0, 0x78b19eb2) // selector for NotHalted()
                revert(0x1c, 0x04)
            }
            if lt(sload(s_activatedOperators.slot), 2) {
                mstore(0, 0x77599fd9) // selector for NotEnoughActivatedOperators()
                revert(0x1c, 0x04)
            }
            mstore(0x00, sload(_OWNER_SLOT))
            mstore(0x20, s_depositAmount.slot)
            let ownerDepositSlot := keccak256(0x00, 0x40)
            let ownerDepositAmount := sload(ownerDepositSlot)
            if gt(callvalue(), 0) {
                ownerDepositAmount := add(ownerDepositAmount, callvalue())
                sstore(ownerDepositSlot, ownerDepositAmount)
            }
            if lt(ownerDepositAmount, sload(s_activationThreshold.slot)) {
                mstore(0, 0xc0013a5a) // selector for LeaderLowDeposit()
                revert(0x1c, 0x04)
            }
            let nextRound := sload(s_currentRound.slot)
            let requestCountMinusOne := sub(sload(s_requestCount.slot), 1)
            let curRound := nextRound
            let requested

            // get next round
            // https://github.com/Uniswap/v4-core/blob/59d3ecf53afa9264a16bba0e38f4c5d2231f80bc/src/libraries/BitMath.sol#L31
            function leastSignificantBit(x) -> r {
                x := and(x, sub(0, x))
                r :=
                    shl(
                        5,
                        shr(
                            252,
                            shl(
                                shl(2, shr(250, mul(x, 0xb6db6db6ddddddddd34d34d349249249210842108c6318c639ce739cffffffff))),
                                0x8040405543005266443200005020610674053026020000107506200176117077
                            )
                        )
                    )
                r :=
                    or(
                        r,
                        byte(
                            and(div(0xd76453e0, shr(r, x)), 0x1f),
                            0x001f0d1e100c1d070f090b19131c1706010e11080a1a141802121b1503160405
                        )
                    )
            }
            function nextRequestedRound(_round) -> _next, _requested {
                let wordPos := shr(8, _round)
                let bitPos := and(_round, 0xff)
                let mask := not(sub(shl(bitPos, 1), 1))
                mstore(0x00, wordPos)
                mstore(0x20, s_roundBitmap.slot)
                let masked := and(sload(keccak256(0x00, 0x40)), mask)
                _requested := gt(masked, 0)
                switch _requested
                case 1 { _next := sub(add(_round, leastSignificantBit(masked)), bitPos) }
                default { _next := sub(add(_round, 255), bitPos) }
            }
            for { let i } lt(i, 10) { i := add(i, 1) } {
                nextRound, requested := nextRequestedRound(nextRound)
                if requested {
                    // Start this requested round
                    mstore(0x00, nextRound)
                    mstore(0x20, s_requestInfo.slot)
                    sstore(add(keccak256(0x00, 0x40), 1), timestamp()) // startTime
                    sstore(s_isInProcess.slot, IN_PROGRESS)
                    mstore(0x40, IN_PROGRESS)
                    switch eq(nextRound, curRound)
                    case 1 {
                        mstore(0x20, s_trialNum.slot)
                        let trialNumSlot := keccak256(0x00, 0x40)
                        let newTrialNum := add(sload(trialNumSlot), 1)
                        sstore(trialNumSlot, newTrialNum)
                        mstore(0x20, newTrialNum)
                        log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
                    }
                    default {
                        sstore(s_currentRound.slot, nextRound)
                        mstore(0x20, 0) // trialNum is 0 for the first trial
                        log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
                    }
                    return(0, 0)
                }
                // If we reach or pass the last round without finding any requested round,
                // mark as COMPLETED and set the current round to the last possible index.
                if iszero(lt(nextRound, requestCountMinusOne)) {
                    sstore(s_isInProcess.slot, COMPLETED)
                    sstore(s_currentRound.slot, requestCountMinusOne)
                    mstore(0x00, requestCountMinusOne)
                    mstore(0x20, s_trialNum.slot)
                    mstore(0x20, sload(keccak256(0x00, 0x40))) // trialNum
                    mstore(0x40, COMPLETED)
                    log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
                    return(0, 0)
                }
                nextRound := add(nextRound, 1)
            }
            sstore(s_currentRound.slot, nextRound)
        }
    }
}
