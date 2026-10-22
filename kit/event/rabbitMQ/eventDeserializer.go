package rabbitMQ

import (
	"app-services-go/kit/event"
	"fmt"
)

type EventDeserializer struct {
	subscriber event.Subscriber
}

func Configure(subscriber event.Subscriber) *EventDeserializer {
	return &EventDeserializer{subscriber: subscriber}
}

type EventPayload struct {
	event.Event
}

func (e *EventDeserializer) Deserialize(message []byte) (event.Event, error) {
	eventType := e.subscriber.SubscribedTo()

	domainEvent, err := eventType.FromJSON(message)
	if err != nil {
		return domainEvent, fmt.Errorf("error deserializing event %s: %w", eventType.Type(), err)
	}

	return domainEvent, nil
}
