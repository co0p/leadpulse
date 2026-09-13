# Team Impact Scorecard — Product Requirements Document (PRD v1)

## 1) Product Summary
The Team Impact Scorecard is a monthly decision-support tool for **Team Leads** to evaluate individual and team health using structured signals, trend analysis, and evidence logs.  
It is designed to reduce bias in people decisions (promotion/support/demotion discussions) by making rationale explicit, repeatable, and auditable.

---

## 2) Goals and Non-Goals

### Goals
1. Provide a monthly, evidence-based view of each team member’s impact.
2. Connect all signals to impact layers:
   - Personal Growth
   - Project
   - Team
   - Organization
3. Detect negative trend changes early via warning signals.
4. Support fair, rational team-lead decisions with transparent formulas.
5. Enable exportable review snapshots.

### Non-Goals (v1)
1. Fully automated HR decisions.
2. Employee self-service portal.
3. Cross-company benchmarking.
4. HR administration suite.
5. Compensation automation.

---

## 3) Target User

### Primary Persona
- **Team Lead** (only persona in v1)

### Permissions
- View all direct reports
- Enter/edit monthly data
- Add evidence notes
- Review trends/alerts/recommendations
- Finalize monthly cycle
- Export summaries

---

## 4) Monthly Data Inputs (Per Member)

## 4.1 Raw Inputs
1. Morale Self-Score (0–5)
2. Billability % (0–100)
3. CSAT (1–5)
4. Net Margin % (suggested range: -20 to +60)
5. Positive Feedback Count (>=0)
6. Critical Feedback Count (>=0)
7. Overtime Hours (>=0)
8. Delivery Reliability % (0–100)
9. Mentoring/Enablement Hours (>=0)
10. Evidence Notes Count (>=0)

## 4.2 Impact Ratings (per signal)
For each signal above, team lead records:
- `IG` (Growth impact) 0–5
- `IP` (Project impact) 0–5
- `IT` (Team impact) 0–5
- `IO` (Organization impact) 0–5

## 4.3 Computed Input Quality
- `CompletenessPct = FilledRequiredFields / 11 * 100`

---

## 5) Scoring System

## 5.1 Utility Functions
- `clamp(x, min, max) = min(max(x, min), max)`
- `norm01(x, min, max) = clamp((x - min)/(max - min), 0, 1)`

## 5.2 Normalization Formulas
- `MoraleN = (Morale / 5) * 100`
- `BillabilityN = clamp(100 - (abs(BillabilityPct - 75) / 15) * 100, 0, 100)`
- `CSATN = ((CSAT - 1) / 4) * 100`
- `MarginN = norm01(NetMarginPct, -20, 60) * 100`
- `PositiveN = min(PositiveFeedbackCount, 8) / 8 * 100`
- `CriticalN = 100 - (min(CriticalFeedbackCount, 6) / 6 * 100)`
- `OvertimeN = 100 - (min(OvertimeHours, 30) / 30 * 100)`
- `DeliveryN = clamp(DeliveryReliabilityPct, 0, 100)`
- `MentoringN = min(MentoringHours, 12) / 12 * 100`
- `EvidenceN = min(EvidenceNotesCount, 6) / 6 * 100`

## 5.3 Impact Weight per Signal
Impact layer default weights:
- Growth 20%
- Project 35%
- Team 25%
- Organization 20%

- `ImpactWeighted = 0.20*IG + 0.35*IP + 0.25*IT + 0.20*IO` (0–5)
- `ImpactWeighted100 = ImpactWeighted * 20` (0–100)

For each signal `s`:
- `Cs = (Ns * Is) / 100`
  - `Ns` = normalized signal score
  - `Is` = impact-weighted score (0–100)

## 5.4 Dimension Scores

### A) Personal Growth (DG)
Weights:
- MoraleN 20
- CriticalN 15
- PositiveN 10
- MentoringN 15
- DeliveryN 10
- BillabilityN 10
- OvertimeN 15
- EvidenceN 5

`DG = Σ(weight_i * C_i) / 100`

### B) Project Impact (DP)
Weights:
- DeliveryN 25
- CSATN 20
- MarginN 20
- BillabilityN 15
- CriticalN 10
- PositiveN 5
- MoraleN 5

