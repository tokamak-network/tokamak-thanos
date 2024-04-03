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

// BlockOracleMetaData contains all meta data concerning the BlockOracle contract.
var BlockOracleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"Hash\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"Timestamp\",\"name\":\"childTimestamp\",\"type\":\"uint64\"}],\"name\":\"Checkpoint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"log_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"val\",\"type\":\"uint256[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"val\",\"type\":\"int256[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"val\",\"type\":\"address[]\"}],\"name\":\"log_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"log_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"log_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"log_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"val\",\"type\":\"address\"}],\"name\":\"log_named_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"val\",\"type\":\"uint256[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"val\",\"type\":\"int256[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"val\",\"type\":\"address[]\"}],\"name\":\"log_named_array\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"val\",\"type\":\"bytes\"}],\"name\":\"log_named_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"val\",\"type\":\"bytes32\"}],\"name\":\"log_named_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"}],\"name\":\"log_named_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"val\",\"type\":\"string\"}],\"name\":\"log_named_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"}],\"name\":\"log_named_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"log_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"logs\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"IS_TEST\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeArtifacts\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"excludedArtifacts_\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"excludedContracts_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"excludeSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"excludedSenders_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"failed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setUp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetArtifactSelectors\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes4[]\",\"name\":\"selectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"name\":\"targetedArtifactSelectors_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetArtifacts\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"targetedArtifacts_\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targetedContracts_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetSelectors\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes4[]\",\"name\":\"selectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structStdInvariant.FuzzSelector[]\",\"name\":\"targetedSelectors_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targetedSenders_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test_checkpointAndLoad_succeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test_load_noBlockHash_reverts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405260008054600160ff19918216811790925560048054909116909117905534801561002d57600080fd5b506118708061003d6000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c8063916a17c61161008c578063cabc1a5411610066578063cabc1a541461016e578063e20c9f7114610176578063f89a51bb1461017e578063fa7626d41461018657600080fd5b8063916a17c614610146578063b5508aa91461014e578063ba414fa61461015657600080fd5b80633f7286f4116100bd5780633f7286f41461011457806366d9a9a01461011c57806385226c811461013157600080fd5b80630a9254e4146100e45780631ed7831c146100ee5780633e5e3c231461010c575b600080fd5b6100ec610193565b005b6100f66102f0565b60405161010391906110c4565b60405180910390f35b6100f661035f565b6100f66103cc565b610124610439565b604051610103919061111e565b61013961054a565b6040516101039190611240565b61012461061a565b610139610722565b61015e6107f2565b6040519015158152602001610103565b6100ec610952565b6100f6610a84565b6100ec610af1565b60005461015e9060ff1681565b60405161019f906110b7565b604051809103906000f0801580156101bb573d6000803e3d6000fd5b50601b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055737109709ecfa91a80626ff3989d68f67f5b1dd12d631f7b4f30610226436001611325565b6040518263ffffffff1660e01b815260040161024491815260200190565b600060405180830381600087803b15801561025e57600080fd5b505af1158015610272573d6000803e3d6000fd5b50737109709ecfa91a80626ff3989d68f67f5b1dd12d925063e5d6bf02915061029e905042600d611325565b6040518263ffffffff1660e01b81526004016102bc91815260200190565b600060405180830381600087803b1580156102d657600080fd5b505af11580156102ea573d6000803e3d6000fd5b50505050565b6060600d80548060200260200160405190810160405280929190818152602001828054801561035557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161032a575b5050505050905090565b6060600f8054806020026020016040519081016040528092919081815260200182805480156103555760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161032a575050505050905090565b6060600e8054806020026020016040519081016040528092919081815260200182805480156103555760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161032a575050505050905090565b60606012805480602002602001604051908101604052809291908181526020016000905b8282101561054157600084815260209081902060408051808201825260028602909201805473ffffffffffffffffffffffffffffffffffffffff16835260018101805483518187028101870190945280845293949193858301939283018282801561052957602002820191906000526020600020906000905b82829054906101000a900460e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190600401906020826003010492830192600103820291508084116104d65790505b5050505050815250508152602001906001019061045d565b50505050905090565b60606011805480602002602001604051908101604052809291908181526020016000905b8282101561054157838290600052602060002001805461058d9061133d565b80601f01602080910402602001604051908101604052809291908181526020018280546105b99061133d565b80156106065780601f106105db57610100808354040283529160200191610606565b820191906000526020600020905b8154815290600101906020018083116105e957829003601f168201915b50505050508152602001906001019061056e565b60606013805480602002602001604051908101604052809291908181526020016000905b8282101561054157600084815260209081902060408051808201825260028602909201805473ffffffffffffffffffffffffffffffffffffffff16835260018101805483518187028101870190945280845293949193858301939283018282801561070a57602002820191906000526020600020906000905b82829054906101000a900460e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190600401906020826003010492830192600103820291508084116106b75790505b5050505050815250508152602001906001019061063e565b60606010805480602002602001604051908101604052809291908181526020016000905b828210156105415783829060005260206000200180546107659061133d565b80601f01602080910402602001604051908101604052809291908181526020018280546107919061133d565b80156107de5780601f106107b3576101008083540402835291602001916107de565b820191906000526020600020905b8154815290600101906020018083116107c157829003601f168201915b505050505081526020019060010190610746565b60008054610100900460ff16156108125750600054610100900460ff1690565b6000737109709ecfa91a80626ff3989d68f67f5b1dd12d3b1561094d5760408051737109709ecfa91a80626ff3989d68f67f5b1dd12d602082018190527f6661696c65640000000000000000000000000000000000000000000000000000828401528251808303840181526060830190935260009290916108b7917f667f9d70ca411d70ead50d8d5c22070dafc36ad75f3dcf5e7237b22ade9aecc491608001611390565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526108ef916113d8565b6000604051808303816000865af19150503d806000811461092c576040519150601f19603f3d011682016040523d82523d6000602084013e610931565b606091505b509150508080602001905181019061094991906113f4565b9150505b919050565b6040517fc31eb0e00000000000000000000000000000000000000000000000000000000081527f37cf2705000000000000000000000000000000000000000000000000000000006004820152737109709ecfa91a80626ff3989d68f67f5b1dd12d9063c31eb0e090602401600060405180830381600087803b1580156109d757600080fd5b505af11580156109eb573d6000803e3d6000fd5b5050601b546040517f99d548aa0000000000000000000000000000000000000000000000000000000081526000600482015273ffffffffffffffffffffffffffffffffffffffff90911692506399d548aa91506024016040805180830381865afa158015610a5d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a81919061141d565b50565b6060600c8054806020026020016040519081016040528092919081815260200182805480156103555760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161032a575050505050905090565b6040517f491cc7c200000000000000000000000000000000000000000000000000000000815260016004820181905260248201819052604482015260006064820152737109709ecfa91a80626ff3989d68f67f5b1dd12d9063491cc7c290608401600060405180830381600087803b158015610b6c57600080fd5b505af1158015610b80573d6000803e3d6000fd5b505050504267ffffffffffffffff16600143610b9c91906114a4565b40610ba86001436114a4565b6040517fb67ff58b33060fd371a35ae2d9f1c3cdaec9b8197969f6efe2594a1ff4ba68c690600090a4601b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c2c4c5c16040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610c40573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c6491906114bb565b506000610c726001436114a4565b601b546040517f99d548aa0000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff909116906399d548aa906024016040805180830381865afa158015610ce6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d0a919061141d565b9050610d1b81600001518340610d37565b610d33816020015167ffffffffffffffff1642610e41565b5050565b808214610d33577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50604051610dc39060208082526025908201527f4572726f723a2061203d3d2062206e6f7420736174697366696564205b62797460408201527f657333325d000000000000000000000000000000000000000000000000000000606082015260800190565b60405180910390a17fafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f9982604051610dfa91906114d4565b60405180910390a17fafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f9981604051610e31919061151f565b60405180910390a1610d33610f3b565b808214610d33577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50604051610ecd9060208082526022908201527f4572726f723a2061203d3d2062206e6f7420736174697366696564205b75696e60408201527f745d000000000000000000000000000000000000000000000000000000000000606082015260800190565b60405180910390a17fb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a882604051610f0491906114d4565b60405180910390a17fb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a881604051610e31919061151f565b737109709ecfa91a80626ff3989d68f67f5b1dd12d3b156110895760408051737109709ecfa91a80626ff3989d68f67f5b1dd12d602082018190527f6661696c656400000000000000000000000000000000000000000000000000009282019290925260016060820152600091907f70ca10bbd0dbfd9020a9f4b13402c16cb120705e0d1c0aeab10fa353ae586fc490608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261100a9291602001611390565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052611042916113d8565b6000604051808303816000865af19150503d806000811461107f576040519150601f19603f3d011682016040523d82523d6000602084013e611084565b606091505b505050505b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055565b6103078061155d83390190565b6020808252825182820181905260009190848201906040850190845b8181101561111257835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016110e0565b50909695505050505050565b60006020808301818452808551808352604092508286019150828160051b8701018488016000805b84811015611205578984037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc00186528251805173ffffffffffffffffffffffffffffffffffffffff168552880151888501889052805188860181905290890190839060608701905b808310156111f05783517fffffffff00000000000000000000000000000000000000000000000000000000168252928b019260019290920191908b01906111ae565b50978a01979550505091870191600101611146565b50919998505050505050505050565b60005b8381101561122f578181015183820152602001611217565b838111156102ea5750506000910152565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156112e9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452815180518087526112ac818989018a8501611214565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01695909501860194509285019290850190600101611267565b5092979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115611338576113386112f6565b500190565b600181811c9082168061135157607f821691505b60208210810361138a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7fffffffff0000000000000000000000000000000000000000000000000000000083168152600082516113ca816004850160208701611214565b919091016004019392505050565b600082516113ea818460208701611214565b9190910192915050565b60006020828403121561140657600080fd5b8151801515811461141657600080fd5b9392505050565b60006040828403121561142f57600080fd5b6040516040810167ffffffffffffffff828210818311171561147a577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b816040528451835260208501519150808216821461149757600080fd5b5060208201529392505050565b6000828210156114b6576114b66112f6565b500390565b6000602082840312156114cd57600080fd5b5051919050565b60408152600061151160408301600a81527f2020202020204c65667400000000000000000000000000000000000000000000602082015260400190565b905082602083015292915050565b60408152600061151160408301600a81527f202020202052696768740000000000000000000000000000000000000000000060208201526040019056fe608060405234801561001057600080fd5b506102e7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806354fd4d501461004657806399d548aa14610098578063c2c4c5c1146100d0575b600080fd5b6100826040518060400160405280600581526020017f302e302e3100000000000000000000000000000000000000000000000000000081525081565b60405161008f9190610210565b60405180910390f35b6100ab6100a6366004610283565b6100e6565b604080518251815260209283015167ffffffffffffffff16928101929092520161008f565b6100d8610165565b60405190815260200161008f565b604080518082018252600080825260209182018190528381528082528281208351808501909452805480855260019091015467ffffffffffffffff169284019290925203610160576040517f37cf270500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b919050565b600061017260014361029c565b60408051808201825282408082524267ffffffffffffffff81811660208086018281526000898152918290528782209651875551600190960180547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016969093169590951790915593519495509093909291849186917fb67ff58b33060fd371a35ae2d9f1c3cdaec9b8197969f6efe2594a1ff4ba68c691a4505090565b600060208083528351808285015260005b8181101561023d57858101830151858201604001528201610221565b8181111561024f576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60006020828403121561029557600080fd5b5035919050565b6000828210156102d5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50039056fea164736f6c634300080f000aa164736f6c634300080f000a",
}

// BlockOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use BlockOracleMetaData.ABI instead.
var BlockOracleABI = BlockOracleMetaData.ABI

// BlockOracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BlockOracleMetaData.Bin instead.
var BlockOracleBin = BlockOracleMetaData.Bin

// DeployBlockOracle deploys a new Ethereum contract, binding an instance of BlockOracle to it.
func DeployBlockOracle(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BlockOracle, error) {
	parsed, err := BlockOracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BlockOracleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BlockOracle{BlockOracleCaller: BlockOracleCaller{contract: contract}, BlockOracleTransactor: BlockOracleTransactor{contract: contract}, BlockOracleFilterer: BlockOracleFilterer{contract: contract}}, nil
}

// BlockOracle is an auto generated Go binding around an Ethereum contract.
type BlockOracle struct {
	BlockOracleCaller     // Read-only binding to the contract
	BlockOracleTransactor // Write-only binding to the contract
	BlockOracleFilterer   // Log filterer for contract events
}

// BlockOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlockOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlockOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlockOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlockOracleSession struct {
	Contract     *BlockOracle      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlockOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlockOracleCallerSession struct {
	Contract *BlockOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BlockOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlockOracleTransactorSession struct {
	Contract     *BlockOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BlockOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlockOracleRaw struct {
	Contract *BlockOracle // Generic contract binding to access the raw methods on
}

// BlockOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlockOracleCallerRaw struct {
	Contract *BlockOracleCaller // Generic read-only contract binding to access the raw methods on
}

// BlockOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlockOracleTransactorRaw struct {
	Contract *BlockOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlockOracle creates a new instance of BlockOracle, bound to a specific deployed contract.
func NewBlockOracle(address common.Address, backend bind.ContractBackend) (*BlockOracle, error) {
	contract, err := bindBlockOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BlockOracle{BlockOracleCaller: BlockOracleCaller{contract: contract}, BlockOracleTransactor: BlockOracleTransactor{contract: contract}, BlockOracleFilterer: BlockOracleFilterer{contract: contract}}, nil
}

// NewBlockOracleCaller creates a new read-only instance of BlockOracle, bound to a specific deployed contract.
func NewBlockOracleCaller(address common.Address, caller bind.ContractCaller) (*BlockOracleCaller, error) {
	contract, err := bindBlockOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlockOracleCaller{contract: contract}, nil
}

// NewBlockOracleTransactor creates a new write-only instance of BlockOracle, bound to a specific deployed contract.
func NewBlockOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockOracleTransactor, error) {
	contract, err := bindBlockOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlockOracleTransactor{contract: contract}, nil
}

