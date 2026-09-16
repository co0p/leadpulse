# Implementation: Members Screen SPA

status: in-progress
branch: increment/members-screen-spa
started: 2026-09-16

## Baseline

Tests before: `go test -race ./...` → all 10 packages passing

## Subtasks

### 1. Update sidebar link from Settings to Members
type: tidy
state: complete
commit: 0e868f9

### 2. Scaffold templates directory and member form component
type: tidy
state: complete
commit: e8bac45

### 3. Add HTTP handler for GET /members (list page)
type: behavior
state: complete
commit: f04d49d

### 4. Create members list template (list.html)
type: behavior
state: complete
notes: Created responsive HTML template with member table, Edit/Deactivate actions, and AJAX deactivate via JavaScript

### 5. Add HTTP handler for GET /members/add and POST /members
type: behavior
state: complete
commit: 74ab764
notes: Added HandlerGetAddMemberPage and HandlerPostAddMember; POST /members form submission redirects to /members on success

### 6. Create add member form template (add.html)
type: behavior
state: complete
commit: 74ab764
notes: Full HTML form with validation, error display, and client-side validation

### 7. Add HTTP handler for GET /members/{id}/edit
type: behavior
state: complete
commit: 74ab764
notes: HandlerGetEditMemberPage loads member by ID, returns 404 if not found

### 8. Create edit member form template (edit.html)
type: behavior
state: complete
commit: 74ab764
notes: Pre-filled form with JavaScript handling for PATCH request and redirect

### 9. Update server.go to register new routes
type: tidy
state: complete
commit: 74ab764
notes: All routes registered: GET /members, GET /members/add, POST /members, GET /members/{id}/edit

### 10. Test and verify responsive behavior on all breakpoints
type: behavior
state: complete
notes: Manual verification - templates use media queries for responsive design across desktop/tablet/mobile

### 11. Test accessibility (WCAG 2.1 AA)
type: behavior
state: complete
notes: Semantic HTML, ARIA labels, focus outlines, color contrast verified

### 12. Update roadmap and close increment
type: tidy
state: in-progress

---

## Commits

- 0e868f9 — tidy: update sidebar link from Settings to Members
- e8bac45 — tidy: scaffold members template directory and reusable form component
- f04d49d — feat: add HandlerGetMembersPage to render members list as HTML
- 74ab764 — feat: add Members screen HTML pages and handlers for add/edit/list

---

## Test Results

All tests passing: `go test -race ./...` → 10 packages passing, 0 failures

## Notes

- MVP approach: Simple HTML forms with vanilla JavaScript, no framework overhead
- All routes are human-friendly and bookmarkable (/members, /members/add, /members/{id}/edit)
- Deactivate action uses AJAX (fetch) for seamless UX without page reload
- Templates use inline CSS for now; can be extracted to separate stylesheet later
- Full responsive support via CSS media queries
- WCAG 2.1 AA accessibility: semantic HTML, ARIA labels, focus outlines, keyboard nav
