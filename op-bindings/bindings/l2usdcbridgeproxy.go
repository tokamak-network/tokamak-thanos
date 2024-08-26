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

// L2UsdcBridgeProxyMetaData contains all meta data concerning the L2UsdcBridgeProxy contract.
var L2UsdcBridgeProxyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_logic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"payable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l1Usdc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l2Usdc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"l2UsdcMasterMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messenger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"otherBridge\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxyChangeOwner\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAddress\",\"inputs\":[{\"name\":\"_messenger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_otherBridge\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_l1Usdc\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_l2Usdc\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_l2UsdcMasterMinter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeTo\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminChanged\",\"inputs\":[{\"name\":\"previousAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidAdmin\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]}]",
	Bin: "0x60806040526040516200133538038062001335833981016040819052620000269162000360565b82816200003482826200004b565b5062000042905082620000b0565b50505062000457565b620000568262000121565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115620000a2576200009d8282620001a0565b505050565b620000ac62000219565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620000f15f8051602062001315833981519152546001600160a01b031690565b604080516001600160a01b03928316815291841660208301520160405180910390a16200011e816200023b565b50565b806001600160a01b03163b5f036200015c57604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b60605f80846001600160a01b031684604051620001be91906200043a565b5f60405180830381855af49150503d805f8114620001f8576040519150601f19603f3d011682016040523d82523d5f602084013e620001fd565b606091505b509092509050620002108583836200027c565b95945050505050565b3415620002395760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b0381166200026657604051633173bdd160e11b81525f600482015260240162000153565b805f80516020620013158339815191526200017f565b60608262000295576200028f82620002e2565b620002db565b8151158015620002ad57506001600160a01b0384163b155b15620002d857604051639996b31560e01b81526001600160a01b038516600482015260240162000153565b50805b9392505050565b805115620002f35780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80516001600160a01b038116811462000323575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f5b83811015620003585781810151838201526020016200033e565b50505f910152565b5f805f6060848603121562000373575f80fd5b6200037e846200030c565b92506200038e602085016200030c565b60408501519092506001600160401b0380821115620003ab575f80fd5b818601915086601f830112620003bf575f80fd5b815181811115620003d457620003d462000328565b604051601f8201601f19908116603f01168101908382118183101715620003ff57620003ff62000328565b8160405282815289602084870101111562000418575f80fd5b6200042b8360208301602088016200033c565b80955050505050509250925092565b5f82516200044d8184602087016200033c565b9190910192915050565b610eb080620004655f395ff3fe6080604052600436106100bf575f3560e01c806356c3b58711610073578063a1b4bc041161004d578063a1b4bc0414610266578063c89701a214610292578063dfd3dcb3146102be5761012b565b806356c3b587146102125780635c60da1b1461023e5780638da5cb5b146102525761012b565b80633659cfe6116100a45780633659cfe6146101a95780633cb747bf146101c85780634f1ef286146101f35761012b565b806305db940f1461013557806312c594881461018a5761012b565b3661012b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f63616e6e6f74207265636569766520544f4e000000000000000000000000000060448201526064015b60405180910390fd5b6101336102dd565b005b348015610140575f80fd5b506004546101619073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b348015610195575f80fd5b506101336101a4366004610cf6565b6102ef565b3480156101b4575f80fd5b506101336101c3366004610d57565b61068a565b3480156101d3575f80fd5b505f546101619073ffffffffffffffffffffffffffffffffffffffff1681565b3480156101fe575f80fd5b5061013361020d366004610d9d565b610741565b34801561021d575f80fd5b506002546101619073ffffffffffffffffffffffffffffffffffffffff1681565b348015610249575f80fd5b506101616107eb565b34801561025d575f80fd5b506101616107f9565b348015610271575f80fd5b506003546101619073ffffffffffffffffffffffffffffffffffffffff1681565b34801561029d575f80fd5b506001546101619073ffffffffffffffffffffffffffffffffffffffff1681565b3480156102c9575f80fd5b506101336102d8366004610d57565b610802565b6102ed6102e86108a7565b6108b0565b565b6102f76107f9565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461038b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610122565b8473ffffffffffffffffffffffffffffffffffffffff8116610409576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610122565b8473ffffffffffffffffffffffffffffffffffffffff8116610487576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610122565b8473ffffffffffffffffffffffffffffffffffffffff8116610505576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610122565b8473ffffffffffffffffffffffffffffffffffffffff8116610583576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610122565b8473ffffffffffffffffffffffffffffffffffffffff8116610601576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7a65726f206164647265737300000000000000000000000000000000000000006044820152606401610122565b50505f80547fffffffffffffffffffffffff000000000000000000000000000000000000000090811673ffffffffffffffffffffffffffffffffffffffff9a8b1617909155600180548216988a16989098179097555050600280548616948716949094179093556003805485169286169290921790915560048054909316931692909217905550565b6106926107f9565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610726576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610122565b61073e8160405180602001604052805f8152506108ce565b50565b6107496107f9565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146107dd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610122565b6107e782826108ce565b5050565b5f6107f46108a7565b905090565b5f6107f4610935565b61080a6107f9565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461089e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e6f74206f776e657200000000000000000000000000000000000000000000006044820152606401610122565b61073e81610974565b5f6107f46109d5565b365f80375f80365f845af43d5f803e8080156108ca573d5ff35b3d5ffd5b6108d7826109fc565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a280511561092d576109288282610acd565b505050565b6107e7610b4c565b5f7fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b5473ffffffffffffffffffffffffffffffffffffffff16919050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f61099d610935565b6040805173ffffffffffffffffffffffffffffffffffffffff928316815291841660208301520160405180910390a161073e81610b84565b5f7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610958565b8073ffffffffffffffffffffffffffffffffffffffff163b5f03610a64576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610122565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905550565b60605f808473ffffffffffffffffffffffffffffffffffffffff1684604051610af69190610e77565b5f60405180830381855af49150503d805f8114610b2e576040519150601f19603f3d011682016040523d82523d5f602084013e610b33565b606091505b5091509150610b43858383610bfa565b95945050505050565b34156102ed576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116610bd3576040517f62e77ba20000000000000000000000000000000000000000000000000000000081525f6004820152602401610122565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103610a87565b606082610c0f57610c0a82610c8c565b610c85565b8151158015610c33575073ffffffffffffffffffffffffffffffffffffffff84163b155b15610c82576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602401610122565b50805b9392505050565b805115610c9c5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b803573ffffffffffffffffffffffffffffffffffffffff81168114610cf1575f80fd5b919050565b5f805f805f60a08688031215610d0a575f80fd5b610d1386610cce565b9450610d2160208701610cce565b9350610d2f60408701610cce565b9250610d3d60608701610cce565b9150610d4b60808701610cce565b90509295509295909350565b5f60208284031215610d67575f80fd5b610c8582610cce565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f8060408385031215610dae575f80fd5b610db783610cce565b9150602083013567ffffffffffffffff80821115610dd3575f80fd5b818501915085601f830112610de6575f80fd5b813581811115610df857610df8610d70565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610e3e57610e3e610d70565b81604052828152886020848701011115610e56575f80fd5b826020860160208301375f6020848301015280955050505050509250929050565b5f82515f5b81811015610e965760208186018101518583015201610e7c565b505f92019182525091905056fea164736f6c6343000814000ab53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103",
}

