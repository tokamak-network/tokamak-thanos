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

// UnsupportedProtocolMetaData contains all meta data concerning the UnsupportedProtocol contract.
var UnsupportedProtocolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"UnsupportedProtocolError\",\"type\":\"error\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"}]",
	Bin: "0x608080604052346013576043908160198239f35b600080fdfe60808060405234603157807fea3559ef0000000000000000000000000000000000000000000000000000000060049252fd5b600080fdfea164736f6c6343000811000a",
}

// UnsupportedProtocolABI is the input ABI used to generate the binding from.
// Deprecated: Use UnsupportedProtocolMetaData.ABI instead.
var UnsupportedProtocolABI = UnsupportedProtocolMetaData.ABI

// UnsupportedProtocolBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UnsupportedProtocolMetaData.Bin instead.
var UnsupportedProtocolBin = UnsupportedProtocolMetaData.Bin

// DeployUnsupportedProtocol deploys a new Ethereum contract, binding an instance of UnsupportedProtocol to it.
func DeployUnsupportedProtocol(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UnsupportedProtocol, error) {
	parsed, err := UnsupportedProtocolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UnsupportedProtocolBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UnsupportedProtocol{UnsupportedProtocolCaller: UnsupportedProtocolCaller{contract: contract}, UnsupportedProtocolTransactor: UnsupportedProtocolTransactor{contract: contract}, UnsupportedProtocolFilterer: UnsupportedProtocolFilterer{contract: contract}}, nil
}

// UnsupportedProtocol is an auto generated Go binding around an Ethereum contract.
type UnsupportedProtocol struct {
	UnsupportedProtocolCaller     // Read-only binding to the contract
	UnsupportedProtocolTransactor // Write-only binding to the contract
	UnsupportedProtocolFilterer   // Log filterer for contract events
}

// UnsupportedProtocolCaller is an auto generated read-only Go binding around an Ethereum contract.
type UnsupportedProtocolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UnsupportedProtocolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UnsupportedProtocolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UnsupportedProtocolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UnsupportedProtocolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UnsupportedProtocolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UnsupportedProtocolSession struct {
	Contract     *UnsupportedProtocol // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// UnsupportedProtocolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UnsupportedProtocolCallerSession struct {
	Contract *UnsupportedProtocolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// UnsupportedProtocolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UnsupportedProtocolTransactorSession struct {
	Contract     *UnsupportedProtocolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// UnsupportedProtocolRaw is an auto generated low-level Go binding around an Ethereum contract.
type UnsupportedProtocolRaw struct {
	Contract *UnsupportedProtocol // Generic contract binding to access the raw methods on
}

// UnsupportedProtocolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UnsupportedProtocolCallerRaw struct {
	Contract *UnsupportedProtocolCaller // Generic read-only contract binding to access the raw methods on
}

// UnsupportedProtocolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UnsupportedProtocolTransactorRaw struct {
	Contract *UnsupportedProtocolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUnsupportedProtocol creates a new instance of UnsupportedProtocol, bound to a specific deployed contract.
func NewUnsupportedProtocol(address common.Address, backend bind.ContractBackend) (*UnsupportedProtocol, error) {
	contract, err := bindUnsupportedProtocol(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UnsupportedProtocol{UnsupportedProtocolCaller: UnsupportedProtocolCaller{contract: contract}, UnsupportedProtocolTransactor: UnsupportedProtocolTransactor{contract: contract}, UnsupportedProtocolFilterer: UnsupportedProtocolFilterer{contract: contract}}, nil
}

// NewUnsupportedProtocolCaller creates a new read-only instance of UnsupportedProtocol, bound to a specific deployed contract.
func NewUnsupportedProtocolCaller(address common.Address, caller bind.ContractCaller) (*UnsupportedProtocolCaller, error) {
	contract, err := bindUnsupportedProtocol(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UnsupportedProtocolCaller{contract: contract}, nil
}

// NewUnsupportedProtocolTransactor creates a new write-only instance of UnsupportedProtocol, bound to a specific deployed contract.
func NewUnsupportedProtocolTransactor(address common.Address, transactor bind.ContractTransactor) (*UnsupportedProtocolTransactor, error) {
	contract, err := bindUnsupportedProtocol(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UnsupportedProtocolTransactor{contract: contract}, nil
}

// NewUnsupportedProtocolFilterer creates a new log filterer instance of UnsupportedProtocol, bound to a specific deployed contract.
func NewUnsupportedProtocolFilterer(address common.Address, filterer bind.ContractFilterer) (*UnsupportedProtocolFilterer, error) {
	contract, err := bindUnsupportedProtocol(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UnsupportedProtocolFilterer{contract: contract}, nil
}

// bindUnsupportedProtocol binds a generic wrapper to an already deployed contract.
func bindUnsupportedProtocol(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UnsupportedProtocolABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UnsupportedProtocol *UnsupportedProtocolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UnsupportedProtocol.Contract.UnsupportedProtocolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UnsupportedProtocol *UnsupportedProtocolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnsupportedProtocol.Contract.UnsupportedProtocolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UnsupportedProtocol *UnsupportedProtocolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UnsupportedProtocol.Contract.UnsupportedProtocolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UnsupportedProtocol *UnsupportedProtocolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UnsupportedProtocol.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UnsupportedProtocol *UnsupportedProtocolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnsupportedProtocol.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UnsupportedProtocol *UnsupportedProtocolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UnsupportedProtocol.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_UnsupportedProtocol *UnsupportedProtocolTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _UnsupportedProtocol.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_UnsupportedProtocol *UnsupportedProtocolSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _UnsupportedProtocol.Contract.Fallback(&_UnsupportedProtocol.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_UnsupportedProtocol *UnsupportedProtocolTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _UnsupportedProtocol.Contract.Fallback(&_UnsupportedProtocol.TransactOpts, calldata)
}
