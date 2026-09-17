# Stage 1: Test
FROM golang:1.25-alpine AS test

WORKDIR /build

# Install build dependencies (gcc for cgo)
RUN apk add --no-cache gcc musl-dev

# Copy module files
COPY go.mod go.sum ./

# Copy source
COPY . .

# Run tests with race detector
RUN go test -race ./...

# Stage 2: Build
FROM golang:1.25-alpine AS build

WORKDIR /build

# Install build dependencies (gcc for cgo)
RUN apk add --no-cache gcc musl-dev

# Copy module files
COPY go.mod go.sum ./

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -o leadpulse-backend .

# Stage 3: Runtime
FROM alpine:latest

# Install required dependencies (sqlite and wget for health checks)
RUN apk add --no-cache ca-certificates libc6-compat wget

WORKDIR /app

# Copy binary from build stage
COPY --from=build /build/leadpulse-backend .

# Create data volume mount point
RUN mkdir -p /data

# Expose port 8080 for API
EXPOSE 8080

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=15s --retries=5 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./leadpulse-backend"]
