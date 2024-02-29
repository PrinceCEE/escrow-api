package push

import (
	"net/smtp"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/jordan-wright/email"
	"github.com/rs/zerolog"
)

type Email struct {
	From    string
	To      []string
	Subject string
	Text    string
	Html    string
}

func SendEmail(data *Email) {
	username := config.Config.Env.EMAIL_USERNAME
	password := config.Config.Env.EMAIL_PASSWORD

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
		config.Config.Logger.Log(zerolog.InfoLevel, "error sending email", nil, err)
	}
}
