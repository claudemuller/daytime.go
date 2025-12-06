package tcp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/claudemuller/daytime/pkg/daytime"
)

type TCPServer struct {
	Host       string
	Port       int
	shutdownCh chan interface{}
	wg         sync.WaitGroup
	connMu     sync.Mutex
	conns      map[net.Conn]interface{}
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
		Host:       "localhost",
		Port:       13,
		shutdownCh: make(chan interface{}),
		conns:      make(map[net.Conn]interface{}),
	}
	for _, option := range options {
		option(&srv)
	}
	return &srv
}

func (s *TCPServer) Shutdown(ctx context.Context) error {
	close(s.shutdownCh)

	s.connMu.Lock()
	{
		for conn := range s.conns {
			conn.Close()
		}
	}
	s.connMu.Unlock()

	doneCh := make(chan interface{})
	go func() {
		s.wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		slog.Info("Connections closed gracefully")
		return nil

	case <-ctx.Done():
		slog.Warn("Shutdown timeout reached, closing open connections")
		return ctx.Err()
	}
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

		s.wg.Add(1)
		s.connMu.Lock()
		{
			s.conns[conn] = struct{}{}
		}
		s.connMu.Unlock()
		go s.handle(conn)
	}
}

func (s *TCPServer) handle(conn net.Conn) {
	defer func(c net.Conn) {
		c.Close()
		s.connMu.Lock()
		{
			delete(s.conns, c)
		}
		s.connMu.Unlock()
		s.wg.Done()
	}(conn)

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
