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

// StdInvariantFuzzSelector is an auto generated low-level Go binding around an user-defined struct.
type StdInvariantFuzzSelector struct {
	Addr      common.Address
	Selectors [][4]byte
}

// DeployerWhitelistMetaData contains all meta data concerning the DeployerWhitelist contract.
var DeployerWhitelistMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"opaqueData\",\"type\":\"bytes\"}],\"name\":\"TransactionDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"log_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"val\",\"type\":\"uint256[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"val\",\"type\":\"int256[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"val\",\"type\":\"address[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"log_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"log_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"log_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"val\",\"type\":\"address\"}],\"name\":\"log_named_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"val\",\"type\":\"uint256[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"val\",\"type\":\"int256[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"val\",\"type\":\"address[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"val\",\"type\":\"bytes\"}],\"name\":\"log_named_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"val\",\"type\":\"bytes32\"}],\"name\":\"log_named_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"}],\"name\":\"log_named_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"val\",\"type\":\"string\"}],\"name\":\"log_named_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"}],\"name\":\"log_named_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"log_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"logs\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"IS_TEST\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeArtifacts\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"excludedArtifacts_\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"excludedContracts_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"excludedSenders_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"failed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setUp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetArtifactSelectors\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes4[]\",\"name\":\"selectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"name\":\"targetedArtifactSelectors_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetArtifacts\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"targetedArtifacts_\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targetedContracts_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetSelectors\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes4[]\",\"name\":\"selectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"name\":\"targetedSelectors_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targetedSenders_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test_owner_succeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test_storageSlots_succeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60008054600160ff1991821681178355600480549092168117909155601b805460806001600160a01b03199182168117909255601c80546101009083168117909155601d80546102009316831790559184905260a092909252606460c05260e09290925261c350909152602061014052600861016052674e4f4e5f5a45524f60c01b610180526060610120527fced1f90d33a6ca7cfbe479a1c2415c4287f559420415e3188c786e36414529be601e5560405260226101a08181529062001e9a6101c039601f90620000d290826200018c565b50348015620000e057600080fd5b5062000258565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200011257607f821691505b6020821081036200013357634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200018757600081815260208120601f850160051c81016020861015620001625750805b601f850160051c820191505b8181101562000183578281556001016200016e565b5050505b505050565b81516001600160401b03811115620001a857620001a8620000e7565b620001c081620001b98454620000fd565b8462000139565b602080601f831160018114620001f85760008415620001df5750858301515b600019600386901b1c1916600185901b17855562000183565b600085815260208120601f198616915b82811015620002295788860151825594840194600190910190840162000208565b5085821015620002485787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c05160e05161010051611c0c6200028e6000396000505060005050600050506000505060005050611c0c6000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c806385226c811161008c578063ba414fa611610066578063ba414fa61461015e578063e20c9f7114610176578063e210dfb31461017e578063fa7626d41461018657600080fd5b806385226c8114610139578063916a17c61461014e578063b5508aa91461015657600080fd5b80633f7286f4116100bd5780633f7286f41461011457806366080bd51461011c57806366d9a9a01461012457600080fd5b80630a9254e4146100e45780631ed7831c146100ee5780633e5e3c231461010c575b600080fd5b6100ec610193565b005b6100f6610203565b6040516101039190610f98565b60405180910390f35b6100f6610272565b6100f66102df565b6100ec61034c565b61012c6105b9565b6040516101039190610ff2565b6101416106ca565b6040516101039190611118565b61012c61079a565b6101416108a2565b610166610972565b6040519015158152602001610103565b6100f6610ad2565b6100ec610b3f565b6000546101669060ff1681565b60405161019f90610f8b565b604051809103906000f0801580156101bb573d6000803e3d6000fd5b50602180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6060600d80548060200260200160405190810160405280929190818152602001828054801561026857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161023d575b5050505050905090565b6060600f8054806020026020016040519081016040528092919081815260200182805480156102685760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161023d575050505050905090565b6060600e8054806020026020016040519081016040528092919081815260200182805480156102685760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161023d575050505050905090565b602154604080517f8da5cb5b0000000000000000000000000000000000000000000000000000000081529051737109709ecfa91a80626ff3989d68f67f5b1dd12d9263ca669fa79273ffffffffffffffffffffffffffffffffffffffff90911691638da5cb5b916004808201926020929091908290030181865afa1580156103d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103fc91906111ce565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401600060405180830381600087803b15801561046257600080fd5b505af1158015610476573d6000803e3d6000fd5b50506021546040517f13af40350000000000000000000000000000000000000000000000000000000081526001600482015273ffffffffffffffffffffffffffffffffffffffff90911692506313af40359150602401600060405180830381600087803b1580156104e657600080fd5b505af11580156104fa573d6000803e3d6000fd5b50506021546040517f667f9d7000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152600060248201526105b7925060019150737109709ecfa91a80626ff3989d68f67f5b1dd12d9063667f9d7090604401602060405180830381865afa15801561058e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105b2919061120b565b610bdb565b565b60606012805480602002602001604051908101604052809291908181526020016000905b828210156106c157600084815260209081902060408051808201825260028602909201805473ffffffffffffffffffffffffffffffffffffffff1683526001810180548351818702810187019094528084529394919385830193928301828280156106a957602002820191906000526020600020906000905b82829054906101000a900460e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190600401906020826003010492830192600103820291508084116106565790505b505050505081525050815260200190600101906105dd565b50505050905090565b60606011805480602002602001604051908101604052809291908181526020016000905b828210156106c157838290600052602060002001805461070d90611224565b80601f016020809104026020016040519081016040528092919081815260200182805461073990611224565b80156107865780601f1061075b57610100808354040283529160200191610786565b820191906000526020600020905b81548152906001019060200180831161076957829003601f168201915b5050505050815260200190600101906106ee565b60606013805480602002602001604051908101604052809291908181526020016000905b828210156106c157600084815260209081902060408051808201825260028602909201805473ffffffffffffffffffffffffffffffffffffffff16835260018101805483518187028101870190945280845293949193858301939283018282801561088a57602002820191906000526020600020906000905b82829054906101000a900460e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190600401906020826003010492830192600103820291508084116108375790505b505050505081525050815260200190600101906107be565b60606010805480602002602001604051908101604052809291908181526020016000905b828210156106c15783829060005260206000200180546108e590611224565b80601f016020809104026020016040519081016040528092919081815260200182805461091190611224565b801561095e5780601f106109335761010080835404028352916020019161095e565b820191906000526020600020905b81548152906001019060200180831161094157829003601f168201915b5050505050815260200190600101906108c6565b60008054610100900460ff16156109925750600054610100900460ff1690565b6000737109709ecfa91a80626ff3989d68f67f5b1dd12d3b15610acd5760408051737109709ecfa91a80626ff3989d68f67f5b1dd12d602082018190527f6661696c6564000000000000000000000000000000000000000000000000000082840152825180830384018152606083019093526000929091610a37917f667f9d70ca411d70ead50d8d5c22070dafc36ad75f3dcf5e7237b22ade9aecc491608001611277565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610a6f916112bf565b6000604051808303816000865af19150503d8060008114610aac576040519150601f19603f3d011682016040523d82523d6000602084013e610ab1565b606091505b5091505080806020019051810190610ac991906112db565b9150505b919050565b6060600c8054806020026020016040519081016040528092919081815260200182805480156102685760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161023d575050505050905090565b602154604080517f8da5cb5b00000000000000000000000000000000000000000000000000000000815290516105b79273ffffffffffffffffffffffffffffffffffffffff1691638da5cb5b9160048083019260209291908290030181865afa158015610bb0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bd491906111ce565b6000610ce9565b808214610ce5577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50604051610c679060208082526025908201527f4572726f723a2061203d3d2062206e6f7420736174697366696564205b62797460408201527f657333325d000000000000000000000000000000000000000000000000000000606082015260800190565b60405180910390a17fafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f9982604051610c9e91906112fd565b60405180910390a17fafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f9981604051610cd59190611348565b60405180910390a1610ce5610e0f565b5050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614610ce5577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50604051610da19060208082526025908201527f4572726f723a2061203d3d2062206e6f7420736174697366696564205b61646460408201527f726573735d000000000000000000000000000000000000000000000000000000606082015260800190565b60405180910390a17f9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f82604051610dd89190611385565b60405180910390a17f9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f81604051610cd591906113e6565b737109709ecfa91a80626ff3989d68f67f5b1dd12d3b15610f5d5760408051737109709ecfa91a80626ff3989d68f67f5b1dd12d602082018190527f6661696c656400000000000000000000000000000000000000000000000000009282019290925260016060820152600091907f70ca10bbd0dbfd9020a9f4b13402c16cb120705e0d1c0aeab10fa353ae586fc490608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610ede9291602001611277565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610f16916112bf565b6000604051808303816000865af19150503d8060008114610f53576040519150601f19603f3d011682016040523d82523d6000602084013e610f58565b606091505b505050505b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055565b6107dc8061142483390190565b6020808252825182820181905260009190848201906040850190845b81811015610fe657835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610fb4565b50909695505050505050565b60006020808301818452808551808352604092508286019150828160051b8701018488016000805b848110156110d9578984037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc00186528251805173ffffffffffffffffffffffffffffffffffffffff168552880151888501889052805188860181905290890190839060608701905b808310156110c45783517fffffffff00000000000000000000000000000000000000000000000000000000168252928b019260019290920191908b0190611082565b50978a0197955050509187019160010161101a565b50919998505050505050505050565b60005b838110156111035781810151838201526020016110eb565b83811115611112576000848401525b50505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156111c1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845281518051808752611184818989018a85016110e8565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169590950186019450928501929085019060010161113f565b5092979650505050505050565b6000602082840312156111e057600080fd5b815173ffffffffffffffffffffffffffffffffffffffff8116811461120457600080fd5b9392505050565b60006020828403121561121d57600080fd5b5051919050565b600181811c9082168061123857607f821691505b602082108103611271577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7fffffffff0000000000000000000000000000000000000000000000000000000083168152600082516112b18160048501602087016110e8565b919091016004019392505050565b600082516112d18184602087016110e8565b9190910192915050565b6000602082840312156112ed57600080fd5b8151801515811461120457600080fd5b60408152600061133a60408301600a81527f2020202020204c65667400000000000000000000000000000000000000000000602082015260400190565b905082602083015292915050565b60408152600061133a60408301600a81527f2020202020526967687400000000000000000000000000000000000000000000602082015260400190565b6040815260006113c260408301600a81527f2020202020204c65667400000000000000000000000000000000000000000000602082015260400190565b905073ffffffffffffffffffffffffffffffffffffffff8316602083015292915050565b6040815260006113c260408301600a81527f202020202052696768740000000000000000000000000000000000000000000060208201526040019056fe608060405234801561001057600080fd5b506107bc806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80638da5cb5b1161005b5780638da5cb5b146100fc5780639b19251a14610141578063b1540a0114610174578063bdc7b54f1461018757600080fd5b806308fd63221461008257806313af40351461009757806354fd4d50146100aa575b600080fd5b6100956100903660046106de565b61018f565b005b6100956100a536600461071a565b6102ef565b6100e66040518060400160405280600581526020017f312e312e3000000000000000000000000000000000000000000000000000000081525081565b6040516100f3919061073c565b60405180910390f35b60005461011c9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f3565b61016461014f36600461071a565b60016020526000908152604090205460ff1681565b60405190151581526020016100f3565b61016461018236600461071a565b610520565b610095610571565b60005473ffffffffffffffffffffffffffffffffffffffff163314610261576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604c60248201527f4465706c6f79657257686974656c6973743a2066756e6374696f6e2063616e2060448201527f6f6e6c792062652063616c6c656420627920746865206f776e6572206f66207460648201527f68697320636f6e74726163740000000000000000000000000000000000000000608482015260a4015b60405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660008181526001602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f8daaf060c3306c38e068a75c054bf96ecd85a3db1252712c4d93632744c42e0d910160405180910390a15050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146103bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604c60248201527f4465706c6f79657257686974656c6973743a2066756e6374696f6e2063616e2060448201527f6f6e6c792062652063616c6c656420627920746865206f776e6572206f66207460648201527f68697320636f6e74726163740000000000000000000000000000000000000000608482015260a401610258565b73ffffffffffffffffffffffffffffffffffffffff8116610485576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604d60248201527f4465706c6f79657257686974656c6973743a2063616e206f6e6c79206265206460448201527f697361626c65642076696120656e61626c65417262697472617279436f6e747260648201527f6163744465706c6f796d656e7400000000000000000000000000000000000000608482015260a401610258565b6000546040805173ffffffffffffffffffffffffffffffffffffffff928316815291831660208301527fb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c910160405180910390a1600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff16158061056b575073ffffffffffffffffffffffffffffffffffffffff821660009081526001602052604090205460ff165b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461063e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604c60248201527f4465706c6f79657257686974656c6973743a2066756e6374696f6e2063616e2060448201527f6f6e6c792062652063616c6c656420627920746865206f776e6572206f66207460648201527f68697320636f6e74726163740000000000000000000000000000000000000000608482015260a401610258565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681527fc0e106cf568e50698fdbde1eff56f5a5c966cc7958e37e276918e9e4ccdf8cd49060200160405180910390a1600080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055565b803573ffffffffffffffffffffffffffffffffffffffff811681146106d957600080fd5b919050565b600080604083850312156106f157600080fd5b6106fa836106b5565b91506020830135801515811461070f57600080fd5b809150509250929050565b60006020828403121561072c57600080fd5b610735826106b5565b9392505050565b600060208083528351808285015260005b818110156107695785810183015185820160400152820161074d565b8181111561077b576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201604001939250505056fea164736f6c634300080f000aa164736f6c634300080f000a0000111122223333444455556666777788889999aaaabbbbccccddddeeeeffff0000",
}

