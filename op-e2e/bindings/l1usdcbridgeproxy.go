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

// L1UsdcBridgeProxyMetaData contains all meta data concerning the L1UsdcBridgeProxy contract.
var L1UsdcBridgeProxyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_logic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"payable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposits\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l1Usdc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l2Usdc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messenger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"otherBridge\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxyChangeOwner\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAddress\",\"inputs\":[{\"name\":\"_messenger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_otherBridge\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_l1Usdc\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_l2Usdc\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidAdmin\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]}]",
	Bin: "0x6080604052604051620013b6380380620013b68339810160408190526200002691620003e9565b82816200003f82826200006060201b620008431760201c565b50506200005782620000d160201b620008ab1760201c565b505050620004e7565b6200006b826200012c565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115620000c357620000be8282620001c360201b6200090c1760201c565b505050565b620000cd62000240565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620000fc62000262565b604080516001600160a01b03928316815291841660208301520160405180910390a162000129816200029b565b50565b806001600160a01b03163b6000036200016857604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b80620001a27f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc60001b620002f060201b6200098f1760201c565b80546001600160a01b0319166001600160a01b039290921691909117905550565b6060600080846001600160a01b031684604051620001e29190620004c9565b600060405180830381855af49150503d80600081146200021f576040519150601f19603f3d011682016040523d82523d6000602084013e62000224565b606091505b50909250905062000237858383620002f3565b95945050505050565b3415620002605760405163b398979f60e01b815260040160405180910390fd5b565b60006200028c6000805160206200139683398151915260001b620002f060201b6200098f1760201c565b546001600160a01b0316919050565b6001600160a01b038116620002c757604051633173bdd160e11b8152600060048201526024016200015f565b80620001a26000805160206200139683398151915260001b620002f060201b6200098f1760201c565b90565b6060826200030c57620003068262000359565b62000352565b81511580156200032457506001600160a01b0384163b155b156200034f57604051639996b31560e01b81526001600160a01b03851660048201526024016200015f565b50805b9392505050565b8051156200036a5780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80516001600160a01b03811681146200039b57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620003d3578181015183820152602001620003b9565b83811115620003e3576000848401525b50505050565b600080600060608486031215620003ff57600080fd5b6200040a8462000383565b92506200041a6020850162000383565b60408501519092506001600160401b03808211156200043857600080fd5b818601915086601f8301126200044d57600080fd5b815181811115620004625762000462620003a0565b604051601f8201601f19908116603f011681019083821181831017156200048d576200048d620003a0565b81604052828152896020848701011115620004a757600080fd5b620004ba836020830160208801620003b6565b80955050505050509250925092565b60008251620004dd818460208701620003b6565b9190910192915050565b610e9f80620004f76000396000f3fe6080604052600436106100c05760003560e01c80638da5cb5b11610074578063a1b4bc041161004e578063a1b4bc041461028a578063c89701a2146102b7578063dfd3dcb3146102e45761012c565b80638da5cb5b1461020f5780638f601f66146102245780639608088c1461026a5761012c565b80634f1ef286116100a55780634f1ef286146101ad57806356c3b587146101cd5780635c60da1b146101fa5761012c565b80633659cfe6146101365780633cb747bf146101565761012c565b3661012c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f63616e6e6f74207265636569766520457468657200000000000000000000000060448201526064015b60405180910390fd5b610134610304565b005b34801561014257600080fd5b50610134610151366004610ca6565b610316565b34801561016257600080fd5b506000546101839073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b3480156101b957600080fd5b506101346101c8366004610cf0565b6103ce565b3480156101d957600080fd5b506002546101839073ffffffffffffffffffffffffffffffffffffffff1681565b34801561020657600080fd5b50610183610478565b34801561021b57600080fd5b50610183610487565b34801561023057600080fd5b5061025c61023f366004610dd0565b600460209081526000928352604080842090915290825290205481565b6040519081526020016101a4565b34801561027657600080fd5b50610134610285366004610e03565b610491565b34801561029657600080fd5b506003546101839073ffffffffffffffffffffffffffffffffffffffff1681565b3480156102c357600080fd5b506001546101839073ffffffffffffffffffffffffffffffffffffffff1681565b3480156102f057600080fd5b506101346102ff366004610ca6565b61079e565b61031461030f610992565b61099c565b565b61031e610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103b2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b6103cb8160405180602001604052806000815250610843565b50565b6103d6610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461046a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b6104748282610843565b5050565b6000610482610992565b905090565b60006104826109c0565b610499610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461052d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff81166105ab576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff8116610629576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff81166106a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff8116610725576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b50506000805473ffffffffffffffffffffffffffffffffffffffff9788167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560018054968816968216969096179095555050600280549285169284169290921790915560038054919093169116179055565b6107a6610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461083a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b6103cb816108ab565b61084c82610a00565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156108a35761089e828261090c565b505050565b610474610ad2565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6108d46109c0565b6040805173ffffffffffffffffffffffffffffffffffffffff928316815291841660208301520160405180910390a16103cb81610b0a565b60606000808473ffffffffffffffffffffffffffffffffffffffff16846040516109369190610e57565b600060405180830381855af49150503d8060008114610971576040519150601f19603f3d011682016040523d82523d6000602084013e610976565b606091505b5091509150610986858383610b81565b95945050505050565b90565b6000610482610c13565b3660008037600080366000845af43d6000803e8080156109bb573d6000f35b3d6000fd5b60007fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b5473ffffffffffffffffffffffffffffffffffffffff16919050565b8073ffffffffffffffffffffffffffffffffffffffff163b600003610a69576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610123565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905550565b3415610314576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116610b5a576040517f62e77ba200000000000000000000000000000000000000000000000000000000815260006004820152602401610123565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103610a8c565b606082610b9657610b9182610c3b565b610c0c565b8151158015610bba575073ffffffffffffffffffffffffffffffffffffffff84163b155b15610c09576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602401610123565b50805b9392505050565b60007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6109e4565b805115610c4b5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b803573ffffffffffffffffffffffffffffffffffffffff81168114610ca157600080fd5b919050565b600060208284031215610cb857600080fd5b610c0c82610c7d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008060408385031215610d0357600080fd5b610d0c83610c7d565b9150602083013567ffffffffffffffff80821115610d2957600080fd5b818501915085601f830112610d3d57600080fd5b813581811115610d4f57610d4f610cc1565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610d9557610d95610cc1565b81604052828152886020848701011115610dae57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60008060408385031215610de357600080fd5b610dec83610c7d565b9150610dfa60208401610c7d565b90509250929050565b60008060008060808587031215610e1957600080fd5b610e2285610c7d565b9350610e3060208601610c7d565b9250610e3e60408601610c7d565b9150610e4c60608601610c7d565b905092959194509250565b6000825160005b81811015610e785760208186018101518583015201610e5e565b81811115610e87576000828501525b50919091019291505056fea164736f6c634300080f000ab53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103",
}

