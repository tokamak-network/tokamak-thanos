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

// LegacyERC20NativeTokenMetaData contains all meta data concerning the LegacyERC20NativeToken contract.
var LegacyERC20NativeTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BRIDGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"REMOTE_TOKEN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"_who\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bridge\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"l1Token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l2Bridge\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"remoteToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"_interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Burn\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Mint\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60e06040523480156200001157600080fd5b5073420000000000000000000000000000000000001060006040518060400160405280600581526020016422ba3432b960d91b8152506040518060400160405280600381526020016208aa8960eb1b81525060128282816003908162000078919062000152565b50600462000087828262000152565b5050506001600160a01b039384166080529390921660a052505060ff1660c0526200021e565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620000d857607f821691505b602082108103620000f957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200014d57600081815260208120601f850160051c81016020861015620001285750805b601f850160051c820191505b81811015620001495782815560010162000134565b5050505b505050565b81516001600160401b038111156200016e576200016e620000ad565b62000186816200017f8454620000c3565b84620000ff565b602080601f831160018114620001be5760008415620001a55750858301515b600019600386901b1c1916600185901b17855562000149565b600085815260208120601f198616915b82811015620001ef57888601518255948401946001909101908401620001ce565b50858210156200020e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c051610b4b6200025c600039600061024401526000818161033d01526103d20152600081816101a901526103630152610b4b6000f3fe608060405234801561001057600080fd5b50600436106101775760003560e01c806370a08231116100d8578063ae1f6aaf1161008c578063dd62ed3e11610066578063dd62ed3e14610387578063e78cea921461033b578063ee9a31a2146103cd57600080fd5b8063ae1f6aaf1461033b578063c01e1bd614610361578063d6c0b2c41461036157600080fd5b80639dc29fac116100bd5780639dc29fac14610302578063a457c2d714610315578063a9059cbb1461032857600080fd5b806370a08231146102d257806395d89b41146102fa57600080fd5b806323b872dd1161012f5780633950935111610114578063395093511461026e57806340c10f191461028157806354fd4d501461029657600080fd5b806323b872dd1461022a578063313ce5671461023d57600080fd5b806306fdde031161016057806306fdde03146101f0578063095ea7b31461020557806318160ddd1461021857600080fd5b806301ffc9a71461017c578063033964be146101a4575b600080fd5b61018f61018a366004610952565b6103f4565b60405190151581526020015b60405180910390f35b6101cb7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161019b565b6101f86104e5565b60405161019b919061099b565b61018f610213366004610a37565b610577565b6002545b60405190815260200161019b565b61018f610238366004610a61565b610607565b60405160ff7f000000000000000000000000000000000000000000000000000000000000000016815260200161019b565b61018f61027c366004610a37565b610692565b61029461028f366004610a37565b61071d565b005b6101f86040518060400160405280600581526020017f312e332e3000000000000000000000000000000000000000000000000000000081525081565b61021c6102e0366004610a9d565b73ffffffffffffffffffffffffffffffffffffffff163190565b6101f86107a5565b610294610310366004610a37565b6107b4565b61018f610323366004610a37565b61083c565b61018f610336366004610a37565b6108c7565b7f00000000000000000000000000000000000000000000000000000000000000006101cb565b7f00000000000000000000000000000000000000000000000000000000000000006101cb565b61021c610395366004610ab8565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b6101cb7f000000000000000000000000000000000000000000000000000000000000000081565b60007f01ffc9a7000000000000000000000000000000000000000000000000000000007f1d1d8b63000000000000000000000000000000000000000000000000000000007fec4fc8e3000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000085168314806104ad57507fffffffff00000000000000000000000000000000000000000000000000000000858116908316145b806104dc57507fffffffff00000000000000000000000000000000000000000000000000000000858116908216145b95945050505050565b6060600380546104f490610aeb565b80601f016020809104026020016040519081016040528092919081815260200182805461052090610aeb565b801561056d5780601f106105425761010080835404028352916020019161056d565b820191906000526020600020905b81548152906001019060200180831161055057829003601f168201915b5050505050905090565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f4c656761637945524332304e6174697665546f6b656e3a20617070726f76652060448201527f69732064697361626c656400000000000000000000000000000000000000000060648201526000906084015b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603060248201527f4c656761637945524332304e6174697665546f6b656e3a207472616e7366657260448201527f46726f6d2069732064697361626c65640000000000000000000000000000000060648201526000906084016105fe565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603560248201527f4c656761637945524332304e6174697665546f6b656e3a20696e63726561736560448201527f416c6c6f77616e63652069732064697361626c6564000000000000000000000060648201526000906084016105fe565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f4c656761637945524332304e6174697665546f6b656e3a206d696e742069732060448201527f64697361626c656400000000000000000000000000000000000000000000000060648201526084016105fe565b6060600480546104f490610aeb565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f4c656761637945524332304e6174697665546f6b656e3a206275726e2069732060448201527f64697361626c656400000000000000000000000000000000000000000000000060648201526084016105fe565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603560248201527f4c656761637945524332304e6174697665546f6b656e3a20646563726561736560448201527f416c6c6f77616e63652069732064697361626c6564000000000000000000000060648201526000906084016105fe565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f4c656761637945524332304e6174697665546f6b656e3a207472616e7366657260448201527f2069732064697361626c6564000000000000000000000000000000000000000060648201526000906084016105fe565b60006020828403121561096457600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461099457600080fd5b9392505050565b600060208083528351808285015260005b818110156109c8578581018301518582016040015282016109ac565b818111156109da576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610a3257600080fd5b919050565b60008060408385031215610a4a57600080fd5b610a5383610a0e565b946020939093013593505050565b600080600060608486031215610a7657600080fd5b610a7f84610a0e565b9250610a8d60208501610a0e565b9150604084013590509250925092565b600060208284031215610aaf57600080fd5b61099482610a0e565b60008060408385031215610acb57600080fd5b610ad483610a0e565b9150610ae260208401610a0e565b90509250929050565b600181811c90821680610aff57607f821691505b602082108103610b38577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b5091905056fea164736f6c634300080f000a",
}

