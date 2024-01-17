package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Bus defines the expected behaviour from an event bus.
type Bus interface {
	// Publish is the method used to publish new events.
	Publish(context.Context, []Event) error
	// Subscribe is the method used to subscribe new event handlers.
	Subscribe(Subscriber)
}

//go:generate mockery --case=snake --outpkg=eventmocks --output=eventmocks --name=Bus

// Subscriber defines the expected behaviour from an event controllers.
type Subscriber interface {
	On(Event) error
	SubscribedTo() Event
}

// Type represents a domain event type.
type Type string

// Event represents a domain command.
type Event interface {
	ID() string
	Type() Type
	FromJSON([]byte) (Event, error)
}

type BaseEvent struct {
	EventID     string    `json:"eventId"`
	AggregateID string    `json:"aggregateId"`
	OccurredOn  time.Time `json:"occurredOn"`
}

func NewBaseEvent(aggregateID string) BaseEvent {
	return BaseEvent{
		EventID:     uuid.New().String(),
		AggregateID: aggregateID,
		OccurredOn:  time.Now(),
	}
}

func (b BaseEvent) ID() string {
	return b.EventID
}
