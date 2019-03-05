package main

import (
	"log"

	"github.com/MagicalTux/seidan/core"
	"github.com/MagicalTux/seidan/node"
)

// initialize
func main() {
	log.Printf("[main] Initializing Seidan...")
	log.Printf("[main] Node ID: %s", node.NodeId())

	core.Shutdown()
}
