package infrastructure

import (
	"errors"
	"sync"

	"github.com/solrac97gr/goarchtest/test/clean_architecture/domain"
)

// InMemoryUserRepository implements the domain.UserRepository interface with an in-memory data store
type InMemoryUserRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// GetByID retrieves a user by their ID
func (r *InMemoryUserRepository) GetByID(id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *InMemoryUserRepository) GetByEmail(email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// Save persists a user
func (r *InMemoryUserRepository) Save(user *domain.User) error {
	if user == nil {
		return errors.New("cannot save nil user")
	}

	if err := user.Validate(); err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if a user with the same email already exists
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email && existingUser.ID != user.ID {
			return errors.New("a user with this email already exists")
		}
	}

	r.users[user.ID] = user
	return nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(user *domain.User) error {
	if user == nil {
		return errors.New("cannot update nil user")
	}

	if err := user.Validate(); err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	// Check if update would create a duplicate email
	for id, existingUser := range r.users {
		if existingUser.Email == user.Email && id != user.ID {
			return errors.New("a user with this email already exists")
		}
	}

	r.users[user.ID] = user
	return nil
}

// Delete removes a user
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}
