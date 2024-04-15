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

// FiatTokenV22MetaData contains all meta data concerning the FiatTokenV22 contract.
var FiatTokenV22MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorizer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"}],\"name\":\"AuthorizationCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"authorizer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"}],\"name\":\"AuthorizationUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"Blacklisted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newBlacklister\",\"type\":\"address\"}],\"name\":\"BlacklisterChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newMasterMinter\",\"type\":\"address\"}],\"name\":\"MasterMinterChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minterAllowedAmount\",\"type\":\"uint256\"}],\"name\":\"MinterConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldMinter\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"PauserChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newRescuer\",\"type\":\"address\"}],\"name\":\"RescuerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"UnBlacklisted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CANCEL_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RECEIVE_WITH_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TRANSFER_WITH_AUTHORIZATION_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authorizer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"}],\"name\":\"authorizationState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"blacklist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blacklister\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authorizer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authorizer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"cancelAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minterAllowedAmount\",\"type\":\"uint256\"}],\"name\":\"configureMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currency\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"decrement\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"increment\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"tokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"tokenSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"tokenCurrency\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"newMasterMinter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newPauser\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newBlacklister\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newName\",\"type\":\"string\"}],\"name\":\"initializeV2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"lostAndFound\",\"type\":\"address\"}],\"name\":\"initializeV2_1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accountsToBlacklist\",\"type\":\"address[]\"},{\"internalType\":\"string\",\"name\":\"newSymbol\",\"type\":\"string\"}],\"name\":\"initializeV2_2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"isBlacklisted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"masterMinter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"minterAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauser\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validAfter\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validBefore\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"receiveWithAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validAfter\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validBefore\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"receiveWithAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"removeMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"tokenContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"rescueERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rescuer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validAfter\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validBefore\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"transferWithAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validAfter\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validBefore\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"transferWithAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unBlacklist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newBlacklister\",\"type\":\"address\"}],\"name\":\"updateBlacklister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newMasterMinter\",\"type\":\"address\"}],\"name\":\"updateMasterMinter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newPauser\",\"type\":\"address\"}],\"name\":\"updatePauser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRescuer\",\"type\":\"address\"}],\"name\":\"updateRescuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805460ff60a01b191690556000600b553480156200002357600080fd5b506200002f3362000035565b62000057565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b615b6880620000676000396000f3fe608060405234801561001057600080fd5b506004361061036d5760003560e01c80638456cb59116101d3578063b7b7289911610104578063e3ee160e116100a2578063ef55bec61161007c578063ef55bec614611122578063f2fde38b1461118e578063f9f92be4146111c1578063fe575a87146111f45761036d565b8063e3ee160e14611075578063e5a6b10f146110e1578063e94a0102146110e95761036d565b8063d505accf116100de578063d505accf14610f64578063d608ea6414610fc2578063d916948714611032578063dd62ed3e1461103a5761036d565b8063b7b7289914610db0578063bd10243014610e78578063cf09299514610e805761036d565b8063a0cc6a6811610171578063aa20e1e41161014b578063aa20e1e414610cd4578063aa271e1a14610d07578063ad38bf2214610d3a578063b2118a8d14610d6d5761036d565b8063a0cc6a6814610c5a578063a457c2d714610c62578063a9059cbb14610c9b5761036d565b80638da5cb5b116101ad5780638da5cb5b14610b6a57806395d89b4114610b725780639fd0506d14610b7a5780639fd5a6cf14610b825761036d565b80638456cb5914610a4b57806388b7ab6314610a535780638a6db9c314610b375761036d565b806338a63183116102ad57806354fd4d501161024b5780635c975abb116102255780635c975abb146109d557806370a08231146109dd5780637ecebe0014610a105780637f2eecc314610a435761036d565b806354fd4d501461094c578063554bab3c146109545780635a049a70146109875761036d565b806340c10f191161028757806340c10f19146107fb57806342966c6814610834578063430239b4146108515780634e44d956146109135761036d565b806338a63183146107b257806339509351146107ba5780633f4ba83a146107f35761036d565b80632fc81e091161031a578063313ce567116102f4578063313ce5671461056f5780633357162b1461058d57806335d99f35146107795780633644e515146107aa5761036d565b80632fc81e09146105015780633092afd51461053457806330adf81f146105675761036d565b80631a8952661161034b5780631a8952661461045657806323b872dd1461048b5780632ab60045146104ce5761036d565b806306fdde0314610372578063095ea7b3146103ef57806318160ddd1461043c575b600080fd5b61037a611227565b6040805160208082528351818301528351919283929083019185019080838360005b838110156103b457818101518382015260200161039c565b50505050905090810190601f1680156103e15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6104286004803603604081101561040557600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356112d3565b604080519115158252519081900360200190f35b610444611374565b60408051918252519081900360200190f35b6104896004803603602081101561046c57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661137a565b005b610428600480360360608110156104a157600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060400135611437565b610489600480360360208110156104e457600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166116f2565b6104896004803603602081101561051757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611853565b6104286004803603602081101561054a57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166118bb565b6104446119b4565b6105776119d8565b6040805160ff9092168252519081900360200190f35b61048960048036036101008110156105a457600080fd5b8101906020810181356401000000008111156105bf57600080fd5b8201836020820111156105d157600080fd5b803590602001918460018302840111640100000000831117156105f357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929594936020810193503591505064010000000081111561064657600080fd5b82018360208201111561065857600080fd5b8035906020019184600183028401116401000000008311171561067a57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092959493602081019350359150506401000000008111156106cd57600080fd5b8201836020820111156106df57600080fd5b8035906020019184600183028401116401000000008311171561070157600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505050813560ff16925050602081013573ffffffffffffffffffffffffffffffffffffffff908116916040810135821691606082013581169160800135166119e1565b610781611d23565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610444611d3f565b610781611d4e565b610428600480360360408110156107d057600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135611d6a565b610489611e02565b6104286004803603604081101561081157600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135611ec5565b6104896004803603602081101561084a57600080fd5b5035612296565b6104896004803603604081101561086757600080fd5b81019060208101813564010000000081111561088257600080fd5b82018360208201111561089457600080fd5b803590602001918460208302840111640100000000831117156108b657600080fd5b9193909290916020810190356401000000008111156108d457600080fd5b8201836020820111156108e657600080fd5b8035906020019184600183028401116401000000008311171561090857600080fd5b509092509050612538565b6104286004803603604081101561092957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356126ef565b61037a612882565b6104896004803603602081101561096a57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166128b9565b610489600480360360a081101561099d57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060208101359060ff6040820135169060608101359060800135612a20565b610428612abe565b610444600480360360208110156109f357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16612adf565b61044460048036036020811015610a2657600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16612af0565b610444612b18565b610489612b3c565b610489600480360360e0811015610a6957600080fd5b73ffffffffffffffffffffffffffffffffffffffff823581169260208101359091169160408201359160608101359160808201359160a08101359181019060e0810160c0820135640100000000811115610ac257600080fd5b820183602082011115610ad457600080fd5b80359060200191846001830284011164010000000083111715610af657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550612c16945050505050565b61044460048036036020811015610b4d57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16612d7a565b610781612da2565b61037a612dbe565b610781612e37565b610489600480360360a0811015610b9857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823581169260208101359091169160408201359160608101359181019060a081016080820135640100000000811115610be557600080fd5b820183602082011115610bf757600080fd5b80359060200191846001830284011164010000000083111715610c1957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550612e53945050505050565b610444612eea565b61042860048036036040811015610c7857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135612f0e565b61042860048036036040811015610cb157600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135612fa6565b61048960048036036020811015610cea57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613109565b61042860048036036020811015610d1d57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613270565b61048960048036036020811015610d5057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661329b565b61048960048036036060811015610d8357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060400135613402565b61048960048036036060811015610dc657600080fd5b73ffffffffffffffffffffffffffffffffffffffff82351691602081013591810190606081016040820135640100000000811115610e0357600080fd5b820183602082011115610e1557600080fd5b80359060200191846001830284011164010000000083111715610e3757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550613498945050505050565b61078161352d565b610489600480360360e0811015610e9657600080fd5b73ffffffffffffffffffffffffffffffffffffffff823581169260208101359091169160408201359160608101359160808201359160a08101359181019060e0810160c0820135640100000000811115610eef57600080fd5b820183602082011115610f0157600080fd5b80359060200191846001830284011164010000000083111715610f2357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550613549945050505050565b610489600480360360e0811015610f7a57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060408101359060608101359060ff6080820135169060a08101359060c001356136a2565b61048960048036036020811015610fd857600080fd5b810190602081018135640100000000811115610ff357600080fd5b82018360208201111561100557600080fd5b8035906020019184600183028401116401000000008311171561102757600080fd5b509092509050613744565b61044461382d565b6104446004803603604081101561105057600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516613851565b610489600480360361012081101561108c57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060408101359060608101359060808101359060a08101359060ff60c0820135169060e0810135906101000135613889565b61037a6139f1565b610428600480360360408110156110ff57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135613a6a565b610489600480360361012081101561113957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060408101359060608101359060808101359060a08101359060ff60c0820135169060e0810135906101000135613aa2565b610489600480360360208110156111a457600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613bfd565b610489600480360360208110156111d757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613d50565b6104286004803603602081101561120a57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613e0d565b6004805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f810184900484028201840190925281815292918301828280156112cb5780601f106112a0576101008083540402835291602001916112cb565b820191906000526020600020905b8154815290600101906020018083116112ae57829003601f168201915b505050505081565b60015460009074010000000000000000000000000000000000000000900460ff161561136057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b61136b338484613e18565b50600192915050565b600b5490565b60025473ffffffffffffffffffffffffffffffffffffffff1633146113ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c81526020018061580d602c913960400191505060405180910390fd5b6113f381613f5f565b60405173ffffffffffffffffffffffffffffffffffffffff8216907f117e3210bb9aa7d9baff172026820255c6f6c30ba8999d1c2fd88e2848137c4e90600090a250565b60015460009074010000000000000000000000000000000000000000900460ff16156114c457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b336114ce81613f6a565b15611524576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b8461152e81613f6a565b15611584576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b8461158e81613f6a565b156115e4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff87166000908152600a6020908152604080832033845290915290205485111561166d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806158fd6028913960400191505060405180910390fd5b611678878787613f98565b73ffffffffffffffffffffffffffffffffffffffff87166000908152600a602090815260408083203384529091529020546116b39086614163565b73ffffffffffffffffffffffffffffffffffffffff88166000908152600a60209081526040808320338452909152902055600193505050509392505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461177857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81166117e4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180615746602a913960400191505060405180910390fd5b600e80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040517fe475e580d85111348e40d8ca33cfdd74c30fe1655c2d8537a13abc10065ffa5a90600090a250565b60125460ff1660011461186557600080fd5b6000611870306141ac565b9050801561188357611883308383613f98565b61188c306141f6565b5050601280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166002179055565b60085460009073ffffffffffffffffffffffffffffffffffffffff16331461192e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806157e46029913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000818152600c6020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055600d909152808220829055517fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb666929190a2506001919050565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c981565b60065460ff1681565b60085474010000000000000000000000000000000000000000900460ff1615611a55576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180615978602a913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8416611ac1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f8152602001806158aa602f913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316611b2d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602981526020018061571d6029913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216611b99576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180615925602e913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116611c05576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180615a656028913960400191505060405180910390fd5b8751611c189060049060208b01906154b6565b508651611c2c9060059060208a01906154b6565b508551611c409060079060208901906154b6565b50600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff8716179055600880547fffffffffffffffffffffffff000000000000000000000000000000000000000090811673ffffffffffffffffffffffffffffffffffffffff8781169190911790925560018054821686841617905560028054909116918416919091179055611cda81614201565b5050600880547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055505050505050565b60085473ffffffffffffffffffffffffffffffffffffffff1681565b6000611d49614248565b905090565b600e5473ffffffffffffffffffffffffffffffffffffffff1690565b60015460009074010000000000000000000000000000000000000000900460ff1615611df757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b61136b33848461433d565b60015473ffffffffffffffffffffffffffffffffffffffff163314611e72576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180615a196022913960400191505060405180910390fd5b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3390600090a1565b60015460009074010000000000000000000000000000000000000000900460ff1615611f5257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b336000908152600c602052604090205460ff16611fba576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806158896021913960400191505060405180910390fd5b33611fc481613f6a565b1561201a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b8361202481613f6a565b1561207a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff85166120e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806156b26023913960400191505060405180910390fd5b6000841161213f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806157956029913960400191505060405180910390fd5b336000908152600d6020526040902054808511156121a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e8152602001806159eb602e913960400191505060405180910390fd5b600b546121b59086614387565b600b556121d4866121cf876121c9836141ac565b90614387565b6143fb565b6121de8186614163565b336000818152600d6020908152604091829020939093558051888152905173ffffffffffffffffffffffffffffffffffffffff8a16937fab8530f87dc9b59234c4623bf917212bb2536d647574c8e7e5da92c2ede0c9f8928290030190a360408051868152905173ffffffffffffffffffffffffffffffffffffffff8816916000917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a350600195945050505050565b60015474010000000000000000000000000000000000000000900460ff161561232057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b336000908152600c602052604090205460ff16612388576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806158896021913960400191505060405180910390fd5b3361239281613f6a565b156123e8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b60006123f3336141ac565b90506000831161244e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806156896029913960400191505060405180910390fd5b828110156124a7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806158636026913960400191505060405180910390fd5b600b546124b49084614163565b600b556124c5336121cf8386614163565b60408051848152905133917fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5919081900360200190a260408051848152905160009133917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a3505050565b60125460ff1660021461254a57600080fd5b61255660058383615534565b5060005b83811015612698576003600086868481811061257257fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff168352508101919091526040016000205460ff166125fa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603d8152602001806155d6603d913960400191505060405180910390fd5b61262b85858381811061260957fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff166141f6565b6003600086868481811061263b57fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905560010161255a565b506126a2306141f6565b505030600090815260036020819052604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009081169091556012805490911690911790555050565b60015460009074010000000000000000000000000000000000000000900460ff161561277c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b60085473ffffffffffffffffffffffffffffffffffffffff1633146127ec576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806157e46029913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000818152600c6020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055600d825291829020859055815185815291517f46980fca912ef9bcdbd36877427b6b90e860769f604e89c0e67720cece530d209281900390910190a250600192915050565b60408051808201909152600181527f3200000000000000000000000000000000000000000000000000000000000000602082015290565b60005473ffffffffffffffffffffffffffffffffffffffff16331461293f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81166129ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806156366028913960400191505060405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117918290556040519116907fb80482a293ca2e013eda8683c9bd7fc8347cfdaeea5ede58cba46df502c2a60490600090a250565b60015474010000000000000000000000000000000000000000900460ff1615612aaa57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b612ab785858585856144fc565b5050505050565b60015474010000000000000000000000000000000000000000900460ff1681565b6000612aea826141ac565b92915050565b73ffffffffffffffffffffffffffffffffffffffff1660009081526011602052604090205490565b7fd099cc98ef71107a616c4f0f941f04c322d8e254fe26b3c6668db87aae413de881565b60015473ffffffffffffffffffffffffffffffffffffffff163314612bac576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180615a196022913960400191505060405180910390fd5b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62590600090a1565b60015474010000000000000000000000000000000000000000900460ff1615612ca057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b86612caa81613f6a565b15612d00576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b86612d0a81613f6a565b15612d60576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b612d6f8989898989898961453c565b505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff166000908152600d602052604090205490565b60005473ffffffffffffffffffffffffffffffffffffffff1690565b6005805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f810184900484028201840190925281815292918301828280156112cb5780601f106112a0576101008083540402835291602001916112cb565b60015473ffffffffffffffffffffffffffffffffffffffff1681565b60015474010000000000000000000000000000000000000000900460ff1615612edd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b612ab7858585858561465d565b7f7c7c6cdb67a18743f49ec6fa9b35f50d52ed05cbed4cc592e13b44501c1a226781565b60015460009074010000000000000000000000000000000000000000900460ff1615612f9b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b61136b338484614921565b60015460009074010000000000000000000000000000000000000000900460ff161561303357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b3361303d81613f6a565b15613093576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b8361309d81613f6a565b156130f3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b6130fe338686613f98565b506001949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461318f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81166131fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f8152602001806158aa602f913960400191505060405180910390fd5b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117918290556040519116907fdb66dfa9c6b8f5226fe9aac7e51897ae8ee94ac31dc70bb6c9900b2574b707e690600090a250565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602052604090205460ff1690565b60005473ffffffffffffffffffffffffffffffffffffffff16331461332157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811661338d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526032815260200180615abb6032913960400191505060405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117918290556040519116907fc67398012c111ce95ecb7429b933096c977380ee6c421175a71a4a4c6c88c06e90600090a250565b600e5473ffffffffffffffffffffffffffffffffffffffff163314613472576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260248152602001806158d96024913960400191505060405180910390fd5b61349373ffffffffffffffffffffffffffffffffffffffff8416838361497d565b505050565b60015474010000000000000000000000000000000000000000900460ff161561352257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b613493838383614a0a565b60025473ffffffffffffffffffffffffffffffffffffffff1681565b60015474010000000000000000000000000000000000000000900460ff16156135d357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b866135dd81613f6a565b15613633576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b8661363d81613f6a565b15613693576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b612d6f89898989898989614b14565b60015474010000000000000000000000000000000000000000900460ff161561372c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b61373b87878787878787614bb2565b50505050505050565b60085474010000000000000000000000000000000000000000900460ff168015613771575060125460ff16155b61377a57600080fd5b61378660048383615534565b506137fb82828080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060408051808201909152600181527f320000000000000000000000000000000000000000000000000000000000000060208201529150614bf49050565b600f555050601280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055565b7f158b0a9edf7a828aad02f63cd515c68ef2f50ba807396f6d12842833a159742981565b73ffffffffffffffffffffffffffffffffffffffff9182166000908152600a6020908152604080832093909416825291909152205490565b60015474010000000000000000000000000000000000000000900460ff161561391357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b8861391d81613f6a565b15613973576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b8861397d81613f6a565b156139d3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b6139e48b8b8b8b8b8b8b8b8b614c0a565b5050505050505050505050565b6007805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f810184900484028201840190925281815292918301828280156112cb5780601f106112a0576101008083540402835291602001916112cb565b73ffffffffffffffffffffffffffffffffffffffff919091166000908152601060209081526040808320938352929052205460ff1690565b60015474010000000000000000000000000000000000000000900460ff1615613b2c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015290519081900360640190fd5b88613b3681613f6a565b15613b8c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b88613b9681613f6a565b15613bec576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615aed6025913960400191505060405180910390fd5b6139e48b8b8b8b8b8b8b8b8b614c4e565b60005473ffffffffffffffffffffffffffffffffffffffff163314613c8357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8116613cef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806156d56026913960400191505060405180910390fd5b6000546040805173ffffffffffffffffffffffffffffffffffffffff9283168152918316602083015280517f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09281900390910190a1613d4d81614201565b50565b60025473ffffffffffffffffffffffffffffffffffffffff163314613dc0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c81526020018061580d602c913960400191505060405180910390fd5b613dc9816141f6565b60405173ffffffffffffffffffffffffffffffffffffffff8216907fffa4e6181777692565cf28528fc88fd1516ea86b56da075235fa575af6a4b85590600090a250565b6000612aea82613f6a565b73ffffffffffffffffffffffffffffffffffffffff8316613e84576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260248152602001806159c76024913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216613ef0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806156fb6022913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8084166000818152600a6020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b613d4d816000614c92565b73ffffffffffffffffffffffffffffffffffffffff1660009081526009602052604090205460ff1c60011490565b73ffffffffffffffffffffffffffffffffffffffff8316614004576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806159a26025913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216614070576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806156136023913960400191505060405180910390fd5b614079836141ac565b8111156140d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806157be6026913960400191505060405180910390fd5b6140e8836121cf836140e2876141ac565b90614163565b6140f9826121cf836121c9866141ac565b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a3505050565b60006141a583836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f770000815250614d1b565b9392505050565b73ffffffffffffffffffffffffffffffffffffffff166000908152600960205260409020547f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b613d4d816001614c92565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6004805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f8101849004840282018401909252818152600093611d4993919290918301828280156142f55780601f106142ca576101008083540402835291602001916142f5565b820191906000526020600020905b8154815290600101906020018083116142d857829003601f168201915b50505050506040518060400160405280600181526020017f3200000000000000000000000000000000000000000000000000000000000000815250614338614dcc565b614dd0565b73ffffffffffffffffffffffffffffffffffffffff8084166000908152600a602090815260408083209386168352929052205461349390849084906143829085614387565b613e18565b6000828201838110156141a557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811115614474576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180615839602a913960400191505060405180910390fd5b61447d82613f6a565b156144d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806157706025913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff909116600090815260096020526040902055565b612ab78585848487604051602001808481526020018381526020018260ff1660f81b81526001019350505050604051602081830303815290604052614a0a565b73ffffffffffffffffffffffffffffffffffffffff861633146145aa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806159536025913960400191505060405180910390fd5b6145b687838686614e44565b604080517fd099cc98ef71107a616c4f0f941f04c322d8e254fe26b3c6668db87aae413de860208083019190915273ffffffffffffffffffffffffffffffffffffffff808b1683850152891660608301526080820188905260a0820187905260c0820186905260e080830186905283518084039091018152610100909201909252805191012061464890889083614f04565b6146528783615082565b61373b878787613f98565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82148061468b5750428210155b6146f657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f46696174546f6b656e56323a207065726d697420697320657870697265640000604482015290519081900360640190fd5b600061479e614703614248565b73ffffffffffffffffffffffffffffffffffffffff80891660008181526011602090815260409182902080546001810190915582517f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c98184015280840194909452938b166060840152608083018a905260a083019390935260c08083018990528151808403909101815260e090920190528051910120615107565b9050734200000000000000000000000000000000000776636ccea6528783856040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561482b578181015183820152602001614813565b50505050905090810190601f1680156148585780820380516001836020036101000a031916815260200191505b5094505050505060206040518083038186803b15801561487757600080fd5b505af415801561488b573d6000803e3d6000fd5b505050506040513d60208110156148a157600080fd5b505161490e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f454950323631323a20696e76616c6964207369676e6174757265000000000000604482015290519081900360640190fd5b614919868686613e18565b505050505050565b613493838361438284604051806060016040528060258152602001615b376025913973ffffffffffffffffffffffffffffffffffffffff808a166000908152600a60209081526040808320938c16835292905220549190614d1b565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052613493908490615141565b614a148383615219565b614a8e837f158b0a9edf7a828aad02f63cd515c68ef2f50ba807396f6d12842833a159742960001b8585604051602001808481526020018373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200193505050506040516020818303038152906040528051906020012083614f04565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260106020908152604080832086845290915280822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055518492917f1cdd46ff242716cdaa72d159d339a485b3438398348d68f09d7c8c0a59353d8191a3505050565b614b2087838686614e44565b604080517f7c7c6cdb67a18743f49ec6fa9b35f50d52ed05cbed4cc592e13b44501c1a226760208083019190915273ffffffffffffffffffffffffffffffffffffffff808b1683850152891660608301526080820188905260a0820187905260c0820186905260e080830186905283518084039091018152610100909201909252805191012061464890889083614f04565b61373b87878787868689604051602001808481526020018381526020018260ff1660f81b8152600101935050505060405160208183030381529060405261465d565b600046614c02848483614dd0565b949350505050565b612d6f89898989898988888b604051602001808481526020018381526020018260ff1660f81b81526001019350505050604051602081830303815290604052614b14565b612d6f89898989898988888b604051602001808481526020018381526020018260ff1660f81b8152600101935050505060405160208183030381529060405261453c565b80614ca557614ca0826141ac565b614cee565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600960205260409020547f8000000000000000000000000000000000000000000000000000000000000000175b73ffffffffffffffffffffffffffffffffffffffff90921660009081526009602052604090209190915550565b60008184841115614dc4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015614d89578181015183820152602001614d71565b50505050905090810190601f168015614db65780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b505050900390565b4690565b8251602093840120825192840192909220604080517f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f8187015280820194909452606084019190915260808301919091523060a0808401919091528151808403909101815260c09092019052805191012090565b814211614e9c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b81526020018061565e602b913960400191505060405180910390fd5b804210614ef4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180615b126025913960400191505060405180910390fd5b614efe8484615219565b50505050565b734200000000000000000000000000000000000776636ccea65284614f30614f2a614248565b86615107565b846040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015614f9f578181015183820152602001614f87565b50505050905090810190601f168015614fcc5780820380516001836020036101000a031916815260200191505b5094505050505060206040518083038186803b158015614feb57600080fd5b505af4158015614fff573d6000803e3d6000fd5b505050506040513d602081101561501557600080fd5b505161349357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f46696174546f6b656e56323a20696e76616c6964207369676e61747572650000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8216600081815260106020908152604080832085845290915280822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055518392917f98de503528ee59b575ef0c0a2576a82497bfc029a5685b209e9ec333479b10a591a35050565b6040517f19010000000000000000000000000000000000000000000000000000000000008152600281019290925260228201526042902090565b60606151a3826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166152a79092919063ffffffff16565b805190915015613493578080602001905160208110156151c257600080fd5b5051613493576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180615a3b602a913960400191505060405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216600090815260106020908152604080832084845290915290205460ff16156152a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180615a8d602e913960400191505060405180910390fd5b5050565b6060614c02848460008560606152bc8561547d565b61532757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015290519081900360640190fd5b600060608673ffffffffffffffffffffffffffffffffffffffff1685876040518082805190602001908083835b6020831061539157805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101615354565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038185875af1925050503d80600081146153f3576040519150601f19603f3d011682016040523d82523d6000602084013e6153f8565b606091505b5091509150811561540c579150614c029050565b80511561541c5780518082602001fd5b6040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201818152865160248401528651879391928392604401919085019080838360008315614d89578181015183820152602001614d71565b6000813f7fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470818114801590614c02575050151592915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106154f757805160ff1916838001178555615524565b82800160010185558215615524579182015b82811115615524578251825591602001919060010190615509565b506155309291506155c0565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10615593578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555615524565b82800160010185558215615524579182015b828111156155245782358255916020019190600101906155a5565b5b8082111561553057600081556001016155c156fe46696174546f6b656e56325f323a20426c61636b6c697374696e672070726576696f75736c7920756e626c61636b6c6973746564206163636f756e742145524332303a207472616e7366657220746f20746865207a65726f20616464726573735061757361626c653a206e65772070617573657220697320746865207a65726f206164647265737346696174546f6b656e56323a20617574686f72697a6174696f6e206973206e6f74207965742076616c696446696174546f6b656e3a206275726e20616d6f756e74206e6f742067726561746572207468616e203046696174546f6b656e3a206d696e7420746f20746865207a65726f20616464726573734f776e61626c653a206e6577206f776e657220697320746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737346696174546f6b656e3a206e65772070617573657220697320746865207a65726f2061646472657373526573637561626c653a206e6577207265736375657220697320746865207a65726f206164647265737346696174546f6b656e56325f323a204163636f756e7420697320626c61636b6c697374656446696174546f6b656e3a206d696e7420616d6f756e74206e6f742067726561746572207468616e203045524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e636546696174546f6b656e3a2063616c6c6572206973206e6f7420746865206d61737465724d696e746572426c61636b6c69737461626c653a2063616c6c6572206973206e6f742074686520626c61636b6c697374657246696174546f6b656e56325f323a2042616c616e636520657863656564732028325e323535202d20312946696174546f6b656e3a206275726e20616d6f756e7420657863656564732062616c616e636546696174546f6b656e3a2063616c6c6572206973206e6f742061206d696e74657246696174546f6b656e3a206e6577206d61737465724d696e74657220697320746865207a65726f2061646472657373526573637561626c653a2063616c6c6572206973206e6f7420746865207265736375657245524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636546696174546f6b656e3a206e657720626c61636b6c697374657220697320746865207a65726f206164647265737346696174546f6b656e56323a2063616c6c6572206d7573742062652074686520706179656546696174546f6b656e3a20636f6e747261637420697320616c726561647920696e697469616c697a656445524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f206164647265737346696174546f6b656e3a206d696e7420616d6f756e742065786365656473206d696e746572416c6c6f77616e63655061757361626c653a2063616c6c6572206973206e6f7420746865207061757365725361666545524332303a204552433230206f7065726174696f6e20646964206e6f74207375636365656446696174546f6b656e3a206e6577206f776e657220697320746865207a65726f206164647265737346696174546f6b656e56323a20617574686f72697a6174696f6e2069732075736564206f722063616e63656c6564426c61636b6c69737461626c653a206e657720626c61636b6c697374657220697320746865207a65726f2061646472657373426c61636b6c69737461626c653a206163636f756e7420697320626c61636b6c697374656446696174546f6b656e56323a20617574686f72697a6174696f6e206973206578706972656445524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa164736f6c634300060c000a",
}

