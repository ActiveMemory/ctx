# Changelog

All notable changes to the **ctx: Persistent Context for AI** extension
will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [0.10.0] - Unreleased

### Added

- **Skill-backed commands**: `/brainstorm`, `/spec`, `/implement`,
  `/next`, `/remember`, `/reflect`, `/wrap-up`, `/blog`, and
  `/consolidate` each run the canonical `ctx-<name>` skill through the
  chat model. The skill text is bundled from
  `internal/assets/claude/skills/` at build time; each request is
  grounded in `ctx agent` output, the read-only ctx output the skill
  relies on, earlier turns of the same skill conversation, and `#file`
  attachments. The model proposes commands and edits; it never claims to
  have run them. A plain reply after a skill answer continues that skill.
- `/decision` and `/learning` list entries (`ctx index`).
- Reminder status bar: `$(bell) ctx` while `ctx remind list` has entries.
- Session start and end events (`ctx system session-event`) on
  activation, after `/init`, and on deactivation.
- **Command-parity guard.** `commandParity.test.ts` checks that
  `package.json` declares exactly the dispatched commands, drives every
  command branch through the chat handler, and records each `ctx` argv in
  `src/ctx-cli-surface.json`. The Go test
  `internal/bootstrap/vscode_surface_test.go` parses every recorded argv
  against the real command tree, runs the entry `add` invocations, and
  checks every bundled skill still ships, so a CLI rename fails CI.

### Fixed

- **Every command targets the current CLI.** All invocations passed the
  removed `--no-color` flag, so every command failed. Also reconciled:
  `recall list` → `journal source`; `add <type>` → `<type> add` with the
  required provenance (`--session-id`, `--branch`, `--commit`);
  `notify` → `hook notify`; `system resources` → `sysinfo`;
  `system message` → `hook message`; `pause`/`resume` →
  `hook pause`/`hook resume`; `/why` with no argument no longer opens the
  CLI's interactive menu.
- **Failures are shown as failures.** A non-zero exit renders the CLI's
  message under "exited with code N" instead of as a normal result;
  cancellation, the 30s timeout, and spawn or buffer failures are
  errors, never partial output.
- **No shell on Windows.** Prompt text reached `cmd.exe` unquoted, so a
  multi-word argument split apart and `&` ran a second command. The CLI
  now runs without a shell, and multi-word values stay one argument
  (`/pad edit`, `/notify <message>`, quoted `/add` flag values).
- stdin is closed, so no command waits on a prompt until the timeout.
- Natural-language routing only reaches read-only commands; a keyword
  match can no longer complete a task or add a reminder.
- In a multi-root window, commands run in the folder of the active
  editor.

### Removed

- `/prompt` and `/dep`: `ctx prompt` and `ctx dep` no longer exist.
- `/reindex`: `ctx reindex` no longer exists; indices are projected on
  demand (`/decision`, `/learning`).
- `/loop`: `ctx loop` writes a shell script that drives a terminal agent,
  not a chat workflow.
- `/site`: `ctx site` is a hidden maintainer command.

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