// DeployerWhitelistABI is the input ABI used to generate the binding from.
// Deprecated: Use DeployerWhitelistMetaData.ABI instead.
var DeployerWhitelistABI = DeployerWhitelistMetaData.ABI

// DeployerWhitelistBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DeployerWhitelistMetaData.Bin instead.
var DeployerWhitelistBin = DeployerWhitelistMetaData.Bin

// DeployDeployerWhitelist deploys a new Ethereum contract, binding an instance of DeployerWhitelist to it.
func DeployDeployerWhitelist(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DeployerWhitelist, error) {
	parsed, err := DeployerWhitelistMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DeployerWhitelistBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DeployerWhitelist{DeployerWhitelistCaller: DeployerWhitelistCaller{contract: contract}, DeployerWhitelistTransactor: DeployerWhitelistTransactor{contract: contract}, DeployerWhitelistFilterer: DeployerWhitelistFilterer{contract: contract}}, nil
}

// DeployerWhitelist is an auto generated Go binding around an Ethereum contract.
type DeployerWhitelist struct {
	DeployerWhitelistCaller     // Read-only binding to the contract
	DeployerWhitelistTransactor // Write-only binding to the contract
	DeployerWhitelistFilterer   // Log filterer for contract events
}

// DeployerWhitelistCaller is an auto generated read-only Go binding around an Ethereum contract.
type DeployerWhitelistCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerWhitelistTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DeployerWhitelistTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerWhitelistFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DeployerWhitelistFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerWhitelistSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DeployerWhitelistSession struct {
	Contract     *DeployerWhitelist // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DeployerWhitelistCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DeployerWhitelistCallerSession struct {
	Contract *DeployerWhitelistCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// DeployerWhitelistTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DeployerWhitelistTransactorSession struct {
	Contract     *DeployerWhitelistTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// DeployerWhitelistRaw is an auto generated low-level Go binding around an Ethereum contract.
type DeployerWhitelistRaw struct {
	Contract *DeployerWhitelist // Generic contract binding to access the raw methods on
}

// DeployerWhitelistCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DeployerWhitelistCallerRaw struct {
	Contract *DeployerWhitelistCaller // Generic read-only contract binding to access the raw methods on
}

// DeployerWhitelistTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DeployerWhitelistTransactorRaw struct {
	Contract *DeployerWhitelistTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeployerWhitelist creates a new instance of DeployerWhitelist, bound to a specific deployed contract.
func NewDeployerWhitelist(address common.Address, backend bind.ContractBackend) (*DeployerWhitelist, error) {
	contract, err := bindDeployerWhitelist(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelist{DeployerWhitelistCaller: DeployerWhitelistCaller{contract: contract}, DeployerWhitelistTransactor: DeployerWhitelistTransactor{contract: contract}, DeployerWhitelistFilterer: DeployerWhitelistFilterer{contract: contract}}, nil
}

// NewDeployerWhitelistCaller creates a new read-only instance of DeployerWhitelist, bound to a specific deployed contract.
func NewDeployerWhitelistCaller(address common.Address, caller bind.ContractCaller) (*DeployerWhitelistCaller, error) {
	contract, err := bindDeployerWhitelist(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistCaller{contract: contract}, nil
}

// NewDeployerWhitelistTransactor creates a new write-only instance of DeployerWhitelist, bound to a specific deployed contract.
func NewDeployerWhitelistTransactor(address common.Address, transactor bind.ContractTransactor) (*DeployerWhitelistTransactor, error) {
	contract, err := bindDeployerWhitelist(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistTransactor{contract: contract}, nil
}

// NewDeployerWhitelistFilterer creates a new log filterer instance of DeployerWhitelist, bound to a specific deployed contract.
func NewDeployerWhitelistFilterer(address common.Address, filterer bind.ContractFilterer) (*DeployerWhitelistFilterer, error) {
	contract, err := bindDeployerWhitelist(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistFilterer{contract: contract}, nil
}

// bindDeployerWhitelist binds a generic wrapper to an already deployed contract.
func bindDeployerWhitelist(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DeployerWhitelistMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployerWhitelist *DeployerWhitelistRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployerWhitelist.Contract.DeployerWhitelistCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployerWhitelist *DeployerWhitelistRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.DeployerWhitelistTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployerWhitelist *DeployerWhitelistRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.DeployerWhitelistTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployerWhitelist *DeployerWhitelistCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployerWhitelist.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployerWhitelist *DeployerWhitelistTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployerWhitelist *DeployerWhitelistTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.contract.Transact(opts, method, params...)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DeployerWhitelist *DeployerWhitelistCaller) ISTEST(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "IS_TEST")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DeployerWhitelist *DeployerWhitelistSession) ISTEST() (bool, error) {
	return _DeployerWhitelist.Contract.ISTEST(&_DeployerWhitelist.CallOpts)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) ISTEST() (bool, error) {
	return _DeployerWhitelist.Contract.ISTEST(&_DeployerWhitelist.CallOpts)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_DeployerWhitelist *DeployerWhitelistCaller) ExcludeArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "excludeArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_DeployerWhitelist *DeployerWhitelistSession) ExcludeArtifacts() ([]string, error) {
	return _DeployerWhitelist.Contract.ExcludeArtifacts(&_DeployerWhitelist.CallOpts)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) ExcludeArtifacts() ([]string, error) {
	return _DeployerWhitelist.Contract.ExcludeArtifacts(&_DeployerWhitelist.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_DeployerWhitelist *DeployerWhitelistCaller) ExcludeContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "excludeContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_DeployerWhitelist *DeployerWhitelistSession) ExcludeContracts() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.ExcludeContracts(&_DeployerWhitelist.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) ExcludeContracts() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.ExcludeContracts(&_DeployerWhitelist.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_DeployerWhitelist *DeployerWhitelistCaller) ExcludeSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "excludeSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_DeployerWhitelist *DeployerWhitelistSession) ExcludeSenders() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.ExcludeSenders(&_DeployerWhitelist.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) ExcludeSenders() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.ExcludeSenders(&_DeployerWhitelist.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_DeployerWhitelist *DeployerWhitelistCaller) TargetArtifactSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "targetArtifactSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_DeployerWhitelist *DeployerWhitelistSession) TargetArtifactSelectors() ([]StdInvariantFuzzSelector, error) {
	return _DeployerWhitelist.Contract.TargetArtifactSelectors(&_DeployerWhitelist.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) TargetArtifactSelectors() ([]StdInvariantFuzzSelector, error) {
	return _DeployerWhitelist.Contract.TargetArtifactSelectors(&_DeployerWhitelist.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_DeployerWhitelist *DeployerWhitelistCaller) TargetArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "targetArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_DeployerWhitelist *DeployerWhitelistSession) TargetArtifacts() ([]string, error) {
	return _DeployerWhitelist.Contract.TargetArtifacts(&_DeployerWhitelist.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) TargetArtifacts() ([]string, error) {
	return _DeployerWhitelist.Contract.TargetArtifacts(&_DeployerWhitelist.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_DeployerWhitelist *DeployerWhitelistCaller) TargetContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "targetContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_DeployerWhitelist *DeployerWhitelistSession) TargetContracts() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.TargetContracts(&_DeployerWhitelist.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) TargetContracts() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.TargetContracts(&_DeployerWhitelist.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_DeployerWhitelist *DeployerWhitelistCaller) TargetSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "targetSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_DeployerWhitelist *DeployerWhitelistSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _DeployerWhitelist.Contract.TargetSelectors(&_DeployerWhitelist.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _DeployerWhitelist.Contract.TargetSelectors(&_DeployerWhitelist.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_DeployerWhitelist *DeployerWhitelistCaller) TargetSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _DeployerWhitelist.contract.Call(opts, &out, "targetSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_DeployerWhitelist *DeployerWhitelistSession) TargetSenders() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.TargetSenders(&_DeployerWhitelist.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_DeployerWhitelist *DeployerWhitelistCallerSession) TargetSenders() ([]common.Address, error) {
	return _DeployerWhitelist.Contract.TargetSenders(&_DeployerWhitelist.CallOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DeployerWhitelist *DeployerWhitelistTransactor) Failed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerWhitelist.contract.Transact(opts, "failed")
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DeployerWhitelist *DeployerWhitelistSession) Failed() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.Failed(&_DeployerWhitelist.TransactOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DeployerWhitelist *DeployerWhitelistTransactorSession) Failed() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.Failed(&_DeployerWhitelist.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_DeployerWhitelist *DeployerWhitelistTransactor) SetUp(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerWhitelist.contract.Transact(opts, "setUp")
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_DeployerWhitelist *DeployerWhitelistSession) SetUp() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.SetUp(&_DeployerWhitelist.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_DeployerWhitelist *DeployerWhitelistTransactorSession) SetUp() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.SetUp(&_DeployerWhitelist.TransactOpts)
}

// TestOwnerSucceeds is a paid mutator transaction binding the contract method 0xe210dfb3.
//
// Solidity: function test_owner_succeeds() returns()
func (_DeployerWhitelist *DeployerWhitelistTransactor) TestOwnerSucceeds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerWhitelist.contract.Transact(opts, "test_owner_succeeds")
}

// TestOwnerSucceeds is a paid mutator transaction binding the contract method 0xe210dfb3.
//
// Solidity: function test_owner_succeeds() returns()
func (_DeployerWhitelist *DeployerWhitelistSession) TestOwnerSucceeds() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.TestOwnerSucceeds(&_DeployerWhitelist.TransactOpts)
}

// TestOwnerSucceeds is a paid mutator transaction binding the contract method 0xe210dfb3.
//
// Solidity: function test_owner_succeeds() returns()
func (_DeployerWhitelist *DeployerWhitelistTransactorSession) TestOwnerSucceeds() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.TestOwnerSucceeds(&_DeployerWhitelist.TransactOpts)
}

// TestStorageSlotsSucceeds is a paid mutator transaction binding the contract method 0x66080bd5.
//
// Solidity: function test_storageSlots_succeeds() returns()
func (_DeployerWhitelist *DeployerWhitelistTransactor) TestStorageSlotsSucceeds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerWhitelist.contract.Transact(opts, "test_storageSlots_succeeds")
}

