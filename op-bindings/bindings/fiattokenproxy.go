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
	_ = abi.ConvertType
)

// FiatTokenProxyMetaData contains all meta data concerning the FiatTokenProxy contract.
var FiatTokenProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementationContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516108803803806108808339818101604052602081101561003357600080fd5b5051808061004081610051565b5061004a336100c3565b5050610123565b610064816100e760201b61042a1760201c565b61009f5760405162461bcd60e51b815260040180806020018281038252603b815260200180610845603b913960400191505060405180910390fd5b7f7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c355565b7f10d6a54a4754c8869d6886b5f5d7fbfa5b4522237ea5c60d11bc4e7a1ff9390b55565b6000813f7fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47081811480159061011b57508115155b949350505050565b610713806101326000396000f3fe60806040526004361061005a5760003560e01c80635c60da1b116100435780635c60da1b146101315780638f2839701461016f578063f851a440146101af5761005a565b80633659cfe6146100645780634f1ef286146100a4575b6100626101c4565b005b34801561007057600080fd5b506100626004803603602081101561008757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166101de565b610062600480360360408110156100ba57600080fd5b73ffffffffffffffffffffffffffffffffffffffff82351691908101906040810160208201356401000000008111156100f257600080fd5b82018360208201111561010457600080fd5b8035906020019184600183028401116401000000008311171561012657600080fd5b509092509050610232565b34801561013d57600080fd5b50610146610309565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b34801561017b57600080fd5b506100626004803603602081101561019257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610318565b3480156101bb57600080fd5b50610146610420565b6101cc610466565b6101dc6101d76104fa565b61051f565b565b6101e6610543565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156102275761022281610568565b61022f565b61022f6101c4565b50565b61023a610543565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156102fc5761027683610568565b60003073ffffffffffffffffffffffffffffffffffffffff16348484604051808383808284376040519201945060009350909150508083038185875af1925050503d80600081146102e3576040519150601f19603f3d011682016040523d82523d6000602084013e6102e8565b606091505b50509050806102f657600080fd5b50610304565b6103046101c4565b505050565b60006103136104fa565b905090565b610320610543565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156102275773ffffffffffffffffffffffffffffffffffffffff81166103bf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260368152602001806106966036913960400191505060405180910390fd5b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6103e8610543565b6040805173ffffffffffffffffffffffffffffffffffffffff928316815291841660208301528051918290030190a1610222816105bd565b6000610313610543565b6000813f7fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47081811480159061045e57508115155b949350505050565b61046e610543565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156104f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260328152602001806106646032913960400191505060405180910390fd5b6101dc6101dc565b7f7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c35490565b3660008037600080366000845af43d6000803e80801561053e573d6000f35b3d6000fd5b7f10d6a54a4754c8869d6886b5f5d7fbfa5b4522237ea5c60d11bc4e7a1ff9390b5490565b610571816105e1565b6040805173ffffffffffffffffffffffffffffffffffffffff8316815290517fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b9181900360200190a150565b7f10d6a54a4754c8869d6886b5f5d7fbfa5b4522237ea5c60d11bc4e7a1ff9390b55565b6105ea8161042a565b61063f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603b8152602001806106cc603b913960400191505060405180910390fd5b7f7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c35556fe43616e6e6f742063616c6c2066616c6c6261636b2066756e6374696f6e2066726f6d207468652070726f78792061646d696e43616e6e6f74206368616e6765207468652061646d696e206f6620612070726f787920746f20746865207a65726f206164647265737343616e6e6f742073657420612070726f787920696d706c656d656e746174696f6e20746f2061206e6f6e2d636f6e74726163742061646472657373a164736f6c634300060c000a43616e6e6f742073657420612070726f787920696d706c656d656e746174696f6e20746f2061206e6f6e2d636f6e74726163742061646472657373",
}

// FiatTokenProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use FiatTokenProxyMetaData.ABI instead.
var FiatTokenProxyABI = FiatTokenProxyMetaData.ABI

// FiatTokenProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FiatTokenProxyMetaData.Bin instead.
var FiatTokenProxyBin = FiatTokenProxyMetaData.Bin

