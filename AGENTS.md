# go_app_template

`go_app_template` is a Go HTTP service template. It exposes a versioned REST API (`/api/v1`) generated from an OpenAPI 3.0 spec using `oapi-codegen`, manages a `users` domain backed by PostgreSQL, and provides a health check endpoint at `GET /health`.

**Core principle: spec-driven development.** The OpenAPI spec at `api/contracts/v1/api-specs.yaml` is the source of truth for the HTTP API surface. Handler stubs are auto-generated; business logic is implemented separately. Never hand-edit generated files.

## Coding Agent Harness

Always use the provided coding harness in `agent-harness/` to lead your coding activities. This is not optional, the code you produce should always be in perfect alignment with the harness.

### Harness Manifesto

The harness is a compiler pipeline — knowledge and standards are optimisation passes, linters and tests are verification passes, and you are the code generation backend. Read the relevant knowledge or skill file at the point of use rather than front-loading everything. When a skill or knowledge is relevant for your task at hand, make sure you invoke and use it.

### Harness Structure

The harness contains the following top-level areas you can draw from:

```
agent-harness/
├── dev-pipeline/        — Development stages the harness guides you through
├── enforcement/         — Programmatic enforcement mechanisms (golangci-lint, etc.); DO NOT READ THESE — the linter applies these on your behalf
├── plans/               — Execution plans for in-progress work
├── quality/             — Quality grades and tech debt tracking
├── knowledge/           — What the system is, why it's built this way, and how to work in it
│   ├── architecture-decision-records/  — ADRs for key design choices
│   ├── code-standards/                 — Coding conventions and patterns (naming, errors, concurrency, testing, API)
│   ├── domain/                         — Domain vocabulary (users, pagination, API fields)
│   ├── infra/                          — Infrastructure and deployment
│   └── repo-architecture/              — Codebase structure, layer rules, dependency matrix
└── skills/              — How-to guides: read before performing specific tasks
    ├── accessing-systems/              — Connecting to external services (GitHub)
    ├── development/                    — Writing code (TDD, adding features, commits)
    ├── quality/                        — Maintenance tasks (auditing harness, refactoring code, refactoring architecture, agent corrections)
    ├── personas/                       — Expert advisory personas to leverage for task mentorship
    ├── planning/                       — Planning and reviewing work (critique, ADRs)
    └── testing/                        — Running checks, writing tests
```

### How to Use

- `**knowledge/**` — The what, why, and how of this system. Consult proactively whenever a task touches architecture, domain concepts, code standards, or conventions, so your work conforms to established patterns.
- `**skills/**` — Mandatory step-by-step playbooks. You **MUST** read the matching skill **before** performing the task, not after. Skills override your defaults; if your default approach disagrees with the skill, follow the skill.
- `**dev-pipeline/**` — The stages the harness guides work through. Skim at the start of a session to know which stage you are in and what is expected.
- `**quality/**` — Quality grades and known tech debt. Consult only when prioritising improvements or auditing health; not needed for routine feature work.
- `**plans/**` — Scratch space for in-progress plans. Read or update only when working on a plan that lives here; not general reference material.
- `**enforcement/**` — Tool-facing rules applied automatically by linters. DO NOT READ — they run on your behalf.

To list the contents of one or more top-level areas, run the harness lister. Pass the areas you care about; with no arguments it lists all agent-facing areas. `enforcement/` is always excluded.

```bash
agent-harness/enforcement/utils/list-harness.sh knowledge skills   # only these areas
agent-harness/enforcement/utils/list-harness.sh                    # all agent-facing areas
```

### Harness check (do this BEFORE acting)

For every user request, before running any non-readonly tool:

1. Identify what the task touches — an action to perform, an architectural choice, a convention, a workflow stage, a known piece of tech debt, etc.
2. Discover what the harness already says about it. Run `agent-harness/enforcement/utils/list-harness.sh` (optionally narrowed by area, e.g. `skills knowledge`) to see what's available.
3. Read every matching file in full and follow it.

This applies even when you "know how" to do the task. Skip only if, after checking, the harness has nothing relevant. Whatever the harness provides — skills, knowledge, workflows, ADRs, debt entries — takes precedence over (and supplements) your defaults and any user-level skills surfaced by the IDE.