// FiatTokenV22ABI is the input ABI used to generate the binding from.
// Deprecated: Use FiatTokenV22MetaData.ABI instead.
var FiatTokenV22ABI = FiatTokenV22MetaData.ABI

// FiatTokenV22Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FiatTokenV22MetaData.Bin instead.
var FiatTokenV22Bin = FiatTokenV22MetaData.Bin

// DeployFiatTokenV22 deploys a new Ethereum contract, binding an instance of FiatTokenV22 to it.
func DeployFiatTokenV22(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FiatTokenV22, error) {
	parsed, err := FiatTokenV22MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FiatTokenV22Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FiatTokenV22{FiatTokenV22Caller: FiatTokenV22Caller{contract: contract}, FiatTokenV22Transactor: FiatTokenV22Transactor{contract: contract}, FiatTokenV22Filterer: FiatTokenV22Filterer{contract: contract}}, nil
}

// FiatTokenV22 is an auto generated Go binding around an Ethereum contract.
type FiatTokenV22 struct {
	FiatTokenV22Caller     // Read-only binding to the contract
	FiatTokenV22Transactor // Write-only binding to the contract
	FiatTokenV22Filterer   // Log filterer for contract events
}

// FiatTokenV22Caller is an auto generated read-only Go binding around an Ethereum contract.
type FiatTokenV22Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatTokenV22Transactor is an auto generated write-only Go binding around an Ethereum contract.
type FiatTokenV22Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatTokenV22Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FiatTokenV22Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatTokenV22Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FiatTokenV22Session struct {
	Contract     *FiatTokenV22     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FiatTokenV22CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FiatTokenV22CallerSession struct {
	Contract *FiatTokenV22Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// FiatTokenV22TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FiatTokenV22TransactorSession struct {
	Contract     *FiatTokenV22Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// FiatTokenV22Raw is an auto generated low-level Go binding around an Ethereum contract.
type FiatTokenV22Raw struct {
	Contract *FiatTokenV22 // Generic contract binding to access the raw methods on
}

// FiatTokenV22CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FiatTokenV22CallerRaw struct {
	Contract *FiatTokenV22Caller // Generic read-only contract binding to access the raw methods on
}

// FiatTokenV22TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FiatTokenV22TransactorRaw struct {
	Contract *FiatTokenV22Transactor // Generic write-only contract binding to access the raw methods on
}

