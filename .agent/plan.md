# Plan: Members Screen (Vue 3)

## Goal

Build the Members screen in Vue 3 with full CRUD operations (list, add, edit, deactivate, reactivate) against the existing Members API, achieving parity with the prior Fyne desktop screen.

## Branch

`increment/members-screen-vue`

## Approach

The Members screen is composed of three Views connected by Vue Router (`/members`, `/members/add`, `/members/:id/edit`) and one Pinia store managing member list state, filters, and error handling. All CRUD operations delegate to the existing JSON API (`/api/members`). Components are tested in isolation with Vitest + Vue Test Utils (no browser); one Playwright acceptance test verifies the happy-path end-to-end flow (add a member, observe it in the list). The shell (sidebar, topbar) remains visible and functional during all navigation — this is guaranteed by Vue Router's persistent `AppShell` layout component wrapping the routed screens.

## Design

### Data Models

**Vue State (Pinia Store: `useMembersStore`)**
```
MembersState {
  members: MemberDTO[]           // List of members returned from API
  loading: boolean               // True while any API call is in-flight
  error: string | null          // Last error message; cleared on successful operation
  currentTab: 'active' | 'deactivated' | 'all'  // Selected filter tab
  
  // Members in current tab (computed from members + currentTab)
  filteredMembers: MemberDTO[]
}

MemberDTO (from API response)
{
  id: string                    // UUID (opaque internal ID; used only for routing/API calls)
  firstName: string             // Non-empty; displayed in UI
  lastName: string              // Non-empty; displayed in UI
  fullName: string              // Computed: `${firstName} ${lastName}`; displayed in table/forms
  seniority: string             // 'Junior' | 'Mid' | 'Senior' | 'Lead'
  status: string                // 'active' | 'deactivated'
  createdAt: string             // RFC3339 timestamp
}
```

**Form State (Component-local)**
```
AddEditFormState {
  form: {
    firstName: string           // Input field value
    lastName: string            // Input field value
    seniority: string           // Dropdown selection
  }
  isSubmitting: boolean         // True while POST/PATCH is in-flight
  fieldErrors: {                // Per-field validation errors
    firstName?: string
    lastName?: string
    seniority?: string
  }
}
```

### Call / Data Flow

**List Members (GET /api/members)**
1. User navigates to `/members` or clicks Members tab
2. `MemberListView.vue` mounts → `useMembersStore().loadMembers(tab)`
3. Store dispatches `setLoading(true)`, `setError(null)`
4. Fetch `GET http://backend:8080/api/members?status=<active|deactivated|all>` via `api/client.ts`
5. Parse JSON response → `{ members: MemberDTO[] }`
6. Store updates: `members = response.members`, `loading = false`
7. Component re-renders list with current tab filter
8. On error: store updates `error = 'Failed to load members'`, `loading = false`; component displays error toast

**Add Member (POST /api/members)**
1. User clicks "Add Member" button → navigate to `/members/add`
2. `AddMemberForm.vue` mounts with empty form
3. User fills firstName, lastName, seniority dropdown and clicks Save
4. Form validation runs (all fields required, seniority in allowed values)
5. If valid: dispatch `store.addMember(formData)`
6. Store: `setLoading(true)`, `setError(null)`, fetch `POST /api/members` with `{ firstName, lastName, seniority }`
7. Server returns `201 Created` with `{ id, firstName, lastName, seniority, status: "active", createdAt }`
8. Store: push new member to `members` array, `loading = false`
9. Component: programmatic navigate to `/members` (router.push)
10. `MemberListView.vue` rerenders with new member in Active tab
11. On error: store updates `error`, form remains visible, user can retry

**Edit Member (PATCH /api/members/{id})**
1. User clicks Edit on a member row → navigate to `/members/{id}/edit`
2. `EditMemberForm.vue` mounts → fetches member by ID from store (pre-loaded in list)
3. Form pre-fills with existing firstName, lastName, seniority
4. User modifies one or more fields and clicks Save
5. Form validation runs
6. If valid: dispatch `store.editMember(id, formData)`
7. Store: `setLoading(true)`, fetch `PATCH /api/members/{id}` with `{ firstName, lastName, seniority }`
8. Server returns `200 OK` with updated member
9. Store: find and update member in `members` array
10. Component: navigate to `/members`
11. List re-renders with updated data

