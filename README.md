# Go Message Broker

An in-memory message broker built to learn Go concurrency: goroutines, channels, mutexes, WaitGroups, and HTTP servers.

## Structure

- `cmd/server` — the broker's HTTP server
- `cmd/loadtest` — concurrent load-testing client
- `internal/broker` — core Broker type: topics, queues, thread-safe Publish/Consume/Count

## Run the server

    go run ./cmd/server

## Endpoints

    POST /publish?topic=orders&message=hello
    GET  /consume?topic=orders
    GET  /count?topic=orders

## Run the load test (server must be running)

    go run ./cmd/loadtest

Results are written to loadtest.log.
