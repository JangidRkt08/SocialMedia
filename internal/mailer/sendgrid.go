package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	// "github.com/sendgrid/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apikey    string
	Client    *sendgrid.Client
}

func NewSendgrid(apikey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apikey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apikey:    apikey,
		Client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{

			Enable: &isSandbox,
		},
	})

	var retryErr error
	for i := 0; i < maxRetries; i++ {
		response, retryErr := m.Client.Send(message)

		if retryErr != nil {

			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			log.Printf(
				"SendGrid rejected email. Status: %d",
				response.StatusCode,
			)
			log.Printf("SendGrid response: %s", response.Body)

			return -1, fmt.Errorf(
				"sendgrid returned status %d: %s",
				response.StatusCode,
				response.Body,
			)
		}
		return response.StatusCode, nil
	}
	return -1, fmt.Errorf("Failed to send email after %d attempts, error: %v", maxRetries, retryErr)
}
