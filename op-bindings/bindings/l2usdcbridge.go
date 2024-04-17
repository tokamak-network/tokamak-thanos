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

// L2UsdcBridgeMetaData contains all meta data concerning the L2UsdcBridge contract.
var L2UsdcBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l2Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"DepositFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l2Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"WithdrawalInitiated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"finalizeDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1Usdc\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2Usdc\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2UsdcMasterMinter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messenger\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"otherBridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_minGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506111d0806100206000396000f3fe6080604052600436106100705760003560e01c806356c3b5871161004e57806356c3b5871461010d578063662a633a1461013a578063a1b4bc041461015a578063c89701a21461018757600080fd5b806305db940f1461007557806332b7006d146100cb5780633cb747bf146100e0575b600080fd5b34801561008157600080fd5b506004546100a29073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6100de6100d9366004610ebf565b6101b4565b005b3480156100ec57600080fd5b506000546100a29073ffffffffffffffffffffffffffffffffffffffff1681565b34801561011957600080fd5b506002546100a29073ffffffffffffffffffffffffffffffffffffffff1681565b34801561014657600080fd5b506100de610155366004610f3b565b61025e565b34801561016657600080fd5b506003546100a29073ffffffffffffffffffffffffffffffffffffffff1681565b34801561019357600080fd5b506001546100a29073ffffffffffffffffffffffffffffffffffffffff1681565b333b15610248576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603760248201527f5374616e646172644272696467653a2066756e6374696f6e2063616e206f6e6c60448201527f792062652063616c6c65642066726f6d20616e20454f4100000000000000000060648201526084015b60405180910390fd5b61025785333387878787610759565b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331480156103335750600154600054604080517f6e296e45000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff9384169390921691636e296e45916004808201926020929091908290030181865afa1580156102f7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031b9190610fd3565b73ffffffffffffffffffffffffffffffffffffffff16145b6103e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604160248201527f5374616e646172644272696467653a2066756e6374696f6e2063616e206f6e6c60448201527f792062652063616c6c65642066726f6d20746865206f7468657220627269646760648201527f6500000000000000000000000000000000000000000000000000000000000000608482015260a40161023f565b600254879073ffffffffffffffffffffffffffffffffffffffff80831691161461046b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f6e6f74204c312075736463000000000000000000000000000000000000000000604482015260640161023f565b600354879073ffffffffffffffffffffffffffffffffffffffff8083169116146104f1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f6e6f74204c322075736463000000000000000000000000000000000000000000604482015260640161023f565b6040517f8a6db9c300000000000000000000000000000000000000000000000000000000815230600482015260009073ffffffffffffffffffffffffffffffffffffffff8a1690638a6db9c390602401602060405180830381865afa15801561055e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105829190610ff0565b90508581101561064357600480546040517fcbf2b8bf0000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9281019290925273ffffffffffffffffffffffffffffffffffffffff169063cbf2b8bf906024016020604051808303816000875af115801561061d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106419190611009565b505b6040517f40c10f1900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8881166004830152602482018890528a16906340c10f1990604401600060405180830381600087803b1580156106b357600080fd5b505af11580156106c7573d6000803e3d6000fd5b505050508773ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff168b73ffffffffffffffffffffffffffffffffffffffff167fb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd898a8a8a8a6040516107459493929190611074565b60405180910390a450505050505050505050565b600354879073ffffffffffffffffffffffffffffffffffffffff8083169116146107df576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f6e6f74204c322075736463000000000000000000000000000000000000000000604482015260640161023f565b6040517f8a6db9c300000000000000000000000000000000000000000000000000000000815230600482015260009073ffffffffffffffffffffffffffffffffffffffff8a1690638a6db9c390602401602060405180830381865afa15801561084c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108709190610ff0565b90508581101561093157600480546040517fcbf2b8bf0000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9281019290925273ffffffffffffffffffffffffffffffffffffffff169063cbf2b8bf906024016020604051808303816000875af115801561090b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061092f9190611009565b505b61095373ffffffffffffffffffffffffffffffffffffffff8a16893089610b72565b6040517f42966c680000000000000000000000000000000000000000000000000000000081526004810187905273ffffffffffffffffffffffffffffffffffffffff8a16906342966c6890602401600060405180830381600087803b1580156109bb57600080fd5b505af11580156109cf573d6000803e3d6000fd5b505060005460015460025460405173ffffffffffffffffffffffffffffffffffffffff9384169550633dbb202b9450918316927fa9f9e6750000000000000000000000000000000000000000000000000000000092610a4092909116908f908f908f908f908e908e906024016110aa565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009485161790525160e085901b9092168252610ad392918a9060040161112b565b600060405180830381600087803b158015610aed57600080fd5b505af1158015610b01573d6000803e3d6000fd5b505060025460405173ffffffffffffffffffffffffffffffffffffffff8c811694508d81169350909116907f73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e90610b5f908c908c908b908b90611074565b60405180910390a4505050505050505050565b6040805173ffffffffffffffffffffffffffffffffffffffff85811660248301528416604482015260648082018490528251808303909101815260849091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd00000000000000000000000000000000000000000000000000000000179052610c07908590610c0d565b50505050565b6000610c2f73ffffffffffffffffffffffffffffffffffffffff841683610ca8565b90508051600014158015610c54575080806020019051810190610c529190611009565b155b15610ca3576040517f5274afe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416600482015260240161023f565b505050565b6060610cb683836000610cbd565b9392505050565b606081471015610cfb576040517fcd78605900000000000000000000000000000000000000000000000000000000815230600482015260240161023f565b6000808573ffffffffffffffffffffffffffffffffffffffff168486604051610d2491906111a7565b60006040518083038185875af1925050503d8060008114610d61576040519150601f19603f3d011682016040523d82523d6000602084013e610d66565b606091505b5091509150610d76868383610d80565b9695505050505050565b606082610d9557610d9082610e0f565b610cb6565b8151158015610db9575073ffffffffffffffffffffffffffffffffffffffff84163b155b15610e08576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260240161023f565b5080610cb6565b805115610e1f5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50565b73ffffffffffffffffffffffffffffffffffffffff81168114610e5157600080fd5b60008083601f840112610e8857600080fd5b50813567ffffffffffffffff811115610ea057600080fd5b602083019150836020828501011115610eb857600080fd5b9250929050565b600080600080600060808688031215610ed757600080fd5b8535610ee281610e54565b945060208601359350604086013563ffffffff81168114610f0257600080fd5b9250606086013567ffffffffffffffff811115610f1e57600080fd5b610f2a88828901610e76565b969995985093965092949392505050565b600080600080600080600060c0888a031215610f5657600080fd5b8735610f6181610e54565b96506020880135610f7181610e54565b95506040880135610f8181610e54565b94506060880135610f9181610e54565b93506080880135925060a088013567ffffffffffffffff811115610fb457600080fd5b610fc08a828b01610e76565b989b979a50959850939692959293505050565b600060208284031215610fe557600080fd5b8151610cb681610e54565b60006020828403121561100257600080fd5b5051919050565b60006020828403121561101b57600080fd5b81518015158114610cb657600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152606060408201526000610d7660608301848661102b565b600073ffffffffffffffffffffffffffffffffffffffff808a1683528089166020840152808816604084015280871660608401525084608083015260c060a08301526110fa60c08301848661102b565b9998505050505050505050565b60005b8381101561112257818101518382015260200161110a565b50506000910152565b73ffffffffffffffffffffffffffffffffffffffff841681526060602082015260008351806060840152611166816080850160208801611107565b63ffffffff93909316604083015250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160160800192915050565b600082516111b9818460208701611107565b919091019291505056fea164736f6c6343000815000a",
}

