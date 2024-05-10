package events

import (
	"app-services-go/kit/event"
	"encoding/json"
)

const UserValidatedEventType event.Type = "events.user.valdiated"

type UserValidatedEvent struct {
	event.BaseEvent
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserValidatedEvent(id, name, email string) UserValidatedEvent {
	return UserValidatedEvent{
		Id:    id,
		Name:  name,
		Email: email,

		BaseEvent: event.NewBaseEvent(id),
	}
}
func (e UserValidatedEvent) ID() string {
	return e.EventID
}
func (e UserValidatedEvent) Type() event.Type {
	return UserValidatedEventType
}

func (e UserValidatedEvent) UserID() string {
	return e.Id
}

func (e UserValidatedEvent) UserName() string {
	return e.Name
}

func (e UserValidatedEvent) UserEmail() string {
	return e.Email
}
func (e UserValidatedEvent) FromJSON(data []byte) (event.Event, error) {
	user := &UserValidatedEvent{}
	err := json.Unmarshal(data, user)
	if err != nil {
		return UserValidatedEvent{}, err
	}
	return *user, nil
}
