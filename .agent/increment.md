# Increment: Members Screen (Vue 3)

## Use Case

When I open the Members screen in the SPA, I want to view, add, edit, and manage the status of team members with a responsive form interface, so that I can keep the team roster current and control which members appear in scoring and alerts.

## Goal

Build the Members screen in Vue 3 with full CRUD operations (list, add, edit, deactivate, reactivate) against the existing Members API (`/api/members`), achieving parity with the prior Fyne desktop screen.

## Branch

`increment/members-screen-vue`

## Acceptance Criteria

1. **Members list displays all three filter tabs** (Active, Deactivated, All) backed by the existing status-filtered API; the shell (sidebar, top bar) remains visible and functional.

2. **Add member flow** — Navigate to `/members/add`, fill form (first name, last name, seniority dropdown), submit; redirects to `/members` list and new member appears immediately in Active tab.

3. **Edit member flow** — Click Edit on an Active member row, form pre-fills with existing data, change one or more fields, save, redirects to list and updates persist.

4. **Deactivate action** — Click Deactivate on an Active member, confirm prompt, member moves to Deactivated tab without a full page reload.

5. **Reactivate action** — Click Reactivate on a Deactivated member, member moves back to Active tab without a full page reload.

6. **Form validation** — Empty or invalid fields show inline error messages; submit is disabled until form is valid.

7. **Error handling** — API errors (400, 500) are caught and displayed as toast notifications or inline messages; user can retry.

8. **Vitest component tests pass** — MemberList, MemberForm, member store cover rendering, props, form submission, validation, and error states without a browser.

9. **One Playwright test passes** — End-to-end flow: add a member with valid data, observe member appear in Active tab, edit seniority, verify change persists on refresh.

10. **Vue Router integration** — Routes `/members` (list) and `/members/add` and `/members/:id/edit` work; redirects on save complete successfully.

## Acceptance-Test Intent

**Main user journey:**
1. Load app → navigate to Members via sidebar link
2. See Active, Deactivated, All tabs; Active tab displays team roster
3. Click "Add Member" → navigate to `/members/add` with form
4. Fill form (first name: "Alice", last name: "Smith", seniority: "Senior")
5. Click Save → redirected to `/members` list; "Alice Smith" appears in Active tab
6. Click Edit on Alice → form prefilled; change seniority to "Lead"; save
7. Verify change persists (edit again to confirm, or refresh page)
8. Click Deactivate → confirm dialog → Alice moves to Deactivated tab
9. Switch to Deactivated tab → see Alice; click Reactivate → Alice returns to Active tab
10. No full-page reloads; shell remains constant throughout

## Out Of Scope

- Member search or filtering (beyond tab status) — filter by name/ID is future work
- Bulk operations (add/remove multiple members at once)
- Import/export member CSV
- Member profile pages or detailed history view (that is Screen C: Member Detail)
- Changing member ID or history reassignment logic (see roadmap open questions on duplicate names)
- Accessibility audit beyond WCAG 2.1 AA semantic HTML (keyboard nav, ARIA labels inline with the shell; detailed audit is a separate increment)

## Constitution Constraints

- **Behavior first:** Form validation and submission logic tested in Vitest before form components render
- **Service is the API contract:** Vue components call only `GET /api/members`, `POST /api/members`, `PATCH /api/members/{id}`, `DELETE /api/members/{id}` — zero direct store queries or coordinator imports
- **Presentation and API separate:** All business logic (member CRUD) lives in backend; Vue components are thin renderers of the API response
- **Testing strategy:** Vitest + Vue Test Utils for component logic (no browser); one Playwright test for the main user flow (with browser, against running docker-compose)
- **Small, focused changes:** This increment is the Members screen only; Monthly Input, Alerts, Overview, etc. are separate future increments

## Roadmap Entry

**Feature:** Frontend Increment 2: Members Screen (Vue)

**Status:** Partial (entering In Progress)

**Job Story:** When I open the Members screen in the SPA, I want to view, add, edit, and manage the status of team members with a responsive form interface, so that I can keep the team roster current and control which members appear in scoring and alerts.

**Acceptance criteria:** All 10 criteria listed above.

**Evidence:** (to be added upon completion)
- Vitest test count and file locations
- Playwright test location and results
- Commit hashes
- `npm test` and `go test -race ./...` passing

---

## Next Action

User approval needed on:
1. Branch name: `increment/members-screen-vue`
2. Acceptance criteria (10 binary, verifiable conditions)
3. Scope boundaries (out-of-scope list)

Once approved, the next step is `4dc-plan` to convert this intent into a detailed technical execution plan with file-level detail.
