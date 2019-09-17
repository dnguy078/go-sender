package main

import (
	"log"

	"github.com/dnguy078/go-sender/config"
	"github.com/dnguy078/go-sender/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("err", err)
	}

	server := server.New(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("unable to start http server: %v", err)
	}
}