// L1UsdcBridgeProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use L1UsdcBridgeProxyMetaData.ABI instead.
var L1UsdcBridgeProxyABI = L1UsdcBridgeProxyMetaData.ABI

// L1UsdcBridgeProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L1UsdcBridgeProxyMetaData.Bin instead.
var L1UsdcBridgeProxyBin = L1UsdcBridgeProxyMetaData.Bin

// DeployL1UsdcBridgeProxy deploys a new Ethereum contract, binding an instance of L1UsdcBridgeProxy to it.
func DeployL1UsdcBridgeProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, initialOwner common.Address, _data []byte) (common.Address, *types.Transaction, *L1UsdcBridgeProxy, error) {
	parsed, err := L1UsdcBridgeProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L1UsdcBridgeProxyBin), backend, _logic, initialOwner, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L1UsdcBridgeProxy{L1UsdcBridgeProxyCaller: L1UsdcBridgeProxyCaller{contract: contract}, L1UsdcBridgeProxyTransactor: L1UsdcBridgeProxyTransactor{contract: contract}, L1UsdcBridgeProxyFilterer: L1UsdcBridgeProxyFilterer{contract: contract}}, nil
}

// L1UsdcBridgeProxy is an auto generated Go binding around an Ethereum contract.
type L1UsdcBridgeProxy struct {
	L1UsdcBridgeProxyCaller     // Read-only binding to the contract
	L1UsdcBridgeProxyTransactor // Write-only binding to the contract
	L1UsdcBridgeProxyFilterer   // Log filterer for contract events
}

