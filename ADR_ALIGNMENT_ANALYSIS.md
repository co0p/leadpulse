# ADR Alignment Analysis - LeadPulse Codebase

**Analysis Date:** September 13, 2026  
**Codebase Version:** Current Development Build  
**Analysis Scope:** 4 ADRs across DDD, Desktop UI, SQLite Driver, and Local Storage  
**Total Lines of Code Analyzed:** ~2,957 lines (domain, store, service, UI layers)

---

## Executive Summary

The LeadPulse codebase demonstrates **strong alignment** with three of the four ADRs:
- **ADR-20260913-ddd-refactor:** 85% aligned - Good implementation with minor gaps
- **ADR-20260913-go-fyne-desktop:** 70% aligned - Basic Fyne structure in place, UI development incomplete
- **ADR-20260913-pure-go-sqlite-driver:** 95% aligned - Properly implemented
- **ADR-20260913-sqlite-local-storage:** 80% aligned - Schema and basic persistence working

**Critical Finding:** The codebase has made excellent progress on DDD architecture and database decisions, but the UI layer is still rudimentary and the integration between layers needs strengthening.

---

## Part 1: ADR-20260913-DDD Refactor — Value Objects, Aggregates, Repositories

### Status: 85% Aligned ✓ (Mostly Implemented)

### What Has Been Implemented Correctly

#### 1. Value Objects (100% ✓)
All domain concepts properly enforce constraints:

**File:** `/engine/domain/value_objects.go`
- **FullName:** Validates non-empty first and last names ✓
- **ImpactRating:** Validates 0–5 range for all dimensions ✓
- **NormalizedScore, ContributionScore:** Validates 0–100 range ✓
- **CompletenessPercent, ConfidenceScore:** Proper validation ✓

**Signal Type Value Objects:** `/engine/domain/signal_types.go`
- **MoraleScore:** 0–5 validation ✓
- **BillabilityPercent:** 0–100 validation ✓
- **CSATScore:** 1–5 validation ✓
- **NetMarginPercent:** -20 to +60 validation ✓
- **FeedbackCount, Hours, Percent:** >= 0 validation ✓

All value objects use constructor functions (`NewXxx`) to enforce constraints immediately.

#### 2. Aggregate Roots (95% ✓)

**TeamMember Aggregate** `/engine/domain/aggregates.go:17–93`
- ✓ Private fields enforcing encapsulation
- ✓ Immutable ID (`TeamMemberID`)
- ✓ Controlled state transitions:
  - `Deactivate()` enforces single deactivation
  - `ChangeSeniority()` validates before change
- ✓ Access via methods (`ID()`, `Name()`, `Seniority()`, `IsActive()`)
- ✓ Aggregate creation enforces all invariants in `NewTeamMember()`

**MonthlyEntry Aggregate** `/engine/domain/aggregates.go:95–359`
- ✓ Private fields with controlled access
- ✓ Composite ID (`MonthlyEntryID` = member + month)
- ✓ Sub-aggregate `MonthlyRawSignals` with validation
  - All 10 signals validated in `NewMonthlyRawSignals()`
  - Signal count tracking
- ✓ Computed scores handling:
  - `SetComputedScores()` validates before persistence
  - Immutable once set (no public mutators)
- ✓ Timestamps properly managed (`CreatedAt`, `ComputedAt`)

**Minor Issue:** Repository methods `SetCreatedAt()` and `SetComputedAt()` break encapsulation slightly, but are marked as internal use only.

#### 3. Repository Interfaces (100% ✓)

**File:** `/engine/domain/repositories.go:1–42`

**TeamMemberRepository Interface:**
```go
type TeamMemberRepository interface {
    Save(member *TeamMember) error
    FindByID(id TeamMemberID) (*TeamMember, error)
    FindActive() ([]*TeamMember, error)
    Delete(id TeamMemberID) error
}
```

**MonthlyEntryRepository Interface:**
```go
type MonthlyEntryRepository interface {
    Save(entry *MonthlyEntry) error
    FindByID(id MonthlyEntryID) (*MonthlyEntry, error)
    FindByMember(memberID TeamMemberID) ([]*MonthlyEntry, error)
    Delete(id MonthlyEntryID) error
}
```

✓ Clean abstraction with no storage-specific concerns  
✓ Domain layer depends on interfaces, not implementations  
✓ Enabling in-memory repository testing

#### 4. Repository Implementations (100% ✓)

**In-Memory Repositories** `/engine/domain/repositories.go:44–178`
- ✓ `InMemoryTeamMemberRepository` - Full implementation
- ✓ `InMemoryMonthlyEntryRepository` - Full implementation
- ✓ Thread-safe with sync.RWMutex
- ✓ Ready for unit testing without database

**SQLite Repositories** `/store/`
- ✓ `SQLiteTeamMemberRepository` `/store/member.go:11–224`
  - Implements full interface
  - Properly reconstructs aggregates from rows
  - Audit logging for changes
- ✓ `SQLiteMonthlyEntryRepository` `/store/monthly_entry.go:12–317`
  - Implements full interface
  - Handles JSON serialization of complex computed scores
  - Proper timestamp management

#### 5. Domain Services (90% ✓)

