---
name: implement-issue
description: Implement a GitHub issue for pratikluitel/antipratik using the gh CLI. Use when the user says "implement issue #N", "fix issue", or "work on issue". Accepts an issue number as argument.
---

Implement GitHub issue **#$ARGUMENTS** from the `pratikluitel/antipratik` repository.

### MAIN RULE
DO NOT UNDER ANY CIRCUMSTANCE USE THE GITHUB MCP. Use only the `gh` CLI and local file reads.

## Workflow

### 1. Read the issue

```bash
gh issue view $ARGUMENTS --repo pratikluitel/antipratik \
  --json title,body,labels,comments \
  --jq '{title, body, labels: [.labels[].name], comments: [.comments[].body]}'
```

Extract: title, body, labels, and comments. Identify affected areas and any specific tokens, components, layers, or files mentioned.

### 2. Determine scope: FE, BE, or both

**Check labels first** (authoritative):
- Labels containing `frontend`, `fe`, `ui` → **FE only**
- Labels containing `backend`, `be`, `api` → **BE only**
- Both present, or no matching labels → scan title + body for keywords:
  - FE signals: component, style, token, CSS, page, layout, UI, design, frontend validation
  - BE signals: endpoint, route, store, migration, database, SQL, Go, API, auth, api/backend validation
- If signals are mixed or absent → treat as **both**

Record the scope decision — it controls every subsequent step.

### 3. Explore the codebase

**Always read the relevant CLAUDE.md(s) first — they contain inviolable rules.**

**If FE (or both):**
- Read `app/antipratik-ui/CLAUDE.md`
- Read `app/antipratik-ui/src/styles/tokens.css` if the issue touches styling or tokens
- Read the relevant component `.module.css` and `.tsx` files under `app/antipratik-ui/src/`

**If BE (or both):**
- Read `app/antipratik-api/CLAUDE.md`
- Read the relevant files under `app/antipratik-api/` — `api/`, `logic/`, `store/`, `models/`

### 4. Plan before coding (use Plan mode)

Enter plan mode. Write a concise plan that covers:
- Root cause / what is wrong / what needs to change
- Exact files to change and what to change in each
- Which scope (FE/BE/both) is being touched and why
- Verification steps

Exit plan mode only after the plan is approved.

### 5. Implement

Make the changes. Apply **only the rules relevant to the scope**:

**FE rules (from `antipratik-ui/CLAUDE.md`):**
- No hardcoded hex/px values — use `var(--token)`
- No Tailwind, no direct `fetch()`, no inline styles
- All data through `src/lib/api.ts`
- Token changes in `tokens.css` only; never in component CSS

**BE rules (from `antipratik-api/CLAUDE.md`):**
- All input parameters must be validated in the logic layer before reaching the store
- Return descriptive `ValidationError` messages; use `logic.IsValidationError` in handlers for 400 vs 500
- JWT middleware on all write endpoints (POST/PUT/DELETE)
- No database access in the API layer — always delegate to logic → store
- Pass `context.Context` through all layers
- Log errors with operation context; never log passwords or tokens
- New error types belong in the `errors.go` of the owning package

### 6. Update CLAUDE.md if needed

If the fix introduces or modifies tokens (FE), updates a Sacred Rule, or records a design decision, update the relevant `CLAUDE.md`. Change only what is required.

### 7. Create a branch and commit locally

```bash
# Create and switch to the branch
git checkout -b <short-description>-$ARGUMENTS

# Stage and commit all changed files
git add <files>
git commit -m "$(cat <<'EOF'
<imperative summary> — closes #$ARGUMENTS

<2–3 sentence description of what was wrong and what was changed>

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
EOF
)"

# Push to remote
git push -u origin <short-description>-$ARGUMENTS
```

### 8. Open a Pull Request

```bash
gh pr create \
  --repo pratikluitel/antipratik \
  --base master \
  --head <short-description>-$ARGUMENTS \
  --title "<short imperative title ≤ 70 chars>" \
  --body "$(cat <<'EOF'
Closes #$ARGUMENTS

## Summary
- <bullet: what changed and why>

## Scope
FE / BE / Both — <one line explaining why>

## Test plan
- [ ] <thing to verify manually>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

Return the PR URL to the user when done.
