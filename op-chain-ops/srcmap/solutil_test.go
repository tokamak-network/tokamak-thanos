package srcmap

import (
	"testing"
)

func TestSourcemap(t *testing.T) {
	t.Skip("TODO(clabby): This test is disabled until source IDs have been added to foundry artifacts.")

	// contractsDir := "../../packages/tokamak/contracts-bedrock"
	// sources := []string{path.Join(contractsDir, "src/cannon/MIPS.sol")}
	// for i, source := range sources {
	// 	sources[i] = path.Join(contractsDir, source)
	// }
	//
	// deployedByteCode := hexutil.MustDecode(bindings.MIPSDeployedBin)
	// srcMap, err := ParseSourceMap(
	// 	sources,
	// 	deployedByteCode,
	// 	bindings.MIPSDeployedSourceMap)
	// require.NoError(t, err)
	//
	// for i := 0; i < len(deployedByteCode); i++ {
	// 	info := srcMap.FormattedInfo(uint64(i))
	// 	if strings.HasPrefix(info, "unknown") {
	// 		t.Fatalf("unexpected info: %q", info)
	// 	}
	// }
}
