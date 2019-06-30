package node

import (
	"log"
	"net"

	"github.com/lucas-clemente/quic-go"
)

func Start() error {
	log.Printf("[node] Starting node ID: %s", NodeId())

	// start quic server
	qcfg := &quic.Config{
		MaxReceiveStreamFlowControlWindow:     16 * 1024 * 1024, // 16MB
		MaxReceiveConnectionFlowControlWindow: 16 * 1024 * 1024, // 16MB
		MaxIncomingStreams:                    10240,            // 10k streams
		MaxIncomingUniStreams:                 0,
		KeepAlive:                             true,
	}

	lis, err := quic.ListenAddr(":4242", srv.getCfg(), qcfg)
	if err != nil {
		return err
	}

	go quicListener(lis)

	// we need to start a tls server locally
	a, err := net.ResolveTCPAddr("tcp", ":65123")
	if err != nil {
		return err
	}

	// store listen addr
	srv.addr = a

	// start http server
	go srv.httpServe()

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return err
	}

	go nodeListener(l)
	// TODO also start node process

	return nil
}
