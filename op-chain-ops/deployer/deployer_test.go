package deployer

import (
	"context"
	"math/big"
	"testing"
)

func TestNewL2Backend(t *testing.T) {
	backend, err := NewL2Backend()
	if err != nil {
		t.Fatalf("NewL2Backend returned error: %v", err)
	}
	defer backend.Close()

	chainID, err := backend.ChainID(context.Background())
	if err != nil {
		t.Fatalf("ChainID returned error: %v", err)
	}
	if chainID.Cmp(ChainID) != 0 {
		t.Fatalf("unexpected chain id: got %s, want %s", chainID, ChainID)
	}
}

func TestNewL2BackendWithChainIDAndPredeploys(t *testing.T) {
	customChainID := big.NewInt(10_001)
	backend, err := NewL2BackendWithChainIDAndPredeploys(customChainID, nil)
	if err != nil {
		t.Fatalf("NewL2BackendWithChainIDAndPredeploys returned error: %v", err)
	}
	defer backend.Close()

	chainID, err := backend.ChainID(context.Background())
	if err != nil {
		t.Fatalf("ChainID returned error: %v", err)
	}
	if chainID.Cmp(customChainID) != 0 {
		t.Fatalf("unexpected chain id: got %s, want %s", chainID, customChainID)
	}
}
