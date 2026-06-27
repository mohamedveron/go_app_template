# Code Standards — Quick Reference

Read the file(s) relevant to your current task. Do not front-load all files.

| I'm about to... | Read |
|---|---|
| Name a package, type, function, file, or organise imports | [`naming-conventions.md`](naming-conventions.md) |
| Add or change `cmd/main.go` (entry point, wiring, shutdown) | [`naming-conventions.md`](naming-conventions.md), [`error-handling.md`](error-handling.md) |
| Add or change the users service (`internal/users/`) | [`naming-conventions.md`](naming-conventions.md), [`error-handling.md`](error-handling.md) |
| Add or change the persistence layer (`internal/users/persistence/`) | [`naming-conventions.md`](naming-conventions.md), [`error-handling.md`](error-handling.md) |
| Add or change domain types (`internal/users/domain/`) | [`naming-conventions.md`](naming-conventions.md), [`api-conventions.md`](api-conventions.md) |
| Add or change HTTP handlers (`internal/transport/http/v1/`) | [`naming-conventions.md`](naming-conventions.md), [`api-conventions.md`](api-conventions.md), [`error-handling.md`](error-handling.md) |
| Add or change the DB or datastore client (`internal/pkg/datastore/`) | [`naming-conventions.md`](naming-conventions.md), [`error-handling.md`](error-handling.md) |
| Add a new API endpoint (OpenAPI → codegen → handler) | skill: [`harness/skills/development/add-api-endpoint.md`](../../skills/development/add-api-endpoint.md) |
| Touch goroutines, channels, mutexes, or context cancellation | [`concurrency.md`](concurrency.md) |
| Handle or propagate errors (return, wrap, log) | [`error-handling.md`](error-handling.md) |
| Write or update tests | [`testing.md`](testing.md), skill: [`harness/skills/testing/write-backend-test.md`](../../skills/testing/write-backend-test.md) |
| Add a new domain field or user attribute | [`naming-conventions.md`](naming-conventions.md) (Constants section), [`api-conventions.md`](api-conventions.md), check `internal/users/domain/user.go` and `api/contracts/v1/schemas/User.yaml` first |
| Add a third-party dependency | [`dependency-rules.md`](../repo-architecture/dependency-rules.md#third-party-dependencies) |
| Check import or layer rules | [`dependency-rules.md`](../repo-architecture/dependency-rules.md) |
| Understand the overall architecture | [`overview.md`](../repo-architecture/overview.md) |
| Make a design decision (when in doubt) | [`design-principles.md`](design-principles.md) |
| Create a git branch | [`branch-naming.md`](branch-naming.md) |
| Write a commit message | skill: [`harness/skills/development/commit-changes.md`](../../skills/development/commit-changes.md) |
