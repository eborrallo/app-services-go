package query

import "context"

// QueryBus is an in-memory implementation of the Bus.
type QueryBus struct {
	handlers map[Type]Handler
}

// NewQueryBus initializes a new instance of QueryBus.
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[Type]Handler),
	}
}

// Ask implements the Bus interface.
func (b *QueryBus) Ask(ctx context.Context, cmd Query) (interface{}, error) {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return "", nil
	}

	return handler.Handle(ctx, cmd)
}

// Register implements the Bus interface.
func (b *QueryBus) Register(cmdType Type, handler Handler) {
	b.handlers[cmdType] = handler
}
