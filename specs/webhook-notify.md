# Webhook Notifications (`ctx notify`)

## Context

ctx sessions are often long-running or autonomous. The human isn't always
watching. A webhook gives skills and loop scripts a way to reach outside the
terminal to get the user's attention.

## Design

### Command

```
ctx notify --event <name> [--session-id <id>] "message"
ctx notify setup                                          # interactive: set webhook URL
ctx notify test                                           # send a test notification
```

- `--event, -e` (required): Event name (e.g., `loop`, `implement`, `nudge`, `relay`)
- `--session-id, -s` (optional): Session ID for multi-agent disambiguation
- Positional arg (required): Notification message

**Behavior:**
- No webhook configured → silent noop (exit 0)
- Webhook set but event not in `events` list → silent noop (exit 0)
- Webhook set and event matches → fire-and-forget HTTP POST
- HTTP errors silently ignored (no retry, no response parsing)

### Webhook URL: Encrypted Storage

The webhook URL is sensitive (contains auth tokens). It MUST NOT be stored
in plaintext in `.ctxrc` or environment variables.

**Solution:** Reuse the existing `internal/crypto/` AES-256-GCM encryption
and the project's encryption key (`.context/.context.key`).

```
.context/.context.key   ← existing, gitignored, mode 0600
.context/.notify.enc       ← NEW, committed (encrypted), safe to share
```

**Setup flow:**
1. User runs `ctx notify setup`
2. Command prompts for webhook URL (stdin)
3. Encrypts URL using existing encryption key (auto-generates if needed)
4. Writes to `.context/.notify.enc`
5. Prints confirmation

**Runtime flow:**
1. `ctx notify` reads `.context/.notify.enc`
2. Decrypts with `.context/.context.key`
3. Uses decrypted URL for HTTP POST
4. If key or encrypted file missing → silent noop

### Encryption Key Rotation

The encryption key has no built-in age tracking. We add a lightweight check:

- **Detection:** `os.Stat(".context/.context.key").ModTime()` gives creation
  time (key is never modified after generation).
- **Threshold:** 90 days (configurable in `.ctxrc` as `key_rotation_days`).
- **Nudge:** Existing `check-version` or new system hook emits a VERBATIM relay:
  ```
  IMPORTANT: Relay this security reminder to the user VERBATIM.

  ┌─ Key Rotation ──────────────────────────────────────┐
  │ Your encryption key is N days old.                   │
  │ Consider rotating: ctx pad rotate-key                │
  └──────────────────────────────────────────────────────┘
  ```
- **Not blocking:** Just a nudge. User can ignore or act.

### Configuration (`.ctxrc`)

```yaml
notify:
  events:          # optional filter; absent/empty = all events pass
    - loop
    - implement
    - nudge
    - relay
  # key_rotation_days: 90    # optional, default 90
```

Note: `webhook` is NOT in `.ctxrc` — it's encrypted in `.context/.notify.enc`.

- `events` absent or empty → all events are sent
- `events` populated → only matching events fire

### Event Types

| Event | Source | Purpose |
|-------|--------|---------|
| `loop` | Loop script (generated bash) | Convergence, max-iteration hit |
| `implement` | ctx-implement skill (future) | Step failure, plan complete |
| `nudge` | System hooks (VERBATIM relays) | User assurance the nudge was emitted |
| `relay` | System hooks (all output types) | Debugging: verify agent received hook output |
| (user-defined) | User skills (e.g., `release`, `backup`) | User wires `ctx notify` in their skills |

**`nudge` event rationale:** The agent is instructed to relay VERBATIM messages,
but it's possible (though unlikely) for the agent to fabricate having seen one
or to skip it entirely. By notifying the webhook when the hook fires, the user
gets an independent channel to verify delivery.

**`relay` event rationale:** Broader than `nudge`. Fires on ALL hook outputs
(VERBATIM + agent directives). For debugging: "Did the agent get a post-commit
nudge? Did it get the QA reminder?" The webhook log becomes an independent
audit trail of what the agent was told.

### Payload (JSON POST)

```json
{
  "event": "loop",
  "message": "Loop completed after 5 iterations",
  "session_id": "abc123-...",
  "timestamp": "2026-02-22T14:30:00Z",
  "project": "ctx"
}
```

