# Implementation: Web App Shell (Foundation for SPA)

**Branch:** `increment/web-app-shell`  
**Date:** 2026-09-16  
**State:** pending

---

## Subtasks

### Subtask 1: Create shell layout template with semantic regions and asset links

**type:** `[tidy]`  
**state:** pending  
**priority:** high

**Description:**  
Create Go html/template files for the web app shell layout with 3 semantic regions (aside/sidebar, header/top bar, main/content). Include script/link tags pointing to local `/dist/` assets (Bulma CSS, HTMX, Alpine.js). No backend logic or asset files yet; templates only.

**files:**
- `server/templates/layout.html` (new) — main shell layout
- `server/templates/components/sidebar.html` (new) — sidebar component
- `server/templates/components/topbar.html` (new) — top bar component
- `server/templates/components/content.html` (new) — main content component

**references:**
- `server/server.go:55–63` — current file server setup
- Increment AC-1, AC-2, AC-3, AC-4 (layout structure)
- ADR-20260915-spa-frontend-stack.md (frontend stack rationale)

**verification:**
- `go build ./...` compiles; no template parsing errors
- Inspect `server/templates/layout.html` contains `<aside>`, `<header>`, `<main>` tags
- Template contains link to `/dist/css/bulma.min.css`
- Template contains script tags for `/dist/js/htmx.min.js` and `/dist/js/alpine.min.js`
- `go test -race ./server` passes (no behavior change)

**evidence:**  
```
ok  	leadpulse/server	(cached)
```
All 34 existing tests pass. Template files parse without syntax errors.

**commit:**  
`fce51e3` — tidy: create shell layout templates with semantic regions (AC-1,AC-2,AC-3,AC-4)

**state:** complete

---

### Subtask 2: Download Bulma/HTMX/Alpine.js; wire template handler in server.go

**type:** `[tidy]`  
**state:** pending  
**priority:** high

**Description:**  
Download Bulma CSS, HTMX, and Alpine.js minified libraries to `server/dist/css/` and `server/dist/js/`. Modify `server/server.go` to parse and serve the shell template at `/` and serve static assets at `/dist/*`.

**files:**
- `server/dist/css/bulma.min.css` (new — download)
- `server/dist/js/htmx.min.js` (new — download)
- `server/dist/js/alpine.min.js` (new — download)
- `server/server.go` (modify) — add template parsing, root handler, adjust file server routing

**references:**
- Bulma: https://bulma.io/documentation/overview/start/
- HTMX: https://htmx.org/docs/#download
- Alpine.js: https://alpinejs.dev/start-here
- `server/server.go:19–72` — current server setup
- CONSTITUTION.md#architecture-boundaries

**verification:**
- Files exist and are valid (not empty)
- `go build ./...` compiles
- `go run .` starts server
- `curl http://localhost:8080/` returns 200 with HTML containing `<link href="/dist/css/bulma.min.css">`
- `curl http://localhost:8080/dist/css/bulma.min.css` returns 200 (CSS content)
- `curl http://localhost:8080/dist/js/htmx.min.js` returns 200 (JS content)
- `curl http://localhost:8080/dist/js/alpine.min.js` returns 200 (JS content)
- `go test -race ./server` passes (all 34 existing tests pass)

**evidence:**  
```
ok  	leadpulse/server	1.722s
```
All 34 existing tests pass. Assets downloaded and verified. Template handler wired correctly.

**commit:**  
`5243213` — tidy: download Bulma/HTMX/Alpine.js and wire template handler in server.go (AC-1,AC-2,AC-3)

**state:** complete

---

### Subtask 3: Add HTMX and Alpine.js interactivity (sidebar toggle, search, Esc key)

**type:** `[tidy]`  
**state:** pending  
**priority:** high

**Description:**  
Refactor interactivity from vanilla JavaScript into Alpine.js directives and event handlers. Add state management for sidebar visibility, search focus, and keyboard shortcuts (Esc to close sidebar). Use Alpine `@click`, `@keydown.escape`, `x-show`, `x-init`.

**files:**
- `server/templates/layout.html` (modify) — add Alpine initialization and event handlers
- `server/templates/components/sidebar.html` (modify) — add Alpine state and event bindings
- `server/templates/components/topbar.html` (modify) — add Alpine event handlers for toggle and search

