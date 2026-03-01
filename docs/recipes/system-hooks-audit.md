---
title: "Auditing System Hooks"
icon: lucide/shield-check
---

![ctx](../images/ctx-banner.png)

## The Problem

`ctx` runs 15 system hooks behind the scenes: nudging your agent to persist
context, warning about resource pressure, gating commits on QA. But these
hooks are **invisible by design**. You never see them fire. You never know
if they stopped working.

**How do you verify your hooks are actually running, audit what they do,
and get alerted when they go silent?**

## TL;DR

```bash
ctx system check-resources # run a hook manually
ls -la .context/logs/      # check hook execution logs
ctx notify setup           # get notified when hooks fire
```

Or ask your agent: *"Are our hooks running?"*

## Commands and Skills Used

| Tool                     | Type          | Purpose                                  |
|--------------------------|---------------|------------------------------------------|
| `ctx system <hook>`      | CLI command   | Run a system hook manually               |
| `ctx system resources`   | CLI command   | Show system resource status              |
| `ctx notify setup`       | CLI command   | Configure webhook for audit trail        |
| `ctx notify test`        | CLI command   | Verify webhook delivery                  |
| `.ctxrc` `notify.events` | Configuration | Subscribe to `relay` for full hook audit |
| `.context/logs/`         | Log files     | Local hook execution ledger              |

---

## What Are System Hooks?

System hooks are **plumbing commands** that `ctx` registers with your AI tool
(*Claude Code, Cursor, etc.*) via the plugin's `hooks.json`. They fire
automatically at specific events during your AI session:

| Event              | When                              | Hooks                               |
|--------------------|-----------------------------------|-------------------------------------|
| `UserPromptSubmit` | Before the agent sees your prompt | 9 check hooks + heartbeat           |
| `PreToolUse`       | Before the agent uses a tool      | `block-non-path-ctx`, `qa-reminder` |
| `PostToolUse`      | After a tool call succeeds        | `post-commit`                       |
| `SessionEnd`       | Session terminates                | `cleanup-tmp`                       |

You never run these manually. Your AI tool runs them for you: That's the
point.

---

## The Complete Hook Catalog

### Prompt-Time Checks (UserPromptSubmit)

These fire before every prompt, but most are throttled to avoid noise.

#### `check-context-size`: Context Capacity Warning

**What**: Adaptive prompt counter. Silent for the first 15 prompts, then
nudges with increasing frequency (*every 5th, then every 3rd*).

**Why**: Long sessions lose coherence. The nudge reminds both you and the
agent to persist context before the window fills up.

**Output**: VERBATIM relay box with prompt count.

```
┌─ Context Checkpoint (prompt #20) ────────────────
│ This session is getting deep. Consider wrapping up
│ soon. If there are unsaved learnings, decisions, or
│ conventions, now is a good time to persist them.
└──────────────────────────────────────────────────
```

---

#### `check-persistence`: Context Staleness Nudge

**What**: Tracks when `.context/*.md` files were last modified. If too many
prompts pass without a write, nudges the agent to persist.

**Why**: Sessions produce insights that evaporate if not recorded. This
catches the "*we talked about it but never wrote it down*" failure mode.

**Output**: VERBATIM relay after 20+ prompts without a context file change.

```
┌─ Persistence Checkpoint (prompt #20) ───────────
│ No context files updated in 20+ prompts.
│ Have you discovered learnings, made decisions,
│ established conventions, or completed tasks
│ worth persisting?
│
│ Run /ctx-wrap-up to capture session context.
└──────────────────────────────────────────────────
```

---

#### `check-ceremonies`: Session Ritual Adoption

**What**: Scans your last 3 journal entries for `/ctx-remember` and
`/ctx-wrap-up` usage. Nudges once per day if missing.

**Why**: Session ceremonies are the highest-leverage habit in `ctx`. This
hook bootstraps the habit until it becomes automatic.

**Output**: Tailored nudge depending on which ceremony is missing.

---

#### `check-journal`: Unexported Session Reminder

**What**: Detects unexported Claude Code sessions and unenriched journal
entries. Fires once per day.

**Why**: Exported sessions become searchable history. Unenriched entries
lack metadata for filtering. Both decay in value over time.

**Output**: VERBATIM relay with counts and exact commands.

```
┌─ Journal Reminder ─────────────────────────────
│ You have 3 new session(s) not yet exported.
│ 5 existing entries need enrichment.
│
│ Export and enrich:
│   ctx recall export --all
│   /ctx-journal-enrich-all
└────────────────────────────────────────────────
```

