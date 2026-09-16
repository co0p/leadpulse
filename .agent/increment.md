# Increment: Hexagonal Architecture (Ports & Adapters) for Member CRUD

## Use Case

When I add a new API endpoint or CLI command for member operations, I want to reuse application logic through explicit use cases that depend on port abstractions, so that business logic is completely independent of storage technology and presentation layer, enabling portability to CLI, mobile, or future clients.

## Goal

Refactor member CRUD to follow Ports & Adapters (Hexagonal Architecture): move business logic to infrastructure-independent `core/` packages (domain aggregates + ports), implement ports as storage adapters in `storage/` (in-memory and SQLite), and replace `MemberAPICoordinator` with explicit use cases, achieving true hexagonal separation of concerns.

## Branch

`increment/clean-architecture-member-crud`

## Acceptance Criteria

1. **Core business logic is infrastructure-independent** — `core/members/` contains aggregate (`member.go`), port/interface (`repository.go`), and use cases (`add_member.go`, `get_members.go`, `edit_member.go`, `deactivate_member.go`) with zero imports from `storage/`, `server/`, or Fyne
2. **Port (interface) defined** — `core/members/repository.go` defines `MemberRepository` port with `Save()`, `FindByID()`, `FindActive()`, `Deactivate()` methods; no implementation logic in port
3. **Storage adapters implemented** — `storage/memory/member_repository.go` (in-memory adapter) and `storage/sqlite/member_repository.go` (SQLite adapter) both implement `MemberRepository` port; adapters depend on core, never vice versa
4. **HTTP adapter calls use cases** — `server/handler_members.go` calls use cases directly (not coordinator); handlers remain thin adapters parsing HTTP and mapping errors; `MemberAPICoordinator` removed entirely
5. **All tests pass and dependencies are clean** — `go test -race ./...` passes with 21+ tests; no circular imports; dependency flow is `server/` → `core/` ← `storage/` (adapters depend inward on core; core has zero outbound dependencies)

## Acceptance-Test Intent

**End-to-end user journey (optional, advisory):**
- User action: Start server, call `POST /api/members` with member data
- Expected: HTTP adapter calls `AddMemberUseCase.Execute()`, which uses injected `MemberRepository` port, adapter implements port using SQLite storage, returns UUID member
- Evidence: Same HTTP behavior as before refactor; integration test passes with in-memory adapter injection; core business logic testable without touching storage or HTTP layers

## Out Of Scope

- Monthly entry CRUD refactoring (separate increment; monthly coordinator and service unchanged for now)
- Domain events, event sourcing, or event-driven architecture (future consideration)
- CQRS (Command/Query Responsibility Segregation) — future optimization
- CLI implementation or CLI argument parsing (separate increment; this refactoring just makes CLI possible)
- Archiving or cleaning up old `service/coordinator/`, `service/member/`, or `store/` code (deletion is a cleanup task after this increment lands)
- Changes to `engine/scoring/` layer (unchanged)
- Removing `engine/domain/` shared value objects; they remain for reuse by other bounded contexts

## Constitution Constraints

- **Hexagonal architecture** — `core/` is infrastructure-agnostic; `storage/` and `server/` are adapters; core has zero dependencies on adapters
- **Port definition** — interfaces (ports) live in `core/` only; implementations (adapters) live in `storage/` or `server/`; adapters depend on core, never vice versa
- **Dependency inversion** — use cases depend on port interfaces, not concrete implementations; ports are injected at composition root (`main.go`)
- **Testing** — use case unit tests inject in-memory adapter (from `storage/memory/`); no database access per use case test; integration tests wire real SQLite adapter; all tests pass with `-race` flag
- **No Fyne imports** — `core/` remains UI-framework-agnostic; Fyne imports only in deprecated `ui/` package
- **Coordinator removal** — `service/coordinator/member_api.go` removed entirely in this increment; no coexistence
- **No infrastructure leakage** — `core/members/` must not import from `storage/`, `server/`, `engine/`, or `service/` packages

## Roadmap Entry

**Feature:** Clean Architecture for Member Operations  
**Job story:** When I add a new API endpoint or CLI command for member operations, I want to reuse application logic through explicit use cases and repositories instead of coordinators, so that business logic is clearly separated from presentation layers and easier to test, extend, and understand.  
**Status:** Partial

