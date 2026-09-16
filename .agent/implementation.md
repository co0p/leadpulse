# Implementation: Members Screen SPA

status: in-progress
branch: increment/members-screen-spa
started: 2026-09-16

## Baseline

Tests before: baseline to be established by first implement skill (tidy-refactor)

## Subtasks

### 1. Update sidebar link from Settings to Members
type: tidy
state: complete
files:
  - server/templates/layout.html
evidence: |
  All tests passing: `go test -race ./...` → 10 packages, 0 failures
  Structural change: sidebar link renamed "Settings" → "Members", href="/settings" → href="/members"
  Test updated: server_test.go line 210, check for "Members" instead of "Settings"
  No behavior change; template rendering unchanged.
commit: 0e868f9 — tidy: update sidebar link from Settings to Members

### 2. Scaffold templates directory and member form component
type: tidy
state: complete
files:
  - server/templates/members/ (new directory)
  - server/templates/components/member_form.html (new file)
evidence: |
  All tests passing: `go test -race ./...` → 10 packages, 0 failures
  Created server/templates/members/ directory (parent for list.html, add.html, edit.html)
  Created server/templates/components/member_form.html with reusable form fields (firstName, lastName, seniority select)
  Form component includes ARIA labels, client-side validation (required, maxlength), and conditional pre-fill for edit mode
  No behavior change; no logic changes.
commit: e8bac45 — tidy: scaffold members template directory and reusable form component

### 3. Add HTTP handler for GET /members (list page)
type: behavior
state: in_progress
files:
  - server/handler_members.go
  - server/handler_members_test.go
tests:
  - id: handler-get-members-success
    name: TestHandlerGetMembersPage_ReturnsHTMLWithMembers
    file: server/handler_members_test.go
    state: red
  - id: handler-get-members-empty
    name: TestHandlerGetMembersPage_EmptyListWorks
    file: server/handler_members_test.go
    state: pending
  - id: handler-get-members-error
    name: TestHandlerGetMembersPage_ReturnsErrorOn500
    file: server/handler_members_test.go
    state: pending
active_test: handler-get-members-success
test_evidence: |
  Test: server/handler_members_test.go:379 — TestHandlerGetMembersPage_ReturnsHTMLWithMembers
  Failure reason: undefined: HandlerGetMembersPage
  Expected behavior: handler calls GetMembersUC.Execute(), renders HTML template with member data
  Status: RED (fails for the right reason — function doesn't exist yet)

### 4. Create members list template (list.html)
type: behavior
state: pending
files:
  - server/templates/members/list.html
tests:
  - id: list-template-renders
    name: TestMembersListTemplate_RendersWithoutError
    file: server/handler_members_test.go
    state: pending
  - id: list-shows-member-names
    name: TestMembersListTemplate_DisplaysMemberNamesAndSeniority
    file: server/handler_members_test.go
    state: pending
  - id: list-edit-link-correct
    name: TestMembersListTemplate_EditLinkPointsToCorrectMember
    file: server/handler_members_test.go
    state: pending
  - id: list-deactivate-button-exists
    name: TestMembersListTemplate_DeactivateButtonPresent
    file: server/handler_members_test.go
    state: pending
active_test: list-template-renders

### 5. Add HTTP handler for GET /members/add and POST /members
type: behavior
state: pending
files:
  - server/handler_members.go
  - server/handler_members_test.go
tests:
  - id: get-add-form-success
    name: TestHandlerGetAddMemberPage_ReturnsEmptyForm
    file: server/handler_members_test.go
    state: pending
  - id: post-add-form-valid
    name: TestHandlerPostAddMember_ValidFormRedirectsToList
    file: server/handler_members_test.go
    state: pending
  - id: post-add-form-invalid
    name: TestHandlerPostAddMember_InvalidFormRerenderWithError
    file: server/handler_members_test.go
    state: pending
  - id: post-add-form-missing-field
    name: TestHandlerPostAddMember_MissingFieldShowsError
    file: server/handler_members_test.go
    state: pending
active_test: get-add-form-success

### 6. Create add member form template (add.html)
type: behavior
state: pending
files:
  - server/templates/members/add.html
tests:
  - id: add-form-renders
    name: TestAddMemberTemplate_RendersWithoutError
    file: server/handler_members_test.go
    state: pending
  - id: add-form-has-fields
    name: TestAddMemberTemplate_HasAllFormFields
    file: server/handler_members_test.go
    state: pending
  - id: add-form-shows-error
    name: TestAddMemberTemplate_DisplaysErrorMessageWhenProvided
    file: server/handler_members_test.go
    state: pending
active_test: add-form-renders

### 7. Add HTTP handler for GET /members/{id}/edit
type: behavior
state: pending
files:
  - server/handler_members.go
  - server/handler_members_test.go
tests:
  - id: get-edit-form-success
    name: TestHandlerGetEditMemberPage_LoadsAndRendersForm
    file: server/handler_members_test.go
    state: pending
  - id: get-edit-form-not-found
    name: TestHandlerGetEditMemberPage_Returns404ForInvalidID
    file: server/handler_members_test.go
    state: pending
  - id: get-edit-form-deactivated
    name: TestHandlerGetEditMemberPage_Returns404ForDeactivatedMember
    file: server/handler_members_test.go
    state: pending
active_test: get-edit-form-success

### 8. Create edit member form template (edit.html)
type: behavior
state: pending
files:
  - server/templates/members/edit.html
tests:
  - id: edit-form-renders
    name: TestEditMemberTemplate_RendersWithoutError
    file: server/handler_members_test.go
    state: pending
  - id: edit-form-prefills
    name: TestEditMemberTemplate_PrefillsCurrentValues
    file: server/handler_members_test.go
    state: pending
  - id: edit-form-shows-error
    name: TestEditMemberTemplate_DisplaysErrorMessageWhenProvided
    file: server/handler_members_test.go
    state: pending
active_test: edit-form-renders

### 9. Update server.go to register new routes
type: tidy
state: pending
files:
  - server/server.go

### 10. Test and verify responsive behavior on all breakpoints
type: behavior
state: pending
files:
  - manual verification
tests:
  - id: responsive-desktop
    name: test Members screen on desktop (≥1024px)
    file: manual
    state: pending
  - id: responsive-tablet
    name: test Members screen on tablet (769–1023px)
    file: manual
    state: pending
  - id: responsive-mobile
    name: test Members screen on mobile (≤768px)
    file: manual
    state: pending
active_test: responsive-desktop

### 11. Test accessibility (WCAG 2.1 AA)
type: behavior
state: pending
files:
  - manual accessibility audit
tests:
  - id: a11y-semantic-html
    name: Verify semantic HTML (form, label, table, headings)
    file: manual
    state: pending
  - id: a11y-keyboard-nav
    name: Verify Tab, Enter, and Esc keys work as expected
    file: manual
    state: pending
  - id: a11y-focus-outline
    name: Verify 2px focus outline on all interactive elements
    file: manual
    state: pending
  - id: a11y-color-contrast
    name: Verify text contrast ≥4.5:1
    file: manual
    state: pending
active_test: a11y-semantic-html

### 12. Update roadmap and close increment
type: tidy
state: pending
files:
  - docs/roadmap.md

---

## Commits

(Populated as subtasks are completed)

---

## Notes

- Baseline tests to be run before subtask 1 begins
- All template changes will be embedded in the binary via embed.FS (already configured in server/server.go)
- Form submission flow: POST to `/members` (HTML form) redirects on success; PATCH/DELETE to `/api/members/{id}` returns JSON for client-side handling
- No new business logic; all validation and persistence delegated to existing use cases in core/members/
