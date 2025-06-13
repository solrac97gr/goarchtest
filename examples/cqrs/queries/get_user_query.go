package queries

import (
	"github.com/solrac97gr/goarchtest/examples/cqrs/domain"
	"github.com/solrac97gr/goarchtest/examples/cqrs/readmodel"
)

// GetUserQuery represents a query to get user information
type GetUserQuery struct {
	UserID string
}

// GetUserHandler handles the GetUserQuery
type GetUserHandler struct {
	userReadRepo readmodel.UserReadRepository
}

// Handle processes the GetUserQuery
func (h *GetUserHandler) Handle(query GetUserQuery) (*domain.User, error) {
	return h.userReadRepo.GetByID(query.UserID)
}