// NewFiatTokenV22 creates a new instance of FiatTokenV22, bound to a specific deployed contract.
func NewFiatTokenV22(address common.Address, backend bind.ContractBackend) (*FiatTokenV22, error) {
	contract, err := bindFiatTokenV22(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22{FiatTokenV22Caller: FiatTokenV22Caller{contract: contract}, FiatTokenV22Transactor: FiatTokenV22Transactor{contract: contract}, FiatTokenV22Filterer: FiatTokenV22Filterer{contract: contract}}, nil
}

// NewFiatTokenV22Caller creates a new read-only instance of FiatTokenV22, bound to a specific deployed contract.
func NewFiatTokenV22Caller(address common.Address, caller bind.ContractCaller) (*FiatTokenV22Caller, error) {
	contract, err := bindFiatTokenV22(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22Caller{contract: contract}, nil
}

// NewFiatTokenV22Transactor creates a new write-only instance of FiatTokenV22, bound to a specific deployed contract.
func NewFiatTokenV22Transactor(address common.Address, transactor bind.ContractTransactor) (*FiatTokenV22Transactor, error) {
	contract, err := bindFiatTokenV22(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22Transactor{contract: contract}, nil
}

// NewFiatTokenV22Filterer creates a new log filterer instance of FiatTokenV22, bound to a specific deployed contract.
func NewFiatTokenV22Filterer(address common.Address, filterer bind.ContractFilterer) (*FiatTokenV22Filterer, error) {
	contract, err := bindFiatTokenV22(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22Filterer{contract: contract}, nil
}

// bindFiatTokenV22 binds a generic wrapper to an already deployed contract.
func bindFiatTokenV22(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FiatTokenV22MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FiatTokenV22 *FiatTokenV22Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FiatTokenV22.Contract.FiatTokenV22Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FiatTokenV22 *FiatTokenV22Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.FiatTokenV22Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FiatTokenV22 *FiatTokenV22Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.FiatTokenV22Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FiatTokenV22 *FiatTokenV22CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FiatTokenV22.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FiatTokenV22 *FiatTokenV22TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FiatTokenV22 *FiatTokenV22TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.contract.Transact(opts, method, params...)
}

// CANCELAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0xd9169487.
//
// Solidity: function CANCEL_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Caller) CANCELAUTHORIZATIONTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "CANCEL_AUTHORIZATION_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CANCELAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0xd9169487.
//
// Solidity: function CANCEL_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Session) CANCELAUTHORIZATIONTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.CANCELAUTHORIZATIONTYPEHASH(&_FiatTokenV22.CallOpts)
}

// CANCELAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0xd9169487.
//
// Solidity: function CANCEL_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22CallerSession) CANCELAUTHORIZATIONTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.CANCELAUTHORIZATIONTYPEHASH(&_FiatTokenV22.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Caller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Session) DOMAINSEPARATOR() ([32]byte, error) {
	return _FiatTokenV22.Contract.DOMAINSEPARATOR(&_FiatTokenV22.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22CallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _FiatTokenV22.Contract.DOMAINSEPARATOR(&_FiatTokenV22.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Caller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "PERMIT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Session) PERMITTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.PERMITTYPEHASH(&_FiatTokenV22.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22CallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.PERMITTYPEHASH(&_FiatTokenV22.CallOpts)
}

// RECEIVEWITHAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0x7f2eecc3.
//
// Solidity: function RECEIVE_WITH_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Caller) RECEIVEWITHAUTHORIZATIONTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "RECEIVE_WITH_AUTHORIZATION_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RECEIVEWITHAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0x7f2eecc3.
//
// Solidity: function RECEIVE_WITH_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Session) RECEIVEWITHAUTHORIZATIONTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.RECEIVEWITHAUTHORIZATIONTYPEHASH(&_FiatTokenV22.CallOpts)
}

// RECEIVEWITHAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0x7f2eecc3.
//
// Solidity: function RECEIVE_WITH_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22CallerSession) RECEIVEWITHAUTHORIZATIONTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.RECEIVEWITHAUTHORIZATIONTYPEHASH(&_FiatTokenV22.CallOpts)
}

// TRANSFERWITHAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0xa0cc6a68.
//
// Solidity: function TRANSFER_WITH_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Caller) TRANSFERWITHAUTHORIZATIONTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "TRANSFER_WITH_AUTHORIZATION_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TRANSFERWITHAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0xa0cc6a68.
//
// Solidity: function TRANSFER_WITH_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22Session) TRANSFERWITHAUTHORIZATIONTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.TRANSFERWITHAUTHORIZATIONTYPEHASH(&_FiatTokenV22.CallOpts)
}

// TRANSFERWITHAUTHORIZATIONTYPEHASH is a free data retrieval call binding the contract method 0xa0cc6a68.
//
// Solidity: function TRANSFER_WITH_AUTHORIZATION_TYPEHASH() view returns(bytes32)
func (_FiatTokenV22 *FiatTokenV22CallerSession) TRANSFERWITHAUTHORIZATIONTYPEHASH() ([32]byte, error) {
	return _FiatTokenV22.Contract.TRANSFERWITHAUTHORIZATIONTYPEHASH(&_FiatTokenV22.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.Allowance(&_FiatTokenV22.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.Allowance(&_FiatTokenV22.CallOpts, owner, spender)
}

// AuthorizationState is a free data retrieval call binding the contract method 0xe94a0102.
//
// Solidity: function authorizationState(address authorizer, bytes32 nonce) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Caller) AuthorizationState(opts *bind.CallOpts, authorizer common.Address, nonce [32]byte) (bool, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "authorizationState", authorizer, nonce)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizationState is a free data retrieval call binding the contract method 0xe94a0102.
//
// Solidity: function authorizationState(address authorizer, bytes32 nonce) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) AuthorizationState(authorizer common.Address, nonce [32]byte) (bool, error) {
	return _FiatTokenV22.Contract.AuthorizationState(&_FiatTokenV22.CallOpts, authorizer, nonce)
}

// AuthorizationState is a free data retrieval call binding the contract method 0xe94a0102.
//
// Solidity: function authorizationState(address authorizer, bytes32 nonce) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22CallerSession) AuthorizationState(authorizer common.Address, nonce [32]byte) (bool, error) {
	return _FiatTokenV22.Contract.AuthorizationState(&_FiatTokenV22.CallOpts, authorizer, nonce)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.BalanceOf(&_FiatTokenV22.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.BalanceOf(&_FiatTokenV22.CallOpts, account)
}

// Blacklister is a free data retrieval call binding the contract method 0xbd102430.
//
// Solidity: function blacklister() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Caller) Blacklister(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "blacklister")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Blacklister is a free data retrieval call binding the contract method 0xbd102430.
//
// Solidity: function blacklister() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Session) Blacklister() (common.Address, error) {
	return _FiatTokenV22.Contract.Blacklister(&_FiatTokenV22.CallOpts)
}

// Blacklister is a free data retrieval call binding the contract method 0xbd102430.
//
// Solidity: function blacklister() view returns(address)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Blacklister() (common.Address, error) {
	return _FiatTokenV22.Contract.Blacklister(&_FiatTokenV22.CallOpts)
}

