package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dnguy078/go-sender/config"
	"github.com/dnguy078/go-sender/daemon"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("err", err)
	}

	daemon, err := daemon.New(cfg)
	if err != nil {
		log.Fatalf("unable to init daemon: %+v", err)
	}

	if err := daemon.Start(); err != nil {
		log.Fatalf("unable to start http server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	// TODO: add a graceful shutdown and clean up connections
}
