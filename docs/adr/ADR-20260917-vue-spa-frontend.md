# ADR-20260917 — Vue 3 SPA as the Frontend, Backend Restricted to a JSON API

**Decision:** Replace the HTMX + Alpine.js + Go `html/template` frontend with a dedicated Vue 3 SPA (Vite build, Pinia for state, Bulma retained for styling). The Go backend becomes API-only — it no longer renders HTML for any feature screen. The compiled SPA (`dist/`) is embedded into the Go binary via `embed.FS` to preserve single-binary distribution.

**Status:** Accepted — supersedes `ADR-20260915-spa-frontend-stack.md`.

---

## Context

`ADR-20260915-spa-frontend-stack.md` chose HTMX + Alpine.js + Go `html/template` to avoid a Node.js build step and to optimize for AI-authoring accuracy. In practice, implementing the Members screen (list with status tabs, add/edit/reactivate) surfaced a structural problem: composing a persistent application shell (sidebar, top bar) with per-screen content required either template partial composition inside Go handlers or HTMX out-of-band swaps. Both approaches produced handlers that built HTML by string concatenation or juggled multiple `ParseFS` calls, which violated the "handlers are thin adapters" rule in `CONSTITUTION.md` and made shell/content composition fragile — a bug shipped where `/members?status=active` rendered without the shell at all.

The PRD (`docs/prd.md`) defines six screens (Overview Dashboard, Monthly Input Workspace, Member Detail, Alerts & Risk Center, Monthly Review & Calibration, Limited Settings), several with live-updating panels (score preview, trend charts, distribution charts) and shared cross-screen state (active member, active cycle). This is exactly the class of "rich client-side interaction" that `ADR-20260915` flagged as a weak point for the HTMX/Alpine approach.

The constraints from `CONSTITUTION.md` that still apply:
- Single binary distribution — no separate server process, no runtime dependency for the end user.
- Privacy by design — no CDN calls, no telemetry, no external network calls at runtime.
- Presentation is replaceable — the backend must expose the same functionality to any client (SPA, CLI, future integrations) through a stable contract.
- Dependency direction is one-way; the backend has zero knowledge of how its output is rendered.

The Members screen shell composition bug is a symptom, not a wart to patch. It is the boundary is wrong: the backend was rendering HTML at all is the underlying issue, not which templating trick is used to stitch it together.

---

## Alternatives

### Vue 3 + Vite SPA, backend as pure JSON API (chosen)
The backend exposes only `/api/*` JSON endpoints. A separate Vue application owns all rendering, client-side routing, and shell composition (layout persists across route changes by construction — Vue Router + a layout component, not server-side template stitching). Built with Vite; output embedded into the Go binary via `embed.FS`.

- **Build pipeline:** `npm run build` produces `dist/`; that directory is embedded into the Go binary before `go build`. Node.js is a build-time dependency only, never a runtime dependency for the shipped binary.
- **Testing:** Vitest + Vue Test Utils for component-level and edge-case coverage; Playwright reserved for a small number of end-to-end acceptance tests, one per feature, covering only the main success flow (e.g., "add a member" happy path) — not exhaustive scenario coverage.
- **Shell composition:** Solved natively by Vue Router's nested routes and a persistent `AppShell.vue` layout component. No server-side template stitching, no HTMX swap targets to keep in sync.
- **Trade-off:** introduces a Node.js/npm toolchain requirement for development and CI (not for the shipped binary). Larger initial learning surface than HTMX for simple forms.

### Keep HTMX + Alpine.js + Go `html/template` (status quo, rejected)
Continue composing shell and screen content via Go templates and HTMX swaps.

- **Trade-off:** the shell/content composition problem is structural, not incidental. Every new screen with shared layout risks the same class of bug. Correctly solving it requires re-implementing client-side routing/layout persistence in Alpine.js by hand, which is exactly what a SPA framework already provides.

### React + Vite (rejected)
Comparable technical fit to Vue for this problem (SPA framework, component-based, client-side routing). Rejected on team preference for Vue's single-file-component model and gentler API surface, not on technical grounds.

### No framework, hand-rolled client-side router (rejected)
Vanilla JS with a manual router and DOM diffing. Rejected — reinvents the core value a SPA framework provides (routing, reactivity, component composition) with more code to author and review, higher risk of subtle bugs, no material benefit over adopting Vue given Node.js is already an accepted build-time dependency once any bundler is introduced.

---

## Rationale

1. **The shell/content composition problem needs a client-side router, not a better templating trick.** Vue Router with a persistent layout component solves "shell always wraps content" as a structural guarantee, not a convention every handler must remember to follow.

2. **Backend-as-API-only restores the dependency direction the constitution already requires.** `server/` handlers currently mix two responsibilities: business-logic adaptation (parse → delegate to use case → serialize) and presentation (render HTML, choose templates, manage shell state). Removing HTML rendering from the backend collapses every handler back to the "15–20 line thin adapter" shape documented in `docs/architecture.md`, and makes "presentation is replaceable" true in practice, not just in principle — a CLI or a future mobile client can call the identical `/api/*` contract the Vue SPA uses.

3. **Playwright scope stays deliberately narrow.** Vitest + Vue Test Utils cover component logic, validation, and error states without a browser. Playwright is reserved for one acceptance test per feature verifying the main success flow end-to-end (e.g., adding a member end-to-end through the real UI and API). This keeps the expensive, slow test layer small and avoids the maintenance burden of a large browser-test suite while still proving the full stack works for the flows that matter most.

4. **Single-binary distribution is preserved, not sacrificed.** The Vue build output is a static `dist/` directory; embedding it via `embed.FS` is mechanically identical to how Bulma/Alpine assets are embedded today. Node.js becomes a build-time dependency (already true for any team using a bundler) but never a runtime dependency for the distributed binary.

---

## Consequences

**Better:**
- Handlers return to pure API-adapter shape: parse request → call use case → serialize JSON. No HTML, no template state, no shell-composition responsibility in `server/`.
- Shell/content composition is a solved problem by construction (Vue Router nested layouts), eliminating the entire class of bug that prompted this decision.
- Client-side state (active member, active cycle, live score preview) has a real state management story (Pinia) instead of ad hoc Alpine.js `x-data` scattered across templates.
- The backend contract is verifiably reusable: any client capable of HTTP + JSON can drive the full application.

**Harder:**
- Node.js/npm/Vite become required in the development and CI toolchain (build-time only). `go build` alone no longer produces a working binary; a `npm run build` step must precede it.
- `server/templates/*` and their handler tests (HTML-assertion style, `strings.Contains` on rendered markup) are retired; new frontend tests live in the Vue project (Vitest/VTU) and are not part of the Go test suite.
- Two codebases (Go backend, Vue frontend) instead of one, with a versioned API contract between them that must be kept in sync manually (no shared-type generation in v1).
- Slightly higher up-front setup cost for each new screen (component + route + store, vs. a single `.html` template file).

---

## Related

- `ADR-20260915-spa-frontend-stack.md` — superseded by this decision. Status updated to `Superseded`.
- `ADR-20260914-ui-controllers.md` — the controller/use-case boundary this decision relies on remains unchanged; use cases are unaffected by this frontend change.
- `ADR-20260916-hexagonal-architecture-member-crud.md` — the hexagonal core/use-case/repository structure is untouched; only the `server/` presentation responsibility changes.
- `docs/roadmap.md` — Frontend Migration section sequences the Vue SPA build-out screen by screen.
