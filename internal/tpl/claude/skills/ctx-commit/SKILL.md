---
name: ctx-commit
description: "Commit with context persistence. Use instead of raw git commit to capture decisions and learnings alongside code changes."
---

Commit code changes, then prompt for decisions and learnings
worth persisting. Bridges the gap between committing code and
recording the context behind it.

## When to Use

- When committing after meaningful work (feature, bugfix,
  refactor)
- When the commit involves a design choice or trade-off that
  future sessions should know about
- When the user says "commit" or "commit this" — prefer this
  over raw git commit to capture context

## When NOT to Use

- For trivial commits (typo, formatting) — just commit normally
- When the user explicitly says "just commit, no context"
- When nothing has changed (no staged or unstaged modifications)

## Usage Examples

```text
/ctx-commit
/ctx-commit "implement session enrichment"
/ctx-commit --skip-qa
```

## Process

### 1. Pre-commit checks

Unless the user says `--skip-qa` or "skip checks":

- Run `git diff --name-only` to see what changed
- If Go files changed, run `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx`
  to verify the build
- If build fails, stop and report — do not commit broken code

### 2. Stage and commit

- Review unstaged changes with `git status`
- Stage relevant files (prefer specific files over `git add -A`)
- Craft a concise commit message:
  - If the user provided a message, use it
  - If not, draft one based on the changes (1-2 sentences,
    "why" not "what")
- Commit with the standard Co-Authored-By trailer

### 3. Context prompt

After a successful commit, ask the user:

> **Any context to capture?**
>
> - **Decision**: Did you make a design choice or trade-off?
>   (I'll record it with `ctx add decision`)
> - **Learning**: Did you hit a gotcha or discover something?
>   (I'll record it with `ctx add learning`)
> - **Neither**: No context to capture — we're done.

Wait for the user's response. If they provide a decision or
learning, record it using the appropriate command:

```bash
ctx add decision "..."
```

```bash
ctx add learning --context "..." --lesson "..." --application "..."
```

### 4. Doc drift check (conditional)

If the committed files include source code that could affect
documentation (Go files in `internal/cli/`, `internal/config/`,
`internal/tpl/`, `cmd/`), remind the user:

> Source files changed — want me to run `/update-docs` to check
> for doc drift?

Skip this prompt if:
- Only non-code files changed (markdown, config, scripts)
- Only test files changed
- The user already ran `/update-docs` this session

### 5. Reflect (optional)

If the commit represents a significant milestone (completing a
feature, finishing a multi-session effort, resolving a complex
bug), suggest a reflection:

> This looks like a good checkpoint. Want me to run a quick
> `/ctx-reflect` to capture the bigger picture?

Only suggest this for substantial commits — not every commit
needs reflection. Signs a reflection is warranted:
- Multiple files changed across different packages
- The commit closes out a task from TASKS.md
- The work spanned discussion of trade-offs or alternatives

## Commit Message Style

Follow the repository's existing commit style. Draft messages
that:
- Focus on **why**, not what (the diff shows what)
- Are concise (1-2 sentences)
- Use lowercase, no period at the end
- End with the Co-Authored-By trailer

Example:
```
add session enrichment to ctx-save skill

The save command produces scaffolds with empty placeholders.
The skill now instructs the agent to fill them in immediately
after saving.

Co-Authored-By: Claude <noreply@anthropic.com>
```

## What NOT to Do

- **Don't commit without asking** — always confirm the commit
  message with the user (or use their provided message)
- **Don't skip the context prompt** — this is the whole point
  of the skill; without it, use raw git commit
- **Don't force reflection** — suggest it only when warranted,
  and accept "no" gracefully
- **Don't commit secrets** — check for `.env`, credentials,
  tokens in the diff

## Quality Checklist

Before committing, verify:
- [ ] Build passes (if Go files changed)
- [ ] Commit message is concise and explains the why
- [ ] No secrets or sensitive files in the staged changes
- [ ] Specific files staged (not blind `git add -A`)

After committing, verify:
- [ ] Context prompt was presented to the user
- [ ] Any decisions/learnings provided were recorded
- [ ] Doc drift check was offered (if source code changed)
- [ ] Reflection was suggested if the commit was substantial
