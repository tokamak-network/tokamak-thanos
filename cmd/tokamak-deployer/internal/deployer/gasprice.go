package deployer

import (
	"context"
	"fmt"
	"log"
	"math/big"
)

// Default gas price bounds applied when DeployConfig does not specify them.
// Floor protects against RPC nodes returning 0 / 1 wei (historical bug).
// Ceil protects against runaway costs on mainnet during congestion.
var (
	defaultGasPriceFloor = big.NewInt(1_000_000_000)       // 1 Gwei
	defaultGasPriceCeil  = big.NewInt(100_000_000_000)     // 100 Gwei
	defaultGasMultiplier = 200                             // 200% = 2x
)

// gasPriceSuggester is the minimal surface of ethclient.Client we need. It
// exists so resolveGasPrice can be unit-tested without a live RPC.
type gasPriceSuggester interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

// resolveGasPrice picks the gas price to reuse for every transaction in a
// single deploy run. This replaces the previous per-TX SuggestGasPrice calls
// that made the deployer re-quote on every one of the 26-32 transactions.
//
// Behaviour:
//   - If cfg.FixedGasPrice is set, use it directly (then clamp to [Floor, Ceil]).
//   - Otherwise, fetch SuggestGasPrice once, multiply by Multiplier/100,
//     and clamp.
//
// The returned price is intentionally conservative so that the bump-on-timeout
// safety net in sendAndWaitMined rarely activates.
func resolveGasPrice(ctx context.Context, sg gasPriceSuggester, cfg DeployConfig) (*big.Int, error) {
	multiplier := cfg.GasPriceMultiplier
	if multiplier <= 0 {
		multiplier = defaultGasMultiplier
	}
	floor := cfg.GasPriceFloor
	if floor == nil || floor.Sign() <= 0 {
		floor = defaultGasPriceFloor
	}
	ceil := cfg.GasPriceCeil
	if ceil == nil || ceil.Sign() <= 0 {
		ceil = defaultGasPriceCeil
	}
	if floor.Cmp(ceil) > 0 {
		return nil, fmt.Errorf("gas price floor (%s) > ceil (%s)", floor, ceil)
	}

	var (
		price     *big.Int
		suggested *big.Int
	)
	if cfg.FixedGasPrice != nil && cfg.FixedGasPrice.Sign() > 0 {
		price = new(big.Int).Set(cfg.FixedGasPrice)
	} else {
		raw, err := sg.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
		if raw == nil || raw.Sign() <= 0 {
			return nil, fmt.Errorf("RPC returned non-positive gas price: %v", raw)
		}
		suggested = new(big.Int).Set(raw)
		price = new(big.Int).Mul(raw, big.NewInt(int64(multiplier)))
		price.Div(price, big.NewInt(100))
	}

	if price.Cmp(floor) < 0 {
		price = new(big.Int).Set(floor)
	}
	if price.Cmp(ceil) > 0 {
		price = new(big.Int).Set(ceil)
	}

	if suggested != nil {
		log.Printf("[deployer] Fixed gas price: %s wei (%s Gwei) — suggested %s wei × %d%%",
			price, gweiStr(price), suggested, multiplier)
	} else {
		log.Printf("[deployer] Fixed gas price: %s wei (%s Gwei) — user-specified",
			price, gweiStr(price))
	}
	return price, nil
}

func gweiStr(wei *big.Int) string {
	gwei := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e9))
	return gwei.Text('f', 3)
}
