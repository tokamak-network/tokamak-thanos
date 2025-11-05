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

// ProxyContractMetaData contains all meta data concerning the ProxyContract contract.
var ProxyContractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeAdmin\",\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"_implementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"_implementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161091f38038061091f83398101604081905261002f916100b5565b6100388161003e565b506100e5565b60006100566000805160206108ff8339815191525490565b6000805160206108ff833981519152838155604080516001600160a01b0380851682528616602082015292935090917f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f910160405180910390a1505050565b6000602082840312156100c757600080fd5b81516001600160a01b03811681146100de57600080fd5b9392505050565b61080b806100f46000396000f3fe60806040526004361061005e5760003560e01c80635c60da1b116100435780635c60da1b146100be5780638f283970146100f8578063f851a440146101185761006d565b80633659cfe6146100755780634f1ef286146100955761006d565b3661006d5761006b61012d565b005b61006b61012d565b34801561008157600080fd5b5061006b6100903660046106dd565b610224565b6100a86100a33660046106f8565b610296565b6040516100b5919061077b565b60405180910390f35b3480156100ca57600080fd5b506100d3610419565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100b5565b34801561010457600080fd5b5061006b6101133660046106dd565b6104b0565b34801561012457600080fd5b506100d3610517565b60006101577f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b905073ffffffffffffffffffffffffffffffffffffffff8116610201576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f50726f78793a20696d706c656d656e746174696f6e206e6f7420696e6974696160448201527f6c697a656400000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b3660008037600080366000845af43d6000803e8061021e573d6000fd5b503d6000f35b7fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035473ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061027d575033155b1561028e5761028b816105a3565b50565b61028b61012d565b60606102c07fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035490565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806102f7575033155b1561040a57610305846105a3565b6000808573ffffffffffffffffffffffffffffffffffffffff16858560405161032f9291906107ee565b600060405180830381855af49150503d806000811461036a576040519150601f19603f3d011682016040523d82523d6000602084013e61036f565b606091505b509150915081610401576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603960248201527f50726f78793a2064656c656761746563616c6c20746f206e657720696d706c6560448201527f6d656e746174696f6e20636f6e7472616374206661696c65640000000000000060648201526084016101f8565b91506104129050565b61041261012d565b9392505050565b60006104437fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035490565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061047a575033155b156104a557507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b6104ad61012d565b90565b7fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035473ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610509575033155b1561028e5761028b8161060c565b60006105417fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035490565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610578575033155b156104a557507fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035490565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81815560405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a25050565b60006106367fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035490565b7fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61038381556040805173ffffffffffffffffffffffffffffffffffffffff80851682528616602082015292935090917f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f910160405180910390a1505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146106d857600080fd5b919050565b6000602082840312156106ef57600080fd5b610412826106b4565b60008060006040848603121561070d57600080fd5b610716846106b4565b9250602084013567ffffffffffffffff8082111561073357600080fd5b818601915086601f83011261074757600080fd5b81358181111561075657600080fd5b87602082850101111561076857600080fd5b6020830194508093505050509250925092565b600060208083528351808285015260005b818110156107a85785810183015185820160400152820161078c565b818111156107ba576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b818382376000910190815291905056fea164736f6c634300080f000ab53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103",
}

// ProxyContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ProxyContractMetaData.ABI instead.
var ProxyContractABI = ProxyContractMetaData.ABI

// ProxyContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProxyContractMetaData.Bin instead.
var ProxyContractBin = ProxyContractMetaData.Bin

// DeployProxyContract deploys a new Ethereum contract, binding an instance of ProxyContract to it.
func DeployProxyContract(auth *bind.TransactOpts, backend bind.ContractBackend, _admin common.Address) (common.Address, *types.Transaction, *ProxyContract, error) {
	parsed, err := ProxyContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProxyContractBin), backend, _admin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ProxyContract{ProxyContractCaller: ProxyContractCaller{contract: contract}, ProxyContractTransactor: ProxyContractTransactor{contract: contract}, ProxyContractFilterer: ProxyContractFilterer{contract: contract}}, nil
}

// ProxyContract is an auto generated Go binding around an Ethereum contract.
type ProxyContract struct {
	ProxyContractCaller     // Read-only binding to the contract
	ProxyContractTransactor // Write-only binding to the contract
	ProxyContractFilterer   // Log filterer for contract events
}

// ProxyContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProxyContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProxyContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProxyContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProxyContractSession struct {
	Contract     *ProxyContract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProxyContractCallerSession struct {
	Contract *ProxyContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ProxyContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProxyContractTransactorSession struct {
	Contract     *ProxyContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ProxyContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProxyContractRaw struct {
	Contract *ProxyContract // Generic contract binding to access the raw methods on
}

// ProxyContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProxyContractCallerRaw struct {
	Contract *ProxyContractCaller // Generic read-only contract binding to access the raw methods on
}

// ProxyContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProxyContractTransactorRaw struct {
	Contract *ProxyContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProxyContract creates a new instance of ProxyContract, bound to a specific deployed contract.
func NewProxyContract(address common.Address, backend bind.ContractBackend) (*ProxyContract, error) {
	contract, err := bindProxyContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProxyContract{ProxyContractCaller: ProxyContractCaller{contract: contract}, ProxyContractTransactor: ProxyContractTransactor{contract: contract}, ProxyContractFilterer: ProxyContractFilterer{contract: contract}}, nil
}

// NewProxyContractCaller creates a new read-only instance of ProxyContract, bound to a specific deployed contract.
func NewProxyContractCaller(address common.Address, caller bind.ContractCaller) (*ProxyContractCaller, error) {
	contract, err := bindProxyContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyContractCaller{contract: contract}, nil
}

// NewProxyContractTransactor creates a new write-only instance of ProxyContract, bound to a specific deployed contract.
func NewProxyContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ProxyContractTransactor, error) {
	contract, err := bindProxyContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyContractTransactor{contract: contract}, nil
}

// NewProxyContractFilterer creates a new log filterer instance of ProxyContract, bound to a specific deployed contract.
func NewProxyContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ProxyContractFilterer, error) {
	contract, err := bindProxyContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProxyContractFilterer{contract: contract}, nil
}

// bindProxyContract binds a generic wrapper to an already deployed contract.
func bindProxyContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProxyContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxyContract *ProxyContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProxyContract.Contract.ProxyContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxyContract *ProxyContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyContract.Contract.ProxyContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxyContract *ProxyContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxyContract.Contract.ProxyContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxyContract *ProxyContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProxyContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxyContract *ProxyContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxyContract *ProxyContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxyContract.Contract.contract.Transact(opts, method, params...)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_ProxyContract *ProxyContractTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyContract.contract.Transact(opts, "admin")
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_ProxyContract *ProxyContractSession) Admin() (*types.Transaction, error) {
	return _ProxyContract.Contract.Admin(&_ProxyContract.TransactOpts)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_ProxyContract *ProxyContractTransactorSession) Admin() (*types.Transaction, error) {
	return _ProxyContract.Contract.Admin(&_ProxyContract.TransactOpts)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address _admin) returns()
