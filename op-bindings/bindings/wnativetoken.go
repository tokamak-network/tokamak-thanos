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

// WNativeTokenMetaData contains all meta data concerning the WNativeToken contract.
var WNativeTokenMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"guy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c0604052601360808190527f57726170706564204e6174697665546f6b656e0000000000000000000000000060a090815261003e9160009190610092565b5060408051808201909152600c8082526b2ba730ba34bb32aa37b5b2b760a11b602090920191825261007291600191610092565b506002805460ff1916601217905534801561008c57600080fd5b5061012d565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100d357805160ff1916838001178555610100565b82800160010185558215610100579182015b828111156101005782518255916020019190600101906100e5565b5061010c929150610110565b5090565b61012a91905b8082111561010c5760008155600101610116565b90565b6107f98061013c6000396000f3fe6080604052600436106100bc5760003560e01c8063313ce56711610074578063a9059cbb1161004e578063a9059cbb146102cb578063d0e30db0146100bc578063dd62ed3e14610311576100bc565b8063313ce5671461024b57806370a082311461027657806395d89b41146102b6576100bc565b806318160ddd116100a557806318160ddd146101aa57806323b872dd146101d15780632e1a7d4d14610221576100bc565b806306fdde03146100c6578063095ea7b314610150575b6100c4610359565b005b3480156100d257600080fd5b506100db6103a8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101155781810151838201526020016100fd565b50505050905090810190601f1680156101425780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015c57600080fd5b506101966004803603604081101561017357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610454565b604080519115158252519081900360200190f35b3480156101b657600080fd5b506101bf6104c7565b60408051918252519081900360200190f35b3480156101dd57600080fd5b50610196600480360360608110156101f457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602081013590911690604001356104cb565b34801561022d57600080fd5b506100c46004803603602081101561024457600080fd5b503561066b565b34801561025757600080fd5b50610260610700565b6040805160ff9092168252519081900360200190f35b34801561028257600080fd5b506101bf6004803603602081101561029957600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610709565b3480156102c257600080fd5b506100db61071b565b3480156102d757600080fd5b50610196600480360360408110156102ee57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610793565b34801561031d57600080fd5b506101bf6004803603604081101561033457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166107a7565b33600081815260036020908152604091829020805434908101909155825190815291517fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9281900390910190a2565b6000805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f8101849004840282018401909252818152929183018282801561044c5780601f106104215761010080835404028352916020019161044c565b820191906000526020600020905b81548152906001019060200180831161042f57829003601f168201915b505050505081565b33600081815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b4790565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600360205260408120548211156104fd57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84163314801590610573575073ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff14155b156105ed5773ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020548211156105b557600080fd5b73ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020805483900390555b73ffffffffffffffffffffffffffffffffffffffff808516600081815260036020908152604080832080548890039055938716808352918490208054870190558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35060019392505050565b3360009081526003602052604090205481111561068757600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f193505050501580156106c6573d6000803e3d6000fd5b5060408051828152905133917f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65919081900360200190a250565b60025460ff1681565b60036020526000908152604090205481565b60018054604080516020600284861615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f8101849004840282018401909252818152929183018282801561044c5780601f106104215761010080835404028352916020019161044c565b60006107a03384846104cb565b9392505050565b60046020908152600092835260408084209091529082529020548156fea265627a7a723158203539bf2477c4c8e49892cbdd9d4f49b11b95f78405b9345f0c445c2733db838364736f6c63430005110032",
}

// WNativeTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use WNativeTokenMetaData.ABI instead.
var WNativeTokenABI = WNativeTokenMetaData.ABI

// WNativeTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WNativeTokenMetaData.Bin instead.
var WNativeTokenBin = WNativeTokenMetaData.Bin

// DeployWNativeToken deploys a new Ethereum contract, binding an instance of WNativeToken to it.
func DeployWNativeToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WNativeToken, error) {
	parsed, err := WNativeTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WNativeTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WNativeToken{WNativeTokenCaller: WNativeTokenCaller{contract: contract}, WNativeTokenTransactor: WNativeTokenTransactor{contract: contract}, WNativeTokenFilterer: WNativeTokenFilterer{contract: contract}}, nil
}

// WNativeToken is an auto generated Go binding around an Ethereum contract.
type WNativeToken struct {
	WNativeTokenCaller     // Read-only binding to the contract
	WNativeTokenTransactor // Write-only binding to the contract
	WNativeTokenFilterer   // Log filterer for contract events
}

// WNativeTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type WNativeTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WNativeTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WNativeTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WNativeTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WNativeTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WNativeTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WNativeTokenSession struct {
	Contract     *WNativeToken     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WNativeTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WNativeTokenCallerSession struct {
	Contract *WNativeTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// WNativeTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WNativeTokenTransactorSession struct {
	Contract     *WNativeTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// WNativeTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type WNativeTokenRaw struct {
	Contract *WNativeToken // Generic contract binding to access the raw methods on
}

// WNativeTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WNativeTokenCallerRaw struct {
	Contract *WNativeTokenCaller // Generic read-only contract binding to access the raw methods on
}

// WNativeTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WNativeTokenTransactorRaw struct {
	Contract *WNativeTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWNativeToken creates a new instance of WNativeToken, bound to a specific deployed contract.
func NewWNativeToken(address common.Address, backend bind.ContractBackend) (*WNativeToken, error) {
	contract, err := bindWNativeToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WNativeToken{WNativeTokenCaller: WNativeTokenCaller{contract: contract}, WNativeTokenTransactor: WNativeTokenTransactor{contract: contract}, WNativeTokenFilterer: WNativeTokenFilterer{contract: contract}}, nil
}

// NewWNativeTokenCaller creates a new read-only instance of WNativeToken, bound to a specific deployed contract.
func NewWNativeTokenCaller(address common.Address, caller bind.ContractCaller) (*WNativeTokenCaller, error) {
	contract, err := bindWNativeToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WNativeTokenCaller{contract: contract}, nil
}

// NewWNativeTokenTransactor creates a new write-only instance of WNativeToken, bound to a specific deployed contract.
func NewWNativeTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*WNativeTokenTransactor, error) {
	contract, err := bindWNativeToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WNativeTokenTransactor{contract: contract}, nil
}

// NewWNativeTokenFilterer creates a new log filterer instance of WNativeToken, bound to a specific deployed contract.
func NewWNativeTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*WNativeTokenFilterer, error) {
	contract, err := bindWNativeToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WNativeTokenFilterer{contract: contract}, nil
}

