package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dnguy078/go-sender/request"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridClient is a wrapper around SendGrid API
type SendGridClient struct {
	client  *http.Client
	apiKey  string
	mailURL string
}

// NewSendGridClient return sa new SendGridClient
func NewSendGridClient(apiKey string) *SendGridClient {
	return &SendGridClient{
		client:  http.DefaultClient,
		mailURL: "https://api.sendgrid.com/v3/mail/send",
		apiKey:  fmt.Sprintf("Bearer %s", apiKey),
	}
}

func (sg *SendGridClient) Type() string {
	return "sendgrid"
}

// Email performs a http request to send emails through SG
func (sgClient *SendGridClient) Email(payload request.EmailRequest) error {
	mail := buildMessage(payload)
	b, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", sgClient.mailURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", sgClient.apiKey)
	req.Header.Set("content-type", "application/json")

	resp, err := sgClient.client.Do(req)
	if err != nil {
		return err
	}

	// return errors.New("sendgrid is down")

	if resp.StatusCode >= http.StatusInternalServerError {
		return errors.New("sendgrid is down")
	}

	if resp.StatusCode <= http.StatusOK || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Printf("sendgrid errors: %s, status: %d", string(body), resp.StatusCode)
	}

	return nil
}

func buildMessage(payload request.EmailRequest) *mail.SGMailV3 {
	m := mail.NewV3Mail()

	from := mail.NewEmail("this is a test", payload.FromEmail)
	m.SetFrom(from)

	to := mail.NewEmail("to test user", payload.ToEmail)
	m.Subject = payload.Subject
	p := mail.NewPersonalization()
	p.AddTos(to)
	m.AddPersonalizations(p)

	content := mail.NewContent("text/plain", payload.Text)
	m.AddContent(content)

	return m
}
