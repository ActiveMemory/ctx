# Session Pause

## Problem

Many sessions involve tasks where context nudges add no value: quick
investigations, one-off questions, small fixes unrelated to the project's
active work. Every hook still fires — ceremony nudges, persistence
reminders, knowledge checks, checkpoint boxes — costing tokens and
attention on work that doesn't benefit from them.

The initial ~8k context injection at session start (CLAUDE.md, playbook,
constitution) is unavoidable — Claude Code loads it before any user
command runs. That's an acceptable cost: it gives the agent the vocabulary
to know `/ctx-pause` exists. The expensive part is the ongoing hooks that
fire repeatedly throughout the session.

## Approach

A session-scoped pause flag stored in `secureTempDir()`, keyed by session
ID. All hook commands check the flag early and no-op when paused. Two new
plumbing commands (`ctx system pause` / `ctx system resume`) set and clear
the flag. Two skills (`/ctx-pause` / `/ctx-resume`) provide discoverability.

### Design Principles

- **Session-scoped, not global.** The flag is keyed to the session ID in
  the temp directory. Other sessions (same project, different terminal)
  are unaffected.
- **Hooks still fire, they just no-op.** No hook registration changes.
  Every hook checks the flag and exits early if set.
- **Explicit commands still work.** `ctx status`, `ctx agent`, etc. are
  CLI commands, not hooks — they are unaffected by pause.
- **Graduated reminder.** Paused hooks emit a minimal indicator so the
  state is never invisible.

## Behavior

### Happy Path

1. User starts session → context loads normally (~8k tokens)
2. User runs `ctx pause` or `/ctx-pause`
3. `ctx system pause` writes pause marker to `secureTempDir()`
4. All subsequent hooks detect the marker and exit early
5. Turns 1–5 while paused: hooks emit `ctx:paused` (2 tokens)
6. Turn 6+: hooks emit `ctx:paused (N turns) — resume with /ctx-resume`
7. User runs `ctx resume` or `/ctx-resume` when ready
8. `ctx system resume` removes the marker
9. Hooks resume normal behavior; turn counter resets

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Pause without active session ID | Use `session-unknown` key (same as other hooks) |
| Resume when not paused | Silent no-op, exit 0 |
| Pause when already paused | Reset turn counter to 0, exit 0 |
| Session ends while paused | Marker is session-scoped in tmpfs; `cleanup-tmp` handles stale files |
| Multiple sessions, one paused | Each session has its own marker keyed by session ID; independent |
| `ctx status` while paused | Works normally — it's a CLI command, not a hook |
| `ctx system bootstrap` while paused | Works normally — not a hook |

### What Gets Paused

All `UserPromptSubmit` hooks:

- `check-context-size` — checkpoint and window warnings
- `check-ceremonies` — session ceremony nudges
- `check-persistence` — context persistence nudges
- `check-journal` — journal maintenance reminders
- `check-reminders` — pending reminders relay
- `check-version` — version update nudges
- `check-resources` — resource pressure warnings
- `check-knowledge` — knowledge file growth nudges
- `check-map-staleness` — architecture map staleness

All `PreToolUse` hooks:

- `context-load-gate` — context file read directives
- `qa-reminder` — QA reminder before edits
- `specs-nudge` — plan-to-specs nudge
- `block-non-path-ctx` — still fires (security, not a nudge)
- `block-dangerous-commands` — still fires (security, not a nudge)

`PostToolUse` and `SessionEnd` hooks:

- `post-commit` — paused (nudge)
- `cleanup-tmp` — still fires (housekeeping, not a nudge)

**Rule:** security and housekeeping hooks always fire. Only nudge/reminder
hooks respect the pause flag.

### Graduated Reminder

The pause marker file stores the turn count (number of hooks that
checked and no-opped). Each hook invocation increments the count.

| Paused turns | Hook output |
|--------------|-------------|
| 1–5 | `ctx:paused` |
| 6+ | `ctx:paused (N turns) — resume with /ctx-resume` |

Only the *first* hook per prompt turn should emit the reminder (to avoid
5+ identical lines when multiple hooks fire on the same prompt). Use a
per-prompt dedup marker: `ctx-pause-emitted-{sessionID}-{promptCount}`.

Alternatively, simpler: only `check-context-size` emits the reminder
(it always fires first in UserPromptSubmit). Other hooks silently no-op.

## Interface

### CLI

```
ctx pause
ctx resume
```

These are user-facing convenience commands that delegate to the plumbing:

```
ctx system pause --session-id <id>
ctx system resume --session-id <id>
```

