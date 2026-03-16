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

// VRFPredeployMetaData contains all meta data concerning the VRFPredeploy contract.
var VRFPredeployMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"coordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequestStatus\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"fulfilled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"randomWords\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_coordinator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestRandomWords\",\"inputs\":[{\"name\":\"numWords\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callbackGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RandomWordsRequested\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"numWords\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"callbackGasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// VRFPredeployABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFPredeployMetaData.ABI instead.
var VRFPredeployABI = VRFPredeployMetaData.ABI

// VRFPredeploy is an auto generated Go binding around an Ethereum contract.
type VRFPredeploy struct {
	VRFPredeployCaller     // Read-only binding to the contract
	VRFPredeployTransactor // Write-only binding to the contract
	VRFPredeployFilterer   // Log filterer for contract events
}

// VRFPredeployCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFPredeployCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFPredeployTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFPredeployTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFPredeployFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFPredeployFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFPredeploySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFPredeploySession struct {
	Contract     *VRFPredeploy     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFPredeployCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFPredeployCallerSession struct {
	Contract *VRFPredeployCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// VRFPredeployTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFPredeployTransactorSession struct {
	Contract     *VRFPredeployTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// VRFPredeployRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFPredeployRaw struct {
	Contract *VRFPredeploy // Generic contract binding to access the raw methods on
}

// VRFPredeployCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFPredeployCallerRaw struct {
	Contract *VRFPredeployCaller // Generic read-only contract binding to access the raw methods on
}

// VRFPredeployTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFPredeployTransactorRaw struct {
	Contract *VRFPredeployTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFPredeploy creates a new instance of VRFPredeploy, bound to a specific deployed contract.
func NewVRFPredeploy(address common.Address, backend bind.ContractBackend) (*VRFPredeploy, error) {
	contract, err := bindVRFPredeploy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFPredeploy{VRFPredeployCaller: VRFPredeployCaller{contract: contract}, VRFPredeployTransactor: VRFPredeployTransactor{contract: contract}, VRFPredeployFilterer: VRFPredeployFilterer{contract: contract}}, nil
}

// NewVRFPredeployCaller creates a new read-only instance of VRFPredeploy, bound to a specific deployed contract.
func NewVRFPredeployCaller(address common.Address, caller bind.ContractCaller) (*VRFPredeployCaller, error) {
	contract, err := bindVRFPredeploy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFPredeployCaller{contract: contract}, nil
}

// NewVRFPredeployTransactor creates a new write-only instance of VRFPredeploy, bound to a specific deployed contract.
func NewVRFPredeployTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFPredeployTransactor, error) {
	contract, err := bindVRFPredeploy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFPredeployTransactor{contract: contract}, nil
}

// NewVRFPredeployFilterer creates a new log filterer instance of VRFPredeploy, bound to a specific deployed contract.
func NewVRFPredeployFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFPredeployFilterer, error) {
	contract, err := bindVRFPredeploy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFPredeployFilterer{contract: contract}, nil
}

// bindVRFPredeploy binds a generic wrapper to an already deployed contract.
func bindVRFPredeploy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFPredeployABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFPredeploy *VRFPredeployRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFPredeploy.Contract.VRFPredeployCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFPredeploy *VRFPredeployRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.VRFPredeployTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFPredeploy *VRFPredeployRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.VRFPredeployTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFPredeploy *VRFPredeployCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFPredeploy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFPredeploy *VRFPredeployTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFPredeploy *VRFPredeployTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.contract.Transact(opts, method, params...)
}

