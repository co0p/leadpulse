---
name: 4dc-refactor
description: "Improve the design of the code that just passed its test, without changing behavior. Tests must stay green throughout. Commits as refactor: <what>. Advances to the next subtask."
---

# Refactor Skill

## One Responsibility

Take the code that just made the current test case green and improve its design — without changing behavior. Tests stay green throughout. Then record the case as complete and either activate the next case in the same subtask or advance.

---

## Foundations

- **Fowler — refactoring.** A behavior-preserving transformation that improves the internal structure of the code. The tests define behavior; if they stay green, behavior is preserved. This is not "cleaning up while adding the feature" — it is a separate discipline, done with a separate hat.
- **Fowler — two hats.** You are either adding behavior (Green) or refactoring. Never both at once. When you refactor, you do not write new tests or change what the code does — you change how it is organized.
- **Beck — the refactor pass.** Once green, ask: can the design be improved? If yes, improve it. If no, skip the commit and advance. The pass is mandatory even when the answer is "nothing to change" — the question must be asked.
- **Poppendieck — eliminate waste.** Refactoring is not gold-plating. It removes duplication, clarifies intent, and keeps the next change cheap. If the change does not serve a future behavior or readability, it is waste.

---

## Expected Input

- `.agent/plan.md` (approved)
- `.agent/implementation.md` with the current `[behavior]` subtask marked `state: green` (test passes, minimal code written, not yet refactored)
- `CONSTITUTION.md` testing strategy

**Narrow context:** load only the files named in the current subtask's `files:` and `references:` fields in `plan.md`. Do not re-scan the codebase — the plan already did that work.

---

## Concrete Output

Updates `.agent/implementation.md` for the current subtask:
- The active test case: `state: complete`
- `evidence:` test output confirming green after refactoring (or a note that no refactoring was needed)
- `refactor:` the commit hash and `refactor: <what changed>` message (omitted if no changes were made)
- If unfinished test cases remain: set `active_test` to the next pending case and leave the subtask `state: in-progress`
- If all test cases are complete: set the behavior subtask `state: complete`

Appends to `.agent/learnings.md` if design decisions or promote candidates emerged.

---

## Scope Boundary

This skill does **one thing**: improve the design of code that already passes its test.

- It does NOT add new behavior.
- It does NOT write new tests.
- It does NOT change what the code does — only how it is organized.
- It does NOT handle `[tidy]` or `[research]` subtasks.

The distinction: `tidy` prepares structure before behavior; `refactor` improves design after behavior. Both are behavior-preserving, but they serve different moments in the cycle.

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
Do NOT change observable behavior. If any test goes red, you changed behavior — revert and try again.
Do NOT add new tests or new production features. This is the refactor hat, not the behavior hat.
Do NOT skip the question "can the design be improved?" — even if the answer is no, the question must be asked and recorded.
Do NOT mark the subtask complete until every planned test case has completed Red → Green → Refactor and all tests are green.
Commit as `refactor: <what changed>` — never `feat:` or `fix:`. If no refactoring was needed, skip the commit and record that decision.
</HARD-GATE>

---

## Process

1. **Read the current subtask** from `.agent/implementation.md` — the `[behavior]` subtask with an active test case in `state: green`.
2. **Run the tests.** Confirm they are green before you start.
3. **Ask the refactor question:** can the design be improved without changing behavior? Look for: duplication, unclear naming, long methods, deep nesting, missing abstraction, poor separation of concerns.
4. **If no improvement is needed:** record `evidence: no refactoring needed — design is sufficient` for the active test case and continue.
5. **If improvement is needed:** make the change. One move at a time — extract method, rename, inline, move. Run the tests after each move. If they go red, revert — you changed behavior.
6. **Run the full test suite** required by the constitution.
7. **Record evidence** in `implementation.md`: set the active test case `state: complete`, with test output confirming green and what was refactored.
8. **Commit** as `refactor: <what changed>`.
9. **Append learnings** if design decisions or promote candidates emerged.
10. **Advance within the subtask:** if a pending test case remains, set it as `active_test` with `state: pending`; `tdd-red` starts its cycle. If all cases are complete, set the behavior subtask `state: complete` and advance to the next subtask.

### When all subtasks are complete

- Run **final verification** — the full test suite or constitution-defined release gate. Confirm all acceptance criteria from `increment.md` are met.
- If `.agent/plan.md` defines acceptance scenarios, run any available scenarios as supplementary evidence. Record each result as `passed`, `failed`, `skipped`, or `not-applicable`; advisory scenarios do not block completion unless an approved project rule makes them required.
- Present the final evidence and remaining risks. Wait for explicit approval.
- On approval, set `implementation.md` top-level `status: complete`. The next skill is `4dc-promote`.

---

## Checklist

### Per subtask
- [ ] Tests green before starting
- [ ] Refactor question asked (duplication, naming, structure, separation)
- [ ] If refactored: tests stayed green after each move
- [ ] If no refactoring needed: decision recorded
- [ ] Active test case marked `state: complete` with evidence
- [ ] Next pending test case activated, or behavior subtask marked complete when all cases are complete
- [ ] Committed as `refactor: <what changed>` (or commit skipped with recorded reason)
- [ ] Learnings appended if decisions or candidates emerged

### When all subtasks complete
- [ ] Final verification run (full suite or release gate)
- [ ] All acceptance criteria from `increment.md` met
- [ ] Optional acceptance scenarios recorded as supplementary evidence, without treating advisory scenarios as blockers
- [ ] User approval received
- [ ] `implementation.md` top-level `status` set to `complete`
- [ ] `learnings.md` has promote candidates listed

---

## Handoff

Updated artifacts: `.agent/implementation.md` (current subtask `state: complete`) + `.agent/learnings.md`
Next skill (if subtasks remain): detected by the orchestrator from the next subtask's type and state
Next skill (if all complete): `4dc-promote` — load `skills/promote/SKILL.md`