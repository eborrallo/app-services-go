package forgot

import (
	"context"

	"app-services-go/internal/auth/domain"
	"app-services-go/kit/crypt"
)

// ForgotUserService is the default ForgotUserService interface
// implementation returned by creating.NewForgotUserService.
type ForgotUserService struct {
	Sender domain.Sender
}

// NewForgotUserService returns the default Service interface implementation.
func NewForgotUserService(sender domain.Sender) ForgotUserService {
	return ForgotUserService{Sender: sender}
}

// ForgotUser implements the creating.ForgotUserService interface.
func (s ForgotUserService) SendForgotEmail(ctx context.Context, userId, userEmail string) error {
	token, err := crypt.Encrypt(userId)
	if err != nil {
		return err
	}

	return s.Sender.Send(userEmail, "Forgot password "+token)
}
