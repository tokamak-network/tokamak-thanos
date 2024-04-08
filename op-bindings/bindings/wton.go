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

// WTONMetaData contains all meta data concerning the WTON contract.
var WTONMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"guy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c0604052600b60808190526a2bb930b83832b2102a27a760a91b60a090815261002c9160009190610078565b50604080518082019091526004808252632baa27a760e11b602090920191825261005891600191610078565b506002805460ff1916601217905534801561007257600080fd5b50610113565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100b957805160ff19168380011785556100e6565b828001600101855582156100e6579182015b828111156100e65782518255916020019190600101906100cb565b506100f29291506100f6565b5090565b61011091905b808211156100f257600081556001016100fc565b90565b6107f9806101226000396000f3fe6080604052600436106100bc5760003560e01c8063313ce56711610074578063a9059cbb1161004e578063a9059cbb146102cb578063d0e30db0146100bc578063dd62ed3e14610311576100bc565b8063313ce5671461024b57806370a082311461027657806395d89b41146102b6576100bc565b806318160ddd116100a557806318160ddd146101aa57806323b872dd146101d15780632e1a7d4d14610221576100bc565b806306fdde03146100c6578063095ea7b314610150575b6100c4610359565b005b3480156100d257600080fd5b506100db6103a8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101155781810151838201526020016100fd565b50505050905090810190601f1680156101425780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015c57600080fd5b506101966004803603604081101561017357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610454565b604080519115158252519081900360200190f35b3480156101b657600080fd5b506101bf6104c7565b60408051918252519081900360200190f35b3480156101dd57600080fd5b50610196600480360360608110156101f457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602081013590911690604001356104cb565b34801561022d57600080fd5b506100c46004803603602081101561024457600080fd5b503561066b565b34801561025757600080fd5b50610260610700565b6040805160ff9092168252519081900360200190f35b34801561028257600080fd5b506101bf6004803603602081101561029957600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610709565b3480156102c257600080fd5b506100db61071b565b3480156102d757600080fd5b50610196600480360360408110156102ee57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610793565b34801561031d57600080fd5b506101bf6004803603604081101561033457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166107a7565b33600081815260036020908152604091829020805434908101909155825190815291517fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9281900390910190a2565b6000805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f8101849004840282018401909252818152929183018282801561044c5780601f106104215761010080835404028352916020019161044c565b820191906000526020600020905b81548152906001019060200180831161042f57829003601f168201915b505050505081565b33600081815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b4790565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600360205260408120548211156104fd57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84163314801590610573575073ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff14155b156105ed5773ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020548211156105b557600080fd5b73ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020805483900390555b73ffffffffffffffffffffffffffffffffffffffff808516600081815260036020908152604080832080548890039055938716808352918490208054870190558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35060019392505050565b3360009081526003602052604090205481111561068757600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f193505050501580156106c6573d6000803e3d6000fd5b5060408051828152905133917f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65919081900360200190a250565b60025460ff1681565b60036020526000908152604090205481565b60018054604080516020600284861615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f8101849004840282018401909252818152929183018282801561044c5780601f106104215761010080835404028352916020019161044c565b60006107a03384846104cb565b9392505050565b60046020908152600092835260408084209091529082529020548156fea265627a7a723158209996572c9803dbed46463160438f5e5c7bd2e0a081ad4bb68099ec90ce3f733064736f6c63430005110032",
}

// WTONABI is the input ABI used to generate the binding from.
// Deprecated: Use WTONMetaData.ABI instead.
var WTONABI = WTONMetaData.ABI

// WTONBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WTONMetaData.Bin instead.
var WTONBin = WTONMetaData.Bin

// DeployWTON deploys a new Ethereum contract, binding an instance of WTON to it.
func DeployWTON(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WTON, error) {
	parsed, err := WTONMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WTONBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WTON{WTONCaller: WTONCaller{contract: contract}, WTONTransactor: WTONTransactor{contract: contract}, WTONFilterer: WTONFilterer{contract: contract}}, nil
}

// WTON is an auto generated Go binding around an Ethereum contract.
type WTON struct {
	WTONCaller     // Read-only binding to the contract
	WTONTransactor // Write-only binding to the contract
	WTONFilterer   // Log filterer for contract events
}

