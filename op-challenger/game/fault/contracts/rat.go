// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
	_ = abi.ConvertType
)

// RATChallengerInfo is an auto generated low-level Go binding around an user-defined struct.
type RATChallengerInfo struct {
	StakingAmount      *big.Int
	TotalSlashedAmount *big.Int
	ValidatorIndex     uint32
	IsValid            bool
}

// RATMetaData contains all meta data concerning the RAT contract.
var RATMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"attentionTests\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bondAmount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"challengerAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"l1BlockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"evidenceSubmitted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"challengers\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"stakingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalSlashedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"disputeGameFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDisputeGameFactory\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"evidenceSubmissionPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChallengerInfo\",\"inputs\":[{\"name\":\"_challenger\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRAT.ChallengerInfo\",\"components\":[{\"name\":\"stakingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalSlashedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidChallengerCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_disputeGameFactory\",\"type\":\"address\",\"internalType\":\"contractIDisputeGameFactory\"},{\"name\":\"_perTestBondAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_evidenceSubmissionPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minimumStakingBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_ratTriggerProbability\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_manager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"minimumStakingBalance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"perTestBondAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxyAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIProxyAdmin\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxyAdminOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ratManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ratTriggerProbability\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resolveClaim\",\"inputs\":[{\"name\":\"_claimant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEvidenceSubmissionPeriod\",\"inputs\":[{\"name\":\"_period\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinimumStakingBalance\",\"inputs\":[{\"name\":\"_balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPerTestBondAmount\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRatTriggerProbability\",\"inputs\":[{\"name\":\"_probability\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"submitCorrectEvidence\",\"inputs\":[{\"name\":\"_gameAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_proofLV\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_proofRV\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"triggerAttentionTest\",\"inputs\":[{\"name\":\"_gameAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_stateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validChallengers\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AttentionTriggered\",\"inputs\":[{\"name\":\"gameAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BondRefunded\",\"inputs\":[{\"name\":\"gameAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"refundedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChallengerStaked\",\"inputs\":[{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CorrectEvidenceSubmitted\",\"inputs\":[{\"name\":\"gameAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"restoredAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AttentionTestNotExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChallengerNotExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EvidenceAlreadySubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EvidenceSubmissionExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientStakingAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChallengerAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoValidChallengers\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotDisputeGameFactory\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRatManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProxyAdminOwnedBase_NotProxyAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProxyAdminOwnedBase_NotProxyAdminOrProxyAdminOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProxyAdminOwnedBase_NotProxyAdminOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProxyAdminOwnedBase_NotResolvedDelegateProxy\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProxyAdminOwnedBase_NotSharedProxyAdminOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProxyAdminOwnedBase_ProxyAdminNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReinitializableBase_ZeroInitVersion\",\"inputs\":[]}]",
}

// RATABI is the input ABI used to generate the binding from.
// Deprecated: Use RATMetaData.ABI instead.
var RATABI = RATMetaData.ABI

// RAT is an auto generated Go binding around an Ethereum contract.
type RAT struct {
	RATCaller     // Read-only binding to the contract
	RATTransactor // Write-only binding to the contract
	RATFilterer   // Log filterer for contract events
}

