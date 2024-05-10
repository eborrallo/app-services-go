package events

import (
	"app-services-go/kit/event"
	"encoding/json"
)

const UserPasswordChangedEventType event.Type = "events.user.password.changed"

type UserPasswordChangedEvent struct {
	event.BaseEvent
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserPasswordChangedEvent(id, email, password string) UserPasswordChangedEvent {
	return UserPasswordChangedEvent{
		Id:       id,
		Email:    email,
		Password: password,

		BaseEvent: event.NewBaseEvent(id),
	}
}
func (e UserPasswordChangedEvent) ID() string {
	return e.EventID
}
func (e UserPasswordChangedEvent) Type() event.Type {
	return UserPasswordChangedEventType
}

func (e UserPasswordChangedEvent) UserID() string {
	return e.Id
}

func (e UserPasswordChangedEvent) UserEmail() string {
	return e.Email
}

func (e UserPasswordChangedEvent) UserPassword() string {
	return e.Password
}
func (e UserPasswordChangedEvent) FromJSON(data []byte) (event.Event, error) {
	user := &UserPasswordChangedEvent{}
	err := json.Unmarshal(data, user)
	if err != nil {
		return UserPasswordChangedEvent{}, err
	}
	return *user, nil
}
