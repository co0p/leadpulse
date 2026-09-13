# Leadpulse

A team member management application built with Go and Fyne.

## Quick Start

### Build

```bash
make build
```

This compiles the application into a `leadpulse` binary.

### Run

```bash
make run
```

This builds and runs the application. The app will launch a GUI using the Fyne framework.

Alternatively, if already built:
```bash
./leadpulse
```

### Test

Run all tests with race detection:
```bash
make test
```

For verbose output:
```bash
make test-v
```

### Clean

Remove the built binary:
```bash
clean
```

### Lint

Run static analysis:
```bash
make lint
```

## Project Structure

- `engine/domain/` — Domain types (TeamMember, Seniority)
- `service/member/` — Business logic and use cases
- `store/` — Data persistence layer (SQLite)
- `ui/` — GUI components and screens (Fyne)
- `main.go` — Application entry point

## Architecture

The application follows a layered architecture:

1. **Domain Layer** (`engine/domain/`) — Immutable value objects
2. **Store Layer** (`store/`) — SQLite persistence with audit logging
3. **Service Layer** (`service/`) — Business logic and validation
4. **UI Layer** (`ui/`) — Fyne-based GUI

## Testing

- **Store tests** — Integration tests using in-memory SQLite
- **Service tests** — Unit tests with mocked store

Run tests:
```bash
go test -race ./...
```

## Making a Change

1. Make your change
2. Run tests: `make test`
3. Build: `make build`
4. Run the app: `./leadpulse`