**File:** `/engine/domain/services.go`

**MonthlyEntryService** enforces cross-aggregate rules:
```go
// CreateEntry enforces:
// - Member must exist and be active
// - Only one entry per member per month
// - All constraints via NewMonthlyEntry()
```

✓ Enforces business rule: "only active members can have entries"  
✓ Enforces business rule: "one entry per member per month"  
✓ Proper error propagation

**Issue:** Service layer only has `MonthlyEntryService`. Could benefit from more explicit domain services for scoring and validation pipelines.

#### 6. Service Layer Dependency Injection (100% ✓)

**Member Service** `/service/member/member.go:1–166`
```go
type Service struct {
    repo domain.TeamMemberRepository  // Injected interface
}

func NewService(repo domain.TeamMemberRepository) *Service {
    return &Service{repo: repo}
}
```

✓ Depends on repository interface, not implementation  
✓ Service methods use domain aggregates and value objects  
✓ No direct database calls

**Monthly Service** `/service/monthly/monthly.go:1–123`
```go
type Service struct {
    memberRepo domain.TeamMemberRepository
    entryRepo  domain.MonthlyEntryRepository
}
```

✓ Both repositories injected as interfaces  
✓ Delegates to domain service for cross-aggregate rules

### What Is Missing or Not Aligned

#### 1. Incomplete Domain Service Coverage
**Gap:** Only `MonthlyEntryService` exists. Missing:
- Scoring pipeline as a domain service
- Trend analysis domain service
- Validation orchestration service

**Recommendation:** Create additional domain services in `/engine/domain/` for:
- `ScoringService` — orchestrates signal scoring
- `TrendAnalysisService` — multi-month analysis
- `EntryValidationService` — data quality checks

**Files to Create:**
- `/engine/domain/scoring_service.go`
- `/engine/domain/trend_service.go`

#### 2. Limited Entity Support in Aggregates
**Gap:** ADR mentions "aggregates are clusters of entities and value objects"

Current state:
- MonthlyRawSignals is a sub-aggregate (good)
- ComputedScores is a value object (good)
- But no Evidence or Action Plan entities within MonthlyEntry

**Recommendation:** If audit/evidence functionality is planned, create:
```go
// Evidence is an entity within MonthlyEntry aggregate
type Evidence struct {
    id        EvidenceID
    entryID   MonthlyEntryID
    text      string
    createdAt time.Time
}
```

#### 3. No Event Sourcing or Audit Events
**Gap:** ADR mentions "domain events (MemberDeactivated, EntryComputed, etc.) could be captured for audit"

Current state:
- Audit logging exists in `/store/audit_log.go` but doesn't model domain events
- No event stream or event-driven domain operations

**Recommendation:** Add domain events:
```go
// /engine/domain/events.go
type DomainEvent interface {
    AggregateID() interface{}
    OccurredAt() time.Time
}

type MemberDeactivated struct {
    MemberID TeamMemberID
    At       time.Time
}
```

#### 4. Weak Aggregate Reconstruction
**Gap:** MonthlyEntry reconstruction in store layer is complex

**File:** `/store/monthly_entry.go:155–202`
- Creates new MonthlyEntry then retrofits timestamps
- Breaks encapsulation with SetCreatedAt/SetComputedAt

**Recommendation:** Add a private factory method in aggregate:
```go
// In domain/aggregates.go (unexported)
func (me *MonthlyEntry) setTimestampsFromStorage(created, computed *time.Time) {
    me.createdAt = created
    me.computedAt = computed
}
```

### Files Implementing DDD Correctly

| File | Purpose | Status |
|------|---------|--------|
| `/engine/domain/value_objects.go` | Value objects with validation | ✓ Complete |
| `/engine/domain/signal_types.go` | Signal type value objects | ✓ Complete |
| `/engine/domain/aggregates.go` | Aggregate roots | ✓ 95% Complete |
| `/engine/domain/repositories.go` | Repository interfaces & in-memory impls | ✓ Complete |
| `/engine/domain/services.go` | Cross-aggregate domain services | ⚠ Partial (only MonthlyEntryService) |
| `/engine/domain/seniority.go` | Seniority value type | ✓ Complete |
| `/engine/domain/scoring.go` | Scoring result types | ✓ Complete |
| `/store/member.go` | SQLite TeamMember repository | ✓ Complete |
| `/store/monthly_entry.go` | SQLite MonthlyEntry repository | ✓ Complete |
| `/service/member/member.go` | Member service with DI | ✓ Complete |
| `/service/monthly/monthly.go` | Monthly service with DI | ✓ Complete |

### DDD Alignment Score: **85/100**

---

## Part 2: ADR-20260913-Go-Fyne Desktop UI

### Status: 70% Aligned ⚠ (Framework in place, features incomplete)

### What Has Been Implemented Correctly

#### 1. Pure Go Single Binary (100% ✓)
- ✓ Fyne v2 dependency properly declared in `go.mod`
- ✓ No Electron, Chromium, or web server bundled
- ✓ Single `go build` produces executable
- ✓ Cross-platform compilation with `GOOS/GOARCH`

