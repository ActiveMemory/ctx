---
name: ctx-history
description: "Browse session history. Use when referencing past discussions or finding context from previous work."
allowed-tools: Bash(ctx:*)
---

Browse, inspect, and import AI session history.

## When to Use

- When the user asks "what did we do last time?" or references
  a past discussion
- When looking for context from previous work sessions
- When importing sessions to the journal for enrichment
- When searching for a specific session by topic, date, or ID

## When NOT to Use

- When the user just wants current context (use `/ctx-status`
  or `/ctx-agent` instead)
- When session data is already loaded in context (no need to
  re-fetch)
- For modifying session content (source browsing is read-only;
  edit journal files directly)

## Usage Examples

```text
/ctx-history
/ctx-history list --limit 5
/ctx-history show <slug-or-id>
/ctx-history import --all
```

## Subcommands

### `ctx journal source`

List recent sessions, newest first.

| Flag             | Short | Default | Purpose                              |
|------------------|-------|---------|--------------------------------------|
| `--limit`        | `-n`  | 20      | Maximum sessions to show             |
| `--project`      | `-p`  | ""      | Filter by project name               |
| `--tool`         | `-t`  | ""      | Filter by tool (e.g., "claude-code") |
| `--all-projects` |       | false   | Include all projects                 |
| `--show`         |       | ""      | Show details of a specific session   |
| `--latest`       |       | false   | Show the most recent session         |
| `--full`         |       | false   | Full conversation (not preview)      |

Output per session: slug, short ID, project, branch, time,
duration, turn count, token count, first message preview.

Use `--show <id>` to inspect a specific session. Accepts a
full UUID, partial UUID prefix, or slug name. Use `--latest`
if no ID is given.

Default output shows metadata and the first 5 user messages.
Use `--full` for the complete conversation.

### `ctx journal import`

Import sessions to the journal directory as markdown.

| Flag                 | Default | Purpose                                          |
|----------------------|---------|--------------------------------------------------|
| `--all`              | false   | Import all sessions (only new files by default)  |
| `--all-projects`     | false   | Include all projects                             |
| `--regenerate`       | false   | Re-import existing files (preserves frontmatter) |
| `--keep-frontmatter` | true    | Preserve enriched YAML frontmatter during regen  |
| `--yes`, `-y`        | false   | Skip confirmation prompt                         |
| `--dry-run`          | false   | Preview what would be imported                   |

Accepts a session ID (always writes), or `--all` to import
everything (safe by default: only new sessions, existing
files skipped). Use `--regenerate` with `--all` to re-import
existing files; YAML frontmatter is preserved by default.
Use `--keep-frontmatter=false` to discard enriched frontmatter.

Locked entries (via `ctx journal lock`) are always skipped.

Large sessions (>200 messages) are automatically split into
parts with navigation links between them.

### `ctx journal lock`

Protect journal entries from import regeneration.

```bash
ctx journal lock <pattern>     # Lock matching entries
ctx journal lock --all         # Lock all entries
```

### `ctx journal unlock`

Remove lock protection from journal entries.

```bash
ctx journal unlock <pattern>   # Unlock matching entries
ctx journal unlock --all       # Unlock all entries
```

### `ctx journal sync`

Sync lock state from journal frontmatter to `.state.json`.

```bash
ctx journal sync
```

Scans all journal markdowns and updates `.state.json` to match
each file's frontmatter. Files with `locked: true` in frontmatter
are marked locked in state; files without a `locked:` line have
their lock cleared. This is the inverse of `ctx journal lock`:
frontmatter drives state instead of state driving frontmatter.
Useful after batch enrichment where you add `locked: true` to
frontmatter manually.

## Data Source

Sessions are read from `~/.claude/projects/` (Claude Code
JSONL files). The system auto-detects and parses session files;
only the current project's sessions are shown by default.

## Process

1. **Determine intent**: does the user want to list, inspect,
   or import?
2. **Run the appropriate subcommand** with relevant flags
3. **Summarize results**: for `list`, highlight notable sessions;
   for `show`, summarize key points; for `import`, report what
   was written and suggest next steps (normalize, enrich)

## Typical Workflows

**"What did we work on recently?"**
```bash
ctx journal source --limit 5
```

**"Show me that session about authentication"**
```bash
ctx journal source --project auth
# then with the slug or ID from the list:
ctx journal source --show <slug>
```

**"Import everything to the journal"**
```bash
ctx journal import --all
```
This only imports new sessions: existing files are skipped.
If the user asks what to do next, mention that `/ctx-journal-enrich-all`
can enrich the imported journals.

**"Re-import sessions after a format improvement"**
```bash
ctx journal import --all --regenerate -y
```

## Quality Checklist

Before reporting results, verify:
- [ ] Used the right subcommand for the user's intent
- [ ] Applied filters if the user mentioned a project, date,
      or topic
- [ ] For import, reminded the user about the normalize/enrich
      pipeline as next steps
- [ ] Used `--all` for bulk import (safe: only new sessions)
- [ ] Suggested `--dry-run` when user seems uncertain
- [ ] Only used `--regenerate` when explicitly needed
