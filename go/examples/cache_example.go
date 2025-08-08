package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Hanasou/news_feed/go/common/cache"
)

// Example models for caching
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	LastSeen time.Time `json:"last_seen"`
}

type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	fmt.Println("=== LRU Cache Examples ===")

	// Example 1: User Cache
	fmt.Println("1. User Cache Example")
	userCache := cache.NewLRUCache[string, *User](100) // Cache up to 100 users

	// Add users to cache
	users := []*User{
		{ID: "user1", Username: "john_doe", Email: "john@example.com", LastSeen: time.Now()},
		{ID: "user2", Username: "jane_smith", Email: "jane@example.com", LastSeen: time.Now().Add(-time.Hour)},
		{ID: "user3", Username: "bob_wilson", Email: "bob@example.com", LastSeen: time.Now().Add(-2 * time.Hour)},
	}

	for _, user := range users {
		userCache.Put(user.ID, user)
		fmt.Printf("Cached user: %s (%s)\n", user.Username, user.ID)
	}

	// Retrieve users from cache
	if user, found := userCache.Get("user1"); found {
		fmt.Printf("Retrieved user from cache: %s (%s)\n", user.Username, user.Email)
	}

	fmt.Printf("Cache size: %d\n", userCache.Size())
	fmt.Printf("Cache stats: %+v\n\n", userCache.Stats())

	// Example 2: Todo Cache with TTL
	fmt.Println("2. Todo Cache with TTL")
	todoCache := cache.NewLRUCache[string, *Todo](50) // Cache up to 50 todos

	// Add todos with TTL
	todos := []*Todo{
		{ID: "todo1", Title: "Buy groceries", Description: "Milk, bread, eggs", UserID: "user1", CreatedAt: time.Now()},
		{ID: "todo2", Title: "Walk the dog", Description: "30 minute walk in the park", UserID: "user1", CreatedAt: time.Now()},
		{ID: "todo3", Title: "Read book", Description: "Continue reading Go programming book", UserID: "user2", CreatedAt: time.Now()},
	}

	// Cache some todos with short TTL (for demonstration)
	todoCache.PutWithTTL("todo1", todos[0], 2*time.Second)
	todoCache.Put("todo2", todos[1]) // No TTL
	todoCache.PutWithTTL("todo3", todos[2], 1*time.Second)

	// Immediate retrieval
	for _, todo := range todos {
		if cachedTodo, found := todoCache.Get(todo.ID); found {
			fmt.Printf("Retrieved todo: %s - %s\n", cachedTodo.ID, cachedTodo.Title)
		}
	}

	fmt.Println("\nWaiting for TTL expiration...")
	time.Sleep(1500 * time.Millisecond)

	// Check what's still available after TTL
	fmt.Println("After TTL expiration:")
	for _, todo := range todos {
		if cachedTodo, found := todoCache.Get(todo.ID); found {
			fmt.Printf("Still cached: %s - %s\n", cachedTodo.ID, cachedTodo.Title)
		} else {
			fmt.Printf("Expired: %s\n", todo.ID)
		}
	}

	// Example 3: LRU Eviction
	fmt.Println("\n3. LRU Eviction Example")
	smallCache := cache.NewLRUCache[string, string](3) // Very small cache

	// Fill cache to capacity
	smallCache.Put("key1", "value1")
	smallCache.Put("key2", "value2")
	smallCache.Put("key3", "value3")
	fmt.Printf("Cache after filling: size=%d\n", smallCache.Size())

	// Access key1 to make it most recently used
	smallCache.Get("key1")
	fmt.Println("Accessed key1 (now most recent)")

	// Add key4 - should evict key2 (least recently used)
	smallCache.Put("key4", "value4")
	fmt.Printf("Added key4, cache size: %d\n", smallCache.Size())

	// Check what's still in cache
	keys := smallCache.Keys()
	fmt.Printf("Remaining keys (in LRU order): %v\n", keys)

	// Example 4: Cache Statistics and Monitoring
	fmt.Println("\n4. Cache Statistics Example")
	monitorCache := cache.NewLRUCache[int, string](10)

	// Generate some cache activity
	for i := 0; i < 20; i++ {
		monitorCache.Put(i, fmt.Sprintf("value-%d", i))
	}

	// Generate hits and misses
	for i := 0; i < 30; i++ {
		monitorCache.Get(i % 15) // Some hits, some misses
	}

	stats := monitorCache.Stats()
	fmt.Printf("Cache Statistics:\n")
	fmt.Printf("  Hits: %d\n", stats.Hits)
	fmt.Printf("  Misses: %d\n", stats.Misses)
	fmt.Printf("  Hit Rate: %.2f%%\n", stats.HitRate*100)
	fmt.Printf("  Size: %d/%d\n", stats.Size, stats.Capacity)

	// Example 5: Cleanup Expired Entries
	fmt.Println("\n5. Cleanup Example")
	cleanupCache := cache.NewLRUCache[string, string](10)

	// Add entries with various TTLs
	cleanupCache.PutWithTTL("expire-fast", "value1", 100*time.Millisecond)
	cleanupCache.PutWithTTL("expire-slow", "value2", 500*time.Millisecond)
	cleanupCache.Put("no-expire", "value3")

	fmt.Printf("Cache size before cleanup: %d\n", cleanupCache.Size())

	// Wait and cleanup
	time.Sleep(200 * time.Millisecond)
	removed := cleanupCache.CleanupExpired()
	fmt.Printf("Removed %d expired entries\n", removed)
	fmt.Printf("Cache size after cleanup: %d\n", cleanupCache.Size())

	// Example 6: Service Integration Pattern
	fmt.Println("\n6. Service Integration Pattern")
	userService := NewUserService()

	// Simulate getting users (with caching)
	for _, userID := range []string{"user1", "user2", "user1", "user3", "user1"} {
		user, err := userService.GetUser(userID)
		if err != nil {
			log.Printf("Error getting user %s: %v", userID, err)
			continue
		}
		fmt.Printf("Got user: %s (from %s)\n", user.Username,
			func() string {
				if userService.cache.Contains(userID) {
					return "cache"
				}
				return "database"
			}())
	}

	// Show final cache stats
	finalStats := userService.cache.Stats()
	fmt.Printf("\nFinal User Service Cache Stats:\n")
	fmt.Printf("  Hit Rate: %.1f%%\n", finalStats.HitRate*100)
	fmt.Printf("  Cache Size: %d\n", finalStats.Size)

	fmt.Println("\n=== Examples Complete ===")
}

// UserService demonstrates how to integrate LRU cache with a service layer
type UserService struct {
	cache *cache.LRUCache[string, *User]
}

func NewUserService() *UserService {
	return &UserService{
		cache: cache.NewLRUCache[string, *User](100),
	}
}

func (s *UserService) GetUser(userID string) (*User, error) {
	// Try cache first
	if user, found := s.cache.Get(userID); found {
		return user, nil
	}

	// Simulate database lookup
	user := s.fetchUserFromDatabase(userID)
	if user == nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	// Cache the result with 10 minute TTL
	s.cache.PutWithTTL(userID, user, 10*time.Minute)

	return user, nil
}

func (s *UserService) fetchUserFromDatabase(userID string) *User {
	// Simulate database lookup delay
	time.Sleep(10 * time.Millisecond)

	// Mock database data
	users := map[string]*User{
		"user1": {ID: "user1", Username: "john_doe", Email: "john@example.com", LastSeen: time.Now()},
		"user2": {ID: "user2", Username: "jane_smith", Email: "jane@example.com", LastSeen: time.Now()},
		"user3": {ID: "user3", Username: "bob_wilson", Email: "bob@example.com", LastSeen: time.Now()},
	}

	return users[userID]
}
