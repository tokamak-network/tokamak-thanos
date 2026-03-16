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

// MultiTokenPaymasterTokenConfig is an auto generated low-level Go binding around an user-defined struct.
type MultiTokenPaymasterTokenConfig struct {
	Enabled  bool
	Oracle   common.Address
	Markup   *big.Int
	Decimals uint8
}

// PackedUserOperation is shared with aaentrypoint.go in this package.

// MultiTokenPaymasterPredeployMetaData contains all meta data concerning the MultiTokenPaymasterPredeploy contract.
var MultiTokenPaymasterPredeployMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addStake\",\"inputs\":[{\"name\":\"unstakeDelaySec\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"contractITokenPriceOracle\"},{\"name\":\"markupPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collectedFees\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"entryPoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"estimateTokenCost\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"estimatedTonGasCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"tokenCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenCostWithMarkup\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"estimateTokenCostPublic\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tonAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDeposit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenConfig\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structMultiTokenPaymaster.TokenConfig\",\"components\":[{\"name\":\"enabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"contractITokenPriceOracle\"},{\"name\":\"markup\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_entryPoint\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minCharge\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postOp\",\"inputs\":[{\"name\":\"mode\",\"type\":\"uint8\",\"internalType\":\"enumIPaymaster.PostOpMode\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actualGasCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualUserOpFeePerGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinCharge\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportedTokens\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"enabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"contractITokenPriceOracle\"},{\"name\":\"markup\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenList\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTokenConfig\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"contractITokenPriceOracle\"},{\"name\":\"markupPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatePaymasterUserOp\",\"inputs\":[{\"name\":\"userOp\",\"type\":\"tuple\",\"internalType\":\"structPackedUserOperation\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"accountGasLimits\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preVerificationGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFees\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"userOpHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validationData\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawCollectedFees\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawStake\",\"inputs\":[{\"name\":\"withdrawAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTo\",\"inputs\":[{\"name\":\"withdrawAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FeesCollected\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeesWithdrawn\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenAdded\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"markup\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenConfigUpdated\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"markup\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRemoved\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// MultiTokenPaymasterPredeployABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiTokenPaymasterPredeployMetaData.ABI instead.
var MultiTokenPaymasterPredeployABI = MultiTokenPaymasterPredeployMetaData.ABI

// MultiTokenPaymasterPredeploy is an auto generated Go binding around an Ethereum contract.
type MultiTokenPaymasterPredeploy struct {
	MultiTokenPaymasterPredeployCaller     // Read-only binding to the contract
	MultiTokenPaymasterPredeployTransactor // Write-only binding to the contract
	MultiTokenPaymasterPredeployFilterer   // Log filterer for contract events
}

// MultiTokenPaymasterPredeployCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiTokenPaymasterPredeployCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiTokenPaymasterPredeployTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiTokenPaymasterPredeployTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiTokenPaymasterPredeployFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiTokenPaymasterPredeployFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiTokenPaymasterPredeploySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiTokenPaymasterPredeploySession struct {
	Contract     *MultiTokenPaymasterPredeploy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// MultiTokenPaymasterPredeployCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiTokenPaymasterPredeployCallerSession struct {
	Contract *MultiTokenPaymasterPredeployCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// MultiTokenPaymasterPredeployTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiTokenPaymasterPredeployTransactorSession struct {
	Contract     *MultiTokenPaymasterPredeployTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// MultiTokenPaymasterPredeployRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiTokenPaymasterPredeployRaw struct {
	Contract *MultiTokenPaymasterPredeploy // Generic contract binding to access the raw methods on
}

// MultiTokenPaymasterPredeployCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiTokenPaymasterPredeployCallerRaw struct {
	Contract *MultiTokenPaymasterPredeployCaller // Generic read-only contract binding to access the raw methods on
}

// MultiTokenPaymasterPredeployTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiTokenPaymasterPredeployTransactorRaw struct {
	Contract *MultiTokenPaymasterPredeployTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiTokenPaymasterPredeploy creates a new instance of MultiTokenPaymasterPredeploy, bound to a specific deployed contract.
