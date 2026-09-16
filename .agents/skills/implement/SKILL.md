---
name: 4dc-implement
description: "Run once after plan.md is approved. Scaffolds implementation.md from the plan, populates the internal todo list with every subtask, and hands off to the first implement skill (tidy, tdd-red, or tdd-green)."
---

# Implement Skill

## One Responsibility

Bootstrap the implementation loop: read the approved plan, create `.agent/implementation.md` with every subtask in its initial state, populate the **internal todo list** so progress is visible throughout the loop, then hand off to the correct first skill. This skill runs exactly once per cycle — it does not loop.

---

## Foundations

- **Beck — make the plan visible.** The implementation loop works from a single tracking artifact. Creating it explicitly, from the approved plan, prevents drift between what was planned and what is being tracked.
- **Poppendieck — small batches, pull.** Each subtask becomes one todo item. Work is pulled one item at a time; nothing is started until the previous item is complete.
- **Fowler — observable progress.** The todo list is the live view of the implementation loop for the user. Update it at every transition — not in batches at the end.

---

## Expected Input

- `.agent/plan.md` (must be approved)
- `CONSTITUTION.md` (for testing strategy reference)

---

## Concrete Output

### `.agent/implementation.md`

Scaffolded from `plan.md`. Every subtask from `## Subtasks` becomes an entry in the initial state:

```markdown
# Implementation: <goal from plan.md>

status: in-progress
branch: <branch from plan.md>
started: <ISO date>

## Baseline
Tests before: establish by running the full suite before any changes

## Subtasks

### 1. <subtask name from plan.md>
type: <tidy | behavior | research>
state: pending
tests:                         # behavior subtasks only
  - id: <id>
    name: <name>
    file: <file>
    state: pending
active_test: <first test id>   # behavior subtasks only

### 2. <subtask name>
type: <type>
state: pending
...
```

Rules:
- Copy subtask names, types, and test lists verbatim from `plan.md`. Do not paraphrase.
- All subtasks start `state: pending`.
- For `[behavior]` subtasks, copy the full `tests` list and set `active_test` to the first test id.
- Do not add fields not present in the template above — the implement skills own those fields.

### Internal todo list

Populate the todo list immediately after writing `implementation.md`. One item per subtask, in plan order. Use the subtask name and type as the label:

```
[tidy]     1. <subtask name>          pending
[behavior] 2. <subtask name>          pending
[behavior] 3. <subtask name>          pending
[research] 4. <subtask name>          pending
```

**Todo list discipline — applies for the entire implementation loop, not just this skill:**

- Mark a todo `in_progress` the moment the skill for that subtask starts. Only one item is `in_progress` at a time.
- Mark a todo `completed` only after `implementation.md` records `state: complete` for that subtask and the commit hash is present.
- Never mark a subtask complete based on intent — only on evidence in `implementation.md`.
- If a subtask is split or re-ordered during the loop, update the todo list to match before continuing.
- The todo list is the user's real-time view of cycle progress. Keep it accurate.

---

## Scope Boundary

This skill does **one thing**: initialise the tracking artifacts and hand off.

- It does NOT write production code.
- It does NOT write tests.
- It does NOT make any structural change to the codebase.
- It does NOT run the full test suite (that baseline is established as the first act of the first implement skill).

---

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
Do NOT create implementation.md until plan.md is confirmed approved.
Do NOT paraphrase subtask names or test ids — copy them verbatim from plan.md.
Do NOT mark any todo in_progress until the corresponding skill has actually started work.
Do NOT skip populating the todo list — it is a required output of this skill, not optional.
</HARD-GATE>

---

## Process

1. **Read `.agent/plan.md`** — extract goal, branch, and the full ordered subtask list including test case ids for behavior subtasks.
2. **Scaffold `.agent/implementation.md`** — create the file using the template above. Subtask names and test lists copied verbatim.
3. **Populate the todo list** — one item per subtask, all `pending`, in plan order, labeled with type and name.
4. **Hand off** — the orchestrator detects the first subtask's type and state and loads the correct skill:
   - First subtask `[tidy]` → load `skills/tidy/SKILL.md`
   - First subtask `[behavior]` → load `skills/tdd-red/SKILL.md`
   - First subtask `[research]` → load `skills/tdd-green/SKILL.md`

---

## Checklist

- [ ] `.agent/plan.md` read and confirmed approved
- [ ] Goal and branch copied from `plan.md` into `implementation.md`
- [ ] Every subtask from `plan.md` present in `implementation.md` with correct type and `state: pending`
- [ ] Every `[behavior]` subtask has its `tests` list and `active_test` set
- [ ] Todo list populated with one item per subtask, all `pending`
- [ ] No code written, no tests written, no structural changes made

---

## Handoff

Terminal artifact: `.agent/implementation.md` (scaffolded, all subtasks `state: pending`)
Todo list: populated, all items `pending`
Next skill: detected by the orchestrator from the first subtask's type and state.