# antipratik-api — Claude Code Reference

Read this file at the start of every session before writing any code.

---

## Project Overview & Philosophy

**antipratik-api** is the backend REST API for antipratik.com, a personal brand site for a Kathmandu-based developer, music tinkerer, and blogger.

### The Backend Philosophy

**Minimalist, Secure, and Robust:**
> "Build for the worst-case user — assume they're trying to break it. Validate everything, fail fast with clear errors, and never trust input."

**Defensive by Design:**
> "Every parameter is suspect. Every request could be malicious. Every response must be safe."

### Three Governing Principles

1. **Input Validation First** — All parameters must be strictly validated before processing. Invalid input is rejected with readable error messages.
2. **Layered Security** — Authentication, authorization, and input sanitization at every level.
3. **Clean Architecture** — Dependency injection with interfaces ensures testability and maintainability.

### Security Principle
> "JWT tokens expire in 7 days. Passwords are bcrypt-hashed. CORS is permissive in development but must be locked down in production."

---

## Tech Stack

| Layer | Choice | Reason |
|-------|--------|--------|
| Language | Go 1.22+ | Compiled, fast, excellent concurrency, strong typing |
| HTTP Framework | stdlib `net/http` | No unnecessary abstractions — direct control over HTTP |
| Database | SQLite via `modernc.org/sqlite` | Pure Go (no CGO), embedded, ACID compliant |
| Config | YAML via `gopkg.in/yaml.v3` | Human-readable, versioned in git |
| Auth | JWT via `github.com/golang-jwt/jwt/v5` | Industry standard, stateless |
| Password Hashing | bcrypt via `golang.org/x/crypto/bcrypt` | Slow, salted, secure |
| UUID | `github.com/google/uuid` | For generating unique IDs |
| JSON | stdlib `encoding/json` | No external dependencies for core functionality |

**Environment Overrides:** Server host and port can be overridden with `ANTIPRATIK_HOST` and `ANTIPRATIK_PORT` environment variables.

---

## Major Architecture Pattern: 3-Layer Factory Pattern

The API follows a **3-Layer Clean Architecture** with dependency injection:

```
┌─────────────────┐
│   API Layer     │  ← HTTP handlers, JSON serialization, data type validation
│   (api/*.go)    │
└─────────────────┘
         ↓
┌─────────────────┐
│  Logic Layer    │  ← Business logic, validation, coordination
│ (logic/*.go)    │
└─────────────────┘
         ↓
┌─────────────────┐
│  Store Layer    │  ← Data persistence, SQL queries
│ (store/*.go)    │
└─────────────────┘

┌─────────────────┐
│   DB Package    │  ← SQLite connection + schema migrations (infrastructure only)
│   (db/*.go)     │     Called from main.go; not a business-logic store
└─────────────────┘
```

Each layer has an **interface** and a **concrete implementation** created via factory functions. Dependencies flow downward through constructor injection.

**Wiring Example:**
```go
postStore := store.NewPostStore(db)
postLogic := logic.NewPostService(postStore)
postH     := api.NewPostHandler(postLogic)
```

**Benefits:**
- Testability: Mock interfaces for unit tests
- Separation of Concerns: Each layer has one responsibility
- Maintainability: Changes in one layer don't affect others

---

## The Sacred Rules

These rules are inviolable. Check them before writing any code.

### Rule 1 — All Parameters Must Be Validated
**Every input parameter must be validated** before processing. Assume the user is malicious.
- Wrong: Accept any string as a slug
- Right: Check for empty strings, invalid characters, length limits
- Use descriptive error messages: `"slug cannot be empty"` not `"invalid input"`

### Rule 2 — Defensive Programming: Assume the User is an Idiot
**Never trust user input.** Validate types, ranges, and business rules.
- Check for negative numbers where only positive make sense
- Validate URLs are well-formed
- Ensure arrays are not empty when required
- Trim whitespace and reject obviously malicious content

### Rule 3 — Readable Error Messages
**Error responses must be user-friendly and actionable.**
- Wrong: `{"error": "internal server error"}`
- Right: `{"error": "title cannot be longer than 200 characters"}`
- Log technical details internally, but expose only safe messages to clients

### Rule 4 — JWT Middleware on Protected Routes
**All write operations require JWT authentication.**
- Use `JWTAuthMiddleware` wrapper for POST/PUT/DELETE endpoints
- Validate tokens before processing any request
- Return `401 Unauthorized` for invalid/missing tokens

