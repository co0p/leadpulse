# Architecture

Team Impact Scorecard — structural overview.

---

## C4 Level 2 — Container View

```
┌─────────────────────────────────────────────────────────────┐
│  User's Machine                                             │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Team Impact Scorecard (Single Binary)               │  │
│  │                                                      │  │
│  │  ┌──────────────┐   calls   ┌──────────────────┐    │  │
│  │  │              │ ────────► │                  │    │  │
│  │  │   ui/        │           │   service/       │    │  │
│  │  │              │           │                  │    │  │
│  │  │  Fyne views, │           │  Use cases:      │    │  │
│  │  │  widgets,    │           │  AddMember,      │    │  │
│  │  │  event       │           │  SubmitMonth,    │    │  │
│  │  │  handlers    │           │  GenerateAlerts, │    │  │
│  │  │              │           │  ExportData.     │    │  │
│  │  │              │           │  No I/O.         │    │  │
│  │  │              │           └────────┬─────────┘    │  │
│  │  │              │                    │ calls        │  │
│  │  │              │                    ▼              │  │
│  │  │              │           ┌──────────────────┐    │  │
│  │  │              │           │                  │    │  │
│  │  │              │   calls   │   store/         │    │  │
│  │  │              │ ────────► │                  │    │  │
│  │  │              │           │  SQLite read/    │    │  │
│  │  │              │           │  write. Domain   │    │  │
│  │  │              │           │  type mapping.   │    │  │
│  │  │              │           └────────┬─────────┘    │  │
│  │  │              │                    │ reads/writes  │  │
│  │  │              │   calls            ▼              │  │
│  │  │              │ ────────► ┌──────────────────┐    │  │
│  │  │              │           │                  │    │  │
│  │  └──────────────┘           │   engine/        │    │  │
│  │                             │                  │    │  │
│  │   (Fyne replaceable         │  Pure functions. │    │  │
│  │    with CLI, web, etc.)     │  Scoring, norms, │    │  │
│  │                             │  alerts, trends. │    │  │
│  │                             │  No I/O.         │    │  │
│  │                             └──────────────────┘    │  │
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

No external network connections. No cloud. No server.
```

---

## Containers

### `service/` — Application Use Case Layer
- **Technology:** Go (pure functions, no I/O)
- **Responsibility:** Application use cases, organized by domain (e.g., `service/member/`, `service/monthly/`). Each use case is a transaction: fetch aggregates via repository, call domain services for cross-aggregate rules, call `engine/` for computation, persist via repository, return the result or error.
- **Dependency model:** All persistence is abstracted via repository interfaces defined in `engine/domain/`. Service layer receives these interfaces as dependencies (not `*sql.DB`). This enables testing with in-memory repositories (no database required).
- **Domain services:** Application services depend on domain services (`ScoringService`, `TrendService`, `ValidationService`, `MonthlyEntryService`) which enforce cross-aggregate rules before persistence. Domain services also accept repository interfaces only; no database coupling.
- **Service aggregation:** `ApplicationServices` struct aggregates all service interfaces (`MemberService`, `MonthlyService`) for injection into the UI layer. The UI depends on this struct, not on individual services or repositories.
- **Constraints:** No UI logic. No direct SQLite access. No direct database dependencies. Delegates persistence to repository interfaces, computation to `engine/`, cross-aggregate rule enforcement to domain services.
- **Implemented use cases:** `service/member/` (AddMember, ListMembers, EditMember, DeactivateMember), `service/monthly/` (CreateEntry, GetEntry, ListEntriesByMember, UpdateEntry, ComputeScores, GetTrends).
- **Planned use cases:** SubmitMonthlyEntry (wrapper around monthly service), GenerateAlerts, ExportTeamOverview.

### `ui/` — Desktop Interface (Replaceable)
- **Technology:** Go, Fyne v2 (`fyne.io/fyne/v2`)
- **Responsibility:** Render all screens (Overview Dashboard, Monthly Input Workspace, Member Detail, Alerts Center, Monthly Review, Settings). Handle user interaction. Delegate all business logic to `service/`.
- **Service dependency:** UI accepts `ApplicationServices` struct containing all service interfaces (`MemberService`, `MonthlyService`). This decouples the UI from persistence concerns; the UI never sees repositories or database connections.
- **Screens:** A (Overview), B (Monthly Input), C (Member Detail), D (Alerts Center), E (Monthly Review), F (Settings)
- **Implemented screens:** F (Settings) — team member list, add/edit/deactivate dialogs; B (Monthly Input) — stub screen registered and reachable via sidebar navigation
- **Navigation shell:** sidebar + header pattern. Left sidebar (fixed width) holds primary nav buttons (Overview, Input, Alerts, Review) and Settings at the bottom. Top header shows the current cycle month. Main content area swaps on nav selection. Shell is defined in `ui/app.go`; screens are registered in `ui/screens/registry.go`.
- **Initialization:** main.go creates repositories → creates services → aggregates services into ApplicationServices → passes to UI. This preserves the separation: UI depends on services only.
- **Constraints:** No business logic. No direct SQLite access. No repository access. No database connections. Calls `service/` for all use cases. Testable by swapping `service/` implementation. Can be replaced with a CLI, web UI, or any other presentation layer that calls the same `service/` interfaces.
- **Future alternative:** A CLI client (`cmd/cli/`) can implement the same `service/` interfaces for headless or scripted workflows.

