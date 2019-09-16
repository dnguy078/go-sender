package adapter

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

// SendGridClient is a wrapper around SendGrid API
type SendGridClient struct {
	HTTPClient *http.Client

	mu         sync.Mutex
	errorCount int
}

// NewSendGridClient return sa new SendGridClient
func NewSendGridClient() *SendGridClient {
	return &SendGridClient{
		HTTPClient: http.DefaultClient,
		errorCount: 0,
	}
}

func (sg *SendGridClient) Name() string {
	return "sendgrid"
}

func (sg *SendGridClient) IncrementError() int {
	sg.mu.Lock()
	defer sg.mu.Unlock()
	fmt.Println("errorcount", sg.errorCount)
	sg.errorCount++
	return sg.errorCount
}

// Email performs a http request to send emails through SG
func (sgClient *SendGridClient) Email() error {
	errorCount := sgClient.IncrementError()
	if errorCount == 5 {
		return nil
	}

	return errors.New("something")
}
