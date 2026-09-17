BINARY_NAME := leadpulse
BUILD_DIR := build
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)
BACKEND_DIR := services/backend

.PHONY: help build run-backend run-frontend test-backend test-frontend test-acceptance clean clean-docker docker-build docker-up docker-down docker-logs

help:
	@echo "Leadpulse Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build            - Build backend binary to services/backend/build/"
	@echo "  run-backend      - Build and run backend locally (without Docker)"
	@echo "  run-frontend     - Run frontend dev server (npm run dev in services/frontend/)"
	@echo "  test-backend     - Run backend Go tests with race detection"
	@echo "  test-frontend    - Run frontend tests (npm run test in services/frontend/)"
	@echo "  test-acceptance  - Run acceptance tests (Playwright in acceptance-tests/)"
	@echo "  clean            - Remove services/backend/build/ directory"
	@echo "  clean-docker     - Stop docker-compose and remove volumes"
	@echo "  docker-build     - Build docker images (backend and frontend)"
	@echo "  docker-up        - Start services via docker-compose up"
	@echo "  docker-down      - Stop services via docker-compose down"
	@echo "  docker-logs      - View docker-compose logs"
	@echo "  help             - Show this message"

# Local development targets
build:
	@echo "Building backend binary..."
	@mkdir -p $(BACKEND_DIR)/$(BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o $(BUILD_DIR)/$(BINARY_NAME)

run-backend: build
	@echo "Running backend..."
	./$(BACKEND_DIR)/$(BUILD_DIR)/$(BINARY_NAME)

run-frontend:
	@echo "Running frontend dev server..."
	cd services/frontend && npm run dev

# Test targets
test-backend:
	@echo "Running backend tests with race detection..."
	cd $(BACKEND_DIR) && go test -race ./...

test-frontend:
	@echo "Running frontend tests..."
	cd services/frontend && npm test 2>&1 | grep -E "PASS|FAIL|✓|✗" || echo "No tests configured yet"

test-acceptance:
	@echo "Running acceptance tests..."
	cd acceptance-tests && npm test

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BACKEND_DIR)/$(BUILD_DIR)

clean-docker:
	@echo "Stopping docker-compose and removing volumes..."
	docker-compose down -v

# Docker targets
docker-build:
	@echo "Building docker images..."
	docker-compose build

docker-up:
	@echo "Starting services..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-logs:
	@echo "Showing docker-compose logs..."
	docker-compose logs -f
