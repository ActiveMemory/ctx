---
name: ctx-recall
description: "Browse session history. Use when referencing past discussions or finding context from previous work."
tools: [bash]
---

Browse, inspect, and import AI session history.

## When to Use

- When the user asks "what did we do last time?"
- When looking for context from previous work sessions
- When importing sessions to the journal for enrichment
- When searching for a specific session by topic or date

## When NOT to Use

- When the user just wants current context (use `ctx-status` or
  `ctx-agent` instead)
- For modifying session content (browsing is read-only)

## Subcommands

### `ctx journal source`

```bash
ctx journal source --limit 5
```

### `ctx journal source --show` / `--latest`

```bash
ctx journal source --show <slug-or-id>
ctx journal source --latest
```

### `ctx journal import`

```bash
ctx journal import --all              # Import new sessions only
ctx journal import --all --regenerate # Re-import all
```

## Typical Workflows

**"What did we work on recently?"**
```bash
ctx journal source --limit 5
```

**"Import everything to the journal"**
```bash
ctx journal import --all
```

Then suggest `ctx-journal-enrich-all` for enrichment.

## Quality Checklist

- [ ] Used the right subcommand for user intent
- [ ] Applied filters if user mentioned project, date, or topic
- [ ] For import, mentioned the normalize/enrich pipeline
