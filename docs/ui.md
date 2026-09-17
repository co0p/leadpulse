# UI Design System

Shared decisions for the web application shell and future screens. Use these patterns and constraints as the foundation for all new features. The frontend is a Vue 3 SPA (see `docs/adr/ADR-20260917-vue-spa-frontend.md`); the visual and interaction decisions below apply regardless of the templating technology used to implement them.

---

## Shell Layout

**Standard 3-region layout:**
- `<header role="banner">` — fixed-height top bar (56px sticky)
- `<aside role="navigation">` — sidebar for primary navigation
- `<main role="main">` — scrollable content area

**Full-height viewport:** All containers together fill `height: 100vh` with no horizontal scroll at any breakpoint.

**Flexbox-based structure:**
```
#app (display: flex, flex-direction: column, height: 100vh)
  └─ .app-shell (display: flex, flex-direction: column, height: 100vh)
       ├─ header.app-topbar (height: 56px, flex-shrink: 0)
       ├─ .app-container (display: flex, flex: 1, overflow: hidden, position: relative)
       │    ├─ aside.app-sidebar (width: 260px desktop, overlay mobile; flex-shrink: 0)
       │    ├─ .app-sidebar-overlay (v-if; visible on mobile when sidebar open)
       │    └─ main.app-main (flex: 1, overflow-y: auto)
       └─ footer.app-footer (flex-shrink: 0)
```

---

## Responsive Breakpoints

**Three breakpoint tiers, mobile-first:**

| Breakpoint | Width | Sidebar | Toggle | Layout |
|------------|-------|---------|--------|--------|
| Desktop | ≥1024px | Fixed 260px, visible | Hidden | Horizontal split (sidebar + content) |
| Tablet | 769–1023px | Overlay, hidden by default | Visible | Full-width content with toggle to show sidebar |
| Mobile | ≤768px | Overlay, hidden by default | Visible | Full-width content with compact spacing |

**CSS Media Queries:**
```css
/* Mobile first: overlay sidebar by default */
@media (max-width: 1023px) {
  .app-sidebar { position: absolute; left: 0; top: 0; bottom: 0; transform: translateX(-100%); }
  .app-sidebar.app-sidebar--open { transform: translateX(0); }
}

/* Tablet/mobile: tighter spacing */
@media (max-width: 767px) {
  .app-topbar { padding: 0 0.75rem; }
  .app-topbar-center .field { max-width: 200px; }
  .app-main { padding: 0.75rem; }
}
```

---

## Top Bar

**Sticky positioning:** Always visible at top; 56px fixed height.

**Three zones:**

| Zone | Content | Responsibility |
|------|---------|-----------------|
| Left (200px min-width) | Sidebar toggle button (mobile/tablet), breadcrumb or context label | Navigation |
| Center (flex: 1) | Centered search input (max-width: 400px desktop, 200px mobile) | Search/discovery |
| Right (120px min-width) | Quick action buttons (notifications, help) | Primary actions |

**Search input:** Centered, flexes to fill available space up to max-width. On mobile, width reduced and label hidden (icon-only).

**Sidebar toggle button:**
- Always 44px × 44px minimum touch target
- Visible only on mobile/tablet (<1024px)
- ARIA: `aria-label="Toggle sidebar"`, `aria-expanded` dynamic
- Focus state: 2px solid outline, offset 2px

---

## Sidebar

**Desktop (≥1024px):** Fixed-width 260px, always visible to the left of content.

**Mobile/Tablet (<1024px):** Overlay drawer, positioned absolutely, hidden by default. On toggle, slides in from left with semi-transparent overlay behind it.

**Sidebar structure:**
- **Header:** Logo/branding (Leadpulse), padding 1.5rem
- **Nav section:** Primary links (Home, Settings, Alerts, Reports), scrollable, flex: 1
- **Collapsible section:** Secondary navigation (e.g., "Tools"), labeled with uppercase menu-label
- **Footer:** Primary action button (e.g., "+ Add Item"), padding 1rem, sticky at bottom

