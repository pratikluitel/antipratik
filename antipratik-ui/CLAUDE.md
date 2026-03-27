# antipratik.com — Claude Code Reference

Read this file at the start of every session before writing any code.

---

## Project Overview & Philosophy

**antipratik.com** is a personal brand site for a Kathmandu-based developer, musician, writer, and photographer.

### The Design Language

**Dark mode — Himalayan Dusk:**
> "The sky above Kathmandu 20 minutes after sunset — deep blue-grey with warmth underneath. Not pure black."

**Light mode — Parchment Morning:**
> "Warm off-white, aged paper. Not clinical white. Feels like a notebook left on a windowsill."

### Three Governing Principles

1. **Organic Minimalism** — Unhurried, textural without being decorative, high contrast but never harsh. Content recedes into the environment.
2. **Stillness** — A visitor should feel stillness 10 seconds after landing. Nothing urgent, nothing demanding attention.
3. **Himalayan / Nepali Cultural Texture** — Prayer flag colour system maps to content types. Mountain landscape palette informs surfaces. Grounded in place.

### Motion Principle
> "Everything breathes." No snapping. Transitions ease in and out at 400–500ms. The site feels alive without feeling animated.

---

## Tech Stack

| Layer | Choice | Reason |
|-------|--------|--------|
| Framework | Next.js 16 (App Router) | SSR for SEO, React Server Components, streaming |
| Styling | CSS Modules | The design system is a hand-crafted token system — utility classes fight it |
| Language | TypeScript (strict) | All interfaces defined in `src/lib/types.ts` |
| Fonts | DM Serif Display + DM Sans | Via Google Fonts preconnect in layout.tsx |
| Component libraries | None | Purpose-built only |
| CSS frameworks | None (no Tailwind) | Every class is a design decision, not a shorthand |

**App Router conventions:** `src/app/` for pages and layouts. `src/components/` for UI components. `src/lib/` for data layer and types.

---

## The 5 Sacred Rules

These rules are inviolable. Check them before writing any component.

### Rule 1 — No hardcoded values
**Never** hardcode a hex colour, an `rgb()`, or a px value that has a corresponding token.
- Wrong: `color: #E03E35;`
- Right: `color: var(--accent-music);`
- Wrong: `padding: 16px;`
- Right: `padding: var(--space-2);`

### Rule 2 — Music player absent from DOM, not hidden
The music player bar must be **genuinely absent from the DOM** when no music is playing.
- Wrong: `visibility: hidden` on a persistent element
- Wrong: `opacity: 0` on a persistent element
- Wrong: `display: none` on a persistent element
- Right: Conditionally render `{currentTrack && <MusicPlayer />}`
- Entry animation: insert into DOM → `translateY(100%)` → `translateY(0)` at 400ms ease-in-out

### Rule 3 — data-theme on `<html>`, not `<body>`
Theme switching is done by setting `data-theme="dark"` or `data-theme="light"` on the root `<html>` element.
All theme-dependent CSS is scoped to `[data-theme="dark"]` and `[data-theme="light"]` selectors.
Default is dark mode. Set via an inline `<script>` in layout.tsx before first paint.

### Rule 4 — All data through api.ts
**Never** call `fetch()` directly in a component or page.
All data fetching goes through `src/lib/api.ts`. This is the repository pattern — it abstracts dummy data vs. real API.

### Rule 5 — data-mode on the Filter Bar element
Pill colours are scoped to `[data-mode="dark"] .pill-essays` etc. Set `data-mode="dark"` or `data-mode="light"` on the filter bar container element itself — not on `<html>`. This prevents CSS specificity conflicts.

---

## CSS Token Categories

All tokens live in `src/styles/tokens.css`. Never import tokens.css anywhere except `src/app/layout.tsx`.

| Prefix | Category |
|--------|----------|
| `--font-serif`, `--font-sans` | Typefaces |
| `--text-*` | Type scale (display/h1/h2/h3/body/ui/meta/label/micro) |
| `--lh-*` | Line heights |
| `--ls-*` | Letter spacing |
| `--measure-*` | Max line lengths (ch/px) |
| `--space-1` through `--space-16` | 8px-grid spacing scale |
| `--gutter-*`, `--margin-*` | Grid gutters and margins |
| `--content-max-width`, `--breakpoint-*` | Layout bounds |
| `--accent-music/essays/short/photos/videos/links/social` | Prayer flag accents |
| `--pill-text-*-dark/light` | Pill text colours per mode |
| `--color-deepest`, `--color-bg-dark`, `--color-surface-dark` | Dark mode surfaces |
| `--color-elevated-dark`, `--color-border-dark` | Dark mode borders |
| `--color-text-primary/sub/body/muted/subtle-dark` | Dark mode text |
| `--color-surface/canvas/border/ink-light` | Light mode surfaces |
| `--color-text-primary/body/muted/subtle-light` | Light mode text |
| `--color-night-sky`, `--color-day-sky`, `--color-stone`, `--color-earth`, `--color-snow` | Mountain landscape |
| `--motion-fast/default/slow/breathe` | Transition durations |
| `--nav-*` | Navbar dimensions and colours |
| `--filter-bar-*`, `--pill-*` | Filter bar and pill specs |
| `--card-*` | Card base values |
| `--essay-*`, `--short-*`, `--music-*`, `--photo-*`, `--video-*`, `--link-*` | Per-content-type card tokens |
| `--player-*`, `--drawer-*`, `--waveform-*` | Music player tokens |
| `--article-*`, `--blockquote-*`, `--progress-*` | Article reading view |
| `--nl-*` | Newsletter block |
| `--link-row-*`, `--link-icon-*` | External links block |
| `--about-*` | About page |

