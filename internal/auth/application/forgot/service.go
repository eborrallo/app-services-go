package forgot

import (
	"app-services-go/internal/auth/domain"
	"context"

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
	token := crypt.Encrypt(userId)

	s.Sender.Send(userEmail, "Forgot password "+token)
	return nil
}
