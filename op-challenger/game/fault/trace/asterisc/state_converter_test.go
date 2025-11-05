package asterisc

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/vm"
)

const testBinary = "./somewhere/asterisc"

func TestStateConverter(t *testing.T) {
	setup := func(t *testing.T) (*StateConverter, *capturingExecutor) {
		vmCfg := vm.Config{
			VmBin: testBinary,
		}
		executor := &capturingExecutor{}
		converter := NewStateConverter(vmCfg)
		converter.cmdExecutor = executor.exec
		return converter, executor
	}

	t.Run("Valid", func(t *testing.T) {
		converter, executor := setup(t)
		// Manually construct JSON with witness field (as hex string)
		// because VMState.Witness has json:"-" tag
		jsonOutput := `{
			"pc": 11,
			"exited": true,
			"step": 42,
			"witness": "0x01020304",
			"stateHash": "0xab00000000000000000000000000000000000000000000000000000000000000"
		}`
		executor.stdOut = jsonOutput
		proof, step, exited, err := converter.ConvertStateToProof(context.Background(), "foo.json")
		require.NoError(t, err)
		require.Equal(t, true, exited)
		require.Equal(t, uint64(42), step)
		require.Equal(t, common.Hash{0xab}, proof.ClaimValue)
		require.Equal(t, []byte{1, 2, 3, 4}, []byte(proof.StateData))
		require.NotNil(t, proof.ProofData, "later validations require this to be non-nil")

		require.Equal(t, testBinary, executor.binary)
		require.Equal(t, []string{"witness", "--input", "foo.json"}, executor.args)
	})

	t.Run("CommandError", func(t *testing.T) {
		converter, executor := setup(t)
		executor.err = errors.New("boom")
		_, _, _, err := converter.ConvertStateToProof(context.Background(), "foo.json")
		require.ErrorIs(t, err, executor.err)
	})

	t.Run("InvalidOutput", func(t *testing.T) {
		converter, executor := setup(t)
		executor.stdOut = "blah blah"
		_, _, _, err := converter.ConvertStateToProof(context.Background(), "foo.json")
		require.ErrorContains(t, err, "failed to parse state data")
	})
}

type capturingExecutor struct {
	binary string
	args   []string

	stdOut string
	stdErr string
	err    error
}

func (c *capturingExecutor) exec(_ context.Context, binary string, args ...string) (string, string, error) {
	c.binary = binary
	c.args = args
	return c.stdOut, c.stdErr, c.err
}