**Deactivate Member (DELETE /api/members/{id})**
1. User clicks Deactivate on an Active member row
2. Component shows confirm dialog: "Are you sure?"
3. User confirms → dispatch `store.deactivateMember(id)`
4. Store: `setLoading(true)`, fetch `DELETE /api/members/{id}` (no body)
5. Server returns `204 No Content`
6. Store: update member status to 'deactivated' in `members` array
7. Component: re-render without full page reload; member moves from Active to Deactivated tab
8. On error: display error toast, member remains in Active tab

**Reactivate Member (PATCH /api/members/{id}/reactivate)**
1. User clicks Reactivate on a Deactivated member row
2. Component shows confirm dialog (optional; could be immediate)
3. Dispatch `store.reactivateMember(id)`
4. Store: `setLoading(true)`, fetch `PATCH /api/members/{id}/reactivate` (no body)
5. Server returns `200 OK` with updated member (status: 'active')
6. Store: update member status to 'active' in `members` array
7. Component: re-render; member moves from Deactivated to Active tab
8. On error (e.g., 409 Conflict if already active): display error toast

### Error / Edge-case Inventory

| Condition | Expected response | Covered by subtask |
|-----------|-------------------|--------------------|
| Network timeout (5s) on any API call | Display error toast, preserve form state if add/edit | Subtask 2, 3, 4 (client-level timeout; API has no control) |
| API returns 400 (validation error) | Extract error message from JSON, display in toast | Subtask 2, 3, 4 (error response parsing) |
| API returns 500 | Display generic "Server error" message, log to console | Subtask 2, 3, 4 (error handling) |
| User navigates away during in-flight request | Request completes in background; store state updates correctly; UI reflects latest state on return | Vitest (component cleanup) |
| Add/edit form with empty firstName or lastName | Form validation prevents submit; inline error displayed | Subtask 1 (form validation) |
| Add/edit form with invalid seniority | Form validation prevents submit; inline error displayed | Subtask 1 (form validation) |
| Attempt to reactivate an already-active member | API returns 409 Conflict; component displays error | Subtask 4 (conditional error handling) |
| Member deactivated while user is editing | Store state unchanged until next refresh; edit still succeeds (PATCH overwrites regardless of status) | N/A (concurrent edit is a future concern; not addressed in v1) |
| Member list is empty | "No members" message displayed for each tab | Subtask 1 (empty state) |

### Observability Intent

| Event | Level | Data |
|-------|-------|------|
| `members.api.list_requested` | debug | `status` (active/deactivated/all), timestamp |
| `members.api.list_success` | info | `count`, `status`, latency_ms |
| `members.api.list_failed` | warn | `status`, error_message, latency_ms |
| `members.api.add_requested` | debug | firstName, lastName, seniority (no PII concern in this tool) |
| `members.api.add_success` | info | new_member_id, latency_ms |
| `members.api.add_failed` | warn | error_message, latency_ms |
| `members.api.edit_requested` | debug | member_id, changed_fields |
| `members.api.edit_success` | info | member_id, latency_ms |
| `members.api.deactivate_requested` | debug | member_id |
| `members.api.deactivate_success` | info | member_id, latency_ms |
| `members.api.reactivate_requested` | debug | member_id |
| `members.api.reactivate_success` | info | member_id, latency_ms |

**Implementation note:** Console logging in dev mode only (guarded by `!import.meta.env.PROD`). No external analytics or telemetry (per CONSTITUTION.md).

### Architecture Delta

**No new containers.** `useMembersStore` is a new Pinia store inside the existing Frontend container. API routes are wired in the existing `server/server.go` (already present); no changes to backend architecture. Vue Router gains three new routes (`/members`, `/members/add`, `/members/:id/edit`); existing routes (`/`, `/alerts`, `/reports`) unchanged. AppShell remains the persistent layout wrapper; Members screen components render within the main content area via `<RouterView />`.

**Container diagram update not needed** — this is a feature within an existing frontend container, not a new service.

