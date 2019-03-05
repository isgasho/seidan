package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MagicalTux/seidan/core"
	"github.com/MagicalTux/seidan/node"
)

var shutdownChannel = make(chan struct{})

func shutdown() {
	log.Println("[main] shutting down...")
	close(shutdownChannel)
}

func setupSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		shutdown()
	}()
}

// initialize
func main() {
	setupSignals()

	log.Printf("[main] Initializing Seidan...")
	node.Start()

	<-shutdownChannel

	core.Shutdown()
}