// RATCaller is an auto generated read-only Go binding around an Ethereum contract.
type RATCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RATTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RATTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RATFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RATFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RATSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RATSession struct {
	Contract     *RAT              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RATCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RATCallerSession struct {
	Contract *RATCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RATTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RATTransactorSession struct {
	Contract     *RATTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RATRaw is an auto generated low-level Go binding around an Ethereum contract.
type RATRaw struct {
	Contract *RAT // Generic contract binding to access the raw methods on
}

// RATCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RATCallerRaw struct {
	Contract *RATCaller // Generic read-only contract binding to access the raw methods on
}

// RATTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RATTransactorRaw struct {
	Contract *RATTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRAT creates a new instance of RAT, bound to a specific deployed contract.
func NewRAT(address common.Address, backend bind.ContractBackend) (*RAT, error) {
	contract, err := bindRAT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RAT{RATCaller: RATCaller{contract: contract}, RATTransactor: RATTransactor{contract: contract}, RATFilterer: RATFilterer{contract: contract}}, nil
}

// NewRATCaller creates a new read-only instance of RAT, bound to a specific deployed contract.
func NewRATCaller(address common.Address, caller bind.ContractCaller) (*RATCaller, error) {
	contract, err := bindRAT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RATCaller{contract: contract}, nil
}

// NewRATTransactor creates a new write-only instance of RAT, bound to a specific deployed contract.
func NewRATTransactor(address common.Address, transactor bind.ContractTransactor) (*RATTransactor, error) {
	contract, err := bindRAT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RATTransactor{contract: contract}, nil
}

// NewRATFilterer creates a new log filterer instance of RAT, bound to a specific deployed contract.
func NewRATFilterer(address common.Address, filterer bind.ContractFilterer) (*RATFilterer, error) {
	contract, err := bindRAT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RATFilterer{contract: contract}, nil
}

// bindRAT binds a generic wrapper to an already deployed contract.
func bindRAT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RATMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RAT *RATRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RAT.Contract.RATCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RAT *RATRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RAT.Contract.RATTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RAT *RATRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RAT.Contract.RATTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RAT *RATCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RAT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RAT *RATTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RAT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RAT *RATTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RAT.Contract.contract.Transact(opts, method, params...)
}

// AttentionTests is a free data retrieval call binding the contract method 0xfba1d220.
//
// Solidity: function attentionTests(address ) view returns(bytes32 stateRoot, uint96 bondAmount, address challengerAddress, uint64 l1BlockNumber, bool evidenceSubmitted)
func (_RAT *RATCaller) AttentionTests(opts *bind.CallOpts, arg0 common.Address) (struct {
	StateRoot         [32]byte
	BondAmount        *big.Int
	ChallengerAddress common.Address
	L1BlockNumber     uint64
	EvidenceSubmitted bool
}, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "attentionTests", arg0)

	outstruct := new(struct {
		StateRoot         [32]byte
		BondAmount        *big.Int
		ChallengerAddress common.Address
		L1BlockNumber     uint64
		EvidenceSubmitted bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StateRoot = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.BondAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ChallengerAddress = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.L1BlockNumber = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.EvidenceSubmitted = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// AttentionTests is a free data retrieval call binding the contract method 0xfba1d220.
//
// Solidity: function attentionTests(address ) view returns(bytes32 stateRoot, uint96 bondAmount, address challengerAddress, uint64 l1BlockNumber, bool evidenceSubmitted)
func (_RAT *RATSession) AttentionTests(arg0 common.Address) (struct {
	StateRoot         [32]byte
	BondAmount        *big.Int
	ChallengerAddress common.Address
	L1BlockNumber     uint64
	EvidenceSubmitted bool
}, error) {
	return _RAT.Contract.AttentionTests(&_RAT.CallOpts, arg0)
}

// AttentionTests is a free data retrieval call binding the contract method 0xfba1d220.
//
// Solidity: function attentionTests(address ) view returns(bytes32 stateRoot, uint96 bondAmount, address challengerAddress, uint64 l1BlockNumber, bool evidenceSubmitted)
func (_RAT *RATCallerSession) AttentionTests(arg0 common.Address) (struct {
	StateRoot         [32]byte
	BondAmount        *big.Int
	ChallengerAddress common.Address
	L1BlockNumber     uint64
	EvidenceSubmitted bool
}, error) {
	return _RAT.Contract.AttentionTests(&_RAT.CallOpts, arg0)
}

// Challengers is a free data retrieval call binding the contract method 0xcfea71c0.
//
// Solidity: function challengers(address ) view returns(uint256 stakingAmount, uint256 totalSlashedAmount, uint32 validatorIndex, bool isValid)
func (_RAT *RATCaller) Challengers(opts *bind.CallOpts, arg0 common.Address) (struct {
	StakingAmount      *big.Int
	TotalSlashedAmount *big.Int
	ValidatorIndex     uint32
	IsValid            bool
}, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "challengers", arg0)

	outstruct := new(struct {
		StakingAmount      *big.Int
		TotalSlashedAmount *big.Int
		ValidatorIndex     uint32
		IsValid            bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakingAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalSlashedAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ValidatorIndex = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.IsValid = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Challengers is a free data retrieval call binding the contract method 0xcfea71c0.
//
// Solidity: function challengers(address ) view returns(uint256 stakingAmount, uint256 totalSlashedAmount, uint32 validatorIndex, bool isValid)
func (_RAT *RATSession) Challengers(arg0 common.Address) (struct {
	StakingAmount      *big.Int
	TotalSlashedAmount *big.Int
	ValidatorIndex     uint32
	IsValid            bool
}, error) {
	return _RAT.Contract.Challengers(&_RAT.CallOpts, arg0)
}

// Challengers is a free data retrieval call binding the contract method 0xcfea71c0.
//
// Solidity: function challengers(address ) view returns(uint256 stakingAmount, uint256 totalSlashedAmount, uint32 validatorIndex, bool isValid)
func (_RAT *RATCallerSession) Challengers(arg0 common.Address) (struct {
	StakingAmount      *big.Int
	TotalSlashedAmount *big.Int
	ValidatorIndex     uint32
	IsValid            bool
}, error) {
	return _RAT.Contract.Challengers(&_RAT.CallOpts, arg0)
}

// DisputeGameFactory is a free data retrieval call binding the contract method 0xf2b4e617.
//
// Solidity: function disputeGameFactory() view returns(address)
func (_RAT *RATCaller) DisputeGameFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "disputeGameFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DisputeGameFactory is a free data retrieval call binding the contract method 0xf2b4e617.
//
// Solidity: function disputeGameFactory() view returns(address)
func (_RAT *RATSession) DisputeGameFactory() (common.Address, error) {
	return _RAT.Contract.DisputeGameFactory(&_RAT.CallOpts)
}

// DisputeGameFactory is a free data retrieval call binding the contract method 0xf2b4e617.
//
// Solidity: function disputeGameFactory() view returns(address)
func (_RAT *RATCallerSession) DisputeGameFactory() (common.Address, error) {
	return _RAT.Contract.DisputeGameFactory(&_RAT.CallOpts)
}

// EvidenceSubmissionPeriod is a free data retrieval call binding the contract method 0xacccb08f.
//
// Solidity: function evidenceSubmissionPeriod() view returns(uint256)
func (_RAT *RATCaller) EvidenceSubmissionPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "evidenceSubmissionPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EvidenceSubmissionPeriod is a free data retrieval call binding the contract method 0xacccb08f.
//
// Solidity: function evidenceSubmissionPeriod() view returns(uint256)
func (_RAT *RATSession) EvidenceSubmissionPeriod() (*big.Int, error) {
	return _RAT.Contract.EvidenceSubmissionPeriod(&_RAT.CallOpts)
}

// EvidenceSubmissionPeriod is a free data retrieval call binding the contract method 0xacccb08f.
//
// Solidity: function evidenceSubmissionPeriod() view returns(uint256)
func (_RAT *RATCallerSession) EvidenceSubmissionPeriod() (*big.Int, error) {
	return _RAT.Contract.EvidenceSubmissionPeriod(&_RAT.CallOpts)
}

// GetChallengerInfo is a free data retrieval call binding the contract method 0x18350f22.
//
// Solidity: function getChallengerInfo(address _challenger) view returns((uint256,uint256,uint32,bool))
func (_RAT *RATCaller) GetChallengerInfo(opts *bind.CallOpts, _challenger common.Address) (RATChallengerInfo, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "getChallengerInfo", _challenger)

	if err != nil {
		return *new(RATChallengerInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(RATChallengerInfo)).(*RATChallengerInfo)

	return out0, err

}

// GetChallengerInfo is a free data retrieval call binding the contract method 0x18350f22.
//
// Solidity: function getChallengerInfo(address _challenger) view returns((uint256,uint256,uint32,bool))
func (_RAT *RATSession) GetChallengerInfo(_challenger common.Address) (RATChallengerInfo, error) {
	return _RAT.Contract.GetChallengerInfo(&_RAT.CallOpts, _challenger)
}

// GetChallengerInfo is a free data retrieval call binding the contract method 0x18350f22.
//
// Solidity: function getChallengerInfo(address _challenger) view returns((uint256,uint256,uint32,bool))
func (_RAT *RATCallerSession) GetChallengerInfo(_challenger common.Address) (RATChallengerInfo, error) {
	return _RAT.Contract.GetChallengerInfo(&_RAT.CallOpts, _challenger)
}

// GetValidChallengerCount is a free data retrieval call binding the contract method 0xcc57f42b.
//
// Solidity: function getValidChallengerCount() view returns(uint256)
func (_RAT *RATCaller) GetValidChallengerCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "getValidChallengerCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValidChallengerCount is a free data retrieval call binding the contract method 0xcc57f42b.
//
// Solidity: function getValidChallengerCount() view returns(uint256)
func (_RAT *RATSession) GetValidChallengerCount() (*big.Int, error) {
	return _RAT.Contract.GetValidChallengerCount(&_RAT.CallOpts)
}

// GetValidChallengerCount is a free data retrieval call binding the contract method 0xcc57f42b.
//
// Solidity: function getValidChallengerCount() view returns(uint256)
func (_RAT *RATCallerSession) GetValidChallengerCount() (*big.Int, error) {
	return _RAT.Contract.GetValidChallengerCount(&_RAT.CallOpts)
}

// InitVersion is a free data retrieval call binding the contract method 0x38d38c97.
//
// Solidity: function initVersion() view returns(uint8)
func (_RAT *RATCaller) InitVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "initVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// InitVersion is a free data retrieval call binding the contract method 0x38d38c97.
//
// Solidity: function initVersion() view returns(uint8)
func (_RAT *RATSession) InitVersion() (uint8, error) {
	return _RAT.Contract.InitVersion(&_RAT.CallOpts)
}

// InitVersion is a free data retrieval call binding the contract method 0x38d38c97.
//
// Solidity: function initVersion() view returns(uint8)
func (_RAT *RATCallerSession) InitVersion() (uint8, error) {
	return _RAT.Contract.InitVersion(&_RAT.CallOpts)
}

// MinimumStakingBalance is a free data retrieval call binding the contract method 0xe3478659.
//
// Solidity: function minimumStakingBalance() view returns(uint256)
func (_RAT *RATCaller) MinimumStakingBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "minimumStakingBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumStakingBalance is a free data retrieval call binding the contract method 0xe3478659.
//
// Solidity: function minimumStakingBalance() view returns(uint256)
func (_RAT *RATSession) MinimumStakingBalance() (*big.Int, error) {
	return _RAT.Contract.MinimumStakingBalance(&_RAT.CallOpts)
}

// MinimumStakingBalance is a free data retrieval call binding the contract method 0xe3478659.
//
// Solidity: function minimumStakingBalance() view returns(uint256)
func (_RAT *RATCallerSession) MinimumStakingBalance() (*big.Int, error) {
	return _RAT.Contract.MinimumStakingBalance(&_RAT.CallOpts)
}

// PerTestBondAmount is a free data retrieval call binding the contract method 0xcf60203e.
//
// Solidity: function perTestBondAmount() view returns(uint256)
func (_RAT *RATCaller) PerTestBondAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "perTestBondAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PerTestBondAmount is a free data retrieval call binding the contract method 0xcf60203e.
//
// Solidity: function perTestBondAmount() view returns(uint256)
func (_RAT *RATSession) PerTestBondAmount() (*big.Int, error) {
	return _RAT.Contract.PerTestBondAmount(&_RAT.CallOpts)
}

// PerTestBondAmount is a free data retrieval call binding the contract method 0xcf60203e.
//
// Solidity: function perTestBondAmount() view returns(uint256)
func (_RAT *RATCallerSession) PerTestBondAmount() (*big.Int, error) {
	return _RAT.Contract.PerTestBondAmount(&_RAT.CallOpts)
}

// ProxyAdmin is a free data retrieval call binding the contract method 0x3e47158c.
//
// Solidity: function proxyAdmin() view returns(address)
func (_RAT *RATCaller) ProxyAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "proxyAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProxyAdmin is a free data retrieval call binding the contract method 0x3e47158c.
//
// Solidity: function proxyAdmin() view returns(address)
func (_RAT *RATSession) ProxyAdmin() (common.Address, error) {
	return _RAT.Contract.ProxyAdmin(&_RAT.CallOpts)
}

// ProxyAdmin is a free data retrieval call binding the contract method 0x3e47158c.
//
// Solidity: function proxyAdmin() view returns(address)
func (_RAT *RATCallerSession) ProxyAdmin() (common.Address, error) {
	return _RAT.Contract.ProxyAdmin(&_RAT.CallOpts)
}

// ProxyAdminOwner is a free data retrieval call binding the contract method 0xdad544e0.
//
// Solidity: function proxyAdminOwner() view returns(address)
func (_RAT *RATCaller) ProxyAdminOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "proxyAdminOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProxyAdminOwner is a free data retrieval call binding the contract method 0xdad544e0.
//
// Solidity: function proxyAdminOwner() view returns(address)
func (_RAT *RATSession) ProxyAdminOwner() (common.Address, error) {
	return _RAT.Contract.ProxyAdminOwner(&_RAT.CallOpts)
}

// ProxyAdminOwner is a free data retrieval call binding the contract method 0xdad544e0.
//
// Solidity: function proxyAdminOwner() view returns(address)
func (_RAT *RATCallerSession) ProxyAdminOwner() (common.Address, error) {
	return _RAT.Contract.ProxyAdminOwner(&_RAT.CallOpts)
}

// RatManager is a free data retrieval call binding the contract method 0x1ecbb09c.
//
// Solidity: function ratManager() view returns(address)
func (_RAT *RATCaller) RatManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "ratManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RatManager is a free data retrieval call binding the contract method 0x1ecbb09c.
//
// Solidity: function ratManager() view returns(address)
func (_RAT *RATSession) RatManager() (common.Address, error) {
	return _RAT.Contract.RatManager(&_RAT.CallOpts)
}

// RatManager is a free data retrieval call binding the contract method 0x1ecbb09c.
//
// Solidity: function ratManager() view returns(address)
func (_RAT *RATCallerSession) RatManager() (common.Address, error) {
	return _RAT.Contract.RatManager(&_RAT.CallOpts)
}

// RatTriggerProbability is a free data retrieval call binding the contract method 0x51567bc2.
//
// Solidity: function ratTriggerProbability() view returns(uint256)
func (_RAT *RATCaller) RatTriggerProbability(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "ratTriggerProbability")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RatTriggerProbability is a free data retrieval call binding the contract method 0x51567bc2.
//
// Solidity: function ratTriggerProbability() view returns(uint256)
func (_RAT *RATSession) RatTriggerProbability() (*big.Int, error) {
	return _RAT.Contract.RatTriggerProbability(&_RAT.CallOpts)
}

// RatTriggerProbability is a free data retrieval call binding the contract method 0x51567bc2.
//
// Solidity: function ratTriggerProbability() view returns(uint256)
func (_RAT *RATCallerSession) RatTriggerProbability() (*big.Int, error) {
	return _RAT.Contract.RatTriggerProbability(&_RAT.CallOpts)
}

// ValidChallengers is a free data retrieval call binding the contract method 0x74343de8.
//
// Solidity: function validChallengers(uint256 ) view returns(address)
func (_RAT *RATCaller) ValidChallengers(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "validChallengers", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ValidChallengers is a free data retrieval call binding the contract method 0x74343de8.
//
// Solidity: function validChallengers(uint256 ) view returns(address)
func (_RAT *RATSession) ValidChallengers(arg0 *big.Int) (common.Address, error) {
	return _RAT.Contract.ValidChallengers(&_RAT.CallOpts, arg0)
}

// ValidChallengers is a free data retrieval call binding the contract method 0x74343de8.
//
// Solidity: function validChallengers(uint256 ) view returns(address)
func (_RAT *RATCallerSession) ValidChallengers(arg0 *big.Int) (common.Address, error) {
	return _RAT.Contract.ValidChallengers(&_RAT.CallOpts, arg0)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_RAT *RATCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RAT.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_RAT *RATSession) Version() (string, error) {
	return _RAT.Contract.Version(&_RAT.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_RAT *RATCallerSession) Version() (string, error) {
	return _RAT.Contract.Version(&_RAT.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x5df5f96f.
//
// Solidity: function initialize(address _disputeGameFactory, uint256 _perTestBondAmount, uint256 _evidenceSubmissionPeriod, uint256 _minimumStakingBalance, uint256 _ratTriggerProbability, address _manager) payable returns()
func (_RAT *RATTransactor) Initialize(opts *bind.TransactOpts, _disputeGameFactory common.Address, _perTestBondAmount *big.Int, _evidenceSubmissionPeriod *big.Int, _minimumStakingBalance *big.Int, _ratTriggerProbability *big.Int, _manager common.Address) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "initialize", _disputeGameFactory, _perTestBondAmount, _evidenceSubmissionPeriod, _minimumStakingBalance, _ratTriggerProbability, _manager)
}

// Initialize is a paid mutator transaction binding the contract method 0x5df5f96f.
//
// Solidity: function initialize(address _disputeGameFactory, uint256 _perTestBondAmount, uint256 _evidenceSubmissionPeriod, uint256 _minimumStakingBalance, uint256 _ratTriggerProbability, address _manager) payable returns()
func (_RAT *RATSession) Initialize(_disputeGameFactory common.Address, _perTestBondAmount *big.Int, _evidenceSubmissionPeriod *big.Int, _minimumStakingBalance *big.Int, _ratTriggerProbability *big.Int, _manager common.Address) (*types.Transaction, error) {
	return _RAT.Contract.Initialize(&_RAT.TransactOpts, _disputeGameFactory, _perTestBondAmount, _evidenceSubmissionPeriod, _minimumStakingBalance, _ratTriggerProbability, _manager)
}

// Initialize is a paid mutator transaction binding the contract method 0x5df5f96f.
//
// Solidity: function initialize(address _disputeGameFactory, uint256 _perTestBondAmount, uint256 _evidenceSubmissionPeriod, uint256 _minimumStakingBalance, uint256 _ratTriggerProbability, address _manager) payable returns()
func (_RAT *RATTransactorSession) Initialize(_disputeGameFactory common.Address, _perTestBondAmount *big.Int, _evidenceSubmissionPeriod *big.Int, _minimumStakingBalance *big.Int, _ratTriggerProbability *big.Int, _manager common.Address) (*types.Transaction, error) {
	return _RAT.Contract.Initialize(&_RAT.TransactOpts, _disputeGameFactory, _perTestBondAmount, _evidenceSubmissionPeriod, _minimumStakingBalance, _ratTriggerProbability, _manager)
}

// ResolveClaim is a paid mutator transaction binding the contract method 0x94d645a8.
//
// Solidity: function resolveClaim(address _claimant) returns()
func (_RAT *RATTransactor) ResolveClaim(opts *bind.TransactOpts, _claimant common.Address) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "resolveClaim", _claimant)
}

// ResolveClaim is a paid mutator transaction binding the contract method 0x94d645a8.
//
// Solidity: function resolveClaim(address _claimant) returns()
func (_RAT *RATSession) ResolveClaim(_claimant common.Address) (*types.Transaction, error) {
	return _RAT.Contract.ResolveClaim(&_RAT.TransactOpts, _claimant)
}

// ResolveClaim is a paid mutator transaction binding the contract method 0x94d645a8.
//
// Solidity: function resolveClaim(address _claimant) returns()
func (_RAT *RATTransactorSession) ResolveClaim(_claimant common.Address) (*types.Transaction, error) {
	return _RAT.Contract.ResolveClaim(&_RAT.TransactOpts, _claimant)
}

// SetEvidenceSubmissionPeriod is a paid mutator transaction binding the contract method 0x36c63d46.
//
// Solidity: function setEvidenceSubmissionPeriod(uint256 _period) returns()
func (_RAT *RATTransactor) SetEvidenceSubmissionPeriod(opts *bind.TransactOpts, _period *big.Int) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "setEvidenceSubmissionPeriod", _period)
}

// SetEvidenceSubmissionPeriod is a paid mutator transaction binding the contract method 0x36c63d46.
//
// Solidity: function setEvidenceSubmissionPeriod(uint256 _period) returns()
func (_RAT *RATSession) SetEvidenceSubmissionPeriod(_period *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetEvidenceSubmissionPeriod(&_RAT.TransactOpts, _period)
}

// SetEvidenceSubmissionPeriod is a paid mutator transaction binding the contract method 0x36c63d46.
//
// Solidity: function setEvidenceSubmissionPeriod(uint256 _period) returns()
func (_RAT *RATTransactorSession) SetEvidenceSubmissionPeriod(_period *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetEvidenceSubmissionPeriod(&_RAT.TransactOpts, _period)
}

// SetMinimumStakingBalance is a paid mutator transaction binding the contract method 0xfd5f97e1.
//
// Solidity: function setMinimumStakingBalance(uint256 _balance) returns()
func (_RAT *RATTransactor) SetMinimumStakingBalance(opts *bind.TransactOpts, _balance *big.Int) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "setMinimumStakingBalance", _balance)
}

// SetMinimumStakingBalance is a paid mutator transaction binding the contract method 0xfd5f97e1.
//
// Solidity: function setMinimumStakingBalance(uint256 _balance) returns()
func (_RAT *RATSession) SetMinimumStakingBalance(_balance *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetMinimumStakingBalance(&_RAT.TransactOpts, _balance)
}

// SetMinimumStakingBalance is a paid mutator transaction binding the contract method 0xfd5f97e1.
//
// Solidity: function setMinimumStakingBalance(uint256 _balance) returns()
func (_RAT *RATTransactorSession) SetMinimumStakingBalance(_balance *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetMinimumStakingBalance(&_RAT.TransactOpts, _balance)
}

// SetPerTestBondAmount is a paid mutator transaction binding the contract method 0xcc447906.
//
// Solidity: function setPerTestBondAmount(uint256 _amount) returns()
func (_RAT *RATTransactor) SetPerTestBondAmount(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "setPerTestBondAmount", _amount)
}

