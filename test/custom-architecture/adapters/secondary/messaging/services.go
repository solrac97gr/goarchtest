package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/domain"
	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/ports"
)

// EmailNotificationService implements the NotificationService interface for email notifications
type EmailNotificationService struct {
	emailServer string
}

// NewEmailNotificationService creates a new EmailNotificationService
func NewEmailNotificationService(emailServer string) ports.NotificationService {
	return &EmailNotificationService{
		emailServer: emailServer,
	}
}

// SendOrderConfirmation sends an order confirmation email
func (s *EmailNotificationService) SendOrderConfirmation(ctx context.Context, order *domain.Order) error {
	// Simulate sending email
	log.Printf("Sending order confirmation email for order %s to customer %s", order.ID, order.CustomerID)
	return nil
}

// SendShippingNotification sends a shipping notification email
func (s *EmailNotificationService) SendShippingNotification(ctx context.Context, order *domain.Order) error {
	// Simulate sending email
	log.Printf("Sending shipping notification email for order %s to customer %s", order.ID, order.CustomerID)
	return nil
}

// MockPaymentService implements the PaymentService interface for testing
type MockPaymentService struct {
}

// NewMockPaymentService creates a new MockPaymentService
func NewMockPaymentService() ports.PaymentService {
	return &MockPaymentService{}
}

// ProcessPayment processes a payment
func (s *MockPaymentService) ProcessPayment(ctx context.Context, orderID string, amount float64) error {
	// Simulate payment processing
	log.Printf("Processing payment of $%.2f for order %s", amount, orderID)

	// Simulate some business logic
	if amount <= 0 {
		return fmt.Errorf("invalid payment amount: %.2f", amount)
	}

	return nil
}

// RefundPayment processes a refund
func (s *MockPaymentService) RefundPayment(ctx context.Context, orderID string, amount float64) error {
	// Simulate refund processing
	log.Printf("Processing refund of $%.2f for order %s", amount, orderID)
	return nil
}
