package domain

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// GetByID retrieves a user by their ID
	GetByID(id string) (*User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(email string) (*User, error)

	// Save persists a user
	Save(user *User) error

	// Update updates an existing user
	Update(user *User) error

	// Delete removes a user
	Delete(id string) error
}
