package main

import (
	"log"

	"github.com/claudemuller/daytime/pkg/tcp"
)

func main() {
	if err := tcp.Listen(); err != nil {
		log.Printf("failed to start server: %v", err)
	}
}
