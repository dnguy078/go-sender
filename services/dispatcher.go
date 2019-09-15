package services

import (
	"github.com/sony/gobreaker"
	breaker "github.com/sony/gobreaker"
)

type Dispatcher struct {
	primarySender  Emailer
	fallbackSender Emailer
	cr             *breaker.CircuitBreaker
}

type Emailer interface {
	Email()
}

func NewDispatcher(primarySender, fallbackSender Emailer) *Dispatcher {
	return &Dispatcher{
		primarySender:  primarySender,
		fallbackSender: fallbackSender,
	}
}

func (d *Dispatcher) RegisterPrimary() error {

	return nil
}

func (d *Dispatcher) RegisterFallBack() error {

	return nil
}

func (d *Dispatcher) Dispatch() error {

	return nil
}

func NewBreaker() *breaker.CircuitBreaker {
	var st gobreaker.Settings
	st.Name = "Email"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	st.OnStateChange = func() {

	}
}
