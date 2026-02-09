---
name: ctx-save
description: "Save session snapshot. Use when significant progress made, complex task completed, or before starting a risky operation."
allowed-tools: Bash(ctx:*), Bash(git diff *)
---

Save the current context state to `.context/sessions/`, then
enrich it with a real summary and concrete next-session tasks.

## When to Use

- When significant progress has been made and it should be
  checkpointed
- After completing a complex task or feature
- Before starting a risky operation (refactor, migration,
  dependency upgrade)
- When context is getting full and the session may end soon
- When the user explicitly asks to save

## When NOT to Use

- After trivial changes (a typo fix does not need a snapshot)
- When nothing meaningful has happened yet in the session
- Immediately after another save with no new work in between

## Usage Examples

```text
/ctx-save
/ctx-save auth-refactor
/ctx-save "database migration" --type feature
```

## What Gets Saved

The snapshot aggregates current context into a single markdown
file:

| Section            | Source                              |
|--------------------|-------------------------------------|
| Date, time, type   | Current timestamp + `--type` flag   |
| Summary            | Filled in by agent (see below)      |
| In Progress tasks  | Extracted from `.context/TASKS.md`  |
| Next Up tasks      | Extracted from `.context/TASKS.md`  |
| Recent decisions   | Last 3 from `.context/DECISIONS.md` |
| Recent learnings   | Last 5 from `.context/LEARNINGS.md` |
| Next session tasks | Filled in by agent (see below)      |
| Files modified     | Added by agent from git diff        |

File is written to:
``.context/sessions/YYYY-MM-DD-HHMMSS-<topic>.md``

## Flags

| Flag     | Short | Default   | Purpose                                  |
|----------|-------|-----------|------------------------------------------|
| `--type` | `-t`  | `session` | Type: feature, bugfix, refactor, session |

The topic argument is optional; defaults to "manual-save".

## Related Subcommands

`ctx session` has other subcommands that complement `save`:

| Subcommand          | Purpose                                          |
|---------------------|--------------------------------------------------|
| `ctx session list`  | List saved sessions (newest first)               |
| `ctx session load`  | Display a saved session by index, date, or topic |
| `ctx session parse` | Convert JSONL transcripts to markdown            |

## Execution

```bash
ctx session save $ARGUMENTS
```

Note the saved file path from the command output.

## Process

1. **Choose a descriptive topic**: if the user did not provide
   one, suggest one based on the session's work (e.g.,
   "cooldown-mechanism", "skill-sweep")
2. **Run the save command** with appropriate topic and type
3. **Enrich the saved file** (see After Saving below)
4. **Report the file path** so the user knows where it is

## After Saving

Read the saved file and replace the placeholder sections:

1. **Summary**: Replace `[Describe what was accomplished...]`
   with 1-2 paragraphs covering: what was built/fixed, key
   decisions made, problems encountered and how they were
   resolved. Write from memory of this session — you have the
   full conversation context.

2. **Tasks for Next Session**: Replace `[List tasks to
   continue...]` with concrete items. Pull from: unfinished
   work this session, new tasks discovered, blocked items that
   need attention. Use `- [ ]` checkbox format.

3. **Files Modified**: Add a `## Files Modified` section after
   Summary listing files changed this session. Run:
   ```bash
   git diff --name-only HEAD~N  # where N covers this session's commits
   ```
   If no commits were made yet, use `git diff --name-only` for
   unstaged changes. Format as a bullet list with a brief
   description of each change.

Write the enriched file back using the Edit tool.

Do NOT add sections that require information you don't have
(like "Key Quotes" or "Participants"). Stick to what you can
synthesize from the session.

### Enrichment Examples

**Good summary**:
> Implemented the `/ctx-remember` skill with template and live
> copies. Build and tests pass. Also fixed the
> block-non-path-ctx.sh hook to guide agents to ask users for
> sudo instead of attempting it directly.

**Bad summary** (this is the failure mode we're fixing):
> `[Describe what was accomplished in this session]`

**Good next-session tasks**:
> - [ ] Add tests for the new session enrichment logic
> - [ ] Update CONVENTIONS.md with the new skill pattern
> - [ ] Investigate flaky test in `internal/cli/session/`

**Bad next-session tasks**:
> `[List tasks to continue in the next session]`

## Quality Checklist

Before saving, verify:
- [ ] Topic is descriptive (not "manual-save" unless truly
      generic)
- [ ] `--type` matches the session's nature (feature, bugfix,
      refactor, or generic session)
- [ ] There is meaningful work to snapshot (do not save empty
      sessions)

After enrichment, verify:
- [ ] Summary contains real content (no placeholders)
- [ ] Tasks for Next Session contains concrete items
- [ ] Files Modified section is present and accurate
- [ ] Reported the file path to the user
