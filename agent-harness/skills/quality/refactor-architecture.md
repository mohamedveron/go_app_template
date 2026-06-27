# Architectural Refactoring

Systematically improve the structural organisation of `go_app_template` — layer boundaries, dependency direction, package placement, and architectural conformance — without changing external behaviour. This is the counterpart to [`refactor-code.md`](refactor-code.md), which handles code-level smells.

**Cardinal rule:** Architectural refactoring and feature work are two separate activities. Wear one hat at a time. If you discover missing behaviour during an architectural pass, stop, switch hats, add a failing test and implement the behaviour, then switch back.

---

## Overview

| Phase | Goal | Output | Requires human? |
|-------|------|--------|-----------------|
| **Phase 1 — Audit & Plan** | Scan for architectural violations, prioritise findings, propose concrete restructuring moves | An exec-plan document in `agent-harness/plans/` | Yes — human reviews and approves before Phase 2 |
| **Phase 2 — Implement** | Execute the approved restructuring in small, tested, committed steps | Clean architecture + commits | No — but human can review commits |

---

# Phase 1 — Audit & Plan

The goal of Phase 1 is *analysis only*. Do not move, rename, or delete any files. The output is a written plan a human can review and approve.

---

## Step 1.1 — Scope

Determine what code to audit. This can come from:

