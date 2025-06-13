package commands

import (
	"fmt"

	"github.com/solrac97gr/goarchtest/examples/cqrs/domain"
	"github.com/solrac97gr/goarchtest/examples/cqrs/queries"   // ❌ VIOLATION: Commands should not depend on queries
	"github.com/solrac97gr/goarchtest/examples/cqrs/readmodel" // ❌ VIOLATION: Commands should use write models, not read models
)

// BadCreateUserCommand demonstrates architectural violations
type BadCreateUserCommand struct {
	Name  string
	Email string
}

// BadCreateUserHandler violates CQRS principles
type BadCreateUserHandler struct {
	userReadRepo readmodel.UserReadRepository // ❌ VIOLATION: Using read repository in command
	getUserQuery queries.GetUserQuery         // ❌ VIOLATION: Commands depending on queries
}

// Handle processes the command but violates CQRS principles
func (h *BadCreateUserHandler) Handle(cmd BadCreateUserCommand) (*domain.User, error) {
	// ❌ VIOLATION: Command is querying data instead of just executing commands
	existingUser, _ := h.userReadRepo.GetByEmail(cmd.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	user := domain.NewUser(cmd.Name, cmd.Email)
	// Missing proper write repository usage - this would be another violation
	return user, nil
}
