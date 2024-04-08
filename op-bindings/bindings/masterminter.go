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

// MasterMinterMetaData contains all meta data concerning the MasterMinter contract.
var MasterMinterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_minterManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_controller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_worker\",\"type\":\"address\"}],\"name\":\"ControllerConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_controller\",\"type\":\"address\"}],\"name\":\"ControllerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decrement\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAllowance\",\"type\":\"uint256\"}],\"name\":\"MinterAllowanceDecremented\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_minter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_increment\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newAllowance\",\"type\":\"uint256\"}],\"name\":\"MinterAllowanceIncremented\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_minter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_allowance\",\"type\":\"uint256\"}],\"name\":\"MinterConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_oldMinterManager\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_newMinterManager\",\"type\":\"address\"}],\"name\":\"MinterManagerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_msgSender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_minter\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_controller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_worker\",\"type\":\"address\"}],\"name\":\"configureController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newAllowance\",\"type\":\"uint256\"}],\"name\":\"configureMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_allowanceDecrement\",\"type\":\"uint256\"}],\"name\":\"decrementMinterAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinterManager\",\"outputs\":[{\"internalType\":\"contractMinterManagementInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_controller\",\"type\":\"address\"}],\"name\":\"getWorker\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_allowanceIncrement\",\"type\":\"uint256\"}],\"name\":\"incrementMinterAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_controller\",\"type\":\"address\"}],\"name\":\"removeController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"removeMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newMinterManager\",\"type\":\"address\"}],\"name\":\"setMinterManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516115cb3803806115cb8339818101604052602081101561003357600080fd5b50518061003f33610065565b600280546001600160a01b0319166001600160a01b039290921691909117905550610087565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b611535806100966000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c8063c011b1c311610081578063ea7215691161005b578063ea72156914610215578063f2fde38b1461021d578063f6a74ed714610250576100c9565b8063c011b1c31461018a578063c4faf7df146101bd578063cbf2b8bf146101f8576100c9565b80637c6b8ef5116100b25780637c6b8ef5146101345780638da5cb5b146101515780639398608b14610182576100c9565b806333db2ad2146100ce578063542fef91146100ff575b600080fd5b6100eb600480360360208110156100e457600080fd5b5035610283565b604080519115158252519081900360200190f35b6101326004803603602081101561011557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610573565b005b6100eb6004803603602081101561014a57600080fd5b5035610687565b61015961098b565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101596109a7565b610159600480360360208110156101a057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166109c3565b610132600480360360408110156101d357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166109ee565b6100eb6004803603602081101561020e57600080fd5b5035610bc8565b6100eb610cb5565b6101326004803603602081101561023357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610e2b565b6101326004803603602081101561026657600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610f7e565b3360009081526001602052604081205473ffffffffffffffffffffffffffffffffffffffff166102fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806114586035913960400191505060405180910390fd5b60008211610357576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a81526020018061148d602a913960400191505060405180910390fd5b336000908152600160209081526040918290205460025483517faa271e1a00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff928316600482018190529451929091169263aa271e1a92602480840193829003018186803b1580156103d957600080fd5b505afa1580156103ed573d6000803e3d6000fd5b505050506040513d602081101561040357600080fd5b505161045a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260398152602001806114b76039913960400191505060405180910390fd5b600254604080517f8a6db9c300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015291516000939290921691638a6db9c391602480820192602092909190829003018186803b1580156104d157600080fd5b505afa1580156104e5573d6000803e3d6000fd5b505050506040513d60208110156104fb57600080fd5b50519050600061050b8286611161565b6040805187815260208101839052815192935073ffffffffffffffffffffffffffffffffffffffff86169233927f3703d23abba1e61f32acc0682fc062ea5c710672c7d100af5ecd08485e983ad0928290030190a361056a83826111d5565b95945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105f957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b60025460405173ffffffffffffffffffffffffffffffffffffffff8084169216907f9992ea32e96992be98be5c833cd5b9fd77314819d2146b6f06ab9cfef957af1290600090a3600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b3360009081526001602052604081205473ffffffffffffffffffffffffffffffffffffffff16610702576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806114586035913960400191505060405180910390fd5b6000821161075b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001806113c2602a913960400191505060405180910390fd5b336000908152600160209081526040918290205460025483517faa271e1a00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff928316600482018190529451929091169263aa271e1a92602480840193829003018186803b1580156107dd57600080fd5b505afa1580156107f1573d6000803e3d6000fd5b505050506040513d602081101561080757600080fd5b505161085e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260398152602001806114f06039913960400191505060405180910390fd5b600254604080517f8a6db9c300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015291516000939290921691638a6db9c391602480820192602092909190829003018186803b1580156108d557600080fd5b505afa1580156108e9573d6000803e3d6000fd5b505050506040513d60208110156108ff57600080fd5b5051905060008482116109125781610914565b845b905060006109228383611287565b6040805184815260208101839052815192935073ffffffffffffffffffffffffffffffffffffffff87169233927f3cc75d3bf58b0100659088c03539964108d5d06342e1bd8085ee43ad8ff6f69a928290030190a361098184826111d5565b9695505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1690565b60025473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff9081166000908152600160205260409020541690565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a7457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8216610ae0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806113ec6025913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116610b4c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806114376021913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82811660008181526001602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055517fa56687ff5096e83f6e2c673cda0b677f56bbfcdf5fe0555d5830c407ede193cb9190a35050565b3360009081526001602052604081205473ffffffffffffffffffffffffffffffffffffffff16610c43576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806114586035913960400191505060405180910390fd5b33600081815260016020908152604091829020548251868152925173ffffffffffffffffffffffffffffffffffffffff90911693849390927f5b0b60a4f757b33d9dcb8bd021b6aa371bb2e6f134086797aefcd8c0afab538c92918290030190a3610cae81846111d5565b9392505050565b3360009081526001602052604081205473ffffffffffffffffffffffffffffffffffffffff16610d30576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806114586035913960400191505060405180910390fd5b3360008181526001602052604080822054905173ffffffffffffffffffffffffffffffffffffffff90911692839290917f4b5ef9a786cf64a7d82ebcf2d5132667edc9faef4ac36260d9a9e52c526b62329190a3600254604080517f3092afd500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015291519190921691633092afd59160248083019260209291908290030181600087803b158015610df957600080fd5b505af1158015610e0d573d6000803e3d6000fd5b505050506040513d6020811015610e2357600080fd5b505191505090565b60005473ffffffffffffffffffffffffffffffffffffffff163314610eb157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8116610f1d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806114116026913960400191505060405180910390fd5b6000546040805173ffffffffffffffffffffffffffffffffffffffff9283168152918316602083015280517f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09281900390910190a1610f7b816112c9565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461100457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8116611070576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806113ec6025913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116600090815260016020526040902054166110ed576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806114376021913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811660008181526001602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055517f33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e811139190a250565b600082820183811015610cae57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600254604080517f4e44d95600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85811660048301526024820185905291516000939290921691634e44d9569160448082019260209290919082900301818787803b15801561125457600080fd5b505af1158015611268573d6000803e3d6000fd5b505050506040513d602081101561127e57600080fd5b50519392505050565b6000610cae83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f770000815250611310565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600081848411156113b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561137e578181015183820152602001611366565b50505050905090810190601f1680156113ab5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50505090039056fe416c6c6f77616e63652064656372656d656e74206d7573742062652067726561746572207468616e2030436f6e74726f6c6c6572206d7573742062652061206e6f6e2d7a65726f20616464726573734f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373576f726b6572206d7573742062652061206e6f6e2d7a65726f20616464726573735468652076616c7565206f6620636f6e74726f6c6c6572735b6d73672e73656e6465725d206d757374206265206e6f6e2d7a65726f416c6c6f77616e636520696e6372656d656e74206d7573742062652067726561746572207468616e203043616e206f6e6c7920696e6372656d656e7420616c6c6f77616e636520666f72206d696e7465727320696e206d696e7465724d616e6167657243616e206f6e6c792064656372656d656e7420616c6c6f77616e636520666f72206d696e7465727320696e206d696e7465724d616e61676572a164736f6c634300060c000a",
}

