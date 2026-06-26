# Cheat Sheets

Quick mental models for key lifecycles and flows in ctx.

## CLI Command Dispatch

Steps:
1. `cmd/ctx/main.go` calls `bootstrap.RootCmd()`
2. RootCmd creates Cobra root with global flags
3. `Initialize()` registers 34 commands in 8 groups
4. PersistentPreRunE fires: boundary check, init check
5. Cobra routes to matched command's Run() handler
6. Handler calls core/ logic, domain packages
7. Output via write/* package, errors via err/* package

Key invariants:
- PersistentPreRunE runs for ALL subcommands
- AnnotationSkipInit bypasses init check (doctor, guide, setup)
- All flag descriptions from YAML (enforced by audit)

Common failure modes:
- Missing .context/ without SkipInit -> error before Run()
- Flag name collision with global flag -> silent shadowing

```
  main.go --> RootCmd() --> PersistentPreRunE --> Run()
                |                  |               |
         Initialize()        boundary +        core/ logic
         34 commands        init guards       --> write/*
```

---

## Context Loading (ctx agent)

Steps:
1. Agent calls `ctx agent --budget N`
2. rc.TokenBudget() resolves budget (flag > env > .ctxrc > 8000)
3. context/load.Do() reads all .md files from .context/
4. Files sorted by config.FileReadOrder priority
5. Each file's tokens estimated (4 chars/token)
6. Files added in priority order until budget exhausted
7. Overflow files listed as "Also Noted" summaries
8. Markdown packet returned to stdout

Key invariants:
- CONSTITUTION always loaded first
- Symlinks rejected (M-2 defense)
- Budget is conservative overestimate (never under-counts)

Common failure modes:
- Very large TASKS.md consumes most of budget -> low-priority
  files (GLOSSARY, PLAYBOOK) never seen by agent
- Empty file detection via EffectivelyEmpty() -> skipped with note

```
  Agent --> rc.TokenBudget() --> load.Do()
               |                    |
         resolve priority     read + estimate tokens
               |                    |
         sort by order       fit to budget
               |                    |
         format packet       "Also Noted" overflow
```

---

## MCP Request Lifecycle

Steps:
1. Client sends JSON-RPC line to stdin
2. Server.Serve() reads line from scanner
3. parse.Request() unmarshals JSON
4. dispatch.Do() routes by method name
5. Handler executes domain logic
6. session.CheckGovernance() appends warnings
7. out.*Response() wraps result
8. io.Writer.WriteJSON() writes to stdout

Key invariants:
- Main loop is single-threaded (sequential processing)
- Governance is advisory only (never blocks)
- Notifications (no ID) produce no response
- Poller runs independently on 5s interval

Common failure modes:
- Slow handler blocks all subsequent requests
- Parse error -> error response, loop continues
- Scanner overflow -> truncated JSON -> parse error

```
  Client                 Server              Handler
  |--JSON-RPC line------>|                   |
  |                      |--parse()          |
  |                      |--dispatch()------>|
  |                      |                   |--domain logic
  |                      |                   |--governance check
  |                      |<--result----------|
  |<--JSON-RPC response--|                   |
```

---

## Journal Import Pipeline

Steps:
1. User runs `ctx journal source --all`
2. parser.FindSessionsForCWD() scans ~/.claude/projects/
3. Auto-detects format (JSONL, Copilot, Copilot CLI, Markdown)
4. Matches sessions by git remote URL and CWD
5. Loads journal state from .state.json
6. Plans each session: new, regen, skip, or locked
7. Formats matched sessions as Markdown
8. Writes to .context/journal/
9. Marks imported in state, saves state

Key invariants:
- Locked entries never regenerated
- State tracks 5 stages: exported -> enriched -> normalized
  -> fences_verified -> locked
- Atomic state writes (temp + rename)

Common failure modes:
- Changed JSONL format -> silent parse failures, empty sessions
- Same project in multiple paths -> duplicate imports
- 1MB buffer limit -> truncated large tool results

```
  [Scan] --> [Detect Format] --> [Match CWD]
     |             |                  |
  4 parsers    auto-detect       git remote
     |             |                  |
  [Plan] ----> [Format MD] ----> [Write]
     |                               |
  state check                  mark imported
```

---

## Hook Lifecycle (UserPromptSubmit)

Steps:
1. Claude Code fires UserPromptSubmit hook
2. hooks.json routes to `ctx system check-*` commands
3. system/input.go reads hook JSON from stdin (2s timeout)
4. Each check runs independently, writes result to stdout
5. Checks: context-size, ceremonies, persistence, journal,
   reminders, version, resources, knowledge, map-staleness,
   memory-drift, freshness, heartbeat
6. Advisory output returned to Claude Code
7. All hooks exit 0 (never block initialization)

Key invariants:
- 2-second stdin read timeout (prevents hanging)
- Daily throttle via marker file date comparison
- Adaptive prompt counter: silent 1-15, periodic 16+
- All hooks exit 0 (never block)

Common failure modes:
- Missing stdin JSON -> timeout, graceful empty response
- Throttle marker file corruption -> check runs every prompt

```
  Claude Code         hooks.json        ctx system check-*
  |--hook fire------->|                 |
  |                   |--route--------->|
  |                   |  (12 checks)    |--read stdin (2s)
  |                   |                 |--check logic
  |                   |                 |--throttle gate
  |<--advisory--------|<--result (0)----|
```

---

## Entry Write Flow

Steps:
1. Caller provides EntryParams (type, content, opts)
2. entry.Validate() checks required fields per type
   - Decision: context, rationale, consequence
   - Learning: context, lesson, application
   - Task/Convention: content only
3. entry.Write() reads existing file
4. Formats entry per type template (from tpl/)
5. Inserts at correct position (tasks: before first unchecked)
6. Writes file back via io.SafeWriteFile()
7. Updates index for decisions/learnings (not tasks/conventions)

Key invariants:
- Entry headers are timestamped: `## [YYYY-MM-DD-HHMMSS] Title`
- Index updated between INDEX:START/END markers
- Three callers: CLI add, MCP handler, watch command

Common failure modes:
- Concurrent writes -> last writer wins (no locking)
- Index update fails after write -> stale index
- Missing required field -> validation error before write

```
  Caller --> Validate() --> Write()
                |              |
          check fields    read-modify-write
                |              |
          type-specific   format + insert
                |              |
          error early     update index
```

---

## Execution Flow Index (enriched 2026-06-09 via GitNexus)

_Auto-detected from the call graph (203 flows in index @ 60d8e823).
Complements the manually written cheat sheets above._

Top cross-community flows (spanning multiple domains). Nearly all
deep flows terminate at desc.Text — text lookup is the universal
last step of every command:

| Flow | Steps | Entry Point |
|------|-------|-------------|
| RunShow -> Text | 10 | internal/cli/journal/core/source/show.go |
| Run -> Text (journal import) | 9 | internal/cli/journal/cmd/importer/run.go |
| Run -> Text (checkversion hook) | 9 | internal/cli/system/cmd/checkversion/run.go |
| Run -> Text (journal site) | 8 | internal/cli/journal/cmd/site/run.go |
| Run -> Text (checkcontextsize hook) | 8 | internal/cli/system/cmd/checkcontextsize/run.go |
| Run -> Text (pad merge) | 8 | internal/cli/pad/cmd/merge/run.go |
| BuildVault -> Text (obsidian export) | 8 | internal/cli/journal/core/obsidian/vault.go |
| Run -> Text (pad tag) | 8 | internal/cli/pad/cmd/tag/run.go |
| Run -> Text (context load gate hook) | 7 | internal/cli/system/cmd/contextloadgate/run.go |
| Run -> Text (setup) | 7 | internal/cli/setup/cmd/root/run.go |
| Run -> Text (add) | 7 | internal/cli/add/core/run/run.go |

### Multi-Flow Hotspots

Symbols participating in 3+ flows (high-impact modification
points), counted via STEP_IN_PROCESS edges:

| Symbol | Flows | Location |
|--------|-------|----------|
| desc.Text | 58 | internal/assets/read/desc/desc.go |
| check.FullPreamble | 40 | internal/cli/system/core/check/full_preamble.go |
| check.Preamble | 39 | internal/cli/system/core/check/input.go |
| load.Do | 28 | internal/context/load/loader.go:34 |
| session.ReadInput | 26 | internal/cli/system/core/session/session.go |
| err/context.StatFailed | 25 | internal/err/context/context.go |
| rc.ContextDir | 25 | internal/rc/rc.go |
| nudge.Paused | 21 | internal/cli/system/core/nudge/pause.go |
| hub.Unmarshal | 20 | internal/hub/types.go |
| rc.RC | 20 | internal/rc/rc.go |
| io.cleanAndValidate | 15 | internal/io/validate.go |
| io.SafeReadUserFile | 14 | internal/io/security.go |
| nudge.PauseMarkerPath | 14 | internal/cli/system/core/nudge/pause.go |

New since 2026-04-03: the system-hook preamble pair
(FullPreamble/Preamble at 40/39 flows) and nudge pause gating
(21/14 flows) are now top-tier integration points — every hook
subcommand routes through them. hub.Unmarshal/Name (20/18) mark
the hub subsystem as a new flow participant.
