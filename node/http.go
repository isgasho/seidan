package node

import (
	"errors"
	"log"
	"net"
	"net/http"
)

type httpServer struct {
	incoming chan net.Conn
	addr     net.Addr
}

var srv = &httpServer{
	incoming: make(chan net.Conn),
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		// redirect to unsecure handler
		serveHttpInsecure(w, r)
		return
	}

	w.Write([]byte("todo..."))
	// TODO
}

func (s *httpServer) Push(n net.Conn) {
	s.incoming <- n
}

func (s *httpServer) Accept() (net.Conn, error) {
	c, ok := <-s.incoming
	if !ok {
		return nil, errors.New("failed to accept: channel closed")
	}

	return c, nil
}

func (s *httpServer) Close() error {
	// do not actually close
	return nil
}

func (s *httpServer) Addr() net.Addr {
	return s.addr
}

func (s *httpServer) httpServe() {
	err := http.Serve(s, s)
	if err != nil {
		log.Printf("[node] failed to serve http: %s", err)
	}
}
