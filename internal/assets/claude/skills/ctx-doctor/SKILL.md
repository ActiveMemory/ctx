---
name: ctx-doctor
description: "Troubleshoot ctx behavior. Runs structural health checks, analyzes event log patterns, and presents findings with suggested actions."
allowed-tools: Bash(ctx:*), Read
---

Diagnose ctx problems by combining structural health checks with
event log analysis.

## When to Use

- User says "doctor", "diagnose", "troubleshoot", "health check"
- User asks "why didn't my hook fire?"
- User says "hooks seem broken" or "context seems stale"
- User says "too many nudges" or "something seems off"
- User asks "what happened last session?"

## When NOT to Use

- User wants a quick status check (use `/ctx-status`)
- User wants to fix drift (use `/ctx-drift`)
- User wants to change hook messages (use `ctx system message`)
- User wants to pause hooks (use `/ctx-pause`)

## Diagnostic Playbook

Follow this triage sequence:

### Phase 1: Structural Baseline

Run `ctx doctor --json` to get the full structural health report.

```bash
ctx doctor --json
```

Parse the JSON output. Note any warnings or errors.

### Phase 2: Event Log Analysis (if available)

If the doctor report shows event logging is enabled, query recent events:

```bash
ctx system events --json --last 100
```

If the user is asking about a specific hook:

```bash
ctx system events --hook <hook-name> --json --last 20
```

If event logging is not enabled, note: "Enable `event_log: true` in
`.ctxrc` for hook-level diagnostics."

### Phase 3: Targeted Investigation

Based on findings, check additional sources:

- **Hook config**: read `.claude/settings.local.json` to verify hook registration
- **Custom messages**: run `ctx system message list` to check for silenced hooks
- **RC config**: read `.ctxrc` to check configuration
- **Reminders**: run `ctx remind list` for pending reminders

### Phase 4: Present Findings

Structure your report as:

```
## Doctor Report

### Structural health
- Summarize ctx doctor results

### Event analysis (if available)
- Patterns, gaps, or anomalies in event data
- Specific hook behavior observations

### Suggested actions
- [ ] Actionable items based on findings
```

### Phase 5: Suggest, Don't Fix

Present actionable next steps but do NOT auto-fix anything.
The user decides what to act on.

## Available Data Sources

| Source | Command | What it reveals |
|--------|---------|-----------------|
| Structural health | `ctx doctor --json` | All mechanical checks |
| Event log | `ctx system events --json --last 100` | Recent hook activity |
| Event log (filtered) | `ctx system events --hook <name> --json` | Specific hook |
| Reminders | `ctx remind list` | Pending reminders |
| Hook messages | `ctx system message list` | Custom vs default |
| RC config | Read `.ctxrc` | Configuration |

## Graceful Degradation

If event logging is not enabled, the skill still works with reduced
capability. Run `ctx doctor` for structural checks and note that
event-level diagnostics require `event_log: true` in `.ctxrc`.
