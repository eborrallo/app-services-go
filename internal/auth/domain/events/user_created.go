package events

import (
	"app-services-go/kit/event"
	"encoding/json"
)

const UserCreatedEventType event.Type = "events.user.created"

type UserCreatedEvent struct {
	event.BaseEvent
	Id      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func NewUserCreatedEvent(id string, name string, email string, address string) UserCreatedEvent {
	return UserCreatedEvent{
		Id:      id,
		Name:    name,
		Email:   email,
		Address: address,

		BaseEvent: event.NewBaseEvent(id),
	}
}
func (e UserCreatedEvent) ID() string {
	return e.EventID
}
func (e UserCreatedEvent) Type() event.Type {
	return UserCreatedEventType
}

func (e UserCreatedEvent) UserID() string {
	return e.Id
}

func (e UserCreatedEvent) UserName() string {
	return e.Name
}

func (e UserCreatedEvent) UserEmail() string {
	return e.Email
}
func (e UserCreatedEvent) FromJSON(data []byte) (event.Event, error) {
	user := &UserCreatedEvent{}
	err := json.Unmarshal(data, user)
	if err != nil {
		return UserCreatedEvent{}, err
	}
	return *user, nil
}
