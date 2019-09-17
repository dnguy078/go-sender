package services

import (
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
		return d.primarySender.Email()
	}, func(err error) error {
		// fall back to fallback sender if primary fails
		return d.fallbackSender.Email()
	})

	return nil
}
