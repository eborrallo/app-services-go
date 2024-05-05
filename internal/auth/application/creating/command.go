package creating

import (
	"context"
	"errors"

	"app-services-go/kit/command"
)

const AuthCommandType command.Type = "command.creating.user"

// UserCommand is the command dispatched to create a new user.
type UserCommand struct {
	id       string
	name     string
	email    string
	password string
}

// NewUserCommand creates a new UserCommand.
func NewUserCommand(id, name, email, password string) UserCommand {
	return UserCommand{
		id:       id,
		name:     name,
		email:    email,
		password: password,
	}
}

func (c UserCommand) Type() command.Type {
	return AuthCommandType
}

// UserCommandHandler is the command controllers
// responsible for creating Users.
type UserCommandHandler struct {
	service UserService
}

// NewUserCommandHandler initializes a new UserCommandHandler.
func NewUserCommandHandler(service UserService) UserCommandHandler {
	return UserCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h UserCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createUserCmd, ok := cmd.(UserCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.CreateUser(
		ctx,
		createUserCmd.id,
		createUserCmd.name,
		createUserCmd.email,
		createUserCmd.password,
	)
}
