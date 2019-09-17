# go-sender
go-sender is an http service that sends basic emails through various email providers (sendgrid, mailgun). If one of the services goes down, go-sender will quickly failover to the secondary configuable email provider.

Failover logic is using hystrix circuit breakers. A hystrix dashboard is provided on http://localhost:7979/hystrix-dashboard/. Must provide configure a stream (http://go-sender:81/hystrix.stream)

## Running locally
To Run:
  ```
  docker-compose build && docker-compose up
  ```
To Test:
    ```
    go test ./...
    ```

Integration test
    ```
    go test -tags=integration -sgkey=sgAPIKey -mgkey=mgAPIKey
    ```

## API
	// curl -X POST localhost:4001/email
	```
	{
        "toEmail": "someemail@email.com",
        "fromEmail": "fromemail@email.com",
        "subject": "somesubject@email.com",
        "text": "sometext"
	}
	```
