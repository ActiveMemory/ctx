---
name: ctx-recall
description: "Browse session history. Use when referencing past discussions or finding context from previous work."
tools: [bash]
---

Browse, inspect, and export AI session history.

## When to Use

- When the user asks "what did we do last time?"
- When looking for context from previous work sessions
- When exporting sessions to the journal for enrichment
- When searching for a specific session by topic or date

## When NOT to Use

- When the user just wants current context (use `ctx-status` or
  `ctx-agent` instead)
- For modifying session content (recall is read-only)

## Subcommands

### `ctx recall list`

```bash
ctx recall list --limit 5
```

### `ctx recall show`

```bash
ctx recall show <slug-or-id>
ctx recall show --latest
```

### `ctx recall export`

```bash
ctx recall export --all        # Export new sessions only
ctx recall export --all --regenerate  # Re-export all
```

## Typical Workflows

**"What did we work on recently?"**
```bash
ctx recall list --limit 5
```

**"Export everything to the journal"**
```bash
ctx recall export --all
```

Then suggest `ctx-journal-enrich-all` for enrichment.

## Quality Checklist

- [ ] Used the right subcommand for user intent
- [ ] Applied filters if user mentioned project, date, or topic
- [ ] For export, mentioned the normalize/enrich pipeline