## Files

| File | Role | Notes |
|------|------|-------|
| `services/frontend/src/stores/members.ts` | new | Pinia store: list, add, edit, deactivate, reactivate; API client methods |
| `services/frontend/src/api/members.ts` | new | Members API client: CRUD methods (fetchMembers, addMember, editMember, deactivateMember, reactivateMember) |
| `services/frontend/src/views/Members.vue` | modify | Replace placeholder with functional Members list view (tabs, table, actions) |
| `services/frontend/src/views/AddMemberView.vue` | new | Form for adding a new member |
| `services/frontend/src/views/EditMemberView.vue` | new | Form for editing an existing member |
| `services/frontend/src/components/MemberList.vue` | new | Table component: render member rows with Edit/Deactivate/Reactivate buttons |
| `services/frontend/src/components/MemberForm.vue` | new | Shared form component (add/edit) with validation |
| `services/frontend/src/components/MemberTabs.vue` | new | Tab selector (Active/Deactivated/All) |
| `services/frontend/src/components/ErrorToast.vue` | new | Toast notification for API errors |
| `services/frontend/src/components/ConfirmDialog.vue` | new | Reusable confirm dialog for deactivate/reactivate actions |
| `services/frontend/src/router.ts` | modify | Add routes for `/members`, `/members/add`, `/members/:id/edit` |
| `services/frontend/src/stores/index.ts` | modify | Export `useMembersStore` |
| `services/frontend/src/stores/__tests__/members.spec.ts` | new | Vitest tests for `useMembersStore` (list, add, edit, deactivate, reactivate) |
| `services/frontend/src/components/__tests__/MemberList.spec.ts` | new | Vitest tests for MemberList component rendering and interactions |
| `services/frontend/src/components/__tests__/MemberForm.spec.ts` | new | Vitest tests for MemberForm validation and submission |
| `services/frontend/src/components/__tests__/MemberTabs.spec.ts` | new | Vitest tests for tab switching |
| `services/frontend/src/components/__tests__/ErrorToast.spec.ts` | new | Vitest tests for error display |
| `acceptance-tests/tests/members-crud.spec.ts` | new | Playwright test: add a member, verify in list, edit, verify, deactivate, verify |
| `docs/ui.md` | touch | Reference — member management UI patterns already documented in prior shell work |

## Subtasks

### 1. [tidy] Create Members Pinia store with state shape and computed filter

**Description:** Create the `useMembersStore()` store with reactive state for the member list, current tab filter, loading state, and error message. Add a computed property that filters members by the current tab. No API calls yet; this is structure only.

**Files:**
- `services/frontend/src/stores/members.ts` (new)
- `services/frontend/src/stores/index.ts` (modify — export store)

**References:**
- `services/frontend/src/stores/health.ts` (pattern: Pinia store composition API, ref/computed)
- `docs/architecture.md#containers` (frontend state management strategy)

**Verification:** `npm test -- stores/members` → store loads without errors; initial state is correct; `filteredMembers` computed property returns empty array initially

**Depends on:** —

**Acceptance criteria covered:** AC-1 (foundation for filtering tabs)

---

### 2. [tidy] Create Members API client module with fetch methods

**Description:** Extract CRUD methods from implicit calls to a dedicated `api/members.ts` module. Implement `fetchMembers(status)`, `addMember(data)`, `editMember(uuid, data)`, `deactivateMember(uuid)`, `reactivateMember(uuid)` functions. Each method calls the appropriate backend endpoint with proper error handling (HTTP status checks, JSON parse, timeout). Does not call store; pure functions.

**Files:**
- `services/frontend/src/api/members.ts` (new)

**References:**
- `services/frontend/src/api/client.ts` (pattern: fetch helpers, error handling, timeout)
- `services/backend/server/handler_members.go:44-339` (API contract: request/response shapes)

**Verification:** `npm test -- api/members` → all functions callable; fetch operations parse JSON correctly; errors throw with descriptive messages

**Depends on:** —

**Acceptance criteria covered:** AC-1, AC-2, AC-3, AC-4, AC-5 (foundation for all CRUD)

---

### 3. [tidy] Create MemberForm component with validation logic

