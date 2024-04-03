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

// DisputeGameFactoryMetaData contains all meta data concerning the DisputeGameFactory contract.
var DisputeGameFactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"disputeProxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"GameType\",\"name\":\"gameType\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"Claim\",\"name\":\"rootClaim\",\"type\":\"bytes32\"}],\"name\":\"DisputeGameCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"impl\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"GameType\",\"name\":\"gameType\",\"type\":\"uint8\"}],\"name\":\"ImplementationSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"GameType\",\"name\":\"_gameType\",\"type\":\"uint8\"},{\"internalType\":\"Claim\",\"name\":\"_rootClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"create\",\"outputs\":[{\"internalType\":\"contractIDisputeGame\",\"name\":\"proxy_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"gameAtIndex\",\"outputs\":[{\"internalType\":\"GameType\",\"name\":\"gameType_\",\"type\":\"uint8\"},{\"internalType\":\"Timestamp\",\"name\":\"timestamp_\",\"type\":\"uint64\"},{\"internalType\":\"contractIDisputeGame\",\"name\":\"proxy_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gameCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"gameCount_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"GameType\",\"name\":\"_gameType\",\"type\":\"uint8\"}],\"name\":\"gameImpls\",\"outputs\":[{\"internalType\":\"contractIDisputeGame\",\"name\":\"impl_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"GameType\",\"name\":\"_gameType\",\"type\":\"uint8\"},{\"internalType\":\"Claim\",\"name\":\"_rootClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"games\",\"outputs\":[{\"internalType\":\"contractIDisputeGame\",\"name\":\"proxy_\",\"type\":\"address\"},{\"internalType\":\"Timestamp\",\"name\":\"timestamp_\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"GameType\",\"name\":\"_gameType\",\"type\":\"uint8\"},{\"internalType\":\"Claim\",\"name\":\"_rootClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"getGameUUID\",\"outputs\":[{\"internalType\":\"Hash\",\"name\":\"uuid_\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"GameType\",\"name\":\"_gameType\",\"type\":\"uint8\"},{\"internalType\":\"contractIDisputeGame\",\"name\":\"_impl\",\"type\":\"address\"}],\"name\":\"setImplementation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DisputeGameFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use DisputeGameFactoryMetaData.ABI instead.
var DisputeGameFactoryABI = DisputeGameFactoryMetaData.ABI

// DisputeGameFactory is an auto generated Go binding around an Ethereum contract.
type DisputeGameFactory struct {
	DisputeGameFactoryCaller     // Read-only binding to the contract
	DisputeGameFactoryTransactor // Write-only binding to the contract
	DisputeGameFactoryFilterer   // Log filterer for contract events
}

// DisputeGameFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DisputeGameFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisputeGameFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DisputeGameFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisputeGameFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DisputeGameFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisputeGameFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DisputeGameFactorySession struct {
	Contract     *DisputeGameFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DisputeGameFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DisputeGameFactoryCallerSession struct {
	Contract *DisputeGameFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DisputeGameFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DisputeGameFactoryTransactorSession struct {
	Contract     *DisputeGameFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DisputeGameFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DisputeGameFactoryRaw struct {
	Contract *DisputeGameFactory // Generic contract binding to access the raw methods on
}

// DisputeGameFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DisputeGameFactoryCallerRaw struct {
	Contract *DisputeGameFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// DisputeGameFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DisputeGameFactoryTransactorRaw struct {
	Contract *DisputeGameFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDisputeGameFactory creates a new instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactory(address common.Address, backend bind.ContractBackend) (*DisputeGameFactory, error) {
	contract, err := bindDisputeGameFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactory{DisputeGameFactoryCaller: DisputeGameFactoryCaller{contract: contract}, DisputeGameFactoryTransactor: DisputeGameFactoryTransactor{contract: contract}, DisputeGameFactoryFilterer: DisputeGameFactoryFilterer{contract: contract}}, nil
}

// NewDisputeGameFactoryCaller creates a new read-only instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactoryCaller(address common.Address, caller bind.ContractCaller) (*DisputeGameFactoryCaller, error) {
	contract, err := bindDisputeGameFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryCaller{contract: contract}, nil
}

// NewDisputeGameFactoryTransactor creates a new write-only instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*DisputeGameFactoryTransactor, error) {
	contract, err := bindDisputeGameFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryTransactor{contract: contract}, nil
}

