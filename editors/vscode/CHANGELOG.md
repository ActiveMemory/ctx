# Changelog

All notable changes to the **ctx: Persistent Context for AI** extension
will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [0.10.0] - 2026-07-05

### Fixed

- **Commands now target the real `ctx` CLI.** The participant dispatched
  a dozen-plus commands to subcommands that don't exist on the shipped
  binary (`ctx tasks`, `ctx complete`, `ctx decisions`, `ctx recall`,
  `ctx add`, `ctx notify`, …), so a large part of the surface was dead on
  arrival. Every invocation is reconciled to the actual command tree:
  `/task`, `/change`, `/decision`, `/learning`, `/permission` (the CLI
  registry is singular); `complete`→`task complete`; `recall`→
  `journal source`; `notify`→`hook notify`; `/add <type>`→`<type> add`;
  tool-config `hook`→`setup`; `system stats`→`usage`, `system resources`→
  `sysinfo`, `system message`→`hook message`, `check-reminders`→
  `check-reminder`.
- **`runCtx` no longer reports failures as success.** Cancelled or
  timed-out runs reject instead of rendering partial output as a clean
  result, and the process exit code is surfaced to callers rather than
  assuming any output means success.
- **`/pause` and `/resume` pause and resume context hooks** (`ctx hook
  pause` / `ctx hook resume`) instead of writing a session-snapshot JSON
  while the hooks — which the command name promises to pause — kept firing.
- **The onboarding gate no longer blocks `/guide`, `/why`, `/config`, and
  `/hook`** in an uninitialized project; these mirror the CLI's init-exempt
  set (the "run /init first" gate previously hid exactly the commands a new
  user reaches for).
- **Reminder `$(bell)` clears when the queue is empty.** The status bar now
  reads `ctx remind list` ("No reminders.") instead of the provenance-only
  `check-reminder` output, which never matched the hide condition and pinned
  the bell on permanently.
- **`/worktree` and `/changelog` git calls are cancellable** and time out
  after 30s, so a git op blocked on an index lock can't hang the request.
- **`saveWatcher` no longer runs against the wrong workspace root** — it
  skips paths starting with `..`, matching its siblings.
- Stopped regenerating the tracked `.github/copilot-instructions.md` on
  every `.context/**` change (git churn / write amplification).

### Changed

- `/verify` and `/wrapup` descriptions corrected to match actual behavior
  (both are read-only: `/verify` runs `ctx doctor` + `drift`; `/wrapup`
  summarizes status, drift, and journal and *suggests* — does not persist —
  decision/learning entries).

### Added

- **Command-parity test** (`commandParity.test.ts`): asserts
  `package.json` ↔ dispatcher ↔ the real `ctx` command tree (built from
  the same commit), so a future CLI rename can never silently strand a
  command again.

### Removed

- `/prompt` and `/deps` slash commands — their `ctx prompt` / `ctx dep`
  backing commands were removed from the CLI with no replacement. The
  command surface is now 43 commands, each mapping to a real `ctx`
  command.
- `/system backup` subcommand (no CLI replacement).
- **Unreachable Command-Palette registrations** — `activate()` registered
  `ctx.*` commands with no matching `contributes.commands`, so none were
  reachable from the palette. Removed pending a deliberate palette design.
- **Violation guardrails** (the terminal-command watcher, the
  sensitive-file watcher, and `.context/state/violations.json`
  recording). Capturing the user's terminal text — credentials
  included — into a file that an MCP server relays into model context
  is a design decision that warrants its own review. It will be
  re-proposed separately with a design note on the capture surface.

## [0.9.0] - 2026-03-19

### Added

- **@ctx chat participant** with 45 slash commands covering context
  lifecycle, task management, session recall, and discovery
- **Natural language routing**: type plain English after `@ctx` and
  the extension maps it to the correct handler
- **Auto-bootstrap**: downloads the ctx CLI binary if not found on PATH
- **Detection ring**: terminal command watcher and file edit watcher
  record governance violations for the MCP engine
- **Status bar reminders**: `$(bell) ctx` indicator for pending reminders
- **Automatic hooks**: file save, git commit, dependency change, and
  context file change handlers
- **Follow-up suggestions**: context-aware buttons after each command
- **`/diag` command**: diagnose extension issues with step-by-step timing

### Configuration

- `ctx.executablePath`: path to the ctx CLI binary (default: `ctx`)

## [Unreleased]

- Marketplace publication
