# API Conventions

`go_app_template` is a spec-driven HTTP service. Its external contract is the OpenAPI spec at `api/contracts/v1/api-specs.yaml`. This file documents conventions for request/response schemas and for the internal domain types that map to them.

## OpenAPI spec is the source of truth

The `api/contracts/v1/api-specs.yaml` (and its referenced sub-files) is the contract with API consumers. Handler stubs and request/response types are auto-generated from it via `make generate-api-specs`. Never add, rename, or remove fields in `spec.gen.go` directly.

## Domain types are internal representations

The internal `domain.User` struct lives in `internal/users/domain/user.go`. It is the service-internal representation and the contract between the service layer and the persistence layer.

```go
// internal/users/domain/user.go
type User struct {
    FirstName string
    LastName  string
    Mobile    string
    Email     string
    CreatedAt *pgtype.Timestamptz
    UpdatedAt *pgtype.Timestamptz
}
```

Rules:

- **Domain types do not import transport or generated types.** The handler (`v1/user.go`) converts between domain types and API types via `domainUserToAPI()`.
- **Required fields are enforced by `Validate()`.** Do not skip validation before persisting.
- **Optional fields are pointer-typed in the API schema.** Check for nil before dereferencing in handlers.

## JSON field names

JSON field names are defined in the OpenAPI schema (`api/contracts/v1/schemas/`). They are camelCase and must match the generated Go struct tags. Do not override struct tags on generated types.

| Go field | JSON wire name |
|----------|---------------|
| `FirstName` | `firstName` |
| `LastName` | `lastName` |
| `Mobile` | `mobile` |
| `Email` | `email` |
| `CreatedAt` | `createdAt` |
| `UpdatedAt` | `updatedAt` |

## Adding a new field

1. Add the field to the appropriate schema YAML (`api/contracts/v1/schemas/User.yaml` or a new schema file).
2. Run `make generate-api-specs` to regenerate `spec.gen.go`.
3. Add the field to `domain.User` in `internal/users/domain/user.go` if it needs to be persisted.
4. Update the persistence layer (`user_postgres.go`) to read/write the new column.
5. Update `domainUserToAPI()` in `v1/user.go` to map the field.
6. Update `internal/schemas/users.sql` if the DB schema changes.

## Adding a new endpoint

See [`agent-harness/skills/development/add-api-endpoint.md`](../../skills/development/add-api-endpoint.md).

## Error response format

All error responses use the `Error` schema:

```go
// api_gen.Error
type Error struct {
    Code    int32  `json:"code"`
    Message string `json:"message"`
}
```

Use `writeError(w, status, msg)` in handlers to produce a consistent error response. The `code` field should match the HTTP status code.

## Versioning

Breaking changes (removing a field, changing a type, renaming a field in a response) require a new API version (`/api/v2`). Additive changes (adding a new optional response field, adding a new endpoint) are backwards-compatible and do not require a version bump.
