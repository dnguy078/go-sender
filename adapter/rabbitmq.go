package adapter

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func newRabbitConnection(user, pass, addr string) (*amqp.Connection, error) {
	dialConfig := amqp.Config{
		Dial: amqp.DefaultDial(1 * time.Minute),
	}
	conn, err := amqp.DialConfig("amqp://guest:guest@rabbitmq:5672/", dialConfig)
	if err != nil {
		return nil, fmt.Errorf("connection.open: %s", err)
	}
	return conn, nil
}

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitClient(user, pass, addr string) (*RabbitClient, error) {
	conn, err := newRabbitConnection(user, pass, addr)
	if err != nil {
		return nil, err
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

func (rc *RabbitClient) Publish(exchange, routingKey string, body []byte) error {
	if err := rc.channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}
	return nil
}

func (rc *RabbitClient) Consume() (<-chan amqp.Delivery, error) {
	msgChan, err := rc.channel.Consume("emailer.incoming.queue", "emailer.incoming.queue", false, false, false, false, nil)
	if err != nil {
		return msgChan, fmt.Errorf("basic.consume: %v", err)
	}
	return msgChan, nil
}
