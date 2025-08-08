package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Entry represents a cache entry with optional TTL
type Entry[V any] struct {
	Key       string
	Value     V
	ExpiresAt *time.Time // nil means no expiration
}

// IsExpired checks if the entry has expired
func (e *Entry[V]) IsExpired() bool {
	return e.ExpiresAt != nil && time.Now().After(*e.ExpiresAt)
}

// LRUCache is a thread-safe LRU cache with optional TTL support
type LRUCache[K comparable, V any] struct {
	mu       sync.RWMutex
	capacity int
	items    map[K]*list.Element
	order    *list.List // doubly linked list for LRU ordering

	// Statistics
	hits   int64
	misses int64
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		panic("cache capacity must be positive")
	}

	return &LRUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves a value from the cache
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V

	element, exists := c.items[key]
	if !exists {
		c.misses++
		return zero, false
	}

	entry := element.Value.(*Entry[V])

	// Check if expired
	if entry.IsExpired() {
		c.removeElement(element)
		c.misses++
		return zero, false
	}

	// Move to front (most recently used)
	c.order.MoveToFront(element)
	c.hits++

	return entry.Value, true
}

// Put adds or updates a value in the cache
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.PutWithTTL(key, value, 0)
}

// PutWithTTL adds or updates a value in the cache with TTL
func (c *LRUCache[K, V]) PutWithTTL(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt *time.Time
	if ttl > 0 {
		expiry := time.Now().Add(ttl)
		expiresAt = &expiry
	}

	entry := &Entry[V]{
		Key:       fmt.Sprintf("%v", key), // Convert key to string for entry
		Value:     value,
		ExpiresAt: expiresAt,
	}

	// If key already exists, update it
	if element, exists := c.items[key]; exists {
		element.Value = entry
		c.order.MoveToFront(element)
		return
	}

	// Add new entry
	element := c.order.PushFront(entry)
	c.items[key] = element

	// Remove oldest entries if over capacity
	for c.order.Len() > c.capacity {
		c.removeOldest()
	}
}

// Delete removes a key from the cache
func (c *LRUCache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, exists := c.items[key]
	if !exists {
		return false
	}

	c.removeElement(element)
	return true
}

// Contains checks if a key exists in the cache (without updating LRU order)
func (c *LRUCache[K, V]) Contains(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	element, exists := c.items[key]
	if !exists {
		return false
	}

	entry := element.Value.(*Entry[V])
	if entry.IsExpired() {
		// Don't remove here to avoid upgrading to write lock
		return false
	}

	return true
}

// Size returns the current number of items in the cache
func (c *LRUCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Capacity returns the maximum capacity of the cache
func (c *LRUCache[K, V]) Capacity() int {
	return c.capacity
}

// Clear removes all items from the cache
func (c *LRUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*list.Element)
	c.order.Init()
	c.hits = 0
	c.misses = 0
}

// Keys returns all keys in the cache (from most to least recently used)
func (c *LRUCache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.items))
	for element := c.order.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*Entry[V])
		if !entry.IsExpired() {
			// Need to find the key from the map since we can't store K in Entry
			for k, elem := range c.items {
				if elem == element {
					keys = append(keys, k)
					break
				}
			}
		}
	}

	return keys
}

// Values returns all values in the cache (from most to least recently used)
func (c *LRUCache[K, V]) Values() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	values := make([]V, 0, len(c.items))
	for element := c.order.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*Entry[V])
		if !entry.IsExpired() {
			values = append(values, entry.Value)
		}
	}

	return values
}

// Stats returns cache statistics
func (c *LRUCache[K, V]) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return CacheStats{
		Hits:     c.hits,
		Misses:   c.misses,
		HitRate:  hitRate,
		Size:     len(c.items),
		Capacity: c.capacity,
	}
}

// CleanupExpired removes all expired entries
func (c *LRUCache[K, V]) CleanupExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	removed := 0
	for element := c.order.Back(); element != nil; {
		entry := element.Value.(*Entry[V])
		prev := element.Prev()

		if entry.IsExpired() {
			c.removeElement(element)
			removed++
		}

		element = prev
	}

	return removed
}

// removeElement removes an element from both the map and list
func (c *LRUCache[K, V]) removeElement(element *list.Element) {
	c.order.Remove(element)

	// Find and remove from map
	for key, elem := range c.items {
		if elem == element {
			delete(c.items, key)
			break
		}
	}
}

// removeOldest removes the least recently used item
func (c *LRUCache[K, V]) removeOldest() {
	if c.order.Len() == 0 {
		return
	}

	oldest := c.order.Back()
	c.removeElement(oldest)
}

// CacheStats holds cache performance statistics
type CacheStats struct {
	Hits     int64   `json:"hits"`
	Misses   int64   `json:"misses"`
	HitRate  float64 `json:"hit_rate"`
	Size     int     `json:"size"`
	Capacity int     `json:"capacity"`
}
