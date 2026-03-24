// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// CommitReveal2StorageCvAndSigRS is an auto generated low-level Go binding around an user-defined struct.
type CommitReveal2StorageCvAndSigRS struct {
	Cv [32]byte
	Rs CommitReveal2StorageSigRS
}

// CommitReveal2StorageSecretAndSigRS is an auto generated low-level Go binding around an user-defined struct.
type CommitReveal2StorageSecretAndSigRS struct {
	Secret [32]byte
	Rs     CommitReveal2StorageSigRS
}

// CommitReveal2StorageSigRS is an auto generated low-level Go binding around an user-defined struct.
type CommitReveal2StorageSigRS struct {
	R [32]byte
	S [32]byte
}

// CommitReveal2L2MetaData contains all meta data concerning the CommitReveal2L2 contract.
var CommitReveal2L2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"activationThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"flatFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"offChainSubmissionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestOrSubmitOrFailDecisionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"offChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"MAX_ACTIVATED_OPERATORS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SET_DELAY_TIME\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"claimSlashReward\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"completeOwnershipHandover\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deactivate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositAndActivate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"estimateRequestPrice\",\"inputs\":[{\"name\":\"callbackGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"estimateRequestPriceWithNumOfOperators\",\"inputs\":[{\"name\":\"callbackGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"numOfOperators\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeSetEconomicParameters\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeSetGasParameters\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"failToRequestSorGenerateRandomNumber\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"failToRequestSubmitCvOrSubmitMerkleRoot\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"failToSubmitCo\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"failToSubmitCv\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"failToSubmitMerkleRootAfterDispute\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"failToSubmitS\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"generateRandomNumber\",\"inputs\":[{\"name\":\"secretSigRSs\",\"type\":\"tuple[]\",\"internalType\":\"structCommitReveal2Storage.SecretAndSigRS[]\",\"components\":[{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rs\",\"type\":\"tuple\",\"internalType\":\"structCommitReveal2Storage.SigRS\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"packedRevealOrders\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"generateRandomNumberWhenSomeCvsAreOnChain\",\"inputs\":[{\"name\":\"allSecrets\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"sigRSsForAllCvsNotOnChain\",\"type\":\"tuple[]\",\"internalType\":\"structCommitReveal2Storage.SigRS[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"packedRevealOrders\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActivatedOperators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActivatedOperatorsLength\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurRoundAndStartTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurRoundAndTrialNum\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurStartTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDepositPlusSlashReward\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisputeInfos\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestedToSubmitCvTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestedToSubmitCvPackedIndicesAscFromLSB\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"zeroBitIfSubmittedCvBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestedToSubmitCoTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestedToSubmitCoPackedIndices\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestedToSubmitCoLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"zeroBitIfSubmittedCoBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"previousSSubmitTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"packedRevealOrders\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestedToSubmitSFromIndexK\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisputeTimestamps\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestedToSubmitCvTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestedToSubmitCoTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"previousSSubmitTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGasParameters\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMerkleRoot\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeriods\",\"inputs\":[],\"outputs\":[{\"name\":\"offChainSubmissionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestOrSubmitOrFailDecisionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"offChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSecrets\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"secrets\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getZeroBitIfSubmittedCoOnChainBitmap\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getZeroBitIfSubmittedCvOnChainBitmap\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"result\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownershipHandoverExpiresAt\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeEconomicParameters\",\"inputs\":[{\"name\":\"activationThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"flatFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proposeGasParameters\",\"inputs\":[{\"name\":\"gasUsedMerkleRootSubAndGenRandNumA\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"maxCallbackGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"getL1UpperBoundGasUsedWhenCalldataSize4\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"failToRequestCvOrSubmitMerkleRootGasUsed\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"failToSubmitMerkleRootAfterDisputeGasUsed\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"failToRequestSOrGenerateRandomNumberGasUsed\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"failToSubmitSGasUsed\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"failToSubmitCoGasUsedBaseA\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"failToSubmitCvGasUsedBaseA\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"failToSubmitGasUsedBaseB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"perOperatorIncreaseGasUsedA\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"perOperatorIncreaseGasUsedB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"perAdditionalDidntSubmitGasUsedA\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"perAdditionalDidntSubmitGasUsedB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"perRequestedIncreaseGasUsed\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"refund\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestRandomNumber\",\"inputs\":[{\"name\":\"callbackGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestToSubmitCo\",\"inputs\":[{\"name\":\"cvRSsForCvsNotOnChainAndReqToSubmitCo\",\"type\":\"tuple[]\",\"internalType\":\"structCommitReveal2Storage.CvAndSigRS[]\",\"components\":[{\"name\":\"cv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rs\",\"type\":\"tuple\",\"internalType\":\"structCommitReveal2Storage.SigRS\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"indicesLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"packedIndicesFirstCvNotOnChainRestCvOnChain\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestToSubmitCv\",\"inputs\":[{\"name\":\"packedIndicesAscendingFromLSB\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestToSubmitS\",\"inputs\":[{\"name\":\"allCos\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"secretsReceivedOffchainInRevealOrder\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sigRSsForAllCvsNotOnChain\",\"type\":\"tuple[]\",\"internalType\":\"structCommitReveal2Storage.SigRS[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"packedRevealOrders\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resume\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"s_activatedOperatorIndex1Based\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_activationThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_cos\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_currentRound\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_cvs\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_depositAmount\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_economicParamsEffectiveTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_flatFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_gasParamsEffectiveTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_isInProcess\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_l1FeeCoefficient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_merkleRootSubmittedTimestamp\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_packedRevealOrders\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_pendingActivationThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_pendingFlatFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_previousSSubmitTimestamp\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestInfo\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"consumer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestedToSubmitCoLength\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestedToSubmitCoPackedIndices\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestedToSubmitCoTimestamp\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestedToSubmitCvPackedIndicesAscFromLSB\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestedToSubmitCvTimestamp\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_requestedToSubmitSFromIndexK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_roundBitmap\",\"inputs\":[{\"name\":\"wordPos\",\"type\":\"uint248\",\"internalType\":\"uint248\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_secrets\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_slashRewardPerOperatorPaidX8\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_slashRewardPerOperatorX8\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_trialNum\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"trialNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_zeroBitIfSubmittedCoBitmap\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setL1FeeCoefficient\",\"inputs\":[{\"name\":\"coefficient\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeriods\",\"inputs\":[{\"name\":\"offChainSubmissionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestOrSubmitOrFailDecisionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"offChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitCo\",\"inputs\":[{\"name\":\"co\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitCv\",\"inputs\":[{\"name\":\"cv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitMerkleRoot\",\"inputs\":[{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Activated\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CoSubmitted\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"co\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CvSubmitted\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cv\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DeActivated\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EconomicParametersProposed\",\"inputs\":[{\"name\":\"activationThreshold\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"flatFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"effectiveTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EconomicParametersSet\",\"inputs\":[{\"name\":\"activationThreshold\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"flatFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GasParametersProposed\",\"inputs\":[{\"name\":\"gasUsedMerkleRootSubAndGenRandNumA\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"gasUsedMerkleRootSubAndGenRandNumB\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"maxCallbackGasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"getL1UpperBoundGasUsedWhenCalldataSize4\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToRequestCvOrSubmitMerkleRootGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToSubmitMerkleRootAfterDisputeGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToRequestSOrGenerateRandomNumberGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToSubmitSGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToSubmitCoGasUsedBaseA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"failToSubmitCvGasUsedBaseA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"failToSubmitGasUsedBaseB\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perOperatorIncreaseGasUsedA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perOperatorIncreaseGasUsedB\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perAdditionalDidntSubmitGasUsedA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perAdditionalDidntSubmitGasUsedB\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perRequestedIncreaseGasUsed\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"effectiveTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GasParametersSet\",\"inputs\":[{\"name\":\"gasUsedMerkleRootSubAndGenRandNumA\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"gasUsedMerkleRootSubAndGenRandNumB\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"maxCallbackGasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"getL1UpperBoundGasUsedWhenCalldataSize4\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToRequestCvOrSubmitMerkleRootGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToSubmitMerkleRootAfterDisputeGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToRequestSOrGenerateRandomNumberGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToSubmitSGasUsed\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"failToSubmitCoGasUsedBaseA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"failToSubmitCvGasUsedBaseA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"failToSubmitGasUsedBaseB\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perOperatorIncreaseGasUsedA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perOperatorIncreaseGasUsedB\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perAdditionalDidntSubmitGasUsedA\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perAdditionalDidntSubmitGasUsedB\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"perRequestedIncreaseGasUsed\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"L1FeeCalculationSet\",\"inputs\":[{\"name\":\"coefficient\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MerkleRootSubmitted\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverCanceled\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverRequested\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"oldOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodsSet\",\"inputs\":[{\"name\":\"offChainSubmissionPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"requestOrSubmitOrFailDecisionPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"offChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"onChainSubmissionPeriodPerOperator\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestedToSubmitCo\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"indicesLength\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"packedIndices\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestedToSubmitCv\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"packedIndicesAscendingFromLSB\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestedToSubmitSFromIndexK\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"indexK\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SSubmitted\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"s\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Status\",\"inputs\":[{\"name\":\"curRound\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"curTrialNum\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"curState\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ActivatedOperatorsLimitReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AllCosNotSubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AllCvsNotSubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AllSubmittedCo\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AllSubmittedCv\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyActivated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyCompleted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyHalted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyRefunded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyRequestedToSubmitCo\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyRequestedToSubmitCv\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyRequestedToSubmitS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadySubmittedMerkleRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadySubmittedS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotRequestWhenHalted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CoNotRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CvNotEqualDoubleHashS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CvNotEqualHashCo\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CvNotRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CvNotRequestedForThisOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CvNotSubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateIndices\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ETHTransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedCallbackGasLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InProcess\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidIndex\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidL1FeeCoefficient\",\"inputs\":[{\"name\":\"coefficient\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidSecretLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L1FeeEstimationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LeaderLowDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthExceedsMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LessThanActivationThreshold\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MerkleRootAlreadySubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MerkleRootIsSubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MerkleRootNotSubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MerkleVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewOwnerCannotBeActivatedOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewOwnerIsZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoCvsOnChain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoHandoverRequest\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentRound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotActivatedOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotConsumer\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughActivatedOperators\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotHalted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnChainCvNotEqualDoubleHashS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActivatedOperatorCanClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotActivate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PendingOwnerCannotBeActivatedOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RevealNotInDescendingOrder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RevealOrderHasDuplicates\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundAlreadyProcessed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundNotInProgress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SNotRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureAndIndexDoNotMatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"SubmitAfterStartTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooEarly\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooLate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooManyRequestsQueued\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawAmountIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongRevealOrder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroLength\",\"inputs\":[]}]",
}

// CommitReveal2L2ABI is the input ABI used to generate the binding from.
// Deprecated: Use CommitReveal2L2MetaData.ABI instead.
var CommitReveal2L2ABI = CommitReveal2L2MetaData.ABI

// CommitReveal2L2 is an auto generated Go binding around an Ethereum contract.
type CommitReveal2L2 struct {
	CommitReveal2L2Caller     // Read-only binding to the contract
	CommitReveal2L2Transactor // Write-only binding to the contract
	CommitReveal2L2Filterer   // Log filterer for contract events
}

// CommitReveal2L2Caller is an auto generated read-only Go binding around an Ethereum contract.
type CommitReveal2L2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommitReveal2L2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type CommitReveal2L2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommitReveal2L2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CommitReveal2L2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommitReveal2L2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CommitReveal2L2Session struct {
	Contract     *CommitReveal2L2  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CommitReveal2L2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CommitReveal2L2CallerSession struct {
	Contract *CommitReveal2L2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CommitReveal2L2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CommitReveal2L2TransactorSession struct {
	Contract     *CommitReveal2L2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CommitReveal2L2Raw is an auto generated low-level Go binding around an Ethereum contract.
type CommitReveal2L2Raw struct {
	Contract *CommitReveal2L2 // Generic contract binding to access the raw methods on
}

// CommitReveal2L2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CommitReveal2L2CallerRaw struct {
	Contract *CommitReveal2L2Caller // Generic read-only contract binding to access the raw methods on
}

// CommitReveal2L2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CommitReveal2L2TransactorRaw struct {
	Contract *CommitReveal2L2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewCommitReveal2L2 creates a new instance of CommitReveal2L2, bound to a specific deployed contract.
func NewCommitReveal2L2(address common.Address, backend bind.ContractBackend) (*CommitReveal2L2, error) {
	contract, err := bindCommitReveal2L2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2{CommitReveal2L2Caller: CommitReveal2L2Caller{contract: contract}, CommitReveal2L2Transactor: CommitReveal2L2Transactor{contract: contract}, CommitReveal2L2Filterer: CommitReveal2L2Filterer{contract: contract}}, nil
}

// NewCommitReveal2L2Caller creates a new read-only instance of CommitReveal2L2, bound to a specific deployed contract.
func NewCommitReveal2L2Caller(address common.Address, caller bind.ContractCaller) (*CommitReveal2L2Caller, error) {
	contract, err := bindCommitReveal2L2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2Caller{contract: contract}, nil
}

// NewCommitReveal2L2Transactor creates a new write-only instance of CommitReveal2L2, bound to a specific deployed contract.
func NewCommitReveal2L2Transactor(address common.Address, transactor bind.ContractTransactor) (*CommitReveal2L2Transactor, error) {
	contract, err := bindCommitReveal2L2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2Transactor{contract: contract}, nil
}

// NewCommitReveal2L2Filterer creates a new log filterer instance of CommitReveal2L2, bound to a specific deployed contract.
func NewCommitReveal2L2Filterer(address common.Address, filterer bind.ContractFilterer) (*CommitReveal2L2Filterer, error) {
	contract, err := bindCommitReveal2L2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2Filterer{contract: contract}, nil
}

// bindCommitReveal2L2 binds a generic wrapper to an already deployed contract.
func bindCommitReveal2L2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CommitReveal2L2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CommitReveal2L2 *CommitReveal2L2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CommitReveal2L2.Contract.CommitReveal2L2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CommitReveal2L2 *CommitReveal2L2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.CommitReveal2L2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CommitReveal2L2 *CommitReveal2L2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.CommitReveal2L2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CommitReveal2L2 *CommitReveal2L2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CommitReveal2L2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CommitReveal2L2 *CommitReveal2L2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CommitReveal2L2 *CommitReveal2L2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.contract.Transact(opts, method, params...)
}

// MAXACTIVATEDOPERATORS is a free data retrieval call binding the contract method 0xd734f455.
//
// Solidity: function MAX_ACTIVATED_OPERATORS() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) MAXACTIVATEDOPERATORS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "MAX_ACTIVATED_OPERATORS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXACTIVATEDOPERATORS is a free data retrieval call binding the contract method 0xd734f455.
//
// Solidity: function MAX_ACTIVATED_OPERATORS() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) MAXACTIVATEDOPERATORS() (*big.Int, error) {
	return _CommitReveal2L2.Contract.MAXACTIVATEDOPERATORS(&_CommitReveal2L2.CallOpts)
}

// MAXACTIVATEDOPERATORS is a free data retrieval call binding the contract method 0xd734f455.
//
// Solidity: function MAX_ACTIVATED_OPERATORS() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) MAXACTIVATEDOPERATORS() (*big.Int, error) {
	return _CommitReveal2L2.Contract.MAXACTIVATEDOPERATORS(&_CommitReveal2L2.CallOpts)
}

// SETDELAYTIME is a free data retrieval call binding the contract method 0x46beb81a.
//
// Solidity: function SET_DELAY_TIME() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SETDELAYTIME(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "SET_DELAY_TIME")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SETDELAYTIME is a free data retrieval call binding the contract method 0x46beb81a.
//
// Solidity: function SET_DELAY_TIME() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SETDELAYTIME() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SETDELAYTIME(&_CommitReveal2L2.CallOpts)
}

// SETDELAYTIME is a free data retrieval call binding the contract method 0x46beb81a.
//
// Solidity: function SET_DELAY_TIME() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SETDELAYTIME() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SETDELAYTIME(&_CommitReveal2L2.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_CommitReveal2L2 *CommitReveal2L2Caller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_CommitReveal2L2 *CommitReveal2L2Session) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _CommitReveal2L2.Contract.Eip712Domain(&_CommitReveal2L2.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _CommitReveal2L2.Contract.Eip712Domain(&_CommitReveal2L2.CallOpts)
}

// EstimateRequestPrice is a free data retrieval call binding the contract method 0x7fb5d19d.
//
// Solidity: function estimateRequestPrice(uint32 callbackGasLimit, uint256 gasPrice) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) EstimateRequestPrice(opts *bind.CallOpts, callbackGasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "estimateRequestPrice", callbackGasLimit, gasPrice)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateRequestPrice is a free data retrieval call binding the contract method 0x7fb5d19d.
//
// Solidity: function estimateRequestPrice(uint32 callbackGasLimit, uint256 gasPrice) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) EstimateRequestPrice(callbackGasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.EstimateRequestPrice(&_CommitReveal2L2.CallOpts, callbackGasLimit, gasPrice)
}

// EstimateRequestPrice is a free data retrieval call binding the contract method 0x7fb5d19d.
//
// Solidity: function estimateRequestPrice(uint32 callbackGasLimit, uint256 gasPrice) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) EstimateRequestPrice(callbackGasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.EstimateRequestPrice(&_CommitReveal2L2.CallOpts, callbackGasLimit, gasPrice)
}