// MasterMinterABI is the input ABI used to generate the binding from.
// Deprecated: Use MasterMinterMetaData.ABI instead.
var MasterMinterABI = MasterMinterMetaData.ABI

// MasterMinterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MasterMinterMetaData.Bin instead.
var MasterMinterBin = MasterMinterMetaData.Bin

// DeployMasterMinter deploys a new Ethereum contract, binding an instance of MasterMinter to it.
func DeployMasterMinter(auth *bind.TransactOpts, backend bind.ContractBackend, _minterManager common.Address) (common.Address, *types.Transaction, *MasterMinter, error) {
	parsed, err := MasterMinterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MasterMinterBin), backend, _minterManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MasterMinter{MasterMinterCaller: MasterMinterCaller{contract: contract}, MasterMinterTransactor: MasterMinterTransactor{contract: contract}, MasterMinterFilterer: MasterMinterFilterer{contract: contract}}, nil
}

// MasterMinter is an auto generated Go binding around an Ethereum contract.
type MasterMinter struct {
	MasterMinterCaller     // Read-only binding to the contract
	MasterMinterTransactor // Write-only binding to the contract
	MasterMinterFilterer   // Log filterer for contract events
}

// MasterMinterCaller is an auto generated read-only Go binding around an Ethereum contract.
type MasterMinterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MasterMinterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MasterMinterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MasterMinterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MasterMinterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MasterMinterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MasterMinterSession struct {
	Contract     *MasterMinter     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MasterMinterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MasterMinterCallerSession struct {
	Contract *MasterMinterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// MasterMinterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MasterMinterTransactorSession struct {
	Contract     *MasterMinterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MasterMinterRaw is an auto generated low-level Go binding around an Ethereum contract.
type MasterMinterRaw struct {
	Contract *MasterMinter // Generic contract binding to access the raw methods on
}

// MasterMinterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MasterMinterCallerRaw struct {
	Contract *MasterMinterCaller // Generic read-only contract binding to access the raw methods on
}

// MasterMinterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MasterMinterTransactorRaw struct {
	Contract *MasterMinterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMasterMinter creates a new instance of MasterMinter, bound to a specific deployed contract.
func NewMasterMinter(address common.Address, backend bind.ContractBackend) (*MasterMinter, error) {
	contract, err := bindMasterMinter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MasterMinter{MasterMinterCaller: MasterMinterCaller{contract: contract}, MasterMinterTransactor: MasterMinterTransactor{contract: contract}, MasterMinterFilterer: MasterMinterFilterer{contract: contract}}, nil
}

// NewMasterMinterCaller creates a new read-only instance of MasterMinter, bound to a specific deployed contract.
func NewMasterMinterCaller(address common.Address, caller bind.ContractCaller) (*MasterMinterCaller, error) {
	contract, err := bindMasterMinter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MasterMinterCaller{contract: contract}, nil
}

// NewMasterMinterTransactor creates a new write-only instance of MasterMinter, bound to a specific deployed contract.
func NewMasterMinterTransactor(address common.Address, transactor bind.ContractTransactor) (*MasterMinterTransactor, error) {
	contract, err := bindMasterMinter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MasterMinterTransactor{contract: contract}, nil
}

// NewMasterMinterFilterer creates a new log filterer instance of MasterMinter, bound to a specific deployed contract.
func NewMasterMinterFilterer(address common.Address, filterer bind.ContractFilterer) (*MasterMinterFilterer, error) {
	contract, err := bindMasterMinter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MasterMinterFilterer{contract: contract}, nil
}

// bindMasterMinter binds a generic wrapper to an already deployed contract.
func bindMasterMinter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MasterMinterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MasterMinter *MasterMinterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MasterMinter.Contract.MasterMinterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MasterMinter *MasterMinterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MasterMinter.Contract.MasterMinterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MasterMinter *MasterMinterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MasterMinter.Contract.MasterMinterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MasterMinter *MasterMinterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MasterMinter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MasterMinter *MasterMinterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MasterMinter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MasterMinter *MasterMinterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MasterMinter.Contract.contract.Transact(opts, method, params...)
}

