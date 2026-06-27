# Exec Plans

This directory holds in-progress and completed coding plans produced by the Planning Agent. Each plan lives in its own subdirectory.

## Naming Convention

All plan directories follow this pattern:

```
NNN-<key>-short-desc/
```

| Segment | Rules |
|---|---|
| `NNN` | Zero-padded three-digit sequence number. Determines creation order. Inspect existing directories and increment from the highest number. |
| `<key>` | Either `issueNNN` matching the GitHub issue number (e.g. `issue125`), or `prNNN` matching the originating PR if no issue exists, or `NOTICKET` if neither applies. |
| `short-desc` | A brief, lowercase, hyphen-separated description of the plan's scope. Keep it under 5 words. |

Examples:

- `007-issue125-exec-priority-flag/`
- `008-pr132-encoding-field-readback/`
- `009-NOTICKET-fix-auth-redirect/`

## Race Policy

Two planners working in parallel may both claim the same `NNN`. This is expected and harmless. The collision is resolved at PR-merge time:

- The **second PR to be merged** must rename its directory with `git mv` before merging.
- Inspect `main` to find the next available number after the first plan merged.
- The renamed plan keeps its original `<key>` and description — only the sequence number changes.

Do not block planning work to avoid a potential collision. Resolve it at merge time.

## Directory Contents

Each plan directory typically contains:

| File | Description |
|---|---|
| `plan.md` | The finalised coding plan approved by the human. |
| `conversation-summaries.md` | Key decisions and rationale from the planning session. |
| `housekeeping_audit.md` | Created by the Quality Agent after merge — deferred or dismissed PR comments. |

Additional files (screenshots, reference images, sketches) may also be present.

## Lifecycle

1. **Created** by the Planning Agent during `001_plan.md`.
2. **Read** by the Coding Agent during `002_build.md` and by the Fixer Agent during `004_apply_fixes.md`.
3. **Archived** by the Quality Agent during `005_quality.md` — moved to `_archive/` with `git mv` after the PR is merged.

## Archive

Completed plans live in `_archive/`. The directory name and contents are preserved — only the location changes.

```bash
git mv harness/plans/007-issue125-exec-priority \
       harness/plans/_archive/007-issue125-exec-priority
```

Do not delete exec-plan directories. The archive provides a searchable history of past decisions.