---

## Architecture & Package Structure

### Current State (To Be Replaced)
```
server/handler_members.go    → service/coordinator/MemberAPICoordinator
                                ↓
                           service/member/Service
                                ↓
                           store/SQLiteTeamMemberRepository
```

### Desired State: Ports & Adapters (Hexagonal Architecture)

**Core (Infrastructure-Independent Business Logic):**
```
core/members/
  member.go             # TeamMember aggregate (domain entity)
  repository.go         # MemberRepository PORT (interface definition only)
  add_member.go         # AddMemberUseCase
  get_members.go        # GetMembersUseCase
  edit_member.go        # EditMemberUseCase
  deactivate_member.go  # DeactivateMemberUseCase
  *_test.go             # Use case unit tests (inject in-memory adapter)
```

**Adapters (Infrastructure Implementations):**
```
storage/
  memory/
    member_repository.go    # InMemoryMemberRepository (adapter; implements MemberRepository port)
  sqlite/
    member_repository.go    # SQLiteTeamMemberRepository (adapter; implements MemberRepository port)
    
server/
  handler_members.go        # HTTP adapter (thin layer; parses HTTP, calls use cases, maps errors)
  *_test.go                 # Handler tests (mock use cases)
```

**Unchanged (Support Layers):**
```
engine/
  scoring/                  # Pure computation
  domain/                   # Shared value objects (Seniority, FullName, etc.)
  
main.go                     # Wires adapters + core + domain
```

### Dependency Flow: Hexagonal/Ports & Adapters

```
┌─────────────────────┐
│   HTTP (Adapter)    │ server/handler_members.go
│  ↓ (calls)          │
├─────────────────────┤
│   CORE (Business)   │ core/members/
│  ┌─────────────────┐│  - Aggregate: member.go
│  │ Use Cases       ││  - Port: repository.go
│  │ (depend on      ││  - Use Cases: add_member.go, etc.
│  │  ports)         ││
│  └─────────────────┘│
├─────────────────────┤
│ Ports (Interfaces)  │ core/members/repository.go
│ Zero impl logic     │
├─────────────────────┤
│  Storage (Adapter)  │ storage/{memory,sqlite}/member_repository.go
│  ↑ (implements)     │
└─────────────────────┘
```

**Rule:** 
- `core/` has ZERO imports from `storage/`, `server/`, or `service/`
- `storage/` and `server/` depend on `core/` (inbound only)
- No circular imports
- Core is reusable: same business logic works with any storage adapter, any presentation adapter (HTTP, CLI, gRPC, etc.)

### Use Case Pattern

**Input/Output Structs (in `core/members/`):**
```go
// core/members/add_member.go

type AddMemberInput struct {
  FirstName string
  LastName  string
  Seniority string
}

type AddMemberOutput struct {
  ID        string    // UUID
  FirstName string
  LastName  string
  Seniority string
  Status    string
  CreatedAt time.Time
}
```

**Port Definition (Interface, in `core/members/repository.go`):**
```go
// core/members/repository.go
// This is a PORT (not an adapter). No implementation here.

type MemberRepository interface {
  Save(member *Member) error
  FindByID(id int64) (*Member, error)
  FindActive() ([]*Member, error)
  Deactivate(id int64) error
}
```

**Use Case (in `core/members/`):**
```go
// core/members/add_member.go

type AddMemberUseCase struct {
  repo MemberRepository  // Injected port (adapter implements this)
}

func NewAddMemberUseCase(repo MemberRepository) *AddMemberUseCase {
  return &AddMemberUseCase{repo: repo}
}

func (uc *AddMemberUseCase) Execute(input AddMemberInput) (*AddMemberOutput, error) {
  // Validation (domain rules)
  // Business logic (create aggregate)
  // Persistence (call port interface)
  err := uc.repo.Save(member)
  if err != nil {
    return nil, err
  }
  return &AddMemberOutput{...}, nil
}
```

