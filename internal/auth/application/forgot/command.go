package forgot

import (
	"context"
	"errors"

	"app-services-go/kit/command"
)

const AuthCommandType command.Type = "command.forgot.user"

// ForgotCommand is the command dispatched to create a new user.
type ForgotCommand struct {
	UserId    string
	UserEmail string
}

// NewForgotCommand creates a new ForgotCommand.
func NewForgotCommand(uerId, userEmail string) ForgotCommand {
	return ForgotCommand{
		UserId:    uerId,
		UserEmail: userEmail,
	}
}

func (c ForgotCommand) Type() command.Type {
	return AuthCommandType
}

// ForgotCommandHandler is the command controllers
// responsible for creating Forgots.
type ForgotCommandHandler struct {
	service ForgotUserService
}

// NewForgotCommandHandler initializes a new ForgotCommandHandler.
func NewForgotCommandHandler(service ForgotUserService) ForgotCommandHandler {
	return ForgotCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h ForgotCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createForgotCmd, ok := cmd.(ForgotCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.SendForgotEmail(
		ctx,
		createForgotCmd.UserId,
		createForgotCmd.UserEmail,
	)
}
