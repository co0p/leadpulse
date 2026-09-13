# Constitution — Team Impact Scorecard

A monthly decision-support tool for Team Leads. Standalone desktop application. No external services, no network calls, all data on-device.

---

## Engineering Principles

- **Small, focused changes.** Each change has one reason and one scope. Large changes are broken into ordered, independently verifiable steps.
- **Behavior first.** Every change to production behavior has a test that captures the intent before the code changes.
- **Determinism is non-negotiable.** The same inputs must always produce the same outputs. Formula engine functions are pure; no side effects.
- **Privacy by design.** No telemetry, no analytics, no external API calls. Data never leaves the device.
- **Dependency direction is a hard rule.** `ui` → `store` → `engine`. No layer may import from a layer above it.
- **Engine is the source of truth.** All scoring, normalization, alert, and trend logic lives in `engine/`. No formula logic in `ui/` or `store/`.
- **No gold-plating.** Build what the current increment requires. Defer generalization until a second use case exists.
- **Human review always.** The tool surfaces recommendations; a human makes every decision. The UI must never present a score as a verdict.

---

## Architecture Boundaries

Three layers. One direction.

```
ui/       — Fyne views, widgets, event handlers. No business logic.
store/    — SQLite persistence. Reads/writes domain types. No formula logic.
engine/   — Pure functions. Scoring, normalization, alerts, trends. No I/O.
```

**Rules:**
- `ui` may import `store` and `engine`.
- `store` may import `engine` types only (no engine computation).
- `engine` imports nothing from this project.
- No circular imports.
- SQLite file lives in the OS user data directory (`os.UserConfigDir()`).
- Single binary distribution — no installer, no runtime dependencies.

**Technology decisions:**
- Language: Go
- UI framework: Fyne v2 (`fyne.io/fyne/v2`)
- Database: SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- Distribution: compiled binary per platform via GitHub Actions

See `docs/architecture.md` for the C4 Level 2 container view.  
See `docs/adr/` for the decisions behind these choices.

---

## Testing Strategy

See `docs/testing.md` for full testing practices, commands, and evidence requirements.

**Summary:**
- `engine/` — unit tests for every formula, normalization, alert threshold, and trend calculation. These are the primary correctness gate.
- `store/` — integration tests against in-memory SQLite. Cover all read/write paths and constraint rules.
- `ui/` — manual verification only. Fyne does not have a reliable headless test driver.
- Gate: all unit and integration tests must pass before any change is considered complete.

---

## Performance Envelope

| Operation | Target |
|---|---|
| Formula recompute (100-member team) | < 100ms |
| Screen render / navigation | < 200ms |
| SQLite read for 24-month history (single member) | < 50ms |
| Export to CSV (full team, 12 months) | < 2s |

Performance is validated manually during development. No automated performance test suite in v1.

---

## Documentation And ADR Policy

- `docs/architecture.md` — C4 Level 2 view. Updated when a container or major dependency changes.
- `docs/domain.md` — domain glossary. Updated when new domain concepts are introduced.
- `docs/ui.md` — UI decisions, interaction patterns, visual conventions. Updated when a new pattern is introduced.
- `docs/testing.md` — testing approach and commands. Updated when strategy changes.
- `docs/deployment.md` — release procedure and rollback. Updated when the release process changes.
- `docs/adr/` — one file per architectural decision. Decisions are recorded when structural, hard-to-reverse, or non-obvious choices are made. Index is the directory listing. See existing ADRs for format.

Do not record implementation details in ADRs. Record the decision, context, alternatives, rationale, and consequences.

---

## Release And Deployment

See `docs/deployment.md` for the full release procedure and rollback guidance.

**Summary:**
- Release unit: a single compiled binary per platform (macOS `.app`, Windows `.exe`, Linux binary).
- Trigger: manual — tag a commit `v<major>.<minor>.<patch>`, push to origin; GitHub Actions builds and publishes to GitHub Releases.
- Versioning: Semantic Versioning. Breaking data-schema changes require a major bump.
- Rollback: user downloads a previous release binary. SQLite schema migrations must be additive-only in v1; no destructive migrations.
- No server to roll back. No deployment target other than the user's own machine.

---

## Delivery and Documentation

- A feature is **done** only when its acceptance criteria pass and evidence is linked in `docs/roadmap.md`.
- Formula changes require a corresponding engine unit test update before the change is merged.
- UI patterns introduced for the first time must be documented in `docs/ui.md`.
- New domain concepts must be added to `docs/domain.md` before or alongside their first use in code.
- The PRD (`docs/prd.md`) is the source of truth for scoring formulas, alert thresholds, and acceptance criteria. Code must match it exactly; discrepancies are bugs.