`DP = Σ(weight_i * C_i) / 100`

### C) Team Impact (DT)
Weights:
- MentoringN 25
- PositiveN 20
- CriticalN 20
- MoraleN 15
- DeliveryN 10
- OvertimeN 10

`DT = Σ(weight_i * C_i) / 100`

### D) Organization Impact (DO)
Weights:
- MarginN 25
- CSATN 20
- BillabilityN 15
- DeliveryN 15
- MentoringN 10
- PositiveN 10
- EvidenceN 5

`DO = Σ(weight_i * C_i) / 100`

## 5.5 Total Impact Index (TII)
Dimension weights:
- DG 20%
- DP 35%
- DT 25%
- DO 20%

`TII = 0.20*DG + 0.35*DP + 0.25*DT + 0.20*DO`

---

## 6) Trend Logic

For month `t`:
- `Delta1 = TII_t - TII_{t-1}`
- `Delta3 = TII_t - TII_{t-3}`
- `MA3 = (TII_t + TII_{t-1} + TII_{t-2}) / 3`
- `Vol3 = stddev(TII_t, TII_{t-1}, TII_{t-2})`

Dimension trend examples:
- `DG_Delta3 = DG_t - DG_{t-3}` (same for DP/DT/DO)

---

## 7) Alert Engine

## 7.1 Individual Alerts

### Performance Deterioration
- Amber: `Delta1 <= -6` OR `Delta3 <= -10`
- Red: `Delta1 <= -10` OR `Delta3 <= -15`

### Morale Risk
- Amber: `MoraleN < 50` for 1 month
- Red: `MoraleN < 50` for 2 consecutive months OR `MoraleN < 35` current month

### Burnout Risk
- Amber: `OvertimeHours >= 20` AND `MoraleN < 60`
- Red: `OvertimeHours >= 25` AND (`Delta1 < 0` OR `DeliveryN < 60`)

### Feedback Risk
- Amber: `CriticalFeedbackCount >= 3`
- Red: `CriticalFeedbackCount >= 4` AND declining 2-month critical trend

### Customer/Business Risk
- Amber: `CSATN < 60` OR `MarginN < 45`
- Red: `CSATN < 50` AND `MarginN < 40`

### Data Quality Risk
- Amber: `CompletenessPct < 85`
- Red: `CompletenessPct < 70`

## 7.2 Team-Level Alerts
- Team Morale Drift Red: >=30% members have Morale Red
- Team Delivery Drift Red: team `Delta3 <= -10`
- Systemic Burnout Red: >=25% members Burnout Red
- Calibration Risk Amber: team TII stddev < 6 for 3 months

---

## 8) Confidence Score

`Confidence = 0.7*CompletenessPct + 0.3*EvidenceN`

Badges:
- High >= 80
- Medium 60–79
- Low < 60

Rule:
- If Confidence < 60, show “Use caution: insufficient evidence quality/completeness.”

---

## 9) Decision Guardrails (Recommendation Only)

### Promotion Consideration
Require all:
1. `MA3 >= 75`
2. `Delta3 >= +5`
3. No Red alerts in last 3 months
4. `DT >= 70` OR `DO >= 70`
5. `Confidence >= 75`

### Support Plan Trigger
If any:
1. `MA3 < 55`
2. 2+ Red alerts current month
3. Morale Red + Burnout Amber/Red
4. Confidence >= 70

### Demotion Review Gate
Only after:
1. 2 documented support cycles
2. `MA3 < 45` for >= 4 months
3. Repeated role mismatch evidence
4. HR/process compliance confirmed

---

## 10) Main Screens

## Screen A: Overview Dashboard
Purpose: monthly team health snapshot.

Components:
- KPI strip (Team Avg TII, Delta1/Delta3, Green/Amber/Red counts)
- Team TII trend chart
- DG/DP/DT/DO heatmap per member
- Alert table (severity, reason, CTA)
- Action queue (review-ready buckets)

Primary CTAs:
- Start Monthly Cycle
- Review Alerts
- Export Snapshot

## Screen B: Monthly Input Workspace
Purpose: capture monthly signals and impact ratings.

Components:
- Member list
- Raw signal form
- IG/IP/IT/IO rating matrix
- Live score/alert preview panel
- Validation + completeness warnings

