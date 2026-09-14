# UI Decisions

Permanent UI decisions for Team Impact Scorecard. Records shared interaction patterns, visual principles, accessibility rules, and recurring design choices with their rationale. Updated when a new pattern is introduced, not when a screen is added.

---

## Framework and Rendering

**Fyne v2** is the UI framework. All views, widgets, and layouts are implemented in Go using Fyne's canvas and widget APIs. There is no HTML, CSS, or JavaScript.

The visual language described below is translated into Fyne widget properties, custom renderers, and layout containers — not HTML/Tailwind snippets. The HTML examples in the original UI notes are retained as visual references for layout intent only.

---

## Layout Architecture

**Application shell: Sidebar + Canvas**

- **Left sidebar (fixed):** Primary navigation — Overview, Input, Alerts, Review — plus a Settings link at the bottom. Light background. Narrow fixed width.
- **Top header (sticky):** Cycle context selector (e.g., "Cycle: August 2026") and Team Lead identity. Slim bar.
- **Main content area (scrollable):** Off-white background. White content cards rendered on top to create depth.

This layout is used on every screen. Navigation between screens replaces the main content area only; sidebar and header do not re-render.

**Rationale:** A single consistent shell reduces cognitive load. The Team Lead is always oriented: navigation is always visible, the current cycle is always visible.

---

## Color System

Severity and confidence are communicated through a consistent color set. These values must not be used for decoration.

| Semantic | Background | Text | Use |
|---|---|---|---|
| Red alert | red-100 equivalent | red-800 equivalent | Critical threshold breach |
| Amber warning | yellow-100 equivalent | yellow-800 equivalent | Warning threshold breach |
| Green / High confidence | emerald-100 equivalent | emerald-800 equivalent | Healthy status, high confidence |
| Neutral tag | slate-100 equivalent | slate-700 equivalent | Dimension labels, neutral categories |

Background: off-white (`slate-50` equivalent) for the content area. White for data cards.

---

## Status Pills

Pills are the only mechanism for communicating severity, confidence level, or dimension membership inline. They are small, rounded, high-contrast labels.

Pill types:
- Red Alert
- Amber Risk
- High Confidence / Medium Confidence / Low Confidence
- Dimension tags: Growth, Project, Team, Org

Pills are never used for decorative purposes. Every pill must map to a domain state.

---

## KPI / Metric Cards (Screen A)

Used in the Overview dashboard KPI strip. White background, subtle border, light shadow. Structure:
- Label (small, muted)
- Primary value (large, bold, dark)
- Trend indicator below the value (colored arrow + delta value)

The trend indicator color matches alert severity: red for negative delta, green for positive.

---

## Alert Triage Rows (Screen D)

Each alert is a single row containing:
- Severity pill (left)
- Member name + alert reason (center)
- Single action CTA button (right)

Rows are not nested. No tabs, no accordions. The full alert reason is visible without interaction.

---

## Split-Screen Input Layout (Screen B)

The Monthly Input Workspace uses a two-column split:
- Left column (narrow): scrollable member list with status badges (Draft / Done)
- Right column (wide): the active member's data entry form

Implementation notes (Screen B):

- Each signal field is a single-line numeric Entry. Empty means "not entered" and is persisted as NULL.
- The right column includes a sticky live preview panel showing TII, completeness percent (progress bar), and a confidence pill.
- The Save button is disabled until completeness >= 70% (policy captured from domain/completeness thresholds). The team lead may still fill fields before enabling Save.
- "Copy from previous month" populates the form but does not auto-save.

Active member is highlighted in the list. Switching members replaces the right column only.

---

## Live Preview Panel (Screen B)

A sticky summary panel in the right column of the input workspace shows real-time computed values as the Team Lead types:
- TII (large, bold)
- Completeness % (progress bar)
- Confidence badge (pill)

The panel updates on every field change. It does not require a save action to refresh.

---

## Decision Support Cards (Screen E)

Each recommendation (Promotion Consideration, Support Plan Trigger) is a card with a colored left border matching the recommendation type (green for promotion, amber for support). Contents:
- Recommendation type label (small caps, colored)
- Member name (prominent)
- Guardrail status (which criteria are met/unmet)
- Override button (requires mandatory rationale on activation)

Cards are never presented as verdicts. The label always reads "Recommendation (human review required)".

---

## Heatmap Grid (Screen A)