**Adapter Implementation (in `storage/sqlite/member_repository.go`):**
```go
// storage/sqlite/member_repository.go
// This is an ADAPTER (implements the MemberRepository port)

package sqlite

import "leadpulse/core/members"  // Depends on port interface

type SQLiteTeamMemberRepository struct {
  db *sql.DB
}

// Implements members.MemberRepository interface
func (r *SQLiteTeamMemberRepository) Save(member *members.Member) error {
  // Actual SQLite logic here
  stmt := r.db.Prepare("INSERT INTO members ...")
  return stmt.Exec(member.FirstName, ...).Error()
}

func (r *SQLiteTeamMemberRepository) FindByID(id int64) (*members.Member, error) {
  // Query logic
}

func (r *SQLiteTeamMemberRepository) FindActive() ([]*members.Member, error) {
  // Query logic
}

func (r *SQLiteTeamMemberRepository) Deactivate(id int64) error {
  // Update logic
}
```

**In-Memory Adapter (in `storage/memory/member_repository.go`):**
```go
// storage/memory/member_repository.go
// ADAPTER: in-memory implementation for testing

package memory

import "leadpulse/core/members"

type InMemoryMemberRepository struct {
  members map[int64]*members.Member
}

// Implements members.MemberRepository interface
func (r *InMemoryMemberRepository) Save(member *members.Member) error {
  r.members[member.ID] = member
  return nil
}

func (r *InMemoryMemberRepository) FindByID(id int64) (*members.Member, error) {
  m, ok := r.members[id]
  if !ok {
    return nil, ErrNotFound
  }
  return m, nil
}
// ... etc
```

**HTTP Adapter (in `server/handler_members.go`):**
```go
// server/handler_members.go
// ADAPTER: HTTP layer (entry point)

package server

import (
  "leadpulse/core/members"
  "encoding/json"
  "net/http"
)

type AddMemberRequest struct {
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Seniority string `json:"seniority"`
}

func HandlerAddMember(w http.ResponseWriter, r *http.Request, addUC *members.AddMemberUseCase) {
  var req AddMemberRequest
  json.NewDecoder(r.Body).Decode(&req)
  
  output, err := addUC.Execute(members.AddMemberInput{
    FirstName: req.FirstName,
    LastName:  req.LastName,
    Seniority: req.Seniority,
  })
  
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
    return
  }
  
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(output)
}
```

**Wiring (in `main.go`):**
```go
// main.go
// Composition root: wire adapters + core + ports

package main

import (
  "leadpulse/core/members"
  "leadpulse/storage/sqlite"
  "leadpulse/server"
)

func main() {
  // Create adapter
  repo := sqlite.NewSQLiteTeamMemberRepository(db)
  
  // Create use cases (inject port)
  addMemberUC := members.NewAddMemberUseCase(repo)
  getMembersUC := members.NewGetMembersUseCase(repo)
  editMemberUC := members.NewEditMemberUseCase(repo)
  deactivateMemberUC := members.NewDeactivateMemberUseCase(repo)
  
  // Wire HTTP adapters
  mux := http.NewServeMux()
  mux.HandleFunc("POST /api/members", 
    func(w http.ResponseWriter, r *http.Request) {
      server.HandlerAddMember(w, r, addMemberUC)
    })
  
  // Start server
  http.ListenAndServe(":8080", mux)
}
```

**Use Case Unit Test (in `core/members/add_member_test.go`):**
```go
// core/members/add_member_test.go
// Unit test: inject in-memory adapter (no database, no HTTP)

package members

import (
  "testing"
  "leadpulse/storage/memory"
)

func TestAddMemberUseCase_Success(t *testing.T) {
  // Inject in-memory adapter
  repo := memory.NewInMemoryMemberRepository()
  
  uc := NewAddMemberUseCase(repo)
  
  output, err := uc.Execute(AddMemberInput{
    FirstName: "Alice",
    LastName:  "Smith",
    Seniority: "Senior",
  })
  
  if err != nil {
    t.Fatalf("unexpected error: %v", err)
  }
  
  if output.FirstName != "Alice" {
    t.Errorf("expected Alice, got %s", output.FirstName)
  }
  
  // Verify persistence via adapter
  retrieved, err := repo.FindByID(output.ID) // Direct repo call, or through use case
  if err != nil {
    t.Fatalf("member not persisted: %v", err)
  }
  
  if retrieved.FirstName != "Alice" {
    t.Errorf("persisted member has wrong name: %s", retrieved.FirstName)
  }
}
```

### Migration Path

