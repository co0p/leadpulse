# Deployment & Operations

Multi-service Docker deployment strategy for local development, CI/CD, and production.

---

## Architecture Overview

```
Local Development (docker-compose):
  → Backend container (port 8080)
  → Frontend container (port 3000, proxies /api/* to backend)
  → Shared docker-compose network (DNS-based service discovery)
  → Persistent SQLite volume (./data/data.db)

CI/CD Pipeline:
  → Build backend image (multi-stage: test → build)
  → Build frontend image (multi-stage: npm build → nginx)
  → Run acceptance tests against built images
  → (Future: push images to registry)

Production (future):
  → Kubernetes or managed container platform
  → Separate database service (PostgreSQL)
  → Environment-specific configuration
  → Health checks and auto-restart policies
```

---

## Local Development: docker-compose

### Start Services

```bash
# Build and start both services
docker-compose up --build

# Or, build separately if images already exist
docker-compose build
docker-compose up

# Run in background (detached mode)
docker-compose up -d

# View logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Access Services

- **Frontend:** http://localhost:3000 (Vue SPA)
- **Backend API:** http://localhost:8080 (JSON endpoints)
- **Health Check:** http://localhost:8080/api/health

### Database

- **Location:** `./data/data.db` (on host machine)
- **Persistence:** Survives `docker-compose down`
- **First run:** Volume created automatically; database initialized fresh
- **Subsequent runs:** Existing database file used (data persists)

### Data Persistence

**When does data persist?**
- `docker-compose up` → `docker-compose down` → `docker-compose up` = **Data persists**
- `docker-compose down -v` (with `-v` flag) → Removes volumes = **Data deleted**
- `make docker-down` (without `-v`) → Data persists

**Typical workflow:**
```bash
# Terminal 1: Keep services running
make docker-up

# Terminal 2: Develop and test
# ... make code changes ...
make test-backend
make test-acceptance
# ... data stays intact across test runs ...

# When done
make docker-down  # Services stop, data persists
```

**To reset data:**
```bash
docker-compose down -v  # Remove volumes
docker-compose up       # Rebuild; database initialized fresh
```

### Restart Services

```bash
# Stop services (data persists)
docker-compose down

# Restart with existing data
docker-compose up

# This is equivalent to redeploying: containers are recreated from images,
# but the data volume is reused.
```

### Common Tasks

```bash
# Rebuild after code changes (frontend or backend)
docker-compose build
docker-compose up

# View resource usage
docker stats

# Execute command inside backend container
docker exec leadpulse-backend go test -race ./...

# View container environment variables
docker exec leadpulse-backend env

# Inspect volume contents
docker volume inspect leadpulse_data

