// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {OperatorManager} from "./OperatorManager.sol";
import {CommitReveal2Storage} from "./CommitReveal2Storage.sol";
import {EIP712} from "@openzeppelin/contracts_v5.0.1/utils/cryptography/EIP712.sol";

contract DisputeLogics is EIP712, OperatorManager, CommitReveal2Storage {
    constructor(string memory name, string memory version) EIP712(name, version) {}

    function requestToSubmitCv(uint256 packedIndicesAscendingFromLSB) external onlyOwner {
        assembly ("memory-safe") {
            // mstore(0x00, sload(s_currentRound.slot))
            // mstore(0x20, s_requestInfo.slot)
            // mstore(0x00, sload(add(keccak256(0x00, 0x40), 1))) // startTime

            let curRound := sload(s_currentRound.slot)
            mstore(0x60, curRound)
            mstore(0x80, s_trialNum.slot)
            mstore(0x20, sload(keccak256(0x60, 0x40))) // trialNum
            // * get requestedToSubmitCvTimestamp
            mstore(0x80, s_requestedToSubmitCvTimestamp.slot)
            mstore(0x40, keccak256(0x60, 0x40))
            let requestedToSubmitCvTimestampSlot := keccak256(0x20, 0x40)
            if gt(sload(requestedToSubmitCvTimestampSlot), 0) {
                mstore(0, 0x899a05f2) // AlreadyRequestedToSubmitCv()
                revert(0x1c, 0x04)
            }
            // * get merkleRootSubmittedTimestamp
            mstore(0x80, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x40, keccak256(0x60, 0x40))
            if gt(sload(keccak256(0x20, 0x40)), 0) {
                mstore(0, 0xf6b442ac) // MerkleRootIsSubmitted()
                revert(0x1c, 0x04)
            }
            let bitSetIfRequestedToSubmitCv
            let maxIndex := sub(sload(s_activatedOperators.slot), 1) // max index
            let previousIndex := and(packedIndicesAscendingFromLSB, 0xff)
            if gt(previousIndex, maxIndex) {
                mstore(0, 0x63df8171) // InvalidIndex()
                revert(0x1c, 0x04)
            }
            bitSetIfRequestedToSubmitCv := or(bitSetIfRequestedToSubmitCv, shl(previousIndex, 1))
            mstore(0x40, packedIndicesAscendingFromLSB)
            for { let i := 1 } true { i := add(i, 1) } {
                let currentIndex := and(mload(sub(0x40, i)), 0xff)
                if gt(currentIndex, maxIndex) {
                    mstore(0, 0x63df8171) // InvalidIndex()
                    revert(0x1c, 0x04)
                }
                if iszero(gt(currentIndex, previousIndex)) { break }
                bitSetIfRequestedToSubmitCv := or(bitSetIfRequestedToSubmitCv, shl(currentIndex, 1))
                previousIndex := currentIndex
            }
            sstore(requestedToSubmitCvTimestampSlot, timestamp())
            sstore(s_requestedToSubmitCvPackedIndicesAscFromLSB.slot, packedIndicesAscendingFromLSB)
            sstore(
                s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot,
                or(shl(128, bitSetIfRequestedToSubmitCv), 0xffffffff)
            ) // set zeroBitIfSubmittedCvBitmap all bits to 1
            mstore(0x00, curRound) // 0x20 already has trialNum, 0x40 already has packedIndicesAscendingFromLSB
            log1(0x00, 0x60, 0x16759d80d11394de93184cfeb4e91cf57282cef239f68ed141c496600454f757) // event RequestedToSubmitCv(uint256 round, uint256 trialNum, uint256 packedIndicesAscendingFromLSB)
        }
    }

    function submitCv(bytes32 cv) external {
        assembly ("memory-safe") {
            mstore(0x00, caller())
            mstore(0x20, s_activatedOperatorIndex1Based.slot)
            let activatedOperatorIndex := sub(sload(keccak256(0x00, 0x40)), 1) // overflows when s_activatedOperatorIndex1Based is 0
            let bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2 :=
                sload(s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot)
            if gt(activatedOperatorIndex, MAX_OPERATOR_INDEX) {
                mstore(0, 0x1b256530) // NotActivatedOperator()
                revert(0x1c, 0x04)
            }
            if iszero(
                and(
                    shr(128, bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2),
                    shl(activatedOperatorIndex, 1)
                )
            ) {
                mstore(0, 0x998cf22e) // CvNotRequestedForThisOperator()
                revert(0x1c, 0x04)
            }
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            mstore(0x20, sload(keccak256(0x40, 0x40))) // trialNum
            // * get merkleRootSubmittedTimestamp
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x40, keccak256(0x40, 0x40))
            // ** can only submit cv if merkleRoot is not submitted
            if gt(sload(keccak256(0x20, 0x40)), 0) {
                mstore(0, 0xf6b442ac) // MerkleRootIsSubmitted()
                revert(0x1c, 0x04)
            }
            sstore(add(s_cvs.slot, activatedOperatorIndex), cv)
            sstore(
                s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot,
                and(bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2, not(shl(activatedOperatorIndex, 1)))
            ) // set to zero
            mstore(0x00, curRound) // 0x20 already has trialNum
            mstore(0x40, cv)
            mstore(0x60, activatedOperatorIndex)
            log1(0x00, 0x80, 0x6a6385c5eaed19d346ec4f9bd0010cfba4ac1d0407e2e55f959cb8fcac30f873) // event CvSubmitted(uint256 round, uint256 trialNum, bytes32 cv, uint256 index)
        }
    }

    function requestToSubmitCo(
        CvAndSigRS[] calldata cvRSsForCvsNotOnChainAndReqToSubmitCo,
        uint256, // packedVsForCvsNotOnChainAndReqToSubmitCo,
        uint256 indicesLength,
        uint256 packedIndicesFirstCvNotOnChainRestCvOnChain
    ) external onlyOwner {
        bytes32 domainSeparator = _domainSeparatorV4();
        assembly ("memory-safe") {
            if iszero(indicesLength) {
                mstore(0, 0xbf557497) // ZeroLength()
                revert(0x1c, 0x04)
            }
            if gt(indicesLength, MAX_ACTIVATED_OPERATORS) {
                mstore(0, 0x12466af8) // LengthExceedsMax()
                revert(0x1c, 0x04)
            }

            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNum := sload(keccak256(0x40, 0x40))
            mstore(0x00, trialNum)
            // * get merkleRootSubmittedTimestamp
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let merkleRootSubmittedTimestamp := sload(keccak256(0x00, 0x40))
            // ** can only request to submit co if merkleRoot is submitted
            if iszero(merkleRootSubmittedTimestamp) {
                mstore(0, 0x8e56b845) // MerkleRootNotSubmitted()
                revert(0x1c, 0x04)
            }
            // ** check if already requested to submit co
            mstore(0x60, s_requestedToSubmitCoTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCoTimestampSlot := keccak256(0x00, 0x40)
            if gt(sload(requestedToSubmitCoTimestampSlot), 0) {
                mstore(0, 0x13efcda2) // AlreadyRequestedToSubmitCo()
                revert(0x1c, 0x04)
            }
            // ** check time window
            if gt(
                timestamp(),
                add(
                    merkleRootSubmittedTimestamp,
                    add(sload(s_offChainSubmissionPeriod.slot), sload(s_requestOrSubmitOrFailDecisionPeriod.slot))
                )
            ) {
                mstore(0, 0xecdd1c29) // TooLate()
                revert(0x1c, 0x04)
            }

            // ** check cv status
            let operatorsLength := sload(s_activatedOperators.slot)
            mstore(0x60, s_requestedToSubmitCvTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCvTimestampSlot := keccak256(0x00, 0x40)
            // if not requested to submit cv, it means no cvs are on-chain
            let zeroBitIfSubmittedCvBitmap
            switch sload(requestedToSubmitCvTimestampSlot)
            case 0 {
                if iszero(eq(indicesLength, cvRSsForCvsNotOnChainAndReqToSubmitCo.length)) {
                    mstore(0, 0xad029eb9)
                    revert(0x1c, 0x04) // AllCvsNotSubmitted()
                }
                sstore(requestedToSubmitCvTimestampSlot, 1) // set to 1 to indicate that cvs are on-chain
                zeroBitIfSubmittedCvBitmap := 0xffffffff // set all bits to 1
            }
            default {
                zeroBitIfSubmittedCvBitmap := sload(s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot)
            }
            let maxIndex := sub(operatorsLength, 1) // max index
            let checkDuplicate

            let fmp := 0x80 // fmp
            mstore(fmp, MESSAGE_TYPEHASH_DIRECT)
            mstore(add(fmp, 0x20), curRound)
            mstore(add(fmp, 0x40), trialNum)
            mstore(add(fmp, 0x80), hex"1901") // prefix and version
            mstore(add(fmp, 0x82), domainSeparator)

            for { let i } lt(i, cvRSsForCvsNotOnChainAndReqToSubmitCo.length) { i := add(i, 1) } {
                // ** check duplicate
                let requestToSubmitCoIndex := and(calldataload(sub(0x64, i)), 0xff) // 0x64: packedIndicesFirstCvNotOnChainRestCvOnChain
                if gt(requestToSubmitCoIndex, maxIndex) {
                    // if greater than max index
                    mstore(0, 0x63df8171) // InvalidIndex()
                    revert(0x1c, 0x04)
                }
                let mask := shl(requestToSubmitCoIndex, 1)
                if gt(and(checkDuplicate, mask), 0) {
                    // if already set
                    mstore(0, 0x7a69f8d3) // DuplicateIndices()
                    revert(0x1c, 0x04)
                }
                checkDuplicate := or(checkDuplicate, mask)

                // ** check signature
                let cvsRSsOffset := add(cvRSsForCvsNotOnChainAndReqToSubmitCo.offset, mul(0x60, i))
                let s := calldataload(add(cvsRSsOffset, 0x40))
                if gt(s, SECP256K1_CURVE_ORDER) {
                    mstore(0, 0xbf4bf5b8) // InvalidSignatureS()
                    revert(0x1c, 0x04)
                }
                mstore(add(fmp, 0x60), calldataload(cvsRSsOffset)) // cv
                mstore(add(fmp, 0xa2), keccak256(fmp, 0x80)) // structHash
                mstore(0x00, keccak256(add(fmp, 0x80), 0x42)) // digest hash
                mstore(0x20, and(calldataload(sub(0x24, i)), 0xff)) // v, 0x24: packedVsForCvsNotOnChainAndReqToSubmitCo offset
                mstore(0x40, calldataload(add(cvsRSsOffset, 0x20))) // r
                mstore(0x60, s)
                let operatorAddress := mload(staticcall(gas(), 1, 0x00, 0x80, 0x01, 0x20))
                // `returndatasize()` will be `0x20` upon success, and `0x00` otherwise.
                if iszero(returndatasize()) {
                    mstore(0x00, 0x8baa579f) // selector for InvalidSignature()
                    revert(0x1c, 0x04)
                }
                mstore(0x00, operatorAddress)
                mstore(0x20, s_activatedOperatorIndex1Based.slot)
                if iszero(eq(add(requestToSubmitCoIndex, 1), sload(keccak256(0x00, 0x40)))) {
                    mstore(0, 0x980c4296) // SignatureAndIndexDoNotMatch()
                    revert(0x1c, 0x04)
                }

                // ** submit cv on-chain
                sstore(add(s_cvs.slot, requestToSubmitCoIndex), calldataload(cvsRSsOffset)) // cv
                zeroBitIfSubmittedCvBitmap := and(zeroBitIfSubmittedCvBitmap, not(shl(requestToSubmitCoIndex, 1))) // set to zero
            }
            sstore(s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot, zeroBitIfSubmittedCvBitmap) // update bitmap

            // ** Operators who already submitted Cv on-chain, simply confirm it exists
            for { let i := cvRSsForCvsNotOnChainAndReqToSubmitCo.length } lt(i, indicesLength) { i := add(i, 1) } {
                let requestToSubmitCoIndex := and(calldataload(sub(0x64, i)), 0xff) // 0x64: packedIndicesFirstCvNotOnChainRestCvOnChain
                // ** check duplicate
                if gt(requestToSubmitCoIndex, maxIndex) {
                    // if greater than max index
                    mstore(0, 0x63df8171) // InvalidIndex()
                    revert(0x1c, 0x04)
                }
                let mask := shl(requestToSubmitCoIndex, 1)
                if gt(and(checkDuplicate, mask), 0) {
                    // if already set
                    mstore(0, 0x7a69f8d3) // DuplicateIndices()
                    revert(0x1c, 0x04)
                }
                checkDuplicate := or(checkDuplicate, mask)
                // ** check cv bitmap
                if gt(and(zeroBitIfSubmittedCvBitmap, mask), 0) {
                    // if bit is still set, meaning no Cv submitted for this operator
                    mstore(0, 0x03798920) // CvNotSubmitted()
                    revert(0x1c, 0x04)
                }
            }

            sstore(s_requestedToSubmitCoPackedIndices.slot, packedIndicesFirstCvNotOnChainRestCvOnChain)
            sstore(s_requestedToSubmitCoLength.slot, indicesLength)
            sstore(requestedToSubmitCoTimestampSlot, timestamp())
            sstore(s_zeroBitIfSubmittedCoBitmap.slot, 0xffffffff) // set all bits to 1

            // ** event
            mstore(0x00, curRound)
            mstore(0x20, trialNum)
            mstore(0x40, indicesLength)
            mstore(0x60, packedIndicesFirstCvNotOnChainRestCvOnChain)
            log1(0x00, 0x80, 0xd4cc5cd95f180f10aaacba0729abc069b8080ec3a7e8e41856decb17bdc28ece) // event RequestedToSubmitCo(uint256 round, uint256 trialNum, uint256 indicesLength, uint256 packedIndices);
        }
    }

    function submitCo(bytes32 co) external {
        assembly ("memory-safe") {
            // ** check co status
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            mstore(0x20, sload(keccak256(0x40, 0x40))) // trialNum
            mstore(0x60, s_requestedToSubmitCoTimestamp.slot)
            mstore(0x40, keccak256(0x40, 0x40))
            if iszero(sload(keccak256(0x20, 0x40))) {
                mstore(0, 0x11974969) // CoNotRequested()
                revert(0x1c, 0x04)
            }
            // ** check cv == hash(co)
            mstore(0x40, caller())
            mstore(0x60, s_activatedOperatorIndex1Based.slot)
            let activatedOperatorIndex := sub(sload(keccak256(0x40, 0x40)), 1) // underflows when s_activatedOperatorIndex1Based is 0
            if gt(activatedOperatorIndex, MAX_OPERATOR_INDEX) {
                mstore(0, 0x1b256530) // NotActivatedOperator()
                revert(0x1c, 0x04)
            }
            let zeroBitIfSubmittedCvBitmap := sload(s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot)
            if gt(and(zeroBitIfSubmittedCvBitmap, shl(activatedOperatorIndex, 1)), 0) {
                // if bit is still set, meaning no Cv submitted for this operator
                // this operator was not requested to submit Co
                mstore(0, 0x03798920) // CvNotSubmitted()
                revert(0x1c, 0x04)
            }
            mstore(0x40, co)
            if iszero(eq(sload(add(s_cvs.slot, activatedOperatorIndex)), keccak256(0x40, 0x20))) {
                mstore(0, 0x67b3c693) // CvNotEqualHashCo()
                revert(0x1c, 0x04)
            }
            // ** bitmap
            sstore(
                s_zeroBitIfSubmittedCoBitmap.slot,
                and(sload(s_zeroBitIfSubmittedCoBitmap.slot), not(shl(activatedOperatorIndex, 1)))
            ) // set to zero bit
            sstore(add(s_cos.slot, activatedOperatorIndex), co)
            // ** event
            mstore(0x00, curRound) // 0x20 already has trialNum, 0x40 already has co
            mstore(0x60, activatedOperatorIndex)
            log1(0x00, 0x80, 0xc294138987faa6e0ebef350caeac5cf5e1eff8dbbe8a158e421601f48674babd) // event CoSubmitted(uint256 round, uint256 trialNum, bytes32 co, uint256 index)
        }
    }

    function requestToSubmitS(
        bytes32[] calldata allCos, // all cos
        bytes32[] calldata secretsReceivedOffchainInRevealOrder, // already received offchain
        uint256, // packedVsForAllCvsNotOnChain
        SigRS[] calldata sigRSsForAllCvsNotOnChain,
        uint256 packedRevealOrders
    ) external onlyOwner {
        bytes32 domainSeparator = _domainSeparatorV4();
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNum := sload(keccak256(0x40, 0x40))
            mstore(0x00, trialNum)
            // ** can only request to submit S if merkleRoot is submitted
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            if iszero(sload(keccak256(0x00, 0x40))) {
                mstore(0, 0x8e56b845) // MerkleRootNotSubmitted()
                revert(0x1c, 0x04)
            }
            // ** check if already requested to submit S
            mstore(0x60, s_previousSSubmitTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let previousSSubmitTimestampSlot := keccak256(0x00, 0x40)
            if gt(sload(previousSSubmitTimestampSlot), 0) {
                mstore(0, 0x0d934196) // AlreadyRequestedToSubmitS()
                revert(0x1c, 0x04)
            }
            // ** check allCos length
            let activatedOperatorsLength := sload(s_activatedOperators.slot)
            if iszero(eq(activatedOperatorsLength, allCos.length)) {
                mstore(0, 0x15467973) // AllCosNotSubmitted()
                revert(0x1c, 0x04)
            }
            // ** check cv status
            mstore(0x60, s_requestedToSubmitCvTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCvTimestampSlot := keccak256(0x00, 0x40)
            let zeroBitIfSubmittedCvBitmap
            switch sload(requestedToSubmitCvTimestampSlot)
            case 0 {
                sstore(requestedToSubmitCvTimestampSlot, 1) // set to 1 to indicate that cvs are on-chain
                zeroBitIfSubmittedCvBitmap := 0xffffffff // set all bits to 1
            }
            default {
                zeroBitIfSubmittedCvBitmap := sload(s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot)
            }

            // ****
            let cos := 0xc0
            let operatorLengthInBytes := mul(activatedOperatorsLength, 0x20)
            calldatacopy(cos, allCos.offset, operatorLengthInBytes) // allCos
            mstore(0x80, keccak256(cos, operatorLengthInBytes)) // rv
            let cvs := add(cos, operatorLengthInBytes) // cvs
            let di := add(cvs, operatorLengthInBytes) // diffs
            let fmp := add(di, operatorLengthInBytes) // fmp
            mstore(fmp, MESSAGE_TYPEHASH_DIRECT)
            mstore(add(fmp, 0x20), curRound)
            mstore(add(fmp, 0x40), trialNum)
            mstore(add(fmp, 0x80), hex"1901") // prefix and version
            mstore(add(fmp, 0x82), domainSeparator)
            let sigCounter
            for { let i } lt(i, activatedOperatorsLength) { i := add(i, 1) } {
                let cv := keccak256(add(cos, shl(5, i)), 0x20)
                mstore(add(cvs, shl(5, i)), cv) // cv
                mstore(0xa0, cv)
                mstore(add(di, shl(5, i)), keccak256(0x80, 0x40)) // hash(rv || cv)
                switch iszero(and(zeroBitIfSubmittedCvBitmap, shl(i, 1)))
                case 1 {
                    // cv is on-chain
                    if iszero(eq(sload(add(s_cvs.slot, i)), cv)) {
                        mstore(0, 0x67b3c693) // CvNotEqualHashCo()
                        revert(0x1c, 0x04)
                    }
                }
                default {
                    // cv is not on-chain
                    // ** check signature
                    let rSOffset := add(sigRSsForAllCvsNotOnChain.offset, shl(6, sigCounter))
                    let s := calldataload(add(rSOffset, 0x20))
                    if gt(s, SECP256K1_CURVE_ORDER) {
                        mstore(0, 0xbf4bf5b8) // InvalidSignatureS()
                        revert(0x1c, 0x04)
                    }
                    mstore(add(fmp, 0x60), cv)
                    mstore(add(fmp, 0xa2), keccak256(fmp, 0x80)) // structHash
                    mstore(0x00, keccak256(add(fmp, 0x80), 0x42)) // digest hash
                    mstore(0x20, and(calldataload(sub(0x44, sigCounter)), 0xff)) // v, 0x44: packedVsForAllCvsNotOnChain offset
                    sigCounter := add(sigCounter, 1)
                    mstore(0x40, calldataload(rSOffset)) // r
                    mstore(0x60, s)
                    let operatorAddress := mload(staticcall(gas(), 1, 0x00, 0x80, 0x01, 0x20))
                    // `returndatasize()` will be `0x20` upon success, and `0x00` otherwise.
                    if iszero(returndatasize()) {
                        mstore(0x00, 0x8baa579f) // selector for InvalidSignature()
                        revert(0x1c, 0x04)
                    }
                    mstore(0x00, operatorAddress)
                    mstore(0x20, s_activatedOperatorIndex1Based.slot)
                    if iszero(eq(add(i, 1), sload(keccak256(0x00, 0x40)))) {
                        mstore(0, 0x980c4296) // SignatureAndIndexDoNotMatch()
                        revert(0x1c, 0x04)
                    }
                    // ** submit cv on-chain
                    sstore(add(s_cvs.slot, i), cv) // cv
                    zeroBitIfSubmittedCvBitmap := and(zeroBitIfSubmittedCvBitmap, not(shl(i, 1))) // set to zero
                }
            }

            // ** verify reveal orders
            let index := and(packedRevealOrders, 0xff) // first reveal index
            let revealBitmap := shl(index, 1)
            let before := mload(add(di, shl(5, index)))
            for { let i := 1 } lt(i, activatedOperatorsLength) { i := add(i, 1) } {
                index := and(calldataload(sub(0x84, i)), 0xff) // 0x84: packedRevealOrders offset
                revealBitmap := or(revealBitmap, shl(index, 1))
                let after := mload(add(di, shl(5, index)))
                if lt(before, after) {
                    mstore(0, 0x24f1948e) // RevealNotInDescendingOrder()
                    revert(0x1c, 0x04)
                }
                before := after
            }
            if iszero(eq(revealBitmap, sub(shl(activatedOperatorsLength, 1), 1))) {
                mstore(0, 0x06efcba4) // selector for RevealOrderHasDuplicates()
                revert(0x1c, 0x04)
            }
            // ** Create Merkle Root and verify it
            let hashCountInBytes := sub(operatorLengthInBytes, 0x20)
            let cvsPosInBytes
            let hashPosInBytes
            for { let i } lt(i, hashCountInBytes) { i := add(i, 0x20) } {
                switch lt(cvsPosInBytes, operatorLengthInBytes)
                case 1 {
                    mstore(0x00, mload(add(cvs, cvsPosInBytes)))
                    cvsPosInBytes := add(cvsPosInBytes, 0x20)
                }
                default {
                    mstore(0x00, mload(add(fmp, hashPosInBytes)))
                    hashPosInBytes := add(hashPosInBytes, 0x20)
                }
                switch lt(cvsPosInBytes, operatorLengthInBytes)
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
            // ** verify the merkle root
            if iszero(eq(mload(add(fmp, sub(hashCountInBytes, 0x20))), sload(s_merkleRoot.slot))) {
                mstore(0, 0x624dc351) // selector for MerkleVerificationFailed()
                revert(0x1c, 0x04)
            }

            // skip updating zeroBitIfSubmittedCvBitmap because it is not used anymore
            sstore(s_packedRevealOrders.slot, packedRevealOrders) // update packedRevealOrders
            sstore(s_requestedToSubmitSFromIndexK.slot, secretsReceivedOffchainInRevealOrder.length)
            mstore(0x00, curRound)
            mstore(0x20, trialNum)
            mstore(0x40, secretsReceivedOffchainInRevealOrder.length)
            log1(0x00, 0x60, 0x583f939e9612a50da8a140b5e7247ff7c3c899c45e4051a5ba045abea6177f08) // event RequestedToSubmitSFromIndexK(uint256 round, uint256 trialNum, uint256 indexK)
            // ** store secrets
            for { let i } lt(i, secretsReceivedOffchainInRevealOrder.length) { i := add(i, 1) } {
                index := and(calldataload(sub(0x84, i)), 0xff) // 0x84: packedRevealOrders offset
                let secret := calldataload(add(secretsReceivedOffchainInRevealOrder.offset, shl(5, i)))
                mstore(0x00, secret)
                mstore(0x00, keccak256(0x00, 0x20)) // co
                if iszero(eq(mload(add(cvs, shl(5, index))), keccak256(0x00, 0x20))) {
                    mstore(0, 0x5bcc2334) // CvNotEqualDoubleHashS()
                    revert(0x1c, 0x04)
                }
                sstore(add(s_secrets.slot, index), secret) // store secret)
            }
            // Record the timestamp of the last S submission
            sstore(previousSSubmitTimestampSlot, timestamp())
        }
    }

    function submitS(bytes32 s) external {
        assembly ("memory-safe") {
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNum := sload(keccak256(0x40, 0x40))
            mstore(0x20, trialNum) // trialNum)
            // ** check if S was requested
            mstore(0x60, s_previousSSubmitTimestamp.slot)
            mstore(0x40, keccak256(0x40, 0x40))
            let previousSSubmitTimestampSlot := keccak256(0x20, 0x40)
            if iszero(sload(previousSSubmitTimestampSlot)) {
                mstore(0, 0x2d37f8d3) // SNotRequested()
                revert(0x1c, 0x04)
            }
            // ** check reveal order
            let fmp := 0x80 // cache fmp
            mstore(0x40, caller())
            mstore(0x60, s_activatedOperatorIndex1Based.slot)
            let activatedOperatorIndex := sub(sload(keccak256(0x40, 0x40)), 1) // underflows when s_activatedOperatorIndex1Based is 0
            if gt(activatedOperatorIndex, MAX_OPERATOR_INDEX) {
                mstore(0, 0x1b256530) // NotActivatedOperator()
                revert(0x1c, 0x04)
            }
            mstore(fmp, sload(s_packedRevealOrders.slot))
            let requestedToSubmitSFromIndexK := sload(s_requestedToSubmitSFromIndexK.slot)
            if iszero(eq(activatedOperatorIndex, and(mload(sub(fmp, requestedToSubmitSFromIndexK)), 0xff))) {
                mstore(0, 0xe3ae7cc0) // WrongRevealOrder()
                revert(0x1c, 0x04)
            }
            // ** check cv = doubleHashS
            mstore(0x40, s)
            mstore(0x60, keccak256(0x40, 0x20)) // co
            if iszero(eq(sload(add(s_cvs.slot, activatedOperatorIndex)), keccak256(0x60, 0x20))) {
                mstore(0, 0x5bcc2334) // CvNotEqualDoubleHashS()
                revert(0x1c, 0x04)
            }
            // ** store S and emit event
            mstore(0x00, curRound) // 0x20 already has trialNum, 0x40 already has s
            mstore(0x60, activatedOperatorIndex)
            log1(0x00, 0x80, 0xfa070a58e2c77080acd5c2b1819669eb194bbeeca6f680a31a2076510be5a7b1) // event SSubmitted(uint256 round, uint256 trialNum, bytes32 s, uint256 index)

            // ** If msg.sender is the last revealer, finalize the random number
            let activatedOperatorsLength := sload(s_activatedOperators.slot)
            switch eq(requestedToSubmitSFromIndexK, sub(activatedOperatorsLength, 1))
            case 1 {
                let storedSLength := sub(activatedOperatorsLength, 1)
                for { let i } lt(i, storedSLength) { i := add(i, 1) } {
                    mstore(add(fmp, shl(5, i)), sload(add(s_secrets.slot, i))) // store secrets, overwrites fmp because it is not used anymore
                }
                mstore(add(fmp, shl(5, storedSLength)), s) // last secret
                let randomNumber := keccak256(fmp, shl(5, activatedOperatorsLength))
                let nextRound := add(curRound, 1)
                let requestCount := sload(s_requestCount.slot)
                switch eq(nextRound, requestCount)
                case 1 {
                    if eq(sload(s_isInProcess.slot), COMPLETED) {
                        mstore(0x00, 0x195332a5) // selector for AlreadyCompleted()
                        revert(0x1c, 0x04)
                    }
                    sstore(s_isInProcess.slot, COMPLETED)
                    // 0x00 already has curRound, 0x20 already has trialNum
                    mstore(0x40, COMPLETED)
                    log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
                }
                default {
                    // get next round
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
                                            shr(
                                                250,
                                                mul(x, 0xb6db6db6ddddddddd34d34d349249249210842108c6318c639ce739cffffffff)
                                            )
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
                mstore(0x00, caller())
                mstore(0x20, s_depositAmount.slot)
                let depositSlot := keccak256(0x00, 0x40) // last revealer
                let flatFee := sload(s_flatFee.slot)
                sstore(depositSlot, add(sload(depositSlot), flatFee))
                // reward sload(add(currentRequestInfoSlot, 2)) - flatFee to the leader
                mstore(0x00, sload(_OWNER_SLOT))
                depositSlot := keccak256(0x00, 0x40) // leader

                mstore(0x00, 0x00fc98b8) // rawFulfillRandomNumber(uint256,uint256) selector
                mstore(0x20, curRound)
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
            default {
                sstore(add(s_secrets.slot, activatedOperatorIndex), s) // store secret
                sstore(s_requestedToSubmitSFromIndexK.slot, add(requestedToSubmitSFromIndexK, 1)) // increment index
            }
        }
    }

    function generateRandomNumberWhenSomeCvsAreOnChain(
        bytes32[] calldata allSecrets,
        SigRS[] calldata sigRSsForAllCvsNotOnChain,
        uint256, // packedVsForAllCvsNotOnChain
        uint256 packedRevealOrders
    ) external {
        bytes32 domainSeparator = _domainSeparatorV4();
        assembly ("memory-safe") {
            // ** check if some cvs are on-chain
            let curRound := sload(s_currentRound.slot)
            mstore(0x40, curRound)
            mstore(0x60, s_trialNum.slot)
            let trialNum := sload(keccak256(0x40, 0x40))
            mstore(0x00, trialNum) // trialNum)
            mstore(0x60, s_requestedToSubmitCvTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            let requestedToSubmitCvTimestampSlot := keccak256(0x00, 0x40)
            if iszero(sload(requestedToSubmitCvTimestampSlot)) {
                mstore(0, 0x96fbee7b) // NoCvsOnChain()
                revert(0x1c, 0x04)
            }
            mstore(0x60, s_merkleRootSubmittedTimestamp.slot)
            mstore(0x20, keccak256(0x40, 0x40))
            if iszero(sload(keccak256(0x00, 0x40))) {
                mstore(0, 0x8e56b845) // MerkleRootNotSubmitted()
                revert(0x1c, 0x04)
            }
            // ** initialize cos and cvs arrays memory, without length data
            let activatedOperatorsLength := sload(s_activatedOperators.slot)
            let activatedOperatorsLengthInBytes := shl(5, activatedOperatorsLength)

            let cos := 0x80
            let cvs := add(cos, activatedOperatorsLengthInBytes)
            let secrets := add(cvs, activatedOperatorsLengthInBytes)
            mstore(0x40, add(secrets, activatedOperatorsLengthInBytes)) // update the free memory pointer

            // ** get cos and cvs
            for { let i } lt(i, activatedOperatorsLengthInBytes) { i := add(i, 0x20) } {
                let secretMemP := add(secrets, i)
                mstore(secretMemP, calldataload(add(allSecrets.offset, i))) // secret
                let cosMemP := add(cos, i)
                mstore(cosMemP, keccak256(secretMemP, 0x20))
                mstore(add(cvs, i), keccak256(cosMemP, 0x20))
            }
            // ** verify reveal order
            mstore(0x00, keccak256(cos, activatedOperatorsLengthInBytes)) // rv
            let index := and(packedRevealOrders, 0xff) // first reveal index
            let revealBitmap := shl(index, 1)
            mstore(0x20, mload(add(cvs, shl(5, index))))
            let before := keccak256(0x00, 0x40)
            // revealOrdersOffset = 0x64
            for { let i := 1 } lt(i, activatedOperatorsLength) { i := add(i, 1) } {
                index := and(calldataload(sub(0x64, i)), 0xff)
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
            // ** verify the merkle root
            if iszero(eq(mload(add(fmp, sub(hashCountInBytes, 0x20))), sload(s_merkleRoot.slot))) {
                mstore(0, 0x624dc351) // selector for MerkleVerificationFailed()
                revert(0x1c, 0x04)
            }

            // ** verify signatures or cvs on-chain
            mstore(fmp, MESSAGE_TYPEHASH_DIRECT) // typehash, overwrite the previous value, which is not used anymore
            mstore(add(fmp, 0x20), curRound)
            mstore(add(fmp, 0x40), trialNum)
            mstore(add(fmp, 0x80), hex"1901") // prefix and version
            mstore(add(fmp, 0x82), domainSeparator)
            let zeroBitIfSubmittedCvBitmap := sload(s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2.slot)
            let sigCounter
            for { let i } lt(i, activatedOperatorsLengthInBytes) { i := add(i, 0x20) } {
                index := shr(5, i)
                switch iszero(and(zeroBitIfSubmittedCvBitmap, shl(index, 1)))
                case 1 {
                    // cv is on-chain
                    if iszero(eq(sload(add(s_cvs.slot, index)), mload(add(cvs, i)))) {
                        mstore(0, 0xa39ecadf) // selector for OnChainCvNotEqualDoubleHashS()
                        revert(0x1c, 0x04)
                    }
                }
                default {
                    // signature malleability prevention
                    let rOffset := add(sigRSsForAllCvsNotOnChain.offset, shl(6, sigCounter))
                    let s := calldataload(add(rOffset, 0x20))
                    if gt(s, SECP256K1_CURVE_ORDER) {
                        mstore(0, 0xbf4bf5b8) // selector for InvalidSignatureS()
                        revert(0x1c, 0x04)
                    }
                    mstore(add(fmp, 0x60), mload(add(cvs, i))) // cv
                    mstore(add(fmp, 0xa2), keccak256(fmp, 0x80)) // structHash
                    mstore(0x00, keccak256(add(fmp, 0x80), 0x42)) // digest hash
                    mstore(0x20, and(calldataload(sub(0x44, sigCounter)), 0xff)) // v, 0x44: packedVsOffset
                    sigCounter := add(sigCounter, 1)
                    mstore(0x40, calldataload(rOffset)) // r
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
            }

            // ** create random number
            let randomNumber := keccak256(secrets, activatedOperatorsLengthInBytes)
            let nextRound := add(curRound, 1)
            let requestCount := sload(s_requestCount.slot)
            switch eq(nextRound, requestCount)
            case 1 {
                // there is no next round
                if eq(sload(s_isInProcess.slot), COMPLETED) {
                    mstore(0x00, 0x195332a5) // selector for AlreadyCompleted()
                    revert(0x1c, 0x04)
                }
                sstore(s_isInProcess.slot, COMPLETED)
                mstore(0x00, curRound)
                mstore(0x20, trialNum)
                mstore(0x40, COMPLETED)
                log1(0x00, 0x60, 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813) // event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
            }
            default {
                // get next round
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
                        and(calldataload(sub(0x64, sub(activatedOperatorsLength, 1))), 0xff) // last revealer index, 0x64: revealOrdersOffset
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
            mstore(0x20, curRound)
            mstore(0x40, s_requestInfo.slot)
            let currentRequestInfoSlot := keccak256(0x20, 0x40)
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
}
