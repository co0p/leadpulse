# Fyne UI Spike — Findings

Learnings from a throwaway Fyne v2.8.1 prototype covering all six PRD screens with stubbed data. Reference this when planning or implementing any `ui/` work.

---

## What Was Verified

All six screens render correctly in Fyne v2 at 1200×800 on macOS arm64. No architecture change is required.

| Layout concern | Verified approach |
|---|---|
| Sidebar + main area | `container.NewHSplit` with `Offset = 0.18` |
| Top header bar | `container.NewBorder` — content right, nil other sides |
| Screen navigation | String-keyed map of constructor functions + `container.NewStack` content swap — instant, no flicker |
| Heatmap grid (Screen A) | `container.NewGridWithColumns(N)` with `canvas.NewRectangle` heat-colored cells |
| Split-pane input (Screen B) | `container.NewHSplit` — member list left, form right |
| Live preview sticky panel (Screen B) | `container.NewVBox` card at the bottom of the right pane |
| Alert triage rows (Screen D) | `container.NewHBox` — pill + label + `layout.NewSpacer()` + buttons |
| Decision support cards (Screen E) | `container.NewBorder` with a narrow `canvas.NewRectangle` left-border strip |
| Evidence timeline (Screen C) | `container.NewVBox` of cards — no native timeline widget needed |
| Settings form (Screen F) | `widget.NewForm` + `widget.NewFormItem` — correct native widget, no custom work |

---

## Shared Widget Patterns

These patterns are needed on every screen. Extract them into `ui/widgets/` before building any screen.

**Card**
```go
bg := canvas.NewRectangle(colCard)   // white
bg.CornerRadius = 6
bg.StrokeColor = colBorder
bg.StrokeWidth = 1
container.NewStack(bg, container.NewPadded(content))
```

**Pill** (severity, confidence, dimension tags)
```go
bg := canvas.NewRectangle(bgColor)
bg.CornerRadius = 8
txt := canvas.NewText(label, fgColor)
txt.TextSize = 11
txt.TextStyle = fyne.TextStyle{Bold: true}
container.NewStack(bg, container.NewPadded(txt))
```

**HeatCell** (heatmap grid cell)
```go
bg := canvas.NewRectangle(heatColor(value))  // green/amber/red bg
txt := canvas.NewText(fmt.Sprintf("%.1f", value), theme.ForegroundColor())
container.NewStack(bg, container.NewCenter(txt))
```

**SectionTitle**
```go
txt := canvas.NewText(label, theme.ForegroundColor())
txt.TextSize = 13
txt.TextStyle = fyne.TextStyle{Bold: true}
```

---

## Color System

Manage all colors explicitly. Do not rely on the Fyne theme for card backgrounds or severity colors — the default theme is dark-biased and will produce incorrect results.

```go
// Backgrounds
colCard   = color.NRGBA{R: 255, G: 255, B: 255, A: 255}  // white card
colBg     = color.NRGBA{R: 248, G: 249, B: 250, A: 255}  // page background
colBorder = color.NRGBA{R: 222, G: 226, B: 230, A: 255}  // card border
colMuted  = color.NRGBA{R: 108, G: 117, B: 125, A: 255}  // muted labels
colSidebar = color.NRGBA{R: 241, G: 243, B: 245, A: 255} // sidebar bg

// Severity foregrounds
colRed   = color.NRGBA{R: 220, G: 53,  B: 69,  A: 255}
colAmber = color.NRGBA{R: 255, G: 193, B: 7,   A: 255}
colGreen = color.NRGBA{R: 25,  G: 135, B: 84,  A: 255}

// Severity backgrounds (pill / heat cell)
colRedBg   = color.NRGBA{R: 248, G: 215, B: 218, A: 255}
colAmberBg = color.NRGBA{R: 255, G: 243, B: 205, A: 255}
colGreenBg = color.NRGBA{R: 209, G: 231, B: 221, A: 255}
```

Heat cell thresholds (from PRD):
- `>= 75` → green background
- `50–74` → amber background
- `< 50` → red background

---

## Gaps and How to Handle Them

**No chart/sparkline in Fyne stdlib**
The 12-month TII trend chart (Screen C) has no native widget. Use `github.com/wcharczuk/go-chart/v2` to render to `image.Image`, then display via `canvas.NewImageFromImage` in a fixed-size container. This is the standard Fyne pattern for custom charts.

**No native pill widget**
Use the `canvas.NewRectangle` + `CornerRadius` + `canvas.NewText` stack pattern above. Extract to `ui/widgets.Pill(label, bg, fg)`.

**No native timeline/feed widget**
A `container.NewVBox` of `Card` objects is sufficient. No custom renderer needed for v1.

**Heatmap minimum sizes**
Fyne grid columns do not enforce minimum widths. At 1200×800 (the planned window default) this is fine. At smaller sizes the heatmap compresses. Set `w.Resize(fyne.NewSize(1200, 800))` as the launch default.

**Progress bar for completeness**
`widget.NewProgressBar` exists and works. No custom work needed.

---

## Theme Setup

Use the Fyne light theme as the base. Apply it at startup before any window is created.

```go
a := app.New()
a.Settings().SetTheme(theme.LightTheme())
w := a.NewWindow("Team Impact Scorecard")
w.Resize(fyne.NewSize(1200, 800))
```

Supplement the theme with explicit background rectangles for cards and the sidebar. Do not override theme colors globally — it causes unpredictable results on other widgets (buttons, inputs, selects).

---

## Screen Navigation

No router library needed. A map of constructor functions and a `container.NewStack` is sufficient.

```go
screens := map[string]func() fyne.CanvasObject{
    "Overview":      screenOverview,
    "Monthly Input": screenInput,
    // ...
}
content := container.NewStack()

setScreen := func(name string) {
    content.Objects = []fyne.CanvasObject{screens[name]()}
    content.Refresh()
}
```

Each screen constructor returns a fresh `fyne.CanvasObject`. State is not preserved between navigations in this model — acceptable for v1 since forms save to the store on submit.

---

## Build Notes

- Fyne v2.8.1 resolves cleanly with `go mod tidy`.
- `go build ./...` produces one harmless macOS linker warning: `ld: warning: ignoring duplicate libraries: '-lobjc'`. Not an error; ignore it.
- No CGO required for Fyne itself on macOS (uses Metal/GLFW via Go bindings).
- Cross-compilation for Windows/Linux from macOS requires `fyne package` or a CI matrix — `go build` alone is insufficient for packaged `.app`/`.exe` bundles, but the binary itself cross-compiles normally.
