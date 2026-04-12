//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package system provides the "ctx system" hidden parent command that
// hosts Claude Code hook plumbing subcommands as native Go binaries,
// replacing the shell scripts previously deployed to .claude/hooks/.
//
// User-facing maintenance commands (backup, bootstrap, event, message,
// prune, resource, stats) have been promoted to top-level commands and
// are registered in internal/bootstrap/group.go.
//
// Plumbing subcommands (used by skills and automation):
//   - mark-journal: Update journal processing state (.state.json)
//   - mark-wrapped-up: Record wrap-up ceremony timestamp
//   - session-event: Record session lifecycle events (start, end)
//   - pause: session-scoped hook suppression
//   - resume: session-scoped hook re-enable
//
// Hook subcommands read JSON from stdin (Claude Code hook contract), perform
// their logic, and exit 0. Block commands output JSON with a "decision" field.
//
// UserPromptSubmit hooks (hidden):
//   - check-context-size: Adaptive prompt counter with checkpoint messages
//   - check-persistence: Context file mtime watcher with persistence nudges
//   - check-ceremony: Session ceremony reminder (remember, wrap-up)
//   - check-journal: Unexported sessions + unenriched entries reminder
//   - check-version: Version update and key rotation nudge
//   - check-resource: Resource pressure monitor (DANGER-only VERBATIM relay)
//   - check-knowledge: Knowledge file growth nudge (daily throttle)
//   - check-map-staleness: Architecture map staleness nudge (daily throttle)
//   - check-memory-drift: Memory bridge drift detection
//   - check-reminder: Session reminder surfacing
//   - check-freshness: Technology-dependent constant staleness check
//   - check-hub-sync: Auto-sync ctx Hub entries on session start
//   - check-backup-age: Backup staleness check (project-local)
//   - check-skill-discovery: One-shot mid-session skill tip nudge
//   - heartbeat: Token telemetry and billing threshold check
//
// PreToolUse hooks (hidden):
//   - block-non-path-ctx: Blocks non-PATH ctx invocations
//   - block-dangerous-command: Blocks dangerous command patterns
//   - context-load-gate: Context injection on tool use with cooldown
//   - qa-reminder: Reminds agent to lint/test before declaring done
//   - specs-nudge: Reminds agent to save plans to specs/
//
// PostToolUse hooks (hidden):
//   - post-commit: Post-commit context capture nudge
//   - check-task-completion: Task completion nudge after edits
package system
