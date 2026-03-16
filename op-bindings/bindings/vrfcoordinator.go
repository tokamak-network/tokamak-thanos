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

// VRFCoordinatorMetaData contains all meta data concerning the VRFCoordinator contract.
var VRFCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterNode\",\"inputs\":[{\"name\":\"node\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fulfillRandomWords\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"randomWords\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRequestStatus\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"fulfilled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"randomWords\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerNode\",\"inputs\":[{\"name\":\"node\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registeredNodes\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRandomWords\",\"inputs\":[{\"name\":\"requester\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"numWords\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callbackGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requests\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requester\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"numWords\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callbackGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fulfilled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setPredeploy\",\"inputs\":[{\"name\":\"_vrfPredeploy\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"vrfPredeploy\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeDeregistered\",\"inputs\":[{\"name\":\"node\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeRegistered\",\"inputs\":[{\"name\":\"node\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RandomWordsFulfilled\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"randomWords\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RandomWordsRequested\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"numWords\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false}]",
}

// VRFCoordinatorABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFCoordinatorMetaData.ABI instead.
var VRFCoordinatorABI = VRFCoordinatorMetaData.ABI

// VRFCoordinator is an auto generated Go binding around an Ethereum contract.
type VRFCoordinator struct {
	VRFCoordinatorCaller     // Read-only binding to the contract
	VRFCoordinatorTransactor // Write-only binding to the contract
	VRFCoordinatorFilterer   // Log filterer for contract events
}

// VRFCoordinatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFCoordinatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFCoordinatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFCoordinatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFCoordinatorSession struct {
	Contract     *VRFCoordinator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFCoordinatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFCoordinatorCallerSession struct {
	Contract *VRFCoordinatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// VRFCoordinatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFCoordinatorTransactorSession struct {
	Contract     *VRFCoordinatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// VRFCoordinatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFCoordinatorRaw struct {
	Contract *VRFCoordinator // Generic contract binding to access the raw methods on
}

// VRFCoordinatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFCoordinatorCallerRaw struct {
	Contract *VRFCoordinatorCaller // Generic read-only contract binding to access the raw methods on
}

// VRFCoordinatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFCoordinatorTransactorRaw struct {
	Contract *VRFCoordinatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFCoordinator creates a new instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinator(address common.Address, backend bind.ContractBackend) (*VRFCoordinator, error) {
	contract, err := bindVRFCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinator{VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

// NewVRFCoordinatorCaller creates a new read-only instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorCaller, error) {
	contract, err := bindVRFCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCaller{contract: contract}, nil
}

// NewVRFCoordinatorTransactor creates a new write-only instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorTransactor, error) {
	contract, err := bindVRFCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTransactor{contract: contract}, nil
}

// NewVRFCoordinatorFilterer creates a new log filterer instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorFilterer, error) {
	contract, err := bindVRFCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorFilterer{contract: contract}, nil
}