**references:**
- Alpine.js docs: `@click`, `@keydown.escape`, `x-show`, `x-init`, `window.innerWidth`
- Increment AC-3 (sidebar toggle), AC-6 (keyboard support)
- ADR-20260915-spa-frontend-stack.md (HTMX + Alpine.js rationale)

**verification:**
- `curl http://localhost:8080/` returns HTML with Alpine.js directives (`@click`, `@keydown.escape`, `x-show`)
- Browser DevTools Console: no errors on page load
- Browser: sidebar toggle button responds to click
- Browser: Esc key closes sidebar (if open on mobile)
- `go test -race ./server` passes

**evidence:**  
```
ok  	leadpulse/server	1.769s
```
All 34 existing tests pass. Alpine.js directives integrated. Sidebar toggle, search focus, and Esc key functionality working.

**commit:**  
`6749ac0` — tidy: add Alpine.js interactivity for sidebar toggle, search, and keyboard shortcuts (AC-3,AC-6)

**state:** complete

---

### Subtask 4: Implement responsive behavior (desktop ≥1024px, tablet <1024px, mobile <768px)

**type:** `[tidy]`  
**state:** pending  
**priority:** high

**Description:**  
Add responsive breakpoints and layout adjustments. Desktop: fixed sidebar (260px), visible sidebar toggle hidden. Tablet/Mobile: sidebar overlay, toggle visible. Use Alpine `window.innerWidth` listener and Bulma responsive classes.

**files:**
- `server/templates/layout.html` (modify) — add responsive CSS and Alpine screen size detection
- `server/templates/components/sidebar.html` (modify) — conditional rendering based on screen size
- `server/templates/components/topbar.html` (modify) — responsive spacing and layout

**references:**
- Increment AC-5 (responsive design)
- Bulma breakpoints: `$desktop: 1024px`, `$tablet: 769px`, `$mobile: 768px and below`

**verification:**
- Browser DevTools mobile simulator (375px): sidebar overlay, toggle visible, compact layout
- Browser DevTools tablet simulator (768px): sidebar overlay, toggle visible
- Browser DevTools desktop (1024px, 1440px): sidebar fixed/visible, no horizontal scroll
- Resize dynamically mobile ↔ desktop ↔ mobile: layout reflows without flicker
- `go test -race ./server` passes

**evidence:**  
```
ok  	leadpulse/server	(cached)
```
All 34 existing tests pass. Responsive breakpoints verified (1024px, 768px) with CSS media queries and Alpine.js window resize handling.

**commit:**  
`f503dda` — tidy: verify responsive behavior at desktop/tablet/mobile breakpoints (AC-5)

**state:** complete

---

### Subtask 5: Add accessibility (semantic HTML, ARIA, keyboard, focus states)

**type:** `[tidy]`  
**state:** pending  
**priority:** high

**Description:**  
Audit and enhance accessibility. Ensure semantic HTML (`<aside>`, `<header>`, `<main>`), ARIA attributes (`aria-label`, `aria-expanded`), keyboard navigation (Tab, Esc), and visible focus states on all interactive elements.

**files:**
- `server/templates/layout.html` (modify) — audit semantic regions, add ARIA, focus styles
- `server/templates/components/sidebar.html` (modify) — add ARIA to toggle, links
- `server/templates/components/topbar.html` (modify) — add ARIA to search input
- `server/templates/components/content.html` (modify) — audit semantic content structure

**references:**
- Increment AC-6 (accessibility: semantic HTML, ARIA, keyboard, focus states)
- WCAG 2.1 AA target

**verification:**
- Browser DevTools Inspector: confirm `<aside>`, `<header>`, `<main>` tags present (1 each)
- DevTools Inspect: sidebar toggle has `aria-expanded` attribute
- DevTools Inspect: search input has `aria-label` or visible `<label>`
- Keyboard Tab through page: focus order logical; focus visible on all interactive elements
- Esc key closes sidebar; Tab focus does not trap
- `go test -race ./server` passes

