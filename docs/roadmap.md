# Roadmap

Product direction and sequencing for Team Impact Scorecard. Each entry explains the user outcome, current confidence, and ordering rationale.

> A feature moves to **Done** only when its user outcome is verified and the evidence is linked here.
> Source of truth: if a feature is not in Done with a passing test link, it is not considered shipped.

---

## Migration Notice

**As of 2026-09-17, the frontend is being rebuilt as a dedicated Vue 3 SPA.** The Go backend (`engine/`, `core/`, `storage/`) is unaffected and becomes API-only: `server/` will expose JSON endpoints exclusively, with no HTML rendering. This supersedes the interim HTMX + Alpine.js + Go `html/template` approach (2026-09-15), which proved awkward for composing a persistent app shell with routed screen content (see `docs/adr/ADR-20260917-vue-spa-frontend.md` for the full rationale).

**What this means for prior "Done" work:**
- Backend/domain work (Team Member Management, Formula Engine, Alert Engine, Members API, Hexagonal Architecture) is unaffected — it is consumed identically by the new Vue frontend.
- The HTML-rendering work ("Web App Shell", "SPA Increment 3: Members Screen") is superseded — it shipped and worked, but is being rebuilt in Vue and the Go HTML template code will be retired. See the **Frontend Migration Path** section below for the rebuild sequence.

See `docs/adr/ADR-20260917-vue-spa-frontend.md` for the current frontend stack decision. See `docs/adr/ADR-20260915-spa-frontend-stack.md` (superseded) for the interim decision's history.

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

### SPA Increment 3: Members Screen (HTMX/Go-template) — Superseded

**Status:** Superseded by the Vue rebuild (see **Frontend Migration Path — Members Screen (Vue)** below). Kept here for evidence continuity; the HTML rendering code this entry describes will be deleted from `server/` as part of the Vue Members screen work.

**Goal:** Build the Members screen (member list, add, edit, deactivate) as a browser-rendered SPA page backed by the existing Members API. This replaced manual member management via the Fyne Settings screen.

**Why superseded:** shell/content composition via Go templates + HTMX proved structurally awkward (a bug shipped where the app shell failed to wrap screen content). See `docs/adr/ADR-20260917-vue-spa-frontend.md` for the full context and decision.

**Acceptance scenarios verified (at the time, against the HTMX/Go-template implementation):**
- ✓ Members list displays active members with Edit/Deactivate buttons
- ✓ Add flow: navigate `/members/add` → fill form → submit → redirected to list, member appears
- ✓ Edit flow: click Edit → pre-filled form → change seniority → save → redirected to list, updates persist
- ✓ Deactivate: click Deactivate → row fades and removes from list (no reload) → member stays deactivated on refresh
- ✓ Form validation: empty fields/invalid data show errors, form re-renders with values preserved
- ✓ Responsive: correct layout on all three breakpoints (desktop, tablet, mobile)
- ✓ Accessibility: semantic HTML, ARIA labels on buttons, focus outlines visible, keyboard nav works

**Evidence (historical):**
- HTTP handlers: `server/handler_members.go` with `HandlerGetMembersPage`, `HandlerGetAddMemberPage`, `HandlerPostAddMember`, `HandlerGetEditMemberPage`
- HTML templates: `server/templates/members/list.html`, `add.html`, `edit.html` with responsive CSS and AJAX deactivate
- All tests passing at the time: `go test -race ./...` → 10 packages, 0 failures

**Key commits:**
- 0e868f9 — tidy: update sidebar link from Settings to Members
- e8bac45 — tidy: scaffold members template directory and reusable form component
- f04d49d — feat: add HandlerGetMembersPage to render members list as HTML
- 74ab764 — feat: add Members screen HTML pages and handlers for add/edit/list

---

## In Progress

### Docker Services Architecture & Acceptance Testing Foundation

**Goal:** Restructure the project into a multi-service architecture with separate frontend and backend services in Docker containers, and establish an acceptance test suite that runs against the full containerized system.

**Job Story:** When I develop features locally or validate them in CI, I want to run the entire system (frontend, backend, database) as containerized services via docker-compose, and verify end-to-end behavior with acceptance tests that exercise both the UI and API, so that I can catch integration bugs early and ship with confidence.

**Branch:** `increment/docker-services-architecture`

**Acceptance criteria:**
1. **Folder structure** — Created `services/frontend/` and `services/backend/` subdirectories; Go backend and npm frontend each have their own space
2. **Backend Dockerfile** — Builds a Go binary and packages it as a container; exposes port 8080; includes health check
3. **Frontend Dockerfile** — Builds Vue 3 SPA via Vite and serves it via a lightweight web server (or embeds in Go); exposes port 3000 (for dev) or as built artifacts
4. **Docker-compose** — Orchestrates both services + PostgreSQL/SQLite database container; sets environment variables for service discovery; all services accessible via standard ports
5. **Acceptance tests** — Playwright tests run against the live docker-compose environment; verify app shell loads, health check works, member CRUD flows end-to-end
6. **CI integration** — Docker build/push steps documented; acceptance tests run in CI after docker-compose services are healthy

