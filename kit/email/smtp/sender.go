package smtp

import (
	"app-services-go/configs"
	"fmt"
	"net/smtp"
)

type Sender struct {
}

// NewSender returns the default Service interface implementation.
func NewSender() Sender {
	return Sender{}
}

func (s Sender) Send(email string, content string) error {
	config, err := configs.GetSmtpConfig()
	if err != nil {
		return fmt.Errorf("error creating smtp config: %w", err)
	}

	// Receiver email address.
	to := []string{email}
	// Build a minimal RFC 5322 message so servers don't reject or mangle it.
	message := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s\r\n",
		config.Email, email, "Notification", content,
	))
	// Authentication.
	auth := smtp.PlainAuth("", config.Email, config.Password, config.Host)
	// Sending email.
	if err := smtp.SendMail(config.Host+":"+config.Port, auth, config.Email, to, message); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
