---
name: review-pr
description: Review a GitHub pull request for pratikluitel/antipratik using the gh CLI. Use when the user says "review PR #N", "check PR", or "review pull request". Accepts a PR number as argument.
---

Review pull request **#$ARGUMENTS** from the `pratikluitel/antipratik` repository and leave a comment with the findings.

## Workflow

### MAIN RULE
DO NOT UNDER ANY CIRCUMSTANCE USE THE GITHUB MCP TO PERFORM THE REVIEW. THIS IS TO ENSURE YOU FULLY ENGAGE WITH THE CHANGES AND UNDERSTAND THEM IN CONTEXT, RATHER THAN RELYING ON AUTOMATED TOOLS.

### 1. Fetch PR metadata and diff

```bash
# PR overview: title, body, source branch, changed files
gh pr view $ARGUMENTS --repo pratikluitel/antipratik \
  --json title,body,headRefName,baseRefName,files \
  --jq '{title, body, branch: .headRefName, files: [.files[].path]}'

# Full diff
gh pr diff $ARGUMENTS --repo pratikluitel/antipratik
```

If the diff is large, fetch per-file:
```bash
# List changed files only (cheap)
gh pr view $ARGUMENTS --repo pratikluitel/antipratik --json files --jq '[.files[].path]'

# Diff a specific file
gh pr diff $ARGUMENTS --repo pratikluitel/antipratik -- path/to/file
```

### 2. Determine scope from changed file paths

Inspect the list of changed files — this is the authoritative scope signal:
- Any path under `app/antipratik-ui/` → **FE**
- Any path under `app/antipratik-api/` → **BE**
- Both present → **FE + BE**

Record the scope — it determines which CLAUDE.md(s) to read and which rules to enforce.

### 3. Read relevant source files locally

**Always read the relevant CLAUDE.md(s) — these are the rules every change must follow.**

**If FE (or both):**
- Read `app/antipratik-ui/CLAUDE.md`
- Read `app/antipratik-ui/src/styles/tokens.css` for any styling changes (verify full token hierarchy)

**If BE (or both):**
- Read `app/antipratik-api/CLAUDE.md`
- Read local copies of changed files under `app/antipratik-api/` for full context

### 4. Review against these criteria

**Correctness**
- Does the change fix or implement what the PR/linked issue describes?
- Are new values/logic sound?

**Completeness**
- Does the PR address everything in the linked issue?
- Are there related cases that were missed?

**Side effects**
- Do changes unintentionally affect other components, tokens, or endpoints?

---

**FE CLAUDE.md compliance** *(check only if scope includes FE):*
- No hardcoded hex/px values — use `var(--token)` (Rule 1)
- No Tailwind classes (Rule 3)
- No direct `fetch()` in components — all data through `src/lib/api.ts` (Rule 4)
- No `--accent-*` tokens on UI chrome (Rule 6)
- Text contrast ≥ 4.5:1 normal / ≥ 3:1 large (Rule 16)
- Typography not crossed: serif = content, sans = interface (Rule 9)
- If tokens were changed: is the Token Additions & Modifications table in `CLAUDE.md` updated?

**BE CLAUDE.md compliance** *(check only if scope includes BE):*
- All input parameters validated in the logic layer before reaching the store (Rule 1 & 2)
- Validation errors return descriptive messages via `ValidationError`; handlers use `logic.IsValidationError` for 400 vs 500 (Rule 3)
- JWT middleware on all POST/PUT/DELETE endpoints (Rule 4)
- No database access in the API layer — must delegate to logic → store (Rule 5)
- `context.Context` passed through all layers (Rule 6)
- Errors logged with operation context; no passwords/tokens in logs (Rule 7)
- New error types defined in `errors.go` of the owning package (Rule 8)
- No `panic` — errors returned instead (Guardrail 9)
- No ignored errors (Guardrail 7)

---

**CLAUDE.md hygiene**
- If a new pattern was introduced, should a rule be added to the relevant `CLAUDE.md`?

**Out-of-scope findings**
- Note pre-existing issues spotted in the diff that are NOT part of this PR — flag as candidates for follow-up, not as blockers.

### 5. Form a verdict

One of:
- **Ready to merge** — no issues found
- **Minor notes** — small things to be aware of, not blocking
- **Needs changes** — something incorrect or non-compliant that must be fixed before merging

### 6. Post the review comment

```bash
gh pr comment $ARGUMENTS --repo pratikluitel/antipratik --body "$(cat <<'EOF'
<verdict line with emoji: ✅ ready / ⚠️ minor notes / ❌ needs changes>

**Scope:** FE / BE / Both

**<filename>** — <one-line summary of changes and assessment>
**<filename>** — <one-line summary of changes and assessment>

<If any issues: numbered list of specific findings with file:line references>

<If any out-of-scope findings: flag them clearly as candidates for a follow-up issue>
EOF
)"
```

Keep the comment concise — full picture in under a minute.
