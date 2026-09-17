# Plan: Frontend Bootstrapping

## Goal

Bootstrap a Vue 3 + Vite SPA with a persistent app shell, health check indicator, and clean JSON-API-only backend; delete all HTML rendering from the Go server.

## Branch

`increment/frontend-bootstrap`

## Approach

This increment scaffolds the full Vue 3 + Vite project (adding Pinia, Vue Router, Vitest, Vue Test Utils, Playwright configuration), replaces the Go template shell with a Vue AppShell component, wires the SPA build output into the Go binary via `embed.FS`, and strips all HTML rendering from the backend (templates deleted, handlers removed). The backend becomes purely JSON-API, and the SPA becomes the UI.

Key architectural boundary: Frontend (Vue SPA) consumes only JSON from backend (`/api/*`); backend has zero knowledge of how the SPA renders responses.

Performance-sensitive path: The health check runs once on app load; subsequent navigations do not hit `/api/health` (no performance regression).

## Design

### Data Models

No new domain models. The health check response shape is extended:

**Backend response: `/api/health`** (currently unchanged)
```json
{
  "status": "ok"
}
```

**Frontend state: `HealthStore` (Pinia)**
```typescript
{
  status: "healthy" | "unhealthy" | "checking"  // state machine
  lastCheckedAt: number                         // Unix milliseconds
  checkError: string | null                    // error message if unhealthy
}
```

**Frontend component prop: `HealthIndicator.vue`**
```typescript
{
  status: "healthy" | "unhealthy" | "checking"  // passed from parent or store
}
```

### Call / Data Flow

1. **App Load:** `main.ts` boots Vue app with Router, Pinia, AppShell layout
2. **AppShell mounts:** `AppShell.vue` mounted hook calls `HealthStore.checkHealth()`
3. **Health Check:** `HealthStore.checkHealth()` → `GET /api/health` (via API client)
4. **Response:** Backend returns `{status: "ok"}` with HTTP 200
5. **UI Update:** HealthStore updates `status` to "healthy"; `HealthIndicator.vue` re-renders (green checkmark ✓)
6. **Route Navigation:** User navigates via sidebar link → Vue Router updates route → `AppShell` stays mounted (shell persists) → main content area updates
7. **No re-check:** Subsequent navigations do not call `/api/health` again (no performance overhead)
8. **Error path:** If `GET /api/health` fails → HealthStore catches error → `status` = "unhealthy"; HealthIndicator shows red ✗

### Error / Edge-case Inventory

| Condition | Expected response | Covered by subtask |
|-----------|-------------------|-------------------|
| Network unreachable on app load | Health check fails silently; badge shows red; app still loads and is usable | Subtask 7 (health check error handling) |
| Backend `/api/health` returns non-200 | Treat as unhealthy; badge shows red | Subtask 7 (health check error handling) |
| Health check request times out (>5s) | Abort request; treat as unhealthy; badge shows red | Subtask 7 (API client timeout config) |
| User navigates before health check completes | Health badge may show "checking" state (spinner) briefly; safe to navigate | Subtask 7 (HealthIndicator component) |
| Multiple simultaneous navigations (fast clicks) | Vue Router prevents race conditions; shell stays mounted; only one route active at a time | Subtask 4 (Vue Router config) |
| Page refresh / F5 | App re-boots; health check runs again; shell re-renders | Subtask 1 (main.ts bootstrap) |

**Gaps:** None identified; all conditions covered.

### Observability Intent

| Event | Level | Data |
|-------|-------|------|
| `app.boot` | info | `timestamp`, environment (dev/production) |
| `health.check_started` | debug | `timestamp`, target URL |
| `health.check_success` | debug | `latency_ms`, `status_code` |
| `health.check_failed` | warn | `error_message`, `latency_ms`, `target_url` |
| `route.navigate` | debug | `from_path`, `to_path`, `latency_ms` |

