# Architecture

Team Impact Scorecard — structural overview (SPA-first, API-driven).

---

## C4 Level 2 — Container View

```
┌─────────────────────────────────────────────────────────────┐
│  User's Machine / Browser                                   │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Team Impact Scorecard (Single Binary + Browser SPA) │  │
│  │                                                      │  │
│  │  ┌──────────────────┐   HTTP calls  ┌────────────┐  │  │
│  │  │                  │ ────────────► │            │  │  │
│  │  │   Browser SPA    │               │  server/   │  │  │
│  │  │                  │               │            │  │  │
│  │  │  HTMX +          │               │  HTTP      │  │  │
│  │  │  Alpine.js +     │               │  handlers  │  │  │
│  │  │  Go templates    │   (JSON API)  │            │  │  │
│  │  │                  │               │            │  │  │
│  │  │  Embedded in     │               └──────┬─────┘  │  │
│  │  │  Go binary       │                      │ calls   │  │
│  │  │  (embed.FS)      │                      ▼         │  │
│  │  │                  │   calls       ┌─────────────┐  │  │
│  │  │                  │ ────────────► │             │  │  │
│  │  │                  │               │ coordinator/│  │  │
│  │  │                  │               │             │  │  │
│  │  └──────────────────┘               │ Orchestrate │  │  │
│  │                                     │ service     │  │  │
│  │   (SPA replaceable with               │ calls      │  │
│  │    CLI, future UIs)                   └──────┬────┘  │  │
│  │                                              │ calls   │  │
│  │                                              ▼         │  │
│  │                                    ┌──────────────┐   │  │
│  │                                    │              │   │  │
│  │                      HTTP handlers │  service/    │   │  │
│  │                      delegate to    │              │   │  │
│  │                      coordinators   │ Use cases:   │   │  │
│  │                                    │ AddMember,   │   │  │
│  │                                    │ SubmitMonth, │   │  │
│  │                                    │ GenerateAlerts
│  │                                    │              │   │  │
│  │                                    └──────┬───────┘   │  │
│  │                                           │ calls      │  │
│  │                                           ▼            │  │
│  │                                    ┌──────────────┐   │  │
│  │                                    │              │   │  │
│  │                                    │   store/     │   │  │
│  │                                    │              │   │  │
│  │                                    │  SQLite      │   │  │
│  │                                    │  read/write. │   │  │
│  │                                    │              │   │  │
│  │                                    └──────┬───────┘   │  │
│  │                                           │ reads/    │  │
│  │                                           │ writes    │  │
│  │                                           ▼           │  │
│  │                                    ┌──────────────┐   │  │
│  │                                    │              │   │  │
│  │                                    │   engine/    │   │  │
│  │                                    │              │   │  │
│  │                                    │  Pure funcs. │   │  │
│  │                                    │  Scoring,    │   │  │
│  │                                    │  norms, etc. │   │  │
│  │                                    │              │   │  │
│  │                                    └──────────────┘   │  │
│  │                                                      │  │
│  └──────────────────────────────────────────────────────┘  │
│                         │                                   │
│                         │ reads/writes                      │
│                         ▼                                   │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  SQLite File                                         │  │
│  │  $UserConfigDir/leadpulse/data.db                    │  │
│  │  — members, monthly entries, evidence, actions,      │  │
│  │    audit log, cycle state                            │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘

No external network connections. No cloud. No server. All data on-device.
```

---

## Containers

### `server/` — HTTP Server & API Layer
- **Technology:** Go `net/http`, Go `html/template`, embedded static assets
- **Responsibility:** Boot an HTTP server on localhost:8080; serve SPA static assets (HTML, CSS, JS); expose JSON REST API endpoints; handle request/response mapping.
- **Handlers:** Receive HTTP requests, parse JSON/form data, call coordinators, map results to HTTP responses. No business logic in handlers; all logic delegated to coordinators.
- **Constraints:** No direct store or engine access; no database connections. All persistence through coordinators and services.

### HTTP Handlers — Request/Response Adapter Pattern

Each HTTP handler is a thin adapter:

1. **Parse:** Extract JSON request body, path parameters, query strings into plain Go values
2. **Validate:** Check request format (no business validation; business validation is in coordinator)
3. **Delegate:** Call coordinator method with application-level inputs (not HTTP types)
4. **Map errors:** Convert coordinator errors to HTTP status codes and error JSON
5. **Serialize:** Convert result to JSON and write HTTP response

**Example:** `HandlerAddMember` in `server/handler_members.go`:
```go
func HandlerAddMember(w http.ResponseWriter, r *http.Request, coord *MemberAPICoordinator) {
  // Parse JSON request
  var req AddMemberRequest
  json.NewDecoder(r.Body).Decode(&req)
  
  // Delegate to coordinator (business logic)
  member, err := coord.AddMember(req.FirstName, req.LastName, req.Seniority)
  
  // Map error to HTTP response
  if err != nil {
    // error is of type coordinator.ValidationError
    w.WriteHeader(http.StatusBadRequest) // or 500 for database errors
    json.NewEncoder(w).Encode(map[string]interface{}{
      "error": err.UserMessage(),
      "kind":  err.Kind,
    })
    return
  }
  
  // Serialize success response
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(map[string]interface{}{
    "id":        memberIDToUUID(member.ID()),
    "firstName": member.Name().First,
    "lastName":  member.Name().Last,
    "seniority": member.Seniority(),
    "status":    "active",
    "createdAt": member.CreatedAt().Format(time.RFC3339),
  })
}
```

