package deployer

import (
	"context"
	"errors"
	"math/big"
	"testing"
)

type stubSuggester struct {
	price *big.Int
	err   error
	calls int
}

func (s *stubSuggester) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	s.calls++
	return s.price, s.err
}

func TestResolveGasPrice_AutoWithDefaultMultiplier(t *testing.T) {
	sg := &stubSuggester{price: big.NewInt(5_000_000_000)} // 5 Gwei
	price, err := resolveGasPrice(context.Background(), sg, DeployConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := big.NewInt(10_000_000_000) // 5 Gwei * 2 = 10 Gwei
	if price.Cmp(want) != 0 {
		t.Errorf("price = %s, want %s", price, want)
	}
	if sg.calls != 1 {
		t.Errorf("SuggestGasPrice called %d times, want 1", sg.calls)
	}
}

func TestResolveGasPrice_AutoWithCustomMultiplier(t *testing.T) {
	sg := &stubSuggester{price: big.NewInt(4_000_000_000)} // 4 Gwei
	cfg := DeployConfig{GasPriceMultiplier: 150}           // 1.5x
	price, err := resolveGasPrice(context.Background(), sg, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := big.NewInt(6_000_000_000) // 4 * 1.5 = 6 Gwei
	if price.Cmp(want) != 0 {
		t.Errorf("price = %s, want %s", price, want)
	}
}

func TestResolveGasPrice_FixedPriceBypassesSuggest(t *testing.T) {
	sg := &stubSuggester{price: big.NewInt(999), err: errors.New("should not be called")}
	cfg := DeployConfig{FixedGasPrice: big.NewInt(7_000_000_000)} // 7 Gwei
	price, err := resolveGasPrice(context.Background(), sg, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price.Cmp(big.NewInt(7_000_000_000)) != 0 {
		t.Errorf("price = %s, want 7 Gwei", price)
	}
	if sg.calls != 0 {
		t.Errorf("SuggestGasPrice should not be called when FixedGasPrice is set")
	}
}

func TestResolveGasPrice_FloorApplied(t *testing.T) {
	// Suggested price 100 wei × 2x = 200 wei, well below 1 Gwei default floor.
	sg := &stubSuggester{price: big.NewInt(100)}
	price, err := resolveGasPrice(context.Background(), sg, DeployConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price.Cmp(defaultGasPriceFloor) != 0 {
		t.Errorf("price = %s, want default floor %s", price, defaultGasPriceFloor)
	}
}

func TestResolveGasPrice_CeilApplied(t *testing.T) {
	// Suggested 100 Gwei × 2x = 200 Gwei, exceeds 100 Gwei ceil.
	sg := &stubSuggester{price: big.NewInt(100_000_000_000)}
	price, err := resolveGasPrice(context.Background(), sg, DeployConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price.Cmp(defaultGasPriceCeil) != 0 {
		t.Errorf("price = %s, want default ceil %s", price, defaultGasPriceCeil)
	}
}

func TestResolveGasPrice_FixedPriceClampedByCeil(t *testing.T) {
	// Fixed 200 Gwei, ceil default 100 Gwei → clamped.
	cfg := DeployConfig{FixedGasPrice: big.NewInt(200_000_000_000)}
	price, err := resolveGasPrice(context.Background(), &stubSuggester{}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price.Cmp(defaultGasPriceCeil) != 0 {
		t.Errorf("price = %s, want ceil clamp %s", price, defaultGasPriceCeil)
	}
}

func TestResolveGasPrice_RPCError(t *testing.T) {
	sg := &stubSuggester{err: errors.New("rpc down")}
	_, err := resolveGasPrice(context.Background(), sg, DeployConfig{})
	if err == nil {
		t.Fatalf("expected error from RPC failure")
	}
}

func TestResolveGasPrice_NonPositiveSuggestion(t *testing.T) {
	sg := &stubSuggester{price: big.NewInt(0)}
	_, err := resolveGasPrice(context.Background(), sg, DeployConfig{})
	if err == nil {
		t.Fatalf("expected error for zero gas price from RPC")
	}
}

func TestResolveGasPrice_CustomFloorAndCeil(t *testing.T) {
	sg := &stubSuggester{price: big.NewInt(3_000_000_000)} // 3 Gwei
	cfg := DeployConfig{
		GasPriceMultiplier: 200,
		GasPriceFloor:      big.NewInt(2_000_000_000),   // 2 Gwei
		GasPriceCeil:       big.NewInt(5_000_000_000),   // 5 Gwei
	}
	price, err := resolveGasPrice(context.Background(), sg, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 3 × 2 = 6 Gwei, ceil 5 Gwei → 5 Gwei.
	if price.Cmp(big.NewInt(5_000_000_000)) != 0 {
		t.Errorf("price = %s, want custom ceil 5 Gwei", price)
	}
}

func TestResolveGasPrice_InvalidFloorGreaterThanCeil(t *testing.T) {
	cfg := DeployConfig{
		GasPriceFloor: big.NewInt(10_000_000_000),
		GasPriceCeil:  big.NewInt(5_000_000_000),
	}
	_, err := resolveGasPrice(context.Background(), &stubSuggester{}, cfg)
	if err == nil {
		t.Fatalf("expected error when floor > ceil")
	}
}
