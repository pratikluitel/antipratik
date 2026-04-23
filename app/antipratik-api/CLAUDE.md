# antipratik-api — Claude Code Reference

## Tech Stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.22+ |
| HTTP Framework | stdlib `net/http` |
| Database | SQLite via `modernc.org/sqlite` (pure Go, no CGO) |
| Config | YAML via `gopkg.in/yaml.v3` |
| Auth | JWT via `github.com/golang-jwt/jwt/v5` (7-day expiry) |
| Password Hashing | bcrypt via `golang.org/x/crypto/bcrypt` |
| UUID | `github.com/google/uuid` |

`ANTIPRATIK_HOST` and `ANTIPRATIK_PORT` env vars override server address.

---

## Architecture: 3-Layer Factory Pattern

```
components/posts/               ← post and link CRUD
  api/ → logic/ → store/

components/auth/                ← authentication, JWT, bootstrap
  api/ → logic/ → store/

components/files/               ← file storage, upload processing, serving
  api/ → logic/ → store/
  services/                     ← StorageService + UploaderService (cross-component injection)

components/broadcaster/         ← newsletter, email broadcast, contact form
  api/ → logic/ → store/
  lib/resend/                   ← Resend SMTP client (EmailSender interface)
  services/                     ← SubscriberService (cross-component injection)
  logic/emails/dist/            ← pre-built React Email HTML (gitignored; built by CI)

handlers/                       ← route registration, CORS, rate-limit middleware
  routes.go                     ← RegisterRoutes
  middleware.go                 ← CORSMiddleware, RateLimitMiddleware

common/db/                      ← SQLite connection + migrations (infrastructure only)
migrations/                     ← SQL files (embedded in main.go via //go:embed)
```

Cross-component calls go through the `services/` interface, never importing another component's `api/`, `logic/`, or `store/` packages. Wiring happens only in `main.go`.

**Wiring example:**
```go
postStore := postsstore.NewPostStore(db)
postLogic := postslogic.NewPostLogic(postStore, storageSvc, logger)
postH     := postsapi.NewPostHandler(postLogic, uploaderSvc, logger)
```

---

## Sacred Rules

### Rule 1 — Validate all input
Every parameter validated before processing. Use descriptive messages: `"slug cannot be empty"` not `"invalid input"`.

### Rule 2 — Defensive programming
Check negative numbers, malformed URLs, empty required arrays. Trim whitespace.

### Rule 3 — Readable error messages
Log technical details internally; expose only safe messages to clients.

### Rule 4 — JWT middleware on protected routes
All POST/PUT/DELETE endpoints use `JWTAuthMiddleware` from `components/auth/api/middleware.go`, applied in `handlers/routes.go`. Return 401 for invalid/missing tokens.

### Rule 5 — No direct DB access in API layer
API handlers call logic layer; logic calls store. No `db.Query()` in handlers.

### Rule 6 — Context propagation
Pass `context.Context` through all layers using `r.Context()`.

### Rule 7 — Structured logging
Use `common/logging.Logger` — never `log.Printf` or `fmt.Print`. Logger constructed once in `main.go` and passed through constructors.
- `INFO` — startup lifecycle only (in `main.go`)
- `ERROR` — internal failures producing 500 (in API layer via `handleLogicError`)
- Never log: validation errors (400), 404s, 401s, passwords, tokens, or personal data.

### Rule 8 — Shared error definitions
- `common/errors` — `ValidationError`, `New`, `Is`, `RequireNonEmpty`, `RequirePositive`
- `api/errors.go` — `handleLogicError`: maps `commonerrors.Is(err)` → 400, else → 500
- Component-specific sentinels stay in the owning package. Never redefine `ValidationError` locally.

### Rule 9 — Consistent JSON format
Success: direct JSON object or array. Error: `{"error": "message"}`. Created: `{"id": "uuid"}`. Use `writeJSON()` and `writeError()` helpers.

### Rule 10 — Config from YAML + env vars
Never hardcode sensitive values. Config loaded from YAML; env vars override.

### Rule 11 — Migration-based schema
Schema changes via SQL files in `migrations/` (embedded in `main.go`). Run on startup via `db.RunMigrations`.

### Rule 12 — Store never called directly from `main`
`main.go` never imports or calls `store` for business operations — delegate to logic layer. `common/db/` (Open, RunMigrations) is infrastructure and may be called from `main.go`.

### Rule 13 — Always rate-limit public POST endpoints
Any `POST` that writes to the DB without JWT must be wrapped with `RateLimitMiddleware` in `handlers/routes.go`. Current: `POST /api/subscribe` → 3 req/hour per IP.

