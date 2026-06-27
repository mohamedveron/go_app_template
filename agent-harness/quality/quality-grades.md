# Quality Grades

Assessment of code quality by area. Grades are qualitative estimates based on codebase exploration. Re-survey on a quarterly cadence or after significant work in an area.

> **Note:** Counts are intentionally omitted — they go stale silently. Re-survey when prioritising improvements or after significant work in an area.

Last updated: 2026-06-17

## Grading Scale

| Grade | Meaning |
|-------|---------|
| A | Excellent — well-tested, clean patterns, race-clean, documented |
| B | Good — solid conventions, minor gaps |
| C | Adequate — functional but has noticeable issues |
| D | Below standard — significant gaps or inconsistencies |
| F | Poor — needs major rework |

## Areas

| Area | Grade | Notes |
|------|-------|-------|
| `internal/users/` | TBD | Core domain: service orchestrator, `UsersService`, `UsersPersistence` interface. |
| `internal/users/domain/` | TBD | `User`, `UserPage`, `Cursor` types. Validate/Sanitize/SetDefaults methods. |
| `internal/users/persistence/` | TBD | Postgres-backed `UserPostgresPersistence`. pgx + squirrel queries. |
| `internal/configs/` | TBD | Env-based config. DB credentials are currently hardcoded — known debt. |
| `internal/pkg/datastore/` | TBD | `pgxpool` Postgres factory and MongoDB client. |
| `internal/pkg/logger/` | TBD | Global zap sugared logger. |
| `internal/observability/` | TBD | OpenTelemetry tracer init + slog integration. |
| `internal/logging/` | TBD | Pretty ANSI slog formatter. |
| `internal/proxy/` | TBD | OpenAI proxy client wrapper. |
| `internal/transport/http/` | TBD | Root HTTP server: chi router, global middleware, health check, v1 mount. |
| `internal/transport/http/v1/` | TBD | Versioned handlers implementing oapi-codegen `ServerInterface`. |
| `internal/api_server_gen/v1/` | TBD | Auto-generated stubs from OpenAPI spec — never hand-edit. |
| `internal/schemas/` | TBD | Raw SQL DDL for `users` table. |
| `cmd/` | TBD | Composition root; wires all dependencies; shutdown order is critical. |
| `api/contracts/` | TBD | OpenAPI 3.0 spec — source of truth for HTTP routes. |

To grade an area, walk it against:

- [`agent-harness/knowledge/code-standards/_index.md`](../knowledge/code-standards/_index.md) — naming, error handling, concurrency, testing.
- [`agent-harness/knowledge/repo-architecture/dependency-rules.md`](../knowledge/repo-architecture/dependency-rules.md) — non-negotiables, import matrix.
- Test coverage (`go test -cover ./internal/<area>/...`).
- Race cleanliness (`go test -race -count=20 ./internal/<area>/...`).
- Open issues mentioning the area (`gh issue list`).