### Rule 5 — No Direct Database Access in API Layer
**API handlers never call the database directly.**
- Wrong: `db.Query()` in `api/posts.go`
- Right: Delegate to logic layer, which calls store layer
- Maintains separation of concerns and enables testing

### Rule 6 — Context Propagation
**Pass `context.Context` through all layers.**
- Use `r.Context()` from HTTP requests
- Enables request tracing, cancellation, and timeouts
- Required for proper database operations

### Rule 7 — Structured Logging
**Use the `logging.Logger` interface from `common/logging`. Never use `log.Printf` or `fmt.Print` in application code.**

The logger is constructed once in `main.go` from `cfg.Logging.Level` and passed through every handler constructor. Level is controlled via `config.yaml` under `logging.level` (debug | info | warn | error), defaulting to `info`.

**What to log and where:**
- `INFO` — startup lifecycle only (server start, migrations). In `main.go`.
- `ERROR` — internal failures that produce a 500 response. In the API layer via `handleLogicError`.
- `WARN` / `DEBUG` — reserved; use sparingly and only for genuinely invisible, non-user-facing events.

**What must never be logged:**
- Validation errors (400) — they are returned to the user; logging them is noise.
- Not found (404) and unauthorized (401) — expected, silent.
- Passwords, tokens, or any personal data.

### Rule 8 — errors.go Per Package
**Each package defines its own `errors.go` for error types specific to that layer.**
- `logic/errors.go` — `ValidationError` and helpers (`requireNonEmpty`, `requirePositive`). Use `logic.IsValidationError` in the API layer to distinguish 400 from 500 responses.
- `api/errors.go` — `handleLogicError` helper that maps `ValidationError` → 400, anything else → 500.
- Add new error types to the `errors.go` of the layer that owns them. Never scatter error definitions across business logic files.

### Rule 9 — Consistent JSON Response Format
**All responses follow the same structure.**
- Success: Direct JSON object or array
- Error: `{"error": "message"}`
- IDs returned as: `{"id": "uuid"}`
- Use `writeJSON()` and `writeError()` helpers

### Rule 10 — Environment-Specific Configuration
**Configuration is loaded from YAML, overridden by environment variables.**
- Default config in `config.yaml`
- Environment overrides for deployment flexibility
- Never hardcode sensitive values

### Rule 11 — Migration-Based Schema Evolution
**Database schema changes via SQL migrations.**
- Versioned migration files in `db/migrations/`
- Run migrations on startup with `db.RunMigrations()`
- Ensures consistent schema across environments

### Rule 12 — Store Is Never Called Directly from `main`
**`main.go` must never import or call the `store` layer for business operations.**
- Wrong: `store.UpsertAdminUser(db, password)` or `store.GetOrCreateJWTSecret(db)` in `main.go`
- Right: Delegate to a `logic.SetupService` which wraps the relevant store interfaces
- The `db/` package (Open, RunMigrations) is infrastructure and may be called from `main.go` directly — it is not a business logic store.
- Any future bootstrapping operations (seeding data, secrets rotation) must go through a logic-layer service.

---

## Layers of the Architecture

### API Layer (`api/*.go`)
**Purpose:** HTTP request/response handling, JSON serialization, routing.

**Responsibilities:**
- Parse HTTP requests into Go structs
- Validate basic request structure (JSON parsing)
- Call logic layer methods
- Serialize responses to JSON
- Handle HTTP status codes and headers

**Key Components:**
- `PostHandler`, `LinkHandler`, `AuthHandler`, `UploadHandler` interfaces
- `CORSMiddleware` for cross-origin requests
- `JWTAuthMiddleware` and `RateLimitMiddleware` in `api/middleware.go`
- Route registration in top-level `routes.go` (package main) — **not** inside `api/`

**Never Does:** Business logic, database queries, complex validation.

### Logic Layer (`logic/*.go`)
**Purpose:** Business rules, input validation, coordination between operations.

**Responsibilities:**
- Validate business rules (e.g., slug uniqueness, tag limits)
- Coordinate multi-step operations
- Transform data between layers
- Handle business logic errors

**Key Components:**
- `PostLogic`, `LinkLogic`, `AuthLogic`, `UploadLogic`, `SetupLogic` interfaces
- `PostService`, `LinkService`, `AuthService`, `UploadService`, `SetupService` implementations
- Input validation and sanitization

