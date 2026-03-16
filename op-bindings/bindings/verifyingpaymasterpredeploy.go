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

// VerifyingPaymasterPredeployMetaData contains all meta data concerning the VerifyingPaymasterPredeploy contract.
var VerifyingPaymasterPredeployMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addStake\",\"inputs\":[{\"name\":\"unstakeDelaySec\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"entryPoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDeposit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHash\",\"inputs\":[{\"name\":\"userOp\",\"type\":\"tuple\",\"internalType\":\"structPackedUserOperation\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"accountGasLimits\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preVerificationGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFees\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"validUntil\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"validAfter\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_entryPoint\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"},{\"name\":\"_verifyingSigner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parsePaymasterAndData\",\"inputs\":[{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"validUntil\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"validAfter\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postOp\",\"inputs\":[{\"name\":\"mode\",\"type\":\"uint8\",\"internalType\":\"enumIPaymaster.PostOpMode\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actualGasCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualUserOpFeePerGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVerifyingSigner\",\"inputs\":[{\"name\":\"_newSigner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatePaymasterUserOp\",\"inputs\":[{\"name\":\"userOp\",\"type\":\"tuple\",\"internalType\":\"structPackedUserOperation\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"accountGasLimits\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preVerificationGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFees\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"userOpHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validationData\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyingSigner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawStake\",\"inputs\":[{\"name\":\"withdrawAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTo\",\"inputs\":[{\"name\":\"withdrawAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifyingSignerUpdated\",\"inputs\":[{\"name\":\"oldSigner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSigner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// VerifyingPaymasterPredeployABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifyingPaymasterPredeployMetaData.ABI instead.
var VerifyingPaymasterPredeployABI = VerifyingPaymasterPredeployMetaData.ABI

// VerifyingPaymasterPredeploy is an auto generated Go binding around an Ethereum contract.
type VerifyingPaymasterPredeploy struct {
	VerifyingPaymasterPredeployCaller     // Read-only binding to the contract
	VerifyingPaymasterPredeployTransactor // Write-only binding to the contract
	VerifyingPaymasterPredeployFilterer   // Log filterer for contract events
}

// VerifyingPaymasterPredeployCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifyingPaymasterPredeployCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifyingPaymasterPredeployTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifyingPaymasterPredeployTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifyingPaymasterPredeployFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifyingPaymasterPredeployFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifyingPaymasterPredeploySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifyingPaymasterPredeploySession struct {
	Contract     *VerifyingPaymasterPredeploy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// VerifyingPaymasterPredeployCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifyingPaymasterPredeployCallerSession struct {
	Contract *VerifyingPaymasterPredeployCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// VerifyingPaymasterPredeployTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifyingPaymasterPredeployTransactorSession struct {
	Contract     *VerifyingPaymasterPredeployTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// VerifyingPaymasterPredeployRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifyingPaymasterPredeployRaw struct {
	Contract *VerifyingPaymasterPredeploy // Generic contract binding to access the raw methods on
}

// VerifyingPaymasterPredeployCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifyingPaymasterPredeployCallerRaw struct {
	Contract *VerifyingPaymasterPredeployCaller // Generic read-only contract binding to access the raw methods on
}

// VerifyingPaymasterPredeployTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifyingPaymasterPredeployTransactorRaw struct {
	Contract *VerifyingPaymasterPredeployTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifyingPaymasterPredeploy creates a new instance of VerifyingPaymasterPredeploy, bound to a specific deployed contract.
func NewVerifyingPaymasterPredeploy(address common.Address, backend bind.ContractBackend) (*VerifyingPaymasterPredeploy, error) {
	contract, err := bindVerifyingPaymasterPredeploy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeploy{VerifyingPaymasterPredeployCaller: VerifyingPaymasterPredeployCaller{contract: contract}, VerifyingPaymasterPredeployTransactor: VerifyingPaymasterPredeployTransactor{contract: contract}, VerifyingPaymasterPredeployFilterer: VerifyingPaymasterPredeployFilterer{contract: contract}}, nil
}

