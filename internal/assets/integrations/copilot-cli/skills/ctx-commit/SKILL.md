---
name: ctx-commit
description: "Commit with context persistence. Use instead of raw git commit to capture decisions and learnings alongside code changes."
tools: [bash, read, write, edit]
---

Commit code changes, then prompt for decisions and learnings
worth persisting. Bridges the gap between committing code and
recording the context behind it.

## When to Use

- When committing after meaningful work (feature, bugfix,
  refactor)
- When the commit involves a design choice or trade-off
- When the user says "commit" or "commit this"

## When NOT to Use

- For trivial commits (typo, formatting): just commit normally
- When the user explicitly says "just commit, no context"
- When nothing has changed

## Process

### 1. Pre-commit checks

Unless the user says "skip checks":

- Run `git diff --name-only` to see what changed
- If Go files changed, run `go build ./cmd/ctx/...`
- If build fails, stop and report: do not commit broken code

### 2. Stage and commit

- Review unstaged changes with `git status`
- Stage relevant files (prefer specific files over `git add -A`)
- Craft a concise commit message:
  - If the user provided a message, use it
  - If not, draft one based on the changes
- Commit with trailers as required by project conventions

### 3. Context prompt

After a successful commit, ask the user:

> **Any context to capture?**
>
> - **Decision**: Did you make a design choice or trade-off?
> - **Learning**: Did you hit a gotcha or discover something?
> - **Neither**: No context to capture.

If they provide a decision or learning, record it:

```bash
ctx add decision "..."
ctx add learning --context "..." --lesson "..." --application "..."
```

### 4. Doc drift check

If committed files include source code that could affect docs,
offer to check for doc drift.

### 5. Reflect (optional)

If the commit represents a significant milestone, suggest:

> This looks like a good checkpoint. Want me to run a quick
> reflection to capture the bigger picture?

## Commit Message Style

- Focus on **why**, not what (the diff shows what)
- Concise (1-2 sentences)
- Follow the repository's existing commit style
- Include required trailers (Spec:, Co-authored-by:, etc.)

## Quality Checklist

Before committing:
- [ ] Build passes (if Go files changed)
- [ ] Commit message explains the why
- [ ] No secrets in staged changes
- [ ] Specific files staged (not blind `git add -A`)

After committing:
- [ ] Context prompt was presented
- [ ] Any decisions/learnings were recorded
- [ ] Doc drift check offered if source code changed
