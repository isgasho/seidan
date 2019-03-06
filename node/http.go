package node

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"log"
	"math/big"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/MagicalTux/seidan/db"
)

type httpServer struct {
	incoming chan net.Conn
	addr     net.Addr

	cfg     *tls.Config
	cfgLock sync.Mutex
}

var srv = &httpServer{
	incoming: make(chan net.Conn),
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// check for TLS client auth
	if len(r.TLS.PeerCertificates) == 0 {
		// redirect to unsecure handler
		serveHttpInsecure(w, r)
		return
	}

	// TODO make sure peer certificate is valid

	w.Write([]byte("todo..."))
	// TODO
}

func (s *httpServer) getCfg() *tls.Config {
	s.cfgLock.Lock()
	defer s.cfgLock.Unlock()

	if s.cfg != nil {
		return s.cfg
	}

	nodeId := NodeId()
	ca := x509.NewCertPool()

	s.cfg = &tls.Config{
		MinVersion:       tls.VersionTLS12,
		RootCAs:          ca,
		ClientCAs:        ca,
		ServerName:       nodeId,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// get CA
	ca_cnt := 0

	if c, err := db.NewDbCursor([]byte("global")); err == nil {
		defer c.Close()
		k, v := c.Seek([]byte("internal:ca:"))
		for {
			if k == nil {
				break
			}
			ca.AppendCertsFromPEM(v)
			ca_cnt++
			k, v = c.Next()
		}
	}

	if ca_cnt == 0 {
		// no CA, need to go with self signed
		crtTpl := &x509.Certificate{
			BasicConstraintsValid: true,
			IsCA:                  true,
			SerialNumber:          big.NewInt(1),
			Issuer:                pkix.Name{CommonName: nodeId},
			Subject:               pkix.Name{CommonName: nodeId},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment,
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour), // 10 year
		}

		k := getNodeKey()
		crt_der, err := x509.CreateCertificate(rand.Reader, crtTpl, crtTpl, k.Public(), k)
		if err != nil {
			// shouldn't happen
			panic(err)
		}

		// get a tls cert ready
		s.cfg.Certificates = []tls.Certificate{tls.Certificate{
			Certificate: [][]byte{crt_der},
			PrivateKey:  k,
		}}

		return s.cfg
	}

	// TODO
	return s.cfg
}

func (s *httpServer) Push(n net.Conn) {
	sc := tls.Server(n, s.getCfg())
	// need to do TLS
	s.incoming <- sc
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