// NewDisputeGameFactoryFilterer creates a new log filterer instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*DisputeGameFactoryFilterer, error) {
	contract, err := bindDisputeGameFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryFilterer{contract: contract}, nil
}

// bindDisputeGameFactory binds a generic wrapper to an already deployed contract.
func bindDisputeGameFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DisputeGameFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DisputeGameFactory *DisputeGameFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DisputeGameFactory.Contract.DisputeGameFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DisputeGameFactory *DisputeGameFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.DisputeGameFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DisputeGameFactory *DisputeGameFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.DisputeGameFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DisputeGameFactory *DisputeGameFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DisputeGameFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DisputeGameFactory *DisputeGameFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DisputeGameFactory *DisputeGameFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.contract.Transact(opts, method, params...)
}

// GameAtIndex is a free data retrieval call binding the contract method 0xbb8aa1fc.
//
// Solidity: function gameAtIndex(uint256 _index) view returns(uint8 gameType_, uint64 timestamp_, address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GameAtIndex(opts *bind.CallOpts, _index *big.Int) (struct {
	GameType  uint8
	Timestamp uint64
	Proxy     common.Address
}, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "gameAtIndex", _index)

	outstruct := new(struct {
		GameType  uint8
		Timestamp uint64
		Proxy     common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameType = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Proxy = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GameAtIndex is a free data retrieval call binding the contract method 0xbb8aa1fc.
//
// Solidity: function gameAtIndex(uint256 _index) view returns(uint8 gameType_, uint64 timestamp_, address proxy_)
func (_DisputeGameFactory *DisputeGameFactorySession) GameAtIndex(_index *big.Int) (struct {
	GameType  uint8
	Timestamp uint64
	Proxy     common.Address
}, error) {
	return _DisputeGameFactory.Contract.GameAtIndex(&_DisputeGameFactory.CallOpts, _index)
}

// GameAtIndex is a free data retrieval call binding the contract method 0xbb8aa1fc.
//
// Solidity: function gameAtIndex(uint256 _index) view returns(uint8 gameType_, uint64 timestamp_, address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GameAtIndex(_index *big.Int) (struct {
	GameType  uint8
	Timestamp uint64
	Proxy     common.Address
}, error) {
	return _DisputeGameFactory.Contract.GameAtIndex(&_DisputeGameFactory.CallOpts, _index)
}

// GameCount is a free data retrieval call binding the contract method 0x4d1975b4.
//
// Solidity: function gameCount() view returns(uint256 gameCount_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GameCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "gameCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GameCount is a free data retrieval call binding the contract method 0x4d1975b4.
//
// Solidity: function gameCount() view returns(uint256 gameCount_)
func (_DisputeGameFactory *DisputeGameFactorySession) GameCount() (*big.Int, error) {
	return _DisputeGameFactory.Contract.GameCount(&_DisputeGameFactory.CallOpts)
}

// GameCount is a free data retrieval call binding the contract method 0x4d1975b4.
//
// Solidity: function gameCount() view returns(uint256 gameCount_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GameCount() (*big.Int, error) {
	return _DisputeGameFactory.Contract.GameCount(&_DisputeGameFactory.CallOpts)
}

