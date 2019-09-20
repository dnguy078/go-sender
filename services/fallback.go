package services

import (
	"fmt"

	"github.com/dnguy078/go-sender/request"
	"github.com/streadway/amqp"
)

type FallBack struct {
	ch       amqp.Channel
	exchange string
	queue    string
	headers  amqp.Table
}

func (f *FallBack) PrimaryFallBack(email request.EmailRequest) {
	// msg := amqp.Publishing{
	// 	DeliveryMode: amqp.Persistent,
	// 	Timestamp:    time.Now(),
	// 	ContentType:  "text/plain",
	// 	Body:         []byte("Go Go AMQP!"),
	// }
	// f.ch.Publish()
	fmt.Println("klsdjflksdjf")
}

func (f *FallBack) SecondaryFallBack(email request.EmailRequest) {

}
