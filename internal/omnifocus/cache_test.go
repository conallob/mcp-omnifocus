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

func TestCacheGetNonExistent(t *testing.T) {
	cache := NewCache(1 * time.Second)

	_, found := cache.Get("non-existent-key")
	if found {
		t.Error("Expected not to find non-existent key")
	}
}

func TestCacheSetOverwrite(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Set("test-key", "initial-value")
	cache.Set("test-key", "updated-value")

	value, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached value")
	}
	if value.(string) != "updated-value" {
		t.Errorf("Expected 'updated-value', got '%s'", value)
	}
}

func TestCacheMultipleTypes(t *testing.T) {
	cache := NewCache(1 * time.Second)

	// Test different value types
	cache.Set("string", "test")
	cache.Set("int", 42)
	cache.Set("bool", true)
	cache.Set("slice", []string{"a", "b", "c"})
	cache.Set("map", map[string]int{"x": 1, "y": 2})

	// Verify all types can be retrieved
	str, found := cache.Get("string")
	if !found || str.(string) != "test" {
		t.Error("Failed to retrieve string value")
	}

	num, found := cache.Get("int")
	if !found || num.(int) != 42 {
		t.Error("Failed to retrieve int value")
	}

	b, found := cache.Get("bool")
	if !found || b.(bool) != true {
		t.Error("Failed to retrieve bool value")
	}

	slice, found := cache.Get("slice")
	if !found || len(slice.([]string)) != 3 {
		t.Error("Failed to retrieve slice value")
	}

	m, found := cache.Get("map")
	if !found || m.(map[string]int)["x"] != 1 {
		t.Error("Failed to retrieve map value")
	}
}

func TestCacheInvalidateNonExistent(t *testing.T) {
	cache := NewCache(1 * time.Second)

	// Invalidating non-existent key should not panic
	cache.Invalidate("non-existent")
}

func TestCacheInvalidatePatternNoMatch(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Set("tasks:all", "all tasks")
	cache.Set("projects:all", "all projects")

	// Invalidate pattern that doesn't match anything
	cache.InvalidatePattern("users:")

	// All original entries should remain
	_, found1 := cache.Get("tasks:all")
	_, found2 := cache.Get("projects:all")

	if !found1 || !found2 {
		t.Error("Expected all entries to remain when pattern doesn't match")
	}
}

func TestCacheInvalidatePatternExactMatch(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Set("tasks", "tasks")
	cache.Set("tasks:all", "all tasks")

	// Invalidate with exact prefix
	cache.InvalidatePattern("tasks")

	// Both should be invalidated
	_, found1 := cache.Get("tasks")
	_, found2 := cache.Get("tasks:all")

	if found1 || found2 {
		t.Error("Expected both entries to be invalidated")
	}
}

func TestCacheCleanupDisabled(t *testing.T) {
	cache := NewCache(0)

	// Should not panic when cache is disabled
	cache.Cleanup()
}

func TestCacheInvalidateAllDisabled(t *testing.T) {
	cache := NewCache(0)

	// Should not panic when cache is disabled
	cache.InvalidateAll()
}

func TestCacheInvalidateDisabled(t *testing.T) {
	cache := NewCache(0)

	// Should not panic when cache is disabled
	cache.Invalidate("test-key")
}

func TestCacheInvalidatePatternDisabled(t *testing.T) {
	cache := NewCache(0)

	// Should not panic when cache is disabled
	cache.InvalidatePattern("test:")
}

func TestCacheNegativeTTL(t *testing.T) {
	cache := NewCache(-1 * time.Second)

	cache.Set("test-key", "test-value")

	_, found := cache.Get("test-key")
	if found {
		t.Error("Expected cache to be disabled with negative TTL")
	}
}

func TestCacheCleanupEmptyCache(t *testing.T) {
	cache := NewCache(1 * time.Second)

	// Should not panic on empty cache
	cache.Cleanup()
}

func TestCacheInvalidateAllEmptyCache(t *testing.T) {
	cache := NewCache(1 * time.Second)

	// Should not panic on empty cache
	cache.InvalidateAll()
}

func TestCacheStartCleanupTimer(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Start cleanup timer with short interval
	cache.StartCleanupTimer(60 * time.Millisecond)

	// Wait for cleanup to potentially run
	time.Sleep(150 * time.Millisecond)

	// Add a new entry
	cache.Set("key3", "value3")

	// Old entries should be cleaned up automatically
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	_, found3 := cache.Get("key3")

	if found1 || found2 {
		t.Error("Expected expired entries to be cleaned up by timer")
	}
	if !found3 {
		t.Error("Expected new entry to remain")
	}
}

func TestCacheStartCleanupTimerDisabled(t *testing.T) {
	cache := NewCache(0)

	// Should not panic when cache is disabled
	cache.StartCleanupTimer(100 * time.Millisecond)
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewCache(1 * time.Second)
	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(n int) {
			for j := 0; j < 100; j++ {
				cache.Set("key", n*100+j)
			}
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				cache.Get("key")
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Should not panic and should have a value
	_, found := cache.Get("key")
	if !found {
		t.Error("Expected to find value after concurrent access")
	}
}

func TestCacheConcurrentInvalidation(t *testing.T) {
	cache := NewCache(1 * time.Second)
	done := make(chan bool)

	// Set initial values
	for i := 0; i < 100; i++ {
		cache.Set("key", i)
	}

	// Concurrent invalidations
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				cache.Invalidate("key")
			}
			done <- true
		}()
	}

	// Concurrent sets
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				cache.Set("key", j)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic
}
