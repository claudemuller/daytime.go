package main

import (
	"flag"
	"log"
	"os"
	"slices"

	"github.com/claudemuller/daytime/pkg/tcp"
)

func main() {
	proto := flag.String("proto", "tcp", "Server protocol (default: tcp)")
	host := flag.String("host", "", "Server host (default: localhost)")
	port := flag.Int("port", 13, "Server port (default: 13)")

	flag.Parse()

	if !slices.Contains([]string{"tcp", "udp"}, *proto) {
		log.Printf("Invalid protocal: %d", port)
		log.Printf("usage: %s --proto=[tcp/udp]", os.Args[0])
		return
	}

	if *port < 0 || *port > 65535 {
		log.Printf("Invalid port: %d", port)
		log.Printf("usage: %s --proto=[tcp/udp]", os.Args[0])
		return
	}

	if *proto == "tcp" {
		srv := tcp.NewServer(tcp.WithPort(*port), tcp.WithHost(*host))

		if err := srv.Listen(); err != nil {
			log.Printf("failed to start server: %v", err)
		}
	} else {
		// udp
	}
}
