package application

import (
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/user/domain"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/pkg/logger"
	// ⚠️  INTENTIONAL DDD VIOLATION ⚠️
	// User domain should not depend on products domain (bounded context isolation)
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/products/domain" // Cross-domain violation!
)

// UserWithProductViolation demonstrates cross-domain dependency violation
type UserWithProductViolation struct {
	userRepo    domain.UserRepository
	userService domain.UserService
	logger      logger.Logger
	// This creates a dependency from user domain to products domain - violation!
	productService productdomain.ProductService // Cross-domain violation!
}

// NewUserWithProductViolation creates a service with cross-domain dependency
func NewUserWithProductViolation(
	userRepo domain.UserRepository,
	userService domain.UserService,
	logger logger.Logger,
	productService productdomain.ProductService, // Cross-domain violation!
) *UserWithProductViolation {
	return &UserWithProductViolation{
		userRepo:       userRepo,
		userService:    userService,
		logger:         logger,
		productService: productService, // User domain accessing products domain - violation!
	}
}
