# ADR-20260914: Alert Engine — Resolving Ambiguous PRD Threshold Rules

## Status
Accepted

## Date
2026-09-14

## Context

The PRD (`docs/prd.md` §7.1–7.2) defines exact numeric thresholds for all Alert Engine conditions, but two rules were not fully specified and required an interpretation before implementation:

1. **Feedback Risk Red condition:** "CriticalFeedbackCount >= 4 AND declining 2-month critical trend." The PRD does not define which direction "declining" refers to. Critical feedback is a negative signal — a *lower* count is an improving situation, so "declining trend" for the *count itself* is directionally ambiguous: it could mean the count is decreasing (literal reading of "declining"), or it could mean the underlying situation is declining/worsening (the count is flat or rising).

2. **Calibration Risk condition:** "team TII stddev < 6 for 3 months." The PRD does not specify whether this is a single stddev computed once across a 3-month window, or a cross-sectional stddev computed once per month and sustained below 6 for 3 consecutive months.

Both rules gate a Red or Amber alert that a Team Lead will act on. Implementing the wrong interpretation would silently produce incorrect alerts — either raising false alarms or missing real risk — with no compiler or test failure to catch the mismatch against product intent, since the code would be internally consistent either way.

## Decision

1. **Feedback Risk "declining 2-month critical trend" means the underlying situation is not improving** — i.e., the current month's `CriticalFeedbackCount` is greater than or equal to the prior month's count (`current >= prior`). A decrease in critical feedback count is treated as improvement and does not qualify for Red, even if the current count is still >= 4. Without prior-month data, Red cannot be confirmed (no trend exists to evaluate); the alert falls back to Amber in that case.

2. **Calibration Risk uses a cross-sectional stddev, computed once per month, and requires all of the 3 most recent monthly values to be below 6.** `EvaluateCalibrationRisk` accepts a 3-element slice of pre-computed monthly stddev values (not raw scores) and returns Amber only if every value in the slice is strictly below 6. Fewer than 3 months of history returns an error — the condition cannot be evaluated with insufficient history.

## Alternatives Considered

### Feedback Risk: literal "declining" = count is decreasing
**Pro:** Matches the literal English meaning of "declining trend" applied directly to the count value.
**Con:** Produces the wrong alert semantics — a member whose critical feedback count is *decreasing* (situation improving) would trigger Red, while a member whose count is flat or rising at a severe level would not. This is the opposite of what a burnout/performance early-warning system should flag.

### Calibration Risk: single stddev over a 3-month window (one value, not three)
**Pro:** Simpler input — one number instead of three.
**Con:** Conflates "3 months of individually-compressed variance" with "variance measured once across a 3-month sample," which are different statistical claims. The PRD's plural framing ("for 3 months") reads as a sustained monthly condition, not a single computation over a wider window.

**Why the chosen approach wins:** Both interpretations align the alert's purpose (early warning for a Team Lead) with its literal trigger. The Feedback Risk interpretation ensures Red fires only when the negative situation is not improving. The Calibration Risk interpretation ensures Amber fires only when score compression has been sustained, not a one-off statistical artifact.

## Consequences

### Positive

1. **Alert semantics match product intent**, not just literal PRD wording — critical for a system that recommends Team Lead intervention.
2. **Conservative default for Feedback Risk** — a member's first tracked month can never reach Red for this condition (no prior-month data to confirm a trend), even with severe current-month feedback. This avoids false Red alerts from insufficient history and is documented in code (`engine/domain/alerts.go`, `EvaluateFeedbackRisk`).
3. **Explicit history requirement for Calibration Risk** — `EvaluateCalibrationRisk` returns an error rather than silently guessing when fewer than 3 months of data exist, preventing incorrect evaluation on partial history.

### Negative (Managed Risks)

1. **Deviates from a strictly literal reading of the PRD.** If the original PRD author intended the literal reading, this ADR documents the deviation and its rationale so it can be revisited.
2. **Feedback Risk Red is unreachable for new members' first month**, regardless of severity. This is an intentional design trade-off, not a bug, but could surprise a Team Lead expecting immediate escalation for a severe first-month case. Mitigation: Amber still fires at `CriticalFeedbackCount >= 3` in the first month, so the condition is not silent — it's capped at Amber until a second month of data exists.

## Implementation Notes

- Both interpretations are implemented as pure functions in `engine/domain/alerts.go`: `EvaluateFeedbackRisk(currentCritical, priorCritical int, hasPriorMonth bool) AlertSeverity` and `EvaluateCalibrationRisk(stddevHistory []float64) (AlertSeverity, error)`.
- Test coverage for both interpretations exists in `engine/domain/alerts_test.go`, including the boundary case where the trend is improving (Feedback Risk stays Amber, not Red) and the case where any one of 3 months fails the stddev threshold (Calibration Risk stays None).
- If the PRD (`docs/prd.md`) is later revised to state either rule literally, this ADR and the corresponding code must be reconciled together — do not update the PRD language without revisiting the implementation, or vice versa.

## Related

- `docs/prd.md` §7.1 (Feedback Risk), §7.2 (Calibration Risk) — the source rules being interpreted
- `docs/domain.md` — Alert and Team-Level Alert glossary entries
- `engine/domain/alerts.go`, `engine/domain/alerts_test.go` — implementation and test evidence

## Questions for Future Work

1. Should the PRD itself be updated to state these interpretations explicitly, removing the ambiguity at the source?
2. When team-level alert input aggregation is built (deferred to a future increment), should the stddev history window be configurable, or fixed at exactly 3 months as implemented here?