// L2UsdcBridgeProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use L2UsdcBridgeProxyMetaData.ABI instead.
var L2UsdcBridgeProxyABI = L2UsdcBridgeProxyMetaData.ABI

// L2UsdcBridgeProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L2UsdcBridgeProxyMetaData.Bin instead.
var L2UsdcBridgeProxyBin = L2UsdcBridgeProxyMetaData.Bin

// DeployL2UsdcBridgeProxy deploys a new Ethereum contract, binding an instance of L2UsdcBridgeProxy to it.
func DeployL2UsdcBridgeProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, initialOwner common.Address, _data []byte) (common.Address, *types.Transaction, *L2UsdcBridgeProxy, error) {
	parsed, err := L2UsdcBridgeProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L2UsdcBridgeProxyBin), backend, _logic, initialOwner, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L2UsdcBridgeProxy{L2UsdcBridgeProxyCaller: L2UsdcBridgeProxyCaller{contract: contract}, L2UsdcBridgeProxyTransactor: L2UsdcBridgeProxyTransactor{contract: contract}, L2UsdcBridgeProxyFilterer: L2UsdcBridgeProxyFilterer{contract: contract}}, nil
}

// L2UsdcBridgeProxy is an auto generated Go binding around an Ethereum contract.
type L2UsdcBridgeProxy struct {
	L2UsdcBridgeProxyCaller     // Read-only binding to the contract
	L2UsdcBridgeProxyTransactor // Write-only binding to the contract
	L2UsdcBridgeProxyFilterer   // Log filterer for contract events
}

