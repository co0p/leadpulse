# Plan: Docker Services Architecture & Acceptance Testing Foundation

## Goal

Restructure the project into a multi-service architecture with separate frontend and backend services in Docker containers, establish an acceptance test suite that runs against the full containerized system, and create a docker-compose orchestration file that boots the entire system locally.

## Branch

`increment/docker-services-architecture`

## Approach

Move the Go backend from root to `services/backend/` and create a new Vue 3 + Vite frontend in `services/frontend/`, both with their own Dockerfiles. Create a `docker-compose.yml` at root that orchestrates backend, frontend, and database services. Establish acceptance tests in `acceptance-tests/` using Playwright that run against the live containerized system. Update the Makefile with targets for building, composing, and testing. The shift moves from a single-binary desktop app to a containerized multi-service system where services communicate over HTTP/JSON API. The hexagonal architecture of the backend is preserved; the frontend is a separate npm project.

## Design

### Data Models

**New: Docker Service Containers**

```
Backend Service (services/backend/)
{
  name: "backend"
  image: "leadpulse-backend:latest"
  buildContext: "services/backend/"
  port: 8080
  environment: {
    DATABASE_URL: "sqlite:///data/data.db"
  }
  healthCheck: "GET /api/health"
  volumes: ["./data:/data"]  // SQLite file persistence
}

Frontend Service (services/frontend/)
{
  name: "frontend"
  image: "leadpulse-frontend:latest"
  buildContext: "services/frontend/"
  port: 3000
  environment: {
    VITE_API_URL: "http://backend:8080"
  }
  dependsOn: ["backend"]
}

Orchestration (docker-compose.yml)
{
  version: "3.8"
  services: [backend, frontend]
  networks: ["leadpulse-network"]  // Service discovery by hostname
  volumes: ["data"]  // Persistent SQLite volume
}
```

### Call / Data Flow

1. Developer runs `docker-compose up` from project root
2. Docker Compose boots backend service (builds from `services/backend/Dockerfile`, runs Go binary on port 8080)
   - Backend initializes with SQLite at `/data/data.db` (mounted via docker-compose volume)
   - If volume is empty (first run), database is created fresh
   - If volume exists from previous run, existing data persists
3. Docker Compose boots frontend service (builds from `services/frontend/Dockerfile`, runs web server on port 3000)
4. Frontend startup script resolves backend hostname (`backend:8080`) via docker-compose network DNS
5. Frontend Vue app loads in browser at `http://localhost:3000`
6. Vue app makes `GET /api/health` call to `http://backend:8080/api/health`
7. Backend responds with `{status: "ok"}` (HTTP 200)
8. Playwright acceptance test connects to `http://localhost:3000`, verifies app shell, verifies health indicator
9. User runs `make docker-down` to stop services (all containers stopped; data volume persists)
10. User runs `make docker-up` to restart with existing data intact (no rebuild needed; data survives restart)
11. User runs `make test:acceptance` to run acceptance tests against live services 

### Error / Edge-case Inventory

| Condition | Expected Response | Covered by |
|-----------|-------------------|------------|
| docker-compose services fail to start | docker-compose exits with error; developer sees logs | Subtask 5 (docker-compose verification) |
| Backend Dockerfile build fails (tests fail) | Docker build stops at test stage; error shown | Subtask 2 (backend Dockerfile) |
| Frontend Dockerfile build fails (npm install fails) | Docker build fails; error shown | Subtask 4 (frontend Dockerfile) |
| Services start but health check fails | docker-compose waits for timeout; logged as unhealthy | Subtask 5 (health check in compose) |
| Frontend cannot resolve backend hostname | Network error in browser console; Playwright test fails | Subtask 7 (acceptance test verifies connectivity) |
| Acceptance tests run but services not healthy | Tests timeout or fail with connection errors | Subtask 7 (test waits for service health) |
| Makefile targets not found | Developer error; make command fails with clear message | Subtask 8 (Makefile targets) |
| User tries to persist data across docker-compose down/up | Data persists in volume; `docker-compose up` restarts with existing data | Subtask 5 (data volume in docker-compose) |

### Observability Intent

| Event | Level | Data |
|-------|-------|------|
| `docker.service.start` | info | service name, port, image tag |
| `docker.service.health_check` | debug | service name, result (pass/fail), response_time_ms |
| `docker.compose.up` | info | timestamp, services started (backend, frontend, database) |
| `docker.build.backend` | info | image tag, build duration_ms, test result |
| `docker.build.frontend` | info | image tag, build duration_ms, bundle size |
| `test.acceptance.start` | info | test count, environment (docker-compose) |
| `test.acceptance.complete` | info | passed count, failed count, duration_ms |

