package services

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/dnguy078/go-sender/request"
	"github.com/streadway/amqp"
)

// this is used to mock
type fakeEmailer struct {
	resp error
}

func (fe *fakeEmailer) Email(req request.EmailRequest) error {
	return fe.resp
}

func (fe *fakeEmailer) Type() string {
	return "fakeemailer"
}

type fakeFailOver struct {
	calledCount int
}

func (fo *fakeFailOver) FallBack(req request.EmailRequest) {
	fo.calledCount++
}

// used to test delivery handlers
type fakeAcknowledger struct{}

func (fa *fakeAcknowledger) Ack(tag uint64, multiple bool) error {
	return nil
}
func (fa *fakeAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	return nil
}
func (fa *fakeAcknowledger) Reject(tag uint64, requeue bool) error {
	return nil
}

func TestEmailWorker_Work(t *testing.T) {
	type fields struct {
		emailer Emailer
	}
	tests := []struct {
		name                   string
		fields                 fields
		expectedFailoverCalled int
	}{
		{
			name: "success",
			fields: fields{
				emailer: &fakeEmailer{},
			},
			expectedFailoverCalled: 0,
		},
		{
			name: "failed",
			fields: fields{
				emailer: &fakeEmailer{resp: errors.New("email service is down")},
			},
			expectedFailoverCalled: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fo := &fakeFailOver{}
			w := &EmailWorker{
				emailer:      tt.fields.emailer,
				queue:        make(chan amqp.Delivery, 2),
				quit:         make(chan bool, 1),
				fallBackFunc: fo.FallBack,
			}
			b, err := json.Marshal(request.EmailRequest{ToEmail: "someemail"})
			if err != nil {
				t.Error(t)
			}
			w.queue <- amqp.Delivery{Body: b, Acknowledger: &fakeAcknowledger{}}
			go w.Work()

			time.Sleep(50 * time.Millisecond)
			close(w.quit)
			if fo.calledCount != tt.expectedFailoverCalled {
				t.Errorf("%s - got %d, expected %d", tt.name, fo.calledCount, tt.expectedFailoverCalled)
			}
		})
	}
}
