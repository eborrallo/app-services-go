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

// ErrUserNotFound is returned by the repository when no user matches the query,
// so callers can tell "does not exist" (404) from an infrastructure error (500).
var ErrUserNotFound = errors.New("user not found")

// NewUserID instantiate the VO for UserID
func NewUserID(value string) (string, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrInvalidUserID, value)
	}

	return v.String(), nil

}

// TokenPayload is the minimal set of claims embedded in a JWT. Never put the
// whole User (which carries the password hash and tokens) into a token.
type TokenPayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// User is the data structure that represents a User.
type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	Address      string `json:"address"`
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
	FetchByAddress(ctx context.Context, address string) (User, error)
}
type UserMessageRepository interface {
	SaveMessage(address string, message string) error
	GetMessage(address string) (string, error)
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../infrastructure/storage/storagemocks --name=UserRepository
//go:generate mockery --case=snake --outpkg=storagemocks --output=../infrastructure/storage/storagemocks --name=UserMessageRepository

// NewUser creates a new User.
func NewUser(id string, name string, email string, password string) (User, error) {
	idVO, err := NewUserID(id)
	if err != nil {
		return User{}, err
	}

	encodePassword, err := crypt.Encode(password)
	if err != nil {
		return User{}, err
	}

	User := User{
		ID:        idVO,
		Name:      name,
		Email:     email,
		Password:  encodePassword,
		Validated: false,
	}
	User.Record(events.NewUserCreatedEvent(idVO, name, email, ""))
	return User, nil
}
func NewWeb3User(id string, address string, token string, refreshToken string) (User, error) {
	idVO, err := NewUserID(id)
	if err != nil {
		return User{}, err
	}

	User := User{
		ID:           idVO,
		Address:      address,
		Validated:    true,
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	// Web3 users are created already validated and have no email, so no
	// UserCreated event is recorded (it would otherwise trigger a verification
	// email to an empty address).
	return User, nil
}

// Record records a new domain event.
func (c *User) Record(evt event.Event) {
	c.events = append(c.events, evt)
}

// Validate marks the user as validated and records the event.
func (c *User) Validate() {
	c.Validated = true
	c.Record(events.NewUserValidatedEvent(c.ID, c.Name, c.Email))

}

// UpdatePassword re-hashes and stores a new password and records the event.
func (c *User) UpdatePassword(newPassword string) error {
	encodePassword, err := crypt.Encode(newPassword)
	if err != nil {
		return err
	}
	c.Password = encodePassword
	c.Record(events.NewUserPasswordChangedEvent(c.ID, c.Email, encodePassword))
	return nil
}

// Login stores the issued tokens and records the event.
func (c *User) Login(token, refreshToken string) {
	c.AccessToken = token
	c.RefreshToken = refreshToken
	c.Record(events.NewUserLoggedEvent(c.ID, c.Name, c.Email, c.AccessToken, c.RefreshToken))
}

// PullEvents returns all the recorded domain events and clears them from the
// aggregate (pointer receiver, so the clearing actually sticks).
func (c *User) PullEvents() []event.Event {
	evt := c.events
	c.events = []event.Event{}

	return evt
}
