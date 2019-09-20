package services

import (
	"log"

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

func NewDispatcher(msgs <-chan amqp.Delivery, maxQueueSize int, emailer Emailer, fallbackFn FallBackFn) *Dispatcher {
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
	log.Println("Starting dispatcher - ", d.emailer.Type())
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
				log.Println(err)
				payload.Reject(false)
			}

			if err := w.emailer.Email(req); err != nil {
				log.Println(err)
				w.fallBackFunc(req)
			}

			if err := payload.Ack(false); err != nil {
				log.Println(err)
			}
		case <-w.quit:
			return
		}
	}
}
