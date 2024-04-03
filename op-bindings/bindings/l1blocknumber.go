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

// L1BlockNumberMetaData contains all meta data concerning the L1BlockNumber contract.
var L1BlockNumberMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"log_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"val\",\"type\":\"uint256[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"val\",\"type\":\"int256[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"val\",\"type\":\"address[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"log_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"log_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"log_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"val\",\"type\":\"address\"}],\"name\":\"log_named_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"val\",\"type\":\"uint256[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"val\",\"type\":\"int256[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"val\",\"type\":\"address[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"val\",\"type\":\"bytes\"}],\"name\":\"log_named_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"val\",\"type\":\"bytes32\"}],\"name\":\"log_named_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"}],\"name\":\"log_named_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"val\",\"type\":\"string\"}],\"name\":\"log_named_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"}],\"name\":\"log_named_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"log_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"logs\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"IS_TEST\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeArtifacts\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"excludedArtifacts_\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"excludedContracts_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"excludedSenders_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"failed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setUp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetArtifactSelectors\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes4[]\",\"name\":\"selectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"name\":\"targetedArtifactSelectors_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetArtifacts\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"targetedArtifacts_\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targetedContracts_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetSelectors\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes4[]\",\"name\":\"selectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"name\":\"targetedSelectors_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targetedSenders_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test_fallback_succeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test_getL1BlockNumber_succeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test_receive_succeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405260008054600160ff19918216811790925560048054909116909117905534801561002d57600080fd5b50611fc48061003d6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806385226c811161008c578063ba414fa611610066578063ba414fa614610171578063e20c9f7114610189578063fa7626d414610191578063fe43c9dd1461019e57600080fd5b806385226c811461014c578063916a17c614610161578063b5508aa91461016957600080fd5b80633f7286f4116100c85780633f7286f41461011f578063573980d01461012757806366b5a0b81461012f57806366d9a9a01461013757600080fd5b80630a9254e4146100ef5780631ed7831c146100f95780633e5e3c2314610117575b600080fd5b6100f76101a6565b005b6101016104fd565b60405161010e91906113b5565b60405180910390f35b61010161056c565b6101016105d9565b6100f7610646565b6100f76106e4565b61013f61077e565b60405161010e919061140f565b61015461088f565b60405161010e919061157b565b61013f61095f565b610154610a67565b610179610b37565b604051901515815260200161010e565b610101610c97565b6000546101799060ff1681565b6100f7610d04565b604051737109709ecfa91a80626ff3989d68f67f5b1dd12d9063b4d6c78290734200000000000000000000000000000000000015906101e49061139b565b604051809103906000f080158015610200573d6000803e3d6000fd5b5073ffffffffffffffffffffffffffffffffffffffff16803b806020016040519081016040528181526000908060200190933c6040518363ffffffff1660e01b81526004016102509291906115fb565b600060405180830381600087803b15801561026a57600080fd5b505af115801561027e573d6000803e3d6000fd5b5050601b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673420000000000000000000000000000000000001517905550506040516102cc906113a8565b604051809103906000f0801580156102e8573d6000803e3d6000fd5b50601c80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601b54604080517fe591b2820000000000000000000000000000000000000000000000000000000081529051737109709ecfa91a80626ff3989d68f67f5b1dd12d9363ca669fa793169163e591b2829160048083019260209291908290030181865afa15801561039e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103c29190611632565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401600060405180830381600087803b15801561042857600080fd5b505af115801561043c573d6000803e3d6000fd5b5050601b546040517f015d8eb90000000000000000000000000000000000000000000000000000000081526063600480830191909152600260248301819052600360448401819052600a60648501526084840192909252600060a484015260c483015260e482015273ffffffffffffffffffffffffffffffffffffffff909116925063015d8eb9915061010401600060405180830381600087803b1580156104e357600080fd5b505af11580156104f7573d6000803e3d6000fd5b50505050565b6060600d80548060200260200160405190810160405280929190818152602001828054801561056257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610537575b5050505050905090565b6060600f8054806020026020016040519081016040528092919081815260200182805480156105625760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610537575050505050905090565b6060600e8054806020026020016040519081016040528092919081815260200182805480156105625760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610537575050505050905090565b601c54604080517fb9b3efe900000000000000000000000000000000000000000000000000000000815290516106e29273ffffffffffffffffffffffffffffffffffffffff169163b9b3efe99160048083019260209291908290030181865afa1580156106b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106db919061166f565b6063610d63565b565b601c54604051600091829173ffffffffffffffffffffffffffffffffffffffff909116908281818181865af19150503d806000811461073f576040519150601f19603f3d011682016040523d82523d6000602084013e610744565b606091505b5091509150610754826001610e6d565b604080516063602082015261077a91839101604051602081830303815290604052611057565b5050565b60606012805480602002602001604051908101604052809291908181526020016000905b8282101561088657600084815260209081902060408051808201825260028602909201805473ffffffffffffffffffffffffffffffffffffffff16835260018101805483518187028101870190945280845293949193858301939283018282801561086e57602002820191906000526020600020906000905b82829054906101000a900460e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001906004019060208260030104928301926001038202915080841161081b5790505b505050505081525050815260200190600101906107a2565b50505050905090565b60606011805480602002602001604051908101604052809291908181526020016000905b828210156108865783829060005260206000200180546108d290611688565b80601f01602080910402602001604051908101604052809291908181526020018280546108fe90611688565b801561094b5780601f106109205761010080835404028352916020019161094b565b820191906000526020600020905b81548152906001019060200180831161092e57829003601f168201915b5050505050815260200190600101906108b3565b60606013805480602002602001604051908101604052809291908181526020016000905b8282101561088657600084815260209081902060408051808201825260028602909201805473ffffffffffffffffffffffffffffffffffffffff168352600181018054835181870281018701909452808452939491938583019392830182828015610a4f57602002820191906000526020600020906000905b82829054906101000a900460e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190600401906020826003010492830192600103820291508084116109fc5790505b50505050508152505081526020019060010190610983565b60606010805480602002602001604051908101604052809291908181526020016000905b82821015610886578382906000526020600020018054610aaa90611688565b80601f0160208091040260200160405190810160405280929190818152602001828054610ad690611688565b8015610b235780601f10610af857610100808354040283529160200191610b23565b820191906000526020600020905b815481529060010190602001808311610b0657829003601f168201915b505050505081526020019060010190610a8b565b60008054610100900460ff1615610b575750600054610100900460ff1690565b6000737109709ecfa91a80626ff3989d68f67f5b1dd12d3b15610c925760408051737109709ecfa91a80626ff3989d68f67f5b1dd12d602082018190527f6661696c6564000000000000000000000000000000000000000000000000000082840152825180830384018152606083019093526000929091610bfc917f667f9d70ca411d70ead50d8d5c22070dafc36ad75f3dcf5e7237b22ade9aecc4916080016116db565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610c3491611723565b6000604051808303816000865af19150503d8060008114610c71576040519150601f19603f3d011682016040523d82523d6000602084013e610c76565b606091505b5091505080806020019051810190610c8e919061173f565b9150505b919050565b6060600c8054806020026020016040519081016040528092919081815260200182805480156105625760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610537575050505050905090565b601c54604051600091829173ffffffffffffffffffffffffffffffffffffffff909116906001908381818185875af1925050503d806000811461073f576040519150601f19603f3d011682016040523d82523d6000602084013e610744565b80821461077a577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50604051610def9060208082526022908201527f4572726f723a2061203d3d2062206e6f7420736174697366696564205b75696e60408201527f745d000000000000000000000000000000000000000000000000000000000000606082015260800190565b60405180910390a17fb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a882604051610e269190611761565b60405180910390a17fb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a881604051610e5d91906117ac565b60405180910390a161077a611061565b8015158215151461077a577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50604051610efd9060208082526022908201527f4572726f723a2061203d3d2062206e6f7420736174697366696564205b626f6f60408201527f6c5d000000000000000000000000000000000000000000000000000000000000606082015260800190565b60405180910390a17f280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf358382610f66576040518060400160405280600581526020017f66616c7365000000000000000000000000000000000000000000000000000000815250610f9d565b6040518060400160405280600481526020017f74727565000000000000000000000000000000000000000000000000000000008152505b604051610faa91906117e9565b60405180910390a17f280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf358381611013576040518060400160405280600581526020017f66616c736500000000000000000000000000000000000000000000000000000081525061104a565b6040518060400160405280600481526020017f74727565000000000000000000000000000000000000000000000000000000008152505b604051610e5d9190611838565b61077a82826111dd565b737109709ecfa91a80626ff3989d68f67f5b1dd12d3b156111af5760408051737109709ecfa91a80626ff3989d68f67f5b1dd12d602082018190527f6661696c656400000000000000000000000000000000000000000000000000009282019290925260016060820152600091907f70ca10bbd0dbfd9020a9f4b13402c16cb120705e0d1c0aeab10fa353ae586fc490608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261113092916020016116db565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261116891611723565b6000604051808303816000865af19150503d80600081146111a5576040519150601f19603f3d011682016040523d82523d6000602084013e6111aa565b606091505b505050505b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055565b6111e782826112de565b61077a577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f506040516112709060208082526023908201527f4572726f723a2061203d3d2062206e6f7420736174697366696564205b62797460408201527f65735d0000000000000000000000000000000000000000000000000000000000606082015260800190565b60405180910390a17fd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18826040516112a791906117e9565b60405180910390a17fd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf1881604051610e5d9190611838565b8051825160019190036113915760005b835181101561138b5782818151811061130957611309611875565b602001015160f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191684828151811061134857611348611875565b01602001517fff00000000000000000000000000000000000000000000000000000000000000161461137957600091505b80611383816118a4565b9150506112ee565b50611395565b5060005b92915050565b61047b8061190483390190565b61023980611d7f83390190565b6020808252825182820181905260009190848201906040850190845b8181101561140357835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016113d1565b50909695505050505050565b60006020808301818452808551808352604092508286019150828160051b8701018488016000805b848110156114f6578984037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc00186528251805173ffffffffffffffffffffffffffffffffffffffff168552880151888501889052805188860181905290890190839060608701905b808310156114e15783517fffffffff00000000000000000000000000000000000000000000000000000000168252928b019260019290920191908b019061149f565b50978a01979550505091870191600101611437565b50919998505050505050505050565b60005b83811015611520578181015183820152602001611508565b838111156104f75750506000910152565b60008151808452611549816020860160208601611505565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156115ee577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526115dc858351611531565b945092850192908501906001016115a2565b5092979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff8316815260406020820152600061162a6040830184611531565b949350505050565b60006020828403121561164457600080fd5b815173ffffffffffffffffffffffffffffffffffffffff8116811461166857600080fd5b9392505050565b60006020828403121561168157600080fd5b5051919050565b600181811c9082168061169c57607f821691505b6020821081036116d5577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7fffffffff000000000000000000000000000000000000000000000000000000008316815260008251611715816004850160208701611505565b919091016004019392505050565b60008251611735818460208701611505565b9190910192915050565b60006020828403121561175157600080fd5b8151801515811461166857600080fd5b60408152600061179e60408301600a81527f2020202020204c65667400000000000000000000000000000000000000000000602082015260400190565b905082602083015292915050565b60408152600061179e60408301600a81527f2020202020526967687400000000000000000000000000000000000000000000602082015260400190565b60408152600061182660408301600a81527f2020202020204c65667400000000000000000000000000000000000000000000602082015260400190565b828103602084015261162a8185611531565b60408152600061182660408301600a81527f2020202020526967687400000000000000000000000000000000000000000000602082015260400190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036118fc577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b506001019056fe608060405234801561001057600080fd5b5061045b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638381f58a11610081578063b80777ea1161005b578063b80777ea146101a4578063e591b282146101c4578063e81b2c6d1461020457600080fd5b80638381f58a1461017e5780638b239f73146101925780639e8c49661461019b57600080fd5b806354fd4d50116100b257806354fd4d50146100ff5780635cf249691461014857806364ca23ef1461015157600080fd5b8063015d8eb9146100ce57806309bd5a60146100e3575b600080fd5b6100e16100dc366004610369565b61020d565b005b6100ec60025481565b6040519081526020015b60405180910390f35b61013b6040518060400160405280600581526020017f312e312e3000000000000000000000000000000000000000000000000000000081525081565b6040516100f691906103db565b6100ec60015481565b6003546101659067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100f6565b6000546101659067ffffffffffffffff1681565b6100ec60055481565b6100ec60065481565b6000546101659068010000000000000000900467ffffffffffffffff1681565b6101df73deaddeaddeaddeaddeaddeaddeaddeaddead000181565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f6565b6100ec60045481565b3373deaddeaddeaddeaddeaddeaddeaddeaddead0001146102b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603b60248201527f4c31426c6f636b3a206f6e6c7920746865206465706f7369746f72206163636f60448201527f756e742063616e20736574204c3120626c6f636b2076616c7565730000000000606482015260840160405180910390fd5b6000805467ffffffffffffffff98891668010000000000000000027fffffffffffffffffffffffffffffffff00000000000000000000000000000000909116998916999099179890981790975560019490945560029290925560038054919094167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009190911617909255600491909155600555600655565b803567ffffffffffffffff8116811461036457600080fd5b919050565b600080600080600080600080610100898b03121561038657600080fd5b61038f8961034c565b975061039d60208a0161034c565b965060408901359550606089013594506103b960808a0161034c565b979a969950949793969560a0850135955060c08501359460e001359350915050565b600060208083528351808285015260005b81811015610408578581018301518582016040015282016103ec565b8181111561041a576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201604001939250505056fea164736f6c634300080f000a608060405234801561001057600080fd5b50610219806100206000396000f3fe60806040526004361061002d5760003560e01c806354fd4d5014610052578063b9b3efe9146100b157610048565b3661004857600061003c6100d4565b90508060005260206000f35b600061003c6100d4565b34801561005e57600080fd5b5061009b6040518060400160405280600581526020017f312e312e3000000000000000000000000000000000000000000000000000000081525081565b6040516100a89190610168565b60405180910390f35b3480156100bd57600080fd5b506100c66100d4565b6040519081526020016100a8565b600073420000000000000000000000000000000000001573ffffffffffffffffffffffffffffffffffffffff16638381f58a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610135573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061015991906101db565b67ffffffffffffffff16905090565b600060208083528351808285015260005b8181101561019557858101830151858201604001528201610179565b818111156101a7576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b6000602082840312156101ed57600080fd5b815167ffffffffffffffff8116811461020557600080fd5b939250505056fea164736f6c634300080f000aa164736f6c634300080f000a",
}

