package smtp

import (
	"app-services-go/configs"
	"fmt"
	"net/smtp"
)

func Send(email string, content string) {

	config, error := configs.GetSmtpConfig()
	if error != nil {
		fmt.Println("Error creating smpt config", error)
		return
	}

	// Receiver email address.
	to := []string{
		"sender@example.com",
	}
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
