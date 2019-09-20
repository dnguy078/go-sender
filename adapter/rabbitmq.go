package adapter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

func newRabbitConnection(user, pass, addr string, port int) (*amqp.Connection, error) {
	dialConfig := amqp.Config{
		Dial: amqp.DefaultDial(1 * time.Minute),
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pass, addr, port)
	conn, err := amqp.DialConfig(url, dialConfig)
	if err != nil {
		return nil, fmt.Errorf("connection.open: %s", err)
	}
	return conn, nil
}

// RabbitClient is a wrapper around rabbitmq
type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitClient returns a rabbitclient
func NewRabbitClient(user, pass, addr string, port int) (*RabbitClient, error) {
	conn, err := newRabbitConnection(user, pass, addr, port)
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

// Close cleanly closes out any connections
func (rc *RabbitClient) Close() {
	if err := rc.conn.Close(); err != nil {
		log.Println(err)
	}
	if err := rc.channel.Close(); err != nil {
		log.Println(err)
	}
}

// Publish sends rabbitmq message
func (rc *RabbitClient) Publish(exchange, routingKey string, body []byte, headers map[string]interface{}) error {
	if err := rc.channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         headers,
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

// Consume returns a channel of amqp.Delivery
func (rc *RabbitClient) Consume(queue, consumerName string) (<-chan amqp.Delivery, error) {
	msgChan, err := rc.channel.Consume(queue, consumerName, false, false, false, false, nil)
	if err != nil {
		return msgChan, fmt.Errorf("basic.consume: %v", err)
	}
	return msgChan, nil
}

// RabbitHealthcheck performs a status check on the rabbitmq management page,
// used to tell if rabbit is up before attempting to establish a amqp.Connection
func RabbitHealthcheck() error {
	req, err := http.NewRequest("GET", fmt.Sprint("http://rabbitmq:15672/api/overview"), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("guest", "guest")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("rabbitmq is down")
	}
	return nil
}