// L1BlockNumberABI is the input ABI used to generate the binding from.
// Deprecated: Use L1BlockNumberMetaData.ABI instead.
var L1BlockNumberABI = L1BlockNumberMetaData.ABI

// L1BlockNumberBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L1BlockNumberMetaData.Bin instead.
var L1BlockNumberBin = L1BlockNumberMetaData.Bin

// DeployL1BlockNumber deploys a new Ethereum contract, binding an instance of L1BlockNumber to it.
func DeployL1BlockNumber(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *L1BlockNumber, error) {
	parsed, err := L1BlockNumberMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L1BlockNumberBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L1BlockNumber{L1BlockNumberCaller: L1BlockNumberCaller{contract: contract}, L1BlockNumberTransactor: L1BlockNumberTransactor{contract: contract}, L1BlockNumberFilterer: L1BlockNumberFilterer{contract: contract}}, nil
}

// L1BlockNumber is an auto generated Go binding around an Ethereum contract.
type L1BlockNumber struct {
	L1BlockNumberCaller     // Read-only binding to the contract
	L1BlockNumberTransactor // Write-only binding to the contract
	L1BlockNumberFilterer   // Log filterer for contract events
}

// L1BlockNumberCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1BlockNumberCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1BlockNumberTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1BlockNumberTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1BlockNumberFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1BlockNumberFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1BlockNumberSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1BlockNumberSession struct {
	Contract     *L1BlockNumber    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1BlockNumberCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1BlockNumberCallerSession struct {
	Contract *L1BlockNumberCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// L1BlockNumberTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1BlockNumberTransactorSession struct {
	Contract     *L1BlockNumberTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// L1BlockNumberRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1BlockNumberRaw struct {
	Contract *L1BlockNumber // Generic contract binding to access the raw methods on
}