// L2UsdcBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use L2UsdcBridgeMetaData.ABI instead.
var L2UsdcBridgeABI = L2UsdcBridgeMetaData.ABI

// L2UsdcBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L2UsdcBridgeMetaData.Bin instead.
var L2UsdcBridgeBin = L2UsdcBridgeMetaData.Bin

// DeployL2UsdcBridge deploys a new Ethereum contract, binding an instance of L2UsdcBridge to it.
func DeployL2UsdcBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *L2UsdcBridge, error) {
	parsed, err := L2UsdcBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L2UsdcBridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L2UsdcBridge{L2UsdcBridgeCaller: L2UsdcBridgeCaller{contract: contract}, L2UsdcBridgeTransactor: L2UsdcBridgeTransactor{contract: contract}, L2UsdcBridgeFilterer: L2UsdcBridgeFilterer{contract: contract}}, nil
}

// L2UsdcBridge is an auto generated Go binding around an Ethereum contract.
type L2UsdcBridge struct {
	L2UsdcBridgeCaller     // Read-only binding to the contract
	L2UsdcBridgeTransactor // Write-only binding to the contract
	L2UsdcBridgeFilterer   // Log filterer for contract events
}

// L2UsdcBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2UsdcBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2UsdcBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2UsdcBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2UsdcBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2UsdcBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2UsdcBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2UsdcBridgeSession struct {
	Contract     *L2UsdcBridge     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L2UsdcBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2UsdcBridgeCallerSession struct {
	Contract *L2UsdcBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// L2UsdcBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2UsdcBridgeTransactorSession struct {
	Contract     *L2UsdcBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// L2UsdcBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2UsdcBridgeRaw struct {
	Contract *L2UsdcBridge // Generic contract binding to access the raw methods on
}