func NewMultiTokenPaymasterPredeploy(address common.Address, backend bind.ContractBackend) (*MultiTokenPaymasterPredeploy, error) {
	contract, err := bindMultiTokenPaymasterPredeploy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeploy{MultiTokenPaymasterPredeployCaller: MultiTokenPaymasterPredeployCaller{contract: contract}, MultiTokenPaymasterPredeployTransactor: MultiTokenPaymasterPredeployTransactor{contract: contract}, MultiTokenPaymasterPredeployFilterer: MultiTokenPaymasterPredeployFilterer{contract: contract}}, nil
}

// NewMultiTokenPaymasterPredeployCaller creates a new read-only instance of MultiTokenPaymasterPredeploy, bound to a specific deployed contract.
func NewMultiTokenPaymasterPredeployCaller(address common.Address, caller bind.ContractCaller) (*MultiTokenPaymasterPredeployCaller, error) {
	contract, err := bindMultiTokenPaymasterPredeploy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployCaller{contract: contract}, nil
}

// NewMultiTokenPaymasterPredeployTransactor creates a new write-only instance of MultiTokenPaymasterPredeploy, bound to a specific deployed contract.
func NewMultiTokenPaymasterPredeployTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiTokenPaymasterPredeployTransactor, error) {
	contract, err := bindMultiTokenPaymasterPredeploy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployTransactor{contract: contract}, nil
}

// NewMultiTokenPaymasterPredeployFilterer creates a new log filterer instance of MultiTokenPaymasterPredeploy, bound to a specific deployed contract.
func NewMultiTokenPaymasterPredeployFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiTokenPaymasterPredeployFilterer, error) {
	contract, err := bindMultiTokenPaymasterPredeploy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployFilterer{contract: contract}, nil
}

// bindMultiTokenPaymasterPredeploy binds a generic wrapper to an already deployed contract.
func bindMultiTokenPaymasterPredeploy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiTokenPaymasterPredeployABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiTokenPaymasterPredeploy.Contract.MultiTokenPaymasterPredeployCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.MultiTokenPaymasterPredeployTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.MultiTokenPaymasterPredeployTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiTokenPaymasterPredeploy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.contract.Transact(opts, method, params...)
}

// CollectedFees is a free data retrieval call binding the contract method 0x1cead9a7.
//
// Solidity: function collectedFees(address ) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) CollectedFees(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "collectedFees", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CollectedFees is a free data retrieval call binding the contract method 0x1cead9a7.
//
// Solidity: function collectedFees(address ) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) CollectedFees(arg0 common.Address) (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.CollectedFees(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// CollectedFees is a free data retrieval call binding the contract method 0x1cead9a7.
//
// Solidity: function collectedFees(address ) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) CollectedFees(arg0 common.Address) (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.CollectedFees(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) EntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) EntryPoint() (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.EntryPoint(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) EntryPoint() (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.EntryPoint(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// EstimateTokenCost is a free data retrieval call binding the contract method 0x8c509739.
//
// Solidity: function estimateTokenCost(address token, uint256 estimatedTonGasCost) view returns(uint256 tokenCost, uint256 tokenCostWithMarkup)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) EstimateTokenCost(opts *bind.CallOpts, token common.Address, estimatedTonGasCost *big.Int) (struct {
	TokenCost           *big.Int
	TokenCostWithMarkup *big.Int
}, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "estimateTokenCost", token, estimatedTonGasCost)

	outstruct := new(struct {
		TokenCost           *big.Int
		TokenCostWithMarkup *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenCost = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TokenCostWithMarkup = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// EstimateTokenCost is a free data retrieval call binding the contract method 0x8c509739.
//
// Solidity: function estimateTokenCost(address token, uint256 estimatedTonGasCost) view returns(uint256 tokenCost, uint256 tokenCostWithMarkup)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) EstimateTokenCost(token common.Address, estimatedTonGasCost *big.Int) (struct {
	TokenCost           *big.Int
	TokenCostWithMarkup *big.Int
}, error) {
	return _MultiTokenPaymasterPredeploy.Contract.EstimateTokenCost(&_MultiTokenPaymasterPredeploy.CallOpts, token, estimatedTonGasCost)
}

// EstimateTokenCost is a free data retrieval call binding the contract method 0x8c509739.
//
// Solidity: function estimateTokenCost(address token, uint256 estimatedTonGasCost) view returns(uint256 tokenCost, uint256 tokenCostWithMarkup)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) EstimateTokenCost(token common.Address, estimatedTonGasCost *big.Int) (struct {
	TokenCost           *big.Int
	TokenCostWithMarkup *big.Int
}, error) {
	return _MultiTokenPaymasterPredeploy.Contract.EstimateTokenCost(&_MultiTokenPaymasterPredeploy.CallOpts, token, estimatedTonGasCost)
}

