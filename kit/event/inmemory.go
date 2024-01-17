package event

// EventBus is an in-memory implementation of the Bus.
type EventBus struct {
	handlers map[Type][]Subscriber
}

// NewEventBus initializes a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[Type][]Subscriber),
	}
}

// Publish implements the Bus interface.
func (b *EventBus) Publish(events []Event) error {
	for _, evt := range events {
		handlers, ok := b.handlers[evt.Type()]
		if !ok {
			return nil
		}

		for _, handler := range handlers {
			handler.On(evt)
		}
	}

	return nil
}

// Subscribe implements the Bus interface.
func (b *EventBus) Subscribe(handler Subscriber) {
	evt := handler.SubscribedTo()
	subscribersForType, ok := b.handlers[evt.Type()]
	if !ok {
		b.handlers[evt.Type()] = []Subscriber{handler}
	}

	subscribersForType = append(subscribersForType, handler)

}
