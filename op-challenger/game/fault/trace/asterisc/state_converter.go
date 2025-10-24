package asterisc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/utils"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/vm"
)

// VMState represents Asterisc VM state
// Supports both Hex (0x...) and Base64 encoding for witness field
type VMState struct {
	PC        uint64        `json:"pc"`
	Exited    bool          `json:"exited"`
	Step      uint64        `json:"step"`
	Witness   hexutil.Bytes `json:"-"` // Custom unmarshal
	StateHash common.Hash   `json:"stateHash"`
}

// UnmarshalJSON implements custom JSON unmarshaling for VMState
// Handles both Base64 (prestate.json files) and Hex (asterisc witness command output)
func (v *VMState) UnmarshalJSON(data []byte) error {
	type Alias VMState
	aux := &struct {
		Witness interface{} `json:"witness"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle witness field (can be Base64 string or Hex string)
	if aux.Witness != nil {
		switch w := aux.Witness.(type) {
		case string:
			// Check if it's Hex (starts with 0x) or Base64
			if strings.HasPrefix(w, "0x") {
				// Hex format (from asterisc witness command)
				decoded, err := hexutil.Decode(w)
				if err != nil {
					return fmt.Errorf("failed to decode hex witness: %w", err)
				}
				v.Witness = decoded
			} else {
				// Base64 format (from prestate.json files)
				decoded, err := base64.StdEncoding.DecodeString(w)
				if err != nil {
					return fmt.Errorf("failed to decode base64 witness: %w", err)
				}
				v.Witness = decoded
			}
		default:
			return fmt.Errorf("witness field has unexpected type: %T", w)
		}
	}

	return nil
}

type StateConverter struct {
	vmConfig    vm.Config
	cmdExecutor func(ctx context.Context, binary string, args ...string) (stdOut string, stdErr string, err error)
}

func NewStateConverter(vmConfig vm.Config) *StateConverter {
	return &StateConverter{
		vmConfig:    vmConfig,
		cmdExecutor: runCmd,
	}
}

func (c *StateConverter) ConvertStateToProof(ctx context.Context, statePath string) (*utils.ProofData, uint64, bool, error) {
	stdOut, stdErr, err := c.cmdExecutor(ctx, c.vmConfig.VmBin, "witness", "--input", statePath)
	if err != nil {
		return nil, 0, false, fmt.Errorf("state conversion failed: %w (%s)", err, stdErr)
	}
	var data VMState
	if err := json.Unmarshal([]byte(stdOut), &data); err != nil {
		return nil, 0, false, fmt.Errorf("failed to parse state data: %w", err)
	}
	// Extend the trace out to the full length using a no-op instruction that doesn't change any state
	// No execution is done, so no proof-data or oracle values are required.
	return &utils.ProofData{
		ClaimValue:   data.StateHash,
		StateData:    data.Witness,
		ProofData:    []byte{},
		OracleKey:    nil,
		OracleValue:  nil,
		OracleOffset: 0,
	}, data.Step, data.Exited, nil
}

func runCmd(ctx context.Context, binary string, args ...string) (stdOut string, stdErr string, err error) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	stdOut = outBuf.String()
	stdErr = errBuf.String()
	return
}

func parseState(path string) (*VMState, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open state file: %w", err)
	}
	defer file.Close()
	var state VMState
	if err := json.NewDecoder(file).Decode(&state); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}
	return &state, nil
}
