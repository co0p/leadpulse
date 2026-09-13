# ADR-20260913 — Use Go and Fyne v2 as the Desktop UI Framework

**Decision:** Build Team Impact Scorecard as a native desktop application using Go and Fyne v2.

**Status:** Accepted

---

## Context

The PRD requires a standalone tool with no external data. The Team Lead runs it on their own machine. The tool handles sensitive performance data that must not leave the device.

The primary deployment question was: web app, Electron, or native desktop?

The key constraints:
- Standalone: no server, no cloud, no network calls
- Privacy: all data on-device, no third-party runtime calling home
- Cross-platform: Team Leads may be on macOS, Windows, or Linux
- Single binary: no installer complexity, no runtime dependency for the user
- The team is building in Go (engine and store are Go)

---

## Alternatives

**Web app (React/Vue + local backend)**
- Requires running a local server process alongside the UI, or bundling a server into an Electron-like shell
- HTML/CSS/JS is a different runtime from the Go core — two languages to maintain
- Browser security model complicates local file access (SQLite)
- More familiar frontend tooling

**Electron (Node.js + Chromium)**
- Familiar web tech for UI
- ~150MB binary baseline; ships Chromium
- Node.js runtime is a large, separate dependency from the Go core
- Chromium makes outbound connections (telemetry, font fetching) by default — requires active suppression
- CGO bridge needed to call Go engine from JS, or duplicate formula logic in JS

**Fyne v2 (Go)**
- Pure Go: UI, engine, and store are the same language and toolchain
- Single binary, no runtime dependency for the user
- Cross-platform (macOS, Windows, Linux) from one codebase
- ~10MB binary baseline
- No background network connections
- Fyne's widget library is adequate for the screens defined in the PRD
- Smaller ecosystem, less UI polish than web frameworks

**Qt (Go bindings via `therecipe/qt`)**
- Mature, polished UI
- Build toolchain is complex; CGO required; binding layer is fragile
- Significant overhead for a single-user tool

---

## Rationale

Fyne wins because it satisfies the hard constraints (standalone, privacy, single binary, Go) without introducing a second language or a background network runtime. The PRD does not require browser-grade UI polish — it requires functional, readable data views. Fyne's widget set covers that.

Electron would satisfy the cross-platform and distribution requirements but introduces Chromium (privacy risk, binary size, network by default) and forces a language split between the UI and the Go engine.

A web app without Electron has no viable path to a truly standalone single-binary distribution without significant wrapper engineering.

---

## Consequences

**Better:**
- Single language and toolchain for the entire application
- Single binary distribution; no runtime to install
- No background network connections by design
- Smaller binary than Electron alternatives

**Harder:**
- Fyne has limited automated UI testing support — UI verification is manual
- Fyne's widget ecosystem is smaller than web component libraries; some custom rendering may be needed (e.g., heatmap, sparkline charts)
- Fyne uses OpenGL/Metal for rendering — requires a GPU driver; unusual display setups (headless VMs, remote desktop) may have issues
- Developers unfamiliar with Fyne will have a learning curve

---

## Related

- [ADR-20260913-sqlite-local-storage.md](ADR-20260913-sqlite-local-storage.md)
- [ADR-20260913-pure-go-sqlite-driver.md](ADR-20260913-pure-go-sqlite-driver.md)