Console logging only (no external telemetry). Use `console.debug`, `console.warn`. Disable in production via `import.meta.env.PROD` checks or `vue.config` logging level.

### Architecture Delta

**Before:**
- Backend serves Go HTML template at `/` (layout.html)
- Frontend is incomplete Vue scaffold (no router, no pinia)
- HTML templates in `server/templates/`
- Backend has knowledge of shell layout, navigation, styling

**After:**
- Backend serves SPA static files at `/` and `/` redirects to `/index.html` (SPA bootstrap)
- Backend has zero HTML; only JSON API endpoints (`/api/*`)
- Frontend Vue 3 app owns shell layout, navigation, styling (AppShell component)
- Frontend calls backend via `/api/health`, `/api/members`, etc. exclusively
- `server/templates/` deleted; `embed.FS` now points to `services/frontend/dist/` only

**Container view:** No new containers. Frontend and backend containers unchanged; their internal structure and communication patterns change.

---

## Files

| File | Role | Notes |
|------|------|-------|
| `services/frontend/package.json` | modify | add Pinia, Vue Router, Vitest, Vue Test Utils, Playwright, @types/node |
| `services/frontend/vite.config.js` | modify | ensure `build.outDir` = `dist`, `build.emptyOutDir` = true |
| `services/frontend/vitest.config.js` | new | config for Vitest (test runner for Vue components) |
| `services/frontend/playwright.config.js` | new | config for Playwright acceptance tests |
| `services/frontend/tsconfig.json` | new | TypeScript config for Vue 3 + Vite project |
| `services/frontend/src/main.ts` | new | Vue app entry point; boots router, pinia, app shell |
| `services/frontend/src/App.vue` | modify | simplify; delegate to router with layout outlet |
| `services/frontend/src/components/AppShell.vue` | new | persistent shell layout (sidebar, topbar, footer, health indicator, router outlet) |
| `services/frontend/src/components/HealthIndicator.vue` | new | health badge component (green ✓ / red ✗ / spinner) |
| `services/frontend/src/stores/health.ts` | new | Pinia store for health check state and logic |
| `services/frontend/src/stores/index.ts` | new | Pinia store barrel export |
| `services/frontend/src/api/client.ts` | new | API client for `/api/health` and future endpoints |
| `services/frontend/src/router.ts` | new | Vue Router config; layout-outlet pattern |
| `services/frontend/src/views/Home.vue` | new | placeholder home screen (empty content area) |
| `services/frontend/src/__tests__/components/HealthIndicator.spec.ts` | new | Vitest tests for HealthIndicator (healthy, unhealthy, checking states) |
| `services/frontend/src/__tests__/stores/health.spec.ts` | new | Vitest tests for HealthStore (check success, check failure, timeout) |
| `services/backend/server/server.go` | modify | change root handler to serve `index.html`; remove template parsing; wire `/` and `/*` to SPA |
| `services/backend/server/server_test.go` | modify | remove HTML template tests; update assertion for SPA bootstrap (index.html served) |
| `services/backend/server/templates/` | delete | entire directory (layout.html, components/, members/) — superseded by Vue |
| `services/backend/server/handler_members.go` | modify | remove HTML-rendering handlers (`HandlerGetMembersPage`, `HandlerGetAddMemberPage`, `HandlerPostAddMember`, `HandlerGetEditMemberPage`); keep only JSON API handlers |
| `services/backend/server/handler_members_test.go` | modify | remove HTML handler tests; keep JSON API tests |
| `services/backend/main.go` | touch | no change required; re-wire not needed (Go binary embed.FS updated in server.go) |
| `docs/architecture.md` | modify | update container diagram to show Vue SPA instead of templates; update frontend responsibility section to mention AppShell, Vue Router, Pinia |
| `docs/ui.md` | touch | reference for shell design patterns (responsive, accessibility) |
| `docs/adr/ADR-20260917-vue-spa-frontend.md` | touch | rationale for this decision; already written but verify it's current |