// Currency is a free data retrieval call binding the contract method 0xe5a6b10f.
//
// Solidity: function currency() view returns(string)
func (_FiatTokenV22 *FiatTokenV22Caller) Currency(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "currency")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Currency is a free data retrieval call binding the contract method 0xe5a6b10f.
//
// Solidity: function currency() view returns(string)
func (_FiatTokenV22 *FiatTokenV22Session) Currency() (string, error) {
	return _FiatTokenV22.Contract.Currency(&_FiatTokenV22.CallOpts)
}

// Currency is a free data retrieval call binding the contract method 0xe5a6b10f.
//
// Solidity: function currency() view returns(string)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Currency() (string, error) {
	return _FiatTokenV22.Contract.Currency(&_FiatTokenV22.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FiatTokenV22 *FiatTokenV22Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FiatTokenV22 *FiatTokenV22Session) Decimals() (uint8, error) {
	return _FiatTokenV22.Contract.Decimals(&_FiatTokenV22.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Decimals() (uint8, error) {
	return _FiatTokenV22.Contract.Decimals(&_FiatTokenV22.CallOpts)
}

// IsBlacklisted is a free data retrieval call binding the contract method 0xfe575a87.
//
// Solidity: function isBlacklisted(address _account) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Caller) IsBlacklisted(opts *bind.CallOpts, _account common.Address) (bool, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "isBlacklisted", _account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBlacklisted is a free data retrieval call binding the contract method 0xfe575a87.
//
// Solidity: function isBlacklisted(address _account) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) IsBlacklisted(_account common.Address) (bool, error) {
	return _FiatTokenV22.Contract.IsBlacklisted(&_FiatTokenV22.CallOpts, _account)
}

// IsBlacklisted is a free data retrieval call binding the contract method 0xfe575a87.
//
// Solidity: function isBlacklisted(address _account) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22CallerSession) IsBlacklisted(_account common.Address) (bool, error) {
	return _FiatTokenV22.Contract.IsBlacklisted(&_FiatTokenV22.CallOpts, _account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Caller) IsMinter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "isMinter", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) IsMinter(account common.Address) (bool, error) {
	return _FiatTokenV22.Contract.IsMinter(&_FiatTokenV22.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_FiatTokenV22 *FiatTokenV22CallerSession) IsMinter(account common.Address) (bool, error) {
	return _FiatTokenV22.Contract.IsMinter(&_FiatTokenV22.CallOpts, account)
}

// MasterMinter is a free data retrieval call binding the contract method 0x35d99f35.
//
// Solidity: function masterMinter() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Caller) MasterMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "masterMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MasterMinter is a free data retrieval call binding the contract method 0x35d99f35.
//
// Solidity: function masterMinter() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Session) MasterMinter() (common.Address, error) {
	return _FiatTokenV22.Contract.MasterMinter(&_FiatTokenV22.CallOpts)
}

// MasterMinter is a free data retrieval call binding the contract method 0x35d99f35.
//
// Solidity: function masterMinter() view returns(address)
func (_FiatTokenV22 *FiatTokenV22CallerSession) MasterMinter() (common.Address, error) {
	return _FiatTokenV22.Contract.MasterMinter(&_FiatTokenV22.CallOpts)
}

// MinterAllowance is a free data retrieval call binding the contract method 0x8a6db9c3.
//
// Solidity: function minterAllowance(address minter) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Caller) MinterAllowance(opts *bind.CallOpts, minter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "minterAllowance", minter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinterAllowance is a free data retrieval call binding the contract method 0x8a6db9c3.
//
// Solidity: function minterAllowance(address minter) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Session) MinterAllowance(minter common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.MinterAllowance(&_FiatTokenV22.CallOpts, minter)
}

