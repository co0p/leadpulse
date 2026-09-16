---
name: 4dc-constitution
description: "Use when CONSTITUTION.md is missing or needs updating. Reads project context, asks focused questions, and produces project guardrails and SDLC standards."
---

# Constitution Skill

## One Responsibility

Create or update `CONSTITUTION.md` — the project's durable engineering guardrails.

The constitution states **principles and boundaries**, not implementations. It answers: what rules do we agree to operate under, what boundaries must not be crossed, and where does each category of concrete detail live? It does not prescribe tools, libraries, frameworks, file names, commands, or configuration — those belong in `docs/` or ADRs.

The generated constitution must describe only the application. Do not include this repository's internal workflow name, phase sequence, agent names, or transient artifact paths.

Before writing it, check that it contains no internal workflow names, skill names, phase names, orchestrator terms, `.agent/` paths, or `.agents/` paths.

---

## Foundations

These three traditions are the source for guardrails. Borrow from them by name — attribute the rule to its author so the team knows why it exists, not just what it says.

- **Beck — team agreements before code.** The constitution is the set of rules the team agrees to operate under. It is written before implementation, not retrofitted after. Rules must be specific enough to resolve disputes; vague principles are not guardrails.
- **Poppendieck — eliminate ambiguity upstream.** Decide the guardrails early so later phases do not rediscover the same constraints. For reversible choices, leave the decision late; for structural rules, fix them now. Pull rules from value; do not mandate process that does not serve delivery.
- **Fowler — evolutionary architecture.** Guardrails, not blueprints. The constitution sets the fitness functions and boundaries that let the design evolve safely. It does not fix the implementation — it defines what must remain true as the implementation changes.

**What belongs in CONSTITUTION.md (principles and boundaries):**
- Engineering principles attributed to their source (Beck, Poppendieck, Fowler, or a project-specific decision)
- Architectural boundaries: which direction dependencies flow, which containers must stay decoupled, what crosses the system boundary
- Testing strategy: what must have tests, what constitutes a green gate before promote, and what kinds of tests are in or out of scope — as rules, not commands
- Performance envelope: latency, throughput, cost, or scale expectations as stated constraints
- Documentation policy: what is permanent, what is transient, where each category lives
- ADR policy: what kinds of decisions require an ADR

**What does NOT belong in CONSTITUTION.md (belongs in `docs/` or ADRs):**
- Test commands, CI scripts, tooling configuration
- Deployment runbooks, environment lists, secret-handling procedures
- Framework choices, library names, file naming conventions
- Coverage numbers, specific thresholds, or tool-specific configuration
- Architecture diagrams or container inventories

---

## Expected Input

- Existing `CONSTITUTION.md` (if present)
- `README.md`
- Current project structure and any existing docs

---

## Concrete Output

### Primary Artifact: `CONSTITUTION.md`

`CONSTITUTION.md` containing only principles and boundaries — no commands, tooling, or concrete procedures:

1. **Engineering principles** — stated as rules, attributed to Beck, Poppendieck, Fowler, or a project decision. Each rule must be justified by this project's context, not copied from generic advice.
2. **Architectural boundaries** — dependency direction, which containers must remain decoupled, performance-critical paths stated as constraints. No container inventory (that lives in `docs/architecture.md`).
3. **Testing strategy** — what must have tests, what constitutes a green gate before promote, and what kinds of tests are in or out of scope — as rules. Points to `docs/testing.md` for procedures and commands.
4. **Performance envelope** — latency, throughput, cost, or scale expectations as stated constraints, or an explicit `N/A` with rationale. No monitoring configuration.
5. **Documentation and ADR policy** — what is permanent, where each category lives, and what decisions trigger an ADR.
6. **Release and deployment** — the release model as a rule (e.g. "every merge to main is releasable"). Points to `docs/deployment.md` for procedures.

Required `CONSTITUTION.md` headings:
- `## Engineering Principles`
- `## Architecture Boundaries`
- `## Testing Strategy`
- `## Performance Envelope`
- `## Documentation And ADR Policy`
- `## Release And Deployment`

### Supporting Documents

