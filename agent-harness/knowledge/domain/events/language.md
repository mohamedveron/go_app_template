# Users Domain — Language

This file is the authoritative vocabulary for the `go_app_template` service's users domain. Use these exact terms in Go identifiers, log fields, JSON field names, and PR comments.

## Core Concepts

### User

The primary domain entity. Represents a person registered in the system.

- Go type: `domain.User` in `internal/users/domain/user.go`
- Fields: `FirstName`, `LastName`, `Mobile`, `Email`, `CreatedAt`, `UpdatedAt`
- API type: `api_gen.User` in `internal/transport/http/api_server_gen/v1/spec.gen.go` (auto-generated from OpenAPI spec)
- Do not confuse the domain type with the API/wire type — the handler converts between them via `domainUserToAPI()`.

### NewUser (creation request)

The payload sent to `POST /api/v1/users`. Contains the fields required to create a user.

- API type: `api_gen.NewUser` (auto-generated)
- Required fields: `firstName`, `email`
- Optional fields: `lastName`, `mobile`

### UserPage (paginated result)

The response type for `GET /api/v1/users`. Contains a slice of users and a cursor for the next page.

- Go type: `domain.UserPage` in `internal/users/domain/user.go`
- Fields: `Users []*User`, `NextCursor *string`
- API type: `api_gen.UserPage` (auto-generated)
- If `NextCursor` is nil, there are no more pages.

### Cursor (pagination state)

A pointer to the current position in the user list. Passed as a query parameter (`after`) on subsequent requests.

- Go type: `domain.Cursor` in `internal/users/domain/user.go`
- Fields: `After *string` (exclusive lower bound, the `createdAt` of the last seen row), `Limit uint`
- Wire parameter name: `after` (query string), `limit` (query string)
- Do not call it "offset" (that is SQL offset pagination) or "token" (ambiguous with auth tokens) or "page" (that is page-number pagination).
- The cursor value is opaque to the client — it is the `createdAt` timestamp of the last returned row encoded as a string.

## Validation Rules (from `domain/user.go`)

| Field | Rule |
|-------|------|
| `FirstName` | Required; trimmed of whitespace |
| `Email` | Required; trimmed; must match email format |
| `LastName` | Optional; trimmed |
| `Mobile` | Optional; trimmed |

`Validate()` returns an error if any required field is missing or malformed.
`Sanitize()` trims whitespace from all string fields.
`SetDefaults()` sets any zero-value fields to sensible defaults before persistence.

## API Field Names (JSON)

These are the wire names defined in `api/contracts/v1/schemas/User.yaml`. They are locked once the API is deployed. Do not rename without bumping the API version.

| Go field | JSON wire name |
|----------|---------------|
| `FirstName` | `firstName` |
| `LastName` | `lastName` |
| `Mobile` | `mobile` |
| `Email` | `email` |
| `CreatedAt` | `createdAt` |
| `UpdatedAt` | `updatedAt` |

## Endpoint Vocabulary

| HTTP method | Path | Operation name | Description |
|-------------|------|----------------|-------------|
| `GET` | `/api/v1/users` | `listUsers` | Paginated user listing |
| `POST` | `/api/v1/users` | `addUser` | Create a new user |
| `GET` | `/api/v1/users/{id}` | `findUserById` | Lookup by ID (stubbed) |

## What This Service is NOT

- **Not an identity provider.** It stores basic user records but does not manage sessions, passwords, or OAuth tokens.
- **Not a soft-delete service.** There is no `deleted_at` field or delete endpoint in the current implementation.
- **Not multi-tenant.** There is no workspace or organisation scoping in the current users table schema.
