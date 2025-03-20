# GoRedis

GoRedis is a lightweight, in-memory key-value data store written in Go, inspired by Redis. It provides high-performance data operations, supports multiple data structures, and includes optional persistence for durability.

## Features
- **In-Memory Storage** – Stores all data in RAM for lightning-fast access.
- **Flexible Data Model** – Supports string, list, and integer-based key-value storage.
- **TTL (Time-To-Live) Expiry** – Enables automatic key expiration after a specified duration.
- **High Concurrency** – Leverages Go’s goroutines and channels for efficient multi-client handling.
- **Persistence** – Supports snapshotting data to disk for long-term storage.
- **Lightweight Protocol** – Offers a simple command-based interface for seamless integration.

## Installation
To start the server, simply run:
```sh
make run  # Starts the server on port 6379. Modify the config.yaml file to change the port.
```

## Usage
Interact with GoRedis using the built-in CLI tool:
```sh
make run-cli  # Connects to the server at localhost:6379 by default.

# To connect to a different port, use:
go run cmd/cli/main.go --a localhost --p 6380
```

## Supported Data Types
```
STRING, LIST, INTEGER
```

## Example Commands
Commands are strongly typed so you need to specify datatypes for set command.
```
SET test(KEY_NAME) value(VALUE) STRING(DATA_TYPE) 300 (TTL)
GET test
DEL test
SET test(KEY_NAME) value(VALUE) INCR(DATA_TYPE) 300 (TTL)
INCR test
DECR test
SET test_list(KEY_NAME) APPLE:BANANA:GRAPES LIST(DATA_TYPE) 300 (TTL)
PUSH test_list MANGO
POP test_list  # OR specify direction and count: POP test_list L (left) | R (right) 1 (number of elements)
```

