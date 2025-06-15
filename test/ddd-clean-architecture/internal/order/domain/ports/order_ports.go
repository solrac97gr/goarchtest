package ports

import (
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/order/domain/models"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// OrderRepository defines the contract for order data persistence
type OrderRepository interface {
	Save(order *models.Order) error
	FindByID(id shared.ID) (*models.Order, error)
	FindByUserID(userID shared.ID) ([]*models.Order, error)
	Delete(id shared.ID) error
	List() ([]*models.Order, error)
}

// PaymentService defines the contract for payment processing
type PaymentService interface {
	ProcessPayment(orderID shared.ID, amount float64) error
	RefundPayment(orderID shared.ID, amount float64) error
	ValidatePaymentMethod(method string) error
}

// InventoryService defines the contract for inventory management
type InventoryService interface {
	ReserveItems(items []models.OrderItem) error
	ReleaseItems(items []models.OrderItem) error
	CheckAvailability(productID shared.ID, quantity int) (bool, error)
}

// OrderNotificationService defines the contract for order notifications
type OrderNotificationService interface {
	SendOrderConfirmation(order *models.Order) error
	SendOrderCancellation(order *models.Order) error
	SendOrderStatusUpdate(order *models.Order) error
}