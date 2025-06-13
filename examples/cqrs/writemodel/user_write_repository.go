package writemodel

import "github.com/solrac97gr/goarchtest/examples/cqrs/domain"

// UserWriteRepository defines the interface for write operations
type UserWriteRepository interface {
	Save(user *domain.User) (*domain.User, error)
	Update(user *domain.User) error
	Delete(userID string) error
}

// UserWriteRepo implements UserWriteRepository
type UserWriteRepo struct {
	// Database connection for write operations
}

// Save saves a user to the write store
func (r *UserWriteRepo) Save(user *domain.User) (*domain.User, error) {
	// Implementation for saving to write database
	return user, nil
}

// Update updates a user in the write store
func (r *UserWriteRepo) Update(user *domain.User) error {
	// Implementation for updating in write database
	return nil
}

// Delete removes a user from the write store
func (r *UserWriteRepo) Delete(userID string) error {
	// Implementation for deleting from write database
	return nil
}
