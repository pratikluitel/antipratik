# Infrastructure — Agent Reference

## GitOps Flow

| Branch   | Environment | Image tag | Deploy dir              | Ports            |
|----------|-------------|-----------|-------------------------|------------------|
| `master` | dev         | `:dev`    | `/opt/antipratik-dev/`  | 8080 (HTTP), 8443 (HTTPS) |
| `prod`   | prod        | `:prod`   | `/opt/antipratik-prod/` | 80 (HTTP), 443 (HTTPS) |

Only `master` and `prod` trigger deployment. Never push breaking changes to `prod` without validating on `master` first.

---

## Directory Layout

```
infrastructure/
  compose/
    docker-compose.dev.yml   # Pulls :dev images from GHCR
    docker-compose.prod.yml  # Pulls :prod images from GHCR
  config/
    config.dev.yaml          # API config for dev (no secrets — credentials are empty placeholders)
    config.prod.yaml         # API config for prod (no secrets — credentials are empty placeholders)
  nginx.conf                 # Nginx config template (uses ${SERVER_NAME} via envsubst)
  deploy.sh                  # Server-side deploy entrypoint (runs via SSH)
  .ssl/
    generate.sh              # Generates self-signed SSL certs for dev/prod

docker/
  Dockerfile.api             # Go multi-stage build (CGO_ENABLED=0, pure Go)
  Dockerfile.ui              # Next.js multi-stage build
  Dockerfile.nginx           # Nginx + envsubst for SERVER_NAME substitution

.github/
  workflows/
    build.yaml               # Build, lint, and image push (triggered on every push)
    deploy.yaml              # Deploy to server (triggered on successful Build run on master/prod)
```

---

## Docker Compose

| Service | Image | Purpose |
|---------|-------|---------|
| `nginx` | `ghcr.io/pratikluitel/antipratik-nginx:{tag}` | Reverse proxy, TLS termination |
| `api`   | `ghcr.io/pratikluitel/antipratik-api:{tag}`   | Go backend |
| `ui`    | `ghcr.io/pratikluitel/antipratik-ui:{tag}`    | Next.js frontend (SSR) |

All services share a single internal Docker network. Only nginx is exposed externally.

**Volumes (server-side paths):**
- SSL certs: `/opt/antipratik-{TAG}/certs/fullchain.pem` + `privkey.pem` → mounted read-only into nginx
- API config: `/opt/antipratik-{TAG}/config.yaml` → mounted at `/root/config.yaml`
- API data: `/opt/antipratik-{TAG}/data` → mounted at `/data` (SQLite lives here)

---

## Nginx

`${SERVER_NAME}` is substituted by `envsubst` at container startup.

| Path prefix | Upstream |
|-------------|----------|
| `/api/`     | `http://api:8080` |
| `/files/`   | `http://api:8080` |
| `/thumbnails/` | `http://api:8080` |
| `/`         | `http://ui:3000` |

- HTTP (port 80) always redirects to HTTPS (301).
- TLS: certs bind-mounted at runtime (never baked into image). TLSv1.2 and TLSv1.3 only.
- HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, X-XSS-Protection set on all responses.
- `client_max_body_size 50M` for media uploads.

---

## SSL Certificates

```bash
./infrastructure/.ssl/generate.sh dev   # generates fullchain-dev.pem + privkey-dev.pem
./infrastructure/.ssl/generate.sh prod  # generates fullchain-prod.pem + privkey-prod.pem
```

Generated `.pem` files are `.gitignore`'d — never commit them. After generation the script prints base64-encoded values for storing in GitHub Secrets; on the server certs go in `/opt/antipratik-{TAG}/certs/`.

---

## API Config (`infrastructure/config/`)

Structure only — no credentials. Fields: `server.host`/`server.port` (always `0.0.0.0:8080`), `db.path`, `admin_password` (placeholder), `static.dir` (empty), `storage.backend` (`"r2"`), R2 fields (empty placeholders), `logging.level` (`"debug"` dev / `"info"` prod).

---

## Secrets

- Never hardcode secrets in workflow files, config files, scripts, or compose files.
- All secrets live in GitHub Secrets and are injected at runtime.
- Server secrets land in `/opt/antipratik-{TAG}/.env` via `printf` — `chmod 600` immediately after.
- `GHCR_TOKEN` is a PAT for `docker login` inside `deploy.sh`. Always pass as env var, never inline in `run:` blocks.

---

## CI/CD Workflows

### `build.yaml` — triggered on every push

```
changes (dorny/paths-filter, base: master)
    │
    ├── lint-api        (if api changed OR master/prod) ──── build-api-image (push to GHCR on master/prod)
    ├── lint-ui         (if ui changed OR master/prod)  ──── build-ui-image  (push to GHCR on master/prod)
    └── build-nginx     (always)                             (push to GHCR on master/prod)
```

- `lint-api`: `go vet` + `golangci-lint` (CGO_ENABLED=0)
- `lint-ui`: `npm run lint` + `npm run typecheck`
- Images tagged `:dev` on `master`, `:prod` on `prod`.

### `deploy.yaml` — triggered on successful Build on `master`/`prod`

1. Write secrets to server (SSH) — `/opt/antipratik-{TAG}/.env`, `chmod 600`
2. Write TLS certs (SSH) — decode base64 → `/opt/antipratik-{TAG}/certs/`
3. SCP files to `/tmp/` — config, compose file, `deploy.sh`
4. SSH and run `deploy.sh`

---

## Deploy Script (`infrastructure/deploy.sh`)

1. Copies compose file and config to `/opt/antipratik-{TAG}/`
2. Logs into GHCR using `$GHCR_TOKEN` / `$GHCR_USER`
3. `docker compose pull`
4. `docker compose up -d --remove-orphans`
5. `docker image prune -f`

Keep this script generic and tag-driven — no secrets or environment-specific logic.

---

## SCP Step Convention

All files are copied in a **single SCP step** using comma-separated sources with `strip_components: 2`. All sources must be exactly 2 path components deep (`infrastructure/config/file`, `infrastructure/compose/file`, `infrastructure/deploy.sh`). Do not add a fourth SCP step — extend the existing `source:` list.

---

## What NOT to Do

- Do not add new sequential build steps inside a single job — add a new parallel job instead.
- Do not add a separate SCP step for each new file — extend the existing combined step.
- Do not reference a job output using the wrong job name (e.g. `needs.filter.outputs.x` when the job is named `changes`).
- Do not store secrets in `infrastructure/config/*.yaml` — these files are committed to the repo.
- Do not bake SSL certs into Docker images — always bind-mount from the server at runtime.
- Do not change `ui`'s `SERVER_API_URL` to an external URL — it must use the internal Docker network address (`http://api:8080`).
