package domain

import (
	"context"
	"errors"
	"fmt"

	"app-services-go/kit/event"

	"github.com/google/uuid"
)

var ErrInvalidUserID = errors.New("invalid User ID")

// UserID represents the User unique identifier.

// NewUserID instantiate the VO for UserID
func NewUserID(value string) (string, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrInvalidUserID, value)
	}

	return v.String(), nil

}

// User is the data structure that represents a User.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"duration"`

	events []event.Event
}

// UserRepository defines the expected behaviour from a User storage.
type UserRepository interface {
	Save(ctx context.Context, User User) error
	FetchById(ctx context.Context, id string) (User, error)
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../infrastructure/storage/storagemocks --name=UserRepository

// NewUser creates a new User.
func NewUser(id, name, email, password string) (User, error) {
	idVO, err := NewUserID(id)
	if err != nil {
		return User{}, err
	}

	User := User{
		ID:       idVO,
		Name:     name,
		Email:    email,
		Password: password,
	}
	User.Record(NewUserCreatedEvent(idVO, name, email, password))
	return User, nil
}

// Record records a new domain event.
func (c *User) Record(evt event.Event) {
	c.events = append(c.events, evt)
}

// PullEvents returns all the recorded domain events.
func (c User) PullEvents() []event.Event {
	evt := c.events
	c.events = []event.Event{}

	return evt
}
