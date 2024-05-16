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

// L2NativeTokenMetaData contains all meta data concerning the L2NativeToken contract.
var L2NativeTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approveAndCall\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"callbackEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enableCallback\",\"inputs\":[{\"name\":\"_callbackEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"faucet\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceMinter\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renouncePauser\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"seigManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractSeigManagerI\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setSeigManager\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractSeigManagerI\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60806040523480156200001157600080fd5b506040518060400160405280601581526020017f546f6b616d616b204e6574776f726b20546f6b656e0000000000000000000000815250604051806040016040528060038152602001622a27a760e91b81525060126000620000786200010160201b60201c565b600380546001600160a01b0319166001600160a01b038316908117909155604051919250906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3506004620000d48482620001aa565b506005620000e38382620001aa565b506006805460ff191660ff9290921691909117905550620002769050565b3390565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200013057607f821691505b6020821081036200015157634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620001a557600081815260208120601f850160051c81016020861015620001805750805b601f850160051c820191505b81811015620001a1578281556001016200018c565b5050505b505050565b81516001600160401b03811115620001c657620001c662000105565b620001de81620001d784546200011b565b8462000157565b602080601f831160018114620002165760008415620001fd5750858301515b600019600386901b1c1916600185901b178555620001a1565b600085815260208120601f198616915b82811015620002475788860151825594840194600190910190840162000226565b5085821015620002665787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b611f8380620002866000396000f3fe608060405234801561001057600080fd5b50600436106101a35760003560e01c80636d435421116100ee5780638f32d59b11610097578063a9059cbb11610071578063a9059cbb146103c4578063cae9ca51146103d7578063dd62ed3e146103ea578063f2fde38b1461043057600080fd5b80638f32d59b1461038957806395d89b41146103a9578063a457c2d7146103b157600080fd5b8063715018a6116100c8578063715018a6146103505780637657f20a146103585780638da5cb5b1461036b57600080fd5b80636d435421146102bd5780636fb7f558146102d057806370a082311461031a57600080fd5b806338bf3cfa11610150578063579158971161012a57806357915897146102715780635f112c6814610284578063633801131461029757600080fd5b806338bf3cfa14610238578063395093511461024b57806341eb24bb1461025e57600080fd5b806323b872dd1161018157806323b872dd146101fb5780633113ed5c1461020e578063313ce5671461022357600080fd5b806306fdde03146101a8578063095ea7b3146101c657806318160ddd146101e9575b600080fd5b6101b0610443565b6040516101bd9190611b7a565b60405180910390f35b6101d96101d4366004611baf565b6104d5565b60405190151581526020016101bd565b6002545b6040519081526020016101bd565b6101d9610209366004611bdb565b6104eb565b61022161021c366004611c2a565b6105ca565b005b60065460405160ff90911681526020016101bd565b610221610246366004611c47565b610696565b6101d9610259366004611baf565b61077a565b61022161026c366004611c47565b6107c2565b61022161027f366004611c64565b61088b565b610221610292366004611c47565b610898565b6006546101d9907501000000000000000000000000000000000000000000900460ff1681565b6102216102cb366004611c7d565b610961565b6006546102f590610100900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101bd565b6101ed610328366004611c47565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b610221610a67565b610221610366366004611c47565b610b57565b60035473ffffffffffffffffffffffffffffffffffffffff166102f5565b60035473ffffffffffffffffffffffffffffffffffffffff1633146101d9565b6101b0610bdf565b6101d96103bf366004611baf565b610bee565b6101d96103d2366004611baf565b610c4a565b6101d96103e5366004611ce5565b610c57565b6101ed6103f8366004611c7d565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b61022161043e366004611c47565b610c82565b60606004805461045290611dd0565b80601f016020809104026020016040519081016040528092919081815260200182805461047e90611dd0565b80156104cb5780601f106104a0576101008083540402835291602001916104cb565b820191906000526020600020905b8154815290600101906020018083116104ae57829003601f168201915b5050505050905090565b60006104e2338484610d0c565b50600192915050565b60003373ffffffffffffffffffffffffffffffffffffffff8516148061052657503373ffffffffffffffffffffffffffffffffffffffff8416145b6105b7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603060248201527f53656967546f6b656e3a206f6e6c792073656e646572206f722072656369706960448201527f656e742063616e207472616e736665720000000000000000000000000000000060648201526084015b60405180910390fd5b6105c2848484610ec0565b949350505050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461064b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ae565b600680549115157501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909216919091179055565b60035473ffffffffffffffffffffffffffffffffffffffff163314610717576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ae565b8073ffffffffffffffffffffffffffffffffffffffff1663715018a66040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561075f57600080fd5b505af1158015610773573d6000803e3d6000fd5b5050505050565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490916104e29185906107bd9086610f27565b610d0c565b60035473ffffffffffffffffffffffffffffffffffffffff163314610843576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ae565b8073ffffffffffffffffffffffffffffffffffffffff16636ef8d66d6040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561075f57600080fd5b6108953382610fa7565b50565b60035473ffffffffffffffffffffffffffffffffffffffff163314610919576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ae565b8073ffffffffffffffffffffffffffffffffffffffff1663986502756040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561075f57600080fd5b60035473ffffffffffffffffffffffffffffffffffffffff1633146109e2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ae565b6040517ff2fde38b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff828116600483015283169063f2fde38b90602401600060405180830381600087803b158015610a4b57600080fd5b505af1158015610a5f573d6000803e3d6000fd5b505050505050565b60035473ffffffffffffffffffffffffffffffffffffffff163314610ae8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ae565b60035460405160009173ffffffffffffffffffffffffffffffffffffffff16907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600380547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f544f4e3a20544f4e20646f65736e277420616c6c6f7720736574536569674d6160448201527f6e6167657200000000000000000000000000000000000000000000000000000060648201526084016105ae565b60606005805461045290611dd0565b60006104e233846107bd85604051806060016040528060258152602001611f526025913933600090815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8d16845290915290205491906110b4565b60006104e2338484611108565b6000610c6384846104d5565b610c6c57600080fd5b610c7833858585611218565b5060019392505050565b60035473ffffffffffffffffffffffffffffffffffffffff163314610d03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ae565b610895816114a9565b73ffffffffffffffffffffffffffffffffffffffff8316610dae576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f726573730000000000000000000000000000000000000000000000000000000060648201526084016105ae565b73ffffffffffffffffffffffffffffffffffffffff8216610e51576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f737300000000000000000000000000000000000000000000000000000000000060648201526084016105ae565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6000610ecd848484611108565b610c7884336107bd85604051806060016040528060288152602001611f2a6028913973ffffffffffffffffffffffffffffffffffffffff8a16600090815260016020908152604080832033845290915290205491906110b4565b600080610f348385611e52565b905083811015610fa0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f77000000000060448201526064016105ae565b9392505050565b610fb182826115da565b6006547501000000000000000000000000000000000000000000900460ff168015610ff85750600654610100900473ffffffffffffffffffffffffffffffffffffffff1615155b156110b0576006546040517f4a3931490000000000000000000000000000000000000000000000000000000081526000600482015273ffffffffffffffffffffffffffffffffffffffff84811660248301526044820184905261010090920490911690634a393149906064016020604051808303816000875af1158015611083573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110a79190611e6a565b6110b057600080fd5b5050565b600081848411156110f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105ae9190611b7a565b5060006110ff8486611e87565b95945050505050565b6111138383836116f9565b6006547501000000000000000000000000000000000000000000900460ff16801561115a5750600654610100900473ffffffffffffffffffffffffffffffffffffffff1615155b15611213576006546040517f4a39314900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116600483015284811660248301526044820184905261010090920490911690634a393149906064016020604051808303816000875af11580156111e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061120a9190611e6a565b61121357600080fd5b505050565b7f4273ca16000000000000000000000000000000000000000000000000000000006112438482611923565b6112cf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603160248201527f45524332304f6e417070726f76653a207370656e64657220646f65736e27742060448201527f737570706f7274206f6e417070726f766500000000000000000000000000000060648201526084016105ae565b6000808573ffffffffffffffffffffffffffffffffffffffff1683888888886040516024016113019493929190611e9e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090941693909317909252905161138a9190611ee7565b6000604051808303816000865af19150503d80600081146113c7576040519150601f19603f3d011682016040523d82523d6000602084013e6113cc565b606091505b509150915081819061140b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105ae9190611b7a565b5060208101519150816114a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f45524332304f6e417070726f76653a206661696c656420746f2063616c6c206f60448201527f6e417070726f766500000000000000000000000000000000000000000000000060648201526084016105ae565b50505050505050565b73ffffffffffffffffffffffffffffffffffffffff811661154c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016105ae565b60035460405173ffffffffffffffffffffffffffffffffffffffff8084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b73ffffffffffffffffffffffffffffffffffffffff8216611657576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f20616464726573730060448201526064016105ae565b6002546116649082610f27565b60025573ffffffffffffffffffffffffffffffffffffffff82166000908152602081905260409020546116979082610f27565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260208181526040808320949094559251848152919290917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b73ffffffffffffffffffffffffffffffffffffffff831661179c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f647265737300000000000000000000000000000000000000000000000000000060648201526084016105ae565b73ffffffffffffffffffffffffffffffffffffffff821661183f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f657373000000000000000000000000000000000000000000000000000000000060648201526084016105ae565b61188981604051806060016040528060268152602001611f046026913973ffffffffffffffffffffffffffffffffffffffff861660009081526020819052604090205491906110b4565b73ffffffffffffffffffffffffffffffffffffffff80851660009081526020819052604080822093909355908416815220546118c59082610f27565b73ffffffffffffffffffffffffffffffffffffffff8381166000818152602081815260409182902094909455518481529092918616917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9101610eb3565b600061192e8361193f565b8015610fa05750610fa083836119a4565b600061196b827f01ffc9a7000000000000000000000000000000000000000000000000000000006119a4565b801561199e575061199c827fffffffff000000000000000000000000000000000000000000000000000000006119a4565b155b92915050565b60008060006119b385856119c7565b915091508180156110ff5750949350505050565b604080517fffffffff00000000000000000000000000000000000000000000000000000000831660248083019190915282518083039091018152604490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a700000000000000000000000000000000000000000000000000000000179052905160009182918290819073ffffffffffffffffffffffffffffffffffffffff881690611a7e908590611ee7565b600060405180830381855afa9150503d8060008114611ab9576040519150601f19603f3d011682016040523d82523d6000602084013e611abe565b606091505b5091509150602081511015611adc5760008094509450505050611af9565b8181806020019051810190611af19190611e6a565b945094505050505b9250929050565b60005b83811015611b1b578181015183820152602001611b03565b83811115611b2a576000848401525b50505050565b60008151808452611b48816020860160208601611b00565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610fa06020830184611b30565b73ffffffffffffffffffffffffffffffffffffffff8116811461089557600080fd5b60008060408385031215611bc257600080fd5b8235611bcd81611b8d565b946020939093013593505050565b600080600060608486031215611bf057600080fd5b8335611bfb81611b8d565b92506020840135611c0b81611b8d565b929592945050506040919091013590565b801515811461089557600080fd5b600060208284031215611c3c57600080fd5b8135610fa081611c1c565b600060208284031215611c5957600080fd5b8135610fa081611b8d565b600060208284031215611c7657600080fd5b5035919050565b60008060408385031215611c9057600080fd5b8235611c9b81611b8d565b91506020830135611cab81611b8d565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600060608486031215611cfa57600080fd5b8335611d0581611b8d565b925060208401359150604084013567ffffffffffffffff80821115611d2957600080fd5b818601915086601f830112611d3d57600080fd5b813581811115611d4f57611d4f611cb6565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715611d9557611d95611cb6565b81604052828152896020848701011115611dae57600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b600181811c90821680611de457607f821691505b602082108103611e1d577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115611e6557611e65611e23565b500190565b600060208284031215611e7c57600080fd5b8151610fa081611c1c565b600082821015611e9957611e99611e23565b500390565b600073ffffffffffffffffffffffffffffffffffffffff808716835280861660208401525083604083015260806060830152611edd6080830184611b30565b9695505050505050565b60008251611ef9818460208701611b00565b919091019291505056fe45524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e636545524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636545524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa164736f6c634300080f000a",
}

