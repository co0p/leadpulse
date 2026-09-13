---
name: 4dc-adr
description: "On-demand. Write one Architecture Decision Record when a structural, hard-to-reverse, or non-obvious choice emerges. Single decision, single file, with rationale and consequences."
---

# ADR Skill

## One Responsibility

Capture one architectural decision in a single ADR file with context, alternatives, rationale, and consequences. Nothing else.

---

## Foundations

- **Fowler — evolutionary architecture.** The system evolves; irreversible decisions need a decision log so future change is informed, not blind.
- **Poppendieck — decide as late as possible, but decide.** Record the decision at the moment commitment becomes necessary, with the options that were live at that moment.
- **Beck — make irreversible decisions visible.** A decision worth recording is one a newcomer would not infer from the code.

---

## Expected Input

- `CONSTITUTION.md` (architectural boundaries the decision must respect)
- The decision itself — what was chosen, what was rejected, why now
- `docs/adr/` (existing ADRs, to link related or superseded decisions)

---

## Concrete Output

`docs/adr/ADR-YYYYMMDD-<slug>.md` containing:
- **Decision**: one sentence stating what was decided
- **Status**: Accepted | Superseded | Deferred
- **Context**: why this decision matters now; what constraint or problem forces it
- **Alternatives**: the options that were live, each with its trade-off
- **Rationale**: why the chosen option wins over the others, grounded in this project's context
- **Consequences**: what gets better and what gets harder (both sides)
- **Related**: links to related or superseded ADRs

Required headings (see the ADR template in `templates/adr.md` for the full structure):
- `# ADR-YYYYMMDD — [Decision Title]`
- `**Decision:**`
- `**Status:**`
- `**Context**`
- `**Alternatives**`
- `**Rationale**`
- `**Consequences**`
- `**Related**`

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
Do NOT write an ADR for implementation details (variable names, argument ordering, small refactors).
Do NOT write an ADR for a decision already captured in `CONSTITUTION.md` or an existing ADR.
Do NOT write the ADR until the user confirms the decision and its rationale.
One ADR per decision — if two decisions are entangled, write two ADRs and cross-link them.
</HARD-GATE>

---

## When to Write an ADR

Write one when the choice is:
- **Structural** — affects multiple parts of the system (new layer, new service, new dependency direction)
- **Hard to reverse** — schema design, language choice, external dependency, data migration
- **Non-obvious** — a newcomer would not infer it from reading the code
- **Trade-off-laden** — performance vs. simplicity, flexibility vs. cost

Do **not** write one for:
- Bug fixes or maintenance patches
- Implementation details that live in code
- Decisions fully captured by `CONSTITUTION.md` guardrails

---

## Process

1. **Name the decision** — one sentence: what was decided. If you cannot state it in one sentence, the decision is not yet crisp.
2. **State the context** — what problem forces this decision now? What constraint from `CONSTITUTION.md` or the codebase applies?
3. **List alternatives** — the options that were genuinely live. Each needs a one-line trade-off, not a strawman.
4. **Record the rationale** — why the chosen option wins, grounded in this project's context (not generic best-practice claims).
5. **State consequences** — what gets better and what gets harder. Both sides.
6. **Check related ADRs** — link any related or superseded decisions in `docs/adr/`. If this supersedes an existing ADR, update the old one's status to `Superseded`.
7. **STOP** — present the ADR. Wait for the user to confirm the decision and rationale.
8. **On approval** — write `docs/adr/ADR-YYYYMMDD-<slug>.md` and link it from `CONSTITUTION.md` or `docs/architecture.md` if appropriate.

---

## Checklist

- [ ] Decision stated in one sentence
- [ ] Context explains why now
- [ ] Alternatives are genuine, not strawmen
- [ ] Rationale grounded in project context
- [ ] Consequences cover both positive and negative
- [ ] Related ADRs cross-linked (superseded ones updated)
- [ ] User confirmed the decision and rationale
- [ ] ADR file written to `docs/adr/ADR-YYYYMMDD-<slug>.md`
- [ ] ADR linked from `CONSTITUTION.md` or `docs/architecture.md` if architectural

---

## Handoff

Terminal artifact: `docs/adr/ADR-YYYYMMDD-<slug>.md` (permanent)
Return to the skill that invoked this one — typically `4dc-plan` or `4dc-tdd-green`.

This skill is a utility, not a sequential phase. It is invoked on demand when a decision worth recording emerges during any phase.