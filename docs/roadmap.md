# Roadmap

Product direction and sequencing for Team Impact Scorecard. Each entry explains the user outcome, current confidence, and ordering rationale.

> A feature moves to **Done** only when its user outcome is verified and the evidence is linked here.
> Source of truth: if a feature is not in Done with a passing test link, it is not considered shipped.

---

## Done

### Team Member Management (CRUD)
- **Job story:** When I set up the tool or my team changes, I want to add, view, edit, and remove team members with their name and seniority level, so that the scorecard always reflects my current team and shows each member's context.
- **Acceptance scenarios verified:**
  - Add member → appears in list with name and seniority ✓
  - Edit member → updated fields persist ✓
  - Deactivate member → soft-deleted, excluded from list, historical data preserved ✓
  - Persistence → members survive app restart (via SQLite) ✓
- **Evidence:**
  - Store layer: 5 integration tests passing (`store/member_test.go`; add, list, edit, deactivate, soft-delete with audit trail)
  - Service layer: 7 unit tests passing (`service/member/member_test.go`; use cases with mocked store, validation rules)
  - UI layer: Settings screen (Screen F) with add/edit/deactivate dialogs; full CRUD via Fyne
  - Database: idempotent SQLite schema in `store/schema.go` with members table and audit_log
  - Build: successful without errors (harmless macOS linker warnings only)
- **Acceptance criteria:** All 5 met (AC-1: add/view/edit/remove, AC-2: seniority, AC-3: edit, AC-4: deactivate, AC-5: persistence)
- **Key commits:** 
  - Store: `ecf0b6e`, `d070add`, `6749c3a`
  - Service: `74b53fa`
  - UI: `173d3b4`, `f920c60` (layout fix + dialogs)
- **Test command:** `go test -race ./...` → all pass

---

## Done

### Formula Engine — Core Scoring
- **Job story:** When I submit a team member's monthly data, I want the system to compute all normalized scores, dimension scores (DG/DP/DT/DO), TII, completeness, and confidence, so that I have an objective, repeatable basis for my review.
- **Acceptance scenarios verified:**
  - All 10 normalization functions compute correctly per PRD 5.2 ✓
  - Impact-weighted signal computation correct per PRD 5.3 ✓
  - Dimension scores (DG/DP/DT/DO) computed correctly per PRD 5.4 ✓
  - Total Impact Index (TII) computed correctly per PRD 5.5 ✓
  - Completeness percentage computed correctly per PRD 4.3 ✓
  - End-to-end scoring works for typical, strong, and struggling members ✓
- **Evidence:**
  - Engine layer: 53 unit and integration tests passing (`engine/scoring/` with 22 normalization, 4 impact, 12 dimension, 9 TII/completeness/e2e tests)
  - Test coverage: 87.5% of scoring engine statements
  - Race detector: all tests pass with `-race` flag
  - Validation: MoraleN, BillabilityN, CSATN, MarginN, PositiveN, CriticalN, OvertimeN, DeliveryN, MentoringN, EvidenceN, ComputeImpactWeighted, ComputeDimensionGrowth/Project/Team/Org, ComputeTII, ComputeCompleteness
- **Acceptance criteria:** All 6 met (AC-1 through AC-6)
- **Key commits:** 
  - Tidy: 7a704ac (domain types)
  - Feat: aeb62c1 (helpers), 6e01c7b (normalization), a866e78 (impact), 8cc6243 (dimensions), 5773707 (TII/completeness)
- **Test command:** `go test -race ./...` → all pass

---

## Done

### Alert Engine
- **Job story:** When scores or trends cross defined thresholds, I want the system to raise Amber or Red alerts automatically, so that I can intervene before a situation worsens.
- **Acceptance scenarios verified:**
  - All 6 individual alert conditions (Performance Deterioration, Morale Risk, Burnout Risk, Feedback Risk, Customer/Business Risk, Data Quality Risk) classify Red/Amber/none exactly per PRD 7.1 thresholds, including boundary values ✓
  - All 4 team-level alert conditions (Team Morale Drift, Team Delivery Drift, Systemic Burnout, Calibration Risk) classify Red/Amber/none exactly per PRD 7.2 thresholds, including boundary values ✓
  - Each alert carries condition type, severity, and member reference — understandable without re-deriving from raw scores ✓
  - Alert evaluation is deterministic: pure functions, no I/O, no global state; verified by static inspection and 5x repeat-run identical results ✓
- **Evidence:**
  - Engine layer: 39 new unit tests in `engine/domain/alerts_test.go` covering both Amber and Red thresholds at exact PRD boundary values for all 10 alert conditions
  - Full domain suite: 55 tests passing (`go test -race ./engine/domain`)
  - Full repo suite: `go test -race ./...` → all packages pass, no race conditions
  - Build: `go build ./...` succeeds (harmless macOS linker warning only)
