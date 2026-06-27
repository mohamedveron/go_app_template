# Run Checks

Run the local verification checks to confirm the code compiles and passes. Referenced by the build stage (`002_build.md`) and the fix stage (`004_apply_fixes.md`).

---

## Check Suite

Five checks total. Four are **required for every change**; one is **conditional**.

| # | Check | Command | When required |
|---|-------|---------|---------------|
| 1 | Format | `make format` | Always. |
| 2 | Lint | `make lint` | Always. |
| 3 | Tests | `make test` | Always. |
| 4 | Build (sanity) | `make build` | Always before pushing. |
| 5 | API spec sync | `make generate-api-specs` | When `api/contracts/v1/api-specs.yaml` or any referenced YAML changes. Commit the regenerated `spec.gen.go`. |

---

### 1. Format

```bash
make format
```

Runs `gofmt -s -w .`. After running, the working tree should still match what you intend to commit — if `make format` produced changes, stage them.

---

### 2. Lint

```bash
make lint
```

Runs `golangci-lint run`. Must exit zero with no findings.

Common failures:

- `ineffassign` — assigned a value that's never read. Usually a missing `:=` vs `=`, or a leftover from refactoring.
- `errcheck` — ignored an error without `_`. Either check it or assign to `_` with a comment.
- `gosimple`, `staticcheck` — simplifications. Apply them.
- `govet` — printf args, lock copies, struct tag mistakes. Treat as P0.

Don't suppress findings with `//nolint` unless you also add a comment explaining the exception.

---

### 3. Tests

```bash
make test
```

Runs `go test -v ./...`. Must exit zero with all tests passing.

Targeted runs for faster feedback:

```bash
go test ./internal/users/...
go test ./internal/transport/http/...
go test ./internal/transport/http/v1/ -run TestAddUser
```

---

### 4. Build (sanity check)

```bash
make build
```

Builds `bin/app`. Catches compilation errors `go test` won't reach (e.g. an unused import in a file with no tests, mismatches between `cmd/` wiring and constructors). Cheap — run before pushing.

---

### 5. API spec sync (conditional)

```bash
make generate-api-specs
```

Must be run whenever `api/contracts/v1/api-specs.yaml` or any referenced schema/resource YAML changes. Regenerates `internal/transport/http/api_server_gen/v1/spec.gen.go`. Commit the updated generated file alongside the spec change.

**If `spec.gen.go` is out of sync with the spec, the build will fail or handlers will silently diverge from the declared API surface.**

---

## Interpreting Failures

### Compile failures

Read the error verbatim. Common issues:

- Unused import or variable — remove it.
- Type mismatch — fix the type at the source; don't paper over with `any`.
- `ServerInterface` method missing — you updated the spec but didn't implement the new method. Add a stub that returns 501 until it's implemented.
- Nil pointer at test time — usually a missing constructor argument when wiring a component in a test.

### Test failures

If a test fails that you didn't touch:

- Did your change alter shared behaviour (e.g., cursor semantics used across tests)?
- Is the test flaky (look at `git log` for the test file)?
- Is there an environmental dependency (Postgres not running)?

If the failure is unrelated to your change, flag it in the PR — don't paper it over.

---

## Rules

- The four always-required checks (`format`, `lint`, `test`, `build`) must pass before pushing.
- Run `make generate-api-specs` and commit the result whenever the OpenAPI spec changes.
- Do not suppress lint or vet findings with `//nolint` without a comment.
- Do not skip or comment out failing tests to make the suite pass.
