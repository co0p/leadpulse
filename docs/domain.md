# LeadPulse Domain Model

This document defines the shared domain language used by LeadPulse. It records durable concepts that cross package boundaries and are referenced by the PRD, architecture, tests, and UI.

## Employee

A person managed by the team lead.

- `firstName`, `secondName`: display name.
- `seniority`: one of Junior, Midlevel, Senior, Principal.
- `startDate`: ISO date string when the employee joined the team.

## Performance Entry

A dated monthly snapshot of an employee's signals. Each entry belongs to a single employee and a single month.

Fields:

- `year` and `month`: the period the entry covers (e.g. 2026-08).
- `morale`: numeric morale score for the month.
- `execution`: numeric execution result for the month.
- `impact`: numeric impact result for the month.
- `growth`: numeric growth result for the month.
- `culture`: numeric culture result for the month.
- `projects`: array of project identifiers the employee worked on during the month.
- `evidenceItems`: array of evidence items backing the entry.

## Evidence Item

A single piece of evidence linked to a performance entry.

- `category`: one of ProjectOutcome, PeerFeedback, BlockerResolution, TechnicalAchievement, Collaboration, Other.
- `source`: one of ManualEntry, GoogleDocImport, ProjectLeadFeedback.
- `description`: free-text reference, example, or note.

## Multiplier Score

The evidence-based multiplier for a month is a 0–5 value. A score of 5 in a given month requires at least three evidence items for that employee in that month. Missing monthly multiplier values are treated as zero in rolling calculations.

## Capacity Risk

An employee is flagged as high capacity risk when their monthly project assignment count is three or more different projects.

## Rolling Three-Month Window

For every metric, the rolling three-month value is the average of the most recent three monthly values, using zero for any missing month.

## Team Pulse Overview

A read model produced for the main screen. It contains:

- Active headcount.
- Rolling average team morale.
- Count of employees with high capacity risk in the latest month.
- Team multiplier index (average of employees' latest rolling multiplier).
- Health dimensions: rolling three-month Execution, Impact, Multiplier, Growth, and Culture results for every employee, aggregated to a team average.
- Actionable anomalies (morale drop ≥ 1.0, high capacity, and 60-day burnout risk when applicable).
- Evidence coverage metadata and evidence-gap state per metric.
