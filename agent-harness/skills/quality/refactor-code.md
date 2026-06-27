# Refactoring

Systematically improve the internal structure of existing code without changing its external behaviour. This skill translates the principles from Martin Fowler's *Refactoring: Improving the Design of Existing Code* into a two-phase process: first audit and plan, then implement with approval.

**Cardinal rule:** Refactoring and feature work are two separate activities. Wear one hat at a time. If you discover missing behaviour during a refactoring pass, stop, switch hats, add a failing test and implement the behaviour, then switch back.

---

## Overview

| Phase | Goal | Output | Requires human? |
|-------|------|--------|-----------------|
| **Phase 1 — Audit & Plan** | Scan the codebase for code smells, prioritise findings, and propose concrete refactorings | An exec-plan document in `agent-harness/plans/` | Yes — human reviews and approves before Phase 2 |
| **Phase 2 — Implement** | Execute the approved refactorings in small, tested, committed steps | Clean code + commits | No — but human can review commits |

---

# Phase 1 — Audit & Plan

The goal of Phase 1 is *analysis only*. Do not change any source code. The output is a written plan that a human can review, adjust, and approve.

## Step 1.1 — Scope

Determine what code to audit. This can come from:

- A human-specified scope (e.g. "audit the exec manager", "look at `internal/transport/http/v1/`")
- The code you are about to change for a feature or fix (campsite rule — audit what is in the way)
- A full-codebase sweep (only when explicitly requested)

If no scope is given, ask the human before proceeding.

## Step 1.2 — Detect Smells

Read the code within the agreed scope. Check against the smell catalog below. Record every smell — do not fix anything.

### Function / Method-Level Smells

| Smell | How to spot it |
|-------|---------------|
| **Long Function** | Body exceeds ~15–20 lines, or multiple levels of abstraction in one body. |
| **Duplicated Code** | Same or near-identical logic in two or more places. |
| **Too Many Parameters** | Function takes 4+ parameters. Often a sign that a concept is missing. |
| **Flag Arguments** | A boolean parameter that makes the function do two different things. |
| **Dead Code** | Functions, branches, or variables never reached. Exports with zero consumers. |

### Data / State Smells

| Smell | How to spot it |
|-------|---------------|
| **Primitive Obsession** | Raw `string`/`int`/`bool` representing a domain concept (e.g. a `string` for a user spec). |
| **Data Clumps** | The same group of fields/parameters travels together across multiple functions or structs. |
| **Mutable Data** | Fields mutated in multiple places without clear ownership. |
| **Temporary Field** | A struct field only meaningful in some states, zero/nil otherwise. |

### Module / Package-Level Smells

| Smell | How to spot it |
|-------|---------------|
| **Large File** | A single file does too many things — more than one reason to change. |
| **Feature Envy** | A function uses more data/methods from another package than its own. |
| **Shotgun Surgery** | A single conceptual change requires edits across many unrelated files. |
| **Divergent Change** | One package is changed for multiple unrelated reasons. |
| **Middle Man** | A struct that delegates almost everything to another, adding no value. |
| **Speculative Generality** | Interfaces or extension points that exist "in case we need them" with no current consumer. |

### Project Code-Standards Violations

1. Open `agent-harness/knowledge/code-standards/_index.md` — it maps task types to the relevant standards files.
2. Read the files relevant to what the scoped code touches (handlers, managers, concurrency, naming, etc.).
3. Record violations as findings, noting which standard was violated and where.

Do not front-load every standards file. Only read those relevant to the audit scope.

## Step 1.3 — Prioritise

Score each finding:

1. **Is it in the way?** Does it make current or upcoming work harder?
2. **Is it frequently changing?** Code that changes often benefits most from clarity.
3. **Is it causing bugs?** Mutable state, unclear conditionals, and duplicated logic are common bug factories.
4. **Is it blocking testability?**

Tier:
- **High** — Actively causing pain or blocking work. Fix in this pass.
- **Medium** — Would improve the codebase but not urgent. Include as optional.
- **Low** — Minor style issue. Note but do not propose a refactoring.

## Step 1.4 — Propose Refactorings

