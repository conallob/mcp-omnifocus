package omnifocus

import (
	"testing"
	"time"
)

func TestCacheGetSet(t *testing.T) {
	cache := NewCache(1 * time.Second)

	// Test setting and getting a value
	cache.Set("test-key", "test-value")

	value, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached value")
	}
	if value.(string) != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", value)
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)

	cache.Set("test-key", "test-value")

	// Value should be present immediately
	_, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached value immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Value should be expired
	_, found = cache.Get("test-key")
	if found {
		t.Error("Expected cached value to be expired")
	}
}

func TestCacheInvalidate(t *testing.T) {
	cache := NewCache(10 * time.Second)

	cache.Set("test-key", "test-value")
	cache.Invalidate("test-key")

	_, found := cache.Get("test-key")
	if found {
		t.Error("Expected cached value to be invalidated")
	}
}

func TestCacheInvalidateAll(t *testing.T) {
	cache := NewCache(10 * time.Second)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	cache.InvalidateAll()

	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	_, found3 := cache.Get("key3")

	if found1 || found2 || found3 {
		t.Error("Expected all cached values to be invalidated")
	}
}

func TestCacheInvalidatePattern(t *testing.T) {
	cache := NewCache(10 * time.Second)

	cache.Set("tasks:all", "all tasks")
	cache.Set("tasks:project:123", "project 123 tasks")
	cache.Set("tasks:project:456", "project 456 tasks")
	cache.Set("projects:all", "all projects")

	// Invalidate all task-related caches
	cache.InvalidatePattern("tasks:")

	_, found1 := cache.Get("tasks:all")
	_, found2 := cache.Get("tasks:project:123")
	_, found3 := cache.Get("tasks:project:456")
	_, found4 := cache.Get("projects:all")

	if found1 || found2 || found3 {
		t.Error("Expected all task caches to be invalidated")
	}
	if !found4 {
		t.Error("Expected projects cache to remain")
	}
}

func TestCacheDisabled(t *testing.T) {
	cache := NewCache(0)

	cache.Set("test-key", "test-value")

	_, found := cache.Get("test-key")
	if found {
		t.Error("Expected cache to be disabled with TTL of 0")
	}
}

func TestCacheCleanup(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Add a new entry that shouldn't expire
	cache.Set("key3", "value3")

	// Run cleanup
	cache.Cleanup()

	// Expired entries should be removed, new entry should remain
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	_, found3 := cache.Get("key3")

	if found1 || found2 {
		t.Error("Expected expired entries to be cleaned up")
	}
	if !found3 {
		t.Error("Expected non-expired entry to remain after cleanup")
	}
}
