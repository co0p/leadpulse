# Testing

Testing practices for Team Impact Scorecard. The service layer is the primary testability boundary: use cases can be tested without Fyne by passing a mock store. The formula engine is the highest-risk component; incorrect scoring or alert thresholds produce bad people decisions. The test strategy gives near-complete confidence in engine and service correctness with fast feedback, while accepting lower automated coverage for the UI layer.

---

## Testing Approach and Rationale

The system has four layers with different risk profiles:

- **`engine/`** — pure functions, deterministic, no I/O. The formulas in the PRD are the specification. Any divergence is a bug. This layer must have comprehensive unit tests because errors here silently corrupt people decisions. Tests run in milliseconds, have no dependencies, and are the primary correctness gate.
- **`service/`** — use cases (transactions). Each use case fetches aggregates via repository, calls domain services and `engine`, and persists results via repository. The service layer is the API contract. Unit tests inject in-memory repository implementations (no database needed) and real `engine` functions. Service tests verify that use cases call repository and domain service methods in the right order with the right data. Service tests prove that UI is replaceable: any presentation layer (Fyne, CLI, web) can call the same `service/` interfaces and behave identically.
- **`store/`** — repository implementations (SQLiteTeamMemberRepository, SQLiteMonthlyEntryRepository) that persist and retrieve aggregates from SQLite. Risk is data loss, constraint violations, incorrect reads of historical data (which affect trend calculations), and audit trail correctness. Integration tests against an in-memory SQLite database (`":memory:"` DSN) cover all CRUD paths and aggregate reconstruction.
- **`ui/`** — Fyne widgets and views. Fyne does not have reliable headless test support. Manual verification is pragmatic for v1. Acceptance criteria in `.agent/increment.md` are manual user stories, not automated suites.

---

## Choosing Test Depth

**Engine unit tests:** Write when:
- A formula, normalization, threshold, or computation rule from the PRD is being implemented or changed.
- A new alert condition is added.
- A trend calculation is introduced (MA3, Delta1, Delta3, Vol3, dimension deltas).
- A business rule is encoded (e.g., guardrail logic).

**Domain service unit tests (engine/domain/):** Write when:
- A new domain service is implemented (e.g., `ValidationService`, `ScoringService`).
- A domain service method enforces a new business rule or cross-aggregate constraint.
- A domain service's rule logic changes.

**Service unit tests:** Write when:
- A new use case is implemented (e.g., `AddMember`, `SubmitMonthlyEntry`, `GenerateAlerts`).
- A use case's logic changes (order of operations, new calls to `store` or `engine`, new validations).
- A use case's error handling is added or changed.
- A use case's integration with domain services changes.

**Store integration tests:** Write when:
- A new store function is added (read, write, update, delete).
- A schema migration is introduced.
- A query aggregates or joins data across months or members.
- Soft-delete, audit trail, or constraint logic is added.

**UI manual verification:** Required when:
- A screen is added or its layout changes.
- A new interaction pattern is introduced (e.g., new form, new table, new dialog).
- See acceptance criteria in `.agent/increment.md` for manual test scenarios.

**No new test needed when:**
- Renaming a variable or extracting a helper with identical behavior.
- Changing UI colors, fonts, or widget positioning (not behavior).
- Adding a log statement or comment.

---

## Test Design Conventions

**Engine tests:**
- Live in `engine/<file>_test.go`, alongside the function under test.
- Function names: `Test<FunctionName>_<scenario>`. Example: `TestNormalizeMorale_midRange`, `TestAlertBurnout_redThreshold`.
- Each case states input, expected output, and the PRD section it covers in a comment.
- Table-driven tests preferred for formula coverage — one table per function, rows per boundary condition.

**Service tests:**
- Live in `service/<domain>/<file>_test.go`, alongside the use case.
- Function names: `Test<UseCaseName>_<scenario>`. Example: `TestAddMember_createsAndReturns`, `TestCreateEntry_validatesActiveMember`.
- Use in-memory repository implementations (InMemoryTeamMemberRepository, InMemoryMonthlyEntryRepository) — no database dependency. Real `engine` functions and domain services (no mocks).
- Each test: create an in-memory repository, inject it into the service, set up fixtures via repository, call the use case, verify repository state and return values.
- Verify error cases: invalid input, repository errors, domain service rule violations (e.g., inactive member, duplicate entry), engine validation failures.
- Example: for CreateMonthlyEntry, test that calling with an inactive member fails before calling repository.Save().

