package domain

// User represents a user in the system
type User struct {
	ID       string
	Username string
	Email    string
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(id string) (*User, error)
	Save(user *User) error
	Delete(id string) error
}
