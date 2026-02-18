package jsonutil

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/tokamak-network/tokamak-thanos/op-service/ioutil"
)

// MergeJSON merges override values into base via JSON round-trip.
func MergeJSON[X any](base X, overrides map[string]interface{}) (X, error) {
	baseBytes, err := json.Marshal(base)
	if err != nil {
		return base, fmt.Errorf("failed to marshal base: %w", err)
	}
	var baseMap map[string]interface{}
	if err := json.Unmarshal(baseBytes, &baseMap); err != nil {
		return base, fmt.Errorf("failed to unmarshal base: %w", err)
	}
	for k, v := range overrides {
		baseMap[k] = v
	}
	merged, err := json.Marshal(baseMap)
	if err != nil {
		return base, fmt.Errorf("failed to marshal merged: %w", err)
	}
	var result X
	if err := json.Unmarshal(merged, &result); err != nil {
		return base, fmt.Errorf("failed to unmarshal merged: %w", err)
	}
	return result, nil
}

// WriteTOML writes value as TOML to the given target.
func WriteTOML[X any](value X, target ioutil.OutputTarget) error {
	w, closer, aborter, err := target()
	if err != nil {
		return err
	}
	defer closer.Close()
	if err := toml.NewEncoder(w).Encode(value); err != nil {
		aborter()
		return err
	}
	return nil
}
