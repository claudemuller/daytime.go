package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/claudemuller/daytime/pkg/tcp"
)

const shutdownTimeout = time.Second * 5

type Server interface {
	Listen() error
	Shutdown(ctx context.Context) error
}

func main() {
	proto := flag.String("proto", "tcp", "Server protocol (default: tcp)")
	host := flag.String("host", "", "Server host (default: localhost)")
	port := flag.Int("port", 13, "Server port (default: 13)")

	flag.Parse()

	if !slices.Contains([]string{"tcp", "udp"}, *proto) {
		slog.Error("Invalid protocal", "protocal", port)
		fmt.Printf("usage: %s --proto=[tcp/udp]", os.Args[0])
		return
	}

	if *port < 0 || *port > 65535 {
		slog.Error("Invalid port", "port", port)
		fmt.Printf("usage: %s --proto=[tcp/udp]", os.Args[0])
		return
	}

	srvErrs := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	var srv Server

	go func() {
		if *proto == "tcp" {
			srv = tcp.NewServer(tcp.WithPort(*port), tcp.WithHost(*host))
			srvErrs <- srv.Listen()
		} else {
			// udp
		}
	}()

	select {
	case err := <-srvErrs:
		slog.Error("failed to start server", "error", err)

	case sig := <-shutdown:
		slog.Info("Server shutdown started", "signal", sig)
		defer slog.Info("Server shutdown completed", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Error shutdowning down server", "error", err)
		}
	}
}
