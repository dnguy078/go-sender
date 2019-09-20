// +build integration

package main

import (
	"flag"
	"testing"

	"github.com/dnguy078/go-sender/adapter"
	"github.com/dnguy078/go-sender/request"
)

var (
	sgAPIKey = flag.String("sgkey", "", "sendgrid api key")
	spAPIKey = flag.String("spkey", "", "sparkpost api key")
	// todo test connections with rabbitmq
)

func TestSendGridClient(t *testing.T) {
	tests := []struct {
		name      string
		expectErr error
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sgClient := adapter.NewSendGridClient(*sgAPIKey)

			err := sgClient.Email(request.EmailRequest{
				ToEmail:   "devnull@test.com",
				FromEmail: "devnull@test.com",
				Subject:   "TESTING SENDGRID INTEGRATION",
				Text:      "This is a test",
			})

			if (err != nil) && err != tt.expectErr {
				t.Errorf("SendGridClient.Email() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// Unable to test sparkpost, I would need to whitelist a domain
// and their "free" only lets me send from whitelist domain.
// They allow a sandbox feature but limited to 50 calls for entire account.
