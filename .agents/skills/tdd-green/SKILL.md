---
name: 4dc-tdd-green
description: "Make the current failing test pass with minimal code. No refactoring. Handles [research] subtasks too (which skip Red and Refactor). Sets state: green and hands off to refactor."
---

# TDD Green Skill

## One Responsibility

Write the minimal production code that makes the current failing `active_test` pass. Nothing more. A behavior subtask may contain several cohesive test cases; complete one case, then hand off to `4dc-refactor` before activating the next case.

- For a `[behavior]` subtask with an active test in `state: red`: write minimal code to pass that test, set the test to `state: green`, and hand off to `4dc-refactor`.
- For a `[research]` subtask in `state: pending`: resolve the blocking unknown, set `state: complete`, advance.

---

## Foundations

- **Beck — Green.** Write the minimal code that makes the test pass. Not the code you would eventually want — the code that passes the test right now. YAGNI.
- **Beck — two hats.** Adding behavior and refactoring are separate disciplines. This skill wears the behavior hat only. The refactor hat is `4dc-refactor`.
- **Poppendieck — eliminate waste.** No speculative generality, no "while I'm here" features, no extra methods. The test defines the spec; write only what satisfies it.
- **Fowler — make it work, then make it right.** First green (make it work), then refactor (make it right). This skill is the "make it work" step.

---

## Expected Input

- `.agent/plan.md` (approved)
- `.agent/implementation.md` with one subtask in `state: red` (behavior) or `state: pending` (research)
- `CONSTITUTION.md` testing strategy

**Narrow context:** load only the files named in the current subtask's `files:` and `references:` fields in `plan.md`. Do not re-scan the codebase — the plan already did that work.

---

## Concrete Output

Updates `.agent/implementation.md` for the current subtask:
- For `[behavior]`: active test `state: green`, `evidence:` test output showing pass, `commit:` hash with `feat:` or `fix:` prefix
- For `[research]`: `state: complete`, `evidence:` finding, `commit:` hash with `research:` prefix (or skipped)

Appends to `.agent/learnings.md` if decisions, deviations, or surprises emerged.

---

## Scope Boundary

This skill does **one thing**: make the test pass with minimal code.

- It does NOT refactor (that is `4dc-refactor`).
- It does NOT write failing tests (that is `4dc-tdd-red`).
- It does NOT handle `[tidy]` subtasks (that is `4dc-tidy`).
- It does NOT write permanent docs (that is `4dc-promote`).
- It does NOT run final verification (that is `4dc-refactor`, when all subtasks are complete).

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
Do NOT write more code than needed to pass the test. No speculative generality, no "while I'm here" features.
Do NOT refactor in this skill. Improving the design is a separate hat — load `4dc-refactor` next.
Do NOT mark a `[behavior]` subtask `state: complete` — set the active test to `state: green` and hand off to refactor. The subtask becomes complete only after every planned test case has completed its own Red → Green → Refactor loop.
Do NOT mix behavior change and refactoring in the same commit.
Commit behavior work as `feat: <what changed>` or `fix: <what changed>`. Never `refactor:` or `tidy:`.
</HARD-GATE>

---

## Process

### For a `[behavior]` subtask in `state: red`

1. **Read the failing `active_test`** from the test file recorded in `implementation.md`.
2. **Write minimal production code** to make the test pass. Add only what the test requires — no extra methods, no speculative abstraction, no "I'll need this later."
3. **Run the test.** Confirm it passes.
4. **Run the narrowest relevant tests**, then the broader suite required by the constitution.
5. **Record evidence** in `implementation.md`: set the active test to `state: green`, with test output showing pass. Leave the other test cases unchanged.
6. **Commit** as `feat: <what changed>` or `fix: <what changed>`.
7. **Append learnings** — decisions, deviations, surprises, promote candidates.
8. **STOP.** The next skill is `4dc-refactor`, which completes this active test case before another case is activated.

### For a `[research]` subtask in `state: pending`

1. **Resolve only the blocking unknown.** Do not build production features.
2. **Record the finding** in `implementation.md`: `state: complete`, evidence of what was learned.
3. **Commit** investigation notes or scratch code as `research: <what was investigated>` (or skip commit if no code changed).
4. **Append learnings** with the finding and its impact on the plan.
5. **Advance** to the next subtask — the orchestrator detects the type and state and loads the right skill.

---

## implementation.md Structure

```markdown
# Implementation: <increment goal>

status: in-progress  <!-- or: complete | blocked -->
started: <ISO date>

## Baseline
Tests before: X passing, Y failing

## Subtasks

### 1. <subtask name>
type: behavior
state: in-progress
tests:
  - id: provider-success
    name: <test name>
    file: <test file>
    state: complete
    evidence: `npm test` — 12 passing, 0 failing
    commit: <hash> — feat: <what changed>
    refactor: <hash> — refactor: <what improved>
  - id: provider-not-found
    name: <test name>
    file: <test file>
    state: red
    evidence: fails — <assertion reason>
active_test: provider-not-found

### 2. <subtask name>
type: tidy
state: complete
evidence: `npm test` — 12 passing, 0 failing (tests unchanged)
commit: <hash> — tidy: <what changed>

### 3. <subtask name>
type: behavior
state: in-progress
tests:
  - id: <case-id>
    name: <test name>
    file: <test file>
    state: green
    evidence: `npm test` — 12 passing, 0 failing
    commit: <hash> — feat: <what changed>
active_test: <case-id>

### 4. <subtask name>
type: behavior
state: in-progress
tests:
  - id: <case-id>
    name: <test name>
    file: <test file>
    state: pending
active_test: <case-id>
```

---

## learnings.md Structure

```markdown
# Learnings: <increment goal>

## Decisions
- <decision>: <rationale>

## Deviations
- Subtask N: <what changed and why>

## Surprises
- <unexpected finding>

## Promote Candidates
- <ADR, architecture update, domain-language update, test pattern, or performance contract worth keeping>
```

---

## Checklist

### Per `[behavior]` subtask
- [ ] Minimal code written to pass the test (no speculative generality)
- [ ] Narrowest relevant tests run, then broader suite
- [ ] Active test case marked `state: green` with evidence
- [ ] Committed as `feat:` or `fix:`
- [ ] Learnings appended (decisions, deviations, surprises, promote candidates)
- [ ] After refactor, activate the next pending test case or mark the behavior subtask complete when all cases are complete

### Per `[research]` subtask
- [ ] Blocking unknown resolved
- [ ] Finding recorded with evidence
- [ ] Subtask marked `state: complete`
- [ ] Committed as `research:` (or commit skipped)
- [ ] Learnings appended

---

## Handoff

Updated artifacts: `.agent/implementation.md` (current `[behavior]` subtask `state: green`, or `[research]` subtask `state: complete`) + `.agent/learnings.md`
Next skill (after `[behavior]` green): `4dc-refactor` — load `skills/refactor/SKILL.md`
Next skill (after `[research]` complete): detected by the orchestrator from the next subtask's type and state