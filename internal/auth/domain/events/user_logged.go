package events

import (
	"app-services-go/kit/event"
	"encoding/json"
)

const UserLoggedEventType event.Type = "events.user.logged"

type UserLoggedEvent struct {
	event.BaseEvent
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"regresh_token"`
}

func NewUserLoggedEvent(id, name, email, access_token, refresh_token string) UserLoggedEvent {
	return UserLoggedEvent{
		Id:           id,
		Name:         name,
		Email:        email,
		AccessToken:  email,
		RefreshToken: refresh_token,

		BaseEvent: event.NewBaseEvent(id),
	}
}
func (e UserLoggedEvent) ID() string {
	return e.EventID
}
func (e UserLoggedEvent) Type() event.Type {
	return UserLoggedEventType
}

func (e UserLoggedEvent) UserID() string {
	return e.Id
}

func (e UserLoggedEvent) UserName() string {
	return e.Name
}

func (e UserLoggedEvent) UserEmail() string {
	return e.Email
}
func (e UserLoggedEvent) FromJSON(data []byte) (event.Event, error) {
	user := &UserLoggedEvent{}
	err := json.Unmarshal(data, user)
	if err != nil {
		return UserLoggedEvent{}, err
	}
	return *user, nil
}