// NewVerifyingPaymasterPredeployCaller creates a new read-only instance of VerifyingPaymasterPredeploy, bound to a specific deployed contract.
func NewVerifyingPaymasterPredeployCaller(address common.Address, caller bind.ContractCaller) (*VerifyingPaymasterPredeployCaller, error) {
	contract, err := bindVerifyingPaymasterPredeploy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeployCaller{contract: contract}, nil
}

// NewVerifyingPaymasterPredeployTransactor creates a new write-only instance of VerifyingPaymasterPredeploy, bound to a specific deployed contract.
func NewVerifyingPaymasterPredeployTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifyingPaymasterPredeployTransactor, error) {
	contract, err := bindVerifyingPaymasterPredeploy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeployTransactor{contract: contract}, nil
}

// NewVerifyingPaymasterPredeployFilterer creates a new log filterer instance of VerifyingPaymasterPredeploy, bound to a specific deployed contract.
func NewVerifyingPaymasterPredeployFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifyingPaymasterPredeployFilterer, error) {
	contract, err := bindVerifyingPaymasterPredeploy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeployFilterer{contract: contract}, nil
}

// bindVerifyingPaymasterPredeploy binds a generic wrapper to an already deployed contract.
func bindVerifyingPaymasterPredeploy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VerifyingPaymasterPredeployABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifyingPaymasterPredeploy.Contract.VerifyingPaymasterPredeployCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.VerifyingPaymasterPredeployTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.VerifyingPaymasterPredeployTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifyingPaymasterPredeploy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.contract.Transact(opts, method, params...)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCaller) EntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifyingPaymasterPredeploy.contract.Call(opts, &out, "entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) EntryPoint() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.EntryPoint(&_VerifyingPaymasterPredeploy.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerSession) EntryPoint() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.EntryPoint(&_VerifyingPaymasterPredeploy.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCaller) GetDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifyingPaymasterPredeploy.contract.Call(opts, &out, "getDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) GetDeposit() (*big.Int, error) {
	return _VerifyingPaymasterPredeploy.Contract.GetDeposit(&_VerifyingPaymasterPredeploy.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerSession) GetDeposit() (*big.Int, error) {
	return _VerifyingPaymasterPredeploy.Contract.GetDeposit(&_VerifyingPaymasterPredeploy.CallOpts)
}

// GetHash is a free data retrieval call binding the contract method 0x5829c5f5.
//
// Solidity: function getHash((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, uint48 validUntil, uint48 validAfter) view returns(bytes32)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCaller) GetHash(opts *bind.CallOpts, userOp PackedUserOperation, validUntil *big.Int, validAfter *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VerifyingPaymasterPredeploy.contract.Call(opts, &out, "getHash", userOp, validUntil, validAfter)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetHash is a free data retrieval call binding the contract method 0x5829c5f5.
//
// Solidity: function getHash((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, uint48 validUntil, uint48 validAfter) view returns(bytes32)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) GetHash(userOp PackedUserOperation, validUntil *big.Int, validAfter *big.Int) ([32]byte, error) {
	return _VerifyingPaymasterPredeploy.Contract.GetHash(&_VerifyingPaymasterPredeploy.CallOpts, userOp, validUntil, validAfter)
}

// GetHash is a free data retrieval call binding the contract method 0x5829c5f5.
//
// Solidity: function getHash((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, uint48 validUntil, uint48 validAfter) view returns(bytes32)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerSession) GetHash(userOp PackedUserOperation, validUntil *big.Int, validAfter *big.Int) ([32]byte, error) {
	return _VerifyingPaymasterPredeploy.Contract.GetHash(&_VerifyingPaymasterPredeploy.CallOpts, userOp, validUntil, validAfter)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifyingPaymasterPredeploy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) Owner() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.Owner(&_VerifyingPaymasterPredeploy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerSession) Owner() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.Owner(&_VerifyingPaymasterPredeploy.CallOpts)
}

