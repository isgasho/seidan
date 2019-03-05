package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s command", os.Args[0])
		return
	}
}
