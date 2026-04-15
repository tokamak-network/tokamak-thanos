package cmd_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGenerateGenesis_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Minimal genesis fixture with key predeploy addresses pre-populated
	// This avoids requiring op-node for the test (--base-genesis flag)
	baseGenesis := `{
		"config": {"chainId": 901, "homesteadBlock": 0},
		"alloc": {
			"0x4200000000000000000000000000000000000015": {
				"code": "0x6080",
				"balance": "0x0"
			},
			"0x4200000000000000000000000000000000000067": {
				"code": "0x6080",
				"balance": "0x0"
			},
			"0x4200000000000000000000000000000000000778": {
				"code": "",
				"balance": "0x0"
			},
			"0x4200000000000000000000000000000000000018": {
				"code": "0x6080",
				"balance": "0x0"
			}
		},
		"number": "0x0",
		"gasLimit": "0x1c9c380",
		"timestamp": "0x0"
	}`
	baseGenesisFile := t.TempDir() + "/base-genesis.json"
	if err := os.WriteFile(baseGenesisFile, []byte(baseGenesis), 0644); err != nil {
		t.Fatalf("failed to write base genesis: %v", err)
	}

	deployOutput := `{
		"l1ChainId": 31337,
		"l2ChainId": 901,
		"ProxyAdmin": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
		"SystemConfigProxy": "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
		"OptimismPortalProxy": "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0",
		"L1StandardBridgeProxy": "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	}`
	deployOutputFile := t.TempDir() + "/deploy-output.json"
	if err := os.WriteFile(deployOutputFile, []byte(deployOutput), 0644); err != nil {
		t.Fatalf("failed to write deploy output: %v", err)
	}

	rollupConfig := `{
		"l2ChainID": 901,
		"l1ChainID": 31337
	}`
	rollupConfigFile := t.TempDir() + "/rollup-config.json"
	if err := os.WriteFile(rollupConfigFile, []byte(rollupConfig), 0644); err != nil {
		t.Fatalf("failed to write rollup config: %v", err)
	}

	genesisOut := t.TempDir() + "/genesis.json"

	cmd := exec.Command("go", "run", ".", "generate-genesis",
		"--deploy-output", deployOutputFile,
		"--config", rollupConfigFile,
		"--base-genesis", baseGenesisFile,
		"--out", genesisOut,
	)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generate-genesis failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(genesisOut)
	if err != nil {
		t.Fatalf("genesis.json not created: %v", err)
	}
	var genesis map[string]interface{}
	if err := json.Unmarshal(data, &genesis); err != nil {
		t.Fatalf("invalid genesis.json: %v", err)
	}

	// L1Block code namespace should have Isthmus-capable bytecode injected
	// Code namespace for 0x4200...0015 is 0xc0d3...0015
	alloc, ok := genesis["alloc"].(map[string]interface{})
	if !ok {
		t.Fatal("genesis missing alloc field")
	}

	// Check L1Block code namespace address has Isthmus-capable bytecode
	// predeployToCodeNamespace(0x4200...0015) = 0xc0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d30015
	l1BlockCodeNS := "0xc0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d30015"
	l1BlockAlloc, exists := alloc[l1BlockCodeNS]
	if !exists {
		// Try without 0x prefix as alloc key format may vary
		l1BlockAlloc, exists = alloc[strings.TrimPrefix(l1BlockCodeNS, "0x")]
		if !exists {
			t.Fatalf("L1Block code namespace not found in genesis alloc (checked %s)", l1BlockCodeNS)
		}
	}
	l1BlockMap := l1BlockAlloc.(map[string]interface{})
	code, _ := l1BlockMap["code"].(string)
	if len(code) < 10 {
		t.Errorf("L1Block code namespace too short, likely not injected: %s", code)
	}
}