**`docs/testing.md`** — testing practices for this project:
- The reasoning behind the testing strategy and the risks it is intended to control
- Guidance for choosing test depth at architectural and user-facing boundaries
- Actual commands for local, CI, acceptance, and release checks, with setup and interpretation guidance
- Evidence required before a change is considered complete
- Known confidence gaps, maintenance practices, and reasons for environment-specific checks
- No inventory of individual tests, test-case lists, or coverage targets unless a project-specific decision genuinely requires one

**`docs/deployment.md`** — deployment and release procedures for this project:
- The deployment model, its rationale, and the operational assumptions it relies on
- Release triggers, versioning decisions, ownership, and required evidence
- Actual deployment and rollback runbooks with verification and recovery guidance
- Configuration and secret-handling principles without secret values
- Health signals, alert actions, and meaningful operational risks
- No historical release log or generic checklist detached from this project's procedure

**`docs/adr/`** — Architecture Decision Records:
- Decisions with rationale and consequences
- The `docs/adr/` directory is the index; do not duplicate that list in `CONSTITUTION.md`
- Created using the ADR template when foundational decisions exist
- Updated when architectural decisions emerge
- Each ADR explains the context, decision, alternatives, trade-offs, and consequences; it is not an implementation diary

**`docs/architecture.md`** — C4 architecture view:
- Required for every project, even when the system is small
- Must contain a current C4 Level 2 container view (or an equivalent explicitly labeled diagram)
- Describes runtime containers, responsibilities, communication paths, and data stores

**`docs/domain.md`** — Domain glossary:
- Required for every project, even when the vocabulary is initially small
- Defines shared concepts, domain events, and system rules in business language
- Must not be replaced by an ADR, README, or implementation-specific notes

**`docs/ui.md`** — Permanent UI decisions (required when the project has a user interface):
- Shared user flows, interaction patterns, visual principles, accessibility rules, and content conventions
- Rationale and consequences of recurring UI decisions
- No component inventory, CSS catalog, or one-off screen notes

### Secondary Artifact

`docs/roadmap.md` (created from the template in the Appendix if it does not exist yet)

### Documentation Baseline

Before creating or updating the constitution, audit the repository for the complete permanent documentation baseline:

- `CONSTITUTION.md`
- `docs/testing.md`
- `docs/deployment.md`
- `docs/architecture.md` with a C4 Level 2 container view
- `docs/domain.md` with the project's glossary
- `docs/ui.md` with the project's UI decisions (if the system has a UI; omit for headless systems)
- `docs/adr/`
- `docs/roadmap.md`

Missing baseline documents are constitution outputs; they are not optional follow-up work. Existing documents must be checked for the required content rather than accepted solely because the path exists.

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
Do NOT write `CONSTITUTION.md` until the user explicitly approves the proposed guardrails.
Do NOT ask more than 5 questions per round.
Do NOT include implementation details — CONSTITUTION.md contains guardrails, not recipes. If a sentence contains a tool name, a command, a file path, a framework, a library, or a configuration value, it belongs in docs/ or an ADR, not in CONSTITUTION.md.
Do NOT copy generic principles from the internet. Every rule must be justified by this project's specific context.
Do NOT add a `## Delivery and Documentation` section — documentation policy belongs under `## Documentation And ADR Policy`.
</HARD-GATE>

---

## Process

1. **Read project context** — scan `README.md`, existing `CONSTITUTION.md`, directory structure, any ADRs or docs, and existing deployment or testing practices
2. **Conversation: Propose the guardrails** — summarize the proposed testing and deployment strategies, identify foundational ADR candidates, and ask whether the direction feels right. Iterate until the user says “looks good” or “proceed.”
3. **On approval:**
   - Write `CONSTITUTION.md` with references to supporting documents
   - Create `docs/testing.md` with project-specific testing practices
   - Create `docs/deployment.md` with project-specific deployment procedures
    - Create or update `docs/architecture.md` with the current C4 Level 2 container view
    - Create `docs/domain.md` with the project's initial glossary, even if only a few concepts are known
    - Create `docs/ui.md` with the project's initial UI decisions when the system has a user interface
   - Create initial `docs/adr/` structure if foundational decisions exist
   - Create `docs/roadmap.md` if not present

