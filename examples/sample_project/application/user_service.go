package application

import (
	"errors"

	"github.com/solrac97gr/goarchtest/examples/sample_project/domain"
)

// UserService handles user-related operations
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	return s.userRepo.GetByID(id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, email string) (*domain.User, error) {
	if username == "" || email == "" {
		return nil, errors.New("username and email cannot be empty")
	}

	user := &domain.User{
		ID:       "user-123", // In a real application, this would be generated
		Username: username,
		Email:    email,
	}

	err := s.userRepo.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("user ID cannot be empty")
	}

	return s.userRepo.Delete(id)
}