// bindVRFCoordinator binds a generic wrapper to an already deployed contract.
func bindVRFCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinator *VRFCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.VRFCoordinatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinator *VRFCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinator *VRFCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinator *VRFCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_VRFCoordinator *VRFCoordinatorCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_VRFCoordinator *VRFCoordinatorSession) Admin() (common.Address, error) {
	return _VRFCoordinator.Contract.Admin(&_VRFCoordinator.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_VRFCoordinator *VRFCoordinatorCallerSession) Admin() (common.Address, error) {
	return _VRFCoordinator.Contract.Admin(&_VRFCoordinator.CallOpts)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFCoordinator *VRFCoordinatorCaller) GetRequestStatus(opts *bind.CallOpts, requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getRequestStatus", requestId)

	outstruct := new(struct {
		Fulfilled   bool
		RandomWords []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFCoordinator *VRFCoordinatorSession) GetRequestStatus(requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _VRFCoordinator.Contract.GetRequestStatus(&_VRFCoordinator.CallOpts, requestId)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFCoordinator *VRFCoordinatorCallerSession) GetRequestStatus(requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _VRFCoordinator.Contract.GetRequestStatus(&_VRFCoordinator.CallOpts, requestId)
}

// RegisteredNodes is a free data retrieval call binding the contract method 0xd3d1fb48.
//
// Solidity: function registeredNodes(address ) view returns(bool)
func (_VRFCoordinator *VRFCoordinatorCaller) RegisteredNodes(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "registeredNodes", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RegisteredNodes is a free data retrieval call binding the contract method 0xd3d1fb48.
//
// Solidity: function registeredNodes(address ) view returns(bool)
func (_VRFCoordinator *VRFCoordinatorSession) RegisteredNodes(arg0 common.Address) (bool, error) {
	return _VRFCoordinator.Contract.RegisteredNodes(&_VRFCoordinator.CallOpts, arg0)
}

// RegisteredNodes is a free data retrieval call binding the contract method 0xd3d1fb48.
//
// Solidity: function registeredNodes(address ) view returns(bool)
func (_VRFCoordinator *VRFCoordinatorCallerSession) RegisteredNodes(arg0 common.Address) (bool, error) {
	return _VRFCoordinator.Contract.RegisteredNodes(&_VRFCoordinator.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(address requester, uint32 numWords, uint256 callbackGasLimit, bool fulfilled)
func (_VRFCoordinator *VRFCoordinatorCaller) Requests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Requester        common.Address
	NumWords         uint32
	CallbackGasLimit *big.Int
	Fulfilled        bool
}, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "requests", arg0)

	outstruct := new(struct {
		Requester        common.Address
		NumWords         uint32
		CallbackGasLimit *big.Int
		Fulfilled        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Requester = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.NumWords = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(address requester, uint32 numWords, uint256 callbackGasLimit, bool fulfilled)
func (_VRFCoordinator *VRFCoordinatorSession) Requests(arg0 *big.Int) (struct {
	Requester        common.Address
	NumWords         uint32
	CallbackGasLimit *big.Int
	Fulfilled        bool
}, error) {
	return _VRFCoordinator.Contract.Requests(&_VRFCoordinator.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(address requester, uint32 numWords, uint256 callbackGasLimit, bool fulfilled)
func (_VRFCoordinator *VRFCoordinatorCallerSession) Requests(arg0 *big.Int) (struct {
	Requester        common.Address
	NumWords         uint32
	CallbackGasLimit *big.Int
	Fulfilled        bool
}, error) {
	return _VRFCoordinator.Contract.Requests(&_VRFCoordinator.CallOpts, arg0)
}

// VrfPredeploy is a free data retrieval call binding the contract method 0x41293ea2.
//
// Solidity: function vrfPredeploy() view returns(address)
func (_VRFCoordinator *VRFCoordinatorCaller) VrfPredeploy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "vrfPredeploy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VrfPredeploy is a free data retrieval call binding the contract method 0x41293ea2.
//
// Solidity: function vrfPredeploy() view returns(address)
func (_VRFCoordinator *VRFCoordinatorSession) VrfPredeploy() (common.Address, error) {
	return _VRFCoordinator.Contract.VrfPredeploy(&_VRFCoordinator.CallOpts)
}

// VrfPredeploy is a free data retrieval call binding the contract method 0x41293ea2.
//
// Solidity: function vrfPredeploy() view returns(address)
func (_VRFCoordinator *VRFCoordinatorCallerSession) VrfPredeploy() (common.Address, error) {
	return _VRFCoordinator.Contract.VrfPredeploy(&_VRFCoordinator.CallOpts)
}

// DeregisterNode is a paid mutator transaction binding the contract method 0xe7658aaf.
//
// Solidity: function deregisterNode(address node) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) DeregisterNode(opts *bind.TransactOpts, node common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "deregisterNode", node)
}

// DeregisterNode is a paid mutator transaction binding the contract method 0xe7658aaf.
//
// Solidity: function deregisterNode(address node) returns()
func (_VRFCoordinator *VRFCoordinatorSession) DeregisterNode(node common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.DeregisterNode(&_VRFCoordinator.TransactOpts, node)
}

// DeregisterNode is a paid mutator transaction binding the contract method 0xe7658aaf.
//
// Solidity: function deregisterNode(address node) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) DeregisterNode(node common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.DeregisterNode(&_VRFCoordinator.TransactOpts, node)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0x38ba4614.
//
// Solidity: function fulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) FulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "fulfillRandomWords", requestId, randomWords)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0x38ba4614.
//
// Solidity: function fulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFCoordinator *VRFCoordinatorSession) FulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomWords(&_VRFCoordinator.TransactOpts, requestId, randomWords)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0x38ba4614.
//
// Solidity: function fulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) FulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomWords(&_VRFCoordinator.TransactOpts, requestId, randomWords)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) Initialize(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "initialize", _admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_VRFCoordinator *VRFCoordinatorSession) Initialize(_admin common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Initialize(&_VRFCoordinator.TransactOpts, _admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) Initialize(_admin common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Initialize(&_VRFCoordinator.TransactOpts, _admin)
}

// RegisterNode is a paid mutator transaction binding the contract method 0x672d7a0d.
//
// Solidity: function registerNode(address node) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) RegisterNode(opts *bind.TransactOpts, node common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "registerNode", node)
}

// RegisterNode is a paid mutator transaction binding the contract method 0x672d7a0d.
//
// Solidity: function registerNode(address node) returns()
func (_VRFCoordinator *VRFCoordinatorSession) RegisterNode(node common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterNode(&_VRFCoordinator.TransactOpts, node)
}

// RegisterNode is a paid mutator transaction binding the contract method 0x672d7a0d.
//
// Solidity: function registerNode(address node) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) RegisterNode(node common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterNode(&_VRFCoordinator.TransactOpts, node)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x83f99538.
//
// Solidity: function requestRandomWords(address requester, uint32 numWords, uint256 callbackGasLimit) returns(uint256 requestId)
func (_VRFCoordinator *VRFCoordinatorTransactor) RequestRandomWords(opts *bind.TransactOpts, requester common.Address, numWords uint32, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestRandomWords", requester, numWords, callbackGasLimit)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x83f99538.
//
// Solidity: function requestRandomWords(address requester, uint32 numWords, uint256 callbackGasLimit) returns(uint256 requestId)
func (_VRFCoordinator *VRFCoordinatorSession) RequestRandomWords(requester common.Address, numWords uint32, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomWords(&_VRFCoordinator.TransactOpts, requester, numWords, callbackGasLimit)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x83f99538.
//
// Solidity: function requestRandomWords(address requester, uint32 numWords, uint256 callbackGasLimit) returns(uint256 requestId)
func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestRandomWords(requester common.Address, numWords uint32, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomWords(&_VRFCoordinator.TransactOpts, requester, numWords, callbackGasLimit)
}

// SetPredeploy is a paid mutator transaction binding the contract method 0x239c77c0.
//
// Solidity: function setPredeploy(address _vrfPredeploy) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) SetPredeploy(opts *bind.TransactOpts, _vrfPredeploy common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setPredeploy", _vrfPredeploy)
}

// SetPredeploy is a paid mutator transaction binding the contract method 0x239c77c0.
//
// Solidity: function setPredeploy(address _vrfPredeploy) returns()
func (_VRFCoordinator *VRFCoordinatorSession) SetPredeploy(_vrfPredeploy common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetPredeploy(&_VRFCoordinator.TransactOpts, _vrfPredeploy)
}

// SetPredeploy is a paid mutator transaction binding the contract method 0x239c77c0.
//
// Solidity: function setPredeploy(address _vrfPredeploy) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetPredeploy(_vrfPredeploy common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetPredeploy(&_VRFCoordinator.TransactOpts, _vrfPredeploy)
}

// VRFCoordinatorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the VRFCoordinator contract.
type VRFCoordinatorInitializedIterator struct {
	Event *VRFCoordinatorInitialized // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorInitialized)
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
		it.Event = new(VRFCoordinatorInitialized)
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
func (it *VRFCoordinatorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorInitialized represents a Initialized event raised by the VRFCoordinator contract.
type VRFCoordinatorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterInitialized(opts *bind.FilterOpts) (*VRFCoordinatorInitializedIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorInitializedIterator{contract: _VRFCoordinator.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorInitialized) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorInitialized)
				if err := _VRFCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseInitialized(log types.Log) (*VRFCoordinatorInitialized, error) {
	event := new(VRFCoordinatorInitialized)
	if err := _VRFCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorNodeDeregisteredIterator is returned from FilterNodeDeregistered and is used to iterate over the raw logs and unpacked data for NodeDeregistered events raised by the VRFCoordinator contract.
type VRFCoordinatorNodeDeregisteredIterator struct {
	Event *VRFCoordinatorNodeDeregistered // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorNodeDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorNodeDeregistered)
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
		it.Event = new(VRFCoordinatorNodeDeregistered)
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
func (it *VRFCoordinatorNodeDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorNodeDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorNodeDeregistered represents a NodeDeregistered event raised by the VRFCoordinator contract.
type VRFCoordinatorNodeDeregistered struct {
	Node common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNodeDeregistered is a free log retrieval operation binding the contract event 0xf80f07c03e488ef2bd8323b9f1b196af964916a2fca042b50c4f1a5180d09508.
//
// Solidity: event NodeDeregistered(address indexed node)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterNodeDeregistered(opts *bind.FilterOpts, node []common.Address) (*VRFCoordinatorNodeDeregisteredIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "NodeDeregistered", nodeRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorNodeDeregisteredIterator{contract: _VRFCoordinator.contract, event: "NodeDeregistered", logs: logs, sub: sub}, nil
}

// WatchNodeDeregistered is a free log subscription operation binding the contract event 0xf80f07c03e488ef2bd8323b9f1b196af964916a2fca042b50c4f1a5180d09508.
//
// Solidity: event NodeDeregistered(address indexed node)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchNodeDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorNodeDeregistered, node []common.Address) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "NodeDeregistered", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorNodeDeregistered)
				if err := _VRFCoordinator.contract.UnpackLog(event, "NodeDeregistered", log); err != nil {
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

// ParseNodeDeregistered is a log parse operation binding the contract event 0xf80f07c03e488ef2bd8323b9f1b196af964916a2fca042b50c4f1a5180d09508.
//
// Solidity: event NodeDeregistered(address indexed node)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseNodeDeregistered(log types.Log) (*VRFCoordinatorNodeDeregistered, error) {
	event := new(VRFCoordinatorNodeDeregistered)
	if err := _VRFCoordinator.contract.UnpackLog(event, "NodeDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorNodeRegisteredIterator is returned from FilterNodeRegistered and is used to iterate over the raw logs and unpacked data for NodeRegistered events raised by the VRFCoordinator contract.
type VRFCoordinatorNodeRegisteredIterator struct {
	Event *VRFCoordinatorNodeRegistered // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorNodeRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorNodeRegistered)
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
		it.Event = new(VRFCoordinatorNodeRegistered)
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
func (it *VRFCoordinatorNodeRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorNodeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorNodeRegistered represents a NodeRegistered event raised by the VRFCoordinator contract.
type VRFCoordinatorNodeRegistered struct {
	Node common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNodeRegistered is a free log retrieval operation binding the contract event 0x564728e6a7c8edd446557d94e0339d5e6ca2e05f42188914efdbdc87bcbbabf6.
//
// Solidity: event NodeRegistered(address indexed node)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterNodeRegistered(opts *bind.FilterOpts, node []common.Address) (*VRFCoordinatorNodeRegisteredIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "NodeRegistered", nodeRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorNodeRegisteredIterator{contract: _VRFCoordinator.contract, event: "NodeRegistered", logs: logs, sub: sub}, nil
}

// WatchNodeRegistered is a free log subscription operation binding the contract event 0x564728e6a7c8edd446557d94e0339d5e6ca2e05f42188914efdbdc87bcbbabf6.
//
// Solidity: event NodeRegistered(address indexed node)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchNodeRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorNodeRegistered, node []common.Address) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "NodeRegistered", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorNodeRegistered)
				if err := _VRFCoordinator.contract.UnpackLog(event, "NodeRegistered", log); err != nil {
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

// ParseNodeRegistered is a log parse operation binding the contract event 0x564728e6a7c8edd446557d94e0339d5e6ca2e05f42188914efdbdc87bcbbabf6.
//
// Solidity: event NodeRegistered(address indexed node)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseNodeRegistered(log types.Log) (*VRFCoordinatorNodeRegistered, error) {
	event := new(VRFCoordinatorNodeRegistered)
	if err := _VRFCoordinator.contract.UnpackLog(event, "NodeRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorRandomWordsFulfilledIterator is returned from FilterRandomWordsFulfilled and is used to iterate over the raw logs and unpacked data for RandomWordsFulfilled events raised by the VRFCoordinator contract.
type VRFCoordinatorRandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorRandomWordsFulfilled // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorRandomWordsFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorRandomWordsFulfilled)
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
func (it *VRFCoordinatorRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorRandomWordsFulfilled represents a RandomWordsFulfilled event raised by the VRFCoordinator contract.
type VRFCoordinatorRandomWordsFulfilled struct {
	RequestId   *big.Int
	RandomWords []*big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRandomWordsFulfilled is a free log retrieval operation binding the contract event 0xf45ee76115b0ed5f4ebe293254449fbe612bad36a53d52b87b6a40687adc48de.
//
// Solidity: event RandomWordsFulfilled(uint256 indexed requestId, uint256[] randomWords)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomWordsFulfilledIterator{contract: _VRFCoordinator.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

// WatchRandomWordsFulfilled is a free log subscription operation binding the contract event 0xf45ee76115b0ed5f4ebe293254449fbe612bad36a53d52b87b6a40687adc48de.
//
// Solidity: event RandomWordsFulfilled(uint256 indexed requestId, uint256[] randomWords)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorRandomWordsFulfilled)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

// ParseRandomWordsFulfilled is a log parse operation binding the contract event 0xf45ee76115b0ed5f4ebe293254449fbe612bad36a53d52b87b6a40687adc48de.
//
// Solidity: event RandomWordsFulfilled(uint256 indexed requestId, uint256[] randomWords)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorRandomWordsFulfilled, error) {
	event := new(VRFCoordinatorRandomWordsFulfilled)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorRandomWordsRequestedIterator is returned from FilterRandomWordsRequested and is used to iterate over the raw logs and unpacked data for RandomWordsRequested events raised by the VRFCoordinator contract.
type VRFCoordinatorRandomWordsRequestedIterator struct {
	Event *VRFCoordinatorRandomWordsRequested // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorRandomWordsRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomWordsRequested)
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
		it.Event = new(VRFCoordinatorRandomWordsRequested)
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
func (it *VRFCoordinatorRandomWordsRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorRandomWordsRequested represents a RandomWordsRequested event raised by the VRFCoordinator contract.
type VRFCoordinatorRandomWordsRequested struct {
	RequestId *big.Int
	Requester common.Address
	NumWords  uint32
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRandomWordsRequested is a free log retrieval operation binding the contract event 0xb1ef6f09514dece37223c1ec7728e9c8c30e319e98425b5a2570f5123587377c.
//
// Solidity: event RandomWordsRequested(uint256 indexed requestId, address indexed requester, uint32 numWords)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, requestId []*big.Int, requester []common.Address) (*VRFCoordinatorRandomWordsRequestedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomWordsRequested", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomWordsRequestedIterator{contract: _VRFCoordinator.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

// WatchRandomWordsRequested is a free log subscription operation binding the contract event 0xb1ef6f09514dece37223c1ec7728e9c8c30e319e98425b5a2570f5123587377c.
//
// Solidity: event RandomWordsRequested(uint256 indexed requestId, address indexed requester, uint32 numWords)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomWordsRequested, requestId []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomWordsRequested", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorRandomWordsRequested)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

// ParseRandomWordsRequested is a log parse operation binding the contract event 0xb1ef6f09514dece37223c1ec7728e9c8c30e319e98425b5a2570f5123587377c.
//
// Solidity: event RandomWordsRequested(uint256 indexed requestId, address indexed requester, uint32 numWords)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorRandomWordsRequested, error) {
	event := new(VRFCoordinatorRandomWordsRequested)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
