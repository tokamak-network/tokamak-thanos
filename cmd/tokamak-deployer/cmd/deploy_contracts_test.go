package cmd_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

func TestDeployContracts_Reuse_Anvil(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	rpcURL, stop := startAnvil(t)
	defer stop()

	// First deploy — no reuse, default behavior.
	out1 := t.TempDir() + "/deploy-output-1.json"
	first := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", rpcURL,
		"--private-key", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"--chain-id", "901",
		"--out", out1,
	)
	first.Dir = ".."
	first.Env = append(os.Environ(), "GOWORK=off")
	if firstOut, err := first.CombinedOutput(); err != nil {
		t.Fatalf("first deploy failed: %v\n%s", err, firstOut)
	}

	// Read first output — extract Implementations map.
	data1, err := os.ReadFile(out1)
	if err != nil {
		t.Fatalf("read out1: %v", err)
	}
	var deployOutput1 struct {
		Implementations map[string]string `json:"implementations"`
	}
	if err := json.Unmarshal(data1, &deployOutput1); err != nil {
		t.Fatalf("parse out1: %v", err)
	}
	if len(deployOutput1.Implementations) < 8 {
		// Without --fault-proof, expect exactly 8 entries (no DisputeGameFactory).
		t.Fatalf("expected ≥8 impl entries from first deploy, got %d: %+v",
			len(deployOutput1.Implementations), deployOutput1.Implementations)
	}

	// Snapshot deployer nonce after first deploy.
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatal(err)
	}
	deployerAddr := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cfFFb92266") // anvil[0]
	nonceAfterFirst, err := client.NonceAt(context.Background(), deployerAddr, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Build a temporary registry from first deploy's Implementations.
	chainID, _ := client.ChainID(context.Background())
	registry := map[string]interface{}{
		"tokamakDeployerVersion": "test",
		"l1ChainId":              chainID.Uint64(),
		"comment":                "integration test",
		"implementations":        deployOutput1.Implementations,
	}
	regBytes, _ := json.MarshalIndent(registry, "", "  ")
	regPath := filepath.Join(t.TempDir(), "reuse-registry.json")
	if err := os.WriteFile(regPath, regBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	// Second deploy — reuse on, override registry.
	out2 := t.TempDir() + "/deploy-output-2.json"
	second := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", rpcURL,
		"--private-key", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"--chain-id", "902",
		"--out", out2,
		"--reuse-deployment",
		"--reuse-impls", regPath,
	)
	second.Dir = ".."
	second.Env = append(os.Environ(), "GOWORK=off")
	secondOut, err := second.CombinedOutput()
	if err != nil {
		t.Fatalf("second deploy failed: %v\n%s", err, secondOut)
	}

	// Verify the log shows reuse (♻ marker).
	if !strings.Contains(string(secondOut), "♻") {
		t.Errorf("expected ♻ reuse marker in log, got:\n%s", secondOut)
	}

	// Nonce delta on second deploy should be 26 - 8 = 18 (no fault-proof).
	nonceAfterSecond, err := client.NonceAt(context.Background(), deployerAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	delta := nonceAfterSecond - nonceAfterFirst
	expected := uint64(26 - 8)
	if delta != expected {
		t.Errorf("expected nonce delta %d (26 base steps - 8 reused impls), got %d",
			expected, delta)
	}

	// Verify second deploy's output records the *reused* impl as L2OutputOracle.
	data2, _ := os.ReadFile(out2)
	var deployOutput2 struct {
		L2OutputOracleProxy string            `json:"L2OutputOracleProxy"`
		Implementations     map[string]string `json:"implementations"`
	}
	if err := json.Unmarshal(data2, &deployOutput2); err != nil {
		t.Fatalf("parse out2: %v", err)
	}
	if !strings.EqualFold(deployOutput2.Implementations["L2OutputOracle"], deployOutput1.Implementations["L2OutputOracle"]) {
		t.Errorf("L2OutputOracle impl differs: first=%s second=%s",
			deployOutput1.Implementations["L2OutputOracle"],
			deployOutput2.Implementations["L2OutputOracle"])
	}
}

func TestDeployContracts_ReuseStrict_Aborts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	rpcURL, stop := startAnvil(t)
	defer stop()

	client, _ := ethclient.Dial(rpcURL)
	chainID, _ := client.ChainID(context.Background())

	// Registry with a bogus address — no code at it.
	bogus := map[string]interface{}{
		"tokamakDeployerVersion": "test",
		"l1ChainId":              chainID.Uint64(),
		"implementations": map[string]string{
			"SystemConfig": "0xdEAD000000000000000000000000000000000000",
		},
	}
	regBytes, _ := json.MarshalIndent(bogus, "", "  ")
	regPath := filepath.Join(t.TempDir(), "bogus.json")
	if err := os.WriteFile(regPath, regBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".", "deploy-contracts",
		"--l1-rpc", rpcURL,
		"--private-key", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"--chain-id", "903",
		"--out", t.TempDir()+"/x.json",
		"--reuse-deployment",
		"--reuse-impls", regPath,
		"--reuse-strict",
	)
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected abort under --reuse-strict with bogus address, got success:\n%s", out)
	}
	if !strings.Contains(string(out), "no code on-chain") && !strings.Contains(string(out), "preflight reuse registry") {
		t.Errorf("expected preflight failure mention, got:\n%s", out)
	}
}