**File:** `/go.mod`
```
require (
    fyne.io/fyne/v2 v2.8.1
    modernc.org/sqlite v1.58.0
)
```

#### 2. Fyne Framework Integration (90% ✓)

**Main Entry Point** `/ui/app.go`
```go
func NewMainWindow(app fyne.App, db *sql.DB) fyne.Window {
    window := app.NewWindow("Team Impact Scorecard")
    window.Resize(fyne.NewSize(1200, 800))
    // ...
}
```

✓ Proper Fyne App creation in `main.go`  
✓ Window lifecycle management  
✓ Container and layout composition

**Main Bootstrap** `/main.go:1–55`
```go
import (
    _ "modernc.org/sqlite"
    "fyne.io/fyne/v2/app"
)

func main() {
    // Fyne app initialization
    fyneApp := app.New()
    window := ui.NewMainWindow(fyneApp, db)
    window.ShowAndRun()
}
```

✓ Correct app initialization flow  
✓ Database passed to UI layer (good coupling boundary)

#### 3. Screen Architecture (70% ✓)

**Screen Registry** `/ui/screens/registry.go`
```go
type ScreenRegistry struct {
    memberService *member.Service
}

func (r *ScreenRegistry) SettingsScreen() fyne.CanvasObject {
    return NewSettingsScreen(r.memberService)
}
```

✓ Decouples screen creation from services  
✓ Allows multiple screens (extensible pattern)

**Settings Screen** `/ui/screens/settings.go:1–141`
```go
func NewSettingsScreen(memberService *member.Service) fyne.CanvasObject {
    // Team member list with edit/deactivate
    // Add member dialog
    // Proper container layout
}
```

✓ Basic CRUD operations for members  
✓ Uses Fyne widgets (List, Button, Entry, Dialog)  
✓ Service layer integration

#### 4. No Background Network Calls (100% ✓)
- ✓ No telemetry or tracking
- ✓ No external API calls in UI
- ✓ SQLite is local only

### What Is Missing or Not Aligned

#### 1. Incomplete Screen Implementations (Critical Gap)

**Current State:**
- Only `SettingsScreen` (Screen F) implemented
- ~141 lines of UI code

**Missing Screens (from PRD):**
- **Screen A:** Dashboard/Overview
- **Screen B:** Team member scorecard detail
- **Screen C:** Monthly entry form
- **Screen D:** Trend visualization
- **Screen E:** Evidence/action plan
- **Screen F:** Settings (only one implemented)

**Impact:** Application is non-functional for user workflows

**Recommendation:** Implement remaining screens:
```
/ui/screens/dashboard.go      (Screen A)
/ui/screens/member_detail.go  (Screen B)
/ui/screens/entry_form.go     (Screen C)
/ui/screens/trends.go         (Screen D)
/ui/screens/evidence.go       (Screen E)
/ui/screens/settings.go       (Already done)
```

#### 2. No Custom Widgets or Visualizations

**Gap:** ADR notes "some custom rendering may be needed (e.g., heatmap, sparkline charts)"

Current state:
- Only standard Fyne widgets (List, Button, Entry, Dialog)
- No heatmap visualization
- No sparkline charts
- No custom scorecards

**Missing Components:**
- **Heatmap Widget** — 24-month color grid per member
- **Sparkline Chart** — Trend visualization
- **Score Card Widget** — Display TII, dimension scores, confidence
- **Impact Grid** — Visual signal impact ratings

**Files to Create:**
```
/ui/widgets/heatmap.go
/ui/widgets/sparkline.go
/ui/widgets/scorecard.go
/ui/widgets/impact_grid.go
```

**Example Custom Widget:**
```go
// /ui/widgets/scorecard.go
type ScoreCard struct {
    widget.BaseWidget
    TII          float64
    Dimensions   DimensionScores
    Confidence   float64
}

func (sc *ScoreCard) CreateRenderer() fyne.WidgetRenderer {
    // Custom rendering for score display
}
```

#### 3. No Data Binding/State Management

**Gap:** UI components don't reflect domain model state changes

Current implementation:
```go
// In settings.go - manual refresh
memberList.Refresh()  // Called after mutations
```

**Issue:** 
- No reactive binding to service updates
- No data change propagation
- UI manually manages state

**Recommendation:** Implement observer pattern:
```go
// /ui/common/observable.go
type Observer interface {
    OnMembersChanged(members []domain.TeamMember)
    OnEntryChanged(entry *domain.MonthlyEntry)
}

type MemberService struct {
    // ...
    observers []Observer
}

func (s *MemberService) OnMemberAdded(m *domain.TeamMember) {
    for _, obs := range s.observers {
        obs.OnMembersChanged(/* ... */)
    }
}
```

#### 4. No Navigation Between Screens

**Gap:** ADR implies multi-screen app, but no navigation implemented

Current state:
- Only SettingsScreen shown: `window.SetContent(screenRegistry.SettingsScreen())`
- No way to navigate to other screens
- No main menu or navigation bar

**Recommendation:** Add navigation:
```go
// /ui/navigation.go
type Navigator struct {
    window fyne.Window
    screens map[string]fyne.CanvasObject
    current string
}

func (n *Navigator) Navigate(screenName string) {
    n.window.SetContent(n.screens[screenName])
}
```