// L1BlockNumberCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1BlockNumberCallerRaw struct {
	Contract *L1BlockNumberCaller // Generic read-only contract binding to access the raw methods on
}

// L1BlockNumberTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1BlockNumberTransactorRaw struct {
	Contract *L1BlockNumberTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1BlockNumber creates a new instance of L1BlockNumber, bound to a specific deployed contract.
func NewL1BlockNumber(address common.Address, backend bind.ContractBackend) (*L1BlockNumber, error) {
	contract, err := bindL1BlockNumber(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1BlockNumber{L1BlockNumberCaller: L1BlockNumberCaller{contract: contract}, L1BlockNumberTransactor: L1BlockNumberTransactor{contract: contract}, L1BlockNumberFilterer: L1BlockNumberFilterer{contract: contract}}, nil
}

// NewL1BlockNumberCaller creates a new read-only instance of L1BlockNumber, bound to a specific deployed contract.
func NewL1BlockNumberCaller(address common.Address, caller bind.ContractCaller) (*L1BlockNumberCaller, error) {
	contract, err := bindL1BlockNumber(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberCaller{contract: contract}, nil
}

// NewL1BlockNumberTransactor creates a new write-only instance of L1BlockNumber, bound to a specific deployed contract.
func NewL1BlockNumberTransactor(address common.Address, transactor bind.ContractTransactor) (*L1BlockNumberTransactor, error) {
	contract, err := bindL1BlockNumber(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberTransactor{contract: contract}, nil
}

// NewL1BlockNumberFilterer creates a new log filterer instance of L1BlockNumber, bound to a specific deployed contract.
func NewL1BlockNumberFilterer(address common.Address, filterer bind.ContractFilterer) (*L1BlockNumberFilterer, error) {
	contract, err := bindL1BlockNumber(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberFilterer{contract: contract}, nil
}

// bindL1BlockNumber binds a generic wrapper to an already deployed contract.
func bindL1BlockNumber(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1BlockNumberMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1BlockNumber *L1BlockNumberRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1BlockNumber.Contract.L1BlockNumberCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1BlockNumber *L1BlockNumberRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1BlockNumber.Contract.L1BlockNumberTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1BlockNumber *L1BlockNumberRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1BlockNumber.Contract.L1BlockNumberTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1BlockNumber *L1BlockNumberCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1BlockNumber.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1BlockNumber *L1BlockNumberTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1BlockNumber.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1BlockNumber *L1BlockNumberTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1BlockNumber.Contract.contract.Transact(opts, method, params...)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_L1BlockNumber *L1BlockNumberCaller) ISTEST(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "IS_TEST")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_L1BlockNumber *L1BlockNumberSession) ISTEST() (bool, error) {
	return _L1BlockNumber.Contract.ISTEST(&_L1BlockNumber.CallOpts)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_L1BlockNumber *L1BlockNumberCallerSession) ISTEST() (bool, error) {
	return _L1BlockNumber.Contract.ISTEST(&_L1BlockNumber.CallOpts)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_L1BlockNumber *L1BlockNumberCaller) ExcludeArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "excludeArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_L1BlockNumber *L1BlockNumberSession) ExcludeArtifacts() ([]string, error) {
	return _L1BlockNumber.Contract.ExcludeArtifacts(&_L1BlockNumber.CallOpts)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_L1BlockNumber *L1BlockNumberCallerSession) ExcludeArtifacts() ([]string, error) {
	return _L1BlockNumber.Contract.ExcludeArtifacts(&_L1BlockNumber.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_L1BlockNumber *L1BlockNumberCaller) ExcludeContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "excludeContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_L1BlockNumber *L1BlockNumberSession) ExcludeContracts() ([]common.Address, error) {
	return _L1BlockNumber.Contract.ExcludeContracts(&_L1BlockNumber.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_L1BlockNumber *L1BlockNumberCallerSession) ExcludeContracts() ([]common.Address, error) {
	return _L1BlockNumber.Contract.ExcludeContracts(&_L1BlockNumber.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_L1BlockNumber *L1BlockNumberCaller) ExcludeSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "excludeSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_L1BlockNumber *L1BlockNumberSession) ExcludeSenders() ([]common.Address, error) {
	return _L1BlockNumber.Contract.ExcludeSenders(&_L1BlockNumber.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_L1BlockNumber *L1BlockNumberCallerSession) ExcludeSenders() ([]common.Address, error) {
	return _L1BlockNumber.Contract.ExcludeSenders(&_L1BlockNumber.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_L1BlockNumber *L1BlockNumberCaller) TargetArtifactSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "targetArtifactSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_L1BlockNumber *L1BlockNumberSession) TargetArtifactSelectors() ([]StdInvariantFuzzSelector, error) {
	return _L1BlockNumber.Contract.TargetArtifactSelectors(&_L1BlockNumber.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_L1BlockNumber *L1BlockNumberCallerSession) TargetArtifactSelectors() ([]StdInvariantFuzzSelector, error) {
	return _L1BlockNumber.Contract.TargetArtifactSelectors(&_L1BlockNumber.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_L1BlockNumber *L1BlockNumberCaller) TargetArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "targetArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_L1BlockNumber *L1BlockNumberSession) TargetArtifacts() ([]string, error) {
	return _L1BlockNumber.Contract.TargetArtifacts(&_L1BlockNumber.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_L1BlockNumber *L1BlockNumberCallerSession) TargetArtifacts() ([]string, error) {
	return _L1BlockNumber.Contract.TargetArtifacts(&_L1BlockNumber.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_L1BlockNumber *L1BlockNumberCaller) TargetContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "targetContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_L1BlockNumber *L1BlockNumberSession) TargetContracts() ([]common.Address, error) {
	return _L1BlockNumber.Contract.TargetContracts(&_L1BlockNumber.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_L1BlockNumber *L1BlockNumberCallerSession) TargetContracts() ([]common.Address, error) {
	return _L1BlockNumber.Contract.TargetContracts(&_L1BlockNumber.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_L1BlockNumber *L1BlockNumberCaller) TargetSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "targetSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_L1BlockNumber *L1BlockNumberSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _L1BlockNumber.Contract.TargetSelectors(&_L1BlockNumber.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_L1BlockNumber *L1BlockNumberCallerSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _L1BlockNumber.Contract.TargetSelectors(&_L1BlockNumber.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_L1BlockNumber *L1BlockNumberCaller) TargetSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _L1BlockNumber.contract.Call(opts, &out, "targetSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_L1BlockNumber *L1BlockNumberSession) TargetSenders() ([]common.Address, error) {
	return _L1BlockNumber.Contract.TargetSenders(&_L1BlockNumber.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_L1BlockNumber *L1BlockNumberCallerSession) TargetSenders() ([]common.Address, error) {
	return _L1BlockNumber.Contract.TargetSenders(&_L1BlockNumber.CallOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_L1BlockNumber *L1BlockNumberTransactor) Failed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1BlockNumber.contract.Transact(opts, "failed")
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_L1BlockNumber *L1BlockNumberSession) Failed() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.Failed(&_L1BlockNumber.TransactOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_L1BlockNumber *L1BlockNumberTransactorSession) Failed() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.Failed(&_L1BlockNumber.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_L1BlockNumber *L1BlockNumberTransactor) SetUp(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1BlockNumber.contract.Transact(opts, "setUp")
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_L1BlockNumber *L1BlockNumberSession) SetUp() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.SetUp(&_L1BlockNumber.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_L1BlockNumber *L1BlockNumberTransactorSession) SetUp() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.SetUp(&_L1BlockNumber.TransactOpts)
}

// TestFallbackSucceeds is a paid mutator transaction binding the contract method 0x66b5a0b8.
//
// Solidity: function test_fallback_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberTransactor) TestFallbackSucceeds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1BlockNumber.contract.Transact(opts, "test_fallback_succeeds")
}

// TestFallbackSucceeds is a paid mutator transaction binding the contract method 0x66b5a0b8.
//
// Solidity: function test_fallback_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberSession) TestFallbackSucceeds() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.TestFallbackSucceeds(&_L1BlockNumber.TransactOpts)
}

// TestFallbackSucceeds is a paid mutator transaction binding the contract method 0x66b5a0b8.
//
// Solidity: function test_fallback_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberTransactorSession) TestFallbackSucceeds() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.TestFallbackSucceeds(&_L1BlockNumber.TransactOpts)
}

// TestGetL1BlockNumberSucceeds is a paid mutator transaction binding the contract method 0x573980d0.
//
// Solidity: function test_getL1BlockNumber_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberTransactor) TestGetL1BlockNumberSucceeds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1BlockNumber.contract.Transact(opts, "test_getL1BlockNumber_succeeds")
}

// TestGetL1BlockNumberSucceeds is a paid mutator transaction binding the contract method 0x573980d0.
//
// Solidity: function test_getL1BlockNumber_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberSession) TestGetL1BlockNumberSucceeds() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.TestGetL1BlockNumberSucceeds(&_L1BlockNumber.TransactOpts)
}

// TestGetL1BlockNumberSucceeds is a paid mutator transaction binding the contract method 0x573980d0.
//
// Solidity: function test_getL1BlockNumber_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberTransactorSession) TestGetL1BlockNumberSucceeds() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.TestGetL1BlockNumberSucceeds(&_L1BlockNumber.TransactOpts)
}

// TestReceiveSucceeds is a paid mutator transaction binding the contract method 0xfe43c9dd.
//
// Solidity: function test_receive_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberTransactor) TestReceiveSucceeds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1BlockNumber.contract.Transact(opts, "test_receive_succeeds")
}

