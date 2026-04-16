# Infrastructure — Agent Reference

This file governs how CI/CD, deployment, and infrastructure work in this repo. Read it before touching any workflow, compose, nginx, or infrastructure file.

---

## GitOps Flow

| Branch   | Environment | Image tag | Deploy dir              | Ports            |
|----------|-------------|-----------|-------------------------|------------------|
| `master` | dev         | `:dev`    | `/opt/antipratik-dev/`  | 8080 (HTTP), 8443 (HTTPS) |
| `prod`   | prod        | `:prod`   | `/opt/antipratik-prod/` | 80 (HTTP), 443 (HTTPS) |

- `master` and `prod` are the only branches that trigger deployment.
- Images are pushed to GHCR (`ghcr.io/pratikluitel/antipratik-*`) and pulled by the server on deploy.
- The `prod` branch is the production gate. Never push breaking changes there without validating on `master` first.

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
  scripts/                   # (empty — deploy.sh moved to infrastructure/)
  workflows/
    build.yaml               # Build, lint, and image push (triggered on every push)
    deploy.yaml              # Deploy to server (triggered on successful Build run on master/prod)
```

---

## Docker Compose

Each environment has its own compose file under `infrastructure/compose/`. Both follow the same structure:

### Services

| Service | Image | Purpose |
|---------|-------|---------|
| `nginx` | `ghcr.io/pratikluitel/antipratik-nginx:{tag}` | Reverse proxy, TLS termination |
| `api`   | `ghcr.io/pratikluitel/antipratik-api:{tag}`   | Go backend |
| `ui`    | `ghcr.io/pratikluitel/antipratik-ui:{tag}`    | Next.js frontend (SSR) |

### Network
All services share a single internal Docker network (`antipratik-dev` / `antipratik-prod`). Only nginx is exposed externally.

### Volumes (server-side paths)
- **SSL certs**: `/opt/antipratik-{TAG}/certs/fullchain.pem` and `privkey.pem` → mounted read-only into nginx at `/etc/ssl/certs/` and `/etc/ssl/private/`
- **API config**: `/opt/antipratik-{TAG}/config.yaml` → mounted into the API container at `/root/config.yaml`
- **API data**: `/opt/antipratik-{TAG}/data` → mounted at `/data` (SQLite database lives here)

### Environment variables
- `nginx`: `SERVER_NAME` (injected from `.env` on server)
- `api`: loaded from `.env` file on server (contains R2 credentials, admin password, etc.)
- `ui`: `SERVER_API_URL=http://api:8080` (internal network address, hardcoded)

---

## Nginx

`infrastructure/nginx.conf` is a template — `${SERVER_NAME}` is substituted by `envsubst` at container startup via `Dockerfile.nginx`.

### Routing

| Path prefix | Upstream |
|-------------|----------|
| `/api/`     | `http://api:8080` |
| `/files/`   | `http://api:8080` |
| `/thumbnails/` | `http://api:8080` |
| `/`         | `http://ui:3000` |

### TLS

- HTTP (port 80) always redirects to HTTPS (301).
- TLS cert/key are bind-mounted from the server at runtime (not baked into the image).
- Protocols: TLSv1.2 and TLSv1.3 only.
- HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, X-XSS-Protection headers set on all responses.
- `client_max_body_size 50M` to support media uploads.

---

## SSL Certificates

`.ssl/generate.sh` generates a self-signed certificate for a given environment:

```bash
./infrastructure/.ssl/generate.sh dev   # generates fullchain-dev.pem + privkey-dev.pem
./infrastructure/.ssl/generate.sh prod  # generates fullchain-prod.pem + privkey-prod.pem
```

- Uses `openssl req -x509 -newkey rsa:4096 -sha256 -days 365`
- SAN includes `DNS:localhost` and `IP:127.0.0.1`
- Generated `.pem` files are `.gitignore`'d — never commit them
- After generation, the script prints base64-encoded cert and key for storing in GitHub Secrets
- On the server, certs are placed at `/opt/antipratik-{TAG}/certs/` and bind-mounted into nginx

---

## API Config (`infrastructure/config/`)

Config files are committed to the repo and contain **structure only** — no credentials. Runtime secrets (R2 keys, admin password) are empty strings filled by environment variables at runtime.

