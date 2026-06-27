# Design Principles

Principles that guide design trade-offs in `go_app_template`. Consult when choosing between approaches or evaluating whether code is well-structured.

## Code quality values

- **Readability over cleverness.** A future maintainer with a stack trace and 10 minutes should be able to find the bug.
- **Boring Go.** Use standard idioms (`if err != nil`, `(T, error)` returns, short receivers, package-scoped types). Don't import generics, code-generation tricks, or reflection unless there is no readable alternative.
- **Composition over interfaces.** This service uses very few interfaces. Most components are concrete and wired at the composition root. Introduce an interface only when there is a real second implementation in sight (e.g. swapping `UserPostgresPersistence` for a fake in tests).
- **Easy to change.** When choosing between two approaches, pick the one that lets the *next* change be smaller.

## Small functions, one level of abstraction

Functions should do one thing. Statements within a function should read at the same level of abstraction.

```go
// Good — each step at one level
func (s *UsersService) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
    u.SetDefaults()
    u.Sanitize()
    if err := u.Validate(); err != nil {
        return nil, fmt.Errorf("validate user: %w", err)
    }
    if err := s.persistence.Create(ctx, u); err != nil {
        return nil, fmt.Errorf("persist user: %w", err)
    }
    return u, nil
}

// Avoid — mixed levels
func (s *UsersService) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
    if u.Email == "" {
        return nil, errors.New("email is required")
    }
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    _, err := s.pool.Exec(ctx, "INSERT INTO users (email, ...) VALUES ($1, ...)", u.Email)
    if err != nil { return nil, err }
    return u, nil
}
```

## Guard clauses

Replace nested conditionals with early returns. Flatten the happy path.

```go
// Good
func (h *HTTP) AddUser(w http.ResponseWriter, r *http.Request) {
    var body api_gen.NewUser
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    user, err := h.usersService.CreateUser(r.Context(), apiNewUserToDomain(&body))
    if err != nil {
        writeError(w, http.StatusUnprocessableEntity, err.Error())
        return
    }
    writeJSON(w, http.StatusCreated, domainUserToAPI(user))
}

// Avoid
func (h *HTTP) AddUser(w http.ResponseWriter, r *http.Request) {
    var body api_gen.NewUser
    if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
        if user, err := h.usersService.CreateUser(r.Context(), apiNewUserToDomain(&body)); err == nil {
            writeJSON(w, http.StatusCreated, domainUserToAPI(user))
        } else {
            writeError(w, http.StatusUnprocessableEntity, err.Error())
        }
    } else {
        writeError(w, http.StatusBadRequest, "invalid request body")
    }
}
```

## YAGNI — no speculative generality

Don't add a parameter, interface, or struct field "in case we need it someday". Add abstractions when there is a concrete second use. Three similar lines is better than a premature shared helper.

## Comments

Default: no comment. Only add one when the *why* is non-obvious:

- A hidden constraint (`// cursor must advance only after confirmed publish to guarantee at-least-once delivery`).
- A subtle invariant (`// caller holds p.cursorMu`).
- A workaround for a specific external limitation.

Don't comment what the code obviously does. Don't reference the current PR or task — that rots.

## Encapsulation

- Managers and pollers expose methods, not fields. Internal state (cursor, mutexes, channels) is unexported.
- Public methods are the contract; tests exercise them. If a test needs to reach into private fields, the method surface is wrong.

## Single responsibility

Each component owns exactly one concern:

- `users.UsersService` — orchestrates business rules, delegates to persistence, validates inputs.
- `users/persistence.UserPostgresPersistence` — executes queries against Postgres, wraps `pgx` errors.
- `transport/http/v1.HTTP` — parses HTTP requests, calls service methods, serialises HTTP responses.
- `pkg/datastore` — constructs and validates the pgxpool connection pool.

If a handler builds SQL, or a service layer parses HTTP headers, something is wrong.

## Dependency direction

Inward only. `cmd/` → domain packages → cross-cutting (`types`, `observability`, `logging`). See [dependency-rules.md](../repo-architecture/dependency-rules.md).

When you find yourself adding an import that violates the matrix, restructure rather than route around it.

## Concurrency is part of the design

The question "who owns this lock / goroutine / channel" must always have a clear answer. See [concurrency.md](concurrency.md).

A goroutine bolted onto a function call is almost always wrong. The right place is a component with a clear lifetime tied to the root context.

## HTTP API correctness is the top constraint

This service is the API surface that callers depend on. Design decisions must not break the API contract:

- A bug in the persistence layer must surface as a 500, not silently corrupt data.
- A new field in the domain must be reflected in the OpenAPI spec (`api/contracts/v1/api-specs.yaml`) before it is deployed.
- Breaking changes (removed fields, changed types) require a new API version (`/api/v2`), not an in-place change.

When in doubt about a design choice, ask: "does this break existing callers or violate the spec?" If yes, the design is wrong.

## Reuse vs duplication

Three similar lines is better than a premature shared helper. Wait until the third use is clear, then extract.

## Avoid frameworks for things stdlib does well

- `log/slog` — structured logging. No replacement needed.
- `context` — cancellation and deadlines. No replacement needed.
- `encoding/json` — JSON. No replacement needed.
- `time.NewTicker` — polling intervals. No replacement needed.

A new dependency should justify itself against the stdlib equivalent.

## Don't bypass safety

- Never hand-edit generated files (`internal/api_server_gen/v1/spec.gen.go`). Regenerate via `make generate-api-specs`.
- Persistence errors must surface to the caller — don't log-and-swallow inside `UserPostgresPersistence`.
- Race tests must pass. Don't skip a race-failing test "until later".
