package sending_email_validation

import (
	"log"
)

type EmailValidatorSenderService struct{}

func NewEmailValidatorSenderService() EmailValidatorSenderService {
	return EmailValidatorSenderService{}
}

func (s EmailValidatorSenderService) Send(email string) error {
	log.Println("Sending email validation", email)
	return nil
}