// L2NativeTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use L2NativeTokenMetaData.ABI instead.
var L2NativeTokenABI = L2NativeTokenMetaData.ABI

// L2NativeTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L2NativeTokenMetaData.Bin instead.
var L2NativeTokenBin = L2NativeTokenMetaData.Bin

// DeployL2NativeToken deploys a new Ethereum contract, binding an instance of L2NativeToken to it.
func DeployL2NativeToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *L2NativeToken, error) {
	parsed, err := L2NativeTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L2NativeTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L2NativeToken{L2NativeTokenCaller: L2NativeTokenCaller{contract: contract}, L2NativeTokenTransactor: L2NativeTokenTransactor{contract: contract}, L2NativeTokenFilterer: L2NativeTokenFilterer{contract: contract}}, nil
}

// L2NativeToken is an auto generated Go binding around an Ethereum contract.
type L2NativeToken struct {
	L2NativeTokenCaller     // Read-only binding to the contract
	L2NativeTokenTransactor // Write-only binding to the contract
	L2NativeTokenFilterer   // Log filterer for contract events
}

// L2NativeTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2NativeTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2NativeTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2NativeTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2NativeTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2NativeTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2NativeTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2NativeTokenSession struct {
	Contract     *L2NativeToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L2NativeTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2NativeTokenCallerSession struct {
	Contract *L2NativeTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// L2NativeTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2NativeTokenTransactorSession struct {
	Contract     *L2NativeTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// L2NativeTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2NativeTokenRaw struct {
	Contract *L2NativeToken // Generic contract binding to access the raw methods on
}

