package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/domain"
	"github.com/solrac97gr/goarchtest/test/custom_architecture/core/ports"
)

// InMemoryOrderRepository implements the OrderRepository interface using in-memory storage
type InMemoryOrderRepository struct {
	orders map[string]*domain.Order
	mutex  sync.RWMutex
}

// NewInMemoryOrderRepository creates a new InMemoryOrderRepository
func NewInMemoryOrderRepository() ports.OrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

// Save stores an order in memory
func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if order.ID == "" {
		order.ID = fmt.Sprintf("order_%d", len(r.orders)+1)
	}

	r.orders[order.ID] = order
	return nil
}

// FindByID retrieves an order by ID
func (r *InMemoryOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("order with ID %s not found", id)
	}

	return order, nil
}

// FindByCustomerID retrieves all orders for a customer
func (r *InMemoryOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var customerOrders []*domain.Order
	for _, order := range r.orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}

	return customerOrders, nil
}

// Update updates an existing order
func (r *InMemoryOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return fmt.Errorf("order with ID %s not found", order.ID)
	}

	r.orders[order.ID] = order
	return nil
}

// Delete removes an order from storage
func (r *InMemoryOrderRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.orders[id]; !exists {
		return fmt.Errorf("order with ID %s not found", id)
	}

	delete(r.orders, id)
	return nil
}
