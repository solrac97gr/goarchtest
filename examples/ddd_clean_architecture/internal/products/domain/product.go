package domain

import (
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/shared"
)

// Product represents a product entity
type Product struct {
	ID          ProductID
	Name        string
	Description string
	Price       Money
	Category    ProductCategory
	Status      shared.EntityStatus // Using shared kernel
}

// ProductID is a value object for product identification
type ProductID struct {
	value string
}

// Money is a value object for monetary amounts
type Money struct {
	amount   int64
	currency string
}

// ProductCategory is a value object for product categorization
type ProductCategory struct {
	name string
}

// ProductRepository defines the contract for product persistence
type ProductRepository interface {
	Save(product *Product) error
	FindByID(id ProductID) (*Product, error)
	FindByCategory(category ProductCategory) ([]*Product, error)
}

// ProductService defines domain services for product operations
type ProductService interface {
	CalculateDiscount(product *Product, discountRate float64) Money
	ValidateProduct(product *Product) error
}