| Field | Type | Required | Source |
|-------|------|----------|--------|
| `event` | string | always | `--event` flag |
| `message` | string | always | positional arg |
| `session_id` | string | omitted if empty | `--session-id` flag |
| `timestamp` | string | always | UTC RFC3339 |
| `project` | string | always | `filepath.Base(os.Getwd())` |

### Who sends notifications?

| Layer | Sends? | How |
|-------|--------|-----|
| Human-authored hooks | Yes | `ctx notify` in hook commands |
| Skills (authored) | Yes, where appropriate | Skill code calls `ctx notify` |
| System hooks (built-in) | Yes, for `nudge`/`relay` events | Inline via shared `internal/notify/` |
| Agent (freeform) | No | Not in playbook, not encouraged |

### Hook Integration for `nudge` and `relay` Events

System hooks in `internal/cli/system/` that emit output also fire notifications
inline using a shared `internal/notify/` package. This avoids spawning a
subprocess.

**VERBATIM relay hooks** (fire both `nudge` and `relay`):
- `check-context-size`
- `check-persistence`
- `check-ceremonies`
- `check-journal`
- `check-resources`
- `check-knowledge`
- `check-version`

**Agent directive hooks** (fire `relay` only):
- `post-commit`
- `qa-reminder`

**Hard gate hooks** (fire `relay` only):
- `block-non-path-ctx`

**Silent hooks** (no notification):
- `cleanup-tmp`

## Architecture

```
internal/notify/               ← NEW: shared notification logic
  notify.go                    ← Send(), EventAllowed(), LoadWebhook()
  notify_test.go               ← Unit tests

internal/cli/notify/           ← NEW: CLI command
  notify.go                    ← Cmd() for "ctx notify"
  setup.go                     ← Cmd() for "ctx notify setup"
  run.go                       ← runNotify(), runSetup()
  notify_test.go               ← CLI integration tests

internal/cli/system/           ← MODIFIED: add notify calls to hooks
  check_context_size.go        ← Add notify.Send() after VERBATIM output
  (... other hook files ...)
```

The split keeps notification logic reusable: CLI command delegates to
`internal/notify/`, hooks import `internal/notify/` directly.

## What this does NOT include

- No retry logic or delivery guarantees (fire-and-forget)
- No response parsing (HTTP response body discarded)
- No payload templating (fixed JSON format)
- No bidirectional communication (no webhook responses trigger agent actions)
- No agent playbook integration (agents don't call `ctx notify`)
- No `--json` output flag (command produces no stdout)

## Files

### New
- `internal/notify/notify.go` — Shared notification logic
- `internal/notify/notify_test.go` — Tests
- `internal/cli/notify/notify.go` — Command definition
- `internal/cli/notify/setup.go` — Setup subcommand
- `internal/cli/notify/run.go` — Execution logic
- `internal/cli/notify/notify_test.go` — CLI tests

### Modified
- `internal/rc/types.go` — Add `NotifyConfig` struct, `Notify` field on `CtxRC`
- `internal/rc/rc.go` — Add `NotifyEvents()` accessor
- `internal/rc/rc_test.go` — Tests for notify config loading
- `internal/bootstrap/bootstrap.go` — Register `notify.Cmd`
- `internal/config/tpl_loop.go` — Add `TplLoopNotify`, update loop templates
- `internal/cli/loop/script.go` — Pass notify into template format calls
- `internal/cli/system/check_context_size.go` — Add relay/nudge notification
- `internal/cli/system/check_persistence.go` — Add relay/nudge notification
- `internal/cli/system/check_ceremonies.go` — Add relay/nudge notification
- `internal/cli/system/check_journal.go` — Add relay/nudge notification
- `internal/cli/system/check_resources.go` — Add relay/nudge notification
- `internal/cli/system/check_knowledge.go` — Add relay/nudge notification
- `internal/cli/system/check_version.go` — Add relay/nudge + key age nudge
- `internal/cli/system/post_commit.go` — Add relay notification
- `internal/cli/system/qa_reminder.go` — Add relay notification
- `internal/cli/system/block_non_path_ctx.go` — Add relay notification
