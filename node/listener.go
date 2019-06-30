package node

import (
	"log"
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
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

func quicListener(l quic.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("[quic] failed to accept: %s", err)
			return // closed?
		}

		// TODO
		c.Close()
	}
}

type quicConn struct {
	quic.Stream
}

func fakeQuicConn(s quic.Stream, err error) (net.Conn, error) {
	if err != nil {
		return nil, err
	}

	return &quicConn{s}, err
}

func (c *quicConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 80}
}

func (c *quicConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 80}
}
