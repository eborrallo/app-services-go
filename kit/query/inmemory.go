package query

import (
	"context"
	"fmt"
	"sync"
)

// QueryBus is an in-memory implementation of the Bus.
type QueryBus struct {
	mu       sync.RWMutex
	handlers map[Type]Handler
}

// Ensure *QueryBus satisfies the Bus interface.
var _ Bus = (*QueryBus)(nil)

// NewQueryBus initializes a new instance of QueryBus.
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[Type]Handler),
	}
}

// Ask implements the Bus interface.
func (b *QueryBus) Ask(ctx context.Context, q Query) (interface{}, error) {
	b.mu.RLock()
	handler, ok := b.handlers[q.Type()]
	b.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handler registered for query %q", q.Type())
	}

	return handler.Handle(ctx, q)
}

// Register implements the Bus interface.
func (b *QueryBus) Register(cmdType Type, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[cmdType] = handler
}