### Architecture Delta

**Before:** Single Go binary embedding Vue SPA via `embed.FS`. User downloads one file, runs it, gets desktop app with persistent local SQLite database.

**After:** Three containerized services (stateless containers, persistent data):
- `leadpulse-backend` container (port 8080): Go hexagonal architecture, SQLite at `/data/data.db` (volume mount)
- `leadpulse-frontend` container (port 3000): Vue 3 SPA served via web server
- `data` volume: Persistent SQLite file; survives `docker-compose down/up` cycle

**Key principle:** Containers are **stateless** (no data inside container image). Data is **persistent** (survives restarts). **No rollback procedures needed** — always deploy forward. **No schema migrations** — v1 ships with single schema that evolves via code changes only, never via migration scripts.

**Deployment model:** `docker-compose up` always idempotent. If volume exists with data, it's used; if volume is empty, database initializes fresh. Never roll back to previous versions. Never run migrations.

**Dependency direction unchanged:** Backend maintains hexagonal inbound dependency (server → core ← storage/engine). Frontend is a pure HTTP client; zero knowledge of backend internals.

**Docker Compose network:** Services discover each other by hostname (`backend`, `frontend`). All communication is HTTP/JSON over 127.0.0.1 (localhost) from developer perspective.

## Files

| File/Directory | Role | Notes |
|---|---|---|
| `services/backend/` | move | Move entire Go backend here (core/, engine/, storage/, server/, main.go, go.mod, go.sum, Makefile, etc.) |
| `services/backend/Dockerfile` | new | Multi-stage build: test stage (go test), build stage (go build binary) |
| `services/backend/main.go` | modify | Update database path to `/data/data.db` (volume mount); ensure startup does not fail if database already exists |
| `services/frontend/` | new | Create new Vue 3 + Vite project scaffold |
| `services/frontend/package.json` | new | npm project with Vite, Vue 3, Router, Pinia, Bulma, Playwright dev dep |
| `services/frontend/vite.config.ts` | new | Vite configuration for SPA build |
| `services/frontend/tsconfig.json` | new | TypeScript config for Vue + Vite |
| `services/frontend/index.html` | new | HTML entry point for Vite |
| `services/frontend/src/main.ts` | new | Vue app entry point |
| `services/frontend/src/App.vue` | new | Router root component (AppShell layout) |
| `services/frontend/src/router.ts` | new | Vue Router configuration |
| `services/frontend/src/components/AppShell.vue` | new | Persistent shell layout (sidebar, top bar, footer) |
| `services/frontend/Dockerfile` | new | Multi-stage: build stage (npm run build), serve stage (nginx or http-server) |
| `docker-compose.yml` | new | Orchestrates backend, frontend services; sets environment variables; defines `data` volume for SQLite persistence |
| `.gitignore` | modify | Add `node_modules/`, `.env`, `data/` (SQLite volume), `.env.local`, `dist/` (build output) |
| `Makefile` | modify | Add targets: `docker-build`, `docker-up`, `docker-down`, `test:acceptance`, `docker-logs` |
| `acceptance-tests/` | new | Directory for Playwright acceptance tests |
| `acceptance-tests/package.json` | new | npm project with Playwright (TypeScript), test runner config |
| `acceptance-tests/playwright.config.ts` | new | Playwright configuration (baseURL: localhost:3000, timeout: 30s) |
| `acceptance-tests/tests/app.spec.ts` | new | First acceptance test: app shell loads, backend reachable |
| `docs/architecture.md` | touch | Reference only (already updated for multi-service in increment proposal) |
| `docs/testing.md` | new | Document testing strategy: backend tests (go test), frontend tests (npm test), acceptance tests (Playwright) |
| `docs/deployment.md` | new | Document docker-compose local dev workflow, CI build procedure, fresh-start philosophy (no rollback, no migration) |
| `CONSTITUTION.md` | touch | Reference (already updated for multi-service Docker architecture) |

## Subtasks

### 1. Move Go backend to services/backend/ directory

**Type:** `[tidy]`

**Description:** Move entire Go backend codebase from project root to `services/backend/` subdirectory. Preserve all directory structure (core/, engine/, storage/, server/, main.go, go.mod, go.sum). This is a pure refactoring; no behavior changes.