// EstimateTokenCostPublic is a free data retrieval call binding the contract method 0x72a9abe2.
//
// Solidity: function estimateTokenCostPublic(address token, uint256 tonAmount) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) EstimateTokenCostPublic(opts *bind.CallOpts, token common.Address, tonAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "estimateTokenCostPublic", token, tonAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateTokenCostPublic is a free data retrieval call binding the contract method 0x72a9abe2.
//
// Solidity: function estimateTokenCostPublic(address token, uint256 tonAmount) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) EstimateTokenCostPublic(token common.Address, tonAmount *big.Int) (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.EstimateTokenCostPublic(&_MultiTokenPaymasterPredeploy.CallOpts, token, tonAmount)
}

// EstimateTokenCostPublic is a free data retrieval call binding the contract method 0x72a9abe2.
//
// Solidity: function estimateTokenCostPublic(address token, uint256 tonAmount) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) EstimateTokenCostPublic(token common.Address, tonAmount *big.Int) (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.EstimateTokenCostPublic(&_MultiTokenPaymasterPredeploy.CallOpts, token, tonAmount)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) GetDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "getDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) GetDeposit() (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.GetDeposit(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) GetDeposit() (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.GetDeposit(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// GetTokenConfig is a free data retrieval call binding the contract method 0xcb67e3b1.
//
// Solidity: function getTokenConfig(address token) view returns((bool,address,uint256,uint8))
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) GetTokenConfig(opts *bind.CallOpts, token common.Address) (MultiTokenPaymasterTokenConfig, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "getTokenConfig", token)

	if err != nil {
		return *new(MultiTokenPaymasterTokenConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiTokenPaymasterTokenConfig)).(*MultiTokenPaymasterTokenConfig)

	return out0, err

}

// GetTokenConfig is a free data retrieval call binding the contract method 0xcb67e3b1.
//
// Solidity: function getTokenConfig(address token) view returns((bool,address,uint256,uint8))
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) GetTokenConfig(token common.Address) (MultiTokenPaymasterTokenConfig, error) {
	return _MultiTokenPaymasterPredeploy.Contract.GetTokenConfig(&_MultiTokenPaymasterPredeploy.CallOpts, token)
}

// GetTokenConfig is a free data retrieval call binding the contract method 0xcb67e3b1.
//
// Solidity: function getTokenConfig(address token) view returns((bool,address,uint256,uint8))
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) GetTokenConfig(token common.Address) (MultiTokenPaymasterTokenConfig, error) {
	return _MultiTokenPaymasterPredeploy.Contract.GetTokenConfig(&_MultiTokenPaymasterPredeploy.CallOpts, token)
}