#### 5. Missing Input Validation in Forms

**Gap:** UI forms don't validate against domain constraints

Current state:
```go
// In settings.go - forms accept any input
firstNameEntry := widget.NewEntry()
// No max length, no character restrictions, no format validation
```

**Recommendation:** Add input validators:
```go
// /ui/common/validators.go
func ValidateFullName(s string) error {
    if len(strings.TrimSpace(s)) == 0 {
        return fmt.Errorf("name cannot be empty")
    }
    return nil
}

// Wire into Entry widgets
firstNameEntry.OnChanged = func(s string) {
    firstNameEntry.SetValidation(ValidateFullName)
}
```

#### 6. No Error Handling UI

**Gap:** Errors are silently ignored in UI callbacks

Current implementation:
```go
// In settings.go
hbox.Objects[4].(*widget.Button).OnTapped = func() {
    _ = memberService.DeactivateMember(int64(m.ID()))  // Error ignored!
}
```

**Recommendation:** Add error dialogs:
```go
dialog.ShowError(
    fmt.Errorf("failed to deactivate member: %w", err),
    window,
)
```

#### 7. No Keyboard Navigation or Accessibility

**Gap:** No keyboard shortcuts or accessible widgets

**Recommendation:** Add keyboard shortcuts:
```go
// /ui/common/shortcuts.go
window.SetOnTypedKey(func(e *fyne.KeyEvent) {
    if e.Name == fyne.KeyN && e.Modifier&fyne.KeyModifierControl != 0 {
        // Ctrl+N: New member
        showAddDialog(memberService, memberList)
    }
})
```

### Files Implementing Fyne Correctly

| File | Purpose | Status |
|------|---------|--------|
| `/main.go` | Fyne app bootstrap | ✓ Complete |
| `/ui/app.go` | Main window factory | ✓ Complete |
| `/ui/screens/registry.go` | Screen registry | ✓ Complete |
| `/ui/screens/settings.go` | Settings screen (partial) | ⚠ Incomplete |
| `/ui/widgets/` | Custom widgets | ✗ Missing |

### Fyne UI Alignment Score: **70/100**

### Critical Path to Full Implementation
1. **Implement Screen A (Dashboard)** - 2-3 days
2. **Implement Screen B (Member Detail)** - 2-3 days
3. **Implement Screen C (Entry Form)** - 3-4 days
4. **Add heatmap widget** - 2-3 days
5. **Add sparkline widget** - 2 days
6. **Navigation & routing** - 1-2 days
7. **Error handling & validation** - 1-2 days

---

## Part 3: ADR-20260913-Pure Go SQLite Driver

### Status: 95% Aligned ✓ (Correctly Implemented)

### What Has Been Implemented Correctly

#### 1. modernc.org/sqlite Driver Usage (100% ✓)

**Dependency Declaration** `/go.mod:7`
```
require (
    modernc.org/sqlite v1.58.0
)
```

✓ Correct driver imported  
✓ Pure Go implementation (no CGO)  
✓ Version pinned appropriately

**Driver Registration** `/main.go:9`
```go
import _ "modernc.org/sqlite"

// Later...
db, err := sql.Open("sqlite", dbPath)
```

✓ Blank import registers driver  
✓ Correct driver name ("sqlite", not "sqlite3")  
✓ Standard Go `database/sql` interface used

#### 2. Cross-Platform Compatibility (100% ✓)

**File:** `/main.go:42–54`
```go
func getDatabasePath() (string, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }
    appDir := filepath.Join(configDir, "leadpulse")
    if err := os.MkdirAll(appDir, 0700); err != nil {
        return "", err
    }
    return filepath.Join(appDir, "data.db"), nil
}
```

✓ Uses platform-independent `os.UserConfigDir()`  
✓ Works on macOS, Windows, Linux  
✓ Single database file (no per-platform variations needed)

#### 3. No CGO Build Complexity (100% ✓)

**Go Build Simplicity:**
- ✓ `go build` works without cross-compiler setup
- ✓ `GOOS=windows GOARCH=amd64 go build` works natively
- ✓ No CGO_ENABLED=0 needed
- ✓ CI/CD build matrix simplified

### What Is Missing or Not Aligned

#### 1. No PRAGMA Configuration for Performance

**Gap:** SQLite performance not optimized

Current state:
- `/store/schema.go` creates tables but doesn't set PRAGMAs
- No WAL (Write-Ahead Logging) enabled
- No cache size optimization
- No journal mode configuration

**Recommendation:** Add PRAGMA configuration:
```go
// /store/schema.go - add after opening DB
func ConfigureDatabase(db *sql.DB) error {
    pragmas := []string{
        "PRAGMA journal_mode = WAL",           // Enable WAL for better concurrency
        "PRAGMA cache_size = -64000",          // 64MB cache
        "PRAGMA foreign_keys = ON",            // Enforce foreign keys
        "PRAGMA synchronous = NORMAL",         // Balance safety and speed
        "PRAGMA temp_store = MEMORY",          // Temp tables in memory
    }
    
    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return fmt.Errorf("failed to set PRAGMA: %w", err)
        }
    }
    return nil
}
```

