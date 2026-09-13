BINARY_NAME := leadpulse
BUILD_DIR := build
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: help build run test clean lint

help:
	@echo "Leadpulse Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build    - Build the leadpulse binary to build/"
	@echo "  run      - Build and run the leadpulse application"
	@echo "  test     - Run all tests with race detection"
	@echo "  test-v   - Run all tests verbosely with race detection"
	@echo "  clean    - Remove the build/ directory"
	@echo "  lint     - Run go vet for static analysis"
	@echo "  help     - Show this message"

build:
	@echo "Building leadpulse..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY_PATH) .

run: build
	@echo "Running leadpulse..."
	./$(BINARY_PATH)

test:
	@echo "Running tests with race detection..."
	go test -race ./...

test-v:
	@echo "Running tests verbosely with race detection..."
	go test -race ./... -v

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

lint:
	@echo "Running go vet..."
	go vet ./...
