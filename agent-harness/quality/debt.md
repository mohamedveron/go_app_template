# Technical Debt

Known technical debt, prioritised by impact. Cross-references ADRs where relevant.

> **Note:** Specific counts are intentionally omitted — they go stale silently and create false precision. Each entry uses a qualitative tier instead. Re-survey when prioritising work.
>
> **Tiers:** `widespread` — affects the majority of files in the area; `scattered` — affects multiple files but not the majority; `isolated` — one or a handful of specific files.

## High Priority

### Architectural-boundary enforcement not yet wired up

- **Tier:** widespread
- **Last surveyed:** 2026-06-26
- **Description:** `make lint` runs `golangci-lint` with its defaults, which catches Go-level issues but **not** the import matrix defined in [`agent-harness/knowledge/repo-architecture/dependency-rules.md`](../knowledge/repo-architecture/dependency-rules.md). Cross-layer violations (a handler reaching into persistence directly, `pkg/datastore` imported outside `cmd/` or `users/persistence/`) must be caught by review. A `depguard` config or a custom analyzer could enforce this at build time.
- **Action when ready:** Add `.golangci.yaml` at the repo root with `depguard` rules encoding the matrix. Validate by intentionally introducing a violation and confirming the lint failure. Update [`agent-harness/enforcement/golint/README.md`](../enforcement/golint/README.md) to point at the live config.

### Database credentials hardcoded in configs

- **Tier:** isolated (`internal/configs/configs.go`)
- **Last surveyed:** 2026-06-26
- **Description:** `Datastore()` returns a hardcoded host, port, username, and password. Safe for local development but must not go to production as-is.
- **Action when ready:** Read `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` from environment variables. Add validation in `New()` that returns an error if required vars are missing.

## Medium Priority

### `FindUserByID` endpoint not implemented

- **Tier:** isolated (`internal/transport/http/v1/user.go:68`)
- **Last surveyed:** 2026-06-26
- **Description:** `GET /api/v1/users/{id}` returns 501. The OpenAPI spec declares the endpoint but there is no service method or persistence query backing it.
- **Action when ready:** Add `ReadByID(ctx, id) (*domain.User, error)` to `UsersPersistence`, implement in `user_postgres.go`, add service method in `users.go`, implement handler in `user.go`.

### No graceful SIGTERM shutdown

- **Tier:** isolated (`cmd/main.go`)
- **Last surveyed:** 2026-06-26
- **Description:** The HTTP server starts with `ListenAndServe` but there is no signal handling. SIGTERM causes an immediate kill; in-flight requests are dropped.
- **Action when ready:** Wrap `main` with `signal.NotifyContext`, call `srv.Shutdown(ctx)` on context cancellation with a configurable drain timeout (e.g. 30s).

### Auth middleware uses hardcoded tokens

- **Tier:** isolated (`internal/transport/http/middleware/auth.go`)
- **Last surveyed:** 2026-06-26
- **Description:** `RequireAuth()` validates against a hardcoded list (`admin-token`, `readonly-token`). This is a testing convenience, not production-ready auth.
- **Action when ready:** Replace with JWT validation or mTLS. The `AuthMiddleware` struct is already injectable — swap the validation logic without changing the middleware interface.

## Low Priority

### MongoDB persistence is a stub

- **Tier:** isolated (`internal/users/persistence/user_mongo.go`)
- **Last surveyed:** 2026-06-26
- **Description:** `UserMongoPersistence` exists but all methods return `errors.New("not implemented")`. Not wired in `cmd/main.go`.
- **Action when ready:** Implement if MongoDB support is needed. The `UsersPersistence` interface is already satisfied; only the methods need bodies.

## Resolved

(none — section reserved for future use)