**Files to Update:**
- `/store/schema.go` - Add ConfigureDatabase call
- `/main.go` - Call ConfigureDatabase after sql.Open

#### 2. No Connection Pooling Configuration

**Gap:** No explicit connection pool settings

Current state:
```go
db, err := sql.Open("sqlite", dbPath)
// No SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime
```

**Recommendation:** Configure connection pool:
```go
db, err := sql.Open("sqlite", dbPath)
if err != nil {
    return nil, err
}

// SQLite is single-threaded; limit connections
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
```

#### 3. No Backup/Recovery Mechanism

**Gap:** ADR mentions "corruption recovery is limited in v1"

Current state:
- No backup functionality
- No WAL mode (which would enable point-in-time recovery)
- No offline backup tools

**Recommendation (Future):**
```go
// /store/backup.go
func BackupDatabase(db *sql.DB, backupPath string) error {
    // Copy data.db file
    // Or use SQLite VACUUM INTO command (requires modernc.org/sqlite support)
}
```

### SQLite Driver Alignment Score: **95/100**

---

## Part 4: ADR-20260913-SQLite Local Storage

### Status: 80% Aligned ✓ (Schema and persistence working, incomplete queries)

### What Has Been Implemented Correctly

#### 1. Single-File Storage (100% ✓)

**Schema Definition** `/store/schema.go:1–73`
```go
// Single database file with three tables
CREATE TABLE members (...)
CREATE TABLE monthly_entries (...)
CREATE TABLE audit_log (...)
```

✓ All data in single `data.db` file  
✓ Simple backup: copy the file  
✓ Inspectable with `sqlite3` CLI

**File Location:**
```go
// /main.go:42–54
// Stored in user's config directory:
// macOS: ~/Library/Application Support/leadpulse/data.db
// Linux: ~/.config/leadpulse/data.db
// Windows: %APPDATA%/leadpulse/data.db
```

#### 2. Relational Schema (100% ✓)

**Members Table:**
```sql
CREATE TABLE members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    seniority TEXT NOT NULL,
    created_at TEXT NOT NULL,
    deactivated_at TEXT,
    UNIQUE (team_id, first_name, last_name)
)
```

✓ Proper primary key  
✓ Seniority as TEXT (enumerated values)  
✓ Soft delete with deactivated_at

**Monthly Entries Table:**
```sql
CREATE TABLE monthly_entries (
    member_id INTEGER NOT NULL,
    month TEXT NOT NULL,  -- YYYY-MM format
    morale INTEGER NOT NULL,
    billability INTEGER NOT NULL,
    -- ... 8 more signal columns
    computed_scores TEXT,  -- JSON serialized
    created_at TEXT NOT NULL,
    computed_at TEXT,
    PRIMARY KEY (member_id, month),
    FOREIGN KEY (member_id) REFERENCES members(id)
)
```

✓ Composite primary key (member_id, month)  
✓ Foreign key constraint  
✓ Signal columns map to domain aggregate

**Audit Log Table:**
```sql
CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    member_id INTEGER NOT NULL,
    field TEXT NOT NULL,
    old_value TEXT,
    new_value TEXT NOT NULL,
    changed_by TEXT NOT NULL DEFAULT 'system',
    changed_at TEXT NOT NULL,
    FOREIGN KEY (member_id) REFERENCES members(id)
)
```

✓ Audit trail for compliance

#### 3. Repository-Based Persistence (100% ✓)

**TeamMemberRepository Implementation** `/store/member.go`
- ✓ Save: Insert or update members
- ✓ FindByID: Single member retrieval
- ✓ FindActive: Filter by deactivation status
- ✓ Delete: Hard delete (matches interface)

**MonthlyEntryRepository Implementation** `/store/monthly_entry.go`
- ✓ Save: Insert or update entries
- ✓ FindByID: Composite key lookup
- ✓ FindByMember: Range query by member
- ✓ Delete: Hard delete

#### 4. Aggregate Reconstruction (90% ✓)

Both repositories properly reconstruct aggregates from rows:

**Example:** `/store/member.go:114–145`
```go
name, err := domain.NewFullName(firstName, lastName)
sen := domain.Seniority(seniority)
member, err := domain.NewTeamMember(int64(id), name, sen)
// If deactivated, apply deactivation state
if deactivatedAtStr.Valid {
    _ = member.Deactivate()
}
```

✓ Creates aggregates via constructors  
✓ Validates on reconstruction  
✓ Restores state correctly

#### 5. Timestamp Management (90% ✓)

**RFC3339 Format for Storage:**
```go
// /store/member.go:40
member.CreatedAt().UTC().Format(time.RFC3339)

// /store/monthly_entry.go:47
entry.ComputedAt().UTC().Format(time.RFC3339)
```

✓ ISO 8601 format (platform-independent)  
✓ UTC normalization

**Parse Function:**
```go
// /store/monthly_entry.go:323–328
func parseTimestamp(s string) (time.Time, error) {
    return time.Parse(time.RFC3339, s)
}
```

✓ Consistent parsing with error handling

### What Is Missing or Not Aligned

#### 1. No Indexing Strategy

**Gap:** No indexes defined for query performance

