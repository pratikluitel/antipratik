# antipratik.com — Claude Code Reference

## Tech Stack

| Layer | Choice |
|-------|--------|
| Framework | Next.js 16 (App Router) |
| Styling | CSS Modules — hand-crafted token system |
| Language | TypeScript (strict) — interfaces in `src/lib/types.ts` |
| Fonts | DM Serif Display + DM Sans (Google Fonts preconnect in layout.tsx) |
| Component/CSS libraries | None |

---

## Design Language

**Dark mode — Himalayan Dusk:** The sky above Kathmandu 20 minutes after sunset — deep blue-grey with warmth underneath. Not pure black.

**Light mode — Parchment Morning:** Warm off-white, aged paper. Not clinical white. Feels like a notebook left on a windowsill.

---

## Sacred Rules

### Rule 1 — No hardcoded values
Never hardcode a hex colour, `rgb()`, or px value that has a token. Use `var(--token-name)`.  
Exception: `--space-*` starts at 8px — values below 8px (2px, 4px, 5px, 6px) have no token and may be hardcoded.

### Rule 2 — Music player absent from DOM, not hidden
Conditionally render `{currentTrack && <MusicPlayer />}`. Never use `visibility: hidden`, `opacity: 0`, or `display: none` on a persistent element.  
Entry animation: insert → `translateY(100%)` → `translateY(0)` at 400ms ease-in-out.

### Rule 3 — data-theme on `<html>`, not `<body>`
Set `data-theme="dark"` or `data-theme="light"` on the root `<html>` element. Default is dark. Set via inline `<script>` in layout.tsx before first paint.

### Rule 4 — All data through api.ts
Never call `fetch()` directly in a component. All data goes through `src/lib/api.ts`.

### Rule 5 — data-mode on the Filter Bar element
Set `data-mode="dark"` or `data-mode="light"` on the filter bar container itself — not on `<html>`. Prevents CSS specificity conflicts.

### Rule 6 — NEVER use accent colors on UI chrome
`--accent-*` tokens are reserved for content type indicators (cards, tags, pills, player). Never on toggles, scrollbars, focus rings, or other UI infrastructure.

### Rule 7 — rgba() in @keyframes is the one token exception
CSS custom properties can't have opacity inside `@keyframes` without `color-mix()`. Raw `rgba()` is acceptable **only** inside `@keyframes`. Document the source token in a comment.

### Rule 8 — The homepage hero hardcode is intentional
`style={{ background: '#0F1118' }}` on the hero div is the only accepted hardcoded hex. It must be theme-resistant. Do not change it to a CSS variable.

### Rule 9 — audio.removeAttribute('src'), never audio.src = ''
`audio.src = ''` resolves to the page URL. Always use `removeAttribute('src')` to clear audio source.

### Rule 10 — params is a Promise in Next.js 16
Always await params in page components and `generateMetadata`:
```typescript
const { slug } = await params;
```

---

## CSS Tokens

All tokens in `src/styles/tokens.css`. Import only in `src/app/layout.tsx` — never elsewhere.

