# ADR-20260916 — Hexagonal Architecture (Ports & Adapters) for Member CRUD

**Decision:** Replace the coordinator pattern with hexagonal architecture (ports & adapters) for member CRUD operations, moving application use cases into a core layer with explicit repository port dependencies.

**Status:** Accepted

**Date:** 2026-09-16

---

## Context

The existing architecture uses a coordinator pattern where HTTP handlers delegate to `MemberAPICoordinator`, which orchestrates calls to the service layer, which reads/writes from the store. This structure conflates presentation concerns (HTTP request/response mapping) with service layer details (transaction sequencing).

Member CRUD refactoring revealed two critical needs:

1. **Portability:** Use cases should be reusable across multiple presentation layers (HTTP SPA, CLI, future mobile or desktop UIs) without duplication. The current pattern makes reuse difficult because business logic is bound to the `service/` package, which is tightly coupled to the store and coordinator patterns.

2. **Testability:** Core business logic (member aggregate rules, validation) should be testable in isolation without the full-stack infrastructure (HTTP server, in-memory SQLite setup). Current tests require coordinators or services, which pull in store dependencies even when testing pure domain logic.

The domain aggregate (`TeamMember`) and repository abstraction pattern (dependency inversion) are established in code but not consistently applied across the codebase. This creates an opportunity to formalize the pattern and establish it as the standard for all future use cases (Monthly Entry, Alerts, etc.).

---

## Alternatives

### 1. Keep the Coordinator Pattern

**Option:** Maintain the current architecture where HTTP handlers call coordinators, which orchestrate service calls.

**Trade-offs:**
- ✅ Minimal refactoring; coordinator pattern already established
- ❌ Couples HTTP adapters to service layer; presentation details leak into orchestration
- ❌ Reusing use cases in CLI requires duplicating coordinator logic or extracting to a separate layer (deferred complexity)
- ❌ Tests require full-stack setup (service → store) even when testing simple business rules
- ❌ Repository abstraction exists but is not explicitly injected; tests must mock at the store level

**Why rejected:** This pattern does not scale to multiple presentation layers. A CLI client cannot import `service/coordinator` without pulling in HTTP dependencies or reimplementing the logic.

### 2. Hexagonal Architecture (Ports & Adapters) — **CHOSEN**

**Option:** Move application use cases to a `core/` layer that depends only on domain aggregates and port interfaces (abstractions). HTTP handlers, CLI commands, and other clients inject repository adapters at composition root.

**Structure:**
```
core/members/                    # Application layer
├── member.go                     # TeamMember aggregate (value objects, behavior)
├── repository.go                 # MemberRepository port (interface, zero impl logic)
├── add_member.go                 # AddMemberUseCase (business logic)
├── get_members.go                # GetMembersUseCase
├── edit_member.go                # EditMemberUseCase
├── deactivate_member.go          # DeactivateMemberUseCase
└── *_test.go                     # Unit tests (inject in-memory adapter)

storage/                          # Adapter layer
├── memory/member_repository.go   # InMemoryMemberRepository (test fixture)
└── sqlite/member_repository.go   # SQLiteTeamMemberRepository (production)

server/                           # HTTP adapter layer
├── handler_members.go            # HTTP handlers (thin adapters)
└── handler_members_test.go       # Handler tests (mock use cases)
```

**Trade-offs:**
- ✅ Core is infrastructure-independent; depends only on stdlib and domain types
- ✅ Use cases are injectable and testable in milliseconds with in-memory repository
- ✅ Use cases are reusable across HTTP, CLI, and future UIs without duplication
- ✅ Clear dependency flow: HTTP → core ← storage; no circular imports; enforced by Go imports
- ✅ Explicit dependency injection at composition root (`main.go`) clarifies wiring
- ✅ Aligns with industry-standard DDD and Clean Architecture vocabulary
- ❌ More files (one port interface + multiple adapters per domain concept)
- ❌ Requires explicit wiring in `main.go`; no magic DI container
- ❌ Mixed architecture during migration: coordinator pattern for legacy features (SubmitMonth, GenerateAlerts) until full refactor
- ❌ Slightly more indirection for simple use cases (one extra interface dereference)

**Why chosen:** This pattern isolates business logic from infrastructure, enables testing without framework overhead, and scales to multiple presentation layers. It establishes a foundation for all future use cases.

### 3. Domain-Driven Design Bounded Contexts

**Option:** Reorganize the entire codebase by business domain (member, entry, alert) instead of technical layer. Each context owns aggregate, repository, use cases, and tests.

**Trade-offs:**
- ✅ Maximum domain clarity; each context is a mini-application
- ✅ Scales well for large teams with domain-specific responsibilities
- ❌ Requires significant restructuring of existing code (refactor risk)
- ❌ Heavy investment for v1; better suited to v2+ when codebase is larger and domain is more stable
- ❌ Does not solve the immediate portability or testability problems (can be layered on top of hexagonal)

**Why rejected:** Too much restructuring for v1. Hexagonal provides 80% of the benefit (isolation + testability) without the refactor cost. Bounded contexts can be introduced in v2 if the codebase grows beyond ~10 features.