// LegacyERC20NativeTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use LegacyERC20NativeTokenMetaData.ABI instead.
var LegacyERC20NativeTokenABI = LegacyERC20NativeTokenMetaData.ABI

// LegacyERC20NativeTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LegacyERC20NativeTokenMetaData.Bin instead.
var LegacyERC20NativeTokenBin = LegacyERC20NativeTokenMetaData.Bin

// DeployLegacyERC20NativeToken deploys a new Ethereum contract, binding an instance of LegacyERC20NativeToken to it.
func DeployLegacyERC20NativeToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LegacyERC20NativeToken, error) {
	parsed, err := LegacyERC20NativeTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LegacyERC20NativeTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LegacyERC20NativeToken{LegacyERC20NativeTokenCaller: LegacyERC20NativeTokenCaller{contract: contract}, LegacyERC20NativeTokenTransactor: LegacyERC20NativeTokenTransactor{contract: contract}, LegacyERC20NativeTokenFilterer: LegacyERC20NativeTokenFilterer{contract: contract}}, nil
}

// LegacyERC20NativeToken is an auto generated Go binding around an Ethereum contract.
type LegacyERC20NativeToken struct {
	LegacyERC20NativeTokenCaller     // Read-only binding to the contract
	LegacyERC20NativeTokenTransactor // Write-only binding to the contract
	LegacyERC20NativeTokenFilterer   // Log filterer for contract events
}

// LegacyERC20NativeTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type LegacyERC20NativeTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LegacyERC20NativeTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LegacyERC20NativeTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LegacyERC20NativeTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LegacyERC20NativeTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LegacyERC20NativeTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LegacyERC20NativeTokenSession struct {
	Contract     *LegacyERC20NativeToken // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// LegacyERC20NativeTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LegacyERC20NativeTokenCallerSession struct {
	Contract *LegacyERC20NativeTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// LegacyERC20NativeTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LegacyERC20NativeTokenTransactorSession struct {
	Contract     *LegacyERC20NativeTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// LegacyERC20NativeTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type LegacyERC20NativeTokenRaw struct {
	Contract *LegacyERC20NativeToken // Generic contract binding to access the raw methods on
}

// LegacyERC20NativeTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LegacyERC20NativeTokenCallerRaw struct {
	Contract *LegacyERC20NativeTokenCaller // Generic read-only contract binding to access the raw methods on
}