// ParsePaymasterAndData is a free data retrieval call binding the contract method 0x94d4ad60.
//
// Solidity: function parsePaymasterAndData(bytes paymasterAndData) pure returns(uint48 validUntil, uint48 validAfter)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCaller) ParsePaymasterAndData(opts *bind.CallOpts, paymasterAndData []byte) (struct {
	ValidUntil *big.Int
	ValidAfter *big.Int
}, error) {
	var out []interface{}
	err := _VerifyingPaymasterPredeploy.contract.Call(opts, &out, "parsePaymasterAndData", paymasterAndData)

	outstruct := new(struct {
		ValidUntil *big.Int
		ValidAfter *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ValidUntil = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ValidAfter = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ParsePaymasterAndData is a free data retrieval call binding the contract method 0x94d4ad60.
//
// Solidity: function parsePaymasterAndData(bytes paymasterAndData) pure returns(uint48 validUntil, uint48 validAfter)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) ParsePaymasterAndData(paymasterAndData []byte) (struct {
	ValidUntil *big.Int
	ValidAfter *big.Int
}, error) {
	return _VerifyingPaymasterPredeploy.Contract.ParsePaymasterAndData(&_VerifyingPaymasterPredeploy.CallOpts, paymasterAndData)
}

// ParsePaymasterAndData is a free data retrieval call binding the contract method 0x94d4ad60.
//
// Solidity: function parsePaymasterAndData(bytes paymasterAndData) pure returns(uint48 validUntil, uint48 validAfter)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerSession) ParsePaymasterAndData(paymasterAndData []byte) (struct {
	ValidUntil *big.Int
	ValidAfter *big.Int
}, error) {
	return _VerifyingPaymasterPredeploy.Contract.ParsePaymasterAndData(&_VerifyingPaymasterPredeploy.CallOpts, paymasterAndData)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifyingPaymasterPredeploy.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) PendingOwner() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.PendingOwner(&_VerifyingPaymasterPredeploy.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerSession) PendingOwner() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.PendingOwner(&_VerifyingPaymasterPredeploy.CallOpts)
}

// VerifyingSigner is a free data retrieval call binding the contract method 0x23d9ac9b.
//
// Solidity: function verifyingSigner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCaller) VerifyingSigner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifyingPaymasterPredeploy.contract.Call(opts, &out, "verifyingSigner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VerifyingSigner is a free data retrieval call binding the contract method 0x23d9ac9b.
//
// Solidity: function verifyingSigner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) VerifyingSigner() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.VerifyingSigner(&_VerifyingPaymasterPredeploy.CallOpts)
}

// VerifyingSigner is a free data retrieval call binding the contract method 0x23d9ac9b.
//
// Solidity: function verifyingSigner() view returns(address)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployCallerSession) VerifyingSigner() (common.Address, error) {
	return _VerifyingPaymasterPredeploy.Contract.VerifyingSigner(&_VerifyingPaymasterPredeploy.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.AcceptOwnership(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.AcceptOwnership(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) AddStake(opts *bind.TransactOpts, unstakeDelaySec uint32) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "addStake", unstakeDelaySec)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) AddStake(unstakeDelaySec uint32) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.AddStake(&_VerifyingPaymasterPredeploy.TransactOpts, unstakeDelaySec)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) AddStake(unstakeDelaySec uint32) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.AddStake(&_VerifyingPaymasterPredeploy.TransactOpts, unstakeDelaySec)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) Deposit() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.Deposit(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) Deposit() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.Deposit(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _entryPoint, address _verifyingSigner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) Initialize(opts *bind.TransactOpts, _entryPoint common.Address, _verifyingSigner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "initialize", _entryPoint, _verifyingSigner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _entryPoint, address _verifyingSigner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) Initialize(_entryPoint common.Address, _verifyingSigner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.Initialize(&_VerifyingPaymasterPredeploy.TransactOpts, _entryPoint, _verifyingSigner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _entryPoint, address _verifyingSigner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) Initialize(_entryPoint common.Address, _verifyingSigner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.Initialize(&_VerifyingPaymasterPredeploy.TransactOpts, _entryPoint, _verifyingSigner)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) PostOp(opts *bind.TransactOpts, mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) PostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.PostOp(&_VerifyingPaymasterPredeploy.TransactOpts, mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) PostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.PostOp(&_VerifyingPaymasterPredeploy.TransactOpts, mode, context, actualGasCost, actualUserOpFeePerGas)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) RenounceOwnership() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.RenounceOwnership(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.RenounceOwnership(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// SetVerifyingSigner is a paid mutator transaction binding the contract method 0xf5cba98c.
//
// Solidity: function setVerifyingSigner(address _newSigner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) SetVerifyingSigner(opts *bind.TransactOpts, _newSigner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "setVerifyingSigner", _newSigner)
}

// SetVerifyingSigner is a paid mutator transaction binding the contract method 0xf5cba98c.
//
// Solidity: function setVerifyingSigner(address _newSigner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) SetVerifyingSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.SetVerifyingSigner(&_VerifyingPaymasterPredeploy.TransactOpts, _newSigner)
}

