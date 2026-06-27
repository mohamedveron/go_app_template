# Commit

Draft a commit message, optionally get user approval, and commit.

This service uses **conventional commits**. Match the style of the recent log (`git log --oneline -10`).

---

## Step 1 — Draft Commit Message

1. Run `git diff --staged` and `git diff` to understand all uncommitted changes.
2. Run `git log --oneline -10` to see recent message style.
3. Compose a commit message in the format below.

### Message Format

```
<type>: <summary>

<optional body — what changed and why>

<optional bullets for non-trivial changes>
```

#### Summary line

- Form: `<type>: <imperative summary>`.
- Type is one of: `feat`, `fix`, `refactor`, `chore`, `docs`, `test`.
- Summary is a single short imperative phrase (no period). Lowercase after the colon, except proper nouns.
- Total line ≤ 72 characters when possible.

Examples:

```
feat: add FindUserByID endpoint
fix: return 422 when user email is invalid
refactor: extract domainUserToAPI into shared helper
chore(deps): bump jackc/pgx from 5.5.0 to 5.6.0
docs: document cursor-based pagination behaviour
test: add handler test for ListUsers with after cursor
```

#### Body (optional)

- One short paragraph if the *why* isn't obvious from the summary.
- Then bullets for concrete changes when there's more than one.
- Reference issue numbers (`#12`, `#34`) when relevant.

#### Format notes

- Don't use double-quotes (`"`) in the summary line.
- Wrap body lines at ~72 characters.
- Don't list every file — group by concept.

Present the drafted message as a single fenced text block.

---

## Step 2 — User Approval (local only)

If you are running **locally** (interactive session with a human):

- Present the drafted commit message to the user.
- Ask the user to approve or request changes.
- Do not proceed until the user explicitly approves.

If you are running in the **cloud** (automated, no human in the loop): skip this step and proceed directly to Step 3.

---

## Step 3 — Pre-commit checks

Before committing, ensure:

- `make format` produced no changes (or you've staged them).
- `make lint` passes.
- `make test` passes.
- `make test-race` passes if your change touches concurrency.

A pre-commit hook may run some of these; do not bypass it with `--no-verify`. If a hook fails, investigate and fix.

---

## Step 4 — Commit

Stage specific files rather than `git add -A` when possible.

```bash
git add internal/events/service.go internal/events/service_test.go
git commit -m "$(cat <<'EOF'
fix: do not advance cursor when kafka publish fails

Cursor now advances only after Publisher.Publish returns nil.
Previously a partial-success path could advance the cursor on
a retriable error, causing events to be silently dropped.
EOF
)"
```

Do NOT append `Authored by Cursor`, `Co-authored-by`, or any AI attribution to the commit message unless explicitly asked.
