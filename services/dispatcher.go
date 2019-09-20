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

// FallBackFn is the fallback function that gets called when email provider returns an error
type FallBackFn func(request.EmailRequest)

// NewDispatcher consumes from an amqp.Delivery channel and pipes messages to a buffer channel where
// pool of workers will process the email request
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

// SetEmailer lets  you change the emailer
func (d *Dispatcher) SetEmailer(emailer Emailer) {
	d.emailer = emailer
}

// Start spins up EmailWorkers
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

// EmailWorker sends emails
type EmailWorker struct {
	emailer      Emailer
	queue        chan amqp.Delivery
	fallBackFunc FallBackFn
	quit         chan bool
}

// Work consumes off a buffer queue and sends emails until the queue is closed or a quit signal is received
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
				payload.Reject(false)
			}

			if err := payload.Ack(false); err != nil {
				log.Println(err)
			}
		case <-w.quit:
			return
		}
	}
}