// EstimateRequestPriceWithNumOfOperators is a free data retrieval call binding the contract method 0x5d4ccbee.
//
// Solidity: function estimateRequestPriceWithNumOfOperators(uint32 callbackGasLimit, uint256 gasPrice, uint256 numOfOperators) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) EstimateRequestPriceWithNumOfOperators(opts *bind.CallOpts, callbackGasLimit uint32, gasPrice *big.Int, numOfOperators *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "estimateRequestPriceWithNumOfOperators", callbackGasLimit, gasPrice, numOfOperators)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateRequestPriceWithNumOfOperators is a free data retrieval call binding the contract method 0x5d4ccbee.
//
// Solidity: function estimateRequestPriceWithNumOfOperators(uint32 callbackGasLimit, uint256 gasPrice, uint256 numOfOperators) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) EstimateRequestPriceWithNumOfOperators(callbackGasLimit uint32, gasPrice *big.Int, numOfOperators *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.EstimateRequestPriceWithNumOfOperators(&_CommitReveal2L2.CallOpts, callbackGasLimit, gasPrice, numOfOperators)
}

// EstimateRequestPriceWithNumOfOperators is a free data retrieval call binding the contract method 0x5d4ccbee.
//
// Solidity: function estimateRequestPriceWithNumOfOperators(uint32 callbackGasLimit, uint256 gasPrice, uint256 numOfOperators) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) EstimateRequestPriceWithNumOfOperators(callbackGasLimit uint32, gasPrice *big.Int, numOfOperators *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.EstimateRequestPriceWithNumOfOperators(&_CommitReveal2L2.CallOpts, callbackGasLimit, gasPrice, numOfOperators)
}

// GetActivatedOperators is a free data retrieval call binding the contract method 0xecd21a7e.
//
// Solidity: function getActivatedOperators() view returns(address[])
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetActivatedOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getActivatedOperators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetActivatedOperators is a free data retrieval call binding the contract method 0xecd21a7e.
//
// Solidity: function getActivatedOperators() view returns(address[])
func (_CommitReveal2L2 *CommitReveal2L2Session) GetActivatedOperators() ([]common.Address, error) {
	return _CommitReveal2L2.Contract.GetActivatedOperators(&_CommitReveal2L2.CallOpts)
}

// GetActivatedOperators is a free data retrieval call binding the contract method 0xecd21a7e.
//
// Solidity: function getActivatedOperators() view returns(address[])
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetActivatedOperators() ([]common.Address, error) {
	return _CommitReveal2L2.Contract.GetActivatedOperators(&_CommitReveal2L2.CallOpts)
}

// GetActivatedOperatorsLength is a free data retrieval call binding the contract method 0x36088f52.
//
// Solidity: function getActivatedOperatorsLength() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetActivatedOperatorsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getActivatedOperatorsLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActivatedOperatorsLength is a free data retrieval call binding the contract method 0x36088f52.
//
// Solidity: function getActivatedOperatorsLength() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetActivatedOperatorsLength() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetActivatedOperatorsLength(&_CommitReveal2L2.CallOpts)
}

// GetActivatedOperatorsLength is a free data retrieval call binding the contract method 0x36088f52.
//
// Solidity: function getActivatedOperatorsLength() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetActivatedOperatorsLength() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetActivatedOperatorsLength(&_CommitReveal2L2.CallOpts)
}

// GetCurRoundAndStartTime is a free data retrieval call binding the contract method 0xdeb05f1d.
//
// Solidity: function getCurRoundAndStartTime() view returns(uint256, uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetCurRoundAndStartTime(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getCurRoundAndStartTime")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetCurRoundAndStartTime is a free data retrieval call binding the contract method 0xdeb05f1d.
//
// Solidity: function getCurRoundAndStartTime() view returns(uint256, uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetCurRoundAndStartTime() (*big.Int, *big.Int, error) {
	return _CommitReveal2L2.Contract.GetCurRoundAndStartTime(&_CommitReveal2L2.CallOpts)
}

// GetCurRoundAndStartTime is a free data retrieval call binding the contract method 0xdeb05f1d.
//
// Solidity: function getCurRoundAndStartTime() view returns(uint256, uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetCurRoundAndStartTime() (*big.Int, *big.Int, error) {
	return _CommitReveal2L2.Contract.GetCurRoundAndStartTime(&_CommitReveal2L2.CallOpts)
}

// GetCurRoundAndTrialNum is a free data retrieval call binding the contract method 0x8b683dc6.
//
// Solidity: function getCurRoundAndTrialNum() view returns(uint256, uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetCurRoundAndTrialNum(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getCurRoundAndTrialNum")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetCurRoundAndTrialNum is a free data retrieval call binding the contract method 0x8b683dc6.
//
// Solidity: function getCurRoundAndTrialNum() view returns(uint256, uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetCurRoundAndTrialNum() (*big.Int, *big.Int, error) {
	return _CommitReveal2L2.Contract.GetCurRoundAndTrialNum(&_CommitReveal2L2.CallOpts)
}

// GetCurRoundAndTrialNum is a free data retrieval call binding the contract method 0x8b683dc6.
//
// Solidity: function getCurRoundAndTrialNum() view returns(uint256, uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetCurRoundAndTrialNum() (*big.Int, *big.Int, error) {
	return _CommitReveal2L2.Contract.GetCurRoundAndTrialNum(&_CommitReveal2L2.CallOpts)
}

// GetCurStartTime is a free data retrieval call binding the contract method 0xd1c36c5e.
//
// Solidity: function getCurStartTime() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetCurStartTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getCurStartTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurStartTime is a free data retrieval call binding the contract method 0xd1c36c5e.
//
// Solidity: function getCurStartTime() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetCurStartTime() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetCurStartTime(&_CommitReveal2L2.CallOpts)
}

// GetCurStartTime is a free data retrieval call binding the contract method 0xd1c36c5e.
//
// Solidity: function getCurStartTime() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetCurStartTime() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetCurStartTime(&_CommitReveal2L2.CallOpts)
}

// GetDepositPlusSlashReward is a free data retrieval call binding the contract method 0x1226f272.
//
// Solidity: function getDepositPlusSlashReward(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetDepositPlusSlashReward(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getDepositPlusSlashReward", operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositPlusSlashReward is a free data retrieval call binding the contract method 0x1226f272.
//
// Solidity: function getDepositPlusSlashReward(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetDepositPlusSlashReward(operator common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetDepositPlusSlashReward(&_CommitReveal2L2.CallOpts, operator)
}

// GetDepositPlusSlashReward is a free data retrieval call binding the contract method 0x1226f272.
//
// Solidity: function getDepositPlusSlashReward(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetDepositPlusSlashReward(operator common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetDepositPlusSlashReward(&_CommitReveal2L2.CallOpts, operator)
}

// GetDisputeInfos is a free data retrieval call binding the contract method 0x38ea288d.
//
// Solidity: function getDisputeInfos(uint256 round, uint256 trialNum) view returns(uint256 requestedToSubmitCvTimestamp, uint256 requestedToSubmitCvPackedIndicesAscFromLSB, uint256 zeroBitIfSubmittedCvBitmap, uint256 requestedToSubmitCoTimestamp, uint256 requestedToSubmitCoPackedIndices, uint256 requestedToSubmitCoLength, uint256 zeroBitIfSubmittedCoBitmap, uint256 previousSSubmitTimestamp, uint256 packedRevealOrders, uint256 requestedToSubmitSFromIndexK)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetDisputeInfos(opts *bind.CallOpts, round *big.Int, trialNum *big.Int) (struct {
	RequestedToSubmitCvTimestamp               *big.Int
	RequestedToSubmitCvPackedIndicesAscFromLSB *big.Int
	ZeroBitIfSubmittedCvBitmap                 *big.Int
	RequestedToSubmitCoTimestamp               *big.Int
	RequestedToSubmitCoPackedIndices           *big.Int
	RequestedToSubmitCoLength                  *big.Int
	ZeroBitIfSubmittedCoBitmap                 *big.Int
	PreviousSSubmitTimestamp                   *big.Int
	PackedRevealOrders                         *big.Int
	RequestedToSubmitSFromIndexK               *big.Int
}, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getDisputeInfos", round, trialNum)

	outstruct := new(struct {
		RequestedToSubmitCvTimestamp               *big.Int
		RequestedToSubmitCvPackedIndicesAscFromLSB *big.Int
		ZeroBitIfSubmittedCvBitmap                 *big.Int
		RequestedToSubmitCoTimestamp               *big.Int
		RequestedToSubmitCoPackedIndices           *big.Int
		RequestedToSubmitCoLength                  *big.Int
		ZeroBitIfSubmittedCoBitmap                 *big.Int
		PreviousSSubmitTimestamp                   *big.Int
		PackedRevealOrders                         *big.Int
		RequestedToSubmitSFromIndexK               *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestedToSubmitCvTimestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RequestedToSubmitCvPackedIndicesAscFromLSB = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ZeroBitIfSubmittedCvBitmap = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.RequestedToSubmitCoTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RequestedToSubmitCoPackedIndices = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.RequestedToSubmitCoLength = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ZeroBitIfSubmittedCoBitmap = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.PreviousSSubmitTimestamp = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.PackedRevealOrders = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.RequestedToSubmitSFromIndexK = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetDisputeInfos is a free data retrieval call binding the contract method 0x38ea288d.
//
// Solidity: function getDisputeInfos(uint256 round, uint256 trialNum) view returns(uint256 requestedToSubmitCvTimestamp, uint256 requestedToSubmitCvPackedIndicesAscFromLSB, uint256 zeroBitIfSubmittedCvBitmap, uint256 requestedToSubmitCoTimestamp, uint256 requestedToSubmitCoPackedIndices, uint256 requestedToSubmitCoLength, uint256 zeroBitIfSubmittedCoBitmap, uint256 previousSSubmitTimestamp, uint256 packedRevealOrders, uint256 requestedToSubmitSFromIndexK)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetDisputeInfos(round *big.Int, trialNum *big.Int) (struct {
	RequestedToSubmitCvTimestamp               *big.Int
	RequestedToSubmitCvPackedIndicesAscFromLSB *big.Int
	ZeroBitIfSubmittedCvBitmap                 *big.Int
	RequestedToSubmitCoTimestamp               *big.Int
	RequestedToSubmitCoPackedIndices           *big.Int
	RequestedToSubmitCoLength                  *big.Int
	ZeroBitIfSubmittedCoBitmap                 *big.Int
	PreviousSSubmitTimestamp                   *big.Int
	PackedRevealOrders                         *big.Int
	RequestedToSubmitSFromIndexK               *big.Int
}, error) {
	return _CommitReveal2L2.Contract.GetDisputeInfos(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// GetDisputeInfos is a free data retrieval call binding the contract method 0x38ea288d.
//
// Solidity: function getDisputeInfos(uint256 round, uint256 trialNum) view returns(uint256 requestedToSubmitCvTimestamp, uint256 requestedToSubmitCvPackedIndicesAscFromLSB, uint256 zeroBitIfSubmittedCvBitmap, uint256 requestedToSubmitCoTimestamp, uint256 requestedToSubmitCoPackedIndices, uint256 requestedToSubmitCoLength, uint256 zeroBitIfSubmittedCoBitmap, uint256 previousSSubmitTimestamp, uint256 packedRevealOrders, uint256 requestedToSubmitSFromIndexK)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetDisputeInfos(round *big.Int, trialNum *big.Int) (struct {
	RequestedToSubmitCvTimestamp               *big.Int
	RequestedToSubmitCvPackedIndicesAscFromLSB *big.Int
	ZeroBitIfSubmittedCvBitmap                 *big.Int
	RequestedToSubmitCoTimestamp               *big.Int
	RequestedToSubmitCoPackedIndices           *big.Int
	RequestedToSubmitCoLength                  *big.Int
	ZeroBitIfSubmittedCoBitmap                 *big.Int
	PreviousSSubmitTimestamp                   *big.Int
	PackedRevealOrders                         *big.Int
	RequestedToSubmitSFromIndexK               *big.Int
}, error) {
	return _CommitReveal2L2.Contract.GetDisputeInfos(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// GetDisputeTimestamps is a free data retrieval call binding the contract method 0xfe65f39d.
//
// Solidity: function getDisputeTimestamps(uint256 round, uint256 trialNum) view returns(uint256 requestedToSubmitCvTimestamp, uint256 requestedToSubmitCoTimestamp, uint256 previousSSubmitTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetDisputeTimestamps(opts *bind.CallOpts, round *big.Int, trialNum *big.Int) (struct {
	RequestedToSubmitCvTimestamp *big.Int
	RequestedToSubmitCoTimestamp *big.Int
	PreviousSSubmitTimestamp     *big.Int
}, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getDisputeTimestamps", round, trialNum)

	outstruct := new(struct {
		RequestedToSubmitCvTimestamp *big.Int
		RequestedToSubmitCoTimestamp *big.Int
		PreviousSSubmitTimestamp     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestedToSubmitCvTimestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RequestedToSubmitCoTimestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PreviousSSubmitTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetDisputeTimestamps is a free data retrieval call binding the contract method 0xfe65f39d.
//
// Solidity: function getDisputeTimestamps(uint256 round, uint256 trialNum) view returns(uint256 requestedToSubmitCvTimestamp, uint256 requestedToSubmitCoTimestamp, uint256 previousSSubmitTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetDisputeTimestamps(round *big.Int, trialNum *big.Int) (struct {
	RequestedToSubmitCvTimestamp *big.Int
	RequestedToSubmitCoTimestamp *big.Int
	PreviousSSubmitTimestamp     *big.Int
}, error) {
	return _CommitReveal2L2.Contract.GetDisputeTimestamps(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// GetDisputeTimestamps is a free data retrieval call binding the contract method 0xfe65f39d.
//
// Solidity: function getDisputeTimestamps(uint256 round, uint256 trialNum) view returns(uint256 requestedToSubmitCvTimestamp, uint256 requestedToSubmitCoTimestamp, uint256 previousSSubmitTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetDisputeTimestamps(round *big.Int, trialNum *big.Int) (struct {
	RequestedToSubmitCvTimestamp *big.Int
	RequestedToSubmitCoTimestamp *big.Int
	PreviousSSubmitTimestamp     *big.Int
}, error) {
	return _CommitReveal2L2.Contract.GetDisputeTimestamps(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// GetGasParameters is a free data retrieval call binding the contract method 0x1d6dc68b.
//
// Solidity: function getGasParameters() view returns(uint128, uint128, uint256, uint48, uint48, uint48, uint48, uint48, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetGasParameters(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getGasParameters")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(uint32), *new(uint32), *new(uint32), *new(uint32), *new(uint32), *new(uint32), *new(uint32), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	out5 := *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	out6 := *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	out7 := *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	out8 := *abi.ConvertType(out[8], new(uint32)).(*uint32)
	out9 := *abi.ConvertType(out[9], new(uint32)).(*uint32)
	out10 := *abi.ConvertType(out[10], new(uint32)).(*uint32)
	out11 := *abi.ConvertType(out[11], new(uint32)).(*uint32)
	out12 := *abi.ConvertType(out[12], new(uint32)).(*uint32)
	out13 := *abi.ConvertType(out[13], new(uint32)).(*uint32)
	out14 := *abi.ConvertType(out[14], new(uint32)).(*uint32)
	out15 := *abi.ConvertType(out[15], new(uint32)).(*uint32)

	return out0, out1, out2, out3, out4, out5, out6, out7, out8, out9, out10, out11, out12, out13, out14, out15, err

}

// GetGasParameters is a free data retrieval call binding the contract method 0x1d6dc68b.
//
// Solidity: function getGasParameters() view returns(uint128, uint128, uint256, uint48, uint48, uint48, uint48, uint48, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetGasParameters() (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32, error) {
	return _CommitReveal2L2.Contract.GetGasParameters(&_CommitReveal2L2.CallOpts)
}

// GetGasParameters is a free data retrieval call binding the contract method 0x1d6dc68b.
//
// Solidity: function getGasParameters() view returns(uint128, uint128, uint256, uint48, uint48, uint48, uint48, uint48, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetGasParameters() (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32, error) {
	return _CommitReveal2L2.Contract.GetGasParameters(&_CommitReveal2L2.CallOpts)
}

// GetMerkleRoot is a free data retrieval call binding the contract method 0xd1f2b5e8.
//
// Solidity: function getMerkleRoot(uint256 round, uint256 trialNum) view returns(bytes32, bool)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetMerkleRoot(opts *bind.CallOpts, round *big.Int, trialNum *big.Int) ([32]byte, bool, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getMerkleRoot", round, trialNum)

	if err != nil {
		return *new([32]byte), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(bool)).(*bool)

	return out0, out1, err

}

