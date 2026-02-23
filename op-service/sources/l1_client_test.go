package sources

import "testing"

func TestL1ClientSimpleConfig_SetsL1BlockRefsCacheSize(t *testing.T) {
	cacheSize := 300
	cfg := L1ClientSimpleConfig(true, RPCKindStandard, cacheSize)

	if cfg.L1BlockRefsCacheSize != cacheSize {
		t.Fatalf("expected L1BlockRefsCacheSize %d, got %d", cacheSize, cfg.L1BlockRefsCacheSize)
	}
	if cfg.ReceiptsCacheSize != cacheSize {
		t.Fatalf("expected ReceiptsCacheSize %d, got %d", cacheSize, cfg.ReceiptsCacheSize)
	}
	if cfg.TransactionsCacheSize != cacheSize {
		t.Fatalf("expected TransactionsCacheSize %d, got %d", cacheSize, cfg.TransactionsCacheSize)
	}
}