// TestReceiveSucceeds is a paid mutator transaction binding the contract method 0xfe43c9dd.
//
// Solidity: function test_receive_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberSession) TestReceiveSucceeds() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.TestReceiveSucceeds(&_L1BlockNumber.TransactOpts)
}

// TestReceiveSucceeds is a paid mutator transaction binding the contract method 0xfe43c9dd.
//
// Solidity: function test_receive_succeeds() returns()
func (_L1BlockNumber *L1BlockNumberTransactorSession) TestReceiveSucceeds() (*types.Transaction, error) {
	return _L1BlockNumber.Contract.TestReceiveSucceeds(&_L1BlockNumber.TransactOpts)
}

// L1BlockNumberLogIterator is returned from FilterLog and is used to iterate over the raw logs and unpacked data for Log events raised by the L1BlockNumber contract.
type L1BlockNumberLogIterator struct {
	Event *L1BlockNumberLog // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLog)
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
		it.Event = new(L1BlockNumberLog)
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
func (it *L1BlockNumberLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLog represents a Log event raised by the L1BlockNumber contract.
type L1BlockNumberLog struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLog is a free log retrieval operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLog(opts *bind.FilterOpts) (*L1BlockNumberLogIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogIterator{contract: _L1BlockNumber.contract, event: "log", logs: logs, sub: sub}, nil
}

// WatchLog is a free log subscription operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLog(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLog) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLog)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLog(log types.Log) (*L1BlockNumberLog, error) {
	event := new(L1BlockNumberLog)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogAddressIterator is returned from FilterLogAddress and is used to iterate over the raw logs and unpacked data for LogAddress events raised by the L1BlockNumber contract.
type L1BlockNumberLogAddressIterator struct {
	Event *L1BlockNumberLogAddress // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogAddress)
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
		it.Event = new(L1BlockNumberLogAddress)
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
func (it *L1BlockNumberLogAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogAddress represents a LogAddress event raised by the L1BlockNumber contract.
type L1BlockNumberLogAddress struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogAddress is a free log retrieval operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogAddress(opts *bind.FilterOpts) (*L1BlockNumberLogAddressIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogAddressIterator{contract: _L1BlockNumber.contract, event: "log_address", logs: logs, sub: sub}, nil
}

// WatchLogAddress is a free log subscription operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogAddress(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogAddress) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogAddress)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_address", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogAddress(log types.Log) (*L1BlockNumberLogAddress, error) {
	event := new(L1BlockNumberLogAddress)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogArrayIterator is returned from FilterLogArray and is used to iterate over the raw logs and unpacked data for LogArray events raised by the L1BlockNumber contract.
type L1BlockNumberLogArrayIterator struct {
	Event *L1BlockNumberLogArray // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogArray)
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
		it.Event = new(L1BlockNumberLogArray)
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
func (it *L1BlockNumberLogArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogArray represents a LogArray event raised by the L1BlockNumber contract.
type L1BlockNumberLogArray struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray is a free log retrieval operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogArray(opts *bind.FilterOpts) (*L1BlockNumberLogArrayIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogArrayIterator{contract: _L1BlockNumber.contract, event: "log_array", logs: logs, sub: sub}, nil
}

// WatchLogArray is a free log subscription operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogArray(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogArray) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogArray)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_array", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogArray(log types.Log) (*L1BlockNumberLogArray, error) {
	event := new(L1BlockNumberLogArray)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogArray0Iterator is returned from FilterLogArray0 and is used to iterate over the raw logs and unpacked data for LogArray0 events raised by the L1BlockNumber contract.
type L1BlockNumberLogArray0Iterator struct {
	Event *L1BlockNumberLogArray0 // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogArray0)
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
		it.Event = new(L1BlockNumberLogArray0)
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
func (it *L1BlockNumberLogArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogArray0 represents a LogArray0 event raised by the L1BlockNumber contract.
type L1BlockNumberLogArray0 struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray0 is a free log retrieval operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogArray0(opts *bind.FilterOpts) (*L1BlockNumberLogArray0Iterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogArray0Iterator{contract: _L1BlockNumber.contract, event: "log_array0", logs: logs, sub: sub}, nil
}

// WatchLogArray0 is a free log subscription operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogArray0(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogArray0) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogArray0)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_array0", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogArray0(log types.Log) (*L1BlockNumberLogArray0, error) {
	event := new(L1BlockNumberLogArray0)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogArray1Iterator is returned from FilterLogArray1 and is used to iterate over the raw logs and unpacked data for LogArray1 events raised by the L1BlockNumber contract.
type L1BlockNumberLogArray1Iterator struct {
	Event *L1BlockNumberLogArray1 // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogArray1)
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
		it.Event = new(L1BlockNumberLogArray1)
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
func (it *L1BlockNumberLogArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogArray1 represents a LogArray1 event raised by the L1BlockNumber contract.
type L1BlockNumberLogArray1 struct {
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray1 is a free log retrieval operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogArray1(opts *bind.FilterOpts) (*L1BlockNumberLogArray1Iterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogArray1Iterator{contract: _L1BlockNumber.contract, event: "log_array1", logs: logs, sub: sub}, nil
}

// WatchLogArray1 is a free log subscription operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogArray1(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogArray1) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogArray1)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_array1", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogArray1(log types.Log) (*L1BlockNumberLogArray1, error) {
	event := new(L1BlockNumberLogArray1)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogBytesIterator is returned from FilterLogBytes and is used to iterate over the raw logs and unpacked data for LogBytes events raised by the L1BlockNumber contract.
type L1BlockNumberLogBytesIterator struct {
	Event *L1BlockNumberLogBytes // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogBytes)
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
		it.Event = new(L1BlockNumberLogBytes)
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
func (it *L1BlockNumberLogBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogBytes represents a LogBytes event raised by the L1BlockNumber contract.
type L1BlockNumberLogBytes struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes is a free log retrieval operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogBytes(opts *bind.FilterOpts) (*L1BlockNumberLogBytesIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogBytesIterator{contract: _L1BlockNumber.contract, event: "log_bytes", logs: logs, sub: sub}, nil
}

// WatchLogBytes is a free log subscription operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogBytes(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogBytes) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogBytes)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_bytes", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogBytes(log types.Log) (*L1BlockNumberLogBytes, error) {
	event := new(L1BlockNumberLogBytes)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogBytes32Iterator is returned from FilterLogBytes32 and is used to iterate over the raw logs and unpacked data for LogBytes32 events raised by the L1BlockNumber contract.
type L1BlockNumberLogBytes32Iterator struct {
	Event *L1BlockNumberLogBytes32 // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogBytes32)
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
		it.Event = new(L1BlockNumberLogBytes32)
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
func (it *L1BlockNumberLogBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogBytes32 represents a LogBytes32 event raised by the L1BlockNumber contract.
type L1BlockNumberLogBytes32 struct {
	Arg0 [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes32 is a free log retrieval operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogBytes32(opts *bind.FilterOpts) (*L1BlockNumberLogBytes32Iterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogBytes32Iterator{contract: _L1BlockNumber.contract, event: "log_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogBytes32 is a free log subscription operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogBytes32(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogBytes32) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogBytes32)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_bytes32", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogBytes32(log types.Log) (*L1BlockNumberLogBytes32, error) {
	event := new(L1BlockNumberLogBytes32)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogIntIterator is returned from FilterLogInt and is used to iterate over the raw logs and unpacked data for LogInt events raised by the L1BlockNumber contract.
type L1BlockNumberLogIntIterator struct {
	Event *L1BlockNumberLogInt // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogInt)
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
		it.Event = new(L1BlockNumberLogInt)
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
func (it *L1BlockNumberLogIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogInt represents a LogInt event raised by the L1BlockNumber contract.
type L1BlockNumberLogInt struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogInt is a free log retrieval operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogInt(opts *bind.FilterOpts) (*L1BlockNumberLogIntIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogIntIterator{contract: _L1BlockNumber.contract, event: "log_int", logs: logs, sub: sub}, nil
}

// WatchLogInt is a free log subscription operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogInt(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogInt) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogInt)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_int", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogInt(log types.Log) (*L1BlockNumberLogInt, error) {
	event := new(L1BlockNumberLogInt)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedAddressIterator is returned from FilterLogNamedAddress and is used to iterate over the raw logs and unpacked data for LogNamedAddress events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedAddressIterator struct {
	Event *L1BlockNumberLogNamedAddress // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedAddress)
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
		it.Event = new(L1BlockNumberLogNamedAddress)
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
func (it *L1BlockNumberLogNamedAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedAddress represents a LogNamedAddress event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedAddress struct {
	Key string
	Val common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedAddress is a free log retrieval operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedAddress(opts *bind.FilterOpts) (*L1BlockNumberLogNamedAddressIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedAddressIterator{contract: _L1BlockNumber.contract, event: "log_named_address", logs: logs, sub: sub}, nil
}

// WatchLogNamedAddress is a free log subscription operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedAddress(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedAddress) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedAddress)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_address", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedAddress(log types.Log) (*L1BlockNumberLogNamedAddress, error) {
	event := new(L1BlockNumberLogNamedAddress)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedArrayIterator is returned from FilterLogNamedArray and is used to iterate over the raw logs and unpacked data for LogNamedArray events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedArrayIterator struct {
	Event *L1BlockNumberLogNamedArray // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedArray)
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
		it.Event = new(L1BlockNumberLogNamedArray)
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
func (it *L1BlockNumberLogNamedArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedArray represents a LogNamedArray event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedArray struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray is a free log retrieval operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedArray(opts *bind.FilterOpts) (*L1BlockNumberLogNamedArrayIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedArrayIterator{contract: _L1BlockNumber.contract, event: "log_named_array", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray is a free log subscription operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedArray(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedArray) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedArray)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_array", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedArray(log types.Log) (*L1BlockNumberLogNamedArray, error) {
	event := new(L1BlockNumberLogNamedArray)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedArray0Iterator is returned from FilterLogNamedArray0 and is used to iterate over the raw logs and unpacked data for LogNamedArray0 events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedArray0Iterator struct {
	Event *L1BlockNumberLogNamedArray0 // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedArray0)
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
		it.Event = new(L1BlockNumberLogNamedArray0)
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
func (it *L1BlockNumberLogNamedArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedArray0 represents a LogNamedArray0 event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedArray0 struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray0 is a free log retrieval operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedArray0(opts *bind.FilterOpts) (*L1BlockNumberLogNamedArray0Iterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedArray0Iterator{contract: _L1BlockNumber.contract, event: "log_named_array0", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray0 is a free log subscription operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedArray0(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedArray0) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedArray0)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_array0", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedArray0(log types.Log) (*L1BlockNumberLogNamedArray0, error) {
	event := new(L1BlockNumberLogNamedArray0)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedArray1Iterator is returned from FilterLogNamedArray1 and is used to iterate over the raw logs and unpacked data for LogNamedArray1 events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedArray1Iterator struct {
	Event *L1BlockNumberLogNamedArray1 // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedArray1)
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
		it.Event = new(L1BlockNumberLogNamedArray1)
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
func (it *L1BlockNumberLogNamedArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedArray1 represents a LogNamedArray1 event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedArray1 struct {
	Key string
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray1 is a free log retrieval operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedArray1(opts *bind.FilterOpts) (*L1BlockNumberLogNamedArray1Iterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedArray1Iterator{contract: _L1BlockNumber.contract, event: "log_named_array1", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray1 is a free log subscription operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedArray1(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedArray1) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedArray1)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_array1", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedArray1(log types.Log) (*L1BlockNumberLogNamedArray1, error) {
	event := new(L1BlockNumberLogNamedArray1)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedBytesIterator is returned from FilterLogNamedBytes and is used to iterate over the raw logs and unpacked data for LogNamedBytes events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedBytesIterator struct {
	Event *L1BlockNumberLogNamedBytes // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedBytes)
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
		it.Event = new(L1BlockNumberLogNamedBytes)
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
func (it *L1BlockNumberLogNamedBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedBytes represents a LogNamedBytes event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedBytes struct {
	Key string
	Val []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes is a free log retrieval operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedBytes(opts *bind.FilterOpts) (*L1BlockNumberLogNamedBytesIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedBytesIterator{contract: _L1BlockNumber.contract, event: "log_named_bytes", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes is a free log subscription operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedBytes(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedBytes) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedBytes)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedBytes(log types.Log) (*L1BlockNumberLogNamedBytes, error) {
	event := new(L1BlockNumberLogNamedBytes)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedBytes32Iterator is returned from FilterLogNamedBytes32 and is used to iterate over the raw logs and unpacked data for LogNamedBytes32 events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedBytes32Iterator struct {
	Event *L1BlockNumberLogNamedBytes32 // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedBytes32)
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
		it.Event = new(L1BlockNumberLogNamedBytes32)
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
func (it *L1BlockNumberLogNamedBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedBytes32 represents a LogNamedBytes32 event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedBytes32 struct {
	Key string
	Val [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes32 is a free log retrieval operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedBytes32(opts *bind.FilterOpts) (*L1BlockNumberLogNamedBytes32Iterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedBytes32Iterator{contract: _L1BlockNumber.contract, event: "log_named_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes32 is a free log subscription operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedBytes32(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedBytes32) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedBytes32)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedBytes32(log types.Log) (*L1BlockNumberLogNamedBytes32, error) {
	event := new(L1BlockNumberLogNamedBytes32)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedDecimalIntIterator is returned from FilterLogNamedDecimalInt and is used to iterate over the raw logs and unpacked data for LogNamedDecimalInt events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedDecimalIntIterator struct {
	Event *L1BlockNumberLogNamedDecimalInt // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedDecimalIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedDecimalInt)
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
		it.Event = new(L1BlockNumberLogNamedDecimalInt)
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
func (it *L1BlockNumberLogNamedDecimalIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedDecimalIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedDecimalInt represents a LogNamedDecimalInt event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedDecimalInt struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalInt is a free log retrieval operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedDecimalInt(opts *bind.FilterOpts) (*L1BlockNumberLogNamedDecimalIntIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedDecimalIntIterator{contract: _L1BlockNumber.contract, event: "log_named_decimal_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalInt is a free log subscription operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedDecimalInt(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedDecimalInt) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedDecimalInt)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedDecimalInt(log types.Log) (*L1BlockNumberLogNamedDecimalInt, error) {
	event := new(L1BlockNumberLogNamedDecimalInt)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedDecimalUintIterator is returned from FilterLogNamedDecimalUint and is used to iterate over the raw logs and unpacked data for LogNamedDecimalUint events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedDecimalUintIterator struct {
	Event *L1BlockNumberLogNamedDecimalUint // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedDecimalUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedDecimalUint)
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
		it.Event = new(L1BlockNumberLogNamedDecimalUint)
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
func (it *L1BlockNumberLogNamedDecimalUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedDecimalUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedDecimalUint represents a LogNamedDecimalUint event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedDecimalUint struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalUint is a free log retrieval operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedDecimalUint(opts *bind.FilterOpts) (*L1BlockNumberLogNamedDecimalUintIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedDecimalUintIterator{contract: _L1BlockNumber.contract, event: "log_named_decimal_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalUint is a free log subscription operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedDecimalUint(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedDecimalUint) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedDecimalUint)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedDecimalUint(log types.Log) (*L1BlockNumberLogNamedDecimalUint, error) {
	event := new(L1BlockNumberLogNamedDecimalUint)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedIntIterator is returned from FilterLogNamedInt and is used to iterate over the raw logs and unpacked data for LogNamedInt events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedIntIterator struct {
	Event *L1BlockNumberLogNamedInt // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedInt)
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
		it.Event = new(L1BlockNumberLogNamedInt)
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
func (it *L1BlockNumberLogNamedIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedInt represents a LogNamedInt event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedInt struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedInt is a free log retrieval operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedInt(opts *bind.FilterOpts) (*L1BlockNumberLogNamedIntIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedIntIterator{contract: _L1BlockNumber.contract, event: "log_named_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedInt is a free log subscription operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedInt(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedInt) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedInt)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_int", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedInt(log types.Log) (*L1BlockNumberLogNamedInt, error) {
	event := new(L1BlockNumberLogNamedInt)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedStringIterator is returned from FilterLogNamedString and is used to iterate over the raw logs and unpacked data for LogNamedString events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedStringIterator struct {
	Event *L1BlockNumberLogNamedString // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedString)
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
		it.Event = new(L1BlockNumberLogNamedString)
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
func (it *L1BlockNumberLogNamedStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedString represents a LogNamedString event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedString struct {
	Key string
	Val string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedString is a free log retrieval operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedString(opts *bind.FilterOpts) (*L1BlockNumberLogNamedStringIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedStringIterator{contract: _L1BlockNumber.contract, event: "log_named_string", logs: logs, sub: sub}, nil
}

// WatchLogNamedString is a free log subscription operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedString(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedString) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedString)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_string", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedString(log types.Log) (*L1BlockNumberLogNamedString, error) {
	event := new(L1BlockNumberLogNamedString)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogNamedUintIterator is returned from FilterLogNamedUint and is used to iterate over the raw logs and unpacked data for LogNamedUint events raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedUintIterator struct {
	Event *L1BlockNumberLogNamedUint // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogNamedUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogNamedUint)
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
		it.Event = new(L1BlockNumberLogNamedUint)
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
func (it *L1BlockNumberLogNamedUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogNamedUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogNamedUint represents a LogNamedUint event raised by the L1BlockNumber contract.
type L1BlockNumberLogNamedUint struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedUint is a free log retrieval operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogNamedUint(opts *bind.FilterOpts) (*L1BlockNumberLogNamedUintIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogNamedUintIterator{contract: _L1BlockNumber.contract, event: "log_named_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedUint is a free log subscription operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogNamedUint(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogNamedUint) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogNamedUint)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_uint", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogNamedUint(log types.Log) (*L1BlockNumberLogNamedUint, error) {
	event := new(L1BlockNumberLogNamedUint)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_named_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogStringIterator is returned from FilterLogString and is used to iterate over the raw logs and unpacked data for LogString events raised by the L1BlockNumber contract.
type L1BlockNumberLogStringIterator struct {
	Event *L1BlockNumberLogString // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogString)
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
		it.Event = new(L1BlockNumberLogString)
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
func (it *L1BlockNumberLogStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogString represents a LogString event raised by the L1BlockNumber contract.
type L1BlockNumberLogString struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogString is a free log retrieval operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogString(opts *bind.FilterOpts) (*L1BlockNumberLogStringIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogStringIterator{contract: _L1BlockNumber.contract, event: "log_string", logs: logs, sub: sub}, nil
}

// WatchLogString is a free log subscription operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogString(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogString) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogString)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_string", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogString(log types.Log) (*L1BlockNumberLogString, error) {
	event := new(L1BlockNumberLogString)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogUintIterator is returned from FilterLogUint and is used to iterate over the raw logs and unpacked data for LogUint events raised by the L1BlockNumber contract.
type L1BlockNumberLogUintIterator struct {
	Event *L1BlockNumberLogUint // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogUint)
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
		it.Event = new(L1BlockNumberLogUint)
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
func (it *L1BlockNumberLogUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogUint represents a LogUint event raised by the L1BlockNumber contract.
type L1BlockNumberLogUint struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogUint is a free log retrieval operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogUint(opts *bind.FilterOpts) (*L1BlockNumberLogUintIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogUintIterator{contract: _L1BlockNumber.contract, event: "log_uint", logs: logs, sub: sub}, nil
}

// WatchLogUint is a free log subscription operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogUint(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogUint) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogUint)
				if err := _L1BlockNumber.contract.UnpackLog(event, "log_uint", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogUint(log types.Log) (*L1BlockNumberLogUint, error) {
	event := new(L1BlockNumberLogUint)
	if err := _L1BlockNumber.contract.UnpackLog(event, "log_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1BlockNumberLogsIterator is returned from FilterLogs and is used to iterate over the raw logs and unpacked data for Logs events raised by the L1BlockNumber contract.
type L1BlockNumberLogsIterator struct {
	Event *L1BlockNumberLogs // Event containing the contract specifics and raw log

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
func (it *L1BlockNumberLogsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1BlockNumberLogs)
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
		it.Event = new(L1BlockNumberLogs)
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
func (it *L1BlockNumberLogsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1BlockNumberLogsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1BlockNumberLogs represents a Logs event raised by the L1BlockNumber contract.
type L1BlockNumberLogs struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogs is a free log retrieval operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) FilterLogs(opts *bind.FilterOpts) (*L1BlockNumberLogsIterator, error) {

	logs, sub, err := _L1BlockNumber.contract.FilterLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return &L1BlockNumberLogsIterator{contract: _L1BlockNumber.contract, event: "logs", logs: logs, sub: sub}, nil
}

// WatchLogs is a free log subscription operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_L1BlockNumber *L1BlockNumberFilterer) WatchLogs(opts *bind.WatchOpts, sink chan<- *L1BlockNumberLogs) (event.Subscription, error) {

	logs, sub, err := _L1BlockNumber.contract.WatchLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1BlockNumberLogs)
				if err := _L1BlockNumber.contract.UnpackLog(event, "logs", log); err != nil {
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
func (_L1BlockNumber *L1BlockNumberFilterer) ParseLogs(log types.Log) (*L1BlockNumberLogs, error) {
	event := new(L1BlockNumberLogs)
	if err := _L1BlockNumber.contract.UnpackLog(event, "logs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
