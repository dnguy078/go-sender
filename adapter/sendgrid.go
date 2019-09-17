package adapter

import (
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
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

func (sg *SendGridClient) Name() string {
	return "sendgrid"
}

// Email performs a http request to send emails through SG
func (sgClient *SendGridClient) Email() error {
	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with Twilio SendGrid is Fun"
	to := mail.NewEmail("Example User", "dnguy078@ucr.edu")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	response, err := sgClient.client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("SendGrid service down error: %s", response.Body)
	}

	return nil
}
