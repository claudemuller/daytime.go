package tcp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/claudemuller/daytime/pkg/daytime"
)

const port = 13

// TODO: implement thread pool

func Listen() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen for connections: %w", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	now := daytime.Get(time.Now)
	n, err := conn.Write([]byte(now))
	if err != nil {
		log.Printf("failed to send data: %v", err)
	}
	if n <= 0 {
		log.Printf("%d bytes sent: %v", n, err)
	}
}