// SetPerTestBondAmount is a paid mutator transaction binding the contract method 0xcc447906.
//
// Solidity: function setPerTestBondAmount(uint256 _amount) returns()
func (_RAT *RATSession) SetPerTestBondAmount(_amount *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetPerTestBondAmount(&_RAT.TransactOpts, _amount)
}

// SetPerTestBondAmount is a paid mutator transaction binding the contract method 0xcc447906.
//
// Solidity: function setPerTestBondAmount(uint256 _amount) returns()
func (_RAT *RATTransactorSession) SetPerTestBondAmount(_amount *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetPerTestBondAmount(&_RAT.TransactOpts, _amount)
}

// SetRatTriggerProbability is a paid mutator transaction binding the contract method 0xd2e5bc72.
//
// Solidity: function setRatTriggerProbability(uint256 _probability) returns()
func (_RAT *RATTransactor) SetRatTriggerProbability(opts *bind.TransactOpts, _probability *big.Int) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "setRatTriggerProbability", _probability)
}

// SetRatTriggerProbability is a paid mutator transaction binding the contract method 0xd2e5bc72.
//
// Solidity: function setRatTriggerProbability(uint256 _probability) returns()
func (_RAT *RATSession) SetRatTriggerProbability(_probability *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetRatTriggerProbability(&_RAT.TransactOpts, _probability)
}

