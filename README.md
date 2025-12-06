# Daytime Server

[![go](https://github.com/claudemuller/daytime.go/actions/workflows/go.yml/badge.svg)](https://github.com/claudemuller/daytime.go/actions/workflows/go.yml)

A TCP and UDP server implementation of the Daytime protocol (RFC867).

## Requirements

- [Go](https://go.dev/)

## Install Project Dependencies

```bash
make tidy
```

## Run

### Flags

| Flag | Type | Default | Note |
|---|---|---|---|
| `--proto/-proto` | `string` | `tcp` | |
| `--host/-host` | `string` | `localhost` | |
| `--port/-port` | `int` | `13` | Elevated privileges are required to run on port `13` |

```bash
make run [ARGS=-proto=tcp -port=9999]

# or
make build
./bin/daytime [-proto=tcp -port=9999]
```

## Build

```bash
make build
```

## Run tests

```bash
make test
```

## TODO

- add integration test for the servers
- load tests