// WTONCaller is an auto generated read-only Go binding around an Ethereum contract.
type WTONCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WTONTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WTONTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WTONFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WTONFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WTONSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WTONSession struct {
	Contract     *WTON             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WTONCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WTONCallerSession struct {
	Contract *WTONCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// WTONTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WTONTransactorSession struct {
	Contract     *WTONTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WTONRaw is an auto generated low-level Go binding around an Ethereum contract.
type WTONRaw struct {
	Contract *WTON // Generic contract binding to access the raw methods on
}

// WTONCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WTONCallerRaw struct {
	Contract *WTONCaller // Generic read-only contract binding to access the raw methods on
}

// WTONTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WTONTransactorRaw struct {
	Contract *WTONTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWTON creates a new instance of WTON, bound to a specific deployed contract.
func NewWTON(address common.Address, backend bind.ContractBackend) (*WTON, error) {
	contract, err := bindWTON(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WTON{WTONCaller: WTONCaller{contract: contract}, WTONTransactor: WTONTransactor{contract: contract}, WTONFilterer: WTONFilterer{contract: contract}}, nil
}

// NewWTONCaller creates a new read-only instance of WTON, bound to a specific deployed contract.
func NewWTONCaller(address common.Address, caller bind.ContractCaller) (*WTONCaller, error) {
	contract, err := bindWTON(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WTONCaller{contract: contract}, nil
}

// NewWTONTransactor creates a new write-only instance of WTON, bound to a specific deployed contract.
func NewWTONTransactor(address common.Address, transactor bind.ContractTransactor) (*WTONTransactor, error) {
	contract, err := bindWTON(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WTONTransactor{contract: contract}, nil
}

// NewWTONFilterer creates a new log filterer instance of WTON, bound to a specific deployed contract.
func NewWTONFilterer(address common.Address, filterer bind.ContractFilterer) (*WTONFilterer, error) {
	contract, err := bindWTON(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WTONFilterer{contract: contract}, nil
}

// bindWTON binds a generic wrapper to an already deployed contract.
func bindWTON(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WTONMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WTON *WTONRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WTON.Contract.WTONCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WTON *WTONRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WTON.Contract.WTONTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WTON *WTONRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WTON.Contract.WTONTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WTON *WTONCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WTON.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WTON *WTONTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WTON.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WTON *WTONTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WTON.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WTON *WTONCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WTON.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WTON *WTONSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WTON.Contract.Allowance(&_WTON.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WTON *WTONCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WTON.Contract.Allowance(&_WTON.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WTON *WTONCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WTON.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WTON *WTONSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WTON.Contract.BalanceOf(&_WTON.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WTON *WTONCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WTON.Contract.BalanceOf(&_WTON.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WTON *WTONCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WTON.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WTON *WTONSession) Decimals() (uint8, error) {
	return _WTON.Contract.Decimals(&_WTON.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WTON *WTONCallerSession) Decimals() (uint8, error) {
	return _WTON.Contract.Decimals(&_WTON.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WTON *WTONCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WTON.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WTON *WTONSession) Name() (string, error) {
	return _WTON.Contract.Name(&_WTON.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WTON *WTONCallerSession) Name() (string, error) {
	return _WTON.Contract.Name(&_WTON.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WTON *WTONCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WTON.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WTON *WTONSession) Symbol() (string, error) {
	return _WTON.Contract.Symbol(&_WTON.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WTON *WTONCallerSession) Symbol() (string, error) {
	return _WTON.Contract.Symbol(&_WTON.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WTON *WTONCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WTON.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WTON *WTONSession) TotalSupply() (*big.Int, error) {
	return _WTON.Contract.TotalSupply(&_WTON.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WTON *WTONCallerSession) TotalSupply() (*big.Int, error) {
	return _WTON.Contract.TotalSupply(&_WTON.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WTON *WTONTransactor) Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.contract.Transact(opts, "approve", guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WTON *WTONSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.Approve(&_WTON.TransactOpts, guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WTON *WTONTransactorSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.Approve(&_WTON.TransactOpts, guy, wad)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WTON *WTONTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WTON.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WTON *WTONSession) Deposit() (*types.Transaction, error) {
	return _WTON.Contract.Deposit(&_WTON.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WTON *WTONTransactorSession) Deposit() (*types.Transaction, error) {
	return _WTON.Contract.Deposit(&_WTON.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WTON *WTONTransactor) Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.contract.Transact(opts, "transfer", dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WTON *WTONSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.Transfer(&_WTON.TransactOpts, dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WTON *WTONTransactorSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.Transfer(&_WTON.TransactOpts, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WTON *WTONTransactor) TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.contract.Transact(opts, "transferFrom", src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WTON *WTONSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.TransferFrom(&_WTON.TransactOpts, src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WTON *WTONTransactorSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.TransferFrom(&_WTON.TransactOpts, src, dst, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WTON *WTONTransactor) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _WTON.contract.Transact(opts, "withdraw", wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WTON *WTONSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.Withdraw(&_WTON.TransactOpts, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WTON *WTONTransactorSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WTON.Contract.Withdraw(&_WTON.TransactOpts, wad)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WTON *WTONTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _WTON.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WTON *WTONSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WTON.Contract.Fallback(&_WTON.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WTON *WTONTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WTON.Contract.Fallback(&_WTON.TransactOpts, calldata)
}

// WTONApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WTON contract.
type WTONApprovalIterator struct {
	Event *WTONApproval // Event containing the contract specifics and raw log

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
func (it *WTONApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WTONApproval)
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
		it.Event = new(WTONApproval)
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
func (it *WTONApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WTONApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WTONApproval represents a Approval event raised by the WTON contract.
type WTONApproval struct {
	Src common.Address
	Guy common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WTON *WTONFilterer) FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*WTONApprovalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WTON.contract.FilterLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return &WTONApprovalIterator{contract: _WTON.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WTON *WTONFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WTONApproval, src []common.Address, guy []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WTON.contract.WatchLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WTONApproval)
				if err := _WTON.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_WTON *WTONFilterer) ParseApproval(log types.Log) (*WTONApproval, error) {
	event := new(WTONApproval)
	if err := _WTON.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WTONDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the WTON contract.
type WTONDepositIterator struct {
	Event *WTONDeposit // Event containing the contract specifics and raw log

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
func (it *WTONDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WTONDeposit)
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
		it.Event = new(WTONDeposit)
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
func (it *WTONDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WTONDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WTONDeposit represents a Deposit event raised by the WTON contract.
type WTONDeposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WTON *WTONFilterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WTONDepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WTON.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &WTONDepositIterator{contract: _WTON.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WTON *WTONFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WTONDeposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WTON.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WTONDeposit)
				if err := _WTON.contract.UnpackLog(event, "Deposit", log); err != nil {
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
func (_WTON *WTONFilterer) ParseDeposit(log types.Log) (*WTONDeposit, error) {
	event := new(WTONDeposit)
	if err := _WTON.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WTONTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WTON contract.
type WTONTransferIterator struct {
	Event *WTONTransfer // Event containing the contract specifics and raw log

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
func (it *WTONTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WTONTransfer)
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
		it.Event = new(WTONTransfer)
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
func (it *WTONTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WTONTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WTONTransfer represents a Transfer event raised by the WTON contract.
type WTONTransfer struct {
	Src common.Address
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WTON *WTONFilterer) FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*WTONTransferIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WTON.contract.FilterLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return &WTONTransferIterator{contract: _WTON.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WTON *WTONFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WTONTransfer, src []common.Address, dst []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WTON.contract.WatchLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WTONTransfer)
				if err := _WTON.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_WTON *WTONFilterer) ParseTransfer(log types.Log) (*WTONTransfer, error) {
	event := new(WTONTransfer)
	if err := _WTON.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WTONWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the WTON contract.
type WTONWithdrawalIterator struct {
	Event *WTONWithdrawal // Event containing the contract specifics and raw log

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
func (it *WTONWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WTONWithdrawal)
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
		it.Event = new(WTONWithdrawal)
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
func (it *WTONWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WTONWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WTONWithdrawal represents a Withdrawal event raised by the WTON contract.
type WTONWithdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WTON *WTONFilterer) FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WTONWithdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WTON.contract.FilterLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &WTONWithdrawalIterator{contract: _WTON.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WTON *WTONFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WTONWithdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WTON.contract.WatchLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WTONWithdrawal)
				if err := _WTON.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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
func (_WTON *WTONFilterer) ParseWithdrawal(log types.Log) (*WTONWithdrawal, error) {
	event := new(WTONWithdrawal)
	if err := _WTON.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