### Rule 14 — Photo post must always have ≥ 1 image
`DELETE /api/posts/{id}/images/{imageID}` returns 400 if the post has only one image remaining. Enforced in `postLogic.DeletePhotoImage`. Never bypass.

### Rule 15 — Never violate component boundaries

Components only expose capabilities via their `services/` package (implementing the root `interface.go`). Other components import the interface, never the concrete implementation.

**Wrong:**
```go
import postslogic "github.com/pratikluitel/antipratik/components/posts/logic"
broadcasterlogic.NewBroadcasterLogic(..., postLogic, ...) // violation
```

**Right:**
```go
import postsservices "github.com/pratikluitel/antipratik/components/posts/services"
postsSvc := postsservices.NewPostsService(postLogic) // returns posts.PostsService
broadcasterlogic.NewBroadcasterLogic(..., postsSvc, ...) // correct
```

Rules: only the component's root package may be imported by other components. Models shared across components live in the root `models.go`. Logic/store/API interfaces are internal.

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/posts` | No | All posts, filtered by `type` and `tag` |
| GET | `/api/posts/{slug}` | No | Single essay by slug |
| POST | `/api/posts/essay` | JWT | Create essay |
| PUT | `/api/posts/essay/{id}` | JWT | Update essay |
| POST | `/api/posts/short` | JWT | Create short post |
| PUT | `/api/posts/short/{id}` | JWT | Update short post |
| POST | `/api/posts/music` | JWT | Create music post |
| PUT | `/api/posts/music/{id}` | JWT | Update music post |
| POST | `/api/posts/photo` | JWT | Create photo post |
| PUT | `/api/posts/photo/{id}` | JWT | Update photo post |
| POST | `/api/posts/video` | JWT | Create video post |
| PUT | `/api/posts/video/{id}` | JWT | Update video post |
| POST | `/api/posts/link` | JWT | Create link post |
| PUT | `/api/posts/link/{id}` | JWT | Update link post |
| DELETE | `/api/posts/{id}` | JWT | Delete any post type |
| GET | `/api/links/featured` | No | Up to 4 featured external links |
| GET | `/api/links` | No | All external links |
| POST | `/api/links` | JWT | Create external link |
| PUT | `/api/links/{id}` | JWT | Update external link |
| DELETE | `/api/links/{id}` | JWT | Delete external link |
| POST | `/api/auth/login` | No | Login, returns JWT |
| GET | `/api/openapi.yaml` | No | OpenAPI 3.0 spec |
| GET | `/api/index.html` | No | Swagger UI |
| GET | `/files/{fileId}` | No | Stream file binary |
| GET | `/thumbnails/{thumbnailId}` | No | Stream thumbnail binary |

### File Upload Contract

File uploads are `multipart/form-data` fields on post endpoints — no separate upload routes.

| Post type | File fields |
|-----------|-------------|
| Music | `audioFile` (required), `albumArtFile` (optional) |
| Photo | `images[]` (one or more, required) |
| Video | `thumbnailFile` (optional) |
| Link | `thumbnailFile` (optional) |

- Allowed photo types: `jpg`, `jpeg`, `png`, `webp`. Audio: `mp3`, `wav`, `ogg`, `m4a`.
- Photo uploads auto-generate 4 thumbnail variants: tiny (20px), small (300px), medium (600px), large (1200px). Widths defined in `logic/uploads.go`.
- Stored keys: `photos/<postId>-<index>.<ext>`, `music/<postId>.<ext>`, `thumbnails/<postId>-<index>-<size>.<ext>`.
- All URL fields in responses are relative (`/files/…`, `/thumbnails/…`) — frontend prefixes with API base URL.
- **Tag handling in multipart:** Use `formTags(r)` helper (`components/posts/api/helpers.go`) in all multipart handlers. Returns: `nil` (key absent, non-multipart → preserve tags), `[]string{}` (absent in multipart or empty → clear all tags), `[]string{…}` (parsed values). Never read `r.Form["tags"]` directly.
- Never add a generic `/uploads/*` endpoint. Exception: per-resource sub-collection endpoints like `POST /api/posts/{id}/images` are acceptable for managing images on an existing post.
- Never expose R2 object URLs — all file access goes through `/files/{fileId}` and `/thumbnails/{thumbnailId}`.

---

## HTTP Status Codes

| Status | When to use |
|--------|-------------|
| `200 OK` | GET, successful POST returning data |
| `201 Created` | POST creating a resource (returns `{"id": "…"}`) |
| `204 No Content` | DELETE |
| `400 Bad Request` | `commonerrors.ValidationError` |
| `401 Unauthorized` | Missing/invalid JWT |
| `404 Not Found` | Resource doesn't exist |
| `429 Too Many Requests` | Rate limited |
| `500 Internal Server Error` | Any non-ValidationError from logic/store |

Use `handleLogicError` for all logic errors — never manually write a 500.

---

## Conventions

### Input Type Naming

| Type | Used for |
|------|----------|
| `EssayPostInput`, `ShortPostInput`, etc. | Create and the merged value inside update |
| `UpdateEssayPost`, `UpdateMusicPost`, etc. | Partial update — all fields are pointers; `nil` = leave unchanged |

Never add pointer fields to `*PostInput` types for partial updates — use the separate `Update*Post` type.

### Post Type Constants
```go
models.PostTypeEssay, models.PostTypeShort, models.PostTypeMusic,
models.PostTypePhoto, models.PostTypeVideo, models.PostTypeLink
```
Use constants everywhere — never bare string literals.

### Concurrency
- `UploadPhotoFiles` processes images concurrently, bounded by `maxConcurrentUploads = 4` (`logic/uploads.go`).
- `DeletePost` deletes the DB record first (post immediately invisible), then cleans up files in a background goroutine using `context.Background()`. File-delete failures are logged but don't affect the response — intentional.

### Transaction Pattern
Multi-step writes that must be atomic (e.g. insert post then insert tags) use `*sql.Tx`. Never split logically atomic operations across two separate store calls from the logic layer.

---

## Database Schema

| Table | Purpose |
|-------|---------|
| `posts` | Base post metadata (id, type, created_at) |
| `post_tags` | Many-to-many post-tag relationships |
| `essay_posts` | Essay-specific data |
| `short_posts` | Short post content |
| `music_posts` | Music metadata |
| `photo_posts` | Photo galleries |
| `video_posts` | Video metadata |
| `link_posts` | Link post data |
| `external_links` | Curated external links |
| `users` | User accounts |
| `settings` | Key-value config storage |

`photo_images` has nullable thumbnail columns added in migrations `004`–`005`: `thumbnail_tiny_url`, `thumbnail_small_url`, `thumbnail_medium_url`, `thumbnail_large_url`. Existing rows may have `NULL`.

---

## Broadcaster Component

Templates written as React Email components in `app/emails/`, built to static HTML by CI, embedded into the Go binary via `//go:embed` in `components/broadcaster/logic/templates.go`.

**Build flow:** `npm run build` in `app/emails/` → copies `dist/` to `components/broadcaster/logic/emails/dist/` → `go build` embeds.  
For local dev: run `npm run build` in `app/emails/` and copy `dist/` manually before running the Go server.

**Token substitution:** Templates use `__TOKEN__` placeholders replaced at send time by `strings.NewReplacer`. `__UNSUBSCRIBE_TOKEN__` and `__POSTS_HTML__` are substituted per-subscriber at dispatch time.

**Post adapter:** `components/broadcaster/logic/post_adapter.go` adapts `posts.PostsService` to the broadcaster's internal `PostService` interface. Wired in `main.go` via `postsservices.NewPostsService(postLogic)`.

**Email image URLs:** All file/thumbnail URLs in the DB are relative. The broadcaster's `absoluteURL` helper prefixes with `cfg.SiteDomain` before writing into email HTML. `site_domain` must be set to the public base URL (e.g. `https://antipratik.com`) — set via `ANTIPRATIK_SITE_DOMAIN` env var in production.

**Email click-through URLs:**

| Post type | Email link destination |
|-----------|----------------------|
| essay | `{site_domain}/{slug}` |
| photo | `{site_domain}/feed?photo={id}` |
| music | `{site_domain}/feed?track={id}` |
| video | video URL directly |
| link | external URL directly |

**Storage backend config** (`config.yaml`):
- `local` — files written to `storage.local_dir` (default `./data/uploads/`)
- `r2` — requires `storage.r2.endpoint`, `bucket`, `access_key_id`, `secret_access_key` (supply via secrets, not committed config)

---

## Code Quality

After every code change:

```bash
CGO_ENABLED=0 go vet ./...
CGO_ENABLED=0 golangci-lint run
```

Both must produce zero errors. False-positive warnings may be suppressed with `//nolint:<linter>` + explanation, but errors must be fixed.

## Running Locally

```bash
cd app/antipratik-api
go run ./main.go
# or with hot reload:
air
```
