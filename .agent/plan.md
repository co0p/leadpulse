# Plan: Members Screen SPA

## Goal

Build the Members screen as a browser-rendered SPA page with human-friendly URLs, backed by the existing Members API, replacing the Fyne Settings screen for member management.

## Branch

`increment/members-screen-spa`

## Approach

The Members screen introduces server-side rendered (SSR) HTML pages via Go templates. Each route (`/members`, `/members/add`, `/members/{id}/edit`) serves a distinct HTML page embedded in the binary. The HTTP handlers remain thin adapters — they load data via existing use cases, pass the data to Go templates, and render the response. Form submissions POST to the existing `/api/members` JSON endpoints; on success, the response triggers a redirect via HTML or client-side JavaScript. No new business logic; all validation and persistence already in core use cases. Performance remains sub-200ms (per CONSTITUTION envelope) because queries and use case execution are unchanged.

## Design

### Data Models

No new domain models. The Members screen renders existing `TeamMember` aggregate from `core/members/member.go`.

**HTTP Request/Response additions:**
- `members/get_members.go` already returns `GetMembersOutput` (list of active members)
- `members/add_member.go` already returns `AddMemberOutput` (created member with ID)
- `members/edit_member.go` already returns `EditMemberOutput` (updated member)
- `members/deactivate_member.go` already returns `DeactivateMemberOutput` (success or error)

**Template data models:**
```go
// Used in Go templates for rendering
type MembersListPageData struct {
  Members []MemberForDisplay  // Active members only
}

type MemberForDisplay struct {
  ID        string  // UUID (for href)
  FirstName string
  LastName  string
  Seniority string
}

type AddMemberPageData struct {
  // Empty on GET; used to render form
}

type EditMemberPageData struct {
  Member MemberForDisplay
}
```

These are template-local types, not domain objects; they reside in `server/templates/` and are not persisted.

### Call / Data Flow

**GET /members (list):**
1. HTTP handler receives GET request
2. Calls `GetMembersUC.Execute()` → returns `GetMembersOutput` (list of active members)
3. Maps `GetMembersOutput` members to template data (convert int64 ID to UUID string)
4. Executes `server/templates/members/list.html` with data
5. Returns rendered HTML

**GET /members/add (add form):**
1. HTTP handler receives GET request
2. Executes `server/templates/members/add.html` with empty data (form only, no API call needed)
3. Returns rendered HTML

**POST /members (form submit for add):**
1. Form submits to POST /members with form-encoded body: `firstName`, `lastName`, `seniority`
2. Handler parses form data
3. Calls existing `AddMemberUC.Execute(input)` → returns `AddMemberOutput` or error
4. If error: re-render `add.html` with error message and form values (client-side validation before submit prevents most errors)
5. If success: redirect to `/members` (HTTP 302 or via JavaScript)
6. Browser navigates to `/members`, sees member in list

**GET /members/{id}/edit (edit form):**
1. Handler extracts member ID from URL path
2. Calls `GetMembersUC.Execute()` to get all active members
3. Finds member with matching ID in response
4. Maps to template data
5. Executes `server/templates/members/edit.html` with member data
6. Returns rendered HTML (form pre-populated with current values)

**PATCH /api/members/{id} (form submit for edit):**
1. Form submits to PATCH /api/members/{id} with JSON body: `firstName`, `lastName`, `seniority`
2. Existing handler processes request, calls `EditMemberUC.Execute(input)` → returns `EditMemberOutput` or error
3. If error: returns JSON error (client-side validation prevents; rare)
4. If success: returns JSON success response
5. Client-side JavaScript receives response, redirects to `/members`

**DELETE /api/members/{id} (deactivate):**
1. List page has "Deactivate" button on each row
2. Button calls DELETE /api/members/{id} via fetch() or HTMX
3. Existing handler processes request, calls `DeactivateMemberUC.Execute(input)` → returns `DeactivateMemberOutput` or error
4. If success: returns JSON success; client-side removes row from table (no page reload)
5. If error: returns JSON error; client-side displays toast/alert