// L2UsdcBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2UsdcBridgeCallerRaw struct {
	Contract *L2UsdcBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// L2UsdcBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2UsdcBridgeTransactorRaw struct {
	Contract *L2UsdcBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2UsdcBridge creates a new instance of L2UsdcBridge, bound to a specific deployed contract.
func NewL2UsdcBridge(address common.Address, backend bind.ContractBackend) (*L2UsdcBridge, error) {
	contract, err := bindL2UsdcBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridge{L2UsdcBridgeCaller: L2UsdcBridgeCaller{contract: contract}, L2UsdcBridgeTransactor: L2UsdcBridgeTransactor{contract: contract}, L2UsdcBridgeFilterer: L2UsdcBridgeFilterer{contract: contract}}, nil
}

// NewL2UsdcBridgeCaller creates a new read-only instance of L2UsdcBridge, bound to a specific deployed contract.
func NewL2UsdcBridgeCaller(address common.Address, caller bind.ContractCaller) (*L2UsdcBridgeCaller, error) {
	contract, err := bindL2UsdcBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeCaller{contract: contract}, nil
}

// NewL2UsdcBridgeTransactor creates a new write-only instance of L2UsdcBridge, bound to a specific deployed contract.
func NewL2UsdcBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*L2UsdcBridgeTransactor, error) {
	contract, err := bindL2UsdcBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeTransactor{contract: contract}, nil
}

// NewL2UsdcBridgeFilterer creates a new log filterer instance of L2UsdcBridge, bound to a specific deployed contract.
func NewL2UsdcBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*L2UsdcBridgeFilterer, error) {
	contract, err := bindL2UsdcBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeFilterer{contract: contract}, nil
}

// bindL2UsdcBridge binds a generic wrapper to an already deployed contract.
func bindL2UsdcBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2UsdcBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2UsdcBridge *L2UsdcBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2UsdcBridge.Contract.L2UsdcBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2UsdcBridge *L2UsdcBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.L2UsdcBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2UsdcBridge *L2UsdcBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.L2UsdcBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2UsdcBridge *L2UsdcBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2UsdcBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2UsdcBridge *L2UsdcBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2UsdcBridge *L2UsdcBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.contract.Transact(opts, method, params...)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCaller) L1Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridge.contract.Call(opts, &out, "l1Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeSession) L1Usdc() (common.Address, error) {
	return _L2UsdcBridge.Contract.L1Usdc(&_L2UsdcBridge.CallOpts)
}

// L1Usdc is a free data retrieval call binding the contract method 0x56c3b587.
//
// Solidity: function l1Usdc() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCallerSession) L1Usdc() (common.Address, error) {
	return _L2UsdcBridge.Contract.L1Usdc(&_L2UsdcBridge.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCaller) L2Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridge.contract.Call(opts, &out, "l2Usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeSession) L2Usdc() (common.Address, error) {
	return _L2UsdcBridge.Contract.L2Usdc(&_L2UsdcBridge.CallOpts)
}

// L2Usdc is a free data retrieval call binding the contract method 0xa1b4bc04.
//
// Solidity: function l2Usdc() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCallerSession) L2Usdc() (common.Address, error) {
	return _L2UsdcBridge.Contract.L2Usdc(&_L2UsdcBridge.CallOpts)
}

// L2UsdcMasterMinter is a free data retrieval call binding the contract method 0x05db940f.
//
// Solidity: function l2UsdcMasterMinter() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCaller) L2UsdcMasterMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridge.contract.Call(opts, &out, "l2UsdcMasterMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2UsdcMasterMinter is a free data retrieval call binding the contract method 0x05db940f.
//
// Solidity: function l2UsdcMasterMinter() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeSession) L2UsdcMasterMinter() (common.Address, error) {
	return _L2UsdcBridge.Contract.L2UsdcMasterMinter(&_L2UsdcBridge.CallOpts)
}

