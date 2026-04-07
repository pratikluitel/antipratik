# antipratik-api — Claude Code Reference

Read this file at the start of every session before writing any code.

---

## Project Overview & Philosophy

**antipratik-api** is the backend REST API for antipratik.com, a personal brand site for a Kathmandu-based developer, musician, writer, and photographer.

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
**Log errors with context, but never sensitive data.**
- Wrong: `log.Printf("error: %v", err)`
- Right: `log.Printf("CreateEssay error: %v", err)` with operation context
- Never log passwords, tokens, or personal data

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
- Versioned migration files in `store/migrations/`
- Run migrations on startup with `store.RunMigrations()`
- Ensures consistent schema across environments

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
- `PostHandler`, `LinkHandler`, `AuthHandler` interfaces
- `CORSMiddleware` for cross-origin requests
- `JWTAuthMiddleware` for authentication
- Route registration in `routes.go`

**Never Does:** Business logic, database queries, complex validation.

### Logic Layer (`logic/*.go`)
**Purpose:** Business rules, input validation, coordination between operations.

**Responsibilities:**
- Validate business rules (e.g., slug uniqueness, tag limits)
- Coordinate multi-step operations
- Transform data between layers
- Handle business logic errors

**Key Components:**
- `PostLogic`, `LinkLogic`, `AuthLogic` interfaces
- `PostService`, `LinkService`, `AuthService` implementations
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
- `PostStore`, `LinkStore`, `UserStore` interfaces
- SQLite implementations with prepared statements
- Migration execution on startup

**Never Does:** Business logic, HTTP responses.

---

## What Claude Code Must Never Do

### Guardrails — Absolute Prohibitions

1. **Never Skip Input Validation** — Every parameter from users must be checked. No exceptions.
2. **Never Return Sensitive Data** — Passwords, tokens, or internal IDs in error messages.
3. **Never Use String Formatting for SQL** — Always use prepared statements to prevent SQL injection.
4. **Never Log Sensitive Information** — Passwords, JWT secrets, or user data in logs.
5. **Never Bypass Authentication** — All write operations must check JWT tokens.
6. **Never Hardcode Secrets** — Use config files and environment variables.
7. **Never Ignore Errors** — Every error must be handled appropriately.
8. **Never Mix Layers** — API layer calls logic, logic calls store. No shortcuts.
9. **Never Use Panic** — Return errors instead of panicking in production code.
10. **Never Trust Client-Side Validation** — Server must validate everything again.

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
- **Rate Limiting:** Consider adding for production (not yet implemented)

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

All endpoints return `application/json`. Errors: `{"error":"message"}`.

### Post Types Supported
- **essay:** Long-form writing with title, slug, excerpt, body, reading time
- **short:** Brief text posts with hashtags
- **music:** Music tracks with album art, audio URL, duration
- **photo:** Photo galleries with metadata
- **video:** Videos with thumbnails and metadata
- **link:** Curated external links as posts

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

---

## Deployment Considerations

- **Environment Variables:** `ANTIPRATIK_HOST`, `ANTIPRATIK_PORT` for configuration
- **Database Path:** Configurable SQLite file location
- **Static Serving:** Serves Next.js build from configurable directory
- **CORS:** Must be configured for production domains
- **HTTPS:** Required for secure cookie/token handling
- **Backups:** SQLite database should be regularly backed up

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