package ports

import (
	"context"

	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/domain"
)

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id string) error
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	SendOrderConfirmation(ctx context.Context, order *domain.Order) error
	SendShippingNotification(ctx context.Context, order *domain.Order) error
}

// PaymentService defines the interface for payment processing
type PaymentService interface {
	ProcessPayment(ctx context.Context, orderID string, amount float64) error
	RefundPayment(ctx context.Context, orderID string, amount float64) error
}
