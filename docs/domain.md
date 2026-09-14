# Domain Glossary

Shared vocabulary for Team Impact Scorecard. Use these terms consistently in code, tests, and documentation. When a concept is introduced in code, it must appear here first or simultaneously.

---

## Design Patterns and Architecture

**Value Object**
An immutable, identity-less object that represents a single domain concept and enforces its own constraints. Examples: `Seniority` (Junior|Mid|Senior|Principal), `FullName` (first + last), `MoraleScore` (0–5). Value objects are created via constructor functions that validate inputs; invalid inputs are rejected immediately. Once created, a value object cannot be modified.

**Aggregate**
A cluster of entities and value objects managed together as a single unit of persistence. The aggregate has one root entity (the Aggregate Root) that controls all access to the aggregate's internal state. All invariants and business rules are enforced within the aggregate. Examples: `TeamMember` (aggregate root managing name, seniority, and deactivation state), `MonthlyEntry` (aggregate root managing raw signals, impact ratings, and computed scores).

**Aggregate Root**
The entry point to an aggregate. The aggregate root has an identity (ID) and controls all interactions with the aggregate. Only the aggregate root is persisted and retrieved directly; other entities within the aggregate are accessed through the root.

**Repository**
A collection-like abstraction that persists and retrieves aggregates from storage. The repository presents a domain-oriented interface (Save, FindByID, FindByMember, Delete) without exposing storage details (SQL, database schema, transaction mechanics). Domain logic depends on repository interfaces, not on concrete storage implementations.

**Domain Service**
A stateless service that enforces business rules spanning multiple aggregates. Domain services operate on aggregates and repositories, but the logic belongs to the domain, not to the application layer. Example: `MonthlyEntryService` enforces that only active members can have entries and that one entry per member per month is allowed.

### Domain Service Implementation Pattern

Every domain service follows a consistent structure for consistency and testability:

```go
// Struct: holds repository dependencies only (no application state)
type DomainService struct {
    memberRepo domain.TeamMemberRepository
    entryRepo  domain.MonthlyEntryRepository
}

// Constructor: factory function for initialization
func NewDomainService(memberRepo, entryRepo) *DomainService {
    return &DomainService{
        memberRepo: memberRepo,
        entryRepo:  entryRepo,
    }
}

// Methods: public methods enforce business rules and delegate to repositories
func (s *DomainService) BusinessRule(aggregate *Aggregate) error {
    // Validate cross-aggregate invariants
    // Call repository methods
    // Return error if rule violated
}
```

**Rules:**
- Domain services accept repository *interfaces*, not concrete storage (`*sql.DB`, store types).
- All methods are testable with in-memory repository implementations.
- No circular dependencies between domain services.
- Error handling preserves error chains using `fmt.Errorf("%w", err)`.

**Testing:** Domain service tests create in-memory repositories, inject them into the service constructor, and verify rule enforcement without database access.

**Examples:** `ScoringService` (computes impact ratings), `TrendService` (calculates moving averages), `ValidationService` (enforces uniqueness and consistency).

---

## Core Concepts

**Team Lead**
The sole user of v1. A person responsible for one team of direct reports. Enters monthly data, reviews scores and alerts, manages the monthly cycle, and exports summaries.

**Team Member**
A direct report of the Team Lead. Represented as an aggregate root with identity (ID), name, and seniority level. All state transitions (e.g., seniority changes, deactivation) are controlled by aggregate methods. Invariants: ID is immutable; once deactivated, cannot be reactivated; seniority must be one of Junior, Mid, Senior, Principal.

**Monthly Entry**
An aggregate root containing all raw signals, impact ratings, and computed scores for one Team Member for one calendar month. Identified by `(member_id, month)`. Contains a sub-aggregate of raw signals and validates all signal constraints and score ranges. Invariants: one entry per member per month; computed scores are validated to be 0–100; only active members can have entries (enforced by MonthlyEntryService).

**Cycle**
The structured monthly operating rhythm. A cycle covers one calendar month. It progresses through phases: Prepare → Data Entry → Auto Compute → Triage → Monthly Review → 1:1 Execution → Follow-through.

**Finalized Cycle**
A cycle that has been locked by the Team Lead after all members have sufficient completeness and all Red alerts have been annotated. Finalization is blocked if any member's completeness is below 70%.

---

## Entry Status: Draft vs Done

- Draft: a monthly entry that exists in the store but has at least one nil (unset) signal. Draft entries are visible in the member list with a "Draft" badge. Draft entries are allowed and represent partially-complete input.
- Done: a monthly entry where all 10 signals are non-nil. Done entries show a "Done" badge and are considered complete for finalization checks.

Notes:
- The UI Save button is disabled until completeness ≥ 70% to encourage sufficient data entry before persisting.
- Completeness is computed as the percentage of required filled fields (FilledRequiredFields / 11 × 100) and is used by confidence and finalization rules.

## Signals (Raw Inputs)

**Morale Self-Score**
A number from 0 to 5 representing the team member's self-reported morale for the month.

**Billability %**
The percentage of the team member's time billed to client projects. Optimal zone is around 75%; deviation in either direction reduces the normalized score.

**CSAT**
Customer Satisfaction score, rated 1–5.

