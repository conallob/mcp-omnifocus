package omnifocus

import (
	"sync"
	"time"
)

// cacheEntry holds cached data with expiration time
type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// Cache provides a simple in-memory cache with TTL
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*cacheEntry
	ttl     time.Duration
	enabled bool
}

// NewCache creates a new cache with the specified TTL
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*cacheEntry),
		ttl:     ttl,
		enabled: ttl > 0, // Disable cache if TTL is 0 or negative
	}
}

// Get retrieves a value from the cache if it exists and hasn't expired
func (c *Cache) Get(key string) (interface{}, bool) {
	if !c.enabled {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.data, true
}

// Set stores a value in the cache with the configured TTL
func (c *Cache) Set(key string, value interface{}) {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &cacheEntry{
		data:      value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Invalidate removes a specific key from the cache
func (c *Cache) Invalidate(key string) {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

// InvalidateAll clears all entries from the cache
func (c *Cache) InvalidateAll() {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*cacheEntry)
}

// InvalidatePattern removes all keys matching a pattern (simple prefix match)
func (c *Cache) InvalidatePattern(prefix string) {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.entries {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.entries, key)
		}
	}
}

// Cleanup removes expired entries from the cache
// This should be called periodically to prevent memory growth
func (c *Cache) Cleanup() {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.expiresAt) {
			delete(c.entries, key)
		}
	}
}

// StartCleanupTimer starts a background goroutine that periodically cleans up expired entries
func (c *Cache) StartCleanupTimer(interval time.Duration) {
	if !c.enabled {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.Cleanup()
		}
	}()
}