// GetMerkleRoot is a free data retrieval call binding the contract method 0xd1f2b5e8.
//
// Solidity: function getMerkleRoot(uint256 round, uint256 trialNum) view returns(bytes32, bool)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetMerkleRoot(round *big.Int, trialNum *big.Int) ([32]byte, bool, error) {
	return _CommitReveal2L2.Contract.GetMerkleRoot(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// GetMerkleRoot is a free data retrieval call binding the contract method 0xd1f2b5e8.
//
// Solidity: function getMerkleRoot(uint256 round, uint256 trialNum) view returns(bytes32, bool)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetMerkleRoot(round *big.Int, trialNum *big.Int) ([32]byte, bool, error) {
	return _CommitReveal2L2.Contract.GetMerkleRoot(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// GetPeriods is a free data retrieval call binding the contract method 0x45c9a558.
//
// Solidity: function getPeriods() view returns(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetPeriods(opts *bind.CallOpts) (struct {
	OffChainSubmissionPeriod            *big.Int
	RequestOrSubmitOrFailDecisionPeriod *big.Int
	OnChainSubmissionPeriod             *big.Int
	OffChainSubmissionPeriodPerOperator *big.Int
	OnChainSubmissionPeriodPerOperator  *big.Int
}, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getPeriods")

	outstruct := new(struct {
		OffChainSubmissionPeriod            *big.Int
		RequestOrSubmitOrFailDecisionPeriod *big.Int
		OnChainSubmissionPeriod             *big.Int
		OffChainSubmissionPeriodPerOperator *big.Int
		OnChainSubmissionPeriodPerOperator  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.OffChainSubmissionPeriod = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RequestOrSubmitOrFailDecisionPeriod = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.OnChainSubmissionPeriod = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.OffChainSubmissionPeriodPerOperator = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.OnChainSubmissionPeriodPerOperator = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPeriods is a free data retrieval call binding the contract method 0x45c9a558.
//
// Solidity: function getPeriods() view returns(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetPeriods() (struct {
	OffChainSubmissionPeriod            *big.Int
	RequestOrSubmitOrFailDecisionPeriod *big.Int
	OnChainSubmissionPeriod             *big.Int
	OffChainSubmissionPeriodPerOperator *big.Int
	OnChainSubmissionPeriodPerOperator  *big.Int
}, error) {
	return _CommitReveal2L2.Contract.GetPeriods(&_CommitReveal2L2.CallOpts)
}

// GetPeriods is a free data retrieval call binding the contract method 0x45c9a558.
//
// Solidity: function getPeriods() view returns(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetPeriods() (struct {
	OffChainSubmissionPeriod            *big.Int
	RequestOrSubmitOrFailDecisionPeriod *big.Int
	OnChainSubmissionPeriod             *big.Int
	OffChainSubmissionPeriodPerOperator *big.Int
	OnChainSubmissionPeriodPerOperator  *big.Int
}, error) {
	return _CommitReveal2L2.Contract.GetPeriods(&_CommitReveal2L2.CallOpts)
}

// GetSecrets is a free data retrieval call binding the contract method 0x1f9b7351.
//
// Solidity: function getSecrets(uint256 length) view returns(bytes32[] secrets)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetSecrets(opts *bind.CallOpts, length *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getSecrets", length)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetSecrets is a free data retrieval call binding the contract method 0x1f9b7351.
//
// Solidity: function getSecrets(uint256 length) view returns(bytes32[] secrets)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetSecrets(length *big.Int) ([][32]byte, error) {
	return _CommitReveal2L2.Contract.GetSecrets(&_CommitReveal2L2.CallOpts, length)
}

// GetSecrets is a free data retrieval call binding the contract method 0x1f9b7351.
//
// Solidity: function getSecrets(uint256 length) view returns(bytes32[] secrets)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetSecrets(length *big.Int) ([][32]byte, error) {
	return _CommitReveal2L2.Contract.GetSecrets(&_CommitReveal2L2.CallOpts, length)
}

// GetZeroBitIfSubmittedCoOnChainBitmap is a free data retrieval call binding the contract method 0x52cfd6c1.
//
// Solidity: function getZeroBitIfSubmittedCoOnChainBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetZeroBitIfSubmittedCoOnChainBitmap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getZeroBitIfSubmittedCoOnChainBitmap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetZeroBitIfSubmittedCoOnChainBitmap is a free data retrieval call binding the contract method 0x52cfd6c1.
//
// Solidity: function getZeroBitIfSubmittedCoOnChainBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetZeroBitIfSubmittedCoOnChainBitmap() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetZeroBitIfSubmittedCoOnChainBitmap(&_CommitReveal2L2.CallOpts)
}

// GetZeroBitIfSubmittedCoOnChainBitmap is a free data retrieval call binding the contract method 0x52cfd6c1.
//
// Solidity: function getZeroBitIfSubmittedCoOnChainBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetZeroBitIfSubmittedCoOnChainBitmap() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetZeroBitIfSubmittedCoOnChainBitmap(&_CommitReveal2L2.CallOpts)
}

// GetZeroBitIfSubmittedCvOnChainBitmap is a free data retrieval call binding the contract method 0x51117e93.
//
// Solidity: function getZeroBitIfSubmittedCvOnChainBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) GetZeroBitIfSubmittedCvOnChainBitmap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "getZeroBitIfSubmittedCvOnChainBitmap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetZeroBitIfSubmittedCvOnChainBitmap is a free data retrieval call binding the contract method 0x51117e93.
//
// Solidity: function getZeroBitIfSubmittedCvOnChainBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) GetZeroBitIfSubmittedCvOnChainBitmap() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetZeroBitIfSubmittedCvOnChainBitmap(&_CommitReveal2L2.CallOpts)
}

// GetZeroBitIfSubmittedCvOnChainBitmap is a free data retrieval call binding the contract method 0x51117e93.
//
// Solidity: function getZeroBitIfSubmittedCvOnChainBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) GetZeroBitIfSubmittedCvOnChainBitmap() (*big.Int, error) {
	return _CommitReveal2L2.Contract.GetZeroBitIfSubmittedCvOnChainBitmap(&_CommitReveal2L2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CommitReveal2L2 *CommitReveal2L2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CommitReveal2L2 *CommitReveal2L2Session) Owner() (common.Address, error) {
	return _CommitReveal2L2.Contract.Owner(&_CommitReveal2L2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) Owner() (common.Address, error) {
	return _CommitReveal2L2.Contract.Owner(&_CommitReveal2L2.CallOpts)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CommitReveal2L2 *CommitReveal2L2Caller) OwnershipHandoverExpiresAt(opts *bind.CallOpts, pendingOwner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "ownershipHandoverExpiresAt", pendingOwner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CommitReveal2L2 *CommitReveal2L2Session) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.OwnershipHandoverExpiresAt(&_CommitReveal2L2.CallOpts, pendingOwner)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.OwnershipHandoverExpiresAt(&_CommitReveal2L2.CallOpts, pendingOwner)
}

// SActivatedOperatorIndex1Based is a free data retrieval call binding the contract method 0xc71854e3.
//
// Solidity: function s_activatedOperatorIndex1Based(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SActivatedOperatorIndex1Based(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_activatedOperatorIndex1Based", operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SActivatedOperatorIndex1Based is a free data retrieval call binding the contract method 0xc71854e3.
//
// Solidity: function s_activatedOperatorIndex1Based(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SActivatedOperatorIndex1Based(operator common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SActivatedOperatorIndex1Based(&_CommitReveal2L2.CallOpts, operator)
}

// SActivatedOperatorIndex1Based is a free data retrieval call binding the contract method 0xc71854e3.
//
// Solidity: function s_activatedOperatorIndex1Based(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SActivatedOperatorIndex1Based(operator common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SActivatedOperatorIndex1Based(&_CommitReveal2L2.CallOpts, operator)
}

// SActivationThreshold is a free data retrieval call binding the contract method 0xa5ab3014.
//
// Solidity: function s_activationThreshold() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SActivationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_activationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SActivationThreshold is a free data retrieval call binding the contract method 0xa5ab3014.
//
// Solidity: function s_activationThreshold() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SActivationThreshold() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SActivationThreshold(&_CommitReveal2L2.CallOpts)
}

// SActivationThreshold is a free data retrieval call binding the contract method 0xa5ab3014.
//
// Solidity: function s_activationThreshold() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SActivationThreshold() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SActivationThreshold(&_CommitReveal2L2.CallOpts)
}

// SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2 is a free data retrieval call binding the contract method 0xe2b4f51e.
//
// Solidity: function s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2 is a free data retrieval call binding the contract method 0xe2b4f51e.
//
// Solidity: function s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2(&_CommitReveal2L2.CallOpts)
}

// SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2 is a free data retrieval call binding the contract method 0xe2b4f51e.
//
// Solidity: function s_bitSetIfRequestedToSubmitCv_zeroBitIfSubmittedCv_bitmap128x2() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SBitSetIfRequestedToSubmitCvZeroBitIfSubmittedCvBitmap128x2(&_CommitReveal2L2.CallOpts)
}

// SCos is a free data retrieval call binding the contract method 0x85b098b9.
//
// Solidity: function s_cos(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SCos(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_cos", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SCos is a free data retrieval call binding the contract method 0x85b098b9.
//
// Solidity: function s_cos(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Session) SCos(arg0 *big.Int) ([32]byte, error) {
	return _CommitReveal2L2.Contract.SCos(&_CommitReveal2L2.CallOpts, arg0)
}

// SCos is a free data retrieval call binding the contract method 0x85b098b9.
//
// Solidity: function s_cos(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SCos(arg0 *big.Int) ([32]byte, error) {
	return _CommitReveal2L2.Contract.SCos(&_CommitReveal2L2.CallOpts, arg0)
}

// SCurrentRound is a free data retrieval call binding the contract method 0xc5c1676d.
//
// Solidity: function s_currentRound() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SCurrentRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_currentRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SCurrentRound is a free data retrieval call binding the contract method 0xc5c1676d.
//
// Solidity: function s_currentRound() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SCurrentRound() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SCurrentRound(&_CommitReveal2L2.CallOpts)
}

// SCurrentRound is a free data retrieval call binding the contract method 0xc5c1676d.
//
// Solidity: function s_currentRound() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SCurrentRound() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SCurrentRound(&_CommitReveal2L2.CallOpts)
}

// SCvs is a free data retrieval call binding the contract method 0x3a0abc9c.
//
// Solidity: function s_cvs(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SCvs(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_cvs", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SCvs is a free data retrieval call binding the contract method 0x3a0abc9c.
//
// Solidity: function s_cvs(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Session) SCvs(arg0 *big.Int) ([32]byte, error) {
	return _CommitReveal2L2.Contract.SCvs(&_CommitReveal2L2.CallOpts, arg0)
}

// SCvs is a free data retrieval call binding the contract method 0x3a0abc9c.
//
// Solidity: function s_cvs(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SCvs(arg0 *big.Int) ([32]byte, error) {
	return _CommitReveal2L2.Contract.SCvs(&_CommitReveal2L2.CallOpts, arg0)
}

// SDepositAmount is a free data retrieval call binding the contract method 0x42875fad.
//
// Solidity: function s_depositAmount(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SDepositAmount(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_depositAmount", operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SDepositAmount is a free data retrieval call binding the contract method 0x42875fad.
//
// Solidity: function s_depositAmount(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SDepositAmount(operator common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SDepositAmount(&_CommitReveal2L2.CallOpts, operator)
}

// SDepositAmount is a free data retrieval call binding the contract method 0x42875fad.
//
// Solidity: function s_depositAmount(address operator) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SDepositAmount(operator common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SDepositAmount(&_CommitReveal2L2.CallOpts, operator)
}

// SEconomicParamsEffectiveTimestamp is a free data retrieval call binding the contract method 0xde52e819.
//
// Solidity: function s_economicParamsEffectiveTimestamp() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SEconomicParamsEffectiveTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_economicParamsEffectiveTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SEconomicParamsEffectiveTimestamp is a free data retrieval call binding the contract method 0xde52e819.
//
// Solidity: function s_economicParamsEffectiveTimestamp() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SEconomicParamsEffectiveTimestamp() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SEconomicParamsEffectiveTimestamp(&_CommitReveal2L2.CallOpts)
}

// SEconomicParamsEffectiveTimestamp is a free data retrieval call binding the contract method 0xde52e819.
//
// Solidity: function s_economicParamsEffectiveTimestamp() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SEconomicParamsEffectiveTimestamp() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SEconomicParamsEffectiveTimestamp(&_CommitReveal2L2.CallOpts)
}

// SFlatFee is a free data retrieval call binding the contract method 0x0c048a81.
//
// Solidity: function s_flatFee() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SFlatFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_flatFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SFlatFee is a free data retrieval call binding the contract method 0x0c048a81.
//
// Solidity: function s_flatFee() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SFlatFee() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SFlatFee(&_CommitReveal2L2.CallOpts)
}

// SFlatFee is a free data retrieval call binding the contract method 0x0c048a81.
//
// Solidity: function s_flatFee() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SFlatFee() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SFlatFee(&_CommitReveal2L2.CallOpts)
}

// SGasParamsEffectiveTimestamp is a free data retrieval call binding the contract method 0x67582523.
//
// Solidity: function s_gasParamsEffectiveTimestamp() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SGasParamsEffectiveTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_gasParamsEffectiveTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SGasParamsEffectiveTimestamp is a free data retrieval call binding the contract method 0x67582523.
//
// Solidity: function s_gasParamsEffectiveTimestamp() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SGasParamsEffectiveTimestamp() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SGasParamsEffectiveTimestamp(&_CommitReveal2L2.CallOpts)
}

// SGasParamsEffectiveTimestamp is a free data retrieval call binding the contract method 0x67582523.
//
// Solidity: function s_gasParamsEffectiveTimestamp() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SGasParamsEffectiveTimestamp() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SGasParamsEffectiveTimestamp(&_CommitReveal2L2.CallOpts)
}

// SIsInProcess is a free data retrieval call binding the contract method 0x7e6f2b50.
//
// Solidity: function s_isInProcess() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SIsInProcess(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_isInProcess")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SIsInProcess is a free data retrieval call binding the contract method 0x7e6f2b50.
//
// Solidity: function s_isInProcess() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SIsInProcess() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SIsInProcess(&_CommitReveal2L2.CallOpts)
}

// SIsInProcess is a free data retrieval call binding the contract method 0x7e6f2b50.
//
// Solidity: function s_isInProcess() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SIsInProcess() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SIsInProcess(&_CommitReveal2L2.CallOpts)
}

// SL1FeeCoefficient is a free data retrieval call binding the contract method 0x90bd5c74.
//
// Solidity: function s_l1FeeCoefficient() view returns(uint8)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SL1FeeCoefficient(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_l1FeeCoefficient")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SL1FeeCoefficient is a free data retrieval call binding the contract method 0x90bd5c74.
//
// Solidity: function s_l1FeeCoefficient() view returns(uint8)
func (_CommitReveal2L2 *CommitReveal2L2Session) SL1FeeCoefficient() (uint8, error) {
	return _CommitReveal2L2.Contract.SL1FeeCoefficient(&_CommitReveal2L2.CallOpts)
}

// SL1FeeCoefficient is a free data retrieval call binding the contract method 0x90bd5c74.
//
// Solidity: function s_l1FeeCoefficient() view returns(uint8)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SL1FeeCoefficient() (uint8, error) {
	return _CommitReveal2L2.Contract.SL1FeeCoefficient(&_CommitReveal2L2.CallOpts)
}

// SMerkleRoot is a free data retrieval call binding the contract method 0xae82de5f.
//
// Solidity: function s_merkleRoot() view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SMerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SMerkleRoot is a free data retrieval call binding the contract method 0xae82de5f.
//
// Solidity: function s_merkleRoot() view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Session) SMerkleRoot() ([32]byte, error) {
	return _CommitReveal2L2.Contract.SMerkleRoot(&_CommitReveal2L2.CallOpts)
}

// SMerkleRoot is a free data retrieval call binding the contract method 0xae82de5f.
//
// Solidity: function s_merkleRoot() view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SMerkleRoot() ([32]byte, error) {
	return _CommitReveal2L2.Contract.SMerkleRoot(&_CommitReveal2L2.CallOpts)
}

// SMerkleRootSubmittedTimestamp is a free data retrieval call binding the contract method 0x7a526bcb.
//
// Solidity: function s_merkleRootSubmittedTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SMerkleRootSubmittedTimestamp(opts *bind.CallOpts, round *big.Int, trialNum *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_merkleRootSubmittedTimestamp", round, trialNum)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SMerkleRootSubmittedTimestamp is a free data retrieval call binding the contract method 0x7a526bcb.
//
// Solidity: function s_merkleRootSubmittedTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SMerkleRootSubmittedTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SMerkleRootSubmittedTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SMerkleRootSubmittedTimestamp is a free data retrieval call binding the contract method 0x7a526bcb.
//
// Solidity: function s_merkleRootSubmittedTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SMerkleRootSubmittedTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SMerkleRootSubmittedTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SPackedRevealOrders is a free data retrieval call binding the contract method 0x78e8aa34.
//
// Solidity: function s_packedRevealOrders() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SPackedRevealOrders(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_packedRevealOrders")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SPackedRevealOrders is a free data retrieval call binding the contract method 0x78e8aa34.
//
// Solidity: function s_packedRevealOrders() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SPackedRevealOrders() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPackedRevealOrders(&_CommitReveal2L2.CallOpts)
}

// SPackedRevealOrders is a free data retrieval call binding the contract method 0x78e8aa34.
//
// Solidity: function s_packedRevealOrders() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SPackedRevealOrders() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPackedRevealOrders(&_CommitReveal2L2.CallOpts)
}

// SPendingActivationThreshold is a free data retrieval call binding the contract method 0x56aa3663.
//
// Solidity: function s_pendingActivationThreshold() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SPendingActivationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_pendingActivationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SPendingActivationThreshold is a free data retrieval call binding the contract method 0x56aa3663.
//
// Solidity: function s_pendingActivationThreshold() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SPendingActivationThreshold() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPendingActivationThreshold(&_CommitReveal2L2.CallOpts)
}

// SPendingActivationThreshold is a free data retrieval call binding the contract method 0x56aa3663.
//
// Solidity: function s_pendingActivationThreshold() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SPendingActivationThreshold() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPendingActivationThreshold(&_CommitReveal2L2.CallOpts)
}

