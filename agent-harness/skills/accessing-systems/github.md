# Accessing GitHub

Use the `gh` CLI for all GitHub operations (PRs, issues, branches, releases, code search) on go_app_template. If a GitHub MCP server is configured in your runtime, prefer it over the CLI for read operations.

## Repository

- **Owner:** `mohamedveron`
- **Repo:** `go_app_template`

| Artifact | URL |
|----------|-----|
| Repo | https://github.com/mohamedveron/go_app_template |
| Pull Requests | https://github.com/mohamedveron/go_app_template/pulls |
| Issues | https://github.com/mohamedveron/go_app_template/issues |
| OpenAPI spec | `api/contracts/v1/api-specs.yaml` (local) |

## Quick reference (`gh` CLI)

```bash
# List recent PRs
gh pr list

# Read a PR (diff, comments, checks)
gh pr view <number>
gh pr diff <number>
gh pr view <number> --comments

# Create a PR
gh pr create --title "..." --body "..."

# Issues
gh issue list
gh issue view <number>

# Branches and commits
gh api repos/mohamedveron/go_app_template/branches | jq -r '.[].name'
git log --oneline -20
```

## Conventions

- This repo uses **conventional commit style** PR titles (`feat:`, `fix:`, `chore:`, `refactor:`, `docs:`, `test:`). Match the recent history (`gh pr list --limit 20`).
- Reference issue numbers in PR bodies with `#NNN`.
- Draft PRs are used during the plan stage; mark `ready for review` when the build stage is complete.

## Reading PR review comments

When applying review fixes, fetch comments at full fidelity:

```bash
# Comment URLs are the canonical reference — they carry persona, severity, file/line, and the conversation thread
gh api repos/mohamedveron/go_app_template/pulls/<number>/comments | jq '.[] | {url: .html_url, body: .body, path: .path, line: .line}'

# Issue-style top-level comments (general PR comments, not file-attached)
gh api repos/mohamedveron/go_app_template/issues/<number>/comments
```

## CI checks

The CI runs:

- `make format` check
- `make lint`
- `make test`
- `make test-race`
- `make build`

A red check on the PR usually means one of those failed. `gh pr checks <number>` shows status; click through to the workflow log for the failure.