**Domain service tests (engine/domain/):**
- Live in `engine/domain/<service>_test.go`, alongside the domain service implementation.
- Function names: `Test<ServiceName>_<scenario>`. Example: `TestValidationService_rejectsDuplicate`, `TestScoringService_validatesRange`.
- Use in-memory repository implementations exclusively (InMemoryTeamMemberRepository, InMemoryMonthlyEntryRepository) — no database or store package imports.
- Each test: create in-memory repositories, inject into domain service constructor, set up test data via repository, call domain service methods, verify rule enforcement without database access.
- Verify business rule enforcement: uniqueness constraints, cross-aggregate invariants, state consistency, value range validation.
- Domain services never require database access during testing. If a domain service test needs a database, the test design is wrong; refactor the dependency injection.
- **Canonical pattern:**
  ```go
  // Create in-memory repositories
  memberRepo := domain.NewInMemoryTeamMemberRepository()
  entryRepo := domain.NewInMemoryMonthlyEntryRepository()
  
  // Create domain service with repository dependencies
  service := domain.NewValidationService(memberRepo, entryRepo)
  
  // Set up test fixtures via repository (no SQL)
  member, _ := domain.NewTeamMember(1, name, seniority)
  memberRepo.Save(member)
  
  // Call domain service and verify rule enforcement
  err := service.ValidateTeamMemberUniqueness(member)
  // Assert error as expected
  ```

**Store integration tests:**
- Live in `store/<file>_test.go`, alongside the function under test.
- Use `":memory:"` as the SQLite DSN. No test writes to disk.
- Each test: create fixtures via store methods, verify state in the database, verify audit trail entries.
- Cover happy path and constraint violations (foreign keys, NOT NULL, unique constraints).
- **Standard helper: `setupTestDB()`** — initializes an in-memory SQLite database with the full schema applied. All store tests must call `setupTestDB()` before creating fixtures. This ensures test isolation and consistent schema versioning across all tests. See `store/member_test.go` for the reference implementation.

**All tests:**
- Must not share mutable state. Each test case sets up its own fixtures.
- Test must be independent: can run in any order, can run in parallel with `-race`.

---

## Running the Checks

**Prerequisites:**
```bash
go version  # requires Go 1.21+
```

**Fast local feedback (engine + store only):**
```bash
go test ./engine/... ./store/...
```

**Full test suite:**
```bash
go test ./...
```

**With race detector (run before any promotion):**
```bash
go test -race ./...
```

**Verbose output (useful when debugging a specific failure):**
```bash
go test -v -run TestAlertBurnout ./engine/...
```

**Build check (catches compile errors across all packages including ui/):**
```bash
go build ./...
```

Interpret failures: a failing engine test is a blocking defect. A failing store test is a blocking defect. A build failure is a blocking defect. All must be resolved before the change is considered complete.

---

## Evidence Required Before Merge

- `go test -race ./engine/... ./service/... ./store/...` passes with no failures and no race conditions.
- `go build ./...` succeeds.
- For any change touching a formula, alert, or trend calculation: the relevant engine test case(s) exist and are named after the PRD section they verify.
- For any new use case: at least one service test covers the happy path and one covers relevant error cases (invalid input, store failures).
- For any new store function: at least one integration test covers the happy path and one covers the relevant constraint or error case.
- For any UI screen or interaction: manual acceptance-test scenarios pass on the developer's machine (see `.agent/increment.md` for the expected user journeys).
- Service tests prove the use case contract is stable, independent of the UI implementation.

---

## Automation and Feedback Loops

- **Local:** `go test ./engine/... ./store/...` — run on every save or before every commit. Completes in seconds.
- **CI (GitHub Actions):** runs `go test -race ./...` and `go build ./...` on every push and pull request. Targets: `ubuntu-latest`, `macos-latest`, `windows-latest`. A failing CI run blocks merge.
- **Release:** same CI gate. No additional test suite for releases in v1.

Manual UI verification is the developer's responsibility before opening a pull request.

---

## Known Risks and Gaps

- **UI layer has no automated tests.** Regressions in Fyne widget behavior require manual detection. Acceptable in v1 given Fyne's limited test tooling. Revisit if Fyne's `test` package matures sufficiently.
- **Export correctness is not automatically tested.** CSV output is verified manually. A regression here would surface quickly in real use.
- **Performance is not automatically measured.** The 100ms / 200ms targets from the constitution are checked manually during development. No benchmark suite exists yet.

---

## Maintenance Guidance

- Engine tests are the most important tests in the project. Do not delete or weaken them to make a change easier.
- If a test is flaky, fix it immediately — do not mark it as skipped.
- When a PRD formula changes, update the relevant test table rows before updating the implementation.
- This document is updated when the testing approach changes, not when individual tests are added.
