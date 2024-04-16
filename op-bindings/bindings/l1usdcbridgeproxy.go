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

// L1UsdcBridgeProxyMetaData contains all meta data concerning the L1UsdcBridgeProxy contract.
var L1UsdcBridgeProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1Usdc\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2Usdc\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messenger\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"otherBridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"proxyChangeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_messenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_otherBridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l1Usdc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l2Usdc\",\"type\":\"address\"}],\"name\":\"setAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040526040516200132c3803806200132c83398101604081905262000026916200036e565b82816200003482826200004b565b5062000042905082620000b1565b5050506200046c565b620000568262000123565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115620000a3576200009e8282620001a3565b505050565b620000ad62000220565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620000f36000805160206200130c833981519152546001600160a01b031690565b604080516001600160a01b03928316815291841660208301520160405180910390a1620001208162000242565b50565b806001600160a01b03163b6000036200015f57604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b6060600080846001600160a01b031684604051620001c291906200044e565b600060405180830381855af49150503d8060008114620001ff576040519150601f19603f3d011682016040523d82523d6000602084013e62000204565b606091505b5090925090506200021785838362000285565b95945050505050565b3415620002405760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b0381166200026e57604051633173bdd160e11b81526000600482015260240162000156565b806000805160206200130c83398151915262000182565b6060826200029e576200029882620002eb565b620002e4565b8151158015620002b657506001600160a01b0384163b155b15620002e157604051639996b31560e01b81526001600160a01b038516600482015260240162000156565b50805b9392505050565b805115620002fc5780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80516001600160a01b03811681146200032d57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620003655781810151838201526020016200034b565b50506000910152565b6000806000606084860312156200038457600080fd5b6200038f8462000315565b92506200039f6020850162000315565b60408501519092506001600160401b0380821115620003bd57600080fd5b818601915086601f830112620003d257600080fd5b815181811115620003e757620003e762000332565b604051601f8201601f19908116603f0116810190838211818310171562000412576200041262000332565b816040528281528960208487010111156200042c57600080fd5b6200043f83602083016020880162000348565b80955050505050509250925092565b600082516200046281846020870162000348565b9190910192915050565b610e90806200047c6000396000f3fe6080604052600436106100c05760003560e01c80638da5cb5b11610074578063a1b4bc041161004e578063a1b4bc041461028a578063c89701a2146102b7578063dfd3dcb3146102e45761012c565b80638da5cb5b1461020f5780638f601f66146102245780639608088c1461026a5761012c565b80634f1ef286116100a55780634f1ef286146101ad57806356c3b587146101cd5780635c60da1b146101fa5761012c565b80633659cfe6146101365780633cb747bf146101565761012c565b3661012c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f63616e6e6f74207265636569766520457468657200000000000000000000000060448201526064015b60405180910390fd5b610134610304565b005b34801561014257600080fd5b50610134610151366004610ca3565b610316565b34801561016257600080fd5b506000546101839073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b3480156101b957600080fd5b506101346101c8366004610ced565b6103ce565b3480156101d957600080fd5b506002546101839073ffffffffffffffffffffffffffffffffffffffff1681565b34801561020657600080fd5b50610183610478565b34801561021b57600080fd5b50610183610487565b34801561023057600080fd5b5061025c61023f366004610dcd565b600460209081526000928352604080842090915290825290205481565b6040519081526020016101a4565b34801561027657600080fd5b50610134610285366004610e00565b610491565b34801561029657600080fd5b506003546101839073ffffffffffffffffffffffffffffffffffffffff1681565b3480156102c357600080fd5b506001546101839073ffffffffffffffffffffffffffffffffffffffff1681565b3480156102f057600080fd5b506101346102ff366004610ca3565b61079e565b61031461030f610843565b61084d565b565b61031e610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103b2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b6103cb8160405180602001604052806000815250610871565b50565b6103d6610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461046a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b6104748282610871565b5050565b6000610482610843565b905090565b60006104826108d9565b610499610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461052d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff81166105ab576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff8116610629576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff81166106a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b8373ffffffffffffffffffffffffffffffffffffffff8116610725576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610123565b50506000805473ffffffffffffffffffffffffffffffffffffffff9788167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560018054968816968216969096179095555050600280549285169284169290921790915560038054919093169116179055565b6107a6610487565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461083a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610123565b6103cb81610919565b600061048261097a565b3660008037600080366000845af43d6000803e80801561086c573d6000f35b3d6000fd5b61087a826109a2565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156108d1576108cc8282610a74565b505050565b610474610af7565b60007fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b5473ffffffffffffffffffffffffffffffffffffffff16919050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6109426108d9565b6040805173ffffffffffffffffffffffffffffffffffffffff928316815291841660208301520160405180910390a16103cb81610b2f565b60007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6108fd565b8073ffffffffffffffffffffffffffffffffffffffff163b600003610a0b576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610123565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905550565b60606000808473ffffffffffffffffffffffffffffffffffffffff1684604051610a9e9190610e54565b600060405180830381855af49150503d8060008114610ad9576040519150601f19603f3d011682016040523d82523d6000602084013e610ade565b606091505b5091509150610aee858383610ba6565b95945050505050565b3415610314576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116610b7f576040517f62e77ba200000000000000000000000000000000000000000000000000000000815260006004820152602401610123565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103610a2e565b606082610bbb57610bb682610c38565b610c31565b8151158015610bdf575073ffffffffffffffffffffffffffffffffffffffff84163b155b15610c2e576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602401610123565b50805b9392505050565b805115610c485780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b803573ffffffffffffffffffffffffffffffffffffffff81168114610c9e57600080fd5b919050565b600060208284031215610cb557600080fd5b610c3182610c7a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008060408385031215610d0057600080fd5b610d0983610c7a565b9150602083013567ffffffffffffffff80821115610d2657600080fd5b818501915085601f830112610d3a57600080fd5b813581811115610d4c57610d4c610cbe565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610d9257610d92610cbe565b81604052828152886020848701011115610dab57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60008060408385031215610de057600080fd5b610de983610c7a565b9150610df760208401610c7a565b90509250929050565b60008060008060808587031215610e1657600080fd5b610e1f85610c7a565b9350610e2d60208601610c7a565b9250610e3b60408601610c7a565b9150610e4960608601610c7a565b905092959194509250565b6000825160005b81811015610e755760208186018101518583015201610e5b565b50600092019182525091905056fea164736f6c6343000814000ab53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103",
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
	parsed, err := L1UsdcBridgeProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// L1UsdcBridgeProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the L1UsdcBridgeProxy contract.
