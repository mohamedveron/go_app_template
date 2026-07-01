# go_app_template

A community Go application template following clean architecture and domain-driven design principles. The HTTP layer is built around **spec-driven development** — you define the API contract first (OpenAPI), generate the server scaffolding, then implement the business logic.

## agent-harness
the harness should stay lean and opinionated on project-specific rules (domain vocabulary, shutdown order, no DB mocks),

## Quick start

Requires Go >= 1.21.

```bash
make all    # generate + test + build + run
make test   # run tests only
make run    # start HTTP server on :9090
```

Docker:
```bash
docker build -t go_app .
docker run -p 9090:9090 go_app
```

---

## Table of contents

1. [Directory structure](#directory-structure)
2. [HTTP layer: spec-driven development](#http-layer-spec-driven-development)
   - [How a request flows](#how-a-request-flows)
   - [How versioning works](#how-versioning-works)
   - [Adding a new endpoint](#adding-a-new-endpoint)
   - [OpenAPI spec conventions](#openapi-spec-conventions)
3. [internal/configs](#internalconfigs)
4. [internal/api](#internalapi)
5. [internal/users](#internalusers)
6. [internal/pkg](#internalpkg)
7. [proxy](#proxy)
8. [docker](#docker)

---

## Directory structure

```
.
├── cmd/
│   └── main.go                          # entry point — wires deps, starts net/http.Server
│
├── api/
│   └── contracts/
│       ├── cfg.yaml                     # oapi-codegen config (chi-server: true)
│       ├── api-specs.yaml               # OpenAPI 3.0 spec — START HERE for new endpoints
│       ├── resources/                   # per-resource path definitions
│       └── schemas/                     # reusable schema definitions
│
├── internal/
│   ├── configs/
│   │   └── configs.go                   # config loading (env vars, defaults)
│   │
│   ├── api/
│   │   └── *.go                         # service-layer façade (shared by HTTP, gRPC, etc.)
│   │
│   ├── users/                           # example business-logic unit
│   │   ├── service.go
│   │   ├── users.go
│   │   ├── users_test.go
│   │   ├── domain/
│   │   └── persistence/
│   │
│   ├── pkg/
│   │   ├── datastore/                   # postgres + mongo initialisation
│   │   └── logger/
│   │
│   └── transport/
│       └── http/
│           ├── server.go                # root chi router — global middleware + version mounts
│           ├── server_test.go
│           ├── middleware/
│           │   ├── auth.go              # Bearer token auth (standard http.Handler)
│           │   ├── cors.go              # go-chi/cors
│           │   └── logging.go           # structured slog request logging
│           ├── testutils/               # shared test helpers
│           ├── api_server_gen/
│           │   └── v1/
│           │       └── spec.gen.go      # GENERATED — do not edit by hand
│           └── v1/
│               ├── server.go            # v1 chi sub-router — v1-scoped middleware
│               ├── user.go              # implements ServerInterface for /users endpoints
│               └── chat.go             # implements ServerInterface for /openai endpoints
│
├── observability/                       # OpenTelemetry setup
├── proxy/                               # third-party HTTP clients (OpenAI, etc.)
├── logging/
├── go.mod
└── Makefile
```

---

## HTTP layer: spec-driven development

All HTTP APIs in this template follow **spec-driven development**: the OpenAPI spec is the single source of truth. You write the spec first, generate the server interface, then fill in your implementation. The generated code is never edited by hand.

### How a request flows

```
Incoming request
      │
      ▼
cmd/main.go
  net/http.Server { Handler: server.GetRouter() }
      │
      ▼
internal/transport/http/server.go   ← root chi router
  Middleware (every request):
    1. chimw.Recoverer          — panic → 500
    2. otelhttp.NewMiddleware   — OpenTelemetry trace span
    3. LoggingMiddleware()      — structured slog: method/path/status/duration
      │
      ├── GET /health           ← no auth, returns version + supported API versions
      │
      └── Mount("/api/v1", ...) ← hands off to the v1 sub-router
              │
              ▼
          internal/transport/http/v1/server.go   ← v1 chi sub-router
            Middleware (v1 requests only):
              1. Recoverer
              2. CorsMiddleware()         — go-chi/cors
              3. authMiddleware.RequireAuth()   — 401 if no/invalid token
              │
              ├── GET  /openai/{topic}
              ├── POST /users
              └── GET  /users/{id}
                        │
                        ▼
                  api_server_gen/v1/spec.gen.go   ← GENERATED
                    Extracts + coerces path/query params (chi.URLParam),
                    calls into your handler struct
                        │
                        ▼
                  v1/user.go, v1/chat.go   ← your business logic
                    Implements ServerInterface methods
```

The full request URL for a v1 endpoint is `/api/v1/users/{id}`. When the request reaches the v1 sub-router via `r.Mount`, the `/api/v1` prefix is already stripped, so the sub-router and generated code only see `/users/{id}`.

### How versioning works

The root server mounts each version as an independent sub-router:

```go
// internal/transport/http/server.go
r.Mount("/api/v1", s.v1Server.GetRouter())
// future: r.Mount("/api/v2", s.v2Server.GetRouter())
```

Each version is a self-contained chi router with its own middleware stack. This means:

- v1 and v2 can have **different auth schemes, CORS rules, or rate limits**
- Adding v2 never touches v1 code
- The `/health` endpoint is intentionally **outside all versions** — it is a global liveness check with no auth

### Adding a new endpoint

Follow these steps every time:

**1. Define the endpoint in the OpenAPI spec**

Edit [api/contracts/api-specs.yaml](api/contracts/api-specs.yaml). Add your path, method, request/response schemas. Follow the schema conventions below.

**2. Regenerate the server code**

```bash
go generate ./...
```

This runs `oapi-codegen` using [api/contracts/cfg.yaml](api/contracts/cfg.yaml) (which sets `chi-server: true`) and overwrites [internal/transport/http/api_server_gen/v1/spec.gen.go](internal/transport/http/api_server_gen/v1/spec.gen.go).

**3. Implement the new method**

The compiler will now tell you that your handler struct no longer satisfies `ServerInterface`. Add the method to the appropriate file under [internal/transport/http/v1/](internal/transport/http/v1/) (or create a new file for a new resource group). The generated wrapper already handles parameter extraction — you receive typed values directly.

```go
// The generated interface (do not edit):
FindUserByID(w http.ResponseWriter, r *http.Request, id int64)

// Your implementation (in v1/user.go):
func (h *HTTP) FindUserByID(w http.ResponseWriter, r *http.Request, id int64) {
    // id is already parsed and typed — focus on business logic
}
```

**Never edit `spec.gen.go` directly.** Changes there will be overwritten on the next `go generate`.

### OpenAPI spec conventions

1. The main spec lives at [api/contracts/api-specs.yaml](api/contracts/api-specs.yaml).
2. Schemas go in [api/contracts/schemas/](api/contracts/schemas/) — one file per resource, indexed via `_index.yaml`. Reuse existing schemas before creating new ones.
3. Paths go in [api/contracts/resources/](api/contracts/resources/) — organised by resource (users, accounts, etc.).
4. All list responses wrap items in a `data` object and reuse a `meta` object for pagination attributes.
5. Name request schemas with a `Request` suffix; name response schemas after the resource itself.
6. All error responses use the shared `Error` schema (`code` + `message`).

---

## internal/configs

Centralises all configuration loading. HTTP port, timeouts, and datastore credentials are read here — either from env vars or defaults. Keeping this isolated means you can swap from hardcoded defaults to etcd or AWS Parameter Store without touching any other package.

## internal/api

A service-layer façade that all transports (HTTP, gRPC, CLI) call into. This guarantees consistent behaviour across transport types — no accidental divergence between what the HTTP handler does vs what a CLI command does.

## internal/users

An example business-logic unit. The pattern applies to any domain entity (orders, accounts, etc.):

- `service.go` — the `Users` struct, initialised via `NewService`, holds all dependencies
- `users.go` — business logic (validation, orchestration)
- `users_test.go` — unit tests for pure functions; API-level tests cover the full stack
- `domain/` — domain types
- `persistence/` — datastore interactions; the `Store` interface enables dependency inversion and makes unit testing without a real DB possible

## internal/pkg

Packages shared across multiple business-logic units (not specific to any one domain). Includes:

- **datastore** — initialises `pgxpool.Pool` (Postgres) and the Mongo driver
- **logger** — wraps `go.uber.org/zap`; define the interface here so you can swap logging libraries without touching business logic

## proxy

All third-party integrations — HTTP clients, SDKs, other protocols. Each integration lives in its own file. Currently includes an OpenAI client.

## docker

```bash
docker build -t go_app .
docker run -p 9090:9090 go_app
```

The Dockerfile produces a minimal image using a multi-stage build.

---

## Available routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | None | Liveness check, returns version info |
| `POST` | `/api/v1/users` | Required | Create a new user |
| `GET` | `/api/v1/users/{id}` | Required | Fetch a user by ID |
| `GET` | `/api/v1/openai/{topic}` | Required | Generate a paragraph on a topic via OpenAI |
