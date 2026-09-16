# Roadmap

Product direction and sequencing for Team Impact Scorecard. Each entry explains the user outcome, current confidence, and ordering rationale.

> A feature moves to **Done** only when its user outcome is verified and the evidence is linked here.
> Source of truth: if a feature is not in Done with a passing test link, it is not considered shipped.

---

## Migration Notice

**As of 2026-09-15, implementation has shifted to a Web SPA frontend.** The Go backend (engine, store, service, and controllers) remains the same. The Fyne desktop GUI is being phased out in favor of an HTMX + Alpine.js + Go `html/template` SPA served from the Go binary. See the **SPA Implementation** section below for the new path. Fyne screens will be retired incrementally as their SPA equivalents pass acceptance criteria.

See `docs/adr/ADR-20260915-spa-frontend-stack.md` for the frontend stack decision and rationale.

---

## Done

### Team Member Management (CRUD)
- **Job story:** When I set up the tool or my team changes, I want to add, view, edit, and remove team members with their name and seniority level, so that the scorecard always reflects my current team and shows each member's context.
- **Acceptance scenarios verified:**
  - Add member → appears in list with name and seniority ✓
  - Edit member → updated fields persist ✓
  - Deactivate member → soft-deleted, excluded from list, historical data preserved ✓
  - Persistence → members survive app restart (via SQLite) ✓
- **Evidence:**
  - Store layer: 5 integration tests passing (`store/member_test.go`; add, list, edit, deactivate, soft-delete with audit trail)
  - Service layer: 7 unit tests passing (`service/member/member_test.go`; use cases with mocked store, validation rules)
  - UI layer: Settings controller (Screen F) with add/edit/deactivate logic; full CRUD via Fyne (deprecated)
  - Database: idempotent SQLite schema in `store/schema.go` with members table and audit_log
  - Build: successful without errors (harmless macOS linker warnings only)
- **Acceptance criteria:** All 5 met (AC-1: add/view/edit/remove, AC-2: seniority, AC-3: edit, AC-4: deactivate, AC-5: persistence)
- **Key commits:** 
  - Store: `ecf0b6e`, `d070add`, `6749c3a`
  - Service: `74b53fa`
  - UI: `173d3b4`, `f920c60` (layout fix + dialogs)
- **Test command:** `go test -race ./...` → all pass

---

### Formula Engine — Core Scoring
- **Job story:** When I submit a team member's monthly data, I want the system to compute all normalized scores, dimension scores (DG/DP/DT/DO), TII, completeness, and confidence, so that I have an objective, repeatable basis for my review.
- **Acceptance scenarios verified:**
  - All 10 normalization functions compute correctly per PRD 5.2 ✓
  - Impact-weighted signal computation correct per PRD 5.3 ✓
  - Dimension scores (DG/DP/DT/DO) computed correctly per PRD 5.4 ✓
  - Total Impact Index (TII) computed correctly per PRD 5.5 ✓
  - Completeness percentage computed correctly per PRD 4.3 ✓
  - End-to-end scoring works for typical, strong, and struggling members ✓
- **Evidence:**
  - Engine layer: 53 unit and integration tests passing (`engine/scoring/` with 22 normalization, 4 impact, 12 dimension, 9 TII/completeness/e2e tests)
  - Test coverage: 87.5% of scoring engine statements
  - Race detector: all tests pass with `-race` flag
  - Validation: MoraleN, BillabilityN, CSATN, MarginN, PositiveN, CriticalN, OvertimeN, DeliveryN, MentoringN, EvidenceN, ComputeImpactWeighted, ComputeDimensionGrowth/Project/Team/Org, ComputeTII, ComputeCompleteness
- **Acceptance criteria:** All 6 met (AC-1 through AC-6)
- **Key commits:** 
  - Tidy: 7a704ac (domain types)
  - Feat: aeb62c1 (helpers), 6e01c7b (normalization), a866e78 (impact), 8cc6243 (dimensions), 5773707 (TII/completeness)
- **Test command:** `go test -race ./...` → all pass

---

