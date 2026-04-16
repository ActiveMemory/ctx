//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package usage implements the ctx usage top-level
// command.
//
// Displays per-session token usage statistics from the
// local telemetry log stored in .context/state/. Each
// AI session records input and output token counts as
// JSONL entries, and usage reads, aggregates, and
// presents them.
//
// # Output Modes
//
//   - Table (default): a human-readable table showing
//     session ID, timestamp, input tokens, output tokens,
//     and total
//   - JSON (--json): structured output for scripting
//   - Follow (--follow): live-streaming mode that tails
//     the telemetry log and prints new entries as they
//     arrive, useful for monitoring active sessions
//
// # Filtering
//
//   - --session: filter to a specific session ID
//   - --last: limit to the N most recent entries
//
// [Cmd] returns the cobra command with --json, --follow,
// --session, and --last flags. [Run] reads the JSONL
// telemetry log, aggregates token counts, and renders
// the result as a table, JSON, or live-tailed stream.
package usage
