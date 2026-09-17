# Constitution — Team Impact Scorecard

---

## Engineering Principles

- **Small, focused changes.** Each change has one reason and one scope. Large changes are broken into ordered, independently verifiable steps.
- **Behavior first.** Every change to production behavior has a test that captures the intent before the code changes.
- **Determinism is non-negotiable.** The same inputs must always produce the same outputs. Formula engine functions are pure; no side effects.
- **Privacy by design.** No telemetry, no analytics, no external API calls. Data never leaves the device.
- **Dependency direction is a hard rule.** `server` (or any client) → `core` → `engine`; adapters (`storage`) depend inward on `core` port interfaces. No layer may import from a layer above it.
- **Service is the API contract.** All use cases live in `core/`. Any client (SPA, CLI, future integrations) can call the same JSON API contract.
- **Presentation and API are separate systems.** The backend exposes data and behavior through a JSON API only. It never renders HTML, never owns layout or navigation state, and never has knowledge of how its output is displayed. Any client (SPA, CLI, future integrations) consumes the same contract.
- **Engine is the source of truth.** All scoring, normalization, alert, and trend logic lives in `engine/`. No formula logic in `service/`, `store/`, or `ui/`.
- **No gold-plating.** Build what the current increment requires. Defer generalization until a second use case exists.
- **Human review always.** The tool surfaces recommendations; a human makes every decision. The UI must never present a score as a verdict.

---

## Multi-Service Architecture

**Backend Service** (`services/backend/`)
- Go application with hexagonal architecture (core, storage, engine, server)
- Exposes JSON REST API; all use cases accessible through stateless HTTP endpoints
- SQLite database (file volume in docker-compose)
- No HTML rendering; no server-side state

**Frontend Service** (`services/frontend/`)
- Vue 3 + Vite SPA (separate npm project with own `package.json`)
- Calls backend exclusively through `/api/*` endpoints
- Served via lightweight web server on standard port
- Static assets built at container build time; no server-side logic

**Orchestration** (`docker-compose.yml`)
- Boots backend, frontend, and database in a local development network
- Services discover each other by hostname (e.g., frontend calls `http://backend:8080/api/*`)
- All services must pass health checks before accepting requests

**Dependency Direction (One-Way Inward):**
- Frontend (client) depends on backend JSON API; zero knowledge of backend internals
- Backend (server) depends on core/storage/engine as before; hexagonal architecture unchanged
- No circular imports or cross-service coupling

---

**Technology decisions:**
- Language: Go
- Database: SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- Distribution: compiled binary per platform via GitHub Actions
- Frontend: a dedicated single-page application, built separately from the backend and communicating exclusively through the JSON API. The compiled frontend is embedded into the Go binary via `embed.FS` to preserve single-binary distribution. See `docs/adr/ADR-20260917-vue-spa-frontend.md` for the framework choice and rationale.
- Architecture: Hexagonal (Ports & Adapters) for application logic; see `docs/adr/ADR-20260916-hexagonal-architecture-member-crud.md`

See `docs/architecture.md` for the C4 Level 2 container view.  
See `docs/adr/` for the decisions behind these choices.

---

## Testing Strategy

See `docs/testing.md` for full testing practices, commands, and evidence requirements.

**Summary:**

**Backend (Go services/backend/):**
- `engine/` — unit tests for every formula, normalization, alert threshold, and trend calculation. These are the primary correctness gate.
- `core/` — unit tests with in-memory repository adapter. All business logic testable without database.
- `storage/` — integration tests against in-memory SQLite. Cover all read/write paths and constraint rules.
- `server/` — handler tests against the JSON contract only (no HTML assertions). The API is the boundary; tests never depend on how a client renders the response.
- Tests run inside container during CI or locally; must pass before container build succeeds.

**Frontend (Vue services/frontend/):**
- Component tests via Vitest + Vue Test Utils. Unit-test rendering, props, emitted events, form validation, error states without a browser.
- API client tests verify request/response mapping.
- Tests run locally; no browser required for component test layer.

**Acceptance Tests (acceptance-tests/):**
- Playwright tests run against live docker-compose environment (all services must be healthy).
- Verify end-to-end user flows: app loads, health check works, CRUD operations complete, data persists.
- Reserved for main success flows only; edge cases and validation errors tested at component/handler layer.
- Each feature ships with at least one acceptance test verifying its primary outcome.

**Gate before merge:**
- Backend: `go test -race ./...` passes
- Frontend: `npm test` passes
- Acceptance: `npm run test:acceptance` passes against running docker-compose environment

**Gate before release:** All three categories pass; acceptance tests verify integration end-to-end.

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

**Summary (Multi-Service Docker Architecture):**

**Local Development:**
- Development: Run `docker-compose up` to boot backend, frontend, and database containers. Services discover each other by hostname.
- Tests: Run backend tests locally (`go test -race ./...`), frontend tests locally (`npm test`), and acceptance tests against running docker-compose environment (`npm run test:acceptance`).

**CI:**
- Build backend container (test inside container, then build binary)
- Build frontend container (npm build Vue SPA, package with web server)
- Run acceptance tests against composed containers
- All tests must pass before merge

**Release (Future; Out of Scope for v1):**
- Tag commit `v<major>.<minor>.<patch>`, push to origin
- CI builds and publishes backend + frontend container images to registry
- Versioning: Semantic Versioning. Breaking data-schema changes require a major bump.

**Rollback:**
- Re-run `docker-compose up` with previous container image tags
- SQLite schema migrations must be additive-only in v1; no destructive migrations
- No stateful containers; data persists only in SQLite volume

**Key Constraints:**
- Both services must pass health checks before accepting requests
- No external service dependencies; all services run on local network
- Database is a file volume (SQLite); no cloud or managed database

---

## Delivery and Documentation

- A feature is **done** only when its acceptance criteria pass and evidence is linked in `docs/roadmap.md`.
- Formula changes require a corresponding engine unit test update before the change is merged.
- UI patterns introduced for the first time must be documented in `docs/ui.md`.
- New domain concepts must be added to `docs/domain.md` before or alongside their first use in code.
- The PRD (`docs/prd.md`) is the source of truth for scoring formulas, alert thresholds, and acceptance criteria. Code must match it exactly; discrepancies are bugs.
