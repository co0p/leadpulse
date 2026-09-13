# ADR-20260913: Domain-Driven Design Refactoring — Value Objects, Aggregates, and Repositories

## Status
Accepted

## Date
2026-09-13

## Context

The engine and service layers were using plain structs (TeamMember, MonthlyEntry) with validation logic scattered across service methods. This approach had several drawbacks:

1. **Constraint enforcement was diffuse.** Validation rules lived in service methods, not at the type level. Invalid states (e.g., morale = 10) could exist in memory if you bypassed the service.

2. **Encapsulation was weak.** All domain fields were public and exported, allowing callers to mutate state directly or bypass invariants.

3. **Persistence was tightly coupled.** Service layer depended on `*sql.DB` directly, making it impossible to test without a database.

4. **Cross-aggregate rules were implicit.** Business rules like "only active members can have entries" and "one entry per member per month" were buried in service code, hard to discover and verify.

5. **Future features (alerts, trends, validation) would inherit these problems,** making the codebase harder to extend safely.

A research prototype validated that Domain-Driven Design patterns could address all of these issues while preserving existing behavior.

## Decision

Refactor the engine and service layers to use DDD patterns:

1. **Value Objects** for every domain concept (signals, scores, composite types like FullName and ImpactRating). Each value object enforces its constraints in the constructor; invalid inputs are rejected immediately.

2. **Aggregate Roots** (TeamMember, MonthlyEntry) that manage identity and state transitions. All fields are private; access is controlled via methods that enforce business rules.

3. **Repository Interfaces** (TeamMemberRepository, MonthlyEntryRepository) that abstract persistence. Domain logic depends on interfaces, not concrete storage.

4. **Domain Services** (MonthlyEntryService) that enforce rules spanning multiple aggregates (e.g., "only active members can have entries") before persistence.

5. **Layering:** Domain logic in `engine/domain/` (value objects, aggregates, domain services, interfaces); storage implementation in `store/` (implements the interfaces); application logic in `service/` (depends on repository interfaces, not on database).

## Alternatives Considered

### 1. Keep the current approach, add more validation layers
**Pro:** Minimal refactoring, lower risk.
**Con:** Validation stays scattered; constraint enforcement remains weak; service tests still require a database; future features inherit the same problems.

### 2. Use an ORM or query builder instead of raw SQL
**Pro:** Reduces boilerplate.
**Con:** Does not solve the constraint enforcement, encapsulation, or decoupling problems.

### 3. Add a type-driven approach at the service level (not domain level)
**Pro:** Avoids changing the domain layer.
**Con:** Service layer becomes cluttered; domain objects remain mutable; repository abstraction still requires service-layer scaffolding.

**Why the chosen approach wins:** Value objects and aggregates move constraints and invariants to the domain layer, where they belong. The repository interface creates a clean boundary for testing and future persistence strategies. Domain services make cross-aggregate rules explicit and testable.

## Consequences

### Positive

1. **Type-level constraint enforcement.** Invalid states are impossible to construct. Callers cannot pass morale = 10 or create a deactivated member twice.

2. **Better encapsulation.** Aggregate fields are private; state transitions are guarded by methods. Invariants cannot be violated.

3. **Testability without a database.** Service and domain logic can be tested with in-memory repository implementations. No database setup required for unit tests.

4. **Explicit cross-aggregate rules.** Business rules (one entry per member per month, only active members can have entries) are now explicit in domain services and easy to verify.

5. **Storage independence.** If persistence needs change in the future (SQLite → PostgreSQL, in-memory → cloud database), only the repository implementation changes. Domain and service logic remain unchanged.

6. **Clearer code for future developers.** Value objects, aggregates, and repositories are well-understood DDD patterns. New features (alerts, trends, validation) will have a consistent pattern to follow.

### Negative (Managed Risks)

1. **More boilerplate in aggregates.** State management in aggregates requires more code than plain structs. This is intentional—the extra code enforces invariants.

2. **Repository implementations are domain-aware.** The SQLite repository must know how to reconstruct aggregates from database rows. If the aggregate changes, the repository must update. (Mitigation: repository is tested; changes are caught quickly.)

3. **Service layer now depends on repository interfaces.** This reverses some of the old dependency direction. (Mitigation: the new direction is cleaner—domain owns the interfaces, not storage.)

## Implementation Notes

- **Backward compatibility:** Old store functions (AddMember, GetMember, etc.) are kept as thin wrappers around the new repository for compatibility with existing code. They will be deprecated.

- **No behavior change.** The refactoring is purely structural. User-facing behavior, data flow, and output remain identical.

- **Performance:** No performance-critical paths are affected. All existing formulas and scoring logic remain unchanged.

- **Testing:** All existing tests pass. New tests exercise value object constraints, aggregate state transitions, and repository operations.

## Related

- `docs/architecture.md` — updated to show repository interface as a boundary
- `docs/domain.md` — updated with value object and aggregate concepts
- `docs/testing.md` — updated with guidance on in-memory repository testing
- `.agent/prototype.md` — research spike that validated the DDD approach (archived)

## Questions for Future Work

1. **Should we add a transaction abstraction?** Currently, each repository operation is independent. If atomicity across aggregates is needed, a Unit of Work pattern could be added later.

2. **Should we add an event log?** Domain events (MemberDeactivated, EntryComputed, etc.) could be captured for audit and future analytics.

3. **Should we expose domain services as a public API?** Currently, domain services are internal to the service layer. If external systems need to consume domain logic, this might change.

All three are reversible decisions; they do not block this refactoring.