// GetMinterManager is a free data retrieval call binding the contract method 0x9398608b.
//
// Solidity: function getMinterManager() view returns(address)
func (_MasterMinter *MasterMinterCaller) GetMinterManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MasterMinter.contract.Call(opts, &out, "getMinterManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMinterManager is a free data retrieval call binding the contract method 0x9398608b.
//
// Solidity: function getMinterManager() view returns(address)
func (_MasterMinter *MasterMinterSession) GetMinterManager() (common.Address, error) {
	return _MasterMinter.Contract.GetMinterManager(&_MasterMinter.CallOpts)
}

// GetMinterManager is a free data retrieval call binding the contract method 0x9398608b.
//
// Solidity: function getMinterManager() view returns(address)
func (_MasterMinter *MasterMinterCallerSession) GetMinterManager() (common.Address, error) {
	return _MasterMinter.Contract.GetMinterManager(&_MasterMinter.CallOpts)
}

// GetWorker is a free data retrieval call binding the contract method 0xc011b1c3.
//
// Solidity: function getWorker(address _controller) view returns(address)
func (_MasterMinter *MasterMinterCaller) GetWorker(opts *bind.CallOpts, _controller common.Address) (common.Address, error) {
	var out []interface{}
	err := _MasterMinter.contract.Call(opts, &out, "getWorker", _controller)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetWorker is a free data retrieval call binding the contract method 0xc011b1c3.
//
// Solidity: function getWorker(address _controller) view returns(address)
func (_MasterMinter *MasterMinterSession) GetWorker(_controller common.Address) (common.Address, error) {
	return _MasterMinter.Contract.GetWorker(&_MasterMinter.CallOpts, _controller)
}

// GetWorker is a free data retrieval call binding the contract method 0xc011b1c3.
//
// Solidity: function getWorker(address _controller) view returns(address)
func (_MasterMinter *MasterMinterCallerSession) GetWorker(_controller common.Address) (common.Address, error) {
	return _MasterMinter.Contract.GetWorker(&_MasterMinter.CallOpts, _controller)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MasterMinter *MasterMinterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MasterMinter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MasterMinter *MasterMinterSession) Owner() (common.Address, error) {
	return _MasterMinter.Contract.Owner(&_MasterMinter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MasterMinter *MasterMinterCallerSession) Owner() (common.Address, error) {
	return _MasterMinter.Contract.Owner(&_MasterMinter.CallOpts)
}

// ConfigureController is a paid mutator transaction binding the contract method 0xc4faf7df.
//
// Solidity: function configureController(address _controller, address _worker) returns()
func (_MasterMinter *MasterMinterTransactor) ConfigureController(opts *bind.TransactOpts, _controller common.Address, _worker common.Address) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "configureController", _controller, _worker)
}

// ConfigureController is a paid mutator transaction binding the contract method 0xc4faf7df.
//
// Solidity: function configureController(address _controller, address _worker) returns()
func (_MasterMinter *MasterMinterSession) ConfigureController(_controller common.Address, _worker common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.ConfigureController(&_MasterMinter.TransactOpts, _controller, _worker)
}

// ConfigureController is a paid mutator transaction binding the contract method 0xc4faf7df.
//
// Solidity: function configureController(address _controller, address _worker) returns()
func (_MasterMinter *MasterMinterTransactorSession) ConfigureController(_controller common.Address, _worker common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.ConfigureController(&_MasterMinter.TransactOpts, _controller, _worker)
}

// ConfigureMinter is a paid mutator transaction binding the contract method 0xcbf2b8bf.
//
// Solidity: function configureMinter(uint256 _newAllowance) returns(bool)
func (_MasterMinter *MasterMinterTransactor) ConfigureMinter(opts *bind.TransactOpts, _newAllowance *big.Int) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "configureMinter", _newAllowance)
}

// ConfigureMinter is a paid mutator transaction binding the contract method 0xcbf2b8bf.
//
// Solidity: function configureMinter(uint256 _newAllowance) returns(bool)
func (_MasterMinter *MasterMinterSession) ConfigureMinter(_newAllowance *big.Int) (*types.Transaction, error) {
	return _MasterMinter.Contract.ConfigureMinter(&_MasterMinter.TransactOpts, _newAllowance)
}

// ConfigureMinter is a paid mutator transaction binding the contract method 0xcbf2b8bf.
//
// Solidity: function configureMinter(uint256 _newAllowance) returns(bool)
func (_MasterMinter *MasterMinterTransactorSession) ConfigureMinter(_newAllowance *big.Int) (*types.Transaction, error) {
	return _MasterMinter.Contract.ConfigureMinter(&_MasterMinter.TransactOpts, _newAllowance)
}

// DecrementMinterAllowance is a paid mutator transaction binding the contract method 0x7c6b8ef5.
//
// Solidity: function decrementMinterAllowance(uint256 _allowanceDecrement) returns(bool)
func (_MasterMinter *MasterMinterTransactor) DecrementMinterAllowance(opts *bind.TransactOpts, _allowanceDecrement *big.Int) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "decrementMinterAllowance", _allowanceDecrement)
}

