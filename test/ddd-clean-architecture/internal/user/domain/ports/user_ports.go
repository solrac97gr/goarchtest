package ports

import (
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/user/domain/models"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// UserRepository defines the contract for user data persistence
type UserRepository interface {
	Save(user *models.User) error
	FindByID(id shared.ID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Delete(id shared.ID) error
	List() ([]*models.User, error)
}

// UserNotificationService defines the contract for user notifications
type UserNotificationService interface {
	SendWelcomeEmail(user *models.User) error
	SendAccountActivation(user *models.User) error
	SendPasswordReset(email string) error
}

// UserValidator defines the contract for user validation
type UserValidator interface {
	ValidateEmail(email string) error
	ValidateName(name string) error
	ValidateUser(user *models.User) error
}