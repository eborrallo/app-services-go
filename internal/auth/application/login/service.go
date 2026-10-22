package login

import (
	"context"

	"app-services-go/internal/auth/domain"
	"app-services-go/kit/crypt"
	"app-services-go/kit/event"
)

// LoginUserService is the default LoginUserService interface
// implementation returned by creating.NewLoginUserService.
type LoginUserService struct {
	UserRepository domain.UserRepository
	eventBus       event.Bus
}

// NewLoginUserService returns the default Service interface implementation.
func NewLoginUserService(UserRepository domain.UserRepository, eventBus event.Bus) LoginUserService {
	return LoginUserService{
		UserRepository: UserRepository,
		eventBus:       eventBus,
	}
}

// LoginUser implements the creating.LoginUserService interface.
func (s LoginUserService) LoginUser(ctx context.Context, token, refreshToken string) error {
	payload := domain.TokenPayload{}
	if err := crypt.GetPayloadFromToken(token, &payload); err != nil {
		return err
	}

	User, err := s.UserRepository.FetchById(ctx, payload.ID)
	if err != nil {
		return err
	}
	User.Login(token, refreshToken)

	if err := s.UserRepository.Update(ctx, User); err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, User.PullEvents())
}
