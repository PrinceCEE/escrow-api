package email

import (
	"net/smtp"

	"github.com/Bupher-Co/bupher-api/config"
	emailClient "github.com/jordan-wright/email"
)

type Email struct {
	From    string
	To      []string
	Subject string
	Text    string
	Html    string
}

func SendEmail(data *Email) error {
	username := config.Config.Env.EMAIL_USERNAME
	password := config.Config.Env.EMAIL_PASSWORD

	e := emailClient.NewEmail()
	if data.From != "" {
		e.From = data.From
	} else {
		e.From = username
	}

	e.To = data.To
	e.Subject = data.Subject
	e.Text = []byte(data.Text)
	e.HTML = []byte(data.Html)

	return e.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth(
			"",
			username,
			password,
			"smtp.gmail.com",
		),
	)
}