# Clean up everything (including volumes)
docker-compose down -v
docker system prune -a --volumes
```

---

## CI/CD Pipeline

### Build Stage

Backend image build:
```dockerfile
# services/backend/Dockerfile
Stage 1 (test):   RUN go test -race ./...     # Must pass
Stage 2 (build):  RUN go build -o binary
Stage 3 (runtime): FROM alpine; COPY binary
```

If tests fail in Stage 1, the build stops immediately (fail-fast). No faulty binary is produced.

Frontend image build:
```dockerfile
# services/frontend/Dockerfile
Stage 1 (build):   npm install; npm run build → dist/
Stage 2 (runtime): COPY dist/ to Nginx; serve + proxy /api/*
```

### Test Stage

```bash
# After successful image builds
docker-compose up -d          # Start services
sleep 10                      # Wait for health checks
npm test acceptance-tests/    # Run Playwright tests

# Capture results
# - test-results/results.json (machine-readable)
# - test-results/index.html (human-readable report)
# - screenshots/ (on failure)
# - videos/ (on failure)
```

### Push Stage (Future)

```bash
# Tag images
docker tag leadpulse-backend:latest myregistry.azurecr.io/leadpulse-backend:v1.0
docker tag leadpulse-frontend:latest myregistry.azurecr.io/leadpulse-frontend:v1.0

# Push to registry
docker push myregistry.azurecr.io/leadpulse-backend:v1.0
docker push myregistry.azurecr.io/leadpulse-frontend:v1.0
```

---

## Deployment Model: Always Deploy Forward

**Philosophy:** No rollback procedures, no schema migrations, no backward compatibility layers.

**How it works:**

1. **Code change** → Create new images with new code
2. **Deploy** → `docker-compose up --build` (or equivalent in production)
3. **Services restart** with fresh state (or migrate from persistent volume if data needed)
4. **If failure** → Rollback by deploying previous code version (fast rebuild from git history)

**No "rollback" scripts** because:
- Images are immutable; rebuilding from prior git commit is rollback
- No database migration scripts to "undo"
- Each deployment is self-contained

**For local dev:**
```bash
# Mistake discovered after deployment
git revert <bad-commit>     # Or check out prior version
make docker-down -v         # Remove old volumes (fresh start)
make docker-up              # Build and deploy from prior code
```

**For production (future):**
```bash
# Via CI/CD:
git revert <bad-commit>
push to main branch
CI pipeline rebuilds and deploys new images automatically
```

---

## Database Strategy

### Version 1 (Current): SQLite File Volume

- **Location:** `/data/data.db` (container) ← `./data/` (host)
- **Persistence:** Volume mounts across restarts
- **Initial state:** Database schema created on first boot (idempotent)
- **No migrations:** Schema defined in Go code (`store/schema.go`); applied once at startup
- **No rollback:** If schema change needed, deploy new code + new database (fresh start)

### Version 2 (Future): PostgreSQL

- **Separate container:** `postgres:latest`
- **Persistence:** PostgreSQL volume for data files
- **Network:** Reachable as `postgres:5432` from backend via docker-compose network
- **Backup/restore:** Use `pg_dump` / `pg_restore`

```yaml
# Future docker-compose.yml addition
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: leadpulse
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - leadpulse-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

---

## Environment Configuration

### Local Development

Environment variables set in docker-compose.yml:

```yaml
services:
  backend:
    environment:
      DATABASE_URL: "sqlite:///data/data.db"
  frontend:
    environment:
      VITE_API_URL: "http://backend:8080"
```

Or via `.env` file (docker-compose reads `.env` automatically):

```bash
# .env
DATABASE_URL=sqlite:///data/data.db
VITE_API_URL=http://backend:8080
```

### CI/CD

Environment variables injected by CI system (e.g., GitHub Actions, Azure Pipelines):

```bash
# Example: GitHub Actions
env:
  DATABASE_URL: sqlite:///data/data.db
  VITE_API_URL: http://backend:8080
  LOG_LEVEL: debug
```

### Production (Future)

Environment variables from:
- Kubernetes Secrets
- Cloud provider configuration (Azure Key Vault, AWS Secrets Manager)
- `.env` file (if using docker-compose on server)

Example for production:
```bash
# production.env
DATABASE_URL=postgresql://user:password@postgres.example.com:5432/leadpulse
VITE_API_URL=https://api.leadpulse.example.com
LOG_LEVEL=info
TLS_CERT_PATH=/etc/certs/server.crt
TLS_KEY_PATH=/etc/certs/server.key
```

---

## Health Checks

### Backend Health Endpoint

```bash
GET http://localhost:8080/api/health

Response (200 OK):
{
  "status": "ok"
}
```

Used by:
- `docker-compose healthcheck` (startup probe)
- Frontend AppShell health indicator (every 5 seconds)
- Load balancer (future, in production)

### Frontend Health Endpoint

```bash
GET http://localhost:3000/

Response (200 OK):
<html>...</html>
```

Used by:
- `docker-compose healthcheck` (startup probe)
- Load balancer (future, in production)

### Startup Probe Behavior

```yaml
# docker-compose.yml
services:
  backend:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
  frontend:
    depends_on:
      backend:
        condition: service_healthy
```

Sequence:
1. Service starts
2. `start_period: 10s` — ignore failures for first 10 seconds
3. `retries: 5` — allow up to 5 failed checks before marking unhealthy
4. `interval: 10s` — check every 10 seconds
5. `timeout: 5s` — if check takes > 5s, consider it failed

---

## Monitoring & Logging

### Local Development

```bash
# View live logs from both services
docker-compose logs -f

# Timestamp format
docker-compose logs -f --timestamps

# Follow only backend logs
docker-compose logs -f backend

# Last 100 lines of frontend logs
docker-compose logs --tail=100 frontend
```

### CI/CD

Logs stored in CI system (GitHub Actions, Azure Pipelines):
- Build logs
- Test results
- Docker build output
- Test coverage reports

### Production (Future)

Centralized logging:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- CloudWatch (AWS) / Application Insights (Azure)
- Datadog, Splunk, etc.

Backend should log to stdout (docker-compose captures automatically):
```go
// Example Go logging to stdout
log.Printf("[INFO] Starting server on %s", addr)
log.Printf("[ERROR] Database error: %v", err)
```

---

## Scaling (Future)

### Horizontal Scaling

Multiple backend instances behind a load balancer:
```yaml
services:
  backend-1:
    build: ./services/backend
  backend-2:
    build: ./services/backend
  backend-3:
    build: ./services/backend
  
  load-balancer:
    image: haproxy:latest
    ports:
      - "8080:8080"
    # ... configuration to round-robin to backend-1/2/3
```

### Kubernetes (Future)

Deployment manifests:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: leadpulse-backend
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: backend
        image: myregistry.azurecr.io/leadpulse-backend:latest
        livenessProbe:
          httpGet:
            path: /api/health
            port: 8080
        readinessProbe:
          httpGet:
            path: /api/health
            port: 8080
```

---

## Troubleshooting

### Services won't start

```bash
# Check docker daemon is running
docker ps

# View detailed error logs
docker-compose logs -f

# Rebuild from scratch
docker-compose down -v
docker-compose build --no-cache
docker-compose up
```

### Port already in use

```bash
# Find process using port 8080
lsof -i :8080

# Kill it (careful!)
kill -9 <PID>

# Or use different ports in docker-compose
docker-compose up -e BACKEND_PORT=8081 -e FRONTEND_PORT=3001
```

### Database corruption

```bash
# Remove volume and restart (fresh database)
docker-compose down -v
docker-compose up --build
```

### Services stuck in unhealthy state

```bash
# Restart services
docker-compose restart

# Or full rebuild
docker-compose down
docker-compose up --build
```

---

## References

- `docker-compose.yml` — Service orchestration configuration
- `CONSTITUTION.md#Release-And-Deployment` — High-level deployment principles
- `docs/testing.md` — How tests run in CI/CD pipeline
- `docs/architecture.md` — Service dependencies and communication
- `services/backend/Dockerfile` — Backend build strategy
- `services/frontend/Dockerfile` — Frontend build strategy