// SetRatTriggerProbability is a paid mutator transaction binding the contract method 0xd2e5bc72.
//
// Solidity: function setRatTriggerProbability(uint256 _probability) returns()
func (_RAT *RATTransactorSession) SetRatTriggerProbability(_probability *big.Int) (*types.Transaction, error) {
	return _RAT.Contract.SetRatTriggerProbability(&_RAT.TransactOpts, _probability)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_RAT *RATTransactor) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "stake")
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_RAT *RATSession) Stake() (*types.Transaction, error) {
	return _RAT.Contract.Stake(&_RAT.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_RAT *RATTransactorSession) Stake() (*types.Transaction, error) {
	return _RAT.Contract.Stake(&_RAT.TransactOpts)
}

// SubmitCorrectEvidence is a paid mutator transaction binding the contract method 0xaacff5f9.
//
// Solidity: function submitCorrectEvidence(address _gameAddress, bytes32 _proofLV, bytes32 _proofRV) returns()
func (_RAT *RATTransactor) SubmitCorrectEvidence(opts *bind.TransactOpts, _gameAddress common.Address, _proofLV [32]byte, _proofRV [32]byte) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "submitCorrectEvidence", _gameAddress, _proofLV, _proofRV)
}

// SubmitCorrectEvidence is a paid mutator transaction binding the contract method 0xaacff5f9.
//
// Solidity: function submitCorrectEvidence(address _gameAddress, bytes32 _proofLV, bytes32 _proofRV) returns()
func (_RAT *RATSession) SubmitCorrectEvidence(_gameAddress common.Address, _proofLV [32]byte, _proofRV [32]byte) (*types.Transaction, error) {
	return _RAT.Contract.SubmitCorrectEvidence(&_RAT.TransactOpts, _gameAddress, _proofLV, _proofRV)
}

