package infrastructure

import (
	"database/sql"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/products/domain"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/pkg/config"
)

// PostgreSQLProductRepository implements ProductRepository using PostgreSQL
type PostgreSQLProductRepository struct {
	db     *sql.DB
	config config.Config
}

// NewPostgreSQLProductRepository creates a new PostgreSQL product repository
func NewPostgreSQLProductRepository(db *sql.DB, config config.Config) *PostgreSQLProductRepository {
	return &PostgreSQLProductRepository{
		db:     db,
		config: config,
	}
}

// Save implements domain.ProductRepository
func (r *PostgreSQLProductRepository) Save(product *domain.Product) error {
	query := "INSERT INTO products (id, name, description, price, category, status) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.Exec(query, product.ID, product.Name, product.Description, product.Price, product.Category, product.Status)
	return err
}

// FindByID implements domain.ProductRepository
func (r *PostgreSQLProductRepository) FindByID(id domain.ProductID) (*domain.Product, error) {
	var product domain.Product
	query := "SELECT id, name, description, price, category, status FROM products WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Status)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindByCategory implements domain.ProductRepository
func (r *PostgreSQLProductRepository) FindByCategory(category domain.ProductCategory) ([]*domain.Product, error) {
	query := "SELECT id, name, description, price, category, status FROM products WHERE category = $1"
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Status)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
