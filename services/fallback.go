package services

import (
	"encoding/json"
	"log"

	"github.com/dnguy078/go-sender/request"
)

type FallBack struct {
	Publisher publisher
}

type publisher interface {
	Publish(exchange, routingKey string, body []byte, headers map[string]interface{}) error
}

func (f *FallBack) PrimaryFallBack(email request.EmailRequest) {
	log.Println("calling fallback primary")
	b, err := json.Marshal(email)
	if err != nil {
		log.Println(err)
	}
	f.Publisher.Publish("emailer.incomingX", "retry", b, nil)
}

func (f *FallBack) SecondaryFallBack(email request.EmailRequest) {
	b, err := json.Marshal(email)
	if err != nil {
		log.Println(err)
	}
	f.Publisher.Publish("emailer.incomingX", "errors", b, nil)
}