// SPendingFlatFee is a free data retrieval call binding the contract method 0x6023619a.
//
// Solidity: function s_pendingFlatFee() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SPendingFlatFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_pendingFlatFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SPendingFlatFee is a free data retrieval call binding the contract method 0x6023619a.
//
// Solidity: function s_pendingFlatFee() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SPendingFlatFee() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPendingFlatFee(&_CommitReveal2L2.CallOpts)
}

// SPendingFlatFee is a free data retrieval call binding the contract method 0x6023619a.
//
// Solidity: function s_pendingFlatFee() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SPendingFlatFee() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPendingFlatFee(&_CommitReveal2L2.CallOpts)
}

// SPreviousSSubmitTimestamp is a free data retrieval call binding the contract method 0xd86abadc.
//
// Solidity: function s_previousSSubmitTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SPreviousSSubmitTimestamp(opts *bind.CallOpts, round *big.Int, trialNum *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_previousSSubmitTimestamp", round, trialNum)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SPreviousSSubmitTimestamp is a free data retrieval call binding the contract method 0xd86abadc.
//
// Solidity: function s_previousSSubmitTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SPreviousSSubmitTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPreviousSSubmitTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SPreviousSSubmitTimestamp is a free data retrieval call binding the contract method 0xd86abadc.
//
// Solidity: function s_previousSSubmitTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SPreviousSSubmitTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SPreviousSSubmitTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SRequestCount is a free data retrieval call binding the contract method 0x557d2e92.
//
// Solidity: function s_requestCount() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestCount is a free data retrieval call binding the contract method 0x557d2e92.
//
// Solidity: function s_requestCount() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestCount() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestCount(&_CommitReveal2L2.CallOpts)
}

// SRequestCount is a free data retrieval call binding the contract method 0x557d2e92.
//
// Solidity: function s_requestCount() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestCount() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestCount(&_CommitReveal2L2.CallOpts)
}

// SRequestInfo is a free data retrieval call binding the contract method 0x39d0c151.
//
// Solidity: function s_requestInfo(uint256 round) view returns(address consumer, uint32 callbackGasLimit, uint256 startTime, uint256 cost)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestInfo(opts *bind.CallOpts, round *big.Int) (struct {
	Consumer         common.Address
	CallbackGasLimit uint32
	StartTime        *big.Int
	Cost             *big.Int
}, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestInfo", round)

	outstruct := new(struct {
		Consumer         common.Address
		CallbackGasLimit uint32
		StartTime        *big.Int
		Cost             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Consumer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.StartTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Cost = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SRequestInfo is a free data retrieval call binding the contract method 0x39d0c151.
//
// Solidity: function s_requestInfo(uint256 round) view returns(address consumer, uint32 callbackGasLimit, uint256 startTime, uint256 cost)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestInfo(round *big.Int) (struct {
	Consumer         common.Address
	CallbackGasLimit uint32
	StartTime        *big.Int
	Cost             *big.Int
}, error) {
	return _CommitReveal2L2.Contract.SRequestInfo(&_CommitReveal2L2.CallOpts, round)
}

// SRequestInfo is a free data retrieval call binding the contract method 0x39d0c151.
//
// Solidity: function s_requestInfo(uint256 round) view returns(address consumer, uint32 callbackGasLimit, uint256 startTime, uint256 cost)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestInfo(round *big.Int) (struct {
	Consumer         common.Address
	CallbackGasLimit uint32
	StartTime        *big.Int
	Cost             *big.Int
}, error) {
	return _CommitReveal2L2.Contract.SRequestInfo(&_CommitReveal2L2.CallOpts, round)
}

// SRequestedToSubmitCoLength is a free data retrieval call binding the contract method 0xab941e9d.
//
// Solidity: function s_requestedToSubmitCoLength() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestedToSubmitCoLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestedToSubmitCoLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestedToSubmitCoLength is a free data retrieval call binding the contract method 0xab941e9d.
//
// Solidity: function s_requestedToSubmitCoLength() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestedToSubmitCoLength() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCoLength(&_CommitReveal2L2.CallOpts)
}

// SRequestedToSubmitCoLength is a free data retrieval call binding the contract method 0xab941e9d.
//
// Solidity: function s_requestedToSubmitCoLength() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestedToSubmitCoLength() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCoLength(&_CommitReveal2L2.CallOpts)
}

// SRequestedToSubmitCoPackedIndices is a free data retrieval call binding the contract method 0xe110bfdb.
//
// Solidity: function s_requestedToSubmitCoPackedIndices() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestedToSubmitCoPackedIndices(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestedToSubmitCoPackedIndices")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestedToSubmitCoPackedIndices is a free data retrieval call binding the contract method 0xe110bfdb.
//
// Solidity: function s_requestedToSubmitCoPackedIndices() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestedToSubmitCoPackedIndices() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCoPackedIndices(&_CommitReveal2L2.CallOpts)
}

// SRequestedToSubmitCoPackedIndices is a free data retrieval call binding the contract method 0xe110bfdb.
//
// Solidity: function s_requestedToSubmitCoPackedIndices() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestedToSubmitCoPackedIndices() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCoPackedIndices(&_CommitReveal2L2.CallOpts)
}

// SRequestedToSubmitCoTimestamp is a free data retrieval call binding the contract method 0xbf81cdea.
//
// Solidity: function s_requestedToSubmitCoTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestedToSubmitCoTimestamp(opts *bind.CallOpts, round *big.Int, trialNum *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestedToSubmitCoTimestamp", round, trialNum)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestedToSubmitCoTimestamp is a free data retrieval call binding the contract method 0xbf81cdea.
//
// Solidity: function s_requestedToSubmitCoTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestedToSubmitCoTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCoTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SRequestedToSubmitCoTimestamp is a free data retrieval call binding the contract method 0xbf81cdea.
//
// Solidity: function s_requestedToSubmitCoTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestedToSubmitCoTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCoTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SRequestedToSubmitCvPackedIndicesAscFromLSB is a free data retrieval call binding the contract method 0x6a2f054b.
//
// Solidity: function s_requestedToSubmitCvPackedIndicesAscFromLSB() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestedToSubmitCvPackedIndicesAscFromLSB(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestedToSubmitCvPackedIndicesAscFromLSB")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestedToSubmitCvPackedIndicesAscFromLSB is a free data retrieval call binding the contract method 0x6a2f054b.
//
// Solidity: function s_requestedToSubmitCvPackedIndicesAscFromLSB() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestedToSubmitCvPackedIndicesAscFromLSB() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCvPackedIndicesAscFromLSB(&_CommitReveal2L2.CallOpts)
}

// SRequestedToSubmitCvPackedIndicesAscFromLSB is a free data retrieval call binding the contract method 0x6a2f054b.
//
// Solidity: function s_requestedToSubmitCvPackedIndicesAscFromLSB() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestedToSubmitCvPackedIndicesAscFromLSB() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCvPackedIndicesAscFromLSB(&_CommitReveal2L2.CallOpts)
}

// SRequestedToSubmitCvTimestamp is a free data retrieval call binding the contract method 0xd033cebb.
//
// Solidity: function s_requestedToSubmitCvTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestedToSubmitCvTimestamp(opts *bind.CallOpts, round *big.Int, trialNum *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestedToSubmitCvTimestamp", round, trialNum)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestedToSubmitCvTimestamp is a free data retrieval call binding the contract method 0xd033cebb.
//
// Solidity: function s_requestedToSubmitCvTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestedToSubmitCvTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCvTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SRequestedToSubmitCvTimestamp is a free data retrieval call binding the contract method 0xd033cebb.
//
// Solidity: function s_requestedToSubmitCvTimestamp(uint256 round, uint256 trialNum) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestedToSubmitCvTimestamp(round *big.Int, trialNum *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitCvTimestamp(&_CommitReveal2L2.CallOpts, round, trialNum)
}

// SRequestedToSubmitSFromIndexK is a free data retrieval call binding the contract method 0xe2bf5dac.
//
// Solidity: function s_requestedToSubmitSFromIndexK() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRequestedToSubmitSFromIndexK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_requestedToSubmitSFromIndexK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestedToSubmitSFromIndexK is a free data retrieval call binding the contract method 0xe2bf5dac.
//
// Solidity: function s_requestedToSubmitSFromIndexK() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRequestedToSubmitSFromIndexK() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitSFromIndexK(&_CommitReveal2L2.CallOpts)
}

// SRequestedToSubmitSFromIndexK is a free data retrieval call binding the contract method 0xe2bf5dac.
//
// Solidity: function s_requestedToSubmitSFromIndexK() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRequestedToSubmitSFromIndexK() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRequestedToSubmitSFromIndexK(&_CommitReveal2L2.CallOpts)
}

// SRoundBitmap is a free data retrieval call binding the contract method 0xc78c1776.
//
// Solidity: function s_roundBitmap(uint248 wordPos) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SRoundBitmap(opts *bind.CallOpts, wordPos *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_roundBitmap", wordPos)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRoundBitmap is a free data retrieval call binding the contract method 0xc78c1776.
//
// Solidity: function s_roundBitmap(uint248 wordPos) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SRoundBitmap(wordPos *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRoundBitmap(&_CommitReveal2L2.CallOpts, wordPos)
}

// SRoundBitmap is a free data retrieval call binding the contract method 0xc78c1776.
//
// Solidity: function s_roundBitmap(uint248 wordPos) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SRoundBitmap(wordPos *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SRoundBitmap(&_CommitReveal2L2.CallOpts, wordPos)
}

// SSecrets is a free data retrieval call binding the contract method 0xffcb420f.
//
// Solidity: function s_secrets(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SSecrets(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_secrets", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SSecrets is a free data retrieval call binding the contract method 0xffcb420f.
//
// Solidity: function s_secrets(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2Session) SSecrets(arg0 *big.Int) ([32]byte, error) {
	return _CommitReveal2L2.Contract.SSecrets(&_CommitReveal2L2.CallOpts, arg0)
}

// SSecrets is a free data retrieval call binding the contract method 0xffcb420f.
//
// Solidity: function s_secrets(uint256 ) view returns(bytes32)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SSecrets(arg0 *big.Int) ([32]byte, error) {
	return _CommitReveal2L2.Contract.SSecrets(&_CommitReveal2L2.CallOpts, arg0)
}

// SSlashRewardPerOperatorPaidX8 is a free data retrieval call binding the contract method 0x36748d6e.
//
// Solidity: function s_slashRewardPerOperatorPaidX8(address ) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SSlashRewardPerOperatorPaidX8(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_slashRewardPerOperatorPaidX8", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SSlashRewardPerOperatorPaidX8 is a free data retrieval call binding the contract method 0x36748d6e.
//
// Solidity: function s_slashRewardPerOperatorPaidX8(address ) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SSlashRewardPerOperatorPaidX8(arg0 common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SSlashRewardPerOperatorPaidX8(&_CommitReveal2L2.CallOpts, arg0)
}

// SSlashRewardPerOperatorPaidX8 is a free data retrieval call binding the contract method 0x36748d6e.
//
// Solidity: function s_slashRewardPerOperatorPaidX8(address ) view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SSlashRewardPerOperatorPaidX8(arg0 common.Address) (*big.Int, error) {
	return _CommitReveal2L2.Contract.SSlashRewardPerOperatorPaidX8(&_CommitReveal2L2.CallOpts, arg0)
}

// SSlashRewardPerOperatorX8 is a free data retrieval call binding the contract method 0x1722fda3.
//
// Solidity: function s_slashRewardPerOperatorX8() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SSlashRewardPerOperatorX8(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_slashRewardPerOperatorX8")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SSlashRewardPerOperatorX8 is a free data retrieval call binding the contract method 0x1722fda3.
//
// Solidity: function s_slashRewardPerOperatorX8() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SSlashRewardPerOperatorX8() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SSlashRewardPerOperatorX8(&_CommitReveal2L2.CallOpts)
}

// SSlashRewardPerOperatorX8 is a free data retrieval call binding the contract method 0x1722fda3.
//
// Solidity: function s_slashRewardPerOperatorX8() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SSlashRewardPerOperatorX8() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SSlashRewardPerOperatorX8(&_CommitReveal2L2.CallOpts)
}

// STrialNum is a free data retrieval call binding the contract method 0xc408b393.
//
// Solidity: function s_trialNum(uint256 round) view returns(uint256 trialNum)
func (_CommitReveal2L2 *CommitReveal2L2Caller) STrialNum(opts *bind.CallOpts, round *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_trialNum", round)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STrialNum is a free data retrieval call binding the contract method 0xc408b393.
//
// Solidity: function s_trialNum(uint256 round) view returns(uint256 trialNum)
func (_CommitReveal2L2 *CommitReveal2L2Session) STrialNum(round *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.STrialNum(&_CommitReveal2L2.CallOpts, round)
}

// STrialNum is a free data retrieval call binding the contract method 0xc408b393.
//
// Solidity: function s_trialNum(uint256 round) view returns(uint256 trialNum)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) STrialNum(round *big.Int) (*big.Int, error) {
	return _CommitReveal2L2.Contract.STrialNum(&_CommitReveal2L2.CallOpts, round)
}

// SZeroBitIfSubmittedCoBitmap is a free data retrieval call binding the contract method 0xef53459d.
//
// Solidity: function s_zeroBitIfSubmittedCoBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Caller) SZeroBitIfSubmittedCoBitmap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CommitReveal2L2.contract.Call(opts, &out, "s_zeroBitIfSubmittedCoBitmap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SZeroBitIfSubmittedCoBitmap is a free data retrieval call binding the contract method 0xef53459d.
//
// Solidity: function s_zeroBitIfSubmittedCoBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) SZeroBitIfSubmittedCoBitmap() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SZeroBitIfSubmittedCoBitmap(&_CommitReveal2L2.CallOpts)
}

// SZeroBitIfSubmittedCoBitmap is a free data retrieval call binding the contract method 0xef53459d.
//
// Solidity: function s_zeroBitIfSubmittedCoBitmap() view returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2CallerSession) SZeroBitIfSubmittedCoBitmap() (*big.Int, error) {
	return _CommitReveal2L2.Contract.SZeroBitIfSubmittedCoBitmap(&_CommitReveal2L2.CallOpts)
}

// Activate is a paid mutator transaction binding the contract method 0x0f15f4c0.
//
// Solidity: function activate() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) Activate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "activate")
}

// Activate is a paid mutator transaction binding the contract method 0x0f15f4c0.
//
// Solidity: function activate() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) Activate() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Activate(&_CommitReveal2L2.TransactOpts)
}

// Activate is a paid mutator transaction binding the contract method 0x0f15f4c0.
//
// Solidity: function activate() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) Activate() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Activate(&_CommitReveal2L2.TransactOpts)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) CancelOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "cancelOwnershipHandover")
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) CancelOwnershipHandover() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.CancelOwnershipHandover(&_CommitReveal2L2.TransactOpts)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.CancelOwnershipHandover(&_CommitReveal2L2.TransactOpts)
}

// ClaimSlashReward is a paid mutator transaction binding the contract method 0xe20d5138.
//
// Solidity: function claimSlashReward() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) ClaimSlashReward(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "claimSlashReward")
}

// ClaimSlashReward is a paid mutator transaction binding the contract method 0xe20d5138.
//
// Solidity: function claimSlashReward() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) ClaimSlashReward() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ClaimSlashReward(&_CommitReveal2L2.TransactOpts)
}

// ClaimSlashReward is a paid mutator transaction binding the contract method 0xe20d5138.
//
// Solidity: function claimSlashReward() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) ClaimSlashReward() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ClaimSlashReward(&_CommitReveal2L2.TransactOpts)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) CompleteOwnershipHandover(opts *bind.TransactOpts, pendingOwner common.Address) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "completeOwnershipHandover", pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.CompleteOwnershipHandover(&_CommitReveal2L2.TransactOpts, pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.CompleteOwnershipHandover(&_CommitReveal2L2.TransactOpts, pendingOwner)
}

// Deactivate is a paid mutator transaction binding the contract method 0x51b42b00.
//
// Solidity: function deactivate() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) Deactivate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "deactivate")
}

// Deactivate is a paid mutator transaction binding the contract method 0x51b42b00.
//
// Solidity: function deactivate() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) Deactivate() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Deactivate(&_CommitReveal2L2.TransactOpts)
}

// Deactivate is a paid mutator transaction binding the contract method 0x51b42b00.
//
// Solidity: function deactivate() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) Deactivate() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Deactivate(&_CommitReveal2L2.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) Deposit() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Deposit(&_CommitReveal2L2.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) Deposit() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Deposit(&_CommitReveal2L2.TransactOpts)
}

// DepositAndActivate is a paid mutator transaction binding the contract method 0x77343032.
//
// Solidity: function depositAndActivate() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) DepositAndActivate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "depositAndActivate")
}

// DepositAndActivate is a paid mutator transaction binding the contract method 0x77343032.
//
// Solidity: function depositAndActivate() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) DepositAndActivate() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.DepositAndActivate(&_CommitReveal2L2.TransactOpts)
}

// DepositAndActivate is a paid mutator transaction binding the contract method 0x77343032.
//
// Solidity: function depositAndActivate() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) DepositAndActivate() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.DepositAndActivate(&_CommitReveal2L2.TransactOpts)
}

// ExecuteSetEconomicParameters is a paid mutator transaction binding the contract method 0xf3d1d7b1.
//
// Solidity: function executeSetEconomicParameters() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) ExecuteSetEconomicParameters(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "executeSetEconomicParameters")
}

// ExecuteSetEconomicParameters is a paid mutator transaction binding the contract method 0xf3d1d7b1.
//
// Solidity: function executeSetEconomicParameters() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) ExecuteSetEconomicParameters() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ExecuteSetEconomicParameters(&_CommitReveal2L2.TransactOpts)
}

