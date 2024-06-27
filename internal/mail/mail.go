package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(receiver string) {
	// Sender data.
	from := "alex71labonne@gmail.com"
	password := os.Getenv("SMTP_PASS")

	// Receiver email address.
	to := []string{
		receiver,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte(
		"To validate your account, please click <a href='http://localhost:8080/login/activate?token=hello'>THIS LINK</a>",
	)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