// L1UsdcBridgeProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1UsdcBridgeProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1UsdcBridgeProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1UsdcBridgeProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1UsdcBridgeProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1UsdcBridgeProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1UsdcBridgeProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1UsdcBridgeProxySession struct {
	Contract     *L1UsdcBridgeProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// L1UsdcBridgeProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1UsdcBridgeProxyCallerSession struct {
	Contract *L1UsdcBridgeProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// L1UsdcBridgeProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1UsdcBridgeProxyTransactorSession struct {
	Contract     *L1UsdcBridgeProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// L1UsdcBridgeProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1UsdcBridgeProxyRaw struct {
	Contract *L1UsdcBridgeProxy // Generic contract binding to access the raw methods on
}

// L1UsdcBridgeProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1UsdcBridgeProxyCallerRaw struct {
	Contract *L1UsdcBridgeProxyCaller // Generic read-only contract binding to access the raw methods on
}

// L1UsdcBridgeProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1UsdcBridgeProxyTransactorRaw struct {
	Contract *L1UsdcBridgeProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1UsdcBridgeProxy creates a new instance of L1UsdcBridgeProxy, bound to a specific deployed contract.
func NewL1UsdcBridgeProxy(address common.Address, backend bind.ContractBackend) (*L1UsdcBridgeProxy, error) {
	contract, err := bindL1UsdcBridgeProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeProxy{L1UsdcBridgeProxyCaller: L1UsdcBridgeProxyCaller{contract: contract}, L1UsdcBridgeProxyTransactor: L1UsdcBridgeProxyTransactor{contract: contract}, L1UsdcBridgeProxyFilterer: L1UsdcBridgeProxyFilterer{contract: contract}}, nil
}

// NewL1UsdcBridgeProxyCaller creates a new read-only instance of L1UsdcBridgeProxy, bound to a specific deployed contract.
func NewL1UsdcBridgeProxyCaller(address common.Address, caller bind.ContractCaller) (*L1UsdcBridgeProxyCaller, error) {
	contract, err := bindL1UsdcBridgeProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeProxyCaller{contract: contract}, nil
}

// NewL1UsdcBridgeProxyTransactor creates a new write-only instance of L1UsdcBridgeProxy, bound to a specific deployed contract.
func NewL1UsdcBridgeProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*L1UsdcBridgeProxyTransactor, error) {
	contract, err := bindL1UsdcBridgeProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeProxyTransactor{contract: contract}, nil
}

// NewL1UsdcBridgeProxyFilterer creates a new log filterer instance of L1UsdcBridgeProxy, bound to a specific deployed contract.
func NewL1UsdcBridgeProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*L1UsdcBridgeProxyFilterer, error) {
	contract, err := bindL1UsdcBridgeProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeProxyFilterer{contract: contract}, nil
}

// bindL1UsdcBridgeProxy binds a generic wrapper to an already deployed contract.
func bindL1UsdcBridgeProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(L1UsdcBridgeProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1UsdcBridgeProxy.Contract.L1UsdcBridgeProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.L1UsdcBridgeProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.L1UsdcBridgeProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1UsdcBridgeProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.contract.Transact(opts, method, params...)
}

// Deposits is a free data retrieval call binding the contract method 0x8f601f66.
//
// Solidity: function deposits(address , address ) view returns(uint256)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCaller) Deposits(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _L1UsdcBridgeProxy.contract.Call(opts, &out, "deposits", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Deposits is a free data retrieval call binding the contract method 0x8f601f66.
//
// Solidity: function deposits(address , address ) view returns(uint256)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) Deposits(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _L1UsdcBridgeProxy.Contract.Deposits(&_L1UsdcBridgeProxy.CallOpts, arg0, arg1)
}