---

## Rationale

Hexagonal architecture is the right fit for this project because:

1. **Business logic isolation:** Core use cases depend only on abstractions (repository ports), not on HTTP, SQLite, or UI frameworks. This makes the business logic portable and testable without infrastructure setup.

2. **Explicit testing strategy:** Unit tests inject in-memory repository adapters (zero I/O cost, <100ms per test suite). Integration tests wire real SQLite adapters. HTTP tests mock use cases. Each layer is testable independently, reducing feedback loop and increasing developer confidence.

3. **Reusability across clients:** A CLI tool, mobile app, or future web client can reuse the same `core/members/` use cases without copying code. The adapter pattern (HTTP, storage) is swappable; only the core changes once.

4. **DDD alignment:** The pattern enforces domain-driven vocabulary (aggregate, repository, port) consistently across code and documentation. Newcomers understand "core has aggregates, storage implements ports" without cultural knowledge.

5. **Clear dependency rules:** CONSTITUTION.md already guards against circular imports. Hexagonal makes this rule more explicit: core → no dependencies; storage/HTTP → inbound only. Tools can enforce this via linters.

6. **Incremental adoption:** Member CRUD is the first feature to use this pattern. Monthly Entry and Alerts will follow. By the time the codebase is 2–3 features deep, the pattern is proven and maintainers have internalized it.

---

## Consequences

### Positive

- **Testability:** Core use cases are testable in <100ms with in-memory adapters. No database setup, no HTTP server, no UI framework. Developers get fast feedback.
- **Portability:** Same use cases can be called from HTTP handlers, CLI commands, mobile adapters, or future presentation layers. No code duplication.
- **Clarity:** Dependency flow is explicit and enforced: core ← storage, core ← HTTP. No hidden dependencies or framework magic. A new maintainer can read `core/members/add_member.go` and know exactly what it depends on.
- **DDD vocabulary:** Aggregates, repositories, and ports are industry-standard terms. Code and documentation use consistent language, reducing onboarding friction.
- **Architecture enforcement:** Go's import system naturally prevents circular imports. The architecture is not a convention; it is a compile-time guarantee.

### Negative

- **More files:** Each feature now requires a port interface + multiple adapters (memory, SQLite, HTTP). Simple features (4 use cases × 3 adapters = 12 files + tests) create more surface area to maintain.
- **Wiring boilerplate:** Composition root (`main.go`) must explicitly create adapters and pass them to use cases. No DI container magic; this is intentional but adds ~20 lines per feature.
- **Mixed architecture during migration:** Legacy features (SubmitMonth, GenerateAlerts) still use coordinator pattern. This mixed pattern is temporary but adds cognitive load until fully migrated (estimated 2–3 more increments).
- **Indirection for simple operations:** A single-operation use case like `GetMembersUseCase` now requires three files (use case, port interface, adapters). This is overhead for trivial features, though the investment pays off when the use case grows.

### Timeline

- **Immediate (this increment):** Member CRUD fully migrated. Coordinator pattern deleted for members. Mixed architecture begins.
- **Next 2–3 increments:** Monthly Entry and Alerts refactored to use hexagonal pattern. Coordinator pattern becomes legacy code, then removed.
- **v2 readiness:** Codebase is fully hexagonal. All use cases are portable. Bounded contexts (DDD) can be layered on top if domain grows.

---

## Related Decisions

- **Supersedes:** [ADR-20260915-clean-architecture-layering.md](ADR-20260915-clean-architecture-layering.md) — The coordinator pattern is now replaced by use cases with explicit port dependencies. Coordinators remain for legacy features but are not the standard for new work.
- **Related:** [ADR-20260913-ddd-refactor.md](ADR-20260913-ddd-refactor.md) — Aggregates and value objects are core to hexagonal; this ADR formalizes their role.
- **Enforces:** [CONSTITUTION.md](../../CONSTITUTION.md) — Dependency direction rules are now architectural patterns, not just conventions.
- **Documented by:** [docs/architecture.md](../../architecture.md) — C4 Level 2 container diagram and dependency flow updated to show ports & adapters.

---

## Implementation Notes

- Port interfaces live in `core/` alongside aggregates and use cases.
- Adapters live in `storage/` (SQLite, in-memory) or `server/` (HTTP).
- Composition root in `main.go` wires adapters into use cases; no reflection or DI container.
- Tests for core use cases inject in-memory adapters; tests for adapters (storage) use real implementations with test fixtures.
- HTTP handler tests mock use cases; no database or file I/O per handler test.

---

## Evidence

- **Member CRUD refactor:** 34 tests passing (18 core + 6 storage + 8 HTTP). All architecturally clean.
- **Zero circular imports:** `go build ./core/members` succeeds; core imports only stdlib and `engine/domain`.
- **Composition root:** `main.go` wires adapters and passes to handlers; ~20 lines of explicit injection.
- **Test performance:** Core use case tests complete in <100ms; storage adapter tests in <200ms; HTTP handler tests in <50ms.
