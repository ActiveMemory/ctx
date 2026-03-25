# Spec: Task Completion Nudge Hook

## Problem

After completing work, agents don't mark tasks as done in TASKS.md.
The user has to remember to ask, leading to:

- Unmarked completed tasks accumulating
- Secondary sweeps in the next session
- Task status that doesn't reflect reality
- `#in-progress` labels rarely applied (historically ~50% compliance)

The AGENT_PLAYBOOK already says to mark tasks done, but the instruction
is buried in a table row and deprioritized under context pressure.

## Decision

Add a PostToolUse hook on Edit that nudges the agent to check TASKS.md
for completed work. The nudge is debounced to avoid noise during active
editing.

## Design

### Hook Type

**PostToolUse** on `Edit` (and `Write`) tool invocations.

### Why Not Other Triggers

| Trigger | Problem |
|---------|---------|
| After `git commit` | User defers commits across sessions; too rare |
| After test/build | Agent sometimes skips tests; unreliable signal |
| UserPromptSubmit | Fires every turn; too noisy |
| Check `#in-progress` | Agent doesn't set labels reliably; circular |

Edit/Write is the most reliable signal that substantive work happened.

### Debounce Mechanism

Counter file at `.context/state/edit-nudge-count`:

1. PostToolUse fires on Edit/Write
2. `ctx system check-task-completion` increments counter
3. Every Nth invocation (configurable, default 5), emit nudge
4. Counter resets after nudge fires

The counter is a plain integer in a file (same pattern as `readCounter`/
`writeCounter` in `system/state.go`).

### Nudge Text

Short, actionable, one line:

```
If you completed a task, mark it [x] in TASKS.md.
```

Output via RESULT channel (agent sees it, user doesn't). The agent
decides whether it's relevant — no TASKS.md parsing in the hook.

### No TASKS.md Dependency

The hook does NOT:
- Parse TASKS.md
- Check for `#in-progress` labels
- Try to match files to tasks
- Make decisions about what's complete

It simply says: "you did work, check if something is done." The agent
connects the dots. This avoids the circular dependency on unreliable
task metadata.

### Configuration

`.ctxrc` option:

```yaml
task_nudge_interval: 5    # edits between nudges (0 = disabled)
```

Default: 5. Set to 0 to disable.

## Implementation

### New Command

`ctx system check-task-completion` — hidden, hook-only.

Reads `--session-id` from stdin (standard hook input).
Reads/writes `.context/state/edit-nudge-count`.

### Hook Registration

Add to `hooks.json`:

```json
{
  "matcher": "Edit",
  "hooks": [{
    "type": "command",
    "command": "ctx system check-task-completion"
  }]
}
```

Also match `Write` (same hook, same counter).

### Files Changed

- `internal/cli/system/check_task_completion.go` — new command
- `internal/cli/system/check_task_completion_test.go` — tests
- `internal/assets/claude/hooks/hooks.json` — add PostToolUse matcher
- `internal/config/file.go` — add state filename constant
- `internal/rc/rc.go` — add `TaskNudgeInterval()` config reader

## Dependencies

- State consolidation spec (init guard ensures `.context/state/` exists)
- No dependency on dir relocation (this is project-scoped state)

## Non-Goals

- Automatic task completion (agent marks, not hook)
- TASKS.md parsing in the hook
- Nudging about specific tasks
- Replacing AGENT_PLAYBOOK instructions (this supplements, not replaces)
