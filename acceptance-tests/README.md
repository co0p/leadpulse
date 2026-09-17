# Acceptance Tests

End-to-end acceptance tests for LeadPulse using Playwright.

## Setup

```bash
npm install
```

## Running Tests

### Run all tests
```bash
npm test
```

### Run with UI (interactive)
```bash
npm run test:ui
```

### Run in debug mode
```bash
npm run test:debug
```

### Run against a custom base URL
```bash
BASE_URL=http://example.com npm test
```

## Environment

Tests expect the following services to be running via `docker-compose up`:
- Backend API at `http://localhost:8080`
- Frontend at `http://localhost:3000`

## Test Structure

- `tests/` - Test files (*.spec.ts)
- `test-results/` - Test results and artifacts (screenshots, videos, HTML report)
- `playwright.config.ts` - Playwright configuration

## CI/CD

In CI environments, tests run headless with:
- 1 worker (sequential)
- 2 retries on failure
- HTML report generation
