//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package system provides the "ctx system" hidden parent
// command that hosts Claude Code hook plumbing subcommands
// as native Go binaries, replacing the shell scripts
// previously deployed to .claude/hooks/.
//
// User-facing maintenance commands (backup, prune, sysinfo,
// usage) are registered as top-level commands in
// internal/bootstrap/group.go. Hook-facing commands (event,
// message, notify, pause, resume) live under the "ctx hook"
// parent, also registered in group.go.
//
// # Agent-Only Subcommands
//
//   - bootstrap: print context location for AI agents
//     (hidden, agent-only)
//
// # Plumbing Subcommands
//
// Used by skills and automation:
//
//   - mark-journal: update journal processing state
//   - mark-wrapped-up: record wrap-up ceremony timestamp
//   - session-event: record session lifecycle events
//   - pause: session-scoped hook suppression
//   - resume: session-scoped hook re-enable
//
// # Hook Subcommands
//
// Hook subcommands read JSON from stdin (Claude Code hook
// contract), perform their logic, and exit 0. Block
// commands output JSON with a "decision" field.
//
// UserPromptSubmit hooks (hidden):
//   - check-anchor-drift: stale CTX_DIR vs project anchor warning
//   - check-context-size: adaptive prompt counter
//   - check-persistence: context file mtime watcher
//   - check-ceremony: session ceremony reminder
//   - check-journal: unexported sessions reminder
//   - check-version: version update nudge
//   - check-resource: resource pressure monitor
//   - check-knowledge: knowledge file growth nudge
//   - check-map-staleness: architecture map nudge
//   - check-memory-drift: memory bridge drift detection
//   - check-reminder: session reminder surfacing
//   - check-freshness: constant staleness check
//   - check-hub-sync: auto-sync Hub entries
//   - check-skill-discovery: skill tip nudge
//   - heartbeat: token telemetry and billing check
//
// PreToolUse hooks (hidden):
//   - block-non-path-ctx: blocks non-PATH ctx calls
//   - context-load-gate: context injection with cooldown
//   - qa-reminder: lint/test before done reminder
//   - specs-nudge: save plans to specs/ reminder
//
// PostToolUse hooks (hidden):
//   - post-commit: post-commit context capture nudge
//   - check-task-completion: task completion nudge
//
// # Subpackages
//
//	cmd/: one subpackage per hook or plumbing command
//	core/: shared helpers (archive, event, health,
//	  resource, state, stats)
package system
