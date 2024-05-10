package remember

import (
	"context"
	"errors"

	"app-services-go/kit/command"
	"app-services-go/kit/crypt"
)

const AuthCommandType command.Type = "command.remember.user"

// PasswordUserCommand is the command dispatched to create a new user.
type PasswordUserCommand struct {
	Token       string
	NewPassword string
}

// NewPasswordUserCommand creates a new PasswordUserCommand.
func NewPasswordUserCommand(token, newPassword string) PasswordUserCommand {
	return PasswordUserCommand{
		Token:       token,
		NewPassword: newPassword,
	}
}

func (c PasswordUserCommand) Type() command.Type {
	return AuthCommandType
}

// PasswordUserCommandHandler is the command controllers
// responsible for creating PasswordUsers.
type PasswordUserCommandHandler struct {
	service PasswordUserService
}

// NewPasswordUserCommandHandler initializes a new PasswordUserCommandHandler.
func NewPasswordUserCommandHandler(service PasswordUserService) PasswordUserCommandHandler {
	return PasswordUserCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h PasswordUserCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	validateUserCmd, ok := cmd.(PasswordUserCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.UpdatePassword(
		ctx,
		crypt.Decrypt(validateUserCmd.Token),
		validateUserCmd.NewPassword,
	)
}
