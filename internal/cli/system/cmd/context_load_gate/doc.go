//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package context_load_gate implements the hidden
// "ctx system context-load-gate" cobra subcommand.
//
// This hook fires on the first PreToolUse event of each
// agent session. It injects the project's CONSTITUTION
// and distilled AGENT_PLAYBOOK_GATE into the agent's
// context window so that hard rules and directives are
// available before any tool executes.
//
// # Behavior
//
// On the first tool call of a session the hook:
//
//   - Reads CONSTITUTION.md and AGENT_PLAYBOOK_GATE.md
//     from the context directory.
//   - Estimates token counts for each file and appends
//     them to a combined payload.
//   - Scans for recent context-file and code changes
//     since the last reference time and appends a
//     changes summary.
//   - Writes the payload to stdout so the agent sees it
//     as a PreToolUse context block.
//   - Sends a webhook notification with file count and
//     total token metadata (never file content).
//   - Writes an oversize flag file when the injected
//     token total exceeds the configured threshold.
//
// A marker file ensures the hook is one-shot per
// session: parallel tool calls in the same session
// are silently skipped after the first fires.
//
// Stale session state files are auto-pruned once per
// session at startup.
//
// # Flags
//
// None. The command reads hook JSON from stdin.
//
// # Output
//
// Emits a PreToolUse context block containing the
// CONSTITUTION, distilled playbook, and a changes
// summary. No output on subsequent calls within the
// same session.
//
// # Delegation
//
// Token estimation is handled by context/token.
// Change detection delegates to change/core/detect
// and change/core/scan. State management uses
// system/core/state and config/load_gate. Webhook
// relay uses system/core/nudge.
package context_load_gate
