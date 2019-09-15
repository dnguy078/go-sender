package adapter

import (
	"fmt"
	"net/http"
)

// MailgunClient is a wrapper around Mailgun's API
type MailgunClient struct {
	HTTPClient *http.Client
}

// NewMailgunClient returns a new MailgunClient
func NewMailgunClient() *MailgunClient {
	return &MailgunClient{
		HTTPClient: http.DefaultClient,
	}
}

func (mg *MailgunClient) Name() string {
	return "mailgun"
}

// Email performs a http request to send emails through MG
func (mg *MailgunClient) Email() error {
	fmt.Println("calliing from mgjsdlkfjsdlkfjsdlkfjsdlkfjsdlkfj ")
	return nil
}
