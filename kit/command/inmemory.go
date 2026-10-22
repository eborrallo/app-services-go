package command

import (
	"context"
	"fmt"
	"sync"
)

// CommandBus is an in-memory implementation of the command.Bus.
type CommandBus struct {
	mu       sync.RWMutex
	handlers map[Type]Handler
}

// Ensure *CommandBus satisfies the Bus interface.
var _ Bus = (*CommandBus)(nil)

// NewCommandBus initializes a new instance of CommandBus.
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[Type]Handler),
	}
}

// Dispatch implements the Bus interface.
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	b.mu.RLock()
	handler, ok := b.handlers[cmd.Type()]
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no handler registered for command %q", cmd.Type())
	}

	return handler.Handle(ctx, cmd)
}

// Register implements the Bus interface.
func (b *CommandBus) Register(cmdType Type, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[cmdType] = handler
}
