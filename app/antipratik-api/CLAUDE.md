# antipratik-api — Claude Code Reference

Go backend for antipratik.com. Provides a JSON REST API and serves the static Next.js frontend build.

---

## Tech Stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.22+ |
| HTTP | stdlib `net/http` — no framework |
| Database | SQLite via `modernc.org/sqlite` (pure Go, no CGO) |
| Config | YAML via `gopkg.in/yaml.v3` |

---

## Project Layout

```
main.go              — entry point: load config → open DB → migrate → seed → route → listen
config/config.go     — Config struct, Load(path) → *Config
models/models.go     — shared Go structs (Post interface + 6 concrete types, ExternalLink, FilterState)
store/store.go       — PostStore, LinkStore interfaces
store/db.go          — Open(path) → *sql.DB, RunMigrations(db)
store/seed.go        — SeedIfEmpty(db) — inserts 12 posts + 8 links if empty
store/posts.go       — SQLitePostStore: NewPostStore, GetPosts, GetPostBySlug
store/links.go       — SQLiteLinkStore: NewLinkStore, GetLinks, GetFeaturedLinks
logic/logic.go       — PostLogic, LinkLogic interfaces
logic/posts.go       — PostService: NewPostService, GetPosts, GetPost
logic/links.go       — LinkService: NewLinkService, GetLinks, GetFeaturedLinks
api/api.go           — PostHandler, LinkHandler interfaces
api/posts.go         — PostHandlerImpl: NewPostHandler, GetPosts, GetPost
api/links.go         — LinkHandlerImpl: NewLinkHandler, GetLinks, GetFeaturedLinks
api/static.go        — NewSPAHandler(dir) for serving Next.js out/
api/middleware.go    — CORSMiddleware
config.yaml          — runtime configuration
data/                — SQLite database file (gitignored)
```

---

## Architecture: 3-Layer Factory Pattern

Each layer has an interface and a concrete implementation created via a factory function.
Dependencies are injected downward: `api` → `logic` → `store`.

```
api.PostHandlerImpl  { logic logic.PostLogic  }
logic.PostService    { store store.PostStore  }
store.SQLitePostStore{ db    *sql.DB          }
```

Wiring in `main.go`:
```go
postStore := store.NewPostStore(db)
postLogic := logic.NewPostService(postStore)
postH     := api.NewPostHandler(postLogic)
```

---

## API Endpoints

| Method | Path | Query Params | Description |
|--------|------|-------------|-------------|
| GET | `/api/posts` | `type` (repeatable), `tag` (repeatable) | All posts, newest first. No filter = all posts. |
| GET | `/api/posts/{slug}` | — | Single essay by slug. 404 if not found. |
| GET | `/api/links/featured` | — | Featured external links, up to 4. |
| GET | `/api/links` | — | All external links. |

All endpoints return `application/json`. Errors: `{"error":"message"}`.

### Post JSON shape

Posts are stored in separate per-type tables. The API merges base fields
(`id`, `type`, `createdAt`, `tags`) with type-specific fields into a single flat object.
JSON keys use camelCase to match the TypeScript types in `antipratik-ui`.

Example essay response:
```json
{
  "id": "essay-001",
  "type": "essay",
  "createdAt": "2026-03-15T08:30:00Z",
  "tags": ["philosophy", "music"],
  "title": "On Impermanence and Code",
  "slug": "on-impermanence-and-code",
  "excerpt": "...",
  "body": "...",
  "readingTimeMinutes": 7
}
```

---

## Database Schema

| Table | Purpose |
|-------|---------|
| `posts` | Base record: `id`, `type`, `created_at` |
| `post_tags` | Normalized tags: `post_id`, `tag` |
| `essay_posts` | Essay-specific fields |
| `short_posts` | Short post body |
| `music_posts` | Track metadata |
| `photo_posts` | Photo post location |
| `photo_images` | Individual images (ordered by `sort_order`) |
| `video_posts` | Video metadata |
| `link_posts` | Curated link metadata |
| `links` | External links (separate entity, not posts) |

---

## Running Locally

```bash
cd app/antipratik-api
go run ./main.go
# custom config:
go run ./main.go --config /path/to/config.yaml
```

For development with hot reloading (automatic restart on code changes):
```bash
cd app/antipratik-api
air
```

The server seeds the database on first run. Restarts are safe — seeding is skipped if posts exist.

Build the frontend first for static serving:
```bash
cd app/antipratik-ui && npm run build
```

---

## Adding a New Post Type

1. Add the new type to the `CHECK` constraint in `store/db.go` and run a migration.
2. Add a new detail table in `store/db.go`.
3. Add the Go struct in `models/models.go` with a `postType()` method.
4. Add fetch and assembly logic in `store/posts.go`.
5. Add seed data in `store/seed.go`.
6. No handler or logic changes needed — the existing `GetPosts` is generic.