**Characteristics:**
- No conditional business logic (e.g., "if member is senior, discount the score")
- No loops or loops over domain objects (that's the coordinator's job)
- Typically 15–20 lines of code (parse, validate format, delegate, map error, serialize)
- All coordination, validation, and computation in coordinator
- Testable without HTTP: test the coordinator directly; test the handler with `httptest` and a mocked coordinator

**Testing:** Use `net/http/httptest.NewRecorder()` to capture responses. Mock the coordinator with an in-memory mock struct. Verify HTTP status codes, JSON structure, and error messages.

### `service/coordinator/` — Orchestration Layer
- **Technology:** Go (pure functions, no I/O)
- **Responsibility:** Orchestrate multi-step workflows. Accept application-level inputs (member IDs, month strings, signal maps — not HTTP types). Validate, sequence service calls, handle cross-service dependencies, return domain error types.
- **Coordinators:** `MonthlyInputCoordinator` (form state, data entry flow), `SettingsCoordinator` (member CRUD). Each coordinator is testable without HTTP or UI framework overhead.
- **Dependency model:** Coordinators call services; they do not call store or engine directly. This enables reuse by HTTP handlers, CLI commands, and future UIs.
- **Key constraint:** Stateless per HTTP request (or CLI invocation). No member variables that persist across calls except service references.

### `service/` — Application Use Case Layer
- **Technology:** Go (pure functions, no I/O)
- **Responsibility:** Application use cases, organized by domain (e.g., `service/member/`, `service/monthly/`). Each use case is a transaction: fetch aggregates via repository, call domain services for cross-aggregate rules, call `engine/` for computation, persist via repository, return the result or error.
- **Dependency model:** All persistence is abstracted via repository interfaces defined in `engine/domain/`. Service layer receives these interfaces as dependencies (not `*sql.DB`). This enables testing with in-memory repositories (no database required).
- **Domain services:** Application services depend on domain services (`ScoringService`, `TrendService`, `ValidationService`, `MonthlyEntryService`) which enforce cross-aggregate rules before persistence. Domain services also accept repository interfaces only; no database coupling.
- **Service aggregation:** Application services are injected into coordinators via constructor.
- **Constraints:** No UI logic. No direct SQLite access. No direct database dependencies. Delegates persistence to repository interfaces, computation to `engine/`, cross-aggregate rule enforcement to domain services.
- **Implemented use cases:** `service/member/` (AddMember, ListMembers, EditMember, DeactivateMember), `service/monthly/` (CreateEntry, GetEntry, ListEntriesByMember, UpdateEntry, ComputeScores, GetTrends).

### `store/` — Persistence Layer
- **Technology:** Go, `modernc.org/sqlite` (pure Go, no CGO)
- **Responsibility:** Implements repository interfaces defined in `engine/domain/`. All SQLite reads and writes. Schema definition and migrations. Reconstructs aggregates from database rows; persists aggregates to database tables.
- **Pattern:** Repository implementations (SQLiteTeamMemberRepository, SQLiteMonthlyEntryRepository) conform to repository interfaces defined in the domain layer. This keeps storage knowledge in the store layer while the domain layer remains storage-agnostic.
- **Data retained:** 24 months of monthly entries per member, evidence notes, action plans, audit trail, cycle finalization state.
- **Constraints:** No formula computation. No persistence of raw domain objects; only aggregates via repository interface. Repositories validate aggregate state before persisting.

### `engine/` — Formula Engine
- **Technology:** Go (pure functions, no external dependencies)
- **Responsibility:** All computation defined in the PRD. Normalization functions, dimension scores (DG/DP/DT/DO), Total Impact Index (TII), trend calculations (MA3/Delta1/Delta3/Vol3), alert evaluation, confidence score, decision guardrail checks.
- **Constraints:** No I/O of any kind. No imports from `store/`, `service/`, `server/`, or `ui/`. Input and output are plain Go structs.

### `SQLite File` — Local Data Store
- **Location:** `$UserConfigDir/leadpulse/data.db`
- **Technology:** SQLite 3
- **Responsibility:** Durable storage of all application data on the user's device.
- **Constraints:** Never accessed directly by HTTP handlers or coordinators. All access is through `store/`.

---

## Dependency Direction

```
HTTP Handler  →  Coordinator  →  Service  →  engine/domain (types only)
                                     ↓
                                repository interfaces
                                     ↓
                                  store  
                                     ↓
                           engine/domain (types only)
```