### Error / Edge-case Inventory

| Condition | Expected response | Covered by |
|-----------|-------------------|------------|
| Member name is empty | Validation error on submit; re-render form with error and form values | Template + handler (AC-5) |
| Member name > 100 chars | Validation error on submit; re-render form with error | Template + handler (AC-5) |
| Seniority field is empty | Validation error on submit; re-render form with error | Template + handler (AC-5) |
| Invalid seniority value | Use case validates; handler returns 400 JSON error | Existing use case (AC-5) |
| Member ID in URL is invalid (not a UUID) | Handler logs error, returns 404 HTML page | Subtask 3 — handler validation |
| Member not found in edit flow | Handler displays "Member not found" message | Subtask 3 — handler validation |
| Database error on fetch members | Existing use case returns error; handler returns 500 HTML page | Existing use case |
| Deactivate a member that is already deactivated | API returns error; client-side displays error toast | Existing use case (AC-4) |
| User tries to edit a deactivated member | Member not in list; if ID guessed in URL, handler returns 404 (deactivated members not returned by GetMembers) | Subtask 3 — handler validation |
| Form double-submit (user clicks button twice before redirect) | Client-side submit button disables on click (JavaScript); second click has no effect | Subtask 2 — client-side interactivity (AC-2, AC-3) |

### Observability Intent

| Event | Level | Data |
|-------|-------|------|
| `handler.get_members_list` | debug | count_members, render_time_ms |
| `handler.add_member_form_displayed` | debug | — (user viewed form) |
| `handler.add_member_submitted` | info | status (success\|error), error_reason |
| `handler.edit_member_form_displayed` | debug | member_id |
| `handler.edit_member_submitted` | info | member_id, status, error_reason |
| `handler.deactivate_member` | info | member_id, status, error_reason |

Logging via Go stdlib `log.Printf` or structured logger (e.g., `slog`). No external telemetry service.

### Architecture Delta

**No container changes.** The Members screen remains within the `server/` (HTTP adapter) container. Templates are new artifacts but not a separate deployable. The `core/members/` use cases are unchanged and continue to be the API contract.

**New data flow:** Previously, only JSON API existed (`/api/members`). Now:
- GET `/members` renders HTML (new)
- GET `/members/add` renders form HTML (new)
- POST `/members` accepts form data, redirects to list (new; routes to same use case as POST `/api/members` but renders HTML instead of JSON)
- GET `/members/{id}/edit` renders form HTML with member data (new)
- Existing PATCH `/api/members/{id}` and DELETE `/api/members/{id}` unchanged

**No reverse dependency:** Server does not import new code from core or storage. Core and storage remain unaware of template rendering.

## Files

| File | Role | Notes |
|------|------|-------|
| `server/server.go` | modify | register new routes: GET `/members`, GET `/members/add`, POST `/members`, GET `/members/{id}/edit` |
| `server/handler_members.go` | modify | add `HandlerGetMembersPage`, `HandlerAddMemberPage`, `HandlerPostAddMember`, `HandlerEditMemberPage` (new functions; existing API handlers unchanged) |
| `server/handler_members_test.go` | modify | add tests for new handlers |
| `server/templates/members/list.html` | new | members list page with table and action buttons |
| `server/templates/members/add.html` | new | add member form page |
| `server/templates/members/edit.html` | new | edit member form page |
| `server/templates/components/member_form.html` | new | reusable form fields (firstName, lastName, seniority) to avoid duplication between add.html and edit.html |
| `server/templates/layout.html` | modify | update sidebar link from `/settings` to `/members` |
| `docs/ui.md` | touch | reference for accessibility, responsive breakpoints, form patterns |
| `docs/architecture.md` | touch | reference for HTTP → use case → repository flow |

## Subtasks

### 1. Update sidebar link from Settings to Members
**type:** [tidy]

**description:** Rename the "Settings" sidebar link to "Members" and point to `/members` instead of `/settings`.

**files:**
- `server/templates/layout.html`