**Description:** Build a shared `MemberForm.vue` component that accepts optional `initialData` prop (for edit mode) and emits a `submit` event with form data. Implement validation: firstName and lastName must be non-empty; seniority must be one of the allowed values ('Junior', 'Mid', 'Senior', 'Lead'). Display inline error messages for each field. Submit button is disabled until form is valid.

**Files:**
- `services/frontend/src/components/MemberForm.vue` (new)
- `services/frontend/src/components/__tests__/MemberForm.spec.ts` (new)

**References:**
- `services/frontend/src/views/Home.vue` (pattern: component structure, styling)
- `docs/ui.md` (form styling, button patterns, accessibility)

**Tests:**
- id: form-render
  file: `services/frontend/src/components/__tests__/MemberForm.spec.ts`
  name: `renders form fields with labels`
  state: pending
- id: form-empty-validation
  file: `services/frontend/src/components/__tests__/MemberForm.spec.ts`
  name: `disables submit when firstName or lastName is empty`
  state: pending
- id: form-invalid-seniority
  file: `services/frontend/src/components/__tests__/MemberForm.spec.ts`
  name: `shows error when seniority is not in allowed values`
  state: pending
- id: form-valid-submit
  file: `services/frontend/src/components/__tests__/MemberForm.spec.ts`
  name: `emits submit event with form data when valid and user clicks Save`
  state: pending
- id: form-pre-fill
  file: `services/frontend/src/components/__tests__/MemberForm.spec.ts`
  name: `pre-fills form fields from initialData prop`
  state: pending

**active_test:** form-render

**Verification:** `npm test -- MemberForm` → all 5 tests pass; form renders, validation blocks submit, submit event fires with correct payload

**Depends on:** —

**Acceptance criteria covered:** AC-2 (add form), AC-3 (edit form), AC-6 (validation)

---

### 4. [behavior] Store actions: loadMembers, addMember, editMember, deactivateMember, reactivateMember

**Description:** Implement async actions in `useMembersStore()` that call the Members API client methods and update store state. Each action sets `loading = true`, clears `error`, calls the API, and on success updates the members list or individual member. On error, sets `error` and leaves `loading = false`. Use optional console logging (dev mode only) for debugging.

**Files:**
- `services/frontend/src/stores/members.ts` (modify)
- `services/frontend/src/stores/__tests__/members.spec.ts` (new)

**References:**
- `services/frontend/src/api/members.ts` (API client methods)
- `services/frontend/src/stores/health.ts:16-43` (pattern: async action with error handling)

**Tests:**
- id: load-members-success
  file: `services/frontend/src/stores/__tests__/members.spec.ts`
  name: `loadMembers updates state with members and sets status to active tab`
  state: pending
- id: load-members-error
  file: `services/frontend/src/stores/__tests__/members.spec.ts`
  name: `loadMembers sets error message on API failure`
  state: pending
- id: add-member-success
  file: `services/frontend/src/stores/__tests__/members.spec.ts`
  name: `addMember appends new member to list and returns ID`
  state: pending
- id: add-member-error
  file: `services/frontend/src/stores/__tests__/members.spec.ts`
  name: `addMember sets error message without modifying list`
  state: pending
- id: edit-member-success
  file: `services/frontend/src/stores/__tests__/members.spec.ts`
  name: `editMember updates existing member in list`
  state: pending
- id: deactivate-member-success
  file: `services/frontend/src/stores/__tests__/members.spec.ts`
  name: `deactivateMember updates member status to deactivated`
  state: pending
- id: reactivate-member-success
  file: `services/frontend/src/stores/__tests__/members.spec.ts`
  name: `reactivateMember updates member status to active`
  state: pending

**active_test:** load-members-success

**Verification:** `npm test -- stores/members` → all 7 tests pass; store state updates correctly after each action; errors are captured and logged

**Depends on:** Subtask 1 (state shape), Subtask 2 (API client)

**Acceptance criteria covered:** AC-1, AC-2, AC-3, AC-4, AC-5, AC-7

---

### 5. [tidy] Create MemberList, MemberTabs, ErrorToast, ConfirmDialog components