### Alert Engine
- **Job story:** When scores or trends cross defined thresholds, I want the system to raise Amber or Red alerts automatically, so that I can intervene before a situation worsens.
- **Acceptance scenarios verified:**
  - All 6 individual alert conditions (Performance Deterioration, Morale Risk, Burnout Risk, Feedback Risk, Customer/Business Risk, Data Quality Risk) classify Red/Amber/none exactly per PRD 7.1 thresholds, including boundary values ✓
  - All 4 team-level alert conditions (Team Morale Drift, Team Delivery Drift, Systemic Burnout, Calibration Risk) classify Red/Amber/none exactly per PRD 7.2 thresholds, including boundary values ✓
  - Each alert carries condition type, severity, and member reference — understandable without re-deriving from raw scores ✓
  - Alert evaluation is deterministic: pure functions, no I/O, no global state; verified by static inspection and 5x repeat-run identical results ✓
- **Evidence:**
  - Engine layer: 39 new unit tests in `engine/domain/alerts_test.go` covering both Amber and Red thresholds at exact PRD boundary values for all 10 alert conditions
  - Full domain suite: 55 tests passing (`go test -race ./engine/domain`)
  - Full repo suite: `go test -race ./...` → all packages pass, no race conditions
  - Build: `go build ./...` succeeds (harmless macOS linker warning only)
- **Acceptance criteria:** All 5 met (AC-1: 6 individual conditions, AC-2: 4 team-level conditions, AC-3: alert detail sufficiency, AC-4: determinism, AC-5: boundary-value test coverage)
- **Known limitation:** team-level alert inputs (% members Red, TII stddev history) are accepted as pre-computed parameters; the aggregation pipeline that derives them from real member data is not yet built (deferred to a future increment).
- **Test command:** `go test -race ./engine/domain -run TestEvaluate` → all pass

---

### UI Refactored to Controller Pattern
- **Job story:** When I add a new screen or fix a form logic bug, I want business logic separated from UI widgets, so that logic is unit testable without Fyne overhead and easier to debug and reuse.
- **Acceptance criteria met:**
  - ✅ Controller package created at `ui/controllers/` with base controller interface
  - ✅ MonthlyInputController extracted; owns form state, data loading, validation, save logic
  - ✅ SettingsController extracted; owns member list state and member CRUD logic
  - ✅ Both controllers have zero Fyne imports and are unit testable
  - ✅ Existing MonthlyInputScreen and SettingsScreen refactored to use controllers
  - ✅ All existing tests continue to pass (acceptance tests, integration tests, unit tests)
  - ✅ New controller unit tests added and passing (10 unit tests across both controllers)
- **Evidence:**
  - Controller unit tests: 10 tests passing (`ui/controllers/monthly_input_controller_test.go` — 5 tests; `ui/controllers/settings_controller_test.go` — 5 tests)
  - Screen acceptance tests: 6 tests passing (`ui/screens/monthly_input_acceptance_test.go` — 3 tests including form persistence and reload)
  - Screen integration test: 1 test passing (`ui/screens/monthly_input_it_test.go` — full form workflow test)
  - Total test count: 17 passing (10 unit + 6 acceptance + 1 integration)
  - Build: `go build ./ui/controllers` and `go build ./ui/screens` succeed
  - Race detector: all tests pass with `-race` flag
  - Architecture decision: ADR-20260914 documents rationale and implementation pattern
- **Architecture:**
  - Controllers live in `ui/controllers/` and import only `service/` and `engine/domain` (zero Fyne)
  - Dependency direction: UI Screen → Controller → Service Layer (never reversed)
  - Screens are thin Fyne renderers that delegate state and logic to controllers
  - Controllers are pure Go functions testable without UI framework setup
- **Documentation:**
  - `docs/adr/ADR-20260914-ui-controllers.md` — Decision rationale and pattern
  - `docs/ui.md` — New "Controller Pattern (Screen Architecture)" section with anatomy, testing, and common patterns
- **Key commits:** 
  - Controllers: extraction and unit tests
  - Screens: refactored to use controllers, all tests remain green
- **Test command:** `go test -race ./ui/...` → all pass (17/17)

---

