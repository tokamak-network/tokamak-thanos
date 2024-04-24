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

// L1UsdcBridgeMetaData contains all meta data concerning the L1UsdcBridge contract.
var L1UsdcBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l2Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"ERC20DepositInitiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l2Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"ERC20WithdrawalFinalized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_minGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"depositERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_minGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"depositERC20To\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"finalizeERC20Withdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1Usdc\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2TokenBridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2Usdc\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messenger\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"otherBridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061113d806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638f601f6611610076578063a1b4bc041161005b578063a1b4bc0414610191578063a9f9e675146101b1578063c89701a2146101c457600080fd5b80638f601f661461013a57806391c49bf81461017357600080fd5b80633cb747bf146100a857806356c3b587146100f257806358a997f614610112578063838b252014610127575b600080fd5b6000546100c89073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6002546100c89073ffffffffffffffffffffffffffffffffffffffff1681565b610125610120366004610d28565b6101e4565b005b610125610135366004610dab565b610290565b610165610148366004610e41565b600460209081526000928352604080842090915290825290205481565b6040519081526020016100e9565b60015473ffffffffffffffffffffffffffffffffffffffff166100c8565b6003546100c89073ffffffffffffffffffffffffffffffffffffffff1681565b6101256101bf366004610e7a565b6102a9565b6001546100c89073ffffffffffffffffffffffffffffffffffffffff1681565b333b15610278576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603760248201527f5374616e646172644272696467653a2066756e6374696f6e2063616e206f6e6c60448201527f792062652063616c6c65642066726f6d20616e20454f4100000000000000000060648201526084015b60405180910390fd5b610288868633338888888861064b565b505050505050565b6102a0878733888888888861064b565b50505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314801561037e5750600154600054604080517f6e296e45000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff9384169390921691636e296e45916004808201926020929091908290030181865afa158015610342573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103669190610ef3565b73ffffffffffffffffffffffffffffffffffffffff16145b610430576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604160248201527f5374616e646172644272696467653a2066756e6374696f6e2063616e206f6e6c60448201527f792062652063616c6c65642066726f6d20746865206f7468657220627269646760648201527f6500000000000000000000000000000000000000000000000000000000000000608482015260a40161026f565b600254879073ffffffffffffffffffffffffffffffffffffffff8083169116146104b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f6e6f74204c312075736463000000000000000000000000000000000000000000604482015260640161026f565b600354879073ffffffffffffffffffffffffffffffffffffffff80831691161461053c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f6e6f74204c322075736463000000000000000000000000000000000000000000604482015260640161026f565b73ffffffffffffffffffffffffffffffffffffffff808a166000908152600460209081526040808320938c168352929052205461057a908690610f3f565b73ffffffffffffffffffffffffffffffffffffffff808b166000818152600460209081526040808320948e16835293905291909120919091556105be908787610990565b8673ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff167f3ceee06c1e37648fcbb6ed52e17b3e1f275a1f8c7b22a84b2b84732431e046b3898989896040516106389493929190610f9f565b60405180910390a4505050505050505050565b600254889073ffffffffffffffffffffffffffffffffffffffff8083169116146106d1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f6e6f74204c312075736463000000000000000000000000000000000000000000604482015260640161026f565b600354889073ffffffffffffffffffffffffffffffffffffffff808316911614610757576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f6e6f74204c322075736463000000000000000000000000000000000000000000604482015260640161026f565b61077973ffffffffffffffffffffffffffffffffffffffff8b16893089610a16565b73ffffffffffffffffffffffffffffffffffffffff808b166000908152600460209081526040808320938d16835292905220546107b7908790610fd5565b73ffffffffffffffffffffffffffffffffffffffff808c1660009081526004602090815260408083208e851684529091528082209390935554600154925190821692633dbb202b9216907f662a633a000000000000000000000000000000000000000000000000000000009061083d908f908f908f908f908f908e908e90602401610fed565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009485161790525160e085901b90921682526108d092918a90600401611076565b600060405180830381600087803b1580156108ea57600080fd5b505af11580156108fe573d6000803e3d6000fd5b505050508773ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff168b73ffffffffffffffffffffffffffffffffffffffff167f718594027abd4eaed59f95162563e0cc6d0e8d5b86b1c7be8b1b0ac3343d03968a8a898960405161097c9493929190610f9f565b60405180910390a450505050505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff838116602483015260448201839052610a1191859182169063a9059cbb906064015b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050610a62565b505050565b60405173ffffffffffffffffffffffffffffffffffffffff8481166024830152838116604483015260648201839052610a5c9186918216906323b872dd906084016109ca565b50505050565b6000610a8473ffffffffffffffffffffffffffffffffffffffff841683610af8565b90508051600014158015610aa9575080806020019051810190610aa791906110f2565b155b15610a11576040517f5274afe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416600482015260240161026f565b6060610b0683836000610b0d565b9392505050565b606081471015610b4b576040517fcd78605900000000000000000000000000000000000000000000000000000000815230600482015260240161026f565b6000808573ffffffffffffffffffffffffffffffffffffffff168486604051610b749190611114565b60006040518083038185875af1925050503d8060008114610bb1576040519150601f19603f3d011682016040523d82523d6000602084013e610bb6565b606091505b5091509150610bc6868383610bd0565b9695505050505050565b606082610be557610be082610c5f565b610b06565b8151158015610c09575073ffffffffffffffffffffffffffffffffffffffff84163b155b15610c58576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260240161026f565b5080610b06565b805115610c6f5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50565b73ffffffffffffffffffffffffffffffffffffffff81168114610ca157600080fd5b803563ffffffff81168114610cda57600080fd5b919050565b60008083601f840112610cf157600080fd5b50813567ffffffffffffffff811115610d0957600080fd5b602083019150836020828501011115610d2157600080fd5b9250929050565b60008060008060008060a08789031215610d4157600080fd5b8635610d4c81610ca4565b95506020870135610d5c81610ca4565b945060408701359350610d7160608801610cc6565b9250608087013567ffffffffffffffff811115610d8d57600080fd5b610d9989828a01610cdf565b979a9699509497509295939492505050565b600080600080600080600060c0888a031215610dc657600080fd5b8735610dd181610ca4565b96506020880135610de181610ca4565b95506040880135610df181610ca4565b945060608801359350610e0660808901610cc6565b925060a088013567ffffffffffffffff811115610e2257600080fd5b610e2e8a828b01610cdf565b989b979a50959850939692959293505050565b60008060408385031215610e5457600080fd5b8235610e5f81610ca4565b91506020830135610e6f81610ca4565b809150509250929050565b600080600080600080600060c0888a031215610e9557600080fd5b8735610ea081610ca4565b96506020880135610eb081610ca4565b95506040880135610ec081610ca4565b94506060880135610ed081610ca4565b93506080880135925060a088013567ffffffffffffffff811115610e2257600080fd5b600060208284031215610f0557600080fd5b8151610b0681610ca4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600082821015610f5157610f51610f10565b500390565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152606060408201526000610bc6606083018486610f56565b60008219821115610fe857610fe8610f10565b500190565b600073ffffffffffffffffffffffffffffffffffffffff808a1683528089166020840152808816604084015280871660608401525084608083015260c060a083015261103d60c083018486610f56565b9998505050505050505050565b60005b8381101561106557818101518382015260200161104d565b83811115610a5c5750506000910152565b73ffffffffffffffffffffffffffffffffffffffff8416815260606020820152600083518060608401526110b181608085016020880161104a565b63ffffffff93909316604083015250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160160800192915050565b60006020828403121561110457600080fd5b81518015158114610b0657600080fd5b6000825161112681846020870161104a565b919091019291505056fea164736f6c634300080f000a",
}

