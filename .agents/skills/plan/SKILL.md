---
name: 4dc-plan
description: "Use after increment.md is approved. Converts increment intent into an ordered, verifiable technical execution plan with file-level detail so implement skills need minimal context loading."
---

# Plan Skill

## One Responsibility

Define HOW to deliver `.agent/increment.md` — an ordered sequence of actionable subtasks with explicit verification points, file-level scope, and context references. The plan is the map; the implement skills follow it without re-discovering the codebase.

---

## Foundations

- **Beck — the planning game.** Tasks are sized for one focused session. The plan is a conversation, not a contract; it is allowed to change when reality arrives, as long as the change is recorded.
- **Poppendieck — pull, don't push.** Subtasks are pulled from acceptance criteria, not pushed from ideas. Small batches flow faster and fail cheaper.
- **Poppendieck — decide as late as possible, but decide.** The plan commits to enough structure that execution is unambiguous, but leaves reversible choices open for the implement skills.
- **Fowler — design for testability.** Every `[behavior]` subtask names its failing test first. The plan surfaces the architectural boundary and any performance-sensitive path so the design can evolve safely.
- **Beck — Tidy First.** When structural preparation makes a behavior change easier, put it in a `[tidy]` subtask that comes first and preserves behavior. The tidy step is not a separate feature; it serves the behavior that follows.

---

## Expected Input

- `CONSTITUTION.md`
- `.agent/increment.md` (must be approved)
- The codebase — read the files that the increment touches, not just the directory listing. The plan must cite specific files, symbols, and line ranges so the implement skills can load narrow context.

---

## Concrete Output

`.agent/plan.md` containing:

- **Goal**: copied from `increment.md` — one sentence
- **Approach**: 2–3 sentences on strategy, including architectural boundary and any performance-sensitive path; no code yet
- **Files**: the complete set of files this increment touches, each labeled by role:
  - `new` — will be created
  - `modify` — existing file that will change
  - `touch` — read for context, unlikely to change (reference only)
  - `delete` — will be removed
