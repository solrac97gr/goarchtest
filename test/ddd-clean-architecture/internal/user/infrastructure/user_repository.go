package infrastructure

import (
	"errors"
	"sync"

	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/user/domain/models"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/user/domain/ports"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// InMemoryUserRepository provides an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	users map[shared.ID]*models.User
	mutex sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() ports.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[shared.ID]*models.User),
	}
}

// Save saves a user to the repository
func (r *InMemoryUserRepository) Save(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.users[user.ID] = user
	return nil
}

// FindByID finds a user by ID
func (r *InMemoryUserRepository) FindByID(id shared.ID) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// FindByEmail finds a user by email
func (r *InMemoryUserRepository) FindByEmail(email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(id shared.ID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.users, id)
	return nil
}

// List returns all users
func (r *InMemoryUserRepository) List() ([]*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}