**Files:**
- `services/backend/` (new directory, created by move)
- All Go source files (move from root to `services/backend/`)
- `go.mod`, `go.sum` (move to `services/backend/`)
- `main.go` (move to `services/backend/`, see subtask references for modifications)

**References:**
- Current root directory structure: `core/`, `engine/`, `storage/`, `server/`, `main.go`, `go.mod`
- CONSTITUTION.md#Architecture-Boundaries — hexagonal architecture unchanged
- `main.go:19` — database path retrieval (will be modified in subtask 2)

**Verification:**
```
cd services/backend && go build ./...  # Must compile without errors
cd services/backend && go test ./...   # All tests pass
```

**Dependencies:** None

**Acceptance criteria covered:** AC-1 (folder structure created)

---

### 2. Create services/backend/Dockerfile with multi-stage build

**Type:** `[tidy]`

**Description:** Create a Dockerfile for the Go backend using a multi-stage build pattern:
- **Stage 1 (test):** Run `go test -race ./...` to verify all tests pass before building
- **Stage 2 (build):** Build the final Go binary from `services/backend/`
- Final image should be minimal (scratch or alpine base), include just the binary
- Update `services/backend/main.go` to use `/data/data.db` as the SQLite database path (for docker-compose volume mounts)

**Files:**
- `services/backend/Dockerfile` (new)
- `services/backend/main.go` (modify — lines 54–64, update database path from `os.UserConfigDir()` to `/data/data.db`)

**References:**
- `services/backend/main.go:54` — `getDatabasePath()` function
- Docker multi-stage build best practices
- CONSTITUTION.md#Testing-Strategy — tests must pass before build

**Verification:**
```
docker build services/backend/ -t leadpulse-backend:latest --target test
# Test stage must pass; if any go test fails, build stops

docker build services/backend/ -t leadpulse-backend:latest
# Final image builds; binary runs and responds to --help or similar verification
```

**Dependencies:** Subtask 1

**Acceptance criteria covered:** AC-2 (backend Dockerfile with health check)

---

### 3. Scaffold services/frontend/ Vue 3 + Vite project

**Type:** `[tidy]`

**Description:** Create a new Vue 3 + Vite project in `services/frontend/` with the following minimal structure:
- `package.json` with dependencies: Vue 3, Vite, Vue Router, Pinia, Bulma CSS, TypeScript
- `vite.config.ts` with appropriate defaults
- `tsconfig.json` for TypeScript
- `index.html` as the HTML entry point
- `src/main.ts` as the Vue app entry point (bootstrap Vue app and Router)
- `src/App.vue` as the root component (renders AppShell layout)
- `src/router.ts` as the Vue Router configuration (basic setup, empty routes for now)
- `src/components/AppShell.vue` as the persistent shell layout (basic structure: sidebar, top bar, footer)
- `public/` directory for static assets (favicon, etc.)

This is scaffolding only; actual screen implementation comes in future increments.

**Files:**
- `services/frontend/package.json` (new)
- `services/frontend/vite.config.ts` (new)
- `services/frontend/tsconfig.json` (new)
- `services/frontend/index.html` (new)
- `services/frontend/src/main.ts` (new)
- `services/frontend/src/App.vue` (new)
- `services/frontend/src/router.ts` (new)
- `services/frontend/src/components/AppShell.vue` (new)
- `services/frontend/public/` (new directory)

**References:**
- Vue 3 setup guide: https://vuejs.org (reference only)
- Vite config docs (reference only)
- CONSTITUTION.md#Multi-Service-Architecture — Frontend Service description

**Verification:**
```
cd services/frontend && npm install  # Dependencies install
cd services/frontend && npm run dev  # Dev server starts on 5173, app loads in browser
```

**Dependencies:** None

**Acceptance criteria covered:** AC-1 (folder structure with frontend service)

---

### 4. Create services/frontend/Dockerfile

**Type:** `[tidy]`

**Description:** Create a Dockerfile for the frontend using a multi-stage build:
- **Stage 1 (build):** `node:18-alpine` base, run `npm install && npm run build` to produce optimized Vue SPA assets in `dist/`
- **Stage 2 (serve):** Use `nginx:alpine` or lightweight HTTP server (e.g., `caddy` or `http-server` via node), copy built assets from stage 1, configure to serve `dist/index.html` for all non-file routes (SPA routing)
- Expose port 3000 (or whatever the frontend service port is)
- Set `VITE_API_URL` environment variable (passed at runtime via docker-compose) to tell frontend where backend is

