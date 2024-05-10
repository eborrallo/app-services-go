package domain

import (
	"context"
	"errors"
	"fmt"

	"app-services-go/internal/auth/domain/events"
	"app-services-go/kit/crypt"
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
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"duration"`
	Validated    bool   `json:"validated"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	events []event.Event
}

// UserRepository defines the expected behaviour from a User storage.
type UserRepository interface {
	Save(ctx context.Context, User User) error
	Update(ctx context.Context, User User) error
	FetchById(ctx context.Context, id string) (User, error)
	FetchByEmail(ctx context.Context, email string) (User, error)
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../infrastructure/storage/storagemocks --name=UserRepository

// NewUser creates a new User.
func NewUser(id, name, email, password string) (User, error) {
	idVO, err := NewUserID(id)
	if err != nil {
		return User{}, err
	}

	encodePassword := crypt.Md5(password)

	User := User{
		ID:        idVO,
		Name:      name,
		Email:     email,
		Password:  encodePassword,
		Validated: false,
	}
	User.Record(events.NewUserCreatedEvent(idVO, name, email, password))
	return User, nil
}

// Record records a new domain event.
func (c *User) Record(evt event.Event) {
	c.events = append(c.events, evt)
}

// Record records a new domain event.
func (c *User) Validate() {
	c.Validated = true
	c.Record(events.NewUserValidatedEvent(c.ID, c.Name, c.Email))

}

// Record records a new domain event.
func (c *User) UpdatePassword(newPassword string) {
	encodePassword := crypt.Md5(newPassword)
	c.Password = encodePassword
	c.Record(events.NewUserPasswordChangedEvent(c.ID, c.Email, encodePassword))
}

// Record records a new domain event.
func (c *User) Login(token, refreshToken string) {
	c.AccessToken = token
	c.RefreshToken = refreshToken
	c.Record(events.NewUserLoggedEvent(c.ID, c.Name, c.Email, c.AccessToken, c.RefreshToken))
}

// PullEvents returns all the recorded domain events.
func (c User) PullEvents() []event.Event {
	evt := c.events
	c.events = []event.Event{}

	return evt
}
