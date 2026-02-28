---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Troubleshooting"
icon: lucide/stethoscope
---

![ctx](../images/ctx-banner.png)

## The Problem

Something isn't working: a hook isn't firing, nudges are too noisy,
context seems stale, or the agent isn't following instructions. The
information to diagnose it exists -- across status, drift, event logs,
hook config, and session history -- but assembling it manually is tedious.

**How do you figure out what's wrong and fix it?**

## TL;DR

```bash
ctx doctor                   # structural health check
ctx system events --last 20  # recent hook activity
# or ask: "something seems off, can you diagnose?"
```

## Commands and Skills Used

| Tool                       | Type        | Purpose                              |
|----------------------------|-------------|--------------------------------------|
| `ctx doctor`               | CLI command | Structural health report             |
| `ctx doctor --json`        | CLI command | Machine-readable health report       |
| `ctx system events`        | CLI command | Query local event log                |
| `/ctx-doctor`              | Skill       | Agent-driven diagnosis with analysis |

---

## The Workflow

### Quick Check: `ctx doctor`

Run `ctx doctor` for an instant structural health report. It checks context
initialization, required files, drift, hook configuration, event logging,
webhooks, reminders, task completion ratio, and context token size -- all in
one pass:

```bash
ctx doctor
```

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

Warnings are non-critical but worth fixing. Errors need attention.
Informational notes (○) flag optional features that aren't enabled.

For scripting:

```bash
ctx doctor --json | jq '.warnings'
```

### Deep Dive: `/ctx-doctor`

When you need the agent to reason about what's wrong, use the skill.
Ask naturally or invoke directly:

```text
Why didn't my hook fire?
Something seems off, can you diagnose?
/ctx-doctor
```

The agent follows a triage sequence:

