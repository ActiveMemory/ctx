---
name: ctx-commit
description: "Commit with context persistence. Use instead of raw git commit to capture decisions and learnings alongside code changes."
---

Commit code changes, then prompt for decisions and learnings
worth persisting. Bridges the gap between committing code and
recording the context behind it.

## When to Use

- For ALL commits. This is the only way to commit in this project.
  Raw `git commit` bypasses spec enforcement and violates CONSTITUTION.
- When the user says "commit", "commit this", "ship it", "let's commit":
  always use this skill, never raw git commit.

## When NOT to Use

- When nothing has changed (no staged or unstaged modifications)

## Usage Examples

```text
/ctx-commit
/ctx-commit "implement session enrichment"
/ctx-commit --skip-qa
```

## Process

### 1. Spec verification (CONSTITUTION requirement)

Every commit MUST reference a spec. Before anything else:

1. Identify which spec covers the current work. Check:
   - The in-progress task in TASKS.md — does it reference a `Spec:` line?
   - Recent specs in `specs/` that match the work being done
2. Verify the spec file exists: `ls specs/<name>.md`
3. If no spec exists:
   - **Stop.** Do not proceed with the commit.
   - Tell the agent: "No spec found for this work. Create a
     retroactive spec in `specs/` or ask the human to scope one.
     This is a CONSTITUTION requirement — no exceptions. Even a
     one-liner fix needs a spec for traceability."
   - If the human explicitly authorizes skipping the spec for this
     commit, note this in the commit message body.

The spec reference goes in the commit message as a `Spec:` trailer
(see Commit Message Format below).

### 2. Pre-commit checks

Unless the user says `--skip-qa` or "skip checks":

- Run `git diff --name-only` to see what changed
- If Go files changed, run `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx`
  to verify the build
- If build fails, stop and report: do not commit broken code

### 3. Stage and commit

- Review unstaged changes with `git status`
- Stage relevant files (prefer specific files over `git add -A`)
- Craft a concise commit message:
  - If the user provided a message, use it
  - If not, draft one based on the changes (1-2 sentences,
    "why" not "what")
- Include the `Spec:` and `Signed-off-by:` trailers (see format below)

### 4. Context prompt

After a successful commit, ask the user:

> **Any context to capture?**
>
> - **Decision**: Did you make a design choice or trade-off?
>   (I'll record it with `ctx add decision`)
> - **Learning**: Did you hit a gotcha or discover something?
>   (I'll record it with `ctx add learning`)
> - **Neither**: No context to capture: we're done.

Wait for the user's response. If they provide a decision or
learning, record it using the appropriate command:

```bash
ctx add decision "..."
```

```bash
ctx add learning --context "..." --lesson "..." --application "..."
```

### 5. Doc drift check (conditional)

If the committed files include source code that could affect
documentation (Go files in `internal/cli/`, `internal/config/`,
`internal/assets/`, `cmd/`), remind the user:

> Source files changed - want me to run `/_ctx-update-docs` to check
> for doc drift?

Skip this prompt if:
- Only non-code files changed (markdown, config, scripts)
- Only test files changed
- The user already ran `/update-docs` this session

### 6. Reflect (optional)

If the commit represents a significant milestone (completing a
feature, finishing a multi-session effort, resolving a complex
bug), suggest a reflection:

> This looks like a good checkpoint. Want me to run a quick
> `/ctx-reflect` to capture the bigger picture?

Only suggest this for substantial commits: not every commit
needs reflection. Signs a reflection is warranted:
- Multiple files changed across different packages
- The commit closes out a task from TASKS.md
- The work spanned discussion of trade-offs or alternatives

## Commit Message Format

Follow the repository's existing commit style. Draft messages
that:
- Focus on **why**, not what (the diff shows what)
- Are concise (1-2 sentences)
- Use lowercase, no period at the end
- Include `Spec:` trailer referencing the spec file (CONSTITUTION requirement)
- Include `Signed-off-by:` trailer

Example:
```
gate checkpoint nudges behind minimum context window percentage

Counter-based checkpoints fire regardless of context usage,
producing noise at 5-8% on 1M windows. Gate behind 20% minimum.

Spec: specs/hook-accountability.md
Signed-off-by: <users-name-configured-in-git> <users-email-configured-in-git>
```

## Commit Discipline

- **Spec trailer is mandatory** — this is the primary gate. If you
  cannot identify a spec, stop and resolve before committing.
- **Confirm the message** with the user before committing (or use
  their provided message)
- **Always present the context prompt**: this is the whole point
  of the skill
- **Suggest reflection only when warranted**: and accept "no"
  gracefully
- **Check for secrets** (`.env`, credentials, tokens) in the diff
  before staging

## Quality Checklist

Before committing, verify:
- [ ] Spec exists and is referenced in the commit message
- [ ] Build passes (if Go files changed)
- [ ] Commit message is concise and explains the why
- [ ] `Spec:` and `Signed-off-by:` trailers are present
- [ ] No secrets or sensitive files in the staged changes
- [ ] Specific files staged (not blind `git add -A`)

After committing, verify:
- [ ] Context prompt was presented to the user
- [ ] Any decisions/learnings provided were recorded
- [ ] Doc drift check was offered (if source code changed)
- [ ] Reflection was suggested if the commit was substantial

## Human Relay

After every successful commit, relay a structured summary to the
human verbatim:

```
┌─ Commit Summary ─────────────────────────
│ Spec: specs/<name>.md
│ Tasks closed: <list or "none">
│ Files changed: <count>
│ Message: <first line of commit message>
└──────────────────────────────────────────
```