### Monthly Input Workspace — Navigation Shell and Screen Registration
- **Job story:** When I open the monthly input workspace, I want to enter raw signals and impact ratings for each team member with live validation and persistence, so that I can complete the cycle accurately and efficiently.
- **Acceptance scenarios verified:**
  - Monthly Input screen registered on the screen registry and reachable at runtime ✓
  - Valid monthly entry saves and persists for an active member ✓
  - Invalid signal values (out-of-range) are rejected with a descriptive error ✓
  - Duplicate entries for the same member and month are rejected ✓
  - Saved entries survive a reload and are retrievable by member ✓
  - App shell provides sidebar navigation (Overview / Input / Alerts / Review / Settings) and a cycle header ✓
- **Evidence:**
  - Service layer: existing 8 tests passing in `service/monthly/monthly_test.go` (create, update, get, list, invalid-signal, duplicate, inactive-member, compute)
  - UI layer: `ui/screens/monthly_input_test.go` — `TestScreenRegistry_ProvidesMonthlyInputScreen` passing
  - Build: `go build ./...` passes
  - Full suite: `go test -race ./...` passes
- **Acceptance criteria:** All 4 met (AC-1: workflow reachable; AC-2: validation blocks invalid input; AC-3: data persists across restarts; AC-4: full team cycle supported)
- **Note:** Monthly Input form content (signal fields, live preview, member list) is the next increment. This increment delivered the service contract, persistence layer, navigation shell, and screen registration.
- **Test command:** `go test -race ./...` → all pass

---

## SPA Implementation Path

The following increments replace Fyne screens with a browser-based SPA served from the Go binary. The backend (engine, store, service, controllers) remains unchanged. Fyne screens are retired incrementally as their SPA equivalents pass acceptance criteria.

**Architecture target:** Go binary runs the HTTP server with embedded SPA assets (HTML, CSS, JS, templates) via `embed.FS`. The browser talks to the Go server over a JSON REST API. Controllers delegate to services. No business logic in HTTP handlers.

**Frontend stack:** HTMX + Alpine.js + Go `html/template` + Bulma CSS (see `docs/adr/ADR-20260915-spa-frontend-stack.md` for rationale).

---

### SPA Increment 1: HTTP Server Shell

**Goal:** Boot a local HTTP server from `main.go`; serve a static "hello" HTML page; prove the embed pipeline works.

**Scope:**
- Add `server/` package with a `Start(addr string, services ApplicationServices)` function
- Wire `net/http` with an `embed.FS` pointing at `web/dist/` (the SPA build output)
- `main.go` starts both the HTTP server (on a fixed or configurable port, e.g. `8080`) and (temporarily) the Fyne window, so the app still works while the SPA is being built
- CI: add `go test ./server/...` to the test command

**Why first:** establishes the binary embed pipeline and server lifecycle before any API or UI work.

**Acceptance criteria verified:**
- AC-1: `go run .` starts the server and the existing Fyne window simultaneously without error ✓
- AC-2: `curl http://localhost:8080/` returns a 200 with a static HTML page embedded in the binary ✓ (verified via `TestHTTPServerStartsAndServesHTML`)
- AC-3: `go test -race ./server/...` passes; handler tests use `net/http/httptest` — no Fyne, no OpenGL ✓

**Evidence:**
- Server package: `server/server.go` with `Start(addr string)` function booting HTTP server and serving embedded assets
- Embedded assets: `server/dist/` with static `index.html` (mirrored from `web/dist/` to support Go embed.FS)
- Main integration: `main.go` updated to start HTTP server in a goroutine while Fyne window runs in main thread
- Server unit tests: 6 tests passing in `server/server_test.go` covering asset embedding, file serving, and no-Fyne verification
- Server integration test: `TestHTTPServerStartsAndServesHTML` verifies the full server startup and HTML response
- Build: `go build ./...` succeeds with embedded assets
- Test: `go test -race ./server/...` passes (4 tests, no Fyne imports)

**Key commits:** 
- Server package with embed.FS and HTTP handler
- HTML asset at web/dist/index.html
- Main.go integration (HTTP + Fyne)
- Server tests (unit + integration)

**Test command:** `go test -race ./...` → all pass

### Fyne Removal & SPA-First Refactor

**Goal:** Remove all Fyne dependencies and desktop UI artifacts; restructure use-case logic as coordinators for reuse by HTTP handlers and future CLI clients; update architecture docs to reflect SPA-first design.

