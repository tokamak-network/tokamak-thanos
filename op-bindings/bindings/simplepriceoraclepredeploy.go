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

// SimplePriceOraclePredeployMetaData contains all meta data concerning the SimplePriceOraclePredeploy contract.
var SimplePriceOraclePredeployMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getPrice\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastUpdated\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePrice\",\"inputs\":[{\"name\":\"newPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceUpdated\",\"inputs\":[{\"name\":\"newPrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// SimplePriceOraclePredeployABI is the input ABI used to generate the binding from.
// Deprecated: Use SimplePriceOraclePredeployMetaData.ABI instead.
var SimplePriceOraclePredeployABI = SimplePriceOraclePredeployMetaData.ABI

// SimplePriceOraclePredeploy is an auto generated Go binding around an Ethereum contract.
type SimplePriceOraclePredeploy struct {
	SimplePriceOraclePredeployCaller     // Read-only binding to the contract
	SimplePriceOraclePredeployTransactor // Write-only binding to the contract
	SimplePriceOraclePredeployFilterer   // Log filterer for contract events
}

// SimplePriceOraclePredeployCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimplePriceOraclePredeployCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimplePriceOraclePredeployTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimplePriceOraclePredeployTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimplePriceOraclePredeployFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimplePriceOraclePredeployFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimplePriceOraclePredeploySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimplePriceOraclePredeploySession struct {
	Contract     *SimplePriceOraclePredeploy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SimplePriceOraclePredeployCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimplePriceOraclePredeployCallerSession struct {
	Contract *SimplePriceOraclePredeployCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// SimplePriceOraclePredeployTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimplePriceOraclePredeployTransactorSession struct {
	Contract     *SimplePriceOraclePredeployTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// SimplePriceOraclePredeployRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimplePriceOraclePredeployRaw struct {
	Contract *SimplePriceOraclePredeploy // Generic contract binding to access the raw methods on
}

// SimplePriceOraclePredeployCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimplePriceOraclePredeployCallerRaw struct {
	Contract *SimplePriceOraclePredeployCaller // Generic read-only contract binding to access the raw methods on
}

// SimplePriceOraclePredeployTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimplePriceOraclePredeployTransactorRaw struct {
	Contract *SimplePriceOraclePredeployTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimplePriceOraclePredeploy creates a new instance of SimplePriceOraclePredeploy, bound to a specific deployed contract.
func NewSimplePriceOraclePredeploy(address common.Address, backend bind.ContractBackend) (*SimplePriceOraclePredeploy, error) {
	contract, err := bindSimplePriceOraclePredeploy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimplePriceOraclePredeploy{SimplePriceOraclePredeployCaller: SimplePriceOraclePredeployCaller{contract: contract}, SimplePriceOraclePredeployTransactor: SimplePriceOraclePredeployTransactor{contract: contract}, SimplePriceOraclePredeployFilterer: SimplePriceOraclePredeployFilterer{contract: contract}}, nil
}

// NewSimplePriceOraclePredeployCaller creates a new read-only instance of SimplePriceOraclePredeploy, bound to a specific deployed contract.
func NewSimplePriceOraclePredeployCaller(address common.Address, caller bind.ContractCaller) (*SimplePriceOraclePredeployCaller, error) {
	contract, err := bindSimplePriceOraclePredeploy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimplePriceOraclePredeployCaller{contract: contract}, nil
}

// NewSimplePriceOraclePredeployTransactor creates a new write-only instance of SimplePriceOraclePredeploy, bound to a specific deployed contract.
func NewSimplePriceOraclePredeployTransactor(address common.Address, transactor bind.ContractTransactor) (*SimplePriceOraclePredeployTransactor, error) {
	contract, err := bindSimplePriceOraclePredeploy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimplePriceOraclePredeployTransactor{contract: contract}, nil
}

// NewSimplePriceOraclePredeployFilterer creates a new log filterer instance of SimplePriceOraclePredeploy, bound to a specific deployed contract.
func NewSimplePriceOraclePredeployFilterer(address common.Address, filterer bind.ContractFilterer) (*SimplePriceOraclePredeployFilterer, error) {
	contract, err := bindSimplePriceOraclePredeploy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimplePriceOraclePredeployFilterer{contract: contract}, nil
}