// MinCharge is a free data retrieval call binding the contract method 0xdf26e26b.
//
// Solidity: function minCharge(address ) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) MinCharge(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "minCharge", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinCharge is a free data retrieval call binding the contract method 0xdf26e26b.
//
// Solidity: function minCharge(address ) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) MinCharge(arg0 common.Address) (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.MinCharge(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// MinCharge is a free data retrieval call binding the contract method 0xdf26e26b.
//
// Solidity: function minCharge(address ) view returns(uint256)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) MinCharge(arg0 common.Address) (*big.Int, error) {
	return _MultiTokenPaymasterPredeploy.Contract.MinCharge(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) Owner() (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.Owner(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) Owner() (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.Owner(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) PendingOwner() (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.PendingOwner(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) PendingOwner() (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.PendingOwner(&_MultiTokenPaymasterPredeploy.CallOpts)
}

// SupportedTokens is a free data retrieval call binding the contract method 0x68c4ac26.
//
// Solidity: function supportedTokens(address ) view returns(bool enabled, address oracle, uint256 markup, uint8 decimals)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) SupportedTokens(opts *bind.CallOpts, arg0 common.Address) (struct {
	Enabled  bool
	Oracle   common.Address
	Markup   *big.Int
	Decimals uint8
}, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "supportedTokens", arg0)

	outstruct := new(struct {
		Enabled  bool
		Oracle   common.Address
		Markup   *big.Int
		Decimals uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Enabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Oracle = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Markup = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Decimals = *abi.ConvertType(out[3], new(uint8)).(*uint8)

	return *outstruct, err

}

// SupportedTokens is a free data retrieval call binding the contract method 0x68c4ac26.
//
// Solidity: function supportedTokens(address ) view returns(bool enabled, address oracle, uint256 markup, uint8 decimals)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) SupportedTokens(arg0 common.Address) (struct {
	Enabled  bool
	Oracle   common.Address
	Markup   *big.Int
	Decimals uint8
}, error) {
	return _MultiTokenPaymasterPredeploy.Contract.SupportedTokens(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// SupportedTokens is a free data retrieval call binding the contract method 0x68c4ac26.
//
// Solidity: function supportedTokens(address ) view returns(bool enabled, address oracle, uint256 markup, uint8 decimals)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) SupportedTokens(arg0 common.Address) (struct {
	Enabled  bool
	Oracle   common.Address
	Markup   *big.Int
	Decimals uint8
}, error) {
	return _MultiTokenPaymasterPredeploy.Contract.SupportedTokens(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// TokenList is a free data retrieval call binding the contract method 0x9ead7222.
//
// Solidity: function tokenList(uint256 ) view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCaller) TokenList(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _MultiTokenPaymasterPredeploy.contract.Call(opts, &out, "tokenList", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenList is a free data retrieval call binding the contract method 0x9ead7222.
//
// Solidity: function tokenList(uint256 ) view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) TokenList(arg0 *big.Int) (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.TokenList(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// TokenList is a free data retrieval call binding the contract method 0x9ead7222.
//
// Solidity: function tokenList(uint256 ) view returns(address)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployCallerSession) TokenList(arg0 *big.Int) (common.Address, error) {
	return _MultiTokenPaymasterPredeploy.Contract.TokenList(&_MultiTokenPaymasterPredeploy.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.AcceptOwnership(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.AcceptOwnership(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) AddStake(opts *bind.TransactOpts, unstakeDelaySec uint32) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "addStake", unstakeDelaySec)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) AddStake(unstakeDelaySec uint32) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.AddStake(&_MultiTokenPaymasterPredeploy.TransactOpts, unstakeDelaySec)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) AddStake(unstakeDelaySec uint32) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.AddStake(&_MultiTokenPaymasterPredeploy.TransactOpts, unstakeDelaySec)
}

// AddToken is a paid mutator transaction binding the contract method 0x9acad299.
//
// Solidity: function addToken(address token, address oracle, uint256 markupPercent, uint8 decimals) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) AddToken(opts *bind.TransactOpts, token common.Address, oracle common.Address, markupPercent *big.Int, decimals uint8) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "addToken", token, oracle, markupPercent, decimals)
}

// AddToken is a paid mutator transaction binding the contract method 0x9acad299.
//
// Solidity: function addToken(address token, address oracle, uint256 markupPercent, uint8 decimals) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) AddToken(token common.Address, oracle common.Address, markupPercent *big.Int, decimals uint8) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.AddToken(&_MultiTokenPaymasterPredeploy.TransactOpts, token, oracle, markupPercent, decimals)
}

// AddToken is a paid mutator transaction binding the contract method 0x9acad299.
//
// Solidity: function addToken(address token, address oracle, uint256 markupPercent, uint8 decimals) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) AddToken(token common.Address, oracle common.Address, markupPercent *big.Int, decimals uint8) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.AddToken(&_MultiTokenPaymasterPredeploy.TransactOpts, token, oracle, markupPercent, decimals)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) Deposit() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.Deposit(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) Deposit() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.Deposit(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _entryPoint, address _owner) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) Initialize(opts *bind.TransactOpts, _entryPoint common.Address, _owner common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "initialize", _entryPoint, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _entryPoint, address _owner) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) Initialize(_entryPoint common.Address, _owner common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.Initialize(&_MultiTokenPaymasterPredeploy.TransactOpts, _entryPoint, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _entryPoint, address _owner) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) Initialize(_entryPoint common.Address, _owner common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.Initialize(&_MultiTokenPaymasterPredeploy.TransactOpts, _entryPoint, _owner)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) PostOp(opts *bind.TransactOpts, mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) PostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.PostOp(&_MultiTokenPaymasterPredeploy.TransactOpts, mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) PostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.PostOp(&_MultiTokenPaymasterPredeploy.TransactOpts, mode, context, actualGasCost, actualUserOpFeePerGas)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address token) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) RemoveToken(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "removeToken", token)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address token) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) RemoveToken(token common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.RemoveToken(&_MultiTokenPaymasterPredeploy.TransactOpts, token)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address token) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) RemoveToken(token common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.RemoveToken(&_MultiTokenPaymasterPredeploy.TransactOpts, token)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) RenounceOwnership() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.RenounceOwnership(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.RenounceOwnership(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// SetMinCharge is a paid mutator transaction binding the contract method 0xc9a5f049.
//
// Solidity: function setMinCharge(address token, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) SetMinCharge(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "setMinCharge", token, amount)
}

