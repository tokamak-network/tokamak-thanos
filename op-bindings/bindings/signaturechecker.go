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

// SignatureCheckerMetaData contains all meta data concerning the SignatureChecker contract.
var SignatureCheckerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"isValidSignatureNow\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x6106a4610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100355760003560e01c80636ccea6521461003a575b600080fd5b6101026004803603606081101561005057600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561008d57600080fd5b82018360208201111561009f57600080fd5b803590602001918460018302840111640100000000831117156100c157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610116945050505050565b604080519115158252519081900360200190f35b600061012184610179565b610164578373ffffffffffffffffffffffffffffffffffffffff16610146848461017f565b73ffffffffffffffffffffffffffffffffffffffff16149050610172565b61016f848484610203565b90505b9392505050565b3b151590565b600081516041146101db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806106296023913960400191505060405180910390fd5b60208201516040830151606084015160001a6101f98682858561042d565b9695505050505050565b60008060608573ffffffffffffffffffffffffffffffffffffffff16631626ba7e60e01b86866040516024018083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561026f578181015183820152602001610257565b50505050905090810190601f16801561029c5780820380516001836020036101000a031916815260200191505b50604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009098169790971787525181519196909550859450925090508083835b6020831061036957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161032c565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855afa9150503d80600081146103c9576040519150601f19603f3d011682016040523d82523d6000602084013e6103ce565b606091505b50915091508180156103e257506020815110155b80156101f9575080517f1626ba7e00000000000000000000000000000000000000000000000000000000906020808401919081101561042057600080fd5b5051149695505050505050565b60007f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08211156104a8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806106726026913960400191505060405180910390fd5b8360ff16601b141580156104c057508360ff16601c14155b15610516576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602681526020018061064c6026913960400191505060405180910390fd5b600060018686868660405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610572573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff811661061f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f45435265636f7665723a20696e76616c6964207369676e617475726500000000604482015290519081900360640190fd5b9594505050505056fe45435265636f7665723a20696e76616c6964207369676e6174757265206c656e67746845435265636f7665723a20696e76616c6964207369676e6174757265202776272076616c756545435265636f7665723a20696e76616c6964207369676e6174757265202773272076616c7565a164736f6c634300060c000a",
}

// SignatureCheckerABI is the input ABI used to generate the binding from.
// Deprecated: Use SignatureCheckerMetaData.ABI instead.
var SignatureCheckerABI = SignatureCheckerMetaData.ABI

// SignatureCheckerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SignatureCheckerMetaData.Bin instead.
var SignatureCheckerBin = SignatureCheckerMetaData.Bin

// DeploySignatureChecker deploys a new Ethereum contract, binding an instance of SignatureChecker to it.
func DeploySignatureChecker(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SignatureChecker, error) {
	parsed, err := SignatureCheckerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SignatureCheckerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SignatureChecker{SignatureCheckerCaller: SignatureCheckerCaller{contract: contract}, SignatureCheckerTransactor: SignatureCheckerTransactor{contract: contract}, SignatureCheckerFilterer: SignatureCheckerFilterer{contract: contract}}, nil
}

// SignatureChecker is an auto generated Go binding around an Ethereum contract.
type SignatureChecker struct {
	SignatureCheckerCaller     // Read-only binding to the contract
	SignatureCheckerTransactor // Write-only binding to the contract
	SignatureCheckerFilterer   // Log filterer for contract events
}

// SignatureCheckerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignatureCheckerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureCheckerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignatureCheckerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureCheckerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignatureCheckerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureCheckerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignatureCheckerSession struct {
	Contract     *SignatureChecker // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SignatureCheckerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignatureCheckerCallerSession struct {
	Contract *SignatureCheckerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SignatureCheckerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignatureCheckerTransactorSession struct {
	Contract     *SignatureCheckerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SignatureCheckerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignatureCheckerRaw struct {
	Contract *SignatureChecker // Generic contract binding to access the raw methods on
}

// SignatureCheckerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignatureCheckerCallerRaw struct {
	Contract *SignatureCheckerCaller // Generic read-only contract binding to access the raw methods on
}

// SignatureCheckerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignatureCheckerTransactorRaw struct {
	Contract *SignatureCheckerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSignatureChecker creates a new instance of SignatureChecker, bound to a specific deployed contract.
func NewSignatureChecker(address common.Address, backend bind.ContractBackend) (*SignatureChecker, error) {
	contract, err := bindSignatureChecker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SignatureChecker{SignatureCheckerCaller: SignatureCheckerCaller{contract: contract}, SignatureCheckerTransactor: SignatureCheckerTransactor{contract: contract}, SignatureCheckerFilterer: SignatureCheckerFilterer{contract: contract}}, nil
}

// NewSignatureCheckerCaller creates a new read-only instance of SignatureChecker, bound to a specific deployed contract.
func NewSignatureCheckerCaller(address common.Address, caller bind.ContractCaller) (*SignatureCheckerCaller, error) {
	contract, err := bindSignatureChecker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignatureCheckerCaller{contract: contract}, nil
}

// NewSignatureCheckerTransactor creates a new write-only instance of SignatureChecker, bound to a specific deployed contract.
func NewSignatureCheckerTransactor(address common.Address, transactor bind.ContractTransactor) (*SignatureCheckerTransactor, error) {
	contract, err := bindSignatureChecker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignatureCheckerTransactor{contract: contract}, nil
}

// NewSignatureCheckerFilterer creates a new log filterer instance of SignatureChecker, bound to a specific deployed contract.
func NewSignatureCheckerFilterer(address common.Address, filterer bind.ContractFilterer) (*SignatureCheckerFilterer, error) {
	contract, err := bindSignatureChecker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignatureCheckerFilterer{contract: contract}, nil
}

// bindSignatureChecker binds a generic wrapper to an already deployed contract.
func bindSignatureChecker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SignatureCheckerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignatureChecker *SignatureCheckerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignatureChecker.Contract.SignatureCheckerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignatureChecker *SignatureCheckerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignatureChecker.Contract.SignatureCheckerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignatureChecker *SignatureCheckerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignatureChecker.Contract.SignatureCheckerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignatureChecker *SignatureCheckerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignatureChecker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignatureChecker *SignatureCheckerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignatureChecker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignatureChecker *SignatureCheckerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignatureChecker.Contract.contract.Transact(opts, method, params...)
}

// IsValidSignatureNow is a free data retrieval call binding the contract method 0x6ccea652.
//
// Solidity: function isValidSignatureNow(address signer, bytes32 digest, bytes signature) view returns(bool)
func (_SignatureChecker *SignatureCheckerCaller) IsValidSignatureNow(opts *bind.CallOpts, signer common.Address, digest [32]byte, signature []byte) (bool, error) {
	var out []interface{}
	err := _SignatureChecker.contract.Call(opts, &out, "isValidSignatureNow", signer, digest, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidSignatureNow is a free data retrieval call binding the contract method 0x6ccea652.
//
// Solidity: function isValidSignatureNow(address signer, bytes32 digest, bytes signature) view returns(bool)
func (_SignatureChecker *SignatureCheckerSession) IsValidSignatureNow(signer common.Address, digest [32]byte, signature []byte) (bool, error) {
	return _SignatureChecker.Contract.IsValidSignatureNow(&_SignatureChecker.CallOpts, signer, digest, signature)
}

// IsValidSignatureNow is a free data retrieval call binding the contract method 0x6ccea652.
//
// Solidity: function isValidSignatureNow(address signer, bytes32 digest, bytes signature) view returns(bool)
func (_SignatureChecker *SignatureCheckerCallerSession) IsValidSignatureNow(signer common.Address, digest [32]byte, signature []byte) (bool, error) {
	return _SignatureChecker.Contract.IsValidSignatureNow(&_SignatureChecker.CallOpts, signer, digest, signature)
}
