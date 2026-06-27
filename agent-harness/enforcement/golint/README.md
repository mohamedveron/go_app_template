# Go Lint Enforcement

This directory is reserved for `golangci-lint` configuration and any auxiliary `go vet` or analysis tools that the harness wants to apply on the agent's behalf.

**Do not read these from the agent loop.** Linting is run via `make lint`; agents see the *output*, not the configuration. Configuration lives here so it can be reviewed, audited, and version-controlled separately from human-facing docs.

## Status

go_app_template does not yet ship a checked-in `.golangci.yaml` at the repo root. `make lint` currently runs `golangci-lint run` with the linter's defaults, which already catches:

- `govet` — printf args, lock copies, struct tag mistakes.
- `errcheck` — unchecked errors.
- `ineffassign` — assignments that are never read.
- `staticcheck`, `gosimple` — simplifications.
- `unused` — unreferenced symbols.

## Planned

When the team adopts a project-specific configuration, place it at `.golangci.yaml` at the repo root and reference it from this README. Configs likely to land:

- `depguard` rules that encode the import matrix from [`harness/knowledge/repo-architecture/dependency-rules.md`](../../knowledge/repo-architecture/dependency-rules.md), so cross-layer violations fail the build.
- `gocritic` for additional code-quality checks.
- `nilerr`, `nilnesserr` to catch `return nil, nil` style mistakes around typed errors.
- `goconst`, `gocyclo`, `funlen` if size/duplication metrics become useful.

## Why no `.golangci.yaml` here

Mirroring the repo-root config under `harness/enforcement/` would risk drift between the two copies. The single source of truth is the file CI and `make lint` actually consume. This directory documents intent; it does not host configuration.

If the team decides the harness should ship a sample config (e.g. for new contributors), drop a `golangci.example.yaml` here with a note that it must be copied to the repo root to take effect.