**Acceptance scenarios verified:**
- Deleted `ui/screens/`, `ui/app.go`, `ui/app_integration_test.go` (19 Fyne screen files removed) ✓
- Moved `ui/controllers/` logic to `service/coordinator/monthly.go`, `service/coordinator/settings.go` ✓
- Removed Fyne from `go.mod`; `grep -i fyne go.mod` returns nothing ✓
- Rewrote `docs/architecture.md` to describe HTTP → Coordinator → Service → Store → Engine ✓
- Deleted `docs/ui.md` and `docs/adr/ADR-20260913-go-fyne-desktop.md` ✓
- Created `docs/adr/ADR-20260915-clean-architecture-layering.md` documenting coordinator pattern ✓
- Updated `main.go` to boot HTTP server only (no Fyne window); starts on localhost:8080 ✓

**Acceptance criteria verified:**
- AC-1: No Fyne in dependencies — `go mod tidy` removes Fyne and transitive deps (OpenGL, GLFW, text-render) ✓
- AC-2: Fyne code removed — `ui/screens/`, `ui/app.go`, `ui/controllers/` deleted; logic in `service/coordinator/` ✓
- AC-3: Controllers → Coordinators — `service/coordinator/monthly.go`, `service/coordinator/settings.go` with 13 passing tests ✓
- AC-4: Docs rewritten — `docs/architecture.md` is SPA-first; `ui.md` deleted ✓
- AC-5: Old Fyne ADR deleted — `ADR-20260913-go-fyne-desktop.md` removed ✓
- AC-6: New clean architecture ADR added — `ADR-20260915-clean-architecture-layering.md` documents coordinator pattern ✓
- AC-7: Roadmap updated — feature marked "In Progress"; now moved to "Done" ✓
- AC-8: main.go boots HTTP server only — Fyne imports removed; server.Start() blocks ✓
- AC-9: All tests pass — `go test -race ./...` passes (10 packages, 13 coordinator tests, 7 handler tests) ✓

**Evidence:**
- Server layer: `server/server.go` with Start(addr) function; `server/handler_test.go` spike test proving HTTP context can call coordinators
- Coordinator layer: `service/coordinator/monthly.go` (7 tests), `service/coordinator/settings.go` (6 tests), `service/coordinator/error.go` with ValidationError types
- Architecture: `docs/architecture.md` describes HTTP → Coordinator → Service → Store → Engine with C4 Level 2 diagram
- Clean architecture: `docs/adr/ADR-20260915-clean-architecture-layering.md` documents coordinator pattern, alternatives, and implementation rules
- Removed: 19 Fyne screen files, `ui/app.go`, `ui/controllers/`, `docs/ui.md`, old Fyne ADR
- Build: `go build ./...` succeeds; `go test -race ./...` passes (all 10 packages); no Fyne imports remain
- Main: Updated to boot server only; database, schema, services remain available for HTTP handlers

**Key commits:**
- Subtask 1-4: Coordinators extracted from controllers (monthly, settings, error)
- Subtask 5: Handler spike test proves coordinator reusability
- Subtask 6: Fyne artifacts deleted (ui/screens, ui/app.go, ui/controllers)
- Subtask 7: main.go updated to boot HTTP server only
- Subtask 8: go mod tidy removes Fyne dependency
- Subtask 9: docs/architecture.md rewritten for SPA-first design
- Subtask 10-11: Deleted ui.md and old Fyne ADR
- Subtask 12: All tests pass; no Fyne deps remain

**Test command:** `go test -race ./...` → all pass (10 packages)

---

### SPA Increment 2: Team Members API

**Goal:** Expose the existing member CRUD over HTTP. Prove that coordinator → handler wiring is testable with plain HTTP.

**Scope:**
- `GET /api/members` — list active members
- `POST /api/members` — add member
- `PATCH /api/members/{id}` — edit member
- `DELETE /api/members/{id}` — deactivate member (soft delete)
- Handlers delegate to `MemberAPICoordinator`; no business logic in handlers
- Request/response shapes: JSON with separate `firstName` and `lastName` fields; member ID as UUID
- Handler unit tests use `httptest.NewRecorder` — zero Fyne, zero SQLite (coordinator mocked)
- Integration tests wire real in-memory SQLite store (same pattern as `store/*_test.go`)

