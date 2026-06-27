# Add a New API Endpoint

Add a new REST endpoint to `go_app_template`. This service is spec-driven — the OpenAPI spec is the source of truth. New endpoints start in the YAML spec, not in Go code.

If your work crosses into another area mid-implementation, stop and re-run `agent-harness/enforcement/utils/list-harness.sh skills knowledge` for that area before continuing.

---

## Prerequisites

Read these before starting:

- [`agent-harness/knowledge/repo-architecture/overview.md`](../../knowledge/repo-architecture/overview.md)
- [`agent-harness/knowledge/repo-architecture/dependency-rules.md`](../../knowledge/repo-architecture/dependency-rules.md)
- [`agent-harness/knowledge/code-standards/api-conventions.md`](../../knowledge/code-standards/api-conventions.md)
- [`agent-harness/knowledge/code-standards/error-handling.md`](../../knowledge/code-standards/error-handling.md)
- [`agent-harness/knowledge/domain/events/language.md`](../../knowledge/domain/events/language.md) — users domain vocabulary

---

## Step 1 — Update the OpenAPI Spec

All HTTP API changes start in `api/contracts/v1/api-specs.yaml` (or the referenced sub-files in `api/contracts/v1/resources/` and `api/contracts/v1/schemas/`).

1. Add the new path and operation to the spec.
2. Define request/response schemas in `api/contracts/v1/schemas/` if they don't already exist.
3. Add parameter definitions if needed.
4. Ensure the operation has an `operationId` (used by `oapi-codegen` to name the generated method).

---

## Step 2 — Regenerate the Server Code

```bash
make generate-api-specs
```

This runs `oapi-codegen` and overwrites `internal/transport/http/api_server_gen/v1/spec.gen.go`. The regenerated file will include your new method on the `ServerInterface` interface.

**Never edit `spec.gen.go` by hand.** If you see a discrepancy, fix the spec and regenerate.

---

## Step 3 — Implement the Handler Method

The `HTTP` struct in `internal/transport/http/v1/` implements `ServerInterface`. Add a new file or extend an existing one (e.g. `user.go` for user-related endpoints).

Rules:
- The handler receives `http.ResponseWriter`, `*http.Request`, and any path/query params auto-extracted by the generated code.
- Translate request types (from `api_server_gen/v1`) to domain types (from `users/domain`), call the relevant service method, then translate the result back.
- Use `writeJSON(w, status, v)` and `writeError(w, status, msg)` helpers (already defined in `user.go`).
- Do not import `pkg/datastore` from the handler.

```go
// Example pattern
func (h *HTTP) GetUserByID(w http.ResponseWriter, r *http.Request, id int64) {
    u, err := h.users.ReadByID(r.Context(), id)
    if err != nil {
        writeError(w, http.StatusNotFound, "user not found")
        return
    }
    writeJSON(w, http.StatusOK, domainUserToAPI(u))
}
```

---

## Step 4 — Add Service Method (if needed)

If the endpoint requires new business logic, add it to `internal/users/users.go` as a method on `*UsersService`.

Rules:
- Call `us.persistence.<Method>(ctx, ...)` to interact with the database.
- Validate and sanitize inputs before calling persistence.
- Return domain types, not DB or API types.

---

## Step 5 — Add Persistence Method (if needed)

If the service method requires a new database query, add it to:
1. `internal/users/persistence/persistence.go` — add the method to the `UsersPersistence` interface.
2. `internal/users/persistence/user_postgres.go` — implement it on `UserPostgresPersistence`.
3. `internal/users/persistence/user_mongo.go` — add a stub that returns `errors.New("not implemented")`.

Rules:
- Use `squirrel` for query building; do not write raw SQL strings.
- Map DB errors to domain errors or sentinel errors at the package boundary.
- Fetch `limit + 1` rows for cursor-paginated queries; set `NextCursor` to the last row's sort key if more rows exist.
- Write a persistence test for every new query.

---

## Step 6 — Write Tests

Write tests before or alongside the implementation (TDD — see [`tdd-based-development.md`](tdd-based-development.md)).

At minimum:
1. **Happy path** — the endpoint returns the expected status and body.
2. **Not found / validation error** — the endpoint returns the correct 4xx.
3. **Auth failure** — request without a valid token returns 401.

See [`write-backend-test.md`](../testing/write-backend-test.md) for the test harness patterns.

---

## Step 7 — Run Checks

```bash
make format
make lint
make test
make build
```

All must pass.

---

## Step 8 — Update Documentation

If you introduced new domain terms, update [`agent-harness/knowledge/domain/events/language.md`](../../knowledge/domain/events/language.md).

If the change adds a new endpoint or changes request/response fields, the OpenAPI spec is the documentation — make sure it is accurate.

If the change requires a new ADR (new external system, breaking schema change), write it before merging. See [`agent-harness/skills/planning/write-adr.md`](../planning/write-adr.md).

---

## Checklist

Before considering the endpoint complete:

- [ ] OpenAPI spec updated in `api/contracts/v1/`.
- [ ] `make generate-api-specs` run; `spec.gen.go` committed.
- [ ] Handler method implemented in `internal/transport/http/v1/`.
- [ ] Handler uses `writeJSON` / `writeError` helpers; does not import `pkg/datastore`.
- [ ] Service method added to `internal/users/users.go` if new business logic needed.
- [ ] Persistence interface and implementations updated if new DB query needed.
- [ ] Tests cover happy path, error path, and auth failure.
- [ ] `make format`, `make lint`, `make test`, `make build` all pass.
- [ ] ADR written if the change introduces a new external dependency or breaks the wire contract.
