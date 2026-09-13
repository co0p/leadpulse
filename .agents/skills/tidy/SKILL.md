---
name: 4dc-tidy
description: "Execute one [tidy] subtask: a behavior-preserving structural change. Tests must stay green. Commits as tidy: <what>. Advances to the next subtask."
---

# Tidy Skill

## One Responsibility

Make one structural change — rename, extract, reorganise, inline — that prepares the code for the behavior work that follows. No new observable behavior. Tests stay green.

---

## Foundations

- **Beck — Tidy First.** Make the code easier to change, then change it. Tidying is structural preparation that serves the behavior change that follows. It is not a separate feature; it is the setup.
- **Beck — small steps.** Each tidy subtask is one structural move, separately committed, separately revertible. If a tidy step is too large to describe in one sentence, it should have been split at the plan.
- **Fowler — behavior-preserving.** A tidy commit must leave every test green. If any test changes behavior, the subtask was mislabeled — it belongs in `[behavior]`, not `[tidy]`.
- **Poppendieck — small batches.** One tidy step, one commit. The batch is small so reversion is cheap and review is fast.

---

## Expected Input

- `.agent/plan.md` (approved)
- `.agent/implementation.md` with the current subtask marked `type: tidy` and `state: pending`
- `CONSTITUTION.md` testing strategy

**Narrow context:** load only the files named in the current subtask's `files:` and `references:` fields in `plan.md`. Do not re-scan the codebase — the plan already did that work.

---

## Concrete Output

Updates `.agent/implementation.md` for the current subtask:
- `state: complete`
- `evidence:` test output confirming green (tests unchanged)
- `commit:` the commit hash and `tidy: <what changed>` message

Appends to `.agent/learnings.md` if design or architecture implications emerged.

---

## Scope Boundary

This skill does **one thing**: one structural, behavior-preserving change.

- It does NOT add new behavior.
- It does NOT write failing tests.
- It does NOT refactor for design quality (that is `4dc-refactor`).
- It does NOT touch `[behavior]` or `[research]` subtasks.

The distinction: `tidy` makes the change easier; `refactor` makes the result cleaner. Tidy comes before behavior; refactor comes after.

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
Do NOT change observable behavior. If any test goes red, stop — the subtask was mislabeled.
Do NOT mix tidy work with behavior change in the same commit.
Do NOT skip the test run. Tests must be green before and after.
Do NOT mark the subtask complete without objective evidence (test output showing green).
Commit as `tidy: <what changed>` — never `feat:` or `fix:`.
</HARD-GATE>

---

## Process

1. **Read the current subtask** from `.agent/implementation.md` — the first subtask with `type: tidy` and `state: pending`.
2. **Run the tests.** Confirm they are green before you start. If they are not green, stop — fix the baseline first.
3. **Make the structural change** — rename, extract, reorganise, inline. One move, one purpose: prepare for the behavior change that follows.
4. **Run the tests again.** Confirm they stay green. If any test changed behavior, the change is not tidy — revert and record the mislabel in `learnings.md`.
5. **Record evidence** in `implementation.md`: `state: complete`, test output confirming green.
6. **Commit** as `tidy: <what changed>`.
7. **Append learnings** if design or architecture implications emerged.
8. **Advance** to the next subtask — the orchestrator detects the type and state and loads the right skill.

---

## Checklist

- [ ] Current subtask read (`type: tidy`, `state: pending`)
- [ ] Tests green before starting
- [ ] One structural change made (rename, extract, reorganise, inline)
- [ ] Tests green after — no behavior change
- [ ] `implementation.md` updated: `state: complete`, evidence
- [ ] Committed as `tidy: <what changed>`
- [ ] Learnings appended if implications emerged

---

## Handoff

Updated artifact: `.agent/implementation.md` (current subtask `state: complete`)
Next skill: detected by the orchestrator from the next subtask's type and state.