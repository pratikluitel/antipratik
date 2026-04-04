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

### 2. Read relevant source files locally

Read for full context:
- `CLAUDE.md` — rules every change must follow
- `src/styles/tokens.css` — for styling changes, verify the full token hierarchy

### 3. Review against these criteria

**Correctness**
- Does the change fix or implement what the PR/linked issue describes?
- Are new values/logic sound? (e.g. for colour tokens: verify contrast ratios)

**Completeness**
- Does the PR address everything in the linked issue?
- Are there related cases that were missed?

**Side effects**
- Do changes unintentionally affect other components or tokens?
- For token changes: trace which components consume the token and confirm none are broken

**CLAUDE.md compliance** — check every Sacred Rule that applies:
- No hardcoded hex/px values (Rule 1)
- No Tailwind classes (Rule 3)
- No direct `fetch()` in components (Rule 4)
- No `--accent-*` on UI chrome (Rule 6)
- Text contrast ≥ 4.5:1 normal / ≥ 3:1 large (Rule 16)
- Typography not crossed: serif = content, sans = interface (Rule 9)

**CLAUDE.md hygiene**
- If tokens were changed, is the Token Additions & Modifications table updated?
- If a new pattern was introduced, should a rule be added?

**Out-of-scope findings**
- Note pre-existing issues spotted in the diff that are NOT part of this PR — flag as candidates for follow-up, not as blockers.

### 4. Form a verdict

One of:
- **Ready to merge** — no issues found
- **Minor notes** — small things to be aware of, not blocking
- **Needs changes** — something incorrect or non-compliant that must be fixed before merging

### 5. Post the review comment

```bash
gh pr comment $ARGUMENTS --repo pratikluitel/antipratik --body "$(cat <<'EOF'
<verdict line with emoji: ✅ ready / ⚠️ minor notes / ❌ needs changes>

**<filename>** — <one-line summary of changes and assessment>
**<filename>** — <one-line summary of changes and assessment>

<If any issues: numbered list of specific findings with file:line references>

<If any out-of-scope findings: flag them clearly as candidates for a follow-up issue>
EOF
)"
```

Keep the comment concise — full picture in under a minute.