// LegacyERC20NativeTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LegacyERC20NativeTokenTransactorRaw struct {
	Contract *LegacyERC20NativeTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLegacyERC20NativeToken creates a new instance of LegacyERC20NativeToken, bound to a specific deployed contract.
func NewLegacyERC20NativeToken(address common.Address, backend bind.ContractBackend) (*LegacyERC20NativeToken, error) {
	contract, err := bindLegacyERC20NativeToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeToken{LegacyERC20NativeTokenCaller: LegacyERC20NativeTokenCaller{contract: contract}, LegacyERC20NativeTokenTransactor: LegacyERC20NativeTokenTransactor{contract: contract}, LegacyERC20NativeTokenFilterer: LegacyERC20NativeTokenFilterer{contract: contract}}, nil
}

// NewLegacyERC20NativeTokenCaller creates a new read-only instance of LegacyERC20NativeToken, bound to a specific deployed contract.
func NewLegacyERC20NativeTokenCaller(address common.Address, caller bind.ContractCaller) (*LegacyERC20NativeTokenCaller, error) {
	contract, err := bindLegacyERC20NativeToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeTokenCaller{contract: contract}, nil
}

// NewLegacyERC20NativeTokenTransactor creates a new write-only instance of LegacyERC20NativeToken, bound to a specific deployed contract.
func NewLegacyERC20NativeTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*LegacyERC20NativeTokenTransactor, error) {
	contract, err := bindLegacyERC20NativeToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeTokenTransactor{contract: contract}, nil
}

// NewLegacyERC20NativeTokenFilterer creates a new log filterer instance of LegacyERC20NativeToken, bound to a specific deployed contract.
func NewLegacyERC20NativeTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*LegacyERC20NativeTokenFilterer, error) {
	contract, err := bindLegacyERC20NativeToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeTokenFilterer{contract: contract}, nil
}

// bindLegacyERC20NativeToken binds a generic wrapper to an already deployed contract.
func bindLegacyERC20NativeToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LegacyERC20NativeTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LegacyERC20NativeToken.Contract.LegacyERC20NativeTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.LegacyERC20NativeTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.LegacyERC20NativeTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LegacyERC20NativeToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.contract.Transact(opts, method, params...)
}

// BRIDGE is a free data retrieval call binding the contract method 0xee9a31a2.
//
// Solidity: function BRIDGE() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) BRIDGE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "BRIDGE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BRIDGE is a free data retrieval call binding the contract method 0xee9a31a2.
//
// Solidity: function BRIDGE() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) BRIDGE() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.BRIDGE(&_LegacyERC20NativeToken.CallOpts)
}

// BRIDGE is a free data retrieval call binding the contract method 0xee9a31a2.
//
// Solidity: function BRIDGE() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) BRIDGE() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.BRIDGE(&_LegacyERC20NativeToken.CallOpts)
}

// REMOTETOKEN is a free data retrieval call binding the contract method 0x033964be.
//
// Solidity: function REMOTE_TOKEN() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) REMOTETOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "REMOTE_TOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// REMOTETOKEN is a free data retrieval call binding the contract method 0x033964be.
//
// Solidity: function REMOTE_TOKEN() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) REMOTETOKEN() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.REMOTETOKEN(&_LegacyERC20NativeToken.CallOpts)
}