CTAs:
- Save Draft
- Submit Member Month
- Copy Previous Month

## Screen C: Member Detail
Purpose: 1:1 and performance discussion support.

Components:
- Current TII + dimensions + confidence + active alerts
- 6–12 month trend timeline
- Signal contribution table
- Evidence notes log
- Action plan tracker

CTAs:
- Add Evidence
- Add Action
- Export Member Summary

## Screen D: Alerts & Risk Center
Purpose: triage worsening trends fast.

Components:
- Filterable alert list
- Root cause summary
- Suggested intervention options
- Due dates / status tracking

CTAs:
- Assign action
- Resolve/Snooze alert

## Screen E: Monthly Review & Calibration
Purpose: fair monthly review finalization.

Components:
- Distribution charts (anti-compression view)
- Member decision-support cards
- Override reason capture (mandatory evidence)

CTAs:
- Finalize Month
- Export Review Packet
- Lock Cycle

## Screen F: Limited Settings
Purpose: team-level threshold tuning.

Allowed:
- Billability target/tolerance
- Alert sensitivity preset
- Reminder dates

Not allowed:
- Core formula redesign
- Audit log edits

---

## 11) Monthly Operating Rhythm

### Phase 0 (Day -3 to 0): Prepare
- New cycle opens
- Reminder sent
- Previous values available for carry-forward

### Phase 1 (Day 1–5): Data Entry
- Team lead enters all member data + impact ratings + evidence

Exit:
- No validation errors
- Completeness target met

### Phase 2 (Day 5–6): Auto Compute
- Scores, trends, confidence, alerts generated

Exit:
- Alerts reviewed once
- Red alerts annotated

### Phase 3 (Day 6–8): Triage
- Team lead assigns actions for risk cases

Exit:
- Every red alert has owner + next step

### Phase 4 (Day 8–10): Monthly Review
- Apply promotion/support guardrails
- Record overrides with evidence

Exit:
- Cycle finalized + exported

### Phase 5 (Day 10–20): 1:1 Execution
- Use member profiles in conversations
- Capture commitments and follow-up

### Phase 6 (Remainder): Follow-through
- Mid-month risk checks
- Ongoing note capture for next cycle

---

## 12) User Stories

1. As a Team Lead, I want one dashboard to see who needs attention this month.
2. As a Team Lead, I want to input monthly signals quickly with validation.
3. As a Team Lead, I want trends and alerts so I can intervene early.
4. As a Team Lead, I want evidence-linked recommendations to reduce bias.
5. As a Team Lead, I want exportable summaries for review meetings.

---

## 13) Functional Requirements

1. System must compute all normalized metrics on save.
2. System must recompute DG/DP/DT/DO/TII in real time.
3. System must generate alerts using exact thresholds.
4. System must compute MA3/Delta1/Delta3 where enough history exists.
5. System must show confidence badge on overview + member pages.
6. System must require reason + evidence when overriding recommendation.
7. System must block finalization if any member completeness < 70%.
8. System must keep audit trail of edits (who, when, what changed).
9. System must allow CSV export for team and member views.

---

## 14) Non-Functional Requirements

1. Page load < 2s for team up to 100 members.
2. Formula determinism: same input => same output always.
3. Full traceability of computed values (drill-down visible).
4. Role-based access (team lead only in v1).
5. Data retention: minimum 24 months trend history.

---

## 15) API Contract (Frontend-Oriented)

## 15.1 POST Monthly Entry
`POST /api/v1/scorecards/{memberId}/{month}`

Request (example):
```json
{
  "raw": {
    "morale": 3.8,
    "billabilityPct": 78,
    "csat": 4.4,
    "netMarginPct": 22,
    "positiveFeedbackCount": 5,
    "criticalFeedbackCount": 2,
    "overtimeHours": 14,
    "deliveryReliabilityPct": 86,
    "mentoringHours": 6,
    "evidenceNotesCount": 4
  },
  "impactRatings": {
    "morale": {"IG":4,"IP":3,"IT":4,"IO":2},
    "billability": {"IG":2,"IP":4,"IT":2,"IO":4},
    "csat": {"IG":3,"IP":5,"IT":3,"IO":4},
    "margin": {"IG":2,"IP":4,"IT":2,"IO":5},
    "positiveFeedback": {"IG":3,"IP":3,"IT":5,"IO":3},
    "criticalFeedback": {"IG":4,"IP":4,"IT":4,"IO":3},
    "overtime": {"IG":4,"IP":3,"IT":3,"IO":2},
    "delivery": {"IG":3,"IP":5,"IT":3,"IO":4},
    "mentoring": {"IG":5,"IP":3,"IT":5,"IO":4},
    "evidence": {"IG":2,"IP":2,"IT":2,"IO":3}
  }
}
```