---

#### `check-resources`: System Resource Pressure

**What**: Monitors memory, swap, disk, and CPU load. Only fires at
**DANGER** severity (memory >= 90%, swap >= 75%, disk >= 95%,
load >= 1.5x CPU count).

**Why**: Resource exhaustion mid-session can corrupt work. This provides
early warning to persist and exit.

**Output**: VERBATIM relay listing critical resources.

---

#### `check-knowledge`: Knowledge File Growth

**What**: Counts entries in `LEARNINGS.md`, `DECISIONS.md`, and lines in
`CONVENTIONS.md`. Fires once per day when thresholds are exceeded.

**Why**: Large knowledge files dilute agent context. 35 learnings compete
for attention; 15 focused ones get applied. Thresholds are configurable
in `.ctxrc`.

**Default thresholds**:

```yaml
# .ctxrc
entry_count_learnings: 30
entry_count_decisions: 20
convention_line_count: 200
```

---

#### `check-version`: Binary/Plugin Version Drift

**What**: Compares the `ctx` binary version against the plugin version.
Fires once per day. Also checks encryption key age for rotation nudge.

**Why**: Version drift means hooks reference features the binary doesn't
have. The key rotation nudge prevents indefinite key reuse.

---

#### `check-reminders`: Pending Reminder Relay

**What**: Reads `.context/reminders.json` and surfaces any due reminders
via VERBATIM relay. **No throttle**: fires every session until dismissed.

**Why**: Reminders are sticky notes to future-you. Unlike nudges (*which
throttle to once per day*), reminders repeat deliberately until the user
dismisses them.

**Output**: VERBATIM relay box listing due reminders.

```
┌─ Reminders ──────────────────────────────────────
│  [1] refactor the swagger definitions
│
│ Dismiss: ctx remind dismiss <id>
│ Dismiss all: ctx remind dismiss --all
└──────────────────────────────────────────────────
```

---

#### `check-map-staleness`: Architecture Map Drift

**What**: Checks whether `map-tracking.json` is older than 30 days and
there are commits touching `internal/` since the last map refresh. Daily
throttle prevents repeated nudges.

**Why**: Architecture documentation drifts silently as code evolves. This
hook detects structural changes that the map hasn't caught up with and
suggests running `/ctx-map` to refresh.

**Output**: VERBATIM relay when stale and modules changed, silent otherwise.

```
┌─ Architecture Map Stale ────────────────────────────
│ ARCHITECTURE.md hasn't been refreshed since 2026-01-15
│ and there are commits touching 12 modules.
│ /ctx-map keeps architecture docs drift-free.
│
│ Want me to run /ctx-map to refresh?
└─────────────────────────────────────────────────────
```

---

#### `heartbeat`: Session Heartbeat Webhook

**What**: Fires on every prompt. Sends a webhook notification with prompt
count, session ID, and whether context files were modified since the last
heartbeat. Never produces stdout.

**Why**: Other hooks only send webhooks when they "speak" (nudge/relay).
When silent, you have no visibility into session activity. The heartbeat
provides a continuous session-alive signal for observability dashboards
or liveness monitoring.

**Output**: None (*webhook + event log only*).

**Payload**:

```json
{
  "event": "heartbeat",
  "message": "heartbeat: prompt #7 (context_modified=false)",
  "detail": {
    "hook": "heartbeat",
    "variant": "pulse",
    "variables": {"prompt_count": 7, "session_id": "abc...", "context_modified": false}
  }
}
```

---

### Tool-Time Hooks (PreToolUse / PostToolUse)

#### `block-non-path-ctx`: PATH Enforcement (Hard Gate)

**What**: Blocks any Bash command that invokes `./ctx`, `./dist/ctx`,
`go run ./cmd/ctx`, or an absolute path to `ctx`. Only PATH invocations
are allowed.

**Why**: Enforces `CONSTITUTION.md`'s invocation invariant. Running a
dev-built binary in production context causes version confusion and
silent behavior drift.

**Output**: Block response (*prevents the tool call*):

```json
{"decision": "block", "reason": "Use 'ctx' from PATH, not './ctx'..."}
```

---

#### `qa-reminder`: Pre-Commit QA Gate

**What**: Fires on every `Edit` tool use. Reminds the agent to lint and
test the **entire** project before committing.