// SetMinCharge is a paid mutator transaction binding the contract method 0xc9a5f049.
//
// Solidity: function setMinCharge(address token, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) SetMinCharge(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.SetMinCharge(&_MultiTokenPaymasterPredeploy.TransactOpts, token, amount)
}

// SetMinCharge is a paid mutator transaction binding the contract method 0xc9a5f049.
//
// Solidity: function setMinCharge(address token, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) SetMinCharge(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.SetMinCharge(&_MultiTokenPaymasterPredeploy.TransactOpts, token, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.TransferOwnership(&_MultiTokenPaymasterPredeploy.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.TransferOwnership(&_MultiTokenPaymasterPredeploy.TransactOpts, newOwner)
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) UnlockStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "unlockStake")
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) UnlockStake() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.UnlockStake(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) UnlockStake() (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.UnlockStake(&_MultiTokenPaymasterPredeploy.TransactOpts)
}

// UpdateTokenConfig is a paid mutator transaction binding the contract method 0x4f089d92.
//
// Solidity: function updateTokenConfig(address token, address oracle, uint256 markupPercent) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) UpdateTokenConfig(opts *bind.TransactOpts, token common.Address, oracle common.Address, markupPercent *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "updateTokenConfig", token, oracle, markupPercent)
}

// UpdateTokenConfig is a paid mutator transaction binding the contract method 0x4f089d92.
//
// Solidity: function updateTokenConfig(address token, address oracle, uint256 markupPercent) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) UpdateTokenConfig(token common.Address, oracle common.Address, markupPercent *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.UpdateTokenConfig(&_MultiTokenPaymasterPredeploy.TransactOpts, token, oracle, markupPercent)
}

// UpdateTokenConfig is a paid mutator transaction binding the contract method 0x4f089d92.
//
// Solidity: function updateTokenConfig(address token, address oracle, uint256 markupPercent) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) UpdateTokenConfig(token common.Address, oracle common.Address, markupPercent *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.UpdateTokenConfig(&_MultiTokenPaymasterPredeploy.TransactOpts, token, oracle, markupPercent)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "validatePaymasterUserOp", userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) ValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.ValidatePaymasterUserOp(&_MultiTokenPaymasterPredeploy.TransactOpts, userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) ValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.ValidatePaymasterUserOp(&_MultiTokenPaymasterPredeploy.TransactOpts, userOp, userOpHash, maxCost)
}

// WithdrawCollectedFees is a paid mutator transaction binding the contract method 0xc5241681.
//
// Solidity: function withdrawCollectedFees(address token, address to, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) WithdrawCollectedFees(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "withdrawCollectedFees", token, to, amount)
}

// WithdrawCollectedFees is a paid mutator transaction binding the contract method 0xc5241681.
//
// Solidity: function withdrawCollectedFees(address token, address to, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) WithdrawCollectedFees(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.WithdrawCollectedFees(&_MultiTokenPaymasterPredeploy.TransactOpts, token, to, amount)
}

// WithdrawCollectedFees is a paid mutator transaction binding the contract method 0xc5241681.
//
// Solidity: function withdrawCollectedFees(address token, address to, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) WithdrawCollectedFees(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.WithdrawCollectedFees(&_MultiTokenPaymasterPredeploy.TransactOpts, token, to, amount)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) WithdrawStake(opts *bind.TransactOpts, withdrawAddress common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "withdrawStake", withdrawAddress)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) WithdrawStake(withdrawAddress common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.WithdrawStake(&_MultiTokenPaymasterPredeploy.TransactOpts, withdrawAddress)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) WithdrawStake(withdrawAddress common.Address) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.WithdrawStake(&_MultiTokenPaymasterPredeploy.TransactOpts, withdrawAddress)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactor) WithdrawTo(opts *bind.TransactOpts, withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.contract.Transact(opts, "withdrawTo", withdrawAddress, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeploySession) WithdrawTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.WithdrawTo(&_MultiTokenPaymasterPredeploy.TransactOpts, withdrawAddress, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployTransactorSession) WithdrawTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MultiTokenPaymasterPredeploy.Contract.WithdrawTo(&_MultiTokenPaymasterPredeploy.TransactOpts, withdrawAddress, amount)
}

