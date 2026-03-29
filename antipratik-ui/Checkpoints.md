# antipratik.com — Checkpoints
**Last updated:** 2026-03-29
**Current Status:** Full frontend complete — all pages (/, /feed, /links, /about, /[slug]), all card types, music player system, lightbox, article view, theme switching. Backend is dummy data. Two token audits pending (ExternalLinksBlock, NewsletterBlock). No known bugs.

This is a growing document. It records every implementation decision, deviation, and discovered rule that emerged during actual development. Periodically, stable knowledge gets compacted into `CLAUDE.md` and removed here — so this file stays lean and reflects only what hasn't yet been absorbed.

When this file and `CLAUDE.md` conflict, this file wins — it reflects what was actually built.

---

## Project State Snapshot

| Phase | Description | Status |
|---|---|---|
| Phase 0 | Scaffold, CLAUDE.md, slash commands, tokens, types, dummy data | ✅ Complete |
| Phase 1 | Navbar, ThemeProvider | ✅ Complete |
| Phase 2 | All 6 feed card components + PostCard dispatcher | ✅ Complete |
| Phase 3 | FilterBar, FeedPageClient, ClusterDivider, feed route | ✅ Complete |
| Phase 4 | MusicPlayer system (MusicProvider, MusicPlayerRoot, MusicPlayer, Waveform) | ✅ Complete |
| Phase 5 | Lightbox | ✅ Complete |
| Phase 6 | ExternalLinksBlock, NewsletterBlock, Homepage, Links, About, Article pages | ✅ Complete |

**Stack confirmed:** Next.js 16.2.1 + React 19.2.4 + TypeScript (strict) + CSS Modules + Turbopack

---

## Component Architecture — Client vs Server

| Component | Final type | Reason |
|---|---|---|
| EssayCard | Server | Display only, no state |
| ShortPostCard | Server | Display only; receives `onTagClick` callback — still server-compatible via prop |
| MusicCard | Client | Uses `useMusicPlayer()` context directly (no onPlay prop) |
| PhotoCard | Client | Needs onClick for lightbox |
| VideoCard | Server | External link only |
| LinkCard | Server | External link only |
| PostCard | Client | Contains client children (MusicCard) |
| FeedPageClient | Client | useReducer, useMemo, lightbox state |
| HomeFeedClient | Client | Simplified FeedPageClient — no FilterBar, direct post map |
| FilterBar | Client | MutationObserver for data-mode sync |
| ArticleClient | Client | Scroll listener, NavbarContext, reading progress |

---

## Decisions Made During Implementation

### Decision 1 — `data-scrolled` attribute bridge

**Problem:** Navbar's scroll detection (JS) needed to communicate compact state to FilterBar's CSS without a shared context.

**Decision:** Navbar sets `document.documentElement.setAttribute('data-scrolled', 'true'/'false')` on scroll. FilterBar uses `[data-scrolled='true'] .bar { top: var(--nav-height-compact); }` CSS selector.

**Why:** Avoids prop drilling or a scroll-specific context. The `<html>` element already carries `data-theme`; `data-scrolled` follows the same pattern.

### Decision 2 — `data-mode` on FilterBar for pill colour scoping

**Problem:** Pill colours need to respond to dark/light mode, but using `[data-theme="dark"]` selectors on pill CSS would conflict with other components' theme selectors.

**Decision:** FilterBar uses a `MutationObserver` to watch `document.documentElement`'s `data-theme` attribute and mirror it as `data-mode` on the FilterBar element itself. Pill CSS is scoped to `[data-mode="dark"] .pill` and `[data-mode="light"] .pill`.

**Why:** Prevents CSS specificity conflicts. Isolated scoping on the component element is cleaner than relying on the global `<html>` attribute for component-internal logic.

### Decision 3 — ClusterDivider simplified to a plain line

**Original spec:** ClusterDivider showed text labels like "↓ photos" or "↓ reading".

**Actual implementation:** Simplified to a plain `0.5px` horizontal rule with no text. `height: 0.5px; background: var(--color-border-dark)`.