For every High and Medium finding, propose a concrete fix using the [Refactoring Catalog](#refactoring-catalog) below. For code-standards violations, reference the specific standard file and rule.

## Step 1.5 — Write the Exec-Plan

Create `agent-harness/plans/NNN-<key>-refactor-<scope-slug>/plan.md` following the naming convention in `agent-harness/plans/README.md` (e.g. `007-issue125-refactor-exec-manager/`). Use this template:

```markdown
# Refactoring Plan: <Scope Description>

Date: <YYYY-MM-DD>
Status: PENDING APPROVAL

## Scope
<What was audited and why.>

## Summary
<2–3 sentence overview of the health of the audited code and the thrust of proposed changes.>

## Test Coverage
<Current state. Note gaps that need characterisation tests before refactoring begins.>

## Proposed Refactorings

### R1 — <Short title>
- **Priority:** High | Medium
- **Smell / Violation:** <Smell name, or "Standards violation: <file>" for code-standards findings>
- **Location:** <File(s) and function(s)>
- **Description:** <What is wrong and why it matters — 2–3 sentences.>
- **Proposed refactoring:** <Catalog name or standards-alignment description>
- **Plan:** <Concrete steps — what will be extracted, moved, renamed.>
- **Risk:** <What could go wrong.>

### R2 — <Short title>
...repeat...

## Items Noted but Not Proposed
<Low-priority smells, one line each.>

## Execution Order
<Numbered list of R-items. Earlier items must not depend on later ones.>
```

## Step 1.6 — Present for Approval

Present the exec-plan and ask:

1. Any refactorings to **remove**?
2. Any to **add**?
3. Any **priority changes**?
4. Is the **execution order** correct?
5. Any areas needing **characterisation tests** first?

Do not proceed to Phase 2 until the human explicitly approves. Set `Status: APPROVED`.

---

# Phase 2 — Implement

Execute the approved exec-plan in execution order. Each R-item is a self-contained unit with its own commit.

## Prerequisites

1. **Confirm tests pass.** Do not start if tests are red.
2. **Commit (or stash) current state.** You need a clean rollback point.
3. **Write characterisation tests** for gaps flagged in the plan. Commit them before starting.

## Step 2.1 — Implement Each Refactoring

### The Rhythm

```
1. Read the R-item
2. Select the refactoring from the catalog (or the referenced code standard)
3. If the R-item touches a code-standards area, read that file now
4. Apply the smallest mechanical step
5. Run tests → green? Continue. Red? Revert the last step.
6. Repeat 4–5 until complete
7. Check for dead code left behind (unused imports, orphaned types)
8. Run the linter — fix anything introduced
9. Commit with type "refactor"
10. Mark the R-item DONE in the exec-plan
11. Move to the next R-item
```

**Never take a step so large that you cannot revert it instantly.**

Commit format: `refactor: <smell addressed> via <refactoring applied>`

## Step 2.2 — Final Verification

1. Run the full test suite — everything must pass.
2. Run the linter — no new warnings.
3. Review naming — do names still make sense with the cleaner structure?
4. Set exec-plan `Status: COMPLETED`.

---

# Refactoring Catalog

## Extracting and Inlining

### Extract Function
**When:** A fragment can be named — loop body, one branch of a conditional, or a block that does one thing inside a longer function.
**Mechanics:** Create a function named after *what* it does. Pass locals as parameters. Replace the fragment with a call. Test.

### Inline Function
**When:** A function body is as clear as its name, or it is a trivial wrapper adding indirection without value.
**Mechanics:** Replace every call site with the body. Remove the function. Test.

### Extract Variable
**When:** An expression is complex and its intent is unclear.
**Mechanics:** Declare a named variable. Assign the expression. Replace the original with the variable. Test.

### Inline Variable
**When:** A variable name says no more than the expression itself, or it gets in the way of further refactoring.

---

## Encapsulation

### Encapsulate Variable
**When:** A package-level variable is accessed directly from multiple sites.
**Mechanics:** Create a getter (and setter if needed). Replace direct accesses. Unexport the variable. Test.

### Replace Primitive with Struct
**When:** A primitive carries domain meaning — validation, formatting, or invariants belong with the value (e.g. a raw `string` used as a user spec).
**Mechanics:** Define a struct. Add a constructor that validates. Add methods for domain operations. Replace raw uses. Test.

### Introduce Parameter Object
**When:** Several parameters always travel together across multiple call sites.
**Mechanics:** Define a struct. Replace the parameter list with the struct at all call sites. Test.

### Extract Type
**When:** A struct does two distinct things, or a set of fields always appears together across multiple structs.
**Mechanics:** Define the new struct/interface. Move relevant fields and methods. Embed or reference it from the original. Update call sites. Test.

---

## Moving Features

### Move Function
**When:** A function references more elements from another package than its own (feature envy).
**Mechanics:** Move the function to the target package. Update imports and call sites. Test.

### Replace Loop with Helper Function
**When:** A loop body processes a collection through multiple stages and the intent is obscured. Go has no pipeline operators — extract a named helper.
**Mechanics:** Extract the loop into a function named after what it produces. Replace the inline loop with a call. Test.

### Split Loop
**When:** A single loop does two unrelated things. Splitting makes each loop easier to extract and reason about.

---

## Simplifying Conditionals

### Decompose Conditional
**When:** A complex `if`/`else if`/`else` has non-trivial tests or bodies that obscure intent.
**Mechanics:** Extract the condition into a named function. Extract each branch body into its own function. Test.

### Replace Nested Conditionals with Guard Clauses
**When:** A function has deeply nested `if`/`else` with a clear normal path. Guard clauses (early returns) are idiomatic Go — use them.
**Mechanics:** For each "exceptional" path, invert the condition and return early. Remove the `else`. Repeat until the happy path is at the bottom with no nesting. Test.

### Replace Conditional with Interface
**When:** A `switch` over a type discriminator appears in multiple places with distinct behaviour per case. Use only when polymorphism is genuine — a plain `switch` is often right for a single site.
**Mechanics:** Define a small interface with the varying method. Create one implementation per case. Replace the switch with a map lookup and interface call. Test.

### Introduce Special Case
**When:** Multiple sites check for `nil` or a sentinel and apply the same default logic.
**Mechanics:** Identify the repeated check and its default. Create a helper or zero-value that encapsulates the behaviour. Replace repeated checks. Test.

---

## Simplifying API Boundaries

### Separate Query from Modifier
**When:** A function both returns a value and produces a side effect.
**Mechanics:** Create a pure query that returns the value. Change the original to return only an error (or nothing). Update call sites. Test.

### Remove Flag Argument
**When:** A boolean parameter makes a function behave in two fundamentally different ways.
**Mechanics:** Create two functions, one per behaviour. Replace call sites. Remove the original. Test.

### Preserve Whole Object
**When:** Several fields are extracted from a struct and passed individually to a function. Pass the struct itself to reduce parameter count and stay resilient to new fields.

---

## Dealing with Delegation

### Replace Delegation with Embedding
**When:** A struct forwards almost every method call to an inner struct with no added logic. Go embedding eliminates the boilerplate.
**Mechanics:** Replace the named field with an embedded type. Remove forwarding methods. Update call sites that used the named field. Test.

### Collapse Middle Man
**When:** A struct exists only to delegate to another, adding no value.
**Mechanics:** Update callers to use the real struct directly. Remove the middle-man struct. Test.

---

# Decision Heuristics

| Situation | Guidance |
|-----------|----------|
| **Refactor now or later?** | If you are about to modify the code, refactor first. If you are not about to change it, leave it. |
| **When to stop?** | When the code clearly expresses its intent and the change you came to make is easy. Aim for "good enough to work with." |
| **Too risky?** | No tests → write tests first. Hot path → smaller steps, more commits. Unsure about target design → try it in a branch. |
| **Performance concern?** | Almost never measurable. Profile first, optimise second. Never sacrifice clarity for speculative performance. |
| **Too many smells?** | Start with what you are about to change. Campsite rule: leave it better than you found it. |
| **Found missing behaviour?** | Stop refactoring. Write a failing test, implement the behaviour, get to green. Then resume. Two hats, never both at once. |

---

# Anti-Patterns

| Anti-pattern | Why it hurts | What to do instead |
|---|---|---|
| Big-bang refactoring | High risk of bugs; hard to revert | Small tested steps, commit after each |
| Refactoring without tests | No safety net | Write characterisation tests first |
| Mixing refactoring with feature work | Can't tell which caused a failure | One hat at a time; commit refactoring before feature |
| Refactoring for its own sake | Wastes time on code nobody touches | Only refactor code in the way of current work |
| Speculative design | Abstractions for future needs that never arrive | Refactor toward the current requirement; Rule of Three before abstracting |
| Renaming without reason | Churn in version control | Rename only when the current name actively misleads |
| Skipping Phase 1 | Wrong work, misaligned priorities | Always audit and get approval before writing code |
