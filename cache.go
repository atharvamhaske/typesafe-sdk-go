package typesafe

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type cacheEntry struct {
	data    []byte
	expires time.Time
}

// ttlCache is a small in-memory cache for SystemOne responses, keyed by the
// hash of the request body. Disabled unless a client is built with
// WithCache: a decision API silently returning a stale verdict is worse
// than a cache miss.
type ttlCache struct {
	mu      sync.Mutex
	ttl     time.Duration
	max     int
	entries map[string]cacheEntry
}

func newTTLCache(max int, ttl time.Duration) *ttlCache {
	return &ttlCache{ttl: ttl, max: max, entries: make(map[string]cacheEntry)}
}

func cacheKey(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func (c *ttlCache) get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok || time.Now().After(e.expires) {
		return nil, false
	}
	return e.data, true
}

func (c *ttlCache) set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= c.max {
		for k := range c.entries { // ponytail: single-entry eviction, swap for LRU if hit rate matters
			delete(c.entries, k)
			break
		}
	}
	c.entries[key] = cacheEntry{data: data, expires: time.Now().Add(c.ttl)}
}
