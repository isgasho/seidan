package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MagicalTux/hsm"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s command", os.Args[0])
		return
	}

	h, err := hsm.New()
	if err != nil {
		log.Printf("failed to initialize HSM: %s", err)
		os.Exit(1)
	}

	ks, err := h.ListKeysByName("seidan:cluster_name")
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
	log.Printf("found key: %s", k)
}