// REMOTETOKEN is a free data retrieval call binding the contract method 0x033964be.
//
// Solidity: function REMOTE_TOKEN() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) REMOTETOKEN() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.REMOTETOKEN(&_LegacyERC20NativeToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LegacyERC20NativeToken.Contract.Allowance(&_LegacyERC20NativeToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LegacyERC20NativeToken.Contract.Allowance(&_LegacyERC20NativeToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _who) view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) BalanceOf(opts *bind.CallOpts, _who common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "balanceOf", _who)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _who) view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) BalanceOf(_who common.Address) (*big.Int, error) {
	return _LegacyERC20NativeToken.Contract.BalanceOf(&_LegacyERC20NativeToken.CallOpts, _who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _who) view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) BalanceOf(_who common.Address) (*big.Int, error) {
	return _LegacyERC20NativeToken.Contract.BalanceOf(&_LegacyERC20NativeToken.CallOpts, _who)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Bridge() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.Bridge(&_LegacyERC20NativeToken.CallOpts)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) Bridge() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.Bridge(&_LegacyERC20NativeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Decimals() (uint8, error) {
	return _LegacyERC20NativeToken.Contract.Decimals(&_LegacyERC20NativeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) Decimals() (uint8, error) {
	return _LegacyERC20NativeToken.Contract.Decimals(&_LegacyERC20NativeToken.CallOpts)
}

// L1Token is a free data retrieval call binding the contract method 0xc01e1bd6.
//
// Solidity: function l1Token() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) L1Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "l1Token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Token is a free data retrieval call binding the contract method 0xc01e1bd6.
//
// Solidity: function l1Token() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) L1Token() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.L1Token(&_LegacyERC20NativeToken.CallOpts)
}

// L1Token is a free data retrieval call binding the contract method 0xc01e1bd6.
//
// Solidity: function l1Token() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) L1Token() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.L1Token(&_LegacyERC20NativeToken.CallOpts)
}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) L2Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "l2Bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) L2Bridge() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.L2Bridge(&_LegacyERC20NativeToken.CallOpts)
}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) L2Bridge() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.L2Bridge(&_LegacyERC20NativeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Name() (string, error) {
	return _LegacyERC20NativeToken.Contract.Name(&_LegacyERC20NativeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) Name() (string, error) {
	return _LegacyERC20NativeToken.Contract.Name(&_LegacyERC20NativeToken.CallOpts)
}

// RemoteToken is a free data retrieval call binding the contract method 0xd6c0b2c4.
//
// Solidity: function remoteToken() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) RemoteToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "remoteToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RemoteToken is a free data retrieval call binding the contract method 0xd6c0b2c4.
//
// Solidity: function remoteToken() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) RemoteToken() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.RemoteToken(&_LegacyERC20NativeToken.CallOpts)
}

// RemoteToken is a free data retrieval call binding the contract method 0xd6c0b2c4.
//
// Solidity: function remoteToken() view returns(address)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) RemoteToken() (common.Address, error) {
	return _LegacyERC20NativeToken.Contract.RemoteToken(&_LegacyERC20NativeToken.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) pure returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) SupportsInterface(opts *bind.CallOpts, _interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "supportsInterface", _interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) pure returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _LegacyERC20NativeToken.Contract.SupportsInterface(&_LegacyERC20NativeToken.CallOpts, _interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) pure returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _LegacyERC20NativeToken.Contract.SupportsInterface(&_LegacyERC20NativeToken.CallOpts, _interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Symbol() (string, error) {
	return _LegacyERC20NativeToken.Contract.Symbol(&_LegacyERC20NativeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) Symbol() (string, error) {
	return _LegacyERC20NativeToken.Contract.Symbol(&_LegacyERC20NativeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) TotalSupply() (*big.Int, error) {
	return _LegacyERC20NativeToken.Contract.TotalSupply(&_LegacyERC20NativeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _LegacyERC20NativeToken.Contract.TotalSupply(&_LegacyERC20NativeToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LegacyERC20NativeToken.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Version() (string, error) {
	return _LegacyERC20NativeToken.Contract.Version(&_LegacyERC20NativeToken.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenCallerSession) Version() (string, error) {
	return _LegacyERC20NativeToken.Contract.Version(&_LegacyERC20NativeToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactor) Approve(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.contract.Transact(opts, "approve", arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Approve(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Approve(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address , uint256 ) returns()
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactor) Burn(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.contract.Transact(opts, "burn", arg0, arg1)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address , uint256 ) returns()
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Burn(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Burn(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address , uint256 ) returns()
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorSession) Burn(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Burn(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.contract.Transact(opts, "decreaseAllowance", arg0, arg1)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) DecreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.DecreaseAllowance(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorSession) DecreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.DecreaseAllowance(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.contract.Transact(opts, "increaseAllowance", arg0, arg1)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) IncreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.IncreaseAllowance(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorSession) IncreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.IncreaseAllowance(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address , uint256 ) returns()
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactor) Mint(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.contract.Transact(opts, "mint", arg0, arg1)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address , uint256 ) returns()
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Mint(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Mint(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address , uint256 ) returns()
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorSession) Mint(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Mint(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactor) Transfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.contract.Transact(opts, "transfer", arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Transfer(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.Transfer(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactor) TransferFrom(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.contract.Transact(opts, "transferFrom", arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.TransferFrom(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenTransactorSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _LegacyERC20NativeToken.Contract.TransferFrom(&_LegacyERC20NativeToken.TransactOpts, arg0, arg1, arg2)
}

// LegacyERC20NativeTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenApprovalIterator struct {
	Event *LegacyERC20NativeTokenApproval // Event containing the contract specifics and raw log

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
func (it *LegacyERC20NativeTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LegacyERC20NativeTokenApproval)
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
		it.Event = new(LegacyERC20NativeTokenApproval)
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
func (it *LegacyERC20NativeTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LegacyERC20NativeTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LegacyERC20NativeTokenApproval represents a Approval event raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LegacyERC20NativeTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeTokenApprovalIterator{contract: _LegacyERC20NativeToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LegacyERC20NativeTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LegacyERC20NativeTokenApproval)
				if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) ParseApproval(log types.Log) (*LegacyERC20NativeTokenApproval, error) {
	event := new(LegacyERC20NativeTokenApproval)
	if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LegacyERC20NativeTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenBurnIterator struct {
	Event *LegacyERC20NativeTokenBurn // Event containing the contract specifics and raw log

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
func (it *LegacyERC20NativeTokenBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LegacyERC20NativeTokenBurn)
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
		it.Event = new(LegacyERC20NativeTokenBurn)
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
func (it *LegacyERC20NativeTokenBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LegacyERC20NativeTokenBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LegacyERC20NativeTokenBurn represents a Burn event raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenBurn struct {
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed account, uint256 amount)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) FilterBurn(opts *bind.FilterOpts, account []common.Address) (*LegacyERC20NativeTokenBurnIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.FilterLogs(opts, "Burn", accountRule)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeTokenBurnIterator{contract: _LegacyERC20NativeToken.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed account, uint256 amount)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *LegacyERC20NativeTokenBurn, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.WatchLogs(opts, "Burn", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LegacyERC20NativeTokenBurn)
				if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed account, uint256 amount)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) ParseBurn(log types.Log) (*LegacyERC20NativeTokenBurn, error) {
	event := new(LegacyERC20NativeTokenBurn)
	if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LegacyERC20NativeTokenMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenMintIterator struct {
	Event *LegacyERC20NativeTokenMint // Event containing the contract specifics and raw log

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
func (it *LegacyERC20NativeTokenMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LegacyERC20NativeTokenMint)
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
		it.Event = new(LegacyERC20NativeTokenMint)
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
func (it *LegacyERC20NativeTokenMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LegacyERC20NativeTokenMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LegacyERC20NativeTokenMint represents a Mint event raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenMint struct {
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed account, uint256 amount)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) FilterMint(opts *bind.FilterOpts, account []common.Address) (*LegacyERC20NativeTokenMintIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.FilterLogs(opts, "Mint", accountRule)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeTokenMintIterator{contract: _LegacyERC20NativeToken.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed account, uint256 amount)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *LegacyERC20NativeTokenMint, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.WatchLogs(opts, "Mint", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LegacyERC20NativeTokenMint)
				if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed account, uint256 amount)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) ParseMint(log types.Log) (*LegacyERC20NativeTokenMint, error) {
	event := new(LegacyERC20NativeTokenMint)
	if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LegacyERC20NativeTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenTransferIterator struct {
	Event *LegacyERC20NativeTokenTransfer // Event containing the contract specifics and raw log

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
func (it *LegacyERC20NativeTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LegacyERC20NativeTokenTransfer)
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
		it.Event = new(LegacyERC20NativeTokenTransfer)
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
func (it *LegacyERC20NativeTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LegacyERC20NativeTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LegacyERC20NativeTokenTransfer represents a Transfer event raised by the LegacyERC20NativeToken contract.
type LegacyERC20NativeTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LegacyERC20NativeTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LegacyERC20NativeTokenTransferIterator{contract: _LegacyERC20NativeToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LegacyERC20NativeTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LegacyERC20NativeToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LegacyERC20NativeTokenTransfer)
				if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LegacyERC20NativeToken *LegacyERC20NativeTokenFilterer) ParseTransfer(log types.Log) (*LegacyERC20NativeTokenTransfer, error) {
	event := new(LegacyERC20NativeTokenTransfer)
	if err := _LegacyERC20NativeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
