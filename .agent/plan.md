# Plan: Hexagonal Architecture for Member CRUD

## Goal

Refactor member CRUD to follow Ports & Adapters (Hexagonal Architecture): move business logic to infrastructure-independent `core/members/` packages (domain aggregates + ports), implement ports as storage adapters in `storage/` (in-memory and SQLite), replace `MemberAPICoordinator` with explicit use cases, and achieve true separation of concerns with zero core dependencies on infrastructure.

## Branch

`increment/clean-architecture-member-crud`

## Approach

Move member aggregates and business logic into a new `core/members/` package that is completely independent of HTTP, database technology, and UI frameworks. Define ports (interfaces) as contract boundaries in `core/members/repository.go`. Implement ports as adapters in `storage/memory/` and `storage/sqlite/` that depend inward on core only. Create four use case classes (AddMemberUseCase, GetMembersUseCase, EditMemberUseCase, DeactivateMemberUseCase) in `core/members/` that orchestrate domain operations and call the injected repository port. Refactor `server/handler_members.go` to call use cases directly instead of the coordinator. Remove `service/coordinator/member_api.go` entirely. Composition root wiring moves to `main.go` where adapters are instantiated and injected into use cases.

This achieves hexagonal separation: core has zero imports from adapters; adapters depend inward on core only. Business logic is fully testable with in-memory adapters; portable to CLI, gRPC, or other clients without change to core logic.

## Design

### Data Models

**Core types (moved/referenced from `engine/domain/aggregates.go`):**

```
TeamMember (aggregate root)
  id:            TeamMemberID (int64)
  name:          FullName { First string, Last string }
  seniority:     Seniority (enum: Junior, Mid, Senior, Lead, Principal)
  createdAt:     time.Time (immutable)
  deactivatedAt: *time.Time (nil if active)
  methods:
    - ID() TeamMemberID
    - Name() FullName
    - Seniority() Seniority
    - CreatedAt() time.Time
    - IsActive() bool
    - DeactivatedAt() *time.Time
    - Deactivate() error
    - ChangeSeniority(newSeniority Seniority) error
    - UpdateName(firstName, lastName string) error  [new for edit use case]
```

**Port interface (new in `core/members/repository.go`):**

```
MemberRepository interface
  - Save(member *Member) error
  - FindByID(id TeamMemberID) (*Member, error)
  - FindActive() ([]*Member, error)
  - Deactivate(id TeamMemberID) error
```

**Use case input/output structs (new in `core/members/`):**

```
AddMemberInput
  FirstName string
  LastName  string
  Seniority string

AddMemberOutput
  ID        int64
  FirstName string
  LastName  string
  Seniority string
  Status    string (Active/Inactive)
  CreatedAt time.Time

GetMembersOutput
  Members []MemberDTO where MemberDTO has same fields as AddMemberOutput

EditMemberInput
  MemberID  int64
  FirstName string (optional, "" means no change)
  LastName  string (optional)
  Seniority string (optional)

EditMemberOutput
  ID        int64
  FirstName string
  LastName  string
  Seniority string
  Status    string
  CreatedAt time.Time

DeactivateMemberInput
  MemberID int64

DeactivateMemberOutput
  Success bool
```

### Call / Data Flow

**AddMemberUseCase.Execute(input)**
1. Validate input (firstName, lastName non-empty; seniority valid)
2. Create new TeamMember aggregate via domain constructor
3. Call repo.Save(member) to persist
4. On success, generate new UUID for external API representation (v5 with namespace)
5. Return AddMemberOutput

**GetMembersUseCase.Execute()**
1. Call repo.FindActive() to load all active members
2. Map each TeamMember to MemberDTO
3. Convert internal IDs to UUIDs
4. Return GetMembersOutput with member list

**EditMemberUseCase.Execute(input)**
1. Validate input (MemberID must be valid; if seniority provided, must be valid)
2. Call repo.FindByID(input.MemberID) to load existing member
3. If not found, return error (member not found)
4. Apply mutations on aggregate: if firstName provided, call member.UpdateName(); if seniority provided, call member.ChangeSeniority()
5. Call repo.Save(member) to persist changes
6. Generate UUID, return EditMemberOutput

**DeactivateMemberUseCase.Execute(input)**
1. Validate input (MemberID must be positive)
2. Call repo.FindByID(input.MemberID) to load member
3. If not found, return error
4. Call member.Deactivate() to apply domain operation
5. Call repo.Save(member) to persist deactivation
6. Return DeactivateMemberOutput { Success: true }

