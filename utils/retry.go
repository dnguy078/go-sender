package utils

import (
	"fmt"
	"time"
)

var DefaultMaxRetry = 5

// RetryHandler wraps a function and performs retry logic
type RetryHandler struct {
	MaxRetries int
	Backoff    BackoffStrategy
}

// NewRetryWrapper returns a new RetryHandler
func NewRetryWrapper(maxRetries int, backoff BackoffStrategy) *RetryHandler {
	return &RetryHandler{
		MaxRetries: maxRetries,
		Backoff:    backoff,
	}
}

// BackoffStrategy is used to determine how long a retry request should wait until attempted
type BackoffStrategy func(retry int) time.Duration

// Work is a function that gets retried if the function returns an error
type Work func() error

// ExponentialBackoff returns increasing backoffs by a power of 2
func ExponentialBackoff(i int) time.Duration {
	return time.Duration(1<<uint(i)) * time.Second
}

// LinearBackoff returns increasing durations, each a second longer than the last
func LinearBackoff(i int) time.Duration {
	return time.Duration(i) * time.Second
}

// WithRetry will retry a function
func (rh *RetryHandler) WithRetry(fn Work) error {
	attempt := 1
	var err error

	for {
		err = fn()
		if err == nil {
			break
		}

		if attempt >= rh.MaxRetries {
			return fmt.Errorf("reached max limit, attempted: %d, maxRetries: %d, error: %s", attempt, rh.MaxRetries, err.Error())
		}

		attempt++
		<-time.After(rh.Backoff(attempt) + 1*time.Microsecond)
	}

	return nil
}
