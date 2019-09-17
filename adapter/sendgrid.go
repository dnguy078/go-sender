package adapter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridClient is a wrapper around SendGrid API
type SendGridClient struct {
	HTTPClient *http.Client
	user       string
	password   string

	mu         sync.Mutex
	errorCount int
}

type SendGridSimpleRequest struct {
	Personalizations []struct {
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
		Subject string `json:"subject"`
	} `json:"personalizations"`
	From struct {
		Email string `json:"email"`
	} `json:"from"`
	Content []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"content"`
}

// NewSendGridClient return sa new SendGridClient
func NewSendGridClient(user, password string) *SendGridClient {
	return &SendGridClient{
		HTTPClient: http.DefaultClient,
		user:       user,
		password:   password,
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
	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with Twilio SendGrid is Fun"
	to := mail.NewEmail("Example User", "dnguy078@ucr.edu")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient(os.Getenv("api key"))
	client := sendgrid.NewSendClient("sdjflsdkfjfdsf")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}
	if response.StatusCode != 202 {
		fmt.Println("i have failed")
		return errors.New("i have failed")
	}
	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
	return nil
}
