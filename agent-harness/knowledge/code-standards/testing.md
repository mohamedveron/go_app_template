# Testing

How `go_app_template` is tested and the conventions any new test should follow.

## Philosophy

- **Tests exercise real components, not mocks by default.** Use a real Postgres instance (via `testcontainers-go` or a local Docker Compose) for persistence and DB tests. External service interactions (e.g. AI proxies) use interfaces with in-memory fakes to avoid network setup, but the fakes must be minimal and documented.
- **Each test creates isolated state.** No shared global DB rows between tests. Use per-test table truncation; seed data in `t.Cleanup`.
- **Race detection is mandatory.** Every package that crosses goroutines must pass `go test -race`. `make test-race` is the gate.
- **Tests are deterministic.** No `time.Sleep` for synchronisation. Use channels, `sync.WaitGroup`, or polling with a bounded timeout.

## File layout

- Tests live next to the code they test: `internal/users/users.go` ↔ `internal/users/users_test.go`.
- Package-internal tests use `package <name>`; black-box tests use `package <name>_test`. Match the file you are editing.
- Shared test helpers belong in `*_test.go` files within the relevant package until there is a genuine need to share across packages; HTTP-level helpers go in `internal/transport/http/testutils/`.

## Running tests

```bash
# All tests
make test

# Targeted
go test ./internal/users/...
go test ./internal/transport/http/...

# A single test
go test ./internal/transport/http/v1/ -run TestAddUser

# Verbose
go test -v ./internal/users/...
```

## Service-level tests

Test `users.UsersService` with a fake `UsersPersistence` implementation:

```go
func TestUsersService_CreateUser(t *testing.T) {
    fake := &fakePersistence{}
    svc, err := users.NewService(fake)
    require.NoError(t, err)

    u := &domain.User{FirstName: "Alice", Email: "alice@example.com"}
    created, err := svc.CreateUser(context.Background(), u)
    require.NoError(t, err)
    require.Equal(t, "Alice", created.FirstName)
}
```

## Persistence tests (real Postgres)

Test `UserPostgresPersistence` directly against a real Postgres instance. Do not mock the database at this layer:

- `Create` → fetches back the row to confirm fields were written correctly.
- `List` with cursor → returns only rows after the given `After` value.
- `ReadByEmail` → returns the correct row or an error for unknown emails.

## Handler-level tests

Use `testutils.SetupTestEnvironment(t, opts)` to create a full HTTP test server:

```go
func TestAddUser_Success(t *testing.T) {
    env := testutils.SetupTestEnvironment(t, testutils.Options{})
    defer env.Cleanup()

    body := strings.NewReader(`{"firstName":"Alice","email":"alice@example.com"}`)
    resp := env.MakeRequestWithToken("POST", "/api/v1/users", body, "admin-token")

    require.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

## Don't

- Don't `time.Sleep` for synchronisation. Poll or signal.
- Don't mock the DB in persistence tests — use a real Postgres instance. Mocked DB tests passed while real migrations failed before; the lesson is applied here.
- Don't share DB state (rows) across tests. Each test owns its data.
- Don't skip a failing test to make the suite pass. If you must, file it as tech debt in [`harness/quality/debt.md`](../../quality/debt.md) and link the issue.
- Don't add a test that passes without `-race` but races otherwise — that is a P0 finding.
