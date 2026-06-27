# Infrastructure

## What go_app_template Is, Operationally

`go_app_template` is a single Go binary HTTP service. It exposes a versioned REST API (`/api/v1`) and a health endpoint. It has one primary external dependency: a PostgreSQL database.

One deployment mode:

- **`bin/app`** — connects to Postgres, starts the chi HTTP server on the configured port, and runs until signalled.

## Local Development

| Concern | How |
|---------|-----|
| Build | `make build` → `bin/app`. |
| Run | `make run` — requires a local Postgres instance (see `docker-compose.yml`). |
| Start DB | `docker-compose up -d` — starts PostgreSQL on the default port. |
| Verify | `make test`, `make lint`, `make format`. |
| Regenerate API | `make generate-api-specs` — rebuilds `spec.gen.go` from the OpenAPI YAML. |

## Runtime Dependencies

| Layer | Dependency |
|-------|-----------|
| Language | Go 1.26+ |
| Database | PostgreSQL (primary store for users) |
| Database (optional) | MongoDB (stub; not wired in production) |
| HTTP framework | `go-chi/chi/v5` |
| API spec | OpenAPI 3.0 via `oapi-codegen` |

## Configuration

The service is configured via environment variables:

| Variable | Purpose | Default |
|----------|---------|---------|
| `PORT` | HTTP server port | `9090` |
| `APP_NAME` | Application name (for logs/health) | `dsp` |
| `APP_VERSION` | Application version string | `v0.0.0` |
| `GOENV` | Environment name (`local`, `staging`, `production`) | `` |

Database connection is currently hardcoded in `internal/configs/configs.go` (`Datastore()`) and should be moved to environment variables before production deployment:

| Config field | Hardcoded default | Recommended env var |
|-------------|-------------------|---------------------|
| Host | `postgres` | `DB_HOST` |
| Port | `5432` | `DB_PORT` |
| StoreName | `go_app` | `DB_NAME` |
| Username | `root` | `DB_USER` |
| Password | `123321` | `DB_PASSWORD` |

## Observability

| Concern | Mechanism |
|---------|-----------|
| Structured logs | `log/slog` via `internal/observability/`; request logs via `LoggingMiddleware`. |
| Tracing | OpenTelemetry. Set `OTEL_EXPORTER_OTLP_ENDPOINT` to ship traces to a collector. OTLP HTTP + stdout exporters are both wired. |
| Metrics | Not yet shipped. |

Key log events to watch:

| Event | Meaning |
|-------|---------|
| `request` (INFO) | Successful HTTP request (2xx/3xx) with method, path, status, latency |
| `request` (WARN) | Client error (4xx) |
| `request` (ERROR) | Server error (5xx) |
| `starting server on :PORT` | Service starting |

## Deployment

`go_app_template` is designed to run as a container (`Dockerfile` provided) or bare process. It does not currently implement graceful SIGTERM shutdown — a future enhancement should add `signal.NotifyContext` and a configurable drain timeout.

**Health check:** `GET /health` returns JSON with `status`, `version`, and `api_supported_versions`. Use this for container readiness probes.

## CI

The CI workflow should run:

- `make format` check
- `make lint`
- `make test`
- `make build`

## Planned / Not Yet Implemented

- Graceful SIGTERM shutdown with drain timeout.
- Database credentials from environment variables (currently hardcoded).
- `FindUserByID` endpoint (`GET /api/v1/users/{id}`) — currently returns 501.
- MongoDB persistence implementation (currently a stub).
- Metrics export (Prometheus / OTEL Metrics).