**Why**: Agents tend to "*I'll test later*" and then commit untested code.
Repetition is **intentional**: the hook reinforces the habit on every edit,
not just before commits.

**Output**: Agent directive with hard QA gate instructions.

---

#### `post-commit`: Context Capture After Commit

**What**: Fires after any `git commit` (excludes `--amend`). Prompts the
agent to offer context capture (decision? learning?) and suggest running
lints/tests before pushing.

**Why**: Commits are natural reflection points. The nudge converts
mechanical git operations into context-capturing opportunities.

---

### Session Lifecycle (SessionEnd)

#### `cleanup-tmp`: Temp File Cleanup

**What**: Removes files older than 15 days from the `ctx` temp directory
(`$XDG_RUNTIME_DIR/ctx/` or `/tmp/ctx-<uid>/`).

**Why**: State files from throttling (*daily markers, prompt counters*)
accumulate. Silent cleanup prevents temp directory bloat.

**Output**: None (*silent side-effect*).

---

## Auditing Hooks via the Local Event Log

If you don't need an external audit trail, enable the **local event log** for
a self-contained record of hook activity:

```yaml
# .ctxrc
event_log: true
```

Once enabled, every hook that fires writes an entry to
`.context/state/events.jsonl`. Query it with `ctx system events`:

```bash
ctx system events                    # last 50 events
ctx system events --hook qa-reminder # filter by hook
ctx system events --session <id>     # filter by session
ctx system events --json | jq '.'    # raw JSONL for processing
```

The event log is local, queryable, and doesn't require any external service.
For a full diagnostic workflow combining event logs with structural health
checks, see [Troubleshooting](troubleshooting.md).

---

## Auditing Hooks via Webhooks

The most powerful audit setup pipes **all** hook output to a webhook,
giving you a real-time external record of what your agent is being told.

### Step 1: Set Up the Webhook

```bash
ctx notify setup
# Enter your webhook URL (Slack, Discord, ntfy.sh, IFTTT, etc.)
```

See [Webhook Notifications](webhook-notifications.md) for service-specific
setup.

### Step 2: Subscribe to `relay` Events

```yaml
# .ctxrc
notify:
  events:
    - relay   # all hook output: VERBATIM relays, directives, blocks
    - nudge   # just the user-facing VERBATIM relays
```

The `relay` event fires for **every** hook that produces output. This
includes:

| Hook                  | Event sent                  |
|-----------------------|-----------------------------|
| `check-context-size`  | `relay` + `nudge`           |
| `check-persistence`   | `relay` + `nudge`           |
| `check-ceremonies`    | `relay` + `nudge`           |
| `check-journal`       | `relay` + `nudge`           |
| `check-resources`     | `relay` + `nudge`           |
| `check-knowledge`     | `relay` + `nudge`           |
| `check-version`       | `relay` + `nudge`           |
| `check-reminders`     | `relay` + `nudge`           |
| `check-map-staleness` | `relay` + `nudge`           |
| `heartbeat`           | `heartbeat` only            |
| `block-non-path-ctx`  | `relay` only                |
| `post-commit`         | `relay` only                |
| `qa-reminder`         | `relay` only                |
| `cleanup-tmp`         | *(silent: no notification)* |

### Step 3: Cross-Reference

With `relay` enabled, your webhook receives a JSON payload every time a
hook fires:

```json
{
  "event": "relay",
  "message": "check-persistence: No context updated in 20+ prompts",
  "session_id": "b854bd9c",
  "timestamp": "2026-02-22T14:30:00Z",
  "project": "my-project"
}
```

This creates an **external audit trail** independent of the agent. You
can now cross-verify: did the agent actually relay the checkpoint the
hook told it to relay?

---

## Verifying Hooks Actually Fire

Hooks are invisible. An invisible thing that breaks is indistinguishable
from an invisible thing that never existed. Three verification methods,
from simplest to most robust:

### Method 1: Ask the Agent

The simplest check. After a few prompts into a session:

```text
"Did you receive any hook output this session? Print the last
context checkpoint or persistence nudge you saw."
```

The agent should be able to recall recent hook output from its context
window. If it says "*I haven't received any hook output*", either:

* The hooks aren't firing (*check installation*);
* The session is too short (*hooks throttle early*);
* The hooks fired but the agent absorbed them silently.

**Limitation**: You are trusting the agent to report accurately. Agents
sometimes confabulate or miss context. Use this as a quick smoke test,
not definitive proof.

### Method 2: Check the Webhook Trail