// ExecuteSetEconomicParameters is a paid mutator transaction binding the contract method 0xf3d1d7b1.
//
// Solidity: function executeSetEconomicParameters() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) ExecuteSetEconomicParameters() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ExecuteSetEconomicParameters(&_CommitReveal2L2.TransactOpts)
}

// ExecuteSetGasParameters is a paid mutator transaction binding the contract method 0x9ad6661f.
//
// Solidity: function executeSetGasParameters() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) ExecuteSetGasParameters(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "executeSetGasParameters")
}

// ExecuteSetGasParameters is a paid mutator transaction binding the contract method 0x9ad6661f.
//
// Solidity: function executeSetGasParameters() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) ExecuteSetGasParameters() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ExecuteSetGasParameters(&_CommitReveal2L2.TransactOpts)
}

// ExecuteSetGasParameters is a paid mutator transaction binding the contract method 0x9ad6661f.
//
// Solidity: function executeSetGasParameters() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) ExecuteSetGasParameters() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ExecuteSetGasParameters(&_CommitReveal2L2.TransactOpts)
}

// FailToRequestSorGenerateRandomNumber is a paid mutator transaction binding the contract method 0x21682bbf.
//
// Solidity: function failToRequestSorGenerateRandomNumber() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) FailToRequestSorGenerateRandomNumber(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "failToRequestSorGenerateRandomNumber")
}

// FailToRequestSorGenerateRandomNumber is a paid mutator transaction binding the contract method 0x21682bbf.
//
// Solidity: function failToRequestSorGenerateRandomNumber() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) FailToRequestSorGenerateRandomNumber() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToRequestSorGenerateRandomNumber(&_CommitReveal2L2.TransactOpts)
}

// FailToRequestSorGenerateRandomNumber is a paid mutator transaction binding the contract method 0x21682bbf.
//
// Solidity: function failToRequestSorGenerateRandomNumber() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) FailToRequestSorGenerateRandomNumber() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToRequestSorGenerateRandomNumber(&_CommitReveal2L2.TransactOpts)
}

// FailToRequestSubmitCvOrSubmitMerkleRoot is a paid mutator transaction binding the contract method 0xf224284d.
//
// Solidity: function failToRequestSubmitCvOrSubmitMerkleRoot() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) FailToRequestSubmitCvOrSubmitMerkleRoot(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "failToRequestSubmitCvOrSubmitMerkleRoot")
}

// FailToRequestSubmitCvOrSubmitMerkleRoot is a paid mutator transaction binding the contract method 0xf224284d.
//
// Solidity: function failToRequestSubmitCvOrSubmitMerkleRoot() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) FailToRequestSubmitCvOrSubmitMerkleRoot() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToRequestSubmitCvOrSubmitMerkleRoot(&_CommitReveal2L2.TransactOpts)
}

// FailToRequestSubmitCvOrSubmitMerkleRoot is a paid mutator transaction binding the contract method 0xf224284d.
//
// Solidity: function failToRequestSubmitCvOrSubmitMerkleRoot() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) FailToRequestSubmitCvOrSubmitMerkleRoot() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToRequestSubmitCvOrSubmitMerkleRoot(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitCo is a paid mutator transaction binding the contract method 0xb2a3cbff.
//
// Solidity: function failToSubmitCo() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) FailToSubmitCo(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "failToSubmitCo")
}

// FailToSubmitCo is a paid mutator transaction binding the contract method 0xb2a3cbff.
//
// Solidity: function failToSubmitCo() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) FailToSubmitCo() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitCo(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitCo is a paid mutator transaction binding the contract method 0xb2a3cbff.
//
// Solidity: function failToSubmitCo() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) FailToSubmitCo() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitCo(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitCv is a paid mutator transaction binding the contract method 0x5ec16a60.
//
// Solidity: function failToSubmitCv() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) FailToSubmitCv(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "failToSubmitCv")
}

// FailToSubmitCv is a paid mutator transaction binding the contract method 0x5ec16a60.
//
// Solidity: function failToSubmitCv() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) FailToSubmitCv() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitCv(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitCv is a paid mutator transaction binding the contract method 0x5ec16a60.
//
// Solidity: function failToSubmitCv() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) FailToSubmitCv() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitCv(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitMerkleRootAfterDispute is a paid mutator transaction binding the contract method 0xc4c0299d.
//
// Solidity: function failToSubmitMerkleRootAfterDispute() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) FailToSubmitMerkleRootAfterDispute(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "failToSubmitMerkleRootAfterDispute")
}

// FailToSubmitMerkleRootAfterDispute is a paid mutator transaction binding the contract method 0xc4c0299d.
//
// Solidity: function failToSubmitMerkleRootAfterDispute() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) FailToSubmitMerkleRootAfterDispute() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitMerkleRootAfterDispute(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitMerkleRootAfterDispute is a paid mutator transaction binding the contract method 0xc4c0299d.
//
// Solidity: function failToSubmitMerkleRootAfterDispute() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) FailToSubmitMerkleRootAfterDispute() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitMerkleRootAfterDispute(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitS is a paid mutator transaction binding the contract method 0x3d8620c3.
//
// Solidity: function failToSubmitS() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) FailToSubmitS(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "failToSubmitS")
}

// FailToSubmitS is a paid mutator transaction binding the contract method 0x3d8620c3.
//
// Solidity: function failToSubmitS() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) FailToSubmitS() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitS(&_CommitReveal2L2.TransactOpts)
}

// FailToSubmitS is a paid mutator transaction binding the contract method 0x3d8620c3.
//
// Solidity: function failToSubmitS() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) FailToSubmitS() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.FailToSubmitS(&_CommitReveal2L2.TransactOpts)
}

// GenerateRandomNumber is a paid mutator transaction binding the contract method 0xa89b873e.
//
// Solidity: function generateRandomNumber((bytes32,(bytes32,bytes32))[] secretSigRSs, uint256 , uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) GenerateRandomNumber(opts *bind.TransactOpts, secretSigRSs []CommitReveal2StorageSecretAndSigRS, arg1 *big.Int, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "generateRandomNumber", secretSigRSs, arg1, packedRevealOrders)
}

// GenerateRandomNumber is a paid mutator transaction binding the contract method 0xa89b873e.
//
// Solidity: function generateRandomNumber((bytes32,(bytes32,bytes32))[] secretSigRSs, uint256 , uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) GenerateRandomNumber(secretSigRSs []CommitReveal2StorageSecretAndSigRS, arg1 *big.Int, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.GenerateRandomNumber(&_CommitReveal2L2.TransactOpts, secretSigRSs, arg1, packedRevealOrders)
}

// GenerateRandomNumber is a paid mutator transaction binding the contract method 0xa89b873e.
//
// Solidity: function generateRandomNumber((bytes32,(bytes32,bytes32))[] secretSigRSs, uint256 , uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) GenerateRandomNumber(secretSigRSs []CommitReveal2StorageSecretAndSigRS, arg1 *big.Int, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.GenerateRandomNumber(&_CommitReveal2L2.TransactOpts, secretSigRSs, arg1, packedRevealOrders)
}

// GenerateRandomNumberWhenSomeCvsAreOnChain is a paid mutator transaction binding the contract method 0x0ab7a70f.
//
// Solidity: function generateRandomNumberWhenSomeCvsAreOnChain(bytes32[] allSecrets, (bytes32,bytes32)[] sigRSsForAllCvsNotOnChain, uint256 , uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) GenerateRandomNumberWhenSomeCvsAreOnChain(opts *bind.TransactOpts, allSecrets [][32]byte, sigRSsForAllCvsNotOnChain []CommitReveal2StorageSigRS, arg2 *big.Int, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "generateRandomNumberWhenSomeCvsAreOnChain", allSecrets, sigRSsForAllCvsNotOnChain, arg2, packedRevealOrders)
}

// GenerateRandomNumberWhenSomeCvsAreOnChain is a paid mutator transaction binding the contract method 0x0ab7a70f.
//
// Solidity: function generateRandomNumberWhenSomeCvsAreOnChain(bytes32[] allSecrets, (bytes32,bytes32)[] sigRSsForAllCvsNotOnChain, uint256 , uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) GenerateRandomNumberWhenSomeCvsAreOnChain(allSecrets [][32]byte, sigRSsForAllCvsNotOnChain []CommitReveal2StorageSigRS, arg2 *big.Int, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.GenerateRandomNumberWhenSomeCvsAreOnChain(&_CommitReveal2L2.TransactOpts, allSecrets, sigRSsForAllCvsNotOnChain, arg2, packedRevealOrders)
}

// GenerateRandomNumberWhenSomeCvsAreOnChain is a paid mutator transaction binding the contract method 0x0ab7a70f.
//
// Solidity: function generateRandomNumberWhenSomeCvsAreOnChain(bytes32[] allSecrets, (bytes32,bytes32)[] sigRSsForAllCvsNotOnChain, uint256 , uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) GenerateRandomNumberWhenSomeCvsAreOnChain(allSecrets [][32]byte, sigRSsForAllCvsNotOnChain []CommitReveal2StorageSigRS, arg2 *big.Int, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.GenerateRandomNumberWhenSomeCvsAreOnChain(&_CommitReveal2L2.TransactOpts, allSecrets, sigRSsForAllCvsNotOnChain, arg2, packedRevealOrders)
}

// ProposeEconomicParameters is a paid mutator transaction binding the contract method 0x124acaa2.
//
// Solidity: function proposeEconomicParameters(uint256 activationThreshold, uint256 flatFee) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) ProposeEconomicParameters(opts *bind.TransactOpts, activationThreshold *big.Int, flatFee *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "proposeEconomicParameters", activationThreshold, flatFee)
}

// ProposeEconomicParameters is a paid mutator transaction binding the contract method 0x124acaa2.
//
// Solidity: function proposeEconomicParameters(uint256 activationThreshold, uint256 flatFee) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) ProposeEconomicParameters(activationThreshold *big.Int, flatFee *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ProposeEconomicParameters(&_CommitReveal2L2.TransactOpts, activationThreshold, flatFee)
}

// ProposeEconomicParameters is a paid mutator transaction binding the contract method 0x124acaa2.
//
// Solidity: function proposeEconomicParameters(uint256 activationThreshold, uint256 flatFee) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) ProposeEconomicParameters(activationThreshold *big.Int, flatFee *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ProposeEconomicParameters(&_CommitReveal2L2.TransactOpts, activationThreshold, flatFee)
}

// ProposeGasParameters is a paid mutator transaction binding the contract method 0xe872de60.
//
// Solidity: function proposeGasParameters(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) ProposeGasParameters(opts *bind.TransactOpts, gasUsedMerkleRootSubAndGenRandNumA *big.Int, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead *big.Int, maxCallbackGasLimit *big.Int, getL1UpperBoundGasUsedWhenCalldataSize4 *big.Int, failToRequestCvOrSubmitMerkleRootGasUsed *big.Int, failToSubmitMerkleRootAfterDisputeGasUsed *big.Int, failToRequestSOrGenerateRandomNumberGasUsed *big.Int, failToSubmitSGasUsed *big.Int, failToSubmitCoGasUsedBaseA uint32, failToSubmitCvGasUsedBaseA uint32, failToSubmitGasUsedBaseB uint32, perOperatorIncreaseGasUsedA uint32, perOperatorIncreaseGasUsedB uint32, perAdditionalDidntSubmitGasUsedA uint32, perAdditionalDidntSubmitGasUsedB uint32, perRequestedIncreaseGasUsed uint32) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "proposeGasParameters", gasUsedMerkleRootSubAndGenRandNumA, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead, maxCallbackGasLimit, getL1UpperBoundGasUsedWhenCalldataSize4, failToRequestCvOrSubmitMerkleRootGasUsed, failToSubmitMerkleRootAfterDisputeGasUsed, failToRequestSOrGenerateRandomNumberGasUsed, failToSubmitSGasUsed, failToSubmitCoGasUsedBaseA, failToSubmitCvGasUsedBaseA, failToSubmitGasUsedBaseB, perOperatorIncreaseGasUsedA, perOperatorIncreaseGasUsedB, perAdditionalDidntSubmitGasUsedA, perAdditionalDidntSubmitGasUsedB, perRequestedIncreaseGasUsed)
}

// ProposeGasParameters is a paid mutator transaction binding the contract method 0xe872de60.
//
// Solidity: function proposeGasParameters(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) ProposeGasParameters(gasUsedMerkleRootSubAndGenRandNumA *big.Int, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead *big.Int, maxCallbackGasLimit *big.Int, getL1UpperBoundGasUsedWhenCalldataSize4 *big.Int, failToRequestCvOrSubmitMerkleRootGasUsed *big.Int, failToSubmitMerkleRootAfterDisputeGasUsed *big.Int, failToRequestSOrGenerateRandomNumberGasUsed *big.Int, failToSubmitSGasUsed *big.Int, failToSubmitCoGasUsedBaseA uint32, failToSubmitCvGasUsedBaseA uint32, failToSubmitGasUsedBaseB uint32, perOperatorIncreaseGasUsedA uint32, perOperatorIncreaseGasUsedB uint32, perAdditionalDidntSubmitGasUsedA uint32, perAdditionalDidntSubmitGasUsedB uint32, perRequestedIncreaseGasUsed uint32) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ProposeGasParameters(&_CommitReveal2L2.TransactOpts, gasUsedMerkleRootSubAndGenRandNumA, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead, maxCallbackGasLimit, getL1UpperBoundGasUsedWhenCalldataSize4, failToRequestCvOrSubmitMerkleRootGasUsed, failToSubmitMerkleRootAfterDisputeGasUsed, failToRequestSOrGenerateRandomNumberGasUsed, failToSubmitSGasUsed, failToSubmitCoGasUsedBaseA, failToSubmitCvGasUsedBaseA, failToSubmitGasUsedBaseB, perOperatorIncreaseGasUsedA, perOperatorIncreaseGasUsedB, perAdditionalDidntSubmitGasUsedA, perAdditionalDidntSubmitGasUsedB, perRequestedIncreaseGasUsed)
}

// ProposeGasParameters is a paid mutator transaction binding the contract method 0xe872de60.
//
// Solidity: function proposeGasParameters(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) ProposeGasParameters(gasUsedMerkleRootSubAndGenRandNumA *big.Int, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead *big.Int, maxCallbackGasLimit *big.Int, getL1UpperBoundGasUsedWhenCalldataSize4 *big.Int, failToRequestCvOrSubmitMerkleRootGasUsed *big.Int, failToSubmitMerkleRootAfterDisputeGasUsed *big.Int, failToRequestSOrGenerateRandomNumberGasUsed *big.Int, failToSubmitSGasUsed *big.Int, failToSubmitCoGasUsedBaseA uint32, failToSubmitCvGasUsedBaseA uint32, failToSubmitGasUsedBaseB uint32, perOperatorIncreaseGasUsedA uint32, perOperatorIncreaseGasUsedB uint32, perAdditionalDidntSubmitGasUsedA uint32, perAdditionalDidntSubmitGasUsedB uint32, perRequestedIncreaseGasUsed uint32) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.ProposeGasParameters(&_CommitReveal2L2.TransactOpts, gasUsedMerkleRootSubAndGenRandNumA, gasUsedMerkleRootSubAndGenRandNumBWithLeaderOverhead, maxCallbackGasLimit, getL1UpperBoundGasUsedWhenCalldataSize4, failToRequestCvOrSubmitMerkleRootGasUsed, failToSubmitMerkleRootAfterDisputeGasUsed, failToRequestSOrGenerateRandomNumberGasUsed, failToSubmitSGasUsed, failToSubmitCoGasUsedBaseA, failToSubmitCvGasUsedBaseA, failToSubmitGasUsedBaseB, perOperatorIncreaseGasUsedA, perOperatorIncreaseGasUsedB, perAdditionalDidntSubmitGasUsedA, perAdditionalDidntSubmitGasUsedB, perRequestedIncreaseGasUsed)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 round) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) Refund(opts *bind.TransactOpts, round *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "refund", round)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 round) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) Refund(round *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Refund(&_CommitReveal2L2.TransactOpts, round)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 round) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) Refund(round *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Refund(&_CommitReveal2L2.TransactOpts, round)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) RenounceOwnership() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RenounceOwnership(&_CommitReveal2L2.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RenounceOwnership(&_CommitReveal2L2.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) RequestOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "requestOwnershipHandover")
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) RequestOwnershipHandover() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestOwnershipHandover(&_CommitReveal2L2.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestOwnershipHandover(&_CommitReveal2L2.TransactOpts)
}

// RequestRandomNumber is a paid mutator transaction binding the contract method 0xb5f3abb0.
//
// Solidity: function requestRandomNumber(uint32 callbackGasLimit) payable returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Transactor) RequestRandomNumber(opts *bind.TransactOpts, callbackGasLimit uint32) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "requestRandomNumber", callbackGasLimit)
}

// RequestRandomNumber is a paid mutator transaction binding the contract method 0xb5f3abb0.
//
// Solidity: function requestRandomNumber(uint32 callbackGasLimit) payable returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2Session) RequestRandomNumber(callbackGasLimit uint32) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestRandomNumber(&_CommitReveal2L2.TransactOpts, callbackGasLimit)
}

// RequestRandomNumber is a paid mutator transaction binding the contract method 0xb5f3abb0.
//
// Solidity: function requestRandomNumber(uint32 callbackGasLimit) payable returns(uint256)
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) RequestRandomNumber(callbackGasLimit uint32) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestRandomNumber(&_CommitReveal2L2.TransactOpts, callbackGasLimit)
}

