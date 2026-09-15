# ADR-20260915 — Clean Architecture Layering: Coordinators Between HTTP and Services

**Decision:** Introduce a `service/coordinator/` layer that orchestrates multi-step workflows and provides a UI-agnostic contract reusable by HTTP handlers, CLI commands, and future presentation layers.

**Status:** Proposed (adopted during Fyne removal and SPA-first refactor)

---

## Context

The original Fyne UI layer contained controllers (`ui/controllers/`) that managed form state, validation sequencing, and service call orchestration. These controllers were testable (zero Fyne imports) but tightly coupled to the desktop UI abstraction.

When migrating to a web SPA (HTMX + Alpine.js), the same business workflows must be reused: form field validation, multi-step data entry, conditional service calls, error handling. Without a dedicated coordinator layer, this logic would be duplicated in HTTP handlers, creating:

- **Maintenance burden:** Bug fixes in one UI must be replicated in another
- **Inconsistency risk:** Different UIs implement the same workflow slightly differently
- **Testability loss:** Handler tests require full HTTP test infrastructure (`httptest`)
- **Future brittleness:** A third UI (CLI, batch import, API client library) requires yet another implementation

---

## Alternatives

### 1. Put orchestration logic in HTTP handlers
```go
// Bad: HTTP handler owns business logic
func HandleSaveEntry(w http.ResponseWriter, r *http.Request) {
    signals := parseSignals(r)
    if !isComplete(signals) {
        http.Error(w, "Incomplete", 400)
        return
    }
    err := monthlyService.CreateEntry(...)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(200)
}
```
**Downside:** Logic is coupled to HTTP semantics. A CLI command or batch import must reimplement the same validation and sequencing.

### 2. Keep controllers in `ui/controllers/`, make them HTTP-agnostic
```go
// Current approach: controllers import both service and Fyne (via ui/screens/)
type MonthlyInputController struct {
    monthlyService *monthly.Service
    // ...
}
```
**Downside:** Controllers were designed around Fyne abstractions (form state, widget refresh). When Fyne is removed, the controller abstraction loses its context. Renaming them to "coordinators" and moving them to `service/` clarifies their true role: service orchestration.

### 3. Add orchestration logic to service layer
```go
// Bad: Service becomes a workflow engine, not a use case
func (s *MonthlyService) SaveEntryWithValidation(signals map[string]int) error {
    // ... validation, state checks, multi-step calls ...
}
```
**Downside:** Services are transaction-focused (atomic read → compute → write). Mixing orchestration into services blurs responsibility. Services are harder to test in isolation when they depend on coordinator-like validation.

### 4. Introduce coordinators in `service/coordinator/`
```go
// Good: Coordinator is UI-agnostic orchestration
type MonthlyInputCoordinator struct {
    monthlyService *monthly.Service
    memberService  *member.Service
}

func (c *MonthlyInputCoordinator) SaveEntry(memberID, month string, signals map[string]int) error {
    if !c.validateCompleteness(signals) {
        return ErrIncomplete
    }
    return c.monthlyService.CreateEntry(memberID, month, signals)
}
```
**Upside:** Pure Go, no HTTP, no Fyne. Testable with in-memory repos. Reusable by HTTP handlers, CLI, batch import, future UIs.

---

## Rationale

**Option 4 (coordinators) is chosen because it:**

1. **Separates concerns cleanly:**
   - **Service layer:** single atomic use cases (AddMember, CreateEntry, GetTrends)
   - **Coordinator layer:** multi-step workflows (form validation → completeness check → service call → error handling)
   - **HTTP handler layer:** request/response mapping (JSON ↔ coordinator calls)

2. **Enables UI reuse without duplication:**
   ```
   HTTP Handler → Coordinator → Service → Store → Engine
   CLI Command  → Coordinator → Service → Store → Engine
   Future Web   → Coordinator → Service → Store → Engine
   ```
   Validation, sequencing, and error logic live once, in the coordinator.

3. **Improves testability:**
   - Coordinator tests use in-memory repos (same pattern as service tests)
   - No HTTP test framework overhead (`httptest`)
   - No UI framework overhead (no Fyne, no browser headless driver)
   - Tests are fast and deterministic

4. **Preserves the existing service contract:**
   - Services remain unchanged; they are the true use-case boundary
   - Coordinators are a thin orchestration layer that adds no new behavior to services
   - Dependency direction is preserved: `coordinator → service → store → engine`