**HTTP adapter flow (server/handler_members.go)**
1. Handler receives HTTP request (e.g., POST /api/members with JSON body)
2. Parse JSON to AddMemberRequest
3. Create AddMemberInput from request fields
4. Call injected AddMemberUseCase.Execute(input)
5. Map AddMemberOutput to AddMemberResponse (convert int64 ID to UUID string)
6. Write HTTP response (201 Created, JSON body)
7. On error from use case: inspect error type (validation vs. not-found vs. database), map to HTTP status + JSON error response

**UUID mapping (unchanged from current)**
- Internal ID: int64 (database primary key)
- External ID: UUID string (API contract)
- Mapping function: memberIDToUUID(domainID int64) string — uses v5 with namespace `6ba7b810-9dad-11d1-80b4-00c04fd430c8`

### Error / Edge-case Inventory

| Condition | Expected response | Covered by |
|-----------|-------------------|------------|
| firstName or lastName empty in AddMemberUseCase | Return validation error (not persisted) | Subtask 2 — AddMemberUseCase test |
| seniority invalid (not in enum) in AddMemberUseCase | Return validation error | Subtask 2 — AddMemberUseCase test |
| Repo.Save fails (database error) in AddMemberUseCase | Return database error wrapped with user message | Subtask 2 — AddMemberUseCase test |
| MemberID not found in EditMemberUseCase | Return not-found error (no partial updates) | Subtask 4 — EditMemberUseCase test |
| seniority invalid in EditMemberUseCase | Return validation error; aggregate unchanged | Subtask 4 — EditMemberUseCase test |
| Repo.FindActive returns empty list in GetMembersUseCase | Return GetMembersOutput with empty Members slice (not an error) | Subtask 3 — GetMembersUseCase test |
| Member already deactivated in DeactivateMemberUseCase | Aggregate.Deactivate() returns error; wrapped and returned | Subtask 5 — DeactivateMemberUseCase test |
| Repo.FindByID returns nil (member not in DB) in DeactivateMemberUseCase | Return not-found error | Subtask 5 — DeactivateMemberUseCase test |
| HTTP handler receives invalid JSON | Parse error caught; return 400 Bad Request | Subtask 7 — handler test |
| HTTP handler receives unknown field in JSON | JSON decoder ignores extra fields (Go default); no error | Subtask 7 — handler test |
| Multiple simultaneous in-memory adapter writes | No race condition (Go map access is single-threaded within use case call); tests use non-concurrent injection | Subtask 6 — no special handling needed |

### Observability Intent

| Event | Level | Data |
|-------|-------|------|
| `core.member.created` | debug | member_id (internal int64), first_name, last_name, seniority |
| `core.member.deactivated` | debug | member_id, deactivated_at (timestamp) |
| `core.member.updated` | debug | member_id, changed_fields (e.g., "seniority, last_name") |
| `adapter.save_failed` | error | member_id, error_message, adapter_type (sqlite/memory) |
| `adapter.find_failed` | error | member_id, error_message, adapter_type |

Implementation: use `log.Printf()` or structured logging package (e.g., `log/slog`). No external telemetry. Logs written to stderr.

### Architecture Delta

Current architecture (from `docs/architecture.md`):
```
server/ (HTTP handlers)
  ↓ delegates to
coordinator/ (MemberAPICoordinator)
  ↓ calls
service/member/ (MemberService)
  ↓ calls
store/ (SQLiteTeamMemberRepository)
```

New architecture:
```
server/ (HTTP handlers)
  ↓ calls
core/members/ (use cases: AddMemberUseCase, etc.)
  ↓ depends on
core/members/ (MemberRepository port/interface)
  ↑ implemented by
storage/memory/ and storage/sqlite/ (adapters)
```

**Container-level changes:**
- New container: `core/members/` (core business logic; reusable, infrastructure-agnostic)
- New containers: `storage/memory/` and `storage/sqlite/` (adapters; depend inward on core)
- Removed container: `coordinator/member_api.go` (replaced by use cases; coordinator pattern deprecated for member CRUD)
- Unchanged: `engine/domain/` (aggregates, value objects remain; referenced by core/members)
- HTTP handlers remain in `server/`, but now call use cases directly instead of coordinator

**Dependency direction change:**
- Old: `server → coordinator → service → store` (one-way downward chain)
- New: `server → core ← storage` (hexagonal: adapters depend inward on core; core has zero outbound dependencies)

The C4 diagram in `docs/architecture.md` will be updated post-increment to reflect the new hexagonal structure. No functional change to API contract (HTTP endpoints, request/response format, UUID behavior).

## Files

