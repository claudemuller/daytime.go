# Daytime Server

[![go](https://github.com/claudemuller/daytime.go/actions/workflows/go.yml/badge.svg)](https://github.com/claudemuller/daytime.go/actions/workflows/go.yml)

A TCP and UDP server implementation of the Daytime protocol (RFC867).

## Requirements

- [Go](https://go.dev/)

## Install Project Dependencies

```bash
go mod tidy
```

## Run

If no `--proto` is specified then TCP is assumed.

```bash
go run cmd/main.go [--proto=[tcp/udp]]
```

## Build

```bash
go build cmd/main.go
```

## Run tests

```bash
go test ./...
```

## TODO

- manage connections with goroutine pool