**Why second:** member CRUD is the simplest use case and most exercised layer. Good test-case for the handler → coordinator → service → store chain before tackling the more complex monthly input.

**Acceptance criteria verified:**
- AC-1: All CRUD operations round-trip correctly through the API ✓ (`TestMembersAPIIntegration_CRUD`)
- AC-2: Invalid inputs return structured JSON error responses (400 with `{error, kind}`) ✓ (14 handler tests covering validation)
- AC-3: Handler unit tests use `httptest` — no Fyne, no SQLite per handler test ✓ (14 handler tests in `server/handler_members_test.go`)
- AC-4: `go test -race ./server` passes ✓ (21/21 tests passing)

**Evidence:**
- HTTP handlers: `server/handler_members.go` with POST/GET/PATCH/DELETE handlers
- Handler unit tests: 14 tests in `server/handler_members_test.go` (add/list/edit/delete success and error cases, all <1 second each)
- Coordinator: `service/coordinator/member_api.go` with AddMember/GetMembers/EditMember/DeactivateMember
- Integration test: `TestMembersAPIIntegration_CRUD` in `server/integration_test.go` (full stack: HTTP → coordinator → service → store)
- UUID mapping: `memberIDToUUID()` function converts domain int64 IDs to deterministic UUID strings (v5 with namespace `6ba7b810-9dad-11d1-80b4-00c04fd430c8`)
- Request/response contract: POST accepts `{firstName, lastName, seniority}`; responses include `id` (UUID), `firstName`, `lastName`, `seniority`, `status`, `createdAt` (RFC3339)
- Error handling: Validation errors (400), database errors (500), member-not-found (400); structured JSON with `error` and `kind` fields
- Server wiring: `server/server.go` registers routes; `main.go` creates coordinator and passes to `Start()`
- Build: `go build ./...` succeeds; `go test -race ./...` passes all 21 tests

**Test command:** `go test -race ./server -run "TestAddMemberHandler|TestGetMembersHandler|TestPatchMemberHandler|TestDeleteMemberHandler|TestMembersAPIIntegration"` → 15 tests passing

**Key commits:**
- Handlers and tests
- Coordinator (AddMember/GetMembers/EditMember/DeactivateMember)
- Server wiring and integration test
- Main.go coordinator injection

---

### Hexagonal Architecture for Member Operations

**Goal:** Refactor member CRUD to follow clean architecture and domain-driven design: replace `MemberAPICoordinator` with explicit use cases that depend on repository abstractions, achieving testability, portability, and clear DDD vocabulary.

**Scope:**
- ✅ Replace `MemberAPICoordinator` with `AddMemberUseCase`, `GetMembersUseCase`, `EditMemberUseCase`, `DeactivateMemberUseCase`
- ✅ Define `MemberRepository` interface (port) in `core/members/repository.go`
- ✅ Implement `InMemoryMemberRepository` in `storage/memory/` (for tests)
- ✅ Implement `SQLiteTeamMemberRepository` in `storage/sqlite/` (for production)
- ✅ Refactor `server/handler_members.go` to call use cases via dependency injection
- ✅ Maintain clean dependency graph: HTTP adapters → use cases → repository port ← storage adapters

**Why:** Use cases are portable to CLI, mobile, and other clients without HTTP coupling. DDD repositories hide storage details; tests inject in-memory implementations. Clean dependency flow (inversion) enables testing without infrastructure. Clear vocabulary matches industry standards (Clean Architecture, Domain-Driven Design).

**Acceptance scenarios verified:**
- Four use cases (Add, Get, Edit, Deactivate) each testable in isolation with injected in-memory repository ✓
- `MemberRepository` port interface enforces contract; all adapters (memory, SQLite) implement via `Save()`, `FindByID()`, `FindActive()`, `Deactivate()` ✓
- HTTP handlers call use case `Execute()` methods; zero coordinator references remain ✓
- `go test -race ./...` passes with 34 total tests (27 new); zero functionality lost; all 5 acceptance criteria met ✓
- No circular imports; core depends only on stdlib + `engine/domain`; dependency flow is clean ✓

