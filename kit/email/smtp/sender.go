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

func (s Sender) Send(email string, content string) {

	config, error := configs.GetSmtpConfig()
	if error != nil {
		fmt.Println("Error creating smpt config", error)
		return
	}

	// Receiver email address.
	to := []string{email}
	// Message.
	message := []byte(content)
	// Authentication.
	auth := smtp.PlainAuth("", config.Email, config.Password, config.Host)
	// Sending email.
	err := smtp.SendMail(config.Host+":"+config.Port, auth, config.Email, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
