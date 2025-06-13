package infrastructure

import (
	"sync"

	"github.com/solrac97gr/goarchtest/examples/sample_project/domain"
)

// UserCache provides caching functionality for users
type UserCache struct {
	cache map[string]*domain.User
	mu    sync.RWMutex
}

// NewUserCache creates a new UserCache
func NewUserCache() *UserCache {
	return &UserCache{
		cache: make(map[string]*domain.User),
	}
}

// Get retrieves a user from cache
func (c *UserCache) Get(id string) (*domain.User, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	user, exists := c.cache[id]
	return user, exists
}

// Set stores a user in cache
func (c *UserCache) Set(id string, user *domain.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[id] = user
}