**Description:** Build four UI components:
- `MemberList.vue`: Renders a table with member rows (firstName, lastName, seniority, status); each row has Edit, Deactivate, or Reactivate buttons (button choice based on status).
- `MemberTabs.vue`: Three buttons (Active, Deactivated, All); emits `tab-changed` event when clicked; highlights the current tab.
- `ErrorToast.vue`: Displays error message in a dismissable toast notification (Bulma alert style).
- `ConfirmDialog.vue`: Modal confirmation dialog with title, message, Cancel/Confirm buttons; emits `confirmed` event.

**Files:**
- `services/frontend/src/components/MemberList.vue` (new)
- `services/frontend/src/components/MemberTabs.vue` (new)
- `services/frontend/src/components/ErrorToast.vue` (new)
- `services/frontend/src/components/ConfirmDialog.vue` (new)

**References:**
- `services/frontend/src/components/AppShell.vue` (pattern: responsive layout, Bulma CSS)
- `docs/ui.md` (styling, button styles, accessibility)

**Verification:** `npm test` → all components render without errors; props and emitted events work correctly

**Depends on:** —

**Acceptance criteria covered:** AC-1 (tabs), AC-4 (deactivate), AC-5 (reactivate), AC-6 (validation errors), AC-7 (error handling)

---

### 6. [tidy] Add component unit tests for MemberList, MemberTabs, ErrorToast, ConfirmDialog

**Description:** Write Vitest + Vue Test Utils tests covering:
- MemberList: renders rows for each member; Edit/Deactivate/Reactivate buttons emit correct events; empty state message
- MemberTabs: all three tabs render and emit event when clicked; active tab is highlighted
- ErrorToast: displays error message; dismiss button clears it
- ConfirmDialog: renders title and message; Cancel and Confirm buttons emit correct events

**Files:**
- `services/frontend/src/components/__tests__/MemberList.spec.ts` (new)
- `services/frontend/src/components/__tests__/MemberTabs.spec.ts` (new)
- `services/frontend/src/components/__tests__/ErrorToast.spec.ts` (new)
- `services/frontend/src/components/__tests__/ConfirmDialog.spec.ts` (new)

**Verification:** `npm test -- components` → 4 files, ~15 tests, all pass

**Depends on:** Subtask 5 (component creation)

**Acceptance criteria covered:** AC-8 (test coverage)

---

### 7. [behavior] Implement Members list view with tab filtering and loading state

**Description:** Replace the placeholder `Members.vue` view with a functional list screen. Wire `MemberTabs`, `MemberList`, and `ErrorToast` components. On mount, call `store.loadMembers('active')`. When user clicks a tab, call `store.setCurrentTab(tab)` and reload. Display loading spinner while `store.loading` is true. Render error toast when `store.error` is set.

**Files:**
- `services/frontend/src/views/Members.vue` (modify)

**References:**
- `services/frontend/src/stores/members.ts` (store actions)
- `services/frontend/src/components/MemberList.vue` (component prop/event contract)
- `services/frontend/src/components/MemberTabs.vue` (component prop/event contract)

**Tests:**
- id: list-load-active
  file: `services/frontend/src/views/__tests__/Members.spec.ts`
  name: `loads and displays active members on mount`
  state: pending
- id: list-tab-switch
  file: `services/frontend/src/views/__tests__/Members.spec.ts`
  name: `switches to deactivated tab and reloads members`
  state: pending
- id: list-error
  file: `services/frontend/src/views/__tests__/Members.spec.ts`
  name: `displays error toast when store.error is set`
  state: pending

**active_test:** list-load-active

**Verification:** `npm test -- views/Members` → 3 tests pass; component loads, tabs switch, errors display

**Depends on:** Subtask 1, 4, 5

**Acceptance criteria covered:** AC-1, AC-8

---

### 8. [behavior] Implement AddMemberView and wire to router

**Description:** Create `AddMemberView.vue` that renders the `MemberForm` component in add mode (no initialData). On form submit, call `store.addMember(formData)`. Show loading spinner and disable form while submitting. On success, navigate to `/members`. On error, display error message and allow retry.

