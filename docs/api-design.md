# API Design

Design principles and patterns for the HTTP API layer. The API is the boundary between the browser (SPA client) and the Go server. This document establishes how handlers, routes, and responses should be structured to keep business logic in the controller and service layers, not in HTTP handlers.

---

## Core Principle: Handlers as Thin Adapters

A handler is a thin adapter:
- **Input:** HTTP request → parse to plain Go values
- **Delegate:** call `ui/controllers/` → call `service/` → call `engine/`
- **Output:** result → serialize to JSON HTTP response

**Rule:** No business logic in handlers. No conditional statements, loops, or computation. A handler is 5–10 lines of code. If it's longer, the logic belongs in a controller.

**Structure:**
```go
func (h *MemberHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
  // Parse query params if needed
  activeOnly := r.URL.Query().Get("active") == "true"

  // Delegate to controller
  members, err := h.controller.GetMembers(activeOnly)

  // Return response
  h.respondJSON(w, http.StatusOK, members)
  // or:
  h.respondError(w, http.StatusBadRequest, err)
}
```

---

## REST Conventions

### Resource-Oriented URLs

The API is resource-oriented, not action-oriented. Resources are nouns. Actions are HTTP verbs.

| Resource | GET | POST | PATCH | DELETE |
|---|---|---|---|---|
| `/api/members` | List all | Create new | — | — |
| `/api/members/{id}` | Get one | — | Update | Deactivate |
| `/api/entries` | List/query | Create | — | — |
| `/api/entries/{id}` | Get one | — | Update | — |
| `/api/entries/preview` | — | Compute (no persist) | — | — |

### Versioning

The API is unversioned for v1. If backwards-incompatible changes are needed in v2, prefix routes with `/api/v2/...`. The controller and service layers remain unchanged; a new handler adapter layer wraps them with a new request/response shape.

### Content Type

All requests and responses are `application/json`. The handler sets:
```go
w.Header().Set("Content-Type", "application/json")
```

---

## Request Shapes

### Query Parameters (for GET, filtering)

Use query parameters for filtering and pagination:
```
GET /api/members?active=true&limit=20&offset=0
GET /api/entries?member_id=123&month=2026-09
GET /api/entries/previous?member_id=123&month=2026-09
```

Parameters are optional. Provide sensible defaults:
- `active=true` by default (list only active members)
- `limit=100`, `offset=0` by default (first 100 results)
- `month=today's month` by default

Parse in the handler:
```go
activeOnly := r.URL.Query().Get("active") != "false"
memberId := r.URL.Query().Get("member_id")
month := r.URL.Query().Get("month")
if month == "" {
  month = time.Now().Format("2006-01")
}
```

### Path Parameters (for specific resources)

Use path parameters for resource identity:
```
GET /api/members/123
PATCH /api/members/123
DELETE /api/members/123
```

Extract using `mux.Vars(r)` or a router library:
```go
id := mux.Vars(r)["id"]
```

### Request Body (for POST, PATCH)

POST and PATCH requests carry a JSON body. Define request types in a `server/requests.go` file:

```go
// server/requests.go
type AddMemberRequest struct {
  FirstName string `json:"first_name"`
  LastName  string `json:"last_name"`
  Seniority string `json:"seniority"`
}

type SaveEntryRequest struct {
  MemberId  int `json:"member_id"`
  Month     string `json:"month"`
  MoraleScore int `json:"morale_score"`
  BillabilityScore int `json:"billability_score"`
  // ... 10 signal fields
}
```

In the handler:
```go
var req AddMemberRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  h.respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
  return
}
```

---

## Response Shapes

### Success Response

All successful responses are JSON objects with a `data` key:

```json
{
  "data": { /* resource or array */ }
}
```

Examples:

**Single resource (GET /api/members/123):**
```json
{
  "data": {
    "id": 123,
    "first_name": "Alice",
    "last_name": "Smith",
    "seniority": "L4",
    "created_at": "2026-01-15T10:30:00Z"
  }
}
```

**List (GET /api/members):**
```json
{
  "data": [
    { "id": 123, "first_name": "Alice", "seniority": "L4" },
    { "id": 124, "first_name": "Bob", "seniority": "L3" }
  ]
}
```

**Create response (POST /api/members):**
```json
{
  "data": {
    "id": 125,
    "first_name": "Charlie",
    "last_name": "Brown",
    "seniority": "L2",
    "created_at": "2026-09-15T14:22:00Z"
  }
}
```

**Stateless computation (POST /api/entries/preview):**
```json
{
  "data": {
    "tii": 78.5,
    "completeness": 85,
    "dimension_growth": 72,
    "dimension_project": 81,
    "dimension_team": 76,
    "dimension_org": 68
  }
}
```

### Error Response