1. **Baseline**: runs `ctx doctor --json` for structural health
2. **Events**: runs `ctx system events --json --last 100` (if event logging enabled)
3. **Correlate**: connects findings across both sources
4. **Present**: structured findings with evidence
5. **Suggest**: actionable next steps (but doesn't auto-fix)

The skill degrades gracefully: without event logging enabled, it still runs
structural checks and notes what you'd gain by enabling it.

### Raw Event Inspection

For power users: `ctx system events` with filters gives direct access to the
event log.

```bash
# Last 50 events (default)
ctx system events

# Events from a specific session
ctx system events --session eb1dc9cd-0163-4853-89d0-785fbfaae3a6

# Only QA reminder events
ctx system events --hook qa-reminder

# Raw JSONL for jq processing
ctx system events --json | jq '.message'

# Include rotated (older) events
ctx system events --all --last 100
```

Filters use AND logic: `--hook qa-reminder --session abc123` returns only
QA reminder events from that specific session.

---

## Common Problems

### "My hook isn't firing"

**Symptoms**: No nudges appearing, webhook silent, event log shows no entries
for the expected hook.

**Diagnosis**:

```bash
# 1. Check if ctx is installed and on PATH
which ctx && ctx --version

# 2. Check if the hook is registered
grep "check-persistence" ~/.claude/plugins/ctx/hooks.json

# 3. Run the hook manually to see if it errors
echo '{"session_id":"test"}' | ctx system check-persistence

# 4. Check event log for the hook (if enabled)
ctx system events --hook check-persistence
```

**Common causes**:

- **Plugin not installed**: run `ctx init --claude` to reinstall
- **PATH issue**: the hook invokes `ctx` from PATH; ensure it resolves
- **Throttle active**: most hooks fire once per day -- check
  `$XDG_RUNTIME_DIR/ctx/` for daily marker files
- **Hook silenced**: a custom message override may be an empty file --
  check `ctx system message list` for overrides

### "Too many nudges"

**Symptoms**: The agent is overwhelmed with hook output. Context checkpoints,
persistence reminders, and QA gates fire constantly.

**Diagnosis**:

```bash
# Check how often hooks fired recently
ctx system events --last 50

# Count fires per hook
ctx system events --json | jq -r '.detail.hook // "unknown"' | sort | uniq -c | sort -rn
```

**Common causes**:

- **QA reminder is noisy by design**: it fires on every `Edit` call with no
  throttle. This is intentional. If it's too much, silence it with an empty
  override: `ctx system message edit qa-reminder gate`, then empty the file
- **Long session**: context checkpoint fires with increasing frequency after
  prompt 15. This is the system telling you the session is getting long --
  consider wrapping up
- **Short throttle window**: if you deleted marker files in
  `$XDG_RUNTIME_DIR/ctx/`, daily-throttled hooks will re-fire

### "Context seems stale"

**Symptoms**: The agent references outdated information, paths that don't
exist, or decisions that were reversed.

**Diagnosis**:

```bash
# Structural drift check
ctx drift

# Full doctor check (includes drift + more)
ctx doctor

# Check when context files were last modified
ctx status --verbose
```

**Common causes**:

- **Drift accumulated**: stale path references in `ARCHITECTURE.md` or
  `CONVENTIONS.md`. Fix with `ctx drift --fix` or ask the agent to clean up
- **Task backlog**: too many completed tasks diluting active context. Archive
  with `ctx tasks archive` or `ctx compact --archive`
- **Large context files**: `LEARNINGS.md` with 40+ entries competes for
  attention. Consolidate with `/ctx-consolidate`
- **Missing session ceremonies**: if `/ctx-remember` and `/ctx-wrap-up` aren't
  being used, context doesn't get refreshed. See
  [Session Ceremonies](session-ceremonies.md)

### "The agent isn't following instructions"

**Symptoms**: The agent ignores conventions, forgets decisions, or acts
contrary to `CONSTITUTION.md` rules.

**Diagnosis**:

```bash
# Check context token size -- is it too large for the model?
ctx doctor --json | jq '.results[] | select(.name == "context_size")'

# Check if context is actually being loaded
ctx system events --hook context-load-gate

# Audit alignment between docs and instructions
# (run inside agent session)
/ctx-alignment-audit
```

**Common causes**:

- **Context too large**: if total tokens exceed the model's effective attention,
  instructions get diluted. Check `ctx doctor` for the size check. Compact with
  `ctx compact --archive`
- **Context not loading**: if `context-load-gate` hasn't fired, the agent
  may not have received context. Verify the hook is registered
- **Conflicting instructions**: `CONVENTIONS.md` says one thing,
  `AGENT_PLAYBOOK.md` says another. Run `/ctx-alignment-audit` to find gaps
- **Agent drift**: the agent's behavior diverges from instructions over long
  sessions. This is normal. Use `/ctx-reflect` to re-anchor, or start a new
  session

---

## Prerequisites

- **Event logging** (*optional but recommended*): `event_log: true` in `.ctxrc`
- **ctx initialized**: `ctx init`

Event logging is not required for `ctx doctor` or `/ctx-doctor` to work. Both
degrade gracefully: structural checks run regardless, and the skill notes when
event data is unavailable.

---

## Tips

* **Start with `ctx doctor`**: It's the fastest way to get a comprehensive
  health picture. Save event log inspection for when you need to understand
  *when* and *how often* something happened.
* **Enable event logging early**: The log is opt-in and low-cost (~250 bytes
  per event, 1MB rotation cap). Enable it before you need it -- diagnosing
  a problem without historical data is much harder.
* **Use the skill for correlation**: `ctx doctor` tells you *what* is wrong.
  `/ctx-doctor` tells you *why* by correlating structural findings with event
  patterns. The agent can spot connections that individual commands miss.
* **Event log is gitignored**: It's machine-local diagnostic data, not project
  context. Different machines produce different event streams.

## Next Up

**[Detecting and Fixing Drift &rarr;](context-health.md)**: Keep context
files accurate as your codebase evolves.

## See Also

* [Auditing System Hooks](system-hooks-audit.md): the complete hook catalog
  and webhook-based audit trails
* [Detecting and Fixing Drift](context-health.md): structural and semantic
  drift detection and repair
* [Webhook Notifications](webhook-notifications.md): push notifications for
  hook activity
* [`ctx doctor` CLI](../cli/doctor.md): full command reference
* [`ctx system events` CLI](../cli/system.md#ctx-system-events): event log
  query reference
* [`/ctx-doctor` skill](../reference/skills.md#ctx-doctor): agent-driven
  diagnosis
