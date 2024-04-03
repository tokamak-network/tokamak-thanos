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

// Attestation is an auto generated low-level Go binding around an user-defined struct.
type Attestation struct {
	Uid            [32]byte
	Schema         [32]byte
	Time           uint64
	ExpirationTime uint64
	RevocationTime uint64
	RefUID         [32]byte
	Recipient      common.Address
	Attester       common.Address
	Revocable      bool
	Data           []byte
}

// AttestationRequest is an auto generated low-level Go binding around an user-defined struct.
type AttestationRequest struct {
	Schema [32]byte
	Data   AttestationRequestData
}

// AttestationRequestData is an auto generated low-level Go binding around an user-defined struct.
type AttestationRequestData struct {
	Recipient      common.Address
	ExpirationTime uint64
	Revocable      bool
	RefUID         [32]byte
	Data           []byte
	Value          *big.Int
}

// DelegatedAttestationRequest is an auto generated low-level Go binding around an user-defined struct.
type DelegatedAttestationRequest struct {
	Schema    [32]byte
	Data      AttestationRequestData
	Signature Signature
	Attester  common.Address
	Deadline  uint64
}

// DelegatedRevocationRequest is an auto generated low-level Go binding around an user-defined struct.
type DelegatedRevocationRequest struct {
	Schema    [32]byte
	Data      RevocationRequestData
	Signature Signature
	Revoker   common.Address
	Deadline  uint64
}

// MultiAttestationRequest is an auto generated low-level Go binding around an user-defined struct.
type MultiAttestationRequest struct {
	Schema [32]byte
	Data   []AttestationRequestData
}

// MultiDelegatedAttestationRequest is an auto generated low-level Go binding around an user-defined struct.
type MultiDelegatedAttestationRequest struct {
	Schema     [32]byte
	Data       []AttestationRequestData
	Signatures []Signature
	Attester   common.Address
	Deadline   uint64
}

// MultiDelegatedRevocationRequest is an auto generated low-level Go binding around an user-defined struct.
type MultiDelegatedRevocationRequest struct {
	Schema     [32]byte
	Data       []RevocationRequestData
	Signatures []Signature
	Revoker    common.Address
	Deadline   uint64
}

// MultiRevocationRequest is an auto generated low-level Go binding around an user-defined struct.
type MultiRevocationRequest struct {
	Schema [32]byte
	Data   []RevocationRequestData
}

// RevocationRequest is an auto generated low-level Go binding around an user-defined struct.
type RevocationRequest struct {
	Schema [32]byte
	Data   RevocationRequestData
}

// RevocationRequestData is an auto generated low-level Go binding around an user-defined struct.
type RevocationRequestData struct {
	Uid   [32]byte
	Value *big.Int
}

// Signature is an auto generated low-level Go binding around an user-defined struct.
type Signature struct {
	V uint8
	R [32]byte
	S [32]byte
}

// EASMetaData contains all meta data concerning the EAS contract.
var EASMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"attester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"schemaUID\",\"type\":\"bytes32\"}],\"name\":\"Attested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"attester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"schemaUID\",\"type\":\"bytes32\"}],\"name\":\"Revoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"revoker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"RevokedOffchain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"Timestamped\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"refUID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structAttestationRequestData\",\"name\":\"data\",\"type\":\"tuple\"}],\"internalType\":\"structAttestationRequest\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"attest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"refUID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structAttestationRequestData\",\"name\":\"data\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structSignature\",\"name\":\"signature\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"attester\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"}],\"internalType\":\"structDelegatedAttestationRequest\",\"name\":\"delegatedRequest\",\"type\":\"tuple\"}],\"name\":\"attestByDelegation\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"}],\"name\":\"getAttestation\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"time\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"revocationTime\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"refUID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"attester\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structAttestation\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"revoker\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"getRevokeOffchain\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSchemaRegistry\",\"outputs\":[{\"internalType\":\"contractISchemaRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"}],\"name\":\"isAttestationValid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"refUID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structAttestationRequestData[]\",\"name\":\"data\",\"type\":\"tuple[]\"}],\"internalType\":\"structMultiAttestationRequest[]\",\"name\":\"multiRequests\",\"type\":\"tuple[]\"}],\"name\":\"multiAttest\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"refUID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structAttestationRequestData[]\",\"name\":\"data\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structSignature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"attester\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"}],\"internalType\":\"structMultiDelegatedAttestationRequest[]\",\"name\":\"multiDelegatedRequests\",\"type\":\"tuple[]\"}],\"name\":\"multiAttestByDelegation\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structRevocationRequestData[]\",\"name\":\"data\",\"type\":\"tuple[]\"}],\"internalType\":\"structMultiRevocationRequest[]\",\"name\":\"multiRequests\",\"type\":\"tuple[]\"}],\"name\":\"multiRevoke\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structRevocationRequestData[]\",\"name\":\"data\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structSignature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"revoker\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"}],\"internalType\":\"structMultiDelegatedRevocationRequest[]\",\"name\":\"multiDelegatedRequests\",\"type\":\"tuple[]\"}],\"name\":\"multiRevokeByDelegation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"data\",\"type\":\"bytes32[]\"}],\"name\":\"multiRevokeOffchain\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"data\",\"type\":\"bytes32[]\"}],\"name\":\"multiTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structRevocationRequestData\",\"name\":\"data\",\"type\":\"tuple\"}],\"internalType\":\"structRevocationRequest\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"revoke\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"schema\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structRevocationRequestData\",\"name\":\"data\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structSignature\",\"name\":\"signature\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"revoker\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"}],\"internalType\":\"structDelegatedRevocationRequest\",\"name\":\"delegatedRequest\",\"type\":\"tuple\"}],\"name\":\"revokeByDelegation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"revokeOffchain\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"timestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EASABI is the input ABI used to generate the binding from.
// Deprecated: Use EASMetaData.ABI instead.
var EASABI = EASMetaData.ABI

