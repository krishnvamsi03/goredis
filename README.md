# GoRedis

A lightweight, in-memory key-value store written in Go, inspired by Redis. GoRedis uses a custom **GRESP** protocol, supports multiple data types with TTL, and offers optional snapshot persistence.

## Features

- **In-memory storage** — All data in RAM for fast access
- **Multiple data types** — Strings, lists, and integers
- **TTL (time-to-live)** — Automatic key expiry after a set duration
- **Concurrent clients** — Goroutines and channels for multi-client handling
- **Persistence** — Configurable snapshots to disk (protobuf-backed)
- **GRESP protocol** — Simple, line-based request/response format

## Prerequisites

- [Go](https://go.dev/) 1.22 or later

## Quick Start

**1. Start the server** (default: `localhost:6379`)

```bash
make run
```

**2. Connect with the CLI**

```bash
make run-cli
```

Type `quit` to exit the client.

## Building

| Target        | Description                    |
|---------------|--------------------------------|
| `make build`  | Build server (Linux amd64)     |
| `make build-windows` | Build server for Windows |
| `make build-cli`    | Build CLI binary         |

Server binary: `./goredis` (or `goredis.exe` on Windows).  
Config path can be set via `GOREDIS_CONFIG_PATH` when running the server.

## Configuration

Config is loaded from `./config/config.yaml` by default. Override with:

```bash
GOREDIS_CONFIG_PATH=/path/to/config.yaml go run cmd/server/main.go
```

Example `config/config.yaml`:

```yaml
log:
  level: info          # log level

srvoptions:
  port: "6379"         # server listen port

persistent:
  interval: 1          # snapshot interval
  unit: "m"            # s | m | h (seconds, minutes, hours)
  path: "kvdb.dmp"     # snapshot file path
```

## CLI Usage

Connect to a different host or port:

```bash
go run cmd/cli/main.go --address localhost --port 6380
# short flags:
go run cmd/cli/main.go -a localhost -p 6380
```

## Supported Data Types & Commands

| Type   | Description        |
|--------|--------------------|
| `STRING` | Single string value |
| `LIST`   | Colon-separated list (e.g. `A:B:C`) |
| `INT`    | Integer (use with INCR/DECR) |

### Command reference

| Command | Description | Example |
|---------|-------------|---------|
| `PING`  | Health check | `PING` |
| `SET`   | Set key with type and optional TTL | `SET mykey myvalue STRING 300` |
| `GET`   | Get value by key | `GET mykey` |
| `DEL`   | Delete key | `DEL mykey` |
| `KEYS`  | List keys (supports `*` wildcard) | `KEYS *` |
| `EXPR`  | Set TTL on existing key | `EXPR mykey 60` |
| `PUSH`  | Append to list | `PUSH mylist ITEM1,ITEM2` |
| `POP`   | Remove from list (optional: L/R and count) | `POP mylist` or `POP mylist L 1` |
| `INCR`  | Increment integer key | `INCR counter` |
| `DECR`  | Decrement integer key | `DECR counter` |

### Examples

```text
# String with TTL (300 seconds)
SET greeting "Hello, World" STRING 300
GET greeting

# Integer (INCR type)
SET counter 0 INT 60
INCR counter
DECR counter

# List
SET fruits APPLE:BANANA:GRAPES LIST 300
PUSH fruits MANGO,ORANGE
POP fruits
POP fruits R 2          # pop 2 from right
```

## Protocol (GRESP)

Clients communicate using **GRESP**, a line-based protocol. Each request starts with a line like:

```text
GRESP OP <COMMAND> KEY <key> [DATA_TYPE <type>] [TTL <seconds>]
```

Commands that send a body (e.g. SET, PUSH) include a `CONTENT_LENGTH` line followed by the payload. See the `GRESP` file in the repo for full examples.

## Project Structure

```text
goredis/
├── cmd/
│   ├── server/          # Server entrypoint
│   └── cli/             # CLI client entrypoint
├── cli/                 # Client implementation
├── common/              # Config, logger
├── internal/
│   ├── command/         # Command handlers (GET, SET, etc.)
│   ├── event_processor/ # Event loop and request processing
│   ├── protocol/        # GRESP parser
│   ├── store/           # Key-value store and persistence
│   └── ...
├── config/              # config.yaml
└── proto/               # Persistence protobufs
```

## Development

```bash
make fmt       # Format code
make lint      # Run golangci-lint
make test      # Run tests
make coverage  # Generate coverage report (output in coverage/)
make check     # fmt + lint
```

## License

See [LICENSE](LICENSE) in the repository.