| File | Role | Notes |
|------|------|-------|
| `core/members/member.go` | new | Moved/refactored from `engine/domain/aggregates.go` TeamMember aggregate; add UpdateName method |
| `core/members/repository.go` | new | MemberRepository port interface definition |
| `core/members/add_member.go` | new | AddMemberUseCase |
| `core/members/get_members.go` | new | GetMembersUseCase |
| `core/members/edit_member.go` | new | EditMemberUseCase |
| `core/members/deactivate_member.go` | new | DeactivateMemberUseCase |
| `core/members/add_member_test.go` | new | Unit tests for AddMemberUseCase |
| `core/members/get_members_test.go` | new | Unit tests for GetMembersUseCase |
| `core/members/edit_member_test.go` | new | Unit tests for EditMemberUseCase |
| `core/members/deactivate_member_test.go` | new | Unit tests for DeactivateMemberUseCase |
| `storage/memory/member_repository.go` | new | InMemoryMemberRepository adapter (for tests) |
| `storage/memory/member_repository_test.go` | new | Tests for in-memory adapter |
| `storage/sqlite/member_repository.go` | new | Extract and refactor from `store/member.go` |
| `storage/sqlite/member_repository_test.go` | new | Tests for SQLite adapter (reference existing `store/member_test.go` for test cases) |
| `server/handler_members.go` | modify | Call use cases instead of coordinator; remain thin HTTP adapter |
| `server/handler_members_test.go` | modify | Mock use cases instead of coordinator |
| `main.go` | modify | Wire adapters + use cases at composition root; inject into handlers |
| `service/coordinator/member_api.go` | delete | Coordinator replaced by use cases |
| `engine/domain/aggregates.go` | touch | Reference for TeamMember aggregate (may remain in engine/domain or move to core/members) |
| `store/member.go` | touch | Reference for current SQLite implementation; code extracted to storage/sqlite/ |
| `CONSTITUTION.md` | touch | Reference for dependency rules and architecture boundaries |
| `docs/architecture.md` | touch | Reference; C4 diagram updated post-increment (not this increment) |

## Subtasks

### 1. Create core/members package structure and move/refactor TeamMember aggregate
**type:** [tidy]

**description:** Extract TeamMember aggregate from `engine/domain/aggregates.go` into `core/members/member.go`. Add `UpdateName(firstName, lastName string) error` method to support edit use case. Ensure all existing methods (ID, Name, Seniority, CreatedAt, IsActive, DeactivatedAt, Deactivate, ChangeSeniority) remain. No behavior change; aggregate interface remains the same.

**files:**
- `core/members/member.go` (new)
- `engine/domain/aggregates.go` (touch — reference; consider keeping FullName and Seniority value objects in engine/domain for reuse by monthly entries)

**references:**
- `engine/domain/aggregates.go:17-93` — TeamMember aggregate definition
- `engine/domain/aggregates.go:320-380` — FullName value object
- `engine/domain/aggregates.go:382-420` — Seniority value object

**verification:** 
- `go build ./core/members` succeeds
- `go test ./core/members -run TestTeamMember` passes (existing aggregate tests moved)
- Aggregate methods all present and call signatures match original

**depends on:** —

**acceptance criteria:** AC-1 (core business logic infrastructure-independent)

---

### 2. Define MemberRepository port interface in core/members/repository.go
**type:** [tidy]

**description:** Create `core/members/repository.go` with MemberRepository interface definition. Define port contract: Save, FindByID, FindActive, Deactivate methods. No implementation; interface only. This port defines what adapters must provide to use cases.

**files:**
- `core/members/repository.go` (new)

**references:**
- `store/member.go:23-126` — Current Save, FindByID, FindActive, Deactivate implementation (reference only for contract)

**verification:**
- `go build ./core/members` succeeds
- Interface compiles with correct method signatures
- No implementation code in repository.go file

**depends on:** 1

**acceptance criteria:** AC-2 (port interface defined in core/members)

---

### 3. Create AddMemberUseCase with unit tests
**type:** [behavior]

**description:** Implement AddMemberUseCase in `core/members/add_member.go`. Constructor takes MemberRepository (injected port). Execute(input AddMemberInput) method: validate firstName, lastName (non-empty), seniority (valid enum); create TeamMember aggregate; call repo.Save(); return AddMemberOutput. Unit tests use in-memory repository from storage/memory/.

**files:**
- `core/members/add_member.go` (new)
- `core/members/add_member_test.go` (new)
- `storage/memory/member_repository.go` (new — created as prerequisite for testing)

**references:**
- `core/members/member.go` — TeamMember aggregate
- `core/members/repository.go` — MemberRepository port
- `engine/domain/aggregates.go:382-420` — Seniority enum and Valid() method
- `service/coordinator/member_api.go:46-80` — Current AddMember validation logic (reference for test cases)

