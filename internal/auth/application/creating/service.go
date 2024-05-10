package creating

import (
	user "app-services-go/internal/auth/domain"
	"context"

	"app-services-go/kit/event"
)

// UserService is the default UserService interface
// implementation returned by creating.NewUserService.
type UserService struct {
	UserRepository user.UserRepository
	eventBus       event.Bus
}

// NewUserService returns the default Service interface implementation.
func NewUserService(UserRepository user.UserRepository, eventBus event.Bus) UserService {
	return UserService{
		UserRepository: UserRepository,
		eventBus:       eventBus,
	}
}

// CreateUser implements the creating.UserService interface.
func (s UserService) CreateUser(ctx context.Context, id, name, email, password string) error {

	User, err := user.NewUser(id, name, email, password)
	if err != nil {
		return err
	}

	if err := s.UserRepository.Save(ctx, User); err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, User.PullEvents())
}