// DeployFiatTokenProxy deploys a new Ethereum contract, binding an instance of FiatTokenProxy to it.
func DeployFiatTokenProxy(auth *bind.TransactOpts, backend bind.ContractBackend, implementationContract common.Address) (common.Address, *types.Transaction, *FiatTokenProxy, error) {
	parsed, err := FiatTokenProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FiatTokenProxyBin), backend, implementationContract)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FiatTokenProxy{FiatTokenProxyCaller: FiatTokenProxyCaller{contract: contract}, FiatTokenProxyTransactor: FiatTokenProxyTransactor{contract: contract}, FiatTokenProxyFilterer: FiatTokenProxyFilterer{contract: contract}}, nil
}

// FiatTokenProxy is an auto generated Go binding around an Ethereum contract.
type FiatTokenProxy struct {
	FiatTokenProxyCaller     // Read-only binding to the contract
	FiatTokenProxyTransactor // Write-only binding to the contract
	FiatTokenProxyFilterer   // Log filterer for contract events
}

// FiatTokenProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type FiatTokenProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatTokenProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FiatTokenProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatTokenProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FiatTokenProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatTokenProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FiatTokenProxySession struct {
	Contract     *FiatTokenProxy   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FiatTokenProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FiatTokenProxyCallerSession struct {
	Contract *FiatTokenProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// FiatTokenProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FiatTokenProxyTransactorSession struct {
	Contract     *FiatTokenProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// FiatTokenProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type FiatTokenProxyRaw struct {
	Contract *FiatTokenProxy // Generic contract binding to access the raw methods on
}

// FiatTokenProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FiatTokenProxyCallerRaw struct {
	Contract *FiatTokenProxyCaller // Generic read-only contract binding to access the raw methods on
}

// FiatTokenProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FiatTokenProxyTransactorRaw struct {
	Contract *FiatTokenProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFiatTokenProxy creates a new instance of FiatTokenProxy, bound to a specific deployed contract.
func NewFiatTokenProxy(address common.Address, backend bind.ContractBackend) (*FiatTokenProxy, error) {
	contract, err := bindFiatTokenProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FiatTokenProxy{FiatTokenProxyCaller: FiatTokenProxyCaller{contract: contract}, FiatTokenProxyTransactor: FiatTokenProxyTransactor{contract: contract}, FiatTokenProxyFilterer: FiatTokenProxyFilterer{contract: contract}}, nil
}

// NewFiatTokenProxyCaller creates a new read-only instance of FiatTokenProxy, bound to a specific deployed contract.
func NewFiatTokenProxyCaller(address common.Address, caller bind.ContractCaller) (*FiatTokenProxyCaller, error) {
	contract, err := bindFiatTokenProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FiatTokenProxyCaller{contract: contract}, nil
}

// NewFiatTokenProxyTransactor creates a new write-only instance of FiatTokenProxy, bound to a specific deployed contract.
func NewFiatTokenProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*FiatTokenProxyTransactor, error) {
	contract, err := bindFiatTokenProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FiatTokenProxyTransactor{contract: contract}, nil
}

// NewFiatTokenProxyFilterer creates a new log filterer instance of FiatTokenProxy, bound to a specific deployed contract.
func NewFiatTokenProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*FiatTokenProxyFilterer, error) {
	contract, err := bindFiatTokenProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FiatTokenProxyFilterer{contract: contract}, nil
}