**tests:**
- id: add-member-success
  file: `core/members/add_member_test.go`
  name: `TestAddMemberUseCase_Success_CreatesAndPersists`
  state: pending
- id: add-member-missing-first-name
  file: `core/members/add_member_test.go`
  name: `TestAddMemberUseCase_EmptyFirstName_ReturnsValidationError`
  state: pending
- id: add-member-missing-last-name
  file: `core/members/add_member_test.go`
  name: `TestAddMemberUseCase_EmptyLastName_ReturnsValidationError`
  state: pending
- id: add-member-invalid-seniority
  file: `core/members/add_member_test.go`
  name: `TestAddMemberUseCase_InvalidSeniority_ReturnsValidationError`
  state: pending
- id: add-member-repo-error
  file: `core/members/add_member_test.go`
  name: `TestAddMemberUseCase_RepositorySaveError_ReturnsError`
  state: pending

**active_test:** add-member-success

**verification:**
- `go test ./core/members -run TestAddMemberUseCase` passes all 5 test cases
- Use case creates new TeamMember with next available ID
- Output includes correct firstName, lastName, seniority, status (Active), createdAt
- Validation errors caught before repo.Save called (tests verify repo not called on validation error)

**depends on:** 1, 2

**acceptance criteria:** AC-1 (AddMemberUseCase exists and tested with injected in-memory adapter)

---

### 4. Create GetMembersUseCase with unit tests
**type:** [behavior]

**description:** Implement GetMembersUseCase in `core/members/get_members.go`. Constructor takes MemberRepository. Execute() method (no input): call repo.FindActive(), map results to output DTOs, return GetMembersOutput. Handle empty list gracefully (return empty Members slice, not error).

**files:**
- `core/members/get_members.go` (new)
- `core/members/get_members_test.go` (new)

**references:**
- `core/members/repository.go:FindActive()` — Port method
- `service/coordinator/member_api.go:42-44` — Current GetMembers (reference)

**tests:**
- id: get-members-empty
  file: `core/members/get_members_test.go`
  name: `TestGetMembersUseCase_NoMembers_ReturnsEmptyList`
  state: pending
- id: get-members-single
  file: `core/members/get_members_test.go`
  name: `TestGetMembersUseCase_SingleMember_ReturnsOneItem`
  state: pending
- id: get-members-multiple
  file: `core/members/get_members_test.go`
  name: `TestGetMembersUseCase_MultipleMembers_ReturnsAll`
  state: pending
- id: get-members-excludes-inactive
  file: `core/members/get_members_test.go`
  name: `TestGetMembersUseCase_ExcludesInactiveMembers`
  state: pending

**active_test:** get-members-empty

**verification:**
- `go test ./core/members -run TestGetMembersUseCase` passes all 4 test cases
- Empty list returns GetMembersOutput with zero-length Members slice (not nil, not error)
- Multiple members returned in order (if deterministic) with correct fields
- Only active members (IsActive() == true) included in output

**depends on:** 1, 2, 3

**acceptance criteria:** AC-1 (GetMembersUseCase exists and tested)

---

### 5. Create EditMemberUseCase with unit tests
**type:** [behavior]

**description:** Implement EditMemberUseCase in `core/members/edit_member.go`. Constructor takes MemberRepository. Execute(input EditMemberInput) method: validate input (memberID positive, seniority valid if provided); call repo.FindByID(); if not found, return error; apply mutations on aggregate (call UpdateName and/or ChangeSeniority); call repo.Save(); return EditMemberOutput. Partial updates supported (empty string means no change).

**files:**
- `core/members/edit_member.go` (new)
- `core/members/edit_member_test.go` (new)
- `core/members/member.go` (modify — add UpdateName method if not present)

**references:**
- `core/members/member.go:UpdateName` — New method to add
- `core/members/repository.go` — MemberRepository port
- `service/coordinator/member_api.go:105-158` — Current EditMember logic (reference)

**tests:**
- id: edit-member-success
  file: `core/members/edit_member_test.go`
  name: `TestEditMemberUseCase_UpdateAllFields_Success`
  state: pending
- id: edit-member-partial-update
  file: `core/members/edit_member_test.go`
  name: `TestEditMemberUseCase_PartialUpdate_OnlyChangedFields`
  state: pending
- id: edit-member-not-found
  file: `core/members/edit_member_test.go`
  name: `TestEditMemberUseCase_MemberNotFound_ReturnsError`
  state: pending
- id: edit-member-invalid-seniority
  file: `core/members/edit_member_test.go`
  name: `TestEditMemberUseCase_InvalidSeniority_ReturnsValidationError`
  state: pending
