// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/dnguy078/go-sender/adapter"
	"github.com/dnguy078/go-sender/config"
	"github.com/dnguy078/go-sender/request"
	"github.com/dnguy078/go-sender/server"
)

var (
	apiPort    = flag.Int("port", 4001, "http server port")
	sgAPIKey   = flag.String("sgkey", "", "sendgrid api key")
	testClient = http.DefaultClient
)

func TestMain(m *testing.M) {
	flag.Parse()

	m.Run()
}

func emailRequest(t *testing.T) (int, error) {
	req := request.EmailRequest{}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("unable to marshal test email request, err %+v", err)
	}
	resp, err := testClient.Post(fmt.Sprintf("http://localhost:%d/email", *apiPort), "", bytes.NewBuffer(b))
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, err
}

func TestSendGridMailSend(t *testing.T) {
	cfg := config.Config{
		SendGridAPIKey: *sgAPIKey,
	}
	sg := adapter.NewSendGridClient(cfg.SendGridAPIKey)
	if err := sg.Email(); err != nil {
		t.Error(err)
	}
}

func xTestAPISend(t *testing.T) {
	cfg := config.Config{
		SenderAPIPort:  *apiPort,
		SendGridAPIKey: *sgAPIKey,
	}

	s := server.New(cfg)
	go s.Start()

	code, err := emailRequest(t)
	if err != nil {
		t.Error(err)
	}
	t.Error(code)

	if code != http.StatusOK {
		t.Errorf("expected bleh , status ")
	}
}
