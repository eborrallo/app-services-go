package validating

import (
	user "app-services-go/internal/auth/domain"
	"context"

	"app-services-go/kit/event"
)

// ValidateUserService is the default ValidateUserService interface
// implementation returned by creating.NewValidateUserService.
type ValidateUserService struct {
	UserRepository user.UserRepository
	eventBus       event.Bus
}

// NewValidateUserService returns the default Service interface implementation.
func NewValidateUserService(UserRepository user.UserRepository, eventBus event.Bus) ValidateUserService {
	return ValidateUserService{
		UserRepository: UserRepository,
		eventBus:       eventBus,
	}
}

// ValidateUser implements the creating.ValidateUserService interface.
func (s ValidateUserService) ValidateUser(ctx context.Context, userId string) error {

	User, err := s.UserRepository.FetchById(ctx, userId)
	if err != nil {
		return err
	}
	User.Validate()

	if err := s.UserRepository.Update(ctx, User); err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, User.PullEvents())
}
