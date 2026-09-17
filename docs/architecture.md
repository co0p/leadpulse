# Architecture

Team Impact Scorecard — multi-service containerized architecture (Vue 3 SPA frontend container, API-only backend container, orchestrated via Docker Compose).

---

## C4 Level 2 — Container View

```
┌──────────────────────────────────────────────────────────────────────────┐
│  Developer's Machine / Local Docker Environment                          │
│                                                                          │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │  Docker Network (leadpulse-network)                              │   │
│  │                                                                  │   │
│  │  ┌─────────────────────────┐         ┌──────────────────────┐   │   │
│  │  │  Frontend Container     │         │  Backend Container   │   │   │
│  │  │  (services/frontend/)   │         │  (services/backend/) │   │   │
│  │  │                         │         │                      │   │   │
│  │  │  ┌─────────────────┐    │         │  ┌────────────────┐  │   │   │
│  │  │  │ Nginx Server    │    │         │  │ Go Binary      │  │   │   │
│  │  │  │ Port 3000       │    │         │  │ Port 8080      │  │   │   │
│  │  │  └────────┬────────┘    │         │  └────────┬───────┘  │   │   │
│  │  │           │             │         │           │          │   │   │
│  │  │  ┌────────▼────────┐    │         │  ┌────────▼────────┐ │   │   │
│  │  │  │ Vue 3 + Vite    │    │         │  │ HTTP API        │ │   │   │
│  │  │  │ SPA (static)    │    │         │  │ Handler Layer   │ │   │   │
│  │  │  │                 │    │         │  │ (server/)       │ │   │   │
│  │  │  │ • AppShell      │    │         │  │                 │ │   │   │
│  │  │  │ • Health Indicator◄──┼─────────┼──┤ /api/health     │ │   │   │
│  │  │  │ • Bulma Styling │    │ HTTP/  │  │ /api/members    │ │   │   │
│  │  │  │                 │    │ JSON   │  │ /api/...        │ │   │   │
│  │  │  └─────────────────┘    │ over   │  └────────┬────────┘ │   │   │
│  │  │                         │ docker-│           │          │   │   │
│  │  │ Proxy: /api/* ──────────┼─ compose         │ depends  │   │   │
│  │  │         to backend      │ network└───────────┤ on port  │   │   │
│  │  └─────────────────────────┘         │          │          │   │   │
│  │                                     │  ┌────────▼────────┐ │   │   │
│  │  Docker Compose                     │  │ core/ members/  │ │   │   │
│  │  • Volumes: data/ (SQLite)          │  │                 │ │   │   │
│  │  • Networks: leadpulse-network      │  │ Use Cases       │ │   │   │
│  │  • Health Checks: both services     │  │ Add, Get, Edit  │ │   │   │
│  │                                     │  │ Deactivate, etc │ │   │   │
│  │                                     │  └────────┬────────┘ │   │   │
│  │                                     │           │          │   │   │
│  │                                     │  ┌────────▼────────┐ │   │   │
│  │                                     │  │ core/members/   │ │   │   │
│  │                                     │  │ PORT (repo      │ │   │   │
│  │                                     │  │ interface)      │ │   │   │
│  │                                     │  └────────┬────────┘ │   │   │
│  │                                     │           │ impl     │   │   │
│  │                                     │  ┌────────▼──────┐   │   │   │
│  │                                     │  │ storage/sqlite│   │   │   │
│  │                                     │  │ Repository    │   │   │   │
│  │                                     │  └────────┬──────┘   │   │   │
│  │                                     │           │          │   │   │
│  │  ┌─────────────────────────────────┼───────────▼──────┐   │   │   │
│  │  │  Volume Mount: /data (in container)                │   │   │   │
│  │  │                                                     │   │   │   │
│  │  │  ┌───────────────────────────────────────────────┐ │   │   │   │
│  │  │  │ SQLite Database File                           │ │   │   │   │
│  │  │  │ /data/data.db (persistent, host: ./data/)     │ │   │   │   │
│  │  │  │ — members, monthly entries, audit log, state  │ │   │   │   │
│  │  │  └───────────────────────────────────────────────┘ │   │   │   │
│  │  └─────────────────────────────────────────────────────┘   │   │   │
│  │                                                            │   │   │
│  └────────────────────────────────────────────────────────────┘   │   │
│                                                                    │   │
│ Browser: http://localhost:3000                                    │   │
│                                                                    │   │
└────────────────────────────────────────────────────────────────────────┘

docker-compose up --build:
  1. Builds leadpulse-backend container (multi-stage: test, build, runtime)
  2. Builds leadpulse-frontend container (multi-stage: build, nginx)
  3. Starts both services on shared docker-compose network
  4. Mounts ./data/ volume for SQLite persistence
  5. Frontend at http://localhost:3000
  6. Backend API at http://localhost:8080 (reachable from frontend as http://backend:8080)
```

