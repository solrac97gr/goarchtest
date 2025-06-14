package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/domain"
	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/ports"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	orderService ports.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService ports.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrderRequest represents the request payload for creating an order
type CreateOrderRequest struct {
	CustomerID string             `json:"customer_id"`
	Items      []domain.OrderItem `json:"items"`
}

// OrderResponse represents the response payload for order operations
type OrderResponse struct {
	ID          string             `json:"id"`
	CustomerID  string             `json:"customer_id"`
	Items       []domain.OrderItem `json:"items"`
	TotalAmount float64            `json:"total_amount"`
	Status      domain.OrderStatus `json:"status"`
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.CreateOrder(context.Background(), req.CustomerID, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := OrderResponse{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Items:       order.Items,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetOrder handles GET /orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// This is a simplified example - in real implementation you'd extract ID from URL
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrder(context.Background(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := OrderResponse{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Items:       order.Items,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
