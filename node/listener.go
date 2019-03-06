package node

import (
	"log"
	"net"
	"time"
)

func nodeListener(l *net.TCPListener) {
	var tempDelay time.Duration

	for {
		c, err := l.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("[node] Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			log.Printf("[node] Fatal error on accept: %v", err)
			// Give up :(
			return
		}
		tempDelay = 0
		go srv.Push(c)
	}
}