**Files:**
- `services/frontend/Dockerfile` (new)
- `services/frontend/.dockerignore` (new, exclude node_modules, .git, .env.local, etc.)

**References:**
- `services/frontend/vite.config.ts` — build output location (usually `dist/`)
- Nginx or http-server configuration for SPA routing
- CONSTITUTION.md#Multi-Service-Architecture — Frontend Service port and environment

**Verification:**
```
docker build services/frontend/ -t leadpulse-frontend:latest
# Build should complete without errors

docker run -p 3000:3000 leadpulse-frontend:latest
# Container starts; accessing http://localhost:3000 in browser loads app
```

**Dependencies:** Subtask 3

**Acceptance criteria covered:** AC-3 (frontend Dockerfile)

---

### 5. Create docker-compose.yml with all services

**Type:** `[behavior]`

**Description:** Create a `docker-compose.yml` at the project root that orchestrates two services:
- **backend:** Built from `services/backend/Dockerfile`, port 8080, environment variable DATABASE_URL, health check `curl http://localhost:8080/api/health`
- **frontend:** Built from `services/frontend/Dockerfile`, port 3000, depends on backend, environment variable VITE_API_URL=http://backend:8080
- **network:** All services on a shared docker-compose network for DNS-based service discovery (frontend can reach `http://backend:8080`)
- **volumes:** `data` volume for SQLite persistence (mounted at `/data` in backend container)

Test that `docker-compose up` successfully boots both services and health checks pass within a timeout.

**Files:**
- `docker-compose.yml` (new, at root)
- `.dockerignore` (new, in root; exclude .git, .agent/, node_modules, etc.)

**References:**
- Subtask 2 (backend Dockerfile health check)
- Subtask 4 (frontend Dockerfile)
- Docker Compose v3.8 reference (networks, volumes, depends_on)

**Tests:**
- id: compose-boots-all-services
  file: Manual verification (no automated test for this; observer runs docker-compose up and verifies output)
  name: `docker-compose up` starts backend, frontend, database; health checks pass
  state: pending
