package adapter

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dnguy078/go-sender/request"

	sp "github.com/SparkPost/gosparkpost"
)

// SparkPostClient is a wrapper around SparkPostClient's API
type SparkPostClient struct {
	client sp.Client
}

// NewSparkPostClient returns a new SparkPostClient
func NewSparkPostClient(apikey string) (*SparkPostClient, error) {
	var sparky sp.Client
	fmt.Println("apikey", apikey)
	err := sparky.Init(&sp.Config{ApiKey: "fced6a2b44ba15eade9539418fc91caab4c6f16d"})
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
		return nil, err
	}

	return &SparkPostClient{
		client: sparky,
	}, nil
}

func (spClient *SparkPostClient) Type() string {
	return "sparkpost"
}

// Email performs a http request to send emails through SP
func (spClient *SparkPostClient) Email(req request.EmailRequest) error {
	sandbox := true
	tx := &sp.Transmission{
		Recipients: []string{req.ToEmail},
		Options:    &sp.TxOptions{Sandbox: &sandbox},
		Content: sp.Content{
			Text:    req.Text,
			From:    "testing@sparkpostbox.com",
			Subject: req.Subject,
		},
	}
	_, resp, err := spClient.client.Send(tx)
	if err != nil {
		return err
	}

	if resp.HTTP.StatusCode > http.StatusInternalServerError {
		return errors.New("sparkpost down")
	}

	return nil
}
