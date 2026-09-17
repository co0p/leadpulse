# Testing Strategy

Multi-service architecture testing across three layers: backend unit tests, frontend tests, and end-to-end acceptance tests.

---

## Testing Pyramid

```
        /\
       /  \  Acceptance Tests (Playwright)
      /----\  Cross-service validation; slow; few scenarios
     /      \
    /        \
   /----------\
  /  Frontend  \ Vue component tests (npm test)
 /   Tests     \ Faster; single-service logic
/              \
/________________\
Backend Unit Tests (go test)
Fastest; core business logic; hexagonal architecture
```

---

## Layer 1: Backend Unit Tests

**Technology:** Go `testing` package with `go test -race`

**Location:** `services/backend/**/*_test.go`

**Scope:**
- Core domain logic (`core/members/`, `engine/scoring/`, etc.)
- Repository adapters (SQLite and in-memory implementations)
- Use cases (business logic orchestration)
- Service coordinators

**Run locally:**
```bash
make test-backend
# or
cd services/backend && go test -race ./...
```

**Run in CI (Docker):**
- Test stage of backend Dockerfile runs `go test -race ./...` before build
- If tests fail, docker build stops immediately (fail-fast)
- No binary produced if tests fail

**Key characteristics:**
- In-memory test repository for speed (no database I/O per test)
- Race detector enabled (`-race` flag)
- Comprehensive coverage of business logic branches
- Fast feedback loop (< 10 seconds typical)

**Test naming convention:**
- Test function: `TestUseCaseName_ScenarioDescription`
- Example: `TestAddMember_Success`, `TestAddMember_InvalidInput`, `TestGetMembers_EmptyList`

---

## Layer 2: Frontend Tests

**Technology:** Vue Testing Library + Vitest (future setup)

**Location:** `services/frontend/src/**/*.spec.ts` (to be implemented in next increment)

**Scope:**
- Vue component rendering (AppShell, etc.)
- User interactions (clicks, form inputs)
- State management (if using Pinia)
- API mocking (mock backend responses)

**Run locally:**
```bash
make test-frontend
# or
cd services/frontend && npm test
```

**Run in CI:**
- Integrated into frontend Docker build (separate stage if desired)
- Can run in parallel with backend tests

**Key characteristics:**
- Component-level unit tests
- Mock API calls (no real backend dependency)
- Fast feedback (< 5 seconds typical)
- Tests isolated from docker-compose services

**Test naming convention:**
- Test name: `'renders AppShell with sidebar'`, `'health indicator shows green'`
- Pattern: human-readable behavior description

**Note:** Frontend tests are scaffolded but not yet written. First acceptance tests (Playwright) will cover frontend behavior end-to-end.

---

## Layer 3: Acceptance Tests (E2E)

**Technology:** Playwright (Chromium + Firefox)

**Location:** `acceptance-tests/tests/**/*.spec.ts`

**Scope:**
- Cross-service integration; frontend + backend working together
- User workflows (load page, check health, navigate, etc.)
- End-to-end scenarios (realistic user journeys)

**Prerequisites:**
- Both services must be running via `docker-compose up`
- Backend at `http://localhost:8080` (reachable from Playwright)
- Frontend at `http://localhost:3000` (Playwright browser target)

**Run locally:**
```bash
# Terminal 1: Start services
make docker-up

# Terminal 2: Run acceptance tests
make test-acceptance

# To use interactive UI:
cd acceptance-tests && npm run test:ui
```

**Run in CI:**
- Services already running in docker-compose
- Acceptance tests run after both services are healthy
- HTML report generated (test-results/)

**Configuration:** `acceptance-tests/playwright.config.ts`
- Base URL: `http://localhost:3000` (configurable via `BASE_URL` env var)
- Browsers: Chromium, Firefox
- Retries: 2 in CI, 0 locally
- Screenshots/videos: captured on failure

**First acceptance test suite:** `acceptance-tests/tests/app.spec.ts`

Three test cases:
1. **App Shell Renders**
   - Navigate to home
   - Verify sidebar, topbar, footer are visible
   - Verify main heading is present
   - Data attributes: `[data-testid="sidebar"]`, `[data-testid="topbar"]`, `[data-testid="footer"]`

2. **Health Indicator Shows Green**
   - Navigate to home
   - Wait for health check to complete (2 seconds)
   - Verify health indicator visible and shows success tag (green)
   - Data attribute: `[data-testid="health-indicator"]`

3. **No Console Errors**
   - Navigate to home
   - Collect console errors during page load
   - Assert no errors in browser console

**Test naming convention:**
- Test name: `'app shell loads with sidebar, top bar, footer'`
- Pattern: imperative behavior description

**Key characteristics:**
- Slow (30–60 seconds typical for small suite)
- High confidence; catches integration issues
- Validates real docker-compose setup
- Uses explicit waits (no implicit waits)
- Retries on timing failures

---

## Testing Gate (Before Merge)

All three test categories must pass:

```bash
# Backend tests (must pass)
make test-backend

# Frontend tests (to be added in future increment; currently no tests)
# make test-frontend

# Acceptance tests (require docker-compose up running)
make docker-up
make test-acceptance
make docker-down
```

Merge is blocked if any category fails.

---

## CI/CD Pipeline (Future)

1. **Backend build:** Dockerfile test stage runs `go test -race ./...`
2. **Frontend build:** npm test (when tests are added)
3. **Service images:** Docker build produces leadpulse-backend and leadpulse-frontend images
4. **Acceptance tests:** Run Playwright tests against running services
5. **Report:** Generate and archive test results (HTML, JSON, screenshots)

---

## Local Development Workflow

```bash
# 1. Make code changes (backend, frontend, or both)

# 2. Run targeted tests
make test-backend          # Quick feedback on backend changes
# make test-frontend        # (when available)

# 3. Run full suite before commit
make docker-down           # Clean up previous run
make docker-build          # Rebuild images
make docker-up             # Start services
sleep 10                   # Wait for services to be healthy
make test-acceptance       # Run Playwright tests
make docker-down           # Clean up

# 4. If all pass, commit and push
git add ...
git commit -m "..."
git push
```

---

## Debugging Tests

### Backend

```bash
cd services/backend

# Run single test
go test -race -run TestAddMember_Success ./core/members

# Verbose output
go test -race -v ./...

# CPU profile
go test -race -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Frontend

```bash
cd services/frontend

# Watch mode (re-run on file change)
npm test -- --watch

# Debug in browser
npm run test:debug
```

### Acceptance Tests

```bash
cd acceptance-tests

# Run single test file
npx playwright test tests/app.spec.ts

# Run single test
npx playwright test -g "app shell loads"

# Debug mode (pause and inspect)
npx playwright test --debug

# UI mode (visual browser)
npm run test:ui

# Generate trace for debugging
npx playwright test --trace on

# View generated trace
npx playwright show-trace test-results/trace.zip
```

---

## Known Limitations

- **Frontend tests:** Not yet implemented; planned for next increment after docker services are stable
- **Acceptance tests:** Manual verification of docker-compose until Playwright can access running services
- **Database state:** Each `docker-compose up` starts with the same initial state (no test data migration)
- **Parallel execution:** Acceptance tests run sequentially (1 worker in CI) to avoid port conflicts

---

## References

- `CONSTITUTION.md#Testing-Strategy` — High-level testing guardrails
- `services/backend/` — Backend test structure and examples
- `services/frontend/` — Frontend test structure (to be added)
- `acceptance-tests/` — Playwright configuration and test examples
- `docs/deployment.md` — How tests fit into CI/CD pipeline
