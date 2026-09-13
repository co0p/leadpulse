---
name: 4dc-promote
description: "Use after implementation.md is marked complete. Reviews all .agent/ artifacts, proposes promotions to permanent docs, and closes the cycle."
---

# Promote Skill

## One Responsibility

Merge durable outcomes from the `.agent/` working set into permanent project artifacts before the branch merges.

---

## Foundations

- **Poppendieck — eliminate waste.** Promote only what was verified, not what was planned. Unverified work is waste; it does not earn a place in permanent docs.
- **Beck — retrospective embedded in delivery.** The cycle's learnings are not an afterthought; they are part of the deliverable. Decisions, deviations, and surprises feed forward into the project's durable knowledge.
- **Fowler — documentation as architecture.** Durable docs are part of the system, not a record about it. When the architecture, domain language, or performance envelope changes, the docs change in the same cycle — otherwise they rot.

---

## Expected Input

- `CONSTITUTION.md`
- `.agent/increment.md`
- `.agent/plan.md`
- `.agent/implementation.md` (status: complete)
- `.agent/learnings.md`
- Existing `docs/`, ADR log, and any project documentation

---

## Concrete Output

One or more of the following, per approval:
- Updated `CONSTITUTION.md` (if guardrails need revision)
- New ADR in `docs/adr/` (for significant architectural decisions)
- Updated `docs/architecture.md` (if runtime structure, dependencies, or performance-critical paths changed)
- Updated `docs/domain.md` (if domain language changed or new concepts appeared)
- Updated `docs/ui.md` (if shared UI, interaction, visual, accessibility, or content decisions changed)
- Updated `README.md` or other docs (for changed behavior or usage)
- Updated `docs/roadmap.md` — feature moved from Partial to Done, acceptance test link added
- Acceptance-scenario evidence — linked when available; advisory scenarios inform confidence but do not block promotion by default
- Deleted or archived `.agent/` files after promotion (keeping `.agent/` clean for next cycle)

Required outputs:
- Promotion candidates are listed individually with destination path, rationale, and approval status.
- The promotion explicitly states whether architecture, domain language, testing guidance, and performance documentation changed or stayed unchanged.

### Permanent Documentation Baseline

Every promotion must verify that the project's permanent documentation baseline exists and is usable:

- `CONSTITUTION.md`
- `docs/testing.md`
- `docs/deployment.md`
- `docs/architecture.md` containing a current C4 Level 2 container view (or an explicitly labeled equivalent)
- `docs/domain.md` containing the current domain glossary
- `docs/ui.md` containing current UI decisions (if the system has a UI; omit for headless systems)
- `docs/adr/`
- `docs/roadmap.md`

The check is semantic, not just a file-existence check. A generic architecture narrative does not satisfy the C4 requirement, and a glossary hidden in an ADR or README does not satisfy the domain-document requirement. Missing or inadequate documents become promotion candidates and must be created or corrected before the cycle can close.

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
Do NOT write permanent docs until each promotion candidate has been individually approved.
Do NOT promote guesses or plans — only promote what was actually built and verified.
Do NOT delete .agent/ files until all promotions are written and confirmed.
Present each candidate separately with destination path and rationale.
Do NOT leave permanent docs stale when the implementation changed architecture, domain language, or performance-critical behavior.
Do NOT close promotion while any required permanent documentation baseline item is missing or inadequate. If runtime structure, domain vocabulary, and UI decisions are unchanged, still verify that `docs/architecture.md`, `docs/domain.md`, and `docs/ui.md` (when applicable) exist and satisfy their requirements.
</HARD-GATE>

---

## Process

1. **Read all `.agent/` artifacts** — full review of increment, plan, implementation, and learnings.
2. **Audit the permanent documentation baseline** — inspect each required path and verify the architecture document contains a C4 Level 2 container view and the domain document contains the glossary. Add missing or inadequate documents to the candidate list.
3. **Conversation: Propose promotions** — identify candidates and state what each is, its destination, and why it is durable. Iterate until the user says to proceed.
4. **On approval** — write each approved permanent artifact.
5. **Re-audit the baseline** — confirm every required document exists and satisfies its content requirement before cleanup.
6. **Clean up** — archive or delete `.agent/` files for this cycle.

---

## Promotion Categories

| Type | Trigger | Destination |
|------|---------|-------------|
| Architecture decision | Non-obvious choice with lasting impact | `docs/adr/ADR-<date>-<slug>.md` |
| Guardrail update | Constitution rule violated, needs clarification | `CONSTITUTION.md` |
| Architecture sync | Runtime containers, dependency direction, or performance-critical paths changed | `docs/architecture.md` |
| Behavior change | Public API, CLI, or user-facing behavior changed | `README.md` |
| Feature shipped | Acceptance tests pass; feature complete | `docs/roadmap.md` — move to Done, add acceptance test link |
| Acceptance evidence | Optional user-journey scenario was run | `docs/roadmap.md` or implementation evidence, linked when useful; not a default gate |
| Test pattern | New testing approach worth standardizing | `CONSTITUTION.md` testing section |
| Performance contract | A latency, throughput, cost, or scaling expectation changed | `CONSTITUTION.md` or `docs/architecture.md` |
| Known issue | Found but not fixed this cycle | `docs/known-issues.md` |
| New domain concept | A concept, event, or rule used in code/tests that has no shared definition | `docs/domain.md` (create using the template in the Appendix if absent) |
| UI decision | Shared UI, interaction, visual, accessibility, or content decision | `docs/ui.md` |
| Structural change | A container added, removed, or re-wired | `docs/architecture.md` (create using the template in the Appendix if absent) |

