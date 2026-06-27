# 003.2 — Review Local Changes

## What This Is

This is a set of instructions you must adopt and execute, not background context. When triggered, you become the Reviewer Agent. You will perform **five review passes**, each under a distinct persona. Follow the execution steps below.

**Subagent dispatch is strongly preferred.** If your runtime supports subagents, launch one subagent per reviewer persona so all five passes run in parallel. Each subagent receives the list of changed files, the full diff, and the persona-specific instructions. Aggregate findings from all subagents before presenting the final report. If subagents are unavailable, run the passes sequentially in the order listed.

## Purpose

Review locally staged and unstaged changes (not yet on a PR) across five focus areas — QA, Code Quality, Security, Architecture Conformance, and Domain Conformance — and present findings triaged by severity.

## Reviewer Personas

You execute all five personas. Each pass has a distinct focus and evaluation source:

| Persona | Focus | Evaluation Source |
|---------|-------|-------------------|
| **QA** | Functional correctness, edge cases, test coverage gaps, regression risks. | The diff itself and your understanding of the feature intent inferred from the changes. If a planning artifact exists in `agent-harness/plans/` on the current branch, use it; otherwise rely on the code context alone. |
| **Code Quality** | Code clarity, maintainability, naming, duplication, adherence to codebase conventions, SOLID principles. | The agent agent-harness in `agent-harness/`. Use `agent-harness/knowledge/code-standards/_index.md` as the canonical lookup table from "what the diff touches" to "which doc to read". Always consult `agent-harness/knowledge/repo-architecture/dependency-rules.md` § Non-Negotiables before flagging anything that looks like duplication or unnecessary indirection. |
| **Security** | Injection vectors, authentication/authorisation gaps, data exposure, dependency vulnerabilities, input validation. go_app_template-specific concerns: workspace path traversal, token leakage in logs, race conditions in auth middleware, command injection in exec args. | Your own internal knowledge of security best practices and common vulnerability patterns, plus `agent-harness/knowledge/domain/shared-language.md` § Tokens / Permissions. |
| **Architecture Conformance** | Violations of the dependency rules and structural conventions that the linter cannot catch. | Start at `agent-harness/knowledge/repo-architecture/overview.md` to orient, then focus on `dependency-rules.md § What is NOT enforced programmatically`. Walk every item in that section against the diff. |
| **Domain Conformance** | Correct use of go_app_template's domain language; no silent term renames; new types/functions/files named after the resource domain they belong to. | `agent-harness/knowledge/domain/<resource>/language.md` for each resource domain (`exec`, `files`, `tasks`, `ports`) touched by the diff, plus `agent-harness/knowledge/domain/shared-language.md`. |

## Evidence Rule

Every **Code Quality**, **Architecture Conformance**, and **Domain Conformance** finding MUST either:

(a) cite the specific agent-harness doc and the rule it contradicts — `<rule-name> — agent-harness/path/to/file.md`, or  
(b) state `"agent-harness silent — general principle: <reasoning>"`.

Findings that fail this rule must not be included.

## Context You Have Access To

### Coding Agent agent-harness and Internal Context

Read the `AGENTS.md` to understand the context and guidance you should draw from provided by the agent agent-harness.

### Planning Artifacts (optional)

If a planning artifact exists at `agent-harness/plans/` on the current branch, use it for additional context. There may not be one — this review mode does not require it.

### Local Changes

- The output of `git status` (to identify all changed, added, and deleted files).
- The output of `git diff` (unstaged changes) and `git diff --cached` (staged changes).
- The full codebase at HEAD.

## Execution Steps

### Step 1 — Gather the Diff

Run the following commands to understand what has changed:

```bash
git status
git diff HEAD
```

Use `git diff HEAD` to capture both staged and unstaged changes in a single diff. Read the diff in full. If a planning artifact exists in `agent-harness/plans/`, cross-reference against it; otherwise infer the intent from the changes themselves.

### Step 1.2 — Map the Diff to agent-harness Coverage (Code Quality persona)

For each meaningful change in the diff:

- Identify what the change touches using this repo's vocabulary
  (application service, DTO, controller, repository, gateway, shared kernel
  type, page route, hook, component, …).
- Look up the relevant agent-harness doc(s) using
  `agent-harness/knowledge/code-standards/_index.md` (the canonical lookup table)
  and `agent-harness/enforcement/utils/list-agent-harness.sh` for discovery.
- Read those docs in full before forming any Code Quality finding on that
  hunk.

Also run `agent-harness/enforcement/utils/list-agent-harness.sh knowledge` and check whether any recently added knowledge file (e.g. a new resource domain language file or a new infra doc) is relevant to the diff.

### Step 2 — Architecture Conformance Pass

Read `agent-harness/knowledge/repo-architecture/dependency-rules.md § What is NOT enforced programmatically`. For each item in that section, inspect the diff for violations:

- Are managers constructed only in `cmd/`, not at package init time or inside other managers?
- Are V1 handlers translation-only — no business logic, no goroutines, no manager-state mutation?
- Are mutexes split per concern (registry / buffer / clients / state) rather than coalesced into one big lock?
- Is the lock order documented and respected where multiple mutexes are involved?
- Are generated files (`*.gen.go`) untouched by hand?
- Does any handler bypass `internal/files/validation.go` to build a filesystem path?
- Is `os.Getenv` read only from `cmd/`, never from a manager?
- Do package names and file names follow `dependency-rules.md § File-Naming Quick Reference`?

Record each violation as a finding. Apply the Evidence Rule.

### Step 3 — Domain Conformance Pass

For each resource domain touched by the diff, read `agent-harness/knowledge/domain/<resource>/language.md`. Then inspect:

- Do new types, functions, and file names use the ubiquitous language of their resource domain?
- Are any domain terms silently renamed in code without a corresponding update to the language file?
- Are terms from one resource domain bleeding into another?

Record each violation as a finding. Apply the Evidence Rule.

### Step 4 — Identify Issues (all personas)

Within each persona's focus area, identify issues that are actionable and material.

### Step 5 — Triage Each Issue

Assign a severity to each finding across all five passes:

| Severity | Meaning | Expectation |
|----------|---------|-------------|
| **P0** | Blocking — must be fixed before merge. | Bug, security vulnerability, broken functionality, data loss risk. |
| **P1** | High — strongly recommended fix. | Significant quality issue, missing test for critical path, architectural concern. |
| **P2** | Medium — should be addressed. | Code clarity, minor test gap, non-critical convention violation. |
| **P3** | Low / nit — optional. | Style preference, minor naming suggestion, non-functional improvement. |

### Step 6 — Present Findings

List all findings grouped by persona, ordered by severity (P0 first). Each finding must include:

- The **persona** that raised it (`QA`, `Code Quality`, `Security`, `Architecture Conformance`, or `Domain Conformance`).
- The **file and line range** where the issue occurs.
- A clear **description** of the issue.
- The **citation** backing up the finding (as required by the Evidence Rule).
- The **severity tag** (`P0`, `P1`, `P2`, or `P3`) with a concise reasoning for the assigned severity level.
- A **suggested fix** or direction (where possible).

If there are no findings, state that the local changes passed review.

## Done

Your work is complete when all findings have been listed and triaged.
