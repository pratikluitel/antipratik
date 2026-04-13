# Deployment — Agent Reference

This file governs how CI/CD works in this repo. Read it before touching any workflow or infrastructure file.

---

## GitOps Flow

| Branch   | Environment | Image tag | Deploy dir              |
|----------|-------------|-----------|-------------------------|
| `master` | dev         | `:dev`    | `/opt/antipratik-dev/`  |
| `prod`   | prod        | `:prod`   | `/opt/antipratik-prod/` |

- `master` and `prod` are the only branches that trigger deployment.
- Images are pushed to GHCR (`ghcr.io/pratikluitel/antipratik-*`) and pulled by the server on deploy.
- The `prod` branch is the production gate. Never push breaking changes there without validating on `master` first.

---

## Secrets — Hard Rules

- **Never hardcode secrets** in workflow files, config files, scripts, or compose files.
- All secrets live in GitHub Secrets and are injected at runtime.
- On the server, secrets land in `/opt/antipratik-{TAG}/.env` via `printf` — that file is `chmod 600` immediately after writing.
- Config files in `infrastructure/config/` store structure only (host, port, log level, storage backend). Credentials are empty placeholders filled by environment variables at runtime.
- Never commit `.env` files. They are server-side only and not in the repo.
- `GHCR_TOKEN` is a PAT used for `docker login` inside `deploy.sh`. Do not inline it in workflow `run:` blocks — always pass it as an env var to the script.

---

## Workflow Architecture

### Job DAG (`deployment.yaml`)

```
sidequest-changes  ──────────────────────────────── deploy-sidequests (master only)
       │
       ├── build-api   (if app/antipratik-api/** changed, any branch) ──┐
       ├── build-ui    (if app/antipratik-ui/** changed, any branch)  ──┤── deploy-app (master/prod only)
       └── build-nginx (always, any branch)                           ──┘
```

- **Build jobs run on every push to any branch** when relevant paths change — on feature branches they build only (no push) to verify the image compiles. Images are pushed to GHCR only on `master` (`:dev`) and `prod` (`:prod`).
- **`deploy-app` only runs on `master`/`prod`**, even if all build jobs were skipped. It re-deploys whatever images are current in GHCR, which picks up config/compose changes too.
- **`deploy-sidequests`** only runs when `sidequests/**` changed and the branch is `master`.
- Build jobs for API and UI are path-filtered to avoid rebuilding unchanged services.
- Build jobs are fully parallel — never collapse them back into a single job.

### Path filters (`sidequest-changes` job)

```yaml
api:       app/antipratik-api/**, docker/Dockerfile.api
ui:        app/antipratik-ui/**, docker/Dockerfile.ui
sidequests: sidequests/**
```

If you add a new service, add its filter here and a corresponding `build-*` job.

### `deploy-app` condition

```yaml
if: always() && !contains(needs.*.result, 'failure') && !contains(needs.*.result, 'cancelled')
```

This runs even when build jobs are skipped, but aborts if any build job failed. Do not change this pattern without understanding the GitHub Actions skipped-vs-failed distinction.

---

## Deploy Script (`deploy.sh`)

`.github/scripts/deploy.sh` runs on the server via SSH. It:
1. Copies compose file and config to `/opt/antipratik-{TAG}/`
2. Logs into GHCR using `$GHCR_TOKEN` / `$GHCR_USER`
3. Pulls latest images
4. Runs `docker compose up -d --remove-orphans`
5. Prunes dangling images

Do not add secrets or environment-specific logic to this script — keep it generic and tag-driven.

---

## Infrastructure Layout

```
infrastructure/
  config/
    config.dev.yaml   # API config for dev (no secrets)
    config.prod.yaml  # API config for prod (no secrets)
  compose/
    docker-compose.dev.yml   # Pulls :dev images from GHCR
    docker-compose.prod.yml  # Pulls :prod images from GHCR
  nginx.conf                 # Nginx config template (uses ${SERVER_NAME})

docker/
  Dockerfile.api    # Go multi-stage build
  Dockerfile.ui     # Next.js multi-stage build
  Dockerfile.nginx  # Nginx + envsubst

.github/
  scripts/
    deploy.sh       # Server-side deploy entrypoint
  workflows/
    deployment.yaml         # Main CI/CD (builds + deploy)
    build-check-api.yaml    # Go build check on PR/push
    build-check.yaml        # Next.js build check on PR/push
```

---

## SCP Step Convention

All files are copied to the server in a **single SCP step** using comma-separated sources with `strip_components: 2`. All sources must be exactly 2 path components deep (e.g. `infrastructure/config/file`, `infrastructure/compose/file`, `.github/scripts/file`). Do not add a fourth SCP step — extend the existing `source:` list.

---

## What NOT to Do

- Do not add new sequential build steps inside a single job — add a new parallel job instead.
- Do not skip path filters on build jobs unless the service has no independent change surface.
- Do not add a separate SCP step for each new file — extend the existing combined step.
- Do not reference a job output using a wrong job name (e.g. `needs.changes.outputs.x` when the job is named `sidequest-changes`).
- Do not push to `prod` branch via force-push or bypass the normal PR flow.
- Do not store secrets in `infrastructure/config/*.yaml` — these files are committed to the repo.