---

## Subtasks

### 1. [tidy] Upgrade frontend `package.json` and add build config

**Files:**
- `services/frontend/package.json`
- `services/frontend/vite.config.js`
- `services/frontend/vitest.config.js` (new)
- `services/frontend/playwright.config.js` (new)
- `services/frontend/tsconfig.json` (new)

**References:**
- Current `package.json`: see it has Vue 3, Vite, but missing router/pinia
- `docs/testing.md` — test strategy and commands
- `docs/constitution.md#testing-strategy` — testing pyramid

**Description:**
Add Pinia, Vue Router, Vitest, Vue Test Utils, Playwright to `package.json`. Create Vitest config (test environment: happy-dom, globals: true). Create Playwright config (baseURL: http://localhost:3000). Create TypeScript config for Vue 3 + Vite. Update `vite.config.js` to ensure `build.outDir` = `dist/`, `build.emptyOutDir` = true.

**Verification:** `npm install` succeeds; `npm run build` produces `dist/` directory; `npm test` runs (empty suite OK); `npx playwright --version` returns version.

**Depends on:** —

**Acceptance criteria:** AC-1 (Vue Project Scaffold), AC-8 (Tests Pass)

---

### 2. [tidy] Create Vue app entry point and router config

**Files:**
- `services/frontend/src/main.ts` (new)
- `services/frontend/src/router.ts` (new)
- `services/frontend/index.html` (modify — ensure entry is `src/main.ts`)

**References:**
- Vue 3 docs: https://vuejs.org (entry point pattern)
- Vue Router docs: https://router.vuejs.org (layout outlet pattern)
- Current `index.html`: check for existing script tag pointing to `main.js`

**Description:**
Create `main.ts` that imports Vue, router, pinia, AppShell; creates Vue app; mounts to `#app`. Create `router.ts` with Vue Router config: define routes (Home, placeholder screens); use layout-outlet pattern (AppShell wraps all routes; main content area is `<router-outlet />`). Update `index.html` to reference `src/main.ts`.

**Verification:** `npm run build` succeeds; `npm run preview` boots SPA locally; navigate to `http://localhost:4173` (or configured dev port), shell renders, navigation links exist.

**Depends on:** 1

**Acceptance criteria:** AC-2 (Persistent Shell), AC-3 (Vue Router Setup)

---

### 3. [tidy] Create Pinia health check store

**Files:**
- `services/frontend/src/stores/health.ts` (new)
- `services/frontend/src/stores/index.ts` (new)
- `services/frontend/src/api/client.ts` (new)

**References:**
- Pinia docs: https://pinia.vuejs.org (state, actions, getters pattern)
- Current backend `/api/health` — see `services/backend/server/server.go:99–103`
- Architecture: SPA should call `/api/health` via exported client function

**Description:**
Create `api/client.ts` with `fetchHealth()` function: makes `GET /api/health` request, returns `{status: "ok"}` or throws error on non-200. Create `stores/health.ts` with Pinia store: state (status: "healthy" | "unhealthy" | "checking", lastCheckedAt, checkError), action `checkHealth()` (calls `fetchHealth()`, updates state, catches errors), getters (isHealthy, statusText). Create `stores/index.ts` barrel.

**Verification:** `npm test -- health.spec.ts` runs (will be empty until subtask 8 adds tests).

**Depends on:** 1

**Acceptance criteria:** AC-4 (Health Endpoint), AC-5 (Health Indicator Component)

---

### 4. [tidy] Create AppShell component with router outlet

**Files:**
- `services/frontend/src/components/AppShell.vue` (new)
- `services/frontend/src/views/Home.vue` (new)

**References:**
- Bulma docs: https://bulma.io (CSS classes)
- `docs/ui.md` — shell design patterns (sidebar, topbar, footer, responsive, a11y)
- `docs/architecture.md#containers` — container view showing AppShell in frontend

**Description:**
Create `AppShell.vue`: renders persistent sidebar, topbar, main content area, footer. Includes `<router-outlet />` for screen content. Sidebar has logo, navigation links (Home, Members, etc.), collapsible sections, footer button. Topbar has hamburger toggle, centered search (placeholder), quick action buttons. Footer has health indicator badge. Responsive: sidebar fixed on desktop (≥1024px), overlay/hamburger on mobile. Keyboard-accessible: semantic HTML, ARIA labels, focus outlines (2px #3273dc per `docs/ui.md`).

Create `Home.vue` as content placeholder.

**Verification:** `npm run build && npm run preview` → navigate to home page → shell renders with all 3 regions (topbar, sidebar, main) visible; shell persists on mobile (sidebar toggles); focus outlines visible when tabbing.

**Depends on:** 2, 3

**Acceptance criteria:** AC-2 (Persistent Shell), AC-10 (Responsive Design)

---

### 5. [tidy] Create HealthIndicator component

**Files:**
- `services/frontend/src/components/HealthIndicator.vue` (new)

**References:**
- `services/frontend/src/stores/health.ts` (from subtask 3) — HealthStore state shape
- `docs/ui.md` — icon conventions, color scheme
- Bulma: icon tag, color classes

**Description:**
Create `HealthIndicator.vue`: small badge component (15–20px) showing:
- Green checkmark ✓ when `status === "healthy"` (title: "Backend OK")
- Red ✗ when `status === "unhealthy"` (title: "Backend unreachable")
- Spinner when `status === "checking"` (loading state)

Component takes optional `status` prop (defaults to reading from HealthStore). Placed in footer (bottom-right corner of AppShell).

**Verification:** Component renders; no console errors; props/state binding works.

**Depends on:** 3

**Acceptance criteria:** AC-5 (Health Indicator Component)

---

### 6. [tidy] Wire AppShell health check on mount; integrate into Vue app

**Files:**
- `services/frontend/src/components/AppShell.vue` (modify)
- `services/frontend/src/App.vue` (modify)

**References:**
- Vue 3 Composition API: onMounted hook
- `services/frontend/src/stores/health.ts` — HealthStore.checkHealth() action
- `services/frontend/src/router.ts` — Router export for App.vue to use

**Description:**
Modify `AppShell.vue` `onMounted()` hook to call `HealthStore.checkHealth()` immediately (async; no await in onMounted). Health check runs once at app load. Modify `App.vue` to render AppShell as layout wrapper; App delegates to router.

**Verification:** `npm run build && npm run preview` → open http://localhost:4173 → browser DevTools Network tab shows one `GET /api/health` request on page load; health badge updates from "checking" to "healthy" (or "unhealthy" if backend not running); subsequent sidebar navigation does not trigger another health check.

**Depends on:** 4, 5

**Acceptance criteria:** AC-5 (Health Indicator Component)

---

### 7. [research] Verify SPA routing and embed.FS integration for Go binary

**Files:**
- `services/backend/server/server.go` (read only)
- `services/frontend/vite.config.js` (read only)
- `docs/architecture.md` (read only)

**References:**
- Go `embed.FS` docs: https://pkg.go.dev/embed
- Vite build output: `dist/` directory structure (index.html, assets/, etc.)

**Description:**
Spike: verify how Go embed.FS will serve Vue SPA. In particular:
1. Does `index.html` need to be at root of `dist/`?
2. How does SPA routing (Vue Router client-side) interact with Go static file serving? (Answer: SPA bootstrap via `index.html` for all non-API routes; Vue Router takes over in browser)
3. Current Dockerfile build order: `npm run build` → `go build` → does Go binary include `dist/` from frontend build?

Record findings in implementation.md.

**Verification:** Findings recorded; no blocker for next subtask.

**Depends on:** 2

**Acceptance criteria:** AC-6 (SPA Build Integrated)

---

### 8. [tidy] Update backend server.go to serve SPA; remove template parsing

**Files:**
- `services/backend/server/server.go` (modify lines 1–130)

**References:**
- Current root handler: line 106–111 (serves template)
- Embed.FS: line 13–14 (currently points to `dist templates`)
- Static file serving: line 113–121 (already serves `/dist/`)
- Go http.FileServer: https://pkg.go.dev/net/http#FileServer

**Description:**
1. Change `embed.FS` comment from `dist templates` to just `dist` (or update path logic)
2. Remove `template.ParseFS()` call (line 93–96)
3. Remove template execution from root handler (line 106–111)
4. New root handler: serves `index.html` for all non-API routes (SPA bootstrap pattern); handles `GET /` and `GET /*` by reading `index.html` from embed.FS and writing to `w`
5. Ensure `GET /api/*` routes are untouched (JSON API handlers remain)
6. Remove reference to `tmpl` variable (no longer used)

**Verification:** `go build ./...` succeeds; no unused variable errors. `go test -race ./...` runs (tests updated in next subtask).

**Depends on:** 7

**Acceptance criteria:** AC-6 (SPA Build Integrated), AC-7 (HTML Rendering Removed)

---

### 9. [tidy] Remove HTML template files and HTML rendering handlers from backend

**Files:**
- `services/backend/server/templates/` (delete entire directory)
- `services/backend/server/handler_members.go` (modify — remove HTML handlers)
- `services/backend/server/handler_members_test.go` (modify — remove HTML handler tests)

**References:**
- `services/backend/server/templates/layout.html` — currently embedded
- `services/backend/server/handler_members.go:68–89` — HTML handlers (HandlerGetMembersPage, HandlerGetAddMemberPage, HandlerPostAddMember, HandlerGetEditMemberPage)
- `services/backend/server/handler_members_test.go` — tests for those handlers

**Description:**
1. Delete `services/backend/server/templates/` and all contents
2. From `handler_members.go`, remove routes:
   - `GET /members` (HandlerGetMembersPage)
   - `GET /members/add` (HandlerGetAddMemberPage)
   - `POST /members` (HandlerPostAddMember)
   - `GET /members/{id}/edit` (HandlerGetEditMemberPage)
   - Keep only JSON API routes: `POST /api/members`, `GET /api/members`, `PATCH /api/members/{id}`, `DELETE /api/members/{id}`, `PATCH /api/members/{id}/reactivate`
3. Remove the handler functions themselves (HandlerGetMembersPage, HandlerGetAddMemberPage, HandlerPostAddMember, HandlerGetEditMemberPage)
4. From `server.go`, remove HTML handler routing (line 67–89)
5. From `handler_members_test.go`, remove tests for HTML handlers
6. Verify no `html/template` imports remain in `server/` package

**Verification:** `grep -r "html/template" services/backend/server/ --include="*.go"` returns nothing (except comments if any). `go build ./...` succeeds. `go test -race ./server/...` runs.

**Depends on:** 8

**Acceptance criteria:** AC-7 (HTML Rendering Removed), AC-8 (Tests Pass)

---

### 10. [tidy] Update server_test.go: replace HTML template tests with SPA bootstrap assertions

**Files:**
- `services/backend/server/server_test.go` (modify)

**References:**
- Current tests: likely check for template rendering or HTML structure
- New assertion: GET `/` returns status 200 and includes `index.html` content (e.g., contains `<div id="app">`, `<script>` references)

**Description:**
Remove or rewrite tests that check for HTML template rendering. Add/update test to verify:
1. `GET /` returns HTTP 200
2. Response contains HTML (content-type: text/html)
3. Response body contains `<div id="app">` (Vue mount point) and script tag referencing `main` or Vue runtime

Verify existing API health check test remains green.

**Verification:** `go test -race ./server/...` passes; all tests green.

**Depends on:** 9

**Acceptance criteria:** AC-8 (Tests Pass), AC-9 (Docker Integration)

---

### 11. [tidy] Add Vitest component tests for HealthIndicator and HealthStore

**Files:**
- `services/frontend/src/__tests__/components/HealthIndicator.spec.ts` (new)
- `services/frontend/src/__tests__/stores/health.spec.ts` (new)

**References:**
- Vitest docs: https://vitest.dev
- Vue Test Utils: https://test-utils.vuejs.org
- `services/frontend/src/components/HealthIndicator.vue` (component under test)
- `services/frontend/src/stores/health.ts` (store under test)

**Description:**
Write component tests for `HealthIndicator.vue`:
- Test: renders green checkmark when `status === "healthy"`
- Test: renders red X when `status === "unhealthy"`
- Test: renders spinner when `status === "checking"`

Write store tests for `HealthStore`:
- Test: `checkHealth()` sets `status` to "checking", then to "healthy" on success
- Test: `checkHealth()` sets `status` to "unhealthy" on network error
- Test: `checkHealth()` sets `checkError` message on error

**Verification:** `npm test` passes; all 6 tests green (or however many written).

**Depends on:** 3, 5

**Acceptance criteria:** AC-8 (Tests Pass)

---

### 12. [behavior] Update docs/architecture.md with Vue SPA design

**Files:**
- `docs/architecture.md` (modify)

**References:**
- Current architecture diagram (line 8–80): shows "Go html/template" → needs to show "Vue 3 SPA"
- Current container descriptions (line 85–126): update Frontend responsibility to mention AppShell, Vue Router, Pinia
- Hexagonal section (line 139–270): no changes needed; backend structure unchanged

**Description:**
Update container diagram to replace "Go html/template + HTMX" with "Vue 3 + Vite SPA". Update Frontend Service description to explain:
- AppShell component owns shell layout (sidebar, topbar, footer, responsive, a11y)
- Vue Router handles client-side navigation (shell persists, no full reload)
- Pinia stores hold client state (health check, member list, etc.)
- Health indicator calls `/api/health` once on load
- All API calls through `/api/*` JSON endpoints

No changes to backend structure or hexagonal architecture (those remain unchanged).

**Verification:** `docs/architecture.md` reads clearly; diagram is accurate; someone new can understand the SPA + API separation.

**Depends on:** 4, 5, 8

**Acceptance criteria:** AC-2 (Persistent Shell), AC-3 (Vue Router Setup), AC-5 (Health Indicator), AC-6 (SPA Build Integrated)

---

### 13. [tidy] Update Dockerfile build order for Vue SPA

**Files:**
- `services/frontend/Dockerfile` (modify)

**References:**
- Current Dockerfile: likely does `npm install` and `npm run build`
- Backend Dockerfile: expects `services/frontend/dist/` to exist after frontend build succeeds

**Description:**
Verify Dockerfile runs `npm install` → `npm run build` → output goes to `dist/`. Ensure build stage completes successfully before final image runs (no npm runtime needed in final image; only static files).

**Verification:** `docker build -f services/frontend/Dockerfile .` succeeds; built image is lean (no Node.js in runtime layer if using multi-stage).

**Depends on:** 2

**Acceptance criteria:** AC-9 (Docker Integration)

---

### 14. [tidy] Update backend Dockerfile to embed Vue SPA build output

**Files:**
- `services/backend/Dockerfile` (modify)

**References:**
- Build order: must build frontend first → then build backend Go binary with embedded assets
- Current Go code: `embed.FS` in `server/server.go` points to `dist/` directory

**Description:**
Update backend Dockerfile (multi-stage):
1. Stage 1: frontend build (copy `services/frontend/` → run `npm install && npm run build` → output to `dist/`)
2. Stage 2: Go build (copy all backend code + frontend `dist/` → `go build` → binary includes embedded assets)
3. Stage 3: runtime (Alpine + binary only)

Ensure Go build can access `services/frontend/dist/` at embed time.

**Verification:** `docker build -f services/backend/Dockerfile .` succeeds; inspect final image size (should include static assets, no Node.js).

**Depends on:** 13

**Acceptance criteria:** AC-9 (Docker Integration)

---

### 15. [behavior] Run full docker-compose end-to-end test

**Files:**
- `docker-compose.yml` (touch — reference for build order)
- `services/frontend/dist/` (generated artifact from subtask 2)
- `services/backend/leadpulse` (generated binary)

**References:**
- `docker-compose.yml`: orchestration of frontend and backend services
- `CONSTITUTION.md#testing-strategy` — acceptance test gate

**Description:**
End-to-end verification:
1. Run `make docker-build` (or `docker-compose build`) → both images build successfully
2. Run `make docker-up` (or `docker-compose up`) → both services start, health checks pass
3. Open browser to `http://localhost:3000` (frontend)
4. Verify shell renders (sidebar, topbar, main content, footer with health indicator)
5. Verify health indicator is green (✓) if backend is reachable
6. Click sidebar link → verify shell persists, content updates (no full reload)
7. Open DevTools Network tab → verify `GET /api/health` appears once on load; no re-request on navigation
8. Stop backend service (docker-compose down) → refresh page → health indicator turns red (✗)

**Verification:** All 8 steps pass. `docker-compose logs` show no errors. Browser console has no JS errors (except expected development warnings).

**Depends on:** 14

**Acceptance criteria:** AC-1–AC-10 (all acceptance criteria)

---

## Context Map

- `CONSTITUTION.md#engineering-principles` — behavior-first, small focused changes, no gold-plating
- `CONSTITUTION.md#testing-strategy` — Vitest for component tests, Playwright for acceptance tests
- `CONSTITUTION.md#performance-envelope` — health check response < 200ms (call once on load, not on nav)
- `docs/architecture.md#containers` — C4 Level 2 container view (frontend/backend boundary)
- `docs/architecture.md#hexagonal-architecture` — backend is unchanged; only frontend changes
- `docs/ui.md` — shell design (sidebar, topbar, footer, responsive, a11y), color scheme, icon conventions
- `docs/adr/ADR-20260917-vue-spa-frontend.md` — rationale for Vue 3 + Vite choice (current decision record)
- `docs/adr/ADR-20260916-hexagonal-architecture-member-crud.md` — backend hexagonal pattern (unchanged)
- `docs/testing.md` — test commands and strategy by layer

---

## Acceptance Scenarios

These scenarios are advisory; the acceptance criteria gate is the test suite passing + end-to-end Docker verification (subtask 15).

### Scenario 1: App loads and health check passes

- **Criterion:** AC-1, AC-2, AC-5, AC-9
- **User action:** Open browser to `http://localhost:3000`
- **Precondition:** Docker services healthy; backend running
- **Expected outcome:** 
  - Shell renders (sidebar visible on desktop, hamburger on mobile)
  - Top bar with search and action buttons visible
  - Main content area (home placeholder) visible
  - Footer with green health checkmark (✓)
  - Browser console clean (no errors, only dev warnings)
- **Evidence:** Manual browser inspection (no automated assertion needed; covered by component tests + docker-compose health checks)
- **Gate:** advisory

### Scenario 2: Navigation persists shell

- **Criterion:** AC-2, AC-3
- **User action:** From home, click "Members" link in sidebar
- **Precondition:** App loaded; sidebar fully visible (desktop) or toggle expanded (mobile)
- **Expected outcome:** 
  - URL changes (e.g., `/` → `/members`)
  - Shell (sidebar, topbar, footer) remains visible
  - Main content area updates with new screen (not yet implemented; placeholder OK)
  - No full page reload (browser tab title does not flash, network tab shows only XHR for data, not HTML)
- **Evidence:** Manual browser + DevTools inspection (covered by component tests for router)
- **Gate:** advisory

### Scenario 3: Mobile responsiveness

- **Criterion:** AC-10
- **User action:** Resize browser to mobile width (≤768px); toggle sidebar with hamburger
- **Precondition:** App loaded
- **Expected outcome:** 
  - Sidebar not visible by default (overlay mode)
  - Hamburger menu icon visible in topbar
  - Click hamburger → sidebar slides in from left with semi-transparent overlay
  - Click overlay or ESC → sidebar closes
  - Topbar remains sticky at top
- **Evidence:** Manual browser resize or device emulation (DevTools)
- **Gate:** advisory

---

## Risks

1. **Embed.FS path confusion:** Go `embed.FS` directive must point to correct path. If `services/frontend/dist/` doesn't exist at Go build time, build fails. Mitigation: run `npm run build` before `go build` in Dockerfile; validate in subtask 7.

2. **SPA routing vs. Go static file serving:** Vue Router client-side navigation can conflict with Go static file handler if misconfigured. For example, if Vue app requests `/members` (client-side route), it must first fetch `index.html` and let Vue Router take over, not serve a 404. Mitigation: root handler must serve `index.html` for all non-API routes (SPA bootstrap pattern). Subtask 8 handles this.

3. **Health check performance:** If health check times out, it could block app bootstrap (bad UX). Mitigation: make health check async and non-blocking; don't wait for it before rendering shell. Subtask 6 ensures this (no await in onMounted).

4. **CORS for cross-container calls:** Frontend (port 3000) calls backend (port 8080 inside container, but hostname is `backend:8080`). If nginx proxy not configured correctly, may fail. Mitigation: verify `docker-compose.yml` sets `VITE_API_URL=http://backend:8080` and nginx proxies `/api/*` to backend. Subtask 14 validates.

5. **TypeScript setup complexity:** Adding TypeScript to Vue project introduces new tooling (tsconfig.json, type definitions). If misconfigured, IDE intellisense breaks. Mitigation: use standard Vue 3 + Vite + TypeScript template; test with `npm run build` early (subtask 1).

---

## Planning Decisions

| Decision | Chosen | Rejected | Reason |
|----------|--------|----------|--------|
| **When to call health check** | Once on app load (onMounted) | On every navigation or periodically | Simplicity; one API call per session; performance envelope met. Polling can be added later if needed. |
| **Health indicator state machine** | Three states: "healthy", "unhealthy", "checking" | Two states (healthy/unhealthy only) | Shows loading feedback during the check (better UX). Checking state brief (usually <100ms), so minimal visual clutter. |
| **Template removal strategy** | Delete all HTML templates; move shell to Vue AppShell | Keep templates and serve side-by-side | Cleaner architecture; no mixed template engines; full separation of concerns (backend JSON only, frontend Vue only). |
| **SPA build integration** | Go embed.FS includes frontend `dist/` output | Frontend service serves SPA separately (no embedding) | Single-binary distribution (per CONSTITUTION). Embedding is complex but critical for portability. |
| **Router layout pattern** | AppShell wraps all routes via layout-outlet | Separate layout component for each page | Simpler; single source of truth for shell; shell persists automatically. Reduces boilerplate per screen. |
| **Pinia for health state** | Yes, use Pinia store | Plain Vue reactive/ref in component | Stores are reusable (other screens may need health); encourages separation of concerns. Slight overhead but pays off as app grows. |
| **Observability logging** | Console.debug/warn only | External service or no logging | Stateless SPA; no backend telemetry (per CONSTITUTION privacy). Console logs are dev/debug aids; can be disabled in production. |

---

## Next Action

User approval of this plan. On approval:
1. Hand off to the implement workflow (orchestrator detects first subtask type and calls appropriate implement skill).
2. First subtask is `[tidy]` → `4dc-tidy` skill handles it.
3. Implement skills load only the files named in their subtask's Files field; they do not re-scan the whole codebase.