**Net Margin %**
The net margin percentage associated with the team member's work. Expected range: -20 to +60.

**Positive Feedback Count**
Number of distinct positive feedback items received during the month.

**Critical Feedback Count**
Number of distinct critical/negative feedback items received during the month.

**Overtime Hours**
Total hours worked beyond the standard schedule during the month.

**Delivery Reliability %**
The percentage of committed deliverables met on time during the month. 0–100.

**Mentoring/Enablement Hours**
Hours spent mentoring others, running workshops, or enabling team capability.

**Evidence Notes Count**
Number of structured evidence notes recorded for the team member during the month.

---

## Impact Dimensions

**Impact Ratings (IG, IP, IT, IO)**
For each signal, the Team Lead assigns four ratings (0–5) describing how much that signal affected each of the four impact layers:
- `IG` — Growth impact
- `IP` — Project impact
- `IT` — Team impact
- `IO` — Organization impact

**Impact Layer**
One of four lenses used to evaluate a team member's contribution:
- **Personal Growth (DG)** — learning, wellbeing, development
- **Project Impact (DP)** — delivery, customer satisfaction, commercial result
- **Team Impact (DT)** — collaboration, mentoring, morale contribution
- **Organization Impact (DO)** — margin, strategic value, organizational contribution

---

## Computed Scores

**Normalized Score (N)**
A signal transformed to a 0–100 scale using a formula defined in the PRD (e.g., `MoraleN`, `CSATN`). Normalization makes signals comparable across different units.

**Impact-Weighted Score (ImpactWeighted100)**
The combined impact rating for a signal across all four layers, converted to a 0–100 scale using the layer weights (Growth 20%, Project 35%, Team 25%, Organization 20%).

**Contribution Score (C)**
For each signal: `C = (NormalizedScore × ImpactWeighted100) / 100`. Represents the signal's contribution to the overall score after accounting for both its measured value and the Team Lead's assessment of its significance.

**Dimension Score (DG / DP / DT / DO)**
A weighted sum of Contribution Scores for the signals relevant to each impact layer. Produced by the formula engine. Range: 0–100.

**Total Impact Index (TII)**
The overall monthly score for a team member: `TII = 0.20×DG + 0.35×DP + 0.25×DT + 0.20×DO`. Range: 0–100.

**Completeness %**
`FilledRequiredFields / 11 × 100`. Measures how fully the Team Lead populated the monthly entry.

**Confidence Score**
`0.7×CompletenessPct + 0.3×EvidenceN`. Indicates how much trust to place in the computed scores. Badge levels: High (≥80), Medium (60–79), Low (<60).

---

## Trend Metrics

**Delta1**
`TII_t − TII_{t-1}`. Month-over-month change.

**Delta3**
`TII_t − TII_{t-3}`. Three-month change.

**MA3**
`(TII_t + TII_{t-1} + TII_{t-2}) / 3`. Three-month moving average.

**Vol3**
Standard deviation of TII over the last three months.

**Dimension Delta3**
The three-month change for a specific dimension score (e.g., `DG_Delta3 = DG_t − DG_{t-3}`).

---

## Alerts

**Alert**
A system-generated flag triggered when a signal, score, or trend crosses a defined threshold. Alerts are Amber (warning) or Red (critical). Defined for: Performance Deterioration, Morale Risk, Burnout Risk, Feedback Risk, Customer/Business Risk, Data Quality Risk. Each alert carries the condition that triggered it, its severity, and the member it applies to — enough detail to act on without re-deriving it from raw scores.

**Team-Level Alert**
An alert triggered by aggregate team conditions: Team Morale Drift, Team Delivery Drift, Systemic Burnout, Calibration Risk. Team-level conditions are evaluated from pre-aggregated team-wide inputs (e.g., percentage of members with a given individual alert, team score standard deviation history) rather than computed directly from each member's raw data within the same step.

---

## Decision Guardrails

**Promotion Consideration**
A system recommendation (not a decision) triggered when MA3 ≥ 75, Delta3 ≥ +5, no Red alerts in 3 months, DT or DO ≥ 70, and Confidence ≥ 75.

**Support Plan Trigger**
A system recommendation triggered when MA3 < 55, or 2+ Red alerts current month, or Morale Red + Burnout Amber/Red, and Confidence ≥ 70.

**Demotion Review Gate**
A system gate (not a recommendation) that requires 2 documented support cycles, MA3 < 45 for ≥ 4 months, repeated role-mismatch evidence, and HR compliance confirmation before surfacing.

**Override**
When a Team Lead records a decision that differs from the system recommendation. Overrides require a mandatory rationale and are recorded in the audit trail.

---

## Operational Terms

**Evidence Note**
A free-text note linked to a team member and a category (e.g., Mentoring, Delivery, Feedback). Used to support scores and justify overrides.

**Action Plan**
A tracked follow-up item assigned to a team member with an owner, due date, and status. Created from the Alerts Center or Member Detail screen.

**Audit Trail**
An append-only log of all edits to monthly entries: who changed what field, from what value, to what value, and when.

**Export**
A CSV file containing team overview or member detail data for a given month range. Produced on demand; not synced anywhere.
