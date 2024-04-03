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

// VmSafeDirEntry is an auto generated low-level Go binding around an user-defined struct.
type VmSafeDirEntry struct {
	ErrorMessage string
	Path         string
	Depth        uint64
	IsDir        bool
	IsSymlink    bool
}

// VmSafeFsMetadata is an auto generated low-level Go binding around an user-defined struct.
type VmSafeFsMetadata struct {
	IsDir     bool
	IsSymlink bool
	Length    *big.Int
	ReadOnly  bool
	Modified  *big.Int
	Accessed  *big.Int
	Created   *big.Int
}

// VmSafeLog is an auto generated low-level Go binding around an user-defined struct.
type VmSafeLog struct {
	Topics  [][32]byte
	Data    []byte
	Emitter common.Address
}

// VmSafeRpc is an auto generated low-level Go binding around an user-defined struct.
type VmSafeRpc struct {
	Key string
	Url string
}

// SafeMetaData contains all meta data concerning the Safe contract.
var SafeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"accesses\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"readSlots\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"writeSlots\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"privateKey\",\"type\":\"uint256\"}],\"name\":\"addr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"keyAddr\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"condition\",\"type\":\"bool\"}],\"name\":\"assume\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"char\",\"type\":\"string\"}],\"name\":\"breakpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"char\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"breakpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"broadcast\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"broadcast\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"privateKey\",\"type\":\"uint256\"}],\"name\":\"broadcast\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"closeFile\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"recursive\",\"type\":\"bool\"}],\"name\":\"createDir\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"mnemonic\",\"type\":\"string\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"deriveKey\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"privateKey\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"mnemonic\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"derivationPath\",\"type\":\"string\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"deriveKey\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"privateKey\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"envAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"}],\"name\":\"envAddress\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"value\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"envBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"}],\"name\":\"envBool\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"value\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"envBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"}],\"name\":\"envBytes\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"value\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"}],\"name\":\"envBytes32\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"value\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"envBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"}],\"name\":\"envInt\",\"outputs\":[{\"internalType\":\"int256[]\",\"name\":\"value\",\"type\":\"int256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"envInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"defaultValue\",\"type\":\"bytes32[]\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"value\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"},{\"internalType\":\"int256[]\",\"name\":\"defaultValue\",\"type\":\"int256[]\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"int256[]\",\"name\":\"value\",\"type\":\"int256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"defaultValue\",\"type\":\"bool\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"defaultValue\",\"type\":\"address\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"defaultValue\",\"type\":\"uint256\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"},{\"internalType\":\"bytes[]\",\"name\":\"defaultValue\",\"type\":\"bytes[]\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"value\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"defaultValue\",\"type\":\"uint256[]\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"value\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"defaultValue\",\"type\":\"string[]\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"value\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"defaultValue\",\"type\":\"bytes\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"defaultValue\",\"type\":\"bytes32\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"defaultValue\",\"type\":\"int256\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"defaultValue\",\"type\":\"address[]\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"value\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"defaultValue\",\"type\":\"string\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"},{\"internalType\":\"bool[]\",\"name\":\"defaultValue\",\"type\":\"bool[]\"}],\"name\":\"envOr\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"value\",\"type\":\"bool[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"}],\"name\":\"envString\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"value\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"envString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"envUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"delim\",\"type\":\"string\"}],\"name\":\"envUint\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"value\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"commandInput\",\"type\":\"string[]\"}],\"name\":\"ffi\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"fsMetadata\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isDir\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isSymlink\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"readOnly\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"modified\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accessed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"created\",\"type\":\"uint256\"}],\"internalType\":\"structVmSafe.FsMetadata\",\"name\":\"metadata\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"artifactPath\",\"type\":\"string\"}],\"name\":\"getCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"creationBytecode\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"artifactPath\",\"type\":\"string\"}],\"name\":\"getDeployedCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"runtimeBytecode\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRecordedLogs\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"}],\"internalType\":\"structVmSafe.Log[]\",\"name\":\"logs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"newLabel\",\"type\":\"string\"}],\"name\":\"label\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"load\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"name\":\"parseAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"parsedValue\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"name\":\"parseBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"parsedValue\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"name\":\"parseBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"parsedValue\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"name\":\"parseBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"parsedValue\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"name\":\"parseInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"parsedValue\",\"type\":\"int256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"name\":\"parseJson\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"abiEncodedData\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"parseJson\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"abiEncodedData\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonAddressArray\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonBoolArray\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonBytes32Array\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonBytesArray\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonIntArray\",\"outputs\":[{\"internalType\":\"int256[]\",\"name\":\"\",\"type\":\"int256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonStringArray\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"parseJsonUintArray\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"name\":\"parseUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"parsedValue\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseGasMetering\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"projectRoot\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"maxDepth\",\"type\":\"uint64\"}],\"name\":\"readDir\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"depth\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isDir\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isSymlink\",\"type\":\"bool\"}],\"internalType\":\"structVmSafe.DirEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"maxDepth\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"followLinks\",\"type\":\"bool\"}],\"name\":\"readDir\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"depth\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isDir\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isSymlink\",\"type\":\"bool\"}],\"internalType\":\"structVmSafe.DirEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"readDir\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"errorMessage\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"depth\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isDir\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isSymlink\",\"type\":\"bool\"}],\"internalType\":\"structVmSafe.DirEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"readFile\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"readFileBinary\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"readLine\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"line\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"linkPath\",\"type\":\"string\"}],\"name\":\"readLink\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"targetPath\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"record\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recordLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"privateKey\",\"type\":\"uint256\"}],\"name\":\"rememberKey\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"keyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"recursive\",\"type\":\"bool\"}],\"name\":\"removeDir\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"removeFile\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resumeGasMetering\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"rpcAlias\",\"type\":\"string\"}],\"name\":\"rpcUrl\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rpcUrlStructs\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"internalType\":\"structVmSafe.Rpc[]\",\"name\":\"urls\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rpcUrls\",\"outputs\":[{\"internalType\":\"string[2][]\",\"name\":\"urls\",\"type\":\"string[2][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"values\",\"type\":\"address[]\"}],\"name\":\"serializeAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"name\":\"serializeAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"bool[]\",\"name\":\"values\",\"type\":\"bool[]\"}],\"name\":\"serializeBool\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"serializeBool\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"}],\"name\":\"serializeBytes\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"serializeBytes\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"values\",\"type\":\"bytes32[]\"}],\"name\":\"serializeBytes32\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"serializeBytes32\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"serializeInt\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"int256[]\",\"name\":\"values\",\"type\":\"int256[]\"}],\"name\":\"serializeInt\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"values\",\"type\":\"string[]\"}],\"name\":\"serializeString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"serializeString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"serializeUint\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"objectKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"serializeUint\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"setEnv\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"privateKey\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"name\":\"sign\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startBroadcast\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"startBroadcast\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"privateKey\",\"type\":\"uint256\"}],\"name\":\"startBroadcast\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stopBroadcast\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"value\",\"type\":\"address\"}],\"name\":\"toString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"toString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"toString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"toString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"toString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"toString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"stringifiedValue\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"writeFile\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"writeFileBinary\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valueKey\",\"type\":\"string\"}],\"name\":\"writeJson\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"json\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"}],\"name\":\"writeJson\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"path\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"writeLine\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SafeABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeMetaData.ABI instead.
var SafeABI = SafeMetaData.ABI

// Safe is an auto generated Go binding around an Ethereum contract.
type Safe struct {
	SafeCaller     // Read-only binding to the contract
	SafeTransactor // Write-only binding to the contract
	SafeFilterer   // Log filterer for contract events
}

// SafeCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeSession struct {
	Contract     *Safe             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeCallerSession struct {
	Contract *SafeCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SafeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeTransactorSession struct {
	Contract     *SafeTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeRaw struct {
	Contract *Safe // Generic contract binding to access the raw methods on
}

// SafeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeCallerRaw struct {
	Contract *SafeCaller // Generic read-only contract binding to access the raw methods on
}

// SafeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeTransactorRaw struct {
	Contract *SafeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafe creates a new instance of Safe, bound to a specific deployed contract.
func NewSafe(address common.Address, backend bind.ContractBackend) (*Safe, error) {
	contract, err := bindSafe(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Safe{SafeCaller: SafeCaller{contract: contract}, SafeTransactor: SafeTransactor{contract: contract}, SafeFilterer: SafeFilterer{contract: contract}}, nil
}

// NewSafeCaller creates a new read-only instance of Safe, bound to a specific deployed contract.
func NewSafeCaller(address common.Address, caller bind.ContractCaller) (*SafeCaller, error) {
	contract, err := bindSafe(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeCaller{contract: contract}, nil
}

// NewSafeTransactor creates a new write-only instance of Safe, bound to a specific deployed contract.
func NewSafeTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeTransactor, error) {
	contract, err := bindSafe(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeTransactor{contract: contract}, nil
}

// NewSafeFilterer creates a new log filterer instance of Safe, bound to a specific deployed contract.
func NewSafeFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeFilterer, error) {
	contract, err := bindSafe(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeFilterer{contract: contract}, nil
}

// bindSafe binds a generic wrapper to an already deployed contract.
func bindSafe(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Safe *SafeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Safe.Contract.SafeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Safe *SafeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.Contract.SafeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Safe *SafeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Safe.Contract.SafeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Safe *SafeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Safe.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Safe *SafeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Safe *SafeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Safe.Contract.contract.Transact(opts, method, params...)
}

// Addr is a free data retrieval call binding the contract method 0xffa18649.
//
// Solidity: function addr(uint256 privateKey) pure returns(address keyAddr)
func (_Safe *SafeCaller) Addr(opts *bind.CallOpts, privateKey *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "addr", privateKey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Addr is a free data retrieval call binding the contract method 0xffa18649.
//
// Solidity: function addr(uint256 privateKey) pure returns(address keyAddr)
func (_Safe *SafeSession) Addr(privateKey *big.Int) (common.Address, error) {
	return _Safe.Contract.Addr(&_Safe.CallOpts, privateKey)
}

// Addr is a free data retrieval call binding the contract method 0xffa18649.
//
// Solidity: function addr(uint256 privateKey) pure returns(address keyAddr)
func (_Safe *SafeCallerSession) Addr(privateKey *big.Int) (common.Address, error) {
	return _Safe.Contract.Addr(&_Safe.CallOpts, privateKey)
}

// Assume is a free data retrieval call binding the contract method 0x4c63e562.
//
// Solidity: function assume(bool condition) pure returns()
func (_Safe *SafeCaller) Assume(opts *bind.CallOpts, condition bool) error {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "assume", condition)

	if err != nil {
		return err
	}

	return err

}

// Assume is a free data retrieval call binding the contract method 0x4c63e562.
//
// Solidity: function assume(bool condition) pure returns()
func (_Safe *SafeSession) Assume(condition bool) error {
	return _Safe.Contract.Assume(&_Safe.CallOpts, condition)
}

// Assume is a free data retrieval call binding the contract method 0x4c63e562.
//
// Solidity: function assume(bool condition) pure returns()
func (_Safe *SafeCallerSession) Assume(condition bool) error {
	return _Safe.Contract.Assume(&_Safe.CallOpts, condition)
}

// DeriveKey is a free data retrieval call binding the contract method 0x6229498b.
//
// Solidity: function deriveKey(string mnemonic, uint32 index) pure returns(uint256 privateKey)
func (_Safe *SafeCaller) DeriveKey(opts *bind.CallOpts, mnemonic string, index uint32) (*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "deriveKey", mnemonic, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DeriveKey is a free data retrieval call binding the contract method 0x6229498b.
//
// Solidity: function deriveKey(string mnemonic, uint32 index) pure returns(uint256 privateKey)
func (_Safe *SafeSession) DeriveKey(mnemonic string, index uint32) (*big.Int, error) {
	return _Safe.Contract.DeriveKey(&_Safe.CallOpts, mnemonic, index)
}

// DeriveKey is a free data retrieval call binding the contract method 0x6229498b.
//
// Solidity: function deriveKey(string mnemonic, uint32 index) pure returns(uint256 privateKey)
func (_Safe *SafeCallerSession) DeriveKey(mnemonic string, index uint32) (*big.Int, error) {
	return _Safe.Contract.DeriveKey(&_Safe.CallOpts, mnemonic, index)
}

// DeriveKey0 is a free data retrieval call binding the contract method 0x6bcb2c1b.
//
// Solidity: function deriveKey(string mnemonic, string derivationPath, uint32 index) pure returns(uint256 privateKey)
func (_Safe *SafeCaller) DeriveKey0(opts *bind.CallOpts, mnemonic string, derivationPath string, index uint32) (*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "deriveKey0", mnemonic, derivationPath, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DeriveKey0 is a free data retrieval call binding the contract method 0x6bcb2c1b.
//
// Solidity: function deriveKey(string mnemonic, string derivationPath, uint32 index) pure returns(uint256 privateKey)
func (_Safe *SafeSession) DeriveKey0(mnemonic string, derivationPath string, index uint32) (*big.Int, error) {
	return _Safe.Contract.DeriveKey0(&_Safe.CallOpts, mnemonic, derivationPath, index)
}

// DeriveKey0 is a free data retrieval call binding the contract method 0x6bcb2c1b.
//
// Solidity: function deriveKey(string mnemonic, string derivationPath, uint32 index) pure returns(uint256 privateKey)
func (_Safe *SafeCallerSession) DeriveKey0(mnemonic string, derivationPath string, index uint32) (*big.Int, error) {
	return _Safe.Contract.DeriveKey0(&_Safe.CallOpts, mnemonic, derivationPath, index)
}

// EnvAddress is a free data retrieval call binding the contract method 0x350d56bf.
//
// Solidity: function envAddress(string name) view returns(address value)
func (_Safe *SafeCaller) EnvAddress(opts *bind.CallOpts, name string) (common.Address, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envAddress", name)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EnvAddress is a free data retrieval call binding the contract method 0x350d56bf.
//
// Solidity: function envAddress(string name) view returns(address value)
func (_Safe *SafeSession) EnvAddress(name string) (common.Address, error) {
	return _Safe.Contract.EnvAddress(&_Safe.CallOpts, name)
}

// EnvAddress is a free data retrieval call binding the contract method 0x350d56bf.
//
// Solidity: function envAddress(string name) view returns(address value)
func (_Safe *SafeCallerSession) EnvAddress(name string) (common.Address, error) {
	return _Safe.Contract.EnvAddress(&_Safe.CallOpts, name)
}

// EnvAddress0 is a free data retrieval call binding the contract method 0xad31b9fa.
//
// Solidity: function envAddress(string name, string delim) view returns(address[] value)
func (_Safe *SafeCaller) EnvAddress0(opts *bind.CallOpts, name string, delim string) ([]common.Address, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envAddress0", name, delim)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// EnvAddress0 is a free data retrieval call binding the contract method 0xad31b9fa.
//
// Solidity: function envAddress(string name, string delim) view returns(address[] value)
func (_Safe *SafeSession) EnvAddress0(name string, delim string) ([]common.Address, error) {
	return _Safe.Contract.EnvAddress0(&_Safe.CallOpts, name, delim)
}

// EnvAddress0 is a free data retrieval call binding the contract method 0xad31b9fa.
//
// Solidity: function envAddress(string name, string delim) view returns(address[] value)
func (_Safe *SafeCallerSession) EnvAddress0(name string, delim string) ([]common.Address, error) {
	return _Safe.Contract.EnvAddress0(&_Safe.CallOpts, name, delim)
}

// EnvBool is a free data retrieval call binding the contract method 0x7ed1ec7d.
//
// Solidity: function envBool(string name) view returns(bool value)
func (_Safe *SafeCaller) EnvBool(opts *bind.CallOpts, name string) (bool, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envBool", name)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EnvBool is a free data retrieval call binding the contract method 0x7ed1ec7d.
//
// Solidity: function envBool(string name) view returns(bool value)
func (_Safe *SafeSession) EnvBool(name string) (bool, error) {
	return _Safe.Contract.EnvBool(&_Safe.CallOpts, name)
}

// EnvBool is a free data retrieval call binding the contract method 0x7ed1ec7d.
//
// Solidity: function envBool(string name) view returns(bool value)
func (_Safe *SafeCallerSession) EnvBool(name string) (bool, error) {
	return _Safe.Contract.EnvBool(&_Safe.CallOpts, name)
}

// EnvBool0 is a free data retrieval call binding the contract method 0xaaaddeaf.
//
// Solidity: function envBool(string name, string delim) view returns(bool[] value)
func (_Safe *SafeCaller) EnvBool0(opts *bind.CallOpts, name string, delim string) ([]bool, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envBool0", name, delim)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// EnvBool0 is a free data retrieval call binding the contract method 0xaaaddeaf.
//
// Solidity: function envBool(string name, string delim) view returns(bool[] value)
func (_Safe *SafeSession) EnvBool0(name string, delim string) ([]bool, error) {
	return _Safe.Contract.EnvBool0(&_Safe.CallOpts, name, delim)
}

// EnvBool0 is a free data retrieval call binding the contract method 0xaaaddeaf.
//
// Solidity: function envBool(string name, string delim) view returns(bool[] value)
func (_Safe *SafeCallerSession) EnvBool0(name string, delim string) ([]bool, error) {
	return _Safe.Contract.EnvBool0(&_Safe.CallOpts, name, delim)
}

// EnvBytes is a free data retrieval call binding the contract method 0x4d7baf06.
//
// Solidity: function envBytes(string name) view returns(bytes value)
func (_Safe *SafeCaller) EnvBytes(opts *bind.CallOpts, name string) ([]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envBytes", name)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EnvBytes is a free data retrieval call binding the contract method 0x4d7baf06.
//
// Solidity: function envBytes(string name) view returns(bytes value)
func (_Safe *SafeSession) EnvBytes(name string) ([]byte, error) {
	return _Safe.Contract.EnvBytes(&_Safe.CallOpts, name)
}

// EnvBytes is a free data retrieval call binding the contract method 0x4d7baf06.
//
// Solidity: function envBytes(string name) view returns(bytes value)
func (_Safe *SafeCallerSession) EnvBytes(name string) ([]byte, error) {
	return _Safe.Contract.EnvBytes(&_Safe.CallOpts, name)
}

// EnvBytes0 is a free data retrieval call binding the contract method 0xddc2651b.
//
// Solidity: function envBytes(string name, string delim) view returns(bytes[] value)
func (_Safe *SafeCaller) EnvBytes0(opts *bind.CallOpts, name string, delim string) ([][]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envBytes0", name, delim)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// EnvBytes0 is a free data retrieval call binding the contract method 0xddc2651b.
//
// Solidity: function envBytes(string name, string delim) view returns(bytes[] value)
func (_Safe *SafeSession) EnvBytes0(name string, delim string) ([][]byte, error) {
	return _Safe.Contract.EnvBytes0(&_Safe.CallOpts, name, delim)
}

// EnvBytes0 is a free data retrieval call binding the contract method 0xddc2651b.
//
// Solidity: function envBytes(string name, string delim) view returns(bytes[] value)
func (_Safe *SafeCallerSession) EnvBytes0(name string, delim string) ([][]byte, error) {
	return _Safe.Contract.EnvBytes0(&_Safe.CallOpts, name, delim)
}

// EnvBytes32 is a free data retrieval call binding the contract method 0x5af231c1.
//
// Solidity: function envBytes32(string name, string delim) view returns(bytes32[] value)
func (_Safe *SafeCaller) EnvBytes32(opts *bind.CallOpts, name string, delim string) ([][32]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envBytes32", name, delim)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// EnvBytes32 is a free data retrieval call binding the contract method 0x5af231c1.
//
// Solidity: function envBytes32(string name, string delim) view returns(bytes32[] value)
func (_Safe *SafeSession) EnvBytes32(name string, delim string) ([][32]byte, error) {
	return _Safe.Contract.EnvBytes32(&_Safe.CallOpts, name, delim)
}

// EnvBytes32 is a free data retrieval call binding the contract method 0x5af231c1.
//
// Solidity: function envBytes32(string name, string delim) view returns(bytes32[] value)
func (_Safe *SafeCallerSession) EnvBytes32(name string, delim string) ([][32]byte, error) {
	return _Safe.Contract.EnvBytes32(&_Safe.CallOpts, name, delim)
}

// EnvBytes320 is a free data retrieval call binding the contract method 0x97949042.
//
// Solidity: function envBytes32(string name) view returns(bytes32 value)
func (_Safe *SafeCaller) EnvBytes320(opts *bind.CallOpts, name string) ([32]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envBytes320", name)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EnvBytes320 is a free data retrieval call binding the contract method 0x97949042.
//
// Solidity: function envBytes32(string name) view returns(bytes32 value)
func (_Safe *SafeSession) EnvBytes320(name string) ([32]byte, error) {
	return _Safe.Contract.EnvBytes320(&_Safe.CallOpts, name)
}

// EnvBytes320 is a free data retrieval call binding the contract method 0x97949042.
//
// Solidity: function envBytes32(string name) view returns(bytes32 value)
func (_Safe *SafeCallerSession) EnvBytes320(name string) ([32]byte, error) {
	return _Safe.Contract.EnvBytes320(&_Safe.CallOpts, name)
}

// EnvInt is a free data retrieval call binding the contract method 0x42181150.
//
// Solidity: function envInt(string name, string delim) view returns(int256[] value)
func (_Safe *SafeCaller) EnvInt(opts *bind.CallOpts, name string, delim string) ([]*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envInt", name, delim)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// EnvInt is a free data retrieval call binding the contract method 0x42181150.
//
// Solidity: function envInt(string name, string delim) view returns(int256[] value)
func (_Safe *SafeSession) EnvInt(name string, delim string) ([]*big.Int, error) {
	return _Safe.Contract.EnvInt(&_Safe.CallOpts, name, delim)
}

// EnvInt is a free data retrieval call binding the contract method 0x42181150.
//
// Solidity: function envInt(string name, string delim) view returns(int256[] value)
func (_Safe *SafeCallerSession) EnvInt(name string, delim string) ([]*big.Int, error) {
	return _Safe.Contract.EnvInt(&_Safe.CallOpts, name, delim)
}

// EnvInt0 is a free data retrieval call binding the contract method 0x892a0c61.
//
// Solidity: function envInt(string name) view returns(int256 value)
func (_Safe *SafeCaller) EnvInt0(opts *bind.CallOpts, name string) (*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envInt0", name)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EnvInt0 is a free data retrieval call binding the contract method 0x892a0c61.
//
// Solidity: function envInt(string name) view returns(int256 value)
func (_Safe *SafeSession) EnvInt0(name string) (*big.Int, error) {
	return _Safe.Contract.EnvInt0(&_Safe.CallOpts, name)
}

// EnvInt0 is a free data retrieval call binding the contract method 0x892a0c61.
//
// Solidity: function envInt(string name) view returns(int256 value)
func (_Safe *SafeCallerSession) EnvInt0(name string) (*big.Int, error) {
	return _Safe.Contract.EnvInt0(&_Safe.CallOpts, name)
}

// EnvString is a free data retrieval call binding the contract method 0x14b02bc9.
//
// Solidity: function envString(string name, string delim) view returns(string[] value)
func (_Safe *SafeCaller) EnvString(opts *bind.CallOpts, name string, delim string) ([]string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envString", name, delim)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// EnvString is a free data retrieval call binding the contract method 0x14b02bc9.
//
// Solidity: function envString(string name, string delim) view returns(string[] value)
func (_Safe *SafeSession) EnvString(name string, delim string) ([]string, error) {
	return _Safe.Contract.EnvString(&_Safe.CallOpts, name, delim)
}

// EnvString is a free data retrieval call binding the contract method 0x14b02bc9.
//
// Solidity: function envString(string name, string delim) view returns(string[] value)
func (_Safe *SafeCallerSession) EnvString(name string, delim string) ([]string, error) {
	return _Safe.Contract.EnvString(&_Safe.CallOpts, name, delim)
}

// EnvString0 is a free data retrieval call binding the contract method 0xf877cb19.
//
// Solidity: function envString(string name) view returns(string value)
func (_Safe *SafeCaller) EnvString0(opts *bind.CallOpts, name string) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envString0", name)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// EnvString0 is a free data retrieval call binding the contract method 0xf877cb19.
//
// Solidity: function envString(string name) view returns(string value)
func (_Safe *SafeSession) EnvString0(name string) (string, error) {
	return _Safe.Contract.EnvString0(&_Safe.CallOpts, name)
}

// EnvString0 is a free data retrieval call binding the contract method 0xf877cb19.
//
// Solidity: function envString(string name) view returns(string value)
func (_Safe *SafeCallerSession) EnvString0(name string) (string, error) {
	return _Safe.Contract.EnvString0(&_Safe.CallOpts, name)
}

// EnvUint is a free data retrieval call binding the contract method 0xc1978d1f.
//
// Solidity: function envUint(string name) view returns(uint256 value)
func (_Safe *SafeCaller) EnvUint(opts *bind.CallOpts, name string) (*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envUint", name)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EnvUint is a free data retrieval call binding the contract method 0xc1978d1f.
//
// Solidity: function envUint(string name) view returns(uint256 value)
func (_Safe *SafeSession) EnvUint(name string) (*big.Int, error) {
	return _Safe.Contract.EnvUint(&_Safe.CallOpts, name)
}

// EnvUint is a free data retrieval call binding the contract method 0xc1978d1f.
//
// Solidity: function envUint(string name) view returns(uint256 value)
func (_Safe *SafeCallerSession) EnvUint(name string) (*big.Int, error) {
	return _Safe.Contract.EnvUint(&_Safe.CallOpts, name)
}

// EnvUint0 is a free data retrieval call binding the contract method 0xf3dec099.
//
// Solidity: function envUint(string name, string delim) view returns(uint256[] value)
func (_Safe *SafeCaller) EnvUint0(opts *bind.CallOpts, name string, delim string) ([]*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "envUint0", name, delim)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// EnvUint0 is a free data retrieval call binding the contract method 0xf3dec099.
//
// Solidity: function envUint(string name, string delim) view returns(uint256[] value)
func (_Safe *SafeSession) EnvUint0(name string, delim string) ([]*big.Int, error) {
	return _Safe.Contract.EnvUint0(&_Safe.CallOpts, name, delim)
}

// EnvUint0 is a free data retrieval call binding the contract method 0xf3dec099.
//
// Solidity: function envUint(string name, string delim) view returns(uint256[] value)
func (_Safe *SafeCallerSession) EnvUint0(name string, delim string) ([]*big.Int, error) {
	return _Safe.Contract.EnvUint0(&_Safe.CallOpts, name, delim)
}

// FsMetadata is a free data retrieval call binding the contract method 0xaf368a08.
//
// Solidity: function fsMetadata(string path) view returns((bool,bool,uint256,bool,uint256,uint256,uint256) metadata)
func (_Safe *SafeCaller) FsMetadata(opts *bind.CallOpts, path string) (VmSafeFsMetadata, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "fsMetadata", path)

	if err != nil {
		return *new(VmSafeFsMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(VmSafeFsMetadata)).(*VmSafeFsMetadata)

	return out0, err

}

// FsMetadata is a free data retrieval call binding the contract method 0xaf368a08.
//
// Solidity: function fsMetadata(string path) view returns((bool,bool,uint256,bool,uint256,uint256,uint256) metadata)
func (_Safe *SafeSession) FsMetadata(path string) (VmSafeFsMetadata, error) {
	return _Safe.Contract.FsMetadata(&_Safe.CallOpts, path)
}

// FsMetadata is a free data retrieval call binding the contract method 0xaf368a08.
//
// Solidity: function fsMetadata(string path) view returns((bool,bool,uint256,bool,uint256,uint256,uint256) metadata)
func (_Safe *SafeCallerSession) FsMetadata(path string) (VmSafeFsMetadata, error) {
	return _Safe.Contract.FsMetadata(&_Safe.CallOpts, path)
}

// GetCode is a free data retrieval call binding the contract method 0x8d1cc925.
//
// Solidity: function getCode(string artifactPath) view returns(bytes creationBytecode)
func (_Safe *SafeCaller) GetCode(opts *bind.CallOpts, artifactPath string) ([]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "getCode", artifactPath)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetCode is a free data retrieval call binding the contract method 0x8d1cc925.
//
// Solidity: function getCode(string artifactPath) view returns(bytes creationBytecode)
func (_Safe *SafeSession) GetCode(artifactPath string) ([]byte, error) {
	return _Safe.Contract.GetCode(&_Safe.CallOpts, artifactPath)
}

// GetCode is a free data retrieval call binding the contract method 0x8d1cc925.
//
// Solidity: function getCode(string artifactPath) view returns(bytes creationBytecode)
func (_Safe *SafeCallerSession) GetCode(artifactPath string) ([]byte, error) {
	return _Safe.Contract.GetCode(&_Safe.CallOpts, artifactPath)
}

// GetDeployedCode is a free data retrieval call binding the contract method 0x3ebf73b4.
//
// Solidity: function getDeployedCode(string artifactPath) view returns(bytes runtimeBytecode)
func (_Safe *SafeCaller) GetDeployedCode(opts *bind.CallOpts, artifactPath string) ([]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "getDeployedCode", artifactPath)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetDeployedCode is a free data retrieval call binding the contract method 0x3ebf73b4.
//
// Solidity: function getDeployedCode(string artifactPath) view returns(bytes runtimeBytecode)
func (_Safe *SafeSession) GetDeployedCode(artifactPath string) ([]byte, error) {
	return _Safe.Contract.GetDeployedCode(&_Safe.CallOpts, artifactPath)
}

// GetDeployedCode is a free data retrieval call binding the contract method 0x3ebf73b4.
//
// Solidity: function getDeployedCode(string artifactPath) view returns(bytes runtimeBytecode)
func (_Safe *SafeCallerSession) GetDeployedCode(artifactPath string) ([]byte, error) {
	return _Safe.Contract.GetDeployedCode(&_Safe.CallOpts, artifactPath)
}

// GetNonce is a free data retrieval call binding the contract method 0x2d0335ab.
//
// Solidity: function getNonce(address account) view returns(uint64 nonce)
func (_Safe *SafeCaller) GetNonce(opts *bind.CallOpts, account common.Address) (uint64, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "getNonce", account)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetNonce is a free data retrieval call binding the contract method 0x2d0335ab.
//
// Solidity: function getNonce(address account) view returns(uint64 nonce)
func (_Safe *SafeSession) GetNonce(account common.Address) (uint64, error) {
	return _Safe.Contract.GetNonce(&_Safe.CallOpts, account)
}

// GetNonce is a free data retrieval call binding the contract method 0x2d0335ab.
//
// Solidity: function getNonce(address account) view returns(uint64 nonce)
func (_Safe *SafeCallerSession) GetNonce(account common.Address) (uint64, error) {
	return _Safe.Contract.GetNonce(&_Safe.CallOpts, account)
}

// Load is a free data retrieval call binding the contract method 0x667f9d70.
//
// Solidity: function load(address target, bytes32 slot) view returns(bytes32 data)
func (_Safe *SafeCaller) Load(opts *bind.CallOpts, target common.Address, slot [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "load", target, slot)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Load is a free data retrieval call binding the contract method 0x667f9d70.
//
// Solidity: function load(address target, bytes32 slot) view returns(bytes32 data)
func (_Safe *SafeSession) Load(target common.Address, slot [32]byte) ([32]byte, error) {
	return _Safe.Contract.Load(&_Safe.CallOpts, target, slot)
}

// Load is a free data retrieval call binding the contract method 0x667f9d70.
//
// Solidity: function load(address target, bytes32 slot) view returns(bytes32 data)
func (_Safe *SafeCallerSession) Load(target common.Address, slot [32]byte) ([32]byte, error) {
	return _Safe.Contract.Load(&_Safe.CallOpts, target, slot)
}

// ParseAddress is a free data retrieval call binding the contract method 0xc6ce059d.
//
// Solidity: function parseAddress(string stringifiedValue) pure returns(address parsedValue)
func (_Safe *SafeCaller) ParseAddress(opts *bind.CallOpts, stringifiedValue string) (common.Address, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseAddress", stringifiedValue)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParseAddress is a free data retrieval call binding the contract method 0xc6ce059d.
//
// Solidity: function parseAddress(string stringifiedValue) pure returns(address parsedValue)
func (_Safe *SafeSession) ParseAddress(stringifiedValue string) (common.Address, error) {
	return _Safe.Contract.ParseAddress(&_Safe.CallOpts, stringifiedValue)
}

// ParseAddress is a free data retrieval call binding the contract method 0xc6ce059d.
//
// Solidity: function parseAddress(string stringifiedValue) pure returns(address parsedValue)
func (_Safe *SafeCallerSession) ParseAddress(stringifiedValue string) (common.Address, error) {
	return _Safe.Contract.ParseAddress(&_Safe.CallOpts, stringifiedValue)
}

// ParseBool is a free data retrieval call binding the contract method 0x974ef924.
//
// Solidity: function parseBool(string stringifiedValue) pure returns(bool parsedValue)
func (_Safe *SafeCaller) ParseBool(opts *bind.CallOpts, stringifiedValue string) (bool, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseBool", stringifiedValue)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ParseBool is a free data retrieval call binding the contract method 0x974ef924.
//
// Solidity: function parseBool(string stringifiedValue) pure returns(bool parsedValue)
func (_Safe *SafeSession) ParseBool(stringifiedValue string) (bool, error) {
	return _Safe.Contract.ParseBool(&_Safe.CallOpts, stringifiedValue)
}

// ParseBool is a free data retrieval call binding the contract method 0x974ef924.
//
// Solidity: function parseBool(string stringifiedValue) pure returns(bool parsedValue)
func (_Safe *SafeCallerSession) ParseBool(stringifiedValue string) (bool, error) {
	return _Safe.Contract.ParseBool(&_Safe.CallOpts, stringifiedValue)
}

// ParseBytes is a free data retrieval call binding the contract method 0x8f5d232d.
//
// Solidity: function parseBytes(string stringifiedValue) pure returns(bytes parsedValue)
func (_Safe *SafeCaller) ParseBytes(opts *bind.CallOpts, stringifiedValue string) ([]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseBytes", stringifiedValue)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ParseBytes is a free data retrieval call binding the contract method 0x8f5d232d.
//
// Solidity: function parseBytes(string stringifiedValue) pure returns(bytes parsedValue)
func (_Safe *SafeSession) ParseBytes(stringifiedValue string) ([]byte, error) {
	return _Safe.Contract.ParseBytes(&_Safe.CallOpts, stringifiedValue)
}

// ParseBytes is a free data retrieval call binding the contract method 0x8f5d232d.
//
// Solidity: function parseBytes(string stringifiedValue) pure returns(bytes parsedValue)
func (_Safe *SafeCallerSession) ParseBytes(stringifiedValue string) ([]byte, error) {
	return _Safe.Contract.ParseBytes(&_Safe.CallOpts, stringifiedValue)
}

// ParseBytes32 is a free data retrieval call binding the contract method 0x087e6e81.
//
// Solidity: function parseBytes32(string stringifiedValue) pure returns(bytes32 parsedValue)
func (_Safe *SafeCaller) ParseBytes32(opts *bind.CallOpts, stringifiedValue string) ([32]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseBytes32", stringifiedValue)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ParseBytes32 is a free data retrieval call binding the contract method 0x087e6e81.
//
// Solidity: function parseBytes32(string stringifiedValue) pure returns(bytes32 parsedValue)
func (_Safe *SafeSession) ParseBytes32(stringifiedValue string) ([32]byte, error) {
	return _Safe.Contract.ParseBytes32(&_Safe.CallOpts, stringifiedValue)
}

// ParseBytes32 is a free data retrieval call binding the contract method 0x087e6e81.
//
// Solidity: function parseBytes32(string stringifiedValue) pure returns(bytes32 parsedValue)
func (_Safe *SafeCallerSession) ParseBytes32(stringifiedValue string) ([32]byte, error) {
	return _Safe.Contract.ParseBytes32(&_Safe.CallOpts, stringifiedValue)
}

// ParseInt is a free data retrieval call binding the contract method 0x42346c5e.
//
// Solidity: function parseInt(string stringifiedValue) pure returns(int256 parsedValue)
func (_Safe *SafeCaller) ParseInt(opts *bind.CallOpts, stringifiedValue string) (*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseInt", stringifiedValue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParseInt is a free data retrieval call binding the contract method 0x42346c5e.
//
// Solidity: function parseInt(string stringifiedValue) pure returns(int256 parsedValue)
func (_Safe *SafeSession) ParseInt(stringifiedValue string) (*big.Int, error) {
	return _Safe.Contract.ParseInt(&_Safe.CallOpts, stringifiedValue)
}

// ParseInt is a free data retrieval call binding the contract method 0x42346c5e.
//
// Solidity: function parseInt(string stringifiedValue) pure returns(int256 parsedValue)
func (_Safe *SafeCallerSession) ParseInt(stringifiedValue string) (*big.Int, error) {
	return _Safe.Contract.ParseInt(&_Safe.CallOpts, stringifiedValue)
}

// ParseJson is a free data retrieval call binding the contract method 0x6a82600a.
//
// Solidity: function parseJson(string json) pure returns(bytes abiEncodedData)
func (_Safe *SafeCaller) ParseJson(opts *bind.CallOpts, json string) ([]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseJson", json)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ParseJson is a free data retrieval call binding the contract method 0x6a82600a.
//
// Solidity: function parseJson(string json) pure returns(bytes abiEncodedData)
func (_Safe *SafeSession) ParseJson(json string) ([]byte, error) {
	return _Safe.Contract.ParseJson(&_Safe.CallOpts, json)
}

// ParseJson is a free data retrieval call binding the contract method 0x6a82600a.
//
// Solidity: function parseJson(string json) pure returns(bytes abiEncodedData)
func (_Safe *SafeCallerSession) ParseJson(json string) ([]byte, error) {
	return _Safe.Contract.ParseJson(&_Safe.CallOpts, json)
}

// ParseJson0 is a free data retrieval call binding the contract method 0x85940ef1.
//
// Solidity: function parseJson(string json, string key) pure returns(bytes abiEncodedData)
func (_Safe *SafeCaller) ParseJson0(opts *bind.CallOpts, json string, key string) ([]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseJson0", json, key)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ParseJson0 is a free data retrieval call binding the contract method 0x85940ef1.
//
// Solidity: function parseJson(string json, string key) pure returns(bytes abiEncodedData)
func (_Safe *SafeSession) ParseJson0(json string, key string) ([]byte, error) {
	return _Safe.Contract.ParseJson0(&_Safe.CallOpts, json, key)
}

// ParseJson0 is a free data retrieval call binding the contract method 0x85940ef1.
//
// Solidity: function parseJson(string json, string key) pure returns(bytes abiEncodedData)
func (_Safe *SafeCallerSession) ParseJson0(json string, key string) ([]byte, error) {
	return _Safe.Contract.ParseJson0(&_Safe.CallOpts, json, key)
}

// ParseUint is a free data retrieval call binding the contract method 0xfa91454d.
//
// Solidity: function parseUint(string stringifiedValue) pure returns(uint256 parsedValue)
func (_Safe *SafeCaller) ParseUint(opts *bind.CallOpts, stringifiedValue string) (*big.Int, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "parseUint", stringifiedValue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParseUint is a free data retrieval call binding the contract method 0xfa91454d.
//
// Solidity: function parseUint(string stringifiedValue) pure returns(uint256 parsedValue)
func (_Safe *SafeSession) ParseUint(stringifiedValue string) (*big.Int, error) {
	return _Safe.Contract.ParseUint(&_Safe.CallOpts, stringifiedValue)
}

// ParseUint is a free data retrieval call binding the contract method 0xfa91454d.
//
// Solidity: function parseUint(string stringifiedValue) pure returns(uint256 parsedValue)
func (_Safe *SafeCallerSession) ParseUint(stringifiedValue string) (*big.Int, error) {
	return _Safe.Contract.ParseUint(&_Safe.CallOpts, stringifiedValue)
}

// ProjectRoot is a free data retrieval call binding the contract method 0xd930a0e6.
//
// Solidity: function projectRoot() view returns(string path)
func (_Safe *SafeCaller) ProjectRoot(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "projectRoot")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ProjectRoot is a free data retrieval call binding the contract method 0xd930a0e6.
//
// Solidity: function projectRoot() view returns(string path)
func (_Safe *SafeSession) ProjectRoot() (string, error) {
	return _Safe.Contract.ProjectRoot(&_Safe.CallOpts)
}

// ProjectRoot is a free data retrieval call binding the contract method 0xd930a0e6.
//
// Solidity: function projectRoot() view returns(string path)
func (_Safe *SafeCallerSession) ProjectRoot() (string, error) {
	return _Safe.Contract.ProjectRoot(&_Safe.CallOpts)
}

// ReadDir is a free data retrieval call binding the contract method 0x1497876c.
//
// Solidity: function readDir(string path, uint64 maxDepth) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeCaller) ReadDir(opts *bind.CallOpts, path string, maxDepth uint64) ([]VmSafeDirEntry, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "readDir", path, maxDepth)

	if err != nil {
		return *new([]VmSafeDirEntry), err
	}

	out0 := *abi.ConvertType(out[0], new([]VmSafeDirEntry)).(*[]VmSafeDirEntry)

	return out0, err

}

// ReadDir is a free data retrieval call binding the contract method 0x1497876c.
//
// Solidity: function readDir(string path, uint64 maxDepth) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeSession) ReadDir(path string, maxDepth uint64) ([]VmSafeDirEntry, error) {
	return _Safe.Contract.ReadDir(&_Safe.CallOpts, path, maxDepth)
}

// ReadDir is a free data retrieval call binding the contract method 0x1497876c.
//
// Solidity: function readDir(string path, uint64 maxDepth) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeCallerSession) ReadDir(path string, maxDepth uint64) ([]VmSafeDirEntry, error) {
	return _Safe.Contract.ReadDir(&_Safe.CallOpts, path, maxDepth)
}

// ReadDir0 is a free data retrieval call binding the contract method 0x8102d70d.
//
// Solidity: function readDir(string path, uint64 maxDepth, bool followLinks) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeCaller) ReadDir0(opts *bind.CallOpts, path string, maxDepth uint64, followLinks bool) ([]VmSafeDirEntry, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "readDir0", path, maxDepth, followLinks)

	if err != nil {
		return *new([]VmSafeDirEntry), err
	}

	out0 := *abi.ConvertType(out[0], new([]VmSafeDirEntry)).(*[]VmSafeDirEntry)

	return out0, err

}

// ReadDir0 is a free data retrieval call binding the contract method 0x8102d70d.
//
// Solidity: function readDir(string path, uint64 maxDepth, bool followLinks) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeSession) ReadDir0(path string, maxDepth uint64, followLinks bool) ([]VmSafeDirEntry, error) {
	return _Safe.Contract.ReadDir0(&_Safe.CallOpts, path, maxDepth, followLinks)
}

// ReadDir0 is a free data retrieval call binding the contract method 0x8102d70d.
//
// Solidity: function readDir(string path, uint64 maxDepth, bool followLinks) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeCallerSession) ReadDir0(path string, maxDepth uint64, followLinks bool) ([]VmSafeDirEntry, error) {
	return _Safe.Contract.ReadDir0(&_Safe.CallOpts, path, maxDepth, followLinks)
}

// ReadDir1 is a free data retrieval call binding the contract method 0xc4bc59e0.
//
// Solidity: function readDir(string path) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeCaller) ReadDir1(opts *bind.CallOpts, path string) ([]VmSafeDirEntry, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "readDir1", path)

	if err != nil {
		return *new([]VmSafeDirEntry), err
	}

	out0 := *abi.ConvertType(out[0], new([]VmSafeDirEntry)).(*[]VmSafeDirEntry)

	return out0, err

}

// ReadDir1 is a free data retrieval call binding the contract method 0xc4bc59e0.
//
// Solidity: function readDir(string path) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeSession) ReadDir1(path string) ([]VmSafeDirEntry, error) {
	return _Safe.Contract.ReadDir1(&_Safe.CallOpts, path)
}

// ReadDir1 is a free data retrieval call binding the contract method 0xc4bc59e0.
//
// Solidity: function readDir(string path) view returns((string,string,uint64,bool,bool)[] entries)
func (_Safe *SafeCallerSession) ReadDir1(path string) ([]VmSafeDirEntry, error) {
	return _Safe.Contract.ReadDir1(&_Safe.CallOpts, path)
}

// ReadFile is a free data retrieval call binding the contract method 0x60f9bb11.
//
// Solidity: function readFile(string path) view returns(string data)
func (_Safe *SafeCaller) ReadFile(opts *bind.CallOpts, path string) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "readFile", path)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ReadFile is a free data retrieval call binding the contract method 0x60f9bb11.
//
// Solidity: function readFile(string path) view returns(string data)
func (_Safe *SafeSession) ReadFile(path string) (string, error) {
	return _Safe.Contract.ReadFile(&_Safe.CallOpts, path)
}

// ReadFile is a free data retrieval call binding the contract method 0x60f9bb11.
//
// Solidity: function readFile(string path) view returns(string data)
func (_Safe *SafeCallerSession) ReadFile(path string) (string, error) {
	return _Safe.Contract.ReadFile(&_Safe.CallOpts, path)
}

// ReadFileBinary is a free data retrieval call binding the contract method 0x16ed7bc4.
//
// Solidity: function readFileBinary(string path) view returns(bytes data)
func (_Safe *SafeCaller) ReadFileBinary(opts *bind.CallOpts, path string) ([]byte, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "readFileBinary", path)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ReadFileBinary is a free data retrieval call binding the contract method 0x16ed7bc4.
//
// Solidity: function readFileBinary(string path) view returns(bytes data)
func (_Safe *SafeSession) ReadFileBinary(path string) ([]byte, error) {
	return _Safe.Contract.ReadFileBinary(&_Safe.CallOpts, path)
}

// ReadFileBinary is a free data retrieval call binding the contract method 0x16ed7bc4.
//
// Solidity: function readFileBinary(string path) view returns(bytes data)
func (_Safe *SafeCallerSession) ReadFileBinary(path string) ([]byte, error) {
	return _Safe.Contract.ReadFileBinary(&_Safe.CallOpts, path)
}

// ReadLine is a free data retrieval call binding the contract method 0x70f55728.
//
// Solidity: function readLine(string path) view returns(string line)
func (_Safe *SafeCaller) ReadLine(opts *bind.CallOpts, path string) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "readLine", path)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ReadLine is a free data retrieval call binding the contract method 0x70f55728.
//
// Solidity: function readLine(string path) view returns(string line)
func (_Safe *SafeSession) ReadLine(path string) (string, error) {
	return _Safe.Contract.ReadLine(&_Safe.CallOpts, path)
}

// ReadLine is a free data retrieval call binding the contract method 0x70f55728.
//
// Solidity: function readLine(string path) view returns(string line)
func (_Safe *SafeCallerSession) ReadLine(path string) (string, error) {
	return _Safe.Contract.ReadLine(&_Safe.CallOpts, path)
}

// ReadLink is a free data retrieval call binding the contract method 0x9f5684a2.
//
// Solidity: function readLink(string linkPath) view returns(string targetPath)
func (_Safe *SafeCaller) ReadLink(opts *bind.CallOpts, linkPath string) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "readLink", linkPath)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ReadLink is a free data retrieval call binding the contract method 0x9f5684a2.
//
// Solidity: function readLink(string linkPath) view returns(string targetPath)
func (_Safe *SafeSession) ReadLink(linkPath string) (string, error) {
	return _Safe.Contract.ReadLink(&_Safe.CallOpts, linkPath)
}

// ReadLink is a free data retrieval call binding the contract method 0x9f5684a2.
//
// Solidity: function readLink(string linkPath) view returns(string targetPath)
func (_Safe *SafeCallerSession) ReadLink(linkPath string) (string, error) {
	return _Safe.Contract.ReadLink(&_Safe.CallOpts, linkPath)
}

// RpcUrl is a free data retrieval call binding the contract method 0x975a6ce9.
//
// Solidity: function rpcUrl(string rpcAlias) view returns(string json)
func (_Safe *SafeCaller) RpcUrl(opts *bind.CallOpts, rpcAlias string) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "rpcUrl", rpcAlias)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RpcUrl is a free data retrieval call binding the contract method 0x975a6ce9.
//
// Solidity: function rpcUrl(string rpcAlias) view returns(string json)
func (_Safe *SafeSession) RpcUrl(rpcAlias string) (string, error) {
	return _Safe.Contract.RpcUrl(&_Safe.CallOpts, rpcAlias)
}

// RpcUrl is a free data retrieval call binding the contract method 0x975a6ce9.
//
// Solidity: function rpcUrl(string rpcAlias) view returns(string json)
func (_Safe *SafeCallerSession) RpcUrl(rpcAlias string) (string, error) {
	return _Safe.Contract.RpcUrl(&_Safe.CallOpts, rpcAlias)
}

// RpcUrlStructs is a free data retrieval call binding the contract method 0x9d2ad72a.
//
// Solidity: function rpcUrlStructs() view returns((string,string)[] urls)
func (_Safe *SafeCaller) RpcUrlStructs(opts *bind.CallOpts) ([]VmSafeRpc, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "rpcUrlStructs")

	if err != nil {
		return *new([]VmSafeRpc), err
	}

	out0 := *abi.ConvertType(out[0], new([]VmSafeRpc)).(*[]VmSafeRpc)

	return out0, err

}

// RpcUrlStructs is a free data retrieval call binding the contract method 0x9d2ad72a.
//
// Solidity: function rpcUrlStructs() view returns((string,string)[] urls)
func (_Safe *SafeSession) RpcUrlStructs() ([]VmSafeRpc, error) {
	return _Safe.Contract.RpcUrlStructs(&_Safe.CallOpts)
}

// RpcUrlStructs is a free data retrieval call binding the contract method 0x9d2ad72a.
//
// Solidity: function rpcUrlStructs() view returns((string,string)[] urls)
func (_Safe *SafeCallerSession) RpcUrlStructs() ([]VmSafeRpc, error) {
	return _Safe.Contract.RpcUrlStructs(&_Safe.CallOpts)
}

// RpcUrls is a free data retrieval call binding the contract method 0xa85a8418.
//
// Solidity: function rpcUrls() view returns(string[2][] urls)
func (_Safe *SafeCaller) RpcUrls(opts *bind.CallOpts) ([][2]string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "rpcUrls")

	if err != nil {
		return *new([][2]string), err
	}

	out0 := *abi.ConvertType(out[0], new([][2]string)).(*[][2]string)

	return out0, err

}

// RpcUrls is a free data retrieval call binding the contract method 0xa85a8418.
//
// Solidity: function rpcUrls() view returns(string[2][] urls)
func (_Safe *SafeSession) RpcUrls() ([][2]string, error) {
	return _Safe.Contract.RpcUrls(&_Safe.CallOpts)
}

// RpcUrls is a free data retrieval call binding the contract method 0xa85a8418.
//
// Solidity: function rpcUrls() view returns(string[2][] urls)
func (_Safe *SafeCallerSession) RpcUrls() ([][2]string, error) {
	return _Safe.Contract.RpcUrls(&_Safe.CallOpts)
}

// Sign is a free data retrieval call binding the contract method 0xe341eaa4.
//
// Solidity: function sign(uint256 privateKey, bytes32 digest) pure returns(uint8 v, bytes32 r, bytes32 s)
func (_Safe *SafeCaller) Sign(opts *bind.CallOpts, privateKey *big.Int, digest [32]byte) (struct {
	V uint8
	R [32]byte
	S [32]byte
}, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "sign", privateKey, digest)

	outstruct := new(struct {
		V uint8
		R [32]byte
		S [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.V = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.R = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.S = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// Sign is a free data retrieval call binding the contract method 0xe341eaa4.
//
// Solidity: function sign(uint256 privateKey, bytes32 digest) pure returns(uint8 v, bytes32 r, bytes32 s)
func (_Safe *SafeSession) Sign(privateKey *big.Int, digest [32]byte) (struct {
	V uint8
	R [32]byte
	S [32]byte
}, error) {
	return _Safe.Contract.Sign(&_Safe.CallOpts, privateKey, digest)
}

// Sign is a free data retrieval call binding the contract method 0xe341eaa4.
//
// Solidity: function sign(uint256 privateKey, bytes32 digest) pure returns(uint8 v, bytes32 r, bytes32 s)
func (_Safe *SafeCallerSession) Sign(privateKey *big.Int, digest [32]byte) (struct {
	V uint8
	R [32]byte
	S [32]byte
}, error) {
	return _Safe.Contract.Sign(&_Safe.CallOpts, privateKey, digest)
}

// ToString is a free data retrieval call binding the contract method 0x56ca623e.
//
// Solidity: function toString(address value) pure returns(string stringifiedValue)
func (_Safe *SafeCaller) ToString(opts *bind.CallOpts, value common.Address) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "toString", value)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ToString is a free data retrieval call binding the contract method 0x56ca623e.
//
// Solidity: function toString(address value) pure returns(string stringifiedValue)
func (_Safe *SafeSession) ToString(value common.Address) (string, error) {
	return _Safe.Contract.ToString(&_Safe.CallOpts, value)
}

// ToString is a free data retrieval call binding the contract method 0x56ca623e.
//
// Solidity: function toString(address value) pure returns(string stringifiedValue)
func (_Safe *SafeCallerSession) ToString(value common.Address) (string, error) {
	return _Safe.Contract.ToString(&_Safe.CallOpts, value)
}

// ToString0 is a free data retrieval call binding the contract method 0x6900a3ae.
//
// Solidity: function toString(uint256 value) pure returns(string stringifiedValue)
func (_Safe *SafeCaller) ToString0(opts *bind.CallOpts, value *big.Int) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "toString0", value)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ToString0 is a free data retrieval call binding the contract method 0x6900a3ae.
//
// Solidity: function toString(uint256 value) pure returns(string stringifiedValue)
func (_Safe *SafeSession) ToString0(value *big.Int) (string, error) {
	return _Safe.Contract.ToString0(&_Safe.CallOpts, value)
}

// ToString0 is a free data retrieval call binding the contract method 0x6900a3ae.
//
// Solidity: function toString(uint256 value) pure returns(string stringifiedValue)
func (_Safe *SafeCallerSession) ToString0(value *big.Int) (string, error) {
	return _Safe.Contract.ToString0(&_Safe.CallOpts, value)
}

// ToString1 is a free data retrieval call binding the contract method 0x71aad10d.
//
// Solidity: function toString(bytes value) pure returns(string stringifiedValue)
func (_Safe *SafeCaller) ToString1(opts *bind.CallOpts, value []byte) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "toString1", value)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ToString1 is a free data retrieval call binding the contract method 0x71aad10d.
//
// Solidity: function toString(bytes value) pure returns(string stringifiedValue)
func (_Safe *SafeSession) ToString1(value []byte) (string, error) {
	return _Safe.Contract.ToString1(&_Safe.CallOpts, value)
}

// ToString1 is a free data retrieval call binding the contract method 0x71aad10d.
//
// Solidity: function toString(bytes value) pure returns(string stringifiedValue)
func (_Safe *SafeCallerSession) ToString1(value []byte) (string, error) {
	return _Safe.Contract.ToString1(&_Safe.CallOpts, value)
}

// ToString2 is a free data retrieval call binding the contract method 0x71dce7da.
//
// Solidity: function toString(bool value) pure returns(string stringifiedValue)
func (_Safe *SafeCaller) ToString2(opts *bind.CallOpts, value bool) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "toString2", value)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ToString2 is a free data retrieval call binding the contract method 0x71dce7da.
//
// Solidity: function toString(bool value) pure returns(string stringifiedValue)
func (_Safe *SafeSession) ToString2(value bool) (string, error) {
	return _Safe.Contract.ToString2(&_Safe.CallOpts, value)
}

// ToString2 is a free data retrieval call binding the contract method 0x71dce7da.
//
// Solidity: function toString(bool value) pure returns(string stringifiedValue)
func (_Safe *SafeCallerSession) ToString2(value bool) (string, error) {
	return _Safe.Contract.ToString2(&_Safe.CallOpts, value)
}

// ToString3 is a free data retrieval call binding the contract method 0xa322c40e.
//
// Solidity: function toString(int256 value) pure returns(string stringifiedValue)
func (_Safe *SafeCaller) ToString3(opts *bind.CallOpts, value *big.Int) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "toString3", value)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ToString3 is a free data retrieval call binding the contract method 0xa322c40e.
//
// Solidity: function toString(int256 value) pure returns(string stringifiedValue)
func (_Safe *SafeSession) ToString3(value *big.Int) (string, error) {
	return _Safe.Contract.ToString3(&_Safe.CallOpts, value)
}

// ToString3 is a free data retrieval call binding the contract method 0xa322c40e.
//
// Solidity: function toString(int256 value) pure returns(string stringifiedValue)
func (_Safe *SafeCallerSession) ToString3(value *big.Int) (string, error) {
	return _Safe.Contract.ToString3(&_Safe.CallOpts, value)
}

// ToString4 is a free data retrieval call binding the contract method 0xb11a19e8.
//
// Solidity: function toString(bytes32 value) pure returns(string stringifiedValue)
func (_Safe *SafeCaller) ToString4(opts *bind.CallOpts, value [32]byte) (string, error) {
	var out []interface{}
	err := _Safe.contract.Call(opts, &out, "toString4", value)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ToString4 is a free data retrieval call binding the contract method 0xb11a19e8.
//
// Solidity: function toString(bytes32 value) pure returns(string stringifiedValue)
func (_Safe *SafeSession) ToString4(value [32]byte) (string, error) {
	return _Safe.Contract.ToString4(&_Safe.CallOpts, value)
}

// ToString4 is a free data retrieval call binding the contract method 0xb11a19e8.
//
// Solidity: function toString(bytes32 value) pure returns(string stringifiedValue)
func (_Safe *SafeCallerSession) ToString4(value [32]byte) (string, error) {
	return _Safe.Contract.ToString4(&_Safe.CallOpts, value)
}

// Accesses is a paid mutator transaction binding the contract method 0x65bc9481.
//
// Solidity: function accesses(address target) returns(bytes32[] readSlots, bytes32[] writeSlots)
func (_Safe *SafeTransactor) Accesses(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "accesses", target)
}

// Accesses is a paid mutator transaction binding the contract method 0x65bc9481.
//
// Solidity: function accesses(address target) returns(bytes32[] readSlots, bytes32[] writeSlots)
func (_Safe *SafeSession) Accesses(target common.Address) (*types.Transaction, error) {
	return _Safe.Contract.Accesses(&_Safe.TransactOpts, target)
}

// Accesses is a paid mutator transaction binding the contract method 0x65bc9481.
//
// Solidity: function accesses(address target) returns(bytes32[] readSlots, bytes32[] writeSlots)
func (_Safe *SafeTransactorSession) Accesses(target common.Address) (*types.Transaction, error) {
	return _Safe.Contract.Accesses(&_Safe.TransactOpts, target)
}

// Breakpoint is a paid mutator transaction binding the contract method 0xf0259e92.
//
// Solidity: function breakpoint(string char) returns()
func (_Safe *SafeTransactor) Breakpoint(opts *bind.TransactOpts, char string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "breakpoint", char)
}

// Breakpoint is a paid mutator transaction binding the contract method 0xf0259e92.
//
// Solidity: function breakpoint(string char) returns()
func (_Safe *SafeSession) Breakpoint(char string) (*types.Transaction, error) {
	return _Safe.Contract.Breakpoint(&_Safe.TransactOpts, char)
}

// Breakpoint is a paid mutator transaction binding the contract method 0xf0259e92.
//
// Solidity: function breakpoint(string char) returns()
func (_Safe *SafeTransactorSession) Breakpoint(char string) (*types.Transaction, error) {
	return _Safe.Contract.Breakpoint(&_Safe.TransactOpts, char)
}

// Breakpoint0 is a paid mutator transaction binding the contract method 0xf7d39a8d.
//
// Solidity: function breakpoint(string char, bool value) returns()
func (_Safe *SafeTransactor) Breakpoint0(opts *bind.TransactOpts, char string, value bool) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "breakpoint0", char, value)
}

// Breakpoint0 is a paid mutator transaction binding the contract method 0xf7d39a8d.
//
// Solidity: function breakpoint(string char, bool value) returns()
func (_Safe *SafeSession) Breakpoint0(char string, value bool) (*types.Transaction, error) {
	return _Safe.Contract.Breakpoint0(&_Safe.TransactOpts, char, value)
}

// Breakpoint0 is a paid mutator transaction binding the contract method 0xf7d39a8d.
//
// Solidity: function breakpoint(string char, bool value) returns()
func (_Safe *SafeTransactorSession) Breakpoint0(char string, value bool) (*types.Transaction, error) {
	return _Safe.Contract.Breakpoint0(&_Safe.TransactOpts, char, value)
}

// Broadcast is a paid mutator transaction binding the contract method 0xafc98040.
//
// Solidity: function broadcast() returns()
func (_Safe *SafeTransactor) Broadcast(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "broadcast")
}

// Broadcast is a paid mutator transaction binding the contract method 0xafc98040.
//
// Solidity: function broadcast() returns()
func (_Safe *SafeSession) Broadcast() (*types.Transaction, error) {
	return _Safe.Contract.Broadcast(&_Safe.TransactOpts)
}

// Broadcast is a paid mutator transaction binding the contract method 0xafc98040.
//
// Solidity: function broadcast() returns()
func (_Safe *SafeTransactorSession) Broadcast() (*types.Transaction, error) {
	return _Safe.Contract.Broadcast(&_Safe.TransactOpts)
}

// Broadcast0 is a paid mutator transaction binding the contract method 0xe6962cdb.
//
// Solidity: function broadcast(address signer) returns()
func (_Safe *SafeTransactor) Broadcast0(opts *bind.TransactOpts, signer common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "broadcast0", signer)
}

// Broadcast0 is a paid mutator transaction binding the contract method 0xe6962cdb.
//
// Solidity: function broadcast(address signer) returns()
func (_Safe *SafeSession) Broadcast0(signer common.Address) (*types.Transaction, error) {
	return _Safe.Contract.Broadcast0(&_Safe.TransactOpts, signer)
}

// Broadcast0 is a paid mutator transaction binding the contract method 0xe6962cdb.
//
// Solidity: function broadcast(address signer) returns()
func (_Safe *SafeTransactorSession) Broadcast0(signer common.Address) (*types.Transaction, error) {
	return _Safe.Contract.Broadcast0(&_Safe.TransactOpts, signer)
}

// Broadcast1 is a paid mutator transaction binding the contract method 0xf67a965b.
//
// Solidity: function broadcast(uint256 privateKey) returns()
func (_Safe *SafeTransactor) Broadcast1(opts *bind.TransactOpts, privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "broadcast1", privateKey)
}

// Broadcast1 is a paid mutator transaction binding the contract method 0xf67a965b.
//
// Solidity: function broadcast(uint256 privateKey) returns()
func (_Safe *SafeSession) Broadcast1(privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.Broadcast1(&_Safe.TransactOpts, privateKey)
}

// Broadcast1 is a paid mutator transaction binding the contract method 0xf67a965b.
//
// Solidity: function broadcast(uint256 privateKey) returns()
func (_Safe *SafeTransactorSession) Broadcast1(privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.Broadcast1(&_Safe.TransactOpts, privateKey)
}

// CloseFile is a paid mutator transaction binding the contract method 0x48c3241f.
//
// Solidity: function closeFile(string path) returns()
func (_Safe *SafeTransactor) CloseFile(opts *bind.TransactOpts, path string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "closeFile", path)
}

// CloseFile is a paid mutator transaction binding the contract method 0x48c3241f.
//
// Solidity: function closeFile(string path) returns()
func (_Safe *SafeSession) CloseFile(path string) (*types.Transaction, error) {
	return _Safe.Contract.CloseFile(&_Safe.TransactOpts, path)
}

// CloseFile is a paid mutator transaction binding the contract method 0x48c3241f.
//
// Solidity: function closeFile(string path) returns()
func (_Safe *SafeTransactorSession) CloseFile(path string) (*types.Transaction, error) {
	return _Safe.Contract.CloseFile(&_Safe.TransactOpts, path)
}

// CreateDir is a paid mutator transaction binding the contract method 0x168b64d3.
//
// Solidity: function createDir(string path, bool recursive) returns()
func (_Safe *SafeTransactor) CreateDir(opts *bind.TransactOpts, path string, recursive bool) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "createDir", path, recursive)
}

// CreateDir is a paid mutator transaction binding the contract method 0x168b64d3.
//
// Solidity: function createDir(string path, bool recursive) returns()
func (_Safe *SafeSession) CreateDir(path string, recursive bool) (*types.Transaction, error) {
	return _Safe.Contract.CreateDir(&_Safe.TransactOpts, path, recursive)
}

// CreateDir is a paid mutator transaction binding the contract method 0x168b64d3.
//
// Solidity: function createDir(string path, bool recursive) returns()
func (_Safe *SafeTransactorSession) CreateDir(path string, recursive bool) (*types.Transaction, error) {
	return _Safe.Contract.CreateDir(&_Safe.TransactOpts, path, recursive)
}

// EnvOr is a paid mutator transaction binding the contract method 0x2281f367.
//
// Solidity: function envOr(string name, string delim, bytes32[] defaultValue) returns(bytes32[] value)
func (_Safe *SafeTransactor) EnvOr(opts *bind.TransactOpts, name string, delim string, defaultValue [][32]byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr", name, delim, defaultValue)
}

// EnvOr is a paid mutator transaction binding the contract method 0x2281f367.
//
// Solidity: function envOr(string name, string delim, bytes32[] defaultValue) returns(bytes32[] value)
func (_Safe *SafeSession) EnvOr(name string, delim string, defaultValue [][32]byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr is a paid mutator transaction binding the contract method 0x2281f367.
//
// Solidity: function envOr(string name, string delim, bytes32[] defaultValue) returns(bytes32[] value)
func (_Safe *SafeTransactorSession) EnvOr(name string, delim string, defaultValue [][32]byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr0 is a paid mutator transaction binding the contract method 0x4700d74b.
//
// Solidity: function envOr(string name, string delim, int256[] defaultValue) returns(int256[] value)
func (_Safe *SafeTransactor) EnvOr0(opts *bind.TransactOpts, name string, delim string, defaultValue []*big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr0", name, delim, defaultValue)
}

// EnvOr0 is a paid mutator transaction binding the contract method 0x4700d74b.
//
// Solidity: function envOr(string name, string delim, int256[] defaultValue) returns(int256[] value)
func (_Safe *SafeSession) EnvOr0(name string, delim string, defaultValue []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr0(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr0 is a paid mutator transaction binding the contract method 0x4700d74b.
//
// Solidity: function envOr(string name, string delim, int256[] defaultValue) returns(int256[] value)
func (_Safe *SafeTransactorSession) EnvOr0(name string, delim string, defaultValue []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr0(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr1 is a paid mutator transaction binding the contract method 0x4777f3cf.
//
// Solidity: function envOr(string name, bool defaultValue) returns(bool value)
func (_Safe *SafeTransactor) EnvOr1(opts *bind.TransactOpts, name string, defaultValue bool) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr1", name, defaultValue)
}

// EnvOr1 is a paid mutator transaction binding the contract method 0x4777f3cf.
//
// Solidity: function envOr(string name, bool defaultValue) returns(bool value)
func (_Safe *SafeSession) EnvOr1(name string, defaultValue bool) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr1(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr1 is a paid mutator transaction binding the contract method 0x4777f3cf.
//
// Solidity: function envOr(string name, bool defaultValue) returns(bool value)
func (_Safe *SafeTransactorSession) EnvOr1(name string, defaultValue bool) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr1(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr10 is a paid mutator transaction binding the contract method 0xc74e9deb.
//
// Solidity: function envOr(string name, string delim, address[] defaultValue) returns(address[] value)
func (_Safe *SafeTransactor) EnvOr10(opts *bind.TransactOpts, name string, delim string, defaultValue []common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr10", name, delim, defaultValue)
}

// EnvOr10 is a paid mutator transaction binding the contract method 0xc74e9deb.
//
// Solidity: function envOr(string name, string delim, address[] defaultValue) returns(address[] value)
func (_Safe *SafeSession) EnvOr10(name string, delim string, defaultValue []common.Address) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr10(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr10 is a paid mutator transaction binding the contract method 0xc74e9deb.
//
// Solidity: function envOr(string name, string delim, address[] defaultValue) returns(address[] value)
func (_Safe *SafeTransactorSession) EnvOr10(name string, delim string, defaultValue []common.Address) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr10(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr11 is a paid mutator transaction binding the contract method 0xd145736c.
//
// Solidity: function envOr(string name, string defaultValue) returns(string value)
func (_Safe *SafeTransactor) EnvOr11(opts *bind.TransactOpts, name string, defaultValue string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr11", name, defaultValue)
}

// EnvOr11 is a paid mutator transaction binding the contract method 0xd145736c.
//
// Solidity: function envOr(string name, string defaultValue) returns(string value)
func (_Safe *SafeSession) EnvOr11(name string, defaultValue string) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr11(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr11 is a paid mutator transaction binding the contract method 0xd145736c.
//
// Solidity: function envOr(string name, string defaultValue) returns(string value)
func (_Safe *SafeTransactorSession) EnvOr11(name string, defaultValue string) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr11(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr12 is a paid mutator transaction binding the contract method 0xeb85e83b.
//
// Solidity: function envOr(string name, string delim, bool[] defaultValue) returns(bool[] value)
func (_Safe *SafeTransactor) EnvOr12(opts *bind.TransactOpts, name string, delim string, defaultValue []bool) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr12", name, delim, defaultValue)
}

// EnvOr12 is a paid mutator transaction binding the contract method 0xeb85e83b.
//
// Solidity: function envOr(string name, string delim, bool[] defaultValue) returns(bool[] value)
func (_Safe *SafeSession) EnvOr12(name string, delim string, defaultValue []bool) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr12(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr12 is a paid mutator transaction binding the contract method 0xeb85e83b.
//
// Solidity: function envOr(string name, string delim, bool[] defaultValue) returns(bool[] value)
func (_Safe *SafeTransactorSession) EnvOr12(name string, delim string, defaultValue []bool) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr12(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr2 is a paid mutator transaction binding the contract method 0x561fe540.
//
// Solidity: function envOr(string name, address defaultValue) returns(address value)
func (_Safe *SafeTransactor) EnvOr2(opts *bind.TransactOpts, name string, defaultValue common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr2", name, defaultValue)
}

// EnvOr2 is a paid mutator transaction binding the contract method 0x561fe540.
//
// Solidity: function envOr(string name, address defaultValue) returns(address value)
func (_Safe *SafeSession) EnvOr2(name string, defaultValue common.Address) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr2(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr2 is a paid mutator transaction binding the contract method 0x561fe540.
//
// Solidity: function envOr(string name, address defaultValue) returns(address value)
func (_Safe *SafeTransactorSession) EnvOr2(name string, defaultValue common.Address) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr2(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr3 is a paid mutator transaction binding the contract method 0x5e97348f.
//
// Solidity: function envOr(string name, uint256 defaultValue) returns(uint256 value)
func (_Safe *SafeTransactor) EnvOr3(opts *bind.TransactOpts, name string, defaultValue *big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr3", name, defaultValue)
}

// EnvOr3 is a paid mutator transaction binding the contract method 0x5e97348f.
//
// Solidity: function envOr(string name, uint256 defaultValue) returns(uint256 value)
func (_Safe *SafeSession) EnvOr3(name string, defaultValue *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr3(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr3 is a paid mutator transaction binding the contract method 0x5e97348f.
//
// Solidity: function envOr(string name, uint256 defaultValue) returns(uint256 value)
func (_Safe *SafeTransactorSession) EnvOr3(name string, defaultValue *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr3(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr4 is a paid mutator transaction binding the contract method 0x64bc3e64.
//
// Solidity: function envOr(string name, string delim, bytes[] defaultValue) returns(bytes[] value)
func (_Safe *SafeTransactor) EnvOr4(opts *bind.TransactOpts, name string, delim string, defaultValue [][]byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr4", name, delim, defaultValue)
}

// EnvOr4 is a paid mutator transaction binding the contract method 0x64bc3e64.
//
// Solidity: function envOr(string name, string delim, bytes[] defaultValue) returns(bytes[] value)
func (_Safe *SafeSession) EnvOr4(name string, delim string, defaultValue [][]byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr4(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr4 is a paid mutator transaction binding the contract method 0x64bc3e64.
//
// Solidity: function envOr(string name, string delim, bytes[] defaultValue) returns(bytes[] value)
func (_Safe *SafeTransactorSession) EnvOr4(name string, delim string, defaultValue [][]byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr4(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr5 is a paid mutator transaction binding the contract method 0x74318528.
//
// Solidity: function envOr(string name, string delim, uint256[] defaultValue) returns(uint256[] value)
func (_Safe *SafeTransactor) EnvOr5(opts *bind.TransactOpts, name string, delim string, defaultValue []*big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr5", name, delim, defaultValue)
}

// EnvOr5 is a paid mutator transaction binding the contract method 0x74318528.
//
// Solidity: function envOr(string name, string delim, uint256[] defaultValue) returns(uint256[] value)
func (_Safe *SafeSession) EnvOr5(name string, delim string, defaultValue []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr5(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr5 is a paid mutator transaction binding the contract method 0x74318528.
//
// Solidity: function envOr(string name, string delim, uint256[] defaultValue) returns(uint256[] value)
func (_Safe *SafeTransactorSession) EnvOr5(name string, delim string, defaultValue []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr5(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr6 is a paid mutator transaction binding the contract method 0x859216bc.
//
// Solidity: function envOr(string name, string delim, string[] defaultValue) returns(string[] value)
func (_Safe *SafeTransactor) EnvOr6(opts *bind.TransactOpts, name string, delim string, defaultValue []string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr6", name, delim, defaultValue)
}

// EnvOr6 is a paid mutator transaction binding the contract method 0x859216bc.
//
// Solidity: function envOr(string name, string delim, string[] defaultValue) returns(string[] value)
func (_Safe *SafeSession) EnvOr6(name string, delim string, defaultValue []string) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr6(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr6 is a paid mutator transaction binding the contract method 0x859216bc.
//
// Solidity: function envOr(string name, string delim, string[] defaultValue) returns(string[] value)
func (_Safe *SafeTransactorSession) EnvOr6(name string, delim string, defaultValue []string) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr6(&_Safe.TransactOpts, name, delim, defaultValue)
}

// EnvOr7 is a paid mutator transaction binding the contract method 0xb3e47705.
//
// Solidity: function envOr(string name, bytes defaultValue) returns(bytes value)
func (_Safe *SafeTransactor) EnvOr7(opts *bind.TransactOpts, name string, defaultValue []byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr7", name, defaultValue)
}

// EnvOr7 is a paid mutator transaction binding the contract method 0xb3e47705.
//
// Solidity: function envOr(string name, bytes defaultValue) returns(bytes value)
func (_Safe *SafeSession) EnvOr7(name string, defaultValue []byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr7(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr7 is a paid mutator transaction binding the contract method 0xb3e47705.
//
// Solidity: function envOr(string name, bytes defaultValue) returns(bytes value)
func (_Safe *SafeTransactorSession) EnvOr7(name string, defaultValue []byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr7(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr8 is a paid mutator transaction binding the contract method 0xb4a85892.
//
// Solidity: function envOr(string name, bytes32 defaultValue) returns(bytes32 value)
func (_Safe *SafeTransactor) EnvOr8(opts *bind.TransactOpts, name string, defaultValue [32]byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr8", name, defaultValue)
}

// EnvOr8 is a paid mutator transaction binding the contract method 0xb4a85892.
//
// Solidity: function envOr(string name, bytes32 defaultValue) returns(bytes32 value)
func (_Safe *SafeSession) EnvOr8(name string, defaultValue [32]byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr8(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr8 is a paid mutator transaction binding the contract method 0xb4a85892.
//
// Solidity: function envOr(string name, bytes32 defaultValue) returns(bytes32 value)
func (_Safe *SafeTransactorSession) EnvOr8(name string, defaultValue [32]byte) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr8(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr9 is a paid mutator transaction binding the contract method 0xbbcb713e.
//
// Solidity: function envOr(string name, int256 defaultValue) returns(int256 value)
func (_Safe *SafeTransactor) EnvOr9(opts *bind.TransactOpts, name string, defaultValue *big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "envOr9", name, defaultValue)
}

// EnvOr9 is a paid mutator transaction binding the contract method 0xbbcb713e.
//
// Solidity: function envOr(string name, int256 defaultValue) returns(int256 value)
func (_Safe *SafeSession) EnvOr9(name string, defaultValue *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr9(&_Safe.TransactOpts, name, defaultValue)
}

// EnvOr9 is a paid mutator transaction binding the contract method 0xbbcb713e.
//
// Solidity: function envOr(string name, int256 defaultValue) returns(int256 value)
func (_Safe *SafeTransactorSession) EnvOr9(name string, defaultValue *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.EnvOr9(&_Safe.TransactOpts, name, defaultValue)
}

// Ffi is a paid mutator transaction binding the contract method 0x89160467.
//
// Solidity: function ffi(string[] commandInput) returns(bytes result)
func (_Safe *SafeTransactor) Ffi(opts *bind.TransactOpts, commandInput []string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "ffi", commandInput)
}

// Ffi is a paid mutator transaction binding the contract method 0x89160467.
//
// Solidity: function ffi(string[] commandInput) returns(bytes result)
func (_Safe *SafeSession) Ffi(commandInput []string) (*types.Transaction, error) {
	return _Safe.Contract.Ffi(&_Safe.TransactOpts, commandInput)
}

// Ffi is a paid mutator transaction binding the contract method 0x89160467.
//
// Solidity: function ffi(string[] commandInput) returns(bytes result)
func (_Safe *SafeTransactorSession) Ffi(commandInput []string) (*types.Transaction, error) {
	return _Safe.Contract.Ffi(&_Safe.TransactOpts, commandInput)
}

// GetLabel is a paid mutator transaction binding the contract method 0x28a249b0.
//
// Solidity: function getLabel(address account) returns(string label)
func (_Safe *SafeTransactor) GetLabel(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "getLabel", account)
}

// GetLabel is a paid mutator transaction binding the contract method 0x28a249b0.
//
// Solidity: function getLabel(address account) returns(string label)
func (_Safe *SafeSession) GetLabel(account common.Address) (*types.Transaction, error) {
	return _Safe.Contract.GetLabel(&_Safe.TransactOpts, account)
}

// GetLabel is a paid mutator transaction binding the contract method 0x28a249b0.
//
// Solidity: function getLabel(address account) returns(string label)
func (_Safe *SafeTransactorSession) GetLabel(account common.Address) (*types.Transaction, error) {
	return _Safe.Contract.GetLabel(&_Safe.TransactOpts, account)
}

// GetRecordedLogs is a paid mutator transaction binding the contract method 0x191553a4.
//
// Solidity: function getRecordedLogs() returns((bytes32[],bytes,address)[] logs)
func (_Safe *SafeTransactor) GetRecordedLogs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "getRecordedLogs")
}

// GetRecordedLogs is a paid mutator transaction binding the contract method 0x191553a4.
//
// Solidity: function getRecordedLogs() returns((bytes32[],bytes,address)[] logs)
func (_Safe *SafeSession) GetRecordedLogs() (*types.Transaction, error) {
	return _Safe.Contract.GetRecordedLogs(&_Safe.TransactOpts)
}

// GetRecordedLogs is a paid mutator transaction binding the contract method 0x191553a4.
//
// Solidity: function getRecordedLogs() returns((bytes32[],bytes,address)[] logs)
func (_Safe *SafeTransactorSession) GetRecordedLogs() (*types.Transaction, error) {
	return _Safe.Contract.GetRecordedLogs(&_Safe.TransactOpts)
}

// Label is a paid mutator transaction binding the contract method 0xc657c718.
//
// Solidity: function label(address account, string newLabel) returns()
func (_Safe *SafeTransactor) Label(opts *bind.TransactOpts, account common.Address, newLabel string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "label", account, newLabel)
}

// Label is a paid mutator transaction binding the contract method 0xc657c718.
//
// Solidity: function label(address account, string newLabel) returns()
func (_Safe *SafeSession) Label(account common.Address, newLabel string) (*types.Transaction, error) {
	return _Safe.Contract.Label(&_Safe.TransactOpts, account, newLabel)
}

// Label is a paid mutator transaction binding the contract method 0xc657c718.
//
// Solidity: function label(address account, string newLabel) returns()
func (_Safe *SafeTransactorSession) Label(account common.Address, newLabel string) (*types.Transaction, error) {
	return _Safe.Contract.Label(&_Safe.TransactOpts, account, newLabel)
}

// ParseJsonAddress is a paid mutator transaction binding the contract method 0x1e19e657.
//
// Solidity: function parseJsonAddress(string , string ) returns(address)
func (_Safe *SafeTransactor) ParseJsonAddress(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonAddress", arg0, arg1)
}

// ParseJsonAddress is a paid mutator transaction binding the contract method 0x1e19e657.
//
// Solidity: function parseJsonAddress(string , string ) returns(address)
func (_Safe *SafeSession) ParseJsonAddress(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonAddress(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonAddress is a paid mutator transaction binding the contract method 0x1e19e657.
//
// Solidity: function parseJsonAddress(string , string ) returns(address)
func (_Safe *SafeTransactorSession) ParseJsonAddress(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonAddress(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonAddressArray is a paid mutator transaction binding the contract method 0x2fce7883.
//
// Solidity: function parseJsonAddressArray(string , string ) returns(address[])
func (_Safe *SafeTransactor) ParseJsonAddressArray(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonAddressArray", arg0, arg1)
}

// ParseJsonAddressArray is a paid mutator transaction binding the contract method 0x2fce7883.
//
// Solidity: function parseJsonAddressArray(string , string ) returns(address[])
func (_Safe *SafeSession) ParseJsonAddressArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonAddressArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonAddressArray is a paid mutator transaction binding the contract method 0x2fce7883.
//
// Solidity: function parseJsonAddressArray(string , string ) returns(address[])
func (_Safe *SafeTransactorSession) ParseJsonAddressArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonAddressArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBool is a paid mutator transaction binding the contract method 0x9f86dc91.
//
// Solidity: function parseJsonBool(string , string ) returns(bool)
func (_Safe *SafeTransactor) ParseJsonBool(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonBool", arg0, arg1)
}

// ParseJsonBool is a paid mutator transaction binding the contract method 0x9f86dc91.
//
// Solidity: function parseJsonBool(string , string ) returns(bool)
func (_Safe *SafeSession) ParseJsonBool(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBool(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBool is a paid mutator transaction binding the contract method 0x9f86dc91.
//
// Solidity: function parseJsonBool(string , string ) returns(bool)
func (_Safe *SafeTransactorSession) ParseJsonBool(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBool(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBoolArray is a paid mutator transaction binding the contract method 0x91f3b94f.
//
// Solidity: function parseJsonBoolArray(string , string ) returns(bool[])
func (_Safe *SafeTransactor) ParseJsonBoolArray(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonBoolArray", arg0, arg1)
}

// ParseJsonBoolArray is a paid mutator transaction binding the contract method 0x91f3b94f.
//
// Solidity: function parseJsonBoolArray(string , string ) returns(bool[])
func (_Safe *SafeSession) ParseJsonBoolArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBoolArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBoolArray is a paid mutator transaction binding the contract method 0x91f3b94f.
//
// Solidity: function parseJsonBoolArray(string , string ) returns(bool[])
func (_Safe *SafeTransactorSession) ParseJsonBoolArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBoolArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytes is a paid mutator transaction binding the contract method 0xfd921be8.
//
// Solidity: function parseJsonBytes(string , string ) returns(bytes)
func (_Safe *SafeTransactor) ParseJsonBytes(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonBytes", arg0, arg1)
}

// ParseJsonBytes is a paid mutator transaction binding the contract method 0xfd921be8.
//
// Solidity: function parseJsonBytes(string , string ) returns(bytes)
func (_Safe *SafeSession) ParseJsonBytes(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytes(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytes is a paid mutator transaction binding the contract method 0xfd921be8.
//
// Solidity: function parseJsonBytes(string , string ) returns(bytes)
func (_Safe *SafeTransactorSession) ParseJsonBytes(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytes(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytes32 is a paid mutator transaction binding the contract method 0x1777e59d.
//
// Solidity: function parseJsonBytes32(string , string ) returns(bytes32)
func (_Safe *SafeTransactor) ParseJsonBytes32(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonBytes32", arg0, arg1)
}

// ParseJsonBytes32 is a paid mutator transaction binding the contract method 0x1777e59d.
//
// Solidity: function parseJsonBytes32(string , string ) returns(bytes32)
func (_Safe *SafeSession) ParseJsonBytes32(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytes32(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytes32 is a paid mutator transaction binding the contract method 0x1777e59d.
//
// Solidity: function parseJsonBytes32(string , string ) returns(bytes32)
func (_Safe *SafeTransactorSession) ParseJsonBytes32(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytes32(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytes32Array is a paid mutator transaction binding the contract method 0x91c75bc3.
//
// Solidity: function parseJsonBytes32Array(string , string ) returns(bytes32[])
func (_Safe *SafeTransactor) ParseJsonBytes32Array(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonBytes32Array", arg0, arg1)
}

// ParseJsonBytes32Array is a paid mutator transaction binding the contract method 0x91c75bc3.
//
// Solidity: function parseJsonBytes32Array(string , string ) returns(bytes32[])
func (_Safe *SafeSession) ParseJsonBytes32Array(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytes32Array(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytes32Array is a paid mutator transaction binding the contract method 0x91c75bc3.
//
// Solidity: function parseJsonBytes32Array(string , string ) returns(bytes32[])
func (_Safe *SafeTransactorSession) ParseJsonBytes32Array(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytes32Array(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytesArray is a paid mutator transaction binding the contract method 0x6631aa99.
//
// Solidity: function parseJsonBytesArray(string , string ) returns(bytes[])
func (_Safe *SafeTransactor) ParseJsonBytesArray(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonBytesArray", arg0, arg1)
}

// ParseJsonBytesArray is a paid mutator transaction binding the contract method 0x6631aa99.
//
// Solidity: function parseJsonBytesArray(string , string ) returns(bytes[])
func (_Safe *SafeSession) ParseJsonBytesArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytesArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonBytesArray is a paid mutator transaction binding the contract method 0x6631aa99.
//
// Solidity: function parseJsonBytesArray(string , string ) returns(bytes[])
func (_Safe *SafeTransactorSession) ParseJsonBytesArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonBytesArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonInt is a paid mutator transaction binding the contract method 0x7b048ccd.
//
// Solidity: function parseJsonInt(string , string ) returns(int256)
func (_Safe *SafeTransactor) ParseJsonInt(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonInt", arg0, arg1)
}

// ParseJsonInt is a paid mutator transaction binding the contract method 0x7b048ccd.
//
// Solidity: function parseJsonInt(string , string ) returns(int256)
func (_Safe *SafeSession) ParseJsonInt(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonInt(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonInt is a paid mutator transaction binding the contract method 0x7b048ccd.
//
// Solidity: function parseJsonInt(string , string ) returns(int256)
func (_Safe *SafeTransactorSession) ParseJsonInt(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonInt(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonIntArray is a paid mutator transaction binding the contract method 0x9983c28a.
//
// Solidity: function parseJsonIntArray(string , string ) returns(int256[])
func (_Safe *SafeTransactor) ParseJsonIntArray(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonIntArray", arg0, arg1)
}

// ParseJsonIntArray is a paid mutator transaction binding the contract method 0x9983c28a.
//
// Solidity: function parseJsonIntArray(string , string ) returns(int256[])
func (_Safe *SafeSession) ParseJsonIntArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonIntArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonIntArray is a paid mutator transaction binding the contract method 0x9983c28a.
//
// Solidity: function parseJsonIntArray(string , string ) returns(int256[])
func (_Safe *SafeTransactorSession) ParseJsonIntArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonIntArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonString is a paid mutator transaction binding the contract method 0x49c4fac8.
//
// Solidity: function parseJsonString(string , string ) returns(string)
func (_Safe *SafeTransactor) ParseJsonString(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonString", arg0, arg1)
}

// ParseJsonString is a paid mutator transaction binding the contract method 0x49c4fac8.
//
// Solidity: function parseJsonString(string , string ) returns(string)
func (_Safe *SafeSession) ParseJsonString(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonString(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonString is a paid mutator transaction binding the contract method 0x49c4fac8.
//
// Solidity: function parseJsonString(string , string ) returns(string)
func (_Safe *SafeTransactorSession) ParseJsonString(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonString(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonStringArray is a paid mutator transaction binding the contract method 0x498fdcf4.
//
// Solidity: function parseJsonStringArray(string , string ) returns(string[])
func (_Safe *SafeTransactor) ParseJsonStringArray(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonStringArray", arg0, arg1)
}

// ParseJsonStringArray is a paid mutator transaction binding the contract method 0x498fdcf4.
//
// Solidity: function parseJsonStringArray(string , string ) returns(string[])
func (_Safe *SafeSession) ParseJsonStringArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonStringArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonStringArray is a paid mutator transaction binding the contract method 0x498fdcf4.
//
// Solidity: function parseJsonStringArray(string , string ) returns(string[])
func (_Safe *SafeTransactorSession) ParseJsonStringArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonStringArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonUint is a paid mutator transaction binding the contract method 0xaddde2b6.
//
// Solidity: function parseJsonUint(string , string ) returns(uint256)
func (_Safe *SafeTransactor) ParseJsonUint(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonUint", arg0, arg1)
}

// ParseJsonUint is a paid mutator transaction binding the contract method 0xaddde2b6.
//
// Solidity: function parseJsonUint(string , string ) returns(uint256)
func (_Safe *SafeSession) ParseJsonUint(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonUint(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonUint is a paid mutator transaction binding the contract method 0xaddde2b6.
//
// Solidity: function parseJsonUint(string , string ) returns(uint256)
func (_Safe *SafeTransactorSession) ParseJsonUint(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonUint(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonUintArray is a paid mutator transaction binding the contract method 0x522074ab.
//
// Solidity: function parseJsonUintArray(string , string ) returns(uint256[])
func (_Safe *SafeTransactor) ParseJsonUintArray(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "parseJsonUintArray", arg0, arg1)
}

// ParseJsonUintArray is a paid mutator transaction binding the contract method 0x522074ab.
//
// Solidity: function parseJsonUintArray(string , string ) returns(uint256[])
func (_Safe *SafeSession) ParseJsonUintArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonUintArray(&_Safe.TransactOpts, arg0, arg1)
}

// ParseJsonUintArray is a paid mutator transaction binding the contract method 0x522074ab.
//
// Solidity: function parseJsonUintArray(string , string ) returns(uint256[])
func (_Safe *SafeTransactorSession) ParseJsonUintArray(arg0 string, arg1 string) (*types.Transaction, error) {
	return _Safe.Contract.ParseJsonUintArray(&_Safe.TransactOpts, arg0, arg1)
}

// PauseGasMetering is a paid mutator transaction binding the contract method 0xd1a5b36f.
//
// Solidity: function pauseGasMetering() returns()
func (_Safe *SafeTransactor) PauseGasMetering(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "pauseGasMetering")
}

// PauseGasMetering is a paid mutator transaction binding the contract method 0xd1a5b36f.
//
// Solidity: function pauseGasMetering() returns()
func (_Safe *SafeSession) PauseGasMetering() (*types.Transaction, error) {
	return _Safe.Contract.PauseGasMetering(&_Safe.TransactOpts)
}

// PauseGasMetering is a paid mutator transaction binding the contract method 0xd1a5b36f.
//
// Solidity: function pauseGasMetering() returns()
func (_Safe *SafeTransactorSession) PauseGasMetering() (*types.Transaction, error) {
	return _Safe.Contract.PauseGasMetering(&_Safe.TransactOpts)
}

// Record is a paid mutator transaction binding the contract method 0x266cf109.
//
// Solidity: function record() returns()
func (_Safe *SafeTransactor) Record(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "record")
}

// Record is a paid mutator transaction binding the contract method 0x266cf109.
//
// Solidity: function record() returns()
func (_Safe *SafeSession) Record() (*types.Transaction, error) {
	return _Safe.Contract.Record(&_Safe.TransactOpts)
}

// Record is a paid mutator transaction binding the contract method 0x266cf109.
//
// Solidity: function record() returns()
func (_Safe *SafeTransactorSession) Record() (*types.Transaction, error) {
	return _Safe.Contract.Record(&_Safe.TransactOpts)
}

// RecordLogs is a paid mutator transaction binding the contract method 0x41af2f52.
//
// Solidity: function recordLogs() returns()
func (_Safe *SafeTransactor) RecordLogs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "recordLogs")
}

// RecordLogs is a paid mutator transaction binding the contract method 0x41af2f52.
//
// Solidity: function recordLogs() returns()
func (_Safe *SafeSession) RecordLogs() (*types.Transaction, error) {
	return _Safe.Contract.RecordLogs(&_Safe.TransactOpts)
}

// RecordLogs is a paid mutator transaction binding the contract method 0x41af2f52.
//
// Solidity: function recordLogs() returns()
func (_Safe *SafeTransactorSession) RecordLogs() (*types.Transaction, error) {
	return _Safe.Contract.RecordLogs(&_Safe.TransactOpts)
}

// RememberKey is a paid mutator transaction binding the contract method 0x22100064.
//
// Solidity: function rememberKey(uint256 privateKey) returns(address keyAddr)
func (_Safe *SafeTransactor) RememberKey(opts *bind.TransactOpts, privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "rememberKey", privateKey)
}

// RememberKey is a paid mutator transaction binding the contract method 0x22100064.
//
// Solidity: function rememberKey(uint256 privateKey) returns(address keyAddr)
func (_Safe *SafeSession) RememberKey(privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.RememberKey(&_Safe.TransactOpts, privateKey)
}

// RememberKey is a paid mutator transaction binding the contract method 0x22100064.
//
// Solidity: function rememberKey(uint256 privateKey) returns(address keyAddr)
func (_Safe *SafeTransactorSession) RememberKey(privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.RememberKey(&_Safe.TransactOpts, privateKey)
}

// RemoveDir is a paid mutator transaction binding the contract method 0x45c62011.
//
// Solidity: function removeDir(string path, bool recursive) returns()
func (_Safe *SafeTransactor) RemoveDir(opts *bind.TransactOpts, path string, recursive bool) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "removeDir", path, recursive)
}

// RemoveDir is a paid mutator transaction binding the contract method 0x45c62011.
//
// Solidity: function removeDir(string path, bool recursive) returns()
func (_Safe *SafeSession) RemoveDir(path string, recursive bool) (*types.Transaction, error) {
	return _Safe.Contract.RemoveDir(&_Safe.TransactOpts, path, recursive)
}

// RemoveDir is a paid mutator transaction binding the contract method 0x45c62011.
//
// Solidity: function removeDir(string path, bool recursive) returns()
func (_Safe *SafeTransactorSession) RemoveDir(path string, recursive bool) (*types.Transaction, error) {
	return _Safe.Contract.RemoveDir(&_Safe.TransactOpts, path, recursive)
}

// RemoveFile is a paid mutator transaction binding the contract method 0xf1afe04d.
//
// Solidity: function removeFile(string path) returns()
func (_Safe *SafeTransactor) RemoveFile(opts *bind.TransactOpts, path string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "removeFile", path)
}

// RemoveFile is a paid mutator transaction binding the contract method 0xf1afe04d.
//
// Solidity: function removeFile(string path) returns()
func (_Safe *SafeSession) RemoveFile(path string) (*types.Transaction, error) {
	return _Safe.Contract.RemoveFile(&_Safe.TransactOpts, path)
}

// RemoveFile is a paid mutator transaction binding the contract method 0xf1afe04d.
//
// Solidity: function removeFile(string path) returns()
func (_Safe *SafeTransactorSession) RemoveFile(path string) (*types.Transaction, error) {
	return _Safe.Contract.RemoveFile(&_Safe.TransactOpts, path)
}

// ResumeGasMetering is a paid mutator transaction binding the contract method 0x2bcd50e0.
//
// Solidity: function resumeGasMetering() returns()
func (_Safe *SafeTransactor) ResumeGasMetering(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "resumeGasMetering")
}

// ResumeGasMetering is a paid mutator transaction binding the contract method 0x2bcd50e0.
//
// Solidity: function resumeGasMetering() returns()
func (_Safe *SafeSession) ResumeGasMetering() (*types.Transaction, error) {
	return _Safe.Contract.ResumeGasMetering(&_Safe.TransactOpts)
}

// ResumeGasMetering is a paid mutator transaction binding the contract method 0x2bcd50e0.
//
// Solidity: function resumeGasMetering() returns()
func (_Safe *SafeTransactorSession) ResumeGasMetering() (*types.Transaction, error) {
	return _Safe.Contract.ResumeGasMetering(&_Safe.TransactOpts)
}

// SerializeAddress is a paid mutator transaction binding the contract method 0x1e356e1a.
//
// Solidity: function serializeAddress(string objectKey, string valueKey, address[] values) returns(string json)
func (_Safe *SafeTransactor) SerializeAddress(opts *bind.TransactOpts, objectKey string, valueKey string, values []common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeAddress", objectKey, valueKey, values)
}

// SerializeAddress is a paid mutator transaction binding the contract method 0x1e356e1a.
//
// Solidity: function serializeAddress(string objectKey, string valueKey, address[] values) returns(string json)
func (_Safe *SafeSession) SerializeAddress(objectKey string, valueKey string, values []common.Address) (*types.Transaction, error) {
	return _Safe.Contract.SerializeAddress(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeAddress is a paid mutator transaction binding the contract method 0x1e356e1a.
//
// Solidity: function serializeAddress(string objectKey, string valueKey, address[] values) returns(string json)
func (_Safe *SafeTransactorSession) SerializeAddress(objectKey string, valueKey string, values []common.Address) (*types.Transaction, error) {
	return _Safe.Contract.SerializeAddress(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeAddress0 is a paid mutator transaction binding the contract method 0x972c6062.
//
// Solidity: function serializeAddress(string objectKey, string valueKey, address value) returns(string json)
func (_Safe *SafeTransactor) SerializeAddress0(opts *bind.TransactOpts, objectKey string, valueKey string, value common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeAddress0", objectKey, valueKey, value)
}

// SerializeAddress0 is a paid mutator transaction binding the contract method 0x972c6062.
//
// Solidity: function serializeAddress(string objectKey, string valueKey, address value) returns(string json)
func (_Safe *SafeSession) SerializeAddress0(objectKey string, valueKey string, value common.Address) (*types.Transaction, error) {
	return _Safe.Contract.SerializeAddress0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeAddress0 is a paid mutator transaction binding the contract method 0x972c6062.
//
// Solidity: function serializeAddress(string objectKey, string valueKey, address value) returns(string json)
func (_Safe *SafeTransactorSession) SerializeAddress0(objectKey string, valueKey string, value common.Address) (*types.Transaction, error) {
	return _Safe.Contract.SerializeAddress0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeBool is a paid mutator transaction binding the contract method 0x92925aa1.
//
// Solidity: function serializeBool(string objectKey, string valueKey, bool[] values) returns(string json)
func (_Safe *SafeTransactor) SerializeBool(opts *bind.TransactOpts, objectKey string, valueKey string, values []bool) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeBool", objectKey, valueKey, values)
}

// SerializeBool is a paid mutator transaction binding the contract method 0x92925aa1.
//
// Solidity: function serializeBool(string objectKey, string valueKey, bool[] values) returns(string json)
func (_Safe *SafeSession) SerializeBool(objectKey string, valueKey string, values []bool) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBool(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeBool is a paid mutator transaction binding the contract method 0x92925aa1.
//
// Solidity: function serializeBool(string objectKey, string valueKey, bool[] values) returns(string json)
func (_Safe *SafeTransactorSession) SerializeBool(objectKey string, valueKey string, values []bool) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBool(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeBool0 is a paid mutator transaction binding the contract method 0xac22e971.
//
// Solidity: function serializeBool(string objectKey, string valueKey, bool value) returns(string json)
func (_Safe *SafeTransactor) SerializeBool0(opts *bind.TransactOpts, objectKey string, valueKey string, value bool) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeBool0", objectKey, valueKey, value)
}

// SerializeBool0 is a paid mutator transaction binding the contract method 0xac22e971.
//
// Solidity: function serializeBool(string objectKey, string valueKey, bool value) returns(string json)
func (_Safe *SafeSession) SerializeBool0(objectKey string, valueKey string, value bool) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBool0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeBool0 is a paid mutator transaction binding the contract method 0xac22e971.
//
// Solidity: function serializeBool(string objectKey, string valueKey, bool value) returns(string json)
func (_Safe *SafeTransactorSession) SerializeBool0(objectKey string, valueKey string, value bool) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBool0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeBytes is a paid mutator transaction binding the contract method 0x9884b232.
//
// Solidity: function serializeBytes(string objectKey, string valueKey, bytes[] values) returns(string json)
func (_Safe *SafeTransactor) SerializeBytes(opts *bind.TransactOpts, objectKey string, valueKey string, values [][]byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeBytes", objectKey, valueKey, values)
}

// SerializeBytes is a paid mutator transaction binding the contract method 0x9884b232.
//
// Solidity: function serializeBytes(string objectKey, string valueKey, bytes[] values) returns(string json)
func (_Safe *SafeSession) SerializeBytes(objectKey string, valueKey string, values [][]byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeBytes is a paid mutator transaction binding the contract method 0x9884b232.
//
// Solidity: function serializeBytes(string objectKey, string valueKey, bytes[] values) returns(string json)
func (_Safe *SafeTransactorSession) SerializeBytes(objectKey string, valueKey string, values [][]byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeBytes0 is a paid mutator transaction binding the contract method 0xf21d52c7.
//
// Solidity: function serializeBytes(string objectKey, string valueKey, bytes value) returns(string json)
func (_Safe *SafeTransactor) SerializeBytes0(opts *bind.TransactOpts, objectKey string, valueKey string, value []byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeBytes0", objectKey, valueKey, value)
}

// SerializeBytes0 is a paid mutator transaction binding the contract method 0xf21d52c7.
//
// Solidity: function serializeBytes(string objectKey, string valueKey, bytes value) returns(string json)
func (_Safe *SafeSession) SerializeBytes0(objectKey string, valueKey string, value []byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeBytes0 is a paid mutator transaction binding the contract method 0xf21d52c7.
//
// Solidity: function serializeBytes(string objectKey, string valueKey, bytes value) returns(string json)
func (_Safe *SafeTransactorSession) SerializeBytes0(objectKey string, valueKey string, value []byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeBytes32 is a paid mutator transaction binding the contract method 0x201e43e2.
//
// Solidity: function serializeBytes32(string objectKey, string valueKey, bytes32[] values) returns(string json)
func (_Safe *SafeTransactor) SerializeBytes32(opts *bind.TransactOpts, objectKey string, valueKey string, values [][32]byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeBytes32", objectKey, valueKey, values)
}

// SerializeBytes32 is a paid mutator transaction binding the contract method 0x201e43e2.
//
// Solidity: function serializeBytes32(string objectKey, string valueKey, bytes32[] values) returns(string json)
func (_Safe *SafeSession) SerializeBytes32(objectKey string, valueKey string, values [][32]byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes32(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeBytes32 is a paid mutator transaction binding the contract method 0x201e43e2.
//
// Solidity: function serializeBytes32(string objectKey, string valueKey, bytes32[] values) returns(string json)
func (_Safe *SafeTransactorSession) SerializeBytes32(objectKey string, valueKey string, values [][32]byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes32(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeBytes320 is a paid mutator transaction binding the contract method 0x2d812b44.
//
// Solidity: function serializeBytes32(string objectKey, string valueKey, bytes32 value) returns(string json)
func (_Safe *SafeTransactor) SerializeBytes320(opts *bind.TransactOpts, objectKey string, valueKey string, value [32]byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeBytes320", objectKey, valueKey, value)
}

// SerializeBytes320 is a paid mutator transaction binding the contract method 0x2d812b44.
//
// Solidity: function serializeBytes32(string objectKey, string valueKey, bytes32 value) returns(string json)
func (_Safe *SafeSession) SerializeBytes320(objectKey string, valueKey string, value [32]byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes320(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeBytes320 is a paid mutator transaction binding the contract method 0x2d812b44.
//
// Solidity: function serializeBytes32(string objectKey, string valueKey, bytes32 value) returns(string json)
func (_Safe *SafeTransactorSession) SerializeBytes320(objectKey string, valueKey string, value [32]byte) (*types.Transaction, error) {
	return _Safe.Contract.SerializeBytes320(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeInt is a paid mutator transaction binding the contract method 0x3f33db60.
//
// Solidity: function serializeInt(string objectKey, string valueKey, int256 value) returns(string json)
func (_Safe *SafeTransactor) SerializeInt(opts *bind.TransactOpts, objectKey string, valueKey string, value *big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeInt", objectKey, valueKey, value)
}

// SerializeInt is a paid mutator transaction binding the contract method 0x3f33db60.
//
// Solidity: function serializeInt(string objectKey, string valueKey, int256 value) returns(string json)
func (_Safe *SafeSession) SerializeInt(objectKey string, valueKey string, value *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeInt(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeInt is a paid mutator transaction binding the contract method 0x3f33db60.
//
// Solidity: function serializeInt(string objectKey, string valueKey, int256 value) returns(string json)
func (_Safe *SafeTransactorSession) SerializeInt(objectKey string, valueKey string, value *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeInt(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeInt0 is a paid mutator transaction binding the contract method 0x7676e127.
//
// Solidity: function serializeInt(string objectKey, string valueKey, int256[] values) returns(string json)
func (_Safe *SafeTransactor) SerializeInt0(opts *bind.TransactOpts, objectKey string, valueKey string, values []*big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeInt0", objectKey, valueKey, values)
}

// SerializeInt0 is a paid mutator transaction binding the contract method 0x7676e127.
//
// Solidity: function serializeInt(string objectKey, string valueKey, int256[] values) returns(string json)
func (_Safe *SafeSession) SerializeInt0(objectKey string, valueKey string, values []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeInt0(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeInt0 is a paid mutator transaction binding the contract method 0x7676e127.
//
// Solidity: function serializeInt(string objectKey, string valueKey, int256[] values) returns(string json)
func (_Safe *SafeTransactorSession) SerializeInt0(objectKey string, valueKey string, values []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeInt0(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeString is a paid mutator transaction binding the contract method 0x561cd6f3.
//
// Solidity: function serializeString(string objectKey, string valueKey, string[] values) returns(string json)
func (_Safe *SafeTransactor) SerializeString(opts *bind.TransactOpts, objectKey string, valueKey string, values []string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeString", objectKey, valueKey, values)
}

// SerializeString is a paid mutator transaction binding the contract method 0x561cd6f3.
//
// Solidity: function serializeString(string objectKey, string valueKey, string[] values) returns(string json)
func (_Safe *SafeSession) SerializeString(objectKey string, valueKey string, values []string) (*types.Transaction, error) {
	return _Safe.Contract.SerializeString(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeString is a paid mutator transaction binding the contract method 0x561cd6f3.
//
// Solidity: function serializeString(string objectKey, string valueKey, string[] values) returns(string json)
func (_Safe *SafeTransactorSession) SerializeString(objectKey string, valueKey string, values []string) (*types.Transaction, error) {
	return _Safe.Contract.SerializeString(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeString0 is a paid mutator transaction binding the contract method 0x88da6d35.
//
// Solidity: function serializeString(string objectKey, string valueKey, string value) returns(string json)
func (_Safe *SafeTransactor) SerializeString0(opts *bind.TransactOpts, objectKey string, valueKey string, value string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeString0", objectKey, valueKey, value)
}

// SerializeString0 is a paid mutator transaction binding the contract method 0x88da6d35.
//
// Solidity: function serializeString(string objectKey, string valueKey, string value) returns(string json)
func (_Safe *SafeSession) SerializeString0(objectKey string, valueKey string, value string) (*types.Transaction, error) {
	return _Safe.Contract.SerializeString0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeString0 is a paid mutator transaction binding the contract method 0x88da6d35.
//
// Solidity: function serializeString(string objectKey, string valueKey, string value) returns(string json)
func (_Safe *SafeTransactorSession) SerializeString0(objectKey string, valueKey string, value string) (*types.Transaction, error) {
	return _Safe.Contract.SerializeString0(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeUint is a paid mutator transaction binding the contract method 0x129e9002.
//
// Solidity: function serializeUint(string objectKey, string valueKey, uint256 value) returns(string json)
func (_Safe *SafeTransactor) SerializeUint(opts *bind.TransactOpts, objectKey string, valueKey string, value *big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeUint", objectKey, valueKey, value)
}

// SerializeUint is a paid mutator transaction binding the contract method 0x129e9002.
//
// Solidity: function serializeUint(string objectKey, string valueKey, uint256 value) returns(string json)
func (_Safe *SafeSession) SerializeUint(objectKey string, valueKey string, value *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeUint(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeUint is a paid mutator transaction binding the contract method 0x129e9002.
//
// Solidity: function serializeUint(string objectKey, string valueKey, uint256 value) returns(string json)
func (_Safe *SafeTransactorSession) SerializeUint(objectKey string, valueKey string, value *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeUint(&_Safe.TransactOpts, objectKey, valueKey, value)
}

// SerializeUint0 is a paid mutator transaction binding the contract method 0xfee9a469.
//
// Solidity: function serializeUint(string objectKey, string valueKey, uint256[] values) returns(string json)
func (_Safe *SafeTransactor) SerializeUint0(opts *bind.TransactOpts, objectKey string, valueKey string, values []*big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "serializeUint0", objectKey, valueKey, values)
}

// SerializeUint0 is a paid mutator transaction binding the contract method 0xfee9a469.
//
// Solidity: function serializeUint(string objectKey, string valueKey, uint256[] values) returns(string json)
func (_Safe *SafeSession) SerializeUint0(objectKey string, valueKey string, values []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeUint0(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SerializeUint0 is a paid mutator transaction binding the contract method 0xfee9a469.
//
// Solidity: function serializeUint(string objectKey, string valueKey, uint256[] values) returns(string json)
func (_Safe *SafeTransactorSession) SerializeUint0(objectKey string, valueKey string, values []*big.Int) (*types.Transaction, error) {
	return _Safe.Contract.SerializeUint0(&_Safe.TransactOpts, objectKey, valueKey, values)
}

// SetEnv is a paid mutator transaction binding the contract method 0x3d5923ee.
//
// Solidity: function setEnv(string name, string value) returns()
func (_Safe *SafeTransactor) SetEnv(opts *bind.TransactOpts, name string, value string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "setEnv", name, value)
}

// SetEnv is a paid mutator transaction binding the contract method 0x3d5923ee.
//
// Solidity: function setEnv(string name, string value) returns()
func (_Safe *SafeSession) SetEnv(name string, value string) (*types.Transaction, error) {
	return _Safe.Contract.SetEnv(&_Safe.TransactOpts, name, value)
}

// SetEnv is a paid mutator transaction binding the contract method 0x3d5923ee.
//
// Solidity: function setEnv(string name, string value) returns()
func (_Safe *SafeTransactorSession) SetEnv(name string, value string) (*types.Transaction, error) {
	return _Safe.Contract.SetEnv(&_Safe.TransactOpts, name, value)
}

// StartBroadcast is a paid mutator transaction binding the contract method 0x7fb5297f.
//
// Solidity: function startBroadcast() returns()
func (_Safe *SafeTransactor) StartBroadcast(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "startBroadcast")
}

// StartBroadcast is a paid mutator transaction binding the contract method 0x7fb5297f.
//
// Solidity: function startBroadcast() returns()
func (_Safe *SafeSession) StartBroadcast() (*types.Transaction, error) {
	return _Safe.Contract.StartBroadcast(&_Safe.TransactOpts)
}

// StartBroadcast is a paid mutator transaction binding the contract method 0x7fb5297f.
//
// Solidity: function startBroadcast() returns()
func (_Safe *SafeTransactorSession) StartBroadcast() (*types.Transaction, error) {
	return _Safe.Contract.StartBroadcast(&_Safe.TransactOpts)
}

// StartBroadcast0 is a paid mutator transaction binding the contract method 0x7fec2a8d.
//
// Solidity: function startBroadcast(address signer) returns()
func (_Safe *SafeTransactor) StartBroadcast0(opts *bind.TransactOpts, signer common.Address) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "startBroadcast0", signer)
}

// StartBroadcast0 is a paid mutator transaction binding the contract method 0x7fec2a8d.
//
// Solidity: function startBroadcast(address signer) returns()
func (_Safe *SafeSession) StartBroadcast0(signer common.Address) (*types.Transaction, error) {
	return _Safe.Contract.StartBroadcast0(&_Safe.TransactOpts, signer)
}

// StartBroadcast0 is a paid mutator transaction binding the contract method 0x7fec2a8d.
//
// Solidity: function startBroadcast(address signer) returns()
func (_Safe *SafeTransactorSession) StartBroadcast0(signer common.Address) (*types.Transaction, error) {
	return _Safe.Contract.StartBroadcast0(&_Safe.TransactOpts, signer)
}

// StartBroadcast1 is a paid mutator transaction binding the contract method 0xce817d47.
//
// Solidity: function startBroadcast(uint256 privateKey) returns()
func (_Safe *SafeTransactor) StartBroadcast1(opts *bind.TransactOpts, privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "startBroadcast1", privateKey)
}

// StartBroadcast1 is a paid mutator transaction binding the contract method 0xce817d47.
//
// Solidity: function startBroadcast(uint256 privateKey) returns()
func (_Safe *SafeSession) StartBroadcast1(privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.StartBroadcast1(&_Safe.TransactOpts, privateKey)
}

// StartBroadcast1 is a paid mutator transaction binding the contract method 0xce817d47.
//
// Solidity: function startBroadcast(uint256 privateKey) returns()
func (_Safe *SafeTransactorSession) StartBroadcast1(privateKey *big.Int) (*types.Transaction, error) {
	return _Safe.Contract.StartBroadcast1(&_Safe.TransactOpts, privateKey)
}

// StopBroadcast is a paid mutator transaction binding the contract method 0x76eadd36.
//
// Solidity: function stopBroadcast() returns()
func (_Safe *SafeTransactor) StopBroadcast(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "stopBroadcast")
}

// StopBroadcast is a paid mutator transaction binding the contract method 0x76eadd36.
//
// Solidity: function stopBroadcast() returns()
func (_Safe *SafeSession) StopBroadcast() (*types.Transaction, error) {
	return _Safe.Contract.StopBroadcast(&_Safe.TransactOpts)
}

// StopBroadcast is a paid mutator transaction binding the contract method 0x76eadd36.
//
// Solidity: function stopBroadcast() returns()
func (_Safe *SafeTransactorSession) StopBroadcast() (*types.Transaction, error) {
	return _Safe.Contract.StopBroadcast(&_Safe.TransactOpts)
}

// WriteFile is a paid mutator transaction binding the contract method 0x897e0a97.
//
// Solidity: function writeFile(string path, string data) returns()
func (_Safe *SafeTransactor) WriteFile(opts *bind.TransactOpts, path string, data string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "writeFile", path, data)
}

// WriteFile is a paid mutator transaction binding the contract method 0x897e0a97.
//
// Solidity: function writeFile(string path, string data) returns()
func (_Safe *SafeSession) WriteFile(path string, data string) (*types.Transaction, error) {
	return _Safe.Contract.WriteFile(&_Safe.TransactOpts, path, data)
}

// WriteFile is a paid mutator transaction binding the contract method 0x897e0a97.
//
// Solidity: function writeFile(string path, string data) returns()
func (_Safe *SafeTransactorSession) WriteFile(path string, data string) (*types.Transaction, error) {
	return _Safe.Contract.WriteFile(&_Safe.TransactOpts, path, data)
}

// WriteFileBinary is a paid mutator transaction binding the contract method 0x1f21fc80.
//
// Solidity: function writeFileBinary(string path, bytes data) returns()
func (_Safe *SafeTransactor) WriteFileBinary(opts *bind.TransactOpts, path string, data []byte) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "writeFileBinary", path, data)
}

// WriteFileBinary is a paid mutator transaction binding the contract method 0x1f21fc80.
//
// Solidity: function writeFileBinary(string path, bytes data) returns()
func (_Safe *SafeSession) WriteFileBinary(path string, data []byte) (*types.Transaction, error) {
	return _Safe.Contract.WriteFileBinary(&_Safe.TransactOpts, path, data)
}

// WriteFileBinary is a paid mutator transaction binding the contract method 0x1f21fc80.
//
// Solidity: function writeFileBinary(string path, bytes data) returns()
func (_Safe *SafeTransactorSession) WriteFileBinary(path string, data []byte) (*types.Transaction, error) {
	return _Safe.Contract.WriteFileBinary(&_Safe.TransactOpts, path, data)
}

// WriteJson is a paid mutator transaction binding the contract method 0x35d6ad46.
//
// Solidity: function writeJson(string json, string path, string valueKey) returns()
func (_Safe *SafeTransactor) WriteJson(opts *bind.TransactOpts, json string, path string, valueKey string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "writeJson", json, path, valueKey)
}

// WriteJson is a paid mutator transaction binding the contract method 0x35d6ad46.
//
// Solidity: function writeJson(string json, string path, string valueKey) returns()
func (_Safe *SafeSession) WriteJson(json string, path string, valueKey string) (*types.Transaction, error) {
	return _Safe.Contract.WriteJson(&_Safe.TransactOpts, json, path, valueKey)
}

// WriteJson is a paid mutator transaction binding the contract method 0x35d6ad46.
//
// Solidity: function writeJson(string json, string path, string valueKey) returns()
func (_Safe *SafeTransactorSession) WriteJson(json string, path string, valueKey string) (*types.Transaction, error) {
	return _Safe.Contract.WriteJson(&_Safe.TransactOpts, json, path, valueKey)
}

// WriteJson0 is a paid mutator transaction binding the contract method 0xe23cd19f.
//
// Solidity: function writeJson(string json, string path) returns()
func (_Safe *SafeTransactor) WriteJson0(opts *bind.TransactOpts, json string, path string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "writeJson0", json, path)
}

// WriteJson0 is a paid mutator transaction binding the contract method 0xe23cd19f.
//
// Solidity: function writeJson(string json, string path) returns()
func (_Safe *SafeSession) WriteJson0(json string, path string) (*types.Transaction, error) {
	return _Safe.Contract.WriteJson0(&_Safe.TransactOpts, json, path)
}

// WriteJson0 is a paid mutator transaction binding the contract method 0xe23cd19f.
//
// Solidity: function writeJson(string json, string path) returns()
func (_Safe *SafeTransactorSession) WriteJson0(json string, path string) (*types.Transaction, error) {
	return _Safe.Contract.WriteJson0(&_Safe.TransactOpts, json, path)
}

// WriteLine is a paid mutator transaction binding the contract method 0x619d897f.
//
// Solidity: function writeLine(string path, string data) returns()
func (_Safe *SafeTransactor) WriteLine(opts *bind.TransactOpts, path string, data string) (*types.Transaction, error) {
	return _Safe.contract.Transact(opts, "writeLine", path, data)
}

// WriteLine is a paid mutator transaction binding the contract method 0x619d897f.
//
// Solidity: function writeLine(string path, string data) returns()
func (_Safe *SafeSession) WriteLine(path string, data string) (*types.Transaction, error) {
	return _Safe.Contract.WriteLine(&_Safe.TransactOpts, path, data)
}

// WriteLine is a paid mutator transaction binding the contract method 0x619d897f.
//
// Solidity: function writeLine(string path, string data) returns()
func (_Safe *SafeTransactorSession) WriteLine(path string, data string) (*types.Transaction, error) {
	return _Safe.Contract.WriteLine(&_Safe.TransactOpts, path, data)
}
