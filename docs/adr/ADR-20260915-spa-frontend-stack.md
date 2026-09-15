# ADR-20260915 — HTMX + Alpine.js + Go html/template + Bulma as the SPA Frontend Stack

**Decision:** Use HTMX, Alpine.js, Go `html/template`, and Bulma CSS to build the browser-based frontend served from the Go binary. React, Vue, and plain-JS alternatives are not adopted.

**Status:** Accepted

---

## Context

The project is migrating from a Fyne desktop GUI to a locally-served web frontend (see `docs/roadmap-spa.md`). A frontend stack must be chosen before Increment 4 (SPA Scaffold). The choice governs:

- How the Go binary embeds and serves UI assets
- How AI-assisted development (Claude models) writes and modifies UI code
- How UI tests are written and how painful they are to maintain
- Whether a Node.js build step is required in the binary pipeline

The existing UI is already decoupled: `ui/controllers/` is pure Go with zero Fyne imports. Controllers expose plain method calls that map cleanly onto HTTP handler inputs. No business logic needs to move.

The primary constraints from `CONSTITUTION.md` that apply:
- Single binary distribution — no separate process, no runtime dependencies
- Privacy by design — no CDN calls, no external JS delivery at runtime
- UI is replaceable — the same `service/` interfaces must remain callable
- Testing must be reliable; the Fyne headless test driver was cited as a friction point

The secondary constraint driving stack choice: Claude models will write most of the UI code. A stack with fewer abstraction layers, less build-time magic, and more markup-visible behavior produces more accurate AI output and easier human review.

---

## Alternatives

### HTMX + Alpine.js + Go `html/template` + Bulma (chosen)
Server renders HTML fragments; HTMX swaps them into the DOM on user actions. Alpine.js handles client-side reactivity (live preview debounce, form field state). Templates are `.html` files embedded in the binary alongside Go handlers.

- **Build pipeline:** none. `go build` embeds assets via `embed.FS`. No Node.js, no Vite, no bundler.
- **Testing:** handler tests use `net/http/httptest`; template output is assertable with `strings.Contains` or `golang.org/x/net/html`; end-to-end flows use Playwright (test runner only, not in the binary).
- **AI authoring:** HTMX attributes (`hx-post`, `hx-target`, `hx-swap`) and Alpine directives (`x-data`, `x-on`, `x-bind`) are declarative and visible in markup. Claude can read a template and a handler together in one context window and produce accurate changes.
- **Trade-off:** Alpine.js becomes awkward for deeply nested shared state. Acceptable for the current screen set (forms, tables, live preview panels).

### React + Vite + Bulma
Full SPA framework with a Node.js build step outputting a `dist/` directory embedded in the binary.

- **Build pipeline:** requires Node.js, `npm install`, Vite config, and a `make web` step before `go build`.
- **Testing:** Vitest + Testing Library + Playwright; excellent ecosystem; tests are easier to read but require more setup.
- **AI authoring:** React component trees are accurate for Claude, but context includes JSX, hooks, and state topology. Larger surface area for subtle bugs in AI-generated code.
- **Trade-off:** adds a Node.js toolchain dependency to every contributor's environment and CI pipeline. Contradicts the "no runtime dependencies" spirit even if Node is dev-only.

### Plain JS + Web Components + Bulma
No framework; vanilla JS modules; custom elements for reusable pieces. No build step.

- **Build pipeline:** none.
- **Testing:** Playwright for end-to-end; no component-level test story without additional tooling.
- **AI authoring:** Claude can write plain JS accurately, but larger interactive components (live preview, member picker) require manual event wiring that is verbose and error-prone to generate.
- **Trade-off:** lowest dependency count, highest authoring friction for interactive components.

---

## Rationale

HTMX + Alpine.js wins for three compounding reasons:

1. **No build step in the binary pipeline.** `go build` stays the single command that produces a working binary. This is the lowest-friction path for both local development and CI, and it keeps `CONSTITUTION.md`'s single-binary constraint trivially satisfied.

2. **AI authoring accuracy.** HTMX behavior is expressed as HTML attributes on the elements that trigger it. A Claude model reading a template and its corresponding Go handler has full context in one pass. There is no component state topology to reconstruct, no hook dependency array to reason about. Errors in AI-generated HTMX code are visible in the markup and easy to spot in review.

3. **Handler-level testability is sufficient.** The controller → handler boundary already enforces that no business logic enters the template layer. `httptest` covers all correctness-critical paths. Playwright covers full-flow acceptance scenarios. React's component test story offers no additional coverage for this project's actual risk surface.

Alpine.js is chosen over plain JS for the live-preview interaction specifically: a debounced `fetch` to `/api/entries/preview` on field change is a natural fit for `x-on:input.debounce` + `x-data`. The alternative (manual `addEventListener` + `fetch` wiring) would be more verbose and more brittle in AI-generated code.

Bulma is a CSS-only framework. It imposes no JS opinions and works identically with HTMX or any other JS approach. No change to this decision is required if Alpine.js is later supplemented or replaced.

---

## Consequences

**Better:**
- `go build` remains the only build command; no Node.js in CI or contributor setup
- Template changes are visible in `.html` files alongside Go handlers — one context window for AI and human review
- Handler unit tests cover all correctness-critical paths without a browser
- Bulma's class-based styling works directly in Go templates with no PostCSS or purge configuration

**Harder:**
- Rich client-side interactions beyond forms and tables (e.g. drag-and-drop calibration, real-time multi-user updates) would require either heavy Alpine.js workarounds or a stack migration
- No hot-module replacement; template changes require a Go recompile and browser refresh (mitigated by fast `go build` times and air/reflex live-reload tooling)
- Alpine.js is a runtime CDN dependency by default — must be self-hosted or inlined into the binary to satisfy the privacy-by-design constraint (resolved at Increment 1 by embedding the minified JS alongside templates)

---

## Related

- `ADR-20260913-go-fyne-desktop.md` — original desktop stack decision; this ADR supersedes it for the UI layer. The Fyne ADR's status remains Accepted until Increment 7 (Remove Fyne Dependency) is complete.
- `ADR-20260914-ui-controllers.md` — controller pattern that makes the handler → controller boundary clean and enables handler-level testability without UI framework setup
- `docs/roadmap-spa.md` — the full migration plan this decision enables