**Colors & spacing:**
- Background: `var(--bulma-scheme-main-bis)` (≈ #f5f5f5)
- Border: `var(--bulma-border)` (≈ #dbdbdb)
- Link padding: 0.75rem 1rem (provided by Bulma `menu-list`)
- Link hover: provided by Bulma `menu-list a:hover`
- Menu label: Bulma `menu-label` class (0.75em, uppercase, `var(--bulma-text-weak)`)

---

## Main Content Area

**Scrollable:** `overflow-y: auto` when content exceeds viewport height.

**Padding:** 1.25rem desktop, 0.75rem mobile.

**Background:** `var(--bulma-scheme-main-ter)` (≈ #fafafa, light gray to distinguish from white shell).

**Sections:** Use Bulma `.section` and `.container` for consistent spacing and max-width management.

---

## Accessibility

**Target standard:** WCAG 2.1 Level AA.

**Semantic HTML:**
- Top bar: `<header role="banner">`
- Sidebar: `<aside role="navigation">` with `<nav>` for link groups
- Content: `<main role="main">`
- Sections: `<section>` for content groupings

**ARIA attributes:**
- Sidebar toggle: `aria-label="Toggle sidebar"`, `aria-expanded` reflects state (true/false)
- Search input: `aria-label="Search members and content"` or visible `<label>`
- Quick action buttons: `aria-label` for icon-only buttons (e.g., "Notifications", "Help")
- Sidebar overlay: `role="presentation"` (non-semantic)

**Keyboard navigation:**
- Tab order flows: toggle → search → quick actions → sidebar links → content
- Esc key closes sidebar on mobile/tablet
- All interactive elements (buttons, links, inputs) receive focus outline on keyboard focus (`:focus-visible`): `2px solid var(--bulma-link)`, offset 2px. Bulma handles this natively on its own form and button elements; custom elements should match the same pattern.
- No negative `tabindex` (all elements naturally focusable)

**Focus states:**
- Buttons: Bulma handles `:focus-visible` natively. Custom buttons: `outline: 2px solid var(--bulma-link); outline-offset: 2px;` on `:focus-visible` only (not `:focus`, to avoid outlines on mouse click).
- Links: same outline pattern on `:focus-visible` + hover background color change
- Inputs: Bulma handles `:focus-visible` natively with border-color change. Custom inputs follow the same rule.
- High contrast: `var(--bulma-link)` (≈ #3273dc) on white background meets AA requirements.

**Color contrast:** All text on backgrounds meets 4.5:1 ratio for normal text, 3:1 for large text (≥18pt or ≥14pt bold).

---

## CSS Framework & Dependencies

**Bulma CSS:** Base responsive framework, bundled with the Vue build and embedded in the binary via the compiled SPA's static assets.
- Used for `button`, `input`, `field`, `control`, `icon`, `icon-text`, `menu`, `menu-list`, `level`, `level-left`, `level-right`, `level-item`, `is-size-*`, `has-text-*`, `has-text-weight-*`, `section`, `container`, `box`, `title` classes
- Media query breakpoints: desktop (≥1024px), tablet (769–1023px), mobile (≤768px)
- Applied as class names directly in Vue single-file component templates
- **Bulma v1 CSS custom properties** — use `var(--bulma-*)` tokens in scoped component CSS instead of hardcoded hex values. Key tokens: `--bulma-scheme-main` (white), `--bulma-scheme-main-bis` (near-white), `--bulma-scheme-main-ter` (light gray), `--bulma-border`, `--bulma-text`, `--bulma-text-weak`, `--bulma-link`. Hardcoded hex values will not respond to theme changes.
- **`menu` component structure** — sidebar navigation must follow the Bulma `menu` pattern exactly: `<aside class="menu"> > <ul class="menu-list">`. Do not apply `menu-list` without the parent `menu` element; Bulma's hover, color, and spacing styles depend on the full structure.
- **`level` component** — use `<nav class="level">` with `level-left`/`level-right`/`level-item` for horizontal layouts that need space-between distribution (e.g., footer, toolbar rows). Do not write bespoke `justify-content: space-between` flexbox when `level` applies.
- **`is-active` modifier** — reserved for Bulma's own semantic active-state (e.g., `menu-list a.is-active` for the current route). Do not use `is-active` to drive show/hide toggling of custom elements; use a descriptive BEM modifier instead (e.g., `app-sidebar--open`).

**Vue Router:** Client-side routing and shell/content composition. A persistent `AppShell.vue` layout component wraps routed screens, guaranteeing the shell is always present around feature content.

**Pinia:** Client-side state management for cross-screen state (e.g., sidebar open/closed, active member selection, active cycle).

**No external CDN calls:** All frontend assets (JS, CSS) are bundled by the Vite build and embedded in the binary. Fonts via Font Awesome CDN only (deferred to v2 for offline-capable font embedding).

---

## Interaction Patterns

**Sidebar toggle (mobile/tablet):**
1. User clicks toggle button
2. Sidebar open/closed state lives in local component state (`sidebarOpen` ref), toggled on click
3. Sidebar modifier class bound reactively (e.g., `:class="{ 'app-sidebar--open': sidebarOpen }"`) triggers CSS transform
4. Overlay rendered conditionally (`v-if="sidebarOpen"`) dims background
5. Clicking overlay or pressing Esc closes sidebar

**Search input:**
- Accepts focus and typing without errors
- Focus state tracked in local component state (for future styling, if needed)
- No current API wired (ready for future implementation)

**Keyboard shortcuts:**
- Esc: closes sidebar (mobile/tablet only)
- Tab: navigates through interactive elements in logical order
- Enter: submits search or activates focused button (default HTML behavior)

---

## Color Palette

**Functional colors:**
- Primary action: `var(--bulma-link)` (≈ #3273dc) — used for active states, focus outlines, primary buttons
- Background: `var(--bulma-scheme-main-ter)` (≈ #fafafa, light gray)
- Borders: `var(--bulma-border)` (≈ #dbdbdb)
- Text: `var(--bulma-text)` (standard), `var(--bulma-text-weak)` (secondary labels)
- Shell background (topbar, footer): `var(--bulma-scheme-main)` (white)
- Sidebar background: `var(--bulma-scheme-main-bis)` (near-white / light gray)

Do not use hardcoded hex values for colors that map to Bulma tokens. Reserve raw hex only for values with no Bulma equivalent (e.g., the overlay `rgba(0, 0, 0, 0.5)`).

**Overlay:** rgba(0, 0, 0, 0.5) — semi-transparent black for sidebar overlay

---

## Typography

**Headings:**
- H1: 2rem, font-weight 600 (page title)
- H2: 1.25rem, font-weight 600 (sidebar logo)
- H3: 1.125rem, font-weight 600 (section headings)

**Body text:** 1rem, font-weight 400, line-height 1.5

**Menu labels:** 0.85rem, uppercase, font-weight 600, letter-spacing 0.5px (secondary navigation groupings)

**Font:** System default (Bulma uses `sans-serif` stack)

---

## Touch Targets

**Minimum touch target size:** 44px × 44px (recommended by WCAG, iOS HIG, Android Material).

- Sidebar toggle: 44px button
- Navigation links: min-height 44px (padding 0.75rem 1rem = ~40px height + border)
- Quick action buttons: 44px min-width
- Search input: min-height 44px
- Close/confirm buttons: 44px × 44px minimum

---

## Future Screens

All new feature screens must:
1. Use the 3-region layout (header + sidebar + main content) via the shared `AppShell.vue` layout and Vue Router nested routes
2. Respect responsive breakpoints (desktop/tablet/mobile)
3. Follow accessibility target (WCAG 2.1 AA)
4. Use Bulma CSS classes and `var(--bulma-*)` CSS custom properties — no hardcoded hex values for colors that map to Bulma tokens
5. Follow Bulma component structure rules: `menu > menu-list` for navigation, `level` for horizontal space-between layouts, `is-active` only for Bulma's semantic active-state
6. Use `:focus-visible` (not `:focus`) for custom focus ring styles
7. Test focus states, keyboard navigation, and screen reader compatibility (component tests plus manual browser verification)

---

## Update Policy

Update this document when:
- A new responsive breakpoint or layout pattern is introduced
- Accessibility standard changes (e.g., upgrade to AAA)
- CSS framework or client-side library is replaced
- New interaction patterns or focus state styles are added

Do NOT update for:
- Individual screen designs or content (covered by separate feature PRD)
- Specific colors for individual components (use the palette, do not invent new colors)
- Implementation details of individual features
