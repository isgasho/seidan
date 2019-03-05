package main

import "log"

// initialize
func main() {
	log.Printf("[main] Initializing Seidan...")
	initDb()

	shutdownDb()
}
