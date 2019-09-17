package services

import (
	"log"

	"github.com/afex/hystrix-go/hystrix"
)

type Dispatcher struct {
	name           string
	primarySender  Emailer
	fallbackSender Emailer
}

type Emailer interface {
	Email() error
	Type() string
}

func NewDispatcher(name string, primarySender, fallbackSender Emailer) *Dispatcher {
	// todo make configuable
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  1,
		RequestVolumeThreshold: 1,
	})

	return &Dispatcher{
		name:           name,
		primarySender:  primarySender,
		fallbackSender: fallbackSender,
	}
}

func (d *Dispatcher) SetPrimary(primarySender Emailer) {
	d.primarySender = primarySender
}

func (d *Dispatcher) SetFallback(fallbackSender Emailer) {
	d.fallbackSender = fallbackSender
}

func (d *Dispatcher) Dispatch() error {
	hystrix.Go(d.name, func() error {
		log.Printf("sending email from %s", d.primarySender.Type())
		return d.primarySender.Email()
	}, func(err error) error {
		// fall back to fallback sender if primary fails
		log.Printf("sending email from %s", d.fallbackSender.Type())
		return d.fallbackSender.Email()
	})

	return nil
}
