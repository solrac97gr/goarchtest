package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/solrac97gr/goarchtest/examples/sample_project/domain"
)

// SQLUserRepository implements the UserRepository interface using SQL
type SQLUserRepository struct {
	db *sql.DB
}

// NewSQLUserRepository creates a new SQLUserRepository
func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{
		db: db,
	}
}

// GetByID retrieves a user by ID from the database
func (r *SQLUserRepository) GetByID(id string) (*domain.User, error) {
	query := "SELECT id, username, email FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, err
	}

	return user, nil
}

// Save saves a user to the database
func (r *SQLUserRepository) Save(user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			username = VALUES(username),
			email = VALUES(email)
	`

	_, err := r.db.Exec(query, user.ID, user.Username, user.Email)
	return err
}

// Delete deletes a user from the database
func (r *SQLUserRepository) Delete(id string) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
