---
name: ctx-doctor
description: "Troubleshoot ctx behavior. Runs structural health checks, analyzes event log patterns, and presents findings with suggested actions."
tools: [bash, read]
---

Diagnose ctx problems by combining structural health checks with
event log analysis.

## When to Use

- User says "doctor", "diagnose", "troubleshoot", "health check"
- User asks "why didn't my hook fire?"
- User says "hooks seem broken" or "context seems stale"

## When NOT to Use

- User wants a quick status check (use `ctx-status`)
- User wants to fix drift (use `ctx-drift`)
- User wants to pause hooks (use `ctx-pause`)

## Diagnostic Playbook

### Phase 1: Structural Baseline

```bash
ctx doctor --json
```

Parse the JSON output. Note any warnings or errors.

### Phase 2: Event Log Analysis (if available)

```bash
ctx system events --json --last 100
```

For specific hooks:
```bash
ctx system events --hook <hook-name> --json --last 20
```

If event logging is not enabled, note: "Enable `event_log: true`
in `.ctxrc` for hook-level diagnostics."

### Phase 3: Targeted Investigation

Based on findings, check:
- Hook config for hook registration
- Custom messages: `ctx system message list`
- RC config: read `.ctxrc`
- Reminders: `ctx remind list`

### Phase 4: Present Findings

```
## Doctor Report

### Structural health
- Summarize ctx doctor results

### Event analysis (if available)
- Patterns, gaps, or anomalies

### Suggested actions
- [ ] Actionable items based on findings
```

### Phase 5: Suggest, Don't Fix

Present actionable next steps but do NOT auto-fix anything.

## Quality Checklist

- [ ] Ran `ctx doctor` for structural checks
- [ ] Checked event log if available
- [ ] Presented findings in structured format
- [ ] Suggested actions without auto-applying
