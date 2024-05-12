package creating

import (
	"context"
	"errors"

	"app-services-go/kit/command"
)

const Web3AuthCommandType command.Type = "command.creating.web3_user"

// Web3UserCommand is the command dispatched to create a new user.
type Web3UserCommand struct {
	address      string
	token        string
	refreshToken string
}

// NewWeb3UserCommand creates a new Web3UserCommand.
func NewWeb3UserCommand(address string, token string, refreshToken string) Web3UserCommand {
	return Web3UserCommand{
		address:      address,
		token:        token,
		refreshToken: refreshToken,
	}
}

func (c Web3UserCommand) Type() command.Type {
	return Web3AuthCommandType
}

// Web3UserCommandHandler is the command controllers
// responsible for creating Web3Users.
type Web3UserCommandHandler struct {
	service UserService
}

// NewWeb3UserCommandHandler initializes a new Web3UserCommandHandler.
func NewWeb3UserCommandHandler(service UserService) Web3UserCommandHandler {
	return Web3UserCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h Web3UserCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createWeb3UserCmd, ok := cmd.(Web3UserCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.CreateWeb3User(
		ctx,
		createWeb3UserCmd.address,
		createWeb3UserCmd.token,
		createWeb3UserCmd.refreshToken,
	)
}
