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

### 1. Check CONSTITUTION for commit rules

Read `.context/CONSTITUTION.md` (if it exists) for commit-specific
rules. Common project rules to look for and enforce:

- **Spec-per-commit**: Add a `Spec:` trailer, verify a spec file exists in
  `specs/` before proceeding. If no spec exists, stop and offer to run
  `/ctx-spec` to scaffold one.
- **Sign-off**: `Signed-off-by:`, include it.
- **Other trailers**: Honor any project-specific trailer requirements.

Read CONSTITUTION fully and apply all relevant rules before
proceeding to pre-commit checks.

### 2. Pre-commit checks

Unless the user says `--skip-qa` or "skip checks":

- Run `git diff --name-only` to see what changed
- Run the project's build and lint commands to verify nothing is broken.
  Check for a Makefile, package.json, or equivalent. If you cannot
  identify the build/lint commands, ask the user before proceeding.
- If the build or lint fails, stop and report: do not commit broken code

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
> - **Learning**: Did you hit a gotcha or discover something?
> - **Neither**: No context to capture: we're done.

Wait for the user's response. If they provide a decision or
learning, record it using the appropriate command:

```bash
ctx add decision "Use PostgreSQL" \
  --context "Need a reliable database" \
  --rationale "ACID compliance and JSON support" \
  --consequence "Team needs training"
```

```bash
ctx add learning "Go embed requires files in same package" \
  --context "..." --lesson "..." --application "..."
```

### 5. Reflect

After every commit, run `/ctx-reflect` to capture the bigger
picture before moving on. This is mandatory — skipping reflection
is how context gets lost between sessions.

## Commit Message Format

Follow the repository's existing commit style. Draft messages
that:
- Focus on **why**, not what (the diff shows what)
- Use lowercase, no period at the end
- Scale detail to match scope: a one-file fix gets 1-2 sentences;
  a multi-package change gets a summary paragraph plus a bulleted
  list of what changed and why
- Include any trailers required by CONSTITUTION (e.g., `Spec:`,
  `Signed-off-by:`)

Example:
```
complete journal-recall merge wiring and cross-cutting cleanup

Wire journal commands through journal/core packages instead of
recall/core. Move importer, lock, unlock, sync cmd packages from
recall/cmd to journal/cmd.

Changes:
- journal/core/{plan,execute,query} are now canonical
- sourcefm/sourceformat renamed to source/frontmatter, source/format
- Magic numbers extracted to config/stats constants
- state.StateDir renamed to state.Dir across 26 callers
- splitLines moved to parse.ByteLines
- /ctx-commit skill generalized to be language-agnostic

Spec: specs/journal-merge-completion.md
Signed-off-by: Jane Doe <jane@example.com>
```

## Commit Discipline

- **Spec trailer is mandatory** — identify the spec that covers
  this work and include `Spec:` in the commit message. If
  CONSTITUTION also requires it, this is non-negotiable.
- **Confirm the message** with the user before committing (or use
  their provided message)
- **Always present the context prompt**: this is the whole point
  of the skill
- **Always reflect**: even a one-sentence reflection prevents
  context loss
- **Check for secrets** (`.env`, credentials, tokens) in the diff
  before staging

## Quality Checklist

Before committing, verify:
- [ ] Spec exists and is referenced in the commit message
- [ ] Build and lint pass
- [ ] Commit message is concise and explains the why
- [ ] `Spec:` and `Signed-off-by:` trailers are present
- [ ] No secrets or sensitive files in the staged changes
- [ ] Specific files staged (not blind `git add -A`)

After committing, verify:
- [ ] Context prompt was presented to the user
- [ ] Any decisions/learnings provided were recorded
- [ ] Reflection was completed

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
