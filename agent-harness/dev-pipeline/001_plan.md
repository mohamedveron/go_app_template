# 001 — Plan

## What This Is

This is a set of instructions you must adopt and execute, not background context. When you are told to read and execute this file, you become the Planning Agent. Follow the execution steps below in order.

## Purpose

Collaborate with the human in an interactive session to produce a complete, unambiguous coding plan that a downstream Coding Agent can execute without further clarification.

## Context You Have Access To

### Coding Agent Harness and Internal Context

Read the `AGENTS.md` to understand the context and guidance you should draw from proveded by the agent harness.

### External Context

Fetch as needed during the planning session:

- **Task definition** — GitHub issues at [`mohamedveron/go_app_template`](https://github.com/mohamedveron/go_app_template/issues), referenced as `#NNN` in PR bodies and exec-plan keys.
- **Documentation for technology and libraries** (Context7 MCP)

See the skills section on accessing these systems

## Execution Steps

Follow these steps in order. Do not skip ahead.

### Step 1 — Verify You Are in a Read-Only Exploration Mode

Confirm that your current runtime mode supports read-only exploration and does not write files. This skill requires iterative, read-only collaboration with the human. If you are in an edit-capable mode, switch to a read-only mode before proceeding.

### Step 2 — Understand the Task

Ask the human for clear details or specifications on what to build. The task may come from:

- A GitHub issue at [`mohamedveron/go_app_template`](https://github.com/mohamedveron/go_app_template/issues) — referenced as `#NNN`.
- A description in the chat.
- A user story or design doc the human pastes.

If additional context (design files, related PRs, RUNBOOK sections) would help and isn't provided, ask for it.

Do not proceed until you have a clear task definition.

### Step 3 — Study the Task Definition

Read the task definition (and everything it references) carefully. Extract:

- Scope and acceptance criteria.
- Linked PRs, issues, or related work.
- Priority, constraints, and deadlines.
- Whether the change touches the OpenAPI spec — if yes, plan the spec edit before the handler change.

### Step 4 — Deep-Dive into the Codebase

Explore the codebase and documentation thoroughly to understand everything relevant to this task. Use `AGENTS.md` as your entry point to find:

- Architecture docs and ADRs that constrain your choices.
- Familiarise yourself with the domain knowledge and ubiquitous language for the domain(s) we are working in, as defined in the harness' domain knowledge section
- Existing patterns, components, and services you must align with.
- Code standards and conventions that apply.
- etc.

⚠️ **Beware "satisfaction of search."** Do not stop exploring after finding the first relevant file or pattern. Actively look for related code, edge cases, and upstream/downstream dependencies. Ask yourself: *"What else might be affected that I haven't looked at yet?"*

### Step 5 — Collaborative Planning

Engage the human in an interactive planning session. Your goal is to fully align on the design concept before writing the plan. Ask about:

- Ambiguities or gaps in the requirements.
- Trade-offs between approaches (present options with pros/cons).
- Architectural decisions that need to be made.
- Test strategy and how acceptance criteria map to tests.
- Any risks or unknowns.
- etc

Iterate until you and the human are aligned on the approach.

### Step 6 — Draft the Plan

Draft the coding plan. The **plan** must at least specify:

- What files will be created or modified.
- The implementation approach (patterns, libraries, architectural decisions).
- How acceptance criteria map to tests.
- Any risks, trade-offs, or open decisions that were resolved during planning.

### Step 7 — Critique the Plan

Read and execute the critique skill at `[agent-harness/skills/planning/critique-coding-plan.md](../skills/planning/critique-coding-plan.md)`. Apply it against the plan you just drafted. This is a self-review — you are both author and critic.

Present the plan **and** the critique to the human together. Let the human decide which critique points to accept, reject, or modify.

### Step 8 — Finalise the Plan

Incorporate the accepted critique feedback into the plan.

Create a new numbered directory under `agent-harness/plans/` following the pattern `NNN-<issue-or-noticket>-short-desc/` (e.g. `007-issue125-exec-priority/` or `008-NOTICKET-fix-auth-redirect/`). See [`agent-harness/plans/README.md`](../plans/README.md) for the full naming convention, including the race policy if two planners collide on the same `NNN`. Determine the next available number by inspecting the existing directories (including `_archive/`). Write the finalised plan and conversation summaries there:


| File                                                                    | Description                                                                                                                                                                                                                                                      |
| ----------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `agent-harness/plans/NNN-<key>-short-desc/plan.md`                 | The finalised coding plan.                                                                                                                                                                                                                                       |
| `agent-harness/plans/NNN-<key>-short-desc/conversation-summaries.md` | Key decisions, tradeoffs, and rationale from the planning session — including collaborative planning, critique feedback, and human triage decisions. It should give someone who will execute the plan context on what outcomes led to the plan we ended up with. |


Ask the human for **explicit approval** of the finalised plan. Iterate if the human requests further changes.

### Step 9 — Branch, Commit, and Draft PR

Once the human approves the plan:

1. **Create a new branch** following the branch naming conventions in your rules.
2. **Commit** the plan files to the branch with a fitting message
3. **Push** the branch to origin.
4. **Create a draft PR** with:
  - A descriptive title matching the task.
  - The label `plan`.
  - Assigned reviewers as directed by the human.
  - NOTE: this should be a draft PR, not a normal ready-for-review PR

## Done

Your work is complete when the draft PR is open with the approved plan committed.