// L1UsdcBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use L1UsdcBridgeMetaData.ABI instead.
var L1UsdcBridgeABI = L1UsdcBridgeMetaData.ABI

// L1UsdcBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L1UsdcBridgeMetaData.Bin instead.
var L1UsdcBridgeBin = L1UsdcBridgeMetaData.Bin

// DeployL1UsdcBridge deploys a new Ethereum contract, binding an instance of L1UsdcBridge to it.
func DeployL1UsdcBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *L1UsdcBridge, error) {
	parsed, err := L1UsdcBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L1UsdcBridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L1UsdcBridge{L1UsdcBridgeCaller: L1UsdcBridgeCaller{contract: contract}, L1UsdcBridgeTransactor: L1UsdcBridgeTransactor{contract: contract}, L1UsdcBridgeFilterer: L1UsdcBridgeFilterer{contract: contract}}, nil
}

// L1UsdcBridge is an auto generated Go binding around an Ethereum contract.
type L1UsdcBridge struct {
	L1UsdcBridgeCaller     // Read-only binding to the contract
	L1UsdcBridgeTransactor // Write-only binding to the contract
	L1UsdcBridgeFilterer   // Log filterer for contract events
}

// L1UsdcBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1UsdcBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1UsdcBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1UsdcBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1UsdcBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1UsdcBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1UsdcBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1UsdcBridgeSession struct {
	Contract     *L1UsdcBridge     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1UsdcBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1UsdcBridgeCallerSession struct {
	Contract *L1UsdcBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// L1UsdcBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1UsdcBridgeTransactorSession struct {
	Contract     *L1UsdcBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// L1UsdcBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1UsdcBridgeRaw struct {
	Contract *L1UsdcBridge // Generic contract binding to access the raw methods on
}

