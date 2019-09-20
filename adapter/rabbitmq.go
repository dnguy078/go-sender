package adapter

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type RabbitClient struct {
	channel *amqp.Channel
	conn    *amqp.Connection
}

func newRabbitConnection(user, pass, addr string) (*RabbitClient, error) {
	dialConfig := amqp.Config{
		Dial: amqp.DefaultDial(1 * time.Minute),
	}
	conn, err := amqp.DialConfig("amqp://guest:guest@rabbitmq:5672/", dialConfig)
	if err != nil {
		return nil, fmt.Errorf("connection.open: %s", err)
	}

	c, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel.open: %s", err)
	}
	return &RabbitClient{
		conn:    conn,
		channel: c,
	}, nil
}

type RabbitClientPublisher struct {
	client *RabbitClient
}

type NewRabbitPublisher(user, pass, addr string) (*RabbitClientPublisher, error) {
	client, 


}