Current state:
- Primary and foreign keys auto-indexed
- No secondary indexes for common queries

**Recommendation:** Add indexes:
```go
// /store/schema.go - add after table creation
func CreateIndexes(db *sql.DB) error {
    indexes := []string{
        "CREATE INDEX IF NOT EXISTS idx_members_seniority ON members(seniority)",
        "CREATE INDEX IF NOT EXISTS idx_monthly_entries_month ON monthly_entries(month)",
        "CREATE INDEX IF NOT EXISTS idx_audit_log_member_id ON audit_log(member_id)",
        "CREATE INDEX IF NOT EXISTS idx_audit_log_changed_at ON audit_log(changed_at)",
    }
    for _, idx := range indexes {
        if _, err := db.Exec(idx); err != nil {
            return fmt.Errorf("failed to create index: %w", err)
        }
    }
    return nil
}
```

#### 2. No Query Helper Functions

**Gap:** Complex SQL queries are duplicated in repository implementations

Current state:
- Trend queries not implemented
- Team aggregation queries missing
- Date range queries absent

**Missing Query Methods:**
```go
// /store/monthly_entry.go - add methods
func (r *SQLiteMonthlyEntryRepository) FindByMonthRange(
    memberID domain.TeamMemberID,
    startMonth, endMonth string,
) ([]*domain.MonthlyEntry, error)

func (r *SQLiteMonthlyEntryRepository) FindAllByMonth(month string) ([]*domain.MonthlyEntry, error)

func (r *SQLiteMonthlyEntryRepository) FindByMemberAndMonthRange(
    memberID domain.TeamMemberID,
    startMonth, endMonth string,
) ([]*domain.MonthlyEntry, error)
```

#### 3. No Transaction Support

**Gap:** Multi-step operations aren't atomic

Example from ADR:
- "One entry per member per month" requires transactional check-then-insert

Current state:
- Check and insert are separate calls in `MonthlyEntryService.CreateEntry()`
- Race condition possible in high-concurrency (unlikely for single-user but still a gap)

**Recommendation:** Add transaction support to domain:
```go
// /engine/domain/repositories.go - add transaction interface
type Transaction interface {
    Commit() error
    Rollback() error
}

type RepositoryFactory interface {
    BeginTx() (Transaction, error)
}

// /store/member.go - implement transactions
func (r *SQLiteTeamMemberRepository) BeginTx() (*sql.Tx, error) {
    return r.db.Begin()
}
```

#### 4. No Schema Migration System

**Gap:** ADR notes "Schema migrations must be managed manually (additive-only in v1)"

Current state:
- No migration tracking
- No version table
- AdditiveMigration logic is hardcoded

**Recommendation:** Add migration table:
```go
// /store/migrations.go
type Migration struct {
    Version int
    Name    string
    SQL     string
}

var migrations = []Migration{
    {
        Version: 1,
        Name:    "initial_schema",
        SQL:     "CREATE TABLE ...",
    },
    {
        Version: 2,
        Name:    "add_audit_log",
        SQL:     "CREATE TABLE audit_log ...",
    },
}
```

#### 5. No Data Validation Constraints in Schema

**Gap:** SQLite doesn't enforce some business rules

Example:
- Morale must be 0–5, but schema allows any integer
- Billability must be 0–100, but schema allows any integer

Current approach:
- Relies on domain layer validation
- No database-level constraints

**Note:** This is actually acceptable per ADR approach (domain validates), but could be enhanced:

```sql
-- Could add CHECK constraints
CREATE TABLE monthly_entries (
    ...
    morale INTEGER NOT NULL CHECK (morale >= 0 AND morale <= 5),
    billability INTEGER NOT NULL CHECK (billability >= 0 AND billability <= 100),
    ...
)
```

#### 6. No Computed Scores Sub-tables

**Gap:** Complex scores stored as JSON string

Current implementation:
```go
// /store/monthly_entry.go:32–42
var computedScoresJSON sql.NullString
if entry.HasComputedScores() {
    computed := entry.ComputedScores()
    jsonBytes, err := json.Marshal(computed)
    computedScoresJSON.String = string(jsonBytes)
}
```

**Issue:**
- JSON stored as opaque string
- Can't query within computed_scores
- Can't aggregate dimension scores across months

**Recommendation (Future Enhancement):**
```sql
-- Normalize computed scores
CREATE TABLE computed_scores (
    member_id INTEGER NOT NULL,
    month TEXT NOT NULL,
    dimension_growth REAL NOT NULL,
    dimension_project REAL NOT NULL,
    dimension_team REAL NOT NULL,
    dimension_org REAL NOT NULL,
    tii REAL NOT NULL,
    completeness REAL NOT NULL,
    confidence REAL NOT NULL,
    PRIMARY KEY (member_id, month),
    FOREIGN KEY (member_id) REFERENCES members(id)
)
```

This would enable:
- Direct SQL queries for trends
- Aggregations without application-level code
- Better normalization

#### 7. Limited Query Interface

**Gap:** Repository doesn't expose query building

Current state:
- Predefined FindByID, FindByMember, etc.
- Can't construct arbitrary queries
- Trend queries require custom methods

