---
name: ctx-drift
description: "Detect and fix context drift. Use to find stale paths, broken references, and constitution violations in context files."
allowed-tools: Bash(ctx:*), Bash(diff:*), Bash(mktemp:*), Bash(rm:cleanup temp), Read
---

Detect context drift at two layers: **structural** (stale paths,
missing files, constitution violations) via `ctx drift`, and
**semantic** (outdated conventions, superseded decisions,
irrelevant learnings) via agent analysis. The semantic layer is
where the real value is — the CLI cannot do it.

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

Drift detection has two layers: **structural** (programmatic) and
**semantic** (agent-driven). Always do both.

### Layer 1: Structural Checks

Run the CLI tool for fast, programmatic checks:

```bash
ctx drift
```

This catches dead paths, missing files, staleness indicators,
and constitution violations. These are necessary but insufficient
— they only detect structural problems.

### Layer 2: Semantic Analysis

After the structural check, read the context files yourself and
compare them to what you know about the codebase. This is where
you add real value — the CLI tool cannot do this.

Check for:

- **Outdated conventions**: Does CONVENTIONS.md describe patterns
  the code no longer follows? Read a few source files in the
  relevant area to verify.
- **Superseded decisions**: Does DECISIONS.md contain entries that
  were implicitly overridden by later work? Look for decisions
  whose rationale no longer applies.
- **Stale architecture descriptions**: Does ARCHITECTURE.md
  describe module purposes that have changed? A path can still
  exist while its description is wrong.
- **Irrelevant learnings**: Does LEARNINGS.md contain entries
  about bugs that were since fixed or patterns that no longer
  apply?
- **Contradictions**: Do any context files contradict each other
  or contradict the actual code?

### Reporting

After both layers, do **not** dump raw output. Instead:

1. **Summarize findings** by severity (structural warnings,
   semantic issues) in plain language
2. **Explain each finding**: what file, what line, why it
   matters
3. **Distinguish structural from semantic**: structural issues
   can be auto-fixed; semantic issues need the user's judgment
4. **Offer to auto-fix** structural issues:
   "I can run `ctx drift --fix` to clean up the dead path
   references. Want me to?"
5. **Propose specific edits** for semantic issues:
   "CONVENTIONS.md still says 'use fmt.Printf for output' but
   we switched to cmd.Printf three weeks ago. Want me to
   update it?"
6. **Suggest follow-up commands** when appropriate:
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

## Permission Drift

After checking skills, verify that `.claude/settings.local.json`
has the expected ctx permissions. This file is gitignored, so it
drifts independently from the codebase.

### Procedure

1. Read `.claude/settings.local.json` and extract the allow list.

2. Check for **missing ctx defaults**. Every entry in
   `DefaultClaudePermissions` (defined in
   `internal/config/file.go`) should be present. The current
   expected set is:

   - `Bash(ctx:*)` — covers all ctx subcommands
   - `Skill(ctx-*)` — one entry per ctx-shipped skill

   To get the authoritative list:

   ```bash
   ctx init --force 2>/dev/null  # in a temp dir
   ```

   Then compare permissions from the generated
   `settings.local.json` against the project's copy.

3. Check for **stale skill permissions**. If a `Skill(ctx-*)`
   entry references a skill that no longer exists in
   `.claude/skills/`, flag it.

4. Check for **missing skill permissions**. If a `ctx-*` skill
   exists in `.claude/skills/` but has no corresponding
   `Skill(ctx-*)` in the allow list, flag it.

### Interpreting Permission Drift

| Finding                           | Action                                          |
|-----------------------------------|-------------------------------------------------|
| Missing `Bash(ctx:*)`            | Suggest adding — required for ctx to work       |
| Missing `Skill(ctx-*)` entry     | Suggest adding — skill will prompt every time   |
| Stale `Skill(ctx-*)` entry       | Suggest removing — dead reference               |
| Granular `Bash(ctx <sub>:*)`     | Suggest consolidating to `Bash(ctx:*)`          |
| One-off / session debris entries  | Note as hygiene issue (see `hack/runbooks/sanitize-permissions.md`) |

### Important

Do **not** edit `settings.local.json` directly. Report findings
and let the user make changes. This file controls agent
permissions — self-modification is a security concern. Refer
users to `hack/runbooks/sanitize-permissions.md` for the manual cleanup
procedure.

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
