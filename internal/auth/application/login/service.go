package login

import (
	"app-services-go/internal/auth/domain"
	user "app-services-go/internal/auth/domain"
	"context"

	"app-services-go/kit/crypt"
	"app-services-go/kit/event"
)

// LoginUserService is the default LoginUserService interface
// implementation returned by creating.NewLoginUserService.
type LoginUserService struct {
	UserRepository user.UserRepository
	eventBus       event.Bus
}

// NewLoginUserService returns the default Service interface implementation.
func NewLoginUserService(UserRepository user.UserRepository, eventBus event.Bus) LoginUserService {
	return LoginUserService{
		UserRepository: UserRepository,
		eventBus:       eventBus,
	}
}

// LoginUser implements the creating.LoginUserService interface.
func (s LoginUserService) LoginUser(ctx context.Context, token, refreshToken string) error {
	user := domain.User{}
	err := crypt.GetPayloadFromToken(token, &user)
	if err != nil {
		return err
	}

	User, err := s.UserRepository.FetchById(ctx, user.ID)
	if err != nil {
		return err
	}
	User.Login(token, refreshToken)

	if err := s.UserRepository.Update(ctx, User); err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, User.PullEvents())
}