Fields:
- `server.host` / `server.port`: listen address (always `0.0.0.0:8080`)
- `db.path`: SQLite path (inside the `/data` volume)
- `admin_password`: placeholder — set via `.env` on server
- `static.dir`: empty (static files served from R2, not local disk)
- `storage.backend`: `"r2"`; `r2.bucket`, `r2.endpoint`, `r2.access_key_id`, `r2.secret_access_key` are empty placeholders
- `logging.level`: `"debug"` for dev, `"info"` for prod

---

## Secrets — Hard Rules

- **Never hardcode secrets** in workflow files, config files, scripts, or compose files.
- All secrets live in GitHub Secrets and are injected at runtime.
- On the server, secrets land in `/opt/antipratik-{TAG}/.env` via `printf` — that file is `chmod 600` immediately after writing.
- `GHCR_TOKEN` is a PAT used for `docker login` inside `deploy.sh`. Do not inline it in workflow `run:` blocks — always pass it as an env var to the script.
- SSL certs are stored as base64-encoded GitHub Secrets (generated by `.ssl/generate.sh`) and written to the server during the deploy workflow.

---

## CI/CD Workflows

There are two workflow files:

### `build.yaml` — Build & Lint (triggered on every push)

```
changes (dorny/paths-filter, base: master)
    │
    ├── lint-api        (if api changed OR master/prod) ──── build-api-image (push to GHCR on master/prod)
    ├── lint-ui         (if ui changed OR master/prod)  ──── build-ui-image  (push to GHCR on master/prod)
    └── build-nginx     (always)                             (push to GHCR on master/prod)
```

- **`changes`** uses `dorny/paths-filter` with `base: master` to detect which services changed.
- **`lint-api`** runs `go vet` and `golangci-lint` (CGO_ENABLED=0). Runs if API files changed, or always on `master`/`prod`.
- **`lint-ui`** runs `npm run lint` and `npm run typecheck`. Runs if UI files changed, or always on `master`/`prod`.
- **`build-api-image`** / **`build-ui-image`** depend on their respective lint jobs passing. They build the Docker image on every branch; push to GHCR only on `master`/`prod`.
- **`build-nginx`** always runs on every push; pushes to GHCR only on `master`/`prod`.
- Images are tagged `:dev` on `master` and `:prod` on `prod`.

### `deploy.yaml` — Deploy (triggered on successful `Build` workflow run on `master`/`prod`)

Steps in `deploy-app`:
1. Write secrets to server (SSH) — creates `/opt/antipratik-{TAG}/.env`, `chmod 600`
2. Write TLS certs to server (SSH) — decodes base64 secrets into `/opt/antipratik-{TAG}/certs/`
3. SCP files to `/tmp/` — config, compose file, `deploy.sh`
4. SSH and run `deploy.sh`

`deploy-sidequests` also runs when the `Build` workflow completes on `master` and `sidequests/**` changed — publishes to Cloudflare Pages.

---

## Deploy Script (`infrastructure/deploy.sh`)

Runs on the server via SSH. Steps:
1. Copies compose file and config to `/opt/antipratik-{TAG}/`
2. Logs into GHCR using `$GHCR_TOKEN` / `$GHCR_USER`
3. Pulls latest images (`docker compose pull`)
4. Starts services (`docker compose up -d --remove-orphans`)
5. Prunes dangling images (`docker image prune -f`)

Do not add secrets or environment-specific logic to this script — keep it generic and tag-driven.

---

## SCP Step Convention

All files are copied to the server in a **single SCP step** using comma-separated sources with `strip_components: 2`. All sources must be exactly 2 path components deep (e.g. `infrastructure/config/file`, `infrastructure/compose/file`, `infrastructure/deploy.sh`). Do not add a fourth SCP step — extend the existing `source:` list.

---

## What NOT to Do

- Do not add new sequential build steps inside a single job — add a new parallel job instead.
- Do not skip path filters on build jobs unless the service has no independent change surface.
- Do not add a separate SCP step for each new file — extend the existing combined step.
- Do not reference a job output using a wrong job name (e.g. `needs.filter.outputs.x` when the job is named `changes`).
- Do not push to `prod` branch via force-push or bypass the normal PR flow.
- Do not store secrets in `infrastructure/config/*.yaml` — these files are committed to the repo.
- Do not bake SSL certs into Docker images — always bind-mount them from the server at runtime.
- Do not change `ui`'s `SERVER_API_URL` to an external URL — it must use the internal Docker network address.
