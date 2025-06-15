package application

import (
	"fmt"

	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/order/domain/models"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/order/domain/ports"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// OrderService provides application services for order management
type OrderService struct {
	orderRepo       ports.OrderRepository
	paymentSvc      ports.PaymentService
	inventorySvc    ports.InventoryService
	notificationSvc ports.OrderNotificationService
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo ports.OrderRepository,
	paymentSvc ports.PaymentService,
	inventorySvc ports.InventoryService,
	notificationSvc ports.OrderNotificationService,
) *OrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		paymentSvc:      paymentSvc,
		inventorySvc:    inventorySvc,
		notificationSvc: notificationSvc,
	}
}

// CreateOrder creates a new order for a user
func (s *OrderService) CreateOrder(userID shared.ID) (*models.Order, error) {
	order := models.NewOrder(userID)
	
	if err := s.orderRepo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

// AddItemToOrder adds an item to an existing order
func (s *OrderService) AddItemToOrder(orderID, productID shared.ID, quantity int, price float64) error {
	// Get the order
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Check inventory availability
	available, err := s.inventorySvc.CheckAvailability(productID, quantity)
	if err != nil {
		return fmt.Errorf("failed to check inventory: %w", err)
	}
	if !available {
		return fmt.Errorf("insufficient inventory for product %s", productID)
	}

	// Add item to order
	order.AddItem(productID, quantity, price)

	// Save updated order
	if err := s.orderRepo.Save(order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}

// ConfirmOrder confirms an order and processes payment
func (s *OrderService) ConfirmOrder(orderID shared.ID) error {
	// Get the order
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Reserve inventory items
	if err := s.inventorySvc.ReserveItems(order.Items); err != nil {
		return fmt.Errorf("failed to reserve items: %w", err)
	}

	// Process payment
	if err := s.paymentSvc.ProcessPayment(orderID, order.Total); err != nil {
		// Release reserved items on payment failure
		if releaseErr := s.inventorySvc.ReleaseItems(order.Items); releaseErr != nil {
			fmt.Printf("Failed to release items after payment failure: %v\n", releaseErr)
		}
		return fmt.Errorf("payment failed: %w", err)
	}

	// Confirm the order
	order.Confirm()

	// Save updated order
	if err := s.orderRepo.Save(order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	// Send confirmation notification
	if err := s.notificationSvc.SendOrderConfirmation(order); err != nil {
		fmt.Printf("Failed to send order confirmation: %v\n", err)
	}

	return nil
}

// CancelOrder cancels an existing order
func (s *OrderService) CancelOrder(orderID shared.ID) error {
	// Get the order
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// If order is confirmed, refund payment and release inventory
	if order.IsConfirmed() {
		if err := s.paymentSvc.RefundPayment(orderID, order.Total); err != nil {
			return fmt.Errorf("failed to refund payment: %w", err)
		}

		if err := s.inventorySvc.ReleaseItems(order.Items); err != nil {
			fmt.Printf("Failed to release items: %v\n", err)
		}
	}

	// Cancel the order
	order.Cancel()

	// Save updated order
	if err := s.orderRepo.Save(order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	// Send cancellation notification
	if err := s.notificationSvc.SendOrderCancellation(order); err != nil {
		fmt.Printf("Failed to send order cancellation: %v\n", err)
	}

	return nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(orderID shared.ID) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	return order, nil
}

// GetUserOrders retrieves all orders for a user
func (s *OrderService) GetUserOrders(userID shared.ID) ([]*models.Order, error) {
	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	return orders, nil
}