package tcp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/claudemuller/daytime/pkg/daytime"
)

type TCPServer struct {
	Host string
	Port int
}

type Option func(*TCPServer)

func WithPort(port int) Option {
	return func(s *TCPServer) {
		s.Port = port
	}
}

func WithHost(host string) Option {
	return func(s *TCPServer) {
		if host != "" {
			s.Host = host
		}
	}
}

func NewServer(options ...Option) *TCPServer {
	srv := TCPServer{
		Host: "localhost",
		Port: 13,
	}
	for _, option := range options {
		option(&srv)
	}
	return &srv
}

func (s *TCPServer) Shutdown(ctx context.Context) error {
	return nil
}

func (s *TCPServer) Listen() error {
	slog.Info("Listening for connections", "host", s.Host, "port", s.Port)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return fmt.Errorf("failed to listen for connections: %w", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "source", err)
			continue
		}

		slog.Info("Accepted connection", "source", conn.RemoteAddr())

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	now := daytime.Get(time.Now)

	n, err := conn.Write([]byte(now))
	if err != nil {
		slog.Error("Failed to send data", "error", err)
	}
	if n <= 0 {
		slog.Error("Zero bytes sent", "error", err)
	}

	slog.Info("Successfully sent datetime", "bytes sent", n, "destination", conn.RemoteAddr())
}