- **Acceptance criteria:** All 5 met (AC-1: 6 individual conditions, AC-2: 4 team-level conditions, AC-3: alert detail sufficiency, AC-4: determinism, AC-5: boundary-value test coverage)
- **Known limitation:** team-level alert inputs (% members Red, TII stddev history) are accepted as pre-computed parameters; the aggregation pipeline that derives them from real member data is not yet built (deferred to a future increment).
- **Test command:** `go test -race ./engine/domain -run TestEvaluate` → all pass

## Done

### UI Refactored to Controller Pattern
- **Job story:** When I add a new screen or fix a form logic bug, I want business logic separated from UI widgets, so that logic is unit testable without Fyne overhead and easier to debug and reuse.
- **Acceptance criteria met:**
  - ✅ Controller package created at `ui/controllers/` with base controller interface
  - ✅ MonthlyInputController extracted; owns form state, data loading, validation, save logic
  - ✅ SettingsController extracted; owns member list state and member CRUD logic
  - ✅ Both controllers have zero Fyne imports and are unit testable
  - ✅ Existing MonthlyInputScreen and SettingsScreen refactored to use controllers
  - ✅ All existing tests continue to pass (acceptance tests, integration tests, unit tests)
  - ✅ New controller unit tests added and passing (10 unit tests across both controllers)
- **Evidence:**
  - Controller unit tests: 10 tests passing (`ui/controllers/monthly_input_controller_test.go` — 5 tests; `ui/controllers/settings_controller_test.go` — 5 tests)
  - Screen acceptance tests: 6 tests passing (`ui/screens/monthly_input_acceptance_test.go` — 3 tests including form persistence and reload)
  - Screen integration test: 1 test passing (`ui/screens/monthly_input_it_test.go` — full form workflow test)
  - Total test count: 17 passing (10 unit + 6 acceptance + 1 integration)
  - Build: `go build ./ui/controllers` and `go build ./ui/screens` succeed
  - Race detector: all tests pass with `-race` flag
  - Architecture decision: ADR-20260914 documents rationale and implementation pattern
- **Architecture:**
  - Controllers live in `ui/controllers/` and import only `service/` and `engine/domain` (zero Fyne)
  - Dependency direction: UI Screen → Controller → Service Layer (never reversed)
  - Screens are thin Fyne renderers that delegate state and logic to controllers
  - Controllers are pure Go functions testable without UI framework setup
- **Documentation:**
  - `docs/adr/ADR-20260914-ui-controllers.md` — Decision rationale and pattern
  - `docs/ui.md` — New "Controller Pattern (Screen Architecture)" section with anatomy, testing, and common patterns
- **Key commits:** 
  - Controllers: extraction and unit tests
  - Screens: refactored to use controllers, all tests remain green
- **Test command:** `go test -race ./ui/...` → all pass (17/17)

## In Progress

_Awaiting next increment._

---

## Open Questions

### Reactivating Deactivated Members
- **Question:** When a previously deactivated member re-joins the team, should they be reactivated in place (restoring their history) or added as a new member (clean slate)?
- **Considerations:** Reactivation preserves audit history and avoids duplicate entries; a new member record is simpler but loses historical context and risks orphaned score data.
- **Decision needed before:** Member Detail (Screen C) and any feature that reads historical data per member.

### Duplicate Member Names
- **Question:** Two or more members can legitimately share the same first and last name. How should the UI help users distinguish between them when selecting or reviewing?
- **Constraints:** The store must allow duplicate names (names are not a unique key). Disambiguation must not require renaming real people.
- **Options to explore:** display seniority + join date inline, require a display alias on add, or show member ID as a tie-breaker.
- **Decision needed before:** Monthly Input Workspace (Screen B) and any picker or dropdown that references members by name.

---

## Done

### Monthly Input Workspace — Navigation Shell and Screen Registration
- **Job story:** When I open the monthly input workspace, I want to enter raw signals and impact ratings for each team member with live validation and persistence, so that I can complete the cycle accurately and efficiently.
- **Acceptance scenarios verified:**
  - Monthly Input screen registered on the screen registry and reachable at runtime ✓
  - Valid monthly entry saves and persists for an active member ✓
  - Invalid signal values (out-of-range) are rejected with a descriptive error ✓
  - Duplicate entries for the same member and month are rejected ✓
  - Saved entries survive a reload and are retrievable by member ✓
  - App shell provides sidebar navigation (Overview / Input / Alerts / Review / Settings) and a cycle header ✓
- **Evidence:**
  - Service layer: existing 8 tests passing in `service/monthly/monthly_test.go` (create, update, get, list, invalid-signal, duplicate, inactive-member, compute)
  - UI layer: `ui/screens/monthly_input_test.go` — `TestScreenRegistry_ProvidesMonthlyInputScreen` passing
  - Build: `go build ./...` passes
  - Full suite: `go test -race ./...` passes