// bindWNativeToken binds a generic wrapper to an already deployed contract.
func bindWNativeToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WNativeTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WNativeToken *WNativeTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WNativeToken.Contract.WNativeTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WNativeToken *WNativeTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WNativeToken.Contract.WNativeTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WNativeToken *WNativeTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WNativeToken.Contract.WNativeTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WNativeToken *WNativeTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WNativeToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WNativeToken *WNativeTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WNativeToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WNativeToken *WNativeTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WNativeToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WNativeToken *WNativeTokenCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WNativeToken.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WNativeToken *WNativeTokenSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WNativeToken.Contract.Allowance(&_WNativeToken.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WNativeToken *WNativeTokenCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WNativeToken.Contract.Allowance(&_WNativeToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WNativeToken *WNativeTokenCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WNativeToken.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WNativeToken *WNativeTokenSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WNativeToken.Contract.BalanceOf(&_WNativeToken.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WNativeToken *WNativeTokenCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WNativeToken.Contract.BalanceOf(&_WNativeToken.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WNativeToken *WNativeTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WNativeToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WNativeToken *WNativeTokenSession) Decimals() (uint8, error) {
	return _WNativeToken.Contract.Decimals(&_WNativeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WNativeToken *WNativeTokenCallerSession) Decimals() (uint8, error) {
	return _WNativeToken.Contract.Decimals(&_WNativeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WNativeToken *WNativeTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WNativeToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WNativeToken *WNativeTokenSession) Name() (string, error) {
	return _WNativeToken.Contract.Name(&_WNativeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WNativeToken *WNativeTokenCallerSession) Name() (string, error) {
	return _WNativeToken.Contract.Name(&_WNativeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WNativeToken *WNativeTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WNativeToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WNativeToken *WNativeTokenSession) Symbol() (string, error) {
	return _WNativeToken.Contract.Symbol(&_WNativeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WNativeToken *WNativeTokenCallerSession) Symbol() (string, error) {
	return _WNativeToken.Contract.Symbol(&_WNativeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WNativeToken *WNativeTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WNativeToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WNativeToken *WNativeTokenSession) TotalSupply() (*big.Int, error) {
	return _WNativeToken.Contract.TotalSupply(&_WNativeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WNativeToken *WNativeTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _WNativeToken.Contract.TotalSupply(&_WNativeToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenTransactor) Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.contract.Transact(opts, "approve", guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.Approve(&_WNativeToken.TransactOpts, guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenTransactorSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.Approve(&_WNativeToken.TransactOpts, guy, wad)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WNativeToken *WNativeTokenTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WNativeToken.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WNativeToken *WNativeTokenSession) Deposit() (*types.Transaction, error) {
	return _WNativeToken.Contract.Deposit(&_WNativeToken.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WNativeToken *WNativeTokenTransactorSession) Deposit() (*types.Transaction, error) {
	return _WNativeToken.Contract.Deposit(&_WNativeToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenTransactor) Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.contract.Transact(opts, "transfer", dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.Transfer(&_WNativeToken.TransactOpts, dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenTransactorSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.Transfer(&_WNativeToken.TransactOpts, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenTransactor) TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.contract.Transact(opts, "transferFrom", src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.TransferFrom(&_WNativeToken.TransactOpts, src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WNativeToken *WNativeTokenTransactorSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.TransferFrom(&_WNativeToken.TransactOpts, src, dst, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WNativeToken *WNativeTokenTransactor) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.contract.Transact(opts, "withdraw", wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WNativeToken *WNativeTokenSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.Withdraw(&_WNativeToken.TransactOpts, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WNativeToken *WNativeTokenTransactorSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WNativeToken.Contract.Withdraw(&_WNativeToken.TransactOpts, wad)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WNativeToken *WNativeTokenTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _WNativeToken.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WNativeToken *WNativeTokenSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WNativeToken.Contract.Fallback(&_WNativeToken.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WNativeToken *WNativeTokenTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WNativeToken.Contract.Fallback(&_WNativeToken.TransactOpts, calldata)
}

// WNativeTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WNativeToken contract.
type WNativeTokenApprovalIterator struct {
	Event *WNativeTokenApproval // Event containing the contract specifics and raw log

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
func (it *WNativeTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WNativeTokenApproval)
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
		it.Event = new(WNativeTokenApproval)
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
func (it *WNativeTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WNativeTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WNativeTokenApproval represents a Approval event raised by the WNativeToken contract.
type WNativeTokenApproval struct {
	Src common.Address
	Guy common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*WNativeTokenApprovalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WNativeToken.contract.FilterLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return &WNativeTokenApprovalIterator{contract: _WNativeToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WNativeTokenApproval, src []common.Address, guy []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WNativeToken.contract.WatchLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WNativeTokenApproval)
				if err := _WNativeToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) ParseApproval(log types.Log) (*WNativeTokenApproval, error) {
	event := new(WNativeTokenApproval)
	if err := _WNativeToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WNativeTokenDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the WNativeToken contract.
type WNativeTokenDepositIterator struct {
	Event *WNativeTokenDeposit // Event containing the contract specifics and raw log

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
func (it *WNativeTokenDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WNativeTokenDeposit)
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
		it.Event = new(WNativeTokenDeposit)
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
func (it *WNativeTokenDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WNativeTokenDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WNativeTokenDeposit represents a Deposit event raised by the WNativeToken contract.
type WNativeTokenDeposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WNativeTokenDepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WNativeToken.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &WNativeTokenDepositIterator{contract: _WNativeToken.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WNativeTokenDeposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WNativeToken.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WNativeTokenDeposit)
				if err := _WNativeToken.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) ParseDeposit(log types.Log) (*WNativeTokenDeposit, error) {
	event := new(WNativeTokenDeposit)
	if err := _WNativeToken.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WNativeTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WNativeToken contract.
type WNativeTokenTransferIterator struct {
	Event *WNativeTokenTransfer // Event containing the contract specifics and raw log

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
func (it *WNativeTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WNativeTokenTransfer)
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
		it.Event = new(WNativeTokenTransfer)
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
func (it *WNativeTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WNativeTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WNativeTokenTransfer represents a Transfer event raised by the WNativeToken contract.
type WNativeTokenTransfer struct {
	Src common.Address
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*WNativeTokenTransferIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WNativeToken.contract.FilterLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return &WNativeTokenTransferIterator{contract: _WNativeToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WNativeTokenTransfer, src []common.Address, dst []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WNativeToken.contract.WatchLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WNativeTokenTransfer)
				if err := _WNativeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) ParseTransfer(log types.Log) (*WNativeTokenTransfer, error) {
	event := new(WNativeTokenTransfer)
	if err := _WNativeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WNativeTokenWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the WNativeToken contract.
type WNativeTokenWithdrawalIterator struct {
	Event *WNativeTokenWithdrawal // Event containing the contract specifics and raw log

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
func (it *WNativeTokenWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WNativeTokenWithdrawal)
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
		it.Event = new(WNativeTokenWithdrawal)
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
func (it *WNativeTokenWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WNativeTokenWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WNativeTokenWithdrawal represents a Withdrawal event raised by the WNativeToken contract.
type WNativeTokenWithdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WNativeTokenWithdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WNativeToken.contract.FilterLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &WNativeTokenWithdrawalIterator{contract: _WNativeToken.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WNativeTokenWithdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WNativeToken.contract.WatchLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WNativeTokenWithdrawal)
				if err := _WNativeToken.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WNativeToken *WNativeTokenFilterer) ParseWithdrawal(log types.Log) (*WNativeTokenWithdrawal, error) {
	event := new(WNativeTokenWithdrawal)
	if err := _WNativeToken.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