// MultiTokenPaymasterPredeployFeesCollectedIterator is returned from FilterFeesCollected and is used to iterate over the raw logs and unpacked data for FeesCollected events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployFeesCollectedIterator struct {
	Event *MultiTokenPaymasterPredeployFeesCollected // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployFeesCollectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployFeesCollected)
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
		it.Event = new(MultiTokenPaymasterPredeployFeesCollected)
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
func (it *MultiTokenPaymasterPredeployFeesCollectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployFeesCollectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployFeesCollected represents a FeesCollected event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployFeesCollected struct {
	Token  common.Address
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeesCollected is a free log retrieval operation binding the contract event 0x9bcb6d1f38f6800906185471a11ede9a8e16200853225aa62558db6076490f2d.
//
// Solidity: event FeesCollected(address indexed token, address indexed sender, uint256 amount)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterFeesCollected(opts *bind.FilterOpts, token []common.Address, sender []common.Address) (*MultiTokenPaymasterPredeployFeesCollectedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "FeesCollected", tokenRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployFeesCollectedIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "FeesCollected", logs: logs, sub: sub}, nil
}

// WatchFeesCollected is a free log subscription operation binding the contract event 0x9bcb6d1f38f6800906185471a11ede9a8e16200853225aa62558db6076490f2d.
//
// Solidity: event FeesCollected(address indexed token, address indexed sender, uint256 amount)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchFeesCollected(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployFeesCollected, token []common.Address, sender []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "FeesCollected", tokenRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployFeesCollected)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "FeesCollected", log); err != nil {
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

// ParseFeesCollected is a log parse operation binding the contract event 0x9bcb6d1f38f6800906185471a11ede9a8e16200853225aa62558db6076490f2d.
//
// Solidity: event FeesCollected(address indexed token, address indexed sender, uint256 amount)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseFeesCollected(log types.Log) (*MultiTokenPaymasterPredeployFeesCollected, error) {
	event := new(MultiTokenPaymasterPredeployFeesCollected)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "FeesCollected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiTokenPaymasterPredeployFeesWithdrawnIterator is returned from FilterFeesWithdrawn and is used to iterate over the raw logs and unpacked data for FeesWithdrawn events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployFeesWithdrawnIterator struct {
	Event *MultiTokenPaymasterPredeployFeesWithdrawn // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployFeesWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployFeesWithdrawn)
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
		it.Event = new(MultiTokenPaymasterPredeployFeesWithdrawn)
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
func (it *MultiTokenPaymasterPredeployFeesWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployFeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployFeesWithdrawn represents a FeesWithdrawn event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployFeesWithdrawn struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeesWithdrawn is a free log retrieval operation binding the contract event 0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8.
//
// Solidity: event FeesWithdrawn(address indexed token, address indexed to, uint256 amount)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterFeesWithdrawn(opts *bind.FilterOpts, token []common.Address, to []common.Address) (*MultiTokenPaymasterPredeployFeesWithdrawnIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "FeesWithdrawn", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployFeesWithdrawnIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFeesWithdrawn is a free log subscription operation binding the contract event 0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8.
//
// Solidity: event FeesWithdrawn(address indexed token, address indexed to, uint256 amount)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployFeesWithdrawn, token []common.Address, to []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "FeesWithdrawn", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployFeesWithdrawn)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
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

// ParseFeesWithdrawn is a log parse operation binding the contract event 0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8.
//
// Solidity: event FeesWithdrawn(address indexed token, address indexed to, uint256 amount)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseFeesWithdrawn(log types.Log) (*MultiTokenPaymasterPredeployFeesWithdrawn, error) {
	event := new(MultiTokenPaymasterPredeployFeesWithdrawn)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiTokenPaymasterPredeployInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployInitializedIterator struct {
	Event *MultiTokenPaymasterPredeployInitialized // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployInitialized)
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
		it.Event = new(MultiTokenPaymasterPredeployInitialized)
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
func (it *MultiTokenPaymasterPredeployInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployInitialized represents a Initialized event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterInitialized(opts *bind.FilterOpts) (*MultiTokenPaymasterPredeployInitializedIterator, error) {

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployInitializedIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployInitialized) (event.Subscription, error) {

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployInitialized)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseInitialized(log types.Log) (*MultiTokenPaymasterPredeployInitialized, error) {
	event := new(MultiTokenPaymasterPredeployInitialized)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiTokenPaymasterPredeployOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployOwnershipTransferStartedIterator struct {
	Event *MultiTokenPaymasterPredeployOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployOwnershipTransferStarted)
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
		it.Event = new(MultiTokenPaymasterPredeployOwnershipTransferStarted)
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
func (it *MultiTokenPaymasterPredeployOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MultiTokenPaymasterPredeployOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployOwnershipTransferStartedIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployOwnershipTransferStarted)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseOwnershipTransferStarted(log types.Log) (*MultiTokenPaymasterPredeployOwnershipTransferStarted, error) {
	event := new(MultiTokenPaymasterPredeployOwnershipTransferStarted)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiTokenPaymasterPredeployOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployOwnershipTransferredIterator struct {
	Event *MultiTokenPaymasterPredeployOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployOwnershipTransferred)
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
		it.Event = new(MultiTokenPaymasterPredeployOwnershipTransferred)
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
func (it *MultiTokenPaymasterPredeployOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployOwnershipTransferred represents a OwnershipTransferred event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MultiTokenPaymasterPredeployOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployOwnershipTransferredIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployOwnershipTransferred)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseOwnershipTransferred(log types.Log) (*MultiTokenPaymasterPredeployOwnershipTransferred, error) {
	event := new(MultiTokenPaymasterPredeployOwnershipTransferred)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiTokenPaymasterPredeployTokenAddedIterator is returned from FilterTokenAdded and is used to iterate over the raw logs and unpacked data for TokenAdded events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployTokenAddedIterator struct {
	Event *MultiTokenPaymasterPredeployTokenAdded // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployTokenAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployTokenAdded)
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
		it.Event = new(MultiTokenPaymasterPredeployTokenAdded)
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
func (it *MultiTokenPaymasterPredeployTokenAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployTokenAdded represents a TokenAdded event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployTokenAdded struct {
	Token    common.Address
	Oracle   common.Address
	Markup   *big.Int
	Decimals uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenAdded is a free log retrieval operation binding the contract event 0x25222b0b8329cba9b0556bf83ca70763ce37f735819a0d3c8690b04a1bcf77b7.
//
// Solidity: event TokenAdded(address indexed token, address oracle, uint256 markup, uint8 decimals)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterTokenAdded(opts *bind.FilterOpts, token []common.Address) (*MultiTokenPaymasterPredeployTokenAddedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "TokenAdded", tokenRule)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployTokenAddedIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "TokenAdded", logs: logs, sub: sub}, nil
}

// WatchTokenAdded is a free log subscription operation binding the contract event 0x25222b0b8329cba9b0556bf83ca70763ce37f735819a0d3c8690b04a1bcf77b7.
//
// Solidity: event TokenAdded(address indexed token, address oracle, uint256 markup, uint8 decimals)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchTokenAdded(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployTokenAdded, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "TokenAdded", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployTokenAdded)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "TokenAdded", log); err != nil {
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

// ParseTokenAdded is a log parse operation binding the contract event 0x25222b0b8329cba9b0556bf83ca70763ce37f735819a0d3c8690b04a1bcf77b7.
//
// Solidity: event TokenAdded(address indexed token, address oracle, uint256 markup, uint8 decimals)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseTokenAdded(log types.Log) (*MultiTokenPaymasterPredeployTokenAdded, error) {
	event := new(MultiTokenPaymasterPredeployTokenAdded)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "TokenAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiTokenPaymasterPredeployTokenConfigUpdatedIterator is returned from FilterTokenConfigUpdated and is used to iterate over the raw logs and unpacked data for TokenConfigUpdated events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployTokenConfigUpdatedIterator struct {
	Event *MultiTokenPaymasterPredeployTokenConfigUpdated // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployTokenConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployTokenConfigUpdated)
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
		it.Event = new(MultiTokenPaymasterPredeployTokenConfigUpdated)
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
func (it *MultiTokenPaymasterPredeployTokenConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployTokenConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployTokenConfigUpdated represents a TokenConfigUpdated event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployTokenConfigUpdated struct {
	Token  common.Address
	Oracle common.Address
	Markup *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenConfigUpdated is a free log retrieval operation binding the contract event 0xae75ec539c80c1222e8c675cfb7f98a0dca2684fa410fe182186cccda491eef1.
//
// Solidity: event TokenConfigUpdated(address indexed token, address oracle, uint256 markup)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterTokenConfigUpdated(opts *bind.FilterOpts, token []common.Address) (*MultiTokenPaymasterPredeployTokenConfigUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "TokenConfigUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployTokenConfigUpdatedIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "TokenConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchTokenConfigUpdated is a free log subscription operation binding the contract event 0xae75ec539c80c1222e8c675cfb7f98a0dca2684fa410fe182186cccda491eef1.
//
// Solidity: event TokenConfigUpdated(address indexed token, address oracle, uint256 markup)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchTokenConfigUpdated(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployTokenConfigUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "TokenConfigUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployTokenConfigUpdated)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "TokenConfigUpdated", log); err != nil {
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

// ParseTokenConfigUpdated is a log parse operation binding the contract event 0xae75ec539c80c1222e8c675cfb7f98a0dca2684fa410fe182186cccda491eef1.
//
// Solidity: event TokenConfigUpdated(address indexed token, address oracle, uint256 markup)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseTokenConfigUpdated(log types.Log) (*MultiTokenPaymasterPredeployTokenConfigUpdated, error) {
	event := new(MultiTokenPaymasterPredeployTokenConfigUpdated)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "TokenConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiTokenPaymasterPredeployTokenRemovedIterator is returned from FilterTokenRemoved and is used to iterate over the raw logs and unpacked data for TokenRemoved events raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployTokenRemovedIterator struct {
	Event *MultiTokenPaymasterPredeployTokenRemoved // Event containing the contract specifics and raw log

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
func (it *MultiTokenPaymasterPredeployTokenRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTokenPaymasterPredeployTokenRemoved)
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
		it.Event = new(MultiTokenPaymasterPredeployTokenRemoved)
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
func (it *MultiTokenPaymasterPredeployTokenRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTokenPaymasterPredeployTokenRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTokenPaymasterPredeployTokenRemoved represents a TokenRemoved event raised by the MultiTokenPaymasterPredeploy contract.
type MultiTokenPaymasterPredeployTokenRemoved struct {
	Token common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTokenRemoved is a free log retrieval operation binding the contract event 0x4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd3.
//
// Solidity: event TokenRemoved(address indexed token)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) FilterTokenRemoved(opts *bind.FilterOpts, token []common.Address) (*MultiTokenPaymasterPredeployTokenRemovedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.FilterLogs(opts, "TokenRemoved", tokenRule)
	if err != nil {
		return nil, err
	}
	return &MultiTokenPaymasterPredeployTokenRemovedIterator{contract: _MultiTokenPaymasterPredeploy.contract, event: "TokenRemoved", logs: logs, sub: sub}, nil
}

// WatchTokenRemoved is a free log subscription operation binding the contract event 0x4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd3.
//
// Solidity: event TokenRemoved(address indexed token)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) WatchTokenRemoved(opts *bind.WatchOpts, sink chan<- *MultiTokenPaymasterPredeployTokenRemoved, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultiTokenPaymasterPredeploy.contract.WatchLogs(opts, "TokenRemoved", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTokenPaymasterPredeployTokenRemoved)
				if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "TokenRemoved", log); err != nil {
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

// ParseTokenRemoved is a log parse operation binding the contract event 0x4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd3.
//
// Solidity: event TokenRemoved(address indexed token)
func (_MultiTokenPaymasterPredeploy *MultiTokenPaymasterPredeployFilterer) ParseTokenRemoved(log types.Log) (*MultiTokenPaymasterPredeployTokenRemoved, error) {
	event := new(MultiTokenPaymasterPredeployTokenRemoved)
	if err := _MultiTokenPaymasterPredeploy.contract.UnpackLog(event, "TokenRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