All error responses are JSON objects with `error` and optional `details`:

```json
{
  "error": "validation_error",
  "message": "First name is required",
  "details": {
    "field": "first_name",
    "reason": "empty"
  }
}
```

Error codes (HTTP status + `error` field):

| HTTP Status | Error Code | Scenario |
|---|---|---|
| 400 | `validation_error` | User input fails validation (out of range, empty required field, duplicate name) |
| 400 | `invalid_request` | Malformed JSON, missing required fields |
| 404 | `not_found` | Resource doesn't exist (member ID, entry month) |
| 409 | `conflict` | Resource already exists (duplicate entry for same member/month) |
| 500 | `internal_error` | Unexpected server error (database corruption, panic) |

Example validation error (signal out of range):
```json
{
  "error": "validation_error",
  "message": "Morale score must be between 1 and 10",
  "details": {
    "field": "morale_score",
    "value": 15,
    "min": 1,
    "max": 10
  }
}
```

Example conflict error (duplicate entry):
```json
{
  "error": "conflict",
  "message": "Monthly entry for member 123 in 2026-09 already exists",
  "details": {
    "member_id": 123,
    "month": "2026-09"
  }
}
```

### HTTP Status Codes

| Status | When |
|---|---|
| 200 | GET, PATCH successful |
| 201 | POST successful (resource created) |
| 204 | DELETE successful (no content to return) |
| 400 | Validation error, malformed request |
| 404 | Resource not found |
| 409 | Conflict (duplicate, constraint violation) |
| 500 | Unexpected server error |

---

## Handler Organization

### File Structure

```
server/
├── server.go          # HTTP server setup, route registration, middleware
├── handlers.go        # Handler struct definitions and method receivers
├── member_handler.go  # Member-related handlers (GET /api/members, POST, etc.)
├── entry_handler.go   # Entry-related handlers (GET /api/entries, POST, etc.)
├── requests.go        # Request type definitions
├── responses.go       # Response helper types and serialization
└── middleware.go      # Logging, error recovery, content-type enforcement
```

### Handler Struct

Handlers are methods on a struct that holds dependencies:

```go
// server/handlers.go
type MemberHandler struct {
  controller *ui.SettingsController
  logger     *log.Logger
}

type EntryHandler struct {
  controller *ui.MonthlyInputController
  logger     *log.Logger
}

// Factory
func NewMemberHandler(ctrl *ui.SettingsController, logger *log.Logger) *MemberHandler {
  return &MemberHandler{
    controller: ctrl,
    logger:     logger,
  }
}
```

### Handler Methods

Each handler method is an HTTP verb on a resource:

```go
// server/member_handler.go

// GET /api/members
func (h *MemberHandler) List(w http.ResponseWriter, r *http.Request) {
  activeOnly := r.URL.Query().Get("active") != "false"
  members, err := h.controller.GetMembers(activeOnly)
  if err != nil {
    h.respondError(w, http.StatusInternalServerError, err)
    return
  }
  h.respondJSON(w, http.StatusOK, members)
}

// GET /api/members/{id}
func (h *MemberHandler) Get(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]
  member, err := h.controller.GetMember(id)
  if err != nil {
    if errors.Is(err, ui.ErrMemberNotFound) {
      h.respondError(w, http.StatusNotFound, err)
      return
    }
    h.respondError(w, http.StatusInternalServerError, err)
    return
  }
  h.respondJSON(w, http.StatusOK, member)
}

// POST /api/members
func (h *MemberHandler) Create(w http.ResponseWriter, r *http.Request) {
  var req AddMemberRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    h.respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
    return
  }

  member, err := h.controller.AddMember(req.FirstName, req.LastName, req.Seniority)
  if err != nil {
    if errors.Is(err, ui.ErrValidation) {
      h.respondError(w, http.StatusBadRequest, err)
      return
    }
    h.respondError(w, http.StatusInternalServerError, err)
    return
  }

  h.respondJSON(w, http.StatusCreated, member)
}

// PATCH /api/members/{id}
func (h *MemberHandler) Update(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]
  var req EditMemberRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    h.respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
    return
  }

  member, err := h.controller.EditMember(id, req.FirstName, req.LastName, req.Seniority)
  if err != nil {
    if errors.Is(err, ui.ErrMemberNotFound) {
      h.respondError(w, http.StatusNotFound, err)
      return
    }
    if errors.Is(err, ui.ErrValidation) {
      h.respondError(w, http.StatusBadRequest, err)
      return
    }
    h.respondError(w, http.StatusInternalServerError, err)
    return
  }

  h.respondJSON(w, http.StatusOK, member)
}

// DELETE /api/members/{id}
func (h *MemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]
  err := h.controller.DeactivateMember(id)
  if err != nil {
    if errors.Is(err, ui.ErrMemberNotFound) {
      h.respondError(w, http.StatusNotFound, err)
      return
    }
    h.respondError(w, http.StatusInternalServerError, err)
    return
  }
  w.WriteHeader(http.StatusNoContent)
}
```

