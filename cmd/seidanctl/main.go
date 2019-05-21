package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/MagicalTux/hsm"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s command\n", os.Args[0])
		return
	}
	clusterName := "test"

	h, err := hsm.New()
	if err != nil {
		log.Printf("failed to initialize HSM: %s", err)
		os.Exit(1)
	}

	ks, err := h.ListKeysByName("seidan:" + clusterName)
	if err != nil {
		log.Printf("failed to list HSM keys: %s", err)
		os.Exit(1)
	} else if len(ks) == 0 {
		// Generate?
		// NOTE: ecdsa, rsa only
		log.Printf("failed to list HSM keys: no keys. Please generate one.")
		os.Exit(1)
	}
	k := ks[0]

	cert, err := h.GetCertificate("seidan:" + clusterName)
	if err == os.ErrNotExist {
		// need to create certificate
		log.Printf("Generating self-signed CA certificate")
		now := time.Now()
		caCrt := &x509.Certificate{
			BasicConstraintsValid: true,
			IsCA:                  true,
			SerialNumber:          big.NewInt(1),
			Issuer:                pkix.Name{CommonName: "seidan:" + clusterName + " CA"},
			Subject:               pkix.Name{CommonName: "seidan:" + clusterName + " CA"},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment,
			NotBefore:             now,
			NotAfter:              now.Add(25 * 365 * 24 * time.Hour), // +25 years (more or less)
			MaxPathLen:            1,
		}
		ca_crt_der, err := x509.CreateCertificate(rand.Reader, caCrt, caCrt, k.Public(), k)
		if err != nil {
			log.Printf("failed to generate CA: %s", err)
			os.Exit(2)
		}
		cert, err = x509.ParseCertificate(ca_crt_der)
		if err != nil {
			log.Printf("failed to parse generated CA: %s", err)
			os.Exit(2)
		}

		// store certificate
		err = h.PutCertificate("seidan:"+clusterName, cert)
		if err != nil {
			log.Printf("failed to store generated CA: %s", err)
			os.Exit(2)
		}
	} else if err != nil {
		log.Printf("failed to get CA certificate: %s", err)
		os.Exit(1)
	}

	log.Printf("found key: %s", k)

	ca_crt_pem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	log.Printf("found cert: %s", ca_crt_pem)
}
