package rabbitMQ

import (
	"app-services-go/kit/event"
	"errors"
	"fmt"
	"reflect"
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
		return domainEvent, errors.New(fmt.Sprintf("Error deserializing event %s", eventType.Type()))
	}

	if err := e.areStructPropertiesPopulated(domainEvent); err != nil {
		return domainEvent, errors.New(fmt.Sprintf("Event %v is not properly populated : %v ", reflect.TypeOf(eventType), err))
	}

	return domainEvent, nil
}

func (e *EventDeserializer) areStructPropertiesPopulated(eventPtr interface{}) error {
	val := reflect.ValueOf(eventPtr)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// Check if the field is zero-valued (unpopulated)
		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			return errors.New(fmt.Sprintf("field %s is empty", val.Type().Field(i).Name))

		}
	}
	return nil
}