- id: edit-member-empty-name
  file: `core/members/edit_member_test.go`
  name: `TestEditMemberUseCase_EmptyFirstNameField_ReturnValidationError`
  state: pending

**active_test:** edit-member-success

**verification:**
- `go test ./core/members -run TestEditMemberUseCase` passes all 5 test cases
- Partial updates work (empty string fields not applied)
- Not-found error returned before repo.Save called
- Validation errors caught before repo.Save called

**depends on:** 1, 2, 4

**acceptance criteria:** AC-1 (EditMemberUseCase exists and tested)

---

### 6. Create DeactivateMemberUseCase with unit tests
**type:** [behavior]

**description:** Implement DeactivateMemberUseCase in `core/members/deactivate_member.go`. Constructor takes MemberRepository. Execute(input DeactivateMemberInput) method: validate input (memberID positive); call repo.FindByID(); if not found, return error; call member.Deactivate() on aggregate; call repo.Save(); return DeactivateMemberOutput { Success: true }. Deactivation is idempotent only if handled correctly (aggregate.Deactivate() returns error if already deactivated).

**files:**
- `core/members/deactivate_member.go` (new)
- `core/members/deactivate_member_test.go` (new)

**references:**
- `core/members/member.go:Deactivate()` — Aggregate method
- `core/members/repository.go` — MemberRepository port
- `service/coordinator/member_api.go:175-206` — Current DeactivateMember logic (reference)

**tests:**
- id: deactivate-member-success
  file: `core/members/deactivate_member_test.go`
  name: `TestDeactivateMemberUseCase_ActiveMember_Success`
  state: pending
- id: deactivate-member-not-found
  file: `core/members/deactivate_member_test.go`
  name: `TestDeactivateMemberUseCase_MemberNotFound_ReturnsError`
  state: pending
- id: deactivate-member-already-inactive
  file: `core/members/deactivate_member_test.go`
  name: `TestDeactivateMemberUseCase_AlreadyInactive_ReturnsError`
  state: pending

**active_test:** deactivate-member-success

**verification:**
- `go test ./core/members -run TestDeactivateMemberUseCase` passes all 3 test cases
- Active member successfully deactivated; deactivatedAt timestamp set
- Not-found error returned before aggregate.Deactivate() called
- Already-inactive error returned from aggregate.Deactivate() and wrapped

**depends on:** 1, 2, 5

**acceptance criteria:** AC-1 (DeactivateMemberUseCase exists and tested)

---

### 7. Create in-memory storage adapter (storage/memory/member_repository.go)
**type:** [tidy]

**description:** Implement InMemoryMemberRepository in `storage/memory/member_repository.go`. Implements MemberRepository port using Go map to store members in memory. Used by all use case unit tests. No persistence; data lost on process exit. Minimal implementation (no fancy data structures; map is sufficient for test use).

**files:**
- `storage/memory/member_repository.go` (new)
- `storage/memory/member_repository_test.go` (new)

**references:**
- `core/members/repository.go` — MemberRepository interface to implement
- `core/members/member.go` — TeamMember type

**verification:**
- `go build ./storage/memory` succeeds
- Implements all MemberRepository methods
- `go test ./storage/memory -run InMemory` passes (CRUD tests)

**depends on:** 1, 2

**acceptance criteria:** AC-3 (in-memory adapter implemented and available for test injection)

---

### 8. Create SQLite storage adapter (storage/sqlite/member_repository.go)
**type:** [tidy]

**description:** Extract and refactor SQLite implementation from `store/member.go` into `storage/sqlite/member_repository.go`. Implements MemberRepository port using SQLite. Minimal refactoring to match port interface; preserve all existing logic (audit log, error handling, ID generation). Ensure no behavior change from current store/member.go implementation.

**files:**
- `storage/sqlite/member_repository.go` (new — extracted from store/member.go)
- `storage/sqlite/member_repository_test.go` (new — reference store/member_test.go)

**references:**
- `core/members/repository.go` — MemberRepository interface
- `store/member.go:23-333` — Current SQLite implementation (extract and adapt)
- `store/member_test.go` — Existing test cases (reference for new tests)

**verification:**
- `go build ./storage/sqlite` succeeds
- All methods implement MemberRepository interface exactly
- `go test ./storage/sqlite -run SQLite` passes (test same cases as store/member_test.go)
- Behavior identical to current store/member.go (no test failures)

**depends on:** 1, 2

**acceptance criteria:** AC-3 (SQLite adapter implemented)

---

### 9. Create use case unit test suite (all four use cases, injecting in-memory adapter)
**type:** [behavior]

**description:** Already covered in subtasks 3–6 (AddMemberUseCase, GetMembersUseCase, EditMemberUseCase, DeactivateMemberUseCase). Each has comprehensive unit tests that inject InMemoryMemberRepository. This subtask confirms all unit tests pass together with `-race` flag and clean imports.

