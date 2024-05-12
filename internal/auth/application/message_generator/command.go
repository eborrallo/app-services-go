package message_generator

import (
	"context"
	"errors"

	"app-services-go/kit/command"
)

const AuthCommandType command.Type = "command.login.user"

type SaveMessageCommand struct {
	Address string
	Message string
}

func NewSaveMessageCommand(address, message string) SaveMessageCommand {
	return SaveMessageCommand{
		Address: address,
		Message: message,
	}
}

func (c SaveMessageCommand) Type() command.Type {
	return AuthCommandType
}

type SaveMessageCommandHandler struct {
	service SaveMessageService
}

func NewSaveMessageCommandHandler(service SaveMessageService) SaveMessageCommandHandler {
	return SaveMessageCommandHandler{
		service: service,
	}
}

func (h SaveMessageCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createSaveMessageCmd, ok := cmd.(SaveMessageCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	h.service.Save(
		createSaveMessageCmd.Address,
		createSaveMessageCmd.Message,
	)
	return nil
}