- `engine/domain` has no project-internal dependencies. Contains value objects, aggregates, domain services, and repository interfaces.
- `store` implements repository interfaces defined in `engine/domain`. May use engine types but must not call engine computation functions.
- `service` depends on repository interfaces (not concrete store implementations). Each use case is a transaction: read from repository, call domain services for cross-aggregate rules, call `engine/` if needed, write to repository, return result.
- `coordinator` depends on service instances. Orchestrates multi-step workflows; accepts application-level inputs; returns domain error types. No HTTP or UI framework knowledge.
- `server` (HTTP handlers) depends on coordinator instances. Maps HTTP request/response to coordinator calls. No business logic in handlers.
- No circular imports. No layer may import from a layer above it.

---

## Key Communication Paths

| From | To | What |
|---|---|---|
| Browser (SPA) | HTTP handler | HTTP request (GET /api/members, POST /api/entries, etc.) |
| HTTP handler | Coordinator | Use case call (e.g., `coordinator.AddMember(firstName, lastName, seniority)`) with application-level inputs |
| Coordinator | Service | Orchestrated use case call (e.g., `service.CreateEntry(...)`) |
| Service | Store | Read/write domain entities via repository interface |
| Service | Engine | Compute scores, alerts, trends, guardrails |
| Store (read) | SQLite file | SQL queries returning domain types |
| Store (write) | SQLite file | SQL inserts/updates, append-only audit log |
| Browser (SPA) | Static assets | HTML, CSS, JS served from embed.FS |

---

## Data Schema (Summary)

Tables:
- `members` — id, name, team_id, created_at
- `monthly_entries` — member_id, month (YYYY-MM), raw inputs (10 fields), impact ratings (40 fields), computed scores (DG/DP/DT/DO/TII/MA3/Delta1/Delta3/Confidence), status, finalized_at
- `evidence_notes` — id, member_id, month, author, category, body, created_at
- `action_plans` — id, member_id, month, description, owner, due_date, status, resolved_at
- `audit_log` — id, member_id, month, field, old_value, new_value, changed_by, changed_at
- `cycle_state` — team_id, month, phase, finalized_at, locked

Schema migrations are additive-only in v1 (see `docs/deployment.md`).

---

## Testing Strategy by Layer

The coordinator layer is the testability boundary for business logic. Tests isolate each layer and verify integration points.

**`engine/` — unit tests (pure functions)**
- No mocks needed. Direct function calls with inputs and assertions on outputs.
- 100% of formulas, alert thresholds, trend calculations, guardrail logic.
- Example: `TestNormalizeMorale_midRange`, `TestAlertBurnout_redThreshold`.
- Gate: all engine tests must pass before any service or coordinator code runs.

**`service/` — unit tests (mock store, real engine)**
- Mock the `store` interface. Real `engine` functions (no mocks).
- Each use case is a test: call the use case with inputs, verify it calls `store` methods in the right order with the right data, verify the result.
- Example: `TestAddMember_createsAndReturns`, `TestSubmitMonthlyEntry_computesAndPersists`.
- Gate: all service tests must pass before coordinator code is written. Services are the API contract.

**`service/coordinator/` — unit tests (in-memory repos, real services)**
- Use in-memory repositories (same pattern as service tests).
- Test coordinator workflows: Load, SelectMember, SaveEntry, AddMember, etc.
- Each coordinator is a test: instantiate with services, call methods, verify service calls are correct.
- Example: `TestMonthlyInputCoordinator_SaveMember_Persists`, `TestSettingsCoordinator_AddMember`.
- Gate: all coordinator tests must pass before HTTP handlers are written. Coordinators are the UI-agnostic contract.

**`server/` — unit and integration tests (in-memory repos + httptest)**
- Handler tests use `httptest.NewRecorder` to test request/response mapping without a live server.
- Handlers are wired to mocked or in-memory coordinators to verify correct delegation.
- Example: `TestHandlerAddMember_createsAndReturnsJSON`, `TestHandlerSaveEntry_validatesAndPersists`.
- Gate: all handler tests must pass. Handlers are the API surface.

**Browser (SPA) — end-to-end tests (live server, real database)**
- Playwright or Cypress tests exercise full workflows: login, load member list, submit entry, view results.
- Tests use a real SQLite database (in-memory or temporary file) and a live HTTP server instance.
- Example: `spec/monthly-input.spec.ts` (submit entry, verify TII updates, save, reload, verify persistence).
- Gate: acceptance criteria verified in browser tests. These are the user-facing correctness gate.

---

## ADR Index

Architectural decisions for this project are recorded in `docs/adr/`. Key decisions:

- [ADR-20260915-clean-architecture-layering.md](adr/ADR-20260915-clean-architecture-layering.md) — Clean architecture with coordinators as the orchestration layer
- [ADR-20260915-spa-frontend-stack.md](adr/ADR-20260915-spa-frontend-stack.md) — HTMX + Alpine.js + Go templates for SPA frontend
- [ADR-20260913-sqlite-local-storage.md](adr/ADR-20260913-sqlite-local-storage.md) — SQLite for local-only persistence
- [ADR-20260913-pure-go-sqlite-driver.md](adr/ADR-20260913-pure-go-sqlite-driver.md) — `modernc.org/sqlite` to avoid CGO

---

**Last updated:** 2026-09-15 — Refactored for SPA-first, HTTP-driven architecture. Removed Fyne desktop UI. Introduced coordinator orchestration layer.