func (_ProxyContract *ProxyContractTransactor) ChangeAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _ProxyContract.contract.Transact(opts, "changeAdmin", _admin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address _admin) returns()
func (_ProxyContract *ProxyContractSession) ChangeAdmin(_admin common.Address) (*types.Transaction, error) {
	return _ProxyContract.Contract.ChangeAdmin(&_ProxyContract.TransactOpts, _admin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address _admin) returns()
func (_ProxyContract *ProxyContractTransactorSession) ChangeAdmin(_admin common.Address) (*types.Transaction, error) {
	return _ProxyContract.Contract.ChangeAdmin(&_ProxyContract.TransactOpts, _admin)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_ProxyContract *ProxyContractTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyContract.contract.Transact(opts, "implementation")
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_ProxyContract *ProxyContractSession) Implementation() (*types.Transaction, error) {
	return _ProxyContract.Contract.Implementation(&_ProxyContract.TransactOpts)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_ProxyContract *ProxyContractTransactorSession) Implementation() (*types.Transaction, error) {
	return _ProxyContract.Contract.Implementation(&_ProxyContract.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _implementation) returns()
func (_ProxyContract *ProxyContractTransactor) UpgradeTo(opts *bind.TransactOpts, _implementation common.Address) (*types.Transaction, error) {
	return _ProxyContract.contract.Transact(opts, "upgradeTo", _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _implementation) returns()
func (_ProxyContract *ProxyContractSession) UpgradeTo(_implementation common.Address) (*types.Transaction, error) {
	return _ProxyContract.Contract.UpgradeTo(&_ProxyContract.TransactOpts, _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _implementation) returns()
func (_ProxyContract *ProxyContractTransactorSession) UpgradeTo(_implementation common.Address) (*types.Transaction, error) {
	return _ProxyContract.Contract.UpgradeTo(&_ProxyContract.TransactOpts, _implementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address _implementation, bytes _data) payable returns(bytes)
func (_ProxyContract *ProxyContractTransactor) UpgradeToAndCall(opts *bind.TransactOpts, _implementation common.Address, _data []byte) (*types.Transaction, error) {
	return _ProxyContract.contract.Transact(opts, "upgradeToAndCall", _implementation, _data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address _implementation, bytes _data) payable returns(bytes)
func (_ProxyContract *ProxyContractSession) UpgradeToAndCall(_implementation common.Address, _data []byte) (*types.Transaction, error) {
	return _ProxyContract.Contract.UpgradeToAndCall(&_ProxyContract.TransactOpts, _implementation, _data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address _implementation, bytes _data) payable returns(bytes)
func (_ProxyContract *ProxyContractTransactorSession) UpgradeToAndCall(_implementation common.Address, _data []byte) (*types.Transaction, error) {
	return _ProxyContract.Contract.UpgradeToAndCall(&_ProxyContract.TransactOpts, _implementation, _data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ProxyContract *ProxyContractTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ProxyContract.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ProxyContract *ProxyContractSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ProxyContract.Contract.Fallback(&_ProxyContract.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ProxyContract *ProxyContractTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ProxyContract.Contract.Fallback(&_ProxyContract.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ProxyContract *ProxyContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyContract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ProxyContract *ProxyContractSession) Receive() (*types.Transaction, error) {
	return _ProxyContract.Contract.Receive(&_ProxyContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ProxyContract *ProxyContractTransactorSession) Receive() (*types.Transaction, error) {
	return _ProxyContract.Contract.Receive(&_ProxyContract.TransactOpts)
}

// ProxyContractAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ProxyContract contract.
type ProxyContractAdminChangedIterator struct {
	Event *ProxyContractAdminChanged // Event containing the contract specifics and raw log

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
func (it *ProxyContractAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyContractAdminChanged)
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
		it.Event = new(ProxyContractAdminChanged)
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
func (it *ProxyContractAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyContractAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyContractAdminChanged represents a AdminChanged event raised by the ProxyContract contract.
type ProxyContractAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ProxyContract *ProxyContractFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ProxyContractAdminChangedIterator, error) {

	logs, sub, err := _ProxyContract.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ProxyContractAdminChangedIterator{contract: _ProxyContract.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ProxyContract *ProxyContractFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ProxyContractAdminChanged) (event.Subscription, error) {

	logs, sub, err := _ProxyContract.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyContractAdminChanged)
				if err := _ProxyContract.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_ProxyContract *ProxyContractFilterer) ParseAdminChanged(log types.Log) (*ProxyContractAdminChanged, error) {
	event := new(ProxyContractAdminChanged)
	if err := _ProxyContract.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProxyContractUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ProxyContract contract.
type ProxyContractUpgradedIterator struct {
	Event *ProxyContractUpgraded // Event containing the contract specifics and raw log

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
func (it *ProxyContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyContractUpgraded)
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
		it.Event = new(ProxyContractUpgraded)
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
func (it *ProxyContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyContractUpgraded represents a Upgraded event raised by the ProxyContract contract.
type ProxyContractUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ProxyContract *ProxyContractFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ProxyContractUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ProxyContract.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ProxyContractUpgradedIterator{contract: _ProxyContract.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ProxyContract *ProxyContractFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ProxyContractUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ProxyContract.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyContractUpgraded)
				if err := _ProxyContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
// Solidity: event Upgraded(address indexed implementation)
func (_ProxyContract *ProxyContractFilterer) ParseUpgraded(log types.Log) (*ProxyContractUpgraded, error) {
	event := new(ProxyContractUpgraded)
	if err := _ProxyContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
