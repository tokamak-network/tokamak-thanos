package cannon

import (
	"encoding/json"
	"fmt"

	"github.com/tokamak-network/tokamak-thanos/cannon/mipsevm"
	"github.com/tokamak-network/tokamak-thanos/op-service/ioutil"
)

func parseState(path string) (*mipsevm.State, error) {
	file, err := ioutil.OpenDecompressed(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open state file (%v): %w", path, err)
	}
	defer file.Close()
	var state mipsevm.State
	err = json.NewDecoder(file).Decode(&state)
	if err != nil {
		return nil, fmt.Errorf("invalid mipsevm state (%v): %w", path, err)
	}
	return &state, nil
}
