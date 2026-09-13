# Roadmap

Product direction and sequencing for Team Impact Scorecard. Each entry explains the user outcome, current confidence, and ordering rationale.

> A feature moves to **Done** only when its user outcome is verified and the evidence is linked here.
> Source of truth: if a feature is not in Done with a passing test link, it is not considered shipped.

---

## Done

### Team Member Management (CRUD)
- **Job story:** When I set up the tool or my team changes, I want to add, view, edit, and remove team members with their name and seniority level, so that the scorecard always reflects my current team and shows each member's context.
- **Evidence:** 
  - Store layer: 5 integration tests passing (add, list, edit, deactivate, soft-delete behavior)
  - Service layer: 7 unit tests passing (create, list, get, edit, delete, validation)
  - UI layer: Settings screen (Screen F) implemented with Fyne v2
  - Database schema: idempotent SQLite schema with members table and audit_log
  - Acceptance criteria: All 5 criteria met (AC-1: add/view/edit/remove, AC-2: seniority, AC-3: edit, AC-4: deactivate, AC-5: persistence)

---

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

## Planned

### Formula Engine — Core Scoring
- **Job story:** When I submit a team member's monthly data, I want the system to compute all normalized scores, dimension scores (DG/DP/DT/DO), TII, completeness, and confidence, so that I have an objective, repeatable basis for my review.
- **Evidence:** pending — unit tests for all normalization and scoring functions against PRD formulas
- **Why first:** everything else (alerts, trends, UI, export) depends on correct formula output. No other increment can be verified without this.

### Alert Engine
- **Job story:** When scores or trends cross defined thresholds, I want the system to raise Amber or Red alerts automatically, so that I can intervene before a situation worsens.
- **Evidence:** pending — unit tests for all six individual alert types and four team-level alerts, covering both Amber and Red thresholds
- **Why now:** alerts are the primary early-warning output of the tool. Without them the dashboard has no actionable signal.

### Trend Calculations (MA3, Delta1, Delta3, Vol3)
- **Job story:** When I view a team member's scorecard, I want to see their trend over time (moving average, deltas, volatility), so that I can distinguish a one-off bad month from a genuine decline.
- **Evidence:** pending — unit tests for all trend functions; integration tests for multi-month history reads
- **Why now:** trends unlock the decision guardrails (promotion, support plan) and are required for several alert conditions.

### Monthly Input Workspace (Screen B)
- **Job story:** When I open the monthly input workspace, I want to enter raw signals and impact ratings for each team member with live formula preview and completeness tracking, so that I can complete the cycle accurately and efficiently.
- **Evidence:** pending — manual verification of all field validations, live preview updates, and copy-previous-month behavior
- **Ordering:** depends on formula engine and alert engine being complete.

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
