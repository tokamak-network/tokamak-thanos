// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {CommitReveal2} from "./CommitReveal2.sol";

contract CommitReveal2L2 is CommitReveal2 {
    /// @dev This is the padding size for unsigned RLP-encoded transaction without the signature data
    /// @dev The padding size was estimated based on hypothetical max RLP-encoded transaction size
    /// @dev Reference: https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/vrf/dev/OptimismL1Fees.sol
    uint256 internal constant L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE = 71;
    /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
    /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
    address internal constant OVM_GASPRICEORACLE_ADDR = 0x420000000000000000000000000000000000000F;

    /// @dev L1 fee coefficient is used to account for the impact of data compression on the l1 fee
    /// getL1FeeUpperBound returns the upper bound of l1 fee so this configurable coefficient will help
    /// charge a predefined percentage of the upper bound.
    uint8 public s_l1FeeCoefficient = 100;

    error InvalidL1FeeCoefficient(uint8 coefficient);

    event L1FeeCalculationSet(uint8 coefficient);

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
    )
        payable
        CommitReveal2(
            activationThreshold,
            flatFee,
            name,
            version,
            offChainSubmissionPeriod,
            requestOrSubmitOrFailDecisionPeriod,
            onChainSubmissionPeriod,
            offChainSubmissionPeriodPerOperator,
            onChainSubmissionPeriodPerOperator
        )
    {}

    function setL1FeeCoefficient(uint8 coefficient) external onlyOwner notInProcess {
        if (coefficient == 0 || coefficient > 100) {
            revert InvalidL1FeeCoefficient(coefficient);
        }
        s_l1FeeCoefficient = coefficient;
        emit L1FeeCalculationSet(coefficient);
    }

    /**
     * @notice Calculates the total fee required for requesting a random number, factoring in:
     *         1. L2 execution costs (based on the callback gas limit and the number of operators).
     *         2. A flat fee (s_flatFee).
     *         3. L1 data costs for sending required transaction calldata.
     * @dev
     *      - The internal gas usage estimation is derived from two operations:
     *        (i) "submitMerkleRoot" (~47,216 L2 gas).
     *        (ii) "generateRandomNumber" (~21,119 × numOfOperators + 134,334 L2 gas).
     * @param callbackGasLimit The gas required by the consumer’s callback execution.
     * @param gasPrice The L2 gas price to be used for cost estimation.
     * @param numOfOperators The number of active operators factored into the total gas cost.
     * @return requestFee The calculated total fee (in wei) needed to cover the request.
     */
    function _calculateRequestPrice(uint32 callbackGasLimit, uint256 gasPrice, uint256 numOfOperators)
        internal
        view
        override
        returns (uint256 requestFee)
    {
        assembly ("memory-safe") {
            mstore(0x00, 0xf1c7a58b) // selector for "getL1FeeUpperBound(uint256 _unsignedTxSize) external view returns (uint256)"
            mstore(0x20, add(MERKLEROOTSUB_CALLDATA_BYTES_SIZE, L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE))
            if iszero(staticcall(gas(), OVM_GASPRICEORACLE_ADDR, 0x1c, 0x24, 0x80, 0x20)) {
                mstore(0, 0xb75f34bf) // selector for L1FeeEstimationFailed()
                revert(0x1c, 0x04)
            }
            mstore(
                0x20,
                add(
                    add(mul(GENRANDNUM_CALLDATA_BYTES_SIZE_A, numOfOperators), GENRANDNUM_CALLDATA_BYTES_SIZE_B),
                    L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE
                )
            )
            if iszero(staticcall(gas(), OVM_GASPRICEORACLE_ADDR, 0x1c, 0x24, 0x20, 0x20)) {
                mstore(0, 0xb75f34bf) // selector for L1FeeEstimationFailed()
                revert(0x1c, 0x04)
            }
            let gasUsedMerkleRootSubAndGenRandNum := sload(s_gasUsedMerkleRootSubAndGenRandNumA.slot)
            requestFee :=
                add(
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
                    ), // l2GasFee
                    div(mul(sload(s_l1FeeCoefficient.slot), add(mload(0x20), mload(0x80))), 100) // L1GasFee
                )
        }
    }

    function _calculateFailGasFee(uint256 gasUsed) internal view override returns (uint256 gasFee) {
        assembly ("memory-safe") {
            mstore(0x00, 0xf1c7a58b) // selector for "getL1FeeUpperBound(uint256 _unsignedTxSize) external view returns (uint256)"
            mstore(0x20, add(FAIL_FUNCTIONS_CALLDATA_BYTES_SIZE, L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE))
            if iszero(staticcall(gas(), OVM_GASPRICEORACLE_ADDR, 0x1c, 0x24, 0x00, 0x20)) {
                mstore(0, 0xb75f34bf) // selector for L1FeeEstimationFailed()
                revert(0x1c, 0x04)
            }
            gasFee := add(mul(gasprice(), gasUsed), div(mul(sload(s_l1FeeCoefficient.slot), mload(0x00)), 100))
        }
    }

    function _getL1FeeUpperBoundOfFailFunction() internal view override returns (uint256 l1GasFee) {
        assembly ("memory-safe") {
            mstore(0x00, 0xf1c7a58b) // selector for "getL1FeeUpperBound(uint256 _unsignedTxSize) external view returns (uint256)"
            mstore(0x20, add(FAIL_FUNCTIONS_CALLDATA_BYTES_SIZE, L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE))
            if iszero(staticcall(gas(), OVM_GASPRICEORACLE_ADDR, 0x1c, 0x24, 0x00, 0x20)) {
                mstore(0, 0xb75f34bf) // selector for L1FeeEstimationFailed()
                revert(0x1c, 0x04)
            }
            l1GasFee := div(mul(sload(s_l1FeeCoefficient.slot), mload(0x00)), 100)
        }
    }

    function _getGetL1UpperBoundGasUsed() internal pure override returns (uint256) {
        return 0;
    }
}
