package main

import (
	"fmt"

	"github.com/dnguy078/go-sender/server"
)

func main() {
	server := server.New()
	err := server.Start()
	if err != nil {
		fmt.Printf("unable to start up %v", err)
	}
}
