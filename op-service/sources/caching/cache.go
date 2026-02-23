package caching

import lru "github.com/hashicorp/golang-lru/v2"

type Metrics interface {
	CacheAdd(label string, cacheSize int, evicted bool)
	CacheGet(label string, hit bool)
}

// LRUCache wraps hashicorp *lru.Cache and tracks cache metrics
type LRUCache[K comparable, V any] struct {
	m     Metrics
	label string
	inner *lru.Cache[K, V]
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	if c.inner != nil {
		value, ok = c.inner.Get(key)
	}
	if c.m != nil {
		c.m.CacheGet(c.label, ok)
	}
	return value, ok
}

func (c *LRUCache[K, V]) Add(key K, value V) (evicted bool) {
	cacheSize := 0
	if c.inner != nil {
		evicted = c.inner.Add(key, value)
		cacheSize = c.inner.Len()
	}
	if c.m != nil {
		c.m.CacheAdd(c.label, cacheSize, evicted)
	}
	return evicted
}

// NewLRUCache creates a LRU cache with the given metrics, labeling the cache adds/gets.
// Metrics are optional: no metrics will be tracked if m == nil.
func NewLRUCache[K comparable, V any](m Metrics, label string, maxSize int) *LRUCache[K, V] {
	cache, err := lru.New[K, V](maxSize)
	if err != nil {
		// Invalid size disables caching while preserving call-sites and metrics behavior.
		cache = nil
	}
	return &LRUCache[K, V]{
		m:     m,
		label: label,
		inner: cache,
	}
}
