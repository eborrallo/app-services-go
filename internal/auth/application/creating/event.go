package creating

import (
	"app-services-go/internal/auth/application/sending_email_validation"
	user "app-services-go/internal/auth/domain"

	"app-services-go/kit/event"
	"errors"
)

type SendEmailVerificationOnUserCreated struct {
	sendVerificationEmailService sending_email_validation.EmailValidatorSenderService
}

func NewSendEmailVerificationOnUserCreated(service sending_email_validation.EmailValidatorSenderService) SendEmailVerificationOnUserCreated {
	return SendEmailVerificationOnUserCreated{
		sendVerificationEmailService: service,
	}
}

func (e SendEmailVerificationOnUserCreated) On(evt event.Event) error {
	userCreatedEvt, ok := evt.(user.UserCreatedEvent)
	if !ok {
		return errors.New("unexpected event")
	}
	return e.sendVerificationEmailService.Send(userCreatedEvt.Email)
}

func (e SendEmailVerificationOnUserCreated) SubscribedTo() event.Event {
	return event.Event(user.UserCreatedEvent{})
}
