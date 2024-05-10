package remember

import (
	"app-services-go/internal/auth/domain"
	"app-services-go/kit/event"
	"context"
)

// PasswordUserService is the default PasswordUserService interface
// implementation returned by creating.NewPasswordUserService.
type PasswordUserService struct {
	UserRepository domain.UserRepository
	eventBus       event.Bus
}

// NewPasswordUserService returns the default Service interface implementation.
func NewPasswordUserService(UserRepository domain.UserRepository, eventBus event.Bus) PasswordUserService {
	return PasswordUserService{
		UserRepository: UserRepository,
		eventBus:       eventBus,
	}
}

// PasswordUser implements the creating.PasswordUserService interface.
func (s PasswordUserService) UpdatePassword(ctx context.Context, userId, newPassword string) error {

	User, err := s.UserRepository.FetchById(ctx, userId)
	if err != nil {
		return err
	}
	User.UpdatePassword(newPassword)

	if err := s.UserRepository.Update(ctx, User); err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, User.PullEvents())

}
