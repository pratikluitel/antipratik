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

## The Sacred Rules

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

### Rule 6 — NEVER use accent colors on UI chrome
`--accent-*` tokens are reserved for content type indicators only (cards, tags, pills, player). Never use them on toggles, scrollbars, focus rings, or other UI infrastructure.

### Rule 7 — rgba() in @keyframes is the one token exception
CSS custom properties cannot have opacity applied inside `@keyframes` without `color-mix()`. Raw `rgba()` values are acceptable **only** inside `@keyframes` blocks. Document the source token in a comment.

### Rule 8 — The homepage hero hardcode is intentional
`style={{ background: '#0F1118' }}` on the hero div is the only accepted hardcoded hex in the codebase. It must be theme-resistant. Do not change it to a CSS variable.

### Rule 9 — audio.removeAttribute('src'), never audio.src = ''
Setting `audio.src = ''` resolves to the page URL. Always use `removeAttribute('src')` to clear audio source.

### Rule 10 — params is a Promise in Next.js 16
Always await params in page components and `generateMetadata`:
```typescript
const { slug } = await params;
```

---

## CSS Token Categories

All tokens live in `src/styles/tokens.css`. Never import tokens.css anywhere except `src/app/layout.tsx`.

| Prefix | Category |
|--------|-----------|
| `--font-serif`, `--font-sans` | Typefaces |
| `--text-*` | Type scale (display/h1/h2/h3/body/ui/meta/label/micro) |
| `--lh-*` | Line heights |
| `--ls-*` | Letter spacing |
| `--measure-*` | Max line lengths (ch/px) |
| `--space-1` through `--space-16` | 8px-grid spacing scale |
| `--gutter-*`, `--margin-*` | Grid gutters and margins |
| `--content-max-width`, `--content-column-width`, `--breakpoint-*` | Layout bounds |
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

### Token Additions & Modifications

These tokens were added or modified during implementation (not in the original design system HTML):

