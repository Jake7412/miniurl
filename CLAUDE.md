# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

All commands go through `make`. The Makefile downloads Go 1.20.3 into `.tools/go/` automatically.

```sh
make unit-test          # Run unit tests with coverage (output: target/cover.out)
make integration-test   # Integration tests (require Docker/Podman)
make benchmark          # Run benchmarks
make fuzz               # Run fuzz tests
make tidy               # go mod tidy
make clean              # Remove target/
```

Run a single test:
```sh
.tools/go/bin/go test ./api/... -run TestAddUrl
```

## Architecture

MiniURL is a URL shortener with three distinct layers:

**Core** (`miniurl.go`) — `Hash(url string) string` generates a short code via MD5. This is the only domain logic; everything else adapts around it.

**API** (`api/api.go`) — HTTP REST layer using `httprouter`. Exposes `POST /api/v1/url` (accepts `{"url":"..."}`, returns `{"url":"...","hash":"..."}`). Defines a `Handler` interface that the caller must implement:
```go
type Handler interface {
    AddUrl(url string) (hash string, err error)
}
```
`api.Bind(router, handler)` wires the handler to the router. The interface is defined here (consumer side), not in the core package — dependency inversion is intentional.

**UI** (`ui/index.html`) — Single static HTML page. Submits the form via `fetch` to `/api/v1/url` and renders the shortened link.

**Data flow:** UI form → `POST /api/v1/url` → `Handler.AddUrl()` → `miniurl.Hash()` → response with hash.

## Key Design Patterns

- **Interface defined by consumer**: `Handler` lives in `api/`, not in the core package — callers inject implementations.
- **External test packages**: Tests use `package miniurl_test` / `package api_test` to test the public surface only.
- **testcontainers** is pulled in for integration tests (Docker/Podman required).
- `docker-compose.yaml` provides a PostgreSQL 15 instance for local integration testing (not yet wired up in code).

## API Spec

`openapi.yaml` is the authoritative REST contract. The Go structs in `api/api.go` must stay in sync with it.
