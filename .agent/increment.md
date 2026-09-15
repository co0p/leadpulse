# Increment: Clean Architecture for Member CRUD

## Use Case

When I add a new API endpoint or CLI command for member operations, I want to reuse application logic through explicit use cases and repositories instead of coordinators, so that business logic is clearly separated from presentation layers and easier to test, extend, and understand.

## Goal

Refactor member CRUD to follow clean architecture and domain-driven design: replace `MemberAPICoordinator` with explicit use cases (`AddMemberUseCase`, `EditMemberUseCase`, etc.) that depend on repository abstractions, achieving testability, portability, and clear DDD vocabulary.

## Branch

`increment/clean-architecture-member-crud`

## Acceptance Criteria

1. **Use cases exist and are tested** — `AddMemberUseCase`, `GetMembersUseCase`, `EditMemberUseCase`, `DeactivateMemberUseCase` all exist in `application/member/` with passing unit tests; each has an `Execute()` method with input/output structs
2. **Repository abstraction defined** — `domain/member/repository.go` defines `MemberRepository` interface with `Save()`, `FindByID()`, `FindActive()`, `Deactivate()` methods; no implementation logic in interface
3. **Storage implementations provided** — `storage/memory/member_repository.go` (in-memory) and `storage/sqlite/member_repository.go` (SQLite) both implement `MemberRepository` interface; in-memory is used for tests, SQLite for production
4. **HTTP handlers call use cases** — `server/handler_members.go` calls `AddMemberUseCase.Execute()`, `GetMembersUseCase.Execute()`, etc.; handlers remain thin adapters parsing HTTP requests and mapping errors
5. **All tests pass and dependencies are clean** — `go test -race ./...` passes with 21+ server tests; no circular imports; dependency flow is presentation → application → domain (never reversed)

## Acceptance-Test Intent

**End-to-end user journey (optional, advisory):**
- User action: Start server, call `POST /api/members` with member data
- Expected: Handler calls `AddMemberUseCase.Execute()`, which uses injected `MemberRepository` to persist, returns UUID member
- Evidence: Same HTTP behavior as before refactor; integration test passes with in-memory repository injection

## Out Of Scope

- Monthly entry CRUD refactoring (separate increment; monthly coordinator unchanged for now)
- Domain events, event sourcing, or event-driven architecture (deferred to v2)
- CQRS (Command/Query Responsibility Segregation) — future optimization, not required for v1
- CLI implementation or CLI argument parsing (future increment; this refactoring just makes CLI possible)
- Removing or archiving old `service/coordinator/` code (can coexist during transition; deletion is a separate cleanup task)
- Changes to `engine/`, `store/`, or `service/member/` layers (they remain as-is for now)

## Constitution Constraints

- **Layer separation** — handlers in `server/` must not contain business logic; business logic lives in `application/` (use cases) and `domain/`
- **Dependency direction** — `server/` → `application/` → `domain/`; no reverse imports; `domain/` has zero dependencies on `application/` or `infrastructure/`
- **Testing** — use case unit tests use in-memory repositories (injected); no database access per use case test; integration tests wire real SQLite storage; all tests pass with `-race` flag
- **Repository interface** — defined in `domain/` only; implementations live in `storage/memory/` and `storage/sqlite/`, never in domain package
- **No Fyne imports** — application and domain layers remain UI-framework-agnostic

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

### Desired State (After Refactoring)
```
Presentation Layer:
  server/handler_members.go  (thin HTTP adapters; parse request → call use case → serialize response)

Application Layer (Use Cases):
  application/member/
    add_member.go            (AddMemberUseCase)
    get_members.go           (GetMembersUseCase)
    edit_member.go           (EditMemberUseCase)
    deactivate_member.go     (DeactivateMemberUseCase)

Domain Layer (Business Rules & Abstractions):
  domain/member/
    aggregate.go             (TeamMember aggregate root; unchanged from engine/domain/)
    repository.go            (MemberRepository interface definition only; no implementation)

Infrastructure Layer (Implementations & External APIs):
  storage/memory/
    member_repository.go     (InMemoryMemberRepository; for tests)
  storage/sqlite/
    member_repository.go     (SQLiteTeamMemberRepository; production)
  server/
    handler_members.go       (HTTP layer; depends on use cases)
```

### Dependency Flow (Strictly One Direction)
```
server/ (HTTP)
  ↓ depends on
application/member/ (Use Cases)
  ↓ depends on
domain/member/ (Repository interface + Aggregate)
  
storage/memory/ & storage/sqlite/ (Implementations)
  ↓ depend on
domain/member/ (interface only)
```

**Rule:** No package above can import packages below. `domain/` has zero imports from `application/`, `storage/`, or `server/`.

### Use Case Pattern

**Input/Output Structs:**
```go
// application/member/input_output.go

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

type GetMembersOutput struct {
  Members []AddMemberOutput
}

type EditMemberInput struct {
  MemberID  int64  // domain ID
  FirstName string // optional
  LastName  string // optional
  Seniority string // optional
}

type EditMemberOutput struct {
  ID        string
  FirstName string
  LastName  string
  Seniority string
  Status    string
  CreatedAt time.Time
}
```

