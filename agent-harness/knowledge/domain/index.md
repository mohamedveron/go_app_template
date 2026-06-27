# Domain Knowledge — Index

This folder describes the `go_app_template` domain — the vocabulary and behaviour of an HTTP service that manages users via a spec-driven REST API.

The service has one primary domain: **users**. Get the terminology right here before writing code that surfaces in API responses or database columns — renaming a field in the OpenAPI spec is a breaking change for consumers.

## When to read what

1. **Always start here:** [context-map.md](context-map.md) — the data flow from HTTP client through the handler, service, and persistence layers.
2. **Working on users domain** (types, validation, field names, cursor pagination): the domain types are in `internal/users/domain/user.go`. Read it before adding new fields.

## Available domain files

| Domain | When to read |
|--------|--------------|
| [`context-map.md`](context-map.md) | System overview: which component calls which, what the auth middleware does, how pagination works. Read first. |

## Why getting terms right matters

The `User` type in `api/contracts/v1/schemas/User.yaml` is a contract with API consumers. Field names that appear in JSON responses — `firstName`, `lastName`, `email`, `createdAt` — are locked once the API goes to production. Renaming them in Go without updating the OpenAPI spec is a silent contract break.

If you introduce a new domain field or endpoint, update the OpenAPI spec first (`api/contracts/v1/api-specs.yaml`), regenerate the server code (`make generate-api-specs`), then implement.