// bindFiatTokenProxy binds a generic wrapper to an already deployed contract.
func bindFiatTokenProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FiatTokenProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FiatTokenProxy *FiatTokenProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FiatTokenProxy.Contract.FiatTokenProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FiatTokenProxy *FiatTokenProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.FiatTokenProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FiatTokenProxy *FiatTokenProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.FiatTokenProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FiatTokenProxy *FiatTokenProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FiatTokenProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FiatTokenProxy *FiatTokenProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FiatTokenProxy *FiatTokenProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_FiatTokenProxy *FiatTokenProxyCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatTokenProxy.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_FiatTokenProxy *FiatTokenProxySession) Admin() (common.Address, error) {
	return _FiatTokenProxy.Contract.Admin(&_FiatTokenProxy.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_FiatTokenProxy *FiatTokenProxyCallerSession) Admin() (common.Address, error) {
	return _FiatTokenProxy.Contract.Admin(&_FiatTokenProxy.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_FiatTokenProxy *FiatTokenProxyCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatTokenProxy.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_FiatTokenProxy *FiatTokenProxySession) Implementation() (common.Address, error) {
	return _FiatTokenProxy.Contract.Implementation(&_FiatTokenProxy.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_FiatTokenProxy *FiatTokenProxyCallerSession) Implementation() (common.Address, error) {
	return _FiatTokenProxy.Contract.Implementation(&_FiatTokenProxy.CallOpts)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_FiatTokenProxy *FiatTokenProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _FiatTokenProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_FiatTokenProxy *FiatTokenProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.ChangeAdmin(&_FiatTokenProxy.TransactOpts, newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_FiatTokenProxy *FiatTokenProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.ChangeAdmin(&_FiatTokenProxy.TransactOpts, newAdmin)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FiatTokenProxy *FiatTokenProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _FiatTokenProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FiatTokenProxy *FiatTokenProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.UpgradeTo(&_FiatTokenProxy.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FiatTokenProxy *FiatTokenProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.UpgradeTo(&_FiatTokenProxy.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FiatTokenProxy *FiatTokenProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FiatTokenProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FiatTokenProxy *FiatTokenProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.UpgradeToAndCall(&_FiatTokenProxy.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FiatTokenProxy *FiatTokenProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.UpgradeToAndCall(&_FiatTokenProxy.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_FiatTokenProxy *FiatTokenProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _FiatTokenProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_FiatTokenProxy *FiatTokenProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.Fallback(&_FiatTokenProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_FiatTokenProxy *FiatTokenProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _FiatTokenProxy.Contract.Fallback(&_FiatTokenProxy.TransactOpts, calldata)
}

// FiatTokenProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the FiatTokenProxy contract.
type FiatTokenProxyAdminChangedIterator struct {
	Event *FiatTokenProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *FiatTokenProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenProxyAdminChanged)
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
		it.Event = new(FiatTokenProxyAdminChanged)
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
func (it *FiatTokenProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenProxyAdminChanged represents a AdminChanged event raised by the FiatTokenProxy contract.
type FiatTokenProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FiatTokenProxy *FiatTokenProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*FiatTokenProxyAdminChangedIterator, error) {

	logs, sub, err := _FiatTokenProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &FiatTokenProxyAdminChangedIterator{contract: _FiatTokenProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FiatTokenProxy *FiatTokenProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *FiatTokenProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _FiatTokenProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenProxyAdminChanged)
				if err := _FiatTokenProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FiatTokenProxy *FiatTokenProxyFilterer) ParseAdminChanged(log types.Log) (*FiatTokenProxyAdminChanged, error) {
	event := new(FiatTokenProxyAdminChanged)
	if err := _FiatTokenProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FiatTokenProxy contract.
type FiatTokenProxyUpgradedIterator struct {
	Event *FiatTokenProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *FiatTokenProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenProxyUpgraded)
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
		it.Event = new(FiatTokenProxyUpgraded)
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
func (it *FiatTokenProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenProxyUpgraded represents a Upgraded event raised by the FiatTokenProxy contract.
type FiatTokenProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address implementation)
func (_FiatTokenProxy *FiatTokenProxyFilterer) FilterUpgraded(opts *bind.FilterOpts) (*FiatTokenProxyUpgradedIterator, error) {

	logs, sub, err := _FiatTokenProxy.contract.FilterLogs(opts, "Upgraded")
	if err != nil {
		return nil, err
	}
	return &FiatTokenProxyUpgradedIterator{contract: _FiatTokenProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address implementation)
func (_FiatTokenProxy *FiatTokenProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FiatTokenProxyUpgraded) (event.Subscription, error) {

	logs, sub, err := _FiatTokenProxy.contract.WatchLogs(opts, "Upgraded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenProxyUpgraded)
				if err := _FiatTokenProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address implementation)
func (_FiatTokenProxy *FiatTokenProxyFilterer) ParseUpgraded(log types.Log) (*FiatTokenProxyUpgraded, error) {
	event := new(FiatTokenProxyUpgraded)
	if err := _FiatTokenProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
