# Increment: Members Screen SPA

## Use Case

When I navigate to the Members screen in the web app, I want to view all active team members, add new members via a dedicated page, edit or deactivate existing members, so that my team roster is always accurate and I can manage members efficiently in the browser.

## Goal

Build the Members screen as a browser-rendered SPA page with human-friendly URLs, backed by the existing Members API, replacing the Fyne Settings screen for member management.

## Branch

`increment/members-screen-spa`

## Acceptance Criteria

1. **Members list page** (`/members`) displays all active members in a responsive table with columns: Name, Seniority, Actions (Edit, Deactivate)
2. **Add member page** (`/members/add`) provides a form to create a new member; on success, redirects to `/members` list
3. **Edit member page** (`/members/{id}/edit`) loads existing member data and allows updates; on success, redirects to `/members` list
4. **Deactivate action** on list removes member from list immediately (soft delete via API); list updates without full page reload
5. **Form validation** blocks submission with descriptive errors for missing/invalid fields (name required, seniority required, name length ≤ 100 chars)
6. **All changes persist** across browser reloads and app restarts (database-backed via SQLite)
7. **Responsive layout** works on desktop (≥1024px), tablet (769–1023px), and mobile (≤768px) per `docs/ui.md` breakpoints
8. **Accessibility** meets WCAG 2.1 AA: semantic HTML, ARIA labels, keyboard navigation, focus states per shell design system

## Acceptance-Test Intent

**User journey:**
1. User opens `/members`, sees list of 3 existing members
2. Clicks "Add Member" button, navigates to `/members/add`
3. Fills form: first name, last name, seniority; submits
4. Redirected to `/members`, sees new member in list
5. Clicks "Edit" on a member, navigates to `/members/{id}/edit`
6. Updates seniority, saves; redirected to list
7. Clicks "Deactivate", member removed from list without page reload
8. Reloads page; member still deactivated (persistence verified)

## Out Of Scope

- Import/export members (CSV)
- Member search or filtering (simple list only)
- Reactivating deactivated members (deactivation is permanent in this increment)
- Member notes, email, or custom fields (name + seniority only)
- Bulk operations (add/edit/deactivate one at a time)
- Role-based access control (all users have full member management access)

## Constitution Constraints

- **Behavior first:** Form validation and persistence logic has unit tests before UI renders
- **Dependency direction:** Routes → HTTP handlers → coordinator → service → store → engine (no reverse imports)
- **Small, focused changes:** Screen is a single increment; logic stays in coordinator reusable by CLI/API
- **Human review:** No auto-generated UX; all navigation human-readable URLs
- **Performance:** List load + add/edit/deactivate round-trip ≤ 200ms on local machine (per performance envelope in CONSTITUTION)
- **Service is API contract:** Handlers delegate to coordinator; SPA is replaceable presentation layer

## Roadmap Entry

**Members Screen (SPA):** When I manage my team roster through the web browser, I want surfable, bookmarkable URLs for adding and editing members, so that the interface feels like a professional web app, not a modal-heavy desktop tool.

---

**Next step:** Load `4dc-plan` skill to convert this increment into a technical execution plan.
