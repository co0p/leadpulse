# Testing

Testing practices for Team Impact Scorecard. The formula engine is the highest-risk component: incorrect scoring or alert thresholds produce bad people decisions. The test strategy is designed to give near-complete confidence in engine correctness with fast feedback, while accepting lower automated coverage for the Fyne UI layer.

---

## Testing Approach and Rationale

The system has three layers with different risk profiles:

- **`engine/`** — pure functions, deterministic, no I/O. The formulas in the PRD are the specification. Any divergence is a bug. This layer must have comprehensive unit tests because errors here silently corrupt people decisions. Tests run in milliseconds, have no dependencies, and are the primary correctness gate.
- **`store/`** — SQLite persistence. Risk is data loss, constraint violations, and incorrect reads of historical data (which affect trend calculations). Integration tests against an in-memory SQLite database cover the meaningful paths.
- **`ui/`** — Fyne widgets and views. Fyne does not have a reliable headless test driver. Manual verification is the pragmatic choice for v1. UI tests would be fragile, slow to maintain, and provide less signal than running the app.

---

## Choosing Test Depth

Write a unit test when:
- A formula, normalization, threshold, or computation rule from the PRD is being implemented or changed.
- A new alert condition is added.
- A trend calculation is introduced (MA3, Delta1, Delta3, Vol3, dimension deltas).
- A business rule is encoded (e.g., finalization block when completeness < 70%).

Write an integration test when:
- A new store function is added (read, write, update, delete).
- A schema migration is introduced.
- A query aggregates or joins data across months.

No new test needed when:
- Renaming a variable or extracting a helper with identical behavior.
- Changing UI layout, colors, or widget positioning.
- Adding a log statement.

---

## Test Design Conventions

- **Engine tests** live in `engine/<file>_test.go`, alongside the function under test.
- Test function names follow `Test<FunctionName>_<scenario>`. Example: `TestNormalizeMorale_midRange`, `TestAlertBurnout_redThreshold`.
- Each test case states its input, expected output, and the PRD section it covers in a comment.
- Table-driven tests are preferred for formula coverage — one table per function, rows per boundary condition.
- Integration tests use `":memory:"` as the SQLite DSN. No test writes to disk.
- Tests must not share mutable state. Each test case sets up its own fixtures.

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

- `go test -race ./engine/... ./store/...` passes with no failures and no race conditions.
- `go build ./...` succeeds.
- For any change touching a formula, alert, or trend calculation: the relevant test case(s) exist and are named after the PRD section they verify.
- For any new store function: at least one integration test covers the happy path and one covers the relevant constraint or error case.
- Manual smoke-test of the affected screen(s) on the developer's machine.

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