**references:**
- `server/templates/layout.html:322` — current Settings link

**verification:** Manual: open browser, verify sidebar shows "Members" link pointing to `/members`

**depends on:** —

**acceptance criteria:** — (structural prep)

---

### 2. Scaffold templates directory and member form component
**type:** [tidy]

**description:** Create the template subdirectories and a reusable `member_form.html` component containing the form fields (firstName, lastName, seniority select). This component will be included in both `add.html` and `edit.html` to avoid duplication.

**files:**
- `server/templates/members/` (new directory)
- `server/templates/components/member_form.html` (new file — form field component)

**references:**
- `server/templates/layout.html` — existing template structure for reference

**verification:** Directory structure created; component file contains form fields without form wrapping (form wrapping added by add.html and edit.html)

**depends on:** —

**acceptance criteria:** — (structural prep)

---

### 3. Add HTTP handler for GET /members (list page)
**type:** [behavior]

**description:** Create `HandlerGetMembersPage` in `server/handler_members.go` that calls `GetMembersUC.Execute()`, maps results to template data, executes `members/list.html` template, and returns rendered HTML. Handle error cases (database error → 500 page).

**files:**
- `server/handler_members.go` (add function)
- `server/handler_members_test.go` (add tests)

**references:**
- `server/handler_members.go:110` — existing `HandlerAddMember` for pattern
- `core/members/get_members.go` — use case interface to understand output
- `docs/ui.md#shell-layout` — responsive table design guidance

**tests:**
- id: handler-get-members-success
  file: `server/handler_members_test.go`
  name: `TestHandlerGetMembersPage_ReturnsHTMLWithMembers`
  state: pending
- id: handler-get-members-empty
  file: `server/handler_members_test.go`
  name: `TestHandlerGetMembersPage_EmptyListWorks`
  state: pending
- id: handler-get-members-error
  file: `server/handler_members_test.go`
  name: `TestHandlerGetMembersPage_ReturnsErrorOn500`
  state: pending

**active_test:** handler-get-members-success

**verification:** `go test -run TestHandlerGetMembersPage` — all tests pass; handler is testable with mocked use case (no database)

**depends on:** 1, 2

**acceptance criteria:** AC-1 (list renders with member data)

---

### 4. Create members list template (list.html)
**type:** [behavior]

**description:** Create `server/templates/members/list.html` — a responsive HTML page showing all active members in a table. Include columns: Name, Seniority, Actions. Actions row has "Edit" link (→ `/members/{id}/edit`) and "Deactivate" button (→ DELETE `/api/members/{id}`). Deactivate button uses fetch() to call API; on success, removes row from DOM and hides it with fade-out animation (Alpine.js or vanilla JS).

**files:**
- `server/templates/members/list.html` (new)

**references:**
- `docs/ui.md#touch-targets` — button min size 44px
- `docs/ui.md#responsive-breakpoints` — table should stack on mobile
- `server/templates/layout.html:159` — main-content area structure for reference

**tests:**
- id: list-template-renders
  file: `server/handler_members_test.go`
  name: `TestMembersListTemplate_RendersWithoutError`
  state: pending
- id: list-shows-member-names
  file: `server/handler_members_test.go`
  name: `TestMembersListTemplate_DisplaysMemberNamesAndSeniority`
  state: pending
- id: list-edit-link-correct
  file: `server/handler_members_test.go`
  name: `TestMembersListTemplate_EditLinkPointsToCorrectMember`
  state: pending
- id: list-deactivate-button-exists
  file: `server/handler_members_test.go`
  name: `TestMembersListTemplate_DeactivateButtonPresent`
  state: pending

**active_test:** list-template-renders

**verification:** `go test -run TestMembersListTemplate` — all tests pass; template renders HTML with no errors; links and buttons present

**depends on:** 3

**acceptance criteria:** AC-1 (list displays members, edit/deactivate buttons present)

---

### 5. Add HTTP handler for GET /members/add and POST /members
**type:** [behavior]