- **Acceptance criteria:** All 4 met (AC-1: workflow reachable; AC-2: validation blocks invalid input; AC-3: data persists across restarts; AC-4: full team cycle supported)
- **Note:** Monthly Input form content (signal fields, live preview, member list) is the next increment. This increment delivered the service contract, persistence layer, navigation shell, and screen registration.
- **Test command:** `go test -race ./...` → all pass

---

## In Progress

### Monthly Input Workspace — Form Content (Screen B)
- **Job story:** When I open the monthly input workspace, I want to enter raw signals and impact ratings for each team member with a live formula preview and completeness tracking, so that I can complete the cycle accurately and efficiently.
- **Acceptance criteria:**
  - AC-1: All 10 signal fields shown with PRD-defined ranges; out-of-range values rejected with a descriptive error.
  - AC-2: Selecting a member loads any previously saved entry for the current cycle month.
  - AC-3: Live preview panel updates TII and completeness on valid input without requiring save.
  - AC-4: Saving a valid entry persists it; completeness indicator reflects the updated count.
  - AC-5: "Copy from previous month" populates form from prior month's entry, or shows a clear message when none exists.
- **Evidence:** pending — manual verification of signal fields, live preview updates, and member list behavior

---

## Planned

### Formula Engine — Core Scoring
- **Job story:** When I submit a team member's monthly data, I want the system to compute all normalized scores, dimension scores (DG/DP/DT/DO), TII, completeness, and confidence, so that I have an objective, repeatable basis for my review.
- **Evidence:** pending — unit tests for all normalization and scoring functions against PRD formulas
- **Why first:** everything else (alerts, trends, UI, export) depends on correct formula output. No other increment can be verified without this.

### Trend Calculations (MA3, Delta1, Delta3, Vol3)
- **Job story:** When I view a team member's scorecard, I want to see their trend over time (moving average, deltas, volatility), so that I can distinguish a one-off bad month from a genuine decline.
- **Evidence:** pending — unit tests for all trend functions; integration tests for multi-month history reads
- **Why now:** trends unlock the decision guardrails (promotion, support plan) and are required for several alert conditions.

### Monthly Input Workspace (Screen B)
- **Job story:** When I open the monthly input workspace, I want to enter raw signals and impact ratings for each team member with live formula preview and completeness tracking, so that I can complete the cycle accurately and efficiently.
- **Evidence:** pending — manual verification of all field validations, live preview updates, and copy-previous-month behavior
- **Ordering:** navigation shell complete; depends on formula engine for live preview.

### Overview Dashboard (Screen A)
- **Job story:** When I start my monthly review, I want one screen showing the KPI strip, dimension heatmap, alert table, and action queue for my whole team, so that I know immediately who needs attention.
- **Evidence:** pending — manual verification of heatmap colors, KPI calculations, and alert list behavior
- **Ordering:** depends on monthly input, formula engine, alerts, and trends.

### Member Detail (Screen C)
- **Job story:** When I'm preparing for a 1:1, I want to see a team member's 12-month trend, signal contribution breakdown, evidence log, and action plan, so that my conversation is grounded in data.
- **Evidence:** pending — manual verification of trend chart, contribution table values, evidence timeline
- **Ordering:** depends on formula engine, trends, and persistence layer.

### Alerts Center (Screen D)
- **Job story:** When I have multiple active alerts, I want to filter, triage, assign, and resolve them with reason capture, so that every risk case has an owner and next step.
- **Evidence:** pending — manual verification of filter, resolve/snooze, and action assignment flows
- **Ordering:** depends on alert engine and persistence layer.

### Monthly Review and Calibration (Screen E)
- **Job story:** When I finalize the monthly cycle, I want to see the distribution of scores, review promotion/support guardrail status for each member, and record overrides with mandatory rationale, so that my decisions are fair, auditable, and locked.
- **Evidence:** pending — manual verification of finalization block (completeness < 70%), override capture, and cycle lock behavior
- **Ordering:** depends on all previous screens and the audit trail.

### CSV Export
- **Job story:** When I need to share a review summary, I want to export team overview and member detail data as CSV, so that I can use it in review meetings or archive it.
- **Evidence:** pending — manual verification of CSV structure and completeness for team and member views
- **Ordering:** depends on scoring, trends, and alerts.

### Settings (Screen F)
- **Job story:** When my team's context changes, I want to adjust the billability target, alert sensitivity preset, and reminder dates, so that the tool reflects how my team actually works.
- **Evidence:** pending — manual verification that settings persist and affect alert thresholds correctly
- **Why last:** lowest risk, least user impact. Core functionality must be stable first.

---

## How This List Works

- Features move: Planned → In Progress → Done. Never skip In Progress.
- A feature enters Done only when its user outcome is verified and the evidence link is present.
- Implementation detail belongs in code and phase artifacts, not here.
- If a planned feature is dropped, remove it and record the reason in a commit message or ADR.