// L1UsdcBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1UsdcBridgeCallerRaw struct {
	Contract *L1UsdcBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// L1UsdcBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1UsdcBridgeTransactorRaw struct {
	Contract *L1UsdcBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1UsdcBridge creates a new instance of L1UsdcBridge, bound to a specific deployed contract.
func NewL1UsdcBridge(address common.Address, backend bind.ContractBackend) (*L1UsdcBridge, error) {
	contract, err := bindL1UsdcBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridge{L1UsdcBridgeCaller: L1UsdcBridgeCaller{contract: contract}, L1UsdcBridgeTransactor: L1UsdcBridgeTransactor{contract: contract}, L1UsdcBridgeFilterer: L1UsdcBridgeFilterer{contract: contract}}, nil
}

// NewL1UsdcBridgeCaller creates a new read-only instance of L1UsdcBridge, bound to a specific deployed contract.
func NewL1UsdcBridgeCaller(address common.Address, caller bind.ContractCaller) (*L1UsdcBridgeCaller, error) {
	contract, err := bindL1UsdcBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeCaller{contract: contract}, nil
}

// NewL1UsdcBridgeTransactor creates a new write-only instance of L1UsdcBridge, bound to a specific deployed contract.
func NewL1UsdcBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*L1UsdcBridgeTransactor, error) {
	contract, err := bindL1UsdcBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeTransactor{contract: contract}, nil
}

// NewL1UsdcBridgeFilterer creates a new log filterer instance of L1UsdcBridge, bound to a specific deployed contract.
func NewL1UsdcBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*L1UsdcBridgeFilterer, error) {
	contract, err := bindL1UsdcBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeFilterer{contract: contract}, nil
}

// bindL1UsdcBridge binds a generic wrapper to an already deployed contract.
func bindL1UsdcBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1UsdcBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1UsdcBridge *L1UsdcBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1UsdcBridge.Contract.L1UsdcBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1UsdcBridge *L1UsdcBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.L1UsdcBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1UsdcBridge *L1UsdcBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.L1UsdcBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1UsdcBridge *L1UsdcBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1UsdcBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1UsdcBridge *L1UsdcBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1UsdcBridge *L1UsdcBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.contract.Transact(opts, method, params...)
}