**Why:** User preference. The labels added visual noise without value.

### Decision 4 — DateMarker removed from rendered feed

**Original spec:** `buildFeedClusters` emits `kind: 'date'` items; FeedPageClient renders `<DateMarker>` for them.

**Actual implementation:** FeedPageClient returns `null` for `kind === 'date'` items. `DateMarker` component exists but is not rendered. Date information is shown inside each card.

**Why:** Floating date text above cards looked disconnected. Dates belong to the cards themselves.

### Decision 5 — ShortPostCard layout: flex footer, not absolute positioning

**Original spec:** Date positioned `absolute; bottom: 14px; right: 16px`.

**Actual implementation:** Flex footer row — date on left, hashtags on right, `justify-content: space-between`.

**Why:** Absolute positioning caused date/hashtag overlap when card content was short. Flex layout handles any card height correctly.

### Decision 6 — MusicCard owns its own context call

**Original spec:** MusicCard received `onPlay: (post: MusicPost) => void` prop, which FeedPageClient provided.

**Actual implementation:** MusicCard calls `useMusicPlayer()` directly. The `onPlay` prop was removed from MusicCard, PostCard, and FeedPageClient.

**Why:** Cleaner — MusicCard is the only component that triggers music. Prop-drilling through PostCard → FeedPageClient added no value. Context is the right tool for global state.

### Decision 8 — `color-mix()` for alpha accent colours

Where the design spec required `rgba(accent-color, 12%)` backgrounds (e.g. icon boxes in ExternalLinksBlock), CSS `color-mix()` was used:

```css
background: color-mix(in srgb, var(--link-accent-music) 12%, transparent);
border-color: color-mix(in srgb, var(--link-accent-music) 25%, transparent);
```

**Why:** Derives alpha variants from the existing accent tokens without hardcoding hex. Token-compliant.

---

## Bugs Found and Fixed — Pattern Library

### React: setState during render error

**Symptom:** `Cannot update a component (MusicProvider) while rendering a different component (MusicPlayer)`

**Cause:** `onStop()` was called inside a `setCurrentTime` updater function. React disallows updating a different component during another component's state update.

**Fix:** Separate into two effects. One effect manages the interval and only updates `currentTime`. A second dedicated `useEffect` watches `currentTime` and calls `onStop()` as a proper side effect.

### CSS: pause icon off-center

**Symptom:** Pause icon appeared off-center in the play button circle.

**Cause:** The pause icon was built as a 2px-wide element with `box-shadow` creating the second bar. `box-shadow` offset doesn't participate in flexbox layout — the 2px element was centered, not the visual span.

**Fix:** Use `linear-gradient` on an 8px-wide element to render both bars. The full visual span is centered correctly by flexbox.

### Audio: track switching broken

**Symptom:** Switching tracks left the previous track playing while the UI showed the new track. Pause did not work.

**Root cause:** `isPlaying` stays `true` when switching tracks (only `activeTrack` changes). The `isPlaying` effect had `if (!audio || !track.audioUrl) return` — this returned early for tracks with no audioUrl, never calling `audio.pause()`.

**Fix:** Split the condition. Always pause (regardless of audioUrl). Only play if `isPlaying && track.audioUrl`:

```typescript
if (isPlaying && track.audioUrl) {
  audio.play().catch(...)
} else {
  audio.pause()
}
```

### CSS: FilterBar pill selected text colour

**Symptom:** Selected pill text was dark/black, hard to read against the coloured background.

**Fix:** Selected pill text uses `var(--color-snow)` with `!important` to override the mode-scoped colour variants. `--color-snow` is `#EEF2F0` — light enough to read on all accent backgrounds.

---

## Pending Token Audits

Two token audits from Phase 6 were not yet run at session end:

- `/check-tokens src/components/ExternalLinksBlock/ExternalLinksBlock.module.css`
- `/check-tokens src/components/NewsletterBlock/NewsletterBlock.module.css`

Run these at the start of the next session.