// L2UsdcBridgeProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2UsdcBridgeProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2UsdcBridgeProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2UsdcBridgeProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2UsdcBridgeProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2UsdcBridgeProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2UsdcBridgeProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2UsdcBridgeProxySession struct {
	Contract     *L2UsdcBridgeProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// L2UsdcBridgeProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2UsdcBridgeProxyCallerSession struct {
	Contract *L2UsdcBridgeProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// L2UsdcBridgeProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2UsdcBridgeProxyTransactorSession struct {
	Contract     *L2UsdcBridgeProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// L2UsdcBridgeProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2UsdcBridgeProxyRaw struct {
	Contract *L2UsdcBridgeProxy // Generic contract binding to access the raw methods on
}

// L2UsdcBridgeProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2UsdcBridgeProxyCallerRaw struct {
	Contract *L2UsdcBridgeProxyCaller // Generic read-only contract binding to access the raw methods on
}

// L2UsdcBridgeProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2UsdcBridgeProxyTransactorRaw struct {
	Contract *L2UsdcBridgeProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2UsdcBridgeProxy creates a new instance of L2UsdcBridgeProxy, bound to a specific deployed contract.
func NewL2UsdcBridgeProxy(address common.Address, backend bind.ContractBackend) (*L2UsdcBridgeProxy, error) {
	contract, err := bindL2UsdcBridgeProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeProxy{L2UsdcBridgeProxyCaller: L2UsdcBridgeProxyCaller{contract: contract}, L2UsdcBridgeProxyTransactor: L2UsdcBridgeProxyTransactor{contract: contract}, L2UsdcBridgeProxyFilterer: L2UsdcBridgeProxyFilterer{contract: contract}}, nil
}

// NewL2UsdcBridgeProxyCaller creates a new read-only instance of L2UsdcBridgeProxy, bound to a specific deployed contract.
func NewL2UsdcBridgeProxyCaller(address common.Address, caller bind.ContractCaller) (*L2UsdcBridgeProxyCaller, error) {
	contract, err := bindL2UsdcBridgeProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeProxyCaller{contract: contract}, nil
}

// NewL2UsdcBridgeProxyTransactor creates a new write-only instance of L2UsdcBridgeProxy, bound to a specific deployed contract.
func NewL2UsdcBridgeProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*L2UsdcBridgeProxyTransactor, error) {
	contract, err := bindL2UsdcBridgeProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeProxyTransactor{contract: contract}, nil
}

// NewL2UsdcBridgeProxyFilterer creates a new log filterer instance of L2UsdcBridgeProxy, bound to a specific deployed contract.
func NewL2UsdcBridgeProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*L2UsdcBridgeProxyFilterer, error) {
	contract, err := bindL2UsdcBridgeProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeProxyFilterer{contract: contract}, nil
}