**Recommendation (if needed later):**
```go
// Optional: Query builder interface for advanced use cases
type QueryBuilder interface {
    Where(field string, op string, value interface{}) QueryBuilder
    OrderBy(field string, desc bool) QueryBuilder
    Limit(n int) QueryBuilder
    Execute(ctx context.Context) ([]*domain.MonthlyEntry, error)
}
```

### Files Implementing SQLite Storage Correctly

| File | Purpose | Status |
|------|---------|--------|
| `/store/schema.go` | Schema definition | ✓ Complete |
| `/store/member.go` | TeamMember repository | ✓ Complete |
| `/store/monthly_entry.go` | MonthlyEntry repository | ✓ Complete |
| `/store/audit_log.go` | Audit trail | ✓ Partial |
| `/main.go` | DB initialization | ✓ Complete |

### SQLite Local Storage Alignment Score: **80/100**

---

## Cross-Layer Analysis

### Layer Separation Quality

```
┌─────────────────────────────────────────┐
│           UI Layer (Fyne)               │
│  /ui/app.go, /ui/screens/              │
│  Status: 70% (incomplete UI)            │
└────────────────┬────────────────────────┘
                 │ (services via DI)
┌────────────────▼────────────────────────┐
│       Application Service Layer         │
│  /service/member/, /service/monthly/    │
│  Status: 90% (good DI, needs tests)     │
└────────────────┬────────────────────────┘
                 │ (repositories via DI)
┌────────────────▼────────────────────────┐
│        Domain Layer (DDD)               │
│  /engine/domain/                        │
│  Status: 85% (strong value objects,     │
│  aggregates, services - missing events) │
└────────────┬──────────────┬─────────────┘
             │              │
        (implements)   (implements)
             │              │
┌────────────▼──┐   ┌───────▼──────────────┐
│In-Memory Repo │   │SQLite Repository    │
│/domain/       │   │/store/              │
│Status: 100%   │   │Status: 90%          │
└───────────────┘   └─────────────────────┘
```

### Dependency Flow Analysis

**Correct Dependencies (Following ADR):**
- ✓ UI → Application Service Layer
- ✓ Application Service Layer → Domain Layer (via interfaces)
- ✓ Domain Layer → Repository Interfaces
- ✓ Concrete Repositories → Domain Layer (aggregate reconstruction)

**Potential Issues:**
1. ⚠ UI has direct `sql.DB` dependency (minor coupling)
   - **Location:** `/ui/app.go:4,13`
   - **Recommendation:** Hide DB behind repository factory

2. ⚠ Service layer creates domain services on-demand
   - **Location:** `/service/monthly/monthly.go:51`
   - **Recommendation:** Inject domain service into service constructor

### Testing Coverage

**Current Test Files:**
- `/store/member_test.go` - Repository tests
- `/store/monthly_entry_test.go` - Repository tests
- `/service/member/member_test.go` - Service tests
- `/service/monthly/monthly_test.go` - Service tests
- `/engine/scoring/*_test.go` - Scoring tests

**Total:** 9 test files

**Gaps:**
- No domain aggregate tests
- No domain service tests
- No value object tests
- No repository interface tests

**Recommendation:** Add tests:
```
/engine/domain/aggregates_test.go
/engine/domain/services_test.go
/engine/domain/value_objects_test.go
```

---

## Summary Table: ADR Compliance

| ADR | Component | Status | Score | Key Issues |
|-----|-----------|--------|-------|-----------|
| **DDD** | Value Objects | ✓ Complete | 100% | None |
| **DDD** | Aggregates | ✓ Complete | 95% | Minor encapsulation via SetCreatedAt/SetComputedAt |
| **DDD** | Repository Interfaces | ✓ Complete | 100% | None |
| **DDD** | Domain Services | ⚠ Partial | 90% | Missing scoring, trend, validation services |
| **DDD** | Domain Events | ✗ Missing | 0% | Not implemented |
| **DDD** | Service Layer DI | ✓ Complete | 100% | None |
| **Fyne** | Framework Integration | ✓ Complete | 100% | None |
| **Fyne** | Screen Architecture | ⚠ Partial | 70% | Only Settings screen; missing 5 screens |
| **Fyne** | Custom Widgets | ✗ Missing | 0% | No heatmap, sparkline, scorecard widgets |
| **Fyne** | Navigation | ✗ Missing | 0% | No multi-screen routing |
| **SQLite Driver** | Driver Selection | ✓ Complete | 100% | None |
| **SQLite Driver** | Cross-Platform | ✓ Complete | 100% | None |
| **SQLite Storage** | Schema Design | ✓ Complete | 95% | Consider CHECK constraints |
| **SQLite Storage** | Repository Pattern | ✓ Complete | 95% | Missing transaction support |
| **SQLite Storage** | Data Queries | ⚠ Partial | 70% | Missing indexes, trend queries |

### Overall ADR Compliance: **82/100**

---

## Recommendations by Priority

### PHASE 1: Critical (Must Fix Before Release)

1. **Implement Missing Screens** (UI Layer)
   - Screen A: Dashboard
   - Screen B: Member detail
   - Screen C: Entry form
   - Effort: 7-10 days
   - Impact: Application becomes functional