---

## Checklist

- [ ] Existing docs read (including any deployment or testing practices)
- [ ] Foundational ADRs identified (if any exist)
- [ ] Proposed testing strategy, deployment strategy, and ADR candidates discussed
- [ ] User approval received
- [ ] `CONSTITUTION.md` written with Testing, Performance, and Release sections populated (with references to supporting docs)
- [ ] `docs/testing.md` created with project-specific practices
- [ ] `docs/deployment.md` created with project-specific procedures
- [ ] `docs/adr/` directory created with index link from `CONSTITUTION.md` (populate with foundational decisions if identified)
- [ ] `docs/roadmap.md` created if not present

---

## Handoff

Terminal artifacts:
- `CONSTITUTION.md` — guardrails and governance
- `docs/testing.md` — testing procedures and practices
- `docs/deployment.md` — deployment and release procedures
- `docs/adr/` structure — architectural decisions
- `docs/roadmap.md` — product roadmap

Future work should reference the testing, deployment, and architecture documents defined by the project constitution. New ADRs should be added when architectural decisions emerge.

After this skill, continue with the next approved work item using the applicable project workflow.

---

## Appendix: Document Templates

Use these verbatim as the starting content when creating a new document for the first time.

### Template: docs/testing.md

```markdown
# Testing

Guide to making reliable testing decisions for this project. Explain why the test strategy is shaped this way, how a developer should choose the cheapest test that gives sufficient confidence, and how to run the relevant checks. Update when the architecture, risk profile, or test workflow changes.

---

## Testing Approach and Rationale

Explain the risks the test strategy is designed to control and the boundaries where each kind of test provides confidence. Prefer principles and decision guidance over inventories. For example, explain why domain rules are tested without infrastructure, why persistence boundaries need integration checks, or why an acceptance test exercises a complete user job story.

---

## Choosing Test Depth

<!--
Describe how to decide whether a change needs a focused check, an integration check, an acceptance scenario, a performance measurement, or no new test. Tie the decision to user risk, architectural boundaries, determinism, and failure cost.
Do not maintain a catalog of individual tests or report a coverage percentage here.
-->

---

## Test Design Conventions

<!--
Describe conventions that make tests communicate behavior: naming, fixture ownership, isolation, determinism, test data, and how user-facing assertions should avoid implementation details. Keep examples small and illustrative rather than listing the suite.
-->

---

## Running the Checks

<!--
Document the actual commands for fast local feedback, the complete pre-merge gate, and any setup required for acceptance or environment-dependent checks. Explain when to use each command and how to interpret failures. Commands must be maintained as executable guidance, not illustrative placeholders.
-->

---

## Evidence Required Before Merge

<!--
State the evidence required before a change is considered complete. Focus on behavior, risk, and reproducibility. Do not use line coverage or a list of passing tests as a substitute for explaining why the evidence is sufficient.
-->

---

## Automation and Feedback Loops

<!--
Explain where checks run (local, CI, release), what feedback each loop provides, and how failures are triaged. Record the rationale for any intentionally manual or environment-specific check.
-->

---

## Known Risks and Gaps

<!--
Document meaningful confidence gaps, why they exist, and what signal would justify changing the approach. Do not turn this section into a test inventory or coverage report.
-->

---

## Maintenance Guidance

<!--
Explain how tests are kept deterministic, how flakiness is handled, when fixtures or helpers should be changed, and how this guide itself is updated when the testing rationale changes.
-->
```

### Template: docs/deployment.md

