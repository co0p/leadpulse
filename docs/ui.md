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
body (display: flex, flex-direction: column, height: 100vh)
  ├─ header.topbar (height: 56px, flex-shrink: 0)
  └─ .app-container (display: flex, flex: 1, overflow: hidden)
     ├─ aside.sidebar (width: 260px desktop, overlay mobile; flex-shrink: 0)
     ├─ .sidebar-overlay (display: none; visible on mobile when sidebar open)
     └─ main.main-content (flex: 1, overflow-y: auto)
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
  .sidebar { position: absolute; left: 0; top: 56px; transform: translateX(-100%); }
  .sidebar.is-active { transform: translateX(0); }
}

/* Tablet/mobile: tighter spacing */
@media (max-width: 767px) {
  .topbar { padding: 0 0.75rem; }
  .topbar-center .field { max-width: 200px; }
  .main-content { padding: 0.75rem; }
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
- Background: #f5f5f5
- Border: 1px solid #e8e8e8
- Link padding: 0.75rem 1rem
- Link hover: background #ebebeb
- Menu label: 0.85rem, uppercase, #7a7a7a

---

## Main Content Area

**Scrollable:** `overflow-y: auto` when content exceeds viewport height.

**Padding:** 1.25rem desktop, 0.75rem mobile.

**Background:** #fafafa (light gray to distinguish from white shell).

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
- All interactive elements (buttons, links, inputs) receive focus outline: `2px solid #3273dc`, offset 2px
- No negative `tabindex` (all elements naturally focusable)

**Focus states:**
- Buttons: `outline: 2px solid #3273dc; outline-offset: 2px;`
- Links: same outline + hover background color change
- Inputs: `outline` + `border-color: #3273dc` on focus
- High contrast color (#3273dc on white background) meets AA requirements

**Color contrast:** All text on backgrounds meets 4.5:1 ratio for normal text, 3:1 for large text (≥18pt or ≥14pt bold).

---

## CSS Framework & Dependencies

**Bulma CSS:** Base responsive framework, bundled with the Vue build and embedded in the binary via the compiled SPA's static assets.
- Used for button, input, level, box, container, section, title classes
- Media query breakpoints: desktop (≥1024px), tablet (769–1023px), mobile (≤768px)
- Applied as class names directly in Vue single-file component templates

**Vue Router:** Client-side routing and shell/content composition. A persistent `AppShell.vue` layout component wraps routed screens, guaranteeing the shell is always present around feature content.

**Pinia:** Client-side state management for cross-screen state (e.g., sidebar open/closed, active member selection, active cycle).

**No external CDN calls:** All frontend assets (JS, CSS) are bundled by the Vite build and embedded in the binary. Fonts via Font Awesome CDN only (deferred to v2 for offline-capable font embedding).

---

## Interaction Patterns

**Sidebar toggle (mobile/tablet):**
1. User clicks toggle button
2. Sidebar open/closed state lives in a Pinia store (or local component state), toggled on click
3. Sidebar class bound reactively (e.g., `:class="{ 'is-active': sidebarOpen }"`) slides in
4. Overlay bound the same way dims background
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
- Primary action: Bulma blue (#3273dc) — used for active states, focus outlines, primary buttons
- Background: #fafafa (light gray)
- Borders: #e8e8e8
- Text: #4a4a4a (standard), #7a7a7a (secondary labels)
- Header/Sidebar: #2c3e50 (dark brand color for logo)

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
4. Use Bulma CSS classes established in the shell
5. Test focus states, keyboard navigation, and screen reader compatibility (component tests plus manual browser verification)

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