### `store/` — Persistence Layer
- **Technology:** Go, `modernc.org/sqlite` (pure Go, no CGO)
- **Responsibility:** Implements repository interfaces defined in `engine/domain/`. All SQLite reads and writes. Schema definition and migrations. Reconstructs aggregates from database rows; persists aggregates to database tables.
- **Pattern:** Repository implementations (SQLiteTeamMemberRepository, SQLiteMonthlyEntryRepository) conform to repository interfaces defined in the domain layer. This keeps storage knowledge in the store layer while the domain layer remains storage-agnostic.
- **Data retained:** 24 months of monthly entries per member, evidence notes, action plans, audit trail, cycle finalization state.
- **Constraints:** No formula computation. No persistence of raw domain objects; only aggregates via repository interface. Repositories validate aggregate state before persisting.

### `engine/` — Formula Engine
- **Technology:** Go (pure functions, no external dependencies)
- **Responsibility:** All computation defined in the PRD. Normalization functions, dimension scores (DG/DP/DT/DO), Total Impact Index (TII), trend calculations (MA3/Delta1/Delta3/Vol3), alert evaluation, confidence score, decision guardrail checks.
- **Constraints:** No I/O of any kind. No imports from `store/` or `ui/`. Input and output are plain Go structs.

### `SQLite File` — Local Data Store
- **Location:** `$UserConfigDir/leadpulse/data.db`
- **Technology:** SQLite 3
- **Responsibility:** Durable storage of all application data on the user's device.
- **Constraints:** Never accessed directly by `ui/`. All access is through `store/`.

---

## Dependency Direction

```
ui  →  ApplicationServices  →  service  →  engine/domain (services + interfaces)
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
- `ApplicationServices` aggregates all service instances for presentation-layer injection.
- `ui` calls `service` via `ApplicationServices` only. No direct `store`, `repository`, or `engine` calls. This makes `ui` replaceable: any presentation layer (Fyne, CLI, web) can call the same `service/` interfaces by receiving `ApplicationServices`.

---

## Key Communication Paths

| From | To | What |
|---|---|---|
| `ui` (event handler) | `service` | Use case call (e.g., `service.AddMember(firstName, lastName, seniority)`) |
| `service` (use case) | `store` | Read/write domain entities |
| `service` (use case) | `engine` | Compute scores, alerts, trends, guardrails |
| `store` (read) | SQLite file | SQL queries returning domain types |
| `store` (write) | SQLite file | SQL inserts/updates, append-only audit log |
| `ui` (live preview) | `service` (stateless) | Call `engine` via `service` for real-time formula preview during input |

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

The service layer is the testability boundary. Tests isolate each layer and verify integration points.

**`engine/` — unit tests (pure functions)**
- No mocks needed. Direct function calls with inputs and assertions on outputs.
- 100% of formulas, alert thresholds, trend calculations, guardrail logic.
- Example: `TestNormalizeMorale_midRange`, `TestAlertBurnout_redThreshold`.
- Gate: all engine tests must pass before any service or UI code runs.

**`service/` — unit tests (mock store, real engine)**
- Mock the `store` interface. Real `engine` functions (no mocks).
- Each use case is a test: call the use case with inputs, verify it calls `store` methods in the right order with the right data, verify the result.
- Example: `TestAddMember_createsAndReturns`, `TestSubmitMonthlyEntry_computesAndPersists`.
- Gate: all service tests must pass before UI code is written. Services are the API contract.

**`store/` — integration tests (in-memory SQLite, real engine)**
- No UI involved. SQLite database is in-memory (`:memory:` DSN).
- Cover all CRUD operations, schema constraints, audit trail logging, soft delete behavior.
- Each test sets up fixtures, calls store methods, verifies state in the database.
- Example: `TestAddMember_persistsAndAuditLogs`, `TestQueryMembers_excludesInactive`.
- Gate: all store tests must pass before services or UI code runs.

**`ui/` — manual verification only**
- Fyne has limited headless test support. Manual smoke-testing on-screen is the pragmatic gate.
- Verify: screens render, buttons trigger service calls, data flows from service to display, navigation works.
- Acceptance tests in `.agent/increment.md` are manual user stories, not automated suites.
- Gate: acceptance criteria pass in manual testing on the target platform (macOS, Windows, Linux).

**UI Replaceability Test**
- If a second presentation layer (CLI, web UI) is added, it must pass the same service-level tests without modification.
- This proves the service layer is the true interface contract, not the Fyne UI.

---

## ADR Index

Architectural decisions for this project are recorded in `docs/adr/`. Key decisions:

- [ADR-20260913-go-fyne-desktop](adr/ADR-20260913-go-fyne-desktop.md) — Go + Fyne as the desktop UI framework
- [ADR-20260913-sqlite-local-storage](adr/ADR-20260913-sqlite-local-storage.md) — SQLite for local-only persistence
- [ADR-20260913-pure-go-sqlite-driver](adr/ADR-20260913-pure-go-sqlite-driver.md) — `modernc.org/sqlite` to avoid CGO

---

**Last updated:** 2026-09-13 — Documented Settings screen (Screen F) implementation and integrated service/member use cases
