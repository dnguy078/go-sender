package main

import (
	"fmt"

	"github.com/dnguy078/go-sender/server"
)

func main() {
	cfg := server.ServerConfig{
		Port: 4001,
	}

	server := server.New(cfg)
	if err := server.Start(); err != nil {
		fmt.Printf("unable to start up %v", err)
	}
}