// SubmitCorrectEvidence is a paid mutator transaction binding the contract method 0xaacff5f9.
//
// Solidity: function submitCorrectEvidence(address _gameAddress, bytes32 _proofLV, bytes32 _proofRV) returns()
func (_RAT *RATTransactorSession) SubmitCorrectEvidence(_gameAddress common.Address, _proofLV [32]byte, _proofRV [32]byte) (*types.Transaction, error) {
	return _RAT.Contract.SubmitCorrectEvidence(&_RAT.TransactOpts, _gameAddress, _proofLV, _proofRV)
}

// TriggerAttentionTest is a paid mutator transaction binding the contract method 0x96e9b641.
//
// Solidity: function triggerAttentionTest(address _gameAddress, bytes32 _stateRoot, bytes32 _blockHash) returns()
func (_RAT *RATTransactor) TriggerAttentionTest(opts *bind.TransactOpts, _gameAddress common.Address, _stateRoot [32]byte, _blockHash [32]byte) (*types.Transaction, error) {
	return _RAT.contract.Transact(opts, "triggerAttentionTest", _gameAddress, _stateRoot, _blockHash)
}

// TriggerAttentionTest is a paid mutator transaction binding the contract method 0x96e9b641.
//
// Solidity: function triggerAttentionTest(address _gameAddress, bytes32 _stateRoot, bytes32 _blockHash) returns()
func (_RAT *RATSession) TriggerAttentionTest(_gameAddress common.Address, _stateRoot [32]byte, _blockHash [32]byte) (*types.Transaction, error) {
	return _RAT.Contract.TriggerAttentionTest(&_RAT.TransactOpts, _gameAddress, _stateRoot, _blockHash)
}

// TriggerAttentionTest is a paid mutator transaction binding the contract method 0x96e9b641.
//
// Solidity: function triggerAttentionTest(address _gameAddress, bytes32 _stateRoot, bytes32 _blockHash) returns()
func (_RAT *RATTransactorSession) TriggerAttentionTest(_gameAddress common.Address, _stateRoot [32]byte, _blockHash [32]byte) (*types.Transaction, error) {
	return _RAT.Contract.TriggerAttentionTest(&_RAT.TransactOpts, _gameAddress, _stateRoot, _blockHash)
}

