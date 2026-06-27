# Concurrency

`go_app_template` is an HTTP service. Concurrency surfaces in the HTTP server (each request runs in its own goroutine) and in any background workers added later. These rules exist to keep the race detector clean and the shutdown path predictable.

## Mutex hygiene

### One concern per mutex

Don't use a single `sync.Mutex` to guard everything a struct owns. Group related fields under their own mutex:

- One mutex for shared **request state** (e.g. a rate-limit counter).
- One mutex for any **in-memory registry** (e.g. a cache added in a future enhancement).

If two operations need to read two separately-guarded fields atomically, that is a sign the design is wrong — not a sign to merge the mutexes.

### Document lock order

When more than one mutex is acquired, document the order in a comment at the top of the file:

```go
// Lock order: p.cursorMu → p.statsMu.
// Never acquire a child's lock while holding a parent's lock.
```

### Prefer `RWMutex` for read-heavy state

Use `sync.RWMutex` where reads dominate. Use `RLock` for reads and `Lock` for writes. Do not upgrade an `RLock` to a `Lock` — release first, then re-acquire.

### Defer the unlock

```go
m.mu.Lock()
defer m.mu.Unlock()
// ... work ...
```

If you can't `defer` because you need to release before a blocking call, extract the critical section into a helper function.

## Channels

### Bounded channels for fan-out

If the publisher ever fans out to multiple sinks, use bounded channels per sink. A slow or unavailable sink must not back-pressure the poll loop.

```go
select {
case sink.ch <- event:
    // delivered
default:
    // sink is full — log and skip; do not block the poller
}
```

### Close ownership

A channel is closed by exactly one goroutine: the *writer*. Readers detect closure via the two-value receive (`v, ok := <-ch`). Closing a channel twice panics; closing from the reader side is wrong.

### Don't leak goroutines

Every goroutine started by a manager must have a clear termination path:

- It returns when its context is cancelled (`ctx.Done()`).
- Or it returns when an owned channel is closed.
- Or its lifetime is bounded by the function that spawned it (the function waits for it before returning).

The manager's `Shutdown(ctx)` method cancels the root context and waits for all goroutines to finish via `sync.WaitGroup`.

## Context

### First parameter, not on a struct

Pass `context.Context` as the first argument to any function that can block or do I/O. Do not store it on a struct field.

```go
// Good
func (p *Poller) Run(ctx context.Context) error

// Avoid
type Poller struct { ctx context.Context }
```

Exception: the manager's own *root* context (set up in `New<Manager>`, cancelled in `Shutdown`) is stored on the manager. Background goroutines launched at construction time use this root context, not a request or caller context.

### Honour cancellation

Any function that can block must check the context:

```go
select {
case <-ctx.Done():
    return ctx.Err()
case event, ok := <-events:
    if !ok {
        return nil
    }
    // ...
}
```

Poll loops must check `ctx.Done()` on every iteration.

## Timers and tickers

Use `time.NewTicker` for any poll loop or recurring background work. Always stop the ticker in a `defer` to avoid a goroutine leak:

```go
ticker := time.NewTicker(p.interval)
defer ticker.Stop()
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-ticker.C:
        if err := p.poll(ctx); err != nil {
            slog.Error("poll cycle failed", "err", err)
            // continue — don't return; log and retry next tick
        }
    }
}
```

## Shutdown order

When graceful shutdown is implemented (see `agent-harness/quality/debt.md`), `cmd/main.go` should orchestrate:

1. **Stop accepting new requests.** Call `srv.Shutdown(ctx)` to drain in-flight HTTP requests.
2. **Close the DB connection.** Releases the Postgres connection pool after in-flight queries complete.
3. **Stop the OTEL tracer provider.** Flushes buffered spans.

Budget: configurable via context timeout; 30 seconds is a reasonable default.

## Race detector

`make test-race` is required for any change touching:

- Any new background goroutine or shared mutable state.
- Any code path that crosses goroutines.

A race detector failure is a P0 finding even if the test "happens to pass" in non-race mode.

## Don't

- Don't start a goroutine without a clear termination path.
- Don't `go func() { ... }()` inside a function called from the poll loop unless it is genuinely fire-and-forget and you've documented why.
- Don't hold a lock across a DB query or external network call. Release before blocking I/O.
- Don't use `sync.Once` to lazily initialise a manager. Construct eagerly in `cmd/`.
- Don't `panic` inside a goroutine without recovering. An un-recovered panic in a background goroutine takes down the whole process.
- Don't use `time.Sleep` for synchronisation in tests. Poll with a bounded timeout or signal via a channel.
