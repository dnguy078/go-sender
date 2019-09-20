package adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dnguy078/go-sender/request"
)

func TestSparkPostClient_Email(t *testing.T) {
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
			name:           "failed",
			fakeRespStatus: http.StatusInternalServerError,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.fakeRespStatus)
				w.Write([]byte(`{"results": { "total_rejected_recipients": 0, "total_accepted_recipients": 1, "id": "697707383584862550"}}`))
			}))
			defer ts.Close()

			spClient, err := NewSparkPostClient("somekey")
			if err != nil {
				t.Error(err)
			}
			spClient.client.Config.BaseUrl = ts.URL

			req := request.EmailRequest{
				ToEmail:   "tosomeonespecial",
				Subject:   "love",
				FromEmail: "yoursecretadmirer",
				Text:      "xoxo!",
			}

			if err := spClient.Email(req); (err != nil) != tt.wantErr {
				t.Errorf("SparkPostClient.Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