type L1UsdcBridgeProxyAdminChangedIterator struct {
	Event *L1UsdcBridgeProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *L1UsdcBridgeProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1UsdcBridgeProxyAdminChanged)
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
		it.Event = new(L1UsdcBridgeProxyAdminChanged)
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
func (it *L1UsdcBridgeProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1UsdcBridgeProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1UsdcBridgeProxyAdminChanged represents a AdminChanged event raised by the L1UsdcBridgeProxy contract.
type L1UsdcBridgeProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*L1UsdcBridgeProxyAdminChangedIterator, error) {

	logs, sub, err := _L1UsdcBridgeProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeProxyAdminChangedIterator{contract: _L1UsdcBridgeProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *L1UsdcBridgeProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _L1UsdcBridgeProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1UsdcBridgeProxyAdminChanged)
				if err := _L1UsdcBridgeProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyFilterer) ParseAdminChanged(log types.Log) (*L1UsdcBridgeProxyAdminChanged, error) {
	event := new(L1UsdcBridgeProxyAdminChanged)
	if err := _L1UsdcBridgeProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1UsdcBridgeProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the L1UsdcBridgeProxy contract.
type L1UsdcBridgeProxyUpgradedIterator struct {
	Event *L1UsdcBridgeProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *L1UsdcBridgeProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1UsdcBridgeProxyUpgraded)
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
		it.Event = new(L1UsdcBridgeProxyUpgraded)
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
func (it *L1UsdcBridgeProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1UsdcBridgeProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1UsdcBridgeProxyUpgraded represents a Upgraded event raised by the L1UsdcBridgeProxy contract.
type L1UsdcBridgeProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*L1UsdcBridgeProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L1UsdcBridgeProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeProxyUpgradedIterator{contract: _L1UsdcBridgeProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *L1UsdcBridgeProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L1UsdcBridgeProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1UsdcBridgeProxyUpgraded)
				if err := _L1UsdcBridgeProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_L1UsdcBridgeProxy *L1UsdcBridgeProxyFilterer) ParseUpgraded(log types.Log) (*L1UsdcBridgeProxyUpgraded, error) {
	event := new(L1UsdcBridgeProxyUpgraded)
	if err := _L1UsdcBridgeProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
