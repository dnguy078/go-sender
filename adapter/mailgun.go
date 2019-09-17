package adapter

import (
	"errors"
	"log"
	"net/http"
)

// MailgunClient is a wrapper around Mailgun's API
type MailgunClient struct {
	HTTPClient *http.Client
}

// curl -s --user 'api:YOUR_API_KEY' \
// https://api.mailgun.net/v3/YOUR_DOMAIN_NAME/messages \
// -F from='Excited User <mailgun@YOUR_DOMAIN_NAME>' \
// -F to=YOU@YOUR_DOMAIN_NAME \
// -F to=bar@example.com \
// -F subject='Hello' \
// -F text='Testing some Mailgun awesomeness!'
type MailGunSimpleRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

// NewMailgunClient returns a new MailgunClient
func NewMailgunClient() *MailgunClient {
	return &MailgunClient{
		HTTPClient: http.DefaultClient,
	}
}

func (mg *MailgunClient) Type() string {
	return "mailgun"
}

// Email performs a http request to send emails through MG
func (mg *MailgunClient) Email() error {
	log.Println("stubbing mail gun is down")
	return errors.New("mailgun is down")
}
