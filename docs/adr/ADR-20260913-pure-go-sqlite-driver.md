# ADR-20260913 — Use modernc.org/sqlite (Pure Go SQLite Driver)

**Decision:** Use `modernc.org/sqlite` as the SQLite driver instead of `mattn/go-sqlite3`.

**Status:** Accepted

---

## Context

Two Go SQLite drivers are in widespread use:
- `github.com/mattn/go-sqlite3` — the dominant choice, uses CGO
- `modernc.org/sqlite` — a pure Go port of SQLite, no CGO

The application targets cross-platform distribution (macOS arm64/amd64, Windows amd64, Linux amd64) compiled via GitHub Actions CI. CGO has implications for cross-compilation.

---

## Alternatives

**`mattn/go-sqlite3` (CGO)**
- Most widely used Go SQLite driver
- Mature, well-tested, large community
- Requires CGO: a C compiler must be present on the build machine for each target platform
- Cross-compilation requires a cross-compiler toolchain (e.g., `zig cc` or platform-specific Docker images) — adds CI complexity
- CGO is disabled by default in some Go build environments

**`modernc.org/sqlite` (pure Go)**
- A machine-translated port of the SQLite C source to Go
- No CGO: standard `GOOS/GOARCH` cross-compilation works without a C toolchain
- Slightly slower than `mattn/go-sqlite3` in benchmarks (single-digit % difference at the scale this app operates)
- Smaller community than `mattn`, but actively maintained
- Full SQLite feature parity for the SQL subset this application uses

**`zombiezen.com/go/sqlite` (CGO with better API)**
- CGO-based; same cross-compilation problem as `mattn`

---

## Rationale

The cross-compilation simplicity is the deciding factor. The application is built on GitHub Actions for three platforms. Using `modernc.org/sqlite` means a standard `GOOS=windows GOARCH=amd64 go build` works without configuring a cross-compiler. The performance difference is negligible for a single-user local app operating on tens of thousands of rows.

The risk of using a less popular driver is low because the SQL used is standard SQLite DDL/DML with no driver-specific extensions.

---

## Consequences

**Better:**
- Cross-compilation works out of the box: `go build` with `GOOS/GOARCH` set, no C toolchain needed
- CI build matrix is simpler — no per-platform C compiler setup
- CGO-related build failures are eliminated

**Harder:**
- Slightly larger binary than `mattn/go-sqlite3` (the C-to-Go translation adds code)
- Smaller community; fewer StackOverflow answers if driver-specific issues arise
- If a SQLite extension requiring native C code is ever needed, migration to `mattn` would be required

---

## Related

- [ADR-20260913-go-fyne-desktop.md](ADR-20260913-go-fyne-desktop.md)
- [ADR-20260913-sqlite-local-storage.md](ADR-20260913-sqlite-local-storage.md)
