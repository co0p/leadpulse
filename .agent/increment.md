# Increment: Frontend Bootstrapping

## Use Case

When I load the web app in my browser, I want to see a working Vue 3 SPA with a persistent Bulma shell (sidebar, top bar, main content area), a green/red health indicator in the footer confirming backend connectivity, and client-side routes that preserve the shell as I navigate, so that I have a professional, polished foundation to build screens on and confidence the backend is accessible.

---

## Goal

Bootstrap a Vue 3 + Vite SPA with a persistent app shell, health check indicator, and clean JSON-API-only backend; delete all HTML rendering from the Go server.

---

## Branch

`increment/frontend-bootstrap`

---

## Acceptance Criteria

1. **Vue Project Scaffold** — Vue 3 + Vite project exists at `services/frontend/src/` with `package.json`, Vitest, Vue Test Utils, Playwright configured; `npm install` succeeds.

2. **Persistent Shell Layout** — Vue AppShell component (sidebar, top bar, main content, footer) renders with Bulma CSS applied; shell persists across client-side route navigation without full-page reload.

3. **Vue Router Setup** — Vue Router configured with a layout-outlet pattern: `AppShell.vue` wraps all routes; navigating between routes keeps the shell visible.

4. **Health Endpoint Extended** — Backend `/api/health` endpoint returns `{status: "ok", version: "..."}` (version optional; status required); HTTP 200; supports CORS for browser requests.

5. **Health Indicator Component** — Footer displays a small health badge (green checkmark ✓ = backend reachable, red ✗ = unreachable); check occurs once on app load (not on every navigation); component tests verify both states.

6. **SPA Build Integrated** — `npm run build` outputs static assets to `services/frontend/dist/`; Go `embed.FS` includes these assets; `go build` produces binary serving SPA at `/` and API at `/api/*`.

7. **HTML Rendering Removed** — `server/templates/` deleted; no `html/template` imports remain in `server/` package; all handler functions remaining in `server/handler_*.go` are JSON-API-only (no HTML rendering).

8. **Tests Pass** — `go test -race ./...` passes (backend); `npm test` passes (frontend component tests); no regression in existing Go test counts.

9. **Docker Integration Works** — Frontend service builds Vue SPA via `npm run build` during container build; backend and frontend discover each other and health checks work across container network; `make docker-up` succeeds.

10. **Responsive Design** — Shell layout adapts correctly to desktop (≥1024px), tablet (769–1023px), and mobile (≤768px) breakpoints; sidebar is fixed on desktop, overlay/hamburger menu on mobile.

---

## Acceptance-Test Intent

**User journey 1: App loads and health check passes**
- Open browser to `http://localhost:8080/`
- Shell renders (sidebar, top bar visible)
- Footer shows green health checkmark (✓)
- No console errors

**User journey 2: Navigation preserves shell**
- From overview screen, click "Members" sidebar link
- Shell remains visible; main content updates in place
- URL changes; no full-page reload

**User journey 3: Backend unavailable, health indicator shows red**
- Stop backend service while app is loaded
- Health badge changes to red (✗) or disappears (per component design)
- App remains usable; error doesn't crash the page

(Journeys 1–2 are required for acceptance; Journey 3 is optional spike/documentation only.)

---

## Out Of Scope

- **Monthly Input Screen, Members Screen, or any other feature screen** — only the shell and health indicator; content placeholder in main area is OK.
- **Form submission or API integration for data entry** — health check is the only API call this increment makes; other API endpoints are not wired to the UI yet.
- **Authentication or authorization** — assume single-user, trusted environment (per CONSTITUTION).
- **Offline mode or service workers** — keep it simple; app requires connectivity.
- **Performance optimization** — use defaults; optimize after behavior is verified.
- **Branding or custom styling** — stick to Bulma defaults and the patterns in `docs/ui.md`; custom CSS deferred.
- **Accessibility audit beyond WCAG 2.1 AA basics** — keyboard nav, semantic HTML, focus outlines; full audit deferred.

---

## Constitution Constraints

- **Small, focused changes** — this increment is one scope (shell + health); no feature screens.
- **Behavior first** — component tests for shell and health indicator; acceptance tests for end-to-end shell persistence and health check.
- **Presentation and API are separate** — frontend is a pure SPA consuming JSON; backend exposes API only.
- **Service is the API contract** — all backend behavior is JSON endpoints; SPA consumes them.
- **No gold-plating** — build what this increment requires; defer animations, advanced layout, custom CSS.
- **Human review always** — health indicator is purely informational (green/red badge); never a decision point.
- **Dependency direction** — frontend depends on backend JSON; zero knowledge of backend internals.

---

## Roadmap Entry

**Feature name:** Frontend Bootstrapping (Vue 3 + Vite SPA with Bulma Shell and Health Indicator)

**Job story:** When I load the web app in my browser, I want to see a working Vue 3 SPA with a persistent Bulma shell, a green/red health indicator confirming backend connectivity, and client-side routes that preserve the shell as I navigate, so that I have a professional, polished foundation to build screens on and confidence the backend is accessible.

**Move from:** Planned (section "Frontend Bootstrap: Bulma Shell + Health Indicator")

**Move to:** Partial (this increment)

---

## Next Action

User approval of this increment definition. On approval:
1. Create branch `increment/frontend-bootstrap` and switch to it.
2. Load `4dc-plan` skill to write `.agent/plan.md` with ordered, verifiable technical execution plan.

