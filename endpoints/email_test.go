package endpoints

import (
	"net/http"
	"testing"
)

func TestEmailHandler_Email(t *testing.T) {
	type fields struct {
		d dispatcher
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EmailHandler{
				d: tt.fields.d,
			}
			e.Email(tt.args.w, tt.args.r)
		})
	}
}
