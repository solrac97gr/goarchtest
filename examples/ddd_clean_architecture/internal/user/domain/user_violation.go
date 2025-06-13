package domain

import (
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/shared"
	// ⚠️  INTENTIONAL DDD VIOLATION ⚠️
	// Domain should not depend on application layer
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/user/application"
)

// UserWithViolation demonstrates a violation where domain depends on application
type UserWithViolation struct {
	ID       UserID
	Email    Email
	Username string
	Status   shared.EntityStatus
	// This creates a dependency from domain to application - violation!
	service *application.UserApplicationService
}

// NewUserWithViolation creates a user with application dependency (violation)
func NewUserWithViolation(service *application.UserApplicationService) *UserWithViolation {
	return &UserWithViolation{
		service: service, // Domain depending on application - violation!
	}
}
