package validating

import (
	"context"
	"errors"

	"app-services-go/kit/command"
	"app-services-go/kit/crypt"
)

const AuthCommandType command.Type = "command.validating.user"

// ValidateUserCommand is the command dispatched to create a new user.
type ValidateUserCommand struct {
	Token string
}

// NewValidateUserCommand creates a new ValidateUserCommand.
func NewValidateUserCommand(token string) ValidateUserCommand {
	return ValidateUserCommand{
		Token: token,
	}
}

func (c ValidateUserCommand) Type() command.Type {
	return AuthCommandType
}

// ValidateUserCommandHandler is the command controllers
// responsible for creating ValidateUsers.
type ValidateUserCommandHandler struct {
	service ValidateUserService
}

// NewValidateUserCommandHandler initializes a new ValidateUserCommandHandler.
func NewValidateUserCommandHandler(service ValidateUserService) ValidateUserCommandHandler {
	return ValidateUserCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h ValidateUserCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	validateUserCmd, ok := cmd.(ValidateUserCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.ValidateUser(
		ctx,
		crypt.Decrypt(validateUserCmd.Token),
	)
}
