# Naming Conventions

Naming rules for Go code in `go_app_template`. Go has strong conventions of its own; this document layers service-specific guidance on top.

## General principles

- **Intention-revealing names.** A reader should understand what a name represents without reading its body. `usersPersistence`, not `up`. `service.CreateUser(ctx, u)`, not `service.Do(ctx, obj)`.
- **Consistent vocabulary.** One word per concept. This service uses: *user*, *cursor*, *service*, *persistence*, *handler*, *middleware*, *server*, *domain*. Don't introduce synonyms (`repo`, `store`, `dao`, `manager` for what is called `persistence`).
- **Domain language wins over generic CS terms** in exported APIs. Inside a package, technical terms (`mutex`, `channel`, `pool`) are fine.

## Packages

- **Lowercase, single-word, no underscores, no camelCase.** `users`, `domain`, `persistence`, `configs`, `datastore`, `logger`, `observability`, `logging`, `proxy`, `middleware`.
- A package named `<x>` exports a `New<X>(...)` constructor for its primary type when appropriate (`users.NewService`, `persistence.NewUserPostgresPersistence`).
- Avoid stutter in the exported name (`users.UsersService` is the current name for historical reasons, but prefer `users.Service` for new domains).

## Files

- **Lowercase, hyphen-free, `_` only for `_test.go` and platform suffixes.** `user.go`, `user_test.go`, `user_postgres.go`.
- One responsibility per file when feasible:
  - `service.go` — the service struct and constructor.
  - `<entity>.go` — the entity's business methods (e.g. `users.go`).
  - `<entity>.go` in `domain/` — the domain type (e.g. `user.go`).
  - `persistence.go` — persistence interface definition (lives in the **consumer** package, e.g. `internal/users/persistence.go`, not in `internal/users/persistence/`).
  - `<entity>_<driver>.go` — persistence implementation (e.g. `user_postgres.go`).
  - `errors.go` — typed error values.

## Types

### Structs

- **PascalCase** for exported, **camelCase** for unexported.

```go
// Exported — used by other packages
type User struct {
    FirstName string
    LastName  string
    Email     string
    CreatedAt *pgtype.Timestamptz
}

// Unexported — internal to the package
type writeResponseWriter struct {
    http.ResponseWriter
    status int
}
```

### Interfaces

- **PascalCase** for exported, **camelCase** for unexported. No `I` prefix, no `Impl` suffix on implementors.
- Define interfaces in the **consumer** package, not the implementation package. The consumer owns the abstraction; the implementor satisfies it implicitly.
- Add an interface only when there is a real abstraction need — to allow swapping a concrete implementation (e.g. swapping Postgres for MongoDB in tests).

```go
// Good — defined in package users (the consumer), not in package persistence (the implementor)
// package users
type usersPersistence interface {
    Create(ctx context.Context, u *domain.User) error
    ReadByEmail(ctx context.Context, email string) (*domain.User, error)
    List(ctx context.Context, cursor domain.Cursor) (*domain.UserPage, error)
}

// Bad — defined in the implementation package
// package persistence
// type UsersPersistence interface { ... }

// Bad — IUsersPersistence, UsersPersistenceInterface
```

### Constants and enums

- **Numeric constants**: `UPPER_SNAKE_CASE` only for cross-package constants; otherwise PascalCase (`DefaultPort`).
- **String sentinel values** (e.g. permission levels in auth middleware): plain untyped string constants.

### Errors

- **Sentinel errors**: `ErrXxx` exported variables.
  ```go
  var ErrUserExists = errors.New("user already exists")
  ```
- **Typed errors**: structs ending in `Error` with `.Error() string`.
  ```go
  type ValidationError struct{ Field string; Reason string }
  func (e *ValidationError) Error() string { return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Reason) }
  ```
- Check with `errors.Is` / `errors.As` at boundaries; never compare error strings.

### Functions and methods

- **PascalCase** for exported, **camelCase** for unexported.
- Methods use verbs that describe what the component does: `service.CreateUser(ctx, u)`, `persistence.List(ctx, cursor)`, `handler.AddUser(w, r)`.

### Variables

- **camelCase** for locals and unexported package-level vars.
- Receiver names: short, consistent across all methods of a type. `us` for `*UsersService`, `p` for `*UserPostgresPersistence`, `h` for `*HTTP`, `s` for `*Server`. Not `self`, not `this`, not the full type name.

```go
func (us *UsersService) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) { ... }
func (p *UserPostgresPersistence) Create(ctx context.Context, u *domain.User) error { ... }
func (h *HTTP) AddUser(w http.ResponseWriter, r *http.Request) { ... }
```

### Acronyms

- Treat as a single word in identifiers: `UserID`, not `UserId`; `userID`, not `userId` (unexported).

## Imports

Standard Go convention:

```go
import (
    // stdlib
    "context"
    "encoding/json"
    "net/http"

    // third-party
    "github.com/go-chi/chi/v5"
    "github.com/Masterminds/squirrel"

    // module-internal
    "github.com/mohamedveron/go_app_template/internal/users/domain"
    "github.com/mohamedveron/go_app_template/internal/pkg/datastore"
)
```

`goimports` enforces the grouping. The module path is `github.com/mohamedveron/go_app_template`.

## Forbidden patterns

| Don't | Do |
|-------|----|
| `repo`, `store`, `dao` for the persistence layer | `persistence` |
| `manager` for a service that orchestrates business logic | `service` or `Service` |
| `IXxx` interface prefix | plain PascalCase |
| `XxxImpl` for the only implementation | name after what it does (e.g. `UserPostgresPersistence`) |
| Underscores in package names | single-word name |
| `package main` outside `cmd/main.go` | one main package per binary |
| Hand-editing `spec.gen.go` | regenerate with `make generate-api-specs` |
