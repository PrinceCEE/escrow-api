package push

import (
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

const (
	ErrSendingEmailMsg = "error sending email"
)

type Email struct {
	From    string
	To      []string
	Subject string
	Text    string
	Html    string
}

func SendEmail(data *Email) error {
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	e := email.NewEmail()
	if data.From != "" {
		e.From = data.From
	} else {
		e.From = username
	}

	e.To = data.To
	e.Subject = data.Subject
	e.Text = []byte(data.Text)
	e.HTML = []byte(data.Html)

	err := e.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth(
			"",
			username,
			password,
			"smtp.gmail.com",
		),
	)

	if err != nil {
		return err
	}

	return nil
}