1. **Move aggregate:** `TeamMember` from `engine/domain/aggregates.go` → `core/members/member.go` (or keep in engine/domain and import from core/members if shared)
2. **Define port:** Create `core/members/repository.go` with `MemberRepository` interface
3. **Create in-memory adapter:** `storage/memory/member_repository.go` (for tests)
4. **Create SQLite adapter:** `storage/sqlite/member_repository.go` (extract from `store/member.go`)
5. **Create use cases:** `core/members/{add,get,edit,deactivate}_member.go` with `Execute()` methods
6. **Create use case tests:** `core/members/*_test.go` (inject in-memory adapter)
7. **Refactor HTTP adapters:** `server/handler_members.go` to call use cases (not coordinator)
8. **Refactor handler tests:** `server/handler_members_test.go` to mock use cases
9. **Update main.go:** Wire adapters + use cases at composition root
10. **Delete coordinator:** Remove `service/coordinator/member_api.go` entirely
11. **Run full test suite:** `go test -race ./...` passes with 21+ tests

## Testing Strategy

**Use case unit tests (in `core/members/*_test.go`):**
- Inject in-memory adapter from `storage/memory/`
- Call use case `Execute()`, verify output and side effects
- No database, no HTTP, no I/O → milliseconds
- Coverage: success path + error cases (validation, not-found, duplicate)

**Storage adapter tests (in `storage/{memory,sqlite}/*_test.go`):**
- Test CRUD operations for each adapter
- Memory adapter: instant, no setup needed
- SQLite adapter: uses `:memory:` DSN or test database
- Verify adapters correctly implement port interface

**HTTP adapter tests (in `server/handler_members_test.go`):**
- Mock use cases (not the real core or storage layer)
- Verify HTTP parsing, status codes, JSON responses
- No database, no storage layer; tests remain fast

**Integration tests (optional, verify end-to-end):**
- Real SQLite adapter + real use cases + HTTP adapter
- Verify full CRUD workflow
- Slower than unit tests; run separately if needed

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Circular imports (core imports storage) | Build failure; architecture violation | Strict review; interfaces (ports) in core only; adapters depend inward; ban any core→storage imports |
| Core becomes too large | Bounded context bleeding | Start with member aggregate only; each use case gets one file; future: split into separate packages per BC |
| Coordinator removal breaks during refactoring | Test failures; incomplete cutover | Adapter wiring in main.go tested before coordinator removal; all tests pass before deletion |
| More adapters = more mock types in tests | Test boilerplate increase | Each adapter is simple; in-memory adapter < 50 lines; minimal mocking needed |
| Package import paths become verbose | Cognitive load | Accepted trade-off; clear architecture outweighs path length; IDE autocomplete helps |
| Storing `MemberRepository` port in `core/members/` alongside use cases | Package size | Intentional design; port definition lives where it's used (core/members); keeps interfaces and implementations coupled at right abstraction level |

## Evidence & References

**Hexagonal Architecture / Ports & Adapters (Alistair Cockburn):**
- Core business logic is independent of infrastructure (database, UI, frameworks)
- Ports are interfaces; adapters implement ports
- Multiple adapters (in-memory, SQLite, Postgres, REST, CLI, gRPC) can work with same core
- Tests inject test adapters; production wires real adapters at composition root

**Domain-Driven Design (Eric Evans):**
- Aggregates (TeamMember) encapsulate business rules
- Repositories abstract persistence (interface in core, implementations external)
- Bounded contexts (member operations = one context) isolate responsibility

**Clean Architecture (Robert C. Martin):**
- Dependency rule: dependencies point inward (adapters → core → nothing)
- Core has zero dependencies on frameworks or libraries
- Testability: core logic tested without touching infrastructure

**Go Community Practices:**
- Kubernetes: core business logic in `pkg/`, adapters in `cmd/` and controllers
- Docker: engine/container (core) vs. cli (adapter) vs. daemon (adapter)
- Standard library: `io.Writer` is a port; `os.File`, `bytes.Buffer` are adapters

**Industry Adoption:**
- Microservices: enables CLI, gRPC, GraphQL, REST all from same core
- Enterprise systems: core portable across UI frameworks, databases, cloud providers
- Testable design: core tests run in milliseconds; integration tests optional

