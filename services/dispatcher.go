package services

import (
	"fmt"

	breaker "github.com/sony/gobreaker"
)

type Dispatcher struct {
	currentSender  Emailer
	primarySender  Emailer
	fallbackSender Emailer
	cb             *breaker.CircuitBreaker
}

type Emailer interface {
	Email() error
	Name() string
}

func NewDispatcher(name string, primarySender, fallbackSender Emailer) *Dispatcher {
	d := &Dispatcher{
		currentSender:  primarySender,
		primarySender:  primarySender,
		fallbackSender: fallbackSender,
	}

	// configure circuit breaker
	var st breaker.Settings
	st.Name = name
	st.ReadyToTrip = d.ReadToTrip
	st.OnStateChange = d.OnStateChange
	// st.Interval = 1 * time.Millisecond

	cb := breaker.NewCircuitBreaker(st)
	d.cb = cb

	return d
}

func (d *Dispatcher) Dispatch() error {
	fmt.Println("klsajdfldksfjlksdjfsdf")

	fmt.Printf("current sender: %s\n", d.currentSender.Name())
	_, err := d.cb.Execute(func() (interface{}, error) {
		if err := d.currentSender.Email(); err != nil {
			if d.currentSender == d.fallbackSender {
				fmt.Println("this is fatal")
				return nil, fmt.Errorf("this is a fatal error")
			}
			fmt.Println("not fatal")
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (d *Dispatcher) OnStateChange(name string, from breaker.State, to breaker.State) {
	if to == breaker.StateOpen || to == breaker.StateHalfOpen {
		d.currentSender = d.fallbackSender
		fmt.Println("setting currenSender", d.currentSender.Name())
	}
}

func (d *Dispatcher) ReadToTrip(counts breaker.Counts) bool {
	fmt.Printf("counts %+v", counts)
	return counts.ConsecutiveFailures > 1
}
