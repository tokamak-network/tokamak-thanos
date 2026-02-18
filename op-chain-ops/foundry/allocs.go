package foundry

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/ethereum/go-ethereum/common"
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