- A human-specified scope (e.g. "audit the users service", "check dependency direction").
- The code you are about to change (campsite rule — audit what's in the way).
- A full-codebase sweep (only when explicitly requested).

If no scope is given, ask the human before proceeding.

---

## Step 1.2 — Load Relevant Architecture Knowledge

Read the architecture knowledge files relevant to the scoped code:

1. [`agent-harness/knowledge/repo-architecture/overview.md`](../../knowledge/repo-architecture/overview.md)
2. [`agent-harness/knowledge/repo-architecture/dependency-rules.md`](../../knowledge/repo-architecture/dependency-rules.md)
3. [`agent-harness/knowledge/code-standards/concurrency.md`](../../knowledge/code-standards/concurrency.md) if the scope touches goroutines, channels, or mutexes.
4. [`agent-harness/knowledge/domain/events/language.md`](../../knowledge/domain/events/language.md) if users domain terminology is in scope.

---

## Step 1.3 — Detect Violations

Walk the scoped code looking for:

- **Import-matrix violations.** A package imports something it shouldn't per [`dependency-rules.md`](../../knowledge/repo-architecture/dependency-rules.md#top-level-import-matrix). Common: `transport/http/v1` importing `pkg/datastore` directly; `users` importing `transport`; `users/persistence` being bypassed by handlers.
- **Composition-root leakage.** A type constructed outside `cmd/main.go`. A package-level `var` that instantiates a service at init time. `os.Getenv` read from inside a service (env reads belong in `cmd/`).
- **SDK leakage.** The `pgx` package imported outside `internal/pkg/datastore/` and `internal/users/persistence/`. The `mongo-driver` imported outside those same packages.
- **Goroutine ownership ambiguity.** A goroutine launched from a one-shot function call with no clear termination path. Long-lived background work belongs on a component launched and shutdown in `cmd/`.
- **Mutex coalescing.** One mutex guarding multiple unrelated concerns — increases contention. See [`concurrency.md`](../../knowledge/code-standards/concurrency.md#mutex-hygiene).
- **Domain-language drift.** A new Go identifier that contradicts `events/language.md` (e.g. using `repo` instead of `persistence`, `store` instead of `persistence`, `manager` instead of `service`).
- **Cross-component coupling.** `transport/http` importing `users/persistence` directly (it should only import `users`).

---

## Step 1.4 — Prioritise

Score each finding:

1. **Is it in the way?** Does the violation make current or upcoming work harder?
2. **Is the area changing often?** Frequently-changed code benefits most from clean structure.
3. **Is it causing bugs?** Concurrency violations and cursor-advancement bugs are the most dangerous.
4. **Is it enforced by tooling?** Import-matrix rules are not yet lint-enforced (planned `depguard` config). Anything not caught by tooling needs more review attention.

Tiers:

- **High** — Actively causing pain, blocking work, or violating a non-negotiable. Fix in this pass.
- **Medium** — Would improve the codebase, not urgent. Include as optional.
- **Low** — Minor placement / naming issue. Note but don't propose action.

---

## Step 1.5 — Propose Restructuring Moves

For every High and Medium finding, propose a concrete move:

- **Move a file** — `git mv internal/<from>/<file>.go internal/<to>/<file>.go`.
- **Split a package** — extract a sub-package or pull related types into their own file.
- **Invert a dependency** — move an interface from where it's implemented to where it's used.
- **Move a composition step into `cmd/`** — pull a manager construction out of init code back into the composition root.

For each, describe:

- What gets moved / split / extracted.
- What imports change.
- What the dependency matrix looks like after.
- What risk exists (broken imports, test compilation failures, behavioural drift).

---

## Step 1.6 — Write the Exec-Plan

Create a directory at `agent-harness/plans/NNN-<key>-refactor-arch-<scope-slug>/` following the convention in [`agent-harness/plans/README.md`](../../plans/README.md). Write the plan at `plan.md` inside:

```markdown
# Architectural Refactoring Plan: <Scope>

Date: <YYYY-MM-DD>
Status: PENDING APPROVAL

## Scope
<What was audited and why.>

## Summary
<2–3 sentence overview of the structural health and the thrust of the moves.>

## Test Coverage
<Current state of test coverage for the scoped code. Note gaps that need characterisation tests before refactoring.>

## Proposed Moves

### M1 — <Short title>

- **Priority:** High | Medium
- **Violation:** <Rule name + harness/path/to/file.md>
- **Location:** <Files affected>
- **Move:** <git mv / split / extract / etc.>
- **After:** <Dependency matrix after the move; import diff>
- **Risk:** <Compile failures, test breakage, performance, mutex order>

### M2 — ...

## Items Noted but Not Proposed
<Low-priority findings, one line each.>

## Execution Order
<Numbered list of M-items. Earlier moves must not depend on later ones.>
```

---

## Step 1.7 — Present for Approval

Ask the human:

1. Anything to **remove** from the plan?
2. Anything to **add**?
3. Any **priority** changes?
4. Is the **execution order** correct?
5. Are **characterisation tests** needed anywhere before moves begin?

Don't proceed to Phase 2 until the human explicitly approves. Update `Status` to `APPROVED`.

---

# Phase 2 — Implement

Execute the approved exec-plan in the recorded order. Each M-item is its own commit.

---

## Prerequisites

1. **Confirm `make test` and `make test-race` pass.** Don't start refactoring red code.
2. **Commit (or stash) the current state.** You need a clean rollback point.
3. **Write characterisation tests** for any gaps flagged in the plan; commit them separately.

---

## Step 2.1 — Implement Each Move

```
1. Read the M-item from the exec-plan.
2. Plan the smallest mechanical step that achieves the move.
3. Apply the step.
4. Update imports across the codebase.
5. Run `make build` — compiles? Continue. Errors? Revert and reduce scope.
6. Run `make test` — green? Continue. Red? Revert the last step.
7. Run `make lint` — clean? Continue. Findings? Fix before proceeding.
8. Run `make test-race` if the move touches concurrency-sensitive code.
9. Commit with type `refactor`.
10. Update the exec-plan: mark M-item as DONE.
11. Move to the next M-item.
```

Never take a step you can't revert quickly. If `make test` goes red and you can't fix it within a minute, revert and take a smaller step.

### Commit Format

```
refactor: isolate pgx imports inside users/persistence

Move all pgx usage out of users service layer into users/persistence/.
Service now depends only on UsersPersistence interface.
No behavioural change.
```

---

## Step 2.2 — Final Verification

After all moves:

1. **Run the full check suite.** `make format && make lint && make test && make test-race`.
2. **Build.** `make build` produces a clean binary.
3. **Review naming.** Do package and file names still make sense given the new structure?
4. **Spot-check imports.** Walk the dependency matrix; ensure no new violations appeared.
5. **Update the exec-plan** to `COMPLETED`.

---

## Common Moves

| Symptom | Move |
|---------|------|
| A type constructed outside `cmd/` | Move construction to `cmd/main.go`; pass it in as a parameter. |
| `pgx` imported outside `users/persistence/` | Extract queries into `internal/users/persistence/`; expose typed methods on `UsersPersistence`. |
| Handler directly importing persistence | Remove the direct import; pass `*users.UsersService` to the handler instead. |
| Goroutine launched from a non-manager function | Move it to a dedicated component; tie its lifetime to a root context. |
| `os.Getenv` read from a service | Pull the read up to `cmd/main.go`; pass the value into the constructor. |
| One mutex protecting unrelated concerns | Split into per-concern mutexes; document the lock order. |

---

## Anti-patterns

| Anti-pattern | Why it hurts | Do instead |
|---|---|---|
| Bundling restructuring and feature work in one PR | Hard to review; bugs from one hide behind the other | Wear one hat at a time |
| Skipping `make test-race` after a concurrency-touching move | Subtle races slip through | Run it; for hot paths, `-count=50` |
| Suppressing lint with `//nolint` to make a move land | Hides the real violation | Restructure or accept the rule |
| Big-bang restructuring instead of M-item-at-a-time | Easy to lose track of what broke | One commit per M-item |
