package node

import (
	"log"
	"net"
)

func Start() error {
	log.Printf("[node] Starting node ID: %s", NodeId())

	// we need to start a tls server locally
	a, err := net.ResolveTCPAddr("tcp", ":65123")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return err
	}

	go nodeListener(l)
	// TODO also start node process

	return nil
}