**description:** Create two handlers:
- `HandlerGetAddMemberPage` — GET /members/add — calls no use case, just renders empty `add.html` form
- `HandlerPostAddMember` — POST /members — parses form data (form-encoded), calls `AddMemberUC.Execute()`, on error re-renders form with error message, on success redirects to `/members`

**files:**
- `server/handler_members.go` (add two functions)
- `server/handler_members_test.go` (add tests)

**references:**
- `server/handler_members.go:58` — existing `HandlerAddMember` (JSON version) for reference; reuse validation logic
- `core/members/add_member.go` — use case interface
- `docs/ui.md#responsive-breakpoints` — form should stack on mobile

**tests:**
- id: get-add-form-success
  file: `server/handler_members_test.go`
  name: `TestHandlerGetAddMemberPage_ReturnsEmptyForm`
  state: pending
- id: post-add-form-valid
  file: `server/handler_members_test.go`
  name: `TestHandlerPostAddMember_ValidFormRedirectsToList`
  state: pending
- id: post-add-form-invalid
  file: `server/handler_members_test.go`
  name: `TestHandlerPostAddMember_InvalidFormRerenderWithError`
  state: pending
- id: post-add-form-missing-field
  file: `server/handler_members_test.go`
  name: `TestHandlerPostAddMember_MissingFieldShowsError`
  state: pending

**active_test:** get-add-form-success

**verification:** `go test -run TestHandlerGetAddMemberPage -run TestHandlerPostAddMember` — all tests pass

**depends on:** 2, 3

**acceptance criteria:** AC-2 (add form renders, submission persists and redirects)

---

### 6. Create add member form template (add.html)
**type:** [behavior]

**description:** Create `server/templates/members/add.html` — a form page with:
- Title: "Add Member"
- Form fields: firstName, lastName, seniority (reused from `member_form.html`)
- Submit button ("Add Member")
- Cancel link (back to `/members`)
- Error message display (if POST returns error, form re-renders with error at top)
- Client-side form validation (required fields, name length < 100) before submit

**files:**
- `server/templates/members/add.html` (new)

**references:**
- `server/templates/components/member_form.html` — included for form fields
- `docs/ui.md#touch-targets` — button sizes
- `CONSTITUTION.md#testing-strategy` — all validation tested before reaching server

**tests:**
- id: add-form-renders
  file: `server/handler_members_test.go`
  name: `TestAddMemberTemplate_RendersWithoutError`
  state: pending
- id: add-form-has-fields
  file: `server/handler_members_test.go`
  name: `TestAddMemberTemplate_HasAllFormFields`
  state: pending
- id: add-form-shows-error
  file: `server/handler_members_test.go`
  name: `TestAddMemberTemplate_DisplaysErrorMessageWhenProvided`
  state: pending

**active_test:** add-form-renders

**verification:** `go test -run TestAddMemberTemplate` — all tests pass; template has form fields, submit button, error display

**depends on:** 5

**acceptance criteria:** AC-2 (form renders, validates, and submits)

---

### 7. Add HTTP handler for GET /members/{id}/edit
**type:** [behavior]

**description:** Create `HandlerGetEditMemberPage` in `server/handler_members.go` that:
- Extracts member ID from URL path (convert string to int64, then to UUID for lookup)
- Calls `GetMembersUC.Execute()` to get all active members
- Finds member with matching ID in response
- If member not found (deactivated or invalid ID): render 404 page
- If found: maps to template data, executes `edit.html` with member data pre-filled

**files:**
- `server/handler_members.go` (add function)
- `server/handler_members_test.go` (add tests)

**references:**
- `server/handler_members.go:132` — existing `HandlerEditMember` (JSON version) for ID parsing logic
- `server/server.go:33` — pattern for extracting path parameter
- `core/members/get_members.go` — use case returns all active members (search within result)

**tests:**
- id: get-edit-form-success
  file: `server/handler_members_test.go`
  name: `TestHandlerGetEditMemberPage_LoadsAndRendersForm`
  state: pending
