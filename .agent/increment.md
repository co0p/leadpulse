# Increment: Web App Shell (Foundation for SPA)

**Date:** 2026-09-16

---

## Use Case

When I load the web app in my browser, I want to see a professional, accessible layout with a persistent sidebar, top navigation bar featuring a centered search bar, and content area, so that I have a foundation for building screens and the app feels polished from the start.

---

## Goal

Build the web app shell with sidebar, top bar (with centered search), and main content area using HTMX + Alpine.js + Go templates + Bulma CSS.

---

## Branch

`increment/web-app-shell`

---

## Acceptance Criteria

1. **AC-1: Shell Layout** — The shell renders with 3 regions (aside/sidebar, header/top bar, main/content section); each region displays placeholder content; full height layout fills viewport (`height: 100vh`)

2. **AC-2: Sidebar Component** — Fixed-width sidebar (260px desktop, 72px collapsed state); displays app logo, primary nav links (Home, Settings, Alerts, Reports placeholder items), collapsible nav section (e.g., "Tools" group), and footer action button; scrolls internally if needed

3. **AC-3: Top Bar Component** — Fixed-height sticky top bar (56px); layout uses `display: flex` with three regions:
   - **Left region:** Sidebar toggle button (visible on mobile), context breadcrumb/label
   - **Center region:** Search input (flexes to fill available space) with search icon inside
   - **Right region:** Quick action icons (optional placeholders for future use)

4. **AC-4: Main Content Area** — Flexible width content section with `overflow-y: auto` and 1rem/1.25rem padding; displays breadcrumb row, title with metadata, action row (primary + secondary buttons), and card/box content blocks; min-width: 0 prevents flex child overflow

5. **AC-5: Responsive Behavior** — Desktop (≥1024px): sidebar fixed/visible, search bar full width in center; Tablet/mobile (<1024px): sidebar becomes overlay drawer with toggle, search bar shrinks but remains centered; Small screens (<768px): tighter spacing, compact search with icon only (label hidden); Layout reflows correctly without horizontal scroll

6. **AC-6: Accessibility** — Semantic HTML regions (aside, header, main); aria-expanded on collapsible controls; search input has accessible label (visible or aria-label); keyboard support (Esc closes drawer/dropdowns, Tab navigates through top bar regions); visible focus states on all interactive elements; tested with keyboard navigation

7. **AC-7: Browser Verification** — `curl http://localhost:8080/` or `http://localhost:8080/shell` serves the shell page; browser loads without errors; Bulma CSS renders correctly; search bar appears centered in top bar; no JavaScript console errors

---

## Acceptance-Test Intent

User journey verification (browser-based):

1. **Load shell** → User opens browser to app URL → sees full shell layout with centered search in top bar (sidebar, top bar with search, content area)
2. **Search interaction** → User clicks in search input → can type without JavaScript errors; placeholder text visible
3. **Mobile sidebar toggle** → User clicks sidebar toggle on mobile → drawer slides in from left; search bar width adjusts but remains centered
4. **Quick actions** → User clicks quick action icon placeholders → no errors (placeholders for future features)
5. **Responsive resize** → User resizes browser window mobile → desktop → mobile → layout responds, search bar re-centers without breakage
6. **Keyboard Esc** → User presses Esc key → any open drawer/dropdown closes
7. **Keyboard Tab** → User tabs through interactive elements → can tab into search bar; focus visible; tab order logical (toggle → search → action icons)

---

## Out of Scope

- Search functionality backend (search API endpoint; that's a future increment)
- Profile management, login, logout (v2)
- Settings screen implementation (next increment after shell)
- Other feature screens (Home, Alerts, Reports)
- Backend API endpoints beyond serving shell HTML
- Persistent sidebar state (collapse preference stored)
- Animation/transitions
- Dark mode (v2)
- Drag-and-drop sidebar customization

---

## Constitution Constraints

- HTTP handler serves shell template (thin adapter, no business logic)
- All SPA assets embedded in binary via `embed.FS`
- Use Go `html/template` for HTML; HTMX for interactivity; Alpine.js for client state
- Bulma CSS framework for styling (no custom CSS beyond layout glue)
- Search input wired to Alpine state (no backend calls in this increment)
- All tests pass with `-race` flag
- No new external dependencies beyond Bulma (already in ADR-20260915-spa-frontend-stack.md)
- Accessibility requirements (WCAG 2.1 AA target)

---

## Roadmap Entry

**Feature:** Web App Shell (Foundation for SPA)

**Job Story:** When I load the web app, I want to see a professional layout with sidebar and centered search bar, so that I have a foundation for all screens

**Status:** Planned → Partial

---

## Next Steps

1. Write `docs/plan.md` with detailed technical execution plan
2. Create branch: `git checkout -b increment/web-app-shell`
3. Begin implementation following plan.md (TDD-Red phase or tidy phase depending on architecture)