// DecrementMinterAllowance is a paid mutator transaction binding the contract method 0x7c6b8ef5.
//
// Solidity: function decrementMinterAllowance(uint256 _allowanceDecrement) returns(bool)
func (_MasterMinter *MasterMinterSession) DecrementMinterAllowance(_allowanceDecrement *big.Int) (*types.Transaction, error) {
	return _MasterMinter.Contract.DecrementMinterAllowance(&_MasterMinter.TransactOpts, _allowanceDecrement)
}

// DecrementMinterAllowance is a paid mutator transaction binding the contract method 0x7c6b8ef5.
//
// Solidity: function decrementMinterAllowance(uint256 _allowanceDecrement) returns(bool)
func (_MasterMinter *MasterMinterTransactorSession) DecrementMinterAllowance(_allowanceDecrement *big.Int) (*types.Transaction, error) {
	return _MasterMinter.Contract.DecrementMinterAllowance(&_MasterMinter.TransactOpts, _allowanceDecrement)
}

// IncrementMinterAllowance is a paid mutator transaction binding the contract method 0x33db2ad2.
//
// Solidity: function incrementMinterAllowance(uint256 _allowanceIncrement) returns(bool)
func (_MasterMinter *MasterMinterTransactor) IncrementMinterAllowance(opts *bind.TransactOpts, _allowanceIncrement *big.Int) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "incrementMinterAllowance", _allowanceIncrement)
}

// IncrementMinterAllowance is a paid mutator transaction binding the contract method 0x33db2ad2.
//
// Solidity: function incrementMinterAllowance(uint256 _allowanceIncrement) returns(bool)
func (_MasterMinter *MasterMinterSession) IncrementMinterAllowance(_allowanceIncrement *big.Int) (*types.Transaction, error) {
	return _MasterMinter.Contract.IncrementMinterAllowance(&_MasterMinter.TransactOpts, _allowanceIncrement)
}

// IncrementMinterAllowance is a paid mutator transaction binding the contract method 0x33db2ad2.
//
// Solidity: function incrementMinterAllowance(uint256 _allowanceIncrement) returns(bool)
func (_MasterMinter *MasterMinterTransactorSession) IncrementMinterAllowance(_allowanceIncrement *big.Int) (*types.Transaction, error) {
	return _MasterMinter.Contract.IncrementMinterAllowance(&_MasterMinter.TransactOpts, _allowanceIncrement)
}

// RemoveController is a paid mutator transaction binding the contract method 0xf6a74ed7.
//
// Solidity: function removeController(address _controller) returns()
func (_MasterMinter *MasterMinterTransactor) RemoveController(opts *bind.TransactOpts, _controller common.Address) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "removeController", _controller)
}

// RemoveController is a paid mutator transaction binding the contract method 0xf6a74ed7.
//
// Solidity: function removeController(address _controller) returns()
func (_MasterMinter *MasterMinterSession) RemoveController(_controller common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.RemoveController(&_MasterMinter.TransactOpts, _controller)
}

// RemoveController is a paid mutator transaction binding the contract method 0xf6a74ed7.
//
// Solidity: function removeController(address _controller) returns()
func (_MasterMinter *MasterMinterTransactorSession) RemoveController(_controller common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.RemoveController(&_MasterMinter.TransactOpts, _controller)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0xea721569.
//
// Solidity: function removeMinter() returns(bool)
func (_MasterMinter *MasterMinterTransactor) RemoveMinter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "removeMinter")
}

// RemoveMinter is a paid mutator transaction binding the contract method 0xea721569.
//
// Solidity: function removeMinter() returns(bool)
func (_MasterMinter *MasterMinterSession) RemoveMinter() (*types.Transaction, error) {
	return _MasterMinter.Contract.RemoveMinter(&_MasterMinter.TransactOpts)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0xea721569.
//
// Solidity: function removeMinter() returns(bool)
func (_MasterMinter *MasterMinterTransactorSession) RemoveMinter() (*types.Transaction, error) {
	return _MasterMinter.Contract.RemoveMinter(&_MasterMinter.TransactOpts)
}

// SetMinterManager is a paid mutator transaction binding the contract method 0x542fef91.
//
// Solidity: function setMinterManager(address _newMinterManager) returns()
func (_MasterMinter *MasterMinterTransactor) SetMinterManager(opts *bind.TransactOpts, _newMinterManager common.Address) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "setMinterManager", _newMinterManager)
}

// SetMinterManager is a paid mutator transaction binding the contract method 0x542fef91.
//
// Solidity: function setMinterManager(address _newMinterManager) returns()
func (_MasterMinter *MasterMinterSession) SetMinterManager(_newMinterManager common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.SetMinterManager(&_MasterMinter.TransactOpts, _newMinterManager)
}