// RATAttentionTriggeredIterator is returned from FilterAttentionTriggered and is used to iterate over the raw logs and unpacked data for AttentionTriggered events raised by the RAT contract.
type RATAttentionTriggeredIterator struct {
	Event *RATAttentionTriggered // Event containing the contract specifics and raw log

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
func (it *RATAttentionTriggeredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RATAttentionTriggered)
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
		it.Event = new(RATAttentionTriggered)
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
func (it *RATAttentionTriggeredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RATAttentionTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RATAttentionTriggered represents a AttentionTriggered event raised by the RAT contract.
type RATAttentionTriggered struct {
	GameAddress common.Address
	Challenger  common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAttentionTriggered is a free log retrieval operation binding the contract event 0x8f5f18c2fab75f3bb8637c4702468685d73e9414ddfa43f842ffb1a63e6dd57a.
//
// Solidity: event AttentionTriggered(address indexed gameAddress, address indexed challenger)
func (_RAT *RATFilterer) FilterAttentionTriggered(opts *bind.FilterOpts, gameAddress []common.Address, challenger []common.Address) (*RATAttentionTriggeredIterator, error) {

	var gameAddressRule []interface{}
	for _, gameAddressItem := range gameAddress {
		gameAddressRule = append(gameAddressRule, gameAddressItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.FilterLogs(opts, "AttentionTriggered", gameAddressRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return &RATAttentionTriggeredIterator{contract: _RAT.contract, event: "AttentionTriggered", logs: logs, sub: sub}, nil
}

// WatchAttentionTriggered is a free log subscription operation binding the contract event 0x8f5f18c2fab75f3bb8637c4702468685d73e9414ddfa43f842ffb1a63e6dd57a.
//
// Solidity: event AttentionTriggered(address indexed gameAddress, address indexed challenger)
func (_RAT *RATFilterer) WatchAttentionTriggered(opts *bind.WatchOpts, sink chan<- *RATAttentionTriggered, gameAddress []common.Address, challenger []common.Address) (event.Subscription, error) {

	var gameAddressRule []interface{}
	for _, gameAddressItem := range gameAddress {
		gameAddressRule = append(gameAddressRule, gameAddressItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.WatchLogs(opts, "AttentionTriggered", gameAddressRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RATAttentionTriggered)
				if err := _RAT.contract.UnpackLog(event, "AttentionTriggered", log); err != nil {
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

// ParseAttentionTriggered is a log parse operation binding the contract event 0x8f5f18c2fab75f3bb8637c4702468685d73e9414ddfa43f842ffb1a63e6dd57a.
//
// Solidity: event AttentionTriggered(address indexed gameAddress, address indexed challenger)
func (_RAT *RATFilterer) ParseAttentionTriggered(log types.Log) (*RATAttentionTriggered, error) {
	event := new(RATAttentionTriggered)
	if err := _RAT.contract.UnpackLog(event, "AttentionTriggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RATBondRefundedIterator is returned from FilterBondRefunded and is used to iterate over the raw logs and unpacked data for BondRefunded events raised by the RAT contract.
type RATBondRefundedIterator struct {
	Event *RATBondRefunded // Event containing the contract specifics and raw log

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
func (it *RATBondRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RATBondRefunded)
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
		it.Event = new(RATBondRefunded)
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
func (it *RATBondRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RATBondRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RATBondRefunded represents a BondRefunded event raised by the RAT contract.
type RATBondRefunded struct {
	GameAddress    common.Address
	Challenger     common.Address
	RefundedAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBondRefunded is a free log retrieval operation binding the contract event 0x4fac74d6b4ed45cb641c8e730951b01c4dbcae0565cf330e909230492bd07d78.
//
// Solidity: event BondRefunded(address indexed gameAddress, address indexed challenger, uint256 refundedAmount)
func (_RAT *RATFilterer) FilterBondRefunded(opts *bind.FilterOpts, gameAddress []common.Address, challenger []common.Address) (*RATBondRefundedIterator, error) {

	var gameAddressRule []interface{}
	for _, gameAddressItem := range gameAddress {
		gameAddressRule = append(gameAddressRule, gameAddressItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.FilterLogs(opts, "BondRefunded", gameAddressRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return &RATBondRefundedIterator{contract: _RAT.contract, event: "BondRefunded", logs: logs, sub: sub}, nil
}

// WatchBondRefunded is a free log subscription operation binding the contract event 0x4fac74d6b4ed45cb641c8e730951b01c4dbcae0565cf330e909230492bd07d78.
//
// Solidity: event BondRefunded(address indexed gameAddress, address indexed challenger, uint256 refundedAmount)
func (_RAT *RATFilterer) WatchBondRefunded(opts *bind.WatchOpts, sink chan<- *RATBondRefunded, gameAddress []common.Address, challenger []common.Address) (event.Subscription, error) {

	var gameAddressRule []interface{}
	for _, gameAddressItem := range gameAddress {
		gameAddressRule = append(gameAddressRule, gameAddressItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.WatchLogs(opts, "BondRefunded", gameAddressRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RATBondRefunded)
				if err := _RAT.contract.UnpackLog(event, "BondRefunded", log); err != nil {
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

// ParseBondRefunded is a log parse operation binding the contract event 0x4fac74d6b4ed45cb641c8e730951b01c4dbcae0565cf330e909230492bd07d78.
//
// Solidity: event BondRefunded(address indexed gameAddress, address indexed challenger, uint256 refundedAmount)
func (_RAT *RATFilterer) ParseBondRefunded(log types.Log) (*RATBondRefunded, error) {
	event := new(RATBondRefunded)
	if err := _RAT.contract.UnpackLog(event, "BondRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RATChallengerStakedIterator is returned from FilterChallengerStaked and is used to iterate over the raw logs and unpacked data for ChallengerStaked events raised by the RAT contract.
type RATChallengerStakedIterator struct {
	Event *RATChallengerStaked // Event containing the contract specifics and raw log

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
func (it *RATChallengerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RATChallengerStaked)
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
		it.Event = new(RATChallengerStaked)
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
func (it *RATChallengerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RATChallengerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RATChallengerStaked represents a ChallengerStaked event raised by the RAT contract.
type RATChallengerStaked struct {
	Challenger common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterChallengerStaked is a free log retrieval operation binding the contract event 0x6f50cc7a01f21c217b2ae66736754596de2004562a26fcc850e70031e5985e8d.
//
// Solidity: event ChallengerStaked(address indexed challenger, uint256 amount)
func (_RAT *RATFilterer) FilterChallengerStaked(opts *bind.FilterOpts, challenger []common.Address) (*RATChallengerStakedIterator, error) {

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.FilterLogs(opts, "ChallengerStaked", challengerRule)
	if err != nil {
		return nil, err
	}
	return &RATChallengerStakedIterator{contract: _RAT.contract, event: "ChallengerStaked", logs: logs, sub: sub}, nil
}

// WatchChallengerStaked is a free log subscription operation binding the contract event 0x6f50cc7a01f21c217b2ae66736754596de2004562a26fcc850e70031e5985e8d.
//
// Solidity: event ChallengerStaked(address indexed challenger, uint256 amount)
func (_RAT *RATFilterer) WatchChallengerStaked(opts *bind.WatchOpts, sink chan<- *RATChallengerStaked, challenger []common.Address) (event.Subscription, error) {

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.WatchLogs(opts, "ChallengerStaked", challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RATChallengerStaked)
				if err := _RAT.contract.UnpackLog(event, "ChallengerStaked", log); err != nil {
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

// ParseChallengerStaked is a log parse operation binding the contract event 0x6f50cc7a01f21c217b2ae66736754596de2004562a26fcc850e70031e5985e8d.
//
// Solidity: event ChallengerStaked(address indexed challenger, uint256 amount)
func (_RAT *RATFilterer) ParseChallengerStaked(log types.Log) (*RATChallengerStaked, error) {
	event := new(RATChallengerStaked)
	if err := _RAT.contract.UnpackLog(event, "ChallengerStaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RATCorrectEvidenceSubmittedIterator is returned from FilterCorrectEvidenceSubmitted and is used to iterate over the raw logs and unpacked data for CorrectEvidenceSubmitted events raised by the RAT contract.
type RATCorrectEvidenceSubmittedIterator struct {
	Event *RATCorrectEvidenceSubmitted // Event containing the contract specifics and raw log

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
func (it *RATCorrectEvidenceSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RATCorrectEvidenceSubmitted)
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
		it.Event = new(RATCorrectEvidenceSubmitted)
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
func (it *RATCorrectEvidenceSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RATCorrectEvidenceSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RATCorrectEvidenceSubmitted represents a CorrectEvidenceSubmitted event raised by the RAT contract.
type RATCorrectEvidenceSubmitted struct {
	GameAddress    common.Address
	Challenger     common.Address
	RestoredAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCorrectEvidenceSubmitted is a free log retrieval operation binding the contract event 0x7954ca3465b792ac0a46ec18f1ed49a20399c01cc510fab6eab44a7f6679f877.
//
// Solidity: event CorrectEvidenceSubmitted(address indexed gameAddress, address indexed challenger, uint256 restoredAmount)
func (_RAT *RATFilterer) FilterCorrectEvidenceSubmitted(opts *bind.FilterOpts, gameAddress []common.Address, challenger []common.Address) (*RATCorrectEvidenceSubmittedIterator, error) {

	var gameAddressRule []interface{}
	for _, gameAddressItem := range gameAddress {
		gameAddressRule = append(gameAddressRule, gameAddressItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.FilterLogs(opts, "CorrectEvidenceSubmitted", gameAddressRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return &RATCorrectEvidenceSubmittedIterator{contract: _RAT.contract, event: "CorrectEvidenceSubmitted", logs: logs, sub: sub}, nil
}

// WatchCorrectEvidenceSubmitted is a free log subscription operation binding the contract event 0x7954ca3465b792ac0a46ec18f1ed49a20399c01cc510fab6eab44a7f6679f877.
//
// Solidity: event CorrectEvidenceSubmitted(address indexed gameAddress, address indexed challenger, uint256 restoredAmount)
func (_RAT *RATFilterer) WatchCorrectEvidenceSubmitted(opts *bind.WatchOpts, sink chan<- *RATCorrectEvidenceSubmitted, gameAddress []common.Address, challenger []common.Address) (event.Subscription, error) {

	var gameAddressRule []interface{}
	for _, gameAddressItem := range gameAddress {
		gameAddressRule = append(gameAddressRule, gameAddressItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _RAT.contract.WatchLogs(opts, "CorrectEvidenceSubmitted", gameAddressRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RATCorrectEvidenceSubmitted)
				if err := _RAT.contract.UnpackLog(event, "CorrectEvidenceSubmitted", log); err != nil {
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

// ParseCorrectEvidenceSubmitted is a log parse operation binding the contract event 0x7954ca3465b792ac0a46ec18f1ed49a20399c01cc510fab6eab44a7f6679f877.
//
// Solidity: event CorrectEvidenceSubmitted(address indexed gameAddress, address indexed challenger, uint256 restoredAmount)
func (_RAT *RATFilterer) ParseCorrectEvidenceSubmitted(log types.Log) (*RATCorrectEvidenceSubmitted, error) {
	event := new(RATCorrectEvidenceSubmitted)
	if err := _RAT.contract.UnpackLog(event, "CorrectEvidenceSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RATInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RAT contract.
type RATInitializedIterator struct {
	Event *RATInitialized // Event containing the contract specifics and raw log

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
func (it *RATInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RATInitialized)
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
		it.Event = new(RATInitialized)
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
func (it *RATInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RATInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RATInitialized represents a Initialized event raised by the RAT contract.
type RATInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RAT *RATFilterer) FilterInitialized(opts *bind.FilterOpts) (*RATInitializedIterator, error) {

	logs, sub, err := _RAT.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RATInitializedIterator{contract: _RAT.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RAT *RATFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RATInitialized) (event.Subscription, error) {

	logs, sub, err := _RAT.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RATInitialized)
				if err := _RAT.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RAT *RATFilterer) ParseInitialized(log types.Log) (*RATInitialized, error) {
	event := new(RATInitialized)
	if err := _RAT.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RATBin contains the bytecode for the RAT contract
var RATBin = "0x60a06040523480156200001157600080fd5b50600260805260018055620000256200002b565b620000ed565b600054610100900460ff1615620000985760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff9081161015620000eb576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b60805161252762000110600039600081816103040152610bd001526125276000f3fe60806040526004361061018b5760003560e01c8063aacff5f9116100d6578063d2e5bc721161007f578063f2b4e61711610059578063f2b4e61714610581578063fba1d220146105ae578063fd5f97e11461069757600080fd5b8063d2e5bc7214610536578063dad544e014610556578063e34786591461056b57600080fd5b8063cc57f42b116100b0578063cc57f42b1461048e578063cf60203e146104a3578063cfea71c0146104b957600080fd5b8063aacff5f914610438578063acccb08f14610458578063cc4479061461046e57600080fd5b806351567bc21161013857806374343de81161011257806374343de8146103d857806394d645a8146103f857806396e9b6411461041857600080fd5b806351567bc21461034b57806354fd4d501461036f5780635df5f96f146103c557600080fd5b806338d38c971161016957806338d38c97146102f05780633a4b66f11461032e5780633e47158c1461033657600080fd5b806318350f22146101905780631ecbb09c1461027c57806336c63d46146102ce575b600080fd5b34801561019c57600080fd5b506102386101ab366004612286565b6040805160808101825260008082526020820181905291810182905260608101919091525073ffffffffffffffffffffffffffffffffffffffff166000908152600860209081526040918290208251608081018452815481526001820154928101929092526002015463ffffffff81169282019290925264010000000090910460ff161515606082015290565b6040516102739190815181526020808301519082015260408083015163ffffffff169082015260609182015115159181019190915260800190565b60405180910390f35b34801561028857600080fd5b506007546102a99073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610273565b3480156102da57600080fd5b506102ee6102e93660046122aa565b6106b7565b005b3480156102fc57600080fd5b5060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610273565b6102ee61079f565b34801561034257600080fd5b506102a96109c3565b34801561035757600080fd5b5061036160065481565b604051908152602001610273565b34801561037b57600080fd5b506103b86040518060400160405280600c81526020017f312e302e302d626574612e31000000000000000000000000000000000000000081525081565b60405161027391906122c3565b6102ee6103d3366004612336565b610bce565b3480156103e457600080fd5b506102a96103f33660046122aa565b610f6a565b34801561040457600080fd5b506102ee610413366004612286565b610fa1565b34801561042457600080fd5b506102ee610433366004612394565b611293565b34801561044457600080fd5b506102ee610453366004612394565b6117a4565b34801561046457600080fd5b5061036160045481565b34801561047a57600080fd5b506102ee6104893660046122aa565b611c1e565b34801561049a57600080fd5b50600a54610361565b3480156104af57600080fd5b5061036160035481565b3480156104c557600080fd5b506105086104d4366004612286565b60086020526000908152604090208054600182015460029092015490919063ffffffff811690640100000000900460ff1684565b6040516102739493929190938452602084019290925263ffffffff1660408301521515606082015260800190565b34801561054257600080fd5b506102ee6105513660046122aa565b611dc3565b34801561056257600080fd5b506102a9611e86565b34801561057757600080fd5b5061036160055481565b34801561058d57600080fd5b506002546102a99073ffffffffffffffffffffffffffffffffffffffff1681565b3480156105ba57600080fd5b5061063c6105c9366004612286565b60096020526000908152604090208054600182015460029092015490916bffffffffffffffffffffffff8116916c0100000000000000000000000090910473ffffffffffffffffffffffffffffffffffffffff169067ffffffffffffffff81169068010000000000000000900460ff1685565b604080519586526bffffffffffffffffffffffff909416602086015273ffffffffffffffffffffffffffffffffffffffff9092169284019290925267ffffffffffffffff90911660608301521515608082015260a001610273565b3480156106a357600080fd5b506102ee6106b23660046122aa565b611f03565b6106bf611fed565b6000811161072e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f506572696f64206d75737420626520706f73697469766500000000000000000060448201526064015b60405180910390fd5b61c4e081111561079a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f506572696f6420746f6f206c6f6e6700000000000000000000000000000000006044820152606401610725565b600455565b60026001540361080b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610725565b600260015534610877576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4d757374207374616b6520706f73697469766520616d6f756e740000000000006044820152606401610725565b33600090815260086020526040812080549091349183919061089a9084906123f8565b90915550506002810154640100000000900460ff161580156108bf5750600354815410155b15610987576002810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff81166401000000009081178355600a805463ffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000909316929092171790915580546001810182556000919091527fc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a80180547fffffffffffffffffffffffff000000000000000000000000000000000000000016331790555b60405134815233907f6f50cc7a01f21c217b2ae66736754596de2004562a26fcc850e70031e5985e8d9060200160405180910390a25060018055565b6000806109ee7fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035490565b905073ffffffffffffffffffffffffffffffffffffffff811615610a1157919050565b6040518060400160405280601a81526020017f4f564d5f4c3143726f7373446f6d61696e4d657373656e676572000000000000815250516002610a549190612410565b604080513060208201526000918101919091527f4f564d5f4c3143726f7373446f6d61696e4d657373656e6765720000000000009190911790610aaf906060015b604051602081830303815290604052805190602001205490565b14610ae6576040517f54e433cd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051306020820152600191810191909152600090610b0890606001610a95565b905073ffffffffffffffffffffffffffffffffffffffff811615610b9c578073ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b71573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b95919061244d565b9250505090565b6040517f332144db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f0000000000000000000000000000000000000000000000000000000000000000600054610100900460ff16158015610c0e575060005460ff8083169116105b610c9a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610725565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660ff8316176101001790556bffffffffffffffffffffffff8610610d66576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f426f6e6420616d6f756e7420657863656564732075696e743936206d6178696d60448201527f756d0000000000000000000000000000000000000000000000000000000000006064820152608401610725565b83861115610df6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603160248201527f426f6e6420616d6f756e742063616e6e6f7420657863656564206d696e696d7560448201527f6d207374616b696e672062616c616e63650000000000000000000000000000006064820152608401610725565b620186a0831115610e63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f496e76616c69642070726f626162696c697479000000000000000000000000006044820152606401610725565b6002805473ffffffffffffffffffffffffffffffffffffffff8981167fffffffffffffffffffffffff00000000000000000000000000000000000000009283161790925560038890556004879055600586905560068590556007805492851692821692909217909155600a805460018101825560009182527fc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a801805490921690915580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16905560405160ff821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150505050505050565b600a8181548110610f7a57600080fd5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905081565b336000908152600960205260409020600101546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16801580159061101457508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b1561128f573360009081526009602052604090206002015468010000000000000000900460ff1661128f573360009081526009602090815260408083206002810180547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1668010000000000000000179055600181015473ffffffffffffffffffffffffffffffffffffffff871685526008909352908320805491936bffffffffffffffffffffffff909316929091839183916110d29084906123f8565b90915550506002810154600354825464010000000090920460ff16911015811580156110fb5750805b156111ef57600283810180546401000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff909116179055600a805473ffffffffffffffffffffffffffffffffffffffff8a16600081815260086020526040812090940180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff909316929092179091558154600181018355919092527fc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a80180547fffffffffffffffffffffffff000000000000000000000000000000000000000016909117905561123c565b8180156111fa575080155b1561123c576002830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff811690915561123c90889063ffffffff16612045565b60405184815273ffffffffffffffffffffffffffffffffffffffff87169033907f4fac74d6b4ed45cb641c8e730951b01c4dbcae0565cf330e909230492bd07d789060200160405180910390a350505050505b5050565b60025473ffffffffffffffffffffffffffffffffffffffff1633146112e4576040517f2a769f9e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6112ec612218565b1561179f57600a54600181111561179d576000816002146113655761131260018361246a565b834260405160200161132e929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c61ffff166113559190612481565b6113609060016123f8565b611368565b60015b90506000600a828154811061137f5761137f6124bc565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1680835260089091526040822080546003549294509092909182106113c9576003546113cb565b815b905060006113d9828461246a565b80855560018501805491925083916000906113f59084906123f8565b909155505067ffffffffffffffff43111561146c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f426c6f636b206e756d62657220746f6f206c61726765000000000000000000006044820152606401610725565b6040518060a001604052808a8152602001836bffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1681526020014367ffffffffffffffff16815260200160001515815250600960008c73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000820151816000015560208201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550604082015181600101600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160020160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060808201518160020160086101000a81548160ff0219169083151502179055509050506000600354821015905060008560020160049054906101000a900460ff1690508115158115151461173a5781156116fd57600286810180546401000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff909116179055600a805473ffffffffffffffffffffffffffffffffffffffff8a16600081815260086020526040812090940180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff909316929092179091558154600181018355919092527fc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a80180547fffffffffffffffffffffffff000000000000000000000000000000000000000016909117905561173a565b6002860180547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff811690915561173a90889063ffffffff16612045565b8673ffffffffffffffffffffffffffffffffffffffff168c73ffffffffffffffffffffffffffffffffffffffff167f8f5f18c2fab75f3bb8637c4702468685d73e9414ddfa43f842ffb1a63e6dd57a60405160405180910390a350505050505050505b505b505050565b73ffffffffffffffffffffffffffffffffffffffff8084166000908152600960205260409020600181015490916c01000000000000000000000000909104168061181a576040517facddc2de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81163314611869576040517f430a2a2b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600282015468010000000000000000900460ff16156118b4576040517fc355107000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60045460028301546000916118d29167ffffffffffffffff166123f8565b600284015490915067ffffffffffffffff1681101561194d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f446561646c696e65206f766572666c6f770000000000000000000000000000006044820152606401610725565b804310611986576040517f01136fa500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8254604080516020810188905290810186905260600160405160208183030381529060405280519060200120146119e9576040517fd611c31800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018301546002840180547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff166801000000000000000017905533600090815260086020526040812080546bffffffffffffffffffffffff9093169290918391839190611a579084906123f8565b90915550506002810154600354825464010000000090920460ff1691101581158015611a805750805b15611b5e57600283810180546401000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff909116179055600a805433600081815260086020526040812090940180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff909316929092179091558154600181018355919092527fc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a80180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169091179055611bab565b818015611b69575080155b15611bab576002830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff8116909155611bab90339063ffffffff16612045565b8573ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff167f7954ca3465b792ac0a46ec18f1ed49a20399c01cc510fab6eab44a7f6679f87786604051611c0a91815260200190565b60405180910390a350505050505050505050565b611c26611fed565b60008111611c90576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f426f6e6420616d6f756e74206d75737420626520706f736974697665000000006044820152606401610725565b6bffffffffffffffffffffffff811115611d2c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f426f6e6420616d6f756e7420657863656564732075696e743936206d6178696d60448201527f756d0000000000000000000000000000000000000000000000000000000000006064820152608401610725565b600554811115611dbe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603160248201527f426f6e6420616d6f756e742063616e6e6f7420657863656564206d696e696d7560448201527f6d207374616b696e672062616c616e63650000000000000000000000000000006064820152608401610725565b600355565b60075473ffffffffffffffffffffffffffffffffffffffff163314611e14576040517fc1de1ba600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620186a0811115611e81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f496e76616c69642070726f626162696c697479000000000000000000000000006044820152606401610725565b600655565b6000611e906109c3565b73ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611eda573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611efe919061244d565b905090565b611f0b611fed565b60008111611f75576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f42616c616e6365206d75737420626520706f73697469766500000000000000006044820152606401610725565b683635c9adc5dea00000811115611fe8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f42616c616e636520746f6f206c617267650000000000000000000000000000006044820152606401610725565b600555565b33611ff6611e86565b73ffffffffffffffffffffffffffffffffffffffff1614612043576040517f7f12c64b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000811180156120565750600a5481105b80156120ab57508173ffffffffffffffffffffffffffffffffffffffff16600a8281548110612087576120876124bc565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16145b1561128f57600a546000906120c29060019061246a565b90508082146121aa576000600a82815481106120e0576120e06124bc565b600091825260209091200154600a805473ffffffffffffffffffffffffffffffffffffffff909216925082918590811061211c5761211c6124bc565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff94851617905592909116815260089091526040902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff84161790555b600a8054806121bb576121bb6124eb565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055505050565b60065460009080820361222d57600091505090565b620186a0811061223f57600191505090565b80620186a061224f60014361246a565b61225a919040612481565b1091505090565b73ffffffffffffffffffffffffffffffffffffffff8116811461228357600080fd5b50565b60006020828403121561229857600080fd5b81356122a381612261565b9392505050565b6000602082840312156122bc57600080fd5b5035919050565b600060208083528351808285015260005b818110156122f0578581018301518582016040015282016122d4565b81811115612302576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60008060008060008060c0878903121561234f57600080fd5b863561235a81612261565b95506020870135945060408701359350606087013592506080870135915060a087013561238681612261565b809150509295509295509295565b6000806000606084860312156123a957600080fd5b83356123b481612261565b95602085013595506040909401359392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561240b5761240b6123c9565b500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615612448576124486123c9565b500290565b60006020828403121561245f57600080fd5b81516122a381612261565b60008282101561247c5761247c6123c9565b500390565b6000826124b7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c634300080f000a"

// DeployRATContract deploys a new RAT contract, binding an instance of RAT to it.
func DeployRATContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, error) {
	parsed, err := RATMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, errors.New("GetABI returned nil")
	}

	address, tx, _, err := bind.DeployContract(auth, *parsed, common.FromHex(RATBin), backend)
	if err != nil {
		return common.Address{}, nil, err
	}
	return address, tx, nil
}
