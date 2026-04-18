package cmd_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func startAnvil(t *testing.T) (rpcURL string, stop func()) {
	t.Helper()
	cmd := exec.Command("anvil", "--port", "18545", "--block-time", "1")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start anvil: %v", err)
	}
	time.Sleep(500 * time.Millisecond)
	return "http://127.0.0.1:18545", func() { _ = cmd.Process.Kill() }
}

func TestDeployContracts_NotImplemented(t *testing.T) {
	rpcURL, stop := startAnvil(t)
	defer stop()

	outFile := t.TempDir() + "/deploy-output.json"
	cmd := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", rpcURL,
		"--private-key", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"--chain-id", "901",
		"--out", outFile,
	)
	cmd.Dir = "../.."
	out, err := cmd.CombinedOutput()
	// Stub returns "not implemented" error — expect non-zero exit
	if err == nil {
		t.Fatalf("expected error from stub, got output: %s", out)
	}
	t.Logf("expected failure: %s", out)
}

func TestDeployContracts_Anvil(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	rpcURL, stop := startAnvil(t)
	defer stop()

	outFile := t.TempDir() + "/deploy-output.json"
	cmd := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", rpcURL,
		"--private-key", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"--chain-id", "901",
		"--out", outFile,
	)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("deploy-contracts failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	var output map[string]interface{}
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	for _, key := range []string{"ProxyAdmin", "SystemConfigProxy", "OptimismPortalProxy", "L1StandardBridgeProxy"} {
		addr, ok := output[key].(string)
		if !ok || addr == "" || addr == "0x0000000000000000000000000000000000000000" {
			t.Errorf("expected non-zero address for %s, got: %v", key, output[key])
		}
	}
}

// TestDeployContracts_FaultProof_Anvil verifies the full producer-side Bug #8
// fix: running deploy-contracts with --fault-proof must execute steps 27-32
// and write non-zero AnchorStateRegistryProxy and DisputeGameFactoryProxy
// addresses to deploy-output.json. Prior to the v0.0.6 release the CLI flag
// did not exist, so cfg.EnableFaultProof stayed false and these addresses
// were always absent.
func TestDeployContracts_FaultProof_Anvil(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	rpcURL, stop := startAnvil(t)
	defer stop()

	outFile := t.TempDir() + "/deploy-output.json"
	cmd := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", rpcURL,
		"--private-key", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"--chain-id", "901",
		"--out", outFile,
		"--fault-proof",
	)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("deploy-contracts --fault-proof failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	var output map[string]interface{}
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Core addresses still populated
	for _, key := range []string{"ProxyAdmin", "SystemConfigProxy"} {
		addr, ok := output[key].(string)
		if !ok || addr == "" || addr == "0x0000000000000000000000000000000000000000" {
			t.Errorf("expected non-zero %s, got: %v", key, output[key])
		}
	}

	// Fault-proof addresses must now also be present
	for _, key := range []string{"AnchorStateRegistryProxy", "DisputeGameFactoryProxy"} {
		addr, ok := output[key].(string)
		if !ok || addr == "" || addr == "0x0000000000000000000000000000000000000000" {
			t.Errorf("fault-proof address %s missing or zero — steps 27-32 did not run: got %v",
				key, output[key])
		}
	}
}

func TestDeployContracts_BadRPC(t *testing.T) {
	outFile := t.TempDir() + "/deploy-output.json"
	cmd := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", "http://127.0.0.1:19999",
		"--private-key", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"--chain-id", "901",
		"--out", outFile,
	)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for unreachable RPC, got output: %s", out)
	}
	output := string(out)
	if !strings.Contains(output, "connect") && !strings.Contains(output, "connection refused") && !strings.Contains(output, "dial") {
		t.Errorf("expected connection error message, got: %s", output)
	}
}

func TestDeployContracts_BadPrivateKey(t *testing.T) {
	rpcURL, stop := startAnvil(t)
	defer stop()

	outFile := t.TempDir() + "/deploy-output.json"
	cmd := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", rpcURL,
		"--private-key", "0xINVALID",
		"--chain-id", "901",
		"--out", outFile,
	)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for invalid private key, got output: %s", out)
	}
	output := string(out)
	if !strings.Contains(output, "private key") && !strings.Contains(output, "hex") && !strings.Contains(output, "invalid") {
		t.Errorf("expected private key error message, got: %s", output)
	}
}
