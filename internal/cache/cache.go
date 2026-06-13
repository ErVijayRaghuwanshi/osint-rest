package cache

import (
	"sync"
	"time"
)

// Cache is the interface for all cache implementations.
type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration)
	Delete(key string)
	Flush()
	Close() error
}

// entry holds a cached value and its expiration time.
type entry struct {
	data      []byte
	expiresAt time.Time
}

// MemoryCache is a simple in-memory cache with per-key TTL.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]entry
	stop  chan struct{}
}

// NewMemoryCache creates a new in-memory cache and starts a background
// goroutine that evicts expired entries every cleanupInterval.
func NewMemoryCache(cleanupInterval time.Duration) *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]entry),
		stop:  make(chan struct{}),
	}

	go c.janitor(cleanupInterval)
	return c
}

// Get returns the cached value for key if it exists and has not expired.
func (c *MemoryCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(e.expiresAt) {
		return nil, false
	}

	return e.data, true
}

// Set stores a value with the given TTL. A TTL of 0 means no expiration.
func (c *MemoryCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	exp := time.Now().Add(ttl)
	if ttl == 0 {
		exp = time.Now().Add(24 * time.Hour) // default max
	}

	c.items[key] = entry{
		data:      value,
		expiresAt: exp,
	}
}

// Delete removes a key from the cache.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all entries from the cache.
func (c *MemoryCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]entry)
}

// Close stops the background janitor goroutine.
func (c *MemoryCache) Close() error {
	close(c.stop)
	return nil
}

// janitor periodically removes expired entries.
func (c *MemoryCache) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.evictExpired()
		case <-c.stop:
			return
		}
	}
}

func (c *MemoryCache) evictExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, e := range c.items {
		if now.After(e.expiresAt) {
			delete(c.items, k)
		}
	}
}
