package infrastructure

import (
	"database/sql"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/user/domain"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/pkg/config"
)

// PostgreSQLUserRepository implements UserRepository using PostgreSQL
type PostgreSQLUserRepository struct {
	db     *sql.DB
	config config.Config
}

// NewPostgreSQLUserRepository creates a new PostgreSQL user repository
func NewPostgreSQLUserRepository(db *sql.DB, config config.Config) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{
		db:     db,
		config: config,
	}
}

// Save implements domain.UserRepository
func (r *PostgreSQLUserRepository) Save(user *domain.User) error {
	// Implementation for saving user to PostgreSQL
	query := "INSERT INTO users (id, email, username, status) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.Status)
	return err
}

// FindByID implements domain.UserRepository
func (r *PostgreSQLUserRepository) FindByID(id domain.UserID) (*domain.User, error) {
	// Implementation for finding user by ID
	var user domain.User
	query := "SELECT id, email, username, status FROM users WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail implements domain.UserRepository
func (r *PostgreSQLUserRepository) FindByEmail(email domain.Email) (*domain.User, error) {
	// Implementation for finding user by email
	var user domain.User
	query := "SELECT id, email, username, status FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
