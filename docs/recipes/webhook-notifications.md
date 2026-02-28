---
title: "Webhook Notifications"
icon: lucide/bell
---

![ctx](../images/ctx-banner.png)

## The Problem

Your agent runs autonomously (*loops, implements, releases*) while you are away
from the terminal. You have no way to know when it finishes, hits a limit, or
when a hook fires a nudge.

**How do you get notified about agent activity without watching the terminal?**

## TL;DR

```bash
ctx notify setup  # configure webhook URL (encrypted)
ctx notify test   # verify delivery
# Hooks auto-notify on: session-end, loop-iteration, resource-danger
```

## Commands and Skills Used

| Tool                              | Type          | Purpose                                 |
|-----------------------------------|---------------|-----------------------------------------|
| `ctx notify setup`                | CLI command   | Configure and encrypt webhook URL       |
| `ctx notify test`                 | CLI command   | Send a test notification                |
| `ctx notify --event <name> "msg"` | CLI command   | Send a notification from scripts/skills |
| `.ctxrc` `notify.events`          | Configuration | Filter which events reach your webhook  |

## The Workflow

### Step 1: Get a Webhook URL

Any service that accepts HTTP POST with JSON works. Common options:

| Service      | How to get a URL                                                       |
|--------------|------------------------------------------------------------------------|
| **IFTTT**    | Create an applet with the "*Webhooks*" trigger                         |
| **Slack**    | Create an [Incoming Webhook](https://api.slack.com/messaging/webhooks) |
| **Discord**  | Channel Settings > Integrations > Webhooks                             |
| **ntfy.sh**  | Use `https://ntfy.sh/your-topic` (no signup)                           |
| **Pushover** | Use API endpoint with your user key                                    |

**The URL contains auth tokens**. `ctx` encrypts it; it never appears in plaintext
in your repo.

### Step 2: Configure the Webhook

```bash
ctx notify setup
# Enter webhook URL: https://maker.ifttt.com/trigger/ctx/json/with/key/YOUR_KEY
# Webhook configured: https://maker.ifttt.com/***
# Encrypted at: .context/.notify.enc
```

This encrypts the URL with AES-256-GCM using the same key as the scratchpad
(`.context/.context.key`). The encrypted file (`.context/.notify.enc`) is
safe to commit. The key is `.gitignore`d.

### Step 3: Test It

```bash
ctx notify test
# Webhook responded: HTTP 200 OK
```

If you see `No webhook configured`, run `ctx notify setup` first.

### Step 4: Configure Events

Notifications are opt-in: no events are sent unless you configure an event
list in `.ctxrc`:

```yaml
# .ctxrc
notify:
  events:
    - loop       # loop completion or max-iteration hit
    - nudge      # VERBATIM relay hooks (context checkpoint, persistence, etc.)
    - relay      # all hook output (verbose, for debugging)
```

Only listed events fire. Omitting an event silently drops it.

### Step 5: Use in Your Own Skills

Add `ctx notify` calls to any skill or script:

```bash
# In a release skill
ctx notify --event release "v1.2.0 released successfully" 2>/dev/null || true

# In a backup script
ctx notify --event backup "Nightly backup completed" 2>/dev/null || true
```

The `2>/dev/null || true` suffix ensures the notification never breaks your
script: If there's no webhook or the HTTP call fails, it's a silent noop.

## Event Types

`ctx` fires these events automatically:

| Event      | Source            | When                                                                                                                  |
|------------|-------------------|-----------------------------------------------------------------------------------------------------------------------|
| `loop`     | Loop script       | Loop completes or hits max iterations                                                                                 |
| `nudge`    | System hooks      | VERBATIM relay nudge is emitted (context checkpoint, persistence, ceremonies, journal, resources, knowledge, version) |
| `relay`    | System hooks      | Any hook output (VERBATIM relays, agent directives, block responses)                                                  |
| `test`     | `ctx notify test` | Manual test notification                                                                                              |
| *(custom)* | Your skills       | You wire `ctx notify --event <name>` in your own scripts                                                              |

**`nudge` vs `relay`**: The `nudge` event fires only for VERBATIM relay hooks
(*the ones the agent is instructed to show verbatim*). The `relay` event fires
for *all* hook output — VERBATIM relays, agent directives, and hard gates.
Subscribe to `relay` for debugging (*"did the agent get the post-commit nudge?"*),
`nudge` for user-facing assurance (*"was the checkpoint emitted?"*).

!!! tip "Webhooks as a Hook Audit Trail"
    Subscribe to `relay` events and you get an external record of every
    hook that fires, independent of the agent. 

    This lets you verify hooks are running and catch cases where the agent 
    absorbs a nudge instead of surfacing it. 

    See [Auditing System Hooks](system-hooks-audit.md) for the full workflow.

## Payload Format

Every notification sends a JSON POST:

```json
{
  "event": "nudge",
  "message": "check-context-size: Context window at 82%",
  "detail": {
    "hook": "check-context-size",
    "variant": "window",
    "variables": {"Percentage": 82, "TokenCount": "164k"}
  },
  "session_id": "abc123-...",
  "timestamp": "2026-02-22T14:30:00Z",
  "project": "ctx"
}
```

The `detail` field is a structured template reference containing the hook
name, variant, and any template variables. This lets receivers filter by
hook or variant without parsing rendered text. The field is omitted when
no template reference applies (e.g. custom `ctx notify` calls).

## Security Model

| Component      | Location                   | Committed?      | Permissions |
|----------------|----------------------------|-----------------|-------------|
| Encryption key | `.context/.context.key`    | No (gitignored) | `0600`      |
| Encrypted URL  | `.context/.notify.enc`     | Yes (safe)      | `0600`      |
| Webhook URL    | Never on disk in plaintext | N/A             | N/A         |

The key is shared with the scratchpad. If you rotate the encryption key,
re-run `ctx notify setup` to re-encrypt the webhook URL with the new key.

## Key Rotation

`ctx` checks the age of `.context/.context.key` once per day. If it's older
than 90 days (*configurable via `notify.key_rotation_days`*), a VERBATIM nudge
is emitted suggesting rotation.

```yaml
# .ctxrc
notify:
  key_rotation_days: 30   # nudge sooner (default: 90)
```

## Worktrees

The webhook URL is encrypted with the same encryption key
(`.context/.context.key`), which is `.gitignore`d. In a git worktree,
the key is absent: Notifications silently do nothing.

This means **agents running in worktrees cannot send webhook alerts**.
For autonomous runs where worktree agents are opaque, monitor them from
the terminal rather than relying on webhooks. Enrich journals and review
results on the main branch after merging.

## Event Log: The Local Complement

Don't need a webhook but want diagnostic visibility? Enable `event_log: true`
in `.ctxrc`. The event log writes the same payload as webhooks to a local
JSONL file (`.context/state/events.jsonl`) that you can query without any
external service:

```bash
ctx system events --last 20          # recent hook activity
ctx system events --hook qa-reminder # filter by hook
```

Webhooks and event logging are independent: you can use either, both, or
neither. Webhooks give you push notifications and an external audit trail.
The event log gives you local queryability and `ctx doctor` integration.

See [Troubleshooting](troubleshooting.md) for how they work together.

---

## Tips

* **Fire-and-forget**: Notifications never block. HTTP errors are silently
  ignored. No retry, no response parsing.
* **No webhook = no cost**: When no webhook is configured, `ctx notify` exits
  immediately. System hooks that call `notify.Send()` add zero overhead.
* **Multiple projects**: Each project has its own `.notify.enc`. You can point
  different projects at different webhooks.
* **Event filter is per-project**: Configure `notify.events` in each project's
  `.ctxrc` independently.

## Next Up

**[Auditing System Hooks →](system-hooks-audit.md)**: Verify your hooks
are running, audit what they do, and get alerted when they go silent.

## See Also

* [CLI Reference: ctx notify](../cli/tools.md#ctx-notify):
  full command reference
* [Configuration](../home/configuration.md): `.ctxrc` settings including
  `notify` options
* [Running an Unattended AI Agent](autonomous-loops.md): how loops work
  and how notifications fit in
* [Hook Output Patterns](hook-output-patterns.md): understanding VERBATIM
  relays, agent directives, and hard gates
* [Auditing System Hooks](system-hooks-audit.md): using webhooks as an
  external audit trail for hook execution