// L2UsdcMasterMinter is a free data retrieval call binding the contract method 0x05db940f.
//
// Solidity: function l2UsdcMasterMinter() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCallerSession) L2UsdcMasterMinter() (common.Address, error) {
	return _L2UsdcBridge.Contract.L2UsdcMasterMinter(&_L2UsdcBridge.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCaller) Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridge.contract.Call(opts, &out, "messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeSession) Messenger() (common.Address, error) {
	return _L2UsdcBridge.Contract.Messenger(&_L2UsdcBridge.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCallerSession) Messenger() (common.Address, error) {
	return _L2UsdcBridge.Contract.Messenger(&_L2UsdcBridge.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCaller) OtherBridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2UsdcBridge.contract.Call(opts, &out, "otherBridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeSession) OtherBridge() (common.Address, error) {
	return _L2UsdcBridge.Contract.OtherBridge(&_L2UsdcBridge.CallOpts)
}

// OtherBridge is a free data retrieval call binding the contract method 0xc89701a2.
//
// Solidity: function otherBridge() view returns(address)
func (_L2UsdcBridge *L2UsdcBridgeCallerSession) OtherBridge() (common.Address, error) {
	return _L2UsdcBridge.Contract.OtherBridge(&_L2UsdcBridge.CallOpts)
}

// FinalizeDeposit is a paid mutator transaction binding the contract method 0x662a633a.
//
// Solidity: function finalizeDeposit(address _l1Token, address _l2Token, address _from, address _to, uint256 _amount, bytes _extraData) returns()
func (_L2UsdcBridge *L2UsdcBridgeTransactor) FinalizeDeposit(opts *bind.TransactOpts, _l1Token common.Address, _l2Token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _L2UsdcBridge.contract.Transact(opts, "finalizeDeposit", _l1Token, _l2Token, _from, _to, _amount, _extraData)
}

// FinalizeDeposit is a paid mutator transaction binding the contract method 0x662a633a.
//
// Solidity: function finalizeDeposit(address _l1Token, address _l2Token, address _from, address _to, uint256 _amount, bytes _extraData) returns()
func (_L2UsdcBridge *L2UsdcBridgeSession) FinalizeDeposit(_l1Token common.Address, _l2Token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.FinalizeDeposit(&_L2UsdcBridge.TransactOpts, _l1Token, _l2Token, _from, _to, _amount, _extraData)
}

// FinalizeDeposit is a paid mutator transaction binding the contract method 0x662a633a.
//
// Solidity: function finalizeDeposit(address _l1Token, address _l2Token, address _from, address _to, uint256 _amount, bytes _extraData) returns()
func (_L2UsdcBridge *L2UsdcBridgeTransactorSession) FinalizeDeposit(_l1Token common.Address, _l2Token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.FinalizeDeposit(&_L2UsdcBridge.TransactOpts, _l1Token, _l2Token, _from, _to, _amount, _extraData)
}

// Withdraw is a paid mutator transaction binding the contract method 0x32b7006d.
//
// Solidity: function withdraw(address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2UsdcBridge *L2UsdcBridgeTransactor) Withdraw(opts *bind.TransactOpts, _l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2UsdcBridge.contract.Transact(opts, "withdraw", _l2Token, _amount, _minGasLimit, _extraData)
}

// Withdraw is a paid mutator transaction binding the contract method 0x32b7006d.
//
// Solidity: function withdraw(address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2UsdcBridge *L2UsdcBridgeSession) Withdraw(_l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.Withdraw(&_L2UsdcBridge.TransactOpts, _l2Token, _amount, _minGasLimit, _extraData)
}

// Withdraw is a paid mutator transaction binding the contract method 0x32b7006d.
//
// Solidity: function withdraw(address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2UsdcBridge *L2UsdcBridgeTransactorSession) Withdraw(_l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2UsdcBridge.Contract.Withdraw(&_L2UsdcBridge.TransactOpts, _l2Token, _amount, _minGasLimit, _extraData)
}

