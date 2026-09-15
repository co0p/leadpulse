# ADR-20260915: UUID Member IDs in REST API

## Status

Accepted

## Context

The Team Members API exposes member CRUD operations over HTTP. Member identifiers must be included in API responses and used in URL path parameters (`/api/members/{id}`).

The domain layer uses `int64` auto-increment primary keys (assigned by SQLite). The HTTP API must use a different identifier format for two reasons:

1. **No collision risk:** Auto-increment IDs are predictable and can collide if multiple instances or databases exist. REST APIs should use globally unique identifiers to avoid name/ID clashing in multi-tenant or federated scenarios.
2. **API contract stability:** UUIDs remain stable across database resets, migrations, or backup/restore operations. An integer ID might be reused if a record is deleted.

## Decision

The Members API uses **deterministic UUID v5** to map domain member IDs (int64) to UUID strings in JSON responses and path parameters.

**Mapping:**
- UUID namespace: `6ba7b810-9dad-11d1-80b4-00c04fd430c8` (SHA1-based, standard v5 namespace)
- Input: domain member ID as int64 (e.g., `1`, `42`)
- Transformation: convert int64 to big-endian byte string; compute SHA1 hash with namespace
- Output: UUID string (e.g., `9ca3fa77-9ce0-595f-aff7-f03880deac96`)
- Idempotent: the same member ID always produces the same UUID; no randomness

**Implementation:** Go function `memberIDToUUID(int64) string` in `server/handler_members.go` using `github.com/google/uuid`.

**API contract:**
- POST `/api/members` response includes `"id": "<uuid>"` (e.g., `"id": "9ca3fa77-9ce0-595f-aff7-f03880deac96"`)
- GET `/api/members` list includes `"id": "<uuid>"` per member
- PATCH `/api/members/{id}` and DELETE `/api/members/{id}` accept UUID string in path
- Handlers parse UUID from path parameter; validate format; convert to int64 for service calls (reverse map via same SHA1 logic)

## Alternatives Considered

1. **Random UUID v4:** Non-idempotent; API responses would contain different UUIDs for the same member on each request. Breaks caching and traceability.

2. **Strict integer IDs in API:** Simpler, but exposes internal auto-increment sequence to clients. Predictable, collides in multi-instance scenarios, and couples API contract to database design.

3. **UUID as domain primary key:** Requires domain refactor; changes store layer schema. Higher risk; deferred to v2.

## Consequences

**Positive:**
- API responses are stable and cacheable (same member ID → same UUID every time)
- UUIDs prevent external clients from assuming sequential ID assignment
- No client-side UUID generation or validation needed; server is source of truth
- Namespace-based v5 is deterministic and testable without cryptographic randomness

**Negative:**
- Adds a bidirectional mapping layer in handlers (int64 ↔ UUID)
- Handlers must parse and validate UUID format from URL path; slightly more complex than integer path parameters
- UUID strings in JSON are longer than integers (36 vs. 1-2 characters for typical IDs)

**Mitigations:**
- Mapping is implemented once in `memberIDToUUID()` and reused by all handlers
- UUID parsing is validated; malformed UUIDs return 400 error with clear message
- Tests verify idempotence: same member ID produces same UUID across requests

## Related Decisions

- **ADR-20260915-clean-architecture-layering:** Handlers are thin; coordinators own validation and business logic. UUID mapping is handler responsibility (presentation layer concern).
- **CONSTITUTION.md#api-design:** API contracts are stable and language-agnostic. UUIDs satisfy this constraint better than integers.

## Evidence

**Test coverage:**
- `TestAddMemberHandler_Success` — verifies UUID is returned in POST response
- `TestGetMembersHandler_WithMembers` — verifies UUIDs are included in GET list
- `TestMembersAPIIntegration_CRUD` — full-stack test verifies UUID stability (same member used in GET, PATCH, DELETE; no UUID change between calls)

**Implementation:**
- `server/handler_members.go:memberIDToUUID()` — v5 UUID generation with fixed namespace
- Handlers parse `r.PathValue("id")`, validate as UUID, then compute domain ID for service calls

## Test Command

```
go test -race ./server -run TestAddMemberHandler_Success
go test -race ./server -run TestMembersAPIIntegration_CRUD
```

Both pass; UUIDs are stable and correctly mapped.