The team dimension heatmap is a table where each cell's background color encodes the score:
- Green (> 75)
- Amber (50–74)
- Red (< 50)

Column headers: Member, Growth (DG), Project (DP), Team (DT), Org (DO). The score value is printed inside each cell. No tooltip required — the value is always visible.

---

## Signal Contribution Table (Screen C)

A table showing the math breakdown for each signal:
- Signal name
- Raw value
- Normalized score
- Associated dimension tags (pills)

Used in Member Detail for 1:1 discussion support. All values are always visible; no expand/collapse.

---

## Evidence Notes Timeline (Screen C)

A chronological feed. Each item contains:
- Author and date (top, small)
- Category pill (top right)
- Note body (full text, no truncation)

Items are sorted newest-first. No pagination in v1 (24-month history per member is manageable in a scrollable list).

---

## Accessibility

- All interactive elements must be reachable via keyboard navigation.
- Severity and confidence states must not be communicated by color alone — pill labels carry the text meaning.
- Minimum contrast for text on colored backgrounds: 4.5:1 (WCAG AA).
- Fyne's default font size must not be reduced below the framework default.

---

## Navigation Shell Implementation

The sidebar + header shell described above is implemented in `ui/app.go` using `container.NewBorder`:

- **Sidebar:** `container.NewStack` wrapping a `canvas.NewRectangle` background (slate-100) and a `container.NewBorder` that places the app title + separator at the top, Settings button + separator at the bottom, and primary nav buttons in the center.
- **Header:** `container.NewBorder` with a separator at the bottom and an `HBox` containing the current cycle label (formatted as "Month YYYY").
- **Content area:** `container.NewStack` whose `Objects` slice is replaced on every nav selection. `content.Refresh()` redraws the area.
- **Nav buttons:** `widget.Button` with `Alignment: widget.ButtonAlignLeading`. Each button's `OnTapped` closure captures its screen loader and calls `setScreen`.
- **Screen registry:** `ui/screens.ScreenRegistry` exposes one method per screen. Screens not yet implemented return a centred placeholder label via `screens.PlaceholderScreen(title)`.

**Rationale:** `container.NewBorder` provides the fixed sidebar and fixed header without manual sizing. Replacing `content.Objects` is the idiomatic Fyne pattern for single-content-area navigation; it avoids tabs (which imply peer screens) and matches the sidebar-driven UX described in the layout architecture section.

---

## Settings Screen Form Dialogs (Screen F)

The Settings screen uses standard Fyne `dialog.ShowForm` patterns for CRUD operations:

### Member List Layout

- **Container:** `container.NewBorder` with:
  - Top: Title ("Settings") + subtitle ("Manage team members") + separator
  - Center: Scrolling member list
  - Bottom: Action button bar ("+ Add Member")
- **List items:** Horizontal box with [Full Name] [Seniority] [spacer] [Edit] [Deactivate]
- **Rationale:** Border layout ensures the list fills available space; button callbacks are wired per row so each item's Edit/Deactivate actions operate on the correct member ID.

### Add Member Dialog

- **Trigger:** "+ Add Member" button
- **Fields:** First Name (text entry), Last Name (text entry), Seniority (dropdown)
- **Defaults:** Seniority defaults to the first value (Junior)
- **Actions:** "Add" (create and refresh list) or "Cancel"
- **Behavior:** Dialog dismisses after successful add; list refreshes to show the new member

### Edit Member Dialog

- **Trigger:** "Edit" button per list item
- **Fields:** Pre-filled with current member data
- **Actions:** "Save" (update and refresh list) or "Cancel"
- **Behavior:** Dialog pre-populates all fields; changes persist only when "Save" is clicked

### Deactivate Member

- **Trigger:** "Deactivate" button per list item
- **Behavior:** No confirmation dialog in v1; soft-delete is immediate. List refreshes. Historical data is preserved in the database.
- **Future:** v2 should add "Are you sure?" confirmation before deactivating

### Error Handling

- Service-layer validation (empty names, invalid seniority) returns errors silently in v1. Errors are logged but not shown to the user.
- Future: v2 should display validation errors in the dialog itself (e.g., "First name is required").

---

## Content Conventions

From the PRD:
- Use "Needs attention" not "Low performer".
- Always show "Recommendation (human review required)" — never "Decision".
- Always pair a score with its trend and confidence level.
- Highlight improvement opportunities alongside deficits.
- Override prompts must ask for evidence, not just confirmation.
