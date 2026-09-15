# PRD Improvement Report — Dimension Model Redesign

Status: Draft for future incorporation into `prd.md`
Related: `docs/prd.md`

## 1) Purpose

This report evaluates a proposed set of scorecard dimensions —
Net Revenue, Billability, Project Impact, Market Relevance, Growth, Satisfaction —
against the dimension model currently defined in the PRD (`DG`/`DP`/`DT`/`DO`
in §5.4), identifies signal-sourcing feasibility, flags blindspots in the
current design, and proposes a consolidated 6-dimension model.

---

## 2) Categorical Mismatch

The proposed 6 dimensions and the PRD's current 4 are not the same *type*
of construct, and merging them naively would create overlap.

| Type | Examples | What it answers |
|---|---|---|
| **Impact layers** (current PRD) | Growth, Project, Team, Organization | *Who benefits* from the work |
| **KPI/outcome metrics** (proposed) | Net Revenue, Billability, Satisfaction | *What was produced/felt* |

`Project Impact` and `Growth` are impact-layer concepts already present in
the PRD. `Net Revenue`, `Billability`, `Satisfaction` are raw outcome
metrics — closer to the PRD's §4.1 "Raw Inputs" than to §5.4 "Dimension
Scores." `Market Relevance` is neither — it does not exist in the PRD
today.

Treating all 6 as top-level dimensions without resolving this would cause
double-counting: e.g., Billability already feeds into `DP`, `DT`, `DO` as a
*signal*, not a dimension.

---

## 3) Signal Sourcing Map

Assessment of how easily each proposed dimension's underlying data can be
gathered, and its current coverage in the PRD.

| Dimension | Signal needed | Source system | Effort to automate | Current PRD coverage |
|---|---|---|---|---|
| **Net Revenue** | Net Margin %, revenue attribution per person | Finance/PSA system (Harvest, NetSuite, custom billing) | Medium — requires per-employee cost/revenue allocation logic | `NetMarginPct` exists, but only in DP/DO |
| **Billability** | Billable % of hours | Time tracking (Harvest, Toggl, PSA tool) | Easy — usually already tracked for invoicing | `BillabilityPct` exists |
| **Project Impact** | Delivery reliability, CSAT, rework/defect rate | PM tool (Jira/Asana/Linear) + client survey tool | Medium — delivery % is easy, quality/rework is not tracked today | Delivery + CSAT exist; rework/quality missing |
| **Market Relevance** | Skill-to-demand alignment, client renewal/expansion tied to person, strategic project weighting | CRM (renewal data), skills/LMS platform, sales/strategy input | Hard — mostly manual/judgment-based today | Not covered |
| **Growth** | Skills acquired, certifications, career-readiness signals | LMS, HRIS, manager 1:1 notes | Medium — training completion is trackable; readiness is subjective | Only a *derived* impact rating (`IG`), not a direct measurement |
| **Satisfaction** | Morale, engagement, burnout risk | Pulse survey tool (Officevibe, Culture Amp) or simple self-score | Easy — self-report; harder to get unbiased/frequent signal | `Morale Self-Score` exists but single-item, self-reported only |

**Takeaway:** Billability, Delivery, and Morale are cheap to automate
(source systems likely already exist). Market Relevance and Net Revenue
attribution are the expensive ones — they require finance integration or
manual strategic judgment, which conflicts with the PRD's goal of
"structured, low-bias signals" (§1).

---

## 4) Blindspots in the Current Design

1. **Single-rater bias is structural, not just a risk.** Every signal in
   the PRD (§4.2 impact ratings, feedback counts, evidence notes) is
   entered by the *team lead alone*. There is no peer, client-direct, or
   upward feedback loop. This undermines the PRD's own stated goal of
   "reducing bias" (§1). A lightweight 360-style input is a real gap.

2. **No quality/rework signal.** `DeliveryReliabilityPct` measures
   *on-time*, not *right-first-time*. A person who ships late-but-correct
   vs. on-time-but-buggy scores identically today.