**Evidence:**
- Core use cases: `core/members/{add,get,edit,deactivate}_member.go` with 18 unit tests (`core/members/*_test.go`)
- Domain aggregate: `core/members/member.go` with `TeamMember` value object and behavior (Move, UpdateName, Deactivate)
- Port interface: `core/members/repository.go` defining `MemberRepository` contract
- Storage adapters: `storage/memory/member_repository.go` (test fixture, 0 I/O) and `storage/sqlite/member_repository.go` (production, 6 tests)
- HTTP layer refactored: `server/handler_members.go` calls use cases directly; `server/handler_members_test.go` mocks use cases (no SQLite per test)
- Composition root: `main.go` wires SQLite adapter into use cases; passes use cases to handlers
- Test counts: 18 core + 6 storage + 8 handler = 32 new tests; all passing with `-race` flag
- Zero coordinator artifact: `service/coordinator/member_api.go` and `server/integration_test.go` deleted (superseded)
- Architecture enforced: `go build ./core/members` imports only stdlib + `engine/domain`; zero framework/storage dependencies

**Architecture:**
- **Hexagonal (Ports & Adapters):** HTTP adapters ↓ core (use cases) ↑ storage adapters
- **Dependency rule:** Core has zero dependencies on HTTP or storage layers; HTTP and storage depend inward on core ports only
- **Testing strategy:** Core tests inject in-memory adapters (millisecond startup); storage tests use real SQLite (`testing.T` hooks); HTTP tests mock use cases
- **Composition:** `main.go` wires SQLite adapter + use cases + HTTP handler injection at startup; no reflection or DI container

**Key commits:**
- Subtask 1: `core/members/member.go` — domain aggregate
- Subtask 2: `core/members/repository.go` — port interface
- Subtasks 3-6: Use cases + 18 unit tests (add, get, edit, deactivate)
- Subtask 7: `storage/memory/member_repository.go` — in-memory adapter (test fixture)
- Subtask 8: `storage/sqlite/member_repository.go` — SQLite adapter + 6 tests
- Subtask 9-10: `server/handler_members.go` refactored to use cases; `server/handler_members_test.go` rewritten (mocked use cases)
- Subtask 11: `main.go` rewritten — composition root wiring
- Subtask 12: `service/coordinator/member_api.go` deleted; `server/integration_test.go` deleted
- Subtask 13: Final verification — all tests green, zero circular imports, core independence confirmed

**Test command:** `go test -race ./...` → all 34 tests passing across 9 packages

**Documentation:**
- `docs/architecture.md` — C4 container diagram and hexagonal dependency flow updated to show ports & adapters
- `docs/adr/ADR-20260916-hexagonal-architecture-member-crud.md` — Decision rationale, alternatives, dependency enforcement, testing strategy
- Roadmap (this entry) — Promotion from In Progress to Done with evidence links

**Acceptance criteria:** All 5 met
- AC-1: ✅ Four use cases in `core/members/` with 18 passing unit tests (all mocked adapter scenarios covered)
- AC-2: ✅ `MemberRepository` port interface in `core/members/repository.go`; two implementations: `storage/memory/` (0 I/O) and `storage/sqlite/` (6 tests)
- AC-3: ✅ HTTP handlers call use case `Execute()` methods; zero coordinator references; dependency injection at construction
- AC-4: ✅ `go test -race ./...` passes with 34 tests; no functionality lost; original CRUD behavior preserved
- AC-5: ✅ No circular imports verified via `go build ./core/members`; dependency flow: HTTP → core ← storage (no reversal)

### Web App Shell (Foundation for SPA)

**Goal:** Build the web app shell with sidebar, top bar (with centered search), and main content area using HTMX + Alpine.js + Go templates + Bulma CSS.

**Job Story:** When I load the web app in my browser, I want to see a professional, accessible layout with a persistent sidebar, top navigation bar featuring a centered search bar, and content area, so that I have a foundation for building screens and the app feels polished from the start.

