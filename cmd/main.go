package main

import (
	"flag"
	"log"
	"os"
	"slices"

	"github.com/claudemuller/daytime/pkg/tcp"
)

func main() {
	proto := flag.String("proto", "tcp", "Server protocol i.e. tcp, udp")
	flag.Parse()

	if !slices.Contains([]string{"tcp", "udp"}, *proto) {
		log.Printf("usage: %s --proto=[tcp/udp]", os.Args[0])
		return
	}

	if *proto == "tcp" {
		if err := tcp.Listen(); err != nil {
			log.Printf("failed to start server: %v", err)
		}
	} else {
		// udp
	}
}