### Response Helpers

Define helper methods to keep handlers short:

```go
// server/responses.go

type JSONResponse struct {
  Data interface{} `json:"data"`
}

type ErrorResponse struct {
  Error   string      `json:"error"`
  Message string      `json:"message"`
  Details interface{} `json:"details,omitempty"`
}

func (h *MemberHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(JSONResponse{Data: data})
}

func (h *MemberHandler) respondError(w http.ResponseWriter, statusCode int, err error) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  
  errCode, message, details := parseError(err)
  json.NewEncoder(w).Encode(ErrorResponse{
    Error:   errCode,
    Message: message,
    Details: details,
  })
}

func parseError(err error) (code, message string, details interface{}) {
  // Map controller/service errors to error codes and messages
  if errors.Is(err, ui.ErrValidation) {
    return "validation_error", err.Error(), nil
  }
  if errors.Is(err, ui.ErrNotFound) {
    return "not_found", err.Error(), nil
  }
  if errors.Is(err, ui.ErrConflict) {
    return "conflict", err.Error(), nil
  }
  return "internal_error", "Unexpected error", nil
}
```

---

## Routing

### Route Registration

Use a router (e.g., `gorilla/mux`, `chi`, or standard library `net/http` with a simple wrapper) to register handlers:

```go
// server/server.go

func NewServer(ctrl *ui.SettingsController, entryCtrl *ui.MonthlyInputController) *http.Server {
  r := mux.NewRouter()

  // Health check
  r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, "OK")
  }).Methods(http.MethodGet)

  // Members API
  memberHandler := NewMemberHandler(ctrl, log.Default())
  r.HandleFunc("/api/members", memberHandler.List).Methods(http.MethodGet)
  r.HandleFunc("/api/members", memberHandler.Create).Methods(http.MethodPost)
  r.HandleFunc("/api/members/{id}", memberHandler.Get).Methods(http.MethodGet)
  r.HandleFunc("/api/members/{id}", memberHandler.Update).Methods(http.MethodPatch)
  r.HandleFunc("/api/members/{id}", memberHandler.Delete).Methods(http.MethodDelete)

  // Entries API
  entryHandler := NewEntryHandler(entryCtrl, log.Default())
  r.HandleFunc("/api/entries", entryHandler.List).Methods(http.MethodGet)
  r.HandleFunc("/api/entries", entryHandler.Create).Methods(http.MethodPost)
  r.HandleFunc("/api/entries/{id}", entryHandler.Get).Methods(http.MethodGet)
  r.HandleFunc("/api/entries/{id}", entryHandler.Update).Methods(http.MethodPatch)
  r.HandleFunc("/api/entries/preview", entryHandler.Preview).Methods(http.MethodPost)
  r.HandleFunc("/api/entries/previous", entryHandler.GetPrevious).Methods(http.MethodGet)

  // SPA static assets (embedded)
  r.PathPrefix("/").Handler(http.FileServer(http.FS(webFS)))

  return &http.Server{
    Addr:    ":8080",
    Handler: r,
  }
}
```

### Middleware

Middleware wraps handlers for cross-cutting concerns:

```go
// server/middleware.go

func loggingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    next.ServeHTTP(w, r)
    log.Printf("%s %s %v", r.Method, r.RequestURI, time.Since(start))
  })
}

func contentTypeMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
      w.WriteHeader(http.StatusUnsupportedMediaType)
      return
    }
    next.ServeHTTP(w, r)
  })
}

// Apply in NewServer:
r.Use(loggingMiddleware)
r.Use(contentTypeMiddleware)
```

---

## Error Handling Strategy

### Controller / Service Errors

Controllers and services raise errors of well-defined types. Handlers catch these and map them to HTTP responses.

**Define error types in the controller package:**

```go
// ui/controllers/errors.go

package ui

import "errors"

var (
  ErrValidation = errors.New("validation error")
  ErrNotFound   = errors.New("not found")
  ErrConflict   = errors.New("conflict")
  ErrInternal   = errors.New("internal error")
)

// Wrap with context:
func NewValidationError(field, reason string) error {
  return fmt.Errorf("%w: %s — %s", ErrValidation, field, reason)
}
```

**In handlers, check for these and map to status codes:**

```go
if errors.Is(err, ui.ErrValidation) {
  h.respondError(w, http.StatusBadRequest, err)
  return
}
if errors.Is(err, ui.ErrNotFound) {
  h.respondError(w, http.StatusNotFound, err)
  return
}
if errors.Is(err, ui.ErrConflict) {
  h.respondError(w, http.StatusConflict, err)
  return
}
// Default to 500
h.respondError(w, http.StatusInternalServerError, err)
```