- **Subtasks**: ordered list, each with:
  - Type: `[research]`, `[tidy]`, or `[behavior]`
  - Description (what, not how)
  - Files: which files this subtask touches (subset of the Files section)
  - References: specific symbols, line ranges, or sections to look at — `path/to/file.ts:42` or `docs/architecture.md#containers`
  - Tests: for `[behavior]` subtasks, a cohesive list of test cases; each case has an id, test file, and test name (Red handles one active case at a time)
  - Verification step (how to confirm it's done)
  - Dependencies on prior subtasks
  - Acceptance criteria covered
- **Acceptance Scenarios** (optional): user-journey examples that provide feature-level evidence. These are advisory unless the constitution or user explicitly makes them a release gate.
- **Context Map**: pointers to the governing docs for this increment — constitution sections, ADRs, architecture sections, domain entries that the implement skills must respect
- **Risks**: known unknowns that could block execution

Required `.agent/plan.md` headings:
- `## Goal`
- `## Approach`
- `## Files`
- `## Subtasks`
- `## Context Map`
- `## Acceptance Scenarios` (optional)
- `## Risks`

`[research]` is only for blocking unknowns.
`[tidy]` is structural only and must preserve observable behavior.
`[behavior]` changes observable behavior and must be verified by a failing test or failing executable check first.

### Why file-level detail belongs in the plan

The plan is the one phase that reads the whole codebase to find the approach. The implement skills (tdd-red, tdd-green, tidy, refactor) should load only the files named in their subtask — not re-scan the tree. This keeps implement context narrow and fast, and makes the plan's intent verifiable: if a subtask lists no files, it is not actionable.

## Execution Contract

- Use plain, direct language. Keep output scannable.
- Prefer short sentences and bullets.
- State only decisions, actions, blockers, and evidence relevant to this task.
- Do not repeat inputs, instructions, or handover contents.
- Do not add motivational language, generic advice, or decorative explanation.
- Explain choices only when they affect the task, risk, or handoff.
- Never copy internal workflow names, skill names, phase names, orchestrator terms, `.agent/` paths, or `.agents/` paths into permanent product artifacts.
- Before writing a permanent artifact, scan it for internal workflow references and remove them.
- Ask one focused question when blocked.
- End with the next action or handoff.
- Produce only the artifact for this phase. Do not leak work from a later phase into this one.
- Treat tests, architecture notes, ADRs, and user-facing docs as first-class communication artifacts.
- Gather only enough context to identify the governing constraints, the target artifact, and the cheapest validation step. Then act.
- Resolve conflicts in this order: explicit user approval, approved prior-phase artifacts, `CONSTITUTION.md`, this skill.
- Low-risk actions: reads, searches, diffs, and local validation commands.
- Medium-risk actions: local reversible edits to phase artifacts.
- High-risk actions: destructive file operations, external side effects, or skipping a stop gate. Require explicit approval first.
- If a required input is missing or contradictory, ask one focused question or stop and wait for explicit approval. Do not invent missing facts.
- Before finishing, run the phase checklist and confirm every required section is present.

---

<HARD-GATE>
Do NOT start implementation during this phase — no code, no file edits.
Do NOT write `plan.md` until the user explicitly approves the proposed plan.
Do NOT list subtasks without verification steps.
Do NOT list subtasks without a Files field — every subtask must name the files it touches.
Every subtask must map to at least one acceptance criterion from increment.md.
Every `[behavior]` subtask must name its test file and test name.
Every `[behavior]` subtask may contain multiple test cases, but they must describe one cohesive behavior. Track each case separately with `id`, `name`, `file`, and `state`; set exactly one `active_test` at a time.
Acceptance scenarios are separate from focused test cases. Use them for larger or cross-boundary increments, but do not make them blocking by default. A scenario may be `planned`, `available`, `passed`, `skipped`, or `not-applicable`.
Do NOT place a `[behavior]` subtask before the `[tidy]` subtasks it depends on.
File paths must be specific (`src/auth/login.ts`), not globs or directory names alone.
References must point to specific symbols, line ranges, or doc sections — not "see the auth module."
</HARD-GATE>

---

## Process

1. **Read inputs** — `CONSTITUTION.md`, `.agent/increment.md`, and any current architecture, design, ADR, or domain docs touched by the change.
2. **Scan the codebase** — find the files the increment touches. Read them well enough to name specific symbols, line ranges, and the boundaries between modules. This is the one phase that reads broadly; the implement skills will read narrowly.
3. **Identify risks** — surface unknowns, performance-sensitive paths, and cross-cutting concerns before drafting.
4. **Draft the Files section** — list every file the increment will create, modify, touch for reference, or delete. Label each by role. Include `docs/ui.md` when shared UI decisions change.
5. **Draft the Context Map** — link the constitution sections, ADRs, architecture sections, and domain entries that govern this change. The implement skills load these by reference, not by re-reading the whole doc.
6. **Draft ordered subtasks** — each with type, description, files, references, test (for behavior), verification, dependencies, and acceptance criteria coverage.
7. **Optionally define acceptance scenarios** — for a larger feature, describe the user action, precondition, expected outcome, criterion covered, likely test location, and whether the scenario is advisory or a required release gate. Default to advisory.
8. **Conversation: Propose the plan** — present the plan and iterate until the user confirms it covers the increment and the file scope is correct.
9. **On approval** — write `.agent/plan.md`.

---

## plan.md Structure

```markdown
# Plan: <increment goal>

## Goal
<one sentence, copied from increment.md>

## Approach
<2–3 sentences: strategy, architectural boundary, performance-sensitive path>

## Files

| File | Role | Notes |
|------|------|-------|
| `src/auth/login.ts` | modify | add token refresh logic |
| `src/auth/session.ts` | modify | extend session type |
| `src/auth/__tests__/login.test.ts` | new | test file for behavior subtasks |
| `docs/architecture.md` | touch | reference — container view for auth |
| `src/legacy/token.ts` | delete | replaced by new refresh logic |

## Subtasks

### 1. Extract token validation to its own module
type: [tidy]
files:
  - `src/auth/login.ts` (extract from lines 42–68)
  - `src/auth/token.ts` (new — extracted code moves here)
references:
  - `src/auth/login.ts:42` — `validateToken` function
  - `docs/architecture.md#containers` — auth container boundary
verification: `npm test` — all existing tests pass (behavior preserved)
depends on: —
acceptance criteria: — (structural prep)

### 2. Refresh token before expiry
type: [behavior]
files:
  - `src/auth/token.ts` (add `refreshIfNeeded`)
  - `src/auth/__tests__/token.test.ts` (new test)
tests:
  - id: refresh-near-expiry
    file: `src/auth/__tests__/token.test.ts`
    name: `refreshes token when within 60s of expiry`
    state: pending
  - id: preserve-valid-token
    file: `src/auth/__tests__/token.test.ts`
    name: `preserves token when more than 60s from expiry`
    state: pending
active_test: refresh-near-expiry
references:
  - `src/auth/token.ts:12` — current token type
  - `CONSTITUTION.md#performance-envelope` — latency budget for auth
verification: `npm test -- token` — fails first, then passes after implementation
depends on: 1
acceptance criteria: AC-2 (token refreshes before expiry)

### 3. Investigate: does the IdP support silent refresh?
type: [research]
files:
  - `src/auth/idp-client.ts` (read only)
references:
  - `src/auth/idp-client.ts:30` — `refreshToken` method
  - `docs/adr/ADR-20260105-idp-choice.md` — IdP selection rationale
verification: finding recorded in implementation.md; no code change
depends on: —
acceptance criteria: informs AC-2

## Context Map

- `CONSTITUTION.md#testing-strategy` — test depth and naming conventions for auth
- `CONSTITUTION.md#performance-envelope` — latency budget
- `docs/architecture.md#containers` — auth container boundary and dependencies
- `docs/domain.md#session` — domain definition of a session
- `docs/adr/ADR-20260105-idp-choice.md` — why this IdP was chosen

## Acceptance Scenarios

These scenarios are optional feature-level evidence. They do not block implementation or promotion unless the `gate` field explicitly says `required` under an approved project rule.

### AT-1: Successful account lookup

- criterion: AC-1
- user action: request account information
- precondition: provider returns a valid account
- expected outcome: account information is shown
- evidence: `tests/acceptance/account-lookup.test.ts`
- gate: advisory
- state: planned

### AT-2: Provider unavailable

- criterion: AC-3
- user action: request account information
- precondition: provider times out
- expected outcome: a temporary service error is shown
- evidence: `tests/acceptance/account-lookup.test.ts`
- gate: advisory
- state: planned

## Risks
- Token refresh may race with in-flight requests (concurrency)
- IdP rate limit on refresh calls (unknown — see subtask 3)
```

---

## Checklist

- [ ] `increment.md` acceptance criteria read
- [ ] Relevant source files scanned (not just directory listing)
- [ ] Files section lists every file the increment touches, labeled by role
- [ ] Every subtask has a Files field naming the files it touches
- [ ] Every subtask has a References field with specific symbols, line ranges, or doc sections
- [ ] Every `[behavior]` subtask has a cohesive `tests` list with `id`, `file`, `name`, and `state` for each case
- [ ] Every `[behavior]` subtask has exactly one `active_test`
- [ ] Related test cases are grouped together; unrelated behavior is split into another subtask
- [ ] Every subtask has a verification step
- [ ] Every acceptance criterion has a covering subtask
- [ ] Optional acceptance scenarios identify user action, expected outcome, criterion, and evidence location
- [ ] Acceptance scenarios are explicitly marked advisory unless made blocking by an approved constitution or explicit user request
- [ ] Context Map links the governing constitution sections, ADRs, architecture, and domain docs
- [ ] Any architectural or performance-sensitive change is reflected in the approach or risks
- [ ] Risks documented
- [ ] User approval received
- [ ] `.agent/plan.md` written

---

## Handoff

Terminal artifact: `.agent/plan.md`
The implement skills (`tidy`, `tdd-red`, `tdd-green`, `refactor`) read the plan and load only the files named in their subtask's Files and References fields — no broad codebase scanning.
Next skill: detected by the orchestrator from the first subtask's type and state.