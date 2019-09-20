package adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dnguy078/go-sender/request"

	sendgrid "github.com/sendgrid/sendgrid-go"
)

func TestSendGridClient_Email(t *testing.T) {
	tests := []struct {
		name           string
		fakeRespStatus int
		fakeRespBody   []byte
		wantErr        bool
	}{
		{
			name:           "success",
			fakeRespStatus: http.StatusAccepted,
		},
		{
			name:           "bad request",
			fakeRespStatus: http.StatusBadRequest,
		},
		{
			name:           "sendgrid down",
			fakeRespStatus: http.StatusInternalServerError,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.fakeRespStatus)
			}))
			defer ts.Close()

			sgClient := &SendGridClient{
				client: sendgrid.NewSendClient("testclient"),
			}
			sgClient.client.BaseURL = ts.URL

			if err := sgClient.Email(request.EmailRequest{
				ToEmail: "test",
			}); (err != nil) != tt.wantErr {
				t.Errorf("SendGridClient.Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
