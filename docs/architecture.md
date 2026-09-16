# Architecture

Team Impact Scorecard — structural overview (SPA-first, API-driven, hexagonal architecture).

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
│  │  │  Alpine.js +     │               │  adapters  │  │  │
│  │  │  Go templates    │   (JSON API)  │            │  │  │
│  │  │                  │               │            │  │  │
│  │  │  Embedded in     │               └──────┬─────┘  │  │
│  │  │  Go binary       │                      │ injects │  │
│  │  │  (embed.FS)      │                      ▼         │  │
│  │  │                  │   calls       ┌─────────────┐  │  │
│  │  │                  │ ────────────► │             │  │  │
│  │  │                  │               │  core/      │  │  │
│  │  │                  │               │  members/   │  │  │
│  │  └──────────────────┘               │             │  │  │
│  │                                     │  Use Cases: │  │  │
│  │   (SPA replaceable with            │  - Add      │  │  │
│  │    CLI, future UIs)                │  - Get      │  │  │
│  │                                     │  - Edit     │  │  │
│  │                                     │  - Deact    │  │  │
│  │                                     └──────┬─────┘  │  │
│  │                                            │ depends │  │
│  │                                            │ on port │  │
│  │                                            ▼         │  │
│  │                                    ┌──────────────┐  │  │
│  │                                    │              │  │  │
│  │                                    │  core/       │  │  │
│  │                                    │  members/    │  │  │
│  │                                    │              │  │  │
│  │                                    │  PORT:       │  │  │
│  │                                    │  Repository  │  │  │
│  │                                    │  (interface) │  │  │
│  │                                    │              │  │  │
│  │                                    └──────┬───────┘  │  │
│  │                                           │ impl     │  │
│  │                    ┌──────────────────────┴──────────┐ │  │
│  │                    ▼ (in-memory for tests)          ▼ │  │
│  │            ┌──────────────────┐         ┌─────────────┐ │  │
│  │            │ storage/memory/  │         │ storage/    │ │  │
│  │            │ InMemory         │         │ sqlite/     │ │  │
│  │            │ Repository       │         │ SQLite      │ │  │
│  │            │ (test fixture)   │         │ Repository  │ │  │
│  │            └──────────────────┘         │ (production)│ │  │
│  │                                         │             │ │  │
│  │                                         └──────┬──────┘ │  │
│  │                                                │        │  │
│  │                                    ┌──────────▼──────┐ │  │
│  │                                    │                 │ │  │
│  │                                    │  engine/        │ │  │
│  │                                    │                 │ │  │
│  │                                    │  Pure funcs.    │ │  │
│  │                                    │  Scoring,       │ │  │
│  │                                    │  norms, etc.    │ │  │
│  │                                    │                 │ │  │
│  │                                    └─────────────────┘ │  │
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

### `server/` — HTTP Adapter Layer
- **Technology:** Go `net/http`, Go `html/template`, embedded static assets
- **Responsibility:** Boot an HTTP server on localhost:8080; serve SPA static assets (HTML, CSS, JS); expose JSON REST API endpoints; handle request/response mapping.
- **Handlers:** Receive HTTP requests, parse JSON/form data, call use cases via dependency injection, map results to HTTP responses. No business logic; all logic in core use cases.
- **Constraints:** No direct store or engine access; no database connections. All persistence through injected use case dependencies.

### HTTP Handlers — Request/Response Adapter Pattern

Each HTTP handler is a thin adapter:

1. **Parse:** Extract JSON request body, path parameters, query strings into plain Go values
2. **Validate:** Check request format (no business validation; business validation in use cases)
3. **Delegate:** Call use case method with application-level inputs (not HTTP types)
4. **Map errors:** Convert use case errors to HTTP status codes and error JSON
5. **Serialize:** Convert result to JSON and write HTTP response

**Example:** `HandlerAddMember` in `server/handler_members.go`:
```go
func HandlerAddMember(w http.ResponseWriter, r *http.Request, addUC AddMemberUC) {
  // Parse JSON request
  var req AddMemberRequest
  json.NewDecoder(r.Body).Decode(&req)
  
  // Delegate to use case (business logic)
  output, err := addUC.Execute(AddMemberInput{
    FirstName: req.FirstName,
    LastName:  req.LastName,
    Seniority: req.Seniority,
  })
  
  // Map error to HTTP response
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]interface{}{
      "error": err.Error(),
      "kind":  "validation_error",
    })
    return
  }
  
  // Serialize success response
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(MemberResponse{
    ID:        memberIDToUUID(output.ID),
    FirstName: output.FirstName,
    LastName:  output.LastName,
    Seniority: output.Seniority,
    Status:    output.Status,
    CreatedAt: output.CreatedAt.Format(time.RFC3339),
  })
}
```

**Characteristics:**
- No conditional business logic
- No loops or loops over domain objects
- Typically 15–20 lines of code (parse, validate format, delegate, map error, serialize)
- All validation and computation in use cases
- Testable without HTTP: test use cases directly; test handler with `httptest` and mocked use cases

**Testing:** Use `net/http/httptest.NewRecorder()` to capture responses. Mock use cases with test doubles. Verify HTTP status codes, JSON structure, and error messages.

### `core/members/` — Application Layer (Business Logic)
- **Technology:** Go (pure functions, no I/O, no framework dependencies)
- **Responsibility:** Implement application use cases (AddMember, GetMembers, EditMember, DeactivateMember); validate business rules; orchestrate persistence via injected port interface.
- **Key invariant:** Zero dependencies on storage, HTTP, or UI frameworks. Depends only on port abstractions and domain value objects.
- **Use Cases:** Each implements `Execute(input) (*output, error)` pattern; injected repository port at construction.