---

## Prayer Flag → Accent Colour Mapping

These colours are **locked** to their content types. Do not use music red for anything except music content.

| Content Type | Hex | Token | Meaning |
|---|---|---|---|
| Music | `#E03E35` | `--accent-music` | Red — fire / energy |
| Essays | `#4A7FBB` | `--accent-essays` | Blue — sky / thought |
| Short Posts | `#D4A832` | `--accent-short` | Yellow — sun / quick flash |
| Photos | `#5E9E6A` | `--accent-photos` | Green — forest / life |
| Videos | `#4A7C6F` | `--accent-videos` | Teal — distinct from photos |
| Links | `#7A8890` | `--accent-links` | Slate — wind / carried word |

---

## Component Naming Conventions

```
src/components/
  Navbar/
    Navbar.tsx        ← React component, typed props interface
    Navbar.module.css ← only var(--token) values, no hardcoded hex/px
    index.ts          ← export { default } from './Navbar'
  FilterBar/
    ...
index.ts              ← barrel: export * from './Navbar'; export * from './FilterBar'; ...
```

Rules:
- PascalCase folder and file name — always
- One `.module.css` per component — never share
- Barrel `index.ts` in each component folder
- Append export to `src/components/index.ts` via `/new-component`
- No inline styles. No `style={{}}` props (except truly dynamic values like animation progress)

---

## Data Layer Contract

```
src/lib/
  types.ts         ← all TypeScript interfaces (single source of truth)
  api.ts           ← repository pattern: dummy data or API depending on env var
  dummy-data/
    posts.ts       ← Post[] sorted newest-first
    links.ts       ← ExternalLink[]
```

**The contract:**
1. Components import from `src/lib/api.ts` only — never from `dummy-data/` directly
2. `api.ts` checks `process.env.NEXT_PUBLIC_API_URL`:
   - Set → fetch from Go backend at `${API_URL}/endpoint`
   - Not set → return from dummy data
3. Function signatures never change when switching from dummy to real API

**Go backend endpoints (future):**
- `GET /api/posts` → `getPosts()`
- `GET /api/posts/:slug` → `getPost(slug)`
- `GET /api/links` → `getLinks()`
- `GET /api/links/featured` → `getFeaturedLinks()`

---

## Build Order

Build in this exact order to respect component dependencies:

| Step | Component/Page | Depends On |
|------|----------------|------------|
| 1 | `tokens.css` | Nothing — build this first |
| 2 | Navbar | tokens.css |
| 3 | Filter Bar | Navbar (shares height token for sticky positioning) |
| 4 | Feed Card components | tokens.css |
| 5 | Lightbox | Photo card (triggered by it) |
| 6 | Music Player Bar + Drawer | Music card (triggers it) |
| 7 | Article Reading View | tokens.css · Navbar (article title behaviour) |
| 8 | Newsletter Block | tokens.css |
| 9 | External Links Block | tokens.css |
| 10 | Feed Page | Navbar · Filter Bar · All feed cards · Music Player |
| 11 | Homepage | Navbar · Feed cards · Links Block · Newsletter |
| 12 | Links Page | Navbar · Links Block |
| 13 | About Page | Navbar · Newsletter |
| 14 | Article Page | Navbar · Article Reading View · Music Player |

---

## What Claude Code Must NEVER Do

1. **Hardcode any colour** — not `#0F1118`, not `rgb(15,17,24)`, not `rgba(...)` without a corresponding token. Use `var(--token-name)`.
2. **Use `visibility: hidden` on the music player** — the player must be removed from DOM, not hidden.
3. **Use Tailwind classes** — no `className="flex items-center gap-4 text-sm"`. This project uses CSS Modules only.
4. **Fetch data in a component** — no `fetch('/api/posts')` in a component. All data goes through `src/lib/api.ts`.
5. **Import tokens.css more than once** — it's imported once in `src/app/layout.tsx`. Never again.
6. **Set data-theme on `<body>`** — always on `<html>`.
7. **Write CSS that doesn't use design tokens** — every property value should come from a `var(--*)`.
8. **Use `px` values in CSS that have corresponding `--space-*` tokens** — use the tokens.
9. **Cross the typography line** — DM Serif Display = content (titles, headings). DM Sans = interface (dates, tags, metadata). Never cross.
10. **Break the prayer flag mapping** — music red on non-music content, essay blue on non-essay content, etc.