**Never Does:** HTTP concerns, direct database access.

### Store Layer (`store/*.go`)
**Purpose:** Data persistence and retrieval.

**Responsibilities:**
- Execute SQL queries
- Map database rows to Go structs
- Handle database transactions
- Manage connections and migrations

**Key Components:**
- `PostStore`, `LinkStore`, `UserStore`, `SettingsStore` interfaces
- `FileStore` interface with `LocalFileStore` and `R2FileStore` implementations
- SQLite implementations with prepared statements

**Never Does:** Business logic, HTTP responses.

---

## What Claude Code Must Never Do

### Guardrails — Absolute Prohibitions

1. **Never Skip Input Validation** — Every parameter from users must be checked. No exceptions.
2. **Never Return Sensitive Data** — Passwords, tokens, or internal IDs in error messages.
3. **Never Use String Formatting for SQL** — Always use prepared statements (`?` placeholders) to prevent SQL injection. This applies to every query in every store file. Never concatenate or `fmt.Sprintf` user input into a SQL string.
4. **Never Log Sensitive Information** — Passwords, JWT secrets, or user data in logs.
5. **Never Bypass Authentication** — All write operations must check JWT tokens.
6. **Never Hardcode Secrets** — Use config files and environment variables.
7. **Never Ignore Errors** — Every error must be handled appropriately.
8. **Never Mix Layers** — API layer calls logic, logic calls store. No shortcuts.
9. **Never Use Panic** — Return errors instead of panicking in production code.
10. **Never Trust Client-Side Validation** — Server must validate everything again.
11. **Never Expose Storage Backend URLs** — File access always goes through `/files/{fileId}` and `/thumbnails/{thumbnailId}`. R2 object URLs must never appear in any response, log, or error message.
12. **Never add a generic upload endpoint** — File uploads belong inside post handlers as `multipart/form-data` fields. Do not create `/uploads/*` routes. _Exception:_ per-resource sub-collection endpoints such as `POST /api/posts/{id}/images` are acceptable when they manage images on an already-existing post (e.g. adding an image to an existing photo post).
13. **Always rate-limit public POST endpoints** — Any `POST` route that writes to the database and requires no JWT must be wrapped with `RateLimitMiddleware` in the top-level `routes.go`. Use `rate.Every(time.Hour/N)` with a matching burst. Currently: `POST /api/subscribe` is rate-limited to 3 req/hour per IP. New public write endpoints must follow the same pattern.
14. **A photo post must always contain at least 1 image** — `DELETE /api/posts/{id}/images/{imageID}` returns a 400 ValidationError if the post has only one image remaining. This is enforced in the logic layer (`PostService.DeletePhotoImage`). Never bypass this check.

---

## Security Patterns Employed

### Authentication & Authorization
- **JWT Bearer Tokens:** Stateless authentication with 7-day expiration
- **Password Hashing:** bcrypt with appropriate cost factor
- **Token Storage:** Database-backed token validation (not just signature)
- **Middleware Protection:** All write endpoints wrapped with JWT validation

### Input Security
- **Parameter Validation:** Strict type, range, and format checking
- **SQL Injection Prevention:** Prepared statements only
- **XSS Prevention:** No direct HTML output (JSON API only)
- **CSRF Protection:** Stateless JWT doesn't require CSRF tokens

### Data Protection
- **No Sensitive Data in Logs:** Structured logging without secrets
- **Environment-Based Config:** Secrets via environment variables
- **SQLite Encryption:** Consider SQLCipher for production if needed

