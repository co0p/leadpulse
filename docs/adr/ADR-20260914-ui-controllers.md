# ADR-20260914 — Extract Screen Logic into Testable Controllers

**Decision:** Create a `ui/controllers/` package to extract business and state logic from Fyne screens into Fyne-independent controller objects.

**Status:** Accepted

---

## Context

The current screen implementations (particularly `monthly_input.go`) are monolithic functions with all state management embedded in widget callbacks and closures. This creates three problems:

1. **Untestability**: Business logic (loading data, validating input, saving) cannot be unit tested without spinning up a Fyne app and simulating UI interactions.
2. **Debugging difficulty**: State changes are scattered across callback closures, making it hard to trace where data corruption happens (e.g., pointer aliasing bugs on line 322).
3. **Scalability**: Each new screen doubles the cognitive load — all logic must be held in mind at once to avoid state bugs.

Per CONSTITUTION.md: "Service is the API contract. The UI is replaceable." Currently, Fyne widgets own the state; the UI is not replaceable and not testable.

---

## Alternatives

**1. Status quo — keep monolithic screens**
- Cost: Continue as-is with untestable logic, hard-to-fix bugs.
- Benefit: No refactor effort.

**2. Extract business logic to `service/` layer**
- Cost: Violates CONSTITUTION ("service" is use cases, not UI-specific logic like form validation or member selection).
- Benefit: Testable.

**3. Extract to `ui/controllers/` (chosen)**
- Cost: New package, ~200 LOC per screen controller initially, refactor existing screens.
- Benefit: Controllers are pure Go, zero Fyne dependencies, testable. UI layer becomes thin.

**4. Global state manager (Redux-like)**
- Cost: Higher complexity, overkill for current screen count.
- Benefit: Easier to coordinate multi-screen state later.

---

## Rationale

Option 3 is the best fit:

- **Testability without Fyne**: Controllers are pure Go with zero Fyne imports. Tests run in Go `_test.go` files without app setup.
- **Clear separation**: Controller owns state and business logic; Fyne layer renders state and dispatches events to controller.
- **Aligned with CONSTITUTION**: Thin UI layer, replaceable screens, service calls from controller (not from UI).
- **Proven pattern**: MVC/MVP/MVVM in desktop and web frameworks use this exact separation.
- **Incremental refactor**: Each screen can be refactored independently.
- **Natural package location**: `ui/controllers/` is a Fyne-specific bridge between `service/` and Fyne widgets.

---

## Consequences

**Better:**
- Business logic testable in isolation (no Fyne, no widget setup).
- State mutations concentrated in one place (controller) — bugs easier to find.
- Controllers are reusable: same controller can drive a CLI, web UI, or desktop UI.
- Clearer dependencies: UI → Controller → Service → Store/Engine.
- New team members can test controllers without learning Fyne.

**Harder:**
- Each screen now requires two files (controller + UI).
- Initial refactor of existing screens (2–3 hours total).
- Developers must discipline themselves to not put state in widgets.

---

## Implementation Pattern

### Controller Structure

```go
// ui/controllers/monthly_input_controller.go
type MonthlyInputController struct {
    memberService  *member.Service
    monthlyService *monthly.Service
    
    // State
    currentMember   *domain.TeamMember
    members         []*domain.TeamMember
    formState       MonthlyInputState
}

type MonthlyInputState struct {
    Morale             *int
    Billability        *int
    CSAT               *int
    // ... other fields
}

// Public API for UI layer to call
func (c *MonthlyInputController) LoadMembers() error
func (c *MonthlyInputController) SelectMember(memberID int) error
func (c *MonthlyInputController) SetFormField(fieldName string, value *int) error
func (c *MonthlyInputController) SaveMember() error
func (c *MonthlyInputController) GetFormState() MonthlyInputState
```

### UI Layer (thin wrapper)

```go
// ui/screens/monthly_input.go
type MonthlyInputScreen struct {
    *container.Split
    controller *controllers.MonthlyInputController
    fields     map[string]*widget.Entry
}

// When user types: update controller, controller updates state
func (s *MonthlyInputScreen) onMoraleChanged(text string) {
    val := parseInt(text)
    s.controller.SetFormField("morale", val)
    s.refreshUI()
}

// When user clicks Save: ask controller to save
func (s *MonthlyInputScreen) onSaveClicked() {
    if err := s.controller.SaveMember(); err != nil {
        s.showError(err)
        return
    }
    s.refreshUI()
}

// Render from controller state
func (s *MonthlyInputScreen) refreshUI() {
    state := s.controller.GetFormState()
    if state.Morale != nil {
        s.fields["morale"].SetText(strconv.Itoa(*state.Morale))
    }
    // ... update other widgets
}
```

### Test Example

```go
// ui/controllers/monthly_input_controller_test.go
func TestMonthlyInputController_SaveMember(t *testing.T) {
    memberRepo := domain.NewInMemoryTeamMemberRepository()
    entryRepo := domain.NewInMemoryMonthlyEntryRepository()
    
    memberSvc := member.NewService(memberRepo, entryRepo)
    monthlySvc := monthly.NewService(memberRepo, entryRepo)
    
    // Create controller (zero Fyne dependencies)
    controller := NewMonthlyInputController(memberSvc, monthlySvc)
    
    // Test business logic
    controller.LoadMembers()
    controller.SelectMember(0)
    controller.SetFormField("morale", intPtr(4))
    controller.SetFormField("billability", intPtr(85))
    
    err := controller.SaveMember()
    if err != nil {
        t.Fatalf("Save failed: %v", err)
    }
    
    // Verify persistence
    state := controller.GetFormState()
    if *state.Morale != 4 {
        t.Errorf("Morale not saved: got %v", state.Morale)
    }
}
```

---

## Related

- [ADR-20260913-go-fyne-desktop.md](ADR-20260913-go-fyne-desktop.md) — Fyne is the UI framework
- CONSTITUTION.md, Section "Architecture Boundaries" — UI layer responsibilities

---

## Next Steps

1. Create `ui/controllers/` package with base controller interface.
2. Extract `MonthlyInputController` from `monthly_input.go` (2 hours).
3. Extract `SettingsController` from `settings.go` (1 hour).
4. Refactor screens to use controllers (1 hour).
5. Add controller unit tests (1 hour).
6. Update `docs/ui.md` with controller pattern documentation.