// EAS is an auto generated Go binding around an Ethereum contract.
type EAS struct {
	EASCaller     // Read-only binding to the contract
	EASTransactor // Write-only binding to the contract
	EASFilterer   // Log filterer for contract events
}

// EASCaller is an auto generated read-only Go binding around an Ethereum contract.
type EASCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EASTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EASTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EASFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EASFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EASSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EASSession struct {
	Contract     *EAS              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EASCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EASCallerSession struct {
	Contract *EASCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EASTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EASTransactorSession struct {
	Contract     *EASTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EASRaw is an auto generated low-level Go binding around an Ethereum contract.
type EASRaw struct {
	Contract *EAS // Generic contract binding to access the raw methods on
}

// EASCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EASCallerRaw struct {
	Contract *EASCaller // Generic read-only contract binding to access the raw methods on
}

// EASTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EASTransactorRaw struct {
	Contract *EASTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEAS creates a new instance of EAS, bound to a specific deployed contract.
func NewEAS(address common.Address, backend bind.ContractBackend) (*EAS, error) {
	contract, err := bindEAS(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EAS{EASCaller: EASCaller{contract: contract}, EASTransactor: EASTransactor{contract: contract}, EASFilterer: EASFilterer{contract: contract}}, nil
}

// NewEASCaller creates a new read-only instance of EAS, bound to a specific deployed contract.
func NewEASCaller(address common.Address, caller bind.ContractCaller) (*EASCaller, error) {
	contract, err := bindEAS(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EASCaller{contract: contract}, nil
}

// NewEASTransactor creates a new write-only instance of EAS, bound to a specific deployed contract.
func NewEASTransactor(address common.Address, transactor bind.ContractTransactor) (*EASTransactor, error) {
	contract, err := bindEAS(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EASTransactor{contract: contract}, nil
}

// NewEASFilterer creates a new log filterer instance of EAS, bound to a specific deployed contract.
func NewEASFilterer(address common.Address, filterer bind.ContractFilterer) (*EASFilterer, error) {
	contract, err := bindEAS(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EASFilterer{contract: contract}, nil
}

// bindEAS binds a generic wrapper to an already deployed contract.
func bindEAS(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EASMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EAS *EASRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EAS.Contract.EASCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EAS *EASRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EAS.Contract.EASTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EAS *EASRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EAS.Contract.EASTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EAS *EASCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EAS.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EAS *EASTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EAS.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EAS *EASTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EAS.Contract.contract.Transact(opts, method, params...)
}

// GetAttestation is a free data retrieval call binding the contract method 0xa3112a64.
//
// Solidity: function getAttestation(bytes32 uid) view returns((bytes32,bytes32,uint64,uint64,uint64,bytes32,address,address,bool,bytes))
func (_EAS *EASCaller) GetAttestation(opts *bind.CallOpts, uid [32]byte) (Attestation, error) {
	var out []interface{}
	err := _EAS.contract.Call(opts, &out, "getAttestation", uid)

	if err != nil {
		return *new(Attestation), err
	}

	out0 := *abi.ConvertType(out[0], new(Attestation)).(*Attestation)

	return out0, err

}

// GetAttestation is a free data retrieval call binding the contract method 0xa3112a64.
//
// Solidity: function getAttestation(bytes32 uid) view returns((bytes32,bytes32,uint64,uint64,uint64,bytes32,address,address,bool,bytes))
func (_EAS *EASSession) GetAttestation(uid [32]byte) (Attestation, error) {
	return _EAS.Contract.GetAttestation(&_EAS.CallOpts, uid)
}

// GetAttestation is a free data retrieval call binding the contract method 0xa3112a64.
//
// Solidity: function getAttestation(bytes32 uid) view returns((bytes32,bytes32,uint64,uint64,uint64,bytes32,address,address,bool,bytes))
func (_EAS *EASCallerSession) GetAttestation(uid [32]byte) (Attestation, error) {
	return _EAS.Contract.GetAttestation(&_EAS.CallOpts, uid)
}

// GetRevokeOffchain is a free data retrieval call binding the contract method 0xb469318d.
//
// Solidity: function getRevokeOffchain(address revoker, bytes32 data) view returns(uint64)
func (_EAS *EASCaller) GetRevokeOffchain(opts *bind.CallOpts, revoker common.Address, data [32]byte) (uint64, error) {
	var out []interface{}
	err := _EAS.contract.Call(opts, &out, "getRevokeOffchain", revoker, data)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetRevokeOffchain is a free data retrieval call binding the contract method 0xb469318d.
//
// Solidity: function getRevokeOffchain(address revoker, bytes32 data) view returns(uint64)
func (_EAS *EASSession) GetRevokeOffchain(revoker common.Address, data [32]byte) (uint64, error) {
	return _EAS.Contract.GetRevokeOffchain(&_EAS.CallOpts, revoker, data)
}

// GetRevokeOffchain is a free data retrieval call binding the contract method 0xb469318d.
//
// Solidity: function getRevokeOffchain(address revoker, bytes32 data) view returns(uint64)
func (_EAS *EASCallerSession) GetRevokeOffchain(revoker common.Address, data [32]byte) (uint64, error) {
	return _EAS.Contract.GetRevokeOffchain(&_EAS.CallOpts, revoker, data)
}

// GetSchemaRegistry is a free data retrieval call binding the contract method 0xf10b5cc8.
//
// Solidity: function getSchemaRegistry() view returns(address)
func (_EAS *EASCaller) GetSchemaRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EAS.contract.Call(opts, &out, "getSchemaRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSchemaRegistry is a free data retrieval call binding the contract method 0xf10b5cc8.
//
// Solidity: function getSchemaRegistry() view returns(address)
func (_EAS *EASSession) GetSchemaRegistry() (common.Address, error) {
	return _EAS.Contract.GetSchemaRegistry(&_EAS.CallOpts)
}

// GetSchemaRegistry is a free data retrieval call binding the contract method 0xf10b5cc8.
//
// Solidity: function getSchemaRegistry() view returns(address)
func (_EAS *EASCallerSession) GetSchemaRegistry() (common.Address, error) {
	return _EAS.Contract.GetSchemaRegistry(&_EAS.CallOpts)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xd45c4435.
//
// Solidity: function getTimestamp(bytes32 data) view returns(uint64)
func (_EAS *EASCaller) GetTimestamp(opts *bind.CallOpts, data [32]byte) (uint64, error) {
	var out []interface{}
	err := _EAS.contract.Call(opts, &out, "getTimestamp", data)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xd45c4435.
//
// Solidity: function getTimestamp(bytes32 data) view returns(uint64)
func (_EAS *EASSession) GetTimestamp(data [32]byte) (uint64, error) {
	return _EAS.Contract.GetTimestamp(&_EAS.CallOpts, data)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xd45c4435.
//
// Solidity: function getTimestamp(bytes32 data) view returns(uint64)
func (_EAS *EASCallerSession) GetTimestamp(data [32]byte) (uint64, error) {
	return _EAS.Contract.GetTimestamp(&_EAS.CallOpts, data)
}

// IsAttestationValid is a free data retrieval call binding the contract method 0xe30bb563.
//
// Solidity: function isAttestationValid(bytes32 uid) view returns(bool)
func (_EAS *EASCaller) IsAttestationValid(opts *bind.CallOpts, uid [32]byte) (bool, error) {
	var out []interface{}
	err := _EAS.contract.Call(opts, &out, "isAttestationValid", uid)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAttestationValid is a free data retrieval call binding the contract method 0xe30bb563.
//
// Solidity: function isAttestationValid(bytes32 uid) view returns(bool)
func (_EAS *EASSession) IsAttestationValid(uid [32]byte) (bool, error) {
	return _EAS.Contract.IsAttestationValid(&_EAS.CallOpts, uid)
}

// IsAttestationValid is a free data retrieval call binding the contract method 0xe30bb563.
//
// Solidity: function isAttestationValid(bytes32 uid) view returns(bool)
func (_EAS *EASCallerSession) IsAttestationValid(uid [32]byte) (bool, error) {
	return _EAS.Contract.IsAttestationValid(&_EAS.CallOpts, uid)
}

// Attest is a paid mutator transaction binding the contract method 0xf17325e7.
//
// Solidity: function attest((bytes32,(address,uint64,bool,bytes32,bytes,uint256)) request) payable returns(bytes32)
func (_EAS *EASTransactor) Attest(opts *bind.TransactOpts, request AttestationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "attest", request)
}

// Attest is a paid mutator transaction binding the contract method 0xf17325e7.
//
// Solidity: function attest((bytes32,(address,uint64,bool,bytes32,bytes,uint256)) request) payable returns(bytes32)
func (_EAS *EASSession) Attest(request AttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.Attest(&_EAS.TransactOpts, request)
}

// Attest is a paid mutator transaction binding the contract method 0xf17325e7.
//
// Solidity: function attest((bytes32,(address,uint64,bool,bytes32,bytes,uint256)) request) payable returns(bytes32)
func (_EAS *EASTransactorSession) Attest(request AttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.Attest(&_EAS.TransactOpts, request)
}

// AttestByDelegation is a paid mutator transaction binding the contract method 0x3c042715.
//
// Solidity: function attestByDelegation((bytes32,(address,uint64,bool,bytes32,bytes,uint256),(uint8,bytes32,bytes32),address,uint64) delegatedRequest) payable returns(bytes32)
func (_EAS *EASTransactor) AttestByDelegation(opts *bind.TransactOpts, delegatedRequest DelegatedAttestationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "attestByDelegation", delegatedRequest)
}

// AttestByDelegation is a paid mutator transaction binding the contract method 0x3c042715.
//
// Solidity: function attestByDelegation((bytes32,(address,uint64,bool,bytes32,bytes,uint256),(uint8,bytes32,bytes32),address,uint64) delegatedRequest) payable returns(bytes32)
func (_EAS *EASSession) AttestByDelegation(delegatedRequest DelegatedAttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.AttestByDelegation(&_EAS.TransactOpts, delegatedRequest)
}

// AttestByDelegation is a paid mutator transaction binding the contract method 0x3c042715.
//
// Solidity: function attestByDelegation((bytes32,(address,uint64,bool,bytes32,bytes,uint256),(uint8,bytes32,bytes32),address,uint64) delegatedRequest) payable returns(bytes32)
func (_EAS *EASTransactorSession) AttestByDelegation(delegatedRequest DelegatedAttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.AttestByDelegation(&_EAS.TransactOpts, delegatedRequest)
}

// MultiAttest is a paid mutator transaction binding the contract method 0x44adc90e.
//
// Solidity: function multiAttest((bytes32,(address,uint64,bool,bytes32,bytes,uint256)[])[] multiRequests) payable returns(bytes32[])
func (_EAS *EASTransactor) MultiAttest(opts *bind.TransactOpts, multiRequests []MultiAttestationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "multiAttest", multiRequests)
}

// MultiAttest is a paid mutator transaction binding the contract method 0x44adc90e.
//
// Solidity: function multiAttest((bytes32,(address,uint64,bool,bytes32,bytes,uint256)[])[] multiRequests) payable returns(bytes32[])
func (_EAS *EASSession) MultiAttest(multiRequests []MultiAttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiAttest(&_EAS.TransactOpts, multiRequests)
}

// MultiAttest is a paid mutator transaction binding the contract method 0x44adc90e.
//
// Solidity: function multiAttest((bytes32,(address,uint64,bool,bytes32,bytes,uint256)[])[] multiRequests) payable returns(bytes32[])
func (_EAS *EASTransactorSession) MultiAttest(multiRequests []MultiAttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiAttest(&_EAS.TransactOpts, multiRequests)
}

// MultiAttestByDelegation is a paid mutator transaction binding the contract method 0x95411525.
//
// Solidity: function multiAttestByDelegation((bytes32,(address,uint64,bool,bytes32,bytes,uint256)[],(uint8,bytes32,bytes32)[],address,uint64)[] multiDelegatedRequests) payable returns(bytes32[])
func (_EAS *EASTransactor) MultiAttestByDelegation(opts *bind.TransactOpts, multiDelegatedRequests []MultiDelegatedAttestationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "multiAttestByDelegation", multiDelegatedRequests)
}

// MultiAttestByDelegation is a paid mutator transaction binding the contract method 0x95411525.
//
// Solidity: function multiAttestByDelegation((bytes32,(address,uint64,bool,bytes32,bytes,uint256)[],(uint8,bytes32,bytes32)[],address,uint64)[] multiDelegatedRequests) payable returns(bytes32[])
func (_EAS *EASSession) MultiAttestByDelegation(multiDelegatedRequests []MultiDelegatedAttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiAttestByDelegation(&_EAS.TransactOpts, multiDelegatedRequests)
}

// MultiAttestByDelegation is a paid mutator transaction binding the contract method 0x95411525.
//
// Solidity: function multiAttestByDelegation((bytes32,(address,uint64,bool,bytes32,bytes,uint256)[],(uint8,bytes32,bytes32)[],address,uint64)[] multiDelegatedRequests) payable returns(bytes32[])
func (_EAS *EASTransactorSession) MultiAttestByDelegation(multiDelegatedRequests []MultiDelegatedAttestationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiAttestByDelegation(&_EAS.TransactOpts, multiDelegatedRequests)
}

// MultiRevoke is a paid mutator transaction binding the contract method 0x4cb7e9e5.
//
// Solidity: function multiRevoke((bytes32,(bytes32,uint256)[])[] multiRequests) payable returns()
func (_EAS *EASTransactor) MultiRevoke(opts *bind.TransactOpts, multiRequests []MultiRevocationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "multiRevoke", multiRequests)
}

// MultiRevoke is a paid mutator transaction binding the contract method 0x4cb7e9e5.
//
// Solidity: function multiRevoke((bytes32,(bytes32,uint256)[])[] multiRequests) payable returns()
func (_EAS *EASSession) MultiRevoke(multiRequests []MultiRevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiRevoke(&_EAS.TransactOpts, multiRequests)
}

// MultiRevoke is a paid mutator transaction binding the contract method 0x4cb7e9e5.
//
// Solidity: function multiRevoke((bytes32,(bytes32,uint256)[])[] multiRequests) payable returns()
func (_EAS *EASTransactorSession) MultiRevoke(multiRequests []MultiRevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiRevoke(&_EAS.TransactOpts, multiRequests)
}

// MultiRevokeByDelegation is a paid mutator transaction binding the contract method 0x0eabf660.
//
// Solidity: function multiRevokeByDelegation((bytes32,(bytes32,uint256)[],(uint8,bytes32,bytes32)[],address,uint64)[] multiDelegatedRequests) payable returns()
func (_EAS *EASTransactor) MultiRevokeByDelegation(opts *bind.TransactOpts, multiDelegatedRequests []MultiDelegatedRevocationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "multiRevokeByDelegation", multiDelegatedRequests)
}

// MultiRevokeByDelegation is a paid mutator transaction binding the contract method 0x0eabf660.
//
// Solidity: function multiRevokeByDelegation((bytes32,(bytes32,uint256)[],(uint8,bytes32,bytes32)[],address,uint64)[] multiDelegatedRequests) payable returns()
func (_EAS *EASSession) MultiRevokeByDelegation(multiDelegatedRequests []MultiDelegatedRevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiRevokeByDelegation(&_EAS.TransactOpts, multiDelegatedRequests)
}

// MultiRevokeByDelegation is a paid mutator transaction binding the contract method 0x0eabf660.
//
// Solidity: function multiRevokeByDelegation((bytes32,(bytes32,uint256)[],(uint8,bytes32,bytes32)[],address,uint64)[] multiDelegatedRequests) payable returns()
func (_EAS *EASTransactorSession) MultiRevokeByDelegation(multiDelegatedRequests []MultiDelegatedRevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.MultiRevokeByDelegation(&_EAS.TransactOpts, multiDelegatedRequests)
}

// MultiRevokeOffchain is a paid mutator transaction binding the contract method 0x13893f61.
//
// Solidity: function multiRevokeOffchain(bytes32[] data) returns(uint64)
func (_EAS *EASTransactor) MultiRevokeOffchain(opts *bind.TransactOpts, data [][32]byte) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "multiRevokeOffchain", data)
}

// MultiRevokeOffchain is a paid mutator transaction binding the contract method 0x13893f61.
//
// Solidity: function multiRevokeOffchain(bytes32[] data) returns(uint64)
func (_EAS *EASSession) MultiRevokeOffchain(data [][32]byte) (*types.Transaction, error) {
	return _EAS.Contract.MultiRevokeOffchain(&_EAS.TransactOpts, data)
}

// MultiRevokeOffchain is a paid mutator transaction binding the contract method 0x13893f61.
//
// Solidity: function multiRevokeOffchain(bytes32[] data) returns(uint64)
func (_EAS *EASTransactorSession) MultiRevokeOffchain(data [][32]byte) (*types.Transaction, error) {
	return _EAS.Contract.MultiRevokeOffchain(&_EAS.TransactOpts, data)
}

// MultiTimestamp is a paid mutator transaction binding the contract method 0xe71ff365.
//
// Solidity: function multiTimestamp(bytes32[] data) returns(uint64)
func (_EAS *EASTransactor) MultiTimestamp(opts *bind.TransactOpts, data [][32]byte) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "multiTimestamp", data)
}

// MultiTimestamp is a paid mutator transaction binding the contract method 0xe71ff365.
//
// Solidity: function multiTimestamp(bytes32[] data) returns(uint64)
func (_EAS *EASSession) MultiTimestamp(data [][32]byte) (*types.Transaction, error) {
	return _EAS.Contract.MultiTimestamp(&_EAS.TransactOpts, data)
}

// MultiTimestamp is a paid mutator transaction binding the contract method 0xe71ff365.
//
// Solidity: function multiTimestamp(bytes32[] data) returns(uint64)
func (_EAS *EASTransactorSession) MultiTimestamp(data [][32]byte) (*types.Transaction, error) {
	return _EAS.Contract.MultiTimestamp(&_EAS.TransactOpts, data)
}

// Revoke is a paid mutator transaction binding the contract method 0x46926267.
//
// Solidity: function revoke((bytes32,(bytes32,uint256)) request) payable returns()
func (_EAS *EASTransactor) Revoke(opts *bind.TransactOpts, request RevocationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "revoke", request)
}

// Revoke is a paid mutator transaction binding the contract method 0x46926267.
//
// Solidity: function revoke((bytes32,(bytes32,uint256)) request) payable returns()
func (_EAS *EASSession) Revoke(request RevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.Revoke(&_EAS.TransactOpts, request)
}

// Revoke is a paid mutator transaction binding the contract method 0x46926267.
//
// Solidity: function revoke((bytes32,(bytes32,uint256)) request) payable returns()
func (_EAS *EASTransactorSession) Revoke(request RevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.Revoke(&_EAS.TransactOpts, request)
}

// RevokeByDelegation is a paid mutator transaction binding the contract method 0xa6d4dbc7.
//
// Solidity: function revokeByDelegation((bytes32,(bytes32,uint256),(uint8,bytes32,bytes32),address,uint64) delegatedRequest) payable returns()
func (_EAS *EASTransactor) RevokeByDelegation(opts *bind.TransactOpts, delegatedRequest DelegatedRevocationRequest) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "revokeByDelegation", delegatedRequest)
}