```markdown
# Deployment

Guide to releasing and operating this project. Explain the deployment model, why it is appropriate, how to execute it safely, and how to recover. Update whenever the release or operational model changes.

Describe the release unit, the target, ownership, and the operational assumptions. Explain why this deployment shape is used.
- This is a Node.js backend service deployed to AWS Lambda
- This is a React frontend deployed to Vercel
- This is a Go CLI tool distributed via Homebrew and GitHub Releases
-->

Explain the release trigger and versioning decision, including who can release and what evidence is required first. Avoid a changelog or release-history list here.
- Manual: Tag a commit with `v<major>.<minor>.<patch>` and push to origin; CI builds and publishes
- Automatic: Merge to `main` triggers a release with semantic versioning based on conventional commits
- Versioning scheme: Semantic Versioning (SemVer) for libraries, CalVer for applications
-->
Describe the environments and their purpose, including the differences that matter for safe verification. Do not use this section as an environment inventory without explaining the deployment model.
Or for CLI:
- macOS: Homebrew tap `example/tap/tool`
- Linux: GitHub Releases, apt repository
- Windows: GitHub Releases, Scoop bucket
-->

---

## Deployment Procedure

Document the actual release procedure as a short runbook, with prerequisites, commands or links, verification signals, and ownership. Explain why the ordering protects users and data. Keep the checklist operational; put rationale in surrounding prose.
Describe rollback triggers, the recovery action, data implications, and who decides. Explain any forward-only or irreversible operation and the recovery alternative.

For CLI releases:
- Yanked versions: Tag with `v<version>` and mark as yanked in release notes
- Users on old version: Keep supporting previous major version for [N] months
-->

---

## Deployment Checklist

Keep only the small set of release decisions and checks that are specific to this project. Do not turn this into a repeated list of every test, deployment, or release ever performed.

Explain configuration ownership, secret handling, safe defaults, and the reason for separating deploy-time configuration from source code. Never record secret values.
- Environment variables are injected at deployment time from [source]
- Configuration file: `.env.production` (not checked in), managed by [process]
- Database connection string: Retrieved from [secrets manager] at startup
Explain how operators know a release is healthy, which signals matter, and what action an alert should trigger. Record thresholds only when they are real, justified, and maintained.
- Error rates are monitored in Datadog; alert if error rate > 1% for 5 min
- Latency p99 is tracked; alert if > [threshold]
- Database connection pool is monitored; alert if exhausted
- Disk space is monitored; alert if < 10% free
- Deployment notifications sent to [Slack channel]
-->

Document constraints that materially affect release safety and the mitigation or follow-up needed. Do not preserve obsolete procedures as historical reference.
- Deployment is not atomic: old and new code may run simultaneously for [duration]
- Secrets rotation requires [manual step]
- Large deployments > 100MB take [N] minutes; monitor for timeout
-->
```

### Template: docs/adr.md

When creating a new Architecture Decision Record, use this template:

```markdown
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

{{SHARED:execution-contract}}

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
```

### Template: docs/roadmap.md

```markdown
# Roadmap

Product direction and sequencing guide. Each entry explains the user outcome, current confidence, and ordering rationale. Keep it concise and decision-oriented; detailed implementation status belongs in phase artifacts and code evidence.

> A feature moves to **Done** only when its user outcome is verified and the evidence is linked here.
> Source of truth: if a feature is not in Done with a passing test link, it is not considered shipped.

---

## Done

<!--
Each entry follows this pattern:

### [Feature name — short, user-visible]
- **Job story:** When [situation], I want to [action], so that [outcome].
- **Evidence:** [link to the smallest durable verification record](path/to/evidence)
- **Use case:** [docs/usecases/use-case-slug.md](docs/usecases/use-case-slug.md) *(if promoted)*
- **Delivered:** [increment slug or YYYY-MM-DD]
-->

---

## In Progress

<!--
### [Feature name]
- **Job story:** When [situation], I want to [action], so that [outcome].
- **Evidence:** pending — define the verification approach before work starts
-->

---

## Planned

<!--
### [Feature name]
- **Job story:** When [situation], I want to [action], so that [outcome].
- **Why now / ordering:** [dependencies, user value, open questions, or sequencing rationale]
-->

---

## How This List Works

- Features move left to right: Planned → In Progress → Done. Never skip In Progress.
- A feature enters In Progress when work begins.
- A feature enters Done only when its user outcome is verified and the evidence link is present.
- Do not add implementation detail here — link to the use case or ADR for that.
- If a planned feature is no longer needed, remove it and record the reason in a code comment, commit message, or ADR.
```