package adapter

import (
	"testing"

	"github.com/sendgrid/sendgrid-go"
)

func TestSendGridClient_Email(t *testing.T) {
	type fields struct {
		client *sendgrid.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sgClient := &SendGridClient{
				client: tt.fields.client,
			}
			if err := sgClient.Email(); (err != nil) != tt.wantErr {
				t.Errorf("SendGridClient.Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