- id: get-edit-form-not-found
  file: `server/handler_members_test.go`
  name: `TestHandlerGetEditMemberPage_Returns404ForInvalidID`
  state: pending
- id: get-edit-form-deactivated
  file: `server/handler_members_test.go`
  name: `TestHandlerGetEditMemberPage_Returns404ForDeactivatedMember`
  state: pending

**active_test:** get-edit-form-success

**verification:** `go test -run TestHandlerGetEditMemberPage` — all tests pass

**depends on:** 2, 3

**acceptance criteria:** AC-3 (edit form loads and pre-fills member data)

---

### 8. Create edit member form template (edit.html)
**type:** [behavior]

**description:** Create `server/templates/members/edit.html` — similar to `add.html` but:
- Title: "Edit Member"
- Form fields: firstName, lastName, seniority (reused from `member_form.html`, pre-filled with current values)
- Submit button ("Save Member")
- Cancel link (back to `/members`)
- Form method: PATCH (to `/api/members/{id}`)
- Error message display (if save fails, form re-renders with error)
- Client-side validation (same as add.html)

**files:**
- `server/templates/members/edit.html` (new)

**references:**
- `server/templates/components/member_form.html` — included for form fields
- `server/templates/members/add.html` — use as pattern
- `docs/ui.md#touch-targets` — button sizes

**tests:**
- id: edit-form-renders
  file: `server/handler_members_test.go`
  name: `TestEditMemberTemplate_RendersWithoutError`
  state: pending
- id: edit-form-prefills
  file: `server/handler_members_test.go`
  name: `TestEditMemberTemplate_PrefillsCurrentValues`
  state: pending
- id: edit-form-shows-error
  file: `server/handler_members_test.go`
  name: `TestEditMemberTemplate_DisplaysErrorMessageWhenProvided`
  state: pending

**active_test:** edit-form-renders

**verification:** `go test -run TestEditMemberTemplate` — all tests pass

**depends on:** 7

**acceptance criteria:** AC-3 (edit form renders with member data, can submit update)

---

### 9. Update server.go to register new routes
**type:** [tidy]

**description:** In `server/Start()`, register the new routes:
- GET `/members` → `HandlerGetMembersPage`
- GET `/members/add` → `HandlerGetAddMemberPage`
- POST `/members` → `HandlerPostAddMember`
- GET `/members/{id}/edit` → `HandlerGetEditMemberPage`

Keep existing JSON API routes unchanged (`POST /api/members`, `GET /api/members`, `PATCH /api/members/{id}`, `DELETE /api/members/{id}`).

**files:**
- `server/server.go`

**references:**
- `server/server.go:20` — Start() function signature
- `server/server.go:23` — existing API route registration pattern

**verification:** `go build ./...` — builds without errors; routes registered in mux

**depends on:** 5, 7

**acceptance criteria:** — (wiring)

---

### 10. Test and verify responsive behavior on all breakpoints
**type:** [behavior]

**description:** Manually verify (or via Playwright acceptance tests) that the Members screen (list, add, edit) renders correctly on:
- Desktop (≥1024px)
- Tablet (769–1023px)
- Mobile (≤768px)

Verify:
- List: table stacks to single column on mobile
- Forms: inputs full-width on mobile, condensed padding
- Buttons: min 44px touch targets on all sizes
- Sidebar: visible on desktop, overlay on mobile
- All text readable, no horizontal scroll

**files:**
- `server/handler_members_test.go` (add acceptance test)
- Manual browser verification

**references:**
- `docs/ui.md#responsive-breakpoints`
- `docs/ui.md#touch-targets`

**tests:**
- id: responsive-desktop
  file: manual/playwright
  name: `test Members screen on desktop (≥1024px)`
  state: pending
- id: responsive-tablet
  file: manual/playwright
  name: `test Members screen on tablet (769–1023px)`
  state: pending
- id: responsive-mobile
  file: manual/playwright
  name: `test Members screen on mobile (≤768px)`
  state: pending

**active_test:** responsive-desktop