// MinterAllowance is a free data retrieval call binding the contract method 0x8a6db9c3.
//
// Solidity: function minterAllowance(address minter) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22CallerSession) MinterAllowance(minter common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.MinterAllowance(&_FiatTokenV22.CallOpts, minter)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FiatTokenV22 *FiatTokenV22Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FiatTokenV22 *FiatTokenV22Session) Name() (string, error) {
	return _FiatTokenV22.Contract.Name(&_FiatTokenV22.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Name() (string, error) {
	return _FiatTokenV22.Contract.Name(&_FiatTokenV22.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Caller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Session) Nonces(owner common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.Nonces(&_FiatTokenV22.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _FiatTokenV22.Contract.Nonces(&_FiatTokenV22.CallOpts, owner)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Session) Owner() (common.Address, error) {
	return _FiatTokenV22.Contract.Owner(&_FiatTokenV22.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Owner() (common.Address, error) {
	return _FiatTokenV22.Contract.Owner(&_FiatTokenV22.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) Paused() (bool, error) {
	return _FiatTokenV22.Contract.Paused(&_FiatTokenV22.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Paused() (bool, error) {
	return _FiatTokenV22.Contract.Paused(&_FiatTokenV22.CallOpts)
}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Caller) Pauser(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "pauser")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Session) Pauser() (common.Address, error) {
	return _FiatTokenV22.Contract.Pauser(&_FiatTokenV22.CallOpts)
}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Pauser() (common.Address, error) {
	return _FiatTokenV22.Contract.Pauser(&_FiatTokenV22.CallOpts)
}

// Rescuer is a free data retrieval call binding the contract method 0x38a63183.
//
// Solidity: function rescuer() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Caller) Rescuer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "rescuer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Rescuer is a free data retrieval call binding the contract method 0x38a63183.
//
// Solidity: function rescuer() view returns(address)
func (_FiatTokenV22 *FiatTokenV22Session) Rescuer() (common.Address, error) {
	return _FiatTokenV22.Contract.Rescuer(&_FiatTokenV22.CallOpts)
}

// Rescuer is a free data retrieval call binding the contract method 0x38a63183.
//
// Solidity: function rescuer() view returns(address)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Rescuer() (common.Address, error) {
	return _FiatTokenV22.Contract.Rescuer(&_FiatTokenV22.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FiatTokenV22 *FiatTokenV22Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FiatTokenV22 *FiatTokenV22Session) Symbol() (string, error) {
	return _FiatTokenV22.Contract.Symbol(&_FiatTokenV22.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Symbol() (string, error) {
	return _FiatTokenV22.Contract.Symbol(&_FiatTokenV22.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22Session) TotalSupply() (*big.Int, error) {
	return _FiatTokenV22.Contract.TotalSupply(&_FiatTokenV22.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FiatTokenV22 *FiatTokenV22CallerSession) TotalSupply() (*big.Int, error) {
	return _FiatTokenV22.Contract.TotalSupply(&_FiatTokenV22.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_FiatTokenV22 *FiatTokenV22Caller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FiatTokenV22.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_FiatTokenV22 *FiatTokenV22Session) Version() (string, error) {
	return _FiatTokenV22.Contract.Version(&_FiatTokenV22.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_FiatTokenV22 *FiatTokenV22CallerSession) Version() (string, error) {
	return _FiatTokenV22.Contract.Version(&_FiatTokenV22.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Approve(&_FiatTokenV22.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Approve(&_FiatTokenV22.TransactOpts, spender, value)
}

// Blacklist is a paid mutator transaction binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address _account) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) Blacklist(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "blacklist", _account)
}

// Blacklist is a paid mutator transaction binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address _account) returns()
func (_FiatTokenV22 *FiatTokenV22Session) Blacklist(_account common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Blacklist(&_FiatTokenV22.TransactOpts, _account)
}

// Blacklist is a paid mutator transaction binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address _account) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Blacklist(_account common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Blacklist(&_FiatTokenV22.TransactOpts, _account)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _amount) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) Burn(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "burn", _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _amount) returns()
func (_FiatTokenV22 *FiatTokenV22Session) Burn(_amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Burn(&_FiatTokenV22.TransactOpts, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _amount) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Burn(_amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Burn(&_FiatTokenV22.TransactOpts, _amount)
}

// CancelAuthorization is a paid mutator transaction binding the contract method 0x5a049a70.
//
// Solidity: function cancelAuthorization(address authorizer, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) CancelAuthorization(opts *bind.TransactOpts, authorizer common.Address, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "cancelAuthorization", authorizer, nonce, v, r, s)
}

// CancelAuthorization is a paid mutator transaction binding the contract method 0x5a049a70.
//
// Solidity: function cancelAuthorization(address authorizer, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Session) CancelAuthorization(authorizer common.Address, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.CancelAuthorization(&_FiatTokenV22.TransactOpts, authorizer, nonce, v, r, s)
}

// CancelAuthorization is a paid mutator transaction binding the contract method 0x5a049a70.
//
// Solidity: function cancelAuthorization(address authorizer, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) CancelAuthorization(authorizer common.Address, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.CancelAuthorization(&_FiatTokenV22.TransactOpts, authorizer, nonce, v, r, s)
}

// CancelAuthorization0 is a paid mutator transaction binding the contract method 0xb7b72899.
//
// Solidity: function cancelAuthorization(address authorizer, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) CancelAuthorization0(opts *bind.TransactOpts, authorizer common.Address, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "cancelAuthorization0", authorizer, nonce, signature)
}

// CancelAuthorization0 is a paid mutator transaction binding the contract method 0xb7b72899.
//
// Solidity: function cancelAuthorization(address authorizer, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Session) CancelAuthorization0(authorizer common.Address, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.CancelAuthorization0(&_FiatTokenV22.TransactOpts, authorizer, nonce, signature)
}

// CancelAuthorization0 is a paid mutator transaction binding the contract method 0xb7b72899.
//
// Solidity: function cancelAuthorization(address authorizer, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) CancelAuthorization0(authorizer common.Address, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.CancelAuthorization0(&_FiatTokenV22.TransactOpts, authorizer, nonce, signature)
}

// ConfigureMinter is a paid mutator transaction binding the contract method 0x4e44d956.
//
// Solidity: function configureMinter(address minter, uint256 minterAllowedAmount) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) ConfigureMinter(opts *bind.TransactOpts, minter common.Address, minterAllowedAmount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "configureMinter", minter, minterAllowedAmount)
}

// ConfigureMinter is a paid mutator transaction binding the contract method 0x4e44d956.
//
// Solidity: function configureMinter(address minter, uint256 minterAllowedAmount) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) ConfigureMinter(minter common.Address, minterAllowedAmount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.ConfigureMinter(&_FiatTokenV22.TransactOpts, minter, minterAllowedAmount)
}

// ConfigureMinter is a paid mutator transaction binding the contract method 0x4e44d956.
//
// Solidity: function configureMinter(address minter, uint256 minterAllowedAmount) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) ConfigureMinter(minter common.Address, minterAllowedAmount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.ConfigureMinter(&_FiatTokenV22.TransactOpts, minter, minterAllowedAmount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 decrement) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, decrement *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "decreaseAllowance", spender, decrement)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 decrement) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) DecreaseAllowance(spender common.Address, decrement *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.DecreaseAllowance(&_FiatTokenV22.TransactOpts, spender, decrement)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 decrement) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) DecreaseAllowance(spender common.Address, decrement *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.DecreaseAllowance(&_FiatTokenV22.TransactOpts, spender, decrement)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 increment) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, increment *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "increaseAllowance", spender, increment)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 increment) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) IncreaseAllowance(spender common.Address, increment *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.IncreaseAllowance(&_FiatTokenV22.TransactOpts, spender, increment)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 increment) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) IncreaseAllowance(spender common.Address, increment *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.IncreaseAllowance(&_FiatTokenV22.TransactOpts, spender, increment)
}

// Initialize is a paid mutator transaction binding the contract method 0x3357162b.
//
// Solidity: function initialize(string tokenName, string tokenSymbol, string tokenCurrency, uint8 tokenDecimals, address newMasterMinter, address newPauser, address newBlacklister, address newOwner) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) Initialize(opts *bind.TransactOpts, tokenName string, tokenSymbol string, tokenCurrency string, tokenDecimals uint8, newMasterMinter common.Address, newPauser common.Address, newBlacklister common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "initialize", tokenName, tokenSymbol, tokenCurrency, tokenDecimals, newMasterMinter, newPauser, newBlacklister, newOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0x3357162b.
//
// Solidity: function initialize(string tokenName, string tokenSymbol, string tokenCurrency, uint8 tokenDecimals, address newMasterMinter, address newPauser, address newBlacklister, address newOwner) returns()
func (_FiatTokenV22 *FiatTokenV22Session) Initialize(tokenName string, tokenSymbol string, tokenCurrency string, tokenDecimals uint8, newMasterMinter common.Address, newPauser common.Address, newBlacklister common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Initialize(&_FiatTokenV22.TransactOpts, tokenName, tokenSymbol, tokenCurrency, tokenDecimals, newMasterMinter, newPauser, newBlacklister, newOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0x3357162b.
//
// Solidity: function initialize(string tokenName, string tokenSymbol, string tokenCurrency, uint8 tokenDecimals, address newMasterMinter, address newPauser, address newBlacklister, address newOwner) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Initialize(tokenName string, tokenSymbol string, tokenCurrency string, tokenDecimals uint8, newMasterMinter common.Address, newPauser common.Address, newBlacklister common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Initialize(&_FiatTokenV22.TransactOpts, tokenName, tokenSymbol, tokenCurrency, tokenDecimals, newMasterMinter, newPauser, newBlacklister, newOwner)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0xd608ea64.
//
// Solidity: function initializeV2(string newName) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) InitializeV2(opts *bind.TransactOpts, newName string) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "initializeV2", newName)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0xd608ea64.
//
// Solidity: function initializeV2(string newName) returns()
func (_FiatTokenV22 *FiatTokenV22Session) InitializeV2(newName string) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.InitializeV2(&_FiatTokenV22.TransactOpts, newName)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0xd608ea64.
//
// Solidity: function initializeV2(string newName) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) InitializeV2(newName string) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.InitializeV2(&_FiatTokenV22.TransactOpts, newName)
}

// InitializeV21 is a paid mutator transaction binding the contract method 0x2fc81e09.
//
// Solidity: function initializeV2_1(address lostAndFound) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) InitializeV21(opts *bind.TransactOpts, lostAndFound common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "initializeV2_1", lostAndFound)
}

// InitializeV21 is a paid mutator transaction binding the contract method 0x2fc81e09.
//
// Solidity: function initializeV2_1(address lostAndFound) returns()
func (_FiatTokenV22 *FiatTokenV22Session) InitializeV21(lostAndFound common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.InitializeV21(&_FiatTokenV22.TransactOpts, lostAndFound)
}

// InitializeV21 is a paid mutator transaction binding the contract method 0x2fc81e09.
//
// Solidity: function initializeV2_1(address lostAndFound) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) InitializeV21(lostAndFound common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.InitializeV21(&_FiatTokenV22.TransactOpts, lostAndFound)
}

// InitializeV22 is a paid mutator transaction binding the contract method 0x430239b4.
//
// Solidity: function initializeV2_2(address[] accountsToBlacklist, string newSymbol) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) InitializeV22(opts *bind.TransactOpts, accountsToBlacklist []common.Address, newSymbol string) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "initializeV2_2", accountsToBlacklist, newSymbol)
}

// InitializeV22 is a paid mutator transaction binding the contract method 0x430239b4.
//
// Solidity: function initializeV2_2(address[] accountsToBlacklist, string newSymbol) returns()
func (_FiatTokenV22 *FiatTokenV22Session) InitializeV22(accountsToBlacklist []common.Address, newSymbol string) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.InitializeV22(&_FiatTokenV22.TransactOpts, accountsToBlacklist, newSymbol)
}

// InitializeV22 is a paid mutator transaction binding the contract method 0x430239b4.
//
// Solidity: function initializeV2_2(address[] accountsToBlacklist, string newSymbol) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) InitializeV22(accountsToBlacklist []common.Address, newSymbol string) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.InitializeV22(&_FiatTokenV22.TransactOpts, accountsToBlacklist, newSymbol)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) Mint(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "mint", _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Mint(&_FiatTokenV22.TransactOpts, _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Mint(&_FiatTokenV22.TransactOpts, _to, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_FiatTokenV22 *FiatTokenV22Session) Pause() (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Pause(&_FiatTokenV22.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Pause() (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Pause(&_FiatTokenV22.TransactOpts)
}

// Permit is a paid mutator transaction binding the contract method 0x9fd5a6cf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "permit", owner, spender, value, deadline, signature)
}

// Permit is a paid mutator transaction binding the contract method 0x9fd5a6cf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Session) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Permit(&_FiatTokenV22.TransactOpts, owner, spender, value, deadline, signature)
}

// Permit is a paid mutator transaction binding the contract method 0x9fd5a6cf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Permit(&_FiatTokenV22.TransactOpts, owner, spender, value, deadline, signature)
}

// Permit0 is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) Permit0(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "permit0", owner, spender, value, deadline, v, r, s)
}

// Permit0 is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Session) Permit0(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Permit0(&_FiatTokenV22.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit0 is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Permit0(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Permit0(&_FiatTokenV22.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// ReceiveWithAuthorization is a paid mutator transaction binding the contract method 0x88b7ab63.
//
// Solidity: function receiveWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) ReceiveWithAuthorization(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "receiveWithAuthorization", from, to, value, validAfter, validBefore, nonce, signature)
}

// ReceiveWithAuthorization is a paid mutator transaction binding the contract method 0x88b7ab63.
//
// Solidity: function receiveWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Session) ReceiveWithAuthorization(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.ReceiveWithAuthorization(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, signature)
}

// ReceiveWithAuthorization is a paid mutator transaction binding the contract method 0x88b7ab63.
//
// Solidity: function receiveWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) ReceiveWithAuthorization(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.ReceiveWithAuthorization(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, signature)
}

// ReceiveWithAuthorization0 is a paid mutator transaction binding the contract method 0xef55bec6.
//
// Solidity: function receiveWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) ReceiveWithAuthorization0(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "receiveWithAuthorization0", from, to, value, validAfter, validBefore, nonce, v, r, s)
}

// ReceiveWithAuthorization0 is a paid mutator transaction binding the contract method 0xef55bec6.
//
// Solidity: function receiveWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Session) ReceiveWithAuthorization0(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.ReceiveWithAuthorization0(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, v, r, s)
}

// ReceiveWithAuthorization0 is a paid mutator transaction binding the contract method 0xef55bec6.
//
// Solidity: function receiveWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) ReceiveWithAuthorization0(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.ReceiveWithAuthorization0(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, v, r, s)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0x3092afd5.
//
// Solidity: function removeMinter(address minter) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) RemoveMinter(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "removeMinter", minter)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0x3092afd5.
//
// Solidity: function removeMinter(address minter) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) RemoveMinter(minter common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.RemoveMinter(&_FiatTokenV22.TransactOpts, minter)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0x3092afd5.
//
// Solidity: function removeMinter(address minter) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) RemoveMinter(minter common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.RemoveMinter(&_FiatTokenV22.TransactOpts, minter)
}

// RescueERC20 is a paid mutator transaction binding the contract method 0xb2118a8d.
//
// Solidity: function rescueERC20(address tokenContract, address to, uint256 amount) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) RescueERC20(opts *bind.TransactOpts, tokenContract common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "rescueERC20", tokenContract, to, amount)
}

// RescueERC20 is a paid mutator transaction binding the contract method 0xb2118a8d.
//
// Solidity: function rescueERC20(address tokenContract, address to, uint256 amount) returns()
func (_FiatTokenV22 *FiatTokenV22Session) RescueERC20(tokenContract common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.RescueERC20(&_FiatTokenV22.TransactOpts, tokenContract, to, amount)
}

// RescueERC20 is a paid mutator transaction binding the contract method 0xb2118a8d.
//
// Solidity: function rescueERC20(address tokenContract, address to, uint256 amount) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) RescueERC20(tokenContract common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.RescueERC20(&_FiatTokenV22.TransactOpts, tokenContract, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Transfer(&_FiatTokenV22.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Transfer(&_FiatTokenV22.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferFrom(&_FiatTokenV22.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FiatTokenV22 *FiatTokenV22TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferFrom(&_FiatTokenV22.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FiatTokenV22 *FiatTokenV22Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferOwnership(&_FiatTokenV22.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferOwnership(&_FiatTokenV22.TransactOpts, newOwner)
}

// TransferWithAuthorization is a paid mutator transaction binding the contract method 0xcf092995.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) TransferWithAuthorization(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "transferWithAuthorization", from, to, value, validAfter, validBefore, nonce, signature)
}

// TransferWithAuthorization is a paid mutator transaction binding the contract method 0xcf092995.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22Session) TransferWithAuthorization(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferWithAuthorization(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, signature)
}

// TransferWithAuthorization is a paid mutator transaction binding the contract method 0xcf092995.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) TransferWithAuthorization(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferWithAuthorization(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, signature)
}

// TransferWithAuthorization0 is a paid mutator transaction binding the contract method 0xe3ee160e.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) TransferWithAuthorization0(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "transferWithAuthorization0", from, to, value, validAfter, validBefore, nonce, v, r, s)
}

// TransferWithAuthorization0 is a paid mutator transaction binding the contract method 0xe3ee160e.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22Session) TransferWithAuthorization0(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferWithAuthorization0(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, v, r, s)
}

// TransferWithAuthorization0 is a paid mutator transaction binding the contract method 0xe3ee160e.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, uint8 v, bytes32 r, bytes32 s) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) TransferWithAuthorization0(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.TransferWithAuthorization0(&_FiatTokenV22.TransactOpts, from, to, value, validAfter, validBefore, nonce, v, r, s)
}

// UnBlacklist is a paid mutator transaction binding the contract method 0x1a895266.
//
// Solidity: function unBlacklist(address _account) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) UnBlacklist(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "unBlacklist", _account)
}