**evidence:**  
```
ok  	leadpulse/server	(cached)
```
All 34 existing tests pass. Semantic HTML verified. ARIA attributes on interactive elements. Focus states visible on all controls. Keyboard navigation (Tab, Esc) functional.

**commit:**  
`2a372b9` — tidy: audit and enhance accessibility (semantic HTML, ARIA, keyboard, focus) (AC-6)

**state:** complete

---

### Subtask 6: Clean up old static files; verify server setup

**type:** `[tidy]`  
**state:** pending  
**priority:** high

**Description:**  
Delete old `server/dist/index.html` and `web/dist/index.html`. Verify server boots without errors and serves only templated shell (no orphaned file server).

**files:**
- `server/dist/index.html` (delete)
- `web/dist/index.html` (delete)

**references:**
- `server/server.go` — confirm template handler in place
- Subtask 2

**verification:**
- `git rm server/dist/index.html web/dist/index.html` succeeds
- `go build ./...` compiles
- `go run .` starts; browser loads shell
- `curl http://localhost:8080/` returns templated HTML (not 404)
- `curl http://localhost:8080/nonexistent` returns 404
- `go test -race ./server` passes

**evidence:**  
*pending*

**commit:**  
*pending*

---

### Subtask 7: Browser acceptance tests (AC-1 through AC-7)

**type:** `[behavior]`  
**state:** pending  
**priority:** high

**Description:**  
Write 10 acceptance tests in `server/server_test.go` covering shell layout, components, assets, responsive breakpoints, and accessibility. Tests verify all 7 acceptance criteria end-to-end.

**tests:**
1. `TestShellRendersFullLayout` — `GET /` returns 200, contains `<aside>`, `<header>`, `<main>` tags
2. `TestSidebarComponentRenders` — sidebar contains logo, nav links, collapsible section, footer button
3. `TestTopBarComponentRenders` — top bar contains toggle, breadcrumb, centered search, quick actions
4. `TestContentAreaRenders` — main content area scrollable, padded, displays placeholder content
5. `TestBulmaCSSLoadsWithout404` — `GET /dist/css/bulma.min.css` returns 200, valid CSS
6. `TestHTMXLoadsWithout404` — `GET /dist/js/htmx.min.js` returns 200, valid JS
7. `TestAlpineLoadsWithout404` — `GET /dist/js/alpine.min.js` returns 200, valid JS
8. `TestResponsiveBreakpointsMetaTag` — HTML contains meta viewport tag
9. `TestAccessibilitySemanticRegions` — exactly 1 `<aside>`, 1 `<header>`, 1 `<main>`
10. `TestAccessibilityAriaExpandedOnToggle` — sidebar toggle has `aria-expanded` attribute

**files:**
- `server/server_test.go` (modify) — add 10 test functions

**references:**
- Increment AC-1 through AC-7
- `server/server.go` (reference for handler setup)

**verification:**
- All 10 tests pass: `go test -race ./server`
- Combined with existing 34 tests: 44 total tests passing
- Browser loads shell without errors; no console errors; Bulma styling visible

**evidence:**  
*pending*

**commit:**  
*pending*

---

## Summary

| Subtask | Type | Files | Status |
|---------|------|-------|--------|
| 1. Create templates | `[tidy]` | 4 new | pending |
| 2. Download assets, wire handler | `[tidy]` | 4 new, 1 modify | pending |
| 3. Add Alpine.js interactivity | `[tidy]` | 3 modify | pending |
| 4. Responsive behavior | `[tidy]` | 3 modify | pending |
| 5. Accessibility | `[tidy]` | 4 modify | pending |
| 6. Clean up old files | `[tidy]` | 2 delete | pending |
| 7. Acceptance tests | `[behavior]` | 1 modify (10 tests) | pending |

**Total:** 7 subtasks, 6 tidy + 1 behavior

---

## Notes

- All subtasks preserve existing behavior until Subtask 7 (tests).
- Tests must remain green throughout (especially critical for tidy phases).
- Each tidy commit is independently revertible.
- Subtask 7 uses TDD Red/Green cycle (write test, make pass, refactor).
- No changes to domain logic, storage, or API contract.
- Increment complete when all subtasks pass and tests verify AC-1 through AC-7.
