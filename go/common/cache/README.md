# Thread-Safe LRU Cache

A high-performance, thread-safe LRU (Least Recently Used) cache implementation in Go using generics.

## Features

- ✅ **Thread-Safe** - Concurrent read/write operations with RWMutex
- ✅ **Generic Types** - Type-safe operations with Go generics
- ✅ **TTL Support** - Optional time-to-live for cache entries
- ✅ **LRU Eviction** - Automatic eviction of least recently used items
- ✅ **Statistics** - Built-in hit/miss tracking and performance metrics
- ✅ **High Performance** - Optimized for speed with minimal allocations

## Quick Start

```go
import "github.com/Hanasou/news_feed/go/common/cache"

// Create a cache that holds up to 100 string-to-int mappings
cache := cache.NewLRUCache[string, int](100)

// Put values
cache.Put("key1", 42)
cache.PutWithTTL("key2", 24, 5*time.Minute) // Expires in 5 minutes

// Get values
if value, found := cache.Get("key1"); found {
    fmt.Printf("Value: %d\n", value) // Value: 42
}

// Check statistics
stats := cache.Stats()
fmt.Printf("Hit rate: %.1f%%\n", stats.HitRate*100)
```

## API Reference

### Creating a Cache

```go
// NewLRUCache creates a new cache with specified capacity
cache := cache.NewLRUCache[KeyType, ValueType](capacity)
```

### Basic Operations

```go
// Put adds or updates a value
cache.Put(key, value)

// PutWithTTL adds a value with expiration time
cache.PutWithTTL(key, value, 10*time.Minute)

// Get retrieves a value (updates LRU order)
value, found := cache.Get(key)

// Contains checks existence (doesn't update LRU order)
exists := cache.Contains(key)

// Delete removes a key
deleted := cache.Delete(key)
```

### Bulk Operations

```go
// Get all keys (in LRU order, most recent first)
keys := cache.Keys()

// Get all values (in LRU order)
values := cache.Values()

// Clear all entries
cache.Clear()
```

### Maintenance

```go
// Remove expired entries manually
removedCount := cache.CleanupExpired()

// Get cache statistics
stats := cache.Stats()
fmt.Printf("Size: %d/%d, Hit Rate: %.1f%%\n", 
    stats.Size, stats.Capacity, stats.HitRate*100)
```

## Performance Benchmarks

```
BenchmarkLRUCache_Put-32           183.4 ns/op    86 B/op    4 allocs/op
BenchmarkLRUCache_Get-32            21.85 ns/op     0 B/op    0 allocs/op  
BenchmarkLRUCache_Concurrent-32    151.4 ns/op    42 B/op    2 allocs/op
```

- **Put operations**: ~183ns with minimal allocations
- **Get operations**: ~22ns with zero allocations (cache hits)
- **Concurrent access**: ~151ns with good scalability

## Usage Patterns

### 1. Simple Key-Value Cache

```go
cache := cache.NewLRUCache[string, string](1000)
cache.Put("user:123", "john_doe")

if username, found := cache.Get("user:123"); found {
    fmt.Println("Username:", username)
}
```

### 2. Struct Caching

```go
type User struct {
    ID   string
    Name string
    Email string
}

userCache := cache.NewLRUCache[string, *User](100)
user := &User{ID: "123", Name: "John", Email: "john@example.com"}

userCache.Put(user.ID, user)
```

### 3. Service Layer Integration

```go
type UserService struct {
    cache *cache.LRUCache[string, *User]
    db    Database
}

func (s *UserService) GetUser(id string) (*User, error) {
    // Try cache first
    if user, found := s.cache.Get(id); found {
        return user, nil
    }
    
    // Fetch from database
    user, err := s.db.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // Cache for 10 minutes
    s.cache.PutWithTTL(id, user, 10*time.Minute)
    return user, nil
}
```

### 4. TTL Cache for Session Management

```go
sessionCache := cache.NewLRUCache[string, *Session](1000)

// Cache session for 30 minutes
sessionCache.PutWithTTL(sessionID, session, 30*time.Minute)

// Periodic cleanup of expired sessions
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        removed := sessionCache.CleanupExpired()
        log.Printf("Cleaned up %d expired sessions", removed)
    }
}()
```

## Thread Safety

The cache is fully thread-safe and supports concurrent operations:

```go
// Safe to use from multiple goroutines
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        cache.Put(fmt.Sprintf("key-%d", id), id)
        cache.Get(fmt.Sprintf("key-%d", id))
    }(i)
}
wg.Wait()
```

## Configuration Best Practices

### Capacity Planning

```go
// Rule of thumb: capacity = expected_entries * 1.3
// For 1000 expected users, use capacity of 1300
userCache := cache.NewLRUCache[string, *User](1300)
```

### TTL Guidelines

```go
// Short TTL for frequently changing data
cache.PutWithTTL("stock-price", price, 1*time.Minute)

// Medium TTL for user profiles
cache.PutWithTTL(userID, profile, 15*time.Minute)

// Long TTL for configuration data
cache.PutWithTTL("config", config, 1*time.Hour)
```

### Monitoring

```go
// Log cache statistics periodically
go func() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        stats := cache.Stats()
        if stats.HitRate < 0.8 { // Alert if hit rate below 80%
            log.Printf("Low cache hit rate: %.1f%%", stats.HitRate*100)
        }
    }
}()
```

## Error Handling

The cache operations are designed to be safe and don't return errors for normal operations:

```go
// Get returns (value, found) - no error
value, found := cache.Get("key")
if !found {
    // Handle cache miss
}

// Put operations don't fail
cache.Put("key", "value") // Always succeeds

// Delete returns boolean indicating if key existed
existed := cache.Delete("key")
```

## Memory Management

The cache automatically manages memory through:

1. **LRU Eviction** - Removes least recently used items when capacity is exceeded
2. **TTL Expiration** - Automatically expires entries after their TTL
3. **Manual Cleanup** - Use `CleanupExpired()` for proactive cleanup

## Comparison with Other Caches

| Feature | LRU Cache | sync.Map | map[K]V |
|---------|-----------|----------|---------|
| Thread Safety | ✅ | ✅ | ❌ |
| Size Limits | ✅ | ❌ | ❌ |
| TTL Support | ✅ | ❌ | ❌ |
| LRU Eviction | ✅ | ❌ | ❌ |
| Type Safety | ✅ | ❌ | ✅ |
| Statistics | ✅ | ❌ | ❌ |

## License

This implementation is part of the news_feed project and follows the same license terms.
