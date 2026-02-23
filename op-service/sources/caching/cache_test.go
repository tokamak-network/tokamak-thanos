package caching

import "testing"

type testMetrics struct {
	addCalls    int
	getCalls    int
	lastLabel   string
	lastSize    int
	lastEvicted bool
	lastHit     bool
}

func (m *testMetrics) CacheAdd(label string, cacheSize int, evicted bool) {
	m.addCalls++
	m.lastLabel = label
	m.lastSize = cacheSize
	m.lastEvicted = evicted
}

func (m *testMetrics) CacheGet(label string, hit bool) {
	m.getCalls++
	m.lastLabel = label
	m.lastHit = hit
}

func TestNewLRUCache_DisablesOnInvalidSize(t *testing.T) {
	metrics := &testMetrics{}
	cache := NewLRUCache[int, string](metrics, "test-cache", 0)

	if cache == nil {
		t.Fatal("expected cache instance")
	}
	if evicted := cache.Add(1, "one"); evicted {
		t.Fatal("did not expect eviction for disabled cache")
	}
	if metrics.addCalls != 1 {
		t.Fatalf("expected 1 add metric call, got %d", metrics.addCalls)
	}
	if metrics.lastSize != 0 {
		t.Fatalf("expected cache size metric 0, got %d", metrics.lastSize)
	}

	if _, ok := cache.Get(1); ok {
		t.Fatal("disabled cache should not return cached values")
	}
	if metrics.getCalls != 1 {
		t.Fatalf("expected 1 get metric call, got %d", metrics.getCalls)
	}
	if metrics.lastHit {
		t.Fatal("disabled cache should report miss")
	}
}

func TestNewLRUCache_StoresValuesWhenEnabled(t *testing.T) {
	cache := NewLRUCache[int, string](nil, "test-cache", 2)

	cache.Add(1, "one")
	cache.Add(2, "two")

	v, ok := cache.Get(1)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v != "one" {
		t.Fatalf("unexpected cached value: %s", v)
	}
}