**files:**
- `core/members/add_member_test.go` (reference — created in subtask 3)
- `core/members/get_members_test.go` (reference — created in subtask 4)
- `core/members/edit_member_test.go` (reference — created in subtask 5)
- `core/members/deactivate_member_test.go` (reference — created in subtask 6)

**references:**
- `CONSTITUTION.md#testing-strategy` — Testing practices

**tests:** (all tests from subtasks 3–6 run together)

**verification:**
- `go test -race ./core/members` passes all 16+ test cases
- No circular imports detected
- No infrastructure dependencies in core packages (grep confirms: no imports from storage/, server/, service/, engine/scoring/)

**depends on:** 3, 4, 5, 6, 7

**acceptance criteria:** AC-1 (use case unit tests pass with in-memory adapter injection)

---

### 10. Refactor HTTP handlers to call use cases directly
**type:** [behavior]

**description:** Update `server/handler_members.go` to call AddMemberUseCase, GetMembersUseCase, EditMemberUseCase, DeactivateMemberUseCase instead of coordinator. Handlers remain thin adapters: parse HTTP request → create use case input → call Execute() → map output to JSON response → write HTTP status. Error handling: inspect use case error type, map to HTTP status (validation 400, not-found 404, database 500).

**files:**
- `server/handler_members.go` (modify)

**references:**
- `server/handler_members.go:36-276` — Current handler implementations (reference)
- `core/members/add_member.go` — AddMemberUseCase interface
- `core/members/get_members.go` — GetMembersUseCase interface
- `core/members/edit_member.go` — EditMemberUseCase interface
- `core/members/deactivate_member.go` — DeactivateMemberUseCase interface

**tests:**
- id: handler-add-member-success
  file: `server/handler_members_test.go`
  name: `TestHandlerAddMember_ValidRequest_Success`
  state: pending
- id: handler-add-member-invalid-json
  file: `server/handler_members_test.go`
  name: `TestHandlerAddMember_InvalidJSON_ReturnsBadRequest`
  state: pending
- id: handler-get-members-success
  file: `server/handler_members_test.go`
  name: `TestHandlerGetMembers_Success`
  state: pending
- id: handler-edit-member-success
  file: `server/handler_members_test.go`
  name: `TestHandlerEditMember_ValidRequest_Success`
  state: pending
- id: handler-deactivate-member-success
  file: `server/handler_members_test.go`
  name: `TestHandlerDeactivateMember_ValidRequest_Success`
  state: pending

**active_test:** handler-add-member-success

**verification:**
- `go test ./server -run TestHandler` passes all 5+ handler tests
- Handlers no longer import coordinator
- HTTP status codes and response formats match current behavior (no API contract change)

**depends on:** 1, 2, 3, 4, 5, 6

**acceptance criteria:** AC-4 (HTTP handlers call use cases; coordinator removed)

---

### 11. Refactor server/handler_members_test.go to mock use cases
**type:** [behavior]

**description:** Update `server/handler_members_test.go` to mock use cases instead of coordinator. Tests remain focused on HTTP parsing and response serialization. Mock use cases using simple test doubles (not a full mocking library). Verify handler tests no longer touch database or call real coordinator.

**files:**
- `server/handler_members_test.go` (modify)

**references:**
- `server/handler_members_test.go` — Existing test structure (reference)
- `core/members/add_member.go`, etc. — Use case interfaces to mock

**verification:**
- `go test ./server -run TestHandler` passes all tests
- Test file imports no coordinator, no store, no database packages
- Mocks are simple test doubles (implement use case interface)

**depends on:** 10

**acceptance criteria:** AC-4, AC-5 (handler tests use injected mock use cases; no infrastructure access)

---

### 12. Update main.go to wire adapters and use cases at composition root
**type:** [tidy]

**description:** Modify `main.go` to instantiate storage adapters (InMemoryMemberRepository for tests or SQLiteTeamMemberRepository for production), create use case instances with injected adapters, and pass use cases to HTTP server. Composition root is the single place where dependency injection happens. Preserve current flow: initialize database, create repositories, create use cases, start server.

**files:**
- `main.go` (modify)

**references:**
- `main.go:18-51` — Current wiring (reference)
- `storage/sqlite/member_repository.go` — SQLite adapter to instantiate
- `core/members/add_member.go`, etc. — Use cases to instantiate

**verification:**
- `go build ./...` succeeds
- No import cycles detected
- main.go imports storage (for adapter), core/members (for use cases), server (for handlers)
- Application starts and HTTP server runs