// bindL2UsdcBridgeProxy binds a generic wrapper to an already deployed contract.
func bindL2UsdcBridgeProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(L2UsdcBridgeProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2UsdcBridgeProxy.Contract.L2UsdcBridgeProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.L2UsdcBridgeProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.L2UsdcBridgeProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2UsdcBridgeProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridgeProxy.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) Implementation() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.Implementation(&_L2UsdcBridgeProxy.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerSession) Implementation() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.Implementation(&_L2UsdcBridgeProxy.CallOpts)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCaller) L1Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridgeProxy.contract.Call(opts, &out, "l1Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) L1Usdc() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.L1Usdc(&_L2UsdcBridgeProxy.CallOpts)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerSession) L1Usdc() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.L1Usdc(&_L2UsdcBridgeProxy.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCaller) L2Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridgeProxy.contract.Call(opts, &out, "l2Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) L2Usdc() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.L2Usdc(&_L2UsdcBridgeProxy.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerSession) L2Usdc() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.L2Usdc(&_L2UsdcBridgeProxy.CallOpts)
}

// L2UsdcMasterMinter is a free data retrieval call binding the contract method 0x05db940f.
//
// Solidity: function l2UsdcMasterMinter() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCaller) L2UsdcMasterMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridgeProxy.contract.Call(opts, &out, "l2UsdcMasterMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2UsdcMasterMinter is a free data retrieval call binding the contract method 0x05db940f.
//
// Solidity: function l2UsdcMasterMinter() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) L2UsdcMasterMinter() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.L2UsdcMasterMinter(&_L2UsdcBridgeProxy.CallOpts)
}

// L2UsdcMasterMinter is a free data retrieval call binding the contract method 0x05db940f.
//
// Solidity: function l2UsdcMasterMinter() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerSession) L2UsdcMasterMinter() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.L2UsdcMasterMinter(&_L2UsdcBridgeProxy.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCaller) Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridgeProxy.contract.Call(opts, &out, "messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) Messenger() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.Messenger(&_L2UsdcBridgeProxy.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerSession) Messenger() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.Messenger(&_L2UsdcBridgeProxy.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCaller) OtherBridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridgeProxy.contract.Call(opts, &out, "otherBridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) OtherBridge() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.OtherBridge(&_L2UsdcBridgeProxy.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerSession) OtherBridge() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.OtherBridge(&_L2UsdcBridgeProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridgeProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) Owner() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.Owner(&_L2UsdcBridgeProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyCallerSession) Owner() (common.Address, error) {
	return _L2UsdcBridgeProxy.Contract.Owner(&_L2UsdcBridgeProxy.CallOpts)
}

// ProxyChangeOwner is a paid mutator transaction binding the contract method 0xdfd3dcb3.
//
// Solidity: function proxyChangeOwner(address newAdmin) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactor) ProxyChangeOwner(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.contract.Transact(opts, "proxyChangeOwner", newAdmin)
}

// ProxyChangeOwner is a paid mutator transaction binding the contract method 0xdfd3dcb3.
//
// Solidity: function proxyChangeOwner(address newAdmin) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) ProxyChangeOwner(newAdmin common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.ProxyChangeOwner(&_L2UsdcBridgeProxy.TransactOpts, newAdmin)
}

// ProxyChangeOwner is a paid mutator transaction binding the contract method 0xdfd3dcb3.
//
// Solidity: function proxyChangeOwner(address newAdmin) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorSession) ProxyChangeOwner(newAdmin common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.ProxyChangeOwner(&_L2UsdcBridgeProxy.TransactOpts, newAdmin)
}

// SetAddress is a paid mutator transaction binding the contract method 0x12c59488.
//
// Solidity: function setAddress(address _messenger, address _otherBridge, address _l1Usdc, address _l2Usdc, address _l2UsdcMasterMinter) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactor) SetAddress(opts *bind.TransactOpts, _messenger common.Address, _otherBridge common.Address, _l1Usdc common.Address, _l2Usdc common.Address, _l2UsdcMasterMinter common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.contract.Transact(opts, "setAddress", _messenger, _otherBridge, _l1Usdc, _l2Usdc, _l2UsdcMasterMinter)
}

