package domain

import "time"

// Order represents a business order
type Order struct {
	ID          string
	CustomerID  string
	Items       []OrderItem
	TotalAmount float64
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// ValidateOrder validates an order before processing
func (o *Order) ValidateOrder() error {
	if o.CustomerID == "" {
		return &DomainError{Message: "customer ID is required"}
	}
	if len(o.Items) == 0 {
		return &DomainError{Message: "order must have at least one item"}
	}
	return nil
}

// CalculateTotal calculates the total amount of the order
func (o *Order) CalculateTotal() {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.TotalAmount = total
}