// L2UsdcBridgeDepositFinalizedIterator is returned from FilterDepositFinalized and is used to iterate over the raw logs and unpacked data for DepositFinalized events raised by the L2UsdcBridge contract.
type L2UsdcBridgeDepositFinalizedIterator struct {
	Event *L2UsdcBridgeDepositFinalized // Event containing the contract specifics and raw log

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
func (it *L2UsdcBridgeDepositFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2UsdcBridgeDepositFinalized)
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
		it.Event = new(L2UsdcBridgeDepositFinalized)
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
func (it *L2UsdcBridgeDepositFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2UsdcBridgeDepositFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2UsdcBridgeDepositFinalized represents a DepositFinalized event raised by the L2UsdcBridge contract.
type L2UsdcBridgeDepositFinalized struct {
	L1Token   common.Address
	L2Token   common.Address
	From      common.Address
	To        common.Address
	Amount    *big.Int
	ExtraData []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositFinalized is a free log retrieval operation binding the contract event 0xb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd89.
//
// Solidity: event DepositFinalized(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L2UsdcBridge *L2UsdcBridgeFilterer) FilterDepositFinalized(opts *bind.FilterOpts, l1Token []common.Address, l2Token []common.Address, from []common.Address) (*L2UsdcBridgeDepositFinalizedIterator, error) {

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

	logs, sub, err := _L2UsdcBridge.contract.FilterLogs(opts, "DepositFinalized", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeDepositFinalizedIterator{contract: _L2UsdcBridge.contract, event: "DepositFinalized", logs: logs, sub: sub}, nil
}

// WatchDepositFinalized is a free log subscription operation binding the contract event 0xb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd89.
//
// Solidity: event DepositFinalized(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L2UsdcBridge *L2UsdcBridgeFilterer) WatchDepositFinalized(opts *bind.WatchOpts, sink chan<- *L2UsdcBridgeDepositFinalized, l1Token []common.Address, l2Token []common.Address, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _L2UsdcBridge.contract.WatchLogs(opts, "DepositFinalized", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2UsdcBridgeDepositFinalized)
				if err := _L2UsdcBridge.contract.UnpackLog(event, "DepositFinalized", log); err != nil {
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

// ParseDepositFinalized is a log parse operation binding the contract event 0xb0444523268717a02698be47d0803aa7468c00acbed2f8bd93a0459cde61dd89.
//
// Solidity: event DepositFinalized(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L2UsdcBridge *L2UsdcBridgeFilterer) ParseDepositFinalized(log types.Log) (*L2UsdcBridgeDepositFinalized, error) {
	event := new(L2UsdcBridgeDepositFinalized)
	if err := _L2UsdcBridge.contract.UnpackLog(event, "DepositFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2UsdcBridgeWithdrawalInitiatedIterator is returned from FilterWithdrawalInitiated and is used to iterate over the raw logs and unpacked data for WithdrawalInitiated events raised by the L2UsdcBridge contract.
type L2UsdcBridgeWithdrawalInitiatedIterator struct {
	Event *L2UsdcBridgeWithdrawalInitiated // Event containing the contract specifics and raw log

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
func (it *L2UsdcBridgeWithdrawalInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2UsdcBridgeWithdrawalInitiated)
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
		it.Event = new(L2UsdcBridgeWithdrawalInitiated)
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
func (it *L2UsdcBridgeWithdrawalInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2UsdcBridgeWithdrawalInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2UsdcBridgeWithdrawalInitiated represents a WithdrawalInitiated event raised by the L2UsdcBridge contract.
type L2UsdcBridgeWithdrawalInitiated struct {
	L1Token   common.Address
	L2Token   common.Address
	From      common.Address
	To        common.Address
	Amount    *big.Int
	ExtraData []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalInitiated is a free log retrieval operation binding the contract event 0x73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e.
//
// Solidity: event WithdrawalInitiated(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L2UsdcBridge *L2UsdcBridgeFilterer) FilterWithdrawalInitiated(opts *bind.FilterOpts, l1Token []common.Address, l2Token []common.Address, from []common.Address) (*L2UsdcBridgeWithdrawalInitiatedIterator, error) {

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

	logs, sub, err := _L2UsdcBridge.contract.FilterLogs(opts, "WithdrawalInitiated", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &L2UsdcBridgeWithdrawalInitiatedIterator{contract: _L2UsdcBridge.contract, event: "WithdrawalInitiated", logs: logs, sub: sub}, nil
}

// WatchWithdrawalInitiated is a free log subscription operation binding the contract event 0x73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e.
//
// Solidity: event WithdrawalInitiated(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L2UsdcBridge *L2UsdcBridgeFilterer) WatchWithdrawalInitiated(opts *bind.WatchOpts, sink chan<- *L2UsdcBridgeWithdrawalInitiated, l1Token []common.Address, l2Token []common.Address, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _L2UsdcBridge.contract.WatchLogs(opts, "WithdrawalInitiated", l1TokenRule, l2TokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2UsdcBridgeWithdrawalInitiated)
				if err := _L2UsdcBridge.contract.UnpackLog(event, "WithdrawalInitiated", log); err != nil {
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

// ParseWithdrawalInitiated is a log parse operation binding the contract event 0x73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e.
//
// Solidity: event WithdrawalInitiated(address indexed l1Token, address indexed l2Token, address indexed from, address to, uint256 amount, bytes extraData)
func (_L2UsdcBridge *L2UsdcBridgeFilterer) ParseWithdrawalInitiated(log types.Log) (*L2UsdcBridgeWithdrawalInitiated, error) {
	event := new(L2UsdcBridgeWithdrawalInitiated)
	if err := _L2UsdcBridge.contract.UnpackLog(event, "WithdrawalInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