// RequestToSubmitCo is a paid mutator transaction binding the contract method 0x2882eb30.
//
// Solidity: function requestToSubmitCo((bytes32,(bytes32,bytes32))[] cvRSsForCvsNotOnChainAndReqToSubmitCo, uint256 , uint256 indicesLength, uint256 packedIndicesFirstCvNotOnChainRestCvOnChain) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) RequestToSubmitCo(opts *bind.TransactOpts, cvRSsForCvsNotOnChainAndReqToSubmitCo []CommitReveal2StorageCvAndSigRS, arg1 *big.Int, indicesLength *big.Int, packedIndicesFirstCvNotOnChainRestCvOnChain *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "requestToSubmitCo", cvRSsForCvsNotOnChainAndReqToSubmitCo, arg1, indicesLength, packedIndicesFirstCvNotOnChainRestCvOnChain)
}

// RequestToSubmitCo is a paid mutator transaction binding the contract method 0x2882eb30.
//
// Solidity: function requestToSubmitCo((bytes32,(bytes32,bytes32))[] cvRSsForCvsNotOnChainAndReqToSubmitCo, uint256 , uint256 indicesLength, uint256 packedIndicesFirstCvNotOnChainRestCvOnChain) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) RequestToSubmitCo(cvRSsForCvsNotOnChainAndReqToSubmitCo []CommitReveal2StorageCvAndSigRS, arg1 *big.Int, indicesLength *big.Int, packedIndicesFirstCvNotOnChainRestCvOnChain *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestToSubmitCo(&_CommitReveal2L2.TransactOpts, cvRSsForCvsNotOnChainAndReqToSubmitCo, arg1, indicesLength, packedIndicesFirstCvNotOnChainRestCvOnChain)
}

// RequestToSubmitCo is a paid mutator transaction binding the contract method 0x2882eb30.
//
// Solidity: function requestToSubmitCo((bytes32,(bytes32,bytes32))[] cvRSsForCvsNotOnChainAndReqToSubmitCo, uint256 , uint256 indicesLength, uint256 packedIndicesFirstCvNotOnChainRestCvOnChain) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) RequestToSubmitCo(cvRSsForCvsNotOnChainAndReqToSubmitCo []CommitReveal2StorageCvAndSigRS, arg1 *big.Int, indicesLength *big.Int, packedIndicesFirstCvNotOnChainRestCvOnChain *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestToSubmitCo(&_CommitReveal2L2.TransactOpts, cvRSsForCvsNotOnChainAndReqToSubmitCo, arg1, indicesLength, packedIndicesFirstCvNotOnChainRestCvOnChain)
}

// RequestToSubmitCv is a paid mutator transaction binding the contract method 0xe08e8349.
//
// Solidity: function requestToSubmitCv(uint256 packedIndicesAscendingFromLSB) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) RequestToSubmitCv(opts *bind.TransactOpts, packedIndicesAscendingFromLSB *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "requestToSubmitCv", packedIndicesAscendingFromLSB)
}

// RequestToSubmitCv is a paid mutator transaction binding the contract method 0xe08e8349.
//
// Solidity: function requestToSubmitCv(uint256 packedIndicesAscendingFromLSB) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) RequestToSubmitCv(packedIndicesAscendingFromLSB *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestToSubmitCv(&_CommitReveal2L2.TransactOpts, packedIndicesAscendingFromLSB)
}

// RequestToSubmitCv is a paid mutator transaction binding the contract method 0xe08e8349.
//
// Solidity: function requestToSubmitCv(uint256 packedIndicesAscendingFromLSB) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) RequestToSubmitCv(packedIndicesAscendingFromLSB *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestToSubmitCv(&_CommitReveal2L2.TransactOpts, packedIndicesAscendingFromLSB)
}

// RequestToSubmitS is a paid mutator transaction binding the contract method 0xae6b251e.
//
// Solidity: function requestToSubmitS(bytes32[] allCos, bytes32[] secretsReceivedOffchainInRevealOrder, uint256 , (bytes32,bytes32)[] sigRSsForAllCvsNotOnChain, uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) RequestToSubmitS(opts *bind.TransactOpts, allCos [][32]byte, secretsReceivedOffchainInRevealOrder [][32]byte, arg2 *big.Int, sigRSsForAllCvsNotOnChain []CommitReveal2StorageSigRS, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "requestToSubmitS", allCos, secretsReceivedOffchainInRevealOrder, arg2, sigRSsForAllCvsNotOnChain, packedRevealOrders)
}

// RequestToSubmitS is a paid mutator transaction binding the contract method 0xae6b251e.
//
// Solidity: function requestToSubmitS(bytes32[] allCos, bytes32[] secretsReceivedOffchainInRevealOrder, uint256 , (bytes32,bytes32)[] sigRSsForAllCvsNotOnChain, uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) RequestToSubmitS(allCos [][32]byte, secretsReceivedOffchainInRevealOrder [][32]byte, arg2 *big.Int, sigRSsForAllCvsNotOnChain []CommitReveal2StorageSigRS, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestToSubmitS(&_CommitReveal2L2.TransactOpts, allCos, secretsReceivedOffchainInRevealOrder, arg2, sigRSsForAllCvsNotOnChain, packedRevealOrders)
}

// RequestToSubmitS is a paid mutator transaction binding the contract method 0xae6b251e.
//
// Solidity: function requestToSubmitS(bytes32[] allCos, bytes32[] secretsReceivedOffchainInRevealOrder, uint256 , (bytes32,bytes32)[] sigRSsForAllCvsNotOnChain, uint256 packedRevealOrders) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) RequestToSubmitS(allCos [][32]byte, secretsReceivedOffchainInRevealOrder [][32]byte, arg2 *big.Int, sigRSsForAllCvsNotOnChain []CommitReveal2StorageSigRS, packedRevealOrders *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.RequestToSubmitS(&_CommitReveal2L2.TransactOpts, allCos, secretsReceivedOffchainInRevealOrder, arg2, sigRSsForAllCvsNotOnChain, packedRevealOrders)
}

// Resume is a paid mutator transaction binding the contract method 0x046f7da2.
//
// Solidity: function resume() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) Resume(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "resume")
}

// Resume is a paid mutator transaction binding the contract method 0x046f7da2.
//
// Solidity: function resume() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) Resume() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Resume(&_CommitReveal2L2.TransactOpts)
}

// Resume is a paid mutator transaction binding the contract method 0x046f7da2.
//
// Solidity: function resume() payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) Resume() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Resume(&_CommitReveal2L2.TransactOpts)
}

// SetL1FeeCoefficient is a paid mutator transaction binding the contract method 0xd26a4ca7.
//
// Solidity: function setL1FeeCoefficient(uint8 coefficient) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) SetL1FeeCoefficient(opts *bind.TransactOpts, coefficient uint8) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "setL1FeeCoefficient", coefficient)
}

// SetL1FeeCoefficient is a paid mutator transaction binding the contract method 0xd26a4ca7.
//
// Solidity: function setL1FeeCoefficient(uint8 coefficient) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) SetL1FeeCoefficient(coefficient uint8) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SetL1FeeCoefficient(&_CommitReveal2L2.TransactOpts, coefficient)
}

// SetL1FeeCoefficient is a paid mutator transaction binding the contract method 0xd26a4ca7.
//
// Solidity: function setL1FeeCoefficient(uint8 coefficient) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) SetL1FeeCoefficient(coefficient uint8) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SetL1FeeCoefficient(&_CommitReveal2L2.TransactOpts, coefficient)
}

// SetPeriods is a paid mutator transaction binding the contract method 0xbb5e657e.
//
// Solidity: function setPeriods(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) SetPeriods(opts *bind.TransactOpts, offChainSubmissionPeriod *big.Int, requestOrSubmitOrFailDecisionPeriod *big.Int, onChainSubmissionPeriod *big.Int, offChainSubmissionPeriodPerOperator *big.Int, onChainSubmissionPeriodPerOperator *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "setPeriods", offChainSubmissionPeriod, requestOrSubmitOrFailDecisionPeriod, onChainSubmissionPeriod, offChainSubmissionPeriodPerOperator, onChainSubmissionPeriodPerOperator)
}

// SetPeriods is a paid mutator transaction binding the contract method 0xbb5e657e.
//
// Solidity: function setPeriods(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) SetPeriods(offChainSubmissionPeriod *big.Int, requestOrSubmitOrFailDecisionPeriod *big.Int, onChainSubmissionPeriod *big.Int, offChainSubmissionPeriodPerOperator *big.Int, onChainSubmissionPeriodPerOperator *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SetPeriods(&_CommitReveal2L2.TransactOpts, offChainSubmissionPeriod, requestOrSubmitOrFailDecisionPeriod, onChainSubmissionPeriod, offChainSubmissionPeriodPerOperator, onChainSubmissionPeriodPerOperator)
}

// SetPeriods is a paid mutator transaction binding the contract method 0xbb5e657e.
//
// Solidity: function setPeriods(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) SetPeriods(offChainSubmissionPeriod *big.Int, requestOrSubmitOrFailDecisionPeriod *big.Int, onChainSubmissionPeriod *big.Int, offChainSubmissionPeriodPerOperator *big.Int, onChainSubmissionPeriodPerOperator *big.Int) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SetPeriods(&_CommitReveal2L2.TransactOpts, offChainSubmissionPeriod, requestOrSubmitOrFailDecisionPeriod, onChainSubmissionPeriod, offChainSubmissionPeriodPerOperator, onChainSubmissionPeriodPerOperator)
}

// SubmitCo is a paid mutator transaction binding the contract method 0x3576690c.
//
// Solidity: function submitCo(bytes32 co) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) SubmitCo(opts *bind.TransactOpts, co [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "submitCo", co)
}

// SubmitCo is a paid mutator transaction binding the contract method 0x3576690c.
//
// Solidity: function submitCo(bytes32 co) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) SubmitCo(co [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitCo(&_CommitReveal2L2.TransactOpts, co)
}

// SubmitCo is a paid mutator transaction binding the contract method 0x3576690c.
//
// Solidity: function submitCo(bytes32 co) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) SubmitCo(co [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitCo(&_CommitReveal2L2.TransactOpts, co)
}

// SubmitCv is a paid mutator transaction binding the contract method 0xda455fc5.
//
// Solidity: function submitCv(bytes32 cv) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) SubmitCv(opts *bind.TransactOpts, cv [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "submitCv", cv)
}

// SubmitCv is a paid mutator transaction binding the contract method 0xda455fc5.
//
// Solidity: function submitCv(bytes32 cv) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) SubmitCv(cv [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitCv(&_CommitReveal2L2.TransactOpts, cv)
}

// SubmitCv is a paid mutator transaction binding the contract method 0xda455fc5.
//
// Solidity: function submitCv(bytes32 cv) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) SubmitCv(cv [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitCv(&_CommitReveal2L2.TransactOpts, cv)
}

// SubmitMerkleRoot is a paid mutator transaction binding the contract method 0xcda5de83.
//
// Solidity: function submitMerkleRoot(bytes32 merkleRoot) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) SubmitMerkleRoot(opts *bind.TransactOpts, merkleRoot [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "submitMerkleRoot", merkleRoot)
}

// SubmitMerkleRoot is a paid mutator transaction binding the contract method 0xcda5de83.
//
// Solidity: function submitMerkleRoot(bytes32 merkleRoot) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) SubmitMerkleRoot(merkleRoot [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitMerkleRoot(&_CommitReveal2L2.TransactOpts, merkleRoot)
}

// SubmitMerkleRoot is a paid mutator transaction binding the contract method 0xcda5de83.
//
// Solidity: function submitMerkleRoot(bytes32 merkleRoot) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) SubmitMerkleRoot(merkleRoot [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitMerkleRoot(&_CommitReveal2L2.TransactOpts, merkleRoot)
}

// SubmitS is a paid mutator transaction binding the contract method 0x5559689f.
//
// Solidity: function submitS(bytes32 s) returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) SubmitS(opts *bind.TransactOpts, s [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "submitS", s)
}

// SubmitS is a paid mutator transaction binding the contract method 0x5559689f.
//
// Solidity: function submitS(bytes32 s) returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) SubmitS(s [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitS(&_CommitReveal2L2.TransactOpts, s)
}

// SubmitS is a paid mutator transaction binding the contract method 0x5559689f.
//
// Solidity: function submitS(bytes32 s) returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) SubmitS(s [32]byte) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.SubmitS(&_CommitReveal2L2.TransactOpts, s)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.TransferOwnership(&_CommitReveal2L2.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.TransferOwnership(&_CommitReveal2L2.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_CommitReveal2L2 *CommitReveal2L2Transactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitReveal2L2.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_CommitReveal2L2 *CommitReveal2L2Session) Withdraw() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Withdraw(&_CommitReveal2L2.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_CommitReveal2L2 *CommitReveal2L2TransactorSession) Withdraw() (*types.Transaction, error) {
	return _CommitReveal2L2.Contract.Withdraw(&_CommitReveal2L2.TransactOpts)
}

// CommitReveal2L2ActivatedIterator is returned from FilterActivated and is used to iterate over the raw logs and unpacked data for Activated events raised by the CommitReveal2L2 contract.
type CommitReveal2L2ActivatedIterator struct {
	Event *CommitReveal2L2Activated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2ActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2Activated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2Activated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2ActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2ActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2Activated represents a Activated event raised by the CommitReveal2L2 contract.
type CommitReveal2L2Activated struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterActivated is a free log retrieval operation binding the contract event 0x0cc43938d137e7efade6a531f663e78c1fc75257b0d65ffda2fdaf70cb49cdf9.
//
// Solidity: event Activated(address operator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterActivated(opts *bind.FilterOpts) (*CommitReveal2L2ActivatedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "Activated")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2ActivatedIterator{contract: _CommitReveal2L2.contract, event: "Activated", logs: logs, sub: sub}, nil
}

// WatchActivated is a free log subscription operation binding the contract event 0x0cc43938d137e7efade6a531f663e78c1fc75257b0d65ffda2fdaf70cb49cdf9.
//
// Solidity: event Activated(address operator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchActivated(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2Activated) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "Activated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2Activated)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "Activated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseActivated is a log parse operation binding the contract event 0x0cc43938d137e7efade6a531f663e78c1fc75257b0d65ffda2fdaf70cb49cdf9.
//
// Solidity: event Activated(address operator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseActivated(log types.Log) (*CommitReveal2L2Activated, error) {
	event := new(CommitReveal2L2Activated)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "Activated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2CoSubmittedIterator is returned from FilterCoSubmitted and is used to iterate over the raw logs and unpacked data for CoSubmitted events raised by the CommitReveal2L2 contract.