// Deposits is a free data retrieval call binding the contract method 0x8f601f66.
//
// Solidity: function deposits(address , address ) view returns(uint256)
func (_L1UsdcBridge *L1UsdcBridgeCaller) Deposits(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _L1UsdcBridge.contract.Call(opts, &out, "deposits", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Deposits is a free data retrieval call binding the contract method 0x8f601f66.
//
// Solidity: function deposits(address , address ) view returns(uint256)
func (_L1UsdcBridge *L1UsdcBridgeSession) Deposits(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _L1UsdcBridge.Contract.Deposits(&_L1UsdcBridge.CallOpts, arg0, arg1)
}

// Deposits is a free data retrieval call binding the contract method 0x8f601f66.
//
// Solidity: function deposits(address , address ) view returns(uint256)
func (_L1UsdcBridge *L1UsdcBridgeCallerSession) Deposits(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _L1UsdcBridge.Contract.Deposits(&_L1UsdcBridge.CallOpts, arg0, arg1)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCaller) L1Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridge.contract.Call(opts, &out, "l1Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeSession) L1Usdc() (common.Address, error) {
	return _L1UsdcBridge.Contract.L1Usdc(&_L1UsdcBridge.CallOpts)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCallerSession) L1Usdc() (common.Address, error) {
	return _L1UsdcBridge.Contract.L1Usdc(&_L1UsdcBridge.CallOpts)
}

// L2TokenBridge is a free data retrieval call binding the contract method 0x91c49bf8.
//
// Solidity: function l2TokenBridge() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCaller) L2TokenBridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridge.contract.Call(opts, &out, "l2TokenBridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2TokenBridge is a free data retrieval call binding the contract method 0x91c49bf8.
//
// Solidity: function l2TokenBridge() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeSession) L2TokenBridge() (common.Address, error) {
	return _L1UsdcBridge.Contract.L2TokenBridge(&_L1UsdcBridge.CallOpts)
}

// L2TokenBridge is a free data retrieval call binding the contract method 0x91c49bf8.
//
// Solidity: function l2TokenBridge() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCallerSession) L2TokenBridge() (common.Address, error) {
	return _L1UsdcBridge.Contract.L2TokenBridge(&_L1UsdcBridge.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCaller) L2Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridge.contract.Call(opts, &out, "l2Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeSession) L2Usdc() (common.Address, error) {
	return _L1UsdcBridge.Contract.L2Usdc(&_L1UsdcBridge.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCallerSession) L2Usdc() (common.Address, error) {
	return _L1UsdcBridge.Contract.L2Usdc(&_L1UsdcBridge.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCaller) Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridge.contract.Call(opts, &out, "messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeSession) Messenger() (common.Address, error) {
	return _L1UsdcBridge.Contract.Messenger(&_L1UsdcBridge.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCallerSession) Messenger() (common.Address, error) {
	return _L1UsdcBridge.Contract.Messenger(&_L1UsdcBridge.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCaller) OtherBridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1UsdcBridge.contract.Call(opts, &out, "otherBridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeSession) OtherBridge() (common.Address, error) {
	return _L1UsdcBridge.Contract.OtherBridge(&_L1UsdcBridge.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L1UsdcBridge *L1UsdcBridgeCallerSession) OtherBridge() (common.Address, error) {
	return _L1UsdcBridge.Contract.OtherBridge(&_L1UsdcBridge.CallOpts)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x58a997f6.
//
// Solidity: function depositERC20(address _l1Token, address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeTransactor) DepositERC20(opts *bind.TransactOpts, _l1Token common.Address, _l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.contract.Transact(opts, "depositERC20", _l1Token, _l2Token, _amount, _minGasLimit, _extraData)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x58a997f6.
//
// Solidity: function depositERC20(address _l1Token, address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeSession) DepositERC20(_l1Token common.Address, _l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.DepositERC20(&_L1UsdcBridge.TransactOpts, _l1Token, _l2Token, _amount, _minGasLimit, _extraData)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x58a997f6.
//
// Solidity: function depositERC20(address _l1Token, address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeTransactorSession) DepositERC20(_l1Token common.Address, _l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.DepositERC20(&_L1UsdcBridge.TransactOpts, _l1Token, _l2Token, _amount, _minGasLimit, _extraData)
}

// DepositERC20To is a paid mutator transaction binding the contract method 0x838b2520.
//
// Solidity: function depositERC20To(address _l1Token, address _l2Token, address _to, uint256 _amount, uint32 _minGasLimit, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeTransactor) DepositERC20To(opts *bind.TransactOpts, _l1Token common.Address, _l2Token common.Address, _to common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.contract.Transact(opts, "depositERC20To", _l1Token, _l2Token, _to, _amount, _minGasLimit, _extraData)
}

// DepositERC20To is a paid mutator transaction binding the contract method 0x838b2520.
//
// Solidity: function depositERC20To(address _l1Token, address _l2Token, address _to, uint256 _amount, uint32 _minGasLimit, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeSession) DepositERC20To(_l1Token common.Address, _l2Token common.Address, _to common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.DepositERC20To(&_L1UsdcBridge.TransactOpts, _l1Token, _l2Token, _to, _amount, _minGasLimit, _extraData)
}

// DepositERC20To is a paid mutator transaction binding the contract method 0x838b2520.
//
// Solidity: function depositERC20To(address _l1Token, address _l2Token, address _to, uint256 _amount, uint32 _minGasLimit, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeTransactorSession) DepositERC20To(_l1Token common.Address, _l2Token common.Address, _to common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.DepositERC20To(&_L1UsdcBridge.TransactOpts, _l1Token, _l2Token, _to, _amount, _minGasLimit, _extraData)
}

// FinalizeERC20Withdrawal is a paid mutator transaction binding the contract method 0xa9f9e675.
//
// Solidity: function finalizeERC20Withdrawal(address _l1Token, address _l2Token, address _from, address _to, uint256 _amount, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeTransactor) FinalizeERC20Withdrawal(opts *bind.TransactOpts, _l1Token common.Address, _l2Token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.contract.Transact(opts, "finalizeERC20Withdrawal", _l1Token, _l2Token, _from, _to, _amount, _extraData)
}

// FinalizeERC20Withdrawal is a paid mutator transaction binding the contract method 0xa9f9e675.
//
// Solidity: function finalizeERC20Withdrawal(address _l1Token, address _l2Token, address _from, address _to, uint256 _amount, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeSession) FinalizeERC20Withdrawal(_l1Token common.Address, _l2Token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.FinalizeERC20Withdrawal(&_L1UsdcBridge.TransactOpts, _l1Token, _l2Token, _from, _to, _amount, _extraData)
}

// FinalizeERC20Withdrawal is a paid mutator transaction binding the contract method 0xa9f9e675.
//
// Solidity: function finalizeERC20Withdrawal(address _l1Token, address _l2Token, address _from, address _to, uint256 _amount, bytes _extraData) returns()
func (_L1UsdcBridge *L1UsdcBridgeTransactorSession) FinalizeERC20Withdrawal(_l1Token common.Address, _l2Token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _L1UsdcBridge.Contract.FinalizeERC20Withdrawal(&_L1UsdcBridge.TransactOpts, _l1Token, _l2Token, _from, _to, _amount, _extraData)
}

// L1UsdcBridgeERC20DepositInitiatedIterator is returned from FilterERC20DepositInitiated and is used to iterate over the raw logs and unpacked data for ERC20DepositInitiated events raised by the L1UsdcBridge contract.
type L1UsdcBridgeERC20DepositInitiatedIterator struct {
	Event *L1UsdcBridgeERC20DepositInitiated // Event containing the contract specifics and raw log

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
func (it *L1UsdcBridgeERC20DepositInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1UsdcBridgeERC20DepositInitiated)
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
		it.Event = new(L1UsdcBridgeERC20DepositInitiated)
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
func (it *L1UsdcBridgeERC20DepositInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1UsdcBridgeERC20DepositInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1UsdcBridgeERC20DepositInitiated represents a ERC20DepositInitiated event raised by the L1UsdcBridge contract.
type L1UsdcBridgeERC20DepositInitiated struct {
	L1Token   common.Address
	L2Token   common.Address
	From      common.Address
	To        common.Address
	Amount    *big.Int
	ExtraData []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterERC20DepositInitiated is a free log retrieval operation binding the contract event 0x718594027abd4eaed59f95162563e0cc6d0e8d5b86b1c7be8b1b0ac3343d0396.
//
// Solidity: event ERC20DepositInitiated(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L1UsdcBridge *L1UsdcBridgeFilterer) FilterERC20DepositInitiated(opts *bind.FilterOpts, l1Token []common.Address, l2Token []common.Address, from []common.Address) (*L1UsdcBridgeERC20DepositInitiatedIterator, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var l2TokenRule []interface{}
	for _, l2TokenItem := range l2Token {
		l2TokenRule = append(l2TokenRule, l2TokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _L1UsdcBridge.contract.FilterLogs(opts, "ERC20DepositInitiated", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeERC20DepositInitiatedIterator{contract: _L1UsdcBridge.contract, event: "ERC20DepositInitiated", logs: logs, sub: sub}, nil
}

// WatchERC20DepositInitiated is a free log subscription operation binding the contract event 0x718594027abd4eaed59f95162563e0cc6d0e8d5b86b1c7be8b1b0ac3343d0396.
//
// Solidity: event ERC20DepositInitiated(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L1UsdcBridge *L1UsdcBridgeFilterer) WatchERC20DepositInitiated(opts *bind.WatchOpts, sink chan<- *L1UsdcBridgeERC20DepositInitiated, l1Token []common.Address, l2Token []common.Address, from []common.Address) (event.Subscription, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var l2TokenRule []interface{}
	for _, l2TokenItem := range l2Token {
		l2TokenRule = append(l2TokenRule, l2TokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _L1UsdcBridge.contract.WatchLogs(opts, "ERC20DepositInitiated", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1UsdcBridgeERC20DepositInitiated)
				if err := _L1UsdcBridge.contract.UnpackLog(event, "ERC20DepositInitiated", log); err != nil {
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

// ParseERC20DepositInitiated is a log parse operation binding the contract event 0x718594027abd4eaed59f95162563e0cc6d0e8d5b86b1c7be8b1b0ac3343d0396.
//
// Solidity: event ERC20DepositInitiated(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L1UsdcBridge *L1UsdcBridgeFilterer) ParseERC20DepositInitiated(log types.Log) (*L1UsdcBridgeERC20DepositInitiated, error) {
	event := new(L1UsdcBridgeERC20DepositInitiated)
	if err := _L1UsdcBridge.contract.UnpackLog(event, "ERC20DepositInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1UsdcBridgeERC20WithdrawalFinalizedIterator is returned from FilterERC20WithdrawalFinalized and is used to iterate over the raw logs and unpacked data for ERC20WithdrawalFinalized events raised by the L1UsdcBridge contract.
type L1UsdcBridgeERC20WithdrawalFinalizedIterator struct {
	Event *L1UsdcBridgeERC20WithdrawalFinalized // Event containing the contract specifics and raw log

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
func (it *L1UsdcBridgeERC20WithdrawalFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1UsdcBridgeERC20WithdrawalFinalized)
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
		it.Event = new(L1UsdcBridgeERC20WithdrawalFinalized)
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
func (it *L1UsdcBridgeERC20WithdrawalFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1UsdcBridgeERC20WithdrawalFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1UsdcBridgeERC20WithdrawalFinalized represents a ERC20WithdrawalFinalized event raised by the L1UsdcBridge contract.
type L1UsdcBridgeERC20WithdrawalFinalized struct {
	L1Token   common.Address
	L2Token   common.Address
	From      common.Address
	To        common.Address
	Amount    *big.Int
	ExtraData []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterERC20WithdrawalFinalized is a free log retrieval operation binding the contract event 0x3ceee06c1e37648fcbb6ed52e17b3e1f275a1f8c7b22a84b2b84732431e046b3.
//
// Solidity: event ERC20WithdrawalFinalized(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L1UsdcBridge *L1UsdcBridgeFilterer) FilterERC20WithdrawalFinalized(opts *bind.FilterOpts, l1Token []common.Address, l2Token []common.Address, from []common.Address) (*L1UsdcBridgeERC20WithdrawalFinalizedIterator, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var l2TokenRule []interface{}
	for _, l2TokenItem := range l2Token {
		l2TokenRule = append(l2TokenRule, l2TokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _L1UsdcBridge.contract.FilterLogs(opts, "ERC20WithdrawalFinalized", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &L1UsdcBridgeERC20WithdrawalFinalizedIterator{contract: _L1UsdcBridge.contract, event: "ERC20WithdrawalFinalized", logs: logs, sub: sub}, nil
}

// WatchERC20WithdrawalFinalized is a free log subscription operation binding the contract event 0x3ceee06c1e37648fcbb6ed52e17b3e1f275a1f8c7b22a84b2b84732431e046b3.
//
// Solidity: event ERC20WithdrawalFinalized(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L1UsdcBridge *L1UsdcBridgeFilterer) WatchERC20WithdrawalFinalized(opts *bind.WatchOpts, sink chan<- *L1UsdcBridgeERC20WithdrawalFinalized, l1Token []common.Address, l2Token []common.Address, from []common.Address) (event.Subscription, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var l2TokenRule []interface{}
	for _, l2TokenItem := range l2Token {
		l2TokenRule = append(l2TokenRule, l2TokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _L1UsdcBridge.contract.WatchLogs(opts, "ERC20WithdrawalFinalized", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1UsdcBridgeERC20WithdrawalFinalized)
				if err := _L1UsdcBridge.contract.UnpackLog(event, "ERC20WithdrawalFinalized", log); err != nil {
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

// ParseERC20WithdrawalFinalized is a log parse operation binding the contract event 0x3ceee06c1e37648fcbb6ed52e17b3e1f275a1f8c7b22a84b2b84732431e046b3.
//
// Solidity: event ERC20WithdrawalFinalized(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L1UsdcBridge *L1UsdcBridgeFilterer) ParseERC20WithdrawalFinalized(log types.Log) (*L1UsdcBridgeERC20WithdrawalFinalized, error) {
	event := new(L1UsdcBridgeERC20WithdrawalFinalized)
	if err := _L1UsdcBridge.contract.UnpackLog(event, "ERC20WithdrawalFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
