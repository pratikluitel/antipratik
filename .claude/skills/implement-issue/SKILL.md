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

Extract: title, body, labels, and comments. Identify affected areas and any specific tokens, components, or files mentioned.

### 2. Explore the codebase

Launch an Explore agent (or read files directly) to understand the affected code:
- Always read `CLAUDE.md` first — it contains inviolable rules and known token deviations
- Read `src/styles/tokens.css` if the issue touches styling or tokens
- Read the relevant component `.module.css` and `.tsx` files
- Check `Checkpoints.md` for prior decisions that may be relevant

### 3. Plan before coding (use Plan mode)

Enter plan mode. Write a concise plan that covers:
- Root cause / what is wrong
- Exact files to change and what to change in each
- Verification steps

Exit plan mode only after the plan is approved.

### 4. Implement

Make the changes. Follow all rules in `CLAUDE.md`:
- No hardcoded hex/px values — use `var(--token)`
- No Tailwind, no direct `fetch()`, no inline styles
- All data through `src/lib/api.ts`
- Token changes in `tokens.css` only; never in component CSS

### 5. Update CLAUDE.md if needed

If the fix introduces or modifies tokens, updates a Sacred Rule, or records a design decision, update the relevant section of `CLAUDE.md`. Change only what is required.

### 6. Create a branch and commit locally

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

### 7. Open a Pull Request

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

## Test plan
- [ ] <thing to verify manually>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

Return the PR URL to the user when done.