If you have `relay` events enabled, check your webhook receiver. Every
hook that fires sends a timestamped notification. No notification =
no fire.

This is the **ground truth**. The webhook is called directly by the `ctx`
binary, not by the agent. The agent cannot fake, suppress, or modify
webhook deliveries.

Compare what the webhook received against what the agent claims to have
relayed. Discrepancies mean the agent is absorbing nudges instead of
surfacing them.

### Method 3: Read the Local Logs

Hooks that support logging write to `.context/logs/`:

```bash
# Check context-size hook activity
cat .context/logs/check-context-size.log

# Sample output:
# [2026-02-22 09:15:00] [session:b854bd9c] prompt#1 silent
# [2026-02-22 09:17:33] [session:b854bd9c] prompt#16 CHECKPOINT
# [2026-02-22 09:20:01] [session:b854bd9c] prompt#20 CHECKPOINT
```

```bash
# Check persistence nudge activity
cat .context/logs/check-persistence.log

# Sample output:
# [2026-02-22 09:15:00] [session:b854bd9c] init count=1 mtime=1770646611
# [2026-02-22 09:20:01] [session:b854bd9c] prompt#20 NUDGE since_nudge=20
```

Logs are append-only and written by the `ctx` binary, not the agent.

---

## Detecting Silent Hook Failures

The hardest failure mode: hooks that **stop firing** without error. The
plugin config changes, a binary update drops a hook, or a PATH issue
silently breaks execution. Nothing errors: The hook just never runs.

### The Staleness Signal

If `.context/logs/check-context-size.log` has no entries newer than
5 days but you've been running sessions daily, something is wrong. The
absence of evidence is evidence of absence: but only if you control for
inactivity.

### False Positive Protection

A naive "*hooks haven't fired in N days*" alert fires incorrectly when
you simply haven't used `ctx`. The correct check needs two inputs:

1. **Last hook fire time**: from `.context/logs/` or webhook history
2. **Last session activity**: from journal entries or `ctx recall list`

If sessions are happening but hooks aren't firing, that's a real
problem. If neither sessions nor hooks are happening, that's a vacation.

### What to Check

When you suspect hooks aren't firing:

```bash
# 1. Verify the plugin is installed
ls ~/.claude/plugins/

# 2. Check hook registration
cat ~/.claude/plugins/ctx/hooks.json | head -20

# 3. Run a hook manually to see if it errors
echo '{"session_id":"test"}' | ctx system check-context-size

# 4. Check for PATH issues
which ctx
ctx --version
```

---

## Tips

* **Start with `nudge`, graduate to `relay`**: The `nudge` event covers
  user-facing VERBATIM relays. Add `relay` when you want full visibility
  into agent directives and hard gates.
* [**Webhooks are your trust anchor**](webhook-notifications.md): 
  The agent can ignore a nudge, but it can't suppress the webhook. 
  If the webhook fired and the agent didn't relay, you have proof of a 
  compliance gap.
* **Hooks are throttled by design**: Most check hooks fire once per day
  or use adaptive frequency. Don't expect a notification every prompt:
  Silence usually means the throttle is working, not that the hook is
  broken.
* **Daily markers live in temp**: Throttle files are stored in
  `$XDG_RUNTIME_DIR/ctx/` and cleaned up by `cleanup-tmp` after 15
  days. If you need to force a hook to re-fire during testing, delete
  the corresponding marker file.
* **The QA reminder is intentionally noisy**: Unlike other hooks,
  `qa-reminder` fires on every `Edit` call with no throttle. This is
  deliberate: The commit quality degrades when the reminder fades from
  salience.
* **Log files are safe to commit**: `.context/logs/` contains only
  timestamps, session IDs, and status keywords. No secrets, no code.

## Next Up

**[Detecting and Fixing Drift →](context-health.md)**: Keep context
files accurate as your codebase evolves.

## See Also

* [Troubleshooting](troubleshooting.md): full diagnostic workflow using
  `ctx doctor`, event logs, and `/ctx-doctor`
* [Customizing Hook Messages](customizing-hook-messages.md): override
  what hooks say without changing what they do
* [Webhook Notifications](webhook-notifications.md): setting up and
  configuring the webhook system
* [Hook Output Patterns](hook-output-patterns.md): understanding
  VERBATIM relays, agent directives, and hard gates
* [Detecting and Fixing Drift](context-health.md): structural checks
  that complement runtime hook auditing
* [CLI Reference](../cli/system.md): full `ctx system`
  command reference
