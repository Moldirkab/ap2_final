package email

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendWelcomeEmail(toEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Welcome!")
	m.SetBody("text/html", "<h1>Welcome to our car rental service!</h1>")

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	return d.DialAndSend(m)
}
