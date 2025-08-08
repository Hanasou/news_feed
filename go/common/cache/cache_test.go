package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache_BasicOperations(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	// Test initial state
	assert.Equal(t, 0, cache.Size())
	assert.Equal(t, 3, cache.Capacity())

	// Test Put and Get
	cache.Put("key1", 1)
	cache.Put("key2", 2)
	cache.Put("key3", 3)

	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, 1, value)

	value, found = cache.Get("key2")
	assert.True(t, found)
	assert.Equal(t, 2, value)

	value, found = cache.Get("key3")
	assert.True(t, found)
	assert.Equal(t, 3, value)

	// Test size
	assert.Equal(t, 3, cache.Size())
}

func TestLRUCache_Eviction(t *testing.T) {
	cache := NewLRUCache[string, int](2)

	// Fill cache to capacity
	cache.Put("key1", 1)
	cache.Put("key2", 2)
	assert.Equal(t, 2, cache.Size())

	// Add one more - should evict least recently used (key1)
	cache.Put("key3", 3)
	assert.Equal(t, 2, cache.Size())

	// key1 should be evicted
	_, found := cache.Get("key1")
	assert.False(t, found)

	// key2 and key3 should still be there
	value, found := cache.Get("key2")
	assert.True(t, found)
	assert.Equal(t, 2, value)

	value, found = cache.Get("key3")
	assert.True(t, found)
	assert.Equal(t, 3, value)
}

func TestLRUCache_LRUOrdering(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	// Add items
	cache.Put("key1", 1)
	cache.Put("key2", 2)
	cache.Put("key3", 3)

	// Access key1 to make it most recently used
	cache.Get("key1")

	// Add key4 - should evict key2 (least recently used)
	cache.Put("key4", 4)

	// key2 should be evicted
	_, found := cache.Get("key2")
	assert.False(t, found)

	// Others should still be there
	_, found = cache.Get("key1")
	assert.True(t, found)
	_, found = cache.Get("key3")
	assert.True(t, found)
	_, found = cache.Get("key4")
	assert.True(t, found)
}

func TestLRUCache_Update(t *testing.T) {
	cache := NewLRUCache[string, int](2)

	// Add item
	cache.Put("key1", 1)

	// Update item
	cache.Put("key1", 10)

	// Should still have size 1
	assert.Equal(t, 1, cache.Size())

	// Should get updated value
	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, 10, value)
}

func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	cache.Put("key1", 1)
	cache.Put("key2", 2)
	assert.Equal(t, 2, cache.Size())

	// Delete existing key
	deleted := cache.Delete("key1")
	assert.True(t, deleted)
	assert.Equal(t, 1, cache.Size())

	// Try to get deleted key
	_, found := cache.Get("key1")
	assert.False(t, found)

	// Delete non-existing key
	deleted = cache.Delete("nonexistent")
	assert.False(t, deleted)
}

func TestLRUCache_Contains(t *testing.T) {
	cache := NewLRUCache[string, int](2)

	cache.Put("key1", 1)

	assert.True(t, cache.Contains("key1"))
	assert.False(t, cache.Contains("nonexistent"))

	// Contains should not affect LRU order
	cache.Put("key2", 2)
	cache.Put("key3", 3) // Should evict key1

	assert.False(t, cache.Contains("key1"))
	assert.True(t, cache.Contains("key2"))
	assert.True(t, cache.Contains("key3"))
}

func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	cache.Put("key1", 1)
	cache.Put("key2", 2)
	cache.Put("key3", 3)
	assert.Equal(t, 3, cache.Size())

	cache.Clear()
	assert.Equal(t, 0, cache.Size())

	// All keys should be gone
	_, found := cache.Get("key1")
	assert.False(t, found)
	_, found = cache.Get("key2")
	assert.False(t, found)
	_, found = cache.Get("key3")
	assert.False(t, found)
}

func TestLRUCache_KeysAndValues(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	cache.Put("key1", 1)
	cache.Put("key2", 2)
	cache.Put("key3", 3)

	// Access key1 to make it most recent
	cache.Get("key1")

	keys := cache.Keys()
	values := cache.Values()

	// Should be in LRU order (most recent first)
	assert.Len(t, keys, 3)
	assert.Len(t, values, 3)

	// key1 should be first (most recent)
	assert.Equal(t, "key1", keys[0])
	assert.Equal(t, 1, values[0])
}

