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

### Flags

| Flag | Type | Default | Note |
|---|---|---|---|
| `--proto/-proto` | `string` | `tcp` | |
| `--host/-host` | `string` | `localhost` | |
| `--port/-port` | `int` | `13` | Elevated privileges are required to run on port `13` |

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

- add integration test for the servers
- load tests