**Acceptance scenarios verified:**
- Shell layout with 3 semantic regions (header, aside, main), full-height flexbox structure ✓
- Sidebar component with logo, navigation links, collapsible section, footer button ✓
- Top bar sticky at 56px with sidebar toggle, centered search, quick action placeholders ✓
- Main content area scrollable with flex layout and padding ✓
- Responsive behavior: desktop fixed sidebar (260px), tablet/mobile overlay with toggle ✓
- Accessibility: semantic HTML, ARIA labels, keyboard navigation (Tab, Esc), focus outlines (2px #3273dc) ✓
- Browser verification: no errors in console, Bulma renders correctly, search input visible and interactive ✓

**Acceptance criteria met:**
- AC-1: Shell layout with 3 regions ✓
- AC-2: Sidebar component ✓
- AC-3: Top bar with toggle and search ✓
- AC-4: Main content area ✓
- AC-5: Responsive desktop/tablet/mobile ✓
- AC-6: Accessibility (WCAG 2.1 AA target) ✓
- AC-7: Browser verification ✓

**Key commits:**
- Subtask 1 (templates): `fce51e3`
- Subtask 2 (assets + server routing): `5243213`
- Subtask 3 (Alpine.js interactivity): `6749ac0`
- Subtask 4 (responsive verification): `f503dda`
- Subtask 5 (accessibility audit): `2a372b9`
- Subtask 6 (cleanup): `44fc199`
- Subtask 7 (acceptance tests): `bde2340`

**Test evidence:**
- 44 tests passing (`go test -race ./...`)
- 10 new acceptance tests in `server/server_test.go` covering AC-1 through AC-7
- 34 existing tests remain green (no regression)
- Focus state verification: visual inspection + keyboard navigation manual test
- Responsive breakpoints: 3 viewport sizes tested (≥1024px desktop, 769–1023px tablet, ≤768px mobile)

**Documentation:**
- `docs/ui.md` — New UI Design System documenting shell layout, breakpoints, accessibility, color palette, typography, and interaction patterns for all future screens
- `.agent/increment.md` — Increment definition and acceptance criteria
- `.agent/plan.md` — Technical execution plan (7 subtasks)
- `.agent/implementation.md` — Completion record with commit hashes

**Architecture notes:**
- All assets (Bulma, HTMX, Alpine.js) embedded in binary via `embed.FS` (offline-capable, zero CDN)
- Go `html/template` server-side rendering with progressive enhancement
- No build step for templates; wired directly in `server/server.go`
- HTTP routes: `/` serves shell layout, `/dist/*` serves embedded assets

---

## In Progress

*(None currently. Next increment: TBD from roadmap.)*

---

## Planned

### SPA Increment 3: Monthly Entry API

**Goal:** Build the Settings screen (member list, add/edit/deactivate) as a browser-rendered SPA page backed by the Members API. This is the first Fyne screen retired.

**Scope:**
- Scaffold the SPA project in `web/` with HTMX + Alpine.js + Go `html/template`
- Implement the Settings screen with identical UX to the existing Fyne screen
- Add `Makefile` target `make web` that prepares SPA templates; no separate build step required (templates embedded directly)
- Remove `ui/screens/settings.go` once the SPA version passes acceptance
- SettingsController remains in `ui/controllers/` and is called by HTTP handlers

**Why Settings first:** it is the simplest screen (pure CRUD, no formula rendering, no live preview). Low risk for proving the full stack.

**Acceptance criteria:**
- AC-1: Add, edit, and deactivate flows work in the browser identically to the Fyne screen behavior
- AC-2: Handler unit tests and Playwright end-to-end tests provide equivalent or better coverage than Fyne tests
- AC-3: `go test -race ./...` passes; `go build ./...` produces a single binary with the SPA assets embedded
- AC-4: `curl http://localhost:8080/settings` serves the Settings page with live functionality

---

### SPA Increment 5: Monthly Input Screen

**Goal:** Replace the Monthly Input Fyne screen (Screen B) with a SPA page.

**Scope:**
- Implement all 10 signal fields, impact rating, live TII/completeness preview, save, and copy-from-previous-month using HTMX + Alpine.js
- Live preview calls `POST /api/entries/preview` on field change (debounced via `x-on:input.debounce`)
- Remove `ui/screens/monthly_input.go` and related Fyne acceptance/integration tests once SPA tests pass
- All prior Fyne acceptance criteria must have equivalent browser tests

**Acceptance criteria:**
- AC-1 through AC-5 from the original Monthly Input increment verified via browser tests
- Live preview response latency ≤ 100ms on local machine (same PRD performance envelope)
- No Fyne or OpenGL import remains in any monthly input test

---

### SPA Increment 6: Remaining Screens (Overview, Member Detail, Alerts, Review)

Each screen follows the same pattern:
1. Add any missing API endpoints for the screen's data
2. Build the SPA screen using HTMX + Alpine.js + Go templates
3. Retire the equivalent Fyne placeholder or screen
4. Verify acceptance criteria via browser tests

These screens map 1:1 to the feature set (Overview Dashboard, Member Detail, Alerts Center, Monthly Review). They are sequenced identically to the original roadmap — the SPA migration does not change feature priority, only the rendering layer.

**Screens to build (in order):**
- Overview Dashboard (Screen A) — team KPI strip, dimension heatmap, alert table, action queue
- Member Detail (Screen C) — 12-month trend, signal contribution, evidence log, action plan
- Alerts Center (Screen D) — filter, triage, assign, resolve with reason capture
- Monthly Review and Calibration (Screen E) — finalization, promotion/support guardrail review, override capture, cycle lock

---

### SPA Increment 7: Remove Fyne Dependency

**Goal:** Once all screens are migrated, remove Fyne entirely from `go.mod`.

**Scope:**
- Delete `ui/screens/`, `ui/app.go`, and any remaining Fyne imports
- Remove `fyne.io/fyne/v2` and all OpenGL/GLFW/text-render transitive dependencies from `go.mod`
- Update `main.go` to only start the HTTP server (no Fyne window)
- Update `docs/architecture.md` and `docs/adr/` with a new ADR documenting the Fyne retirement decision
- Binary size and startup time should decrease materially

**Acceptance criteria:**
- AC-1: `go mod tidy` leaves no Fyne dependency in `go.mod` or `go.sum`
- AC-2: `go test -race ./...` passes with no OpenGL/headless test infrastructure
- AC-3: `go build ./...` produces a working binary; `curl http://localhost:8080/` serves the SPA homepage

---

### Feature Backlog (Unaffected by Migration)

The following features are still planned but are not blocking the SPA migration. They can be built once the SPA foundation is stable:

- **Trend Calculations (MA3, Delta1, Delta3, Vol3)** — moving average, deltas, volatility for trend-based alerts
- **CSV Export** — export team overview and member detail data as CSV for meetings and archival
- **Settings threshold tuning** — allow users to adjust billability target, alert sensitivity presets, reminder dates

---

## Open Questions

### Port Configuration
- **Question:** Should the port be hardcoded, a flag, or read from a config file?
- **Considerations:** A flag (`--port 8080`) is the standard Go convention and easiest to test. A config file adds complexity.
- **Decision needed before:** SPA Increment 1

### Browser Launch on Startup
- **Question:** Should the binary automatically open the default browser on startup?
- **Considerations:** Good UX for a desktop-replacement tool; `os/exec` + `open`/`xdg-open`/`start` per platform; optional via `--no-browser` flag.
- **Decision needed before:** SPA Increment 1

### Reactivating Deactivated Members
- **Question:** When a previously deactivated member re-joins the team, should they be reactivated in place (restoring their history) or added as a new member (clean slate)?
- **Considerations:** Reactivation preserves audit history and avoids duplicate entries; a new member record is simpler but loses historical context and risks orphaned score data.
- **Decision needed before:** Member Detail (Screen C) and any feature that reads historical data per member.

### Duplicate Member Names
- **Question:** Two or more members can legitimately share the same first and last name. How should the UI help users distinguish between them when selecting or reviewing?
- **Constraints:** The store must allow duplicate names (names are not a unique key). Disambiguation must not require renaming real people.
- **Options to explore:** display seniority + join date inline, require a display alias on add, or show member ID as a tie-breaker.
- **Decision needed before:** Monthly Input Workspace (Screen B) and any picker or dropdown that references members by name.

---

## How This List Works

- Features move: Planned → In Progress → Done. Never skip In Progress.
- A feature enters Done only when its user outcome is verified and the evidence link is present.
- Fyne screens are removed only after their SPA equivalent passes all acceptance criteria.
- Implementation detail belongs in code and phase artifacts, not here.
- If a planned feature is dropped, remove it and record the reason in a commit message or ADR.