func TestLRUCache_TTL(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	// Add item with short TTL
	cache.PutWithTTL("key1", 1, 50*time.Millisecond)
	cache.Put("key2", 2) // No TTL

	// Should be available immediately
	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, 1, value)

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// key1 should be expired
	_, found = cache.Get("key1")
	assert.False(t, found)

	// key2 should still be available
	value, found = cache.Get("key2")
	assert.True(t, found)
	assert.Equal(t, 2, value)
}

func TestLRUCache_CleanupExpired(t *testing.T) {
	cache := NewLRUCache[string, int](5)

	// Add items with different TTLs
	cache.PutWithTTL("expire1", 1, 30*time.Millisecond)
	cache.PutWithTTL("expire2", 2, 30*time.Millisecond)
	cache.Put("keep1", 3)
	cache.Put("keep2", 4)

	assert.Equal(t, 4, cache.Size())

	// Wait for expiration
	time.Sleep(40 * time.Millisecond)

	// Cleanup expired items
	removed := cache.CleanupExpired()
	assert.Equal(t, 2, removed)
	assert.Equal(t, 2, cache.Size())

	// Only non-expired items should remain
	_, found := cache.Get("expire1")
	assert.False(t, found)
	_, found = cache.Get("expire2")
	assert.False(t, found)
	_, found = cache.Get("keep1")
	assert.True(t, found)
	_, found = cache.Get("keep2")
	assert.True(t, found)
}

func TestLRUCache_Stats(t *testing.T) {
	cache := NewLRUCache[string, int](3)

	// Initial stats
	stats := cache.Stats()
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
	assert.Equal(t, float64(0), stats.HitRate)

	// Add some data
	cache.Put("key1", 1)
	cache.Put("key2", 2)

	// Generate some hits and misses
	cache.Get("key1")    // hit
	cache.Get("key1")    // hit
	cache.Get("key2")    // hit
	cache.Get("missing") // miss
	cache.Get("missing") // miss

	stats = cache.Stats()
	assert.Equal(t, int64(3), stats.Hits)
	assert.Equal(t, int64(2), stats.Misses)
	assert.Equal(t, 0.6, stats.HitRate) // 3/5 = 0.6
	assert.Equal(t, 2, stats.Size)
	assert.Equal(t, 3, stats.Capacity)
}

func TestLRUCache_ThreadSafety(t *testing.T) {
	cache := NewLRUCache[int, string](100)

	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup

	// Start multiple goroutines performing concurrent operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j

				// Put
				cache.Put(key, fmt.Sprintf("value-%d", key))

				// Get
				cache.Get(key)

				// Contains
				cache.Contains(key)

				// Maybe delete
				if j%10 == 0 {
					cache.Delete(key)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Cache should still be in a valid state
	assert.True(t, cache.Size() <= cache.Capacity())

	stats := cache.Stats()
	assert.True(t, stats.Hits >= 0)
	assert.True(t, stats.Misses >= 0)
}

func TestLRUCache_DifferentTypes(t *testing.T) {
	// Test with different key and value types

	// String keys, struct values
	type User struct {
		ID   int
		Name string
	}

	userCache := NewLRUCache[string, User](2)
	userCache.Put("user1", User{ID: 1, Name: "John"})
	userCache.Put("user2", User{ID: 2, Name: "Jane"})

	user, found := userCache.Get("user1")
	assert.True(t, found)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.Name)

	// Integer keys, string values
	intCache := NewLRUCache[int, string](2)
	intCache.Put(1, "one")
	intCache.Put(2, "two")

	value, found := intCache.Get(1)
	assert.True(t, found)
	assert.Equal(t, "one", value)
}

func BenchmarkLRUCache_Put(b *testing.B) {
	cache := NewLRUCache[int, string](1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(i%1000, fmt.Sprintf("value-%d", i))
	}
}

func BenchmarkLRUCache_Get(b *testing.B) {
	cache := NewLRUCache[int, string](1000)

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		cache.Put(i, fmt.Sprintf("value-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % 1000)
	}
}

func BenchmarkLRUCache_Concurrent(b *testing.B) {
	cache := NewLRUCache[int, string](1000)

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		cache.Put(i, fmt.Sprintf("value-%d", i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				cache.Get(i % 1000)
			} else {
				cache.Put(i%1000, fmt.Sprintf("value-%d", i))
			}
			i++
		}
	})
}
