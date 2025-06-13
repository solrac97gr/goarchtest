package domain

// User represents the domain entity
type User struct {
	ID    string
	Name  string
	Email string
}

// NewUser creates a new User
func NewUser(name, email string) *User {
	return &User{
		ID:    generateID(), // Assume this function exists
		Name:  name,
		Email: email,
	}
}

func generateID() string {
	return "user-123" // Simplified for example
}
