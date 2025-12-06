package tcp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/claudemuller/daytime/pkg/daytime"
)

type Server struct {
	Host string
	Port int
}

type Option func(*Server)

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithHost(host string) Option {
	return func(s *Server) {
		if host != "" {
			s.Host = host
		}
	}
}

func NewServer(options ...Option) *Server {
	srv := &Server{
		Host: "localhost",
		Port: 13,
	}
	for _, option := range options {
		option(srv)
	}
	return srv
}

func (s *Server) Listen() error {
	log.Printf("Listening on %s:%d\n", s.Host, s.Port)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
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