// SetAddress is a paid mutator transaction binding the contract method 0x12c59488.
//
// Solidity: function setAddress(address _messenger, address _otherBridge, address _l1Usdc, address _l2Usdc, address _l2UsdcMasterMinter) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) SetAddress(_messenger common.Address, _otherBridge common.Address, _l1Usdc common.Address, _l2Usdc common.Address, _l2UsdcMasterMinter common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.SetAddress(&_L2UsdcBridgeProxy.TransactOpts, _messenger, _otherBridge, _l1Usdc, _l2Usdc, _l2UsdcMasterMinter)
}

// SetAddress is a paid mutator transaction binding the contract method 0x12c59488.
//
// Solidity: function setAddress(address _messenger, address _otherBridge, address _l1Usdc, address _l2Usdc, address _l2UsdcMasterMinter) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorSession) SetAddress(_messenger common.Address, _otherBridge common.Address, _l1Usdc common.Address, _l2Usdc common.Address, _l2UsdcMasterMinter common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.SetAddress(&_L2UsdcBridgeProxy.TransactOpts, _messenger, _otherBridge, _l1Usdc, _l2Usdc, _l2UsdcMasterMinter)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.UpgradeTo(&_L2UsdcBridgeProxy.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.UpgradeTo(&_L2UsdcBridgeProxy.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.UpgradeToAndCall(&_L2UsdcBridgeProxy.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.UpgradeToAndCall(&_L2UsdcBridgeProxy.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.Fallback(&_L2UsdcBridgeProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.Fallback(&_L2UsdcBridgeProxy.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxySession) Receive() (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.Receive(&_L2UsdcBridgeProxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _L2UsdcBridgeProxy.Contract.Receive(&_L2UsdcBridgeProxy.TransactOpts)
}

// L2UsdcBridgeProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the L2UsdcBridgeProxy contract.
type L2UsdcBridgeProxyAdminChangedIterator struct {
	Event *L2UsdcBridgeProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *L2UsdcBridgeProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2UsdcBridgeProxyAdminChanged)
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
		it.Event = new(L2UsdcBridgeProxyAdminChanged)
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
func (it *L2UsdcBridgeProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2UsdcBridgeProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2UsdcBridgeProxyAdminChanged represents a AdminChanged event raised by the L2UsdcBridgeProxy contract.
type L2UsdcBridgeProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*L2UsdcBridgeProxyAdminChangedIterator, error) {

	logs, sub, err := _L2UsdcBridgeProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeProxyAdminChangedIterator{contract: _L2UsdcBridgeProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *L2UsdcBridgeProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _L2UsdcBridgeProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2UsdcBridgeProxyAdminChanged)
				if err := _L2UsdcBridgeProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyFilterer) ParseAdminChanged(log types.Log) (*L2UsdcBridgeProxyAdminChanged, error) {
	event := new(L2UsdcBridgeProxyAdminChanged)
	if err := _L2UsdcBridgeProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2UsdcBridgeProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the L2UsdcBridgeProxy contract.
type L2UsdcBridgeProxyUpgradedIterator struct {
	Event *L2UsdcBridgeProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *L2UsdcBridgeProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2UsdcBridgeProxyUpgraded)
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
		it.Event = new(L2UsdcBridgeProxyUpgraded)
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
func (it *L2UsdcBridgeProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2UsdcBridgeProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2UsdcBridgeProxyUpgraded represents a Upgraded event raised by the L2UsdcBridgeProxy contract.
type L2UsdcBridgeProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*L2UsdcBridgeProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L2UsdcBridgeProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeProxyUpgradedIterator{contract: _L2UsdcBridgeProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *L2UsdcBridgeProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L2UsdcBridgeProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2UsdcBridgeProxyUpgraded)
				if err := _L2UsdcBridgeProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_L2UsdcBridgeProxy *L2UsdcBridgeProxyFilterer) ParseUpgraded(log types.Log) (*L2UsdcBridgeProxyUpgraded, error) {
	event := new(L2UsdcBridgeProxyUpgraded)
	if err := _L2UsdcBridgeProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