// L2NativeTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2NativeTokenCallerRaw struct {
	Contract *L2NativeTokenCaller // Generic read-only contract binding to access the raw methods on
}

// L2NativeTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2NativeTokenTransactorRaw struct {
	Contract *L2NativeTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2NativeToken creates a new instance of L2NativeToken, bound to a specific deployed contract.
func NewL2NativeToken(address common.Address, backend bind.ContractBackend) (*L2NativeToken, error) {
	contract, err := bindL2NativeToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2NativeToken{L2NativeTokenCaller: L2NativeTokenCaller{contract: contract}, L2NativeTokenTransactor: L2NativeTokenTransactor{contract: contract}, L2NativeTokenFilterer: L2NativeTokenFilterer{contract: contract}}, nil
}

// NewL2NativeTokenCaller creates a new read-only instance of L2NativeToken, bound to a specific deployed contract.
func NewL2NativeTokenCaller(address common.Address, caller bind.ContractCaller) (*L2NativeTokenCaller, error) {
	contract, err := bindL2NativeToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2NativeTokenCaller{contract: contract}, nil
}

// NewL2NativeTokenTransactor creates a new write-only instance of L2NativeToken, bound to a specific deployed contract.
func NewL2NativeTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*L2NativeTokenTransactor, error) {
	contract, err := bindL2NativeToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2NativeTokenTransactor{contract: contract}, nil
}

