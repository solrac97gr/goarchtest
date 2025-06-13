package infrastructure

import (
	"errors"
	"fmt"
	"sync"

	"github.com/solrac97gr/goarchtest/examples/sample_project/domain"
)

// InMemoryUserRepository implements the UserRepository interface using an in-memory store
type InMemoryUserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// GetByID retrieves a user by ID from the in-memory store
func (r *InMemoryUserRepository) GetByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}

	// Return a copy to prevent mutation of the stored object
	return &domain.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// Save stores a user in the in-memory store
func (r *InMemoryUserRepository) Save(user *domain.User) error {
	if user.ID == "" {
		return errors.New("user ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent mutation of the original object
	r.users[user.ID] = &domain.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return nil
}

// Delete removes a user from the in-memory store
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user with ID %s not found", id)
	}

	delete(r.users, id)
	return nil
}
