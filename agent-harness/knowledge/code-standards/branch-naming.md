# Branch Naming Conventions

## Format

```
<type>/<short-description>
```

Use lowercase kebab-case for the description. Keep it concise but descriptive enough to identify the work at a glance.

## Branch Types

| Prefix | Use When |
|--------|----------|
| `feat/` | Adding a new feature or capability |
| `fix/` | Fixing a bug |
| `chore/` | Maintenance, dependency updates, tooling changes |
| `docs/` | Documentation-only changes |
| `refactor/` | Restructuring code without changing behaviour |
| `test/` | Adding or updating tests only |

## Examples

```
feat/find-user-by-id
feat/cursor-pagination
fix/list-users-empty-cursor
chore/bump-otel-1.41
docs/openapi-spec-conventions
refactor/extract-users-service
test/persistence-create-coverage
```

## Rules

1. **Always use a type prefix.** Bare names like `users-fix` are not allowed.
2. **Keep descriptions to 3–5 words.** Long branch names are hard to work with.
3. **Use kebab-case.** No underscores, camelCase, or spaces.
4. **No PR numbers in the branch name.** The PR description carries the link; the branch name describes the work.
5. **One branch per task.** Do not bundle unrelated changes on a single branch.
