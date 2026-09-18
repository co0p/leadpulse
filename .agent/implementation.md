# Implementation: Members Screen (Vue 3)

**Branch:** `increment/members-screen-vue`

**Goal:** Build the Members screen in Vue 3 with full CRUD operations (list, add, edit, deactivate, reactivate) against the existing Members API.

---

## Subtasks

### Subtask 1: [tidy] Create Members Pinia store with state shape and computed filter

**Status:** complete

**Plan reference:** `.agent/plan.md` — Subtask 1

**Description:** Create the `useMembersStore()` store with reactive state for the member list, current tab filter, loading state, and error message. Add a computed property that filters members by the current tab. No API calls yet; this is structure only.

**Files to create/modify:**
- `services/frontend/src/stores/members.ts` (new) ✅
- `services/frontend/src/stores/index.ts` (modify — export store) ✅

**Verification:** `npm test -- stores/members` → store loads without errors; initial state is correct; `filteredMembers` computed property returns empty array initially ✅

**Evidence:**
- state: complete
- tests: 11 tests passing
- commits: b63a16a (tidy: create Members Pinia store with state shape and computed filter)

---

### Subtask 2: [tidy] Create Members API client module with fetch methods

**Status:** complete

**Plan reference:** `.agent/plan.md` — Subtask 2

**Description:** Extract CRUD methods from implicit calls to a dedicated `api/members.ts` module. Implement `fetchMembers(status)`, `addMember(data)`, `editMember(uuid, data)`, `deactivateMember(uuid)`, `reactivateMember(uuid)` functions. Each method calls the appropriate backend endpoint with proper error handling (HTTP status checks, JSON parse, timeout). Does not call store; pure functions.

**Files to create/modify:**
- `services/frontend/src/api/members.ts` (new) ✅

**Verification:** `npm test -- api/members` → all functions callable; fetch operations parse JSON correctly; errors throw with descriptive messages ✅

**Evidence:**
- state: complete
- tests: All existing tests still pass (26 tests)
- commits: 41ae115 (tidy: create Members API client module with fetch methods)

---

### Subtask 3: [tidy] Create MemberForm component with validation logic

**Status:** complete

**Plan reference:** `.agent/plan.md` — Subtask 3

**Description:** Build a reusable form component in `components/MemberForm.vue` that renders firstName (text), lastName (text), and seniority (dropdown with Junior/Mid/Senior/Lead). Implement field-level validation that shows error messages on blur. Provide props: `initialData` (optional MemberFormData), `isSubmitting` (bool), `submitButtonLabel` (string, default "Save"). Emit `@submit` with form data (trimmed) and `@cancel` events. Do not call API or store; pure UI component with local form state.

**Files to create/modify:**
- `services/frontend/src/components/MemberForm.vue` (new) ✅
- `services/frontend/src/components/__tests__/MemberForm.spec.ts` (new) ✅

**Verification:** `npm test -- MemberForm` → 14 tests passing; form validates all fields; emits correct events; pre-fills from initialData ✅

**Evidence:**
- state: complete
- tests: 14 tests passing (form field rendering, validation errors, submit/cancel events, pre-fill behavior, whitespace trimming)
- commits: aad1aa0 (tidy: create MemberForm component with validation logic)

---

### Subtask 4: [behavior] Store actions: loadMembers, addMember, editMember, deactivateMember, reactivateMember

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 4

**Description:** Implement async actions in `useMembersStore()` that call the Members API client methods and update store state. Each action sets `loading = true`, clears `error`, calls the API, and on success updates the members list or individual member. On error, sets `error` and leaves `loading = false`. Use optional console logging (dev mode only) for debugging.

**Files to create/modify:**
- `services/frontend/src/stores/members.ts` (modify)
- `services/frontend/src/stores/__tests__/members.spec.ts` (new)

**Verification:** `npm test -- stores/members` → all 7 tests pass; store state updates correctly after each action; errors are captured and logged

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 5: [tidy] Create MemberList, MemberTabs, ErrorToast, ConfirmDialog components

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 5

