package daemon

import (
	"fmt"
	"time"

	"github.com/dnguy078/go-sender/adapter"
	"github.com/dnguy078/go-sender/config"
	"github.com/dnguy078/go-sender/services"
	"github.com/dnguy078/go-sender/utils"

	"github.com/streadway/amqp"
)

type Daemon struct {
	cfg config.Config
}

func New(cfg config.Config) (*Daemon, error) {
	rw := utils.NewRetryWrapper(10, utils.ExponentialBackoff)
	d := &Daemon{
		cfg: cfg,
	}

	if err := rw.WithRetry(d.connectRabbitMQ); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Daemon) Start() error {
	primaryEmailsChan, err := d.channel.Consume("emailer.incoming.queue", "emailer.incoming.queue", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("basic.consume: %v", err)
	}

	sgClient := adapter.NewSendGridClient(d.cfg.SendGridAPIKey)
	// spClient, err := adapter.NewSparkPostClient(d.cfg.SparkPostKey)
	// if err != nil {
	// 	return err
	// }
	fb := &services.FallBack{}

	// primaryDispatcher := services.NewDispatcher(primaryEmailsChan, 10, spClient, fb.PrimaryFallBack)
	primaryDispatcher := services.NewDispatcher(d.cfg, 10, sgClient, fb.PrimaryFallBack)
	primaryDispatcher.Start()

	return nil
}

func (d *Daemon) connectRabbitMQ() error {
	dialConfig := amqp.Config{
		Dial: amqp.DefaultDial(1 * time.Minute),
	}
	conn, err := amqp.DialConfig("amqp://guest:guest@rabbitmq:5672/", dialConfig)
	if err != nil {
		return fmt.Errorf("connection.open: %s", err)
	}

	c, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel.open: %s", err)
	}
	d.conn = conn
	d.channel = c
	return nil
}
