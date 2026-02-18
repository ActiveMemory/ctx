---
name: ctx-recall
description: "Browse session history. Use when referencing past discussions or finding context from previous work."
allowed-tools: Bash(ctx:*)
---

Browse, inspect, and export AI session history.

## When to Use

- When the user asks "what did we do last time?" or references
  a past discussion
- When looking for context from previous work sessions
- When exporting sessions to the journal for enrichment
- When searching for a specific session by topic, date, or ID

## When NOT to Use

- When the user just wants current context (use `/ctx-status`
  or `/ctx-agent` instead)
- When session data is already loaded in context (no need to
  re-fetch)
- For modifying session content (recall is read-only; edit
  journal files directly)

## Usage Examples

```text
/ctx-recall
/ctx-recall list --limit 5
/ctx-recall show <slug-or-id>
/ctx-recall export --all
```

## Subcommands

### `ctx recall list`

List recent sessions, newest first.

| Flag             | Short | Default | Purpose                              |
|------------------|-------|---------|--------------------------------------|
| `--limit`        | `-n`  | 20      | Maximum sessions to show             |
| `--project`      | `-p`  | ""      | Filter by project name               |
| `--tool`         | `-t`  | ""      | Filter by tool (e.g., "claude-code") |
| `--all-projects` |       | false   | Include all projects                 |

Output per session: slug, short ID, project, branch, time,
duration, turn count, token count, first message preview.

### `ctx recall show`

Show details of a specific session.

| Flag              | Default | Purpose                          |
|-------------------|---------|----------------------------------|
| `--latest`        | false   | Show the most recent session     |
| `--full`          | false   | Full conversation (not preview)  |
| `--all-projects`  | false   | Search across all projects       |

Accepts a session identifier: full UUID, partial UUID prefix,
or slug name. Use `--latest` if no ID is given.

Default output shows metadata and the first 5 user messages.
Use `--full` for the complete conversation.

### `ctx recall export`

Export sessions to `.context/journal/` as markdown.

| Flag              | Default | Purpose                          |
|-------------------|---------|----------------------------------|
| `--all`           | false   | Export all sessions              |
| `--all-projects`  | false   | Include all projects             |
| `--force`         | false   | Overwrite existing files         |
| `--skip-existing` | false   | Skip files that already exist    |

Accepts a session ID, or use `--all` to export everything.
Default behavior preserves YAML frontmatter from previous
exports (enrichment data is not lost).

Large sessions (>200 messages) are automatically split into
parts with navigation links between them.

## Data Source

Sessions are read from `~/.claude/projects/` (Claude Code
JSONL files). The system auto-detects and parses session files;
only the current project's sessions are shown by default.

## Process

1. **Determine intent**: does the user want to list, inspect,
   or export?
2. **Run the appropriate subcommand** with relevant flags
3. **Summarize results**: for `list`, highlight notable sessions;
   for `show`, summarize key points; for `export`, report what
   was written and suggest next steps (normalize, enrich)

## Typical Workflows

**"What did we work on recently?"**
```bash
ctx recall list --limit 5
```

**"Show me that session about authentication"**
```bash
ctx recall list --project auth
# then with the slug or ID from the list:
ctx recall show <slug>
```

**"Export everything to the journal"**
```bash
ctx recall export --all --skip-existing
```
Then suggest: normalize (`/ctx-journal-normalize`) and enrich
(`/ctx-journal-enrich`) as next steps.

## Quality Checklist

Before reporting results, verify:
- [ ] Used the right subcommand for the user's intent
- [ ] Applied filters if the user mentioned a project, date,
      or topic
- [ ] For export, reminded the user about the normalize/enrich
      pipeline as next steps
- [ ] Did not re-export sessions the user already has (use
      `--skip-existing` by default unless asked otherwise)
