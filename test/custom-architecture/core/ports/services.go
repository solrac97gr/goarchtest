package ports

import (
	"context"

	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/domain"
)

// OrderService defines the primary port for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, customerID string, items []domain.OrderItem) (*domain.Order, error)
	GetOrder(ctx context.Context, orderID string) (*domain.Order, error)
	GetOrdersByCustomer(ctx context.Context, customerID string) ([]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus) error
	CancelOrder(ctx context.Context, orderID string) error
}

// OrderUseCase implements the business logic for orders
type OrderUseCase struct {
	orderRepo       OrderRepository
	paymentService  PaymentService
	notificationSvc NotificationService
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(
	orderRepo OrderRepository,
	paymentService PaymentService,
	notificationSvc NotificationService,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:       orderRepo,
		paymentService:  paymentService,
		notificationSvc: notificationSvc,
	}
}

// CreateOrder creates a new order
func (uc *OrderUseCase) CreateOrder(ctx context.Context, customerID string, items []domain.OrderItem) (*domain.Order, error) {
	order := &domain.Order{
		CustomerID: customerID,
		Items:      items,
		Status:     domain.OrderStatusPending,
	}

	order.CalculateTotal()

	if err := order.ValidateOrder(); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (uc *OrderUseCase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	return uc.orderRepo.FindByID(ctx, orderID)
}

// GetOrdersByCustomer retrieves all orders for a customer
func (uc *OrderUseCase) GetOrdersByCustomer(ctx context.Context, customerID string) ([]*domain.Order, error) {
	return uc.orderRepo.FindByCustomerID(ctx, customerID)
}

// UpdateOrderStatus updates the status of an order
func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = status
	return uc.orderRepo.Update(ctx, order)
}

// CancelOrder cancels an order
func (uc *OrderUseCase) CancelOrder(ctx context.Context, orderID string) error {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = domain.OrderStatusCancelled
	return uc.orderRepo.Update(ctx, order)
}
