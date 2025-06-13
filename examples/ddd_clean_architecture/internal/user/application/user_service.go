package application

import (
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/user/domain"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/pkg/logger"
)

// UserApplicationService handles user use cases
type UserApplicationService struct {
	userRepo    domain.UserRepository
	userService domain.UserService
	logger      logger.Logger
}

// NewUserApplicationService creates a new user application service
func NewUserApplicationService(
	userRepo domain.UserRepository,
	userService domain.UserService,
	logger logger.Logger,
) *UserApplicationService {
	return &UserApplicationService{
		userRepo:    userRepo,
		userService: userService,
		logger:      logger,
	}
}

// CreateUser creates a new user
func (s *UserApplicationService) CreateUser(email, username string) (*domain.User, error) {
	// Use case logic for creating a user
	user := &domain.User{
		Email:    domain.Email{},
		Username: username,
	}

	if err := s.userService.ValidateUser(user); err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(user); err != nil {
		s.logger.Error("Failed to save user", err)
		return nil, err
	}

	s.logger.Info("User created successfully", user.ID)
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserApplicationService) GetUser(userID domain.UserID) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}
