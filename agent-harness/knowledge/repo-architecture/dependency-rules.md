# Dependency Rules

This document defines the allowed import relationships between `go_app_template`'s layers and packages. Every new file and every refactor must comply.

For the broader layout, see [overview.md](./overview.md).

## Layer Vocabulary

`go_app_template` is a **clean-architecture HTTP service**. Layers are:

- **Entry point (`cmd/main.go`)** — the process root. Constructs config, database pool, auth middleware, persistence, service, and HTTP server; wires them together; starts listening.
- **Users domain (`internal/users/`)** — the core business logic layer. Owns the `UsersService` struct, methods (`CreateUser`, `ListUsers`, `ReadByEmail`), and the `UsersPersistence` interface. Does **not** own the database driver or HTTP logic directly.
  - **`internal/users/domain/`** — pure domain types: `User`, `UserPage`, `Cursor`. No external dependencies.
  - **`internal/users/persistence/`** — interface definition + PostgreSQL and MongoDB implementations. Imports `pgx` and `squirrel`; isolated here so a storage swap touches only this sub-package.
- **Configuration (`internal/configs/`)** — reads env vars, returns typed config structs. May import `internal/pkg/datastore` (for `datastore.Config`) and `internal/pkg/logger`.
- **Datastore (`internal/pkg/datastore/`)** — infrastructure factories: `NewPostgresService` (pgxpool), `NewMongoService` (mongo). Only `internal/users/persistence/` and `cmd/` may import this.
- **Logger (`internal/pkg/logger/`)** — global zap sugared logger. May be imported by any package.
- **Observability (`internal/observability/`)** — OpenTelemetry init and slog integration. Imports `internal/logging/` for the pretty formatter.
- **Logging (`internal/logging/`)** — pretty ANSI formatter for slog. A leaf; imports nothing internal.
- **Proxy (`internal/proxy/`)** — thin wrappers for external AI services (OpenAI). No imports from domain or persistence layers.
- **Transport (`internal/transport/http/`)** — HTTP root server, versioned sub-servers, generated API stubs, middleware, and test utilities.
  - **`internal/transport/http/server.go`** — root chi router; imports `users`, `middleware`, `v1`.
  - **`internal/transport/http/v1/`** — V1 sub-router and handlers; imports `users`, `users/domain`, `api_server_gen/v1`, `middleware`.
  - **`internal/transport/http/middleware/`** — Auth, CORS, logging. May import `pkg/logger`.
  - **`internal/transport/http/api_server_gen/v1/`** — auto-generated; no manual imports should be added here.
  - **`internal/transport/http/testutils/`** — test helpers; imports `middleware`.
- **Legacy API (`internal/api/`)** — thin wrapper used for health helpers. Not part of the primary HTTP transport. Does not add new features here.

## Top-Level Import Matrix

Rows import from columns. ✓ = allowed. ✗ = forbidden.

| From ↓ / To →           | `cmd` | `users` | `users/domain` | `users/persistence` | `configs` | `pkg/datastore` | `pkg/logger` | `observability` | `logging` | `proxy` | `transport/http` | `transport/http/v1` | `transport/http/middleware` | `api_server_gen/v1` | `api` |
|--------------------------|-------|---------|----------------|---------------------|-----------|-----------------|--------------|-----------------|-----------|---------|------------------|---------------------|-----------------------------|---------------------|-------|
| `cmd`                    | ✓     | ✓       | ✓              | ✓                   | ✓         | ✓               | ✓            | ✓               | ✓         | ✓       | ✓                | ✗                   | ✓                           | ✗                   | ✓     |
| `users`                  | ✗     | ✓       | ✓              | ✓                   | ✗         | ✗               | ✓            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `users/domain`           | ✗     | ✗       | ✓              | ✗                   | ✗         | ✗               | ✗            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `users/persistence`      | ✗     | ✓       | ✓              | ✓                   | ✗         | ✓               | ✓            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `configs`                | ✗     | ✗       | ✗              | ✗                   | ✓         | ✓               | ✓            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `pkg/datastore`          | ✗     | ✗       | ✗              | ✗                   | ✗         | ✓               | ✗            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `pkg/logger`             | ✗     | ✗       | ✗              | ✗                   | ✗         | ✗               | ✓            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `observability`          | ✗     | ✗       | ✗              | ✗                   | ✗         | ✗               | ✗            | ✓               | ✓         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `logging`                | ✗     | ✗       | ✗              | ✗                   | ✗         | ✗               | ✗            | ✗               | ✓         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `proxy`                  | ✗     | ✗       | ✗              | ✗                   | ✗         | ✗               | ✓            | ✓               | ✗         | ✓       | ✗                | ✗                   | ✗                           | ✗                   | ✗     |
| `transport/http`         | ✗     | ✓       | ✗              | ✗                   | ✗         | ✗               | ✓            | ✓               | ✗         | ✗       | ✓                | ✓                   | ✓                           | ✗                   | ✗     |
| `transport/http/v1`      | ✗     | ✓       | ✓              | ✗                   | ✗         | ✗               | ✓            | ✗               | ✗         | ✗       | ✗                | ✓                   | ✓                           | ✓                   | ✗     |
| `transport/http/middleware` | ✗  | ✗       | ✗              | ✗                   | ✗         | ✗               | ✓            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✓                           | ✗                   | ✗     |
| `api_server_gen/v1`      | ✗     | ✗       | ✗              | ✗                   | ✗         | ✗               | ✗            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✓                   | ✗     |
| `api`                    | ✗     | ✓       | ✓              | ✗                   | ✗         | ✗               | ✓            | ✗               | ✗         | ✗       | ✗                | ✗                   | ✗                           | ✗                   | ✓     |

