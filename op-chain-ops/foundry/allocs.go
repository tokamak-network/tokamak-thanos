package foundry

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

// ArtifactsFS provides access to foundry build artifacts.
type ArtifactsFS struct {
	FS fs.FS
}

// ForgeAllocs represents the account allocations from a forge state dump.
type ForgeAllocs struct {
	Accounts types.GenesisAlloc
}

// Copy creates a deep copy of the ForgeAllocs.
func (fa *ForgeAllocs) Copy() *ForgeAllocs {
	out := &ForgeAllocs{Accounts: make(types.GenesisAlloc, len(fa.Accounts))}
	for addr, acc := range fa.Accounts {
		out.Accounts[addr] = acc
	}
	return out
}

// UnmarshalJSON implements json.Unmarshaler.
func (fa *ForgeAllocs) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &fa.Accounts)
}

// MarshalJSON implements json.Marshaler.
func (fa *ForgeAllocs) MarshalJSON() ([]byte, error) {
	return json.Marshal(fa.Accounts)
}

// LoadForgeAllocs loads a ForgeAllocs from a JSON file.
func LoadForgeAllocs(path string) (*ForgeAllocs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading allocs file: %w", err)
	}
	var allocs ForgeAllocs
	if err := json.Unmarshal(data, &allocs); err != nil {
		return nil, fmt.Errorf("parsing allocs: %w", err)
	}
	return &allocs, nil
}

// FromGenesisAlloc creates ForgeAllocs from genesis alloc data.
func FromGenesisAlloc(alloc types.GenesisAlloc) *ForgeAllocs {
	return &ForgeAllocs{Accounts: alloc}
}

// Addresses returns all account addresses in the allocs.
func (fa *ForgeAllocs) Addresses() []common.Address {
	addrs := make([]common.Address, 0, len(fa.Accounts))
	for addr := range fa.Accounts {
		addrs = append(addrs, addr)
	}
	return addrs
}

// ReadArtifact reads a foundry artifact from the FS.
func (a *ArtifactsFS) ReadArtifact(fileName, contractName string) (*Artifact, error) {
	if a.FS == nil {
		return nil, fmt.Errorf("no artifacts FS configured")
	}
	path := fileName + "/" + contractName + ".json"
	data, err := fs.ReadFile(a.FS, path)
	if err != nil {
		return nil, fmt.Errorf("reading artifact %s: %w", path, err)
	}
	var artifact Artifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		return nil, fmt.Errorf("parsing artifact %s: %w", path, err)
	}
	return &artifact, nil
}

// FromState populates ForgeAllocs from a state database by iterating all accounts.
func (a *ForgeAllocs) FromState(stateDB *state.StateDB) {
	a.Accounts = make(core.GenesisAlloc)
	// Iterate all accounts via state dump
	dumpConf := &state.DumpConfig{
		SkipCode:          false,
		SkipStorage:       false,
		OnlyWithAddresses: false,
	}
	dump := stateDB.RawDump(dumpConf)
	for addrStr, account := range dump.Accounts {
		addr := common.HexToAddress(addrStr)
		balance, _ := new(big.Int).SetString(account.Balance, 10)
		ga := core.GenesisAccount{
			Nonce:   account.Nonce,
			Balance: balance,
			Code:    []byte(account.Code),
			Storage: make(map[common.Hash]common.Hash),
		}
		for k, v := range account.Storage {
			ga.Storage[k] = common.HexToHash(v)
		}
		a.Accounts[addr] = ga
	}
}