**Use Case Interface (for testability):**
```go
// application/member/add_member.go

type AddMemberUseCase interface {
  Execute(input AddMemberInput) (*AddMemberOutput, error)
}

type addMemberUseCase struct {
  repo MemberRepository
  // domain services if needed
}

func NewAddMemberUseCase(repo domain.MemberRepository) AddMemberUseCase {
  return &addMemberUseCase{repo: repo}
}

func (uc *addMemberUseCase) Execute(input AddMemberInput) (*AddMemberOutput, error) {
  // Validation (via domain value objects or domain services)
  // Business logic (create aggregate, call domain services)
  // Persistence (call repo.Save())
  // Return output
}
```

**HTTP Handler (thin adapter):**
```go
// server/handler_members.go

func HandlerAddMember(w http.ResponseWriter, r *http.Request, addUC application.AddMemberUseCase) {
  var req AddMemberRequest
  json.NewDecoder(r.Body).Decode(&req)
  
  output, err := addUC.Execute(application.AddMemberInput{
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

**Use Case Unit Test (no database, no HTTP):**
```go
// application/member/add_member_test.go

func TestAddMemberUseCase_Success(t *testing.T) {
  // Inject in-memory repository (from storage/memory/)
  repo := storage.NewInMemoryMemberRepository()
  
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
}
```

### Repository Interface (Domain Package)

**Location:** `domain/member/repository.go`

```go
package domain

// MemberRepository is the abstraction for member persistence.
// All implementations (SQLite, in-memory) must conform to this interface.
type MemberRepository interface {
  Save(member *TeamMember) error
  FindByID(id TeamMemberID) (*TeamMember, error)
  FindActive() ([]*TeamMember, error)
  Deactivate(id TeamMemberID) error
}
```

### Storage Implementations

**In-Memory (for tests):**
```
storage/
  memory/
    member_repository.go  (InMemoryMemberRepository implements domain.MemberRepository)
```

**SQLite (for production):**
```
storage/
  sqlite/
    member_repository.go  (SQLiteTeamMemberRepository implements domain.MemberRepository)
```

Both packages:
- Import only `domain/` (interfaces and aggregates)
- Do not import `application/`, `service/`, or `server/`
- Provide a factory function: `New<Impl>MemberRepository() domain.MemberRepository`

### Migration Path

1. Create `domain/member/repository.go` with interface
2. Create `storage/memory/member_repository.go` with in-memory implementation
3. Create `storage/sqlite/member_repository.go` by extracting/refactoring from `store/member_test.go` and `store/member.go`
4. Create `application/member/` use cases with `Execute()` methods
5. Refactor `server/handler_members.go` to call use cases instead of coordinator
6. Update `main.go` to inject in-memory repository for tests, SQLite repository for production
7. Keep old `service/coordinator/member_api.go` and `store/` untouched during transition; can be archived later

## Testing Strategy

**Use case unit tests:**
- Location: `application/member/<use_case>_test.go`
- Setup: Create in-memory repository, create use case with it, call `Execute()`
- Assertions: Check output, verify repository state
- Speed: Milliseconds; no database, no HTTP, no I/O
- Coverage: success path + error cases (validation, not-found, duplicate)

**Storage integration tests (per implementation):**
- Location: `storage/memory/<impl>_test.go` and `storage/sqlite/<impl>_test.go`
- Setup: Create repository (in-memory or SQLite)
- Assertions: Verify CRUD operations work correctly
- Speed: Milliseconds (memory), milliseconds–seconds (SQLite with `:memory:` DSN)

**HTTP handler tests (unchanged from current):**
- Location: `server/handler_members_test.go`
- Setup: Mock use case with test fixtures
- Assertions: HTTP status, JSON response body
- Speed: Milliseconds; no real use case or database

**Full integration tests (optional, verify end-to-end):**
- Location: `server/integration_test.go`
- Setup: Real SQLite repository + real use cases + HTTP handler
- Assertions: POST → GET → PATCH → DELETE workflow

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Circular imports if not careful | Build failure; architecture violation | Strict package review; interfaces in domain/ only; no application→infrastructure imports |
| More layers = more mocking in tests | Test complexity increase | Use cases only need repository mock (1 interface); simpler than coordinator tests with multiple dependencies |
| Package rename churn | Refactoring noise | Move files once; tests verify new imports work; git history preserved |
| Coordinator not removed immediately | Temporary code duplication | Coordinator left as-is during increment; removal is separate "cleanup" task in future increment |
| Domain layer becomes too large | Bounded context bleeding | Start with member aggregate only; if domain grows, split into separate packages per bounded context |

## Evidence & References

**Clean Architecture (Robert C. Martin):**
- Presentation (HTTP) → Application (Use Cases) → Domain (Business Rules)
- Domain has no dependencies; Application depends on Domain; Presentation depends on both

**Domain-Driven Design (Eric Evans):**
- Bounded contexts (member operations = one bounded context)
- Aggregates (TeamMember = aggregate root)
- Repositories (abstraction for persistence)

**Go Community Practices:**
- Kubernetes: `pkg/` directory with interfaces + implementations
- Docker: clear separation of API, runtime, daemon layers
- CockroachDB: application → storage (abstraction) → engine (implementation)

**Industry Examples:**
- Successful microservices use this pattern to enable CLI, gRPC, and REST from same business logic
- Testing becomes faster and clearer with use case boundaries and injected dependencies

