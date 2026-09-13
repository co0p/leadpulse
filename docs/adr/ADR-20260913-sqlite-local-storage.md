# ADR-20260913 — Use SQLite for Local-Only Persistence

**Decision:** Use SQLite as the sole data store, stored as a single file on the user's device.

**Status:** Accepted

---

## Context

The application must persist 24 months of monthly entries per team member, evidence notes, action plans, an audit trail, and cycle state. All data must remain on the user's device with no external service dependency.

The primary question was: how to store structured relational data locally in a desktop Go application?

---

## Alternatives

**SQLite (local file)**
- Embedded, zero-server, single file
- Full SQL with foreign keys, constraints, and indexes
- Proven at the scale required (tens of thousands of rows over 24 months for a typical team)
- Backup is a file copy
- Go ecosystem has mature drivers

**JSON flat files (one file per member per month)**
- No dependency
- No query language — all filtering, sorting, and aggregation in Go
- Trend queries (MA3, Delta3) require loading multiple files and joining in memory
- Referential integrity (e.g., evidence linked to a month) is manual
- File proliferation with a large team and long history

**BoltDB / bbolt (embedded key-value)**
- Single binary file, no CGO
- No SQL — all range queries and joins in Go
- Schema evolution requires manual key-format migration
- Less expressive for the relational data this application needs

**PostgreSQL / MySQL (local server)**
- Full relational power
- Requires the user to install and run a database server — unacceptable for a standalone tool
- Operational overhead is incompatible with the deployment model

---

## Rationale

SQLite is the correct fit. The data model is relational (members, monthly entries, evidence, actions, audit log all cross-reference each other). SQL handles the trend and aggregation queries that would otherwise require custom in-memory logic. A single file is easy to back up and understand. There is no user-visible operational overhead.

Flat files are viable for simple write/read but become painful for the queries this application requires (multi-month trend history, team aggregations, audit log scans).

bbolt would work but forces all query logic into Go, duplicating what SQL provides for free, with more complex schema migration.

---

## Consequences

**Better:**
- Relational model fits the domain naturally
- SQL handles trend queries, aggregations, and joins
- Single-file backup is simple to explain to users
- Schema is inspectable with standard SQLite tools
- 24-month history for a 100-member team is well within SQLite's performance envelope

**Harder:**
- Schema migrations must be managed manually (additive-only in v1; see `docs/deployment.md`)
- SQLite's write concurrency model (one writer at a time) is fine for a single-user app but would not scale to multi-user
- Corruption recovery is limited — no WAL-based point-in-time recovery in v1

---

## Related

- [ADR-20260913-go-fyne-desktop.md](ADR-20260913-go-fyne-desktop.md)
- [ADR-20260913-pure-go-sqlite-driver.md](ADR-20260913-pure-go-sqlite-driver.md)
