package sending_email_validation

import (
	"app-services-go/internal/auth/domain"
	"app-services-go/kit/crypt"
	"log"
)

type EmailValidatorSenderService struct {
	Sender domain.Sender
}

func NewEmailValidatorSenderService(Sender domain.Sender) EmailValidatorSenderService {
	return EmailValidatorSenderService{Sender: Sender}
}

func (s EmailValidatorSenderService) Send(email string, userId string) error {
	log.Println("Sending email validation", email)
	token := crypt.Encrypt(userId)

	s.Sender.Send(email, "Email validation "+token)
	return nil
}
