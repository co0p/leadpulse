# Constitution — Team Impact Scorecard

A monthly decision-support tool for Team Leads. Standalone desktop application. No external services, no network calls, all data on-device.

---

## Engineering Principles

- **Small, focused changes.** Each change has one reason and one scope. Large changes are broken into ordered, independently verifiable steps.
- **Behavior first.** Every change to production behavior has a test that captures the intent before the code changes.
- **Determinism is non-negotiable.** The same inputs must always produce the same outputs. Formula engine functions are pure; no side effects.
- **Privacy by design.** No telemetry, no analytics, no external API calls. Data never leaves the device.
- **Dependency direction is a hard rule.** `ui` → `service` → `store` → `engine`. No layer may import from a layer above it.
- **Service is the API contract.** All use cases live in `service/`. The UI is replaceable; a CLI, web server, or other client can call the same `service/` interfaces.
- **Engine is the source of truth.** All scoring, normalization, alert, and trend logic lives in `engine/`. No formula logic in `service/`, `store/`, or `ui/`.
- **No gold-plating.** Build what the current increment requires. Defer generalization until a second use case exists.
- **Human review always.** The tool surfaces recommendations; a human makes every decision. The UI must never present a score as a verdict.

---

## Architecture Boundaries

Hexagonal architecture (Ports & Adapters). Business logic independent of infrastructure.

```
server/       — HTTP handlers (adapters). Request/response mapping. Injects use cases. Zero business logic.
core/         — Application use cases (e.g., AddMemberUseCase). Business logic, aggregates, port interfaces. Zero infrastructure dependencies.
                └─ Ports (interfaces): MemberRepository (defined in core/, implemented elsewhere)
storage/      — Adapters. Concrete implementations of ports (e.g., SQLiteTeamMemberRepository). Persistence only.
engine/       — Pure functions. Scoring, normalization, alerts, trends. No I/O.
```

**Dependency Direction (One-Way Inbound):**
- `core/` depends only on `engine/` types and Go stdlib. Zero HTTP, storage, or framework imports.
- `server/` (HTTP) depends on `core/` use cases. Injects repository adapters at construction.
- `storage/` adapters depend on `core/` port interfaces. Implement persistence logic only.
- `engine/` imports nothing from this project.
- No circular imports.

**Rules:**
- `core/` defines aggregates, value objects, use cases, and port interfaces. All business rules live here. Testable in isolation with in-memory adapters.
- `server/` handlers are thin adapters: parse HTTP request → call use case → map result/error to HTTP response. Unit testable with `httptest` and mocked use cases.
- `storage/` adapters implement port interfaces from `core/`. Multiple adapters can coexist (e.g., in-memory for tests, SQLite for production). Swappable at composition root.
- `service/` and coordinator patterns are retired in favor of use cases with explicit dependencies (ports). Existing `service/` functions for future features (e.g., `ScoringService`, `TrendService`) continue to use the service layer pattern.
- Composition root in `main.go` wires adapters + use cases + handlers. No reflection or DI container. Explicit injection.
- SQLite file lives in the OS user data directory (`os.UserConfigDir()`).
- Single binary distribution — no installer, no runtime dependencies.
- **Presentation is replaceable:** any presentation layer (HTTP SPA, CLI, future UI) can inject the same use cases and port adapters.

**Technology decisions:**
- Language: Go
- Database: SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- Distribution: compiled binary per platform via GitHub Actions
- Web frontend: HTMX + Alpine.js + Go `html/template`, served from Go binary via `embed.FS`
- Architecture: Hexagonal (Ports & Adapters) for application logic; see `docs/adr/ADR-20260916-hexagonal-architecture-member-crud.md`

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
