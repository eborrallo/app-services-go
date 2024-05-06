package sending_email_validation

import (
	smtp "app-services-go/kit/email/smpt"
	"log"
)

type EmailValidatorSenderService struct{}

func NewEmailValidatorSenderService() EmailValidatorSenderService {
	return EmailValidatorSenderService{}
}

func (s EmailValidatorSenderService) Send(email string) error {
	log.Println("Sending email validation", email)
	smtp.Send(email, "Email validation")
	return nil
}