- id: backend-reachable-from-frontend
  file: Manual verification (docker exec frontend curl http://backend:8080/api/health)
  name: Frontend container can reach backend by hostname
  state: pending

**Active test:** compose-boots-all-services

**Verification:**
```
docker-compose up --build  # Builds and starts both services
# Wait 30 seconds, then verify in another terminal:
curl http://localhost:8080/api/health  # Backend responds
curl http://localhost:3000  # Frontend HTML loads
docker-compose logs backend  # Logs show "Starting HTTP server on http://0.0.0.0:8080"
docker-compose logs frontend  # Logs show server started
docker-compose down  # Stop all services; volume persists
docker-compose up  # Restart without rebuild; backend uses existing /data/data.db; data intact
```

**Dependencies:** Subtasks 2, 4

**Acceptance criteria covered:** AC-4 (docker-compose orchestration)

---

### 6. Set up acceptance-tests/ directory with Playwright scaffold

**Type:** `[tidy]`

**Description:** Create `acceptance-tests/` directory at project root with a Playwright TypeScript project scaffold:
- `package.json` with Playwright, TypeScript dev dependencies
- `playwright.config.ts` with baseURL set to `http://localhost:3000` and timeout 30s
- `tests/` directory (empty, ready for tests in subtask 7)
- `.gitignore` to exclude `node_modules/`, `.env`, `results/`, `playwright/.auth/`, etc.

This is scaffolding only; actual tests are written in subtask 7.

**Files:**
- `acceptance-tests/package.json` (new)
- `acceptance-tests/playwright.config.ts` (new)
- `acceptance-tests/tests/` (new directory)
- `acceptance-tests/.gitignore` (new)

**References:**
- Playwright docs: https://playwright.dev (reference only)
- CONSTITUTION.md#Testing-Strategy — Acceptance tests run against live environment

**Verification:**
```
cd acceptance-tests && npm install  # Dependencies install
cd acceptance-tests && npx playwright --version  # Playwright installed
```

**Dependencies:** None

**Acceptance criteria covered:** AC-5 (acceptance test foundation)

---

### 7. Write first acceptance test (app shell loads + backend reachable)

**Type:** `[behavior]`

**Description:** Write a Playwright test that verifies:
1. App shell loads at `http://localhost:3000`
2. Sidebar is visible (has expected class/content)
3. Top bar is visible
4. Footer is visible
5. Health indicator in footer shows green (✓) when backend is reachable
6. No console errors during page load

This test serves as the template for future acceptance tests and verifies end-to-end integration.

**Files:**
- `acceptance-tests/tests/app.spec.ts` (new)

**References:**
- `services/frontend/src/components/AppShell.vue` — expected DOM structure
- Playwright test best practices (use data-testid attributes for selection)
- CONSTITUTION.md#Testing-Strategy — Acceptance test scope

**Tests:**
- id: app-shell-renders
  file: `acceptance-tests/tests/app.spec.ts`
  name: `test('app shell loads with sidebar, top bar, footer')`
  state: pending
- id: backend-reachable
  file: `acceptance-tests/tests/app.spec.ts`
  name: `test('health indicator shows green when backend responds')`
  state: pending
- id: no-console-errors
  file: `acceptance-tests/tests/app.spec.ts`
  name: `test('no console errors on page load')`
  state: pending

**Active test:** app-shell-renders

**Verification:**
```
cd acceptance-tests && npx playwright test
# Tests run against http://localhost:3000 (assumes docker-compose up is running)
# All three tests pass; output shows test results and duration
```

**Dependencies:** Subtasks 3, 5, 6

**Acceptance criteria covered:** AC-5 (acceptance test suite verifies app shell)

---

### 8. Update Makefile with docker and test targets

**Type:** `[tidy]`

**Description:** Update `Makefile` at root with new targets for managing docker-compose and running tests:
- `make docker-build` — build both backend and frontend container images
- `make docker-up` — start docker-compose (runs `docker-compose up -d`)
- `make docker-down` — stop docker-compose (runs `docker-compose down`)
- `make docker-logs` — tail docker-compose logs
- `make test:backend` — run Go backend tests (cd services/backend && go test -race ./...)
- `make test:frontend` — run Vue frontend tests (cd services/frontend && npm test, if added later)
- `make test:acceptance` — run Playwright acceptance tests against live docker-compose (assumes `docker-compose up` is running)
- `make test:all` — run all three test categories in sequence
- `make clean` — clean up: docker-compose down, remove volumes, clear build artifacts

Keep existing targets (if any) and add these new ones.

**Files:**
- `Makefile` (modify — add new targets)

**References:**
- Current `Makefile` structure
- Make syntax for phony targets and comments

**Verification:**
```
make docker-build  # Builds images
make docker-up  # Starts services
make test:acceptance  # Runs acceptance tests
make docker-down  # Stops services
make clean  # Cleans up
```

**Dependencies:** Subtasks 2, 4, 5, 7

**Acceptance criteria covered:** AC-6 (make targets for common tasks)

---

### 9. Update .gitignore for docker volumes and frontend artifacts

**Type:** `[tidy]`

**Description:** Update `.gitignore` at root to exclude:
- `data/` — SQLite data volume (persistent storage, should not be committed)
- `services/frontend/node_modules/` — npm dependencies
- `services/frontend/dist/` — built Vue SPA
- `services/backend/node_modules/` (if any)
- `.env`, `.env.local` — environment files
- `acceptance-tests/node_modules/` — Playwright dependencies
- `acceptance-tests/test-results/` — test output
- Docker-related: `.dockerignore` already in place

**Files:**
- `.gitignore` (modify)

**References:**
- Current `.gitignore` contents
- Docker and npm best practices

**Verification:**
```
git status  # Should not show node_modules, data/, dist/ as untracked
git check-ignore data/  # Should return data/ (is ignored)
```

**Dependencies:** None

**Acceptance criteria covered:** AC-1 (folder structure properly managed)

---

### 10. Verify docs/architecture.md reflects multi-service design

**Type:** `[tidy]`

**Description:** Confirm `docs/architecture.md` accurately describes the new multi-service architecture:
- Explicitly document the three containers (backend, frontend, database)
- Update C4 Level 2 diagram to show services in separate containers
- Describe communication paths: frontend → backend via HTTP/JSON API over docker-compose network
- Describe data persistence: SQLite volume at `/data`
- Confirm dependency direction diagram shows frontend (HTTP client) → backend (API) with hexagonal core/storage/engine inside backend

No changes required if the diagram already reflects this; this is a verification step.

**Files:**
- `docs/architecture.md` (touch — reference only; review for accuracy)

**References:**
- `docs/architecture.md#C4-Level-2` — container diagram
- CONSTITUTION.md#Multi-Service-Architecture — architectural guardrails
- Subtask 5 (docker-compose.yml structure)

**Verification:**
```
# Visual inspection: Open docs/architecture.md and confirm:
# 1. Three containers shown (backend, frontend, database)
# 2. Frontend → Backend communication labeled as HTTP/JSON
# 3. Backend contains hexagonal layers (core, storage, engine, server)
# 4. Database described as SQLite file volume
```

**Dependencies:** Subtasks 2, 4, 5

**Acceptance criteria covered:** AC-1, AC-4 (architecture documented)

---

### 11. Create docs/testing.md

**Type:** `[research]`

**Description:** Create `docs/testing.md` documenting the multi-service testing strategy:
- Backend tests: `go test -race ./...` in `services/backend/`; verifies API contracts, business logic, storage
- Frontend tests: `npm test` in `services/frontend/` (Vitest, to be set up in future increment); verifies components, routing, API client
- Acceptance tests: `npm run test:acceptance` in `acceptance-tests/`; runs Playwright against live `http://localhost:3000`; verifies end-to-end flows
- Test running procedure: How to run each category locally, in Docker, and in CI
- Evidence required before merge: All three categories pass
- Confidence model: What each test layer provides confidence in

Reference the existing `CONSTITUTION.md#Testing-Strategy` and expand with procedures and commands.

**Files:**
- `docs/testing.md` (new)

**References:**
- `CONSTITUTION.md#Testing-Strategy` — guardrails
- `services/backend/Makefile` (if it exists) — test commands
- `services/frontend/package.json` — future npm test command
- `acceptance-tests/playwright.config.ts` — test config

**Verification:**
```
# Review docs/testing.md for clarity and completeness:
# 1. Explains backend test approach and commands
# 2. Explains frontend test approach (future setup)
# 3. Explains acceptance test approach and commands
# 4. Shows how to run all tests locally
# 5. Shows how to run tests in CI (inside docker build)
```

**Dependencies:** Subtasks 2, 4, 7

**Acceptance criteria covered:** AC-6 (testing procedures documented)

---

### 12. Create docs/deployment.md

**Type:** `[research]`

**Description:** Create `docs/deployment.md` documenting the multi-service deployment strategy:
- Local development workflow: `docker-compose up` starts backend and frontend; data persists in `data/` volume
- `docker-compose down` stops services; volume preserves all state
- `docker-compose up` again restarts with existing data intact
- CI workflow: Build backend and frontend containers in separate CI jobs; run acceptance tests against composed containers
- Release procedure (future; out of scope for v1): How to tag, build, and push container images
- Configuration: How environment variables are passed to services (VITE_API_URL, DATABASE_URL, etc.)
- **No rollback procedure:** Always deploy forward; no mechanism to revert to older database schema or code
- **No database migrations:** v1 deploys with single schema; schema evolution handled via code deployments only, never via migration scripts
- **Data persistence:** SQLite volume mounts; data survives `docker-compose down/up` cycles and restarts

Reference the existing `CONSTITUTION.md#Release-And-Deployment` and expand with procedures.

**Files:**
- `docs/deployment.md` (new)

**References:**
- `CONSTITUTION.md#Release-And-Deployment` — guardrails
- `docker-compose.yml` — environment variables and volumes
- `services/backend/main.go` — database initialization

**Verification:**
```
# Review docs/deployment.md for clarity and completeness:
# 1. Explains docker-compose local dev workflow (data persists)
# 2. Explains CI build and test workflow
# 3. Explains fresh-start principle (no rollback, no migration)
# 4. Documents all environment variables used
# 5. Explains that docker-compose down preserves volume, docker-compose up restarts with existing data
```

**Dependencies:** Subtasks 5, 8

**Acceptance criteria covered:** AC-6 (deployment procedures documented)

---

## Context Map

- `CONSTITUTION.md#Multi-Service-Architecture` — Guardrails for backend, frontend, orchestration, dependency direction
- `CONSTITUTION.md#Testing-Strategy` — Testing categories: backend (go test), frontend (npm test), acceptance (Playwright)
- `CONSTITUTION.md#Release-And-Deployment` — Release model (docker-compose for local dev, CI, future container registry)
- `docs/architecture.md` — C4 Level 2 diagram (to be updated to show multi-service containers)
- Existing backend code structure (`core/`, `engine/`, `storage/`, `server/`, `main.go`) — to be moved unchanged to `services/backend/`

## Acceptance Scenarios

**Scenario 1: Local developer boots full system**

- **User action:** Run `make docker-up` from project root
- **Precondition:** Docker and Docker Compose installed; code is up-to-date
- **Expected outcome:** All three services start; health checks pass; accessing `http://localhost:3000` shows app shell with green health indicator
- **Criterion:** AC-4 (docker-compose works), AC-5 (acceptance test verifies it)
- **Evidence:** `docker-compose logs` shows all services healthy; browser shows app
- **Gate:** required (main success flow)
- **State:** planned

**Scenario 2: Acceptance tests pass against running containers**

- **User action:** Run `make test:acceptance` from project root (assumes `docker-compose up` is running)
- **Precondition:** Docker-compose services are healthy
- **Expected outcome:** All acceptance tests pass; app shell loads, backend reachable, no console errors
- **Criterion:** AC-5 (acceptance test suite)
- **Evidence:** Playwright test output shows 3/3 passed
- **Gate:** required (acceptance gate before merge)
- **State:** planned

**Scenario 3: Developer rebuilds after code changes**

- **User action:** Edit code in `services/backend/` or `services/frontend/`; run `make docker-build && make docker-down && make docker-up`
- **Precondition:** Changes are saved
- **Expected outcome:** Containers rebuild with new code; services restart; app reloads with changes
- **Criterion:** AC-1 (services rebuild), AC-4 (docker-compose restarts)
- **Evidence:** Docker build logs show successful builds; app in browser reflects changes
- **Gate:** advisory (developer workflow)
- **State:** planned

**Scenario 4: Backend tests pass inside docker build**

- **User action:** Run `docker build services/backend/ -t leadpulse-backend:latest`
- **Precondition:** Go code is valid
- **Expected outcome:** Docker multi-stage build runs `go test` first; if tests fail, build stops; if tests pass, binary is built
- **Criterion:** AC-2 (backend Dockerfile)
- **Evidence:** Docker build log shows "PASS" for all tests; final image created
- **Gate:** required (no image without passing tests)
- **State:** planned

**Scenario 5: Acceptance test detects backend failure**

- **User action:** Stop backend service manually (`docker stop leadpulse-backend`); run `make test:acceptance`
- **Precondition:** Frontend is still running; acceptance tests are waiting
- **Expected outcome:** Acceptance test fails because health indicator is red or backend is unreachable
- **Criterion:** AC-5 (acceptance test catches integration issues)
- **Evidence:** Playwright test fails with clear error (e.g., "Expected green health indicator, got red")
- **Gate:** advisory (error detection)
- **State:** planned

**Scenario 6: Makefile targets work end-to-end**

- **User action:** Run `make docker-build`, then `make docker-up`, then `make test:acceptance`, then `make docker-down`
- **Precondition:** Makefile has all targets; code is valid
- **Expected outcome:** All make targets execute without error; full workflow succeeds
- **Criterion:** AC-6 (make targets)
- **Evidence:** All four commands complete; acceptance tests pass
- **Gate:** required (make workflow is developer daily tool)
- **State:** planned

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-----------|
| Docker build layer caching causes stale tests | Medium | Medium | Add `--no-cache` flag to docker build in Makefile; document when to use it. Monitor build times. |
| Frontend build takes too long (npm install in container) | Medium | Low | Consider caching node_modules layer or using a builder image; document expected build time (target < 2 min). |
| Service discovery fails (frontend cannot resolve `backend` hostname) | Low | High | Test DNS resolution inside frontend container during docker-compose up; add verbose logging to docker-compose.yml (logging driver). If fails, document troubleshooting steps in docs/deployment.md. |
| Playwright tests flake (timing issues waiting for services to be healthy) | Medium | Low | Add explicit wait loops in Playwright tests (wait for backend health endpoint to succeed); set generous timeouts (30s). Document flakiness handling in docs/testing.md. |
| Port conflicts (8080 or 3000 already in use) | Low | Medium | Make docker-compose port mapping configurable via .env file; add troubleshooting to docs/deployment.md (show how to find and kill process using port). |
| SQLite volume persists across docker-compose down (confusing for clean slate) | Low | Low | Document that `docker-compose down -v` removes volumes; add `make clean` target that includes `-v` flag; add warning comment to docker-compose.yml. |
| Makefile targets conflict with existing targets | Low | Medium | Review current Makefile before adding new targets; use descriptive names (`docker-build` vs. `build`) to avoid collisions. Prefix docker targets consistently. |

## Planning Decisions

| Decision | Chosen | Rejected | Reason |
|----------|--------|----------|--------|
| Multi-service architecture | Backend + frontend in separate containers | Single-binary approach | Allows independent scaling, development, and deployment. Aligns with modern containerized workflows. Frontend and backend can evolve independently. |
| Data persistence model | SQLite file volumes; persistent across restarts | In-memory databases (ephemeral) | Data survives `docker-compose down/up` cycles. Developers can restart services without losing work. Aligns with stateless containers + persistent volumes pattern. |
| Database per container | SQLite at `/data/data.db` (volume mount) | Shared PostgreSQL container | Simplicity for v1; no separate database container. SQLite file-based works in docker-compose volume. Migrate to PostgreSQL container in future if needed. |
| Frontend port | 3000 | 8080 (same as backend) | Avoids port collision; 3000 is conventional for npm dev servers. Easier to run both locally simultaneously. |
| Frontend web server | Nginx or http-server (lightweight) | Node.js express or similar | Keeps container small and fast; no runtime Node.js needed in final image. Static asset serving is built-in. |
| Service discovery | Docker-compose internal DNS (hostname: `backend`, `frontend`) | Hardcoded IPs or environment variables | Simplest approach; hostname-based discovery is standard docker-compose pattern. Frontend environment variable `VITE_API_URL=http://backend:8080` is injected at container start. |
| Acceptance test location | `acceptance-tests/` at root (own npm project) | Inside `services/frontend/` as `tests/e2e/` | Separate project keeps concerns clear; acceptance tests are cross-service validation, not frontend-only. Independent from frontend build cycle. Easier to run separately. |
| Makefile targets | Prefix with `docker-` and `test:` | Mix prefix styles | Consistency and clarity; `docker-` targets are container-related, `test:` targets are test-related. Easy to discover with `make` (tab-complete). |
| Test gate before merge | All three categories (backend, frontend, acceptance) must pass | Acceptance tests only | Comprehensive validation across layers. Backend + frontend unit tests catch logic errors; acceptance tests catch integration issues. Highest confidence before ship. |
| Backend Dockerfile test stage | `go test -race ./...` inside docker build | Tests run separately in CI | Builds confidence that tests pass *before* binary is built. Fails fast if tests break. Enforces test-first discipline. No binary without passing tests. |
| Rollback strategy | No rollback procedure; always deploy forward | Keep old container versions available for rollback | Simpler operational model; no need to maintain old database schemas or code versions. If deployment fails, fix and redeploy. Aligns with immutable infrastructure principles. |
| Database migrations | No migration scripts; v1 deploys with single schema | Run migrations as part of deployment | Simpler; eliminates migration bugs and rollback complexity. Schema evolution handled via code changes only. If schema must change, deploy new code + new database; never run "migration" scripts. |

## Verification Checkpoints

| Subtask | Pass Condition | Fail Condition |
|---------|---|---|
| 1. Move backend to services/backend/ | `cd services/backend && go build ./...` succeeds | Import errors or missing files |
| 2. Backend Dockerfile | `docker build services/backend/ -t leadpulse-backend:latest` completes; `go test` passes inside container | Docker build fails; tests fail inside container |
| 3. Frontend scaffold | `cd services/frontend && npm install` completes; `npm run dev` starts dev server | npm install fails; dev server won't start |
| 4. Frontend Dockerfile | `docker build services/frontend/ -t leadpulse-frontend:latest` completes; container serves assets on port 3000 | Docker build fails; container won't start or port unreachable |
| 5. docker-compose.yml | `docker-compose up --build` starts all services; health checks pass within 30s | Services fail to start; health checks timeout |
| 6. Acceptance tests scaffold | `cd acceptance-tests && npm install && npx playwright --version` works | npm install fails; Playwright not found |
| 7. First acceptance test | `docker-compose up` running; `cd acceptance-tests && npm run test` passes all 3 tests | Tests fail; app shell not visible or backend unreachable |
| 8. Makefile targets | `make docker-build`, `make docker-up`, `make test:acceptance`, `make docker-down` all work | Any make target fails with unclear error |
| 9. .gitignore | `git status` does not show `data/`, `node_modules/`, or `dist/` | Untracked files appear in git status |
| 10. docs/architecture.md | Diagram shows three containers; frontend → backend communication labeled; hexagonal core inside backend | Diagram outdated; missing multi-service view |
| 11. docs/testing.md | Document explains backend, frontend, acceptance test strategies with commands | Missing procedures or commands unclear |
| 12. docs/deployment.md | Document explains docker-compose workflow, CI procedure, fresh-start philosophy (no rollback, no migration) | Missing deployment procedures or unclear |

