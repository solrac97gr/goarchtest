package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// User represents a user entity in our domain
type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User entity with validation
func NewUser(username, email string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	// Simple email validation
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email format")
	}

	now := time.Now()
	return &User{
		ID:        generateID(),
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// generateID creates a simple ID for the user
// In a real application, you might use UUID or another ID generation mechanism
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Validate checks if the User entity is valid
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username cannot be empty")
	}

	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}

	return nil
}
