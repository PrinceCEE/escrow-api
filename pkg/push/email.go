package push

import (
	"fmt"
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
	from := os.Getenv("EMAIL_FROM")
	port := os.Getenv("EMAIL_PORT")
	host := os.Getenv("EMAIL_HOST")

	e := email.NewEmail()
	if data.From != "" {
		e.From = data.From
	} else {
		e.From = from
	}

	e.To = data.To
	e.Subject = data.Subject
	e.Text = []byte(data.Text)
	e.HTML = []byte(data.Html)

	err := e.Send(
		fmt.Sprintf("%s:%s", host, port),
		smtp.PlainAuth(
			"",
			username,
			password,
			host,
		),
	)

	if err != nil {
		return err
	}

	return nil
}
