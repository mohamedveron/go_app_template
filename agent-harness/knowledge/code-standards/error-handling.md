# Error Handling

Go error-handling conventions for `go_app_template`.

## Core principles

### Errors are values

Functions that can fail return `(T, error)`. The error is checked at the call site. There is no try/catch, no panic-based control flow.

### Wrap with context

When propagating an error up through layers, wrap it with `fmt.Errorf("...: %w", err)` so the cause is preserved for `errors.Is` / `errors.As`.

```go
// Good
if err := s.persistence.Create(ctx, u); err != nil {
    return nil, fmt.Errorf("persist user: %w", err)
}

// Avoid — discards the cause
if err := s.persistence.Create(ctx, u); err != nil {
    return nil, errors.New("create failed")
}
```

### Don't return both a value and an error

When `err != nil`, the value should be the zero value. Callers must not try to use a partial result.

### Don't ignore errors

Every `err` must be checked. `_` is only acceptable when the error truly cannot inform behaviour — and even then, leave a comment explaining why.

```go
// Good — closing on shutdown; nothing useful to do with the error
_ = rows.Close() // best-effort: already shutting down

// Avoid — silent
pub.Publish(ctx, event) // returned error ignored
```

### Define typed errors at boundaries

Inside a package, use sentinel errors or typed errors. At cross-package boundaries, callers use `errors.Is` / `errors.As` to distinguish them.

```go
// internal/users/persistence/errors.go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
)

type PersistenceError struct{ Cause error }
func (e *PersistenceError) Error() string { return "persistence: " + e.Cause.Error() }
func (e *PersistenceError) Unwrap() error { return e.Cause }
```

## Layer-specific rules

### `internal/users/` (service layer)

- Validate and sanitize inputs before calling persistence. Return a descriptive error (not a 500) for invalid inputs so the handler can map it to a 400/422.
- Log with full context at the point of first observation; don't re-log the same error as it bubbles up.

```go
// In the service — validate before persisting
if err := u.Validate(); err != nil {
    return nil, fmt.Errorf("validate user: %w", err)
}
```

### `internal/transport/http/v1/` (HTTP handlers)

- Map domain errors to HTTP status codes at the handler boundary. Do not let raw `pgx` errors or internal messages reach the client.
- Use `writeError(w, status, msg)` consistently; never write error detail that could leak implementation internals.

```go
if err != nil {
    writeError(w, http.StatusUnprocessableEntity, err.Error())
    return
}
```

### `internal/users/persistence/`

- Return typed errors for known conditions (`ErrUserExists`, `ErrUserNotFound`, connection errors). Wrap raw `pgx` errors before returning so the persistence package is the only place that imports `pgx`.

### `cmd/main.go`

- Startup errors (DB connection refused, missing env var) cause the process to exit non-zero with a clear message. Use `log.Fatal` or return a non-zero exit code directly — this is the composition root, not a library.
- Shutdown errors should be logged but must not prevent the remaining shutdown steps from running.

## Retry policy

HTTP handlers do not retry — they fail fast and return an appropriate status code. If transient errors (DB timeout, connection loss) need to be handled, surface them as 503 and let the client retry. Do not add retry loops inside service or persistence methods unless there is a well-defined transient error class with exponential backoff and jitter.

## Panic policy

- **Never panic for control flow.** Errors are values.
- Panics are reserved for programmer errors (impossible states, contract violations).
- An un-recovered panic in a background goroutine takes down the whole process. Always recover in long-lived goroutines:

```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            slog.Error("background worker panic", "recover", r)
        }
    }()
    s.runLoop(ctx)
}()
```

## Logging

- Use `slog.Default()` configured by `internal/observability/InitLogger`.
- Use `pkg/logger` (zap-based) for structured logging in request handlers.
- Levels:
  - `slog.Debug` — verbose details useful only for debugging (query parameters, cursor values).
  - `slog.Info` — significant lifecycle events (server started, database connected).
  - `slog.Warn` — recoverable issues (slow query, retrying a transient error).
  - `slog.Error` — errors the operator needs to know about.
- Include structured context as key/value pairs:
  ```go
  slog.Error("create user", "err", err, "email", u.Email)
  ```
- **Never log credentials, connection strings, or secrets.**
- **Never swallow an error by logging it without propagating** when the caller needs to act. Logging is an *additional* action, not a substitute for returning the error.

## Context

Pass `context.Context` as the first argument to any function that may block or do I/O. Honour cancellation — return promptly when `ctx.Done()` fires.

```go
func (p *UserPostgresPersistence) Create(ctx context.Context, u *domain.User) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    // ... execute query ...
}
```

## Don't

- Don't `fmt.Errorf("%s", err.Error())` — that drops the wrap chain. Use `%w`.
- Don't compare errors by string (`err.Error() == "not found"`). Use sentinel or typed errors with `errors.Is` / `errors.As`.
- Don't suppress errors during shutdown — log them. A silent shutdown failure is harder to diagnose.
- Don't return `nil` for "no result"; return `(nil, ErrUserNotFound)` so the caller can map it to a 404 rather than treat it as a success.
