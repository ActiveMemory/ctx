---
name: ctx-status
description: "Show context summary. Use at session start or when unclear about current project state."
tools: [bash]
---

Show the current context status: files, token budget, tasks,
and recent activity.

## When to Use

- At session start to orient before doing work
- When confused about what's being worked on
- To check token usage and context health
- When the user asks "what's the state of the project?"

## When NOT to Use

- When you already loaded context via `ctx-agent` in this session
- Repeatedly within the same session without changes

## Flags

| Flag        | Default | Purpose                        |
|-------------|---------|--------------------------------|
| `--json`    | false   | Output as JSON (for scripting) |
| `--verbose` | false   | Include file content previews  |

## Execution

```bash
ctx status
```

After running, summarize the key points:
- How many active tasks remain
- Whether any context files are empty
- Token budget usage
- What was recently modified

## Interpreting Results

| Observation             | Suggestion                                      |
|-------------------------|-------------------------------------------------|
| Many empty files        | Populate core files (TASKS, CONVENTIONS)         |
| High token count (>30k) | Consider `ctx compact` or archiving tasks       |
| No recent activity      | Context may be stale; check if files need update |
| TASKS.md has 0 active   | All work done, or tasks need to be added        |

## Quality Checklist

- [ ] Summarized the output (do not just dump raw output)
- [ ] Flagged empty core files
- [ ] Noted token budget if high or low
