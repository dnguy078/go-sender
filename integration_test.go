// +build integration

package main

import (
	"flag"
	"net/http"
)

var (
	apiPort    = flag.Int("port", 4001, "http server port")
	sgAPIKey   = flag.String("sgkey", "", "sendgrid api key")
	testClient = http.DefaultClient
)

//TODO publish events to queue and assert all queues are empty after startup