**depends on:** 8, 10

**acceptance criteria:** AC-5 (composition root wires adapters; dependency graph is acyclic)

---

### 13. Delete service/coordinator/member_api.go entirely
**type:** [tidy]

**description:** Remove `service/coordinator/member_api.go`. Coordinator has been replaced by use cases. No other file should import member_api after step 10; confirm with grep before deletion.

**files:**
- `service/coordinator/member_api.go` (delete)

**references:**
- `service/coordinator/member_api.go` — File to delete

**verification:**
- `grep -r "coordinator.MemberAPICoordinator" ./` returns zero results (no remaining imports)
- `go build ./...` succeeds after deletion
- No test failures

**depends on:** 10, 12

**acceptance criteria:** AC-4 (coordinator removed entirely)

---

### 14. Run full test suite and verify no circular imports
**type:** [behavior]

**description:** Execute `go test -race ./...` to confirm all 21+ tests pass. Verify no circular imports with `go build ./...`. Check that core/members/ has zero imports from storage/, server/, service/, or ui/ (use grep or go list -d). Confirm dependency flow is hexagonal: server → core ← storage.

**files:**
- All files (indirect verification)

**references:**
- `CONSTITUTION.md#architecture-boundaries` — Dependency rules
- `.agent/increment.md#acceptance-criteria` — Binary criteria

**tests:**
- id: full-test-suite
  file: (implicit)
  name: `All 21+ tests pass with -race`
  state: pending
- id: no-circular-imports
  file: (implicit)
  name: `go build ./... succeeds; no import cycles`
  state: pending
- id: core-independence
  file: (implicit)
  name: `core/members/ has zero outbound infrastructure dependencies`
  state: pending

**active_test:** full-test-suite

**verification:**
- `go test -race ./...` output: all tests PASS
- `go build ./...` output: no errors or warnings
- `go list -d ./core/members/...` output: imports are only stdlib and engine/domain
- `grep -r "import.*storage/" ./core/members/` returns zero results
- `grep -r "import.*server/" ./core/members/` returns zero results
- `grep -r "import.*service/" ./core/members/` returns zero results
- `grep -r "import.*ui/" ./core/members/` returns zero results

**depends on:** 1–13 (all prior subtasks)

**acceptance criteria:** AC-1, AC-2, AC-3, AC-4, AC-5 (all binary criteria met; no circular imports)

---

## Context Map

- `CONSTITUTION.md#architecture-boundaries` — Dependency rules: `server → service → store → engine`; **this increment changes to `server → core ← storage` (hexagonal)**; must verify no violation
- `CONSTITUTION.md#coordinator-pattern` — Coordinator documentation; coordinator for members is being removed; this increment closes the gap with use cases
- `docs/architecture.md#c4-level-2` — Current container view; post-increment, update to show hexagonal structure (not part of this increment)
- `docs/adr/ADR-20260915-uuid-member-ids.md` — UUID v5 namespace and mapping rules; unchanged in this increment; use cases generate UUIDs as before
- `engine/domain/aggregates.go:1-93` — TeamMember aggregate (reference; will be refactored into core/members/)
- `service/coordinator/member_api.go:1-207` — Current coordinator (reference for test cases and error handling; to be deleted)
- `store/member.go:1-333` — Current SQLite implementation (reference for adapter extraction)
- `server/handler_members.go:1-276` — Current handlers (to be refactored)

## Acceptance Scenarios

### AT-1: Add member via HTTP API
- criterion: AC-1, AC-4
- user action: POST /api/members with { firstName: "Alice", lastName: "Smith", seniority: "Senior" }
- precondition: server running, database initialized
- expected outcome: HTTP 201 Created, response includes id (UUID), firstName, lastName, seniority, status (Active), createdAt
- evidence: `server/handler_members_test.go::TestHandlerAddMember_Success` and `server/integration_test.go::TestMembersAPIIntegration_CRUD`
- gate: advisory
- state: planned

### AT-2: Edit member via HTTP API
- criterion: AC-1, AC-4
- user action: PATCH /api/members/{id} with { lastName: "Jones", seniority: "Lead" }
- precondition: member exists with id
- expected outcome: HTTP 200 OK, response includes updated lastName and seniority; createdAt unchanged
- evidence: `server/handler_members_test.go::TestHandlerEditMember_Success`
- gate: advisory
- state: planned

### AT-3: Deactivate member via HTTP API
- criterion: AC-1, AC-4
- user action: DELETE /api/members/{id}
- precondition: member is active
- expected outcome: HTTP 200 OK; member marked inactive; subsequent GET /api/members excludes deactivated member
- evidence: `server/handler_members_test.go::TestHandlerDeactivateMember_Success` and integration test
- gate: advisory
- state: planned