// SetMinterManager is a paid mutator transaction binding the contract method 0x542fef91.
//
// Solidity: function setMinterManager(address _newMinterManager) returns()
func (_MasterMinter *MasterMinterTransactorSession) SetMinterManager(_newMinterManager common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.SetMinterManager(&_MasterMinter.TransactOpts, _newMinterManager)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MasterMinter *MasterMinterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MasterMinter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MasterMinter *MasterMinterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.TransferOwnership(&_MasterMinter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MasterMinter *MasterMinterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MasterMinter.Contract.TransferOwnership(&_MasterMinter.TransactOpts, newOwner)
}

// MasterMinterControllerConfiguredIterator is returned from FilterControllerConfigured and is used to iterate over the raw logs and unpacked data for ControllerConfigured events raised by the MasterMinter contract.
type MasterMinterControllerConfiguredIterator struct {
	Event *MasterMinterControllerConfigured // Event containing the contract specifics and raw log

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
func (it *MasterMinterControllerConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterControllerConfigured)
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
		it.Event = new(MasterMinterControllerConfigured)
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
func (it *MasterMinterControllerConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterControllerConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterControllerConfigured represents a ControllerConfigured event raised by the MasterMinter contract.
type MasterMinterControllerConfigured struct {
	Controller common.Address
	Worker     common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterControllerConfigured is a free log retrieval operation binding the contract event 0xa56687ff5096e83f6e2c673cda0b677f56bbfcdf5fe0555d5830c407ede193cb.
//
// Solidity: event ControllerConfigured(address indexed _controller, address indexed _worker)
func (_MasterMinter *MasterMinterFilterer) FilterControllerConfigured(opts *bind.FilterOpts, _controller []common.Address, _worker []common.Address) (*MasterMinterControllerConfiguredIterator, error) {

	var _controllerRule []interface{}
	for _, _controllerItem := range _controller {
		_controllerRule = append(_controllerRule, _controllerItem)
	}
	var _workerRule []interface{}
	for _, _workerItem := range _worker {
		_workerRule = append(_workerRule, _workerItem)
	}

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "ControllerConfigured", _controllerRule, _workerRule)
	if err != nil {
		return nil, err
	}
	return &MasterMinterControllerConfiguredIterator{contract: _MasterMinter.contract, event: "ControllerConfigured", logs: logs, sub: sub}, nil
}

// WatchControllerConfigured is a free log subscription operation binding the contract event 0xa56687ff5096e83f6e2c673cda0b677f56bbfcdf5fe0555d5830c407ede193cb.
//
// Solidity: event ControllerConfigured(address indexed _controller, address indexed _worker)
func (_MasterMinter *MasterMinterFilterer) WatchControllerConfigured(opts *bind.WatchOpts, sink chan<- *MasterMinterControllerConfigured, _controller []common.Address, _worker []common.Address) (event.Subscription, error) {

	var _controllerRule []interface{}
	for _, _controllerItem := range _controller {
		_controllerRule = append(_controllerRule, _controllerItem)
	}
	var _workerRule []interface{}
	for _, _workerItem := range _worker {
		_workerRule = append(_workerRule, _workerItem)
	}

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "ControllerConfigured", _controllerRule, _workerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterControllerConfigured)
				if err := _MasterMinter.contract.UnpackLog(event, "ControllerConfigured", log); err != nil {
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

// ParseControllerConfigured is a log parse operation binding the contract event 0xa56687ff5096e83f6e2c673cda0b677f56bbfcdf5fe0555d5830c407ede193cb.
//
// Solidity: event ControllerConfigured(address indexed _controller, address indexed _worker)
func (_MasterMinter *MasterMinterFilterer) ParseControllerConfigured(log types.Log) (*MasterMinterControllerConfigured, error) {
	event := new(MasterMinterControllerConfigured)
	if err := _MasterMinter.contract.UnpackLog(event, "ControllerConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MasterMinterControllerRemovedIterator is returned from FilterControllerRemoved and is used to iterate over the raw logs and unpacked data for ControllerRemoved events raised by the MasterMinter contract.
type MasterMinterControllerRemovedIterator struct {
	Event *MasterMinterControllerRemoved // Event containing the contract specifics and raw log

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
func (it *MasterMinterControllerRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterControllerRemoved)
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
		it.Event = new(MasterMinterControllerRemoved)
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
func (it *MasterMinterControllerRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterControllerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterControllerRemoved represents a ControllerRemoved event raised by the MasterMinter contract.
type MasterMinterControllerRemoved struct {
	Controller common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterControllerRemoved is a free log retrieval operation binding the contract event 0x33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e81113.
//
// Solidity: event ControllerRemoved(address indexed _controller)
func (_MasterMinter *MasterMinterFilterer) FilterControllerRemoved(opts *bind.FilterOpts, _controller []common.Address) (*MasterMinterControllerRemovedIterator, error) {

	var _controllerRule []interface{}
	for _, _controllerItem := range _controller {
		_controllerRule = append(_controllerRule, _controllerItem)
	}

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "ControllerRemoved", _controllerRule)
	if err != nil {
		return nil, err
	}
	return &MasterMinterControllerRemovedIterator{contract: _MasterMinter.contract, event: "ControllerRemoved", logs: logs, sub: sub}, nil
}

// WatchControllerRemoved is a free log subscription operation binding the contract event 0x33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e81113.
//
// Solidity: event ControllerRemoved(address indexed _controller)
func (_MasterMinter *MasterMinterFilterer) WatchControllerRemoved(opts *bind.WatchOpts, sink chan<- *MasterMinterControllerRemoved, _controller []common.Address) (event.Subscription, error) {

	var _controllerRule []interface{}
	for _, _controllerItem := range _controller {
		_controllerRule = append(_controllerRule, _controllerItem)
	}

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "ControllerRemoved", _controllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterControllerRemoved)
				if err := _MasterMinter.contract.UnpackLog(event, "ControllerRemoved", log); err != nil {
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

// ParseControllerRemoved is a log parse operation binding the contract event 0x33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e81113.
//
// Solidity: event ControllerRemoved(address indexed _controller)
func (_MasterMinter *MasterMinterFilterer) ParseControllerRemoved(log types.Log) (*MasterMinterControllerRemoved, error) {
	event := new(MasterMinterControllerRemoved)
	if err := _MasterMinter.contract.UnpackLog(event, "ControllerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MasterMinterMinterAllowanceDecrementedIterator is returned from FilterMinterAllowanceDecremented and is used to iterate over the raw logs and unpacked data for MinterAllowanceDecremented events raised by the MasterMinter contract.
type MasterMinterMinterAllowanceDecrementedIterator struct {
	Event *MasterMinterMinterAllowanceDecremented // Event containing the contract specifics and raw log

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
func (it *MasterMinterMinterAllowanceDecrementedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterMinterAllowanceDecremented)
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
		it.Event = new(MasterMinterMinterAllowanceDecremented)
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
func (it *MasterMinterMinterAllowanceDecrementedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterMinterAllowanceDecrementedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterMinterAllowanceDecremented represents a MinterAllowanceDecremented event raised by the MasterMinter contract.
type MasterMinterMinterAllowanceDecremented struct {
	MsgSender    common.Address
	Minter       common.Address
	Decrement    *big.Int
	NewAllowance *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinterAllowanceDecremented is a free log retrieval operation binding the contract event 0x3cc75d3bf58b0100659088c03539964108d5d06342e1bd8085ee43ad8ff6f69a.
//
// Solidity: event MinterAllowanceDecremented(address indexed msgSender, address indexed minter, uint256 decrement, uint256 newAllowance)
func (_MasterMinter *MasterMinterFilterer) FilterMinterAllowanceDecremented(opts *bind.FilterOpts, msgSender []common.Address, minter []common.Address) (*MasterMinterMinterAllowanceDecrementedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "MinterAllowanceDecremented", msgSenderRule, minterRule)
	if err != nil {
		return nil, err
	}
	return &MasterMinterMinterAllowanceDecrementedIterator{contract: _MasterMinter.contract, event: "MinterAllowanceDecremented", logs: logs, sub: sub}, nil
}

// WatchMinterAllowanceDecremented is a free log subscription operation binding the contract event 0x3cc75d3bf58b0100659088c03539964108d5d06342e1bd8085ee43ad8ff6f69a.
//
// Solidity: event MinterAllowanceDecremented(address indexed msgSender, address indexed minter, uint256 decrement, uint256 newAllowance)
func (_MasterMinter *MasterMinterFilterer) WatchMinterAllowanceDecremented(opts *bind.WatchOpts, sink chan<- *MasterMinterMinterAllowanceDecremented, msgSender []common.Address, minter []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "MinterAllowanceDecremented", msgSenderRule, minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterMinterAllowanceDecremented)
				if err := _MasterMinter.contract.UnpackLog(event, "MinterAllowanceDecremented", log); err != nil {
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

// ParseMinterAllowanceDecremented is a log parse operation binding the contract event 0x3cc75d3bf58b0100659088c03539964108d5d06342e1bd8085ee43ad8ff6f69a.
//
// Solidity: event MinterAllowanceDecremented(address indexed msgSender, address indexed minter, uint256 decrement, uint256 newAllowance)
func (_MasterMinter *MasterMinterFilterer) ParseMinterAllowanceDecremented(log types.Log) (*MasterMinterMinterAllowanceDecremented, error) {
	event := new(MasterMinterMinterAllowanceDecremented)
	if err := _MasterMinter.contract.UnpackLog(event, "MinterAllowanceDecremented", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MasterMinterMinterAllowanceIncrementedIterator is returned from FilterMinterAllowanceIncremented and is used to iterate over the raw logs and unpacked data for MinterAllowanceIncremented events raised by the MasterMinter contract.
type MasterMinterMinterAllowanceIncrementedIterator struct {
	Event *MasterMinterMinterAllowanceIncremented // Event containing the contract specifics and raw log

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
func (it *MasterMinterMinterAllowanceIncrementedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterMinterAllowanceIncremented)
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
		it.Event = new(MasterMinterMinterAllowanceIncremented)
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
func (it *MasterMinterMinterAllowanceIncrementedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterMinterAllowanceIncrementedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterMinterAllowanceIncremented represents a MinterAllowanceIncremented event raised by the MasterMinter contract.
type MasterMinterMinterAllowanceIncremented struct {
	MsgSender    common.Address
	Minter       common.Address
	Increment    *big.Int
	NewAllowance *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinterAllowanceIncremented is a free log retrieval operation binding the contract event 0x3703d23abba1e61f32acc0682fc062ea5c710672c7d100af5ecd08485e983ad0.
//
// Solidity: event MinterAllowanceIncremented(address indexed _msgSender, address indexed _minter, uint256 _increment, uint256 _newAllowance)
func (_MasterMinter *MasterMinterFilterer) FilterMinterAllowanceIncremented(opts *bind.FilterOpts, _msgSender []common.Address, _minter []common.Address) (*MasterMinterMinterAllowanceIncrementedIterator, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _minterRule []interface{}
	for _, _minterItem := range _minter {
		_minterRule = append(_minterRule, _minterItem)
	}

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "MinterAllowanceIncremented", _msgSenderRule, _minterRule)
	if err != nil {
		return nil, err
	}
	return &MasterMinterMinterAllowanceIncrementedIterator{contract: _MasterMinter.contract, event: "MinterAllowanceIncremented", logs: logs, sub: sub}, nil
}

// WatchMinterAllowanceIncremented is a free log subscription operation binding the contract event 0x3703d23abba1e61f32acc0682fc062ea5c710672c7d100af5ecd08485e983ad0.
//
// Solidity: event MinterAllowanceIncremented(address indexed _msgSender, address indexed _minter, uint256 _increment, uint256 _newAllowance)
func (_MasterMinter *MasterMinterFilterer) WatchMinterAllowanceIncremented(opts *bind.WatchOpts, sink chan<- *MasterMinterMinterAllowanceIncremented, _msgSender []common.Address, _minter []common.Address) (event.Subscription, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _minterRule []interface{}
	for _, _minterItem := range _minter {
		_minterRule = append(_minterRule, _minterItem)
	}

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "MinterAllowanceIncremented", _msgSenderRule, _minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterMinterAllowanceIncremented)
				if err := _MasterMinter.contract.UnpackLog(event, "MinterAllowanceIncremented", log); err != nil {
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

// ParseMinterAllowanceIncremented is a log parse operation binding the contract event 0x3703d23abba1e61f32acc0682fc062ea5c710672c7d100af5ecd08485e983ad0.
//
// Solidity: event MinterAllowanceIncremented(address indexed _msgSender, address indexed _minter, uint256 _increment, uint256 _newAllowance)
func (_MasterMinter *MasterMinterFilterer) ParseMinterAllowanceIncremented(log types.Log) (*MasterMinterMinterAllowanceIncremented, error) {
	event := new(MasterMinterMinterAllowanceIncremented)
	if err := _MasterMinter.contract.UnpackLog(event, "MinterAllowanceIncremented", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MasterMinterMinterConfiguredIterator is returned from FilterMinterConfigured and is used to iterate over the raw logs and unpacked data for MinterConfigured events raised by the MasterMinter contract.
type MasterMinterMinterConfiguredIterator struct {
	Event *MasterMinterMinterConfigured // Event containing the contract specifics and raw log

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
func (it *MasterMinterMinterConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterMinterConfigured)
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
		it.Event = new(MasterMinterMinterConfigured)
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
func (it *MasterMinterMinterConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterMinterConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterMinterConfigured represents a MinterConfigured event raised by the MasterMinter contract.
type MasterMinterMinterConfigured struct {
	MsgSender common.Address
	Minter    common.Address
	Allowance *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMinterConfigured is a free log retrieval operation binding the contract event 0x5b0b60a4f757b33d9dcb8bd021b6aa371bb2e6f134086797aefcd8c0afab538c.
//
// Solidity: event MinterConfigured(address indexed _msgSender, address indexed _minter, uint256 _allowance)
func (_MasterMinter *MasterMinterFilterer) FilterMinterConfigured(opts *bind.FilterOpts, _msgSender []common.Address, _minter []common.Address) (*MasterMinterMinterConfiguredIterator, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _minterRule []interface{}
	for _, _minterItem := range _minter {
		_minterRule = append(_minterRule, _minterItem)
	}

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "MinterConfigured", _msgSenderRule, _minterRule)
	if err != nil {
		return nil, err
	}
	return &MasterMinterMinterConfiguredIterator{contract: _MasterMinter.contract, event: "MinterConfigured", logs: logs, sub: sub}, nil
}

// WatchMinterConfigured is a free log subscription operation binding the contract event 0x5b0b60a4f757b33d9dcb8bd021b6aa371bb2e6f134086797aefcd8c0afab538c.
//
// Solidity: event MinterConfigured(address indexed _msgSender, address indexed _minter, uint256 _allowance)
func (_MasterMinter *MasterMinterFilterer) WatchMinterConfigured(opts *bind.WatchOpts, sink chan<- *MasterMinterMinterConfigured, _msgSender []common.Address, _minter []common.Address) (event.Subscription, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _minterRule []interface{}
	for _, _minterItem := range _minter {
		_minterRule = append(_minterRule, _minterItem)
	}

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "MinterConfigured", _msgSenderRule, _minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterMinterConfigured)
				if err := _MasterMinter.contract.UnpackLog(event, "MinterConfigured", log); err != nil {
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

// ParseMinterConfigured is a log parse operation binding the contract event 0x5b0b60a4f757b33d9dcb8bd021b6aa371bb2e6f134086797aefcd8c0afab538c.
//
// Solidity: event MinterConfigured(address indexed _msgSender, address indexed _minter, uint256 _allowance)
func (_MasterMinter *MasterMinterFilterer) ParseMinterConfigured(log types.Log) (*MasterMinterMinterConfigured, error) {
	event := new(MasterMinterMinterConfigured)
	if err := _MasterMinter.contract.UnpackLog(event, "MinterConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MasterMinterMinterManagerSetIterator is returned from FilterMinterManagerSet and is used to iterate over the raw logs and unpacked data for MinterManagerSet events raised by the MasterMinter contract.
type MasterMinterMinterManagerSetIterator struct {
	Event *MasterMinterMinterManagerSet // Event containing the contract specifics and raw log

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
func (it *MasterMinterMinterManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterMinterManagerSet)
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
		it.Event = new(MasterMinterMinterManagerSet)
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
func (it *MasterMinterMinterManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterMinterManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterMinterManagerSet represents a MinterManagerSet event raised by the MasterMinter contract.
type MasterMinterMinterManagerSet struct {
	OldMinterManager common.Address
	NewMinterManager common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterMinterManagerSet is a free log retrieval operation binding the contract event 0x9992ea32e96992be98be5c833cd5b9fd77314819d2146b6f06ab9cfef957af12.
//
// Solidity: event MinterManagerSet(address indexed _oldMinterManager, address indexed _newMinterManager)
func (_MasterMinter *MasterMinterFilterer) FilterMinterManagerSet(opts *bind.FilterOpts, _oldMinterManager []common.Address, _newMinterManager []common.Address) (*MasterMinterMinterManagerSetIterator, error) {

	var _oldMinterManagerRule []interface{}
	for _, _oldMinterManagerItem := range _oldMinterManager {
		_oldMinterManagerRule = append(_oldMinterManagerRule, _oldMinterManagerItem)
	}
	var _newMinterManagerRule []interface{}
	for _, _newMinterManagerItem := range _newMinterManager {
		_newMinterManagerRule = append(_newMinterManagerRule, _newMinterManagerItem)
	}

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "MinterManagerSet", _oldMinterManagerRule, _newMinterManagerRule)
	if err != nil {
		return nil, err
	}
	return &MasterMinterMinterManagerSetIterator{contract: _MasterMinter.contract, event: "MinterManagerSet", logs: logs, sub: sub}, nil
}

// WatchMinterManagerSet is a free log subscription operation binding the contract event 0x9992ea32e96992be98be5c833cd5b9fd77314819d2146b6f06ab9cfef957af12.
//
// Solidity: event MinterManagerSet(address indexed _oldMinterManager, address indexed _newMinterManager)
func (_MasterMinter *MasterMinterFilterer) WatchMinterManagerSet(opts *bind.WatchOpts, sink chan<- *MasterMinterMinterManagerSet, _oldMinterManager []common.Address, _newMinterManager []common.Address) (event.Subscription, error) {

	var _oldMinterManagerRule []interface{}
	for _, _oldMinterManagerItem := range _oldMinterManager {
		_oldMinterManagerRule = append(_oldMinterManagerRule, _oldMinterManagerItem)
	}
	var _newMinterManagerRule []interface{}
	for _, _newMinterManagerItem := range _newMinterManager {
		_newMinterManagerRule = append(_newMinterManagerRule, _newMinterManagerItem)
	}

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "MinterManagerSet", _oldMinterManagerRule, _newMinterManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterMinterManagerSet)
				if err := _MasterMinter.contract.UnpackLog(event, "MinterManagerSet", log); err != nil {
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

// ParseMinterManagerSet is a log parse operation binding the contract event 0x9992ea32e96992be98be5c833cd5b9fd77314819d2146b6f06ab9cfef957af12.
//
// Solidity: event MinterManagerSet(address indexed _oldMinterManager, address indexed _newMinterManager)
func (_MasterMinter *MasterMinterFilterer) ParseMinterManagerSet(log types.Log) (*MasterMinterMinterManagerSet, error) {
	event := new(MasterMinterMinterManagerSet)
	if err := _MasterMinter.contract.UnpackLog(event, "MinterManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MasterMinterMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the MasterMinter contract.
type MasterMinterMinterRemovedIterator struct {
	Event *MasterMinterMinterRemoved // Event containing the contract specifics and raw log

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
func (it *MasterMinterMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterMinterRemoved)
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
		it.Event = new(MasterMinterMinterRemoved)
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
func (it *MasterMinterMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterMinterRemoved represents a MinterRemoved event raised by the MasterMinter contract.
type MasterMinterMinterRemoved struct {
	MsgSender common.Address
	Minter    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0x4b5ef9a786cf64a7d82ebcf2d5132667edc9faef4ac36260d9a9e52c526b6232.
//
// Solidity: event MinterRemoved(address indexed _msgSender, address indexed _minter)
func (_MasterMinter *MasterMinterFilterer) FilterMinterRemoved(opts *bind.FilterOpts, _msgSender []common.Address, _minter []common.Address) (*MasterMinterMinterRemovedIterator, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _minterRule []interface{}
	for _, _minterItem := range _minter {
		_minterRule = append(_minterRule, _minterItem)
	}

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "MinterRemoved", _msgSenderRule, _minterRule)
	if err != nil {
		return nil, err
	}
	return &MasterMinterMinterRemovedIterator{contract: _MasterMinter.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0x4b5ef9a786cf64a7d82ebcf2d5132667edc9faef4ac36260d9a9e52c526b6232.
//
// Solidity: event MinterRemoved(address indexed _msgSender, address indexed _minter)
func (_MasterMinter *MasterMinterFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *MasterMinterMinterRemoved, _msgSender []common.Address, _minter []common.Address) (event.Subscription, error) {

	var _msgSenderRule []interface{}
	for _, _msgSenderItem := range _msgSender {
		_msgSenderRule = append(_msgSenderRule, _msgSenderItem)
	}
	var _minterRule []interface{}
	for _, _minterItem := range _minter {
		_minterRule = append(_minterRule, _minterItem)
	}

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "MinterRemoved", _msgSenderRule, _minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterMinterRemoved)
				if err := _MasterMinter.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// ParseMinterRemoved is a log parse operation binding the contract event 0x4b5ef9a786cf64a7d82ebcf2d5132667edc9faef4ac36260d9a9e52c526b6232.
//
// Solidity: event MinterRemoved(address indexed _msgSender, address indexed _minter)
func (_MasterMinter *MasterMinterFilterer) ParseMinterRemoved(log types.Log) (*MasterMinterMinterRemoved, error) {
	event := new(MasterMinterMinterRemoved)
	if err := _MasterMinter.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MasterMinterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MasterMinter contract.
type MasterMinterOwnershipTransferredIterator struct {
	Event *MasterMinterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MasterMinterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MasterMinterOwnershipTransferred)
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
		it.Event = new(MasterMinterOwnershipTransferred)
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
func (it *MasterMinterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MasterMinterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MasterMinterOwnershipTransferred represents a OwnershipTransferred event raised by the MasterMinter contract.
type MasterMinterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address previousOwner, address newOwner)
func (_MasterMinter *MasterMinterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts) (*MasterMinterOwnershipTransferredIterator, error) {

	logs, sub, err := _MasterMinter.contract.FilterLogs(opts, "OwnershipTransferred")
	if err != nil {
		return nil, err
	}
	return &MasterMinterOwnershipTransferredIterator{contract: _MasterMinter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address previousOwner, address newOwner)
func (_MasterMinter *MasterMinterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MasterMinterOwnershipTransferred) (event.Subscription, error) {

	logs, sub, err := _MasterMinter.contract.WatchLogs(opts, "OwnershipTransferred")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MasterMinterOwnershipTransferred)
				if err := _MasterMinter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
// Solidity: event OwnershipTransferred(address previousOwner, address newOwner)
func (_MasterMinter *MasterMinterFilterer) ParseOwnershipTransferred(log types.Log) (*MasterMinterOwnershipTransferred, error) {
	event := new(MasterMinterOwnershipTransferred)
	if err := _MasterMinter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
