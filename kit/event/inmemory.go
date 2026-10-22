package event

import (
	"context"
	"errors"
	"sync"
)

// EventBus is an in-memory implementation of the Bus.
type EventBus struct {
	mu       sync.RWMutex
	handlers map[Type][]Subscriber
}

// Ensure *EventBus satisfies the Bus interface.
var _ Bus = (*EventBus)(nil)

// NewEventBus initializes a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[Type][]Subscriber),
	}
}

// Publish implements the Bus interface.
func (b *EventBus) Publish(_ context.Context, events []Event) error {
	var err error
	for _, evt := range events {
		b.mu.RLock()
		handlers, ok := b.handlers[evt.Type()]
		b.mu.RUnlock()
		if !ok {
			continue
		}

		for _, handler := range handlers {
			if onErr := handler.On(evt); onErr != nil {
				err = errors.Join(err, onErr)
			}
		}
	}

	return err
}

// Subscribe implements the Bus interface.
func (b *EventBus) Subscribe(handler Subscriber) {
	evt := handler.SubscribedTo()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[evt.Type()] = append(b.handlers[evt.Type()], handler)
}