The top-level `ctx pause` / `ctx resume` commands read the session ID from
stdin (same as hooks) or accept `--session-id` flag.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--session-id` | string | from stdin | Session ID to pause/resume |

### Skills

```
/ctx-pause
/ctx-resume
```

Trigger phrases:
- "pause ctx" / "pause context" / "stop the nudges" / "quiet mode"
- "resume ctx" / "resume context" / "turn nudges back on" / "unpause"

Skill behavior:
- `/ctx-pause`: runs `ctx pause`, confirms with "Context hooks paused
  for this session. Resume with /ctx-resume."
- `/ctx-resume`: runs `ctx resume`, confirms with "Context hooks
  resumed."

## Implementation

### State File

Path: `secureTempDir()/ctx-paused-{sessionID}`

Contents: integer turn count (same format as existing counter files).

```
0     ← just paused
5     ← 5 hooks have checked and no-opped
12    ← 12 hooks have checked
```

Presence of the file = paused. Absence = not paused.

### Shared Helper

Add to `state.go`:

```go
// pauseMarkerPath returns the path to the session pause marker file.
func pauseMarkerPath(sessionID string) string {
    return filepath.Join(secureTempDir(), "ctx-paused-"+sessionID)
}

// paused checks if the session is paused. If paused, increments the
// turn counter and returns the current count. Returns 0 if not paused.
func paused(sessionID string) int {
    path := pauseMarkerPath(sessionID)
    data, err := os.ReadFile(path)
    if err != nil {
        return 0
    }
    count, _ := strconv.Atoi(strings.TrimSpace(string(data)))
    count++
    writeCounter(path, count)
    return count
}

// pausedMessage returns the appropriate pause indicator for the given
// turn count, or empty string if not paused.
func pausedMessage(turns int) string {
    if turns == 0 {
        return ""
    }
    if turns <= 5 {
        return "ctx:paused"
    }
    return fmt.Sprintf("ctx:paused (%d turns) — resume with /ctx-resume", turns)
}
```

### Hook Integration

Each pausable hook adds an early-return block after reading input:

```go
sessionID := input.SessionID
if sessionID == "" {
    sessionID = sessionUnknown
}
if turns := paused(sessionID); turns > 0 {
    // Only check-context-size emits the reminder to avoid duplication
    if emitPauseReminder {
        cmd.Println(pausedMessage(turns))
    }
    return nil
}
```

To avoid every hook emitting the reminder, designate `check-context-size`
as the single emitter (it's the first UserPromptSubmit hook). All other
hooks call `paused()` to increment the counter but don't print.

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/cli/system/state.go` | Add `pauseMarkerPath()`, `paused()`, `pausedMessage()` |
| `internal/cli/system/pause.go` | New: `ctx system pause` plumbing command |
| `internal/cli/system/resume.go` | New: `ctx system resume` plumbing command |
| `internal/cli/system/pause_test.go` | Tests for pause/resume/counter |
| `internal/cli/system/system.go` | Register `pauseCmd()`, `resumeCmd()` |
| `internal/cli/system/check_context_size.go` | Add pause check + reminder emission |
| `internal/cli/system/check_ceremonies.go` | Add pause check (silent) |
| `internal/cli/system/check_persistence.go` | Add pause check (silent) |
| `internal/cli/system/check_journal.go` | Add pause check (silent) |
| `internal/cli/system/check_reminders.go` | Add pause check (silent) |
| `internal/cli/system/check_version.go` | Add pause check (silent) |
| `internal/cli/system/check_resources.go` | Add pause check (silent) |
| `internal/cli/system/check_knowledge.go` | Add pause check (silent) |
| `internal/cli/system/check_map_staleness.go` | Add pause check (silent) |
| `internal/cli/system/context_load_gate.go` | Add pause check (silent) |
| `internal/cli/system/qa_reminder.go` | Add pause check (silent) |
| `internal/cli/system/post_commit.go` | Add pause check (silent) |
| `internal/cli/system/specs_nudge.go` | Add pause check (silent) |
| `internal/cli/pause.go` | New: top-level `ctx pause` command |
| `internal/cli/resume.go` | New: top-level `ctx resume` command |
| `internal/bootstrap/bootstrap.go` | Register top-level pause/resume commands |
| `.claude/skills/ctx-pause/SKILL.md` | New skill |
| `.claude/skills/ctx-resume/SKILL.md` | New skill |
| `internal/assets/claude/skills/ctx-pause/SKILL.md` | New skill template |
| `internal/assets/claude/skills/ctx-resume/SKILL.md` | New skill template |

## Configuration

None. No `.ctxrc` keys, no environment variables. The feature is
self-contained and session-scoped.

## Testing

- **Unit**: `pause()` creates file and increments counter; `paused()`
  returns 0 when no file; `pausedMessage()` returns correct strings
  for turns 1, 5, 6, 100
- **Unit**: `ctx system pause` creates marker; `ctx system resume`
  removes it; double-pause resets counter; resume when not paused
  is no-op
- **Integration**: run `check-context-size` with pause marker present,
  verify it emits `ctx:paused` and increments counter instead of
  the normal checkpoint
- **Integration**: verify `block-non-path-ctx` and
  `block-dangerous-commands` still fire when paused

## Non-Goals

- **Preventing initial context load.** CLAUDE.md injection is controlled
  by Claude Code, not ctx. The ~8k startup cost is accepted.
- **Granular hook selection.** No "pause only ceremony nudges." It's
  all-or-nothing for nudge hooks. Complexity isn't justified.
- **Persistent pause.** The flag lives in tmpfs and is session-scoped.
  No `.ctxrc` setting to start paused.
- **"Light mode."** A middle tier that suppresses some hooks but keeps
  others was considered and rejected for v1. Start simple.
