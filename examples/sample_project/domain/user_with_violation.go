package domain

import (
	// ⚠️  INTENTIONAL CLEAN ARCHITECTURE VIOLATION ⚠️
	// This import creates a Clean Architecture violation!
	// Domain layer should NOT depend on infrastructure layer
	// This is included to demonstrate how GoArchTest detects violations
	"github.com/solrac97gr/goarchtest/examples/sample_project/infrastructure"
)

// UserWithViolation demonstrates a Clean Architecture violation
// by having the domain layer depend on the infrastructure layer
type UserWithViolation struct {
	ID       string
	Username string
	Email    string
	// This field creates a dependency from domain to infrastructure
	// which violates Clean Architecture principles
	cache *infrastructure.UserCache
}

// NewUserWithViolation creates a user with a direct dependency on infrastructure
// This is a violation of Clean Architecture
func NewUserWithViolation(id, username, email string, cache *infrastructure.UserCache) *UserWithViolation {
	return &UserWithViolation{
		ID:       id,
		Username: username,
		Email:    email,
		cache:    cache,
	}
}