**verification:** Open `/members`, `/members/add`, `/members/{id}/edit` in browser at three viewport sizes; verify layout adjusts correctly, no horizontal scroll, touch targets ≥44px

**depends on:** 4, 6, 8

**acceptance criteria:** AC-7 (responsive on all breakpoints)

---

### 11. Test accessibility (WCAG 2.1 AA)
**type:** [behavior]

**description:** Verify accessibility of the Members screen:
- Semantic HTML: `<form>`, `<label>`, `<table>`, heading hierarchy (H1, H2, H3)
- ARIA labels on form inputs, buttons
- Keyboard navigation: Tab through all interactive elements, enter submits forms
- Focus outlines: 2px solid #3273dc, visible on all buttons and inputs
- Color contrast: 4.5:1 for normal text (verified via browser dev tools)
- Screen reader: test with NVDA or macOS VoiceOver (quick spot-check)

**files:**
- Manual accessibility audit (no code changes; documentation in implementation.md)

**references:**
- `docs/ui.md#accessibility`
- `CONSTITUTION.md` — no explicit accessibility testing requirement, but `docs/ui.md` documents WCAG 2.1 AA target

**tests:**
- id: a11y-semantic-html
  file: manual
  name: `Verify semantic HTML (form, label, table, headings)`
  state: pending
- id: a11y-keyboard-nav
  file: manual
  name: `Verify Tab, Enter, and Esc keys work as expected`
  state: pending
- id: a11y-focus-outline
  file: manual
  name: `Verify 2px focus outline on all interactive elements`
  state: pending
- id: a11y-color-contrast
  file: manual
  name: `Verify text contrast ≥4.5:1`
  state: pending

**active_test:** a11y-semantic-html

**verification:** Manual browser inspection + quick screen reader test (read form fields aloud)

**depends on:** 6, 8

**acceptance criteria:** AC-8 (accessibility WCAG 2.1 AA)

---

### 12. Update roadmap and close increment
**type:** [tidy]

**description:** Update `docs/roadmap.md`:
- Move "Members Screen (SPA)" from Planned to Done
- Add acceptance criteria met, evidence links (test files), and key commits
- Mark Fyne Settings screen as deprecated (note when retiring)

**files:**
- `docs/roadmap.md`

**references:**
- `docs/roadmap.md#planned` — current section
- `docs/roadmap.md#done` — where entry should go

**verification:** Manual: read roadmap, verify entry in Done section with evidence links

**depends on:** 11 (all functionality complete)

**acceptance criteria:** — (documentation)

---

## Context Map

- `CONSTITUTION.md#engineering-principles` — behavior-first, small focused changes, dependency direction
- `CONSTITUTION.md#performance-envelope` — screen render ≤ 200ms target
- `CONSTITUTION.md#testing-strategy` — all layers tested; browser tests optional in v1
- `docs/ui.md` — shell layout, responsive breakpoints, accessibility, interaction patterns
- `docs/architecture.md#containers` — HTTP adapter pattern, use case delegation, no business logic in handlers
- `docs/architecture.md#communication-paths` — browser → HTTP handler → use case → storage
- `docs/roadmap.md` — prior screens (Members API, Web App Shell), SPA increment sequencing

## Acceptance Scenarios

These scenarios are advisory (not release-blocking per CONSTITUTION v1) but guide manual testing and Playwright E2E test writing.

### AT-1: View all members

- criterion: AC-1
- precondition: 3 active members exist in database
- user action: navigate to `/members` in browser
- expected outcome: list page renders with all 3 members (name, seniority), each with Edit and Deactivate buttons
- evidence: screenshot or Playwright test
- gate: advisory
- state: planned

### AT-2: Add a new member

- criterion: AC-2
- precondition: user is on `/members`
- user action: click "Add Member" button, fill form (first name, last name, seniority), click Save
- expected outcome: form submits, page redirects to `/members`, new member appears in list
- evidence: Playwright test in `tests/acceptance/members.spec.ts`
- gate: advisory
- state: planned

### AT-3: Edit a member

