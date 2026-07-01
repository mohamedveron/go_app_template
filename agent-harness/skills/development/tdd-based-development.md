# Test-Driven Development

Use TDD when implementing new behaviour in `go_app_template`. This skill describes how to write tests that protect behaviour, enable safe refactoring, and avoid brittleness from testing implementation details.

The principles here are drawn from Kent Beck's *Test Driven Development: By Example* and Ian Cooper's talk *TDD, Where Did It All Go Wrong*.

---

## Core Principle

**Test behaviours, not implementation details.**

The trigger for writing a new test is a *requirement you want to implement* — never a method you want to add to a struct. A test should express what the system does for its consumers, not how it does it internally.

- **Behaviour:** "When the user email is missing, `CreateUser` returns a validation error and nothing is written to the database."
- **Not behaviour:** "The `validateEmail` helper inside `UsersService` is called once."

The system under test is the **public API of a package** — its exported types and functions — not the unexported helpers behind it.

---

## When to Use TDD

Apply TDD when implementing **new behaviour**:

- A new service method or persistence query.
- A new API endpoint or domain field.
- A bug fix where a failing test reproduces the defect.

**Skip TDD** for:

- Pure wiring (constructing the service in `cmd/`, adding an import).
- Trivial getters with no logic.
- Refactoring that does not change behaviour — existing tests already cover you.

---

## The Red-Green-Refactor Cycle

### 1. Red — Write a Failing Test

Write a test that describes the behaviour you are about to implement. The test name should read like a requirement:

```go
func TestService_CursorDoesNotAdvanceOnPublishFailure(t *testing.T) { ... }
```

Not like an implementation detail:

```go
func TestService_ValidateEmailCalledInternally(t *testing.T) { ... }
```

**Run the test and confirm it fails.** A test you have never seen fail proves nothing.

### 2. Green — Make It Pass Quickly

Get the test to pass by the fastest route possible. Write direct, inline, ugly code. Hard-coded return values. The goal is to **understand how to solve the problem**, not to write production-quality code.

> *"For this brief moment, speed trumps design."* — Kent Beck

### 3. Refactor — Clean Up Without New Tests

Now improve the code. Extract helpers, remove duplication, rename, restructure.

**Critical rule: do not write new tests during the refactoring step.**

Your behaviour test from step 1 already covers you. If you extract an unexported helper, do not write a separate test for it — it's an implementation detail.

If during refactoring you discover a genuinely new behaviour (a new conditional, a new public-facing capability), that signals a new requirement. Stop, go back to Red, and write a new behaviour test.

---

## What to Test

### Test the Public Contract

The public contract of a package is its exported types, functions, and methods.

| Layer | Public contract |
|-------|----------------|
| `internal/users/` | `UsersService` — `NewService`, `CreateUser`, `ListUsers`, `ReadByEmail`. Interface: `UsersPersistence`. |
| `internal/users/domain/` | Plain data types — `User`, `UserPage`, `Cursor`. Exported methods: `Validate`, `Sanitize`, `SetDefaults`. |
| `internal/users/persistence/` | `UserPostgresPersistence` — `Create`, `List`, `ReadByEmail`. |
| `internal/transport/http/v1/` | HTTP handlers — `AddUser`, `ListUsers`, `FindUserByID`. Tested via `testutils.SetupTestEnvironment`. |
| `internal/pkg/datastore/` | Connection factory — `NewPostgresService`. Typically exercised via integration tests. |
| `cmd/` | No direct tests — exercise via component tests. |

### Do Not Test Internals

Implementation details are everything behind the public API: unexported helpers, internal data layouts, the specific sequence of calls between collaborators.

**Practical checks:**

- If you find yourself making a function exported solely to test it — stop. It's an implementation detail.
- If a test describes *how* the code works rather than *what* it achieves — rewrite it.

---

## Mocking Guidance

### When to Use a Fake Persistence Layer

The `usersPersistence` interface (defined in `internal/users/persistence.go`) is the main seam in this service. Use a fake (in-memory) implementation in `users.UsersService` tests to avoid requiring a database:

```go
type fakePersistence struct {
    users     []*domain.User
    createErr error
}

func (f *fakePersistence) Create(_ context.Context, u *domain.User) error {
    if f.createErr != nil {
        return f.createErr
    }
    f.users = append(f.users, u)
    return nil
}
```

A fake that returns errors on demand lets you test error paths cheaply without a real DB.

### Use Real Postgres for Persistence Tests

Use a real Postgres instance (via `testcontainers-go` or a local Docker Compose) for all `internal/users/persistence/` tests. Mocked DB tests can diverge from real schema migrations; use the real thing. See [`write-backend-test.md`](../testing/write-backend-test.md).

### When Not to Mock

- The Postgres persistence layer in persistence tests — use a real Postgres instance.
- `internal/users/domain/` — it's plain data, no mocking needed.
- The HTTP server in handler tests — use `testutils.SetupTestEnvironment`.

---

## Test Structure

Use the given-when-then pattern:

```go
func TestUsersService_CreateUser_RejectsInvalidEmail(t *testing.T) {
    // Given: a fake persistence layer and a user with no email
    fake := &fakePersistence{}
    svc, err := users.NewService(fake)
    require.NoError(t, err)
    u := &domain.User{FirstName: "Alice", Email: ""}

    // When: CreateUser is called
    _, err = svc.CreateUser(context.Background(), u)

    // Then: the error is surfaced and nothing was persisted
    require.Error(t, err)
    require.Len(t, fake.users, 0)
}
```

### Test Location

Tests live next to the code they test:

```
internal/users/users.go
internal/users/users_test.go

internal/users/persistence/user_postgres.go
internal/users/persistence/user_postgres_test.go
```

No separate `tests/` tree.

---

## Anti-Patterns to Avoid

| Anti-pattern | Why it hurts | What to do instead |
|---|---|---|
| Writing a test for every helper extracted during refactoring | Couples tests to implementation; refactoring breaks tests | Test behaviour through the public API only |
| Mocking the DB with an in-memory fake | Diverges from real Postgres behaviour; missed migration bugs | Use a real Postgres instance |
| Exporting an unexported function to test it | Breaks encapsulation | Test through the package's exported API |
| Test names that describe methods (`TestCreateUserCallsPersistence`) | Unreadable; doesn't document requirements | Name after behaviours (`TestUsersService_ReturnsErrorForInvalidEmail`) |
| Skipping the red step | Can't trust a test you've never seen fail | Always confirm it fails before implementing |
| `time.Sleep` to wait for state | Causes flake under load | Poll with a bounded timeout or use a channel signal |

---

## Reference

- Kent Beck — *Test Driven Development: By Example* (2002)
- Ian Cooper — [*TDD, Where Did It All Go Wrong*](https://www.youtube.com/watch?v=EZ05e7EMOLM)
- Martin Fowler — *Refactoring: Improving the Design of Existing Code*