**Files:**
- `services/frontend/src/views/AddMemberView.vue` (new)
- `services/frontend/src/views/__tests__/AddMemberView.spec.ts` (new)
- `services/frontend/src/router.ts` (modify — add route)

**References:**
- `services/frontend/src/router.ts` (route structure)
- `services/frontend/src/components/MemberForm.vue` (submit event shape)

**Tests:**
- id: add-form-submit
  file: `services/frontend/src/views/__tests__/AddMemberView.spec.ts`
  name: `submits form data to store and navigates to members list on success`
  state: pending
- id: add-form-error
  file: `services/frontend/src/views/__tests__/AddMemberView.spec.ts`
  name: `displays error and allows retry on API failure`
  state: pending

**active_test:** add-form-submit

**Verification:** `npm test -- AddMemberView` → 2 tests pass; form submits, navigation works, error handling works

**Depends on:** Subtask 3, 4

**Acceptance criteria covered:** AC-2, AC-6, AC-7, AC-10

---

### 9. [behavior] Implement EditMemberView and wire to router

**Description:** Create `EditMemberView.vue` that routes on `/members/:id/edit`, loads the member data from the store (by matching ID to member in `members` array), pre-fills the form, and on submit calls `store.editMember(id, formData)`. Handle case where member is not found (display error, redirect). Show loading spinner during submission.

**Files:**
- `services/frontend/src/views/EditMemberView.vue` (new)
- `services/frontend/src/views/__tests__/EditMemberView.spec.ts` (new)
- `services/frontend/src/router.ts` (modify — add route with param)

**References:**
- `services/frontend/src/router.ts:4-25` (route param syntax)
- `services/frontend/src/views/AddMemberView.vue` (similar form wiring pattern)

**Tests:**
- id: edit-form-prefill
  file: `services/frontend/src/views/__tests__/EditMemberView.spec.ts`
  name: `loads member data and pre-fills form fields`
  state: pending
- id: edit-form-submit
  file: `services/frontend/src/views/__tests__/EditMemberView.spec.ts`
  name: `submits edited data to store and navigates to members list on success`
  state: pending
- id: edit-member-not-found
  file: `services/frontend/src/views/__tests__/EditMemberView.spec.ts`
  name: `displays error when member ID is not found in store`
  state: pending

**active_test:** edit-form-prefill

**Verification:** `npm test -- EditMemberView` → 3 tests pass; form pre-fills, submits, navigation works

**Depends on:** Subtask 3, 4

**Acceptance criteria covered:** AC-3, AC-6, AC-7, AC-10

---

### 10. [behavior] Wire deactivate and reactivate actions in MemberList component

**Description:** Update `MemberList.vue` to emit `deactivate` and `reactivate` events (with member ID) when buttons are clicked. Update the parent `Members.vue` view to listen to these events, show a confirm dialog, and on user confirmation call `store.deactivateMember()` or `store.reactivateMember()`. Show loading spinner during action. Update member row immediately after success (status changes, button options change) without full page reload.

**Files:**
- `services/frontend/src/components/MemberList.vue` (modify)
- `services/frontend/src/views/Members.vue` (modify — handle deactivate/reactivate events)
- `services/frontend/src/components/__tests__/MemberList.spec.ts` (modify — add event tests)

**References:**
- `services/frontend/src/components/ConfirmDialog.vue` (confirm flow)
- `services/frontend/src/stores/members.ts` (deactivate/reactivate actions)

**Tests:**
- id: deactivate-button-emit
  file: `services/frontend/src/components/__tests__/MemberList.spec.ts`
  name: `emits deactivate event when Deactivate button clicked`
  state: pending
- id: reactivate-button-emit
  file: `services/frontend/src/components/__tests__/MemberList.spec.ts`
  name: `emits reactivate event when Reactivate button clicked`
  state: pending
- id: deactivate-confirm-action
  file: `services/frontend/src/views/__tests__/Members.spec.ts`
  name: `calls store.deactivateMember after user confirms`
  state: pending
- id: reactivate-confirm-action
  file: `services/frontend/src/views/__tests__/Members.spec.ts`
  name: `calls store.reactivateMember after user confirms`
  state: pending

