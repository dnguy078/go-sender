package daemon

import (
	"github.com/dnguy078/go-sender/adapter"
	"github.com/dnguy078/go-sender/config"
	"github.com/dnguy078/go-sender/services"
	"github.com/dnguy078/go-sender/utils"
)

type Daemon struct {
	cfg config.Config
}

func New(cfg config.Config) (*Daemon, error) {
	d := &Daemon{
		cfg: cfg,
	}

	return d, nil
}

func (d *Daemon) Start() error {
	// performing exponential backoff healthcheck on rabbitmq before attempting to establish connection
	// docker-compose not great for determining start up order :(
	rh := utils.NewRetryWrapper(10, utils.ExponentialBackoff)
	err := rh.WithRetry(adapter.RabbitHealthcheck)
	if err != nil {
		return err
	}

	sgClient := adapter.NewSendGridClient(d.cfg.SendGridAPIKey)
	spClient, err := adapter.NewSparkPostClient(d.cfg.SparkPostKey)
	if err != nil {
		return err
	}

	publisherClient, err := adapter.NewRabbitClient(
		d.cfg.RabbitUsername,
		d.cfg.RabbitPassword,
		d.cfg.RabbitAddr,
		d.cfg.RabbitPort,
	)
	if err != nil {
		return err
	}

	fb := &services.FallBack{
		Publisher: publisherClient,
	}

	consumerClient, err := adapter.NewRabbitClient(
		d.cfg.RabbitUsername,
		d.cfg.RabbitPassword,
		d.cfg.RabbitAddr,
		d.cfg.RabbitPort,
	)
	if err != nil {
		return err
	}

	sgMsg, err := consumerClient.Consume("emailer.incoming.queue", "sendgrid")
	if err != nil {
		return err
	}

	spMsg, err := consumerClient.Consume("emailer.retry.queue", "sparkpost")
	if err != nil {
		return err
	}

	// primaryDispatcher := services.NewDispatcher(primaryEmailsChan, 10, spClient, fb.PrimaryFallBack)
	primaryDispatcher := services.NewDispatcher(sgMsg, 10, sgClient, fb.PrimaryFallBack)
	primaryDispatcher.Start()

	secondaryDispatcher := services.NewDispatcher(spMsg, 10, spClient, fb.SecondaryFallBack)
	secondaryDispatcher.Start()

	return nil
}
