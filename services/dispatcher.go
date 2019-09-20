package services

import (
	"fmt"

	"github.com/dnguy078/go-sender/config"
	"github.com/dnguy078/go-sender/request"
	"github.com/streadway/amqp"
)

type Dispatcher struct {
	name       string
	jobQueue   chan amqp.Delivery
	emailer    Emailer
	fallBackFn FallBackFn
	quit       chan bool
}

type Emailer interface {
	Email(request.EmailRequest) error
	Type() string
}

type DispatcherConfigs struct {
	maxQueueSize int
	maxWorker    int
}

type FallBackFn func(request.EmailRequest)

func NewDispatcher(cfg config.Config, maxQueueSize int, emailer Emailer, fallbackFn FallBackFn) *Dispatcher {
	primaryEmailsChan, err := d.channel.Consume("emailer.incoming.queue", "emailer.incoming.queue", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("basic.consume: %v", err)
	}

	jobQueue := make(chan amqp.Delivery, maxQueueSize)

	go func() {
		for j := range msgs {
			jobQueue <- j
		}
	}()

	return &Dispatcher{
		jobQueue:   jobQueue,
		emailer:    emailer,
		fallBackFn: fallbackFn,
		quit:       make(chan bool, 1),
	}
}

func (d *Dispatcher) SetEmailer(emailer Emailer) {
	d.emailer = emailer
}

func (d *Dispatcher) Start() {
	fmt.Println("starting dispatcher")
	for i := 0; i < 100; i++ {
		worker := &EmailWorker{
			emailer:      d.emailer,
			queue:        d.jobQueue,
			fallBackFunc: d.fallBackFn,
			quit:         d.quit,
		}
		go worker.Work()
	}
}

type EmailWorker struct {
	emailer      Emailer
	queue        chan amqp.Delivery
	fallBackFunc FallBackFn
	quit         chan bool
}

func (w *EmailWorker) Work() {
	for {
		select {
		case payload, ok := <-w.queue:
			if !ok {
				return
			}
			req, err := request.Validate(payload.Body)
			if err != nil {
				fmt.Printf("%+v", req)
				payload.Reject(false)
			}

			if err := w.emailer.Email(req); err != nil {
				fmt.Println(err)
				w.fallBackFunc(req)
			}

			if err := payload.Ack(false); err != nil {
				// log
			}
		case <-w.quit:
			return
		}
	}
}
