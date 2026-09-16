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
- **Branch**: copied from `increment.md` — the `increment/<slug>` branch for this cycle
- **Approach**: 2–3 sentences on strategy, including architectural boundary and any performance-sensitive path; no code yet
- **Design** (required when data shapes, call flow, or architecture change; omit only for purely structural tidying):
  - **Data Models**: new or modified types, structs, schemas, or domain objects — named fields with types and invariants. Include before/after when modifying existing shapes.
  - **Call / Data Flow**: the sequence of calls or data transformations this increment introduces, from entry point to persistence or output. Prose or a numbered sequence. Names must match actual symbols in the codebase.
  - **Error / Edge-case Inventory**: every error condition and boundary case the increment must handle, derived by walking the call flow and asking "what can fail here?" For each: the condition, the expected response or recovery, and which behavior subtask covers it. Cases not covered by any subtask are gaps — resolve them before approving the plan.
  - **Observability Intent**: the events, metrics, or log lines this increment should emit, stated as a short list. For each: event name, severity/level, and the data it carries. Omit only when the project's constitution explicitly exempts observability for this type of change.
  - **Architecture Delta**: how the container view in `docs/architecture.md` changes — new containers, removed containers, new communication paths, changed dependency direction, or "no change" if none.
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
- **Planning Decisions**: choices made during planning that a future reader would not infer from the plan alone. For each: the option chosen, the alternatives that were live, and the reason. Captures the "why this approach and not another" before it is lost.

Required `.agent/plan.md` headings:
- `## Goal`
- `## Branch`
- `## Approach`
- `## Design` (omit only when no data shape, call flow, or architecture boundary changes)
- `## Files`
- `## Subtasks`
- `## Context Map`
- `## Acceptance Scenarios` (optional)
- `## Risks`
- `## Planning Decisions`

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
Do NOT omit `## Design` when the increment adds or changes data shapes, call flow, or architecture boundaries. Symbol names in the call flow must match actual codebase symbols, not invented names.
Do NOT write "no change" in Architecture Delta without checking `docs/architecture.md` first.
Do NOT leave error conditions undiscovered — walk every step in the Call/Data Flow and ask "what can fail here?" before the plan is approved. Any uncovered case is a gap and must appear in Risks or be added to a subtask.
Do NOT omit Observability Intent unless the constitution explicitly exempts this type of change.
Do NOT leave `## Planning Decisions` empty — if the approach was obvious with no alternatives considered, state that explicitly rather than omitting the section.
</HARD-GATE>

---

## Process

1. **Read inputs** — `CONSTITUTION.md`, `.agent/increment.md`, and any current architecture, design, ADR, or domain docs touched by the change.
2. **Scan the codebase** — find the files the increment touches. Read them well enough to name specific symbols, line ranges, and the boundaries between modules. This is the one phase that reads broadly; the implement skills will read narrowly.
3. **Identify risks** — surface unknowns, performance-sensitive paths, and cross-cutting concerns before drafting.
4. **Draft `## Design`** — for any increment that changes data shapes, call flow, or architecture: name new or modified types with fields and invariants, trace the call sequence from entry point to output using actual codebase symbols, walk every step in the flow and record every failure condition and edge case with its expected response and covering subtask (gaps go to Risks), list the observability events the change should emit, and state the architecture delta against `docs/architecture.md`. Omit only for pure structural tidying with no behavior change.
5. **Draft the Files section** — list every file the increment will create, modify, touch for reference, or delete. Label each by role. Include `docs/ui.md` when shared UI decisions change.
6. **Draft the Context Map** — link the constitution sections, ADRs, architecture sections, and domain entries that govern this change. The implement skills load these by reference, not by re-reading the whole doc.
7. **Draft ordered subtasks** — each with type, description, files, references, test (for behavior), verification, dependencies, and acceptance criteria coverage.
8. **Optionally define acceptance scenarios** — for a larger feature, describe the user action, precondition, expected outcome, criterion covered, likely test location, and whether the scenario is advisory or a required release gate. Default to advisory.
9. **Record `## Planning Decisions`** — before proposing the plan, capture every non-obvious choice made during planning: approach selection, subtask sequencing decisions, trade-offs accepted. If the approach was unambiguous with no real alternatives, state that.
10. **Conversation: Propose the plan** — present the plan and iterate until the user confirms it covers the increment and the file scope is correct.
11. **On approval** — write `.agent/plan.md`. Then load `skills/implement/SKILL.md` to scaffold `.agent/implementation.md` and populate the todo list.

---

## plan.md Structure

