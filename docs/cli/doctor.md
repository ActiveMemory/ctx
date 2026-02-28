---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Doctor
icon: lucide/stethoscope
---

### `ctx doctor`

Structural health check across context, hooks, and configuration. Runs
mechanical checks that don't require semantic analysis. Think of it as
`ctx status` + `ctx drift` + configuration audit in one pass.

```bash
ctx doctor [flags]
```

**Flags**:

| Flag     | Short | Type | Default | Description                    |
|----------|-------|------|---------|--------------------------------|
| `--json` | `-j`  | bool | `false` | Machine-readable JSON output   |

---

#### What It Checks

| Check                    | Category  | What it verifies                                          |
|--------------------------|-----------|-----------------------------------------------------------|
| Context initialized      | Structure | `.context/` directory exists                              |
| Required files present   | Structure | All required context files exist (`TASKS.md`, etc.)       |
| Drift detected           | Quality   | Stale paths, missing files, constitution violations       |
| Event logging status     | Hooks     | Whether `event_log: true` is set in `.ctxrc`              |
| Webhook configured       | Hooks     | `.notify.enc` file exists                                 |
| Pending reminders        | State     | Count of entries in `reminders.json`                      |
| Task completion ratio    | State     | Pending vs completed tasks in `TASKS.md`                  |
| Context token size       | Size      | Estimated token count across all context files            |
| Recent event activity    | Events    | Last event timestamp (only when event logging is enabled) |

---

#### Output Format (Human)

```
ctx doctor
==========

Structure
  ✓ Context initialized (.context/)
  ✓ Required files present (4/4)

Quality
  ⚠ Drift: 2 warnings (stale path in ARCHITECTURE.md, high entry count in LEARNINGS.md)

Hooks
  ✓ hooks.json valid (14 hooks registered)
  ○ Event logging disabled (enable with event_log: true in .ctxrc)

State
  ✓ No pending reminders
  ⚠ Task completion ratio high (18/22 = 82%) — consider archiving

Size
  ✓ Context size: ~4200 tokens (budget: 8000)

Summary: 2 warnings, 0 errors
```

Status indicators:

| Icon | Status  | Meaning                         |
|------|---------|---------------------------------|
| ✓    | ok      | Check passed                    |
| ⚠    | warning | Non-critical issue worth fixing |
| ✗    | error   | Problem that needs attention    |
| ○    | info    | Informational note              |

---

#### Output Format (JSON)

```json
{
  "results": [
    {
      "name": "context_initialized",
      "category": "Structure",
      "status": "ok",
      "message": "Context initialized (.context/)"
    },
    {
      "name": "required_files",
      "category": "Structure",
      "status": "ok",
      "message": "Required files present (4/4)"
    },
    {
      "name": "drift",
      "category": "Quality",
      "status": "warning",
      "message": "Drift: 2 warnings"
    },
    {
      "name": "event_logging",
      "category": "Hooks",
      "status": "info",
      "message": "Event logging disabled (enable with event_log: true in .ctxrc)"
    },
    {
      "name": "webhook",
      "category": "Hooks",
      "status": "ok",
      "message": "Webhook configured"
    },
    {
      "name": "reminders",
      "category": "State",
      "status": "ok",
      "message": "No pending reminders"
    },
    {
      "name": "task_completion",
      "category": "State",
      "status": "warning",
      "message": "Tasks: 18/22 completed (82%) — consider archiving with ctx tasks archive"
    },
    {
      "name": "context_size",
      "category": "Size",
      "status": "ok",
      "message": "Context size: ~4200 tokens (budget: 8000)"
    }
  ],
  "warnings": 2,
  "errors": 0
}
```

---

#### Examples

```bash
# Quick structural health check
ctx doctor

# Machine-readable output for scripting
ctx doctor --json

# Count warnings
ctx doctor --json | jq '.warnings'

# Check for errors only
ctx doctor --json | jq '[.results[] | select(.status == "error")]'
```

---

#### When to Use What

| Tool             | When                                                     |
|------------------|----------------------------------------------------------|
| `ctx status`     | Quick glance at files, tokens, and drift                 |
| `ctx doctor`     | Thorough structural checkup (hooks, config, events too)  |
| `/ctx-doctor`    | Agent-driven diagnosis with event log pattern analysis   |

`ctx status` tells you *what's there*. `ctx doctor` tells you *what's wrong*.
`/ctx-doctor` tells you *why it's wrong* and *what to do about it*.

---

#### What It Does Not Do

- **No event pattern analysis** — that's the `/ctx-doctor` skill's job
- **No auto-fixing** — reports findings, doesn't modify anything
- **No external service checks** — doesn't verify webhook endpoint availability

---

**See also**: [Troubleshooting](../recipes/troubleshooting.md) |
[`ctx system events`](system.md#ctx-system-events) |
[`/ctx-doctor` skill](../reference/skills.md#ctx-doctor) |
[Detecting and Fixing Drift](../recipes/context-health.md)