// RevokeByDelegation is a paid mutator transaction binding the contract method 0xa6d4dbc7.
//
// Solidity: function revokeByDelegation((bytes32,(bytes32,uint256),(uint8,bytes32,bytes32),address,uint64) delegatedRequest) payable returns()
func (_EAS *EASSession) RevokeByDelegation(delegatedRequest DelegatedRevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.RevokeByDelegation(&_EAS.TransactOpts, delegatedRequest)
}

// RevokeByDelegation is a paid mutator transaction binding the contract method 0xa6d4dbc7.
//
// Solidity: function revokeByDelegation((bytes32,(bytes32,uint256),(uint8,bytes32,bytes32),address,uint64) delegatedRequest) payable returns()
func (_EAS *EASTransactorSession) RevokeByDelegation(delegatedRequest DelegatedRevocationRequest) (*types.Transaction, error) {
	return _EAS.Contract.RevokeByDelegation(&_EAS.TransactOpts, delegatedRequest)
}

// RevokeOffchain is a paid mutator transaction binding the contract method 0xcf190f34.
//
// Solidity: function revokeOffchain(bytes32 data) returns(uint64)
func (_EAS *EASTransactor) RevokeOffchain(opts *bind.TransactOpts, data [32]byte) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "revokeOffchain", data)
}