// NewL2NativeTokenFilterer creates a new log filterer instance of L2NativeToken, bound to a specific deployed contract.
func NewL2NativeTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*L2NativeTokenFilterer, error) {
	contract, err := bindL2NativeToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2NativeTokenFilterer{contract: contract}, nil
}

// bindL2NativeToken binds a generic wrapper to an already deployed contract.
func bindL2NativeToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(L2NativeTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2NativeToken *L2NativeTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2NativeToken.Contract.L2NativeTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2NativeToken *L2NativeTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2NativeToken.Contract.L2NativeTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2NativeToken *L2NativeTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2NativeToken.Contract.L2NativeTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2NativeToken *L2NativeTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2NativeToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2NativeToken *L2NativeTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2NativeToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2NativeToken *L2NativeTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2NativeToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_L2NativeToken *L2NativeTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_L2NativeToken *L2NativeTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _L2NativeToken.Contract.Allowance(&_L2NativeToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_L2NativeToken *L2NativeTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _L2NativeToken.Contract.Allowance(&_L2NativeToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_L2NativeToken *L2NativeTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_L2NativeToken *L2NativeTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _L2NativeToken.Contract.BalanceOf(&_L2NativeToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_L2NativeToken *L2NativeTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _L2NativeToken.Contract.BalanceOf(&_L2NativeToken.CallOpts, account)
}

// CallbackEnabled is a free data retrieval call binding the contract method 0x63380113.
//
// Solidity: function callbackEnabled() view returns(bool)
func (_L2NativeToken *L2NativeTokenCaller) CallbackEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "callbackEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CallbackEnabled is a free data retrieval call binding the contract method 0x63380113.
//
// Solidity: function callbackEnabled() view returns(bool)
func (_L2NativeToken *L2NativeTokenSession) CallbackEnabled() (bool, error) {
	return _L2NativeToken.Contract.CallbackEnabled(&_L2NativeToken.CallOpts)
}

// CallbackEnabled is a free data retrieval call binding the contract method 0x63380113.
//
// Solidity: function callbackEnabled() view returns(bool)
func (_L2NativeToken *L2NativeTokenCallerSession) CallbackEnabled() (bool, error) {
	return _L2NativeToken.Contract.CallbackEnabled(&_L2NativeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_L2NativeToken *L2NativeTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_L2NativeToken *L2NativeTokenSession) Decimals() (uint8, error) {
	return _L2NativeToken.Contract.Decimals(&_L2NativeToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_L2NativeToken *L2NativeTokenCallerSession) Decimals() (uint8, error) {
	return _L2NativeToken.Contract.Decimals(&_L2NativeToken.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_L2NativeToken *L2NativeTokenCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_L2NativeToken *L2NativeTokenSession) IsOwner() (bool, error) {
	return _L2NativeToken.Contract.IsOwner(&_L2NativeToken.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_L2NativeToken *L2NativeTokenCallerSession) IsOwner() (bool, error) {
	return _L2NativeToken.Contract.IsOwner(&_L2NativeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_L2NativeToken *L2NativeTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_L2NativeToken *L2NativeTokenSession) Name() (string, error) {
	return _L2NativeToken.Contract.Name(&_L2NativeToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_L2NativeToken *L2NativeTokenCallerSession) Name() (string, error) {
	return _L2NativeToken.Contract.Name(&_L2NativeToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2NativeToken *L2NativeTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2NativeToken *L2NativeTokenSession) Owner() (common.Address, error) {
	return _L2NativeToken.Contract.Owner(&_L2NativeToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2NativeToken *L2NativeTokenCallerSession) Owner() (common.Address, error) {
	return _L2NativeToken.Contract.Owner(&_L2NativeToken.CallOpts)
}

// SeigManager is a free data retrieval call binding the contract method 0x6fb7f558.
//
// Solidity: function seigManager() view returns(address)
func (_L2NativeToken *L2NativeTokenCaller) SeigManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "seigManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SeigManager is a free data retrieval call binding the contract method 0x6fb7f558.
//
// Solidity: function seigManager() view returns(address)
func (_L2NativeToken *L2NativeTokenSession) SeigManager() (common.Address, error) {
	return _L2NativeToken.Contract.SeigManager(&_L2NativeToken.CallOpts)
}

// SeigManager is a free data retrieval call binding the contract method 0x6fb7f558.
//
// Solidity: function seigManager() view returns(address)
func (_L2NativeToken *L2NativeTokenCallerSession) SeigManager() (common.Address, error) {
	return _L2NativeToken.Contract.SeigManager(&_L2NativeToken.CallOpts)
}

// SetSeigManager is a free data retrieval call binding the contract method 0x7657f20a.
//
// Solidity: function setSeigManager(address ) pure returns()
func (_L2NativeToken *L2NativeTokenCaller) SetSeigManager(opts *bind.CallOpts, arg0 common.Address) error {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "setSeigManager", arg0)

	if err != nil {
		return err
	}

	return err

}

// SetSeigManager is a free data retrieval call binding the contract method 0x7657f20a.
//
// Solidity: function setSeigManager(address ) pure returns()
func (_L2NativeToken *L2NativeTokenSession) SetSeigManager(arg0 common.Address) error {
	return _L2NativeToken.Contract.SetSeigManager(&_L2NativeToken.CallOpts, arg0)
}

// SetSeigManager is a free data retrieval call binding the contract method 0x7657f20a.
//
// Solidity: function setSeigManager(address ) pure returns()
func (_L2NativeToken *L2NativeTokenCallerSession) SetSeigManager(arg0 common.Address) error {
	return _L2NativeToken.Contract.SetSeigManager(&_L2NativeToken.CallOpts, arg0)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_L2NativeToken *L2NativeTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_L2NativeToken *L2NativeTokenSession) Symbol() (string, error) {
	return _L2NativeToken.Contract.Symbol(&_L2NativeToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_L2NativeToken *L2NativeTokenCallerSession) Symbol() (string, error) {
	return _L2NativeToken.Contract.Symbol(&_L2NativeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_L2NativeToken *L2NativeTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2NativeToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_L2NativeToken *L2NativeTokenSession) TotalSupply() (*big.Int, error) {
	return _L2NativeToken.Contract.TotalSupply(&_L2NativeToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_L2NativeToken *L2NativeTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _L2NativeToken.Contract.TotalSupply(&_L2NativeToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.Approve(&_L2NativeToken.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.Approve(&_L2NativeToken.TransactOpts, spender, amount)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 amount, bytes data) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactor) ApproveAndCall(opts *bind.TransactOpts, spender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "approveAndCall", spender, amount, data)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 amount, bytes data) returns(bool)
func (_L2NativeToken *L2NativeTokenSession) ApproveAndCall(spender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _L2NativeToken.Contract.ApproveAndCall(&_L2NativeToken.TransactOpts, spender, amount, data)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 amount, bytes data) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactorSession) ApproveAndCall(spender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _L2NativeToken.Contract.ApproveAndCall(&_L2NativeToken.TransactOpts, spender, amount, data)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_L2NativeToken *L2NativeTokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.DecreaseAllowance(&_L2NativeToken.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.DecreaseAllowance(&_L2NativeToken.TransactOpts, spender, subtractedValue)
}

// EnableCallback is a paid mutator transaction binding the contract method 0x3113ed5c.
//
// Solidity: function enableCallback(bool _callbackEnabled) returns()
func (_L2NativeToken *L2NativeTokenTransactor) EnableCallback(opts *bind.TransactOpts, _callbackEnabled bool) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "enableCallback", _callbackEnabled)
}

// EnableCallback is a paid mutator transaction binding the contract method 0x3113ed5c.
//
// Solidity: function enableCallback(bool _callbackEnabled) returns()
func (_L2NativeToken *L2NativeTokenSession) EnableCallback(_callbackEnabled bool) (*types.Transaction, error) {
	return _L2NativeToken.Contract.EnableCallback(&_L2NativeToken.TransactOpts, _callbackEnabled)
}

// EnableCallback is a paid mutator transaction binding the contract method 0x3113ed5c.
//
// Solidity: function enableCallback(bool _callbackEnabled) returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) EnableCallback(_callbackEnabled bool) (*types.Transaction, error) {
	return _L2NativeToken.Contract.EnableCallback(&_L2NativeToken.TransactOpts, _callbackEnabled)
}

// Faucet is a paid mutator transaction binding the contract method 0x57915897.
//
// Solidity: function faucet(uint256 _amount) returns()
func (_L2NativeToken *L2NativeTokenTransactor) Faucet(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "faucet", _amount)
}

// Faucet is a paid mutator transaction binding the contract method 0x57915897.
//
// Solidity: function faucet(uint256 _amount) returns()
func (_L2NativeToken *L2NativeTokenSession) Faucet(_amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.Faucet(&_L2NativeToken.TransactOpts, _amount)
}

// Faucet is a paid mutator transaction binding the contract method 0x57915897.
//
// Solidity: function faucet(uint256 _amount) returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) Faucet(_amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.Faucet(&_L2NativeToken.TransactOpts, _amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_L2NativeToken *L2NativeTokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.IncreaseAllowance(&_L2NativeToken.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.IncreaseAllowance(&_L2NativeToken.TransactOpts, spender, addedValue)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x5f112c68.
//
// Solidity: function renounceMinter(address target) returns()
func (_L2NativeToken *L2NativeTokenTransactor) RenounceMinter(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "renounceMinter", target)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x5f112c68.
//
// Solidity: function renounceMinter(address target) returns()
func (_L2NativeToken *L2NativeTokenSession) RenounceMinter(target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenounceMinter(&_L2NativeToken.TransactOpts, target)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x5f112c68.
//
// Solidity: function renounceMinter(address target) returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) RenounceMinter(target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenounceMinter(&_L2NativeToken.TransactOpts, target)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x38bf3cfa.
//
// Solidity: function renounceOwnership(address target) returns()
func (_L2NativeToken *L2NativeTokenTransactor) RenounceOwnership(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "renounceOwnership", target)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x38bf3cfa.
//
// Solidity: function renounceOwnership(address target) returns()
func (_L2NativeToken *L2NativeTokenSession) RenounceOwnership(target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenounceOwnership(&_L2NativeToken.TransactOpts, target)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x38bf3cfa.
//
// Solidity: function renounceOwnership(address target) returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) RenounceOwnership(target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenounceOwnership(&_L2NativeToken.TransactOpts, target)
}

// RenounceOwnership0 is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2NativeToken *L2NativeTokenTransactor) RenounceOwnership0(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "renounceOwnership0")
}

// RenounceOwnership0 is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2NativeToken *L2NativeTokenSession) RenounceOwnership0() (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenounceOwnership0(&_L2NativeToken.TransactOpts)
}

// RenounceOwnership0 is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) RenounceOwnership0() (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenounceOwnership0(&_L2NativeToken.TransactOpts)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x41eb24bb.
//
// Solidity: function renouncePauser(address target) returns()
func (_L2NativeToken *L2NativeTokenTransactor) RenouncePauser(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "renouncePauser", target)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x41eb24bb.
//
// Solidity: function renouncePauser(address target) returns()
func (_L2NativeToken *L2NativeTokenSession) RenouncePauser(target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenouncePauser(&_L2NativeToken.TransactOpts, target)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x41eb24bb.
//
// Solidity: function renouncePauser(address target) returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) RenouncePauser(target common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.RenouncePauser(&_L2NativeToken.TransactOpts, target)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.Transfer(&_L2NativeToken.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.Transfer(&_L2NativeToken.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.TransferFrom(&_L2NativeToken.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_L2NativeToken *L2NativeTokenTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L2NativeToken.Contract.TransferFrom(&_L2NativeToken.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0x6d435421.
//
// Solidity: function transferOwnership(address target, address newOwner) returns()
func (_L2NativeToken *L2NativeTokenTransactor) TransferOwnership(opts *bind.TransactOpts, target common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "transferOwnership", target, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0x6d435421.
//
// Solidity: function transferOwnership(address target, address newOwner) returns()
func (_L2NativeToken *L2NativeTokenSession) TransferOwnership(target common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.TransferOwnership(&_L2NativeToken.TransactOpts, target, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0x6d435421.
//
// Solidity: function transferOwnership(address target, address newOwner) returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) TransferOwnership(target common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.TransferOwnership(&_L2NativeToken.TransactOpts, target, newOwner)
}

// TransferOwnership0 is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2NativeToken *L2NativeTokenTransactor) TransferOwnership0(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L2NativeToken.contract.Transact(opts, "transferOwnership0", newOwner)
}

// TransferOwnership0 is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2NativeToken *L2NativeTokenSession) TransferOwnership0(newOwner common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.TransferOwnership0(&_L2NativeToken.TransactOpts, newOwner)
}

// TransferOwnership0 is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2NativeToken *L2NativeTokenTransactorSession) TransferOwnership0(newOwner common.Address) (*types.Transaction, error) {
	return _L2NativeToken.Contract.TransferOwnership0(&_L2NativeToken.TransactOpts, newOwner)
}

// L2NativeTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the L2NativeToken contract.
type L2NativeTokenApprovalIterator struct {
	Event *L2NativeTokenApproval // Event containing the contract specifics and raw log

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
func (it *L2NativeTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2NativeTokenApproval)
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
		it.Event = new(L2NativeTokenApproval)
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
func (it *L2NativeTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2NativeTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2NativeTokenApproval represents a Approval event raised by the L2NativeToken contract.
type L2NativeTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_L2NativeToken *L2NativeTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*L2NativeTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _L2NativeToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &L2NativeTokenApprovalIterator{contract: _L2NativeToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_L2NativeToken *L2NativeTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *L2NativeTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _L2NativeToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2NativeTokenApproval)
				if err := _L2NativeToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_L2NativeToken *L2NativeTokenFilterer) ParseApproval(log types.Log) (*L2NativeTokenApproval, error) {
	event := new(L2NativeTokenApproval)
	if err := _L2NativeToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2NativeTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L2NativeToken contract.
type L2NativeTokenOwnershipTransferredIterator struct {
	Event *L2NativeTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L2NativeTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2NativeTokenOwnershipTransferred)
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
		it.Event = new(L2NativeTokenOwnershipTransferred)
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
func (it *L2NativeTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2NativeTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2NativeTokenOwnershipTransferred represents a OwnershipTransferred event raised by the L2NativeToken contract.
type L2NativeTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2NativeToken *L2NativeTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L2NativeTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2NativeToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L2NativeTokenOwnershipTransferredIterator{contract: _L2NativeToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2NativeToken *L2NativeTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L2NativeTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2NativeToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2NativeTokenOwnershipTransferred)
				if err := _L2NativeToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L2NativeToken *L2NativeTokenFilterer) ParseOwnershipTransferred(log types.Log) (*L2NativeTokenOwnershipTransferred, error) {
	event := new(L2NativeTokenOwnershipTransferred)
	if err := _L2NativeToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2NativeTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the L2NativeToken contract.
type L2NativeTokenTransferIterator struct {
	Event *L2NativeTokenTransfer // Event containing the contract specifics and raw log

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
func (it *L2NativeTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2NativeTokenTransfer)
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
		it.Event = new(L2NativeTokenTransfer)
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
func (it *L2NativeTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2NativeTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2NativeTokenTransfer represents a Transfer event raised by the L2NativeToken contract.
type L2NativeTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_L2NativeToken *L2NativeTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*L2NativeTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _L2NativeToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &L2NativeTokenTransferIterator{contract: _L2NativeToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_L2NativeToken *L2NativeTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *L2NativeTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _L2NativeToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2NativeTokenTransfer)
				if err := _L2NativeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_L2NativeToken *L2NativeTokenFilterer) ParseTransfer(log types.Log) (*L2NativeTokenTransfer, error) {
	event := new(L2NativeTokenTransfer)
	if err := _L2NativeToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
