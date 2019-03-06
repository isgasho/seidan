package node

import (
	"log"
	"net"
	"time"

	"github.com/MagicalTux/seidan/db"
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
		go handleClient(c)
	}
}

func handleClient(c *net.TCPConn) {
	// if we have a CA, initialize a tls connection with said CA
	// if not, generate & return a CSR
	ca, err := db.SimpleGet([]byte("global"), []byte("ca"))
	if err != nil {
		srv.Push(c)
		return
	}

	// TODO
	_ = ca
	c.Close()
}