// GameImpls is a free data retrieval call binding the contract method 0xdfa162d3.
//
// Solidity: function gameImpls(uint8 _gameType) view returns(address impl_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GameImpls(opts *bind.CallOpts, _gameType uint8) (common.Address, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "gameImpls", _gameType)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GameImpls is a free data retrieval call binding the contract method 0xdfa162d3.
//
// Solidity: function gameImpls(uint8 _gameType) view returns(address impl_)
func (_DisputeGameFactory *DisputeGameFactorySession) GameImpls(_gameType uint8) (common.Address, error) {
	return _DisputeGameFactory.Contract.GameImpls(&_DisputeGameFactory.CallOpts, _gameType)
}

// GameImpls is a free data retrieval call binding the contract method 0xdfa162d3.
//
// Solidity: function gameImpls(uint8 _gameType) view returns(address impl_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GameImpls(_gameType uint8) (common.Address, error) {
	return _DisputeGameFactory.Contract.GameImpls(&_DisputeGameFactory.CallOpts, _gameType)
}

// Games is a free data retrieval call binding the contract method 0xc49d5271.
//
// Solidity: function games(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) view returns(address proxy_, uint64 timestamp_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) Games(opts *bind.CallOpts, _gameType uint8, _rootClaim [32]byte, _extraData []byte) (struct {
	Proxy     common.Address
	Timestamp uint64
}, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "games", _gameType, _rootClaim, _extraData)

	outstruct := new(struct {
		Proxy     common.Address
		Timestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Proxy = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// Games is a free data retrieval call binding the contract method 0xc49d5271.
//
// Solidity: function games(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) view returns(address proxy_, uint64 timestamp_)
func (_DisputeGameFactory *DisputeGameFactorySession) Games(_gameType uint8, _rootClaim [32]byte, _extraData []byte) (struct {
	Proxy     common.Address
	Timestamp uint64
}, error) {
	return _DisputeGameFactory.Contract.Games(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// Games is a free data retrieval call binding the contract method 0xc49d5271.
//
// Solidity: function games(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) view returns(address proxy_, uint64 timestamp_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) Games(_gameType uint8, _rootClaim [32]byte, _extraData []byte) (struct {
	Proxy     common.Address
	Timestamp uint64
}, error) {
	return _DisputeGameFactory.Contract.Games(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// GetGameUUID is a free data retrieval call binding the contract method 0x26daafbe.
//
// Solidity: function getGameUUID(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) pure returns(bytes32 uuid_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GetGameUUID(opts *bind.CallOpts, _gameType uint8, _rootClaim [32]byte, _extraData []byte) ([32]byte, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "getGameUUID", _gameType, _rootClaim, _extraData)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetGameUUID is a free data retrieval call binding the contract method 0x26daafbe.
//
// Solidity: function getGameUUID(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) pure returns(bytes32 uuid_)
func (_DisputeGameFactory *DisputeGameFactorySession) GetGameUUID(_gameType uint8, _rootClaim [32]byte, _extraData []byte) ([32]byte, error) {
	return _DisputeGameFactory.Contract.GetGameUUID(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// GetGameUUID is a free data retrieval call binding the contract method 0x26daafbe.
//
// Solidity: function getGameUUID(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) pure returns(bytes32 uuid_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GetGameUUID(_gameType uint8, _rootClaim [32]byte, _extraData []byte) ([32]byte, error) {
	return _DisputeGameFactory.Contract.GetGameUUID(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// Create is a paid mutator transaction binding the contract method 0x3142e55e.
//
// Solidity: function create(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryTransactor) Create(opts *bind.TransactOpts, _gameType uint8, _rootClaim [32]byte, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "create", _gameType, _rootClaim, _extraData)
}

// Create is a paid mutator transaction binding the contract method 0x3142e55e.
//
// Solidity: function create(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactorySession) Create(_gameType uint8, _rootClaim [32]byte, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.Create(&_DisputeGameFactory.TransactOpts, _gameType, _rootClaim, _extraData)
}

// Create is a paid mutator transaction binding the contract method 0x3142e55e.
//
// Solidity: function create(uint8 _gameType, bytes32 _rootClaim, bytes _extraData) returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) Create(_gameType uint8, _rootClaim [32]byte, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.Create(&_DisputeGameFactory.TransactOpts, _gameType, _rootClaim, _extraData)
}

// SetImplementation is a paid mutator transaction binding the contract method 0x45583b7a.
//
// Solidity: function setImplementation(uint8 _gameType, address _impl) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactor) SetImplementation(opts *bind.TransactOpts, _gameType uint8, _impl common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "setImplementation", _gameType, _impl)
}

// SetImplementation is a paid mutator transaction binding the contract method 0x45583b7a.
//
// Solidity: function setImplementation(uint8 _gameType, address _impl) returns()
func (_DisputeGameFactory *DisputeGameFactorySession) SetImplementation(_gameType uint8, _impl common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.SetImplementation(&_DisputeGameFactory.TransactOpts, _gameType, _impl)
}

// SetImplementation is a paid mutator transaction binding the contract method 0x45583b7a.
//
// Solidity: function setImplementation(uint8 _gameType, address _impl) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) SetImplementation(_gameType uint8, _impl common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.SetImplementation(&_DisputeGameFactory.TransactOpts, _gameType, _impl)
}

// DisputeGameFactoryDisputeGameCreatedIterator is returned from FilterDisputeGameCreated and is used to iterate over the raw logs and unpacked data for DisputeGameCreated events raised by the DisputeGameFactory contract.
type DisputeGameFactoryDisputeGameCreatedIterator struct {
	Event *DisputeGameFactoryDisputeGameCreated // Event containing the contract specifics and raw log

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
func (it *DisputeGameFactoryDisputeGameCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DisputeGameFactoryDisputeGameCreated)
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
		it.Event = new(DisputeGameFactoryDisputeGameCreated)
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
func (it *DisputeGameFactoryDisputeGameCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DisputeGameFactoryDisputeGameCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DisputeGameFactoryDisputeGameCreated represents a DisputeGameCreated event raised by the DisputeGameFactory contract.
type DisputeGameFactoryDisputeGameCreated struct {
	DisputeProxy common.Address
	GameType     uint8
	RootClaim    [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDisputeGameCreated is a free log retrieval operation binding the contract event 0xfad0599ff449d8d9685eadecca8cb9e00924c5fd8367c1c09469824939e1ffec.
//
// Solidity: event DisputeGameCreated(address indexed disputeProxy, uint8 indexed gameType, bytes32 indexed rootClaim)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) FilterDisputeGameCreated(opts *bind.FilterOpts, disputeProxy []common.Address, gameType []uint8, rootClaim [][32]byte) (*DisputeGameFactoryDisputeGameCreatedIterator, error) {

	var disputeProxyRule []interface{}
	for _, disputeProxyItem := range disputeProxy {
		disputeProxyRule = append(disputeProxyRule, disputeProxyItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}
	var rootClaimRule []interface{}
	for _, rootClaimItem := range rootClaim {
		rootClaimRule = append(rootClaimRule, rootClaimItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.FilterLogs(opts, "DisputeGameCreated", disputeProxyRule, gameTypeRule, rootClaimRule)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryDisputeGameCreatedIterator{contract: _DisputeGameFactory.contract, event: "DisputeGameCreated", logs: logs, sub: sub}, nil
}

// WatchDisputeGameCreated is a free log subscription operation binding the contract event 0xfad0599ff449d8d9685eadecca8cb9e00924c5fd8367c1c09469824939e1ffec.
//
// Solidity: event DisputeGameCreated(address indexed disputeProxy, uint8 indexed gameType, bytes32 indexed rootClaim)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) WatchDisputeGameCreated(opts *bind.WatchOpts, sink chan<- *DisputeGameFactoryDisputeGameCreated, disputeProxy []common.Address, gameType []uint8, rootClaim [][32]byte) (event.Subscription, error) {

	var disputeProxyRule []interface{}
	for _, disputeProxyItem := range disputeProxy {
		disputeProxyRule = append(disputeProxyRule, disputeProxyItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}
	var rootClaimRule []interface{}
	for _, rootClaimItem := range rootClaim {
		rootClaimRule = append(rootClaimRule, rootClaimItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.WatchLogs(opts, "DisputeGameCreated", disputeProxyRule, gameTypeRule, rootClaimRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DisputeGameFactoryDisputeGameCreated)
				if err := _DisputeGameFactory.contract.UnpackLog(event, "DisputeGameCreated", log); err != nil {
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

// ParseDisputeGameCreated is a log parse operation binding the contract event 0xfad0599ff449d8d9685eadecca8cb9e00924c5fd8367c1c09469824939e1ffec.
//
// Solidity: event DisputeGameCreated(address indexed disputeProxy, uint8 indexed gameType, bytes32 indexed rootClaim)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) ParseDisputeGameCreated(log types.Log) (*DisputeGameFactoryDisputeGameCreated, error) {
	event := new(DisputeGameFactoryDisputeGameCreated)
	if err := _DisputeGameFactory.contract.UnpackLog(event, "DisputeGameCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DisputeGameFactoryImplementationSetIterator is returned from FilterImplementationSet and is used to iterate over the raw logs and unpacked data for ImplementationSet events raised by the DisputeGameFactory contract.
type DisputeGameFactoryImplementationSetIterator struct {
	Event *DisputeGameFactoryImplementationSet // Event containing the contract specifics and raw log

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
func (it *DisputeGameFactoryImplementationSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DisputeGameFactoryImplementationSet)
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
		it.Event = new(DisputeGameFactoryImplementationSet)
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
func (it *DisputeGameFactoryImplementationSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DisputeGameFactoryImplementationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DisputeGameFactoryImplementationSet represents a ImplementationSet event raised by the DisputeGameFactory contract.
type DisputeGameFactoryImplementationSet struct {
	Impl     common.Address
	GameType uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterImplementationSet is a free log retrieval operation binding the contract event 0x623713f72f6e427a8044bb8b3bd6834357cf285decbaa21bcc73c1d0632c4d84.
//
// Solidity: event ImplementationSet(address indexed impl, uint8 indexed gameType)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) FilterImplementationSet(opts *bind.FilterOpts, impl []common.Address, gameType []uint8) (*DisputeGameFactoryImplementationSetIterator, error) {

	var implRule []interface{}
	for _, implItem := range impl {
		implRule = append(implRule, implItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.FilterLogs(opts, "ImplementationSet", implRule, gameTypeRule)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryImplementationSetIterator{contract: _DisputeGameFactory.contract, event: "ImplementationSet", logs: logs, sub: sub}, nil
}

// WatchImplementationSet is a free log subscription operation binding the contract event 0x623713f72f6e427a8044bb8b3bd6834357cf285decbaa21bcc73c1d0632c4d84.
//
// Solidity: event ImplementationSet(address indexed impl, uint8 indexed gameType)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) WatchImplementationSet(opts *bind.WatchOpts, sink chan<- *DisputeGameFactoryImplementationSet, impl []common.Address, gameType []uint8) (event.Subscription, error) {

	var implRule []interface{}
	for _, implItem := range impl {
		implRule = append(implRule, implItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.WatchLogs(opts, "ImplementationSet", implRule, gameTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DisputeGameFactoryImplementationSet)
				if err := _DisputeGameFactory.contract.UnpackLog(event, "ImplementationSet", log); err != nil {
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

// ParseImplementationSet is a log parse operation binding the contract event 0x623713f72f6e427a8044bb8b3bd6834357cf285decbaa21bcc73c1d0632c4d84.
//
// Solidity: event ImplementationSet(address indexed impl, uint8 indexed gameType)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) ParseImplementationSet(log types.Log) (*DisputeGameFactoryImplementationSet, error) {
	event := new(DisputeGameFactoryImplementationSet)
	if err := _DisputeGameFactory.contract.UnpackLog(event, "ImplementationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
