---
name: 4dc-tdd-red
description: "Write exactly one failing test for the current [behavior] subtask. Confirm it fails for the right reason. Do not write production code. Do not refactor."
---

# TDD Red Skill

## One Responsibility

Write one failing test for the current `active_test` in the `[behavior]` subtask and confirm it fails for the right reason. Then stop. A subtask may contain several cohesive test cases; this skill advances one case at a time.

---

## Foundations

- **Beck — Red.** The test is a specification written as code. It must fail, and it must fail for the right reason — not a syntax error, not a missing import, but the absence of the behavior.
- **Poppendieck — eliminate waste.** The test is the spec. There is no separate spec document to maintain.
- **Fowler — test-first makes design testable.** Writing the test first forces the public interface to exist before the implementation, which keeps the design loosely coupled.

---

## Expected Input

- `.agent/plan.md` (approved)
- `.agent/implementation.md` with the current subtask marked `state: pending` and `type: behavior`
- `CONSTITUTION.md` testing strategy (test depth, naming, isolation conventions)

**Narrow context:** load only the files named in the current subtask's `files:` and `references:` fields in `plan.md`. Do not re-scan the codebase — the plan already did that work.

---

## Concrete Output

Updates `.agent/implementation.md` for the current subtask:
- The active test case: `state: red`
- `test:` the failing test name and location
- `evidence:` the test runner output showing the failure and the assertion reason
- The subtask remains in progress until every planned test case is complete

Appends to `.agent/learnings.md` only if a decision or surprise emerged while writing the test (e.g. the subtask needs splitting, or the interface is unclear).

---

## Scope Boundary

This skill does **one thing**: write the failing test.

- It does NOT write production code.
- It does NOT refactor.
- It does NOT make the test pass.
- It does NOT touch other subtasks.
- It does NOT handle `[tidy]` or `[research]` subtasks — those skip Red and go straight to `4dc-tdd-green`.

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
Do NOT write production code in this skill — not even a stub that makes the test pass.
Do NOT skip the failure-confirmation step. The test must run and fail for the right reason.
Do NOT write tests for more than one active test case per invocation.
Do NOT proceed to `4dc-tdd-green` if the test fails for the wrong reason (syntax error, missing import, wrong assertion). Fix the test first.
</HARD-GATE>

---

## Process

1. **Read the current subtask** from `.agent/implementation.md` — the first `[behavior]` subtask with an unfinished `tests` list. Select its `active_test`, or the first test case with `state: pending`.
2. **Read the testing strategy** in `CONSTITUTION.md` and `docs/testing.md` — choose the cheapest test depth that gives sufficient confidence at this boundary.
3. **Write one test case** for `active_test` that specifies one example of the behavior. Name it in domain language so the test reads as a specification. Do not implement or activate the other cases yet.
4. **Run the test.** Confirm it fails. Read the failure message — it must fail because the behavior does not exist, not because of a setup or import error.
5. **Record evidence** in `.agent/implementation.md`:
   - Set the active test case `state: red` (the subtask itself remains in progress)
   - Record `test:` the test name and file
   - Record `evidence:` the failure output (test name + assertion reason)
6. **STOP.** The next skill is `4dc-tdd-green`, which handles this active test case.

If the test reveals the subtask is too large or the interface is unclear, record the finding in `learnings.md` and ask one focused question before proceeding.

---

## Checklist

- [ ] Current `[behavior]` subtask and `active_test` read from `implementation.md`
- [ ] Testing strategy consulted (depth, naming, isolation)
- [ ] One `active_test` written that specifies one example of the behavior
- [ ] Test run and confirmed failing for the right reason
- [ ] Active test case updated: `state: red`, test name, evidence
- [ ] `active_test` remains the only active test case
- [ ] No production code written
- [ ] Any interface surprise recorded in `learnings.md`

---

## Handoff

Updated artifact: `.agent/implementation.md` (current subtask `state: red`)
Next skill: `4dc-tdd-green` — load `skills/tdd-green/SKILL.md`