5. **Future-proofs the architecture:**
   - If a CLI is added, it reuses coordinators without reimplementing validation
   - If batch import is needed, it calls the same coordinators
   - If a third web framework is chosen (Vue, React, etc.), it calls the same HTTP endpoints that delegate to coordinators

---

## Consequences

### Better
- **No logic duplication** across HTTP handlers, CLI, or other UIs
- **Cleaner handler code** — handlers are thin request/response adapters
- **Better testability** — coordinators tested without HTTP or UI framework overhead
- **Explicit workflow orchestration** — readers understand the multi-step flow immediately
- **Consistency across UIs** — all presentation layers follow the same coordinator contracts

### Harder
- **One more layer to navigate** — developers must understand when to add logic to service vs. coordinator
- **Naming clarity** — "coordinator" is less familiar than "controller"; requires documentation (this ADR)

---

## Implementation Rules

### Coordinator Responsibility
- **Accept:** application-level inputs (member IDs, month strings, signal maps)
- **Validate:** pre-service checks (completeness, range, business rules that must pass before any service call)
- **Orchestrate:** call services in correct order, handle cross-service dependencies
- **Return:** result or error; error types are domain types (not HTTP codes)

### Coordinator Constraints
- **No HTTP imports** — no `net/http`, no `http.Request`, no status codes
- **No UI framework imports** — no Fyne, no web framework
- **No I/O beyond service calls** — no direct database, file, or network access
- **No side effects** — deterministic functions; idempotent where possible
- **Stateless** — no member variables that persist across calls (except service references)

### HTTP Handler Responsibility
- **Parse request** — JSON body, query params, path variables
- **Call coordinator** — pass application-level inputs
- **Map result to response** — coordinator result → JSON response, error → HTTP status code
- **Log and monitor** — handler is the request/response boundary

**Example:**
```go
// Handler layer (server/handler.go)
func HandleSaveEntry(w http.ResponseWriter, r *http.Request) {
    signals := parseSignals(r.Body)  // HTTP parsing
    
    // Delegate to coordinator
    err := coordinator.SaveEntry(memberID, month, signals)
    
    // Map result to HTTP response
    if err == coordinator.ErrIncomplete {
        http.Error(w, "Incomplete", 400)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(200)
}

// Coordinator layer (service/coordinator/monthly.go)
func (c *MonthlyCoordinator) SaveEntry(memberID, month string, signals map[string]int) error {
    if !c.validateCompleteness(signals) {
        return ErrIncomplete  // Domain error, not HTTP
    }
    return c.monthlyService.CreateEntry(memberID, month, signals)
}

// Service layer (service/monthly/monthly.go)
func (s *Service) CreateEntry(memberID, month string, signals map[string]int) error {
    // Atomic operation: validate aggregate, compute scores, persist
    return s.repo.SaveEntry(...)
}
```

---

## Related Decisions

- **ADR-20260915-spa-frontend-stack.md** — HTMX + Alpine.js + Go templates; coordinators are called via HTTP handlers from the SPA
- **CONSTITUTION.md — Service is the API contract** — Coordinators are above services; handlers call coordinators, not services directly
- **CONSTITUTION.md — Dependency direction is a hard rule** — Handlers → Coordinator → Service → Store → Engine; no reverse imports

---

## Migration from Controllers

**Old pattern (Fyne UI):**
```
Fyne Screen (ui/screens/monthly_input.go)
    ↓ (calls)
    ↓
Controller (ui/controllers/monthly_input_controller.go)
    ↓ (calls)
    ↓
Service (service/monthly/monthly.go)
```

**New pattern (SPA + CLI-ready):**
```
HTTP Handler (server/handler.go)
    ↓ (calls)
    ↓
Coordinator (service/coordinator/monthly.go)
    ↓ (calls)
    ↓
Service (service/monthly/monthly.go)
```

Controllers are renamed to coordinators and moved from `ui/controllers/` to `service/coordinator/`. Logic is preserved; abstraction context changes from "UI-aware form management" to "UI-agnostic workflow orchestration."

---

## Checklist for New Coordinators

When adding a new coordinator:
- [ ] Lives in `service/coordinator/`
- [ ] Has zero HTTP imports
- [ ] Has zero UI framework imports
- [ ] Tests use in-memory repos (same pattern as service tests)
- [ ] Accepts application-level inputs (not HTTP types)
- [ ] Returns domain error types (not HTTP status codes)
- [ ] Documented in `docs/coordinator.md` (if pattern is new) or this ADR

---

**Last updated:** 2026-09-15 — Proposed during Fyne removal and SPA-first refactor
