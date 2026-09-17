# Increment: Docker Services Architecture & Acceptance Testing Foundation

## Use Case

When I develop features locally or need to validate the system in CI, I want to run the entire application (frontend, backend, database) as containerized services using docker-compose, and verify end-to-end behavior with acceptance tests that exercise both the UI and API together, so that I can catch integration bugs early and ship with confidence.

## Goal

Restructure the project into a multi-service architecture with separate frontend and backend services in Docker containers, establish an acceptance test suite that runs against the full containerized system, and create a docker-compose orchestration file that boots the entire system locally.

## Branch

`increment/docker-services-architecture`

## Acceptance Criteria

1. **Services Folder Structure** — Project reorganized with `services/frontend/` and `services/backend/` subdirectories; each contains its own build artifacts and tooling (package.json for frontend, go.mod for backend)

2. **Backend Dockerfile** — Builds Go binary from `services/backend/`, packages it as a container, exposes HTTP port 8080, includes a health check command that queries `GET /api/health`

3. **Frontend Dockerfile** — Builds Vue 3 SPA via npm/Vite from `services/frontend/`, produces static assets, serves them via a lightweight web server (nginx or similar), exposes port 3000 or appropriate port

4. **Docker-Compose Orchestration** — `docker-compose.yml` in project root defines three services: frontend, backend, database (SQLite file volume or PostgreSQL container); services discover each other by hostname; frontend can reach backend at `http://backend:8080`

5. **Acceptance Tests** — Playwright test suite in `tests/acceptance/` that runs against live docker-compose environment; verifies system boots, app shell loads, backend is reachable, health check works

6. **Local Development** — Running `docker-compose up` boots the entire system; app is accessible at `http://localhost:3000` (frontend) or `http://localhost:8080` (direct backend); all tests can run against containerized services

## Acceptance-Test Intent

User journey: Run `docker-compose up` → wait for services to be healthy → open browser to `http://localhost:3000` → app loads with shell visible → health indicator in footer shows green ✓ → run `npm run test:acceptance` → tests pass, verifying end-to-end integration.

## Out Of Scope

- Kubernetes deployment or orchestration (docker-compose for local dev and CI only)
- Production-grade multi-region or load-balanced deployment
- Container registry, image push, or image versioning strategy (build locally for now)
- Database migrations or schema management beyond SQLite file persistence
- Frontend build optimization (minification, code splitting) — Vite defaults are acceptable
- Environment variable secrets management (hard-code for local dev; CI uses GitHub Secrets)
- Hot reload or live development server inside containers (use docker-compose override file for dev, standard Dockerfile for CI)

## Constitution Constraints

- **Single-binary principle affected:** Backend is now a Docker container, not a single redistributable binary. This increments marks the shift from desktop app → development of a multi-container system. Future "release build" may return to single-binary if needed.
- **No external network:** All services run locally (localhost); no cloud or external API calls. Database is a local volume or in-memory SQLite.
- **Dependency direction:** Backend (`services/backend/`) remains hexagonal; frontend (`services/frontend/`) is a separate project; both talk only through JSON API.
- **Testing strategy:** Acceptance tests run against the full containerized system (new); unit and integration tests for individual services run inside their respective containers during build.

## Roadmap Entry

**Docker Services Architecture & Acceptance Testing Foundation**

Restructures the project into a multi-service architecture with separate frontend and backend containers, establishes docker-compose orchestration for local development, and creates an acceptance test suite that validates end-to-end behavior across the containerized system.

Job story: When I develop features locally or validate them in CI, I want to run the entire system (frontend, backend, database) as containerized services via docker-compose, and verify end-to-end behavior with acceptance tests that exercise both the UI and API, so that I can catch integration bugs early and ship with confidence.

Evidence: `docker-compose up` boots all services successfully; acceptance tests in `tests/acceptance/` pass against the running containers; health checks are verified in both frontend and backend containers.
