package deployer

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestBytecodeEqual(t *testing.T) {
	cases := []struct {
		name string
		a, b []byte
		want bool
	}{
		{"identical empty", []byte{}, []byte{}, true},
		{"identical non-empty", []byte{0x60, 0x80, 0x60}, []byte{0x60, 0x80, 0x60}, true},
		{"length differs", []byte{0x60}, []byte{0x60, 0x60}, false},
		{"same length, last byte diff", []byte{0x60, 0x80}, []byte{0x60, 0x81}, false},
		{"same length, first byte diff", []byte{0x61, 0x80}, []byte{0x60, 0x80}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := bytecodeEqual(tc.a, tc.b)
			if got != tc.want {
				t.Fatalf("bytecodeEqual(%x, %x) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestLoadRegistry_DisabledReturnsEmpty(t *testing.T) {
	cfg := DeployConfig{ReuseDeployment: false}
	r, err := loadRegistry(cfg, 11155111, fstest.MapFS{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil || len(r.Implementations) != 0 {
		t.Fatalf("want empty registry, got %+v", r)
	}
}

func TestLoadRegistry_NoEmbedNoOverride(t *testing.T) {
	cfg := DeployConfig{ReuseDeployment: true}
	r, err := loadRegistry(cfg, 99999, fstest.MapFS{})
	if err != nil {
		t.Fatalf("expected silent fallback to empty, got: %v", err)
	}
	if r == nil || len(r.Implementations) != 0 {
		t.Fatalf("want empty registry for unknown chain, got %+v", r)
	}
}

func TestLoadRegistry_EmbedFound(t *testing.T) {
	body := `{
		"tokamakDeployerVersion": "v0.0.6",
		"l1ChainId": 11155111,
		"implementations": { "SystemConfig": "0x1111111111111111111111111111111111111111" }
	}`
	mfs := fstest.MapFS{"registry/11155111.json": &fstest.MapFile{Data: []byte(body)}}
	cfg := DeployConfig{ReuseDeployment: true}
	r, err := loadRegistry(cfg, 11155111, mfs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Implementations["SystemConfig"] != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("want SystemConfig address loaded, got %+v", r.Implementations)
	}
}

func TestLoadRegistry_OverridePrecedence(t *testing.T) {
	embedBody := `{"tokamakDeployerVersion":"v0.0.6","l1ChainId":11155111,"implementations":{"X":"0xaa"}}`
	overrideBody := `{"tokamakDeployerVersion":"v0.0.6","l1ChainId":11155111,"implementations":{"Y":"0xbb"}}`
	mfs := fstest.MapFS{"registry/11155111.json": &fstest.MapFile{Data: []byte(embedBody)}}

	dir := t.TempDir()
	overridePath := filepath.Join(dir, "custom.json")
	if err := os.WriteFile(overridePath, []byte(overrideBody), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := DeployConfig{ReuseDeployment: true, RegistryPath: overridePath}
	r, err := loadRegistry(cfg, 11155111, mfs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Implementations["X"]; ok {
		t.Fatalf("override should hide embed: %+v", r.Implementations)
	}
	if r.Implementations["Y"] != "0xbb" {
		t.Fatalf("want override entry, got %+v", r.Implementations)
	}
}

func TestLoadRegistry_MalformedJSON(t *testing.T) {
	mfs := fstest.MapFS{"registry/11155111.json": &fstest.MapFile{Data: []byte("not json")}}
	cfg := DeployConfig{ReuseDeployment: true}
	_, err := loadRegistry(cfg, 11155111, mfs)
	if err == nil {
		t.Fatal("want parse error, got nil")
	}
}

func TestLoadRegistry_ChainIDMismatch(t *testing.T) {
	body := `{"tokamakDeployerVersion":"v0.0.6","l1ChainId":1,"implementations":{}}`
	mfs := fstest.MapFS{"registry/11155111.json": &fstest.MapFile{Data: []byte(body)}}
	cfg := DeployConfig{ReuseDeployment: true}
	_, err := loadRegistry(cfg, 11155111, mfs)
	if err == nil {
		t.Fatal("want chainId mismatch error, got nil")
	}
}

func TestLoadRegistry_OverrideMissingFile(t *testing.T) {
	cfg := DeployConfig{ReuseDeployment: true, RegistryPath: "/nonexistent/path/x.json"}
	_, err := loadRegistry(cfg, 11155111, fstest.MapFS{})
	if err == nil {
		t.Fatal("want read error, got nil")
	}
}

// mockRPC returns an httptest server that responds to eth_getCode with codeMap[address]
// and to all other methods with a generic "method not implemented" error.
func mockRPC(t *testing.T, codeMap map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     int           `json:"id"`
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case "eth_getCode":
			addr := req.Params[0].(string)
			code, ok := codeMap[addr]
			if !ok {
				code = "0x"
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"jsonrpc": "2.0", "id": req.ID, "result": code,
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"jsonrpc": "2.0", "id": req.ID,
				"error": map[string]interface{}{"code": -32601, "message": "method not supported"},
			})
		}
	}))
}

// mustLoadDeployedBytecode pulls the embedded artifact's deployedBytecode.object.
func mustLoadDeployedBytecode(t *testing.T, artifactsFS fs.FS, name string) []byte {
	t.Helper()
	a, err := loadArtifact(artifactsFS, name)
	if err != nil {
		t.Fatalf("loadArtifact(%s): %v", name, err)
	}
	return common.FromHex(a.DeployedBytecode.Object)
}

func TestRegistryVerify_HitMatchingBytecode(t *testing.T) {
	// os.DirFS("../../cmd") points at cmd/tokamak-deployer/cmd/, which contains deploy-artifacts/.
	artifactsFS := os.DirFS("../../cmd")
	hexCode := hexutil.Encode(mustLoadDeployedBytecode(t, artifactsFS, "SystemConfig"))

	addr := "0x1111111111111111111111111111111111111111"
	srv := mockRPC(t, map[string]string{addr: hexCode})
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	r := &Registry{
		L1ChainID:       11155111,
		Implementations: map[string]string{"SystemConfig": addr},
	}
	table, err := r.verify(context.Background(), client, artifactsFS, false)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	got := table.lookup("SystemConfig")
	if got.Hex() != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("want hit, got %v", got)
	}
}

func TestRegistryVerify_NoCode_LenientSkips(t *testing.T) {
	artifactsFS := os.DirFS("../../cmd")
	srv := mockRPC(t, map[string]string{}) // no code
	defer srv.Close()
	client, _ := ethclient.Dial(srv.URL)

	addr := "0x2222222222222222222222222222222222222222"
	r := &Registry{
		L1ChainID:       11155111,
		Implementations: map[string]string{"SystemConfig": addr},
	}
	table, err := r.verify(context.Background(), client, artifactsFS, false)
	if err != nil {
		t.Fatalf("lenient mode should not error: %v", err)
	}
	if table.size() != 0 {
		t.Fatalf("want empty table, got size %d", table.size())
	}
}

func TestRegistryVerify_NoCode_StrictAborts(t *testing.T) {
	artifactsFS := os.DirFS("../../cmd")
	srv := mockRPC(t, map[string]string{})
	defer srv.Close()
	client, _ := ethclient.Dial(srv.URL)

	r := &Registry{
		L1ChainID:       11155111,
		Implementations: map[string]string{"SystemConfig": "0x2222222222222222222222222222222222222222"},
	}
	_, err := r.verify(context.Background(), client, artifactsFS, true)
	if err == nil {
		t.Fatal("strict mode should abort, got nil")
	}
}

func TestRegistryVerify_HashMismatch_LenientSkips(t *testing.T) {
	artifactsFS := os.DirFS("../../cmd")
	addr := "0x3333333333333333333333333333333333333333"
	wrongCode := "0xdeadbeef" // not the SystemConfig deployedBytecode
	srv := mockRPC(t, map[string]string{addr: wrongCode})
	defer srv.Close()
	client, _ := ethclient.Dial(srv.URL)

	r := &Registry{
		L1ChainID:       11155111,
		Implementations: map[string]string{"SystemConfig": addr},
	}
	table, err := r.verify(context.Background(), client, artifactsFS, false)
	if err != nil {
		t.Fatalf("lenient: %v", err)
	}
	if table.size() != 0 {
		t.Fatalf("want empty table due to hash mismatch, got size %d", table.size())
	}
}

func TestRegistryVerify_InvalidHex_LenientSkips(t *testing.T) {
	artifactsFS := os.DirFS("../../cmd")
	srv := mockRPC(t, map[string]string{})
	defer srv.Close()
	client, _ := ethclient.Dial(srv.URL)
	r := &Registry{
		L1ChainID:       11155111,
		Implementations: map[string]string{"SystemConfig": "not-a-hex-address"},
	}
	table, err := r.verify(context.Background(), client, artifactsFS, false)
	if err != nil {
		t.Fatalf("lenient: %v", err)
	}
	if table.size() != 0 {
		t.Fatalf("want empty table, got size %d", table.size())
	}
}

func TestReuseTable_LookupMiss(t *testing.T) {
	var t0 *reuseTable // nil
	if got := t0.lookup("X"); (got != common.Address{}) {
		t.Fatalf("nil table lookup should be zero, got %v", got)
	}
	t1 := &reuseTable{addrs: map[string]common.Address{}}
	if got := t1.lookup("missing"); (got != common.Address{}) {
		t.Fatalf("empty table lookup should be zero, got %v", got)
	}
}

func TestDeployOrReuse_HitReturnsCachedAddress(t *testing.T) {
	addr := common.HexToAddress("0xaaaa000000000000000000000000000000000001")
	table := &reuseTable{addrs: map[string]common.Address{"SystemConfig": addr}}
	nonceBefore := uint64(7)
	nonce := nonceBefore

	// On hit: deployContract is never called, so passing nil for client/auth/gasPrice/artifact is safe.
	got, err := deployOrReuse(context.Background(), nil, nil, &nonce, nil, "SystemConfig", nil, table)
	if err != nil {
		t.Fatalf("hit path should not error: %v", err)
	}
	if got != addr {
		t.Fatalf("want cached %v, got %v", addr, got)
	}
	if nonce != nonceBefore {
		t.Fatalf("nonce should not advance on hit, before=%d after=%d", nonceBefore, nonce)
	}
}
