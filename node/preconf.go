package node

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"net/http"
)

func serveHttpInsecure(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/csr":
		// got no CA, make a csr
		tpl := &x509.CertificateRequest{
			Subject: pkix.Name{CommonName: NodeId()},
		}
		csr, err := x509.CreateCertificateRequest(rand.Reader, tpl, getNodeKey())
		if err != nil {
			log.Printf("[node] failed to generate CSR: %s", err)
			return // will close connection
		}

		// encode csr to PEM
		d := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})
		w.Write(d)
		return
	}
}
