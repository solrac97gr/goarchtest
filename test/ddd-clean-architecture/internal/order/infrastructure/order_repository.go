package infrastructure

import (
	"errors"
	"sync"

	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/order/domain/models"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/internal/order/domain/ports"
	"github.com/solrac97gr/goarchtest/test/ddd_clean_architecture/shared"
)

// InMemoryOrderRepository provides an in-memory implementation of OrderRepository
type InMemoryOrderRepository struct {
	orders map[shared.ID]*models.Order
	mutex  sync.RWMutex
}

// NewInMemoryOrderRepository creates a new in-memory order repository
func NewInMemoryOrderRepository() ports.OrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[shared.ID]*models.Order),
	}
}

// Save saves an order to the repository
func (r *InMemoryOrderRepository) Save(order *models.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.orders[order.ID] = order
	return nil
}

// FindByID finds an order by ID
func (r *InMemoryOrderRepository) FindByID(id shared.ID) (*models.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

// FindByUserID finds all orders for a user
func (r *InMemoryOrderRepository) FindByUserID(userID shared.ID) ([]*models.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var userOrders []*models.Order
	for _, order := range r.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}
	return userOrders, nil
}

// Delete deletes an order by ID
func (r *InMemoryOrderRepository) Delete(id shared.ID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.orders, id)
	return nil
}

// List returns all orders
func (r *InMemoryOrderRepository) List() ([]*models.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	orders := make([]*models.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	return orders, nil
}