# internal/cli/system — Hook Plumbing

The `ctx system` command hosts 7 visible subcommands and 26 hidden
subcommands that implement Claude Code hook logic. See `doc.go`
for the full subcommand catalog.

## Hook Protocol

All hook subcommands follow the same contract:

1. Claude Code fires a lifecycle event (UserPromptSubmit,
   PreToolUse, PostToolUse)
2. `hooks.json` routes to `ctx system <check-name>`
3. The subcommand reads JSON from stdin (2-second timeout)
4. Runs its check logic
5. Writes advisory output to stdout (or JSON for block commands)
6. Exits 0 — hooks never block initialization

### stdin JSON (from Claude Code)

```json
{
  "tool_name": "Bash",
  "tool_input": {"command": "..."},
  "session_id": "..."
}
```

Read via `core/input.go` with 2-second timeout. Missing or
malformed stdin produces an empty `HookInput` (graceful no-op).

### Block command output (PreToolUse only)

```json
{"decision": "block", "reason": "..."}
```

or

```json
{"decision": "allow"}
```

## Throttle Patterns

Most checks use daily throttle to avoid running every prompt:

```
marker file: .context/state/<check-name>-marker
logic: compare marker mtime date to today
  same day → skip
  different day → run check, touch marker
```

### Adaptive counters

Some checks use prompt counters instead of daily throttle:

- **check-context-size**: silent 1-15, every 5th 16-30,
  every 3rd 30+
- **check-persistence**: silent 1-10, nudge at #20, then every 15

Counter state lives in `.context/state/` via `core/counter/`.

## Subcommand Categories

### Visible (user-invocable)

| Command | Purpose |
|---------|---------|
| `backup` | Timestamped tar.gz of context + Claude data |
| `bootstrap` | Print context dir path (for agent init) |
| `events` | Display event log entries |
| `message` | Manage hook message templates (list/show/edit/reset) |
| `prune` | Remove stale state files |
| `resources` | System resource usage with threshold alerts |
| `stats` | Session token statistics |

### Plumbing (hidden, used by skills/automation)

| Command | Purpose |
|---------|---------|
| `mark-journal` | Update journal .state.json |
| `mark-wrapped-up` | Record wrap-up ceremony timestamp |
| `session-event` | Record session start/end lifecycle |

### UserPromptSubmit hooks (hidden, 14 checks)

| Command | Trigger | Throttle |
|---------|---------|----------|
| `check-context-size` | Every prompt | Adaptive counter |
| `check-persistence` | Every prompt | Adaptive counter |
| `check-ceremony` | Every prompt | Per-ceremony cooldown |
| `check-journal` | Every prompt | Daily |
| `check-version` | Every prompt | Daily |
| `check-resources` | Every prompt | Daily |
| `check-knowledge` | Every prompt | Daily |
| `check-map-staleness` | Every prompt | Daily |
| `check-memory-drift` | Every prompt | Daily |
| `check-reminder` | Every prompt | None (always runs) |
| `check-freshness` | Every prompt | Daily |
| `check-backup-age` | Every prompt | Daily |
| `check-skill-discovery` | Every prompt | One-shot |
| `heartbeat` | Every prompt | None (telemetry) |

### PreToolUse hooks (hidden, 6 matchers)

| Command | Matches | Action |
|---------|---------|--------|
| `block-non-path-ctx` | Bash | Block bare `./ctx` invocations |
| `block-dangerous-command` | Bash | Block destructive patterns |
| `context-load-gate` | All tools | Inject context with cooldown |
| `qa-reminder` | Bash | Lint/test reminder before commits |
| `specs-nudge` | EnterPlanMode | Save plans to specs/ |
| `pause`/`resume` | All | Session-scoped hook suppression |

### PostToolUse hooks (hidden, 2 matchers)

| Command | Matches | Action |
|---------|---------|--------|
| `post-commit` | Bash (git commit) | Context capture nudge |
| `check-task-completion` | Edit/Write | Task completion nudge |

## Shared Infrastructure

| Package | Purpose |
|---------|---------|
| `core/input/` | Hook stdin JSON reading (2s timeout) |
| `core/counter/` | Prompt counter with adaptive thresholds |
| `core/session/` | Session token info extraction |
| `core/persistence/` | Persistence state tracking |
| `core/heartbeat/` | Heartbeat mtime management |
| `core/load/` | Context-load-gate state |
| `core/archive/` | Backup to SMB shares |

## Adding a New Hook

1. Create `cmd/<hook_name>/` with `cmd.go` and `run.go`
2. `run.go`: read stdin via `input.Read()`, implement check logic,
   write output to cmd
3. Register in `system.go` as a hidden subcommand
4. Add entry to `internal/assets/claude/hooks/hooks.json` with the
   appropriate lifecycle event and matcher
5. Add hook message template in `internal/assets/hooks/messages/`
   if the hook produces configurable output

## Caution

These hidden subcommands are effectively **public API** for agent
integrations. Changing throttle timing, output format, or exit
behavior affects every connected agent silently. Treat changes
as breaking changes.
