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

const poolSize = 1000

type TCPServer struct {
	Host       string
	Port       int
	shutdownCh chan interface{}
	connCh     chan net.Conn
	wg         sync.WaitGroup
	listener   net.Listener
	listenerMu sync.Mutex
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
		connCh:     make(chan net.Conn, poolSize),
	}
	for _, option := range options {
		option(&srv)
	}
	return &srv
}

func (s *TCPServer) Listen() error {
	slog.Info("Listening for connections", "host", s.Host, "port", s.Port)

	var err error

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return fmt.Errorf("failed to listen for connections: %w", err)
	}
	defer s.listener.Close()

	slog.Info("Setting up connection pool", "size", poolSize)
	for i := range poolSize {
		s.wg.Add(1)
		go s.worker(i)
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdownCh:
				slog.Info("Server shutting down, stop accepting connections")
				close(s.connCh)
				s.wg.Wait()
				return nil

			default:
				slog.Error("Failed to accept connection", "remote", err)
				continue
			}
		}

		select {
		case s.connCh <- conn:
			slog.Info("Connection sent to worker", "remote", conn.RemoteAddr())

		case <-s.shutdownCh:
			conn.Close()
			slog.Info("Server shutting down, rejected connection", "remote", conn.RemoteAddr())
		}
	}
}

func (s *TCPServer) worker(id int) {
	defer s.wg.Done()

	for conn := range s.connCh {
		slog.Info("Worker handling connection", "worker", id, "remote", conn.RemoteAddr())
		s.handle(conn)
	}

	slog.Info("Worker done", "id", id)
}

func (s *TCPServer) handle(conn net.Conn) {
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

func (s *TCPServer) Shutdown(ctx context.Context) error {
	close(s.shutdownCh)

	s.listenerMu.Lock()
	{
		s.listener.Close()
	}
	s.listenerMu.Unlock()

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