**Description:** Build four UI components:
- `MemberList.vue`: Renders a table with member rows (firstName, lastName, seniority, status); each row has Edit, Deactivate, or Reactivate buttons (button choice based on status).
- `MemberTabs.vue`: Three buttons (Active, Deactivated, All); emits `tab-changed` event when clicked; highlights the current tab.
- `ErrorToast.vue`: Displays error message in a dismissable toast notification (Bulma alert style).
- `ConfirmDialog.vue`: Modal confirmation dialog with title, message, Cancel/Confirm buttons; emits `confirmed` event.

**Files to create/modify:**
- `services/frontend/src/components/MemberList.vue` (new)
- `services/frontend/src/components/MemberTabs.vue` (new)
- `services/frontend/src/components/ErrorToast.vue` (new)
- `services/frontend/src/components/ConfirmDialog.vue` (new)

**Verification:** `npm test` → all components render without errors; props and emitted events work correctly

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 6: [tidy] Add component unit tests for MemberList, MemberTabs, ErrorToast, ConfirmDialog

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 6

**Description:** Write Vitest + Vue Test Utils tests covering:
- MemberList: renders rows for each member; Edit/Deactivate/Reactivate buttons emit correct events; empty state message
- MemberTabs: all three tabs render and emit event when clicked; active tab is highlighted
- ErrorToast: displays error message; dismiss button clears it
- ConfirmDialog: renders title and message; Cancel and Confirm buttons emit correct events

**Files to create/modify:**
- `services/frontend/src/components/__tests__/MemberList.spec.ts` (new)
- `services/frontend/src/components/__tests__/MemberTabs.spec.ts` (new)
- `services/frontend/src/components/__tests__/ErrorToast.spec.ts` (new)
- `services/frontend/src/components/__tests__/ConfirmDialog.spec.ts` (new)

**Verification:** `npm test -- components` → 4 files, ~15 tests, all pass

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 7: [behavior] Implement Members list view with tab filtering and loading state

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 7

**Description:** Replace the placeholder `Members.vue` view with a functional list screen. Wire `MemberTabs`, `MemberList`, and `ErrorToast` components. On mount, call `store.loadMembers('active')`. When user clicks a tab, call `store.setCurrentTab(tab)` and reload. Display loading spinner while `store.loading` is true. Render error toast when `store.error` is set.

**Files to create/modify:**
- `services/frontend/src/views/Members.vue` (modify)

**Verification:** `npm test -- views/Members` → 3 tests pass; component loads, tabs switch, errors display

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 8: [behavior] Implement AddMemberView and wire to router

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 8

**Description:** Create `AddMemberView.vue` that renders the `MemberForm` component in add mode (no initialData). On form submit, call `store.addMember(formData)`. Show loading spinner and disable form while submitting. On success, navigate to `/members`. On error, display error message and allow retry.

**Files to create/modify:**
- `services/frontend/src/views/AddMemberView.vue` (new)
- `services/frontend/src/views/__tests__/AddMemberView.spec.ts` (new)
- `services/frontend/src/router.ts` (modify — add route)

**Verification:** `npm test -- AddMemberView` → 2 tests pass; form submits, navigation works, error handling works

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 9: [behavior] Implement EditMemberView and wire to router

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 9

**Description:** Create `EditMemberView.vue` that routes on `/members/:id/edit`, loads the member data from the store (by matching ID to member in `members` array), pre-fills the form, and on submit calls `store.editMember(id, formData)`. Handle case where member is not found (display error, redirect). Show loading spinner during submission.

**Files to create/modify:**
- `services/frontend/src/views/EditMemberView.vue` (new)
- `services/frontend/src/views/__tests__/EditMemberView.spec.ts` (new)
- `services/frontend/src/router.ts` (modify — add route with param)

**Verification:** `npm test -- EditMemberView` → 3 tests pass; form pre-fills, submits, navigation works

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 10: [behavior] Wire deactivate and reactivate actions in MemberList component

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 10

