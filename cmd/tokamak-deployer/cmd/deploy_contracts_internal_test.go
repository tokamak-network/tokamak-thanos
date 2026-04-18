package cmd

import (
	"testing"
)

// TestDeployContractsCmd_FaultProofFlag verifies that `--fault-proof` is
// registered as a bool flag on the deploy-contracts command and defaults to
// false. Regression test for Bug #8 (fault-proof contracts never deployed
// because the CLI flag was missing, leaving cfg.EnableFaultProof = false).
func TestDeployContractsCmd_FaultProofFlag(t *testing.T) {
	flag := deployContractsCmd.Flags().Lookup("fault-proof")
	if flag == nil {
		t.Fatal("expected --fault-proof flag to be registered on deploy-contracts command")
	}
	if flag.Value.Type() != "bool" {
		t.Errorf("expected --fault-proof to be bool, got %q", flag.Value.Type())
	}
	if flag.DefValue != "false" {
		t.Errorf("expected --fault-proof default=false, got %q", flag.DefValue)
	}
}