Key rules in words:

1. **`cmd/` is the composition root.** It constructs all top-level objects — Postgres pool, auth middleware, persistence, service, and HTTP server — and is the only entry point that wires them together.
2. **`users` owns interfaces, not implementations.** The `usersPersistence` interface lives in `internal/users/persistence.go` (consumer-side, package `users`). Concrete implementations (`UserPostgresPersistence`, `UserMongoPersistence`) live in `internal/users/persistence/`; `cmd/` injects them.
3. **`users/domain` is a pure leaf.** `User`, `UserPage`, `Cursor` and their validation methods have no internal imports.
4. **`users/persistence` is the only layer that imports `pkg/datastore`.** `pgxpool.Pool`, `squirrel`, and `mongo.Client` must not leak into `users`, `cmd/`, or `transport/http`.
5. **`transport/http/v1` does not import `users/persistence` directly.** Handlers talk to `*users.UsersService`, not to the DB layer.
6. **`api_server_gen/v1` is auto-generated — do not add manual imports.** Regenerate it with `make generate-api-specs`.
7. **`pkg/datastore` is a leaf infrastructure package.** It imports no internal packages. Only `cmd/` and `users/persistence/` construct datastore clients.
8. **`pkg/logger` is universally importable.** Any package may use the global zap logger. It is a leaf with no internal imports.
9. **`logging` is a pure leaf.** `PrettyFormatter` and helpers have no internal imports. Only `observability` imports it.
10. **`proxy` is isolated.** External AI client wrappers do not import domain or persistence layers.

## Composition Root Rules

Top-level construction happens in `cmd/main.go`. When you add a new component:

1. Define its interface in the **consumer** package (e.g. `internal/users/persistence.go` for the `users` package).
2. Implement it in the appropriate sub-package (e.g. `internal/users/persistence/user_postgres.go`).
3. Construct it in `cmd/main.go`, inject its dependencies from already-constructed peers.
4. Update the import matrix above.

## Third-Party Dependencies

| Library | Where it may appear |
|---------|-------------------|
| `github.com/jackc/pgx/v5` | `internal/pkg/datastore/` and `internal/users/persistence/` only |
| `github.com/Masterminds/squirrel` | `internal/users/persistence/` only |
| `go.mongodb.org/mongo-driver` | `internal/pkg/datastore/` and `internal/users/persistence/user_mongo.go` only |
| `github.com/go-chi/chi/v5` | `internal/transport/http/` only |
| `github.com/go-chi/cors` | `internal/transport/http/middleware/` only |
| `github.com/sashabaranov/go-openai` | `internal/proxy/` only |
| `go.opentelemetry.io/otel/*` | `internal/observability/` and any package needing spans |
| `go.uber.org/zap` | `internal/pkg/logger/` only |
| `log/slog` | anywhere (via `internal/observability/` helpers) |
| `github.com/google/uuid` | anywhere |
| `github.com/stretchr/testify` | test files only |

The swap test: *if this library were replaced, would the change ripple through multiple packages?* If yes, it belongs only in the package that wraps it. `pgx`, `squirrel`, `mongo-driver`, and `go-openai` are classic examples — they must not leak beyond their wrapping packages.

## File-Naming Quick Reference

| Role | File pattern | Example |
|------|-------------|---------|
| Process entry point | `main.go` in `cmd/` | `cmd/main.go` |
| Domain types | `<entity>.go` in `domain/` | `internal/users/domain/user.go` |
| Domain service | `service.go` | `internal/users/service.go` |
| Domain business logic | `<entity>.go` co-located with service | `internal/users/users.go` |
| Persistence interface | `persistence.go` | `internal/users/persistence.go` (package `users`) |
| Persistence implementation | `<entity>_<driver>.go` | `internal/users/persistence/user_postgres.go` |
| HTTP handlers | `<entity>.go` in `v1/` | `internal/transport/http/v1/user.go` |
| Middleware | `<concern>.go` in `middleware/` | `internal/transport/http/middleware/auth.go` |
| Datastore factory | `<driver>.go` in `pkg/datastore/` | `internal/pkg/datastore/postgres.go` |
| Config | `configs.go` | `internal/configs/configs.go` |
| SQL DDL | `<table>.sql` in `schemas/` | `internal/schemas/users.sql` |
| Test files | `*_test.go` co-located with subject | `internal/users/users_test.go` |
| Generated files | `*.gen.go` | `internal/transport/http/api_server_gen/v1/spec.gen.go` |

## Non-Negotiables

1. **Never hand-edit `spec.gen.go`.** Regenerate with `make generate-api-specs` after changing the OpenAPI spec.
2. **`pkg/datastore` wraps the drivers** — nothing else imports `pgx` or `mongo-driver` directly (except `users/persistence` which wraps the pool).
3. **HTTP handlers must not import `pkg/datastore`.** Handler code talks to `*users.UsersService`, not to the DB.
4. **Graceful shutdown.** Close the DB pool (`defer pgPool.Close()`) before the process exits.
5. **Race tests must pass.** `make test` is non-negotiable for any change touching the service or persistence layers.

## How These Rules Are Enforced

| Mechanism | What it guards |
|-----------|---------------|
| `gofmt -s` | Formatting, simplifications (`make format`) |
| `golangci-lint` | Lint rules, errcheck, ineffective assigns (`make lint`) |
| `go vet` | Static analysis |
| `go test` | Correctness (`make test`) |
| `make generate-api-specs` | Keeps `spec.gen.go` in sync with the OpenAPI spec |

Import boundary rules are not yet lint-enforced. A `depguard` configuration under `agent-harness/enforcement/golint/` is planned to enforce the matrix above programmatically.
