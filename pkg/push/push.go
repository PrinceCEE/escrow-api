package push

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

type IPush interface {
	SendSMS(data *Sms)
	SendEmail(data *Email) error
}

type Push struct{}

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

func (p *Push) SendEmail(data *Email) error {
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

type Sms struct {
	Phone   string
	Message string
}

func (p *Push) SendSMS(data *Sms) {}