### `core/members/repository.go` — Port (Interface)
- **Technology:** Go interface (abstract contract)
- **Responsibility:** Define the contract for member persistence without implementation details.
- **Methods:** `Save()`, `FindByID()`, `FindActive()`, `Deactivate()`
- **Key rule:** Interface lives in core, implementations (adapters) live in `storage/`; adapters depend inward, never vice versa.

### `storage/` — Adapter Layer
- **Responsibility:** Implement the `MemberRepository` port for specific storage technologies.
- **Implementation:** Two adapters provided:
  - `storage/memory/`: In-memory map-based repository (test fixture, zero I/O)
  - `storage/sqlite/`: SQLite repository (production storage, ACID compliance)
- **Key invariant:** Each adapter implements the port interface; adapters depend on core, never vice versa.

### `engine/` — Pure Computation Layer
- **Technology:** Go (pure functions, no I/O)
- **Responsibility:** Scoring algorithms, norms, analysis logic. Shared value objects (Seniority, FullName, etc.)
- **Note:** Unchanged by hexagonal refactor; remains a pure computation layer used by core and services.

---

## Dependency Architecture

### Hexagonal (Ports & Adapters)

```
┌─────────────────────────────────────────────────────┐
│                    HTTP (Adapter)                   │
│              server/handler_members.go              │
│                   ↓ (calls)                         │
├─────────────────────────────────────────────────────┤
│              CORE (Business Logic)                  │
│            core/members/use_cases                   │
│          • Aggregate (TeamMember)                   │
│          • Port interface (MemberRepository)        │
│          • Use cases (Add, Get, Edit, Deactivate)  │
│          • Unit tests (inject in-memory adapter)   │
│                   ↓ (depends on)                    │
│      ┌──────────────────────────────────┐          │
│      │ Port: MemberRepository           │          │
│      │ (interface, zero impl logic)     │          │
│      └──────────────────────────────────┘          │
├─────────────────────────────────────────────────────┤
│        ADAPTERS (Implementations)                   │
│  ┌──────────────────────────────────────────────┐  │
│  │ storage/memory:      In-memory adapter       │  │
│  │ storage/sqlite:      SQLite adapter          │  │
│  │ ↑ (implements port)  ↑ (implements port)     │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘

Rules:
• core/ has ZERO imports from storage/, server/, engine/, or service/
• storage/ imports only from core/ (inbound only)
• server/ imports only from core/ (inbound only)
• No circular imports
```

### Wiring (Composition Root)

```go
// main.go: dependency injection at startup

// 1. Create storage adapter (SQLite)
repo := sqlite.NewSQLiteTeamMemberRepository(db)

// 2. Create use cases with injected port
addMemberUC := members.NewAddMemberUseCase(repo)
getMembersUC := members.NewGetMembersUseCase(repo)
editMemberUC := members.NewEditMemberUseCase(repo)
deactivateMemberUC := members.NewDeactivateMemberUseCase(repo)

// 3. Wire HTTP adapters with use cases
mux.HandleFunc("POST /api/members", 
  func(w http.ResponseWriter, r *http.Request) {
    server.HandlerAddMember(w, r, addMemberUC)
  })

// 4. Start server
http.ListenAndServe(":8080", mux)
```

---

## Communication Paths

| From | To | Protocol | Notes |
|------|----|----------|-------|
| Browser SPA | server/ | HTTP (localhost:8080) | JSON API |
| server/ (handlers) | core/members/ (use cases) | Function calls (interfaces) | Dependency injection |
| core/members/ (use cases) | MemberRepository port | Function calls (interface) | Adapter injected at composition root |
| MemberRepository port | storage/{memory,sqlite}/ | Adapter implementations | SQLiteRepository for production; InMemoryRepository for tests |
| storage/sqlite/ | SQLite database | SQL queries | ACID-compliant persistence |
| core/ (any) | engine/ (value objects) | Function calls | FullName, Seniority, etc. (shared domain language) |

---

## Data Stores

| Store | Type | Owned by | Purpose |
|-------|------|----------|---------|
| SQLite (data.db) | SQLite file | storage/sqlite/ | Persistent member data, monthly entries, audit log |
| In-memory map | Go map[int64]*TeamMember | storage/memory/ | Test fixture; no persistence |

---

## Key Constraints

- **No external network:** All computation local to single binary
- **Single binary:** Go application; database is a file (SQLite)
- **Hexagonal boundaries:** Core business logic is infrastructure-agnostic
- **Dependency inversion:** Dependencies point inward (HTTP → core ← storage); core depends on abstractions (ports), not implementations (adapters)
- **No framework bloat:** No ORM, no DI container, no reflection; explicit wiring in main.go
- **Testability:** Core use cases testable in milliseconds with in-memory adapter; no database setup required

---

## Out of Scope

- Internal component structure of each container (C4 Level 3)
- Deployment topology
- CI/CD pipeline
- Monthly entry operations (separate container; uses similar hexagonal pattern)

---

## Update Policy

Update this file when:
- A container is added, removed, or its technology changes
- A communication path between containers changes
- A new external system dependency is added
- A critical constraint changes

Do NOT update for internal refactors, new features within an existing container, or test changes.

**Last updated:** 2026-09-16 — Hexagonal architecture (ports & adapters) implemented for member CRUD; coordinator pattern replaced by use cases with explicit port dependencies.

---

## Architecture Decision

See [ADR-20260916 — Hexagonal Architecture (Ports & Adapters) for Member CRUD](./adr/ADR-20260916-hexagonal-architecture-member-crud.md) for the decision rationale, alternatives considered, and consequences.
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
