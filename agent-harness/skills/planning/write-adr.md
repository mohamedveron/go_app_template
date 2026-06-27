# Write ADR

Write an Architecture Decision Record (ADR) to capture a significant design decision and its rationale.

---

## Prerequisites

Read `harness/knowledge/architecture-decision-records/ABOUT.md` for the purpose and lifecycle of ADRs.

---

## When to Write an ADR

Write an ADR when a decision:

- Constrains how future code is written (e.g. choosing a library, defining a pattern)
- Would be non-obvious to someone joining the project later
- Was debated between multiple viable alternatives
- Changes or supersedes an existing architectural approach

Do **not** write an ADR for routine implementation choices, bug fixes, or decisions that are already captured in code standards.

---

## Step 1 — Determine the ADR Number

List existing ADRs to find the next available number:

```bash
ls harness/knowledge/architecture-decision-records/
```

Use the pattern `NNN-short-description` (no `.md` extension, matching existing ADRs — e.g. `ADR-002-cursor-storage`).

---

## Step 2 — Write the ADR

Create the file at `harness/knowledge/architecture-decision-records/NNN-short-description` using this template:

```markdown
## Context and Problem Statement

<What problem or need prompted this decision? What constraints exist? What goals must
the chosen option satisfy? Keep it to 2-4 paragraphs.>

## Considered Options

### Option #1: <Short name>

<Describe the option.>

**Pros:**

1. <pro>
2. <pro>

**Cons:**

1. <con>
2. <con>

### Option #2: <Short name>

<Same format as above.>

## Decision Outcome

<State the chosen option and the primary reason for the choice. Reference the guiding
principle (e.g. "billing should observe, never gate") if one applies.
Keep it to 1-3 sentences.>
```

---

## Step 3 — Follow the Inverted Pyramid

Structure the ADR so the most important information comes first:

1. **Context and Problem Statement** — the problem being solved and constraints
2. **Considered Options** — alternatives with explicit pros and cons
3. **Decision Outcome** — what was chosen and why

Keep the ADR to roughly one page. If there is extensive supporting material, link to it rather than inlining it.

---

## Step 4 — Handle Superseding

If this ADR supersedes an existing one:

1. Add a note at the top of the new ADR: `Supersedes: ADR-NNN`.
2. Add a note at the top of the old ADR: `Superseded by: ADR-NNN`.
3. Do **not** modify the body of the superseded ADR — it is a historical record.

---

## Step 5 — Update the Index

Add a row to the index table in `harness/knowledge/architecture-decision-records/ABOUT.md`:

```markdown
| [ADR-NNN-short-description](ADR-NNN-short-description) | <one-line summary of the decision> |
```

---

## Step 6 — Review with the Team

Present the ADR to the human for review. ADRs are collaborative — the act of writing surfaces disagreements. Iterate on the content based on feedback before the code that implements it is merged.

---

## Rules

- Keep ADRs short — one page is ideal, two pages maximum.
- Never modify an existing ADR's body. Supersede it with a new ADR instead.
- Write in plain language. Avoid jargon that would not be clear to someone joining the team.
- One decision per ADR is preferred, but tightly coupled decisions can share an ADR.
- ADRs live permanently in the repo — they are historical records, not living documents.
- No `.md` extension on ADR files.