// UnBlacklist is a paid mutator transaction binding the contract method 0x1a895266.
//
// Solidity: function unBlacklist(address _account) returns()
func (_FiatTokenV22 *FiatTokenV22Session) UnBlacklist(_account common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UnBlacklist(&_FiatTokenV22.TransactOpts, _account)
}

// UnBlacklist is a paid mutator transaction binding the contract method 0x1a895266.
//
// Solidity: function unBlacklist(address _account) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) UnBlacklist(_account common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UnBlacklist(&_FiatTokenV22.TransactOpts, _account)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_FiatTokenV22 *FiatTokenV22Session) Unpause() (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Unpause(&_FiatTokenV22.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) Unpause() (*types.Transaction, error) {
	return _FiatTokenV22.Contract.Unpause(&_FiatTokenV22.TransactOpts)
}

// UpdateBlacklister is a paid mutator transaction binding the contract method 0xad38bf22.
//
// Solidity: function updateBlacklister(address _newBlacklister) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) UpdateBlacklister(opts *bind.TransactOpts, _newBlacklister common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "updateBlacklister", _newBlacklister)
}

// UpdateBlacklister is a paid mutator transaction binding the contract method 0xad38bf22.
//
// Solidity: function updateBlacklister(address _newBlacklister) returns()
func (_FiatTokenV22 *FiatTokenV22Session) UpdateBlacklister(_newBlacklister common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdateBlacklister(&_FiatTokenV22.TransactOpts, _newBlacklister)
}

// UpdateBlacklister is a paid mutator transaction binding the contract method 0xad38bf22.
//
// Solidity: function updateBlacklister(address _newBlacklister) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) UpdateBlacklister(_newBlacklister common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdateBlacklister(&_FiatTokenV22.TransactOpts, _newBlacklister)
}

// UpdateMasterMinter is a paid mutator transaction binding the contract method 0xaa20e1e4.
//
// Solidity: function updateMasterMinter(address _newMasterMinter) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) UpdateMasterMinter(opts *bind.TransactOpts, _newMasterMinter common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "updateMasterMinter", _newMasterMinter)
}

// UpdateMasterMinter is a paid mutator transaction binding the contract method 0xaa20e1e4.
//
// Solidity: function updateMasterMinter(address _newMasterMinter) returns()
func (_FiatTokenV22 *FiatTokenV22Session) UpdateMasterMinter(_newMasterMinter common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdateMasterMinter(&_FiatTokenV22.TransactOpts, _newMasterMinter)
}

// UpdateMasterMinter is a paid mutator transaction binding the contract method 0xaa20e1e4.
//
// Solidity: function updateMasterMinter(address _newMasterMinter) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) UpdateMasterMinter(_newMasterMinter common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdateMasterMinter(&_FiatTokenV22.TransactOpts, _newMasterMinter)
}

// UpdatePauser is a paid mutator transaction binding the contract method 0x554bab3c.
//
// Solidity: function updatePauser(address _newPauser) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) UpdatePauser(opts *bind.TransactOpts, _newPauser common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "updatePauser", _newPauser)
}

// UpdatePauser is a paid mutator transaction binding the contract method 0x554bab3c.
//
// Solidity: function updatePauser(address _newPauser) returns()
func (_FiatTokenV22 *FiatTokenV22Session) UpdatePauser(_newPauser common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdatePauser(&_FiatTokenV22.TransactOpts, _newPauser)
}

// UpdatePauser is a paid mutator transaction binding the contract method 0x554bab3c.
//
// Solidity: function updatePauser(address _newPauser) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) UpdatePauser(_newPauser common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdatePauser(&_FiatTokenV22.TransactOpts, _newPauser)
}

// UpdateRescuer is a paid mutator transaction binding the contract method 0x2ab60045.
//
// Solidity: function updateRescuer(address newRescuer) returns()
func (_FiatTokenV22 *FiatTokenV22Transactor) UpdateRescuer(opts *bind.TransactOpts, newRescuer common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.contract.Transact(opts, "updateRescuer", newRescuer)
}

// UpdateRescuer is a paid mutator transaction binding the contract method 0x2ab60045.
//
// Solidity: function updateRescuer(address newRescuer) returns()
func (_FiatTokenV22 *FiatTokenV22Session) UpdateRescuer(newRescuer common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdateRescuer(&_FiatTokenV22.TransactOpts, newRescuer)
}

// UpdateRescuer is a paid mutator transaction binding the contract method 0x2ab60045.
//
// Solidity: function updateRescuer(address newRescuer) returns()
func (_FiatTokenV22 *FiatTokenV22TransactorSession) UpdateRescuer(newRescuer common.Address) (*types.Transaction, error) {
	return _FiatTokenV22.Contract.UpdateRescuer(&_FiatTokenV22.TransactOpts, newRescuer)
}

