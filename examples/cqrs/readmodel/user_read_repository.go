package readmodel

import "github.com/solrac97gr/goarchtest/examples/cqrs/domain"

// UserReadRepository defines the interface for read operations
type UserReadRepository interface {
	GetByID(userID string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
}

// UserReadRepo implements UserReadRepository
type UserReadRepo struct {
	// Database connection for read operations (could be different from write)
}

// GetByID retrieves a user by ID from the read store
func (r *UserReadRepo) GetByID(userID string) (*domain.User, error) {
	// Implementation for reading from read database/cache
	return &domain.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil
}

// GetAll retrieves all users from the read store
func (r *UserReadRepo) GetAll() ([]*domain.User, error) {
	// Implementation for reading all users
	return []*domain.User{}, nil
}

// GetByEmail retrieves a user by email from the read store
func (r *UserReadRepo) GetByEmail(email string) (*domain.User, error) {
	// Implementation for reading by email
	return nil, nil
}
