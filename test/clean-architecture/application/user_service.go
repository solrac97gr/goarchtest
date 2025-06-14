package application

import (
	"errors"

	"github.com/solrac97gr/goarchtest/test/clean_architecture/domain"
)

// UserService contains the business logic for user management
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUserRequest represents the data needed to create a user
type CreateUserRequest struct {
	Username string
	Email    string
}

// CreateUserResponse contains the data returned after creating a user
type CreateUserResponse struct {
	ID       string
	Username string
	Email    string
}

// CreateUser handles user creation
func (s *UserService) CreateUser(req CreateUserRequest) (*CreateUserResponse, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user, err := domain.NewUser(req.Username, req.Email)
	if err != nil {
		return nil, err
	}

	// Save user
	err = s.userRepo.Save(user)
	if err != nil {
		return nil, err
	}

	// Return response
	return &CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// GetUserByIDRequest represents the data needed to get a user by ID
type GetUserByIDRequest struct {
	ID string
}

// GetUserByIDResponse contains the data returned after getting a user
type GetUserByIDResponse struct {
	ID       string
	Username string
	Email    string
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(req GetUserByIDRequest) (*GetUserByIDResponse, error) {
	if req.ID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	user, err := s.userRepo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &GetUserByIDResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
