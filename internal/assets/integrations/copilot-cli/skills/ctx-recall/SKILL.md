---
name: ctx-recall
description: "Browse session history. Use when referencing past discussions or finding context from previous work."
---

Browse, inspect, and export AI session history.

## When to Use

- When the user asks "what did we do last time?"
- When looking for context from previous work sessions
- When exporting sessions to the journal
- When searching for a specific session by topic or date

## When NOT to Use

- When the user just wants current context (use ctx-status instead)
- For modifying session content (recall is read-only)

## Execution

List recent sessions:

```bash
ctx recall list --limit 5
```

Show details of a specific session:

```bash
ctx recall show --latest
ctx recall show <session-id>
```

Export sessions to journal markdown:

```bash
ctx recall export --all
```

After listing sessions, summarize relevant findings rather than
dumping raw output.