---

## Containers

## Services

### Frontend Service (`services/frontend/`)

**Container:** `leadpulse-frontend:latest`
- **Technology:** Node.js 20 + Vue 3 + Vite + Nginx
- **Port:** 3000 (exposed to localhost)
- **Build:** Multi-stage Dockerfile
  - Stage 1 (build): `npm install`, `npm run build` → `dist/`
  - Stage 2 (runtime): Nginx serves static `dist/` files; proxies `/api/*` to backend
- **Startup:** Nginx server on 0.0.0.0:3000
- **Environment:** `VITE_API_URL=http://backend:8080` (injected at runtime, used by AppShell for health checks)
- **Health Check:** HTTP 200 on `http://localhost:3000/`
- **Network:** Joined to `leadpulse-network` docker-compose network; can reach backend as `http://backend:8080`

**Frontend Stack:**
- **Vue 3 + Vite:** Build tool-driven single-page app (SPA)
- **Vue Router:** Client-side routing for /members, /alerts, /reports, etc. Root handler (/) serves index.html; Vue Router handles rest
- **Pinia:** State management (health check status, lastCheckedAt, checkError)
- **Bulma CSS:** Responsive design (fixed sidebar ≥1024px, overlay mobile ≤768px)
- **Components:**
  - **AppShell:** Persistent layout wrapper (sidebar, topbar, main content, footer); renders `<RouterView />` for screen content
  - **HealthIndicator:** Badge in footer; three states: healthy (green ✓), unhealthy (red ✗), checking (spinner)
- **Health Check:** Async `checkHealth()` action runs on app mount (non-blocking); fetches `/api/health` with 5s timeout; updates store state (healthy | unhealthy | checking)

**Responsibility:**
- Serve Vue 3 SPA with persistent AppShell layout
- Route static assets (CSS, JS, images)
- Proxy `/api/*` requests to backend service
- Display real-time health status of backend via HealthIndicator badge
- Execute client-side navigation without full page reloads (Vue Router handles routing)
- No server-side rendering or templating

---

## Health Check Flow (Vue SPA to Backend)

**Sequence Diagram:**

```
┌──────────┐                          ┌─────────────────┐          ┌──────────┐
│ Browser  │                          │ Vue App (Pinia) │          │ Backend  │
│ (User)   │                          │ + HealthIndicator          │ /api/    │
└────┬─────┘                          └────────┬────────┘          └────┬─────┘
     │                                         │                       │
     │ 1. Load http://localhost:3000           │                       │
     ├────────────────────────────────────────►│                       │
     │                                         │                       │
     │ 2. Render AppShell                      │                       │
     │ (HealthIndicator in footer,             │                       │
     │  status: 'checking', spinner visible)   │                       │
     │◄────────────────────────────────────────┤                       │
     │                                         │                       │
     │                                  3. onMounted() fires           │
     │                                  checkHealth() action           │
     │                                         │                       │
     │                                  4. Pinia dispatch             │
     │                                  (status: 'checking')           │
     │                                         │                       │
     │                                         │ 5. Fetch /api/health  │
     │                                         ├──────────────────────►│
     │                                         │                       │
     │                                         │                       │
     │                                    6. HTTP 200 {status: "ok"}   │
     │                                         │◄──────────────────────┤
     │                                         │                       │
     │                                  7. Update store               │
     │                                  (status: 'healthy',           │
     │                                   lastCheckedAt: now)          │
     │                                         │                       │
     │ 8. HealthIndicator watches store      │                       │
     │ and re-renders (green ✓)               │                       │
     │◄────────────────────────────────────────┤                       │
     │                                         │                       │
     │ (No re-check on subsequent navigation)  │                       │
     │                                         │                       │
```