// RevokeOffchain is a paid mutator transaction binding the contract method 0xcf190f34.
//
// Solidity: function revokeOffchain(bytes32 data) returns(uint64)
func (_EAS *EASSession) RevokeOffchain(data [32]byte) (*types.Transaction, error) {
	return _EAS.Contract.RevokeOffchain(&_EAS.TransactOpts, data)
}

// RevokeOffchain is a paid mutator transaction binding the contract method 0xcf190f34.
//
// Solidity: function revokeOffchain(bytes32 data) returns(uint64)
func (_EAS *EASTransactorSession) RevokeOffchain(data [32]byte) (*types.Transaction, error) {
	return _EAS.Contract.RevokeOffchain(&_EAS.TransactOpts, data)
}

// Timestamp is a paid mutator transaction binding the contract method 0x4d003070.
//
// Solidity: function timestamp(bytes32 data) returns(uint64)
func (_EAS *EASTransactor) Timestamp(opts *bind.TransactOpts, data [32]byte) (*types.Transaction, error) {
	return _EAS.contract.Transact(opts, "timestamp", data)
}

// Timestamp is a paid mutator transaction binding the contract method 0x4d003070.
//
// Solidity: function timestamp(bytes32 data) returns(uint64)
func (_EAS *EASSession) Timestamp(data [32]byte) (*types.Transaction, error) {
	return _EAS.Contract.Timestamp(&_EAS.TransactOpts, data)
}

