# go-sender
go-sender is an event-based microservice that sends basic emails through various email providers (sendgrid, sparkpost). If one of the services goes down, go-sender will quickly failover to the secondary configuable email provider. It consumes off rabbitmq for email events and hands the event to a pool of email workers that calls out to the apppropriate email provider. 

## Running locally
To Run (replace docker-compose.yml with appropriate keys):
  ```
  docker-compose build && docker-compose up
  ```
To Test:
    ```
    go test ./...
    ```

Integration test
    ```
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

## Architecture
