package queries

import (
	"fmt"

	"github.com/solrac97gr/goarchtest/examples/cqrs/commands" // ❌ VIOLATION: Queries should not depend on commands
	"github.com/solrac97gr/goarchtest/examples/cqrs/domain"
	"github.com/solrac97gr/goarchtest/examples/cqrs/writemodel" // ❌ VIOLATION: Queries should use read models, not write models
)

// BadGetUserQuery demonstrates architectural violations
type BadGetUserQuery struct {
	UserID string
}

// BadGetUserHandler violates CQRS principles
type BadGetUserHandler struct {
	userWriteRepo writemodel.UserWriteRepository // ❌ VIOLATION: Using write repository in query
	createCommand commands.CreateUserCommand     // ❌ VIOLATION: Queries depending on commands
}

// Handle processes the query but violates CQRS principles
func (h *BadGetUserHandler) Handle(query BadGetUserQuery) (*domain.User, error) {
	// ❌ VIOLATION: Query is trying to use write repository
	// This violates the separation between read and write models
	return nil, fmt.Errorf("bad implementation - using write repo in query")
}
