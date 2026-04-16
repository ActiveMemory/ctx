//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package load_gate defines constants for the context
// load gate hook, which controls when and how project
// context is injected into a new Claude Code session.
//
// # Session Marker Files
//
// Each session that has loaded context writes a marker
// file with the PrefixCtxLoaded prefix ("ctx-loaded-")
// followed by the session ID. The load gate checks for
// this marker to avoid double-loading in the same
// session and to auto-prune stale markers older than
// AutoPruneStaleDays (7 days).
//
// # Event Logging
//
// Every load gate invocation emits an event named
// EventContextLoadGate ("context-load-gate") to the
// event log. The event payload includes a timestamp
// extracted from log lines via JSONKeyTimestamp.
//
// # Output Formatting
//
// The load gate wraps its injected context block with
// separator lines built from ContextLoadSeparatorChar
// ("=") repeated ContextLoadSeparatorWidth (80) times,
// producing a visible header and footer in the
// session transcript.
//
// # Key Constants
//
//   - PrefixCtxLoaded -- marker file prefix for
//     per-session load tracking.
//   - AutoPruneStaleDays -- days before stale markers
//     are auto-pruned (7).
//   - EventContextLoadGate -- event name logged on
//     each context injection.
//   - ContextLoadSeparatorChar / Width -- visual
//     separator for the injected block.
//   - JSONKeyTimestamp -- JSON key used to extract
//     timestamps from event log lines.
package load_gate
