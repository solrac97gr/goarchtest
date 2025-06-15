package models

import (
	"time"

	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// User represents a user entity in the user bounded context
type User struct {
	ID       shared.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"created_at"`
}

// NewUser creates a new user with the given details
func NewUser(name, email string) *User {
	return &User{
		ID:       shared.NewID(),
		Name:     name,
		Email:    email,
		Status:   "active",
		CreateAt: time.Now(),
	}
}

// Activate sets the user status to active
func (u *User) Activate() {
	u.Status = "active"
}

// Deactivate sets the user status to inactive
func (u *User) Deactivate() {
	u.Status = "inactive"
}

// IsActive returns true if the user is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}
