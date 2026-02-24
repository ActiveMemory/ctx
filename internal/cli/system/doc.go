//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package system provides the "ctx system" command for resource monitoring
// and hidden subcommands that implement Claude Code hook logic as native Go
// binaries, replacing the shell scripts previously deployed to .claude/hooks/.
//
// The parent command (ctx system) displays memory, swap, disk, and load
// metrics with WARNING/DANGER threshold alerts.
//
// Plumbing subcommands (hidden, used by skills and automation):
//   - mark-journal: Update journal processing state (.state.json)
//
// Hook subcommands read JSON from stdin (Claude Code hook contract), perform
// their logic, and exit 0. Block commands output JSON with a "decision" field.
//
// Hook subcommands (hidden):
//   - check-context-size: Adaptive prompt counter with checkpoint messages
//   - check-persistence: Context file mtime watcher with persistence nudges
//   - check-journal: Unexported sessions + unenriched entries reminder
//   - check-version: Version update nudge
//   - check-resources: Resource pressure monitor (DANGER-only VERBATIM relay)
//   - check-knowledge: Knowledge file growth nudge (daily throttle)
//   - check-map-staleness: Architecture map staleness nudge (daily throttle)
//   - block-non-path-ctx: Blocks non-PATH ctx invocations
//   - block-dangerous-commands: Blocks dangerous command patterns (project-local)
//   - check-backup-age: Backup staleness check (project-local)
//   - post-commit: Post-commit context capture nudge
//   - cleanup-tmp: Removes stale temp files on session end
//   - qa-reminder: Reminds agent to lint/test full project before declaring done
package system