**active_test:** deactivate-button-emit

**Verification:** `npm test` → all tests pass; deactivate/reactivate buttons work, confirm dialog flows, members update without reload

**Depends on:** Subtask 5, 7

**Acceptance criteria covered:** AC-4, AC-5, AC-8

---

### 11. [tidy] Add "Add Member" button to AppShell sidebar footer

**Description:** Update the placeholder "Add Item" button in `AppShell.vue` footer to navigate to `/members/add` (RouterLink or programmatic navigation). Button text: "Add Member". Icon: fa-plus (already present). Only show button when user is on `/members` path (optional; or always show and let it navigate from anywhere).

**Files:**
- `services/frontend/src/components/AppShell.vue` (modify)

**References:**
- `services/frontend/src/components/AppShell.vue:114-120` (current button)

**Verification:** `npm test` → AppShell still renders; "Add Member" button navigates to `/members/add`

**Depends on:** Subtask 8 (AddMemberView exists)

**Acceptance criteria covered:** AC-2 (accessibility to add form)

---

### 12. [behavior] Write Playwright acceptance test: add member happy path

**Description:** Create `acceptance-tests/tests/members-crud.spec.ts` with one test covering the main user flow:
1. Navigate to app, wait for shell to load
2. Click Members in sidebar
3. Verify Active tab is selected and members list is visible
4. Click "Add Member" button
5. Fill form (firstName: "Alice", lastName: "Smith", seniority: "Senior")
6. Click Save
7. Verify redirected to /members list
8. Verify "Alice Smith" appears in Active tab with correct seniority
9. Click Edit on Alice's row
10. Verify form pre-filled with correct values
11. Change seniority to "Lead"
12. Click Save
13. Verify "Alice Smith" shows "Lead" in the list
14. Click Deactivate on Alice's row
15. Confirm dialog
16. Verify Alice moves to Deactivated tab
17. Click Reactivate on Alice
18. Verify Alice moves back to Active tab

**Files:**
- `acceptance-tests/tests/members-crud.spec.ts` (new)

**References:**
- `acceptance-tests/` (Playwright setup, existing test structure if any)
- `docs/testing.md` (acceptance test scope and strategy)

**Verification:** `npm run test:acceptance -- members-crud` → test passes against running docker-compose environment (backend healthy, frontend loaded)

**Depends on:** Subtask 7, 8, 9, 10 (all UI features complete)

**Acceptance criteria covered:** AC-1, AC-2, AC-3, AC-4, AC-5, AC-9

---

### 13. [tidy] Update docs/ui.md: add Members screen section

**Description:** Document the Members screen UI pattern in `docs/ui.md`:
- Three-tab navigation (Active/Deactivated/All)
- Table layout with firstName, lastName, seniority, actions columns
- Form layout for add/edit (fields, validation, error display)
- Confirm dialogs for deactivate/reactivate
- Empty state ("No members in this view")
- Loading spinner state
- Error toast display
- Responsive behavior on mobile

**Files:**
- `docs/ui.md` (modify)

**Verification:** `docs/ui.md` updated; no build or test required; manual review by team

**Depends on:** Subtask 7, 8, 9, 10 (UI complete)

**Acceptance criteria covered:** AC-10 (responsive), AC-6 (validation)

---

## Context Map

- `CONSTITUTION.md#testing-strategy` — component tests (Vitest), handler tests skip HTML rendering, acceptance tests narrow scope to one main flow per feature
- `CONSTITUTION.md#performance-envelope` — screen render < 200ms (verify with dev tools; not an automated gate in v1)
- `CONSTITUTION.md#delivery-and-documentation` — feature done when acceptance criteria pass and roadmap evidence is linked
- `docs/architecture.md#containers` — frontend container: Vue 3 SPA, Pinia state, Vue Router, calls `/api/*` endpoints
- `docs/architecture.md#communication-paths` — Vue → HTTP → backend; JSON contract only
- `docs/testing.md` — Vitest component tests at unit layer; Playwright acceptance tests at feature layer
- `docs/ui.md` — existing Bulma framework decisions, button/form/modal patterns (reference, no changes needed unless documenting new patterns)
- `docs/adr/ADR-20260917-vue-spa-frontend.md` — Vue 3 SPA stack decision and rationale