// Deposits is a free data retrieval call binding the contract method 0x8f601f66.
//
// Solidity: function deposits(address , address ) view returns(uint256)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerSession) Deposits(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _L1UsdcBridgeProxy.Contract.Deposits(&_L1UsdcBridgeProxy.CallOpts, arg0, arg1)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridgeProxy.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) Implementation() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.Implementation(&_L1UsdcBridgeProxy.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerSession) Implementation() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.Implementation(&_L1UsdcBridgeProxy.CallOpts)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCaller) L1Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridgeProxy.contract.Call(opts, &out, "l1Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) L1Usdc() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.L1Usdc(&_L1UsdcBridgeProxy.CallOpts)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerSession) L1Usdc() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.L1Usdc(&_L1UsdcBridgeProxy.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCaller) L2Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridgeProxy.contract.Call(opts, &out, "l2Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) L2Usdc() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.L2Usdc(&_L1UsdcBridgeProxy.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerSession) L2Usdc() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.L2Usdc(&_L1UsdcBridgeProxy.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCaller) Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridgeProxy.contract.Call(opts, &out, "messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) Messenger() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.Messenger(&_L1UsdcBridgeProxy.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerSession) Messenger() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.Messenger(&_L1UsdcBridgeProxy.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCaller) OtherBridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridgeProxy.contract.Call(opts, &out, "otherBridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) OtherBridge() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.OtherBridge(&_L1UsdcBridgeProxy.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerSession) OtherBridge() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.OtherBridge(&_L1UsdcBridgeProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridgeProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) Owner() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.Owner(&_L1UsdcBridgeProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyCallerSession) Owner() (common.Address, error) {
	return _L1UsdcBridgeProxy.Contract.Owner(&_L1UsdcBridgeProxy.CallOpts)
}

// ProxyChangeOwner is a paid mutator transaction binding the contract method 0xdfd3dcb3.
//
// Solidity: function proxyChangeOwner(address newAdmin) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactor) ProxyChangeOwner(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.contract.Transact(opts, "proxyChangeOwner", newAdmin)
}

// ProxyChangeOwner is a paid mutator transaction binding the contract method 0xdfd3dcb3.
//
// Solidity: function proxyChangeOwner(address newAdmin) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) ProxyChangeOwner(newAdmin common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.ProxyChangeOwner(&_L1UsdcBridgeProxy.TransactOpts, newAdmin)
}

// ProxyChangeOwner is a paid mutator transaction binding the contract method 0xdfd3dcb3.
//
// Solidity: function proxyChangeOwner(address newAdmin) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorSession) ProxyChangeOwner(newAdmin common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.ProxyChangeOwner(&_L1UsdcBridgeProxy.TransactOpts, newAdmin)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9608088c.
//
// Solidity: function setAddress(address _messenger, address _otherBridge, address _l1Usdc, address _l2Usdc) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactor) SetAddress(opts *bind.TransactOpts, _messenger common.Address, _otherBridge common.Address, _l1Usdc common.Address, _l2Usdc common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.contract.Transact(opts, "setAddress", _messenger, _otherBridge, _l1Usdc, _l2Usdc)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9608088c.
//
// Solidity: function setAddress(address _messenger, address _otherBridge, address _l1Usdc, address _l2Usdc) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) SetAddress(_messenger common.Address, _otherBridge common.Address, _l1Usdc common.Address, _l2Usdc common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.SetAddress(&_L1UsdcBridgeProxy.TransactOpts, _messenger, _otherBridge, _l1Usdc, _l2Usdc)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9608088c.
//
// Solidity: function setAddress(address _messenger, address _otherBridge, address _l1Usdc, address _l2Usdc) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorSession) SetAddress(_messenger common.Address, _otherBridge common.Address, _l1Usdc common.Address, _l2Usdc common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.SetAddress(&_L1UsdcBridgeProxy.TransactOpts, _messenger, _otherBridge, _l1Usdc, _l2Usdc)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.UpgradeTo(&_L1UsdcBridgeProxy.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.UpgradeTo(&_L1UsdcBridgeProxy.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.UpgradeToAndCall(&_L1UsdcBridgeProxy.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.UpgradeToAndCall(&_L1UsdcBridgeProxy.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.Fallback(&_L1UsdcBridgeProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.Fallback(&_L1UsdcBridgeProxy.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxySession) Receive() (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.Receive(&_L1UsdcBridgeProxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _L1UsdcBridgeProxy.Contract.Receive(&_L1UsdcBridgeProxy.TransactOpts)
}
