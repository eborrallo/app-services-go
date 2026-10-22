package sending_email_validation

import (
	"log"

	"app-services-go/internal/auth/domain"
	"app-services-go/kit/crypt"
)

type EmailValidatorSenderService struct {
	Sender domain.Sender
}

func NewEmailValidatorSenderService(Sender domain.Sender) EmailValidatorSenderService {
	return EmailValidatorSenderService{Sender: Sender}
}

func (s EmailValidatorSenderService) Send(email string, userId string) error {
	log.Println("Sending email validation", email)
	token, err := crypt.Encrypt(userId)
	if err != nil {
		return err
	}

	return s.Sender.Send(email, "Email validation "+token)
}