| Prefix | Category |
|--------|-----------|
| `--font-serif`, `--font-sans` | Typefaces |
| `--text-*`, `--lh-*`, `--ls-*`, `--measure-*` | Type scale, line heights, letter spacing, line lengths |
| `--space-1` through `--space-16` | 8px-grid spacing scale |
| `--border-thin` (1px), `--border-hairline` (0.5px) | Border widths |
| `--icon-xs` through `--icon-xl` | Icon and circular button sizes (20–44px) |
| `--gutter-*`, `--margin-*`, `--content-max-width`, `--breakpoint-*` | Layout bounds |
| `--accent-music/essays/short/photos/videos/links/social` | Prayer flag accents |
| `--color-deepest`, `--color-bg-dark`, `--color-surface-dark`, etc. | Dark mode surfaces/text |
| `--color-surface/canvas/border/ink-light`, etc. | Light mode surfaces/text |
| `--color-night-sky`, `--color-day-sky`, `--color-stone`, `--color-earth`, `--color-snow` | Mountain landscape |
| `--motion-fast/default/slow/breathe`, `--motion-theme` | Transition durations |
| `--nav-*` | Navbar dimensions and colours |
| `--filter-bar-*`, `--pill-*` | Filter bar and pill specs |
| `--card-*`, `--essay-*`, `--short-*`, `--music-*`, `--photo-*`, `--video-*`, `--link-*` | Card tokens |
| `--player-*`, `--drawer-*`, `--waveform-*` | Music player |
| `--article-*`, `--blockquote-*`, `--progress-*` | Article reading view |
| `--nl-*` | Newsletter block |
| `--link-row-*`, `--link-icon-*`, `--about-*` | External links, about page |
| `--admin-*` | Admin panel tokens |

**Breakpoints:** CSS custom properties can't be used inside `@media` queries — use pixel values directly, but they must match the token exactly:

| Token | Value | Use for |
|---|---|---|
| `--breakpoint-desktop` | `1280px` | Desktop+ |
| `--breakpoint-tablet` | `768px` | Tablet+ (min-width) |
| `--breakpoint-mobile` | `767px` | Mobile (max-width) |
| `--breakpoint-small` | `640px` | Narrow mobile (max-width) |

Always use `767px` (not `768px`) for `max-width` mobile media queries — matched pairs. Never introduce a new breakpoint without a corresponding `--breakpoint-*` token.

---

## Prayer Flag → Accent Colour Mapping

Locked to content types. Do not use these colours for anything else.

| Content Type | Hex | Token |
|---|---|---|
| Music | `#E03E35` | `--accent-music` |
| Essays | `#4A7FBB` | `--accent-essays` |
| Short Posts | `#D4A832` | `--accent-short` |
| Photos | `#5E9E6A` | `--accent-photos` |
| Videos | `#4A7C6F` | `--accent-videos` |
| Links | `#7A8890` | `--accent-links` |

---

## Component Naming Conventions

```
src/components/
  Navbar/
    Navbar.tsx        ← React component, typed props interface
    Navbar.module.css ← only var(--token) values
    index.ts          ← export { default } from './Navbar'
index.ts              ← barrel: export * from './Navbar'; ...
```

- PascalCase folder and file name — always
- One `.module.css` per component — never share
- Barrel `index.ts` in each component folder, appended to `src/components/index.ts`
- No inline styles except truly dynamic values (e.g. animation progress)

---

## Data Layer Contract

```
src/lib/
  types.ts    ← all TypeScript interfaces (single source of truth)
  api.ts      ← repository pattern: returns dummy data or real API data
```

`api.ts` checks `process.env.NEXT_PUBLIC_API_URL`: set → fetch from `${API_URL}/endpoint`; not set → return dummy data. Function signatures never change when switching.

**URL prefixing:** Backend stores relative URLs for all file fields (`/files/…`, `/thumbnails/…`). `api.ts` must prefix every file-backed URL with `API_URL` via the `prefixPost()` helper in `getPosts()`. Affected fields: `MusicPost.albumArt`, `MusicPost.audioUrl`, `PhotoImage.url`, `PhotoImage.thumbnailSmallUrl/thumbnailMediumUrl/thumbnailLargeUrl`, `VideoPost.videoUrl`, `VideoPost.thumbnailUrl`, `LinkPost.thumbnailUrl`.

**Type notes:**
- `ShortPost` has no `hashtags` field — use `post.tags` from `BasePost`
- `LinkPost.category` is optional; `ExternalLink.category` is required
- `PhotoImage` includes optional thumbnail fields: `thumbnailSmallUrl?`, `thumbnailMediumUrl?`, `thumbnailLargeUrl?` — use `thumbnailSmallUrl` for card previews, `thumbnailLargeUrl` for lightbox, falling back to `url`

