package application

import (
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/products/domain"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/pkg/logger"
	// ⚠️  INTENTIONAL DDD VIOLATION ⚠️
	// Application should not depend on shared directly (only domain should)
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/shared" // Shared violation!
)

// ProductWithSharedViolation demonstrates application layer using shared directly
type ProductWithSharedViolation struct {
	productRepo    domain.ProductRepository
	productService domain.ProductService
	logger         logger.Logger
	// This creates a direct dependency from application to shared - violation!
	status shared.EntityStatus // Application using shared directly - violation!
}

// NewProductWithSharedViolation creates a service that violates shared usage
func NewProductWithSharedViolation(
	productRepo domain.ProductRepository,
	productService domain.ProductService,
	logger logger.Logger,
) *ProductWithSharedViolation {
	return &ProductWithSharedViolation{
		productRepo:    productRepo,
		productService: productService,
		logger:         logger,
		status:         shared.StatusActive, // Application using shared directly - violation!
	}
}
