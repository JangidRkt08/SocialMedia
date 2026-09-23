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

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{

			Enable: &isSandbox,
		},
	})

	for i := 0; i < maxRetries; i++ {
		response, err := m.Client.Send(message)

		if err != nil {
			log.Printf("Failed to send emails to %v, atteempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

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

			return fmt.Errorf(
				"sendgrid returned status %d: %s",
				response.StatusCode,
				response.Body,
			)
		}
		log.Printf("Email sent with status code %v", response.StatusCode)
		return nil
	}
	return fmt.Errorf("Failed to send email after %d attempts", maxRetries)
}
