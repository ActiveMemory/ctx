---
name: ctx-drift
description: "Detect and fix context drift. Use to find stale paths, broken references, and constitution violations in .context/ files."
allowed-tools: Bash(ctx:*), Bash(diff:*), Bash(mktemp:*), Bash(rm:cleanup temp), Read
---

Detect context drift: stale path references, missing files,
constitution violations, bloated task lists, and outdated
skills. Explain findings in plain language and offer to fix
what can be automated.

## When to Use

- At session start to verify context health before working
- After refactors, renames, or major structural changes
- When the user asks "is our context clean?", "anything
  stale?", or "check for drift"
- Proactively when you notice a path in ARCHITECTURE.md or
  CONVENTIONS.md that does not match the actual file tree
- Before a release or milestone to ensure context is accurate

## When NOT to Use

- When you just ran `/ctx-status` and everything looked fine
  (status already shows drift warnings)
- Repeatedly in the same session without changes in between
- When the user is mid-flow on a task; do not interrupt with
  unsolicited maintenance

## Usage Examples

```text
/ctx-drift
/ctx-drift (after the refactor)
```

## Execution

Run drift detection:

```bash
ctx drift
```

After running, do **not** dump raw output. Instead:

1. **Summarize findings** by severity (warnings, violations,
   staleness) in plain language
2. **Explain each finding**: what file, what line, why it
   matters
3. **Offer to auto-fix** if fixable issues were found:
   "I can run `ctx drift --fix` to clean up the dead path
   references. Want me to?"
4. **Suggest follow-up commands** when appropriate:
   - Many stale paths after a refactor → suggest `ctx sync`
   - Heavy task clutter → suggest `ctx compact --archive`
   - Old files untouched for weeks → suggest reviewing content

## Interpreting Results

| Finding                        | What It Means                          | Suggested Action                       |
|--------------------------------|----------------------------------------|----------------------------------------|
| Path does not exist            | Context references a deleted file/dir  | Remove reference or update path        |
| Directory is empty             | Referenced dir exists but has no files | Remove reference or populate directory |
| Many completed tasks           | TASKS.md is cluttered                  | Run `ctx compact --archive`            |
| File not modified in 30+ days  | Content may be outdated                | Review and update or confirm current   |
| Constitution violation         | A hard rule may be broken              | Fix immediately                        |
| Required file missing          | A core context file does not exist     | Create it with `ctx init` or manually  |

## Auto-Fix

When the user agrees to auto-fix:

```bash
ctx drift --fix
```

After fixing, run `ctx drift` again to confirm remaining
issues need manual attention. Report what was fixed and what
still needs the user's judgment.

## Skill Template Drift

After running `ctx drift`, check whether the project's
installed skills (`.claude/skills/`) match the canonical
templates shipped with `ctx`.

### Procedure

1. Create a temp directory and run `ctx init --force` inside
   it to get the latest templates:

   ```bash
   CTX_TPL_DIR=$(mktemp -d)
   cd "$CTX_TPL_DIR" && ctx init --force 2>/dev/null
   ```

2. Compare each skill in the project against the template:

   ```bash
   diff -ru "$CTX_TPL_DIR/.claude/skills/" .claude/skills/ 2>/dev/null
   ```

3. Clean up the temp directory:

   ```bash
   rm -rf "$CTX_TPL_DIR"
   ```

### Interpreting Skill Drift

| Finding                              | Action                                              |
|--------------------------------------|-----------------------------------------------------|
| Skill missing from project           | Offer to install: copy from template                |
| Skill differs from template          | Show the diff; offer to update to latest template   |
| Project has extra skills (no match)  | These are custom — leave them alone                 |
| No differences                       | Skills are up to date; report clean                 |

When reporting skill drift, distinguish between:

- **ctx-managed skills** (present in the template): these
  should generally match; differences mean the user's copy
  is outdated or was customized intentionally
- **Custom skills** (only in the project): these are user
  additions and should not be flagged as drift

If a skill was intentionally customized, note it and move on.
Offer to update only ctx-managed skills, and always show the
diff before overwriting.

## Proactive Use

Run drift detection without being asked when:

- You load context at session start and notice a path
  reference that does not match the file tree
- The user just completed a refactor that renamed or moved
  files
- TASKS.md has obviously heavy clutter (20+ completed items
  visible when you read it)

When running proactively, keep the report brief:

> I ran a quick drift check after the refactor. Two stale
> path references in ARCHITECTURE.md — want me to clean
> them up?

## Quality Checklist

After running drift detection, verify:
- [ ] Summarized findings in plain language (did not just
      paste raw CLI output)
- [ ] Explained why each finding matters
- [ ] Offered auto-fix for fixable issues before running it
- [ ] Suggested appropriate follow-up commands
- [ ] Did not run `--fix` without user confirmation