## Acceptance Scenarios

These scenarios describe user actions that span multiple components and would be difficult to test at the unit level. They are advisory unless marked `required`.

### AS-1: Add and view new member

- criterion: AC-2
- user action: Click "Add Member", fill form, save
- precondition: User is on `/members` view
- expected outcome: New member appears in Active tab immediately after redirect
- evidence: `acceptance-tests/tests/members-crud.spec.ts`
- gate: advisory
- state: planned

### AS-2: Edit and verify persistence

- criterion: AC-3
- user action: Click Edit, change seniority, save
- precondition: User is on `/members` view with at least one member
- expected outcome: Changes persist in the list; refresh page shows updated values
- evidence: `acceptance-tests/tests/members-crud.spec.ts`
- gate: advisory
- state: planned

### AS-3: Deactivate and reactivate member

- criterion: AC-4, AC-5
- user action: Click Deactivate, confirm; later click Reactivate, confirm
- precondition: User is on `/members` view
- expected outcome: Member moves between tabs without page reload; status reflected in sidebar/overview (future)
- evidence: `acceptance-tests/tests/members-crud.spec.ts`
- gate: advisory
- state: planned

## Risks

1. **Concurrent edit:** If member is edited by another browser tab while user is editing, the update will overwrite without warning. This is out of scope for v1.
   - Mitigation: Document as known limitation; future increment will add optimistic locking or conflict detection.

2. **API timeout on slow networks:** 5-second fetch timeout may be too aggressive. If users are on slow connections, the app may show errors prematurely.
   - Mitigation: Console logs will surface timeouts in dev mode; manual testing in slow-network conditions recommended before release.

3. **Form state reset on navigation:** If user fills form, navigates away, then returns to the form, the form will be empty (form state is component-local, not persistent).
   - Mitigation: Expected behavior (forms don't auto-save); document in future UX review if needed.

4. **Empty member list handling:** Three tabs with potentially empty arrays. Component must render empty state for each tab.
   - Mitigation: Subtask 5 and 7 explicitly cover empty state rendering.

## Planning Decisions

| Decision | Chosen | Rejected | Reason |
|----------|--------|----------|--------|
| Store library | Pinia | Vuex, Redux | Pinia is lightweight, composition-API-first (modern Vue 3 pattern), and already a dependency (from health store setup) |
| Form library | Native Vue + validation logic | VeeValidate, Formik | Keep dependencies minimal; form is simple (3 fields, 4 validation rules); native Vue is maintainable |
| API client location | Dedicated `api/members.ts` module | Inline in store | Separation of concerns; client is testable independently; reusable if CLI client added later |
| ID display vs. storage | Use UUID as opaque internal ID; display full name in UI | UUID-heavy UI | UUIDs are the contract with backend; no conversion logic; UI friendly to users (shows names, not IDs); routes use UUID (proper REST semantics) |
| Confirm dialogs | Reusable `ConfirmDialog.vue` component | Browser `confirm()` | Reusable, styled with Bulma, testable; browser confirm is not accessible |
| Tab navigation | Re-fetch on tab click | Filter client-side | Keeps store single source of truth; ensures server and client agree on what "deactivated" means; no sorting/filtering ambiguity |
| Error messages | Toast notifications + store error field | Inline validation | Inline is for form validation; toast for API errors (network, 500, etc.) is user-friendly and dismissable |
| Loading state | Global `store.loading` | Per-action spinners | Simpler; single loading spinner during any CRUD operation; acceptable UX for this tool (not a high-frequency CRUD app) |
| Router lazy loading | Yes (`() => import(...)`) | Eager imports | Performance: members screen is not needed until user navigates; lazy loading reduces initial bundle size |
| Component co-location | `views/` for screens, `components/` for reusables | Flat structure | Standard Vue project organization; clear intent; screens render in RouterView, reusables used across multiple screens |

---

## Next Action

On user approval, this plan will be handed off to the implement skill to execute the subtasks in order. The orchestrator will detect the first subtask's type (tidy) and load the appropriate implement skill.
