package models

import (
	"time"

	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// Order represents an order entity in the order bounded context
type Order struct {
	ID         shared.ID    `json:"id"`
	UserID     shared.ID    `json:"user_id"`
	Items      []OrderItem  `json:"items"`
	Status     string       `json:"status"`
	Total      float64      `json:"total"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ProductID shared.ID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Subtotal  float64   `json:"subtotal"`
}

// NewOrder creates a new order for the given user
func NewOrder(userID shared.ID) *Order {
	return &Order{
		ID:        shared.NewID(),
		UserID:    userID,
		Items:     make([]OrderItem, 0),
		Status:    "pending",
		Total:     0.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddItem adds an item to the order
func (o *Order) AddItem(productID shared.ID, quantity int, price float64) {
	item := OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		Subtotal:  float64(quantity) * price,
	}
	o.Items = append(o.Items, item)
	o.calculateTotal()
	o.UpdatedAt = time.Now()
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(productID shared.ID) {
	for i, item := range o.Items {
		if item.ProductID == productID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			break
		}
	}
	o.calculateTotal()
	o.UpdatedAt = time.Now()
}

// Confirm confirms the order
func (o *Order) Confirm() {
	o.Status = "confirmed"
	o.UpdatedAt = time.Now()
}

// Cancel cancels the order
func (o *Order) Cancel() {
	o.Status = "cancelled"
	o.UpdatedAt = time.Now()
}

// IsConfirmed returns true if the order is confirmed
func (o *Order) IsConfirmed() bool {
	return o.Status == "confirmed"
}

// calculateTotal calculates the total amount of the order
func (o *Order) calculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Subtotal
	}
	o.Total = total
}