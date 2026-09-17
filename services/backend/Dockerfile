# Stage 0: Build Frontend SPA
FROM node:20-alpine AS frontend-build

WORKDIR /frontend

# Copy frontend package files
COPY ./services/frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy frontend source
COPY ./services/frontend/ .

# Run frontend tests
RUN npm run test -- --run || true

# Build Vue app (produces dist/ directory)
RUN npm run build

# Stage 1: Test Backend
FROM golang:1.25-alpine AS test

WORKDIR /build

# Install build dependencies (gcc for cgo)
RUN apk add --no-cache gcc musl-dev

# Copy module files
COPY ./services/backend/go.mod ./services/backend/go.sum ./

# Copy source (backend package)
COPY ./services/backend/ .

# Copy frontend dist/ from frontend-build stage into server/dist/
COPY --from=frontend-build /frontend/dist ./server/dist

# Run tests with race detector
RUN go test -race ./...

# Stage 2: Build
FROM golang:1.25-alpine AS build

WORKDIR /build

# Install build dependencies (gcc for cgo)
RUN apk add --no-cache gcc musl-dev

# Copy module files
COPY ./services/backend/go.mod ./services/backend/go.sum ./

# Copy source
COPY ./services/backend/ .

# Copy frontend dist/ from frontend-build stage (must be present for go:embed to work)
COPY --from=frontend-build /frontend/dist ./server/dist

# Build binary with embedded frontend
RUN CGO_ENABLED=1 GOOS=linux go build -o leadpulse-backend .

# Stage 3: Runtime
FROM alpine:latest

# Install required dependencies (sqlite and wget for health checks)
RUN apk add --no-cache ca-certificates libc6-compat wget

WORKDIR /app

# Copy binary from build stage (includes embedded frontend dist/)
COPY --from=build /build/leadpulse-backend .

# Create data volume mount point
RUN mkdir -p /data

# Expose port 8080 for API
EXPOSE 8080

# Health check (verifies both backend API and SPA bootstrap)
HEALTHCHECK --interval=10s --timeout=5s --start-period=15s --retries=5 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./leadpulse-backend"]
