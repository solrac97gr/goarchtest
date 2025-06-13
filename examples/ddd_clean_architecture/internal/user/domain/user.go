package domain

import (
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/shared"
)

// User represents a user entity in the system
type User struct {
	ID       UserID
	Email    Email
	Username string
	Status   shared.EntityStatus // Using shared kernel
}

// UserID is a value object for user identification
type UserID struct {
	value string
}

// Email is a value object for email addresses
type Email struct {
	value string
}

// UserRepository defines the contract for user persistence
type UserRepository interface {
	Save(user *User) error
	FindByID(id UserID) (*User, error)
	FindByEmail(email Email) (*User, error)
}

// UserService defines domain services for user operations
type UserService interface {
	ValidateUser(user *User) error
	CanUserPerformAction(userID UserID, action string) bool
}
