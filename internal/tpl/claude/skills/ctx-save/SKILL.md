---
name: ctx-save
description: "Save session snapshot. Use when significant progress made, complex task completed, or before starting a risky operation."
allowed-tools: Bash(ctx:*)
---

Save the current context state to `.context/sessions/`.

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
| Summary            | Placeholder for user to fill in     |
| In Progress tasks  | Extracted from `.context/TASKS.md`  |
| Next Up tasks      | Extracted from `.context/TASKS.md`  |
| Recent decisions   | Last 3 from `.context/DECISIONS.md` |
| Recent learnings   | Last 5 from `.context/LEARNINGS.md` |
| Next session tasks | Placeholder for follow-up work      |

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

Report the saved file path to the user.

## Process

1. **Choose a descriptive topic**: if the user did not provide
   one, suggest one based on the session's work (e.g.,
   "cooldown-mechanism", "skill-sweep")
2. **Run the save command** with appropriate topic and type
3. **Report the file path** so the user knows where it is
4. **Suggest filling in the Summary section** if the user wants
   richer context for future sessions

## Quality Checklist

Before saving, verify:
- [ ] Topic is descriptive (not "manual-save" unless truly
      generic)
- [ ] `--type` matches the session's nature (feature, bugfix,
      refactor, or generic session)
- [ ] There is meaningful work to snapshot (do not save empty
      sessions)
- [ ] Reported the file path after saving