// bindSimplePriceOraclePredeploy binds a generic wrapper to an already deployed contract.
func bindSimplePriceOraclePredeploy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimplePriceOraclePredeployABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimplePriceOraclePredeploy.Contract.SimplePriceOraclePredeployCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.SimplePriceOraclePredeployTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.SimplePriceOraclePredeployTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimplePriceOraclePredeploy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.contract.Transact(opts, method, params...)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployCaller) GetPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimplePriceOraclePredeploy.contract.Call(opts, &out, "getPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeploySession) GetPrice() (*big.Int, error) {
	return _SimplePriceOraclePredeploy.Contract.GetPrice(&_SimplePriceOraclePredeploy.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployCallerSession) GetPrice() (*big.Int, error) {
	return _SimplePriceOraclePredeploy.Contract.GetPrice(&_SimplePriceOraclePredeploy.CallOpts)
}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns(uint256)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployCaller) LastUpdated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimplePriceOraclePredeploy.contract.Call(opts, &out, "lastUpdated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns(uint256)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeploySession) LastUpdated() (*big.Int, error) {
	return _SimplePriceOraclePredeploy.Contract.LastUpdated(&_SimplePriceOraclePredeploy.CallOpts)
}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns(uint256)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployCallerSession) LastUpdated() (*big.Int, error) {
	return _SimplePriceOraclePredeploy.Contract.LastUpdated(&_SimplePriceOraclePredeploy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimplePriceOraclePredeploy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeploySession) Owner() (common.Address, error) {
	return _SimplePriceOraclePredeploy.Contract.Owner(&_SimplePriceOraclePredeploy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployCallerSession) Owner() (common.Address, error) {
	return _SimplePriceOraclePredeploy.Contract.Owner(&_SimplePriceOraclePredeploy.CallOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeploySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.TransferOwnership(&_SimplePriceOraclePredeploy.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.TransferOwnership(&_SimplePriceOraclePredeploy.TransactOpts, newOwner)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x8d6cc56d.
//
// Solidity: function updatePrice(uint256 newPrice) returns()
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployTransactor) UpdatePrice(opts *bind.TransactOpts, newPrice *big.Int) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.contract.Transact(opts, "updatePrice", newPrice)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x8d6cc56d.
//
// Solidity: function updatePrice(uint256 newPrice) returns()
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeploySession) UpdatePrice(newPrice *big.Int) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.UpdatePrice(&_SimplePriceOraclePredeploy.TransactOpts, newPrice)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x8d6cc56d.
//
// Solidity: function updatePrice(uint256 newPrice) returns()
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployTransactorSession) UpdatePrice(newPrice *big.Int) (*types.Transaction, error) {
	return _SimplePriceOraclePredeploy.Contract.UpdatePrice(&_SimplePriceOraclePredeploy.TransactOpts, newPrice)
}

// SimplePriceOraclePredeployOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SimplePriceOraclePredeploy contract.
type SimplePriceOraclePredeployOwnershipTransferredIterator struct {
	Event *SimplePriceOraclePredeployOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SimplePriceOraclePredeployOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimplePriceOraclePredeployOwnershipTransferred)
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
		it.Event = new(SimplePriceOraclePredeployOwnershipTransferred)
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
func (it *SimplePriceOraclePredeployOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimplePriceOraclePredeployOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimplePriceOraclePredeployOwnershipTransferred represents a OwnershipTransferred event raised by the SimplePriceOraclePredeploy contract.
type SimplePriceOraclePredeployOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SimplePriceOraclePredeployOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimplePriceOraclePredeploy.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SimplePriceOraclePredeployOwnershipTransferredIterator{contract: _SimplePriceOraclePredeploy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimplePriceOraclePredeployOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimplePriceOraclePredeploy.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimplePriceOraclePredeployOwnershipTransferred)
				if err := _SimplePriceOraclePredeploy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployFilterer) ParseOwnershipTransferred(log types.Log) (*SimplePriceOraclePredeployOwnershipTransferred, error) {
	event := new(SimplePriceOraclePredeployOwnershipTransferred)
	if err := _SimplePriceOraclePredeploy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimplePriceOraclePredeployPriceUpdatedIterator is returned from FilterPriceUpdated and is used to iterate over the raw logs and unpacked data for PriceUpdated events raised by the SimplePriceOraclePredeploy contract.
type SimplePriceOraclePredeployPriceUpdatedIterator struct {
	Event *SimplePriceOraclePredeployPriceUpdated // Event containing the contract specifics and raw log

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
func (it *SimplePriceOraclePredeployPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimplePriceOraclePredeployPriceUpdated)
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
		it.Event = new(SimplePriceOraclePredeployPriceUpdated)
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
func (it *SimplePriceOraclePredeployPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimplePriceOraclePredeployPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimplePriceOraclePredeployPriceUpdated represents a PriceUpdated event raised by the SimplePriceOraclePredeploy contract.
type SimplePriceOraclePredeployPriceUpdated struct {
	NewPrice  *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPriceUpdated is a free log retrieval operation binding the contract event 0x945c1c4e99aa89f648fbfe3df471b916f719e16d960fcec0737d4d56bd696838.
//
// Solidity: event PriceUpdated(uint256 newPrice, uint256 timestamp)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployFilterer) FilterPriceUpdated(opts *bind.FilterOpts) (*SimplePriceOraclePredeployPriceUpdatedIterator, error) {

	logs, sub, err := _SimplePriceOraclePredeploy.contract.FilterLogs(opts, "PriceUpdated")
	if err != nil {
		return nil, err
	}
	return &SimplePriceOraclePredeployPriceUpdatedIterator{contract: _SimplePriceOraclePredeploy.contract, event: "PriceUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceUpdated is a free log subscription operation binding the contract event 0x945c1c4e99aa89f648fbfe3df471b916f719e16d960fcec0737d4d56bd696838.
//
// Solidity: event PriceUpdated(uint256 newPrice, uint256 timestamp)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployFilterer) WatchPriceUpdated(opts *bind.WatchOpts, sink chan<- *SimplePriceOraclePredeployPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _SimplePriceOraclePredeploy.contract.WatchLogs(opts, "PriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimplePriceOraclePredeployPriceUpdated)
				if err := _SimplePriceOraclePredeploy.contract.UnpackLog(event, "PriceUpdated", log); err != nil {
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

// ParsePriceUpdated is a log parse operation binding the contract event 0x945c1c4e99aa89f648fbfe3df471b916f719e16d960fcec0737d4d56bd696838.
//
// Solidity: event PriceUpdated(uint256 newPrice, uint256 timestamp)
func (_SimplePriceOraclePredeploy *SimplePriceOraclePredeployFilterer) ParsePriceUpdated(log types.Log) (*SimplePriceOraclePredeployPriceUpdated, error) {
	event := new(SimplePriceOraclePredeployPriceUpdated)
	if err := _SimplePriceOraclePredeploy.contract.UnpackLog(event, "PriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
