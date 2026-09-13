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
│  │  │   ui/        │           │   store/         │    │  │
│  │  │              │           │                  │    │  │
│  │  │  Fyne views, │   calls   │  SQLite read/    │    │  │
│  │  │  widgets,    │ ────────► │  write. Domain   │    │  │
│  │  │  event       │           │  type mapping.   │    │  │
│  │  │  handlers    │           │  No formulas.    │    │  │
│  │  │              │           └────────┬─────────┘    │  │
│  │  │              │                    │ reads/writes  │  │
│  │  │              │   calls            ▼              │  │
│  │  │              │ ────────► ┌──────────────────┐    │  │
│  │  │              │           │                  │    │  │
│  │  └──────────────┘           │   engine/        │    │  │
│  │                             │                  │    │  │
│  │                             │  Pure functions. │    │  │
│  │                             │  Scoring, norms, │    │  │
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

### `ui/` — Desktop Interface
- **Technology:** Go, Fyne v2 (`fyne.io/fyne/v2`)
- **Responsibility:** Render all screens (Overview Dashboard, Monthly Input Workspace, Member Detail, Alerts Center, Monthly Review, Settings). Handle user interaction. Delegate computation to `engine/`, delegate persistence to `store/`.
- **Screens:** A (Overview), B (Monthly Input), C (Member Detail), D (Alerts Center), E (Monthly Review), F (Settings)
- **Constraints:** No business logic. No direct SQLite access. Calls `engine/` for live formula preview during input.

### `store/` — Persistence Layer
- **Technology:** Go, `modernc.org/sqlite` (pure Go, no CGO)
- **Responsibility:** All SQLite reads and writes. Schema definition and migrations. Returns domain types; never raw SQL results.
- **Data retained:** 24 months of monthly entries per member, evidence notes, action plans, audit trail, cycle finalization state.
- **Constraints:** No formula computation. Receives computed values from the caller and persists them.

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
ui  →  store  →  engine  (types only)
ui  →  engine
```

`engine` has no project-internal dependencies.  
`store` may use engine types for its function signatures but must not call engine computation functions.  
`ui` may call both `store` and `engine` directly (e.g., for live preview during input).

---

## Key Communication Paths

| From | To | What |
|---|---|---|
| `ui` (input form) | `engine` | Live formula recompute as user types |
| `ui` (save action) | `store` | Persist monthly entry + computed scores |
| `ui` (load screen) | `store` | Fetch member list, historical entries, alerts, evidence |
| `store` (read) | SQLite file | SQL queries |
| `store` (write) | SQLite file | SQL inserts/updates, append-only audit log |

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

## ADR Index

Architectural decisions for this project are recorded in `docs/adr/`. Key decisions:

- [ADR-20260913-go-fyne-desktop](adr/ADR-20260913-go-fyne-desktop.md) — Go + Fyne as the desktop UI framework
- [ADR-20260913-sqlite-local-storage](adr/ADR-20260913-sqlite-local-storage.md) — SQLite for local-only persistence
- [ADR-20260913-pure-go-sqlite-driver](adr/ADR-20260913-pure-go-sqlite-driver.md) — `modernc.org/sqlite` to avoid CGO