- criterion: AC-3
- precondition: member exists in list
- user action: click Edit on a member, change seniority, click Save
- expected outcome: form submits, page redirects to `/members`, member's seniority updated in list
- evidence: Playwright test
- gate: advisory
- state: planned

### AT-4: Deactivate a member

- criterion: AC-4
- precondition: member is visible in list
- user action: click Deactivate button on a member
- expected outcome: row fades out and disappears from list (no page reload), database confirms member is deactivated
- evidence: Playwright test or manual browser verification
- gate: advisory
- state: planned

### AT-5: Persistence across reload

- criterion: AC-6
- precondition: members have been added/edited/deactivated in the session
- user action: press F5 to reload the page
- expected outcome: changes persist; list shows the same members as before reload
- evidence: manual browser verification
- gate: advisory
- state: planned

## Risks

1. **Template rendering performance:** Go's `html/template` parsing and executing on every request could exceed 200ms performance envelope if templates are large or members list is huge (100+). Mitigation: measure latency; if > 100ms, consider template caching or precompilation.

2. **Form double-submit:** User clicks Submit twice before page redirects. Mitigation: client-side button disable on submit (JavaScript in form template).

3. **Stale member data on edit:** Between GET `/members/{id}/edit` and PATCH `/api/members/{id}`, another user edits the same member. Our form overwrites their change (lost update). Mitigation: per CONSTITUTION "human review always" — no optimistic locking in v1; document as v2 feature.

4. **Deactivated member race condition:** User views `/members` list, another user deactivates a member, current user tries to edit that member (ID in URL). Mitigation: handler validates member is in active list; if not, returns 404.

5. **Mobile form usability:** Input fields on mobile may be small or keyboard-hidden. Mitigation: test on iOS/Android; use `autofocus` sparingly; ensure 44px min-height buttons.

## Planning Decisions

| Decision | Chosen | Rejected | Reason |
|----------|--------|----------|--------|
| **URL structure** | Routable: `/members`, `/members/add`, `/members/{id}/edit` | Modal overlays | User requested surfable, bookmarkable URLs. Routable URLs feel like a "real" web app and support browser back button. |
| **Form rendering** | Server-side rendered (SSR) Go templates embedded in binary | Client-side SPA (React/Vue) | Go templates are simple, zero JS build step, all code in Go, files embedded in binary. SPA adds complexity (npm, build, bundle size). SSR sufficient for member CRUD. |
| **Form submission** | POST to `/members` (HTML form) + PATCH to `/api/members/{id}` (JSON) | Single form method | HTML form submission is semantically correct for SSR; existing JSON API used for edit to leverage already-built error handling. Hybrid approach minimizes code duplication. |
| **List deactivate** | AJAX delete (fetch) without page reload | Full page reload on delete | No page reload improves UX (faster, no flicker). Fetch + DOM removal is simple with vanilla JavaScript or HTMX. |
| **Deactivated member data in list** | Exclude from `GetMembersUC` (already filters to active members) | Include with "deactivated" badge | GetMembers already returns active only (per CONSTITUTION: "deactivate = soft delete, hidden from list"). Simplifies template; no conditional rendering needed. |
| **Accessibility approach** | Manual verification + WCAG 2.1 AA target per docs/ui.md | Automated axe-core tests | CONSTITUTION v1 does not require automated a11y testing. Manual verification plus semantic HTML, ARIA labels, and keyboard nav per `docs/ui.md` is sufficient. Axe-core deferred to v2. |
| **Error handling in forms** | Re-render form with error message and form values pre-filled | Redirect to error page | Pre-filled re-render allows user to fix single field and resubmit (UX-friendly). Existing use case validation already structured for this. |
| **Response after add/edit** | Redirect to `/members` (show success in context) | Redirect to detail page (show member solo) | Redirect to list allows user to confirm member appears and see context (other team members). Simpler navigation flow. Detail page deferred to future increment. |

---

**Next step:** Wait for approval of plan.md. On approval, load the `4dc-implement` skill to scaffold implementation.md and populate the task queue.