### Network Security
- **CORS Configuration:** Permissive in dev, locked down in production
- **HTTPS Enforcement:** Required for production deployments
- **Rate Limiting:** Per-IP rate limiting via `RateLimitMiddleware` (`api/middleware.go`) — required on all public POST endpoints that write to the database

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/posts` | No | All posts, filtered by `type` and `tag` query params |
| GET | `/api/posts/{slug}` | No | Single essay by slug |
| POST | `/api/posts/essay` | JWT | Create new essay |
| PUT | `/api/posts/essay/{id}` | JWT | Update existing essay |
| POST | `/api/posts/short` | JWT | Create new short post |
| PUT | `/api/posts/short/{id}` | JWT | Update existing short post |
| POST | `/api/posts/music` | JWT | Create new music post |
| PUT | `/api/posts/music/{id}` | JWT | Update existing music post |
| POST | `/api/posts/photo` | JWT | Create new photo post |
| PUT | `/api/posts/photo/{id}` | JWT | Update existing photo post |
| POST | `/api/posts/video` | JWT | Create new video post |
| PUT | `/api/posts/video/{id}` | JWT | Update existing video post |
| POST | `/api/posts/link` | JWT | Create new link post |
| PUT | `/api/posts/link/{id}` | JWT | Update existing link post |
| DELETE | `/api/posts/{id}` | JWT | Delete any post type |
| GET | `/api/links/featured` | No | Up to 4 featured external links |
| GET | `/api/links` | No | All external links |
| POST | `/api/links` | JWT | Create new external link |
| PUT | `/api/links/{id}` | JWT | Update existing external link |
| DELETE | `/api/links/{id}` | JWT | Delete external link |
| POST | `/api/auth/login` | No | Login with username/password, returns JWT |
| GET | `/api/openapi.yaml` | No | OpenAPI 3.0 specification |
| GET | `/api/index.html` | No | Swagger UI for API documentation |
| GET | `/files/{fileId}` | No | Stream original uploaded file binary |
| GET | `/thumbnails/{thumbnailId}` | No | Stream photo thumbnail binary |

Most API endpoints return `application/json`. Exceptions:
- `GET /files/{fileId}` streams original audio/image binary content with the proper MIME type.
- `GET /thumbnails/{thumbnailId}` streams binary thumbnail images.
Errors: `{"error":"message"}`.

### File Upload Contract
File uploads are embedded in the existing post create endpoints as `multipart/form-data` (not separate upload endpoints).

| Post type | Endpoint | File fields |
|-----------|----------|-------------|
| Music | `POST/PUT /api/posts/music` | `audioFile` (required), `albumArtFile` (optional) |
| Photo | `POST/PUT /api/posts/photo` | `images[]` (one or more, required) |
| Video | `POST/PUT /api/posts/video` | `thumbnailFile` (optional) |
| Link | `POST/PUT /api/posts/link` | `thumbnailFile` (optional) |

- Allowed photo types: `jpg`, `jpeg`, `png`, `webp`. Allowed audio types: `mp3`, `wav`, `ogg`, `m4a`.
- `POST /api/posts/music` supports optional `albumArtFile`.
- `POST /api/posts/video` and `POST /api/posts/link` support optional `thumbnailFile`.
- Link URLs must be absolute; the server derives `domain` automatically from the submitted `url`.
- Photo uploads auto-generate 4 thumbnail variants: tiny (20px), small (300px), medium (600px), large (1200px) wide. Widths are defined as constants in `logic/uploads.go` — change them there.
- Stored file keys: `photos/<postId>-<index>.<ext>`, `music/<postId>.<ext>`, `thumbnails/<postId>-<index>-<size>.<ext>`.
- All URL fields in responses are **relative** (`/files/…`, `/thumbnails/…`). The frontend must prefix them with the API base URL.
- File URLs always route through the backend's own `/files/` and `/thumbnails/` endpoints — the storage backend (local or R2) is never exposed.
- **Tag handling in multipart requests:** Use the `formTags(r)` helper (defined in `api/helpers.go`) in all multipart handlers — both CREATE and UPDATE. It reads either `tags` or `tags[]` from the form, splits comma-separated values, and returns:
  - `nil` — key was absent in a non-multipart request → preserve existing tags (UPDATE only)
  - `[]string{}` — key was absent in a multipart request, or present but empty → **clear all tags**
  - `[]string{…}` — the parsed tag values
  Never read `r.Form["tags"]` or `r.Form["tags[]"]` directly in handlers; always call `formTags(r)`.

### Post Types Supported
- **essay:** Long-form writing with title, slug, excerpt, body, reading time
- **short:** Brief text posts with hashtags
- **music:** Music tracks with album art, audio URL, duration
- **photo:** Photo galleries with metadata
- **video:** Videos with thumbnails and metadata
- **link:** Curated external links as posts

---

## Conventions

### Input Type Naming
Input types (passed into logic/store layers) follow the `*PostInput` pattern — never `Create*Post`.

| Input type | Used for |
|-----------|----------|
| `EssayPostInput`, `ShortPostInput` | create **and** the merged value inside update |
| `MusicPostInput`, `VideoPostInput`, `LinkPostInput`, `PhotoPostInput` | same |
| `Update*Post` (e.g. `UpdateEssayPost`) | **partial** update inputs only — all fields are pointers; `nil` means "leave unchanged" |

Rule: if a field is required on create but optional on update, the Create type uses a value type and the Update type uses a pointer. Never add pointer fields to `*PostInput` types just to support partial updates — keep a separate `Update*Post` for that.

### Post Type Constants
Use `models.PostType*` constants instead of bare string literals everywhere:
```go
models.PostTypeEssay, models.PostTypeShort, models.PostTypeMusic,
models.PostTypePhoto, models.PostTypeVideo, models.PostTypeLink
```

### Concurrency
- `UploadPhotoFiles` processes images concurrently, bounded by `maxConcurrentUploads = 4` (defined in `logic/uploads.go`). Do not raise this without considering memory and file-store rate limits.
- `DeletePost` deletes the database record first (post immediately invisible to readers), then cleans up files in a background goroutine using `context.Background()`. File-delete failures are logged but do not affect the response — this is intentional.

---

## HTTP Status Code Contract

Every handler must return the appropriate status using these rules. Do not deviate.

| Status | Meaning | When to use |
|--------|---------|-------------|
| `200 OK` | Success with body | GET, successful POST that returns data |
| `201 Created` | Resource created | POST that creates a new resource (returns `{"id": "…"}`) |
| `204 No Content` | Success, no body | DELETE |
| `400 Bad Request` | Client error | `logic.ValidationError` — bad input, failed validation, duplicate entry |
| `401 Unauthorized` | Auth failure | Missing/invalid JWT token |
| `404 Not Found` | Missing resource | Resource with given ID/slug does not exist |
| `429 Too Many Requests` | Rate limited | IP exceeded the per-endpoint rate limit |
| `500 Internal Server Error` | Server error | Any non-ValidationError from the logic or store layer |

The `handleLogicError` helper in `api/errors.go` maps `ValidationError → 400` and everything else `→ 500`. Use it for all logic layer errors. Never manually write a 500 — let the helper do it so the log entry is consistent.

---

## Transaction Pattern

Multi-step database writes that must be atomic (e.g. insert a post row then insert its tags) use `*sql.Tx`.
Never split a logically atomic operation across two separate store calls from the logic layer — if one succeeds and the other fails, the database will be in a partial state.

---

## Database Schema

SQLite database with foreign key constraints and cascading deletes.

### Core Tables
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
| `users` | User accounts for authentication |
| `settings` | Key-value configuration storage |

`photo_images` has nullable thumbnail columns added across migrations `004`–`005`: `thumbnail_tiny_url`, `thumbnail_small_url`, `thumbnail_medium_url`, `thumbnail_large_url`. Existing rows may have `NULL`; rows created via the upload endpoint have all four populated.

### Key Relationships
- All post types reference `posts.id` with CASCADE DELETE
- Tags stored separately to allow efficient filtering
- Users have optional current JWT token with expiration

---

## Development Workflow

1. **Start with Tests:** Write failing tests first
2. **Validate Inputs:** Add validation in logic layer
3. **Handle Errors:** Return descriptive error messages
4. **Log Appropriately:** Context-rich logs for debugging
5. **Test Endpoints:** Use OpenAPI spec and Swagger UI
6. **Migrate Schema:** Add SQL migrations for schema changes

## Code Quality — Required Before Every Change

After every code change, run both checks and fix any errors before considering the work done:

```bash
CGO_ENABLED=0 go vet ./...
CGO_ENABLED=0 golangci-lint run
```

Both commands must produce **zero errors**. Warnings from `golangci-lint` that are false positives may be suppressed with a targeted `//nolint:<linter>` comment with a brief explanation, but errors must be fixed.

---

## Deployment Considerations

- **Environment Variables:** `ANTIPRATIK_HOST`, `ANTIPRATIK_PORT` for configuration
- **Database Path:** Configurable SQLite file location
- **Static Serving:** Serves Next.js build from configurable directory
- **CORS:** Must be configured for production domains
- **HTTPS:** Required for secure cookie/token handling
- **Backups:** SQLite database should be regularly backed up
- **File Storage:** Configure `storage.backend` in `config.yaml`:
  - `local` — files written to `storage.local_dir` (default `./data/uploads/`)
  - `r2` — files stored in Cloudflare R2; requires `storage.r2.endpoint`, `bucket`, `access_key_id`, `secret_access_key`
  - R2 credentials should be supplied via environment or secrets management, not committed to `config.yaml`

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