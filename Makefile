.PHONY: help build run test clean lint

help:
	@echo "Leadpulse Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build    - Build the leadpulse binary"
	@echo "  run      - Run the leadpulse application"
	@echo "  test     - Run all tests with race detection"
	@echo "  test-v   - Run all tests verbosely with race detection"
	@echo "  clean    - Remove the built binary and temp files"
	@echo "  lint     - Run go vet for static analysis"
	@echo "  help     - Show this message"

build:
	@echo "Building leadpulse..."
	go build -o leadpulse .

run: build
	@echo "Running leadpulse..."
	./leadpulse

test:
	@echo "Running tests with race detection..."
	go test -race ./...

test-v:
	@echo "Running tests verbosely with race detection..."
	go test -race ./... -v

clean:
	@echo "Cleaning build artifacts..."
	rm -f leadpulse

lint:
	@echo "Running go vet..."
	go vet ./...
