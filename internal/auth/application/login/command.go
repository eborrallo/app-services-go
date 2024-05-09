package login

import (
	"context"
	"errors"

	"app-services-go/kit/command"
)

const AuthCommandType command.Type = "command.creating.user"

// LoginCommand is the command dispatched to create a new user.
type LoginCommand struct {
	token        string
	refreshToken string
}

// NewLoginCommand creates a new LoginCommand.
func NewLoginCommand(token, refreshToken string) LoginCommand {
	return LoginCommand{
		token:        token,
		refreshToken: refreshToken,
	}
}

func (c LoginCommand) Type() command.Type {
	return AuthCommandType
}

// LoginCommandHandler is the command controllers
// responsible for creating Logins.
type LoginCommandHandler struct {
	service LoginUserService
}

// NewLoginCommandHandler initializes a new LoginCommandHandler.
func NewLoginCommandHandler(service LoginUserService) LoginCommandHandler {
	return LoginCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h LoginCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createLoginCmd, ok := cmd.(LoginCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.LoginUser(
		ctx,
		createLoginCmd.token,
		createLoginCmd.refreshToken,
	)
}