// NewBlockOracleFilterer creates a new log filterer instance of BlockOracle, bound to a specific deployed contract.
func NewBlockOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*BlockOracleFilterer, error) {
	contract, err := bindBlockOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlockOracleFilterer{contract: contract}, nil
}

// bindBlockOracle binds a generic wrapper to an already deployed contract.
func bindBlockOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BlockOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockOracle *BlockOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockOracle.Contract.BlockOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockOracle *BlockOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockOracle.Contract.BlockOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockOracle *BlockOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockOracle.Contract.BlockOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockOracle *BlockOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockOracle *BlockOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockOracle *BlockOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockOracle.Contract.contract.Transact(opts, method, params...)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_BlockOracle *BlockOracleCaller) ISTEST(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "IS_TEST")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_BlockOracle *BlockOracleSession) ISTEST() (bool, error) {
	return _BlockOracle.Contract.ISTEST(&_BlockOracle.CallOpts)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_BlockOracle *BlockOracleCallerSession) ISTEST() (bool, error) {
	return _BlockOracle.Contract.ISTEST(&_BlockOracle.CallOpts)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_BlockOracle *BlockOracleCaller) ExcludeArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "excludeArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_BlockOracle *BlockOracleSession) ExcludeArtifacts() ([]string, error) {
	return _BlockOracle.Contract.ExcludeArtifacts(&_BlockOracle.CallOpts)
}

// ExcludeArtifacts is a free data retrieval call binding the contract method 0xb5508aa9.
//
// Solidity: function excludeArtifacts() view returns(string[] excludedArtifacts_)
func (_BlockOracle *BlockOracleCallerSession) ExcludeArtifacts() ([]string, error) {
	return _BlockOracle.Contract.ExcludeArtifacts(&_BlockOracle.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_BlockOracle *BlockOracleCaller) ExcludeContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "excludeContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_BlockOracle *BlockOracleSession) ExcludeContracts() ([]common.Address, error) {
	return _BlockOracle.Contract.ExcludeContracts(&_BlockOracle.CallOpts)
}

// ExcludeContracts is a free data retrieval call binding the contract method 0xe20c9f71.
//
// Solidity: function excludeContracts() view returns(address[] excludedContracts_)
func (_BlockOracle *BlockOracleCallerSession) ExcludeContracts() ([]common.Address, error) {
	return _BlockOracle.Contract.ExcludeContracts(&_BlockOracle.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_BlockOracle *BlockOracleCaller) ExcludeSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "excludeSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_BlockOracle *BlockOracleSession) ExcludeSenders() ([]common.Address, error) {
	return _BlockOracle.Contract.ExcludeSenders(&_BlockOracle.CallOpts)
}

// ExcludeSenders is a free data retrieval call binding the contract method 0x1ed7831c.
//
// Solidity: function excludeSenders() view returns(address[] excludedSenders_)
func (_BlockOracle *BlockOracleCallerSession) ExcludeSenders() ([]common.Address, error) {
	return _BlockOracle.Contract.ExcludeSenders(&_BlockOracle.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_BlockOracle *BlockOracleCaller) TargetArtifactSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "targetArtifactSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_BlockOracle *BlockOracleSession) TargetArtifactSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BlockOracle.Contract.TargetArtifactSelectors(&_BlockOracle.CallOpts)
}

// TargetArtifactSelectors is a free data retrieval call binding the contract method 0x66d9a9a0.
//
// Solidity: function targetArtifactSelectors() view returns((address,bytes4[])[] targetedArtifactSelectors_)
func (_BlockOracle *BlockOracleCallerSession) TargetArtifactSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BlockOracle.Contract.TargetArtifactSelectors(&_BlockOracle.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_BlockOracle *BlockOracleCaller) TargetArtifacts(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "targetArtifacts")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_BlockOracle *BlockOracleSession) TargetArtifacts() ([]string, error) {
	return _BlockOracle.Contract.TargetArtifacts(&_BlockOracle.CallOpts)
}

// TargetArtifacts is a free data retrieval call binding the contract method 0x85226c81.
//
// Solidity: function targetArtifacts() view returns(string[] targetedArtifacts_)
func (_BlockOracle *BlockOracleCallerSession) TargetArtifacts() ([]string, error) {
	return _BlockOracle.Contract.TargetArtifacts(&_BlockOracle.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_BlockOracle *BlockOracleCaller) TargetContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "targetContracts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_BlockOracle *BlockOracleSession) TargetContracts() ([]common.Address, error) {
	return _BlockOracle.Contract.TargetContracts(&_BlockOracle.CallOpts)
}

// TargetContracts is a free data retrieval call binding the contract method 0x3f7286f4.
//
// Solidity: function targetContracts() view returns(address[] targetedContracts_)
func (_BlockOracle *BlockOracleCallerSession) TargetContracts() ([]common.Address, error) {
	return _BlockOracle.Contract.TargetContracts(&_BlockOracle.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_BlockOracle *BlockOracleCaller) TargetSelectors(opts *bind.CallOpts) ([]StdInvariantFuzzSelector, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "targetSelectors")

	if err != nil {
		return *new([]StdInvariantFuzzSelector), err
	}

	out0 := *abi.ConvertType(out[0], new([]StdInvariantFuzzSelector)).(*[]StdInvariantFuzzSelector)

	return out0, err

}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_BlockOracle *BlockOracleSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BlockOracle.Contract.TargetSelectors(&_BlockOracle.CallOpts)
}