**Status:** ✅ Done

**Delivered Subtasks:**
- tidy: Move Go backend to `services/backend/` (75 files reorganized, tests pass)
- tidy: Create backend Dockerfile (multi-stage: test → build → runtime)
- tidy: Scaffold Vue 3 + Vite frontend (AppShell component, health indicator, Bulma styling)
- tidy: Create frontend Dockerfile (multi-stage: npm build → nginx)
- feat: Create `docker-compose.yml` (backend + frontend orchestration, volume persistence, health checks)
- tidy: Set up `acceptance-tests/` directory with Playwright scaffold
- feat: Write first acceptance test suite (`acceptance-tests/tests/app.spec.ts` with 3 test cases)
- tidy: Update Makefile with docker-build, docker-up, docker-down, docker-logs, test-backend, test-frontend, test-acceptance targets
- tidy: Update .gitignore for docker volumes, node_modules, dist, test-results
- tidy: Update docs/architecture.md with multi-service C4 diagram
- feat: Create docs/testing.md (testing pyramid, strategy, commands)
- feat: Create docs/deployment.md (docker-compose guide, CI/CD, health checks, scaling)

**Evidence:**
- `make docker-build` ✅ both images build successfully (backend:latest, frontend:latest)
- `make docker-up` ✅ services start; backend marked healthy; frontend responding on port 3000
- `make test-backend` ✅ all 10 Go packages pass with race detection
- `GET /api/health` ✅ returns JSON `{"status":"ok"}`
- `docker-compose.yml` ✅ validated; persistent data volume works across down/up cycles
- Acceptance test scaffold ✅ ready at `acceptance-tests/tests/app.spec.ts` (3 test cases: app-shell-renders, health-indicator-green, no-console-errors)
- All Makefile targets ✅ verified working (help, build, run-backend, run-frontend, test-backend, docker-build, docker-up, docker-down, docker-logs, clean, clean-docker)

**Acceptance test location:** `acceptance-tests/tests/app.spec.ts` — ready to run once services are healthy.

---

## Planned

### Frontend Bootstrap: Bulma Shell + Health Indicator

**Goal:** Bootstrap a Vue 3 + Vite SPA with a persistent Bulma shell, serve it from the frontend service, and display a green/red health indicator in the footer that checks backend connectivity on page load.

**Job Story:** When I load the web app in my browser, I want to see a working frontend with a Bulma shell layout and a small health icon in the footer that confirms the backend service is accessible, so that I have confidence the app is running and the backend is reachable.

**Branch:** `increment/spa-bootstrap-health`

**Scope:** This increment assumes the Docker services infrastructure (from **Docker Services Architecture & Acceptance Testing Foundation** above) is already in place. It adds the Vue SPA to `services/frontend/src/` with Bulma shell and health indicator component. The backend's `/api/health` endpoint is extended to return build version info (or kept simple).

**Acceptance criteria:**
1. Shell Layout — Vue app loads in browser with Bulma CSS applied; layout includes sidebar, top bar, main content area, and footer
2. Shell Persistence — Sidebar and top bar remain visible as the user navigates between client-side routes (no full-page reload)
3. Health Endpoint — Backend exposes `GET /api/health` returning `{status: "ok"}` with HTTP 200
4. Health Indicator — Footer displays a green checkmark icon (✓) on successful backend connection, red X (✗) if connection fails; check occurs once when the page loads
5. Docker Integration — Frontend service's npm build output is served by the frontend container; health checks work across container network

**Next:** Execute plan.md for Docker services, then this increment.

---

### Feature Backlog (Unaffected by Migration)

The following increments build a dedicated Vue 3 SPA against the existing (and, where noted, extended) JSON API. The backend (`engine/`, `core/`, `storage/`) is unaffected. `server/` is trimmed to API-only handlers as each screen's Go-template equivalent (if any) is retired. See `docs/adr/ADR-20260917-vue-spa-frontend.md` for the stack decision and rationale, and `docs/testing.md` for the testing boundary between Vitest/Vue Test Utils and Playwright.

**Sequencing principle:** each increment ships one screen, backed by its API (existing or newly built), with component tests for all behavior and exactly one Playwright test for the screen's main success flow.

---

#### Frontend Increment 1: Vue Project Scaffold and API-only Backend Cleanup

**Goal:** Stand up the Vue 3 + Vite project, wire the build into the Go embed pipeline, and strip `server/` of HTML rendering so it is API-only.

**Scope:**
- Scaffold a Vue 3 + Vite project (own `package.json`, Vitest, Vue Test Utils, Playwright configured but empty)
- Add Pinia; add Vue Router with a persistent `AppShell.vue` layout component reproducing the shell decisions in `docs/ui.md` (sidebar, top bar, responsive breakpoints)
- Wire `npm run build` output into the Go binary's `embed.FS`; update `main.go`/`server/server.go` to serve the built SPA for non-API routes
- Delete `server/templates/*` and the HTML-rendering handlers (`HandlerGetMembersPage`, `HandlerGetAddMemberPage`, `HandlerPostAddMember`, `HandlerGetEditMemberPage`) — their JSON-API equivalents already exist and are retained
- Update `docs/deployment.md`-documented build order in CI (`npm run build` → `go build`)