2. **Add Custom Widgets** (UI Layer)
   - Heatmap for monthly data
   - Sparkline for trends
   - Score card display
   - Effort: 5-7 days
   - Impact: User experience improvement

3. **Add Navigation System** (UI Layer)
   - Multi-screen routing
   - Keyboard shortcuts
   - Back button handling
   - Effort: 2-3 days
   - Impact: User can navigate app

4. **Add Error Handling UI** (UI Layer)
   - Error dialogs
   - Form validation
   - User feedback
   - Effort: 2-3 days
   - Impact: User aware of problems

### PHASE 2: Important (Improve Architecture)

5. **Create Additional Domain Services** (Domain Layer)
   - ScoringService for pipeline orchestration
   - TrendAnalysisService for multi-month analysis
   - EntryValidationService for data quality
   - Effort: 3-4 days
   - Impact: Cleaner domain model

6. **Add Database PRAGMA Configuration** (Store Layer)
   - WAL mode for concurrency
   - Cache optimization
   - Foreign key enforcement
   - Effort: 1 day
   - Impact: Better performance

7. **Add Database Indexes** (Store Layer)
   - Index common queries
   - Improve query performance
   - Effort: 1-2 days
   - Impact: Faster data retrieval

8. **Add Domain Tests** (Testing)
   - Aggregate tests
   - Value object tests
   - Domain service tests
   - Effort: 3-4 days
   - Impact: Higher confidence in domain logic

### PHASE 3: Nice to Have (Future Enhancement)

9. **Implement Domain Events** (Domain Layer)
   - Event emission from aggregates
   - Event sourcing preparation
   - Audit trail enhancement
   - Effort: 3-5 days
   - Impact: Future extensibility

10. **Add Transaction Support** (Store Layer)
    - Atomic multi-step operations
    - Distributed transaction safety
    - Effort: 2-3 days
    - Impact: Data integrity

11. **Implement Schema Migrations** (Store Layer)
    - Version tracking
    - Automated migrations
    - Rollback capability
    - Effort: 2-3 days
    - Impact: Future schema evolution

12. **Normalize Computed Scores** (Store Layer)
    - Create scores sub-table
    - Enable trend queries in SQL
    - Effort: 2-3 days
    - Impact: Better query performance

---

## Code Quality Observations

### Strengths
1. **Clean Separation of Concerns** - Each layer has clear responsibility
2. **Strong Type System** - Value objects prevent invalid states
3. **Good Encapsulation** - Aggregates properly hide internals
4. **Proper Dependency Injection** - Services depend on interfaces, not implementations
5. **Error Handling** - Most functions properly return errors
6. **Code Organization** - Clear directory structure matching ADRs

### Areas for Improvement
1. **Test Coverage** - Missing tests for domain layer
2. **Error Context** - Some errors lack context information
3. **Documentation** - No architecture documentation beyond ADRs
4. **Logging** - No structured logging implementation
5. **UI Polish** - UI is minimal, needs refinement
6. **Configuration** - No config file support
7. **Validation** - Some validation duplicated between domain and service layers

---

## Files Requiring Attention

### High Priority Changes

| File | Change Type | Reason |
|------|-------------|--------|
| `/ui/screens/settings.go` | Enhancement | Add error dialog handling |
| `/ui/app.go` | Refactor | Remove sql.DB dependency |
| `/engine/domain/services.go` | Enhancement | Add more domain services |
| `/store/schema.go` | Enhancement | Add PRAGMA configuration |

### New Files Needed

| File | Purpose | Priority |
|------|---------|----------|
| `/ui/screens/dashboard.go` | Screen A | Critical |
| `/ui/screens/member_detail.go` | Screen B | Critical |
| `/ui/screens/entry_form.go` | Screen C | Critical |
| `/ui/screens/trends.go` | Screen D | Critical |
| `/ui/screens/evidence.go` | Screen E | Critical |
| `/ui/widgets/heatmap.go` | Custom widget | Critical |
| `/ui/widgets/sparkline.go` | Custom widget | Critical |
| `/ui/widgets/scorecard.go` | Custom widget | Critical |
| `/ui/navigation.go` | Navigation system | Critical |
| `/engine/domain/scoring_service.go` | Domain service | Important |
| `/engine/domain/events.go` | Domain events | Nice to Have |
| `/store/migrations.go` | Schema versioning | Nice to Have |

---

## Conclusion

The LeadPulse codebase demonstrates **strong architectural foundation** with good implementation of:
- Domain-Driven Design patterns (85% complete)
- Go + Fyne desktop framework (70% complete - framework present, UI incomplete)
- Pure Go SQLite driver (95% correctly implemented)
- SQLite local storage (80% correctly implemented)

**The main gap is in the UI layer**, where only 1 of 6 screens is implemented and custom widgets are missing. The domain and storage layers are well-architected and could support full UI implementation without major changes.

**Estimated effort to full compliance:**
- UI screens and widgets: 15-20 days
- Domain services and tests: 10-12 days
- Database optimization: 5-7 days
- **Total: 30-40 days for senior developer**

The ADR decisions are sound and implementation is progressing well. With focused effort on the UI layer and some domain enhancements, the application can achieve 95%+ ADR alignment.