// TestStorageSlotsSucceeds is a paid mutator transaction binding the contract method 0x66080bd5.
//
// Solidity: function test_storageSlots_succeeds() returns()
func (_DeployerWhitelist *DeployerWhitelistSession) TestStorageSlotsSucceeds() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.TestStorageSlotsSucceeds(&_DeployerWhitelist.TransactOpts)
}

// TestStorageSlotsSucceeds is a paid mutator transaction binding the contract method 0x66080bd5.
//
// Solidity: function test_storageSlots_succeeds() returns()
func (_DeployerWhitelist *DeployerWhitelistTransactorSession) TestStorageSlotsSucceeds() (*types.Transaction, error) {
	return _DeployerWhitelist.Contract.TestStorageSlotsSucceeds(&_DeployerWhitelist.TransactOpts)
}

// DeployerWhitelistOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DeployerWhitelist contract.
type DeployerWhitelistOwnershipTransferredIterator struct {
	Event *DeployerWhitelistOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistOwnershipTransferred)
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
		it.Event = new(DeployerWhitelistOwnershipTransferred)
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
func (it *DeployerWhitelistOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistOwnershipTransferred represents a OwnershipTransferred event raised by the DeployerWhitelist contract.
type DeployerWhitelistOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DeployerWhitelistOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistOwnershipTransferredIterator{contract: _DeployerWhitelist.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistOwnershipTransferred)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseOwnershipTransferred(log types.Log) (*DeployerWhitelistOwnershipTransferred, error) {
	event := new(DeployerWhitelistOwnershipTransferred)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistTransactionDepositedIterator is returned from FilterTransactionDeposited and is used to iterate over the raw logs and unpacked data for TransactionDeposited events raised by the DeployerWhitelist contract.
type DeployerWhitelistTransactionDepositedIterator struct {
	Event *DeployerWhitelistTransactionDeposited // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistTransactionDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistTransactionDeposited)
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
		it.Event = new(DeployerWhitelistTransactionDeposited)
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
func (it *DeployerWhitelistTransactionDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistTransactionDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistTransactionDeposited represents a TransactionDeposited event raised by the DeployerWhitelist contract.
type DeployerWhitelistTransactionDeposited struct {
	From       common.Address
	To         common.Address
	Version    *big.Int
	OpaqueData []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransactionDeposited is a free log retrieval operation binding the contract event 0xb3813568d9991fc951961fcb4c784893574240a28925604d09fc577c55bb7c32.
//
// Solidity: event TransactionDeposited(address indexed from, address indexed to, uint256 indexed version, bytes opaqueData)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterTransactionDeposited(opts *bind.FilterOpts, from []common.Address, to []common.Address, version []*big.Int) (*DeployerWhitelistTransactionDepositedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "TransactionDeposited", fromRule, toRule, versionRule)
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistTransactionDepositedIterator{contract: _DeployerWhitelist.contract, event: "TransactionDeposited", logs: logs, sub: sub}, nil
}

// WatchTransactionDeposited is a free log subscription operation binding the contract event 0xb3813568d9991fc951961fcb4c784893574240a28925604d09fc577c55bb7c32.
//
// Solidity: event TransactionDeposited(address indexed from, address indexed to, uint256 indexed version, bytes opaqueData)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchTransactionDeposited(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistTransactionDeposited, from []common.Address, to []common.Address, version []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "TransactionDeposited", fromRule, toRule, versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistTransactionDeposited)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "TransactionDeposited", log); err != nil {
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

// ParseTransactionDeposited is a log parse operation binding the contract event 0xb3813568d9991fc951961fcb4c784893574240a28925604d09fc577c55bb7c32.
//
// Solidity: event TransactionDeposited(address indexed from, address indexed to, uint256 indexed version, bytes opaqueData)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseTransactionDeposited(log types.Log) (*DeployerWhitelistTransactionDeposited, error) {
	event := new(DeployerWhitelistTransactionDeposited)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "TransactionDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogIterator is returned from FilterLog and is used to iterate over the raw logs and unpacked data for Log events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogIterator struct {
	Event *DeployerWhitelistLog // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLog)
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
		it.Event = new(DeployerWhitelistLog)
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
func (it *DeployerWhitelistLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLog represents a Log event raised by the DeployerWhitelist contract.
type DeployerWhitelistLog struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLog is a free log retrieval operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLog(opts *bind.FilterOpts) (*DeployerWhitelistLogIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogIterator{contract: _DeployerWhitelist.contract, event: "log", logs: logs, sub: sub}, nil
}

// WatchLog is a free log subscription operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLog(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLog) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLog)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log", log); err != nil {
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

// ParseLog is a log parse operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLog(log types.Log) (*DeployerWhitelistLog, error) {
	event := new(DeployerWhitelistLog)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogAddressIterator is returned from FilterLogAddress and is used to iterate over the raw logs and unpacked data for LogAddress events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogAddressIterator struct {
	Event *DeployerWhitelistLogAddress // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogAddress)
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
		it.Event = new(DeployerWhitelistLogAddress)
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
func (it *DeployerWhitelistLogAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogAddress represents a LogAddress event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogAddress struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogAddress is a free log retrieval operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogAddress(opts *bind.FilterOpts) (*DeployerWhitelistLogAddressIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogAddressIterator{contract: _DeployerWhitelist.contract, event: "log_address", logs: logs, sub: sub}, nil
}

// WatchLogAddress is a free log subscription operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogAddress(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogAddress) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogAddress)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_address", log); err != nil {
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

// ParseLogAddress is a log parse operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogAddress(log types.Log) (*DeployerWhitelistLogAddress, error) {
	event := new(DeployerWhitelistLogAddress)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogArrayIterator is returned from FilterLogArray and is used to iterate over the raw logs and unpacked data for LogArray events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogArrayIterator struct {
	Event *DeployerWhitelistLogArray // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogArray)
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
		it.Event = new(DeployerWhitelistLogArray)
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
func (it *DeployerWhitelistLogArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogArray represents a LogArray event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogArray struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray is a free log retrieval operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogArray(opts *bind.FilterOpts) (*DeployerWhitelistLogArrayIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogArrayIterator{contract: _DeployerWhitelist.contract, event: "log_array", logs: logs, sub: sub}, nil
}

// WatchLogArray is a free log subscription operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogArray(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogArray) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogArray)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_array", log); err != nil {
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

// ParseLogArray is a log parse operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogArray(log types.Log) (*DeployerWhitelistLogArray, error) {
	event := new(DeployerWhitelistLogArray)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogArray0Iterator is returned from FilterLogArray0 and is used to iterate over the raw logs and unpacked data for LogArray0 events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogArray0Iterator struct {
	Event *DeployerWhitelistLogArray0 // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogArray0)
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
		it.Event = new(DeployerWhitelistLogArray0)
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
func (it *DeployerWhitelistLogArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogArray0 represents a LogArray0 event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogArray0 struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray0 is a free log retrieval operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogArray0(opts *bind.FilterOpts) (*DeployerWhitelistLogArray0Iterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogArray0Iterator{contract: _DeployerWhitelist.contract, event: "log_array0", logs: logs, sub: sub}, nil
}

// WatchLogArray0 is a free log subscription operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogArray0(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogArray0) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogArray0)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_array0", log); err != nil {
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

// ParseLogArray0 is a log parse operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogArray0(log types.Log) (*DeployerWhitelistLogArray0, error) {
	event := new(DeployerWhitelistLogArray0)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogArray1Iterator is returned from FilterLogArray1 and is used to iterate over the raw logs and unpacked data for LogArray1 events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogArray1Iterator struct {
	Event *DeployerWhitelistLogArray1 // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogArray1)
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
		it.Event = new(DeployerWhitelistLogArray1)
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
func (it *DeployerWhitelistLogArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogArray1 represents a LogArray1 event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogArray1 struct {
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray1 is a free log retrieval operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogArray1(opts *bind.FilterOpts) (*DeployerWhitelistLogArray1Iterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogArray1Iterator{contract: _DeployerWhitelist.contract, event: "log_array1", logs: logs, sub: sub}, nil
}

// WatchLogArray1 is a free log subscription operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogArray1(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogArray1) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogArray1)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_array1", log); err != nil {
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

// ParseLogArray1 is a log parse operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogArray1(log types.Log) (*DeployerWhitelistLogArray1, error) {
	event := new(DeployerWhitelistLogArray1)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogBytesIterator is returned from FilterLogBytes and is used to iterate over the raw logs and unpacked data for LogBytes events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogBytesIterator struct {
	Event *DeployerWhitelistLogBytes // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogBytes)
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
		it.Event = new(DeployerWhitelistLogBytes)
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
func (it *DeployerWhitelistLogBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogBytes represents a LogBytes event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogBytes struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes is a free log retrieval operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogBytes(opts *bind.FilterOpts) (*DeployerWhitelistLogBytesIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogBytesIterator{contract: _DeployerWhitelist.contract, event: "log_bytes", logs: logs, sub: sub}, nil
}

// WatchLogBytes is a free log subscription operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogBytes(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogBytes) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogBytes)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_bytes", log); err != nil {
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

// ParseLogBytes is a log parse operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogBytes(log types.Log) (*DeployerWhitelistLogBytes, error) {
	event := new(DeployerWhitelistLogBytes)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogBytes32Iterator is returned from FilterLogBytes32 and is used to iterate over the raw logs and unpacked data for LogBytes32 events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogBytes32Iterator struct {
	Event *DeployerWhitelistLogBytes32 // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogBytes32)
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
		it.Event = new(DeployerWhitelistLogBytes32)
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
func (it *DeployerWhitelistLogBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogBytes32 represents a LogBytes32 event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogBytes32 struct {
	Arg0 [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes32 is a free log retrieval operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogBytes32(opts *bind.FilterOpts) (*DeployerWhitelistLogBytes32Iterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogBytes32Iterator{contract: _DeployerWhitelist.contract, event: "log_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogBytes32 is a free log subscription operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogBytes32(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogBytes32) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogBytes32)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_bytes32", log); err != nil {
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

// ParseLogBytes32 is a log parse operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogBytes32(log types.Log) (*DeployerWhitelistLogBytes32, error) {
	event := new(DeployerWhitelistLogBytes32)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogIntIterator is returned from FilterLogInt and is used to iterate over the raw logs and unpacked data for LogInt events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogIntIterator struct {
	Event *DeployerWhitelistLogInt // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogInt)
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
		it.Event = new(DeployerWhitelistLogInt)
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
func (it *DeployerWhitelistLogIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogInt represents a LogInt event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogInt struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogInt is a free log retrieval operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogInt(opts *bind.FilterOpts) (*DeployerWhitelistLogIntIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogIntIterator{contract: _DeployerWhitelist.contract, event: "log_int", logs: logs, sub: sub}, nil
}

// WatchLogInt is a free log subscription operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogInt(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogInt) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogInt)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_int", log); err != nil {
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

// ParseLogInt is a log parse operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogInt(log types.Log) (*DeployerWhitelistLogInt, error) {
	event := new(DeployerWhitelistLogInt)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedAddressIterator is returned from FilterLogNamedAddress and is used to iterate over the raw logs and unpacked data for LogNamedAddress events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedAddressIterator struct {
	Event *DeployerWhitelistLogNamedAddress // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedAddress)
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
		it.Event = new(DeployerWhitelistLogNamedAddress)
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
func (it *DeployerWhitelistLogNamedAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedAddress represents a LogNamedAddress event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedAddress struct {
	Key string
	Val common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedAddress is a free log retrieval operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedAddress(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedAddressIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedAddressIterator{contract: _DeployerWhitelist.contract, event: "log_named_address", logs: logs, sub: sub}, nil
}

// WatchLogNamedAddress is a free log subscription operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedAddress(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedAddress) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedAddress)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_address", log); err != nil {
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

// ParseLogNamedAddress is a log parse operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedAddress(log types.Log) (*DeployerWhitelistLogNamedAddress, error) {
	event := new(DeployerWhitelistLogNamedAddress)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedArrayIterator is returned from FilterLogNamedArray and is used to iterate over the raw logs and unpacked data for LogNamedArray events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedArrayIterator struct {
	Event *DeployerWhitelistLogNamedArray // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedArray)
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
		it.Event = new(DeployerWhitelistLogNamedArray)
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
func (it *DeployerWhitelistLogNamedArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedArray represents a LogNamedArray event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedArray struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray is a free log retrieval operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedArray(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedArrayIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedArrayIterator{contract: _DeployerWhitelist.contract, event: "log_named_array", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray is a free log subscription operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedArray(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedArray) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedArray)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_array", log); err != nil {
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

// ParseLogNamedArray is a log parse operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedArray(log types.Log) (*DeployerWhitelistLogNamedArray, error) {
	event := new(DeployerWhitelistLogNamedArray)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedArray0Iterator is returned from FilterLogNamedArray0 and is used to iterate over the raw logs and unpacked data for LogNamedArray0 events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedArray0Iterator struct {
	Event *DeployerWhitelistLogNamedArray0 // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedArray0)
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
		it.Event = new(DeployerWhitelistLogNamedArray0)
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
func (it *DeployerWhitelistLogNamedArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedArray0 represents a LogNamedArray0 event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedArray0 struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray0 is a free log retrieval operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedArray0(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedArray0Iterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedArray0Iterator{contract: _DeployerWhitelist.contract, event: "log_named_array0", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray0 is a free log subscription operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedArray0(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedArray0) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedArray0)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_array0", log); err != nil {
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

// ParseLogNamedArray0 is a log parse operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedArray0(log types.Log) (*DeployerWhitelistLogNamedArray0, error) {
	event := new(DeployerWhitelistLogNamedArray0)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedArray1Iterator is returned from FilterLogNamedArray1 and is used to iterate over the raw logs and unpacked data for LogNamedArray1 events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedArray1Iterator struct {
	Event *DeployerWhitelistLogNamedArray1 // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedArray1)
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
		it.Event = new(DeployerWhitelistLogNamedArray1)
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
func (it *DeployerWhitelistLogNamedArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedArray1 represents a LogNamedArray1 event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedArray1 struct {
	Key string
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray1 is a free log retrieval operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedArray1(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedArray1Iterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedArray1Iterator{contract: _DeployerWhitelist.contract, event: "log_named_array1", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray1 is a free log subscription operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedArray1(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedArray1) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedArray1)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_array1", log); err != nil {
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

// ParseLogNamedArray1 is a log parse operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedArray1(log types.Log) (*DeployerWhitelistLogNamedArray1, error) {
	event := new(DeployerWhitelistLogNamedArray1)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedBytesIterator is returned from FilterLogNamedBytes and is used to iterate over the raw logs and unpacked data for LogNamedBytes events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedBytesIterator struct {
	Event *DeployerWhitelistLogNamedBytes // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedBytes)
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
		it.Event = new(DeployerWhitelistLogNamedBytes)
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
func (it *DeployerWhitelistLogNamedBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedBytes represents a LogNamedBytes event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedBytes struct {
	Key string
	Val []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes is a free log retrieval operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedBytes(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedBytesIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedBytesIterator{contract: _DeployerWhitelist.contract, event: "log_named_bytes", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes is a free log subscription operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedBytes(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedBytes) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedBytes)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
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

// ParseLogNamedBytes is a log parse operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedBytes(log types.Log) (*DeployerWhitelistLogNamedBytes, error) {
	event := new(DeployerWhitelistLogNamedBytes)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedBytes32Iterator is returned from FilterLogNamedBytes32 and is used to iterate over the raw logs and unpacked data for LogNamedBytes32 events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedBytes32Iterator struct {
	Event *DeployerWhitelistLogNamedBytes32 // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedBytes32)
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
		it.Event = new(DeployerWhitelistLogNamedBytes32)
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
func (it *DeployerWhitelistLogNamedBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedBytes32 represents a LogNamedBytes32 event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedBytes32 struct {
	Key string
	Val [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes32 is a free log retrieval operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedBytes32(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedBytes32Iterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedBytes32Iterator{contract: _DeployerWhitelist.contract, event: "log_named_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes32 is a free log subscription operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedBytes32(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedBytes32) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedBytes32)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
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

// ParseLogNamedBytes32 is a log parse operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedBytes32(log types.Log) (*DeployerWhitelistLogNamedBytes32, error) {
	event := new(DeployerWhitelistLogNamedBytes32)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedDecimalIntIterator is returned from FilterLogNamedDecimalInt and is used to iterate over the raw logs and unpacked data for LogNamedDecimalInt events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedDecimalIntIterator struct {
	Event *DeployerWhitelistLogNamedDecimalInt // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedDecimalIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedDecimalInt)
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
		it.Event = new(DeployerWhitelistLogNamedDecimalInt)
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
func (it *DeployerWhitelistLogNamedDecimalIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedDecimalIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedDecimalInt represents a LogNamedDecimalInt event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedDecimalInt struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalInt is a free log retrieval operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedDecimalInt(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedDecimalIntIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedDecimalIntIterator{contract: _DeployerWhitelist.contract, event: "log_named_decimal_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalInt is a free log subscription operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedDecimalInt(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedDecimalInt) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedDecimalInt)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
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

// ParseLogNamedDecimalInt is a log parse operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedDecimalInt(log types.Log) (*DeployerWhitelistLogNamedDecimalInt, error) {
	event := new(DeployerWhitelistLogNamedDecimalInt)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedDecimalUintIterator is returned from FilterLogNamedDecimalUint and is used to iterate over the raw logs and unpacked data for LogNamedDecimalUint events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedDecimalUintIterator struct {
	Event *DeployerWhitelistLogNamedDecimalUint // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedDecimalUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedDecimalUint)
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
		it.Event = new(DeployerWhitelistLogNamedDecimalUint)
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
func (it *DeployerWhitelistLogNamedDecimalUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedDecimalUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedDecimalUint represents a LogNamedDecimalUint event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedDecimalUint struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalUint is a free log retrieval operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedDecimalUint(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedDecimalUintIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedDecimalUintIterator{contract: _DeployerWhitelist.contract, event: "log_named_decimal_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalUint is a free log subscription operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedDecimalUint(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedDecimalUint) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedDecimalUint)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
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

// ParseLogNamedDecimalUint is a log parse operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedDecimalUint(log types.Log) (*DeployerWhitelistLogNamedDecimalUint, error) {
	event := new(DeployerWhitelistLogNamedDecimalUint)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedIntIterator is returned from FilterLogNamedInt and is used to iterate over the raw logs and unpacked data for LogNamedInt events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedIntIterator struct {
	Event *DeployerWhitelistLogNamedInt // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedInt)
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
		it.Event = new(DeployerWhitelistLogNamedInt)
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
func (it *DeployerWhitelistLogNamedIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedInt represents a LogNamedInt event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedInt struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedInt is a free log retrieval operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedInt(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedIntIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedIntIterator{contract: _DeployerWhitelist.contract, event: "log_named_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedInt is a free log subscription operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedInt(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedInt) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedInt)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_int", log); err != nil {
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

// ParseLogNamedInt is a log parse operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedInt(log types.Log) (*DeployerWhitelistLogNamedInt, error) {
	event := new(DeployerWhitelistLogNamedInt)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedStringIterator is returned from FilterLogNamedString and is used to iterate over the raw logs and unpacked data for LogNamedString events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedStringIterator struct {
	Event *DeployerWhitelistLogNamedString // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedString)
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
		it.Event = new(DeployerWhitelistLogNamedString)
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
func (it *DeployerWhitelistLogNamedStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedString represents a LogNamedString event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedString struct {
	Key string
	Val string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedString is a free log retrieval operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedString(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedStringIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedStringIterator{contract: _DeployerWhitelist.contract, event: "log_named_string", logs: logs, sub: sub}, nil
}

// WatchLogNamedString is a free log subscription operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedString(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedString) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedString)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_string", log); err != nil {
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

// ParseLogNamedString is a log parse operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedString(log types.Log) (*DeployerWhitelistLogNamedString, error) {
	event := new(DeployerWhitelistLogNamedString)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogNamedUintIterator is returned from FilterLogNamedUint and is used to iterate over the raw logs and unpacked data for LogNamedUint events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedUintIterator struct {
	Event *DeployerWhitelistLogNamedUint // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogNamedUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogNamedUint)
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
		it.Event = new(DeployerWhitelistLogNamedUint)
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
func (it *DeployerWhitelistLogNamedUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogNamedUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogNamedUint represents a LogNamedUint event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogNamedUint struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedUint is a free log retrieval operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogNamedUint(opts *bind.FilterOpts) (*DeployerWhitelistLogNamedUintIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogNamedUintIterator{contract: _DeployerWhitelist.contract, event: "log_named_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedUint is a free log subscription operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogNamedUint(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogNamedUint) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogNamedUint)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_uint", log); err != nil {
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

// ParseLogNamedUint is a log parse operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogNamedUint(log types.Log) (*DeployerWhitelistLogNamedUint, error) {
	event := new(DeployerWhitelistLogNamedUint)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_named_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogStringIterator is returned from FilterLogString and is used to iterate over the raw logs and unpacked data for LogString events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogStringIterator struct {
	Event *DeployerWhitelistLogString // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogString)
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
		it.Event = new(DeployerWhitelistLogString)
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
func (it *DeployerWhitelistLogStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogString represents a LogString event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogString struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogString is a free log retrieval operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogString(opts *bind.FilterOpts) (*DeployerWhitelistLogStringIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogStringIterator{contract: _DeployerWhitelist.contract, event: "log_string", logs: logs, sub: sub}, nil
}

// WatchLogString is a free log subscription operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogString(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogString) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogString)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_string", log); err != nil {
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

// ParseLogString is a log parse operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogString(log types.Log) (*DeployerWhitelistLogString, error) {
	event := new(DeployerWhitelistLogString)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogUintIterator is returned from FilterLogUint and is used to iterate over the raw logs and unpacked data for LogUint events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogUintIterator struct {
	Event *DeployerWhitelistLogUint // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogUint)
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
		it.Event = new(DeployerWhitelistLogUint)
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
func (it *DeployerWhitelistLogUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogUint represents a LogUint event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogUint struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogUint is a free log retrieval operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogUint(opts *bind.FilterOpts) (*DeployerWhitelistLogUintIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogUintIterator{contract: _DeployerWhitelist.contract, event: "log_uint", logs: logs, sub: sub}, nil
}

// WatchLogUint is a free log subscription operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogUint(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogUint) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogUint)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "log_uint", log); err != nil {
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

// ParseLogUint is a log parse operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogUint(log types.Log) (*DeployerWhitelistLogUint, error) {
	event := new(DeployerWhitelistLogUint)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "log_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerWhitelistLogsIterator is returned from FilterLogs and is used to iterate over the raw logs and unpacked data for Logs events raised by the DeployerWhitelist contract.
type DeployerWhitelistLogsIterator struct {
	Event *DeployerWhitelistLogs // Event containing the contract specifics and raw log

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
func (it *DeployerWhitelistLogsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerWhitelistLogs)
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
		it.Event = new(DeployerWhitelistLogs)
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
func (it *DeployerWhitelistLogsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerWhitelistLogsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerWhitelistLogs represents a Logs event raised by the DeployerWhitelist contract.
type DeployerWhitelistLogs struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogs is a free log retrieval operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) FilterLogs(opts *bind.FilterOpts) (*DeployerWhitelistLogsIterator, error) {

	logs, sub, err := _DeployerWhitelist.contract.FilterLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return &DeployerWhitelistLogsIterator{contract: _DeployerWhitelist.contract, event: "logs", logs: logs, sub: sub}, nil
}

// WatchLogs is a free log subscription operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) WatchLogs(opts *bind.WatchOpts, sink chan<- *DeployerWhitelistLogs) (event.Subscription, error) {

	logs, sub, err := _DeployerWhitelist.contract.WatchLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerWhitelistLogs)
				if err := _DeployerWhitelist.contract.UnpackLog(event, "logs", log); err != nil {
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

// ParseLogs is a log parse operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DeployerWhitelist *DeployerWhitelistFilterer) ParseLogs(log types.Log) (*DeployerWhitelistLogs, error) {
	event := new(DeployerWhitelistLogs)
	if err := _DeployerWhitelist.contract.UnpackLog(event, "logs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
