# go-sender
go-sender is an event-based microservice that sends basic emails through various email providers (sendgrid, sparkpost). If one of the services goes down, go-sender will quickly failover to the secondary configuable email provider. It consumes off rabbitmq for email events and hands the event to a pool of email workers that calls out to the apppropriate email provider.

## Running locally
To Run (replace docker-compose.yml with appropriate keys):
  ```
  docker-compose build && docker-compose up
  ```
Unit Test:
    ```
    go test ./...
    ```

Integration test
    ```
    docker-compose build && docker-compose up -d
    go test -tags=integration -sgkey=sgAPIKey -spkey=spAPIKey
    ```

## Email Request Event
	// Users of this service will need to send this email event payload to rabbitmq. `amqp://guest:guest@rabbitmq:5672/` (local development),
  ** Alternatively, for testing purposes, one can submit the payload through the rabbitmq management page (http://localhost:15672/#/) **

	```
	{
        "toEmail": "someemail@email.com",
        "fromEmail": "fromemail@email.com",
        "subject": "somesubject@email.com",
        "text": "sometext"
	}
	```

## External Dependencies
1. RabbitMQ - consumes/publishes messages to `emailer.incoming.queue`, `emailer.retry.queue`, `emailer.errors.queue`
2. SendGrid - primary email provider
3. Sparkpost - secondary email provider

## Architecture limitations
An event base architecture for sending out email has multiple advantages in terms of scale and reliability. Events could be easily be replayed in case of downtime from external APIs. Initially I had written this service as an HTTP API with a circuit breaker than fell back to a secondary email provider. While this was fine, it also poses the issue if both email providers were down. Essentially we'd have to store those failed events to be retried somewhere. Storing these request in RabbitMQ or Kafka allows us to replay events without much difficulty.

This email system has its limitations in that end users would know necessarily know if their request failed/succeeded. I feel like in most email systems, emails are sent asynchronously. Given time restrictions, it would like to

1. Fix logging in the application
2. Remove some hard coded values
3. Basic request validation

**Additional Notes**
- Sparkpost does not allow you send from a non whitelisted domain. Their free account only lets you send from a sandbox account setting up to 50 emails. I used all 50 calls.
