package commands

import (
	"github.com/solrac97gr/goarchtest/examples/cqrs/domain"
	"github.com/solrac97gr/goarchtest/examples/cqrs/writemodel"
)

// CreateUserCommand represents a command to create a new user
type CreateUserCommand struct {
	Name  string
	Email string
}

// CreateUserHandler handles the CreateUserCommand
type CreateUserHandler struct {
	userWriteRepo writemodel.UserWriteRepository
}

// Handle processes the CreateUserCommand
func (h *CreateUserHandler) Handle(cmd CreateUserCommand) (*domain.User, error) {
	user := domain.NewUser(cmd.Name, cmd.Email)
	return h.userWriteRepo.Save(user)
}
