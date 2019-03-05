package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"

	"github.com/MagicalTux/seidan/db"
)

func getNodeKey() *ecdsa.PrivateKey {
	// get node key from db
	v, err := db.SimpleGet([]byte("local"), []byte("node_key"))
	if err != nil {
		// generate new key
		k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			// error with random generator? that's bad
			panic(err)
		}
		v, err = x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			// shouldn't happen
			panic(err)
		}

		// store key
		err = db.SimpleSet([]byte("local"), []byte("node_key"), v)
		if err != nil {
			panic(err)
		}

		return k
	}

	// decode key
	k, err := x509.ParsePKCS8PrivateKey(v)
	if err != nil {
		panic(err)
	}

	// TODO support other key types?

	return k.(*ecdsa.PrivateKey)
}

func GetKeyHash() []byte {
	k := getNodeKey()
	v, err := x509.MarshalPKIXPublicKey(k.Public())
	if err != nil {
		// shouldn't happen
		panic(err)
	}
	h := sha256.Sum256(v)
	return h[:]
}