3. **Market Relevance has no natural data owner.** Unlike other signals,
   this requires input from sales/strategy roles the team lead doesn't
   have visibility into. Without a defined data source, this dimension
   risks becoming a subjective "gut feel" slider — the opposite of the
   PRD's evidence-based philosophy.

4. **Growth is entirely derived, never directly measured.** `IG` is a 0–5
   *rating* the team lead assigns per signal — there is no actual growth
   data (skills gained, scope expansion, certifications). This makes
   "Growth" circular: it is an opinion about impact, not evidence of
   growth.

5. **No innovation/initiative signal.** Process improvements, proactive
   problem-solving, and IP contributions are not captured anywhere in raw
   inputs.

6. **No attrition/flight-risk signal distinct from morale.** Morale is a
   snapshot; retention risk (tenure trend, market comp gap, engagement
   trajectory) is a different and increasingly important construct for a
   "holistic view."

7. **Net Revenue attribution is fragile for team-based work.** Margin is
   easy to compute per-project but hard to fairly attribute per-individual
   on shared deliverables — this could introduce more gaming risk than it
   removes (contradicts Risk #1 in §18).

---

## 5) Proposed Consolidated 6-Dimension Model

To stay within a 6-dimension limit while resolving the categorical
mismatch and closing the most important blindspots, restructure as
**outcome-oriented dimensions**, each backed by a defined signal source,
replacing `DG`/`DP`/`DT`/`DO`:

| # | Dimension | Absorbs (from current PRD) | New signals needed | Primary source |
|---|---|---|---|---|
| 1 | **Financial Contribution** | Net Margin, Billability, Overtime | — | Finance/PSA (mostly automatable) |
| 2 | **Delivery & Quality** (renamed Project Impact) | Delivery Reliability, CSAT | Rework/defect rate | PM tool + client survey |
| 3 | **Team Collaboration** | Mentoring, Positive/Critical Feedback | Peer/360 input (lightweight) | Manager + peer pulse |
| 4 | **Market & Strategic Relevance** (new) | — | Client renewal/expansion tied to person, skill-market alignment tag | Sales/CRM + skills registry |
| 5 | **Growth & Capability** | (was `IG` derived rating) | Certifications/training completed, scope expansion evidence | LMS/HRIS + evidence notes |
| 6 | **Satisfaction & Well-being** | Morale | Engagement pulse (multi-item, not single self-score), burnout/retention indicator | Pulse survey tool |

This keeps the model at exactly 6 dimensions, gives each dimension a
**named data source**, and folds Net Revenue + Billability into a single
Financial Contribution dimension — combined because they are both
financial-efficiency signals, and attributing them separately adds
complexity without much decision value.

### Trade-off to resolve before incorporation

"Organization Impact" (`DO`) as a standalone concept disappears under this
model — its signals get redistributed into Financial Contribution, Market
Relevance, and Satisfaction. The PRD's promotion guardrail in §9 currently
checks `DO >= 70`. A decision is needed on which of the 6 dimensions
substitutes for that gate — **Market & Strategic Relevance** is the
suggested candidate.

---

## 6) Open Questions for Incorporation

1. Should Net Revenue and Billability remain merged into one "Financial
   Contribution" dimension, or be tracked separately despite the added
   complexity?
2. What is the intended data source/owner for Market & Strategic
   Relevance in this org — sales, delivery leadership, or the team lead's
   judgment?
3. Should the promotion guardrail (§9) be rewritten around the new
   6-dimension names, and which dimension(s) replace the `DO >= 70`
   check?
4. Is a lightweight peer/360 input in scope for v1, or deferred to a later
   version given the PRD's current v1 scope is team-lead-only (§3)?
5. Should quality/rework tracking be added as a new raw input (§4.1) to
   support the "Delivery & Quality" dimension, and if so, what is the
   source system?
</content>
</invoke>