**Key characteristics:**
- Health check runs **once** on app mount, **not** on every navigation
- Status stored in Pinia (survives client-side route changes)
- Async/non-blocking: UI renders while fetch is in-flight
- 5-second timeout on fetch (AbortController)
- Error states stored (`checkError` field for display/logging)

**State transitions:**
| Event | Status Before | Status After | Display |
|-------|---------------|--------------|---------|
| App mounts | (none) | checking | Spinner |
| /api/health returns 200 | checking | healthy | Green ✓ |
| /api/health fails or timeout | checking | unhealthy | Red ✗ + error message |
| User navigates (Vue Router) | healthy/unhealthy | (unchanged) | Persists from last check |

### Backend Service (`services/backend/`)

**Container:** `leadpulse-backend:latest`
- **Technology:** Go 1.23 + hexagonal architecture
- **Port:** 8080 (exposed to localhost and docker network)
- **Build:** Multi-stage Dockerfile
  - Stage 1 (test): `go test -race ./...` (must pass before build proceeds)
  - Stage 2 (build): `go build` → binary
  - Stage 3 (runtime): Alpine-based; binary + health check script
- **Startup:** Go binary on 0.0.0.0:8080
- **Environment:** `DATABASE_URL="sqlite:///data/data.db"` (or fallback to ~/.config/leadpulse/data.db if /data not available)
- **Volume Mount:** `/data` → persistent SQLite database
- **Health Check:** HTTP 200 on `http://localhost:8080/api/health`
- **Network:** Joined to `leadpulse-network`; reachable from frontend as `backend:8080`

**Responsibility:**
- Expose JSON REST API endpoints (`/api/health`, `/api/members`, etc.)
- Apply business logic via use cases (hexagonal core)
- Persist data to SQLite `/data/data.db`
- No HTML rendering; no server-side templates; API-only

### Database Volume

**Mount:** `data:` volume mapped to `/data/` inside backend container
- **Persistence:** Survives `docker-compose down` cycles
- **Host Path:** `./data/` (relative to project root)
- **Contents:** `data.db` (SQLite file)
- **Lifecycle:** First `docker-compose up` creates empty volume; subsequent runs preserve data
- **Cleanup:** `docker-compose down -v` deletes volume

---

## Hexagonal Architecture (Inside Backend)

**Same as before; unchanged by containerization:**
- **Handlers:** Receive HTTP requests, parse JSON/path/query data, call use cases via dependency injection, map results to JSON HTTP responses. No business logic; all logic in core use cases. No template rendering, no layout state, no knowledge of how the response is displayed.
- **Constraints:** No direct store or engine access; no database connections. All persistence through injected use case dependencies. No external network calls (all API responses are JSON; static SPA assets are embedded, not fetched from a CDN).

### Vue SPA — Frontend Container
- **Technology:** Vue 3, Vite (build tool), Vue Router (client-side routing), Pinia (state management), Bulma (CSS).
- **Repository location:** a separate frontend project (own `package.json`, own test suite), built independently of the Go module.
- **Responsibility:** All rendering, navigation, and client-side state. Owns the persistent application shell (sidebar, top bar) as a layout component wrapping routed screens — shell/content composition is a Vue Router concern, not a backend concern. Calls the backend exclusively through `/api/*` JSON endpoints.
- **Build integration:** `npm run build` produces a static `dist/` directory. That directory is embedded into the Go binary via `embed.FS` before `go build`. The Go binary serves these static files directly; no Node.js process runs alongside the shipped binary.
- **Constraints:** No direct database or filesystem access. No knowledge of Go types or internal package structure. Communicates with the backend only through the documented JSON API contract — the same contract any other client (CLI, future integrations) would use.
- **Decision record:** see `docs/adr/ADR-20260917-vue-spa-frontend.md` for the full rationale and alternatives considered.

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
| Vue SPA | server/ | HTTP (localhost:8080) | JSON API |
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

## Architecture Decision

See [ADR-20260916 — Hexagonal Architecture (Ports & Adapters) for Member CRUD](./adr/ADR-20260916-hexagonal-architecture-member-crud.md) for the core/use-case/repository decision, and [ADR-20260917 — Vue SPA Frontend, API-only Backend](./adr/ADR-20260917-vue-spa-frontend.md) for the frontend/backend boundary decision.