Response:
```json
{
  "memberId": "u_123",
  "month": "2026-08",
  "computed": {
    "DG": 69.4,
    "DP": 73.1,
    "DT": 71.8,
    "DO": 68.9,
    "TII": 71.2,
    "MA3": 69.5,
    "Delta1": 2.1,
    "Delta3": 6.8,
    "CompletenessPct": 100,
    "Confidence": 90
  },
  "alerts": []
}
```

## 15.2 GET Team Overview
`GET /api/v1/teams/{teamId}/overview?month=YYYY-MM`

Response:
```json
{
  "teamId": "team_1",
  "month": "2026-08",
  "kpis": {
    "teamAvgTII": 68.4,
    "teamDelta1": -1.8,
    "teamDelta3": -4.9,
    "greenCount": 7,
    "amberCount": 3,
    "redCount": 2
  },
  "members": [
    {
      "memberId": "u_123",
      "name": "Alex",
      "DG": 69.4,
      "DP": 73.1,
      "DT": 71.8,
      "DO": 68.9,
      "TII": 71.2,
      "severity": "green"
    }
  ],
  "teamAlerts": []
}
```

## 15.3 GET Member Detail
`GET /api/v1/scorecards/{memberId}?months=12`

Returns:
- monthly raw inputs
- computed scores
- trends
- alerts
- evidence timeline
- actions timeline

## 15.4 POST Finalize Cycle
`POST /api/v1/teams/{teamId}/cycles/{month}/finalize`

Rules:
- reject if any member completeness < 70
- return finalization summary + export links

---

## 16) Acceptance Criteria by Screen

## A) Overview Dashboard
- [ ] Shows KPI strip with current month values
- [ ] Heatmap colors match thresholds
- [ ] Alert list sortable by severity and member
- [ ] Clicking alert opens member detail in 1 click

## B) Monthly Input
- [ ] All fields validated by allowed range
- [ ] Live formula recomputation on input change
- [ ] Shows completeness in real time
- [ ] Can copy prior month data

## C) Member Detail
- [ ] Shows 12-month trend chart
- [ ] Shows contribution table with raw + normalized + weighted values
- [ ] Can add/edit evidence notes
- [ ] Shows actionable next-step section

## D) Alerts Center
- [ ] Supports filtering by type/severity
- [ ] Displays exact threshold breach reason
- [ ] Supports resolve/snooze with reason capture

## E) Monthly Review
- [ ] Shows recommendation guardrail status
- [ ] Override requires mandatory rationale
- [ ] Finalize blocked by data-quality rule

---

## 17) UX Content Guidelines

- Prefer “Needs attention” over “Low performer”
- Show “Recommendation (human review required)”
- Always pair score with trend and confidence
- Highlight improvement opportunities, not only deficits

---

## 18) Risks and Mitigations

1. **Metric gaming**
   - Mitigation: mixed signals + evidence requirement + trend emphasis

2. **Bias in manual ratings**
   - Mitigation: calibration screen + audit logs + distribution checks

3. **Over-indexing on financial metrics**
   - Mitigation: explicit Growth/Team/Org dimensions and minimum thresholds

4. **Data incompleteness**
   - Mitigation: completeness alerts + finalize guardrails

---

## 19) Implementation Plan (Frontend MVP)

### Sprint 1
- Data model + input forms
- Formula engine
- Member detail baseline

### Sprint 2
- Overview dashboard
- Alert engine + risk center
- CSV export

### Sprint 3
- Monthly review + calibration
- Finalization workflow
- Audit trail UI hooks

---

## 20) Definition of Done (v1)
- End-to-end monthly cycle works for team lead
- All formulas implemented and verified
- Alerts generated correctly
- Export works for overview and member detail
- Finalization + guardrails enforced
- Basic auditability in place