package application

import (
	"fmt"

	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/user/domain/models"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/user/domain/ports"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// UserService provides application services for user management
type UserService struct {
	userRepo         ports.UserRepository
	notificationSvc  ports.UserNotificationService
	validator        ports.UserValidator
}

// NewUserService creates a new user service
func NewUserService(
	userRepo ports.UserRepository,
	notificationSvc ports.UserNotificationService,
	validator ports.UserValidator,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		notificationSvc: notificationSvc,
		validator:       validator,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name, email string) (*models.User, error) {
	// Validate input
	if err := s.validator.ValidateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}
	
	if err := s.validator.ValidateName(name); err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Create new user
	user := models.NewUser(name, email)
	
	// Validate user
	if err := s.validator.ValidateUser(user); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// Save user
	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Send welcome email
	if err := s.notificationSvc.SendWelcomeEmail(user); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to send welcome email: %v\n", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id shared.ID) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(id shared.ID) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Activate()
	
	if err := s.userRepo.Save(user); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Send activation email
	if err := s.notificationSvc.SendAccountActivation(user); err != nil {
		fmt.Printf("Failed to send activation email: %v\n", err)
	}

	return nil
}