type CommitReveal2L2CoSubmittedIterator struct {
	Event *CommitReveal2L2CoSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2CoSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2CoSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2CoSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2CoSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2CoSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2CoSubmitted represents a CoSubmitted event raised by the CommitReveal2L2 contract.
type CommitReveal2L2CoSubmitted struct {
	Round    *big.Int
	TrialNum *big.Int
	Co       [32]byte
	Index    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCoSubmitted is a free log retrieval operation binding the contract event 0xc294138987faa6e0ebef350caeac5cf5e1eff8dbbe8a158e421601f48674babd.
//
// Solidity: event CoSubmitted(uint256 round, uint256 trialNum, bytes32 co, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterCoSubmitted(opts *bind.FilterOpts) (*CommitReveal2L2CoSubmittedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "CoSubmitted")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2CoSubmittedIterator{contract: _CommitReveal2L2.contract, event: "CoSubmitted", logs: logs, sub: sub}, nil
}

// WatchCoSubmitted is a free log subscription operation binding the contract event 0xc294138987faa6e0ebef350caeac5cf5e1eff8dbbe8a158e421601f48674babd.
//
// Solidity: event CoSubmitted(uint256 round, uint256 trialNum, bytes32 co, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchCoSubmitted(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2CoSubmitted) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "CoSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2CoSubmitted)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "CoSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCoSubmitted is a log parse operation binding the contract event 0xc294138987faa6e0ebef350caeac5cf5e1eff8dbbe8a158e421601f48674babd.
//
// Solidity: event CoSubmitted(uint256 round, uint256 trialNum, bytes32 co, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseCoSubmitted(log types.Log) (*CommitReveal2L2CoSubmitted, error) {
	event := new(CommitReveal2L2CoSubmitted)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "CoSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2CvSubmittedIterator is returned from FilterCvSubmitted and is used to iterate over the raw logs and unpacked data for CvSubmitted events raised by the CommitReveal2L2 contract.
type CommitReveal2L2CvSubmittedIterator struct {
	Event *CommitReveal2L2CvSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2CvSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2CvSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2CvSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2CvSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2CvSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2CvSubmitted represents a CvSubmitted event raised by the CommitReveal2L2 contract.
type CommitReveal2L2CvSubmitted struct {
	Round    *big.Int
	TrialNum *big.Int
	Cv       [32]byte
	Index    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCvSubmitted is a free log retrieval operation binding the contract event 0x6a6385c5eaed19d346ec4f9bd0010cfba4ac1d0407e2e55f959cb8fcac30f873.
//
// Solidity: event CvSubmitted(uint256 round, uint256 trialNum, bytes32 cv, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterCvSubmitted(opts *bind.FilterOpts) (*CommitReveal2L2CvSubmittedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "CvSubmitted")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2CvSubmittedIterator{contract: _CommitReveal2L2.contract, event: "CvSubmitted", logs: logs, sub: sub}, nil
}

// WatchCvSubmitted is a free log subscription operation binding the contract event 0x6a6385c5eaed19d346ec4f9bd0010cfba4ac1d0407e2e55f959cb8fcac30f873.
//
// Solidity: event CvSubmitted(uint256 round, uint256 trialNum, bytes32 cv, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchCvSubmitted(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2CvSubmitted) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "CvSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2CvSubmitted)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "CvSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCvSubmitted is a log parse operation binding the contract event 0x6a6385c5eaed19d346ec4f9bd0010cfba4ac1d0407e2e55f959cb8fcac30f873.
//
// Solidity: event CvSubmitted(uint256 round, uint256 trialNum, bytes32 cv, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseCvSubmitted(log types.Log) (*CommitReveal2L2CvSubmitted, error) {
	event := new(CommitReveal2L2CvSubmitted)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "CvSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2DeActivatedIterator is returned from FilterDeActivated and is used to iterate over the raw logs and unpacked data for DeActivated events raised by the CommitReveal2L2 contract.
type CommitReveal2L2DeActivatedIterator struct {
	Event *CommitReveal2L2DeActivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2DeActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2DeActivated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2DeActivated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2DeActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2DeActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2DeActivated represents a DeActivated event raised by the CommitReveal2L2 contract.
type CommitReveal2L2DeActivated struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeActivated is a free log retrieval operation binding the contract event 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d.
//
// Solidity: event DeActivated(address operator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterDeActivated(opts *bind.FilterOpts) (*CommitReveal2L2DeActivatedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "DeActivated")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2DeActivatedIterator{contract: _CommitReveal2L2.contract, event: "DeActivated", logs: logs, sub: sub}, nil
}

// WatchDeActivated is a free log subscription operation binding the contract event 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d.
//
// Solidity: event DeActivated(address operator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchDeActivated(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2DeActivated) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "DeActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2DeActivated)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "DeActivated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeActivated is a log parse operation binding the contract event 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d.
//
// Solidity: event DeActivated(address operator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseDeActivated(log types.Log) (*CommitReveal2L2DeActivated, error) {
	event := new(CommitReveal2L2DeActivated)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "DeActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2EIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the CommitReveal2L2 contract.
type CommitReveal2L2EIP712DomainChangedIterator struct {
	Event *CommitReveal2L2EIP712DomainChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2EIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2EIP712DomainChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2EIP712DomainChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2EIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2EIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2EIP712DomainChanged represents a EIP712DomainChanged event raised by the CommitReveal2L2 contract.
type CommitReveal2L2EIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*CommitReveal2L2EIP712DomainChangedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2EIP712DomainChangedIterator{contract: _CommitReveal2L2.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2EIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2EIP712DomainChanged)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseEIP712DomainChanged(log types.Log) (*CommitReveal2L2EIP712DomainChanged, error) {
	event := new(CommitReveal2L2EIP712DomainChanged)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2EconomicParametersProposedIterator is returned from FilterEconomicParametersProposed and is used to iterate over the raw logs and unpacked data for EconomicParametersProposed events raised by the CommitReveal2L2 contract.
type CommitReveal2L2EconomicParametersProposedIterator struct {
	Event *CommitReveal2L2EconomicParametersProposed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2EconomicParametersProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2EconomicParametersProposed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2EconomicParametersProposed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2EconomicParametersProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2EconomicParametersProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2EconomicParametersProposed represents a EconomicParametersProposed event raised by the CommitReveal2L2 contract.
type CommitReveal2L2EconomicParametersProposed struct {
	ActivationThreshold *big.Int
	FlatFee             *big.Int
	EffectiveTimestamp  *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterEconomicParametersProposed is a free log retrieval operation binding the contract event 0xdcf23dfc5bc14859d1943fd156abd0fb732347e70c61c56215bbd728307234e2.
//
// Solidity: event EconomicParametersProposed(uint256 activationThreshold, uint256 flatFee, uint256 effectiveTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterEconomicParametersProposed(opts *bind.FilterOpts) (*CommitReveal2L2EconomicParametersProposedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "EconomicParametersProposed")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2EconomicParametersProposedIterator{contract: _CommitReveal2L2.contract, event: "EconomicParametersProposed", logs: logs, sub: sub}, nil
}

// WatchEconomicParametersProposed is a free log subscription operation binding the contract event 0xdcf23dfc5bc14859d1943fd156abd0fb732347e70c61c56215bbd728307234e2.
//
// Solidity: event EconomicParametersProposed(uint256 activationThreshold, uint256 flatFee, uint256 effectiveTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchEconomicParametersProposed(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2EconomicParametersProposed) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "EconomicParametersProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2EconomicParametersProposed)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "EconomicParametersProposed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEconomicParametersProposed is a log parse operation binding the contract event 0xdcf23dfc5bc14859d1943fd156abd0fb732347e70c61c56215bbd728307234e2.
//
// Solidity: event EconomicParametersProposed(uint256 activationThreshold, uint256 flatFee, uint256 effectiveTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseEconomicParametersProposed(log types.Log) (*CommitReveal2L2EconomicParametersProposed, error) {
	event := new(CommitReveal2L2EconomicParametersProposed)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "EconomicParametersProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2EconomicParametersSetIterator is returned from FilterEconomicParametersSet and is used to iterate over the raw logs and unpacked data for EconomicParametersSet events raised by the CommitReveal2L2 contract.
type CommitReveal2L2EconomicParametersSetIterator struct {
	Event *CommitReveal2L2EconomicParametersSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2EconomicParametersSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2EconomicParametersSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2EconomicParametersSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2EconomicParametersSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2EconomicParametersSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2EconomicParametersSet represents a EconomicParametersSet event raised by the CommitReveal2L2 contract.
type CommitReveal2L2EconomicParametersSet struct {
	ActivationThreshold *big.Int
	FlatFee             *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterEconomicParametersSet is a free log retrieval operation binding the contract event 0x08f0774e7eb69e2d6a7cf2192cbf9c6f519a40bcfa16ff60d3f18496585e46dc.
//
// Solidity: event EconomicParametersSet(uint256 activationThreshold, uint256 flatFee)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterEconomicParametersSet(opts *bind.FilterOpts) (*CommitReveal2L2EconomicParametersSetIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "EconomicParametersSet")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2EconomicParametersSetIterator{contract: _CommitReveal2L2.contract, event: "EconomicParametersSet", logs: logs, sub: sub}, nil
}

// WatchEconomicParametersSet is a free log subscription operation binding the contract event 0x08f0774e7eb69e2d6a7cf2192cbf9c6f519a40bcfa16ff60d3f18496585e46dc.
//
// Solidity: event EconomicParametersSet(uint256 activationThreshold, uint256 flatFee)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchEconomicParametersSet(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2EconomicParametersSet) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "EconomicParametersSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2EconomicParametersSet)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "EconomicParametersSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEconomicParametersSet is a log parse operation binding the contract event 0x08f0774e7eb69e2d6a7cf2192cbf9c6f519a40bcfa16ff60d3f18496585e46dc.
//
// Solidity: event EconomicParametersSet(uint256 activationThreshold, uint256 flatFee)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseEconomicParametersSet(log types.Log) (*CommitReveal2L2EconomicParametersSet, error) {
	event := new(CommitReveal2L2EconomicParametersSet)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "EconomicParametersSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2GasParametersProposedIterator is returned from FilterGasParametersProposed and is used to iterate over the raw logs and unpacked data for GasParametersProposed events raised by the CommitReveal2L2 contract.
type CommitReveal2L2GasParametersProposedIterator struct {
	Event *CommitReveal2L2GasParametersProposed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2GasParametersProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2GasParametersProposed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2GasParametersProposed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2GasParametersProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2GasParametersProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2GasParametersProposed represents a GasParametersProposed event raised by the CommitReveal2L2 contract.
type CommitReveal2L2GasParametersProposed struct {
	GasUsedMerkleRootSubAndGenRandNumA          *big.Int
	GasUsedMerkleRootSubAndGenRandNumB          *big.Int
	MaxCallbackGasLimit                         *big.Int
	GetL1UpperBoundGasUsedWhenCalldataSize4     *big.Int
	FailToRequestCvOrSubmitMerkleRootGasUsed    *big.Int
	FailToSubmitMerkleRootAfterDisputeGasUsed   *big.Int
	FailToRequestSOrGenerateRandomNumberGasUsed *big.Int
	FailToSubmitSGasUsed                        *big.Int
	FailToSubmitCoGasUsedBaseA                  uint32
	FailToSubmitCvGasUsedBaseA                  uint32
	FailToSubmitGasUsedBaseB                    uint32
	PerOperatorIncreaseGasUsedA                 uint32
	PerOperatorIncreaseGasUsedB                 uint32
	PerAdditionalDidntSubmitGasUsedA            uint32
	PerAdditionalDidntSubmitGasUsedB            uint32
	PerRequestedIncreaseGasUsed                 uint32
	EffectiveTimestamp                          *big.Int
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterGasParametersProposed is a free log retrieval operation binding the contract event 0xac29dedddb8466e143ff09a21b0181b73354eae633cc2787fb6dd4c3b50dfbe2.
//
// Solidity: event GasParametersProposed(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumB, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed, uint256 effectiveTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterGasParametersProposed(opts *bind.FilterOpts) (*CommitReveal2L2GasParametersProposedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "GasParametersProposed")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2GasParametersProposedIterator{contract: _CommitReveal2L2.contract, event: "GasParametersProposed", logs: logs, sub: sub}, nil
}

// WatchGasParametersProposed is a free log subscription operation binding the contract event 0xac29dedddb8466e143ff09a21b0181b73354eae633cc2787fb6dd4c3b50dfbe2.
//
// Solidity: event GasParametersProposed(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumB, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed, uint256 effectiveTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchGasParametersProposed(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2GasParametersProposed) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "GasParametersProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2GasParametersProposed)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "GasParametersProposed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGasParametersProposed is a log parse operation binding the contract event 0xac29dedddb8466e143ff09a21b0181b73354eae633cc2787fb6dd4c3b50dfbe2.
//
// Solidity: event GasParametersProposed(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumB, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed, uint256 effectiveTimestamp)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseGasParametersProposed(log types.Log) (*CommitReveal2L2GasParametersProposed, error) {
	event := new(CommitReveal2L2GasParametersProposed)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "GasParametersProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2GasParametersSetIterator is returned from FilterGasParametersSet and is used to iterate over the raw logs and unpacked data for GasParametersSet events raised by the CommitReveal2L2 contract.
type CommitReveal2L2GasParametersSetIterator struct {
	Event *CommitReveal2L2GasParametersSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2GasParametersSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2GasParametersSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2GasParametersSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2GasParametersSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2GasParametersSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2GasParametersSet represents a GasParametersSet event raised by the CommitReveal2L2 contract.
type CommitReveal2L2GasParametersSet struct {
	GasUsedMerkleRootSubAndGenRandNumA          *big.Int
	GasUsedMerkleRootSubAndGenRandNumB          *big.Int
	MaxCallbackGasLimit                         *big.Int
	GetL1UpperBoundGasUsedWhenCalldataSize4     *big.Int
	FailToRequestCvOrSubmitMerkleRootGasUsed    *big.Int
	FailToSubmitMerkleRootAfterDisputeGasUsed   *big.Int
	FailToRequestSOrGenerateRandomNumberGasUsed *big.Int
	FailToSubmitSGasUsed                        *big.Int
	FailToSubmitCoGasUsedBaseA                  uint32
	FailToSubmitCvGasUsedBaseA                  uint32
	FailToSubmitGasUsedBaseB                    uint32
	PerOperatorIncreaseGasUsedA                 uint32
	PerOperatorIncreaseGasUsedB                 uint32
	PerAdditionalDidntSubmitGasUsedA            uint32
	PerAdditionalDidntSubmitGasUsedB            uint32
	PerRequestedIncreaseGasUsed                 uint32
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterGasParametersSet is a free log retrieval operation binding the contract event 0x8d09171105499771f96d6d39dcdda061a70fd18e5eafd65881c2158c55f94e1d.
//
// Solidity: event GasParametersSet(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumB, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterGasParametersSet(opts *bind.FilterOpts) (*CommitReveal2L2GasParametersSetIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "GasParametersSet")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2GasParametersSetIterator{contract: _CommitReveal2L2.contract, event: "GasParametersSet", logs: logs, sub: sub}, nil
}

// WatchGasParametersSet is a free log subscription operation binding the contract event 0x8d09171105499771f96d6d39dcdda061a70fd18e5eafd65881c2158c55f94e1d.
//
// Solidity: event GasParametersSet(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumB, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchGasParametersSet(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2GasParametersSet) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "GasParametersSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2GasParametersSet)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "GasParametersSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGasParametersSet is a log parse operation binding the contract event 0x8d09171105499771f96d6d39dcdda061a70fd18e5eafd65881c2158c55f94e1d.
//
// Solidity: event GasParametersSet(uint128 gasUsedMerkleRootSubAndGenRandNumA, uint128 gasUsedMerkleRootSubAndGenRandNumB, uint256 maxCallbackGasLimit, uint48 getL1UpperBoundGasUsedWhenCalldataSize4, uint48 failToRequestCvOrSubmitMerkleRootGasUsed, uint48 failToSubmitMerkleRootAfterDisputeGasUsed, uint48 failToRequestSOrGenerateRandomNumberGasUsed, uint48 failToSubmitSGasUsed, uint32 failToSubmitCoGasUsedBaseA, uint32 failToSubmitCvGasUsedBaseA, uint32 failToSubmitGasUsedBaseB, uint32 perOperatorIncreaseGasUsedA, uint32 perOperatorIncreaseGasUsedB, uint32 perAdditionalDidntSubmitGasUsedA, uint32 perAdditionalDidntSubmitGasUsedB, uint32 perRequestedIncreaseGasUsed)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseGasParametersSet(log types.Log) (*CommitReveal2L2GasParametersSet, error) {
	event := new(CommitReveal2L2GasParametersSet)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "GasParametersSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2L1FeeCalculationSetIterator is returned from FilterL1FeeCalculationSet and is used to iterate over the raw logs and unpacked data for L1FeeCalculationSet events raised by the CommitReveal2L2 contract.
type CommitReveal2L2L1FeeCalculationSetIterator struct {
	Event *CommitReveal2L2L1FeeCalculationSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2L1FeeCalculationSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2L1FeeCalculationSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2L1FeeCalculationSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2L1FeeCalculationSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2L1FeeCalculationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2L1FeeCalculationSet represents a L1FeeCalculationSet event raised by the CommitReveal2L2 contract.
type CommitReveal2L2L1FeeCalculationSet struct {
	Coefficient uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterL1FeeCalculationSet is a free log retrieval operation binding the contract event 0x8b20b84893eb600b867c893a944643ee6c0ce967aa98367fee46d84c56eec022.
//
// Solidity: event L1FeeCalculationSet(uint8 coefficient)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterL1FeeCalculationSet(opts *bind.FilterOpts) (*CommitReveal2L2L1FeeCalculationSetIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "L1FeeCalculationSet")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2L1FeeCalculationSetIterator{contract: _CommitReveal2L2.contract, event: "L1FeeCalculationSet", logs: logs, sub: sub}, nil
}

// WatchL1FeeCalculationSet is a free log subscription operation binding the contract event 0x8b20b84893eb600b867c893a944643ee6c0ce967aa98367fee46d84c56eec022.
//
// Solidity: event L1FeeCalculationSet(uint8 coefficient)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchL1FeeCalculationSet(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2L1FeeCalculationSet) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "L1FeeCalculationSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2L1FeeCalculationSet)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "L1FeeCalculationSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseL1FeeCalculationSet is a log parse operation binding the contract event 0x8b20b84893eb600b867c893a944643ee6c0ce967aa98367fee46d84c56eec022.
//
// Solidity: event L1FeeCalculationSet(uint8 coefficient)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseL1FeeCalculationSet(log types.Log) (*CommitReveal2L2L1FeeCalculationSet, error) {
	event := new(CommitReveal2L2L1FeeCalculationSet)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "L1FeeCalculationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2MerkleRootSubmittedIterator is returned from FilterMerkleRootSubmitted and is used to iterate over the raw logs and unpacked data for MerkleRootSubmitted events raised by the CommitReveal2L2 contract.
type CommitReveal2L2MerkleRootSubmittedIterator struct {
	Event *CommitReveal2L2MerkleRootSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2MerkleRootSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2MerkleRootSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2MerkleRootSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2MerkleRootSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2MerkleRootSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2MerkleRootSubmitted represents a MerkleRootSubmitted event raised by the CommitReveal2L2 contract.
type CommitReveal2L2MerkleRootSubmitted struct {
	Round      *big.Int
	TrialNum   *big.Int
	MerkleRoot [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootSubmitted is a free log retrieval operation binding the contract event 0x45b19880b523c6750f7f39fca8d77d51101b315495adc482994a4fa2a8294466.
//
// Solidity: event MerkleRootSubmitted(uint256 round, uint256 trialNum, bytes32 merkleRoot)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterMerkleRootSubmitted(opts *bind.FilterOpts) (*CommitReveal2L2MerkleRootSubmittedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "MerkleRootSubmitted")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2MerkleRootSubmittedIterator{contract: _CommitReveal2L2.contract, event: "MerkleRootSubmitted", logs: logs, sub: sub}, nil
}

// WatchMerkleRootSubmitted is a free log subscription operation binding the contract event 0x45b19880b523c6750f7f39fca8d77d51101b315495adc482994a4fa2a8294466.
//
// Solidity: event MerkleRootSubmitted(uint256 round, uint256 trialNum, bytes32 merkleRoot)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchMerkleRootSubmitted(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2MerkleRootSubmitted) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "MerkleRootSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2MerkleRootSubmitted)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "MerkleRootSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMerkleRootSubmitted is a log parse operation binding the contract event 0x45b19880b523c6750f7f39fca8d77d51101b315495adc482994a4fa2a8294466.
//
// Solidity: event MerkleRootSubmitted(uint256 round, uint256 trialNum, bytes32 merkleRoot)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseMerkleRootSubmitted(log types.Log) (*CommitReveal2L2MerkleRootSubmitted, error) {
	event := new(CommitReveal2L2MerkleRootSubmitted)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "MerkleRootSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2OwnershipHandoverCanceledIterator is returned from FilterOwnershipHandoverCanceled and is used to iterate over the raw logs and unpacked data for OwnershipHandoverCanceled events raised by the CommitReveal2L2 contract.
type CommitReveal2L2OwnershipHandoverCanceledIterator struct {
	Event *CommitReveal2L2OwnershipHandoverCanceled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2OwnershipHandoverCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2OwnershipHandoverCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2OwnershipHandoverCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2OwnershipHandoverCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2OwnershipHandoverCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2OwnershipHandoverCanceled represents a OwnershipHandoverCanceled event raised by the CommitReveal2L2 contract.
type CommitReveal2L2OwnershipHandoverCanceled struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverCanceled is a free log retrieval operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterOwnershipHandoverCanceled(opts *bind.FilterOpts, pendingOwner []common.Address) (*CommitReveal2L2OwnershipHandoverCanceledIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2OwnershipHandoverCanceledIterator{contract: _CommitReveal2L2.contract, event: "OwnershipHandoverCanceled", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverCanceled is a free log subscription operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchOwnershipHandoverCanceled(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2OwnershipHandoverCanceled, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2OwnershipHandoverCanceled)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipHandoverCanceled is a log parse operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseOwnershipHandoverCanceled(log types.Log) (*CommitReveal2L2OwnershipHandoverCanceled, error) {
	event := new(CommitReveal2L2OwnershipHandoverCanceled)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2OwnershipHandoverRequestedIterator is returned from FilterOwnershipHandoverRequested and is used to iterate over the raw logs and unpacked data for OwnershipHandoverRequested events raised by the CommitReveal2L2 contract.
type CommitReveal2L2OwnershipHandoverRequestedIterator struct {
	Event *CommitReveal2L2OwnershipHandoverRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2OwnershipHandoverRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2OwnershipHandoverRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2OwnershipHandoverRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2OwnershipHandoverRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2OwnershipHandoverRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2OwnershipHandoverRequested represents a OwnershipHandoverRequested event raised by the CommitReveal2L2 contract.
type CommitReveal2L2OwnershipHandoverRequested struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverRequested is a free log retrieval operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterOwnershipHandoverRequested(opts *bind.FilterOpts, pendingOwner []common.Address) (*CommitReveal2L2OwnershipHandoverRequestedIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2OwnershipHandoverRequestedIterator{contract: _CommitReveal2L2.contract, event: "OwnershipHandoverRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverRequested is a free log subscription operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchOwnershipHandoverRequested(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2OwnershipHandoverRequested, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2OwnershipHandoverRequested)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipHandoverRequested is a log parse operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseOwnershipHandoverRequested(log types.Log) (*CommitReveal2L2OwnershipHandoverRequested, error) {
	event := new(CommitReveal2L2OwnershipHandoverRequested)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CommitReveal2L2 contract.
type CommitReveal2L2OwnershipTransferredIterator struct {
	Event *CommitReveal2L2OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2OwnershipTransferred represents a OwnershipTransferred event raised by the CommitReveal2L2 contract.
type CommitReveal2L2OwnershipTransferred struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*CommitReveal2L2OwnershipTransferredIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2OwnershipTransferredIterator{contract: _CommitReveal2L2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2OwnershipTransferred, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2OwnershipTransferred)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseOwnershipTransferred(log types.Log) (*CommitReveal2L2OwnershipTransferred, error) {
	event := new(CommitReveal2L2OwnershipTransferred)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2PeriodsSetIterator is returned from FilterPeriodsSet and is used to iterate over the raw logs and unpacked data for PeriodsSet events raised by the CommitReveal2L2 contract.
type CommitReveal2L2PeriodsSetIterator struct {
	Event *CommitReveal2L2PeriodsSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2PeriodsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2PeriodsSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2PeriodsSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2PeriodsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2PeriodsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2PeriodsSet represents a PeriodsSet event raised by the CommitReveal2L2 contract.
type CommitReveal2L2PeriodsSet struct {
	OffChainSubmissionPeriod            *big.Int
	RequestOrSubmitOrFailDecisionPeriod *big.Int
	OnChainSubmissionPeriod             *big.Int
	OffChainSubmissionPeriodPerOperator *big.Int
	OnChainSubmissionPeriodPerOperator  *big.Int
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterPeriodsSet is a free log retrieval operation binding the contract event 0xe0fd8eabd2cc23ea87b43a00ac588c61789ad28d3edfeb76613f623fa1f6bd08.
//
// Solidity: event PeriodsSet(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterPeriodsSet(opts *bind.FilterOpts) (*CommitReveal2L2PeriodsSetIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "PeriodsSet")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2PeriodsSetIterator{contract: _CommitReveal2L2.contract, event: "PeriodsSet", logs: logs, sub: sub}, nil
}

// WatchPeriodsSet is a free log subscription operation binding the contract event 0xe0fd8eabd2cc23ea87b43a00ac588c61789ad28d3edfeb76613f623fa1f6bd08.
//
// Solidity: event PeriodsSet(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchPeriodsSet(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2PeriodsSet) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "PeriodsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2PeriodsSet)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "PeriodsSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePeriodsSet is a log parse operation binding the contract event 0xe0fd8eabd2cc23ea87b43a00ac588c61789ad28d3edfeb76613f623fa1f6bd08.
//
// Solidity: event PeriodsSet(uint256 offChainSubmissionPeriod, uint256 requestOrSubmitOrFailDecisionPeriod, uint256 onChainSubmissionPeriod, uint256 offChainSubmissionPeriodPerOperator, uint256 onChainSubmissionPeriodPerOperator)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParsePeriodsSet(log types.Log) (*CommitReveal2L2PeriodsSet, error) {
	event := new(CommitReveal2L2PeriodsSet)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "PeriodsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2RequestedToSubmitCoIterator is returned from FilterRequestedToSubmitCo and is used to iterate over the raw logs and unpacked data for RequestedToSubmitCo events raised by the CommitReveal2L2 contract.
type CommitReveal2L2RequestedToSubmitCoIterator struct {
	Event *CommitReveal2L2RequestedToSubmitCo // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2RequestedToSubmitCoIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2RequestedToSubmitCo)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2RequestedToSubmitCo)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2RequestedToSubmitCoIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2RequestedToSubmitCoIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2RequestedToSubmitCo represents a RequestedToSubmitCo event raised by the CommitReveal2L2 contract.
type CommitReveal2L2RequestedToSubmitCo struct {
	Round         *big.Int
	TrialNum      *big.Int
	IndicesLength *big.Int
	PackedIndices *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRequestedToSubmitCo is a free log retrieval operation binding the contract event 0xd4cc5cd95f180f10aaacba0729abc069b8080ec3a7e8e41856decb17bdc28ece.
//
// Solidity: event RequestedToSubmitCo(uint256 round, uint256 trialNum, uint256 indicesLength, uint256 packedIndices)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterRequestedToSubmitCo(opts *bind.FilterOpts) (*CommitReveal2L2RequestedToSubmitCoIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "RequestedToSubmitCo")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2RequestedToSubmitCoIterator{contract: _CommitReveal2L2.contract, event: "RequestedToSubmitCo", logs: logs, sub: sub}, nil
}

// WatchRequestedToSubmitCo is a free log subscription operation binding the contract event 0xd4cc5cd95f180f10aaacba0729abc069b8080ec3a7e8e41856decb17bdc28ece.
//
// Solidity: event RequestedToSubmitCo(uint256 round, uint256 trialNum, uint256 indicesLength, uint256 packedIndices)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchRequestedToSubmitCo(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2RequestedToSubmitCo) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "RequestedToSubmitCo")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2RequestedToSubmitCo)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "RequestedToSubmitCo", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestedToSubmitCo is a log parse operation binding the contract event 0xd4cc5cd95f180f10aaacba0729abc069b8080ec3a7e8e41856decb17bdc28ece.
//
// Solidity: event RequestedToSubmitCo(uint256 round, uint256 trialNum, uint256 indicesLength, uint256 packedIndices)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseRequestedToSubmitCo(log types.Log) (*CommitReveal2L2RequestedToSubmitCo, error) {
	event := new(CommitReveal2L2RequestedToSubmitCo)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "RequestedToSubmitCo", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2RequestedToSubmitCvIterator is returned from FilterRequestedToSubmitCv and is used to iterate over the raw logs and unpacked data for RequestedToSubmitCv events raised by the CommitReveal2L2 contract.
type CommitReveal2L2RequestedToSubmitCvIterator struct {
	Event *CommitReveal2L2RequestedToSubmitCv // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2RequestedToSubmitCvIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2RequestedToSubmitCv)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2RequestedToSubmitCv)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2RequestedToSubmitCvIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2RequestedToSubmitCvIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2RequestedToSubmitCv represents a RequestedToSubmitCv event raised by the CommitReveal2L2 contract.
type CommitReveal2L2RequestedToSubmitCv struct {
	Round                         *big.Int
	TrialNum                      *big.Int
	PackedIndicesAscendingFromLSB *big.Int
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterRequestedToSubmitCv is a free log retrieval operation binding the contract event 0x16759d80d11394de93184cfeb4e91cf57282cef239f68ed141c496600454f757.
//
// Solidity: event RequestedToSubmitCv(uint256 round, uint256 trialNum, uint256 packedIndicesAscendingFromLSB)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterRequestedToSubmitCv(opts *bind.FilterOpts) (*CommitReveal2L2RequestedToSubmitCvIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "RequestedToSubmitCv")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2RequestedToSubmitCvIterator{contract: _CommitReveal2L2.contract, event: "RequestedToSubmitCv", logs: logs, sub: sub}, nil
}

// WatchRequestedToSubmitCv is a free log subscription operation binding the contract event 0x16759d80d11394de93184cfeb4e91cf57282cef239f68ed141c496600454f757.
//
// Solidity: event RequestedToSubmitCv(uint256 round, uint256 trialNum, uint256 packedIndicesAscendingFromLSB)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchRequestedToSubmitCv(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2RequestedToSubmitCv) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "RequestedToSubmitCv")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2RequestedToSubmitCv)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "RequestedToSubmitCv", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestedToSubmitCv is a log parse operation binding the contract event 0x16759d80d11394de93184cfeb4e91cf57282cef239f68ed141c496600454f757.
//
// Solidity: event RequestedToSubmitCv(uint256 round, uint256 trialNum, uint256 packedIndicesAscendingFromLSB)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseRequestedToSubmitCv(log types.Log) (*CommitReveal2L2RequestedToSubmitCv, error) {
	event := new(CommitReveal2L2RequestedToSubmitCv)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "RequestedToSubmitCv", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2RequestedToSubmitSFromIndexKIterator is returned from FilterRequestedToSubmitSFromIndexK and is used to iterate over the raw logs and unpacked data for RequestedToSubmitSFromIndexK events raised by the CommitReveal2L2 contract.
type CommitReveal2L2RequestedToSubmitSFromIndexKIterator struct {
	Event *CommitReveal2L2RequestedToSubmitSFromIndexK // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2RequestedToSubmitSFromIndexKIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2RequestedToSubmitSFromIndexK)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2RequestedToSubmitSFromIndexK)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2RequestedToSubmitSFromIndexKIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2RequestedToSubmitSFromIndexKIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2RequestedToSubmitSFromIndexK represents a RequestedToSubmitSFromIndexK event raised by the CommitReveal2L2 contract.
type CommitReveal2L2RequestedToSubmitSFromIndexK struct {
	Round    *big.Int
	TrialNum *big.Int
	IndexK   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRequestedToSubmitSFromIndexK is a free log retrieval operation binding the contract event 0x583f939e9612a50da8a140b5e7247ff7c3c899c45e4051a5ba045abea6177f08.
//
// Solidity: event RequestedToSubmitSFromIndexK(uint256 round, uint256 trialNum, uint256 indexK)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterRequestedToSubmitSFromIndexK(opts *bind.FilterOpts) (*CommitReveal2L2RequestedToSubmitSFromIndexKIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "RequestedToSubmitSFromIndexK")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2RequestedToSubmitSFromIndexKIterator{contract: _CommitReveal2L2.contract, event: "RequestedToSubmitSFromIndexK", logs: logs, sub: sub}, nil
}

// WatchRequestedToSubmitSFromIndexK is a free log subscription operation binding the contract event 0x583f939e9612a50da8a140b5e7247ff7c3c899c45e4051a5ba045abea6177f08.
//
// Solidity: event RequestedToSubmitSFromIndexK(uint256 round, uint256 trialNum, uint256 indexK)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchRequestedToSubmitSFromIndexK(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2RequestedToSubmitSFromIndexK) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "RequestedToSubmitSFromIndexK")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2RequestedToSubmitSFromIndexK)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "RequestedToSubmitSFromIndexK", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestedToSubmitSFromIndexK is a log parse operation binding the contract event 0x583f939e9612a50da8a140b5e7247ff7c3c899c45e4051a5ba045abea6177f08.
//
// Solidity: event RequestedToSubmitSFromIndexK(uint256 round, uint256 trialNum, uint256 indexK)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseRequestedToSubmitSFromIndexK(log types.Log) (*CommitReveal2L2RequestedToSubmitSFromIndexK, error) {
	event := new(CommitReveal2L2RequestedToSubmitSFromIndexK)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "RequestedToSubmitSFromIndexK", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2SSubmittedIterator is returned from FilterSSubmitted and is used to iterate over the raw logs and unpacked data for SSubmitted events raised by the CommitReveal2L2 contract.
type CommitReveal2L2SSubmittedIterator struct {
	Event *CommitReveal2L2SSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2SSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2SSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2SSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2SSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2SSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2SSubmitted represents a SSubmitted event raised by the CommitReveal2L2 contract.
type CommitReveal2L2SSubmitted struct {
	Round    *big.Int
	TrialNum *big.Int
	S        [32]byte
	Index    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSSubmitted is a free log retrieval operation binding the contract event 0xfa070a58e2c77080acd5c2b1819669eb194bbeeca6f680a31a2076510be5a7b1.
//
// Solidity: event SSubmitted(uint256 round, uint256 trialNum, bytes32 s, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterSSubmitted(opts *bind.FilterOpts) (*CommitReveal2L2SSubmittedIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "SSubmitted")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2SSubmittedIterator{contract: _CommitReveal2L2.contract, event: "SSubmitted", logs: logs, sub: sub}, nil
}

// WatchSSubmitted is a free log subscription operation binding the contract event 0xfa070a58e2c77080acd5c2b1819669eb194bbeeca6f680a31a2076510be5a7b1.
//
// Solidity: event SSubmitted(uint256 round, uint256 trialNum, bytes32 s, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchSSubmitted(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2SSubmitted) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "SSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2SSubmitted)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "SSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSSubmitted is a log parse operation binding the contract event 0xfa070a58e2c77080acd5c2b1819669eb194bbeeca6f680a31a2076510be5a7b1.
//
// Solidity: event SSubmitted(uint256 round, uint256 trialNum, bytes32 s, uint256 index)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseSSubmitted(log types.Log) (*CommitReveal2L2SSubmitted, error) {
	event := new(CommitReveal2L2SSubmitted)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "SSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommitReveal2L2StatusIterator is returned from FilterStatus and is used to iterate over the raw logs and unpacked data for Status events raised by the CommitReveal2L2 contract.
type CommitReveal2L2StatusIterator struct {
	Event *CommitReveal2L2Status // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CommitReveal2L2StatusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitReveal2L2Status)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CommitReveal2L2Status)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CommitReveal2L2StatusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitReveal2L2StatusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitReveal2L2Status represents a Status event raised by the CommitReveal2L2 contract.
type CommitReveal2L2Status struct {
	CurRound    *big.Int
	CurTrialNum *big.Int
	CurState    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStatus is a free log retrieval operation binding the contract event 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813.
//
// Solidity: event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) FilterStatus(opts *bind.FilterOpts) (*CommitReveal2L2StatusIterator, error) {

	logs, sub, err := _CommitReveal2L2.contract.FilterLogs(opts, "Status")
	if err != nil {
		return nil, err
	}
	return &CommitReveal2L2StatusIterator{contract: _CommitReveal2L2.contract, event: "Status", logs: logs, sub: sub}, nil
}

// WatchStatus is a free log subscription operation binding the contract event 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813.
//
// Solidity: event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) WatchStatus(opts *bind.WatchOpts, sink chan<- *CommitReveal2L2Status) (event.Subscription, error) {

	logs, sub, err := _CommitReveal2L2.contract.WatchLogs(opts, "Status")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitReveal2L2Status)
				if err := _CommitReveal2L2.contract.UnpackLog(event, "Status", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStatus is a log parse operation binding the contract event 0xd42cacab4700e77b08a2d33cc97d95a9cb985cdfca3a206cfa4990da46dd1813.
//
// Solidity: event Status(uint256 curRound, uint256 curTrialNum, uint256 curState)
func (_CommitReveal2L2 *CommitReveal2L2Filterer) ParseStatus(log types.Log) (*CommitReveal2L2Status, error) {
	event := new(CommitReveal2L2Status)
	if err := _CommitReveal2L2.contract.UnpackLog(event, "Status", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
