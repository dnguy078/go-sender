package adapter

import (
	"fmt"
	"net/http"

	"github.com/dnguy078/go-sender/request"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridClient is a wrapper around SendGrid API
type SendGridClient struct {
	client *sendgrid.Client
}

// NewSendGridClient return sa new SendGridClient
func NewSendGridClient(apiKey string) *SendGridClient {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridClient{
		client: client,
	}
}

func (sg *SendGridClient) Type() string {
	return "sendgrid"
}

// Email performs a http request to send emails through SG
func (sgClient *SendGridClient) Email(payload request.EmailRequest) error {
	mail := buildMessage(payload)
	response, err := sgClient.client.Send(mail)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("SendGrid service down error: %s", response.Body)
	}

	return nil
}

func buildMessage(payload request.EmailRequest) *mail.SGMailV3 {
	from := mail.NewEmail(payload.FromEmail, payload.FromEmail)
	subject := payload.Subject
	to := mail.NewEmail(payload.ToEmail, payload.ToEmail)
	plainTextContent := payload.Text
	//not sure this is needed
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"

	return mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
}
