package adapter

import (
	"errors"
	"fmt"
	"net/http"
)

// SendGridClient is a wrapper around SendGrid API
type SendGridClient struct {
	HTTPClient *http.Client
}

// NewSendGridClient return sa new SendGridClient
func NewSendGridClient() *SendGridClient {
	return &SendGridClient{
		HTTPClient: http.DefaultClient,
	}
}

func (sg *SendGridClient) Name() string {
	return "sendgrid"
}

// Email performs a http request to send emails through SG
func (sgClient *SendGridClient) Email() error {
	fmt.Println("calliing from sg")
	return errors.New("something")
}