// SetVerifyingSigner is a paid mutator transaction binding the contract method 0xf5cba98c.
//
// Solidity: function setVerifyingSigner(address _newSigner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) SetVerifyingSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.SetVerifyingSigner(&_VerifyingPaymasterPredeploy.TransactOpts, _newSigner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.TransferOwnership(&_VerifyingPaymasterPredeploy.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.TransferOwnership(&_VerifyingPaymasterPredeploy.TransactOpts, newOwner)
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) UnlockStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "unlockStake")
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) UnlockStake() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.UnlockStake(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) UnlockStake() (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.UnlockStake(&_VerifyingPaymasterPredeploy.TransactOpts)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "validatePaymasterUserOp", userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) ValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.ValidatePaymasterUserOp(&_VerifyingPaymasterPredeploy.TransactOpts, userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) ValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.ValidatePaymasterUserOp(&_VerifyingPaymasterPredeploy.TransactOpts, userOp, userOpHash, maxCost)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) WithdrawStake(opts *bind.TransactOpts, withdrawAddress common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "withdrawStake", withdrawAddress)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) WithdrawStake(withdrawAddress common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.WithdrawStake(&_VerifyingPaymasterPredeploy.TransactOpts, withdrawAddress)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) WithdrawStake(withdrawAddress common.Address) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.WithdrawStake(&_VerifyingPaymasterPredeploy.TransactOpts, withdrawAddress)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactor) WithdrawTo(opts *bind.TransactOpts, withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.contract.Transact(opts, "withdrawTo", withdrawAddress, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeploySession) WithdrawTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.WithdrawTo(&_VerifyingPaymasterPredeploy.TransactOpts, withdrawAddress, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployTransactorSession) WithdrawTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VerifyingPaymasterPredeploy.Contract.WithdrawTo(&_VerifyingPaymasterPredeploy.TransactOpts, withdrawAddress, amount)
}

// VerifyingPaymasterPredeployInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployInitializedIterator struct {
	Event *VerifyingPaymasterPredeployInitialized // Event containing the contract specifics and raw log

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
func (it *VerifyingPaymasterPredeployInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifyingPaymasterPredeployInitialized)
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
		it.Event = new(VerifyingPaymasterPredeployInitialized)
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
func (it *VerifyingPaymasterPredeployInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifyingPaymasterPredeployInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifyingPaymasterPredeployInitialized represents a Initialized event raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) FilterInitialized(opts *bind.FilterOpts) (*VerifyingPaymasterPredeployInitializedIterator, error) {

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeployInitializedIterator{contract: _VerifyingPaymasterPredeploy.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VerifyingPaymasterPredeployInitialized) (event.Subscription, error) {

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifyingPaymasterPredeployInitialized)
				if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) ParseInitialized(log types.Log) (*VerifyingPaymasterPredeployInitialized, error) {
	event := new(VerifyingPaymasterPredeployInitialized)
	if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifyingPaymasterPredeployOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployOwnershipTransferStartedIterator struct {
	Event *VerifyingPaymasterPredeployOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *VerifyingPaymasterPredeployOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifyingPaymasterPredeployOwnershipTransferStarted)
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
		it.Event = new(VerifyingPaymasterPredeployOwnershipTransferStarted)
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
func (it *VerifyingPaymasterPredeployOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifyingPaymasterPredeployOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifyingPaymasterPredeployOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VerifyingPaymasterPredeployOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeployOwnershipTransferStartedIterator{contract: _VerifyingPaymasterPredeploy.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *VerifyingPaymasterPredeployOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifyingPaymasterPredeployOwnershipTransferStarted)
				if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) ParseOwnershipTransferStarted(log types.Log) (*VerifyingPaymasterPredeployOwnershipTransferStarted, error) {
	event := new(VerifyingPaymasterPredeployOwnershipTransferStarted)
	if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifyingPaymasterPredeployOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployOwnershipTransferredIterator struct {
	Event *VerifyingPaymasterPredeployOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VerifyingPaymasterPredeployOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifyingPaymasterPredeployOwnershipTransferred)
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
		it.Event = new(VerifyingPaymasterPredeployOwnershipTransferred)
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
func (it *VerifyingPaymasterPredeployOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifyingPaymasterPredeployOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifyingPaymasterPredeployOwnershipTransferred represents a OwnershipTransferred event raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VerifyingPaymasterPredeployOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeployOwnershipTransferredIterator{contract: _VerifyingPaymasterPredeploy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifyingPaymasterPredeployOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifyingPaymasterPredeployOwnershipTransferred)
				if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) ParseOwnershipTransferred(log types.Log) (*VerifyingPaymasterPredeployOwnershipTransferred, error) {
	event := new(VerifyingPaymasterPredeployOwnershipTransferred)
	if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifyingPaymasterPredeployVerifyingSignerUpdatedIterator is returned from FilterVerifyingSignerUpdated and is used to iterate over the raw logs and unpacked data for VerifyingSignerUpdated events raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployVerifyingSignerUpdatedIterator struct {
	Event *VerifyingPaymasterPredeployVerifyingSignerUpdated // Event containing the contract specifics and raw log

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
func (it *VerifyingPaymasterPredeployVerifyingSignerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifyingPaymasterPredeployVerifyingSignerUpdated)
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
		it.Event = new(VerifyingPaymasterPredeployVerifyingSignerUpdated)
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
func (it *VerifyingPaymasterPredeployVerifyingSignerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifyingPaymasterPredeployVerifyingSignerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifyingPaymasterPredeployVerifyingSignerUpdated represents a VerifyingSignerUpdated event raised by the VerifyingPaymasterPredeploy contract.
type VerifyingPaymasterPredeployVerifyingSignerUpdated struct {
	OldSigner common.Address
	NewSigner common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVerifyingSignerUpdated is a free log retrieval operation binding the contract event 0xdeee0b46266dc0734b009b2ca66e2ac28d5fd72d38506a3fe462865f00c5df4e.
//
// Solidity: event VerifyingSignerUpdated(address indexed oldSigner, address indexed newSigner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) FilterVerifyingSignerUpdated(opts *bind.FilterOpts, oldSigner []common.Address, newSigner []common.Address) (*VerifyingPaymasterPredeployVerifyingSignerUpdatedIterator, error) {

	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.FilterLogs(opts, "VerifyingSignerUpdated", oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return &VerifyingPaymasterPredeployVerifyingSignerUpdatedIterator{contract: _VerifyingPaymasterPredeploy.contract, event: "VerifyingSignerUpdated", logs: logs, sub: sub}, nil
}

// WatchVerifyingSignerUpdated is a free log subscription operation binding the contract event 0xdeee0b46266dc0734b009b2ca66e2ac28d5fd72d38506a3fe462865f00c5df4e.
//
// Solidity: event VerifyingSignerUpdated(address indexed oldSigner, address indexed newSigner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) WatchVerifyingSignerUpdated(opts *bind.WatchOpts, sink chan<- *VerifyingPaymasterPredeployVerifyingSignerUpdated, oldSigner []common.Address, newSigner []common.Address) (event.Subscription, error) {

	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}

	logs, sub, err := _VerifyingPaymasterPredeploy.contract.WatchLogs(opts, "VerifyingSignerUpdated", oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifyingPaymasterPredeployVerifyingSignerUpdated)
				if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "VerifyingSignerUpdated", log); err != nil {
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

// ParseVerifyingSignerUpdated is a log parse operation binding the contract event 0xdeee0b46266dc0734b009b2ca66e2ac28d5fd72d38506a3fe462865f00c5df4e.
//
// Solidity: event VerifyingSignerUpdated(address indexed oldSigner, address indexed newSigner)
func (_VerifyingPaymasterPredeploy *VerifyingPaymasterPredeployFilterer) ParseVerifyingSignerUpdated(log types.Log) (*VerifyingPaymasterPredeployVerifyingSignerUpdated, error) {
	event := new(VerifyingPaymasterPredeployVerifyingSignerUpdated)
	if err := _VerifyingPaymasterPredeploy.contract.UnpackLog(event, "VerifyingSignerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