**Why first:** every subsequent screen increment depends on the shell, router, and build pipeline existing. This also resolves the shell-composition bug that triggered the Vue migration.

**Acceptance criteria:**
- AC-1: `npm run build` followed by `go build` produces a working binary; loading `http://localhost:8080/` in a browser shows the Vue shell
- AC-2: Navigating to any client-side route keeps the shell (sidebar, top bar) visible — no full-page reload, no missing-shell regressions
- AC-3: `server/templates/*` is deleted; no `html/template` import remains in `server/`
- AC-4: `go test -race ./...` passes with only API-contract handler tests remaining in `server/`
- AC-5: `npm test` (Vitest) passes for the shell component

---

#### Frontend Increment 2: Members Screen (Vue)

**Goal:** Rebuild the Members screen (list with Active/Deactivated/All tabs, add, edit, deactivate, reactivate) as Vue components against the existing Members API.

**Scope:**
- `MemberList.vue` — table with status tabs, calling `GET /api/members?status=...`
- `MemberForm.vue` — shared add/edit form component
- Pinia store for member state (list, filters, loading/error states)
- Wire deactivate (`DELETE /api/members/{id}`) and reactivate (`PATCH /api/members/{id}/reactivate`) actions
- Delete the superseded Go-template Members screen artifacts (see **SPA Increment 3 — Superseded** above) once the Vue screen passes acceptance criteria

**Why second:** the API already exists in full (status filtering, reactivation) from prior increments; this is the lowest-risk screen to prove the Vue pattern end-to-end.

**Acceptance criteria:**
- AC-1: Members list displays Active/Deactivated/All tabs backed by the existing status-filtered API; the shell remains visible at all times
- AC-2: Add member form validates required fields and persists a new member
- AC-3: Edit member form loads existing data, persists changes
- AC-4: Deactivate and reactivate actions work without a full page reload and reflect immediately in the list
- AC-5: Vitest component tests cover list rendering, tab switching, form validation, and error states
- AC-6: One Playwright test covers the main success flow: add a member, see it appear in the Active tab

---

#### Frontend Increment 3: Monthly Entry API

**Goal:** Build the backend API for monthly signal entry (currently only the service/coordinator layer exists) so the Monthly Input screen has a contract to build against.

**Scope:**
- `POST /api/entries`, `GET /api/entries/{id}`, `PATCH /api/entries/{id}`, `GET /api/entries?member={id}` following the same core/use-case/repository pattern as members
- `POST /api/entries/preview` — computes TII/completeness/alerts for in-progress (unsaved) input, without persisting
- Handler tests only (JSON contract) — no UI work in this increment

**Acceptance criteria:**
- AC-1: All CRUD endpoints round-trip correctly (core use-case tests + handler tests)
- AC-2: Preview endpoint returns computed scores without persisting any data
- AC-3: Invalid signal values (out of range) return structured 400 errors
- AC-4: `go test -race ./...` passes

---

#### Frontend Increment 4: Monthly Input Screen (Vue)

**Goal:** Build the Monthly Input Workspace (Screen B) in Vue against the API from Increment 3.

**Scope:**
- Member list, raw signal form, IG/IP/IT/IO rating matrix, live score/alert preview panel, validation and completeness warnings
- Live preview calls `POST /api/entries/preview` on field change (debounced client-side)
- Save Draft / Submit Member Month / Copy Previous Month actions

**Acceptance criteria:**
- AC-1: All 10 signal fields and impact ratings can be entered and saved
- AC-2: Live preview updates within the PRD's 100ms performance envelope
- AC-3: Validation and completeness warnings display inline, matching PRD 4.3
- AC-4: Vitest component tests cover the form, validation, and preview panel; one Playwright test covers "enter and submit a complete monthly entry"

---

#### Frontend Increment 5: Remaining Screens (Overview, Member Detail, Alerts, Review, Settings)

Each screen follows the same pattern:
1. Build or extend the JSON API endpoints the screen needs (own increment if the API doesn't exist yet)
2. Build the Vue screen (components + Pinia store + route)
3. Cover behavior and edge cases with Vitest/Vue Test Utils; add one Playwright test for the screen's main success flow
4. Mark the screen Done in this roadmap with evidence

**Screens to build (in PRD order):**
- Overview Dashboard (Screen A) — KPI strip, team TII trend chart, DG/DP/DT/DO heatmap, alert table, action queue
- Member Detail (Screen C) — current TII/dimensions/confidence/alerts, 6–12 month trend timeline, signal contribution table, evidence log, action plan tracker
- Alerts & Risk Center (Screen D) — filterable alert list, root cause summary, suggested interventions, assign/resolve/snooze
- Monthly Review & Calibration (Screen E) — distribution charts, member decision-support cards, override reason capture, finalize/lock cycle
- Limited Settings (Screen F) — billability target/tolerance, alert sensitivity preset, reminder dates

Each screen is its own increment with its own acceptance criteria, defined when that increment starts (per `docs/prd.md` §10 for component scope and CTAs).

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
