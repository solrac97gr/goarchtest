package application

import (
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/internal/products/domain"
	"github.com/solrac97gr/goarchtest/examples/ddd_clean_architecture/pkg/logger"
)

// ProductApplicationService handles product use cases
type ProductApplicationService struct {
	productRepo    domain.ProductRepository
	productService domain.ProductService
	logger         logger.Logger
}

// NewProductApplicationService creates a new product application service
func NewProductApplicationService(
	productRepo domain.ProductRepository,
	productService domain.ProductService,
	logger logger.Logger,
) *ProductApplicationService {
	return &ProductApplicationService{
		productRepo:    productRepo,
		productService: productService,
		logger:         logger,
	}
}

// CreateProduct creates a new product
func (s *ProductApplicationService) CreateProduct(name, description string, price domain.Money) (*domain.Product, error) {
	product := &domain.Product{
		Name:        name,
		Description: description,
		Price:       price,
	}

	if err := s.productService.ValidateProduct(product); err != nil {
		return nil, err
	}

	if err := s.productRepo.Save(product); err != nil {
		s.logger.Error("Failed to save product", err)
		return nil, err
	}

	s.logger.Info("Product created successfully", product.ID)
	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductApplicationService) GetProduct(productID domain.ProductID) (*domain.Product, error) {
	return s.productRepo.FindByID(productID)
}