- **Key constraint:** Stateless per HTTP request. No member variables that persist across calls except injected use case/repository references.

## Testing Strategy by Layer

**`engine/` — unit tests (pure functions)**
- No mocks needed. Direct function calls with inputs and assertions on outputs.
- 100% of formulas, alert thresholds, trend calculations, guardrail logic.
- Example: `TestNormalizeMorale_midRange`, `TestAlertBurnout_redThreshold`.
- Gate: all engine tests must pass before any core use case code runs.

**`core/` — unit tests (in-memory repository, real engine)**
- Inject an in-memory repository implementation (`storage/memory/`). Real `engine` functions (no mocks).
- Each use case is a test: call `Execute()` with inputs, verify repository state and return values.
- Example: `TestAddMemberUseCase_createsAndReturns`, `TestReactivateMemberUseCase_reactivatesInactiveMember`.
- Gate: all core tests must pass before HTTP handlers are written. Use cases are the API contract.

**`storage/` — integration tests (real SQLite, in-memory DSN)**
- Use `":memory:"` SQLite DSN. Verify each repository method persists and reconstructs aggregates correctly.
- Example: `TestSQLiteTeamMemberRepository_Save_Inserts_NewMember`.
- Gate: all storage tests must pass. Storage adapters are swappable only if they pass the same behavioral contract as the in-memory adapter.

**`server/` — unit tests (mocked use cases + httptest)**
- Handler tests use `httptest.NewRecorder` to test request/response mapping without a live server.
- Handlers are wired to mocked use cases (hand-written test doubles) to verify correct delegation and JSON mapping.
- Example: `TestHandlerAddMember_createsAndReturnsJSON`, `TestHandlerGetMembers_filtersbyStatus`.
- Gate: all handler tests must pass. Handlers are the JSON API surface; tests assert on JSON structure only, never on HTML.

**Vue SPA — component and edge-case tests (Vitest + Vue Test Utils)**
- Component-level tests cover rendering, props, emitted events, form validation, and error states without a browser.
- Live in the frontend project's own test suite, isolated from the Go test suite.
- Gate: component tests pass before a screen is considered complete.

**Vue SPA — acceptance tests (Playwright, narrow scope)**
- Reserved for one end-to-end test per feature, covering only the main success flow (e.g., "add a member" happy path) against a live server and real database.
- Not used for edge cases, validation errors, or exhaustive interaction coverage — those belong to Vitest/Vue Test Utils.
- Gate: the main success flow of each shipped feature has a passing Playwright test before the feature is marked Done in `docs/roadmap.md`.

---

## ADR Index

Architectural decisions for this project are recorded in `docs/adr/`. Key decisions:

- [ADR-20260917-vue-spa-frontend.md](adr/ADR-20260917-vue-spa-frontend.md) — Vue 3 SPA frontend, API-only backend (supersedes the HTMX/Alpine decision)
- [ADR-20260916-hexagonal-architecture-member-crud.md](adr/ADR-20260916-hexagonal-architecture-member-crud.md) — Hexagonal architecture (ports & adapters) for member CRUD
- [ADR-20260915-clean-architecture-layering.md](adr/ADR-20260915-clean-architecture-layering.md) — Clean architecture layering (historical; superseded in part by the hexagonal ADR above)
- [ADR-20260915-spa-frontend-stack.md](adr/ADR-20260915-spa-frontend-stack.md) — HTMX + Alpine.js + Go templates (superseded by ADR-20260917)
- [ADR-20260913-sqlite-local-storage.md](adr/ADR-20260913-sqlite-local-storage.md) — SQLite for local-only persistence
- [ADR-20260913-pure-go-sqlite-driver.md](adr/ADR-20260913-pure-go-sqlite-driver.md) — `modernc.org/sqlite` to avoid CGO

---

## Update Policy

Update this file when:
- A service (container) is added, removed, or its technology changes
- A communication path between services changes
- A new external system dependency is added
- A critical constraint changes

Do NOT update for internal refactors, new features within an existing service, or test changes.

**Last updated:** 2026-09-17 — Transitioned from single-binary desktop app to multi-service Docker architecture. Frontend and backend now in separate containers, orchestrated via docker-compose. Database persistence via volume mounts. All communication over HTTP/JSON network API.