```markdown
# Plan: <increment goal>

## Goal
<one sentence, copied from increment.md>

## Branch
`increment/<slug>` — copied from increment.md

## Approach
<2–3 sentences: strategy, architectural boundary, performance-sensitive path>

## Design

### Data Models

<!-- New or modified types. Include before/after for modifications. -->

**New: `TokenClaims`**
```
TokenClaims {
  subject:   string      // user identifier
  expiresAt: timestamp   // Unix seconds; must be > issuedAt
  issuedAt:  timestamp
  scopes:    string[]    // non-empty
}
```

**Modified: `Session`** (adds `claims` field)
```
Before: Session { id, userId, createdAt }
After:  Session { id, userId, createdAt, claims: TokenClaims }
```

### Call / Data Flow

<!-- Numbered sequence from entry point to output. Symbol names must match the codebase. -->

1. `AuthMiddleware.handle(request)` — extracts bearer token from `Authorization` header
2. `TokenValidator.validate(token) → TokenClaims` — decodes and verifies signature; throws `TokenExpiredError` if past `expiresAt`
3. `SessionStore.load(claims.subject) → Session` — looks up or creates session
4. If `claims.expiresAt - now < 60s`: `TokenRefresher.refresh(session) → TokenClaims` — calls IdP silent-refresh endpoint, updates `session.claims`
5. `request.context.session = session` — attaches session to request context for downstream handlers

### Error / Edge-case Inventory

<!-- Walk the call flow step by step and state every failure condition. Each entry must name the covering subtask or flag as a gap. -->

| Condition | Expected response | Covered by |
|-----------|-------------------|------------|
| Token signature invalid | Reject with 401, do not create session | Subtask 2 |
| Token expired and IdP refresh fails | Reject with 401, log warning with `subject` | Subtask 3 |
| Token near-expiry but IdP is unreachable | Allow request with current session; schedule background retry | Subtask 3 — **gap: retry not yet planned; add to Risks** |
| `expiresAt` before `issuedAt` in claims | Reject with 401; `TokenClaims` invariant violation | Subtask 2 |
| `SessionStore.load` returns null (new user) | Create new session; continue | Subtask 2 |

### Observability Intent

<!-- One line per event. Omit only if the constitution explicitly exempts observability for this change type. -->

| Event | Level | Data |
|-------|-------|------|
| `auth.token.validated` | debug | `subject`, `expiresAt`, latency_ms |
| `auth.token.refresh_attempted` | info | `subject`, `idp`, success: bool |
| `auth.token.refresh_failed` | warn | `subject`, `idp`, error_code |
| `auth.token.invalid` | warn | reason, token_prefix (first 8 chars) |

### Architecture Delta

<!-- Describe container-level changes. State "no change" explicitly if none. -->

No new containers. `TokenRefresher` is a new module inside the existing `auth` container — not a separate deployable. The auth container gains a new outbound call to the IdP refresh endpoint (already present in the container diagram; no diagram update needed). If the IdP endpoint proves unreachable in tests, an ADR is required before promote.

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
- Near-expiry background retry not yet planned (flagged in Error Inventory)

## Planning Decisions

<!-- Choices made during planning that would not be obvious from the plan alone. Captured here before they are lost. -->

| Decision | Chosen | Rejected | Reason |
|----------|--------|----------|--------|
| Where to put refresh logic | New `TokenRefresher` module in auth container | Inside `AuthMiddleware` | Keeps middleware thin; `TokenRefresher` is independently testable |
| Refresh trigger threshold | 60 seconds before expiry | 30s, 120s | Matches IdP's recommended pre-refresh window from ADR-20260105 |
| Error on IdP unreachable | Allow request, warn | Hard reject | Availability over strict freshness; consistent with constitution performance envelope |
```

---

## Checklist

- [ ] `increment.md` acceptance criteria read
- [ ] Relevant source files scanned (not just directory listing)
- [ ] `## Design` written when data shapes, call flow, or architecture boundaries change:
  - [ ] Data Models: new/modified types named with fields and invariants; before/after shown for modifications
  - [ ] Call / Data Flow: numbered sequence from entry point to output; symbol names match codebase
  - [ ] Error / Edge-case Inventory: every failure condition walked from the call flow; each entry names the covering subtask or is flagged as a gap in Risks
  - [ ] Observability Intent: events, levels, and data listed; omitted only with constitutional justification
  - [ ] Architecture Delta: container-level changes stated, or "no change" confirmed against `docs/architecture.md`
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
- [ ] `## Planning Decisions` written: each non-obvious choice records the option chosen, alternatives considered, and reason
- [ ] User approval received
- [ ] `.agent/plan.md` written

---

## Handoff

Terminal artifact: `.agent/plan.md`

Next phase: **implement**

Load `skills/implement/SKILL.md`. The implement skill scaffolds `.agent/implementation.md` from this plan, populates the todo list, and hands off to the first implement skill (`tidy`, `tdd-red`, or `tdd-green`). Only after `implementation.md` exists can the orchestrator route to the appropriate implementation skill.