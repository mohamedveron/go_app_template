# Write a Test

Conventions and patterns for writing a test for `go_app_template`. Read [`agent-harness/knowledge/code-standards/testing.md`](../../knowledge/code-standards/testing.md) first — this skill assumes you've internalised the philosophy there.

## When to TDD

For new feature development, follow [`skills/development/tdd-based-development.md`](../development/tdd-based-development.md). For bug fixes, write the failing test that reproduces the bug *first*, then fix.

---

## Step 1 — Identify the right layer

| Where you changed code | Where the test goes |
|------------------------|---------------------|
| `internal/users/` | `internal/users/users_test.go` |
| `internal/users/persistence/` | `internal/users/persistence/user_postgres_test.go` |
| `internal/transport/http/v1/` | `internal/transport/http/v1/user_test.go` |
| `internal/transport/http/` (root server) | `internal/transport/http/server_test.go` |
| `cmd/` | No direct test — exercise via component test. |

Prefer the layer closest to the change. A bug in user validation is a `users_test.go` test. A bug in JSON serialisation is a handler-level test.

---

## Step 2 — Spin up the right harness

### Handler-level tests (real HTTP server)

Use `testutils.SetupTestEnvironment(t, opts)` from `internal/transport/http/testutils/`, which creates an `httptest.Server` backed by the full chi router:

```go
func TestAddUser_Success(t *testing.T) {
    env := testutils.SetupTestEnvironment(t, testutils.Options{})
    defer env.Cleanup()

    body := `{"firstName":"Alice","email":"alice@example.com"}`
    resp := env.MakeRequestWithToken("POST", "/api/v1/users", strings.NewReader(body), "admin-token")

    require.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

### Service-level tests (mock persistence)

Define a fake `UsersPersistence` in the test file and inject it via `users.NewService(fakePersistence)`:

```go
type fakePersistence struct {
    users  []*domain.User
    createErr error
}

func (f *fakePersistence) Create(_ context.Context, u *domain.User) error {
    if f.createErr != nil {
        return f.createErr
    }
    f.users = append(f.users, u)
    return nil
}

func (f *fakePersistence) ReadByEmail(_ context.Context, email string) (*domain.User, error) {
    for _, u := range f.users {
        if u.Email == email {
            return u, nil
        }
    }
    return nil, errors.New("not found")
}

func (f *fakePersistence) List(_ context.Context, cursor domain.Cursor) (*domain.UserPage, error) {
    return &domain.UserPage{Users: f.users}, nil
}
```

### Persistence tests (real Postgres)

Use a real Postgres instance for persistence tests — do not mock the database at the persistence layer:

```go
func TestUserPostgresPersistence_Create(t *testing.T) {
    pool := testutils.NewTestDB(t) // starts Postgres via docker or uses local instance
    p, err := persistence.NewUserPostgresPersistence(pool)
    require.NoError(t, err)

    u := &domain.User{FirstName: "Alice", Email: "alice@example.com"}
    err = p.Create(context.Background(), u)
    require.NoError(t, err)
}
```

---

## Step 3 — Cover the right cases

For a new endpoint or feature, write at minimum:

1. **Happy path** — the feature works as specified.
2. **Validation failure** — invalid input returns a 400/422 with a meaningful message.
3. **Auth failure** — missing or invalid token returns 401.
4. **Not found** — looking up a non-existent resource returns 404.
5. **DB failure** — persistence returns an error; verify the handler returns 500 and does not panic.

For a bug fix:
1. A test that fails on the *unfixed* code and passes on the fixed code.

---

## Step 4 — Auth in tests

The auth middleware checks for hardcoded tokens. Use `testutils.MakeRequestWithToken` to pass the right token:

```go
// Admin token — can call write endpoints
resp := env.MakeRequestWithToken("POST", "/api/v1/users", body, "admin-token")

// Read-only token — can call read endpoints, rejected on writes
resp := env.MakeRequestWithToken("GET", "/api/v1/users", nil, "readonly-token")

// No token — triggers 401
resp := env.MakeRequest("GET", "/api/v1/users", nil)
```

---

## Step 5 — Run all checks

```bash
make format
make lint
make test
```

---

## Don't

- Don't mock the DB at the persistence layer — use a real Postgres instance.
- Don't `time.Sleep` for synchronisation.
- Don't share DB state (rows, schemas) across tests — each test owns its data.
- Don't skip a test with `t.Skip` to make CI green.
- Don't put shared test helpers inline in one test file — extend `internal/transport/http/testutils/` if they're reusable.
- Don't hand-edit `spec.gen.go` to make tests compile — fix the spec and regenerate.