// Coordinator is a free data retrieval call binding the contract method 0x0a009097.
//
// Solidity: function coordinator() view returns(address)
func (_VRFPredeploy *VRFPredeployCaller) Coordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFPredeploy.contract.Call(opts, &out, "coordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Coordinator is a free data retrieval call binding the contract method 0x0a009097.
//
// Solidity: function coordinator() view returns(address)
func (_VRFPredeploy *VRFPredeploySession) Coordinator() (common.Address, error) {
	return _VRFPredeploy.Contract.Coordinator(&_VRFPredeploy.CallOpts)
}

// Coordinator is a free data retrieval call binding the contract method 0x0a009097.
//
// Solidity: function coordinator() view returns(address)
func (_VRFPredeploy *VRFPredeployCallerSession) Coordinator() (common.Address, error) {
	return _VRFPredeploy.Contract.Coordinator(&_VRFPredeploy.CallOpts)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFPredeploy *VRFPredeployCaller) GetRequestStatus(opts *bind.CallOpts, requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	var out []interface{}
	err := _VRFPredeploy.contract.Call(opts, &out, "getRequestStatus", requestId)

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
func (_VRFPredeploy *VRFPredeploySession) GetRequestStatus(requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _VRFPredeploy.Contract.GetRequestStatus(&_VRFPredeploy.CallOpts, requestId)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFPredeploy *VRFPredeployCallerSession) GetRequestStatus(requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _VRFPredeploy.Contract.GetRequestStatus(&_VRFPredeploy.CallOpts, requestId)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _coordinator) returns()
func (_VRFPredeploy *VRFPredeployTransactor) Initialize(opts *bind.TransactOpts, _coordinator common.Address) (*types.Transaction, error) {
	return _VRFPredeploy.contract.Transact(opts, "initialize", _coordinator)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _coordinator) returns()
func (_VRFPredeploy *VRFPredeploySession) Initialize(_coordinator common.Address) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.Initialize(&_VRFPredeploy.TransactOpts, _coordinator)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _coordinator) returns()
func (_VRFPredeploy *VRFPredeployTransactorSession) Initialize(_coordinator common.Address) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.Initialize(&_VRFPredeploy.TransactOpts, _coordinator)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0xb561b1bf.
//
// Solidity: function requestRandomWords(uint32 numWords, uint256 callbackGasLimit) returns(uint256 requestId)
func (_VRFPredeploy *VRFPredeployTransactor) RequestRandomWords(opts *bind.TransactOpts, numWords uint32, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _VRFPredeploy.contract.Transact(opts, "requestRandomWords", numWords, callbackGasLimit)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0xb561b1bf.
//
// Solidity: function requestRandomWords(uint32 numWords, uint256 callbackGasLimit) returns(uint256 requestId)
func (_VRFPredeploy *VRFPredeploySession) RequestRandomWords(numWords uint32, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.RequestRandomWords(&_VRFPredeploy.TransactOpts, numWords, callbackGasLimit)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0xb561b1bf.
//
// Solidity: function requestRandomWords(uint32 numWords, uint256 callbackGasLimit) returns(uint256 requestId)
func (_VRFPredeploy *VRFPredeployTransactorSession) RequestRandomWords(numWords uint32, callbackGasLimit *big.Int) (*types.Transaction, error) {
	return _VRFPredeploy.Contract.RequestRandomWords(&_VRFPredeploy.TransactOpts, numWords, callbackGasLimit)
}

// VRFPredeployInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the VRFPredeploy contract.
type VRFPredeployInitializedIterator struct {
	Event *VRFPredeployInitialized // Event containing the contract specifics and raw log

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
func (it *VRFPredeployInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFPredeployInitialized)
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
		it.Event = new(VRFPredeployInitialized)
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
func (it *VRFPredeployInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFPredeployInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFPredeployInitialized represents a Initialized event raised by the VRFPredeploy contract.
type VRFPredeployInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_VRFPredeploy *VRFPredeployFilterer) FilterInitialized(opts *bind.FilterOpts) (*VRFPredeployInitializedIterator, error) {

	logs, sub, err := _VRFPredeploy.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VRFPredeployInitializedIterator{contract: _VRFPredeploy.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_VRFPredeploy *VRFPredeployFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFPredeployInitialized) (event.Subscription, error) {

	logs, sub, err := _VRFPredeploy.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFPredeployInitialized)
				if err := _VRFPredeploy.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_VRFPredeploy *VRFPredeployFilterer) ParseInitialized(log types.Log) (*VRFPredeployInitialized, error) {
	event := new(VRFPredeployInitialized)
	if err := _VRFPredeploy.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFPredeployRandomWordsRequestedIterator is returned from FilterRandomWordsRequested and is used to iterate over the raw logs and unpacked data for RandomWordsRequested events raised by the VRFPredeploy contract.
type VRFPredeployRandomWordsRequestedIterator struct {
	Event *VRFPredeployRandomWordsRequested // Event containing the contract specifics and raw log

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
func (it *VRFPredeployRandomWordsRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFPredeployRandomWordsRequested)
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
		it.Event = new(VRFPredeployRandomWordsRequested)
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
func (it *VRFPredeployRandomWordsRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFPredeployRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFPredeployRandomWordsRequested represents a RandomWordsRequested event raised by the VRFPredeploy contract.
type VRFPredeployRandomWordsRequested struct {
	RequestId        *big.Int
	Requester        common.Address
	NumWords         uint32
	CallbackGasLimit *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRandomWordsRequested is a free log retrieval operation binding the contract event 0x713c193000a808cf126f88fc26d1d408aa7f5ab384b408c7679c71ea222f7954.
//
// Solidity: event RandomWordsRequested(uint256 indexed requestId, address indexed requester, uint32 numWords, uint256 callbackGasLimit)
func (_VRFPredeploy *VRFPredeployFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, requestId []*big.Int, requester []common.Address) (*VRFPredeployRandomWordsRequestedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFPredeploy.contract.FilterLogs(opts, "RandomWordsRequested", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFPredeployRandomWordsRequestedIterator{contract: _VRFPredeploy.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

// WatchRandomWordsRequested is a free log subscription operation binding the contract event 0x713c193000a808cf126f88fc26d1d408aa7f5ab384b408c7679c71ea222f7954.
//
// Solidity: event RandomWordsRequested(uint256 indexed requestId, address indexed requester, uint32 numWords, uint256 callbackGasLimit)
func (_VRFPredeploy *VRFPredeployFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFPredeployRandomWordsRequested, requestId []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFPredeploy.contract.WatchLogs(opts, "RandomWordsRequested", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFPredeployRandomWordsRequested)
				if err := _VRFPredeploy.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

// ParseRandomWordsRequested is a log parse operation binding the contract event 0x713c193000a808cf126f88fc26d1d408aa7f5ab384b408c7679c71ea222f7954.
//
// Solidity: event RandomWordsRequested(uint256 indexed requestId, address indexed requester, uint32 numWords, uint256 callbackGasLimit)
func (_VRFPredeploy *VRFPredeployFilterer) ParseRandomWordsRequested(log types.Log) (*VRFPredeployRandomWordsRequested, error) {
	event := new(VRFPredeployRandomWordsRequested)
	if err := _VRFPredeploy.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