// TargetSelectors is a free data retrieval call binding the contract method 0x916a17c6.
//
// Solidity: function targetSelectors() view returns((address,bytes4[])[] targetedSelectors_)
func (_BlockOracle *BlockOracleCallerSession) TargetSelectors() ([]StdInvariantFuzzSelector, error) {
	return _BlockOracle.Contract.TargetSelectors(&_BlockOracle.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_BlockOracle *BlockOracleCaller) TargetSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BlockOracle.contract.Call(opts, &out, "targetSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_BlockOracle *BlockOracleSession) TargetSenders() ([]common.Address, error) {
	return _BlockOracle.Contract.TargetSenders(&_BlockOracle.CallOpts)
}

// TargetSenders is a free data retrieval call binding the contract method 0x3e5e3c23.
//
// Solidity: function targetSenders() view returns(address[] targetedSenders_)
func (_BlockOracle *BlockOracleCallerSession) TargetSenders() ([]common.Address, error) {
	return _BlockOracle.Contract.TargetSenders(&_BlockOracle.CallOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_BlockOracle *BlockOracleTransactor) Failed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockOracle.contract.Transact(opts, "failed")
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_BlockOracle *BlockOracleSession) Failed() (*types.Transaction, error) {
	return _BlockOracle.Contract.Failed(&_BlockOracle.TransactOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_BlockOracle *BlockOracleTransactorSession) Failed() (*types.Transaction, error) {
	return _BlockOracle.Contract.Failed(&_BlockOracle.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_BlockOracle *BlockOracleTransactor) SetUp(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockOracle.contract.Transact(opts, "setUp")
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_BlockOracle *BlockOracleSession) SetUp() (*types.Transaction, error) {
	return _BlockOracle.Contract.SetUp(&_BlockOracle.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_BlockOracle *BlockOracleTransactorSession) SetUp() (*types.Transaction, error) {
	return _BlockOracle.Contract.SetUp(&_BlockOracle.TransactOpts)
}

// TestCheckpointAndLoadSucceeds is a paid mutator transaction binding the contract method 0xf89a51bb.
//
// Solidity: function test_checkpointAndLoad_succeeds() returns()
func (_BlockOracle *BlockOracleTransactor) TestCheckpointAndLoadSucceeds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockOracle.contract.Transact(opts, "test_checkpointAndLoad_succeeds")
}

// TestCheckpointAndLoadSucceeds is a paid mutator transaction binding the contract method 0xf89a51bb.
//
// Solidity: function test_checkpointAndLoad_succeeds() returns()
func (_BlockOracle *BlockOracleSession) TestCheckpointAndLoadSucceeds() (*types.Transaction, error) {
	return _BlockOracle.Contract.TestCheckpointAndLoadSucceeds(&_BlockOracle.TransactOpts)
}

// TestCheckpointAndLoadSucceeds is a paid mutator transaction binding the contract method 0xf89a51bb.
//
// Solidity: function test_checkpointAndLoad_succeeds() returns()
func (_BlockOracle *BlockOracleTransactorSession) TestCheckpointAndLoadSucceeds() (*types.Transaction, error) {
	return _BlockOracle.Contract.TestCheckpointAndLoadSucceeds(&_BlockOracle.TransactOpts)
}

// TestLoadNoBlockHashReverts is a paid mutator transaction binding the contract method 0xcabc1a54.
//
// Solidity: function test_load_noBlockHash_reverts() returns()
func (_BlockOracle *BlockOracleTransactor) TestLoadNoBlockHashReverts(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockOracle.contract.Transact(opts, "test_load_noBlockHash_reverts")
}

// TestLoadNoBlockHashReverts is a paid mutator transaction binding the contract method 0xcabc1a54.
//
// Solidity: function test_load_noBlockHash_reverts() returns()
func (_BlockOracle *BlockOracleSession) TestLoadNoBlockHashReverts() (*types.Transaction, error) {
	return _BlockOracle.Contract.TestLoadNoBlockHashReverts(&_BlockOracle.TransactOpts)
}

// TestLoadNoBlockHashReverts is a paid mutator transaction binding the contract method 0xcabc1a54.
//
// Solidity: function test_load_noBlockHash_reverts() returns()
func (_BlockOracle *BlockOracleTransactorSession) TestLoadNoBlockHashReverts() (*types.Transaction, error) {
	return _BlockOracle.Contract.TestLoadNoBlockHashReverts(&_BlockOracle.TransactOpts)
}

// BlockOracleCheckpointIterator is returned from FilterCheckpoint and is used to iterate over the raw logs and unpacked data for Checkpoint events raised by the BlockOracle contract.
type BlockOracleCheckpointIterator struct {
	Event *BlockOracleCheckpoint // Event containing the contract specifics and raw log

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
func (it *BlockOracleCheckpointIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleCheckpoint)
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
		it.Event = new(BlockOracleCheckpoint)
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
func (it *BlockOracleCheckpointIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleCheckpointIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleCheckpoint represents a Checkpoint event raised by the BlockOracle contract.
type BlockOracleCheckpoint struct {
	BlockNumber    *big.Int
	BlockHash      [32]byte
	ChildTimestamp uint64
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCheckpoint is a free log retrieval operation binding the contract event 0xb67ff58b33060fd371a35ae2d9f1c3cdaec9b8197969f6efe2594a1ff4ba68c6.
//
// Solidity: event Checkpoint(uint256 indexed blockNumber, bytes32 indexed blockHash, uint64 indexed childTimestamp)
func (_BlockOracle *BlockOracleFilterer) FilterCheckpoint(opts *bind.FilterOpts, blockNumber []*big.Int, blockHash [][32]byte, childTimestamp []uint64) (*BlockOracleCheckpointIterator, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}
	var blockHashRule []interface{}
	for _, blockHashItem := range blockHash {
		blockHashRule = append(blockHashRule, blockHashItem)
	}
	var childTimestampRule []interface{}
	for _, childTimestampItem := range childTimestamp {
		childTimestampRule = append(childTimestampRule, childTimestampItem)
	}

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "Checkpoint", blockNumberRule, blockHashRule, childTimestampRule)
	if err != nil {
		return nil, err
	}
	return &BlockOracleCheckpointIterator{contract: _BlockOracle.contract, event: "Checkpoint", logs: logs, sub: sub}, nil
}

// WatchCheckpoint is a free log subscription operation binding the contract event 0xb67ff58b33060fd371a35ae2d9f1c3cdaec9b8197969f6efe2594a1ff4ba68c6.
//
// Solidity: event Checkpoint(uint256 indexed blockNumber, bytes32 indexed blockHash, uint64 indexed childTimestamp)
func (_BlockOracle *BlockOracleFilterer) WatchCheckpoint(opts *bind.WatchOpts, sink chan<- *BlockOracleCheckpoint, blockNumber []*big.Int, blockHash [][32]byte, childTimestamp []uint64) (event.Subscription, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}
	var blockHashRule []interface{}
	for _, blockHashItem := range blockHash {
		blockHashRule = append(blockHashRule, blockHashItem)
	}
	var childTimestampRule []interface{}
	for _, childTimestampItem := range childTimestamp {
		childTimestampRule = append(childTimestampRule, childTimestampItem)
	}

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "Checkpoint", blockNumberRule, blockHashRule, childTimestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleCheckpoint)
				if err := _BlockOracle.contract.UnpackLog(event, "Checkpoint", log); err != nil {
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

// ParseCheckpoint is a log parse operation binding the contract event 0xb67ff58b33060fd371a35ae2d9f1c3cdaec9b8197969f6efe2594a1ff4ba68c6.
//
// Solidity: event Checkpoint(uint256 indexed blockNumber, bytes32 indexed blockHash, uint64 indexed childTimestamp)
func (_BlockOracle *BlockOracleFilterer) ParseCheckpoint(log types.Log) (*BlockOracleCheckpoint, error) {
	event := new(BlockOracleCheckpoint)
	if err := _BlockOracle.contract.UnpackLog(event, "Checkpoint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogIterator is returned from FilterLog and is used to iterate over the raw logs and unpacked data for Log events raised by the BlockOracle contract.
type BlockOracleLogIterator struct {
	Event *BlockOracleLog // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLog)
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
		it.Event = new(BlockOracleLog)
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
func (it *BlockOracleLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLog represents a Log event raised by the BlockOracle contract.
type BlockOracleLog struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLog is a free log retrieval operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLog(opts *bind.FilterOpts) (*BlockOracleLogIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogIterator{contract: _BlockOracle.contract, event: "log", logs: logs, sub: sub}, nil
}

// WatchLog is a free log subscription operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLog(opts *bind.WatchOpts, sink chan<- *BlockOracleLog) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLog)
				if err := _BlockOracle.contract.UnpackLog(event, "log", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLog(log types.Log) (*BlockOracleLog, error) {
	event := new(BlockOracleLog)
	if err := _BlockOracle.contract.UnpackLog(event, "log", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogAddressIterator is returned from FilterLogAddress and is used to iterate over the raw logs and unpacked data for LogAddress events raised by the BlockOracle contract.
type BlockOracleLogAddressIterator struct {
	Event *BlockOracleLogAddress // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogAddress)
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
		it.Event = new(BlockOracleLogAddress)
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
func (it *BlockOracleLogAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogAddress represents a LogAddress event raised by the BlockOracle contract.
type BlockOracleLogAddress struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogAddress is a free log retrieval operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLogAddress(opts *bind.FilterOpts) (*BlockOracleLogAddressIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogAddressIterator{contract: _BlockOracle.contract, event: "log_address", logs: logs, sub: sub}, nil
}

// WatchLogAddress is a free log subscription operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLogAddress(opts *bind.WatchOpts, sink chan<- *BlockOracleLogAddress) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogAddress)
				if err := _BlockOracle.contract.UnpackLog(event, "log_address", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogAddress(log types.Log) (*BlockOracleLogAddress, error) {
	event := new(BlockOracleLogAddress)
	if err := _BlockOracle.contract.UnpackLog(event, "log_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogArrayIterator is returned from FilterLogArray and is used to iterate over the raw logs and unpacked data for LogArray events raised by the BlockOracle contract.
type BlockOracleLogArrayIterator struct {
	Event *BlockOracleLogArray // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogArray)
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
		it.Event = new(BlockOracleLogArray)
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
func (it *BlockOracleLogArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogArray represents a LogArray event raised by the BlockOracle contract.
type BlockOracleLogArray struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray is a free log retrieval operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_BlockOracle *BlockOracleFilterer) FilterLogArray(opts *bind.FilterOpts) (*BlockOracleLogArrayIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogArrayIterator{contract: _BlockOracle.contract, event: "log_array", logs: logs, sub: sub}, nil
}

// WatchLogArray is a free log subscription operation binding the contract event 0xfb102865d50addddf69da9b5aa1bced66c80cf869a5c8d0471a467e18ce9cab1.
//
// Solidity: event log_array(uint256[] val)
func (_BlockOracle *BlockOracleFilterer) WatchLogArray(opts *bind.WatchOpts, sink chan<- *BlockOracleLogArray) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogArray)
				if err := _BlockOracle.contract.UnpackLog(event, "log_array", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogArray(log types.Log) (*BlockOracleLogArray, error) {
	event := new(BlockOracleLogArray)
	if err := _BlockOracle.contract.UnpackLog(event, "log_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogArray0Iterator is returned from FilterLogArray0 and is used to iterate over the raw logs and unpacked data for LogArray0 events raised by the BlockOracle contract.
type BlockOracleLogArray0Iterator struct {
	Event *BlockOracleLogArray0 // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogArray0)
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
		it.Event = new(BlockOracleLogArray0)
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
func (it *BlockOracleLogArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogArray0 represents a LogArray0 event raised by the BlockOracle contract.
type BlockOracleLogArray0 struct {
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray0 is a free log retrieval operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_BlockOracle *BlockOracleFilterer) FilterLogArray0(opts *bind.FilterOpts) (*BlockOracleLogArray0Iterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogArray0Iterator{contract: _BlockOracle.contract, event: "log_array0", logs: logs, sub: sub}, nil
}

// WatchLogArray0 is a free log subscription operation binding the contract event 0x890a82679b470f2bd82816ed9b161f97d8b967f37fa3647c21d5bf39749e2dd5.
//
// Solidity: event log_array(int256[] val)
func (_BlockOracle *BlockOracleFilterer) WatchLogArray0(opts *bind.WatchOpts, sink chan<- *BlockOracleLogArray0) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogArray0)
				if err := _BlockOracle.contract.UnpackLog(event, "log_array0", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogArray0(log types.Log) (*BlockOracleLogArray0, error) {
	event := new(BlockOracleLogArray0)
	if err := _BlockOracle.contract.UnpackLog(event, "log_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogArray1Iterator is returned from FilterLogArray1 and is used to iterate over the raw logs and unpacked data for LogArray1 events raised by the BlockOracle contract.
type BlockOracleLogArray1Iterator struct {
	Event *BlockOracleLogArray1 // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogArray1)
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
		it.Event = new(BlockOracleLogArray1)
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
func (it *BlockOracleLogArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogArray1 represents a LogArray1 event raised by the BlockOracle contract.
type BlockOracleLogArray1 struct {
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogArray1 is a free log retrieval operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_BlockOracle *BlockOracleFilterer) FilterLogArray1(opts *bind.FilterOpts) (*BlockOracleLogArray1Iterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogArray1Iterator{contract: _BlockOracle.contract, event: "log_array1", logs: logs, sub: sub}, nil
}

// WatchLogArray1 is a free log subscription operation binding the contract event 0x40e1840f5769073d61bd01372d9b75baa9842d5629a0c99ff103be1178a8e9e2.
//
// Solidity: event log_array(address[] val)
func (_BlockOracle *BlockOracleFilterer) WatchLogArray1(opts *bind.WatchOpts, sink chan<- *BlockOracleLogArray1) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogArray1)
				if err := _BlockOracle.contract.UnpackLog(event, "log_array1", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogArray1(log types.Log) (*BlockOracleLogArray1, error) {
	event := new(BlockOracleLogArray1)
	if err := _BlockOracle.contract.UnpackLog(event, "log_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogBytesIterator is returned from FilterLogBytes and is used to iterate over the raw logs and unpacked data for LogBytes events raised by the BlockOracle contract.
type BlockOracleLogBytesIterator struct {
	Event *BlockOracleLogBytes // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogBytes)
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
		it.Event = new(BlockOracleLogBytes)
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
func (it *BlockOracleLogBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogBytes represents a LogBytes event raised by the BlockOracle contract.
type BlockOracleLogBytes struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes is a free log retrieval operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLogBytes(opts *bind.FilterOpts) (*BlockOracleLogBytesIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogBytesIterator{contract: _BlockOracle.contract, event: "log_bytes", logs: logs, sub: sub}, nil
}

// WatchLogBytes is a free log subscription operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLogBytes(opts *bind.WatchOpts, sink chan<- *BlockOracleLogBytes) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogBytes)
				if err := _BlockOracle.contract.UnpackLog(event, "log_bytes", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogBytes(log types.Log) (*BlockOracleLogBytes, error) {
	event := new(BlockOracleLogBytes)
	if err := _BlockOracle.contract.UnpackLog(event, "log_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogBytes32Iterator is returned from FilterLogBytes32 and is used to iterate over the raw logs and unpacked data for LogBytes32 events raised by the BlockOracle contract.
type BlockOracleLogBytes32Iterator struct {
	Event *BlockOracleLogBytes32 // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogBytes32)
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
		it.Event = new(BlockOracleLogBytes32)
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
func (it *BlockOracleLogBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogBytes32 represents a LogBytes32 event raised by the BlockOracle contract.
type BlockOracleLogBytes32 struct {
	Arg0 [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes32 is a free log retrieval operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLogBytes32(opts *bind.FilterOpts) (*BlockOracleLogBytes32Iterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogBytes32Iterator{contract: _BlockOracle.contract, event: "log_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogBytes32 is a free log subscription operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLogBytes32(opts *bind.WatchOpts, sink chan<- *BlockOracleLogBytes32) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogBytes32)
				if err := _BlockOracle.contract.UnpackLog(event, "log_bytes32", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogBytes32(log types.Log) (*BlockOracleLogBytes32, error) {
	event := new(BlockOracleLogBytes32)
	if err := _BlockOracle.contract.UnpackLog(event, "log_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogIntIterator is returned from FilterLogInt and is used to iterate over the raw logs and unpacked data for LogInt events raised by the BlockOracle contract.
type BlockOracleLogIntIterator struct {
	Event *BlockOracleLogInt // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogInt)
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
		it.Event = new(BlockOracleLogInt)
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
func (it *BlockOracleLogIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogInt represents a LogInt event raised by the BlockOracle contract.
type BlockOracleLogInt struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogInt is a free log retrieval operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLogInt(opts *bind.FilterOpts) (*BlockOracleLogIntIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogIntIterator{contract: _BlockOracle.contract, event: "log_int", logs: logs, sub: sub}, nil
}

// WatchLogInt is a free log subscription operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLogInt(opts *bind.WatchOpts, sink chan<- *BlockOracleLogInt) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogInt)
				if err := _BlockOracle.contract.UnpackLog(event, "log_int", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogInt(log types.Log) (*BlockOracleLogInt, error) {
	event := new(BlockOracleLogInt)
	if err := _BlockOracle.contract.UnpackLog(event, "log_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedAddressIterator is returned from FilterLogNamedAddress and is used to iterate over the raw logs and unpacked data for LogNamedAddress events raised by the BlockOracle contract.
type BlockOracleLogNamedAddressIterator struct {
	Event *BlockOracleLogNamedAddress // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedAddress)
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
		it.Event = new(BlockOracleLogNamedAddress)
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
func (it *BlockOracleLogNamedAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedAddress represents a LogNamedAddress event raised by the BlockOracle contract.
type BlockOracleLogNamedAddress struct {
	Key string
	Val common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedAddress is a free log retrieval operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedAddress(opts *bind.FilterOpts) (*BlockOracleLogNamedAddressIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedAddressIterator{contract: _BlockOracle.contract, event: "log_named_address", logs: logs, sub: sub}, nil
}

// WatchLogNamedAddress is a free log subscription operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedAddress(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedAddress) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedAddress)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_address", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedAddress(log types.Log) (*BlockOracleLogNamedAddress, error) {
	event := new(BlockOracleLogNamedAddress)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedArrayIterator is returned from FilterLogNamedArray and is used to iterate over the raw logs and unpacked data for LogNamedArray events raised by the BlockOracle contract.
type BlockOracleLogNamedArrayIterator struct {
	Event *BlockOracleLogNamedArray // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedArrayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedArray)
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
		it.Event = new(BlockOracleLogNamedArray)
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
func (it *BlockOracleLogNamedArrayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedArrayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedArray represents a LogNamedArray event raised by the BlockOracle contract.
type BlockOracleLogNamedArray struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray is a free log retrieval operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedArray(opts *bind.FilterOpts) (*BlockOracleLogNamedArrayIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedArrayIterator{contract: _BlockOracle.contract, event: "log_named_array", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray is a free log subscription operation binding the contract event 0x00aaa39c9ffb5f567a4534380c737075702e1f7f14107fc95328e3b56c0325fb.
//
// Solidity: event log_named_array(string key, uint256[] val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedArray(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedArray) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_array")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedArray)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_array", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedArray(log types.Log) (*BlockOracleLogNamedArray, error) {
	event := new(BlockOracleLogNamedArray)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_array", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedArray0Iterator is returned from FilterLogNamedArray0 and is used to iterate over the raw logs and unpacked data for LogNamedArray0 events raised by the BlockOracle contract.
type BlockOracleLogNamedArray0Iterator struct {
	Event *BlockOracleLogNamedArray0 // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedArray0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedArray0)
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
		it.Event = new(BlockOracleLogNamedArray0)
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
func (it *BlockOracleLogNamedArray0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedArray0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedArray0 represents a LogNamedArray0 event raised by the BlockOracle contract.
type BlockOracleLogNamedArray0 struct {
	Key string
	Val []*big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray0 is a free log retrieval operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedArray0(opts *bind.FilterOpts) (*BlockOracleLogNamedArray0Iterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedArray0Iterator{contract: _BlockOracle.contract, event: "log_named_array0", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray0 is a free log subscription operation binding the contract event 0xa73eda09662f46dde729be4611385ff34fe6c44fbbc6f7e17b042b59a3445b57.
//
// Solidity: event log_named_array(string key, int256[] val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedArray0(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedArray0) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_array0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedArray0)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_array0", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedArray0(log types.Log) (*BlockOracleLogNamedArray0, error) {
	event := new(BlockOracleLogNamedArray0)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_array0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedArray1Iterator is returned from FilterLogNamedArray1 and is used to iterate over the raw logs and unpacked data for LogNamedArray1 events raised by the BlockOracle contract.
type BlockOracleLogNamedArray1Iterator struct {
	Event *BlockOracleLogNamedArray1 // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedArray1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedArray1)
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
		it.Event = new(BlockOracleLogNamedArray1)
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
func (it *BlockOracleLogNamedArray1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedArray1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedArray1 represents a LogNamedArray1 event raised by the BlockOracle contract.
type BlockOracleLogNamedArray1 struct {
	Key string
	Val []common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedArray1 is a free log retrieval operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedArray1(opts *bind.FilterOpts) (*BlockOracleLogNamedArray1Iterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedArray1Iterator{contract: _BlockOracle.contract, event: "log_named_array1", logs: logs, sub: sub}, nil
}

// WatchLogNamedArray1 is a free log subscription operation binding the contract event 0x3bcfb2ae2e8d132dd1fce7cf278a9a19756a9fceabe470df3bdabb4bc577d1bd.
//
// Solidity: event log_named_array(string key, address[] val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedArray1(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedArray1) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_array1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedArray1)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_array1", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedArray1(log types.Log) (*BlockOracleLogNamedArray1, error) {
	event := new(BlockOracleLogNamedArray1)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_array1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedBytesIterator is returned from FilterLogNamedBytes and is used to iterate over the raw logs and unpacked data for LogNamedBytes events raised by the BlockOracle contract.
type BlockOracleLogNamedBytesIterator struct {
	Event *BlockOracleLogNamedBytes // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedBytes)
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
		it.Event = new(BlockOracleLogNamedBytes)
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
func (it *BlockOracleLogNamedBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedBytes represents a LogNamedBytes event raised by the BlockOracle contract.
type BlockOracleLogNamedBytes struct {
	Key string
	Val []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes is a free log retrieval operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedBytes(opts *bind.FilterOpts) (*BlockOracleLogNamedBytesIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedBytesIterator{contract: _BlockOracle.contract, event: "log_named_bytes", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes is a free log subscription operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedBytes(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedBytes) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedBytes)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedBytes(log types.Log) (*BlockOracleLogNamedBytes, error) {
	event := new(BlockOracleLogNamedBytes)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedBytes32Iterator is returned from FilterLogNamedBytes32 and is used to iterate over the raw logs and unpacked data for LogNamedBytes32 events raised by the BlockOracle contract.
type BlockOracleLogNamedBytes32Iterator struct {
	Event *BlockOracleLogNamedBytes32 // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedBytes32)
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
		it.Event = new(BlockOracleLogNamedBytes32)
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
func (it *BlockOracleLogNamedBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedBytes32 represents a LogNamedBytes32 event raised by the BlockOracle contract.
type BlockOracleLogNamedBytes32 struct {
	Key string
	Val [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes32 is a free log retrieval operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedBytes32(opts *bind.FilterOpts) (*BlockOracleLogNamedBytes32Iterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedBytes32Iterator{contract: _BlockOracle.contract, event: "log_named_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes32 is a free log subscription operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedBytes32(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedBytes32) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedBytes32)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedBytes32(log types.Log) (*BlockOracleLogNamedBytes32, error) {
	event := new(BlockOracleLogNamedBytes32)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedDecimalIntIterator is returned from FilterLogNamedDecimalInt and is used to iterate over the raw logs and unpacked data for LogNamedDecimalInt events raised by the BlockOracle contract.
type BlockOracleLogNamedDecimalIntIterator struct {
	Event *BlockOracleLogNamedDecimalInt // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedDecimalIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedDecimalInt)
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
		it.Event = new(BlockOracleLogNamedDecimalInt)
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
func (it *BlockOracleLogNamedDecimalIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedDecimalIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedDecimalInt represents a LogNamedDecimalInt event raised by the BlockOracle contract.
type BlockOracleLogNamedDecimalInt struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalInt is a free log retrieval operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedDecimalInt(opts *bind.FilterOpts) (*BlockOracleLogNamedDecimalIntIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedDecimalIntIterator{contract: _BlockOracle.contract, event: "log_named_decimal_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalInt is a free log subscription operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedDecimalInt(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedDecimalInt) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedDecimalInt)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedDecimalInt(log types.Log) (*BlockOracleLogNamedDecimalInt, error) {
	event := new(BlockOracleLogNamedDecimalInt)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedDecimalUintIterator is returned from FilterLogNamedDecimalUint and is used to iterate over the raw logs and unpacked data for LogNamedDecimalUint events raised by the BlockOracle contract.
type BlockOracleLogNamedDecimalUintIterator struct {
	Event *BlockOracleLogNamedDecimalUint // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedDecimalUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedDecimalUint)
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
		it.Event = new(BlockOracleLogNamedDecimalUint)
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
func (it *BlockOracleLogNamedDecimalUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedDecimalUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedDecimalUint represents a LogNamedDecimalUint event raised by the BlockOracle contract.
type BlockOracleLogNamedDecimalUint struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalUint is a free log retrieval operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedDecimalUint(opts *bind.FilterOpts) (*BlockOracleLogNamedDecimalUintIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedDecimalUintIterator{contract: _BlockOracle.contract, event: "log_named_decimal_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalUint is a free log subscription operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedDecimalUint(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedDecimalUint) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedDecimalUint)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedDecimalUint(log types.Log) (*BlockOracleLogNamedDecimalUint, error) {
	event := new(BlockOracleLogNamedDecimalUint)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedIntIterator is returned from FilterLogNamedInt and is used to iterate over the raw logs and unpacked data for LogNamedInt events raised by the BlockOracle contract.
type BlockOracleLogNamedIntIterator struct {
	Event *BlockOracleLogNamedInt // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedInt)
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
		it.Event = new(BlockOracleLogNamedInt)
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
func (it *BlockOracleLogNamedIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedInt represents a LogNamedInt event raised by the BlockOracle contract.
type BlockOracleLogNamedInt struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedInt is a free log retrieval operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedInt(opts *bind.FilterOpts) (*BlockOracleLogNamedIntIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedIntIterator{contract: _BlockOracle.contract, event: "log_named_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedInt is a free log subscription operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedInt(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedInt) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedInt)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_int", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedInt(log types.Log) (*BlockOracleLogNamedInt, error) {
	event := new(BlockOracleLogNamedInt)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedStringIterator is returned from FilterLogNamedString and is used to iterate over the raw logs and unpacked data for LogNamedString events raised by the BlockOracle contract.
type BlockOracleLogNamedStringIterator struct {
	Event *BlockOracleLogNamedString // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedString)
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
		it.Event = new(BlockOracleLogNamedString)
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
func (it *BlockOracleLogNamedStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedString represents a LogNamedString event raised by the BlockOracle contract.
type BlockOracleLogNamedString struct {
	Key string
	Val string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedString is a free log retrieval operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedString(opts *bind.FilterOpts) (*BlockOracleLogNamedStringIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedStringIterator{contract: _BlockOracle.contract, event: "log_named_string", logs: logs, sub: sub}, nil
}

// WatchLogNamedString is a free log subscription operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedString(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedString) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedString)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_string", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedString(log types.Log) (*BlockOracleLogNamedString, error) {
	event := new(BlockOracleLogNamedString)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogNamedUintIterator is returned from FilterLogNamedUint and is used to iterate over the raw logs and unpacked data for LogNamedUint events raised by the BlockOracle contract.
type BlockOracleLogNamedUintIterator struct {
	Event *BlockOracleLogNamedUint // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogNamedUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogNamedUint)
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
		it.Event = new(BlockOracleLogNamedUint)
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
func (it *BlockOracleLogNamedUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogNamedUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogNamedUint represents a LogNamedUint event raised by the BlockOracle contract.
type BlockOracleLogNamedUint struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedUint is a free log retrieval operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_BlockOracle *BlockOracleFilterer) FilterLogNamedUint(opts *bind.FilterOpts) (*BlockOracleLogNamedUintIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogNamedUintIterator{contract: _BlockOracle.contract, event: "log_named_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedUint is a free log subscription operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_BlockOracle *BlockOracleFilterer) WatchLogNamedUint(opts *bind.WatchOpts, sink chan<- *BlockOracleLogNamedUint) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogNamedUint)
				if err := _BlockOracle.contract.UnpackLog(event, "log_named_uint", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogNamedUint(log types.Log) (*BlockOracleLogNamedUint, error) {
	event := new(BlockOracleLogNamedUint)
	if err := _BlockOracle.contract.UnpackLog(event, "log_named_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogStringIterator is returned from FilterLogString and is used to iterate over the raw logs and unpacked data for LogString events raised by the BlockOracle contract.
type BlockOracleLogStringIterator struct {
	Event *BlockOracleLogString // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogString)
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
		it.Event = new(BlockOracleLogString)
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
func (it *BlockOracleLogStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogString represents a LogString event raised by the BlockOracle contract.
type BlockOracleLogString struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogString is a free log retrieval operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLogString(opts *bind.FilterOpts) (*BlockOracleLogStringIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogStringIterator{contract: _BlockOracle.contract, event: "log_string", logs: logs, sub: sub}, nil
}

// WatchLogString is a free log subscription operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLogString(opts *bind.WatchOpts, sink chan<- *BlockOracleLogString) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogString)
				if err := _BlockOracle.contract.UnpackLog(event, "log_string", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogString(log types.Log) (*BlockOracleLogString, error) {
	event := new(BlockOracleLogString)
	if err := _BlockOracle.contract.UnpackLog(event, "log_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogUintIterator is returned from FilterLogUint and is used to iterate over the raw logs and unpacked data for LogUint events raised by the BlockOracle contract.
type BlockOracleLogUintIterator struct {
	Event *BlockOracleLogUint // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogUint)
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
		it.Event = new(BlockOracleLogUint)
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
func (it *BlockOracleLogUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogUint represents a LogUint event raised by the BlockOracle contract.
type BlockOracleLogUint struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogUint is a free log retrieval operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLogUint(opts *bind.FilterOpts) (*BlockOracleLogUintIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogUintIterator{contract: _BlockOracle.contract, event: "log_uint", logs: logs, sub: sub}, nil
}

// WatchLogUint is a free log subscription operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLogUint(opts *bind.WatchOpts, sink chan<- *BlockOracleLogUint) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogUint)
				if err := _BlockOracle.contract.UnpackLog(event, "log_uint", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogUint(log types.Log) (*BlockOracleLogUint, error) {
	event := new(BlockOracleLogUint)
	if err := _BlockOracle.contract.UnpackLog(event, "log_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockOracleLogsIterator is returned from FilterLogs and is used to iterate over the raw logs and unpacked data for Logs events raised by the BlockOracle contract.
type BlockOracleLogsIterator struct {
	Event *BlockOracleLogs // Event containing the contract specifics and raw log

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
func (it *BlockOracleLogsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockOracleLogs)
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
		it.Event = new(BlockOracleLogs)
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
func (it *BlockOracleLogsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockOracleLogsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockOracleLogs represents a Logs event raised by the BlockOracle contract.
type BlockOracleLogs struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogs is a free log retrieval operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_BlockOracle *BlockOracleFilterer) FilterLogs(opts *bind.FilterOpts) (*BlockOracleLogsIterator, error) {

	logs, sub, err := _BlockOracle.contract.FilterLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return &BlockOracleLogsIterator{contract: _BlockOracle.contract, event: "logs", logs: logs, sub: sub}, nil
}

// WatchLogs is a free log subscription operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_BlockOracle *BlockOracleFilterer) WatchLogs(opts *bind.WatchOpts, sink chan<- *BlockOracleLogs) (event.Subscription, error) {

	logs, sub, err := _BlockOracle.contract.WatchLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockOracleLogs)
				if err := _BlockOracle.contract.UnpackLog(event, "logs", log); err != nil {
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
func (_BlockOracle *BlockOracleFilterer) ParseLogs(log types.Log) (*BlockOracleLogs, error) {
	event := new(BlockOracleLogs)
	if err := _BlockOracle.contract.UnpackLog(event, "logs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