**Description:** Update `MemberList.vue` to emit `deactivate` and `reactivate` events (with member ID) when buttons are clicked. Update the parent `Members.vue` view to listen to these events, show a confirm dialog, and on user confirmation call `store.deactivateMember()` or `store.reactivateMember()`. Show loading spinner during action. Update member row immediately after success (status changes, button options change) without full page reload.

**Files to create/modify:**
- `services/frontend/src/components/MemberList.vue` (modify)
- `services/frontend/src/views/Members.vue` (modify — handle deactivate/reactivate events)
- `services/frontend/src/components/__tests__/MemberList.spec.ts` (modify — add event tests)

**Verification:** `npm test` → all tests pass; deactivate/reactivate buttons work, confirm dialog flows, members update without reload

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 11: [tidy] Add "Add Member" button to AppShell sidebar footer

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 11

**Description:** Update the placeholder "Add Item" button in `AppShell.vue` footer to navigate to `/members/add` (RouterLink or programmatic navigation). Button text: "Add Member". Icon: fa-plus (already present). Only show button when user is on `/members` path (optional; or always show and let it navigate from anywhere).

**Files to create/modify:**
- `services/frontend/src/components/AppShell.vue` (modify)

**Verification:** `npm test` → AppShell still renders; "Add Member" button navigates to `/members/add`

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 12: [behavior] Write Playwright acceptance test: add member happy path

**Status:** pending

**Plan reference:** `.agent/plan.md` — Subtask 12

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

**Files to create/modify:**
- `acceptance-tests/tests/members-crud.spec.ts` (new)

**Verification:** `npm run test:acceptance -- members-crud` → test passes against running docker-compose environment

**Evidence:**
- state: pending
- tests: []
- commits: []

---

### Subtask 14: [refactor] Restyle Members screen components to use Bulma CSS classes

**Status:** complete

**Plan reference:** Discovered during implementation

**Description:** Refactor all Members screen components to replace custom scoped CSS with Bulma framework classes. This eliminates style drift and leverages battle-tested CSS:
- MemberForm.vue: Use Bulma `.field`, `.control`, `.input`, `.select`, `.button` utilities
- MemberList.vue: Use Bulma `.table`, `.tag`, `.button`, `.box` utilities  
- ErrorToast.vue: Use Bulma `.notification`, `.is-danger` positioning utilities
- ConfirmDialog.vue: Use Bulma `.modal`, `.modal-card`, `.modal-background` classes
- MemberTabs.vue: Use Bulma `.tabs` utilities
- Members.vue, AddMemberView.vue, EditMemberView.vue: Use Bulma `.hero`, `.section`, `.container`, `.box` classes

**Files to create/modify:**
- `services/frontend/src/components/MemberForm.vue` (modify) ✅
- `services/frontend/src/components/MemberList.vue` (modify) ✅
- `services/frontend/src/components/ErrorToast.vue` (modify) ✅
- `services/frontend/src/components/ConfirmDialog.vue` (modify) ✅
- `services/frontend/src/components/MemberTabs.vue` (modify) ✅
- `services/frontend/src/views/Members.vue` (modify) ✅
- `services/frontend/src/views/AddMemberView.vue` (modify) ✅
- `services/frontend/src/views/EditMemberView.vue` (modify) ✅
- `services/frontend/src/components/__tests__/ErrorToast.spec.ts` (modify — update selectors) ✅

**Verification:** `npm run test` → all 79 tests passing; `npm run build` → no build errors ✅

**Evidence:**
- state: complete
- tests: 79 tests passing (including ErrorToast refactored selectors)
- commits: (ready to commit)
- build: success (688.11 kB css total)

---

## Progress Summary

- **Total subtasks:** 14
- **Pending:** 11
- **In progress:** 0
- **Complete:** 3

---

## Next Action

Commit refactor changes, then start Subtask 3: `[tidy] Create MemberForm component with validation logic`