---

## Routing

| Route | File |
|---|---|
| `/` | `src/app/page.tsx` |
| `/feed` | `src/app/feed/page.tsx` |
| `/links` | `src/app/links/page.tsx` |
| `/about` | `src/app/about/page.tsx` |
| `/[slug]` | `src/app/[slug]/page.tsx` (essay catch-all) |

`EssayCard` links to `/${post.slug}` — not `/post/${post.slug}`.

**Deep-link params** (used by newsletter email click-throughs):

| Param | Effect |
|---|---|
| `/feed?photo=<postId>` | Opens lightbox for that PhotoPost on mount |
| `/feed?track=<postId>` | Auto-plays that MusicPost on mount |
| `/feed?video=<postId>` | Opens VideoPlayer modal for that video-category LinkPost on mount |

All use `Promise.resolve().then(setState)` microtask to satisfy `react-hooks/set-state-in-effect`. Don't remove param handling without updating the Go broadcaster too.

---

## Context Architecture

**Provider nesting (layout.tsx):**
```
ThemeProvider
  └── NavbarProvider
        └── MusicProvider
              ├── Navbar
              ├── {children}
              └── MusicPlayerRoot
```

| Context | Purpose |
|---|---|
| `ThemeContext` | Dark/light mode, `localStorage` persistence, sets `data-theme` on `<html>` |
| `NavbarContext` | `articleTitle` string — written by `ArticleClient`, read by `Navbar` |
| `MusicContext` | `activeTrack`, `isPlaying`, `isExiting` — global playback state |

---

## What Claude Code Must NEVER Do

1. **Hardcode any colour** — not `#0F1118`, not `rgb(…)`, not `rgba(…)` without a token. Use `var(--token-name)`.
2. **Use `visibility: hidden` on the music player** — remove from DOM, not hidden.
3. **Use Tailwind classes** — CSS Modules only.
4. **Fetch data in a component** — all data through `src/lib/api.ts`.
5. **Import tokens.css more than once** — imported once in `src/app/layout.tsx`.
6. **Set data-theme on `<body>`** — always on `<html>`.
7. **Use `px` values that have corresponding `--space-*` tokens** — use the tokens.
8. **Cross the typography line** — DM Serif Display = content (titles, headings). DM Sans = interface (dates, tags, metadata).
9. **Break the prayer flag mapping** — music red on non-music content, etc.
10. **Use `box-shadow` for animated outlines** — animate `border-color` instead.
11. **Set `audio.src = ''`** — use `removeAttribute('src')`.
12. **Call setState on a different component inside a state updater function** — use a separate `useEffect`.
13. **Link essays to `/post/${slug}`** — correct path is `/${slug}`.
14. **Use `--accent-*` colors on UI chrome elements** — accent tokens are for content only.
15. **Set text color tokens below minimum contrast** — all text tokens must achieve ≥ 4.5:1 (WCAG AA) against expected background. Verify before changing any `--color-text-*`, `--link-*-color`, or `--nl-*-color` token.
16. **Manipulate theme-sensitive DOM attributes directly in effects** — drive `data-mode` and similar via `useState` + `useLayoutEffect` so values survive re-renders.
17. **Write critical CSS only under `[data-theme]` selectors** — if `data-theme` is absent the element is unstyled. Put dark-mode default on the base selector; light mode overrides with `[data-theme='light']`.
18. **Define a component-level `transition` that omits properties transitioned by the global `*` catch-all** — component `transition` overrides the globals.css catch-all entirely. Include every property that needs to animate, or properties not listed will snap on theme switch.

---

## Code Quality

After every change:

```bash
npm run lint
npm run typecheck
```

Both must produce zero errors. ESLint warnings for unavoidable false positives (e.g. `@next/next/no-img-element` for external API images) may be suppressed with `// eslint-disable-next-line` + explanation.
