// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {DisputeLogics} from "./DisputeLogics.sol";

contract FailLogics is DisputeLogics {
    // if eq(sload(s_isInProcess.slot), HALTED) {
    modifier notHalted() {
        assembly ("memory-safe") {
            // ** check if the contract is HALTED
            if eq(sload(s_isInProcess.slot), HALTED) {
                mstore(0, 0xd6c912e6) // selector for AlreadyHalted()
                revert(0x1c, 0x04)
            }
        }
        _;
    }

    //  if iszero(eq(sload(s_isInProcess.slot), IN_PROGRESS)) {
    modifier inProgress() {
        assembly ("memory-safe") {
            // ** check if the contract is COMPLETED or HALTED
            if iszero(eq(sload(s_isInProcess.slot), IN_PROGRESS)) {
                mstore(0, 0x6b4bc078) // RoundNotInProgress()
                revert(0x1c, 0x04)
            }
        }
        _;
    }

    constructor(string memory name, string memory version) DisputeLogics(name, version) {}

    function failToRequestSubmitCvOrSubmitMerkleRoot() external notHalted {
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNum := sload(keccak256(0x40, 0x40))
            mstore(0x00, trialNum)
            // ** Not requested to submit cv
            mstore(0x60, s_requestedToSubmitCvTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            if gt(sload(keccak256(0x00, 0x40)), 0) {
                mstore(0, 0x899a05f2) // AlreadyRequestedToSubmitCv()
                revert(0x1c, 0x04)
            }
            // ** MerkleRoot Not Submitted
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            if gt(sload(keccak256(0x00, 0x40)), 0) {
                mstore(0, 0x1c044d8b) // AlreadySubmittedMerkleRoot()
                revert(0x1c, 0x04)
            }
            // ** check time window
            mstore(0x00, curRound)
            mstore(0x20, s_requestInfo.slot)
            let startTime := sload(add(keccak256(0x00, 0x40), 1))
            if lt(
                timestamp(),
                add(
                    add(startTime, sload(s_offChainSubmissionPeriod.slot)),
                    sload(s_requestOrSubmitOrFailDecisionPeriod.slot)
                )
            ) {
                mstore(0, 0x085de625) // TooEarly()
                revert(0x1c, 0x04)
            }

            // ** Halt the round
            sstore(s_isInProcess.slot, HALTED)
            // 0x00 already has curRound
            mstore(0x20, trialNum)
            mstore(0x40, HALTED)
            log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
        }
        _executeSlashLeaderAndDistribute(FAILTOREQUESTSUBMITCV_OR_SUBMITMEKRLEROOT_OFFSET);
    }

    function failToSubmitMerkleRootAfterDispute() external notHalted {
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNum := sload(keccak256(0x40, 0x40))
            mstore(0x00, trialNum)
            // ** check if it is requested to submit cv
            mstore(0x60, s_requestedToSubmitCvTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCvTimestamp := sload(keccak256(0x00, 0x40))
            if iszero(requestedToSubmitCvTimestamp) {
                mstore(0, 0xd3e6c959) // CvNotRequested()
                revert(0x1c, 0x04)
            }
            // ** check time window
            if lt(
                timestamp(),
                add(
                    add(requestedToSubmitCvTimestamp, sload(s_onChainSubmissionPeriod.slot)),
                    sload(s_requestOrSubmitOrFailDecisionPeriod.slot)
                )
            ) {
                mstore(0, 0x085de625) // TooEarly()
                revert(0x1c, 0x04)
            }
            // ** MerkleRoot Not Submitted
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            if gt(sload(keccak256(0x00, 0x40)), 0) {
                mstore(0, 0x1c044d8b) // AlreadySubmittedMerkleRoot()
                revert(0x1c, 0x04)
            }
            // ** Halt the round
            mstore(0x00, curRound)
            mstore(0x20, trialNum)
            sstore(s_isInProcess.slot, HALTED)
            mstore(0x40, HALTED)
            log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
        }
        _executeSlashLeaderAndDistribute(FAILTOSUBMITMERKLEROOTAFTERDISPUTE_OFFSET);
    }

    function failToSubmitCv() external notHalted {
        uint256 returnGasFee = _getL1FeeUpperBoundOfFailFunction();
        uint256 getL1UpperBoundGasUsed = _getGetL1UpperBoundGasUsed();
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNumSlot := keccak256(0x40, 0x40)
            let trialNum := sload(trialNumSlot)
            mstore(0x00, trialNum)
            // ** check if it is requested to submit cv
            mstore(0x60, s_requestedToSubmitCvTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCvTimestamp := sload(keccak256(0x00, 0x40))
            if iszero(requestedToSubmitCvTimestamp) {
                mstore(0, 0xd3e6c959) // CvNotRequested()
                revert(0x1c, 0x04)
            }
            // ** check time window
            if lt(timestamp(), add(requestedToSubmitCvTimestamp, sload(s_onChainSubmissionPeriod.slot))) {
                mstore(0, 0x085de625) // TooEarly()
                revert(0x1c, 0x04)
            }
            // ** MerkleRoot Not Submitted
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            if gt(sload(keccak256(0x00, 0x40)), 0) {
                mstore(0, 0x1c044d8b) // AlreadySubmittedMerkleRoot()
                revert(0x1c, 0x04)
            }

            // ** who didn't submit cv even though requested
            let didntSubmitCvLength
            let addressToDeactivatesPtr := 0x80 // fmp
            let zeroBitIfSubmittedCvBitmap := sload(s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot)
            mstore(0x20, s_activatedOperators.slot)
            let firstActivatedOperatorSlot := keccak256(0x20, 0x20)
            mstore(0x20, sload(s_requestedToSubmitCvPackedIndicesAscFromLSB.slot))
            // Handle first iteration separately to avoid checking previousIndex
            let operatorIndex := and(mload(0x20), 0xff)
            if gt(and(zeroBitIfSubmittedCvBitmap, shl(operatorIndex, 1)), 0) {
                // if bit is still set, meaning no Cv submitted for this operator
                mstore(
                    add(addressToDeactivatesPtr, shl(5, didntSubmitCvLength)),
                    sload(add(firstActivatedOperatorSlot, operatorIndex))
                )
                didntSubmitCvLength := add(didntSubmitCvLength, 1)
            }
            let previousIndex := operatorIndex
            let requestedToSubmitLength
            // Continue with remaining iterations
            for { let i := 1 } true { i := add(i, 1) } {
                operatorIndex := and(mload(sub(0x20, i)), 0xff)
                if iszero(gt(operatorIndex, previousIndex)) {
                    requestedToSubmitLength := i
                    break
                }
                if gt(and(zeroBitIfSubmittedCvBitmap, shl(operatorIndex, 1)), 0) {
                    // if bit is still set, meaning no Cv submitted for this operator
                    mstore(
                        add(addressToDeactivatesPtr, shl(5, didntSubmitCvLength)),
                        sload(add(firstActivatedOperatorSlot, operatorIndex))
                    )
                    didntSubmitCvLength := add(didntSubmitCvLength, 1)
                }
                previousIndex := operatorIndex
            }
            if iszero(didntSubmitCvLength) {
                mstore(0, 0x7d39a81b) // AllSubmittedCv()
                revert(0x1c, 0x04)
            }
            let activatedOperatorLength := sload(s_activatedOperators.slot)

            // ** return gas fee to the caller()
            let dynamicFailToSubmitGasUsed := sload(s_failToSubmitCoGasUsedBaseA.slot)
            switch eq(requestedToSubmitLength, activatedOperatorLength)
            case 1 {
                returnGasFee :=
                    add(
                        returnGasFee,
                        add(
                            sub(
                                and(
                                    shr(FAILTOSUBMITCVGASUSEDBASEA_OFFSET, dynamicFailToSubmitGasUsed),
                                    DYNAMICFAILTOSUBMIT_MASK
                                ),
                                getL1UpperBoundGasUsed
                            ),
                            add(
                                mul(
                                    and(
                                        shr(PEROPERATORINCREASEGASUSEDA_OFFSET, dynamicFailToSubmitGasUsed),
                                        DYNAMICFAILTOSUBMIT_MASK
                                    ),
                                    activatedOperatorLength
                                ),
                                mul(
                                    and(
                                        shr(PERADDITIONALDIDNTSUBMITGASUSEDA_OFFSET, dynamicFailToSubmitGasUsed),
                                        DYNAMICFAILTOSUBMIT_MASK
                                    ),
                                    sub(didntSubmitCvLength, 1)
                                )
                            )
                        )
                    )
            }
            default {
                returnGasFee :=
                    add(
                        returnGasFee,
                        add(
                            sub(
                                and(
                                    shr(FAILTOSUBMITGASUSEDBASEB_OFFSET, dynamicFailToSubmitGasUsed),
                                    DYNAMICFAILTOSUBMIT_MASK
                                ),
                                getL1UpperBoundGasUsed
                            ),
                            add(
                                mul(
                                    and(
                                        shr(PEROPERATORINCREASEGASUSEDB_OFFSET, dynamicFailToSubmitGasUsed),
                                        DYNAMICFAILTOSUBMIT_MASK
                                    ),
                                    activatedOperatorLength
                                ),
                                add(
                                    mul(
                                        and(
                                            shr(PERREQUESTEDINCREASEGASUSED_OFFSET, dynamicFailToSubmitGasUsed),
                                            DYNAMICFAILTOSUBMIT_MASK
                                        ),
                                        requestedToSubmitLength
                                    ),
                                    mul(
                                        and(
                                            shr(PERADDITIONALDIDNTSUBMITGASUSEDB_OFFSET, dynamicFailToSubmitGasUsed),
                                            DYNAMICFAILTOSUBMIT_MASK
                                        ),
                                        sub(didntSubmitCvLength, 1)
                                    )
                                )
                            )
                        )
                    )
            }
            let activationThreshold := sload(s_activationThreshold.slot)
            if gt(returnGasFee, activationThreshold) { returnGasFee := activationThreshold } // if returnGasFee is greater than one operator's activationThreshold, set returnGasFee to activationThreshold
            mstore(0x20, caller())
            mstore(0x40, s_depositAmount.slot)
            let depositSlot := keccak256(0x20, 0x40) // msg.sender
            sstore(depositSlot, add(sload(depositSlot), returnGasFee))

            // ** cache slash rewards
            let slashRewardPerOperatorX8 := sload(s_slashRewardPerOperatorX8.slot)
            let distributeAmount := sub(mul(activationThreshold, didntSubmitCvLength), returnGasFee)
            let updatedSlashRewardPerOperatorX8 := slashRewardPerOperatorX8
            if gt(distributeAmount, 0) {
                // if distributeAmount is not zero
                updatedSlashRewardPerOperatorX8 :=
                    add(
                        slashRewardPerOperatorX8,
                        div(
                            shl(8, distributeAmount),
                            add(sub(activatedOperatorLength, didntSubmitCvLength), 1) // 1 for owner
                        )
                    )
                // ** update global slash reward
                sstore(s_slashRewardPerOperatorX8.slot, updatedSlashRewardPerOperatorX8)
            }

            // ** update slash reward and deactivate for non cv submitters
            let fmp := add(addressToDeactivatesPtr, shl(5, didntSubmitCvLength)) // traverse in reverse order
            for { let i } lt(i, didntSubmitCvLength) { i := add(i, 1) } {
                addressToDeactivatesPtr := sub(fmp, 0x20)
                // ** update slashRewardPerOperatorPaid
                mstore(fmp, s_slashRewardPerOperatorPaidX8.slot)
                let slotToUpdate := keccak256(addressToDeactivatesPtr, 0x40) // s_slashRewardPerOperatorPaidX8[operator]
                let accumulatedReward := shr(8, sub(slashRewardPerOperatorX8, sload(slotToUpdate)))
                sstore(slotToUpdate, updatedSlashRewardPerOperatorX8)
                // ** update deposit Amount
                mstore(fmp, s_depositAmount.slot)
                slotToUpdate := keccak256(addressToDeactivatesPtr, 0x40) // s_depositAmount[operator]
                sstore(slotToUpdate, add(sub(sload(slotToUpdate), activationThreshold), accumulatedReward))

                // ** deactivate operator
                mstore(fmp, s_activatedOperatorIndex1Based.slot)
                let operatorToDeactivateIndex := sub(sload(keccak256(addressToDeactivatesPtr, 0x40)), 1)
                let operatorToDeactivate := mload(addressToDeactivatesPtr)
                activatedOperatorLength := sub(activatedOperatorLength, 1)
                let lastOperatorIndex := activatedOperatorLength
                let lastOperatorAddress := sload(add(firstActivatedOperatorSlot, lastOperatorIndex))
                // ** activatedOperatorIndex1Based = 0
                sstore(keccak256(addressToDeactivatesPtr, 0x40), 0)
                log1(addressToDeactivatesPtr, 0x20, 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d) // `DeActivated(address operator)`.

                if iszero(eq(lastOperatorAddress, operatorToDeactivate)) {
                    sstore(add(firstActivatedOperatorSlot, operatorToDeactivateIndex), lastOperatorAddress)
                    mstore(addressToDeactivatesPtr, lastOperatorAddress) // overwrite because it is not used anymore
                    sstore(keccak256(addressToDeactivatesPtr, 0x40), add(operatorToDeactivateIndex, 1)) // activatedOperatorIndex1Based
                }

                // ** update addressToDeactivatesPtr
                fmp := sub(fmp, 0x20)
            }
            // ** update activatedOperators
            sstore(s_activatedOperators.slot, activatedOperatorLength)

            // ** restart or end this round
            mstore(0x00, curRound)
            switch gt(sload(s_activatedOperators.slot), 1)
            case 1 {
                mstore(0x20, s_requestInfo.slot)
                sstore(add(keccak256(0x00, 0x40), 1), timestamp()) // startTime
                let newTrialNum := add(trialNum, 1) // trialNum++
                sstore(trialNumSlot, newTrialNum)
                // 0x00 already has curRound
                mstore(0x20, newTrialNum)
                mstore(0x40, IN_PROGRESS)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
            default {
                sstore(s_isInProcess.slot, HALTED)
                // 0x00 already has curRound
                mstore(0x20, trialNum)
                mstore(0x40, HALTED)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
        }
    }

    function failToSubmitCo() external inProgress {
        uint256 returnGasFee = _getL1FeeUpperBoundOfFailFunction();
        uint256 getL1UpperBoundGasUsed = _getGetL1UpperBoundGasUsed();
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNumSlot := keccak256(0x40, 0x40)
            let trialNum := sload(trialNumSlot)
            mstore(0x00, trialNum)
            // ** check if it is requested to submit co
            mstore(0x60, s_requestedToSubmitCoTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCoTimestamp := sload(keccak256(0x00, 0x40))
            if iszero(requestedToSubmitCoTimestamp) {
                mstore(0, 0x11974969) // CoNotRequested()
                revert(0x1c, 0x04)
            }
            // ** check time window
            if lt(timestamp(), add(requestedToSubmitCoTimestamp, sload(s_onChainSubmissionPeriod.slot))) {
                mstore(0, 0x085de625) // TooEarly()
                revert(0x1c, 0x04)
            }

            // ** who didn't submit co even though requested
            let requestedToSubmitCoLength := sload(s_requestedToSubmitCoLength.slot)
            let didntSubmitCoLength
            let addressToDeactivatesPtr := 0x80 // fmp
            let zeroBitIfSubmittedCoBitmap := sload(s_zeroBitIfSubmittedCoBitmap.slot)
            mstore(0x20, s_activatedOperators.slot)
            let firstActivatedOperatorSlot := keccak256(0x20, 0x20)
            mstore(0x20, sload(s_requestedToSubmitCoPackedIndices.slot))
            for { let i } lt(i, requestedToSubmitCoLength) { i := add(i, 1) } {
                let operatorIndex := and(mload(sub(0x20, i)), 0xff)
                if gt(and(zeroBitIfSubmittedCoBitmap, shl(operatorIndex, 1)), 0) {
                    // if bit is still set, meaning no Co submitted for this operator
                    mstore(
                        add(addressToDeactivatesPtr, shl(5, didntSubmitCoLength)),
                        sload(add(firstActivatedOperatorSlot, operatorIndex))
                    )
                    didntSubmitCoLength := add(didntSubmitCoLength, 1)
                }
            }
            if iszero(didntSubmitCoLength) {
                mstore(0, 0x1c7f7cc9) // AllSubmittedCo()
                revert(0x1c, 0x04)
            }

            // ** return gas fee to the caller()
            let activatedOperatorLength := sload(s_activatedOperators.slot)
            let dynamicFailToSubmitGasUsed := sload(s_failToSubmitCoGasUsedBaseA.slot)
            switch eq(requestedToSubmitCoLength, activatedOperatorLength)
            case 1 {
                returnGasFee :=
                    add(
                        returnGasFee,
                        add(
                            sub(and(dynamicFailToSubmitGasUsed, DYNAMICFAILTOSUBMIT_MASK), getL1UpperBoundGasUsed),
                            add(
                                mul(
                                    and(
                                        shr(PEROPERATORINCREASEGASUSEDA_OFFSET, dynamicFailToSubmitGasUsed),
                                        DYNAMICFAILTOSUBMIT_MASK
                                    ),
                                    activatedOperatorLength
                                ),
                                mul(
                                    and(
                                        shr(PERADDITIONALDIDNTSUBMITGASUSEDA_OFFSET, dynamicFailToSubmitGasUsed),
                                        DYNAMICFAILTOSUBMIT_MASK
                                    ),
                                    sub(didntSubmitCoLength, 1)
                                )
                            )
                        )
                    )
            }
            default {
                returnGasFee :=
                    add(
                        returnGasFee,
                        add(
                            sub(
                                and(
                                    shr(FAILTOSUBMITGASUSEDBASEB_OFFSET, dynamicFailToSubmitGasUsed),
                                    DYNAMICFAILTOSUBMIT_MASK
                                ),
                                getL1UpperBoundGasUsed
                            ),
                            add(
                                mul(
                                    and(
                                        shr(PEROPERATORINCREASEGASUSEDB_OFFSET, dynamicFailToSubmitGasUsed),
                                        DYNAMICFAILTOSUBMIT_MASK
                                    ),
                                    activatedOperatorLength
                                ),
                                add(
                                    mul(
                                        and(
                                            shr(PERREQUESTEDINCREASEGASUSED_OFFSET, dynamicFailToSubmitGasUsed),
                                            DYNAMICFAILTOSUBMIT_MASK
                                        ),
                                        requestedToSubmitCoLength
                                    ),
                                    mul(
                                        and(
                                            shr(PERADDITIONALDIDNTSUBMITGASUSEDB_OFFSET, dynamicFailToSubmitGasUsed),
                                            DYNAMICFAILTOSUBMIT_MASK
                                        ),
                                        sub(didntSubmitCoLength, 1)
                                    )
                                )
                            )
                        )
                    )
            }
            let activationThreshold := sload(s_activationThreshold.slot)
            if gt(returnGasFee, activationThreshold) { returnGasFee := activationThreshold } // if returnGasFee is greater than one operator's activationThreshold, set returnGasFee to activationThreshold
            mstore(0x20, caller())
            mstore(0x40, s_depositAmount.slot)
            let depositSlot := keccak256(0x20, 0x40) // msg.sender
            sstore(depositSlot, add(sload(depositSlot), returnGasFee))

            // ** cache slash rewards
            let slashRewardPerOperatorX8 := sload(s_slashRewardPerOperatorX8.slot)
            let distributeAmount := sub(mul(activationThreshold, didntSubmitCoLength), returnGasFee)
            let updatedSlashRewardPerOperatorX8 := slashRewardPerOperatorX8
            if gt(distributeAmount, 0) {
                updatedSlashRewardPerOperatorX8 :=
                    add(
                        slashRewardPerOperatorX8,
                        div(
                            shl(8, sub(mul(activationThreshold, didntSubmitCoLength), returnGasFee)),
                            add(sub(activatedOperatorLength, didntSubmitCoLength), 1) // 1 for owner
                        )
                    )
            }
            // ** update global slash reward
            sstore(s_slashRewardPerOperatorX8.slot, updatedSlashRewardPerOperatorX8)

            // ** update slash reward and deactivate for non co submitters
            let fmp := add(addressToDeactivatesPtr, shl(5, didntSubmitCoLength)) // traverse in reverse order
            for { let i } lt(i, didntSubmitCoLength) { i := add(i, 1) } {
                addressToDeactivatesPtr := sub(fmp, 0x20)
                // ** update slashRewardPerOperatorPaid
                mstore(fmp, s_slashRewardPerOperatorPaidX8.slot)
                let slotToUpdate := keccak256(addressToDeactivatesPtr, 0x40) // s_slashRewardPerOperatorPaidX8[operator]
                let accumulatedReward := shr(8, sub(slashRewardPerOperatorX8, sload(slotToUpdate)))
                sstore(slotToUpdate, updatedSlashRewardPerOperatorX8)
                // ** update deposit Amount
                mstore(fmp, s_depositAmount.slot)
                slotToUpdate := keccak256(addressToDeactivatesPtr, 0x40) // s_depositAmount[operator]
                sstore(slotToUpdate, add(sub(sload(slotToUpdate), activationThreshold), accumulatedReward))

                // ** deactivate operator
                mstore(fmp, s_activatedOperatorIndex1Based.slot)
                let operatorToDeactivateIndex := sub(sload(keccak256(addressToDeactivatesPtr, 0x40)), 1)
                let operatorToDeactivate := mload(addressToDeactivatesPtr)
                activatedOperatorLength := sub(activatedOperatorLength, 1)
                let lastOperatorIndex := activatedOperatorLength
                let lastOperatorAddress := sload(add(firstActivatedOperatorSlot, lastOperatorIndex))
                // ** activatedOperatorIndex1Based = 0
                sstore(keccak256(addressToDeactivatesPtr, 0x40), 0)
                log1(addressToDeactivatesPtr, 0x20, 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d) // `DeActivated(address operator)`.

                if iszero(eq(lastOperatorAddress, operatorToDeactivate)) {
                    sstore(add(firstActivatedOperatorSlot, operatorToDeactivateIndex), lastOperatorAddress)
                    mstore(addressToDeactivatesPtr, lastOperatorAddress) // overwrite because it is not used anymore
                    sstore(keccak256(addressToDeactivatesPtr, 0x40), add(operatorToDeactivateIndex, 1)) // activatedOperatorIndex1Based
                }

                // ** update addressToDeactivatesPtr
                fmp := sub(fmp, 0x20)
            }
            // ** update activatedOperatorLength
            sstore(s_activatedOperators.slot, activatedOperatorLength)

            // ** restart or end this round
            mstore(0x00, curRound)
            switch gt(sload(s_activatedOperators.slot), 1)
            case 1 {
                mstore(0x20, s_requestInfo.slot)
                sstore(add(keccak256(0x00, 0x40), 1), timestamp()) // startTime
                let newTrialNum := add(trialNum, 1) // trialNum++
                sstore(trialNumSlot, newTrialNum)
                // 0x00 already has curRound
                mstore(0x20, newTrialNum)
                mstore(0x40, IN_PROGRESS)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
            default {
                sstore(s_isInProcess.slot, HALTED)
                // 0x00 already has curRound
                mstore(0x20, trialNum)
                mstore(0x40, HALTED)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // emit Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
        }
    }

    function failToSubmitS() external inProgress {
        uint256 returnGasFee = _calculateFailGasFee(FAILTOSUBMITS_OFFSET);
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNumSlot := keccak256(0x40, 0x40)
            let trialNum := sload(trialNumSlot)
            mstore(0x00, trialNum)
            // ** Ensure S was requested
            mstore(0x60, s_previousSSubmitTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let previousSSubmitTimestamp := sload(keccak256(0x00, 0x40))
            if iszero(previousSSubmitTimestamp) {
                mstore(0, 0x2d37f8d3) // SNotRequested()
                revert(0x1c, 0x04)
            }
            // ** check time window
            if lt(timestamp(), add(previousSSubmitTimestamp, sload(s_onChainSubmissionPeriodPerOperator.slot))) {
                mstore(0, 0x085de625) // TooEarly()
                revert(0x1c, 0x04)
            }

            // ** Refund gas fee to the caller()
            let activationThreshold := sload(s_activationThreshold.slot)
            if gt(returnGasFee, activationThreshold) { returnGasFee := activationThreshold }
            mstore(0x20, caller())
            mstore(0x40, s_depositAmount.slot)
            let depositSlot := keccak256(0x20, 0x40) // msg.sender
            sstore(depositSlot, add(sload(depositSlot), returnGasFee))

            // ** Update slash reward
            let slashRewardPerOperatorX8 := sload(s_slashRewardPerOperatorX8.slot)
            let activatedOperatorLength := sload(s_activatedOperators.slot)
            let distributeAmount := sub(activationThreshold, returnGasFee)
            let updatedSlashRewardPerOperatorX8 := slashRewardPerOperatorX8
            if gt(distributeAmount, 0) {
                updatedSlashRewardPerOperatorX8 :=
                    add(
                        slashRewardPerOperatorX8,
                        div(
                            shl(8, sub(activationThreshold, returnGasFee)),
                            activatedOperatorLength // 1 for owner
                        )
                    )
                sstore(s_slashRewardPerOperatorX8.slot, updatedSlashRewardPerOperatorX8)
            }

            // ** s_revealOrders[s_requestedToSubmitSFromIndexK] is the index of the operator who didn't submit S
            mstore(0x20, sload(s_packedRevealOrders.slot))
            let operatorToDeactivateIndex := and(mload(sub(0x20, sload(s_requestedToSubmitSFromIndexK.slot))), 0xff)
            mstore(0x20, s_activatedOperators.slot)
            let firstActivatedOperatorSlot := keccak256(0x20, 0x20)
            let operatorToDeactivate := sload(add(firstActivatedOperatorSlot, operatorToDeactivateIndex))
            // ** update deposit amount
            mstore(0x20, operatorToDeactivate)
            mstore(0x40, s_depositAmount.slot)
            depositSlot := keccak256(0x20, 0x40) // operatorToDeactivate
            mstore(0x40, s_slashRewardPerOperatorPaidX8.slot)
            let slashRewardPerOperatorPaidX8Slot := keccak256(0x20, 0x40) // s_slashRewardPerOperatorPaid[operatorToDeactivate]
            sstore(
                depositSlot,
                add(
                    sub(sload(depositSlot), activationThreshold),
                    shr(8, sub(slashRewardPerOperatorX8, sload(slashRewardPerOperatorPaidX8Slot)))
                )
            )
            sstore(slashRewardPerOperatorPaidX8Slot, updatedSlashRewardPerOperatorX8)
            // ** deactivate operator
            activatedOperatorLength := sub(activatedOperatorLength, 1)
            let lastOperatorIndex := activatedOperatorLength
            let lastOperatorAddress := sload(add(firstActivatedOperatorSlot, lastOperatorIndex))
            // ** activatedOperatorIndex1Based = 0
            mstore(0x40, s_activatedOperatorIndex1Based.slot)
            sstore(keccak256(0x20, 0x40), 0)
            if iszero(eq(lastOperatorAddress, operatorToDeactivate)) {
                sstore(add(firstActivatedOperatorSlot, operatorToDeactivateIndex), lastOperatorAddress)
                mstore(0x20, lastOperatorAddress)
                sstore(keccak256(0x20, 0x40), add(operatorToDeactivateIndex, 1)) // activatedOperatorIndex1Based
            }
            // ** update activatedOperatorLength
            sstore(s_activatedOperators.slot, activatedOperatorLength)

            // ** restart or end this round
            mstore(0x00, curRound)
            switch gt(sload(s_activatedOperators.slot), 1)
            case 1 {
                mstore(0x20, s_requestInfo.slot)
                sstore(add(keccak256(0x00, 0x40), 1), timestamp()) // startTime
                let newTrialNum := add(trialNum, 1) // trialNum++
                sstore(trialNumSlot, newTrialNum)
                // 0x00 already has curRound
                mstore(0x20, newTrialNum)
                mstore(0x40, IN_PROGRESS)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
            default {
                sstore(s_isInProcess.slot, HALTED)
                // 0x00 already has curRound
                mstore(0x20, trialNum)
                mstore(0x40, HALTED)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
        }
    }

    function failToRequestSorGenerateRandomNumber() external inProgress {
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNum := sload(keccak256(0x40, 0x40))
            mstore(0x00, trialNum)
            // ** Ensure S was not requested
            mstore(0x60, s_previousSSubmitTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            if gt(sload(keccak256(0x00, 0x40)), 0) {
                mstore(0, 0x53489cf9) // SRequested()
                revert(0x1c, 0x04)
            }
            // ** Ensure Merkle Root is submitted
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let merkleRootSubmittedTimestamp := sload(keccak256(0x00, 0x40))
            if iszero(merkleRootSubmittedTimestamp) {
                mstore(0, 0x8e56b845) // MerkleRootNotSubmitted()
                revert(0x1c, 0x04)
            }
            // ** check time window
            mstore(0x60, s_requestedToSubmitCoTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCoTimestamp := sload(keccak256(0x00, 0x40))
            let activatedOperatorLength := sload(s_activatedOperators.slot)
            switch gt(requestedToSubmitCoTimestamp, 0)
            case 1 {
                if lt(
                    timestamp(),
                    add(
                        add(
                            add(requestedToSubmitCoTimestamp, sload(s_onChainSubmissionPeriod.slot)),
                            mul(sload(s_offChainSubmissionPeriodPerOperator.slot), activatedOperatorLength)
                        ),
                        sload(s_requestOrSubmitOrFailDecisionPeriod.slot)
                    )
                ) {
                    mstore(0, 0x085de625) // TooEarly()
                    revert(0x1c, 0x04)
                }
            }
            default {
                if lt(
                    timestamp(),
                    add(
                        add(
                            add(merkleRootSubmittedTimestamp, sload(s_offChainSubmissionPeriod.slot)),
                            mul(sload(s_offChainSubmissionPeriodPerOperator.slot), activatedOperatorLength)
                        ),
                        sload(s_requestOrSubmitOrFailDecisionPeriod.slot)
                    )
                ) {
                    mstore(0, 0x085de625) // TooEarly()
                    revert(0x1c, 0x04)
                }
            }
            // ** Halt the round
            sstore(s_isInProcess.slot, HALTED)
            mstore(0x00, curRound)
            mstore(0x20, trialNum)
            mstore(0x40, HALTED)
            log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
        }
        _executeSlashLeaderAndDistribute(FAILTOREQUESTS_OR_GENERATERANDOMNUMBER_OFFSET);
    }

    function _executeSlashLeaderAndDistribute(uint256 bitsToShiftRight) internal {
        uint256 returnGasFee = _calculateFailGasFee(bitsToShiftRight);
        assembly ("memory-safe") {
            let activationThreshold := sload(s_activationThreshold.slot)
            mstore(0x20, sload(_OWNER_SLOT))
            // ** Distribute remainder among operators
            if gt(activationThreshold, returnGasFee) {
                let delta := div(shl(8, sub(activationThreshold, returnGasFee)), sload(s_activatedOperators.slot))
                sstore(s_slashRewardPerOperatorX8.slot, add(sload(s_slashRewardPerOperatorX8.slot), delta))
                mstore(0x40, s_slashRewardPerOperatorPaidX8.slot)
                let slashRewardPerOperatorPaidX8Slot := keccak256(0x20, 0x40) // owner
                sstore(slashRewardPerOperatorPaidX8Slot, add(sload(slashRewardPerOperatorPaidX8Slot), delta))
            }
            if gt(returnGasFee, activationThreshold) { returnGasFee := activationThreshold }
            // ** slash the leadernode(owner)
            mstore(0x40, s_depositAmount.slot)
            let depositSlot := keccak256(0x20, 0x40) // owner
            sstore(depositSlot, sub(sload(depositSlot), activationThreshold))
            // ** return gas fee to the caller()
            mstore(0x20, caller())
            depositSlot := keccak256(0x20, 0x40) // msg.sender
            sstore(depositSlot, add(sload(depositSlot), returnGasFee))
        }
    }

    function _calculateFailGasFee(uint256 bitsToShiftRight) internal view virtual returns (uint256 gasFee) {
        assembly ("memory-safe") {
            let failgasUsed := sload(s_getL1UpperBoundGasUsedWhenCalldataSize4.slot)
            gasFee :=
                mul(
                    gasprice(),
                    sub(and(shr(bitsToShiftRight, failgasUsed), FAILTOSUBMIT_MASK), and(failgasUsed, FAILTOSUBMIT_MASK))
                )
        }
    }

    function _getL1FeeUpperBoundOfFailFunction() internal view virtual returns (uint256 l1GasFee) {
        return 0;
    }

    function _getGetL1UpperBoundGasUsed() internal view virtual returns (uint256 getL1UpperBoundGasUsed) {
        assembly ("memory-safe") {
            getL1UpperBoundGasUsed := and(sload(s_getL1UpperBoundGasUsedWhenCalldataSize4.slot), FAILTOSUBMIT_MASK)
        }
    }
}