// FiatTokenV22ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the FiatTokenV22 contract.
type FiatTokenV22ApprovalIterator struct {
	Event *FiatTokenV22Approval // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22Approval)
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
		it.Event = new(FiatTokenV22Approval)
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
func (it *FiatTokenV22ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22Approval represents a Approval event raised by the FiatTokenV22 contract.
type FiatTokenV22Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*FiatTokenV22ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22ApprovalIterator{contract: _FiatTokenV22.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *FiatTokenV22Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22Approval)
				if err := _FiatTokenV22.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseApproval(log types.Log) (*FiatTokenV22Approval, error) {
	event := new(FiatTokenV22Approval)
	if err := _FiatTokenV22.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22AuthorizationCanceledIterator is returned from FilterAuthorizationCanceled and is used to iterate over the raw logs and unpacked data for AuthorizationCanceled events raised by the FiatTokenV22 contract.
type FiatTokenV22AuthorizationCanceledIterator struct {
	Event *FiatTokenV22AuthorizationCanceled // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22AuthorizationCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22AuthorizationCanceled)
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
		it.Event = new(FiatTokenV22AuthorizationCanceled)
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
func (it *FiatTokenV22AuthorizationCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22AuthorizationCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22AuthorizationCanceled represents a AuthorizationCanceled event raised by the FiatTokenV22 contract.
type FiatTokenV22AuthorizationCanceled struct {
	Authorizer common.Address
	Nonce      [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAuthorizationCanceled is a free log retrieval operation binding the contract event 0x1cdd46ff242716cdaa72d159d339a485b3438398348d68f09d7c8c0a59353d81.
//
// Solidity: event AuthorizationCanceled(address indexed authorizer, bytes32 indexed nonce)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterAuthorizationCanceled(opts *bind.FilterOpts, authorizer []common.Address, nonce [][32]byte) (*FiatTokenV22AuthorizationCanceledIterator, error) {

	var authorizerRule []interface{}
	for _, authorizerItem := range authorizer {
		authorizerRule = append(authorizerRule, authorizerItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "AuthorizationCanceled", authorizerRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22AuthorizationCanceledIterator{contract: _FiatTokenV22.contract, event: "AuthorizationCanceled", logs: logs, sub: sub}, nil
}

// WatchAuthorizationCanceled is a free log subscription operation binding the contract event 0x1cdd46ff242716cdaa72d159d339a485b3438398348d68f09d7c8c0a59353d81.
//
// Solidity: event AuthorizationCanceled(address indexed authorizer, bytes32 indexed nonce)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchAuthorizationCanceled(opts *bind.WatchOpts, sink chan<- *FiatTokenV22AuthorizationCanceled, authorizer []common.Address, nonce [][32]byte) (event.Subscription, error) {

	var authorizerRule []interface{}
	for _, authorizerItem := range authorizer {
		authorizerRule = append(authorizerRule, authorizerItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "AuthorizationCanceled", authorizerRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22AuthorizationCanceled)
				if err := _FiatTokenV22.contract.UnpackLog(event, "AuthorizationCanceled", log); err != nil {
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

// ParseAuthorizationCanceled is a log parse operation binding the contract event 0x1cdd46ff242716cdaa72d159d339a485b3438398348d68f09d7c8c0a59353d81.
//
// Solidity: event AuthorizationCanceled(address indexed authorizer, bytes32 indexed nonce)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseAuthorizationCanceled(log types.Log) (*FiatTokenV22AuthorizationCanceled, error) {
	event := new(FiatTokenV22AuthorizationCanceled)
	if err := _FiatTokenV22.contract.UnpackLog(event, "AuthorizationCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22AuthorizationUsedIterator is returned from FilterAuthorizationUsed and is used to iterate over the raw logs and unpacked data for AuthorizationUsed events raised by the FiatTokenV22 contract.
type FiatTokenV22AuthorizationUsedIterator struct {
	Event *FiatTokenV22AuthorizationUsed // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22AuthorizationUsedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22AuthorizationUsed)
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
		it.Event = new(FiatTokenV22AuthorizationUsed)
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
func (it *FiatTokenV22AuthorizationUsedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22AuthorizationUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22AuthorizationUsed represents a AuthorizationUsed event raised by the FiatTokenV22 contract.
type FiatTokenV22AuthorizationUsed struct {
	Authorizer common.Address
	Nonce      [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAuthorizationUsed is a free log retrieval operation binding the contract event 0x98de503528ee59b575ef0c0a2576a82497bfc029a5685b209e9ec333479b10a5.
//
// Solidity: event AuthorizationUsed(address indexed authorizer, bytes32 indexed nonce)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterAuthorizationUsed(opts *bind.FilterOpts, authorizer []common.Address, nonce [][32]byte) (*FiatTokenV22AuthorizationUsedIterator, error) {

	var authorizerRule []interface{}
	for _, authorizerItem := range authorizer {
		authorizerRule = append(authorizerRule, authorizerItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "AuthorizationUsed", authorizerRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22AuthorizationUsedIterator{contract: _FiatTokenV22.contract, event: "AuthorizationUsed", logs: logs, sub: sub}, nil
}

// WatchAuthorizationUsed is a free log subscription operation binding the contract event 0x98de503528ee59b575ef0c0a2576a82497bfc029a5685b209e9ec333479b10a5.
//
// Solidity: event AuthorizationUsed(address indexed authorizer, bytes32 indexed nonce)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchAuthorizationUsed(opts *bind.WatchOpts, sink chan<- *FiatTokenV22AuthorizationUsed, authorizer []common.Address, nonce [][32]byte) (event.Subscription, error) {

	var authorizerRule []interface{}
	for _, authorizerItem := range authorizer {
		authorizerRule = append(authorizerRule, authorizerItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "AuthorizationUsed", authorizerRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22AuthorizationUsed)
				if err := _FiatTokenV22.contract.UnpackLog(event, "AuthorizationUsed", log); err != nil {
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

// ParseAuthorizationUsed is a log parse operation binding the contract event 0x98de503528ee59b575ef0c0a2576a82497bfc029a5685b209e9ec333479b10a5.
//
// Solidity: event AuthorizationUsed(address indexed authorizer, bytes32 indexed nonce)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseAuthorizationUsed(log types.Log) (*FiatTokenV22AuthorizationUsed, error) {
	event := new(FiatTokenV22AuthorizationUsed)
	if err := _FiatTokenV22.contract.UnpackLog(event, "AuthorizationUsed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22BlacklistedIterator is returned from FilterBlacklisted and is used to iterate over the raw logs and unpacked data for Blacklisted events raised by the FiatTokenV22 contract.
type FiatTokenV22BlacklistedIterator struct {
	Event *FiatTokenV22Blacklisted // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22BlacklistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22Blacklisted)
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
		it.Event = new(FiatTokenV22Blacklisted)
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
func (it *FiatTokenV22BlacklistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22BlacklistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22Blacklisted represents a Blacklisted event raised by the FiatTokenV22 contract.
type FiatTokenV22Blacklisted struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBlacklisted is a free log retrieval operation binding the contract event 0xffa4e6181777692565cf28528fc88fd1516ea86b56da075235fa575af6a4b855.
//
// Solidity: event Blacklisted(address indexed _account)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterBlacklisted(opts *bind.FilterOpts, _account []common.Address) (*FiatTokenV22BlacklistedIterator, error) {

	var _accountRule []interface{}
	for _, _accountItem := range _account {
		_accountRule = append(_accountRule, _accountItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "Blacklisted", _accountRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22BlacklistedIterator{contract: _FiatTokenV22.contract, event: "Blacklisted", logs: logs, sub: sub}, nil
}

// WatchBlacklisted is a free log subscription operation binding the contract event 0xffa4e6181777692565cf28528fc88fd1516ea86b56da075235fa575af6a4b855.
//
// Solidity: event Blacklisted(address indexed _account)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchBlacklisted(opts *bind.WatchOpts, sink chan<- *FiatTokenV22Blacklisted, _account []common.Address) (event.Subscription, error) {

	var _accountRule []interface{}
	for _, _accountItem := range _account {
		_accountRule = append(_accountRule, _accountItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "Blacklisted", _accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22Blacklisted)
				if err := _FiatTokenV22.contract.UnpackLog(event, "Blacklisted", log); err != nil {
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

// ParseBlacklisted is a log parse operation binding the contract event 0xffa4e6181777692565cf28528fc88fd1516ea86b56da075235fa575af6a4b855.
//
// Solidity: event Blacklisted(address indexed _account)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseBlacklisted(log types.Log) (*FiatTokenV22Blacklisted, error) {
	event := new(FiatTokenV22Blacklisted)
	if err := _FiatTokenV22.contract.UnpackLog(event, "Blacklisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22BlacklisterChangedIterator is returned from FilterBlacklisterChanged and is used to iterate over the raw logs and unpacked data for BlacklisterChanged events raised by the FiatTokenV22 contract.
type FiatTokenV22BlacklisterChangedIterator struct {
	Event *FiatTokenV22BlacklisterChanged // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22BlacklisterChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22BlacklisterChanged)
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
		it.Event = new(FiatTokenV22BlacklisterChanged)
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
func (it *FiatTokenV22BlacklisterChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22BlacklisterChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22BlacklisterChanged represents a BlacklisterChanged event raised by the FiatTokenV22 contract.
type FiatTokenV22BlacklisterChanged struct {
	NewBlacklister common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBlacklisterChanged is a free log retrieval operation binding the contract event 0xc67398012c111ce95ecb7429b933096c977380ee6c421175a71a4a4c6c88c06e.
//
// Solidity: event BlacklisterChanged(address indexed newBlacklister)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterBlacklisterChanged(opts *bind.FilterOpts, newBlacklister []common.Address) (*FiatTokenV22BlacklisterChangedIterator, error) {

	var newBlacklisterRule []interface{}
	for _, newBlacklisterItem := range newBlacklister {
		newBlacklisterRule = append(newBlacklisterRule, newBlacklisterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "BlacklisterChanged", newBlacklisterRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22BlacklisterChangedIterator{contract: _FiatTokenV22.contract, event: "BlacklisterChanged", logs: logs, sub: sub}, nil
}

// WatchBlacklisterChanged is a free log subscription operation binding the contract event 0xc67398012c111ce95ecb7429b933096c977380ee6c421175a71a4a4c6c88c06e.
//
// Solidity: event BlacklisterChanged(address indexed newBlacklister)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchBlacklisterChanged(opts *bind.WatchOpts, sink chan<- *FiatTokenV22BlacklisterChanged, newBlacklister []common.Address) (event.Subscription, error) {

	var newBlacklisterRule []interface{}
	for _, newBlacklisterItem := range newBlacklister {
		newBlacklisterRule = append(newBlacklisterRule, newBlacklisterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "BlacklisterChanged", newBlacklisterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22BlacklisterChanged)
				if err := _FiatTokenV22.contract.UnpackLog(event, "BlacklisterChanged", log); err != nil {
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

// ParseBlacklisterChanged is a log parse operation binding the contract event 0xc67398012c111ce95ecb7429b933096c977380ee6c421175a71a4a4c6c88c06e.
//
// Solidity: event BlacklisterChanged(address indexed newBlacklister)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseBlacklisterChanged(log types.Log) (*FiatTokenV22BlacklisterChanged, error) {
	event := new(FiatTokenV22BlacklisterChanged)
	if err := _FiatTokenV22.contract.UnpackLog(event, "BlacklisterChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22BurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the FiatTokenV22 contract.
type FiatTokenV22BurnIterator struct {
	Event *FiatTokenV22Burn // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22BurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22Burn)
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
		it.Event = new(FiatTokenV22Burn)
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
func (it *FiatTokenV22BurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22BurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22Burn represents a Burn event raised by the FiatTokenV22 contract.
type FiatTokenV22Burn struct {
	Burner common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed burner, uint256 amount)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterBurn(opts *bind.FilterOpts, burner []common.Address) (*FiatTokenV22BurnIterator, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "Burn", burnerRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22BurnIterator{contract: _FiatTokenV22.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed burner, uint256 amount)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *FiatTokenV22Burn, burner []common.Address) (event.Subscription, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "Burn", burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22Burn)
				if err := _FiatTokenV22.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed burner, uint256 amount)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseBurn(log types.Log) (*FiatTokenV22Burn, error) {
	event := new(FiatTokenV22Burn)
	if err := _FiatTokenV22.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22MasterMinterChangedIterator is returned from FilterMasterMinterChanged and is used to iterate over the raw logs and unpacked data for MasterMinterChanged events raised by the FiatTokenV22 contract.
type FiatTokenV22MasterMinterChangedIterator struct {
	Event *FiatTokenV22MasterMinterChanged // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22MasterMinterChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22MasterMinterChanged)
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
		it.Event = new(FiatTokenV22MasterMinterChanged)
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
func (it *FiatTokenV22MasterMinterChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22MasterMinterChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22MasterMinterChanged represents a MasterMinterChanged event raised by the FiatTokenV22 contract.
type FiatTokenV22MasterMinterChanged struct {
	NewMasterMinter common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMasterMinterChanged is a free log retrieval operation binding the contract event 0xdb66dfa9c6b8f5226fe9aac7e51897ae8ee94ac31dc70bb6c9900b2574b707e6.
//
// Solidity: event MasterMinterChanged(address indexed newMasterMinter)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterMasterMinterChanged(opts *bind.FilterOpts, newMasterMinter []common.Address) (*FiatTokenV22MasterMinterChangedIterator, error) {

	var newMasterMinterRule []interface{}
	for _, newMasterMinterItem := range newMasterMinter {
		newMasterMinterRule = append(newMasterMinterRule, newMasterMinterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "MasterMinterChanged", newMasterMinterRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22MasterMinterChangedIterator{contract: _FiatTokenV22.contract, event: "MasterMinterChanged", logs: logs, sub: sub}, nil
}

// WatchMasterMinterChanged is a free log subscription operation binding the contract event 0xdb66dfa9c6b8f5226fe9aac7e51897ae8ee94ac31dc70bb6c9900b2574b707e6.
//
// Solidity: event MasterMinterChanged(address indexed newMasterMinter)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchMasterMinterChanged(opts *bind.WatchOpts, sink chan<- *FiatTokenV22MasterMinterChanged, newMasterMinter []common.Address) (event.Subscription, error) {

	var newMasterMinterRule []interface{}
	for _, newMasterMinterItem := range newMasterMinter {
		newMasterMinterRule = append(newMasterMinterRule, newMasterMinterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "MasterMinterChanged", newMasterMinterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22MasterMinterChanged)
				if err := _FiatTokenV22.contract.UnpackLog(event, "MasterMinterChanged", log); err != nil {
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

// ParseMasterMinterChanged is a log parse operation binding the contract event 0xdb66dfa9c6b8f5226fe9aac7e51897ae8ee94ac31dc70bb6c9900b2574b707e6.
//
// Solidity: event MasterMinterChanged(address indexed newMasterMinter)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseMasterMinterChanged(log types.Log) (*FiatTokenV22MasterMinterChanged, error) {
	event := new(FiatTokenV22MasterMinterChanged)
	if err := _FiatTokenV22.contract.UnpackLog(event, "MasterMinterChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22MintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the FiatTokenV22 contract.
type FiatTokenV22MintIterator struct {
	Event *FiatTokenV22Mint // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22MintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22Mint)
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
		it.Event = new(FiatTokenV22Mint)
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
func (it *FiatTokenV22MintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22MintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22Mint represents a Mint event raised by the FiatTokenV22 contract.
type FiatTokenV22Mint struct {
	Minter common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0xab8530f87dc9b59234c4623bf917212bb2536d647574c8e7e5da92c2ede0c9f8.
//
// Solidity: event Mint(address indexed minter, address indexed to, uint256 amount)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterMint(opts *bind.FilterOpts, minter []common.Address, to []common.Address) (*FiatTokenV22MintIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "Mint", minterRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22MintIterator{contract: _FiatTokenV22.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0xab8530f87dc9b59234c4623bf917212bb2536d647574c8e7e5da92c2ede0c9f8.
//
// Solidity: event Mint(address indexed minter, address indexed to, uint256 amount)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchMint(opts *bind.WatchOpts, sink chan<- *FiatTokenV22Mint, minter []common.Address, to []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "Mint", minterRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22Mint)
				if err := _FiatTokenV22.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0xab8530f87dc9b59234c4623bf917212bb2536d647574c8e7e5da92c2ede0c9f8.
//
// Solidity: event Mint(address indexed minter, address indexed to, uint256 amount)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseMint(log types.Log) (*FiatTokenV22Mint, error) {
	event := new(FiatTokenV22Mint)
	if err := _FiatTokenV22.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22MinterConfiguredIterator is returned from FilterMinterConfigured and is used to iterate over the raw logs and unpacked data for MinterConfigured events raised by the FiatTokenV22 contract.
type FiatTokenV22MinterConfiguredIterator struct {
	Event *FiatTokenV22MinterConfigured // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22MinterConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22MinterConfigured)
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
		it.Event = new(FiatTokenV22MinterConfigured)
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
func (it *FiatTokenV22MinterConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22MinterConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22MinterConfigured represents a MinterConfigured event raised by the FiatTokenV22 contract.
type FiatTokenV22MinterConfigured struct {
	Minter              common.Address
	MinterAllowedAmount *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterMinterConfigured is a free log retrieval operation binding the contract event 0x46980fca912ef9bcdbd36877427b6b90e860769f604e89c0e67720cece530d20.
//
// Solidity: event MinterConfigured(address indexed minter, uint256 minterAllowedAmount)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterMinterConfigured(opts *bind.FilterOpts, minter []common.Address) (*FiatTokenV22MinterConfiguredIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "MinterConfigured", minterRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22MinterConfiguredIterator{contract: _FiatTokenV22.contract, event: "MinterConfigured", logs: logs, sub: sub}, nil
}

// WatchMinterConfigured is a free log subscription operation binding the contract event 0x46980fca912ef9bcdbd36877427b6b90e860769f604e89c0e67720cece530d20.
//
// Solidity: event MinterConfigured(address indexed minter, uint256 minterAllowedAmount)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchMinterConfigured(opts *bind.WatchOpts, sink chan<- *FiatTokenV22MinterConfigured, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "MinterConfigured", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22MinterConfigured)
				if err := _FiatTokenV22.contract.UnpackLog(event, "MinterConfigured", log); err != nil {
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

// ParseMinterConfigured is a log parse operation binding the contract event 0x46980fca912ef9bcdbd36877427b6b90e860769f604e89c0e67720cece530d20.
//
// Solidity: event MinterConfigured(address indexed minter, uint256 minterAllowedAmount)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseMinterConfigured(log types.Log) (*FiatTokenV22MinterConfigured, error) {
	event := new(FiatTokenV22MinterConfigured)
	if err := _FiatTokenV22.contract.UnpackLog(event, "MinterConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22MinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the FiatTokenV22 contract.
type FiatTokenV22MinterRemovedIterator struct {
	Event *FiatTokenV22MinterRemoved // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22MinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22MinterRemoved)
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
		it.Event = new(FiatTokenV22MinterRemoved)
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
func (it *FiatTokenV22MinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22MinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22MinterRemoved represents a MinterRemoved event raised by the FiatTokenV22 contract.
type FiatTokenV22MinterRemoved struct {
	OldMinter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed oldMinter)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterMinterRemoved(opts *bind.FilterOpts, oldMinter []common.Address) (*FiatTokenV22MinterRemovedIterator, error) {

	var oldMinterRule []interface{}
	for _, oldMinterItem := range oldMinter {
		oldMinterRule = append(oldMinterRule, oldMinterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "MinterRemoved", oldMinterRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22MinterRemovedIterator{contract: _FiatTokenV22.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed oldMinter)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *FiatTokenV22MinterRemoved, oldMinter []common.Address) (event.Subscription, error) {

	var oldMinterRule []interface{}
	for _, oldMinterItem := range oldMinter {
		oldMinterRule = append(oldMinterRule, oldMinterItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "MinterRemoved", oldMinterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22MinterRemoved)
				if err := _FiatTokenV22.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// ParseMinterRemoved is a log parse operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed oldMinter)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseMinterRemoved(log types.Log) (*FiatTokenV22MinterRemoved, error) {
	event := new(FiatTokenV22MinterRemoved)
	if err := _FiatTokenV22.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FiatTokenV22 contract.
type FiatTokenV22OwnershipTransferredIterator struct {
	Event *FiatTokenV22OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22OwnershipTransferred)
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
		it.Event = new(FiatTokenV22OwnershipTransferred)
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
func (it *FiatTokenV22OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22OwnershipTransferred represents a OwnershipTransferred event raised by the FiatTokenV22 contract.
type FiatTokenV22OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address previousOwner, address newOwner)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts) (*FiatTokenV22OwnershipTransferredIterator, error) {

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "OwnershipTransferred")
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22OwnershipTransferredIterator{contract: _FiatTokenV22.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address previousOwner, address newOwner)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FiatTokenV22OwnershipTransferred) (event.Subscription, error) {

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "OwnershipTransferred")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22OwnershipTransferred)
				if err := _FiatTokenV22.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseOwnershipTransferred(log types.Log) (*FiatTokenV22OwnershipTransferred, error) {
	event := new(FiatTokenV22OwnershipTransferred)
	if err := _FiatTokenV22.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22PauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the FiatTokenV22 contract.
type FiatTokenV22PauseIterator struct {
	Event *FiatTokenV22Pause // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22PauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22Pause)
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
		it.Event = new(FiatTokenV22Pause)
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
func (it *FiatTokenV22PauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22PauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22Pause represents a Pause event raised by the FiatTokenV22 contract.
type FiatTokenV22Pause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterPause(opts *bind.FilterOpts) (*FiatTokenV22PauseIterator, error) {

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22PauseIterator{contract: _FiatTokenV22.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchPause(opts *bind.WatchOpts, sink chan<- *FiatTokenV22Pause) (event.Subscription, error) {

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22Pause)
				if err := _FiatTokenV22.contract.UnpackLog(event, "Pause", log); err != nil {
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

// ParsePause is a log parse operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_FiatTokenV22 *FiatTokenV22Filterer) ParsePause(log types.Log) (*FiatTokenV22Pause, error) {
	event := new(FiatTokenV22Pause)
	if err := _FiatTokenV22.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22PauserChangedIterator is returned from FilterPauserChanged and is used to iterate over the raw logs and unpacked data for PauserChanged events raised by the FiatTokenV22 contract.
type FiatTokenV22PauserChangedIterator struct {
	Event *FiatTokenV22PauserChanged // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22PauserChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22PauserChanged)
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
		it.Event = new(FiatTokenV22PauserChanged)
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
func (it *FiatTokenV22PauserChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22PauserChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22PauserChanged represents a PauserChanged event raised by the FiatTokenV22 contract.
type FiatTokenV22PauserChanged struct {
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPauserChanged is a free log retrieval operation binding the contract event 0xb80482a293ca2e013eda8683c9bd7fc8347cfdaeea5ede58cba46df502c2a604.
//
// Solidity: event PauserChanged(address indexed newAddress)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterPauserChanged(opts *bind.FilterOpts, newAddress []common.Address) (*FiatTokenV22PauserChangedIterator, error) {

	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "PauserChanged", newAddressRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22PauserChangedIterator{contract: _FiatTokenV22.contract, event: "PauserChanged", logs: logs, sub: sub}, nil
}

// WatchPauserChanged is a free log subscription operation binding the contract event 0xb80482a293ca2e013eda8683c9bd7fc8347cfdaeea5ede58cba46df502c2a604.
//
// Solidity: event PauserChanged(address indexed newAddress)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchPauserChanged(opts *bind.WatchOpts, sink chan<- *FiatTokenV22PauserChanged, newAddress []common.Address) (event.Subscription, error) {

	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "PauserChanged", newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22PauserChanged)
				if err := _FiatTokenV22.contract.UnpackLog(event, "PauserChanged", log); err != nil {
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

// ParsePauserChanged is a log parse operation binding the contract event 0xb80482a293ca2e013eda8683c9bd7fc8347cfdaeea5ede58cba46df502c2a604.
//
// Solidity: event PauserChanged(address indexed newAddress)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParsePauserChanged(log types.Log) (*FiatTokenV22PauserChanged, error) {
	event := new(FiatTokenV22PauserChanged)
	if err := _FiatTokenV22.contract.UnpackLog(event, "PauserChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22RescuerChangedIterator is returned from FilterRescuerChanged and is used to iterate over the raw logs and unpacked data for RescuerChanged events raised by the FiatTokenV22 contract.
type FiatTokenV22RescuerChangedIterator struct {
	Event *FiatTokenV22RescuerChanged // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22RescuerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22RescuerChanged)
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
		it.Event = new(FiatTokenV22RescuerChanged)
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
func (it *FiatTokenV22RescuerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22RescuerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22RescuerChanged represents a RescuerChanged event raised by the FiatTokenV22 contract.
type FiatTokenV22RescuerChanged struct {
	NewRescuer common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRescuerChanged is a free log retrieval operation binding the contract event 0xe475e580d85111348e40d8ca33cfdd74c30fe1655c2d8537a13abc10065ffa5a.
//
// Solidity: event RescuerChanged(address indexed newRescuer)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterRescuerChanged(opts *bind.FilterOpts, newRescuer []common.Address) (*FiatTokenV22RescuerChangedIterator, error) {

	var newRescuerRule []interface{}
	for _, newRescuerItem := range newRescuer {
		newRescuerRule = append(newRescuerRule, newRescuerItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "RescuerChanged", newRescuerRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22RescuerChangedIterator{contract: _FiatTokenV22.contract, event: "RescuerChanged", logs: logs, sub: sub}, nil
}

// WatchRescuerChanged is a free log subscription operation binding the contract event 0xe475e580d85111348e40d8ca33cfdd74c30fe1655c2d8537a13abc10065ffa5a.
//
// Solidity: event RescuerChanged(address indexed newRescuer)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchRescuerChanged(opts *bind.WatchOpts, sink chan<- *FiatTokenV22RescuerChanged, newRescuer []common.Address) (event.Subscription, error) {

	var newRescuerRule []interface{}
	for _, newRescuerItem := range newRescuer {
		newRescuerRule = append(newRescuerRule, newRescuerItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "RescuerChanged", newRescuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22RescuerChanged)
				if err := _FiatTokenV22.contract.UnpackLog(event, "RescuerChanged", log); err != nil {
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

// ParseRescuerChanged is a log parse operation binding the contract event 0xe475e580d85111348e40d8ca33cfdd74c30fe1655c2d8537a13abc10065ffa5a.
//
// Solidity: event RescuerChanged(address indexed newRescuer)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseRescuerChanged(log types.Log) (*FiatTokenV22RescuerChanged, error) {
	event := new(FiatTokenV22RescuerChanged)
	if err := _FiatTokenV22.contract.UnpackLog(event, "RescuerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the FiatTokenV22 contract.
type FiatTokenV22TransferIterator struct {
	Event *FiatTokenV22Transfer // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22Transfer)
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
		it.Event = new(FiatTokenV22Transfer)
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
func (it *FiatTokenV22TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22Transfer represents a Transfer event raised by the FiatTokenV22 contract.
type FiatTokenV22Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FiatTokenV22TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22TransferIterator{contract: _FiatTokenV22.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *FiatTokenV22Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22Transfer)
				if err := _FiatTokenV22.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseTransfer(log types.Log) (*FiatTokenV22Transfer, error) {
	event := new(FiatTokenV22Transfer)
	if err := _FiatTokenV22.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22UnBlacklistedIterator is returned from FilterUnBlacklisted and is used to iterate over the raw logs and unpacked data for UnBlacklisted events raised by the FiatTokenV22 contract.
type FiatTokenV22UnBlacklistedIterator struct {
	Event *FiatTokenV22UnBlacklisted // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22UnBlacklistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22UnBlacklisted)
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
		it.Event = new(FiatTokenV22UnBlacklisted)
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
func (it *FiatTokenV22UnBlacklistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22UnBlacklistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22UnBlacklisted represents a UnBlacklisted event raised by the FiatTokenV22 contract.
type FiatTokenV22UnBlacklisted struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnBlacklisted is a free log retrieval operation binding the contract event 0x117e3210bb9aa7d9baff172026820255c6f6c30ba8999d1c2fd88e2848137c4e.
//
// Solidity: event UnBlacklisted(address indexed _account)
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterUnBlacklisted(opts *bind.FilterOpts, _account []common.Address) (*FiatTokenV22UnBlacklistedIterator, error) {

	var _accountRule []interface{}
	for _, _accountItem := range _account {
		_accountRule = append(_accountRule, _accountItem)
	}

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "UnBlacklisted", _accountRule)
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22UnBlacklistedIterator{contract: _FiatTokenV22.contract, event: "UnBlacklisted", logs: logs, sub: sub}, nil
}

// WatchUnBlacklisted is a free log subscription operation binding the contract event 0x117e3210bb9aa7d9baff172026820255c6f6c30ba8999d1c2fd88e2848137c4e.
//
// Solidity: event UnBlacklisted(address indexed _account)
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchUnBlacklisted(opts *bind.WatchOpts, sink chan<- *FiatTokenV22UnBlacklisted, _account []common.Address) (event.Subscription, error) {

	var _accountRule []interface{}
	for _, _accountItem := range _account {
		_accountRule = append(_accountRule, _accountItem)
	}

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "UnBlacklisted", _accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22UnBlacklisted)
				if err := _FiatTokenV22.contract.UnpackLog(event, "UnBlacklisted", log); err != nil {
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

// ParseUnBlacklisted is a log parse operation binding the contract event 0x117e3210bb9aa7d9baff172026820255c6f6c30ba8999d1c2fd88e2848137c4e.
//
// Solidity: event UnBlacklisted(address indexed _account)
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseUnBlacklisted(log types.Log) (*FiatTokenV22UnBlacklisted, error) {
	event := new(FiatTokenV22UnBlacklisted)
	if err := _FiatTokenV22.contract.UnpackLog(event, "UnBlacklisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatTokenV22UnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the FiatTokenV22 contract.
type FiatTokenV22UnpauseIterator struct {
	Event *FiatTokenV22Unpause // Event containing the contract specifics and raw log

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
func (it *FiatTokenV22UnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatTokenV22Unpause)
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
		it.Event = new(FiatTokenV22Unpause)
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
func (it *FiatTokenV22UnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatTokenV22UnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatTokenV22Unpause represents a Unpause event raised by the FiatTokenV22 contract.
type FiatTokenV22Unpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_FiatTokenV22 *FiatTokenV22Filterer) FilterUnpause(opts *bind.FilterOpts) (*FiatTokenV22UnpauseIterator, error) {

	logs, sub, err := _FiatTokenV22.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &FiatTokenV22UnpauseIterator{contract: _FiatTokenV22.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_FiatTokenV22 *FiatTokenV22Filterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *FiatTokenV22Unpause) (event.Subscription, error) {

	logs, sub, err := _FiatTokenV22.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatTokenV22Unpause)
				if err := _FiatTokenV22.contract.UnpackLog(event, "Unpause", log); err != nil {
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

// ParseUnpause is a log parse operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_FiatTokenV22 *FiatTokenV22Filterer) ParseUnpause(log types.Log) (*FiatTokenV22Unpause, error) {
	event := new(FiatTokenV22Unpause)
	if err := _FiatTokenV22.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
