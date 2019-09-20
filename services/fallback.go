package services

import (
	"encoding/json"

	"github.com/dnguy078/go-sender/request"
)

type FallBack struct {
	Publisher publisher
}

type publisher interface {
	Publish(exchange, routingKey string, body []byte, headers map[string]interface{}) error
}

func (f *FallBack) PrimaryFallBack(email request.EmailRequest) {
	b, err := generateRabbitPayload(email)
	if err != nil {
		return
	}
	f.Publisher.Publish("emailer.incomingX", "retry", b, nil)
}

func (f *FallBack) SecondaryFallBack(email request.EmailRequest) {
	b, err := generateRabbitPayload(email)
	if err != nil {
		return
	}
	f.Publisher.Publish("emailer.incomingX", "errors", b, nil)
}

func generateRabbitPayload(email request.EmailRequest) ([]byte, error) {
	b, err := json.Marshal(email)
	if err != nil {
		return nil, err
	}
	return b, nil
}