---

## Checklist

- [ ] All `.agent/` artifacts read
- [ ] Promotion candidates identified and categorized
- [ ] Permanent documentation baseline audited for existence and required content
- [ ] User approval received per candidate
- [ ] Each approved artifact written to permanent location
- [ ] Final baseline audit passes: glossary, C4 architecture view, and UI decisions (when applicable) are present and current
- [ ] Architecture, domain language, testing guidance, and performance documentation either updated or explicitly marked unchanged
- [ ] Acceptance scenarios, if present, have results recorded; advisory failures or unavailable scenarios are documented without blocking by default
- [ ] `.agent/` files cleaned up

---

## Handoff

Terminal artifacts: permanent docs updated, `.agent/` clean
Cycle complete. Next action: `4dc-increment` for the next cycle — load `skills/increment/SKILL.md`

---

## Appendix: Document Templates

Use these verbatim as the starting content when creating a new document for the first time.

### Template: docs/domain.md

```markdown
# Domain Vocabulary

Shared business language for this project. Use this as a reference when requirements, code, tests, and conversations use the same concept. Update when meaning, rules, or relationships change; do not use it as a data dictionary or implementation catalog.

---

## Concepts

### [ConceptName]

**Definition:** [One sentence in domain terms].

**Meaningful state:**
- [Only state that changes how people understand or use the concept]

**Rules:**
- [Invariant or constraint in domain language]

**Related:**
- [Concepts or events that clarify the meaning]

---

## Domain Events

### [EventName]

**When:** [What triggers this event in domain terms]
**Information carried:** [Business information, not an implementation payload schema]
**Consumers:** [Who or what reacts in the domain]

---

## Rules and Constraints

[System-wide invariants, state transitions, and cross-concept rules described in domain language.]

---

## Writing Rules

- Define terms in language a product owner and implementer can both use.
- Record rules and relationships only when they affect decisions or outcomes.
- Do not list database columns, code paths, test cases, APIs, or historical introductions.
- If a term is still uncertain, record the ambiguity in the project decision record instead of inventing a definition.
```

### Template: docs/architecture.md

```markdown
# Architecture — [Project Name]

C4 Level 2: Container diagram. Updated when structural boundaries change.

> This is a durable orientation guide, not a component inventory. It explains the system boundary, runtime containers, important communication paths, and constraints so a reader can reason about change safely.

---

## Context (C4 Level 1 summary)

**System:** [Project Name]
**Users:** [Who uses it — one line each]
**Purpose:** [One sentence: what problem does this system solve?]

---

## Containers

A container is any separately runnable or deployable unit: a process, a script, a service, a database.

| Container | Technology | Responsibility |
|-----------|-----------|----------------|
| [Container A] | [e.g. Go binary, bash script, Node.js process] | [What it does — one line] |
| [Container B] | [technology] | [responsibility] |
| [Container C] | [technology] | [responsibility] |

---

## Container Diagram

```
┌─────────────────────────────────────────────────────────┐
│  [System Name]                                          │
│                                                         │
│  ┌──────────────────┐         ┌──────────────────────┐  │
│  │  [Container A]   │──────▶  │  [Container B]       │  │
│  │                  │  [how]  │                      │  │
│  │  [technology]    │         │  [technology]        │  │
│  └──────────────────┘         └──────────────────────┘  │
│           │                             │                │
│           ▼                             ▼                │
│  ┌──────────────────┐         ┌──────────────────────┐  │
│  │  [Container C]   │         │  [External System]   │  │
│  │  [technology]    │         │  (out of scope)      │  │
│  └──────────────────┘         └──────────────────────┘  │
└─────────────────────────────────────────────────────────┘

         [User / Actor]
              │
              ▼ [interaction description]
         [Container A]
```

---

## Communication

| From | To | Protocol / Mechanism | Notes |
|------|----|----------------------|-------|
| [Container A] | [Container B] | [e.g. stdout pipe, HTTP, file read] | [any constraint] |
| [User] | [Container A] | [e.g. CLI args, browser] | |

---

## Data Stores

| Store | Type | Owned by | Schema / Format |
|-------|------|----------|-----------------|
| [Store A] | [e.g. SQLite file, CSV, in-memory] | [Container A] | [brief description] |

If no persistent data store exists, state that explicitly:
> This system is stateless. No persistent data store.

---

## Key Constraints

Constraints that affect all containers and must not be violated:

- [e.g. "No network calls — runs entirely offline"]
- [e.g. "Single binary, no install step"]
- [e.g. "All state is held in browser memory; nothing is written to a server"]

## Reading and Update Guidance

Explain the architectural reasoning that matters to contributors: why the containers are separated, which boundaries must remain stable, and what kinds of changes require an ADR or an update to this document. Do not duplicate class lists, endpoint lists, or deployment instructions.

---

## Out of Scope

Explicitly name what this diagram does NOT cover:

- Internal component structure of each container (C4 Level 3 — not written unless needed)
- Deployment topology
- CI/CD pipeline

---

## Update Policy

Update this file when:
- A container is added, removed, or its technology changes
- A communication path between containers changes
- A new external system dependency is added

Do NOT update for internal refactors, new features within an existing container, or test changes.

The diagram is a current model, not a historical record. Remove stale paths and obsolete containers rather than preserving them for context.

**Last updated:** [YYYY-MM-DD] — [brief reason]
```