---

## Architectural Boundaries in the HTTP Layer

### Hexagonal / Clean Architecture Principles

The API layer (handlers) sits at the boundary between the external (browser) and the internal (service). It should:

1. **Translate external format to internal format.** Browser sends JSON. Handlers parse JSON to Go values (request structs). Controllers and services work with domain types.

2. **Never introduce new business logic.** If a handler has conditional logic beyond parsing/error checking, move it to a controller.

3. **Never access repositories or database directly.** The handler calls controllers. Controllers call services. Services call repositories. No shortcuts.

4. **Keep the service contract the same.** If the HTTP API changes shape (new fields, renamed endpoints), the controller and service remain unchanged. Only the handler adapter is rewritten.

---

## Dependency Flow

```
Browser (HTTP)
    ↓
Handler (server/member_handler.go)
    ↓
Controller (ui/controllers/settings_controller.go)
    ↓
Service (service/member/service.go)
    ↓
Store (store/member.go)
    ↓
SQLite + Engine
```

Each layer has a single responsibility:
- **Handler:** Parse HTTP request, call controller, serialize response
- **Controller:** Manage form/list state, call service methods
- **Service:** Enforce use cases, call repository, call engine
- **Store:** Persist aggregates, implement repository interface
- **Engine:** Compute formulas, pure functions

---

## Testing

### Handler Unit Tests

Use `net/http/httptest` to test handlers without a real browser or database:

```go
// server/member_handler_test.go

func TestMemberHandler_ListActive(t *testing.T) {
  // Mock controller
  mockCtrl := &mockSettingsController{
    members: []*Member{{ID: 1, FirstName: "Alice"}},
  }
  handler := NewMemberHandler(mockCtrl, log.New(io.Discard, "", 0))

  // Prepare request
  req, _ := http.NewRequest("GET", "/api/members?active=true", nil)
  rec := httptest.NewRecorder()

  // Call handler
  handler.List(rec, req)

  // Assert
  if rec.Code != http.StatusOK {
    t.Errorf("got %d, want 200", rec.Code)
  }
  var resp JSONResponse
  json.NewDecoder(rec.Body).Decode(&resp)
  if len(resp.Data.([]interface{})) != 1 {
    t.Errorf("got %d members, want 1", len(resp.Data.([]interface{})))
  }
}
```

### Integration Tests

Once handlers are tested, wire a real `ApplicationServices` (with in-memory SQLite) and test end-to-end HTTP flows:

```go
// server/integration_test.go

func TestMemberAPI_CreateAndList(t *testing.T) {
  // Setup in-memory database and real services
  db := setupTestDB(t)
  services := setupServices(db)
  
  // Create HTTP server with real services
  server := NewServer(services.SettingsController, services.MonthlyController)
  client := &http.Client{}

  // POST /api/members
  addReq := AddMemberRequest{FirstName: "Alice", LastName: "Smith", Seniority: "L4"}
  body, _ := json.Marshal(addReq)
  resp, _ := client.Post("http://localhost:8080/api/members", "application/json", bytes.NewReader(body))
  if resp.StatusCode != http.StatusCreated {
    t.Fatalf("create failed with %d", resp.StatusCode)
  }

  // GET /api/members
  resp, _ = client.Get("http://localhost:8080/api/members")
  var listResp JSONResponse
  json.NewDecoder(resp.Body).Decode(&listResp)
  if len(listResp.Data.([]interface{})) != 1 {
    t.Fatalf("got %d, want 1", len(listResp.Data.([]interface{})))
  }
}
```

---

## Checklist for New Endpoints

Before adding a new endpoint:

- [ ] **Resource-oriented.** Is the URL a noun (resource) or a verb (action)? If action, reconsider the design.
- [ ] **Handler is thin.** Can the handler be written in ≤10 lines of code?
- [ ] **Controller method exists.** Does a `ui/controllers/` method exist for this use case?
- [ ] **Error types defined.** Are the possible errors documented and mapped to HTTP status codes?
- [ ] **Request type defined.** Is there a struct in `server/requests.go` for the request body?
- [ ] **Response shape documented.** Is the success response shape documented in this file or in `docs/prd.md` Section 6?
- [ ] **Handler unit tests written.** Do handler unit tests exist (with mocked controller)?
- [ ] **Integration test planned.** Is there an end-to-end HTTP test with real services?

---

## References

- `docs/architecture.md` — Overall system layers and dependency direction
- `docs/ui.md` — Controller pattern and state management
- `CONSTITUTION.md` — Architecture boundaries and service as API contract
- `ADR-20260915-spa-frontend-stack.md` — Frontend technology choice (HTMX, Alpine, Go templates)