// Timestamp is a paid mutator transaction binding the contract method 0x4d003070.
//
// Solidity: function timestamp(bytes32 data) returns(uint64)
func (_EAS *EASTransactorSession) Timestamp(data [32]byte) (*types.Transaction, error) {
	return _EAS.Contract.Timestamp(&_EAS.TransactOpts, data)
}

// EASAttestedIterator is returned from FilterAttested and is used to iterate over the raw logs and unpacked data for Attested events raised by the EAS contract.
type EASAttestedIterator struct {
	Event *EASAttested // Event containing the contract specifics and raw log

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
func (it *EASAttestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EASAttested)
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
		it.Event = new(EASAttested)
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
func (it *EASAttestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EASAttestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EASAttested represents a Attested event raised by the EAS contract.
type EASAttested struct {
	Recipient common.Address
	Attester  common.Address
	Uid       [32]byte
	SchemaUID [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAttested is a free log retrieval operation binding the contract event 0x8bf46bf4cfd674fa735a3d63ec1c9ad4153f033c290341f3a588b75685141b35.
//
// Solidity: event Attested(address indexed recipient, address indexed attester, bytes32 uid, bytes32 indexed schemaUID)
func (_EAS *EASFilterer) FilterAttested(opts *bind.FilterOpts, recipient []common.Address, attester []common.Address, schemaUID [][32]byte) (*EASAttestedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}

	var schemaUIDRule []interface{}
	for _, schemaUIDItem := range schemaUID {
		schemaUIDRule = append(schemaUIDRule, schemaUIDItem)
	}

	logs, sub, err := _EAS.contract.FilterLogs(opts, "Attested", recipientRule, attesterRule, schemaUIDRule)
	if err != nil {
		return nil, err
	}
	return &EASAttestedIterator{contract: _EAS.contract, event: "Attested", logs: logs, sub: sub}, nil
}

// WatchAttested is a free log subscription operation binding the contract event 0x8bf46bf4cfd674fa735a3d63ec1c9ad4153f033c290341f3a588b75685141b35.
//
// Solidity: event Attested(address indexed recipient, address indexed attester, bytes32 uid, bytes32 indexed schemaUID)
func (_EAS *EASFilterer) WatchAttested(opts *bind.WatchOpts, sink chan<- *EASAttested, recipient []common.Address, attester []common.Address, schemaUID [][32]byte) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}

	var schemaUIDRule []interface{}
	for _, schemaUIDItem := range schemaUID {
		schemaUIDRule = append(schemaUIDRule, schemaUIDItem)
	}

	logs, sub, err := _EAS.contract.WatchLogs(opts, "Attested", recipientRule, attesterRule, schemaUIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EASAttested)
				if err := _EAS.contract.UnpackLog(event, "Attested", log); err != nil {
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

// ParseAttested is a log parse operation binding the contract event 0x8bf46bf4cfd674fa735a3d63ec1c9ad4153f033c290341f3a588b75685141b35.
//
// Solidity: event Attested(address indexed recipient, address indexed attester, bytes32 uid, bytes32 indexed schemaUID)
func (_EAS *EASFilterer) ParseAttested(log types.Log) (*EASAttested, error) {
	event := new(EASAttested)
	if err := _EAS.contract.UnpackLog(event, "Attested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EASRevokedIterator is returned from FilterRevoked and is used to iterate over the raw logs and unpacked data for Revoked events raised by the EAS contract.
type EASRevokedIterator struct {
	Event *EASRevoked // Event containing the contract specifics and raw log

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
func (it *EASRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EASRevoked)
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
		it.Event = new(EASRevoked)
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
func (it *EASRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EASRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EASRevoked represents a Revoked event raised by the EAS contract.
type EASRevoked struct {
	Recipient common.Address
	Attester  common.Address
	Uid       [32]byte
	SchemaUID [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRevoked is a free log retrieval operation binding the contract event 0xf930a6e2523c9cc298691873087a740550b8fc85a0680830414c148ed927f615.
//
// Solidity: event Revoked(address indexed recipient, address indexed attester, bytes32 uid, bytes32 indexed schemaUID)
func (_EAS *EASFilterer) FilterRevoked(opts *bind.FilterOpts, recipient []common.Address, attester []common.Address, schemaUID [][32]byte) (*EASRevokedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}

	var schemaUIDRule []interface{}
	for _, schemaUIDItem := range schemaUID {
		schemaUIDRule = append(schemaUIDRule, schemaUIDItem)
	}

	logs, sub, err := _EAS.contract.FilterLogs(opts, "Revoked", recipientRule, attesterRule, schemaUIDRule)
	if err != nil {
		return nil, err
	}
	return &EASRevokedIterator{contract: _EAS.contract, event: "Revoked", logs: logs, sub: sub}, nil
}

// WatchRevoked is a free log subscription operation binding the contract event 0xf930a6e2523c9cc298691873087a740550b8fc85a0680830414c148ed927f615.
//
// Solidity: event Revoked(address indexed recipient, address indexed attester, bytes32 uid, bytes32 indexed schemaUID)
func (_EAS *EASFilterer) WatchRevoked(opts *bind.WatchOpts, sink chan<- *EASRevoked, recipient []common.Address, attester []common.Address, schemaUID [][32]byte) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}

	var schemaUIDRule []interface{}
	for _, schemaUIDItem := range schemaUID {
		schemaUIDRule = append(schemaUIDRule, schemaUIDItem)
	}

	logs, sub, err := _EAS.contract.WatchLogs(opts, "Revoked", recipientRule, attesterRule, schemaUIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EASRevoked)
				if err := _EAS.contract.UnpackLog(event, "Revoked", log); err != nil {
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

// ParseRevoked is a log parse operation binding the contract event 0xf930a6e2523c9cc298691873087a740550b8fc85a0680830414c148ed927f615.
//
// Solidity: event Revoked(address indexed recipient, address indexed attester, bytes32 uid, bytes32 indexed schemaUID)
func (_EAS *EASFilterer) ParseRevoked(log types.Log) (*EASRevoked, error) {
	event := new(EASRevoked)
	if err := _EAS.contract.UnpackLog(event, "Revoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EASRevokedOffchainIterator is returned from FilterRevokedOffchain and is used to iterate over the raw logs and unpacked data for RevokedOffchain events raised by the EAS contract.
type EASRevokedOffchainIterator struct {
	Event *EASRevokedOffchain // Event containing the contract specifics and raw log

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
func (it *EASRevokedOffchainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EASRevokedOffchain)
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
		it.Event = new(EASRevokedOffchain)
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
func (it *EASRevokedOffchainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EASRevokedOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EASRevokedOffchain represents a RevokedOffchain event raised by the EAS contract.
type EASRevokedOffchain struct {
	Revoker   common.Address
	Data      [32]byte
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRevokedOffchain is a free log retrieval operation binding the contract event 0x92a1f7a41a7c585a8b09e25b195e225b1d43248daca46b0faf9e0792777a2229.
//
// Solidity: event RevokedOffchain(address indexed revoker, bytes32 indexed data, uint64 indexed timestamp)
func (_EAS *EASFilterer) FilterRevokedOffchain(opts *bind.FilterOpts, revoker []common.Address, data [][32]byte, timestamp []uint64) (*EASRevokedOffchainIterator, error) {

	var revokerRule []interface{}
	for _, revokerItem := range revoker {
		revokerRule = append(revokerRule, revokerItem)
	}
	var dataRule []interface{}
	for _, dataItem := range data {
		dataRule = append(dataRule, dataItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _EAS.contract.FilterLogs(opts, "RevokedOffchain", revokerRule, dataRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return &EASRevokedOffchainIterator{contract: _EAS.contract, event: "RevokedOffchain", logs: logs, sub: sub}, nil
}

// WatchRevokedOffchain is a free log subscription operation binding the contract event 0x92a1f7a41a7c585a8b09e25b195e225b1d43248daca46b0faf9e0792777a2229.
//
// Solidity: event RevokedOffchain(address indexed revoker, bytes32 indexed data, uint64 indexed timestamp)
func (_EAS *EASFilterer) WatchRevokedOffchain(opts *bind.WatchOpts, sink chan<- *EASRevokedOffchain, revoker []common.Address, data [][32]byte, timestamp []uint64) (event.Subscription, error) {

	var revokerRule []interface{}
	for _, revokerItem := range revoker {
		revokerRule = append(revokerRule, revokerItem)
	}
	var dataRule []interface{}
	for _, dataItem := range data {
		dataRule = append(dataRule, dataItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _EAS.contract.WatchLogs(opts, "RevokedOffchain", revokerRule, dataRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EASRevokedOffchain)
				if err := _EAS.contract.UnpackLog(event, "RevokedOffchain", log); err != nil {
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

// ParseRevokedOffchain is a log parse operation binding the contract event 0x92a1f7a41a7c585a8b09e25b195e225b1d43248daca46b0faf9e0792777a2229.
//
// Solidity: event RevokedOffchain(address indexed revoker, bytes32 indexed data, uint64 indexed timestamp)
func (_EAS *EASFilterer) ParseRevokedOffchain(log types.Log) (*EASRevokedOffchain, error) {
	event := new(EASRevokedOffchain)
	if err := _EAS.contract.UnpackLog(event, "RevokedOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EASTimestampedIterator is returned from FilterTimestamped and is used to iterate over the raw logs and unpacked data for Timestamped events raised by the EAS contract.
type EASTimestampedIterator struct {
	Event *EASTimestamped // Event containing the contract specifics and raw log

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
func (it *EASTimestampedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EASTimestamped)
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
		it.Event = new(EASTimestamped)
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
func (it *EASTimestampedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EASTimestampedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EASTimestamped represents a Timestamped event raised by the EAS contract.
type EASTimestamped struct {
	Data      [32]byte
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTimestamped is a free log retrieval operation binding the contract event 0x5aafceeb1c7ad58e4a84898bdee37c02c0fc46e7d24e6b60e8209449f183459f.
//
// Solidity: event Timestamped(bytes32 indexed data, uint64 indexed timestamp)
func (_EAS *EASFilterer) FilterTimestamped(opts *bind.FilterOpts, data [][32]byte, timestamp []uint64) (*EASTimestampedIterator, error) {

	var dataRule []interface{}
	for _, dataItem := range data {
		dataRule = append(dataRule, dataItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _EAS.contract.FilterLogs(opts, "Timestamped", dataRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return &EASTimestampedIterator{contract: _EAS.contract, event: "Timestamped", logs: logs, sub: sub}, nil
}

// WatchTimestamped is a free log subscription operation binding the contract event 0x5aafceeb1c7ad58e4a84898bdee37c02c0fc46e7d24e6b60e8209449f183459f.
//
// Solidity: event Timestamped(bytes32 indexed data, uint64 indexed timestamp)
func (_EAS *EASFilterer) WatchTimestamped(opts *bind.WatchOpts, sink chan<- *EASTimestamped, data [][32]byte, timestamp []uint64) (event.Subscription, error) {

	var dataRule []interface{}
	for _, dataItem := range data {
		dataRule = append(dataRule, dataItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _EAS.contract.WatchLogs(opts, "Timestamped", dataRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EASTimestamped)
				if err := _EAS.contract.UnpackLog(event, "Timestamped", log); err != nil {
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

// ParseTimestamped is a log parse operation binding the contract event 0x5aafceeb1c7ad58e4a84898bdee37c02c0fc46e7d24e6b60e8209449f183459f.
//
// Solidity: event Timestamped(bytes32 indexed data, uint64 indexed timestamp)
func (_EAS *EASFilterer) ParseTimestamped(log types.Log) (*EASTimestamped, error) {
	event := new(EASTimestamped)
	if err := _EAS.contract.UnpackLog(event, "Timestamped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