### AT-4: Validation error on missing firstName
- criterion: AC-1
- user action: POST /api/members with { firstName: "", lastName: "Smith", seniority: "Senior" }
- precondition: server running
- expected outcome: HTTP 400 Bad Request, JSON error response with error message
- evidence: `server/handler_members_test.go::TestHandlerAddMember_ValidationError`
- gate: advisory
- state: planned

## Risks

1. **Circular imports during refactoring** — If any file in storage/ or server/ accidentally imports from core/members (or vice versa), build will fail. Mitigation: strict code review before approval; run `go build ./...` frequently during implementation.

2. **Partial migration leaves old coordinator imported** — If handler refactoring incomplete, handlers still call coordinator; old and new code paths coexist. Mitigation: subtask 13 (delete coordinator) cannot complete until all handlers refactored and tests pass.

3. **UUID mapping inconsistency** — If new use cases generate UUIDs differently than current coordinator, API contract changes. Mitigation: reuse existing memberIDToUUID() function from current handler_members.go; no new ID generation logic.

4. **In-memory adapter too simplistic for concurrent tests** — Go map access is not thread-safe. If tests run in parallel and share in-memory adapter, race condition occurs. Mitigation: each test creates its own InMemoryMemberRepository instance; tests run sequentially (Go default); `-race` flag used to detect any unexpected concurrency.

5. **Test coverage gap: storage adapter edge cases** — Current store/member_test.go has specific test cases (audit log, error handling) that must be preserved in storage/sqlite/member_repository_test.go. If not transferred, behavior regression. Mitigation: reference store/member_test.go explicitly in subtask 8; verify same test cases run on new adapter.

6. **Performance regression** — Hexagonal architecture adds an extra abstraction layer (port interface). If adapters are slow or allocate heavily, HTTP latency could exceed CONSTITUTION.md envelope (< 200ms for screen render). Mitigation: in-memory adapter is trivially fast; SQLite adapter unchanged from store/member.go; performance validated manually during testing.

## Planning Decisions

| Decision | Chosen | Rejected | Reason |
|----------|--------|----------|--------|
| **Where to put core aggregates and ports** | `core/members/` (co-located) | Separate `domain/` and `application/` directories | Simpler navigation for small codebase (2 bounded contexts); easier to move entire feature to separate module later if needed; "core" signals infrastructure-independent |
| **Repository port location** | `core/members/repository.go` (same package as use cases) | Separate `core/members/interfaces/` package | Less package fragmentation; port and use cases live together; adapters clearly depend on this one package |
| **Adapter organization** | `storage/memory/` and `storage/sqlite/` (separate subdirs per tech) | Flat `storage/` with `InMemoryMemberRepository` and `SQLiteMemberRepository` in same file | Scaling: when monthly entry adapters added, subdirectories keep storage/ organized; easy to see which tech-specific adapters exist |
| **Coordinator deletion timing** | Delete immediately (subtask 13) after handlers refactored | Leave for cleanup increment | Cleaner cutover; avoid maintaining two parallel paths; no risk of old code being called once handlers refactored and tests pass |
| **TeamMember aggregate location** | Move to `core/members/member.go` | Keep in `engine/domain/` and import from there | core/members/ needs to be self-contained; future: if multiple bounded contexts use TeamMember, refactor to shared package (not yet needed) |
| **Test doubles for handlers** | Simple hand-written mocks (test doubles) | Use `github.com/golang/mock` or other mocking library | Minimal dependencies; clearer test intent; no code generation; appropriate for 4–5 mock methods |
| **In-memory adapter persistence** | Stateless per test (new instance per test) | Global map (shared across tests) | Test isolation; no hidden dependencies between tests; standard Go testing practice |
| **UUID generation in use cases** | Call memberIDToUUID() from current handler code (unchanged) | Move to use cases and make configurable | No change to API contract; reuse existing logic; ports don't need to know about UUID scheme |
| **Error types for validation vs. not-found** | Use Go built-in `error` interface; use error message prefix or custom type wrapper | Create domain error enum or error type hierarchy | Simpler for v1; can refactor to richer error types later if needed; tests can assert on error message or error type with type assertion |
| **Async/concurrency in use cases** | Sequential (no goroutines in use case logic) | Parallel processing for batch operations | Not needed for CRUD operations; keep simple; add concurrency only if performance testing shows need |

---

## Implementation Readiness

All planning decisions recorded. File-level scope clear (14 new files, 4 modified, 4 touched, 1 deleted). Dependency graph acyclic. Test strategy defined per subtask. Ready for implementation using tidy, tdd-red, tdd-green, refactor skills.