| Token | Value | Note |
|---|---|---|
| `--content-column-width` | `860px` | Main content column max-width — widened from 680px to use 2 more grid columns (1 each side) per issue #5 |
| `--article-column-width` | `680px` | Article reading view max-width — intentionally narrower than feed (860px) for prose comfort (~65ch at 16px body) |
| `--measure-article` | `680px` | Matches `--article-column-width`; kept in sync |
| `--nav-height` | `48px` | FilterBar needs this for sticky top positioning |
| `--nav-height-compact` | `40px` | Navbar compact state on scroll |
| `--nav-height-default` | `48px` | Alias for clarity in FilterBar CSS |
| `--z-navbar` | `100` | Navbar z-index |
| `--z-filterbar` | `90` | FilterBar z-index — must sit below navbar |
| `--color-text-muted-dark` | `#7A9AB4` | **Changed from `#3A5068` → `#607D96` → `#7A9AB4`** — previous values failed contrast on dark bg; affects navbar links, hero text, section meta |
| `--color-text-subtle-dark` | `#4A6A84` | **Changed from `#2E3748`** — was identical to `--color-border-dark`, near-invisible; affects ExternalLinksBlock labels and arrows |
| `--link-domain-color-dark` | `#4A6A84` | **Changed from `#2E3748`** — domain text in ExternalLinksBlock (dark mode) was unreadable |
| `--nl-legal-color` | `#607D96` (root) / `#7A9AB4` (per theme) | **Root changed from `#1E2535`**; then per-theme overrides added in `[data-theme]` blocks — `#607D96` fails WCAG AA on elevated newsletter bg (#8); both themes now use `var(--color-text-muted-dark)` = `#7A9AB4` |
| `--nl-bg` | `#0A0E14` (root fallback) / `#181D28` dark / `#1C2B3A` light | **Per-theme overrides added** — `#0A0E14` was indistinguishable from dark page bg and jarring in light mode; dark now uses `--color-surface-dark`, light uses `--color-stone` (#8) |
| `--nl-input-bg` | `#181D28` (root) / `#0F1118` dark | **Dark mode override added** — input bg must be darker than new block bg `#181D28` to remain visible; overridden to `--color-bg-dark` in dark mode (#8) |
| `--color-text-subtle-light` | `#7A7268` | **Changed from `#C8C0B4`** — was too close to parchment bg, near-invisible |
| `--link-domain-color-light` | `#8A8070` | **Changed from `#C8C0B4`** — domain text in ExternalLinksBlock (light mode) was unreadable |
| `--color-text-muted-light` | `#4A5860` | **Changed from `#7A8890`** — navbar links and section meta on parchment bg were below WCAG AA |
| `--pill-hover-bg-dark` | `var(--color-elevated-dark)` | FilterBar pill hover state bg in dark mode |
| `--pill-hover-bg-light` | `#D8D4CA` | FilterBar pill hover state bg in light mode (no existing token matches) |
| `--pill-all-selected-bg-dark` | `var(--color-border-dark)` | FilterBar "All" pill when all tags selected, dark mode |
| `--pill-all-selected-bg-light` | `var(--color-stone)` | FilterBar "All" pill when all tags selected, light mode |
| `--nav-toggle-width` | `44px` | **Updated from `32px`** — glassy pill toggle track width |
| `--nav-toggle-height` | `24px` | **Updated from `18px`** — glassy pill toggle track height |
| `--nav-toggle-radius` | `12px` | Toggle track border-radius |
| `--nav-toggle-padding-x` | `2px` | Toggle track horizontal padding (thumb inset) |
| `--nav-toggle-thumb-size` | `20px` | Toggle thumb diameter |
| `--nav-toggle-thumb-offset-dark` | `2px` | Thumb `left` position in dark mode |
| `--nav-toggle-thumb-offset-light` | `22px` | Thumb `left` position in light mode |
| `--nav-toggle-icon-size` | `12px` | Sun/moon icon size inside thumb |
| `--toggle-track-bg-dark` | `rgba(255, 255, 255, 0.08)` | Glassy pill toggle track bg in dark mode |
| `--toggle-track-bg-light` | `rgba(0, 0, 0, 0.06)` | Glassy pill toggle track bg in light mode |
| `--toggle-track-hover-bg-dark` | `rgba(255, 255, 255, 0.12)` | Toggle track hover bg in dark mode |
| `--toggle-track-hover-bg-light` | `rgba(0, 0, 0, 0.10)` | Toggle track hover bg in light mode |
| `--toggle-thumb-bg-dark` | `rgba(255, 255, 255, 0.15)` | Toggle thumb bg in dark mode |
| `--toggle-thumb-bg-light` | `rgba(0, 0, 0, 0.12)` | Toggle thumb bg in light mode |
| `--color-text-primary` | *(theme alias)* | Theme-switching alias for `--color-text-primary-dark/light` — set in `[data-theme]` blocks |
| `--color-text-sub` | *(theme alias)* | Theme-switching alias for `--color-text-sub-dark` / `--color-text-primary-light` |
| `--color-text-body` | *(theme alias)* | Theme-switching alias for `--color-text-body-dark/light` |
| `--color-text-muted` | *(theme alias)* | Theme-switching alias for `--color-text-muted-dark/light` |
| `--color-text-subtle` | *(theme alias)* | Theme-switching alias for `--color-text-subtle-dark/light` |
| `--admin-content-max` | `860px` | Admin panel main content max-width |
| `--admin-tab-height` | `44px` | Admin dashboard tab bar item height |
| `--admin-top-bar-height` | `52px` | Admin top bar height |
| `--admin-input-radius` | `6px` | Border-radius for admin form inputs, buttons, and rows |
| `--admin-input-padding-x` | `var(--space-2)` | Horizontal padding for admin form inputs |
| `--admin-input-padding-y` | `10px` | Vertical padding for admin form inputs (between `--space-1`=8px and `--space-2`=16px) |
| `--admin-danger` | `#C0392B` | Admin destructive/error colour (delete buttons, error states, required asterisk) — separate from `--accent-music` (Rule 6) |
| `--admin-form-section-radius` | `8px` | Admin form section card border-radius — wider than `--admin-input-radius` (6px) |
| `--admin-btn-radius` | `4px` | Smaller radius for inline action buttons (edit, delete) and tag chips |
| `--nav-bg-dark-blur` | `rgba(15, 17, 24, 0.85)` | Navbar bg in dark mode with backdrop blur — replaces `--nav-bg-dark` (solid) |
| `--nav-bg-light-blur` | `rgba(247, 245, 240, 0.85)` | Navbar bg in light mode with backdrop blur — replaces `--nav-bg-light` (solid) |
| `--video-overlay-bg` | `rgba(255, 255, 255, 0.92)` | VideoCard play-button overlay background |
| `--video-overlay-text` | `rgba(255, 255, 255, 0.55)` | VideoCard duration text overlay color |
| `--lightbox-btn-bg` | `rgba(238, 242, 240, 0.12)` | Lightbox close/nav button background (`--color-snow` @ 12%) |
| `--lightbox-btn-border` | `rgba(238, 242, 240, 0.20)` | Lightbox button border (`--color-snow` @ 20%) |
| `--lightbox-counter-color` | `rgba(238, 242, 240, 0.50)` | Lightbox image counter text (`--color-snow` @ 50%) |
| `--lightbox-caption-color` | `rgba(238, 242, 240, 0.55)` | Lightbox caption text (`--color-snow` @ 55%) |
| `--lightbox-dot-bg` | `rgba(238, 242, 240, 0.30)` | Lightbox inactive dot background (`--color-snow` @ 30%) |
| `--lightbox-nav-icon-size` | `18px` | Lightbox prev/next navigation icon size |
| `--player-close-btn-bg` | `rgba(238, 242, 240, 0.08)` | MusicPlayer drawer close button background (`--color-snow` @ 8%) |
| `--essay-card-padding-x` | `22px` | EssayCard inner horizontal padding |
| `--music-content-padding-x` | `18px` | MusicCard content area horizontal padding |
| `--photo-card-body-padding` | `12px var(--space-2) 14px` | PhotoCard body padding (top / sides / bottom) |

These contrast fixes are global — every component using these tokens benefits automatically.

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

### Type System Notes (deviations from original plan)

- **ShortPost** has NO `hashtags` field. Use `post.tags` from `BasePost`.
- **LinkPost.category** is optional: `'music' | 'writing' | 'video' | 'social'`
- **ExternalLink.category** is REQUIRED: `'music' | 'writing' | 'video' | 'social'`
- **FilterState** includes `sortOrder: 'newest' | 'oldest'`
- **Additional types added post-scaffold:** `LayoutMode`, `LAYOUT_MODE`, `FeedItem`, `FilterAction` — see `src/lib/types.ts`

**Go backend endpoints:**
- `GET /api/posts` → `getPosts()`
- `GET /api/posts/:slug` → `getPost(slug)`
- `GET /api/links` → `getLinks()`
- `GET /api/links/featured` → `getFeaturedLinks()`

**URL prefixing (important):** The backend stores relative URLs for all file fields (`/files/…`, `/thumbnails/…`). `api.ts` must prefix every file-backed URL with `API_URL` before returning data to components. This is done via the `prefixPost()` helper in `api.ts` — called inside `getPosts()`. Affected fields: `MusicPost.albumArt`, `MusicPost.audioUrl`, `PhotoImage.url`, `PhotoImage.thumbnailSmallUrl/thumbnailMediumUrl/thumbnailLargeUrl`, `VideoPost.thumbnailUrl`, `LinkPost.thumbnailUrl`.

**PhotoImage type** now includes optional thumbnail fields: `thumbnailSmallUrl?`, `thumbnailMediumUrl?`, `thumbnailLargeUrl?` — present when the image was uploaded through the backend. Use `thumbnailSmallUrl` for card previews and `thumbnailLargeUrl` for lightbox when available, falling back to `url`.

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

## Routing

| Route | File |
|---|---|
| `/` | `src/app/page.tsx` |
| `/feed` | `src/app/feed/page.tsx` |
| `/links` | `src/app/links/page.tsx` |
| `/about` | `src/app/about/page.tsx` |
| `/[slug]` | `src/app/[slug]/page.tsx` (essay catch-all) |

**Important:** `EssayCard` links to `/${post.slug}`, NOT `/post/${post.slug}`.
Next.js App Router prioritises named routes (`/feed`, `/links`, `/about`) over the `[slug]` catch-all. The catch-all only activates for paths that don't match a named route.

---

## Context Architecture

### Provider nesting order (layout.tsx)
```
ThemeProvider
  └── NavbarProvider
        └── MusicProvider
              ├── Navbar
              ├── {children}
              └── MusicPlayerRoot
```

### Contexts

| Context | Provider | Purpose |
|---|---|---|
| `ThemeContext` | `ThemeProvider` | Dark/light mode, `localStorage` persistence, sets `data-theme` on `<html>` |
| `NavbarContext` | `NavbarProvider` | `articleTitle` string — written by `ArticleClient`, read by `Navbar` |
| `MusicContext` | `MusicProvider` | `activeTrack`, `isPlaying`, `isExiting` — global playback state |

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
11. **Use `box-shadow` for animated outlines on cards** — animate `border-color` instead. `box-shadow` draws outside the border box, creating a visible gap artifact.
12. **Set `audio.src = ''`** — use `removeAttribute('src')` instead. Empty string resolves to the page URL.
13. **Call setState on a different component inside a state updater function** — use a separate `useEffect` to trigger cross-component state changes.
14. **Link essays to `/post/${slug}`** — correct path is `/${slug}`. Next.js catch-all handles it.
15. **Use `--accent-*` colors on UI chrome elements** — accent tokens are for content only.
16. **Set text color tokens below minimum contrast** — all text tokens must achieve ≥ 4.5:1 (WCAG AA normal text) or ≥ 3:1 (large/bold text) against their expected background. "Subtle" is for de-emphasis, not invisibility. Verify contrast before changing any `--color-text-*`, `--link-*-color`, or `--nl-*-color` token.
